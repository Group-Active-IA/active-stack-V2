package reports

import (
	"encoding/json"
	"time"
)

// manifestVersion versiona el esquema del manifiesto, no el contenido
// capturado. Subir esta constante es el único lugar que toca un lector de
// manifiestos viejos el día que el esquema cambie.
const manifestVersion = 1

// manifest es el documento estructurado que acompaña a cada informe
// capturado ("manifiesto estructurado" del propose). Se serializa a
// manifest.json dentro del directorio de la identidad.
type manifest struct {
	Version int `json:"version"`

	Key    string `json:"key"`
	Change string `json:"change"`
	Phase  string `json:"phase"`
	Round  int    `json:"round"`
	Role   string `json:"role"`

	DispatchID string `json:"dispatch_id,omitempty"`
	SHA        string `json:"sha,omitempty"`
	DiffRef    string `json:"diff_ref,omitempty"`

	Contract string `json:"contract"`
	Agent    string `json:"agent,omitempty"`
	Session  string `json:"session,omitempty"`
	Turn     string `json:"turn,omitempty"`
	// Provider y Model sólo se vuelcan si el Source los traía (EVIDENCIA de
	// origen); nunca se infieren del texto del informe.
	Provider string `json:"provider,omitempty"`
	Model    string `json:"model,omitempty"`

	RawHash      string `json:"raw_hash,omitempty"`
	RawSizeBytes int    `json:"raw_size_bytes,omitempty"`
	ContentHash  string `json:"content_hash"`
	SizeBytes    int    `json:"size_bytes"`

	Status     Status `json:"-"`
	StatusText string `json:"status"`
	Cause      string `json:"cause,omitempty"`
	// SecretKinds nunca incluye el valor detectado, sólo el tipo de patrón.
	SecretKinds []string `json:"secret_kinds,omitempty"`
	// RetainedPath es relativo a RepoDir, siempre bajo una ruta ignorada por
	// git. Vacío salvo en Status == StatusRetained.
	RetainedPath string `json:"retained_path,omitempty"`

	CapturedAt time.Time `json:"captured_at"`
}

func (m manifest) encode() ([]byte, error) {
	m.StatusText = m.Status.String()
	return json.MarshalIndent(m, "", "  ")
}

func decodeManifest(b []byte) (manifest, error) {
	var m manifest
	if err := json.Unmarshal(b, &m); err != nil {
		return manifest{}, err
	}
	switch m.StatusText {
	case "COMPLETE":
		m.Status = StatusComplete
	case "RETAINED":
		m.Status = StatusRetained
	case "CONFLICT":
		m.Status = StatusConflict
	case "UNSUPPORTED":
		m.Status = StatusUnsupported
	default:
		m.Status = StatusPending
	}
	return m, nil
}
