package verify

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/Group-Active-IA/active-stack/internal/harness/config"
	"github.com/Group-Active-IA/active-stack/internal/model"
)

func composedDeliveryCheck(h model.Harness, adapter Adapter, home string) []Check {
	variant, ok := adapter.(interface{ VariantKey() string })
	if !ok || h.ID != "sdd-orchestrator" {
		return nil
	}
	return []Check{{ID: "config:compose:" + string(adapter.Agent()), Description: "contenido instalado igual a Compose para variante y toggles efectivos", Run: func(context.Context) error {
		expected, err := config.Compose(h.Toggles, variant.VariantKey())
		if err != nil {
			return err
		}
		path := adapter.InstructionsPath(home)
		if adapter.ConfigDelivery() == model.ConfigDeliveryPrimaryAgent {
			path = adapter.SettingsPath(home)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		var actual string
		if adapter.ConfigDelivery() == model.ConfigDeliveryPrimaryAgent {
			var root struct {
				Agent map[string]struct {
					Prompt string `json:"prompt"`
				} `json:"agent"`
			}
			if err := json.Unmarshal(data, &root); err != nil {
				return err
			}
			actual = root.Agent[h.ID].Prompt
		} else {
			open := "<!-- active-stack:sdd-orchestrator -->"
			close := "<!-- /active-stack:sdd-orchestrator -->"
			text := string(data)
			if strings.Count(text, open) != 1 || strings.Count(text, close) != 1 {
				return fmt.Errorf("markers del orquestador incompletos o duplicados en %q", path)
			}
			_, tail, _ := strings.Cut(text, open)
			actual, _, ok = strings.Cut(tail, close)
			if !ok {
				return fmt.Errorf("marker de cierre ausente en %q", path)
			}
		}
		if strings.TrimSpace(actual) != strings.TrimSpace(expected) {
			return fmt.Errorf("contenido instalado en %q difiere de Compose(%v, %q)", path, h.Toggles, variant.VariantKey())
		}
		return nil
	}}}
}
