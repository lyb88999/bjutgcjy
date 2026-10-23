package core

import (
	"os"
	"path/filepath"
	"testing"
)

func clearAgentEnvironment(t *testing.T) {
	t.Helper()
	for _, name := range []string{"EXPERT_AGENT_API_KEY", "EXPERT_AGENT_MODEL", "EXPERT_AGENT_BASE_URL"} {
		name := name
		previous, exists := os.LookupEnv(name)
		if err := os.Unsetenv(name); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() {
			if exists {
				_ = os.Setenv(name, previous)
			} else {
				_ = os.Unsetenv(name)
			}
		})
	}
}

func TestLoadExpertAgentEnv_PreservesProcessSettingsAndLoadsOnlyAllowlist(t *testing.T) {
	clearAgentEnvironment(t)
	t.Setenv("EXPERT_AGENT_API_KEY", "process-test-key")
	t.Setenv("EXPERT_AGENT_UNRELATED", "unchanged")
	dir := t.TempDir()
	text := "# local settings\nEXPERT_AGENT_API_KEY=file-test-key\nexport EXPERT_AGENT_MODEL='glm-4.7-flash'\nEXPERT_AGENT_BASE_URL=\"https://example.invalid/api\"\nEXPERT_AGENT_UNRELATED=unexpected\n"
	if err := os.WriteFile(filepath.Join(dir, ".env.local"), []byte(text), 0600); err != nil {
		t.Fatal(err)
	}
	if err := loadExpertAgentEnv(filepath.Join(dir, "custom.yaml")); err != nil {
		t.Fatal(err)
	}
	if os.Getenv("EXPERT_AGENT_API_KEY") != "process-test-key" || os.Getenv("EXPERT_AGENT_UNRELATED") != "unchanged" {
		t.Fatal("local settings must not replace process values or load unrelated variables")
	}
	if os.Getenv("EXPERT_AGENT_MODEL") != "glm-4.7-flash" || os.Getenv("EXPERT_AGENT_BASE_URL") != "https://example.invalid/api" {
		t.Fatal("quoted assistant settings should be loaded beside the selected config")
	}
}

func TestLoadExpertAgentEnv_MalformedFileDoesNotPartiallyLoad(t *testing.T) {
	clearAgentEnvironment(t)
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, ".env.local"), []byte("EXPERT_AGENT_API_KEY=test-key\nEXPERT_AGENT_MODEL=\"unterminated\n"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := loadExpertAgentEnv(filepath.Join(dir, "config.yaml")); err == nil {
		t.Fatal("malformed quoted settings should fail")
	}
	if _, exists := os.LookupEnv("EXPERT_AGENT_API_KEY"); exists {
		t.Fatal("a malformed file must not load credentials partially")
	}
}

func TestLoadExpertAgentEnv_MissingFileIsOptional(t *testing.T) {
	if err := loadExpertAgentEnv(filepath.Join(t.TempDir(), "config.yaml")); err != nil {
		t.Fatal(err)
	}
}
