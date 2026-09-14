// Package roles convierte una decisión en un dato.
//
// El bloque del orquestador dice "convocá al revisor que la capa toca". Está
// bien escrito y no funciona, por un motivo medido sobre tres corridas reales:
// el que decide si una capa "se toca" es el mismo que se ahorra el trabajo de
// convocar, y nadie revisa esa decisión. De diez roles disponibles el
// orquestador convocó siempre los mismos cuatro, y el revisor senior no corrió
// ni una vez en tres corridas. No convocar no deja rastro, así que nadie se
// entera.
//
// El contrato de este paquete es la tabla hecho-del-diff → revisor de
// `internal/harness/config/assets/roles.md`, y `## The Dispatch` del base es lo
// que obliga a justificar cada convocatoria con un hecho del diff. Esto produce
// ese hecho.
//
// NO CONVOCA A NADIE. Devuelve la lista. La convocatoria sigue siendo del
// orquestador, que es el único que puede hablar con la persona.
//
// LA REGLA QUE GOBIERNA TODO EL DISEÑO: ANTE LA DUDA, CONVOCA.
//
// La clasificación es por ruta, y toda heurística por ruta se equivoca. Lo que
// decide el diseño es hacia qué lado se equivoca: un revisor de más cuesta
// tokens; un revisor de menos cuesta un bug en producción. Por eso las listas
// de abajo son generosas y no finas, y por eso un empate se resuelve
// convocando. Si dentro de seis meses alguien viene a "optimizar" esto
// recortando términos, está cambiando el costo de un error de tokens a
// producción, y eso es la decisión que se está tomando.
package roles

import (
	"path"
	"strings"
)

// Los siete roles, con el identificador del sub-agente que los implementa.
const (
	RoleDBA      = "role-dba"
	RoleBackend  = "role-backend"
	RoleFrontend = "role-frontend"
	RoleDevOps   = "role-devops"
	RoleQA       = "role-qa"
	RoleArch     = "role-architect"
	RoleAnalyst  = "role-analyst"
)

// ordenCanonico es el orden en que los roles salen en el informe dentro de cada
// grupo (convocados primero, después los que no). Fijo a propósito: dos
// corridas sobre el mismo diff tienen que devolver el mismo texto, o el informe
// no sirve para pegarlo en un registro y compararlo después.
var ordenCanonico = []string{
	RoleDBA, RoleBackend, RoleFrontend, RoleDevOps, RoleQA, RoleArch, RoleAnalyst,
}

// Decision es la respuesta del comando sobre un rol: convocado o no, y por qué.
//
// Los dos casos salen en el informe. Hoy un rol que no se convoca simplemente
// no aparece en ningún lado, y esa ausencia es indistinguible de haberlo
// olvidado: que salga con el motivo es la mitad del valor del comando.
type Decision struct {
	Role       string   `json:"role"`
	Convened   bool     `json:"convened"`
	Reason     string   `json:"reason"`
	ReasonCode string   `json:"reason_code"`
	Files      []string `json:"files,omitempty"`
}

// Unclear es la tercera respuesta: ni "corresponde" ni "no corresponde", sino
// "no lo puedo saber".
//
// Se reserva para archivos de código cuyo NOMBRE es genérico (utils, helpers,
// common). Esa restricción es deliberada: si cada archivo sin rol saliera con
// reserva, la sección aparecería en toda corrida y se saltearía como se saltea
// todo lo que aparece siempre. El precio de la restricción se declara acá y no
// se esconde: un `parser.go` con la validación de tokens adentro no cae en esta
// red, y el comando no puede verlo.
type Unclear struct {
	File   string `json:"file"`
	Reason string `json:"reason"`
}

// Floor registra si el diff activó el piso y qué archivo lo activó.
type Floor struct {
	Triggered bool     `json:"triggered"`
	File      string   `json:"file,omitempty"`
	Fact      string   `json:"fact,omitempty"`
	Files     []string `json:"files,omitempty"`
}

// Report es el resultado de clasificar un diff.
type Report struct {
	Files     []string   `json:"files"`
	Floor     Floor      `json:"floor"`
	Decisions []Decision `json:"roles"`
	Unclear   []Unclear  `json:"unclear"`
}

// Convened devuelve los identificadores de los roles convocados, en el orden
// canónico del informe.
func (r Report) Convened() []string {
	var out []string
	for _, d := range r.Decisions {
		if d.Convened {
			out = append(out, d.Role)
		}
	}
	return out
}

// Decision devuelve la decisión tomada para un rol.
func (r Report) Decision(role string) (Decision, bool) {
	for _, d := range r.Decisions {
		if d.Role == role {
			return d, true
		}
	}
	return Decision{}, false
}

// ── El piso ───────────────────────────────────────────────────────────────────
//
// Cuando el diff toca autenticación, pagos, seguridad, logs de auditoría o
// datos de usuario, se convoca el piso entero: backend, QA, DBA y arquitecto.
//
// ESE CASO NO SE RELAJA POR TAMAÑO. Un diff de tres líneas que toca auth.py
// convoca el piso igual que uno de trescientas. Un piso que se salta cuando el
// cambio parece menor se salta exactamente en el cambio donde alguien tenía
// apuro, y esa es la única corrida en la que importaba. Por eso en `Classify`
// no hay ninguna rama que mire len(files) antes de decidir el piso.
var rolesDelPiso = []string{RoleBackend, RoleQA, RoleDBA, RoleArch}

// hechosDelPiso son los tokens de ruta que activan el piso, agrupados por el
// hecho que se reporta. Se comparan como palabra completa contra los tokens de
// la ruta, no como subcadena: "auth" como subcadena convierte cada `authors.go`
// de un blog en un cambio de autenticación, y un piso que se activa siempre es
// un piso que se ignora.
var hechosDelPiso = []struct {
	hecho  string
	tokens []string
}{
	{"autenticación", []string{
		"auth", "authn", "authz", "authentication", "authorization", "authenticate",
		"login", "logout", "signin", "signup", "oauth", "oauth2", "oidc", "saml",
		"jwt", "session", "sessions", "rbac", "acl", "permission", "permissions",
		"identity", "sso", "mfa", "otp", "2fa", "password", "passwords", "passwd",
	}},
	{"secretos o variables de entorno", []string{
		"env", "dotenv", "secret", "secrets", "credential", "credentials",
		"vault", "keyring", "keystore", "apikey", "privatekey", "token", "tokens",
	}},
	{"pagos", []string{
		"billing", "payment", "payments", "pay", "checkout", "invoice", "invoices",
		"stripe", "paypal", "subscription", "subscriptions", "charge", "charges",
		"refund", "refunds", "payout", "merchant", "pci",
	}},
	{"datos de usuario", []string{
		"user", "users", "account", "accounts", "profile", "profiles",
		"customer", "customers", "member", "members", "pii", "gdpr", "personal",
	}},
	{"seguridad o logs de auditoría", []string{
		"security", "audit", "auditlog", "auditlogs",
	}},
}

// tokensDBA: migración, esquema o consulta.
var tokensDBA = []string{
	"migration", "migrations", "migrate", "schema", "schemas", "alembic",
	"flyway", "liquibase", "seeds", "seed", "fixtures", "query", "queries",
	"repository", "repositories", "dao", "models", "db", "database", "sqlc",
}

// extDBA: extensiones que son esquema o consulta por sí solas.
var extDBA = []string{".sql", ".prisma"}

// tokensFrontend: el árbol de UI.
var tokensFrontend = []string{
	"ui", "tui", "frontend", "front", "web", "webapp", "client", "components",
	"component", "views", "view", "pages", "page", "styles", "style",
	"templates", "template", "public", "static", "widgets", "widget", "layouts",
	"layout", "screens", "screen",
}

var extFrontend = []string{
	".tsx", ".jsx", ".vue", ".svelte", ".css", ".scss", ".sass", ".less",
	".html", ".htm", ".astro",
}

// tokensBackend: un endpoint, un handler, una validación en el borde.
var tokensBackend = []string{
	"api", "apis", "endpoint", "endpoints", "handler", "handlers", "controller",
	"controllers", "route", "routes", "router", "routers", "server", "servers",
	"service", "services", "middleware", "middlewares", "resolver", "resolvers",
	"rpc", "grpc", "graphql", "rest", "validation", "validations", "validator",
	"validators", "serializer", "serializers", "usecase", "usecases",
}

// tokensDevOps: deploy, CI o contenedores.
var tokensDevOps = []string{
	"deploy", "deployment", "deployments", "ci", "cd", "pipeline", "pipelines",
	"docker", "dockerfile", "k8s", "kubernetes", "helm", "chart",
	"charts", "terraform", "ansible", "jenkins", "jenkinsfile", "makefile",
	"nginx", "systemd", "infra", "infrastructure", "provision", "release",
	"buildkite", "circleci", "travis", "vagrantfile", "procfile", "environments",
}

var extDevOps = []string{".tf", ".tfvars"}

// esCompose reconoce un archivo de compose de verdad.
//
// "compose" era un token suelto de DevOps hasta que la primera corrida sobre el
// diff real de este repo (2026-09-13) convocó a DevOps por
// `testdata/compose_all_toggles_claude.golden`. "Ante la duda, convoca" es sobre
// la duda; un falso positivo obvio que aparece en toda corrida no es duda, y un
// revisor que se convoca siempre deja de ser señal, que cuesta lo mismo que no
// convocarlo.
func esCompose(base, ext string) bool {
	if ext != ".yml" && ext != ".yaml" {
		return false
	}
	return base == "compose" || base == "docker-compose" || base == "podman-compose"
}

// prefijosDevOps son rutas que son CI por estar donde están, sin importar cómo
// se llame el archivo.
var prefijosDevOps = []string{
	".github/", ".gitlab/", ".circleci/", ".buildkite/", "ci/", "deploy/",
	"k8s/", "helm/", "charts/", "terraform/", "ansible/", "infra/",
}

// nombresGenericos son los basenames que admiten no decir qué hay adentro. Ver
// el comentario de `Unclear` para por qué la red es angosta a propósito.
var nombresGenericos = map[string]bool{
	"util": true, "utils": true, "helper": true, "helpers": true,
	"common": true, "commons": true, "misc": true, "shared": true,
	"core": true, "lib": true, "base": true, "internal": true, "mod": true,
	"aux": true, "stuff": true, "tools": true,
}

// "main" e "index" NO están arriba a propósito. Son nombres POSICIONALES —el
// punto de entrada, el barrel de un módulo—, no cajones de sastre, y aparecen
// en casi todo diff de un CLI o de un front. La primera corrida sobre el diff
// real de este repo (2026-09-13) marcó `cmd/active-stack/main.go` con reserva, y
// una reserva que sale siempre se saltea igual que cualquier cosa que aparece
// siempre. El precio se declara: un main.go con la validación de tokens adentro
// no cae en esta red.

// extCodigo: sólo un archivo de código puede esconder lógica que la ruta no
// anuncia. Un .md genérico no guarda la validación de tokens de nadie.
var extCodigo = []string{
	".go", ".py", ".js", ".ts", ".tsx", ".jsx", ".mjs", ".cjs", ".rb", ".java",
	".php", ".rs", ".kt", ".swift", ".cs", ".c", ".cc", ".cpp", ".h", ".hpp",
	".scala", ".ex", ".exs", ".dart", ".m", ".pl", ".lua", ".sh",
}

// Normalize deja una ruta comparable: barras normales y minúsculas.
//
// La barra invertida se reemplaza a mano y no con filepath.ToSlash: en Linux
// ToSlash es la identidad, así que una ruta de Windows pasada por --files
// llegaría entera a la clasificación y `web\components\Boton.tsx` sería un solo
// token ilegible. El CI cruza Linux y GOOS=windows justamente porque ya cazó
// bugs de rutas de esta familia.
func Normalize(p string) string {
	p = strings.ReplaceAll(p, "\\", "/")
	p = strings.TrimPrefix(p, "./")
	return strings.ToLower(p)
}

// tokens parte una ruta en palabras comparables: por separador, por puntuación
// y por camelCase.
func tokens(p string) []string {
	var out []string
	var actual strings.Builder
	runes := []rune(p)
	for i, r := range runes {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			// camelCase: la minúscula anterior con mayúscula actual ya se
			// resolvió abajo; acá sólo acumulamos.
			actual.WriteRune(r)
		case r >= 'A' && r <= 'Z':
			if i > 0 && runes[i-1] >= 'a' && runes[i-1] <= 'z' {
				out = appendTok(out, actual.String())
				actual.Reset()
			}
			actual.WriteRune(r - 'A' + 'a')
		default:
			out = appendTok(out, actual.String())
			actual.Reset()
		}
	}
	return appendTok(out, actual.String())
}

func appendTok(out []string, s string) []string {
	if s == "" {
		return out
	}
	return append(out, s)
}

func contiene(lista []string, s string) bool {
	for _, x := range lista {
		if x == s {
			return true
		}
	}
	return false
}

func algunToken(toks, lista []string) bool {
	for _, t := range toks {
		if contiene(lista, t) {
			return true
		}
	}
	return false
}

// Classify decide qué revisores corresponden para la lista de archivos dada.
func Classify(files []string) Report {
	rep := Report{Files: files}

	porRol := map[string][]string{}
	var pisoFiles []string

	for _, raw := range files {
		p := Normalize(raw)
		if p == "" {
			continue
		}
		toks := tokens(p)
		ext := path.Ext(p)
		base := strings.TrimSuffix(path.Base(p), ext)

		clasificado := false

		// El piso primero: no depende de nada más, y en particular no depende
		// del tamaño del diff.
		if hecho, ok := hechoDelPiso(toks); ok {
			if !rep.Floor.Triggered {
				rep.Floor.Triggered = true
				rep.Floor.File = raw
				rep.Floor.Fact = hecho
			}
			pisoFiles = append(pisoFiles, raw)
			clasificado = true
		}

		if algunToken(toks, tokensDBA) || contiene(extDBA, ext) {
			porRol[RoleDBA] = append(porRol[RoleDBA], raw)
			clasificado = true
		}
		if algunToken(toks, tokensFrontend) || contiene(extFrontend, ext) {
			porRol[RoleFrontend] = append(porRol[RoleFrontend], raw)
			clasificado = true
		}
		if algunToken(toks, tokensBackend) {
			porRol[RoleBackend] = append(porRol[RoleBackend], raw)
			clasificado = true
		}
		if algunToken(toks, tokensDevOps) || contiene(extDevOps, ext) ||
			tienePrefijoDevOps(p) || esCompose(base, ext) {
			porRol[RoleDevOps] = append(porRol[RoleDevOps], raw)
			clasificado = true
		}

		if !clasificado && contiene(extCodigo, ext) && nombresGenericos[base] {
			rep.Unclear = append(rep.Unclear, Unclear{
				File: raw,
				Reason: "nombre genérico: el comando clasifica por ruta y no puede ver " +
					"qué hay adentro. No dice que no corresponda, dice que no lo sabe",
			})
		}
	}

	rep.Floor.Files = pisoFiles
	rep.Decisions = decidir(rep.Floor, porRol, len(files) > 0)
	return rep
}

func tienePrefijoDevOps(p string) bool {
	for _, pre := range prefijosDevOps {
		if strings.HasPrefix(p, pre) {
			return true
		}
	}
	return false
}

func hechoDelPiso(toks []string) (string, bool) {
	for _, h := range hechosDelPiso {
		if algunToken(toks, h.tokens) {
			return h.hecho, true
		}
	}
	return "", false
}

// motivos describe, por rol, qué se dice cuando corresponde y qué cuando no.
var motivos = map[string]struct{ si, no, codigoSi, codigoNo string }{
	RoleDBA: {
		si: "migración, esquema o consulta", codigoSi: "schema",
		no: "no corresponde: ninguna migración, esquema ni consulta en el diff", codigoNo: "no-schema",
	},
	RoleBackend: {
		si: "endpoint, handler o validación en el borde", codigoSi: "endpoint",
		no: "no corresponde: ningún endpoint, handler ni validación en el borde", codigoNo: "no-endpoint",
	},
	RoleFrontend: {
		si: "archivo bajo el árbol de UI", codigoSi: "ui-tree",
		no: "no corresponde: ningún archivo bajo el árbol de UI", codigoNo: "no-ui-tree",
	},
	RoleDevOps: {
		si: "deploy, CI o contenedores", codigoSi: "deploy",
		no: "no corresponde: deploy y configuración sin cambios", codigoNo: "no-deploy",
	},
}

func decidir(piso Floor, porRol map[string][]string, hayArchivos bool) []Decision {
	out := make([]Decision, 0, len(ordenCanonico))

	for _, rol := range ordenCanonico {
		switch rol {
		case RoleDBA, RoleBackend, RoleFrontend, RoleDevOps:
			m := motivos[rol]
			archivos := porRol[rol]
			d := Decision{Role: rol, Files: archivos}
			switch {
			case piso.Triggered && contiene(rolesDelPiso, rol):
				d.Convened = true
				d.ReasonCode = "floor"
				d.Reason = motivoDelPiso(piso)
				d.Files = unir(archivos, piso.Files)
			case len(archivos) > 0:
				d.Convened = true
				d.ReasonCode = m.codigoSi
				d.Reason = m.si
			default:
				d.ReasonCode = m.codigoNo
				d.Reason = m.no
			}
			out = append(out, d)

		case RoleQA:
			// QA sale siempre que haya algo que revisar. En `roles.md` es el
			// "always" de verify y también la revisión chica de la última fila:
			// las dos lecturas terminan en QA.
			d := Decision{Role: rol}
			switch {
			case piso.Triggered:
				d.Convened, d.ReasonCode, d.Reason, d.Files = true, "floor", motivoDelPiso(piso), piso.Files
			case hayArchivos:
				d.Convened, d.ReasonCode, d.Reason = true, "always", "siempre"
			default:
				d.ReasonCode = "nothing-to-review"
				d.Reason = "no corresponde: el diff está vacío, no hay nada que revisar"
			}
			out = append(out, d)

		case RoleArch:
			d := Decision{Role: rol}
			if piso.Triggered {
				d.Convened, d.ReasonCode, d.Reason, d.Files = true, "floor", motivoDelPiso(piso), piso.Files
			} else {
				// Ni "corresponde" ni "no corresponde". Una frontera no se lee
				// por ruta: afirmar que no la hay sería afirmar algo que el
				// comando no miró.
				d.ReasonCode = "not-decided-by-diff"
				d.Reason = "no lo decide el diff: una frontera no se lee por ruta. " +
					"Se convoca por fase (propose siempre) o si el cambio movió una frontera"
			}
			out = append(out, d)

		case RoleAnalyst:
			out = append(out, Decision{
				Role:       rol,
				ReasonCode: "not-decided-by-diff",
				Reason: "no lo decide el diff: el analista funcional se convoca por fase " +
					"(explore y archive), no por un hecho del diff",
			})
		}
	}
	return out
}

func motivoDelPiso(piso Floor) string {
	return "piso: " + piso.Fact + " (" + piso.File + "). No se relaja por tamaño del diff"
}

func unir(a, b []string) []string {
	visto := map[string]bool{}
	var out []string
	for _, s := range append(append([]string{}, a...), b...) {
		if !visto[s] {
			visto[s] = true
			out = append(out, s)
		}
	}
	return out
}
