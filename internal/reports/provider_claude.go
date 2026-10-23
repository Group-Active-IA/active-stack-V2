package reports

import (
	"bufio"
	"encoding/json"
	"strings"
)

// Sólo eventos explícitos de un despacho: nunca sigue transcript_path ni
// agent_transcript_path y nunca lee el HOME del proveedor.
func resolveClaude(session string, raw []byte, expectedAgent []string) Source {
	s := Source{Contract: ContractClaude271, Raw: string(raw), Incomplete: true}
	if len(raw) > MaxContentBytes || session == "" || len(expectedAgent) != 1 || expectedAgent[0] == "" {
		return s
	}
	scanner := bufio.NewScanner(strings.NewReader(string(raw)))
	scanner.Buffer(make([]byte, 4096), MaxContentBytes+1)
	state := 0
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		var e struct {
			Session string `json:"session_id"`
			Agent   string `json:"agent_id"`
			Event   string `json:"hook_event_name"`
			Tool    string `json:"tool_name"`
			ToolID  string `json:"tool_use_id"`
			Input   struct {
				Message *string `json:"message"`
			} `json:"tool_input"`
			Last       *string           `json:"last_assistant_message"`
			Active     bool              `json:"stop_hook_active"`
			Background []json.RawMessage `json:"background_tasks"`
			Crons      []json.RawMessage `json:"session_crons"`
		}
		if json.Unmarshal([]byte(line), &e) != nil || e.Session != session || e.Agent != expectedAgent[0] {
			return s
		}
		s.Session = e.Session
		s.Agent = e.Agent
		switch e.Event {
		case "PostToolUse":
			if state != 0 || e.Tool != "SubagentHandback" || e.ToolID == "" || e.Input.Message == nil {
				return s
			}
			s.Text = *e.Input.Message
			state = 1
		case "SubagentStop":
			if state != 1 || e.Active || len(e.Background) > 0 || len(e.Crons) > 0 {
				return s
			}
			state = 2
		default:
			return s
		}
	}
	s.Raw = string(raw)
	s.Incomplete = scanner.Err() != nil || state != 2
	return s
}
