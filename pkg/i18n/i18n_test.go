package i18n

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadLocalesAndTranslate(t *testing.T) {
	t.Parallel()

	prev := translations
	translations = make(map[string]map[string]interface{})
	t.Cleanup(func() {
		translations = prev
	})

	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "en.json"), []byte(`{"greetings":{"hello":"Hello {{name}}"}}`), 0644); err != nil {
		t.Fatalf("WriteFile(en.json): %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "vi.json"), []byte(`{"greetings":{"hello":"Xin chao {{name}}"}}`), 0644); err != nil {
		t.Fatalf("WriteFile(vi.json): %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "README.txt"), []byte("ignore"), 0644); err != nil {
		t.Fatalf("WriteFile(README.txt): %v", err)
	}

	if err := LoadLocales(dir); err != nil {
		t.Fatalf("LoadLocales() error = %v", err)
	}

	if got := Translate("vi", "greetings.hello", map[string]interface{}{"name": "An"}); got != "Xin chao An" {
		t.Fatalf("expected translated greeting, got %q", got)
	}
	if got := Translate("fr", "greetings.hello", map[string]interface{}{"name": "Lan"}); got != "Hello Lan" {
		t.Fatalf("expected english fallback, got %q", got)
	}
}

func TestTranslateFallbacks(t *testing.T) {
	t.Parallel()

	prev := translations
	translations = map[string]map[string]interface{}{
		"en": {
			"simple": "Hello",
			"nested": map[string]interface{}{
				"value": "World",
			},
		},
		"vi": {},
	}
	t.Cleanup(func() {
		translations = prev
	})

	if got := Translate("vi", "nested.value"); got != "World" {
		t.Fatalf("expected english nested fallback, got %q", got)
	}
	if got := Translate("en", "missing.key"); got != "missing.key" {
		t.Fatalf("expected missing key fallback, got %q", got)
	}
	if got := Translate("en", "nested"); got != "nested" {
		t.Fatalf("expected non-string fallback to key, got %q", got)
	}
}
