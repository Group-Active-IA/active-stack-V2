package approval

import (
	"fmt"
	"strings"
)

// propuestaRevisadaTrailer es la línea que el commit de aprobación (C) tiene
// que llevar en su mensaje para referenciar P (J1 revisada: "la decisión
// referencia el SHA completo P de la propuesta revisada"). La tabla de
// Aprobación no tiene columna para P — agregarla reabriría el formato que
// Juani ya firmó — así que la referencia vive en el commit documental mismo,
// que es la propia "decisión". Ver informe del constructor, sección JUICIO:
// es una interpretación, no algo que el propose deletree carácter a carácter.
const propuestaRevisadaTrailer = "Propuesta-revisada:"

func verify(opts Options) (Result, error) {
	dir := opts.RepoDir

	cFull, err := resolveCommit(dir, opts.Commit)
	if err != nil {
		return blocked(StatusInvalid, "el commit de aprobación no resuelve a un commit válido (SHA ambiguo o inexistente)",
			"--commit", "corregí --commit a un SHA o referencia que exista en este repositorio"), nil
	}

	dirty, err := isDirty(dir)
	if err != nil {
		return Result{}, fmt.Errorf("git status: %w", err)
	}
	if dirty {
		return blocked(StatusInvalid, "árbol documental sucio: hay cambios sin commitear",
			dir, "commiteá o descartá los cambios pendientes antes de validar"), nil
	}

	changed, err := changedPaths(dir, cFull)
	if err != nil {
		return Result{}, fmt.Errorf("git diff-tree de %s: %w", cFull, err)
	}
	touchesProposal := false
	var others []string
	for _, p := range changed {
		if p == opts.ProposalPath {
			touchesProposal = true
		} else {
			others = append(others, p)
		}
	}
	if !touchesProposal {
		return blocked(StatusInvalid, "el commit de aprobación no modifica la tabla de Aprobación",
			opts.ProposalPath, "--commit debe apuntar al commit que agrega la fila de Aprobación en "+opts.ProposalPath), nil
	}
	if len(others) > 0 {
		return blocked(StatusInvalid, "el commit de aprobación mezcla la tabla con cambios de código, no es un commit documental separado",
			strings.Join(others, ", "), "separá la aprobación en su propio commit, sin otros archivos"), nil
	}

	parent, err := parentOf(dir, cFull)
	if err != nil {
		return Result{}, fmt.Errorf("resolver el padre de %s: %w", cFull, err)
	}

	var before string
	if parent != "" {
		content, ok, errShow := showFile(dir, parent, opts.ProposalPath)
		if errShow != nil {
			return Result{}, fmt.Errorf("leer %s en %s: %w", opts.ProposalPath, parent, errShow)
		}
		if ok {
			before = content
		}
	}

	after, ok, err := showFile(dir, cFull, opts.ProposalPath)
	if err != nil {
		return Result{}, fmt.Errorf("leer %s en %s: %w", opts.ProposalPath, cFull, err)
	}
	if !ok {
		return blocked(StatusMissing, "la propuesta no existe en el commit de aprobación",
			opts.ProposalPath, "revisá la ruta pasada en --change"), nil
	}

	beforeRow, beforeFound := parseApprovalTable(before)
	afterRow, afterFound := parseApprovalTable(after)
	if !afterFound {
		return blocked(StatusMissing, "no hay tabla de Aprobación, o no tiene ninguna fila de datos",
			"## Aprobación", "completá la tabla de Aprobación en "+opts.ProposalPath+" con elección, persona, fecha y alcance, en un commit documental separado"), nil
	}
	if beforeFound && beforeRow.complete() {
		return blocked(StatusInvalid, "la tabla de Aprobación ya estaba completa antes del commit de aprobación: no cumple \"commit documental separado\"",
			opts.ProposalPath, "la fila de Aprobación tiene que agregarse en el commit que --commit identifica, no antes"), nil
	}
	if !afterRow.complete() {
		missing := afterRow.missingFields()
		return blocked(StatusInvalid, "faltan campos en la fila de Aprobación: "+strings.Join(missing, ", "),
			strings.Join(missing, ", "), "completá elección, persona, fecha y alcance en la fila de Aprobación"), nil
	}

	msg, err := commitMessage(dir, cFull)
	if err != nil {
		return Result{}, fmt.Errorf("leer el mensaje de %s: %w", cFull, err)
	}
	pRef := extractProposalRef(msg)
	if pRef == "" {
		return blocked(StatusInvalid, "el commit de aprobación no referencia el SHA completo (P) de la propuesta revisada",
			"mensaje del commit", "agregá una línea \""+propuestaRevisadaTrailer+" <sha completo de P>\" al mensaje del commit de aprobación"), nil
	}
	pFull, err := resolveCommit(dir, pRef)
	if err != nil {
		return blocked(StatusInvalid, "P no resuelve a un commit válido en este repositorio",
			propuestaRevisadaTrailer+" "+pRef, "corregí la referencia a P en el mensaje del commit de aprobación"), nil
	}
	if parent == "" {
		return blocked(StatusInvalid, "el commit de aprobación no tiene padre: P no puede ser su ancestro",
			opts.ProposalPath, "la aprobación tiene que ser un commit aparte, encima de la propuesta revisada"), nil
	}
	anc, err := isAncestor(dir, pFull, parent)
	if err != nil {
		return Result{}, fmt.Errorf("git merge-base --is-ancestor %s %s: %w", pFull, parent, err)
	}
	if !anc {
		return blocked(StatusInvalid, "P no es ancestro del commit de aprobación",
			propuestaRevisadaTrailer+" "+pFull, "P tiene que ser la propuesta revisada que el commit de aprobación tiene por debajo, no un commit ajeno"), nil
	}

	informes := referencedInformes(after, opts.ProposalPath)
	for _, inf := range informes {
		for _, p := range others {
			if p == inf {
				// Inalcanzable hoy: others ya está vacío cuando llegamos acá
				// (se bloqueó arriba). Se deja explícito por si esa regla
				// cambia y deja de implicar esta.
				return blocked(StatusInvalid, "el commit de aprobación modifica un informe referenciado, no sólo la tabla",
					p, "los informes no se tocan en el commit de aprobación"), nil
			}
		}
	}
	manifest, missingInforme, err := buildManifest(dir, pFull, opts.ProposalPath, informes)
	if err != nil {
		return Result{}, fmt.Errorf("construir el manifiesto en %s: %w", pFull, err)
	}
	if missingInforme != "" {
		return blocked(StatusInvalid, "falta un informe referenciado por la propuesta en la revisión P (falta o fue alterado)",
			missingInforme, "versioná el informe en P o corregí la referencia en la propuesta"), nil
	}

	currentRef := opts.CurrentRef
	if currentRef == "" {
		currentRef = "HEAD"
	}
	currentFull, err := resolveCommit(dir, currentRef)
	if err != nil {
		return Result{}, fmt.Errorf("resolver %s: %w", currentRef, err)
	}
	cHash, ok, err := blobHash(dir, cFull, opts.ProposalPath)
	if err != nil || !ok {
		return Result{}, fmt.Errorf("hash de %s en %s: %w", opts.ProposalPath, cFull, err)
	}
	curHash, ok, err := blobHash(dir, currentFull, opts.ProposalPath)
	if err != nil {
		return Result{}, fmt.Errorf("hash de %s en %s: %w", opts.ProposalPath, currentFull, err)
	}
	if !ok || curHash != cHash {
		return blocked(StatusStale, "la propuesta cambió después del commit de aprobación: la aprobación es para otro contenido",
			opts.ProposalPath, "revalidar requiere una aprobación nueva sobre el contenido vigente"), nil
	}

	return Result{
		Status:      StatusWrittenUnverified,
		Message:     writtenUnverifiedMessage,
		Scope:       afterRow.Alcance,
		ProposalSHA: pFull,
		CommitSHA:   cFull,
		Manifest:    manifest,
	}, nil
}

// extractProposalRef busca "Propuesta-revisada: <ref>" en el mensaje del
// commit y devuelve <ref> recortado, o "" si no está.
func extractProposalRef(msg string) string {
	for _, line := range strings.Split(msg, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, propuestaRevisadaTrailer) {
			return strings.TrimSpace(strings.TrimPrefix(line, propuestaRevisadaTrailer))
		}
	}
	return ""
}
