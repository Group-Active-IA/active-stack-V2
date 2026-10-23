package reports

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

// MaxContentBytes es el límite de tamaño del informe. El propose exige que
// "límites de tamaño/timeout/retención (TBD) deben quedar fijados antes de
// habilitar captura" — 8 MiB es un valor PROVISORIO de este constructor, no
// una decisión de producto, para que exceder el límite falle visiblemente en
// vez de truncar en silencio. DEUDA: fijar el valor real y timeout/retención
// (sin implementar: no hay ningún límite de tiempo ni de borrado acá) es
// trabajo pendiente, explícito en el informe de apply.
const MaxContentBytes = 8 * 1024 * 1024

func capture(in CaptureInput) (CaptureResult, error) {
	root, err := requireCanonicalRoot(in.RepoDir)
	if err != nil {
		return CaptureResult{}, err
	}

	segs, err := dirSegments(in.Change, in.Phase, in.Round, in.Role)
	if err != nil {
		return CaptureResult{}, err
	}
	targetDir, err := resolveWithinRoot(root, segs)
	if err != nil {
		return CaptureResult{}, err
	}
	keyRound := in.Round
	if keyRound < 0 {
		keyRound = 0
	}
	key := formatKeyUnchecked(in.Change, in.Phase, keyRound, in.Role)
	relDir, err := filepath.Rel(root, targetDir)
	if err != nil {
		return CaptureResult{}, fmt.Errorf("ruta relativa de %s: %w", targetDir, err)
	}

	if !knownContract(in.Source.Contract) {
		return CaptureResult{
			Status: StatusUnsupported,
			Key:    key,
			Dir:    relDir,
			Cause:  fmt.Sprintf("contrato %q no está registrado: ningún adapter lo resolvió todavía", in.Source.Contract),
		}, nil
	}

	content := in.Source.Text
	raw := in.Source.Raw
	if raw == "" {
		raw = content
	}
	if len(raw) > MaxContentBytes {
		return CaptureResult{}, fmt.Errorf("reports.Capture: bruto excede MaxContentBytes")
	}
	if len(content) > MaxContentBytes {
		return CaptureResult{}, fmt.Errorf("reports.Capture: el informe ocupa %d bytes, supera el límite de %d (MaxContentBytes, provisorio); se rechaza, nunca se trunca", len(content), MaxContentBytes)
	}
	hash := sha256hex(content)
	blank := strings.TrimSpace(content) == ""

	metadata := strings.Join([]string{in.Source.Agent, in.Source.Session, in.Source.Turn, in.Source.Provider, in.Source.Model}, "\n")
	secretKinds := scanSecrets(content + "\n" + raw + "\n" + metadata)
	if len(scanSecrets(metadata)) > 0 {
		// Nunca copiar metadatos sensibles al manifiesto público.
		in.Source.Agent = ""
		in.Source.Session = ""
		in.Source.Turn = ""
		in.Source.Provider = ""
		in.Source.Model = ""
	}
	ambiguous := in.DispatchID == "" || in.SHA == "" || in.Round <= 0

	var status Status
	var cause string
	switch {
	case len(secretKinds) > 0:
		status = StatusRetained
		cause = "escaneo de secretos detectó: " + strings.Join(secretKinds, ", ")
	case blank:
		status = StatusPending
		cause = "informe vacío o en blanco"
	case ambiguous || in.Source.Incomplete:
		status = StatusPending
		cause = "evento incompleto o identidad ambigua: revisar cierre, sesión, dispatch id, sha y vuelta"
	default:
		status = StatusComplete
	}

	m := manifest{
		Version:      manifestVersion,
		Key:          key,
		Change:       in.Change,
		Phase:        in.Phase,
		Round:        in.Round,
		Role:         in.Role,
		DispatchID:   in.DispatchID,
		SHA:          in.SHA,
		DiffRef:      in.DiffRef,
		Contract:     in.Source.Contract,
		Agent:        in.Source.Agent,
		Session:      in.Source.Session,
		Turn:         in.Source.Turn,
		Provider:     in.Source.Provider,
		Model:        in.Source.Model,
		ContentHash:  hash,
		RawHash:      sha256hex(raw),
		RawSizeBytes: len(raw),
		SizeBytes:    len(content),
		Status:       status,
		Cause:        cause,
		SecretKinds:  secretKinds,
		CapturedAt:   time.Now().UTC(),
	}

	result, err := persistOrReuse(root, targetDir, m, raw, content)
	if err != nil {
		return CaptureResult{}, err
	}

	if err := regenerateLEEME(root, in.Change); err != nil {
		return CaptureResult{}, err
	}
	return result, nil
}

// persistOrReuse decide, mirando lo que YA hay en dir, si esta captura es un
// no-op (mismo hash: reintento, no incrementa nada), un conflicto (misma
// identidad, bytes distintos: el original existente se preserva sin pisar,
// este intento queda en un subdirectorio de conflicto propio) o una
// escritura nueva.
func persistOrReuse(root, dir string, m manifest, rawContent, text string) (CaptureResult, error) {
	existing, ok, err := readManifest(dir)
	if err != nil {
		return CaptureResult{}, err
	}
	if ok {
		if existing.ContentHash == m.ContentHash && (existing.RawHash == m.RawHash || existing.RawHash == "" && m.RawHash == m.ContentHash) && existing.Status == m.Status {
			return resultFromManifest(root, dir, existing), nil
		}
		conflictDir := filepath.Join(dir, "conflict-"+sha256hex(m.RawHash + m.ContentHash + m.Status.String())[:12])
		if m.Status != StatusRetained {
			m.Status = StatusConflict
		}
		relDir, _ := filepath.Rel(root, dir)
		m.Cause += "; " + fmt.Sprintf("misma identidad (%s), bytes distintos a lo ya persistido en %s; original preservado sin pisar", m.Key, relDir)
		// Retención conserva causa y tipos, incluso ante colisión.
		return persistOrReuse(root, conflictDir, m, rawContent, text)
	}

	written, err := writeFiles(root, dir, m, rawContent, text)
	if err != nil {
		return CaptureResult{}, err
	}
	return resultFromManifest(root, dir, written), nil
}

func resultFromManifest(root, dir string, m manifest) CaptureResult {
	rel, _ := filepath.Rel(root, dir)
	return CaptureResult{
		Status:       m.Status,
		Key:          m.Key,
		Dir:          rel,
		ContentHash:  m.ContentHash,
		Cause:        m.Cause,
		SecretKinds:  m.SecretKinds,
		RetainedPath: m.RetainedPath,
	}
}
