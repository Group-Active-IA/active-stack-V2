package headless_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/Group-Active-IA/active-stack/cmd/active-stack/headless"
)

// Review 2026-10-29 (7.1): toggles lived only in harnesses.yaml, so dropping
// ship or KB, or going back to the legacy composition, meant editing the
// catalog and rebuilding. --toggles picks them at install time.
func TestParseInstallFlags_Toggles(t *testing.T) {
	got, err := headless.ParseInstallFlags([]string{"--toggles", "lean-workflow, governance,evidence"})
	if err != nil {
		t.Fatal(err)
	}
	if got.TUI {
		t.Fatal("--toggles alone must select headless, like --tier")
	}
	want := []string{"lean-workflow", "governance", "evidence"}
	if !reflect.DeepEqual(got.Intent.Toggles, want) {
		t.Fatalf("Intent.Toggles = %v, want %v", got.Intent.Toggles, want)
	}
}

func TestParseInstallFlags_UnknownToggleListsKnownOnes(t *testing.T) {
	_, err := headless.ParseInstallFlags([]string{"--mode", "full", "--toggles", "lean-workflow,typo"})
	if err == nil {
		t.Fatal("unknown toggle accepted")
	}
	if !strings.Contains(err.Error(), "typo") || !strings.Contains(err.Error(), "lean-workflow") {
		t.Fatalf("error must name the bad toggle and the valid ones: %v", err)
	}
}

func TestParseInstallFlags_NoTogglesKeepsCatalogDefault(t *testing.T) {
	got, err := headless.ParseInstallFlags([]string{"--mode", "full"})
	if err != nil {
		t.Fatal(err)
	}
	if got.Intent.Toggles != nil {
		t.Fatalf("without --toggles the catalog decides, got %v", got.Intent.Toggles)
	}
}
