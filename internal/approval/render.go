package approval

import (
	"encoding/json"
	"fmt"
	"strings"
)

// RenderText arma la salida legible del CLI. Para WRITTEN_UNVERIFIED repite
// el mensaje de contrato literal: es la parte que no se puede resumir ni
// reformular sin volver a abrir la pregunta de identidad acreditada.
func RenderText(r Result) string {
	var b strings.Builder
	fmt.Fprintf(&b, "estado: %s (código %d)\n", r.Status, r.Status.ExitCode())

	if r.Status == StatusWrittenUnverified {
		fmt.Fprintf(&b, "%s\n", r.Message)
		fmt.Fprintf(&b, "alcance: %s\n", r.Scope)
		fmt.Fprintf(&b, "propuesta (P): %s\n", r.ProposalSHA)
		fmt.Fprintf(&b, "aprobación (C): %s\n", r.CommitSHA)
		fmt.Fprintf(&b, "manifiesto: %d archivo(s)\n", len(r.Manifest))
		return b.String()
	}

	fmt.Fprintf(&b, "causa: %s\n", r.Cause)
	if r.Field != "" {
		fmt.Fprintf(&b, "campo/archivo: %s\n", r.Field)
	}
	if r.Recovery != "" {
		fmt.Fprintf(&b, "recuperación: %s\n", r.Recovery)
	}
	return b.String()
}

// RenderJSON es la salida para consumidores que no parsean texto. El campo
// "status" viaja como la cadena estable (MISSING/INVALID/STALE/
// WRITTEN_UNVERIFIED), nunca como el int de iota — ese es un detalle interno
// de Go, no un contrato.
func RenderJSON(r Result) ([]byte, error) {
	out := struct {
		Status      string          `json:"status"`
		ExitCode    int             `json:"exit_code"`
		Message     string          `json:"message,omitempty"`
		Cause       string          `json:"cause,omitempty"`
		Field       string          `json:"field,omitempty"`
		Recovery    string          `json:"recovery,omitempty"`
		Scope       string          `json:"scope,omitempty"`
		ProposalSHA string          `json:"proposal_sha,omitempty"`
		CommitSHA   string          `json:"commit_sha,omitempty"`
		Manifest    []ManifestEntry `json:"manifest,omitempty"`
	}{
		Status:      r.Status.String(),
		ExitCode:    r.Status.ExitCode(),
		Message:     r.Message,
		Cause:       r.Cause,
		Field:       r.Field,
		Recovery:    r.Recovery,
		Scope:       r.Scope,
		ProposalSHA: r.ProposalSHA,
		CommitSHA:   r.CommitSHA,
		Manifest:    r.Manifest,
	}
	b, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("serializar el resultado de approval: %w", err)
	}
	return append(b, '\n'), nil
}
