package generator

import (
	"bytes"
	"fmt"
	"text/template"
)

// processTemplate processes the raw template data using the Go text/template engine.
func processTemplate(name string, content []byte, config ProjectConfig) ([]byte, error) {
	tmpl, err := template.New(name).Parse(string(content))
	if err != nil {
		return nil, fmt.Errorf("error parsing template %s: %w", name, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, config); err != nil {
		return nil, fmt.Errorf("error executing template %s: %w", name, err)
	}

	return buf.Bytes(), nil
}
