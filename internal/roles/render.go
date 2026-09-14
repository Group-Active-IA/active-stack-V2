package roles

import (
	"encoding/json"
	"fmt"
	"strings"
)

// anchoRol deja alineada la segunda columna. "role-architect" es el más largo.
const anchoRol = 16

// maxArchivos es cuántos archivos se muestran como hecho del diff antes de
// resumir. Tres alcanzan para justificar una convocatoria; el resto es el mismo
// dato repetido, y una salida que hay que scrollear se lee por arriba.
const maxArchivos = 3

// RenderText arma la salida legible.
//
// Los dos lados salen: el que corresponde con el hecho del diff que lo
// justifica, y el que NO corresponde con el motivo. Lo segundo es la mitad del
// valor del comando. Hoy un rol que no se convoca simplemente no aparece en
// ningún lado, y esa ausencia es indistinguible de haberlo olvidado.
func RenderText(rep Report) string {
	var b strings.Builder

	if len(rep.Files) == 0 {
		// No es un fallo. Un diff vacío es una respuesta.
		b.WriteString("sin archivos cambiados: nada que revisar\n")
		return b.String()
	}

	if rep.Floor.Triggered {
		fmt.Fprintf(&b, "PISO ACTIVADO: %s (%s)\n", rep.Floor.Fact, rep.Floor.File)
		b.WriteString("  se convocan backend, QA, DBA y arquitecto. " +
			"No se relaja por tamaño del diff: un cambio de tres líneas convoca el " +
			"piso igual que uno de trescientas\n\n")
	}

	for _, d := range ordenar(rep.Decisions) {
		fmt.Fprintf(&b, "%-*s%s\n", anchoRol, d.Role, columnaDerecha(d))
	}

	if len(rep.Unclear) > 0 {
		fmt.Fprintf(&b, "\nNO CLASIFICABLE POR RUTA (%d)\n", len(rep.Unclear))
		for _, u := range rep.Unclear {
			fmt.Fprintf(&b, "  %s — %s\n", u.File, u.Reason)
		}
		b.WriteString("  Ante la duda se convoca: si alguno de estos toca algo del piso, " +
			"el comando no lo sabe y vos sí podés abrirlo.\n")
	}

	return b.String()
}

// columnaDerecha es el hecho cuando el rol corresponde y el motivo cuando no.
func columnaDerecha(d Decision) string {
	if !d.Convened {
		return d.Reason
	}
	if len(d.Files) == 0 {
		return d.Reason
	}
	if d.ReasonCode == "floor" {
		// El motivo completo ya está en el encabezado PISO ACTIVADO; acá va
		// sólo el hecho, o la línea repite "piso: piso:" y se vuelve ilegible.
		return listar(d.Files) + "  [piso]"
	}
	return listar(d.Files)
}

func listar(files []string) string {
	if len(files) <= maxArchivos {
		return strings.Join(files, ", ")
	}
	return strings.Join(files[:maxArchivos], ", ") +
		fmt.Sprintf(" (+%d más)", len(files)-maxArchivos)
}

// ordenar pone los convocados primero, conservando el orden canónico dentro de
// cada grupo. Determinista a propósito: el informe se pega en el registro del
// reparto, y dos corridas sobre el mismo diff tienen que dar el mismo texto.
func ordenar(ds []Decision) []Decision {
	out := make([]Decision, 0, len(ds))
	for _, d := range ds {
		if d.Convened {
			out = append(out, d)
		}
	}
	for _, d := range ds {
		if !d.Convened {
			out = append(out, d)
		}
	}
	return out
}

// RenderJSON arma la salida para el orquestador, que la consume sin parsear
// texto. Cada decisión viaja con `reason_code`, estable, además de la prosa:
// un consumidor que tuviera que parsear el español volvería a depender de la
// redacción, que es lo que el JSON viene a evitar.
func RenderJSON(rep Report) ([]byte, error) {
	if rep.Files == nil {
		rep.Files = []string{}
	}
	if rep.Unclear == nil {
		rep.Unclear = []Unclear{}
	}
	out, err := json.MarshalIndent(rep, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("serializar el informe de roles: %w", err)
	}
	return append(out, '\n'), nil
}
