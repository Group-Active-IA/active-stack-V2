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

// TestMergeJSONObjects_ExistingKeysKeepTheirPlace asserts a key present in both
// keeps the base's position and takes the overlay's value.
//
// Without this the merge would shuffle the user's file on every reinstall, and
// the diff of a reinstall would be unreadable.
func TestMergeJSONObjects_ExistingKeysKeepTheirPlace(t *testing.T) {
	// Claves cuyo orden alfabetico NO coincide con el de insercion: si el
	// resultado saliera ordenado, el test pasaria sin probar nada. Con
	// first/second/third pasaba por coincidencia.
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
	pf, ps, pt := strings.Index(s, `"zebra"`), strings.Index(s, `"mid"`), strings.Index(s, `"alpha"`)
	if !(pf < ps && ps < pt) {
		t.Errorf("una clave que ya existia cambio de lugar:\n%s", s)
	}
}
