package reports

import (
	"bufio"
	"encoding/json"
	"strings"
)

// ContractCodexExec132 fija el formato público JSONL de codex exec 0.132.0.
const ContractCodexExec132 = "codex-exec-0.132.0-jsonl-v1"

const ContractClaude271 = "claude-2.1.271-hooks-jsonl-v1"

func init() { RegisterContract(ContractCodexExec132); RegisterContract(ContractClaude271) }

// ResolveProvider consume un archivo por despacho, obtenido de stdout de
// codex exec --json (no los rollouts privados). expectedSession debe venir
// del registro del despacho. El formato no proporciona turn ID ni modelo.
// Claude usa su contrato de hooks. OpenCode no tiene contrato habilitado.
func ResolveProvider(provider, version, expectedSession string, raw []byte, expectedAgent ...string) Source {
	if provider == "claude" && version == "2.1.271" {
		return resolveClaude(expectedSession, raw, expectedAgent)
	}
	if provider != "codex" || version != "0.132.0" {
		return Source{}
	}
	s := Source{Contract: ContractCodexExec132, Agent: "codex", Incomplete: true}
	if len(raw) > MaxContentBytes {
		return s
	}
	scanner := bufio.NewScanner(strings.NewReader(string(raw)))
	scanner.Buffer(make([]byte, 4096), MaxContentBytes+1)
	state := 0
	messages := 0
	var original []string
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		var e struct {
			Type     string           `json:"type"`
			Usage    *json.RawMessage `json:"usage"`
			ThreadID string           `json:"thread_id"`
			Item     *struct {
				ID   string  `json:"id"`
				Type string  `json:"type"`
				Text *string `json:"text"`
			} `json:"item"`
		}
		if json.Unmarshal([]byte(line), &e) != nil {
			return s
		}
		switch e.Type {
		case "thread.started":
			if state != 0 || e.ThreadID == "" {
				return s
			}
			s.Session = e.ThreadID
			state = 1
			original = append(original, line)
		case "turn.started":
			if state != 1 {
				return s
			}
			state = 2
			original = append(original, line)
		case "item.completed", "item.started", "item.updated":
			if state != 2 || e.Item == nil {
				return s
			}
			if e.Type == "item.completed" && e.Item.Type == "agent_message" {
				if e.Item.ID == "" || e.Item.Text == nil {
					return s
				}
				messages++
				s.Text = *e.Item.Text
				original = append(original, line)
			}
		case "turn.completed":
			if state != 2 || e.Usage == nil {
				return s
			}
			state = 3
			original = append(original, line)
		default:
			// Error, turn.failed y eventos desconocidos nunca acreditan cierre.
			return s
		}
	}
	s.Raw = strings.Join(original, "\n") + "\n"
	s.Incomplete = scanner.Err() != nil || state != 3 || messages != 1 || expectedSession == "" || s.Session != expectedSession
	return s
}
