package renderer

import (
	"bytes"
	"html/template"

	"github.com/martinrue/ofteco/assets"
)

func renderTemplate(path string, data interface{}) (string, error) {
	tmpl, err := assets.FSString(false, path)
	if err != nil {
		return "", err
	}

	funcs := template.FuncMap{
		"pluralise": func(word string, amount int) string {
			if amount != 1 {
				return word + "j"
			}

			return word
		},
	}

	t, err := template.New("index").Funcs(funcs).Parse(tmpl)
	if err != nil {
		return "", err
	}

	buffer := &bytes.Buffer{}
	t.Execute(buffer, data)

	return buffer.String(), nil
}
