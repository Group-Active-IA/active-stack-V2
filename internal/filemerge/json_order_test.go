package filemerge

import (
	"strings"
	"testing"
)

// ──────────────────────────────────────────────────────────────────────────────
// El merge tiene que preservar el orden de las claves
//
// Encontrado el 2026-09-08 por dos revisiones independientes sobre
// permissions/overlays.go, y confirmado corriendo el merge.
//
// El overlay de opencode construye su bloque `bash` con las reglas de denegacion
// AL FINAL, porque opencode evalua last-match-wins: la ultima regla que matchea
// es la que gana. Para eso overlays.go tiene un MarshalJSON custom sobre un slice
// de pares, con diez lineas de comentario declarandolo "the ONLY safe way to
// guarantee deny-last".
//
// Ese mecanismo no llegaba al disco. MergeJSONObjects deserializaba a
// map[string]any y json.MarshalIndent ordena las claves de un mapa
// alfabeticamente, asi que el orden lo decidia el ASCII del patron. Funcionaba de
// casualidad: "*" es 0x2A y los deny empiezan con "r".
//
// Un usuario con "sudo *": "allow" en su opencode.json terminaba con ese allow
// DESPUES del piso de denegacion, y con last-match-wins el allow ganaba. Ningun
// test se ponia rojo, porque el que custodiaba el orden leia el archivo escrito y
// validaba la coincidencia alfabetica.
// ──────────────────────────────────────────────────────────────────────────────

// TestMergeJSONObjects_PreservesKeyOrder asserts merged output keeps insertion
// order instead of sorting alphabetically.
func TestMergeJSONObjects_PreservesKeyOrder(t *testing.T) {
	// Claves elegidas para que el orden alfabetico NO coincida con el de entrada:
	// si el resultado sale ordenado, "alpha" aparece antes que "zeta".
	overlay := []byte(`{"zeta":1,"alpha":2,"mid":3}`)

	out, err := MergeJSONObjects(nil, overlay)
	if err != nil {
		t.Fatalf("MergeJSONObjects: %v", err)
	}
	s := string(out)

	posZeta, posAlpha := strings.Index(s, `"zeta"`), strings.Index(s, `"alpha"`)
	if posZeta > posAlpha {
		t.Errorf("las claves salieron alfabetizadas, no en orden de entrada:\n%s", s)
	}
}

// TestMergeJSONObjects_OverlayKeysLandAfterBaseKeys asserts a key the overlay
// introduces ends up after the keys the user already had.
//
// This is the security-relevant case: opencode's deny floor only holds if it is
// the last thing that matches.
func TestMergeJSONObjects_OverlayKeysLandAfterBaseKeys(t *testing.T) {
	// Lo que el usuario ya tenia.
	base := []byte(`{"permission":{"bash":{"sudo *":"allow"}}}`)
	// Lo que emite composeOpencodeOverlay: el deny va ultimo, por construccion.
	overlay := []byte(`{"permission":{"bash":{"*":"ask","rm -rf *":"deny","rm -rf /*":"deny"}}}`)

	out, err := MergeJSONObjects(base, overlay)
	if err != nil {
		t.Fatalf("MergeJSONObjects: %v", err)
	}
	s := string(out)

	posSudo := strings.Index(s, `"sudo *"`)
	posDeny := strings.Index(s, `"rm -rf *"`)
	if posSudo < 0 || posDeny < 0 {
		t.Fatalf("falta alguna de las dos claves en el resultado:\n%s", s)
	}
	if posDeny < posSudo {
		t.Errorf("el piso de denegacion quedo ANTES de una regla del usuario.\n"+
			"Con last-match-wins gana el allow del usuario y el piso no existe:\n%s", s)
	}
}

// TestMergeJSONObjects_OverlayKeyMovesToTheEnd asserts a key present in both
// takes the overlay's value AND the overlay's position, at the end.
//
// The first version of this test asserted the opposite: that an existing key
// kept its place, so a reinstall would not churn the user's file. That was the
// nicer property and it was wrong. Keeping the place guarantees "a NEW key lands
// last", which is not "the deny floor lands last": switching from estricto to
// balanceado appended the new allow rules after the floor.
//
// Order is legibility. The floor being last is a security property. When they
// disagree, legibility loses — and the file barely churns anyway, because
// reinstalling the same overlay is byte-identical.
func TestMergeJSONObjects_OverlayKeyMovesToTheEnd(t *testing.T) {
	// Claves cuyo orden alfabetico no coincide con el de insercion, para que el
	// test no pueda pasar por coincidencia.
	base := []byte(`{"zebra":"a","mid":"b","alpha":"c"}`)
	overlay := []byte(`{"mid":"CHANGED"}`)

	out, err := MergeJSONObjects(base, overlay)
	if err != nil {
		t.Fatalf("MergeJSONObjects: %v", err)
	}
	s := string(out)

	if !strings.Contains(s, `"CHANGED"`) {
		t.Errorf("el overlay no piso el valor:\n%s", s)
	}
	pz, pa, pm := strings.Index(s, `"zebra"`), strings.Index(s, `"alpha"`), strings.Index(s, `"mid"`)
	if !(pz < pa && pa < pm) {
		t.Errorf("la clave del overlay no quedo al final.\n"+
			"Sin eso el piso de denegacion puede quedar tapado por lo que el usuario ya tenia:\n%s", s)
	}
}

// TestMergeJSONObjects_TierChangeKeepsDenyLast asserts the deny floor is still
// the last thing that matches after switching tiers.
//
// The first fix (d2fe658) guaranteed "a new key lands at the end", which is NOT
// the same as "the deny floor lands at the end". Installing estricto and then
// balanceado appended the new allow rules AFTER the floor:
//
//	"*", "rm -rf *", "rm -rf /*", "go test *", "git diff *"
//
// It did not blow up only because "go test *" does not overlap with "rm -rf".
// Right for the wrong reason, which is the same family of failure as the bug it
// was fixing. Found by the auditor, confirmed by running the merge.
func TestMergeJSONObjects_TierChangeKeepsDenyLast(t *testing.T) {
	estricto := []byte(`{"permission":{"bash":{"*":"ask","rm -rf *":"deny","rm -rf /*":"deny"}}}`)
	installed, err := MergeJSONObjects(nil, estricto)
	if err != nil {
		t.Fatalf("primera instalacion: %v", err)
	}

	balanceado := []byte(`{"permission":{"bash":{"*":"ask","go test *":"allow","git diff *":"allow","rm -rf *":"deny","rm -rf /*":"deny"}}}`)
	out, err := MergeJSONObjects(installed, balanceado)
	if err != nil {
		t.Fatalf("cambio de tier: %v", err)
	}
	s := string(out)

	lastDeny := strings.LastIndex(s, `"deny"`)
	lastAllow := strings.LastIndex(s, `"allow"`)
	if lastAllow > lastDeny {
		t.Errorf("tras cambiar de tier hay un allow DESPUES del piso de denegacion.\n"+
			"Con last-match-wins el piso dejo de ser lo ultimo que matchea:\n%s", s)
	}
}

// TestMergeJSONObjects_ReinstallIsIdempotent asserts installing the same overlay
// twice produces the same bytes.
//
// This is what keeps "overlay keys move to the end" from churning the file: with
// an unchanged overlay the result is stable. It only reorders when the overlay
// itself changed, which is exactly when reordering is what you want.
func TestMergeJSONObjects_ReinstallIsIdempotent(t *testing.T) {
	base := []byte(`{"permission":{"bash":{"sudo *":"allow"}}}`)
	overlay := []byte(`{"permission":{"bash":{"*":"ask","rm -rf *":"deny"}}}`)

	once, err := MergeJSONObjects(base, overlay)
	if err != nil {
		t.Fatalf("primera: %v", err)
	}
	twice, err := MergeJSONObjects(once, overlay)
	if err != nil {
		t.Fatalf("segunda: %v", err)
	}
	if string(once) != string(twice) {
		t.Errorf("reinstalar cambio el archivo.\nprimera:\n%s\nsegunda:\n%s", once, twice)
	}
}
