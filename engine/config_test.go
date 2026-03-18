package engine

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Zyko0/go-sdl3/sdl"
)

func TestParseColor(t *testing.T) {
	tests := []struct {
		input    string
		expected sdl.Color
	}{
		{"255,128,64,255", sdl.Color{R: 255, G: 128, B: 64, A: 255}},
		{"0,0,0,0", sdl.Color{R: 0, G: 0, B: 0, A: 0}},
		{"100,100,100", sdl.Color{R: 100, G: 100, B: 100, A: 255}},
		{"", sdl.Color{R: 0, G: 0, B: 0, A: 0}},
	}

	for _, tt := range tests {
		got := ParseColor(tt.input)
		if got != tt.expected {
			t.Errorf("ParseColor(%q) = %v, expected %v", tt.input, got, tt.expected)
		}
	}
}

func TestConfig_SaveLoad(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "test_config.toml")

	content := `
subject_id = "test-subject"
font_size = 42
bg_color = "10,20,30,255"
last_dir = "/tmp/test"
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	cfg := DefaultConfig()
	if err := cfg.LoadFromFile(configPath); err != nil {
		t.Fatalf("LoadFromFile failed: %v", err)
	}

	if cfg.SubjectID != "test-subject" {
		t.Errorf("expected SubjectID 'test-subject', got %q", cfg.SubjectID)
	}
	if cfg.FontSize != 42 {
		t.Errorf("expected FontSize 42, got %d", cfg.FontSize)
	}
	expectedColor := sdl.Color{R: 10, G: 20, B: 30, A: 255}
	if cfg.BGColor != expectedColor {
		t.Errorf("expected BGColor %v, got %v", expectedColor, cfg.BGColor)
	}
	if cfg.LastDir != "/tmp/test" {
		t.Errorf("expected LastDir '/tmp/test', got %q", cfg.LastDir)
	}
}
