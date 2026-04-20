package mailer

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestTemplateRendererRender(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "welcome.html"), []byte(`<p>Hello {{.Name}}</p>`), 0644); err != nil {
		t.Fatalf("WriteFile(): %v", err)
	}

	renderer := NewTemplateRenderer(dir)
	got, err := renderer.Render("welcome", map[string]string{"Name": "Alice"})
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	if !strings.Contains(got, "Hello Alice") {
		t.Fatalf("expected rendered output to contain template data, got %q", got)
	}
}

func TestTemplateRendererRenderErrors(t *testing.T) {
	t.Parallel()

	renderer := NewTemplateRenderer(t.TempDir())
	if _, err := renderer.Render("missing", nil); err == nil {
		t.Fatal("expected parse error for missing template")
	}
}
