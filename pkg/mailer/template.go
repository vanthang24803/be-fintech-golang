package mailer

import (
	"bytes"
	"fmt"
	"html/template"
	"path/filepath"
)

// TemplateRenderer handles rendering of HTML templates
type TemplateRenderer struct {
	basePath string
}

// NewTemplateRenderer initializes a new renderer
func NewTemplateRenderer(basePath string) *TemplateRenderer {
	return &TemplateRenderer{basePath: basePath}
}

// Render renders a template with the provided data
func (r *TemplateRenderer) Render(templateName string, data interface{}) (string, error) {
	tmplPath := filepath.Join(r.basePath, templateName+".html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}
