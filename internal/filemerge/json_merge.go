package filemerge

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func MergeJSONObjects(baseJSON []byte, overlayJSON []byte) ([]byte, error) {
	base, err := decodeOrderedOrEmpty(baseJSON)
	if err != nil {
		// Malformed base is treated as {} — overlay prevalence is safe and
		// preferable to aborting the install.
		base = newOrderedObject()
	}

	overlay, err := decodeOrderedOrEmpty(overlayJSON)
	if err != nil {
		return nil, fmt.Errorf("unmarshal overlay json: %w", err)
	}

	merged := mergeOrdered(base, overlay)

	// Marshal first, then Indent: MarshalIndent does not reach a custom
	// MarshalJSON, so indenting has to happen on the produced bytes.
	encoded, err := json.Marshal(merged)
	if err != nil {
		return nil, fmt.Errorf("marshal merged json: %w", err)
	}
	var pretty bytes.Buffer
	if err := json.Indent(&pretty, encoded, "", "  "); err != nil {
		return nil, fmt.Errorf("indent merged json: %w", err)
	}

	return append(pretty.Bytes(), '\n'), nil
}

// mergeOrdered merges overlay into base preserving key order.
//
// Same rules as the map version, plus one that carries security weight: a key
// the overlay introduces lands AFTER everything the base already had. opencode
// evaluates last-match-wins, so the deny floor only holds while it is the last
// thing that matches.
func mergeOrdered(base, overlay *orderedObject) *orderedObject {
	result := newOrderedObject()
	for _, key := range base.keys {
		value, _ := base.get(key)
		result.set(key, value)
	}

	for _, key := range overlay.keys {
		overlayValue, _ := overlay.get(key)

		// Replace sentinel: use the sentinel's value verbatim, no deep merge.
		if replacement, isSentinel := asOrderedSentinel(overlayValue); isSentinel {
			result.set(key, replacement)
			continue
		}

		baseValue, exists := result.get(key)
		if !exists {
			// Recurse into overlay objects even with no base value, so nested
			// __replace__ sentinels are unwrapped before reaching the output.
			if overlayObj, isObj := overlayValue.(*orderedObject); isObj {
				result.set(key, mergeOrdered(newOrderedObject(), overlayObj))
			} else {
				result.set(key, overlayValue)
			}
			continue
		}

		baseObj, baseIsObj := baseValue.(*orderedObject)
		overlayObj, overlayIsObj := overlayValue.(*orderedObject)
		if baseIsObj && overlayIsObj {
			result.set(key, mergeOrdered(baseObj, overlayObj))
			continue
		}

		result.set(key, overlayValue)
	}

	return result
}

// asOrderedSentinel reports whether value is {"__replace__": X} and returns X.
func asOrderedSentinel(value any) (any, bool) {
	obj, ok := value.(*orderedObject)
	if !ok || len(obj.keys) != 1 || obj.keys[0] != replacesentinel {
		return nil, false
	}
	replacement, _ := obj.get(replacesentinel)
	return replacement, true
}

func unmarshalJSONObject(raw []byte) (map[string]any, error) {
	object := map[string]any{}
	if len(bytes.TrimSpace(raw)) == 0 {
		return object, nil
	}

	if err := json.Unmarshal(raw, &object); err == nil {
		return object, nil
	}

	normalized := normalizeJSON(raw)
	if err := json.Unmarshal(normalized, &object); err != nil {
		return nil, err
	}

	return object, nil
}

func normalizeJSON(raw []byte) []byte {
	withoutComments := stripJSONComments(raw)
	return stripTrailingCommas(withoutComments)
}

func stripJSONComments(raw []byte) []byte {
	out := make([]byte, 0, len(raw))
	inString := false
	escaped := false
	inLineComment := false
	inBlockComment := false

	for i := 0; i < len(raw); i++ {
		ch := raw[i]

		if inLineComment {
			if ch == '\n' {
				inLineComment = false
				out = append(out, ch)
			}
			continue
		}

		if inBlockComment {
			if ch == '*' && i+1 < len(raw) && raw[i+1] == '/' {
				inBlockComment = false
				i++
			}
			continue
		}

		if inString {
			out = append(out, ch)
			if escaped {
				escaped = false
				continue
			}
			if ch == '\\' {
				escaped = true
				continue
			}
			if ch == '"' {
				inString = false
			}
			continue
		}

		if ch == '"' {
			inString = true
			out = append(out, ch)
			continue
		}

		if ch == '/' && i+1 < len(raw) {
			next := raw[i+1]
			if next == '/' {
				inLineComment = true
				i++
				continue
			}
			if next == '*' {
				inBlockComment = true
				i++
				continue
			}
		}

		out = append(out, ch)
	}

	return out
}

func stripTrailingCommas(raw []byte) []byte {
	out := make([]byte, 0, len(raw))
	inString := false
	escaped := false

	for i := 0; i < len(raw); i++ {
		ch := raw[i]

		if inString {
			out = append(out, ch)
			if escaped {
				escaped = false
				continue
			}
			if ch == '\\' {
				escaped = true
				continue
			}
			if ch == '"' {
				inString = false
			}
			continue
		}

		if ch == '"' {
			inString = true
			out = append(out, ch)
			continue
		}

		if ch == ',' {
			j := i + 1
			for j < len(raw) {
				next := raw[j]
				if next == ' ' || next == '\t' || next == '\n' || next == '\r' {
					j++
					continue
				}
				if next == '}' || next == ']' {
					ch = 0
				}
				break
			}
		}

		if ch != 0 {
			out = append(out, ch)
		}
	}

	return out
}

// replacesentinel is the key used in an overlay map to signal that the parent
// key should be replaced atomically rather than deep-merged. When mergeObjects
// encounters a nested map whose only key is "__replace__", the value stored
// under that key is used verbatim as the replacement — the corresponding base
// value is discarded entirely.
const replacesentinel = "__replace__"

// asSentinel checks if v is a map with exactly one key "__replace__".
// If so, it returns the replacement value and true. Otherwise it returns nil, false.
func asSentinel(v any) (any, bool) {
	m, isMap := v.(map[string]any)
	if !isMap {
		return nil, false
	}
	if replacement, hasSentinel := m[replacesentinel]; hasSentinel && len(m) == 1 {
		return replacement, true
	}
	return nil, false
}

func mergeObjects(base map[string]any, overlay map[string]any) map[string]any {
	result := make(map[string]any, len(base)+len(overlay))
	for key, value := range base {
		result[key] = value
	}

	for key, overlayValue := range overlay {
		// Check for the replace sentinel: if the overlay value is a map with
		// exactly one key "__replace__", use the sentinel's value verbatim.
		if replacement, isSentinel := asSentinel(overlayValue); isSentinel {
			result[key] = replacement
			continue
		}

		baseValue, ok := result[key]
		if !ok {
			// Even when there is no base value, recurse into overlay maps so
			// that any nested __replace__ sentinels are unwrapped before
			// they reach the output.
			if overlayMap, isMap := overlayValue.(map[string]any); isMap {
				result[key] = mergeObjects(map[string]any{}, overlayMap)
			} else {
				result[key] = overlayValue
			}
			continue
		}

		baseMap, baseIsMap := baseValue.(map[string]any)
		overlayMap, overlayIsMap := overlayValue.(map[string]any)
		if baseIsMap && overlayIsMap {
			result[key] = mergeObjects(baseMap, overlayMap)
			continue
		}

		result[key] = overlayValue
	}

	return result
}
