package filemerge

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

// orderedObject is a JSON object that remembers the order its keys appeared in.
//
// Why this exists: json.Marshal sorts map keys alphabetically, and for at least
// one consumer that order IS the semantics. opencode's `permission.bash` block
// is evaluated last-match-wins, so a deny rule only holds while nothing that
// matches comes after it.
//
// On 2026-09-08 the overlay built that block with the deny rules last, using a
// custom MarshalJSON declared in overlays.go as "the ONLY safe way to guarantee
// deny-last" — and then the merge round-tripped it through map[string]any and
// alphabetised it. It worked by coincidence: "*" is 0x2A and the deny patterns
// start with "r". A user whose opencode.json already had `"sudo *": "allow"`
// ended up with that allow AFTER the deny floor, and won.
//
// Order preservation also keeps a reinstall from reshuffling the user's file,
// which made every reinstall diff unreadable.
type orderedObject struct {
	keys   []string
	values map[string]any
}

func newOrderedObject() *orderedObject {
	return &orderedObject{values: map[string]any{}}
}

func (o *orderedObject) get(key string) (any, bool) {
	v, ok := o.values[key]
	return v, ok
}

// set stores a value at the end, moving the key if it was already present.
//
// "Already present keeps its place" was the first version, and it was not enough:
// it guarantees "a new key lands last", which is NOT "the deny floor lands last".
// Installing estricto and then balanceado appended the new allow rules after the
// floor, and the floor stopped being the last thing that matches.
//
// Moving instead of updating in place makes the overlay's own order win, so a
// floor emitted last by overlays.go arrives last on disk whatever the user had.
// Reinstalling the same overlay is still byte-identical, because the resulting
// order only depends on the overlay — see TestMergeJSONObjects_ReinstallIsIdempotent.
func (o *orderedObject) set(key string, value any) {
	if _, exists := o.values[key]; exists {
		for i, k := range o.keys {
			if k == key {
				o.keys = append(o.keys[:i], o.keys[i+1:]...)
				break
			}
		}
	}
	o.keys = append(o.keys, key)
	o.values[key] = value
}

// setKeepingPlace stores a value without moving an existing key. Used when
// rebuilding the base side of a merge, where the base's own order is the
// starting point.
func (o *orderedObject) setKeepingPlace(key string, value any) {
	if _, exists := o.values[key]; !exists {
		o.keys = append(o.keys, key)
	}
	o.values[key] = value
}

// MarshalJSON emits the keys in insertion order.
func (o *orderedObject) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, key := range o.keys {
		if i > 0 {
			buf.WriteByte(',')
		}
		encodedKey, err := json.Marshal(key)
		if err != nil {
			return nil, err
		}
		buf.Write(encodedKey)
		buf.WriteByte(':')
		encodedValue, err := json.Marshal(o.values[key])
		if err != nil {
			return nil, err
		}
		buf.Write(encodedValue)
	}
	buf.WriteByte('}')
	return buf.Bytes(), nil
}

// decodeOrdered parses JSON keeping object key order. Objects become
// *orderedObject; everything else decodes as encoding/json normally would.
func decodeOrdered(raw []byte) (*orderedObject, error) {
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber() // no float64 round-trip: 1 must not come back as 1e+00

	tok, err := dec.Token()
	if err != nil {
		return nil, err
	}
	delim, ok := tok.(json.Delim)
	if !ok || delim != '{' {
		return nil, fmt.Errorf("expected a JSON object, got %v", tok)
	}
	return decodeOrderedBody(dec)
}

// decodeOrderedBody reads an object's entries; the opening brace is consumed.
func decodeOrderedBody(dec *json.Decoder) (*orderedObject, error) {
	out := newOrderedObject()
	for {
		tok, err := dec.Token()
		if err != nil {
			return nil, err
		}
		if delim, ok := tok.(json.Delim); ok && delim == '}' {
			return out, nil
		}
		key, ok := tok.(string)
		if !ok {
			return nil, fmt.Errorf("expected an object key, got %v", tok)
		}
		value, err := decodeOrderedValue(dec)
		if err != nil {
			return nil, err
		}
		out.setKeepingPlace(key, value)
	}
}

func decodeOrderedValue(dec *json.Decoder) (any, error) {
	tok, err := dec.Token()
	if err != nil {
		return nil, err
	}
	delim, isDelim := tok.(json.Delim)
	if !isDelim {
		return tok, nil
	}
	switch delim {
	case '{':
		return decodeOrderedBody(dec)
	case '[':
		items := []any{}
		for {
			next, err := dec.Token()
			if err != nil {
				return nil, err
			}
			if d, ok := next.(json.Delim); ok && d == ']' {
				return items, nil
			}
			// Push the token back by decoding from it: arrays can hold objects.
			item, err := valueFromToken(next, dec)
			if err != nil {
				return nil, err
			}
			items = append(items, item)
		}
	default:
		return nil, fmt.Errorf("unexpected delimiter %v", delim)
	}
}

// valueFromToken resolves a token already read from the stream.
func valueFromToken(tok json.Token, dec *json.Decoder) (any, error) {
	delim, isDelim := tok.(json.Delim)
	if !isDelim {
		return tok, nil
	}
	switch delim {
	case '{':
		return decodeOrderedBody(dec)
	case '[':
		items := []any{}
		for {
			next, err := dec.Token()
			if err != nil {
				return nil, err
			}
			if d, ok := next.(json.Delim); ok && d == ']' {
				return items, nil
			}
			item, err := valueFromToken(next, dec)
			if err != nil {
				return nil, err
			}
			items = append(items, item)
		}
	default:
		return nil, fmt.Errorf("unexpected delimiter %v", delim)
	}
}

// decodeOrderedOrEmpty parses raw, tolerating comments and trailing commas, and
// returns an empty object for empty input. Mirrors unmarshalJSONObject.
func decodeOrderedOrEmpty(raw []byte) (*orderedObject, error) {
	if len(bytes.TrimSpace(raw)) == 0 {
		return newOrderedObject(), nil
	}
	if obj, err := decodeOrdered(raw); err == nil {
		return obj, nil
	}
	obj, err := decodeOrdered(normalizeJSON(raw))
	if err != nil {
		if err == io.EOF {
			return nil, fmt.Errorf("unexpected end of json")
		}
		return nil, err
	}
	return obj, nil
}
