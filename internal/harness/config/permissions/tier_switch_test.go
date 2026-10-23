package permissions_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Group-Active-IA/active-stack/internal/harness/config/permissions"
	"github.com/Group-Active-IA/active-stack/internal/model"
)

// Bajar de tier tiene que dejar el archivo describiendo SOLO el tier nuevo.
// Hasta el 2026-09-08 no habia un solo test que instalara un tier y despues
// otro distinto: cada caso arrancaba de un TempDir limpio. Por ese hueco vivio
// el bug de que el tier viejo sobrevivia al nuevo.
//
// El invariante tampoco se cumplia para permission.bash de opencode hasta
// 2026-10-23: TestOpencodeEveryTierDeclaresBashWildcard sostenia la unica
// razon por la que el residuo era inofensivo (el "*" del tier nuevo quedaba
// despues); el fix (centinela __replace__ sobre "bash") lo elimino del todo,
// ver TestOpencodeTierSwitchFullyReplacesBash.

// installTiers instala los tiers en orden sobre el mismo home y devuelve el home.
func installTiers(t *testing.T, adapter permissions.PermissionsAdapter, tiers ...model.PermissionTier) string {
	t.Helper()
	home := t.TempDir()
	for _, tier := range tiers {
		if _, err := permissions.Install(home, []permissions.PermissionsAdapter{adapter}, tier); err != nil {
			t.Fatalf("Install(%s) error = %v", tier, err)
		}
	}
	return home
}

func allTiers() []model.PermissionTier {
	return []model.PermissionTier{model.TierEstricto, model.TierBalanceado, model.TierBypass}
}

// ── Claude Code ─────────────────────────────────────────────────────────────

// TestClaudeTierDowngradeDropsPreviousAllowList: balanceado deja una allow-list
// que auto-aprueba Read/Edit/Bash(go test:*). Instalar estricto despues tiene
// que borrarla; si sobrevive, el usuario cree que endurecio y sigue igual.
func TestClaudeTierDowngradeDropsPreviousAllowList(t *testing.T) {
	home := installTiers(t, claudeAdapter(), model.TierBalanceado, model.TierEstricto)

	perms := readClaudePermissions(t, filepath.Join(home, ".claude", "settings.json"))

	if mode, _ := perms["defaultMode"].(string); mode != "default" {
		t.Errorf("defaultMode = %q; se esperaba \"default\" para estricto", mode)
	}
	// Sin early return ni assert dentro de un if: la clave tiene que estar Y
	// tiene que estar vacia. Un `if !present { return }` es la forma exacta de
	// pasar sin ejercitar nada.
	allow := claudeAllowEntries(t, perms)
	if len(allow) != 0 {
		t.Errorf("estricto quedo con la allow-list de balanceado viva: %#v", allow)
	}
}

// TestClaudeBalanceadoThenBypassDropsAllowList cubre el otro orden posible.
// (El nombre anterior decia "BypassThenEstricto" y era al reves de lo que hace.)
func TestClaudeBalanceadoThenBypassDropsAllowList(t *testing.T) {
	home := installTiers(t, claudeAdapter(), model.TierBalanceado, model.TierBypass)

	perms := readClaudePermissions(t, filepath.Join(home, ".claude", "settings.json"))

	if mode, _ := perms["defaultMode"].(string); mode != "bypassPermissions" {
		t.Errorf("defaultMode = %q; se esperaba \"bypassPermissions\"", mode)
	}
	allow := claudeAllowEntries(t, perms)
	if len(allow) != 0 {
		t.Errorf("bypass quedo con la allow-list de balanceado viva: %#v", allow)
	}
}

// TestClaudeTiersDeclareTheSameKeys: mismo invariante de forma para el bloque
// "permissions" de Claude Code.
func TestClaudeTiersDeclareTheSameKeys(t *testing.T) {
	want := []string{"defaultMode", "allow", "deny"}

	for _, tier := range allTiers() {
		home := installTiers(t, claudeAdapter(), tier)
		perms := readClaudePermissions(t, filepath.Join(home, ".claude", "settings.json"))

		if len(perms) != len(want) {
			t.Errorf("%s escribio %d claves en permissions, se esperaban %d: %#v", tier, len(perms), len(want), perms)
		}
		for _, key := range want {
			if _, present := perms[key]; !present {
				t.Errorf("%s no declara permissions.%q: al cambiar de tier queda el valor viejo", tier, key)
			}
		}
	}
}

// ── opencode: external_directory ────────────────────────────────────────────

// externalDirectoryByTier es la politica, escrita una sola vez. Cambiar un
// valor en overlays.go sin cambiarlo aca tiene que poner algo en rojo: es una
// clave de permisos, y "que no sea deny" no es un assert de politica.
var externalDirectoryByTier = map[model.PermissionTier]string{
	model.TierEstricto:   "deny",
	model.TierBalanceado: "ask",
	model.TierBypass:     "allow",
}

// TestOpencodeExternalDirectoryValuePerTier fija el valor exacto de cada tier
// en una instalacion limpia.
func TestOpencodeExternalDirectoryValuePerTier(t *testing.T) {
	for _, tier := range allTiers() {
		t.Run(string(tier), func(t *testing.T) {
			home := installTiers(t, openCodeAdapter(), tier)
			perm := readOpencodePermission(t, opencodePath(home))

			got, present := perm["external_directory"]
			if !present {
				t.Fatalf("%s no declara external_directory: al cambiar de tier queda el valor viejo", tier)
			}
			if got != externalDirectoryByTier[tier] {
				t.Errorf("%s: external_directory = %v, se esperaba %q", tier, got, externalDirectoryByTier[tier])
			}
		})
	}
}

// TestOpencodeTierSwitchAdoptsNewTierExternalDirectory: external_directory =
// deny lo escribia solo estricto, asi que quedaba pegado al pasar a otro tier.
// El assert es el valor del tier NUEVO, en todas las transiciones.
func TestOpencodeTierSwitchAdoptsNewTierExternalDirectory(t *testing.T) {
	for _, from := range allTiers() {
		for _, to := range allTiers() {
			if from == to {
				continue
			}
			t.Run(string(from)+"→"+string(to), func(t *testing.T) {
				home := installTiers(t, openCodeAdapter(), from, to)
				perm := readOpencodePermission(t, opencodePath(home))

				got, present := perm["external_directory"]
				if !present {
					t.Fatalf("tras %s → %s no quedo external_directory declarado", from, to)
				}
				if got != externalDirectoryByTier[to] {
					t.Errorf("tras %s → %s: external_directory = %v, manda %s con %q",
						from, to, got, to, externalDirectoryByTier[to])
				}
			})
		}
	}
}

// TestOpencodeTiersDeclareTheSameKeys es la red contra la reincidencia. Los
// casos de arriba cubren external_directory porque ya sabiamos que fallaba; este
// falla con la PROXIMA clave que alguien agregue en un solo tier. Un tier que
// omite una clave hereda la respuesta del tier anterior, asi que el invariante
// no es "external_directory esta bien": es "todos los tiers declaran lo mismo".
func TestOpencodeTiersDeclareTheSameKeys(t *testing.T) {
	reference := map[string]bool{}
	var referenceTier model.PermissionTier

	for i, tier := range allTiers() {
		home := installTiers(t, openCodeAdapter(), tier)
		perm := readOpencodePermission(t, opencodePath(home))

		keys := map[string]bool{}
		for key := range perm {
			keys[key] = true
		}
		if i == 0 {
			reference, referenceTier = keys, tier
			continue
		}
		for key := range reference {
			if !keys[key] {
				t.Errorf("%s no declara %q, que si declara %s: al cambiar de tier queda el valor viejo", tier, key, referenceTier)
			}
		}
		for key := range keys {
			if !reference[key] {
				t.Errorf("%s declara %q, que no declara %s: al cambiar de tier queda el valor viejo", tier, key, referenceTier)
			}
		}
	}
}

// TestOpencodeNeverWritesEmptyPermissionVerb es una regresion real encontrada
// en el uso real (binario compilado, home falso con opencode.json
// preexistente), no anticipada por ningun test anterior: la reescritura de
// composeOpencodeOverlay con un map[string]any en vez del struct con
// `omitempty` perdio ese omitempty para "task" (el unico campo que ningun
// tier fija), y escribio literalmente "task": "" en el archivo real. El
// comentario original de opencodePermissionShape ya advertia que el string
// vacio "seria peor que omitir la clave": "" no es un verbo valido de
// opencode. Ningun tier declara Task todavia, asi que el invariante que
// importa es "ninguna clave de permission vale string vacio", no un valor
// fijo para "task".
func TestOpencodeNeverWritesEmptyPermissionVerb(t *testing.T) {
	for _, tier := range allTiers() {
		t.Run(string(tier), func(t *testing.T) {
			home := installTiers(t, openCodeAdapter(), tier)
			perm := readOpencodePermission(t, opencodePath(home))

			for key, value := range perm {
				if key == "bash" {
					continue // objeto, no verbo escalar
				}
				if value == "" {
					t.Errorf("%s: permission.%s = \"\" (string vacio no es un verbo valido de opencode; "+
						"se esperaba que la clave no se escriba en absoluto)", tier, key)
				}
			}
		})
	}
}

// ── opencode: permission.bash ───────────────────────────────────────────────

// TestOpencodeEveryTierDeclaresBashWildcard ya no sostiene un residuo (el
// centinela __replace__ sobre "bash" lo elimino, ver
// TestOpencodeTierSwitchFullyReplacesBash): queda como defensa en profundidad,
// porque "*" sigue siendo lo que hace inofensivo CUALQUIER patron ajeno que
// una instalacion futura deje en el objeto que se reemplaza (por ejemplo, si
// algun dia un tier deja de reemplazar "bash" entero por error).
func TestOpencodeEveryTierDeclaresBashWildcard(t *testing.T) {
	for _, tier := range allTiers() {
		home := installTiers(t, openCodeAdapter(), tier)
		order := opencodeBashOrder(t, opencodePath(home))

		found := false
		for _, key := range order {
			if key == "*" {
				found = true
			}
		}
		if !found {
			t.Errorf("%s no declara \"*\" en permission.bash", tier)
		}
	}
}

// expectedBashPairsByTier es la politica de permission.bash, escrita una sola
// vez, en el mismo orden que composeOpencodeOverlay en overlays.go. Cambiar un
// patron ahi sin cambiarlo aca tiene que poner esto en rojo.
func expectedBashPairsByTier(tier model.PermissionTier) [][2]string {
	switch tier {
	case model.TierEstricto:
		return [][2]string{
			{"*", "ask"},
			{"rm -rf *", "deny"},
			{"rm -rf /*", "deny"},
		}
	case model.TierBypass:
		return [][2]string{
			{"*", "allow"},
			{"rm -rf /*", "deny"},
			{"rm -rf ~/", "deny"},
		}
	default: // TierBalanceado
		return [][2]string{
			{"*", "ask"},
			{"go test *", "allow"},
			{"go build *", "allow"},
			{"git status", "allow"},
			{"git diff *", "allow"},
			{"rm -rf *", "deny"},
			{"rm -rf /*", "deny"},
		}
	}
}

// TestOpencodeTierSwitchFullyReplacesBash es el invariante FUERTE que el test
// que este reemplaza ("TestOpencodeBashResidueSurvivesTierSwitch") pedia
// afirmar el dia que permission.bash pasara a reemplazarse entero: ese dia
// llego con el centinela __replace__ sobre la clave "bash" (igual al patron ya
// usado para las entradas MCP en internal/harness/external/mcp.go).
//
// Matriz 3x3 CON la diagonal (mismo tier → mismo tier incluido, a proposito:
// reinstalar el tier que ya esta instalado tiene que ser un no-op real, no
// solo "no empeora"). Para cada (from, to), instalar from y despues to tiene
// que dejar permission.bash describiendo EXACTAMENTE los pares de "to", en el
// mismo orden — ningun patron de "from" sobrevive, sin importar si "from" es
// mas permisivo o mas estricto que "to".
func TestOpencodeTierSwitchFullyReplacesBash(t *testing.T) {
	for _, from := range allTiers() {
		for _, to := range allTiers() {
			t.Run(string(from)+"→"+string(to), func(t *testing.T) {
				home := installTiers(t, openCodeAdapter(), from, to)

				raw, err := os.ReadFile(opencodePath(home))
				if err != nil {
					t.Fatalf("read opencode.json: %v", err)
				}
				got, err := extractBashOrderedPairs(raw)
				if err != nil {
					t.Fatalf("extractBashOrderedPairs: %v", err)
				}

				want := expectedBashPairsByTier(to)
				if len(got) != len(want) {
					t.Fatalf("%s → %s: permission.bash = %v, se esperaban exactamente %v (sin residuo de %s)",
						from, to, got, want, from)
				}
				for i := range want {
					if got[i] != want[i] {
						t.Errorf("%s → %s: permission.bash[%d] = %v, se esperaba %v", from, to, i, got[i], want[i])
					}
				}
			})
		}
	}
}

// TestOpencodeTierSwitchPreservesPersonalSiblingKeyInPermission es la
// diferencia deliberada contra la alternativa rechazada del propose (el fix
// ddeb94f de otra rama, que envolvia el objeto "permission" ENTERO en
// __replace__). Esta implementacion envuelve solo "permission.bash": una
// clave hermana de bash dentro de permission que ningun tier administra tiene
// que sobrevivir a cualquier instalacion/reinstalacion.
func TestOpencodeTierSwitchPreservesPersonalSiblingKeyInPermission(t *testing.T) {
	home := t.TempDir()
	settingsPath := opencodePath(home)
	if err := os.MkdirAll(filepath.Dir(settingsPath), 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	// "mcp_tool_x" no es una clave que ningun tier declare: es una regla
	// personal que el usuario (o un agent/MCP futuro) agrego a mano dentro de
	// "permission", hermana de "bash", no dentro de el.
	seed := `{"permission": {"mcp_tool_x": "allow", "bash": {"*": "ask"}}}`
	if err := os.WriteFile(settingsPath, []byte(seed), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	if _, err := permissions.Install(home, []permissions.PermissionsAdapter{openCodeAdapter()}, model.TierBypass); err != nil {
		t.Fatalf("Install(bypass) error = %v", err)
	}

	perm := readOpencodePermission(t, settingsPath)
	if perm["mcp_tool_x"] != "allow" {
		t.Errorf("mcp_tool_x = %v, se esperaba que sobreviva: reemplazar solo permission.bash no debe "+
			"tocar sus hermanas dentro de permission", perm["mcp_tool_x"])
	}
}

// TestOpencodeTierSwitchReplacesPersonalBashRule documenta la perdida
// ACEPTADA y deliberada: una regla que el usuario agrego a mano DENTRO de
// permission.bash (no hermana de bash, sino dentro del objeto que bash
// reemplaza entero) no sobrevive a una reinstalacion. El propose la acepta
// explicitamente ("Perdida deliberada de bash personal exige diff/aprobacion
// y snapshot") y deja pendiente el gate de mostrar el diff (J6, TBD) — ver
// deuda en el informe de apply. Este test fija el comportamiento actual para
// que no se lo lea como un bug nuevo si alguien lo redescubre.
func TestOpencodeTierSwitchReplacesPersonalBashRule(t *testing.T) {
	home := t.TempDir()
	settingsPath := opencodePath(home)
	if err := os.MkdirAll(filepath.Dir(settingsPath), 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	seed := `{"permission": {"bash": {"my-deploy-script.sh *": "allow", "*": "ask"}}}`
	if err := os.WriteFile(settingsPath, []byte(seed), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	if _, err := permissions.Install(home, []permissions.PermissionsAdapter{openCodeAdapter()}, model.TierBalanceado); err != nil {
		t.Fatalf("Install(balanceado) error = %v", err)
	}

	perm := readOpencodePermission(t, settingsPath)
	bash, ok := perm["bash"].(map[string]any)
	if !ok {
		t.Fatalf("permission.bash ausente o con tipo raro: %#v", perm["bash"])
	}
	if _, present := bash["my-deploy-script.sh *"]; present {
		t.Errorf("la regla personal sobrevivio: se esperaba que el reemplazo de bash la descarte " +
			"(perdida aceptada por el propose, no un bug)")
	}
}

// ── codex ───────────────────────────────────────────────────────────────────

// TestCodexEveryTierDeclaresNetworkAccess: network_access lo escribia solo
// balanceado, dentro de [sandbox_workspace_write], asi que al cambiar de tier
// quedaba en el archivo el valor de un tier que ya no estaba.
//
// Los tres tiers escriben false, asi que un test de "instalar A y despues B"
// seria vacuo: el residuo coincide con el valor esperado y pasaria con el bug
// presente. Lo que muerde es sembrar un valor que NINGUN tier produce y exigir
// que cada tier lo pise.
func TestCodexEveryTierDeclaresNetworkAccess(t *testing.T) {
	for _, tier := range allTiers() {
		t.Run(string(tier), func(t *testing.T) {
			path := filepath.Join(t.TempDir(), ".codex", "config.toml")
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				t.Fatal(err)
			}
			seed := "model = \"keep-me\"\n\n[sandbox_workspace_write]\nnetwork_access = true\n"
			if err := os.WriteFile(path, []byte(seed), 0o644); err != nil {
				t.Fatal(err)
			}

			installTiers(t, codexPermissionsAdapter{path}, tier)

			raw, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			got := string(raw)
			if !strings.Contains(got, "network_access = false") {
				t.Errorf("%s no piso network_access: quedo el valor que ya estaba en el archivo:\n%s", tier, got)
			}
			if strings.Contains(got, "network_access = true") {
				t.Errorf("%s dejo network_access = true en el archivo:\n%s", tier, got)
			}
			if !strings.Contains(got, `model = "keep-me"`) {
				t.Errorf("%s perdio config no relacionada:\n%s", tier, got)
			}
		})
	}
}

// ── helpers ─────────────────────────────────────────────────────────────────

func opencodePath(home string) string {
	return filepath.Join(home, ".config", "opencode", "opencode.json")
}

func readClaudePermissions(t *testing.T, path string) map[string]any {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	var settings map[string]any
	if err := json.Unmarshal(content, &settings); err != nil {
		t.Fatalf("unmarshal %s: %v", path, err)
	}
	perms, ok := settings["permissions"].(map[string]any)
	if !ok {
		t.Fatalf("nodo permissions ausente o con tipo raro: %#v", settings["permissions"])
	}
	return perms
}

// claudeAllowEntries exige que la clave exista y sea una lista, y devuelve sus
// entradas. Falla en vez de saltear: la clave presente en todo tier ES el fix.
func claudeAllowEntries(t *testing.T, perms map[string]any) []any {
	t.Helper()
	raw, present := perms["allow"]
	if !present {
		t.Fatalf("permissions.allow ausente: es la clave que pisa la allow-list del tier anterior")
	}
	list, ok := raw.([]any)
	if !ok {
		t.Fatalf("permissions.allow no es una lista: %#v", raw)
	}
	return list
}

func readOpencodePermission(t *testing.T, path string) map[string]any {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	var settings map[string]any
	if err := json.Unmarshal(content, &settings); err != nil {
		t.Fatalf("unmarshal %s: %v", path, err)
	}
	perm, ok := settings["permission"].(map[string]any)
	if !ok {
		t.Fatalf("nodo permission ausente o con tipo raro: %#v", settings["permission"])
	}
	return perm
}

// opencodeBashOrder devuelve las claves de permission.bash EN ORDEN. Hace falta
// leerlas del texto crudo porque json.Unmarshal a un map pierde el orden, y en
// opencode el orden es la semantica (last-match-wins).
func opencodeBashOrder(t *testing.T, path string) []string {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}

	marker := []byte(`"bash":`)
	start := bytes.Index(raw, marker)
	if start < 0 {
		t.Fatalf("permission.bash ausente en %s:\n%s", path, raw)
	}
	start += len(marker)
	for start < len(raw) && raw[start] != '{' {
		start++
	}
	depth, end := 0, -1
	for i := start; i < len(raw) && end < 0; i++ {
		switch raw[i] {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				end = i + 1
			}
		}
	}
	if end < 0 {
		t.Fatalf("bloque permission.bash sin cerrar en %s:\n%s", path, raw)
	}

	dec := json.NewDecoder(bytes.NewReader(raw[start:end]))
	if _, err := dec.Token(); err != nil { // '{'
		t.Fatalf("decode permission.bash: %v", err)
	}
	var keys []string
	for dec.More() {
		key, err := dec.Token()
		if err != nil {
			t.Fatalf("decode clave de permission.bash: %v", err)
		}
		name, ok := key.(string)
		if !ok {
			t.Fatalf("clave de permission.bash con tipo raro: %#v", key)
		}
		if _, err := dec.Token(); err != nil {
			t.Fatalf("decode valor de %q: %v", name, err)
		}
		keys = append(keys, name)
	}
	return keys
}
