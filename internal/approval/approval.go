// Package approval lee, nunca escribe, una decisión humana ya tomada.
//
// POR QUÉ EXISTE (J1 revisada, propose 2026-10-23-fase-2-pasos-1-4.md §1).
// El plan original contemplaba una firma criptográfica (GPG/SSH) como prueba
// de que una persona aprobó un cambio. Juani la rechazó: un agente con acceso
// al repo puede imitar una tabla escrita tan fácil como puede imitar un
// commit firmado con llaves que también están en esa máquina. Lo que este
// paquete ofrece es CONTROL DE PROCESO — una tabla completa, vigente, en un
// commit documental separado del código — nunca IDENTIDAD ACREDITADA. Por eso
// no existe un estado VERIFIED: el código más alto que este paquete emite es
// WrittenUnverified, y su mensaje lo dice literal.
//
// Este paquete NUNCA ESCRIBE. No completa la tabla de Aprobación, no firma,
// no toca el árbol de trabajo. Es lectura: resuelve commits, lee contenido en
// revisiones de git y compara.
package approval

import "fmt"

// Status es el resultado de verificar una aprobación. El cero-valor es
// StatusMissing a propósito: un consumidor que ignora el error y lee el
// cero-valor de Result cae en el estado MÁS restrictivo, nunca en uno que se
// pueda confundir con aprobado. "Cualquier retorno no fatal es aprobación" es
// exactamente el error que este diseño hace imposible por default.
type Status int

const (
	// StatusMissing: no hay tabla de Aprobación, o no tiene ninguna fila de
	// datos. Código de salida 1.
	StatusMissing Status = iota
	// StatusInvalid: hay estructura pero algo no cierra — campo vacío, commit
	// que no resuelve, commit que mezcla aprobación con código, tabla que ya
	// estaba completa antes del commit de aprobación, referencia a P ausente
	// o que no es ancestro, informe referenciado faltante o alterado. Código 1.
	StatusInvalid
	// StatusStale: la propuesta cambió después del commit de aprobación. La
	// aprobación fue válida para OTRO contenido. Código 1.
	StatusStale
	// StatusWrittenUnverified: estructura completa, coincidente y vigente.
	// Código 2. NO es identidad acreditada — ver Result.Message.
	StatusWrittenUnverified
)

// ExitCode es el código de proceso que el CLI usa. 2 únicamente para
// WrittenUnverified; 1 para cualquier otro estado, incluido el cero-valor.
// No existe 0: este paquete nunca certifica éxito genérico.
func (s Status) ExitCode() int {
	if s == StatusWrittenUnverified {
		return 2
	}
	return 1
}

func (s Status) String() string {
	switch s {
	case StatusMissing:
		return "MISSING"
	case StatusInvalid:
		return "INVALID"
	case StatusStale:
		return "STALE"
	case StatusWrittenUnverified:
		return "WRITTEN_UNVERIFIED"
	default:
		return "MISSING"
	}
}

// writtenUnverifiedMessage es literal por contrato (§1 del propose): dice
// explícitamente que la estructura es válida y que la identidad humana NO
// está acreditada. Ningún consumidor debe poder confundir código 2 con
// identidad verificada, y la forma de evitarlo es que el propio texto lo diga.
const writtenUnverifiedMessage = "estructura válida; identidad humana no acreditada"

// Result es la salida de Verify. Nunca tiene un campo booleano "OK" u similar:
// forzar el switch sobre Status es la protección contra el consumidor que
// trata cualquier retorno no fatal como aprobación.
type Result struct {
	Status Status

	// Message es el texto de contrato para WrittenUnverified, vacío en el resto.
	Message string

	// Cause, Field y Recovery sólo se completan en los estados que bloquean
	// (Missing/Invalid/Stale): causa, archivo/campo afectado y recuperación
	// sugerida. Nunca vuelcan contenido original ni valores sensibles.
	Cause    string
	Field    string
	Recovery string

	// Scope es el texto de la columna Alcance. Sólo tiene sentido en
	// WrittenUnverified: "código 2 habilita únicamente el alcance escrito".
	Scope string

	// ProposalSHA es P: el commit de la propuesta revisada que C aprueba.
	// CommitSHA es C: el commit documental de aprobación, resuelto completo.
	ProposalSHA string
	CommitSHA   string

	// Manifest es el manifiesto determinista de archivos/hashes de la
	// propuesta y sus informes referenciados, en la revisión P.
	Manifest []ManifestEntry
}

// Options son los parámetros de Verify.
type Options struct {
	// RepoDir es el repositorio git a inspeccionar.
	RepoDir string
	// ProposalPath es la ruta del archivo de propuesta, relativa a RepoDir
	// (ej. "docs/propuestas/2026-10-23-fase-2-pasos-1-4.md").
	ProposalPath string
	// Commit es C: la referencia o SHA del commit de aprobación.
	Commit string
	// CurrentRef es contra qué revisión se compara la propuesta para detectar
	// STALE (cambios normativos posteriores). Vacío = HEAD.
	CurrentRef string
}

func blocked(status Status, cause, field, recovery string) Result {
	return Result{Status: status, Cause: cause, Field: field, Recovery: recovery}
}

// Verify lee, no escribe. Devuelve error sólo por fallos de infraestructura
// (git ausente, opciones vacías); cualquier desacuerdo de contrato vuelve
// como Result con su Status, nunca como error — approval es lectura y el
// dominio de fallas humanas no es una falla de Go.
func Verify(opts Options) (Result, error) {
	if opts.RepoDir == "" || opts.ProposalPath == "" || opts.Commit == "" {
		return Result{}, fmt.Errorf("approval.Verify: RepoDir, ProposalPath y Commit son obligatorios")
	}
	return verify(opts)
}
