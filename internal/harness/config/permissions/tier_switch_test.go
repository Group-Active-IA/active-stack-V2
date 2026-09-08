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
// El invariante no se cumple todavia para permission.bash de opencode: ver
// TestOpencodeBashResidueSurvivesTierSwitch, que documenta el residuo, y
// TestOpencodeEveryTierDeclaresBashWildcard, que sostiene lo unico que hoy lo
// vuelve inofensivo.

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

// ── opencode: permission.bash ───────────────────────────────────────────────

// TestOpencodeEveryTierDeclaresBashWildcard sostiene la unica razon por la que
// el residuo documentado abajo es inofensivo.
//
// permission.bash NO cumple el invariante "el archivo describe solo el tier
// nuevo": el merge fusiona el objeto en vez de reemplazarlo, asi que los
// patrones de un tier anterior siguen ahi. Hoy no hacen nada porque opencode es
// last-match-wins y el "*" del tier nuevo llega despues y los tapa.
//
// Esa inercia depende enteramente de que TODO tier declare "*". Un tier nuevo
// sin "*", o un "*" restringido a un subconjunto, activa el residuo en silencio
// y sin que falle nada. Este test es lo que impide que eso pase callado.
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
			t.Errorf("%s no declara \"*\" en permission.bash: sin ese comodin, los patrones "+
				"que dejo el tier anterior dejan de estar tapados y vuelven a aplicar. Ver "+
				"TestOpencodeBashResidueSurvivesTierSwitch", tier)
		}
	}
}

// TestOpencodeBashResidueSurvivesTierSwitch documenta un residuo REAL que este
// fix no elimina, para que quede visible y no se redescubra en seis meses.
//
// Instalar balanceado y despues estricto deja los cuatro allow de balanceado
// ("go test *", "go build *", "git status", "git diff *") dentro de
// permission.bash, porque MergeJSONObjects fusiona objetos en vez de
// reemplazarlos. El arreglo estructural es el centinela __replace__ sobre la
// clave bash; es un cambio mas grande que este fix y no entra aca.
//
// Mientras tanto el residuo es inerte por UNA razon y este test la fija: el "*"
// del tier nuevo queda DESPUES de los patrones heredados, y opencode es
// last-match-wins. Si algun dia el orden se invierte, esto se pone rojo.
func TestOpencodeBashResidueSurvivesTierSwitch(t *testing.T) {
	home := installTiers(t, openCodeAdapter(), model.TierBalanceado, model.TierEstricto)
	order := opencodeBashOrder(t, opencodePath(home))

	position := map[string]int{}
	for i, key := range order {
		position[key] = i
	}

	wildcard, hasWildcard := position["*"]
	if !hasWildcard {
		t.Fatalf("permission.bash quedo sin \"*\" tras balanceado → estricto: %v", order)
	}

	// El residuo existe. No es una regresion: es el estado conocido.
	residue := []string{"go test *", "go build *", "git status", "git diff *"}
	for _, key := range residue {
		at, present := position[key]
		if !present {
			t.Fatalf("el residuo de balanceado ya no esta: %q desaparecio de %v. Si permission.bash "+
				"paso a reemplazarse entero, borra este test y afirma el invariante fuerte", key, order)
		}
		if at > wildcard {
			t.Errorf("%q quedo DESPUES del \"*\" de estricto (posiciones %d > %d): con last-match-wins "+
				"el allow heredado de balanceado vuelve a aplicar", key, at, wildcard)
		}
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
