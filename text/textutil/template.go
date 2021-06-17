package textutil

import (
	"bytes"
	"net/url"
	"os"
	"text/template"
)

// process applies the data structure 'vars' onto an already
// parsed template 't', and returns the resulting string.
func render(t *template.Template, vars interface{}) string {
	var tmplBytes bytes.Buffer

	err := t.Execute(&tmplBytes, vars)
	if err != nil {
		panic(err)
	}

	return tmplBytes.String()
}

// RenderString should be commented
func RenderString(str string, vars interface{}) string {
	tmpl, err := template.New("tmpl").Funcs(FuncMap()).Parse(str)

	if err != nil {
		panic(err)
	}

	return render(tmpl, vars)
}

// FuncMap should be commented
func FuncMap() template.FuncMap {
	return template.FuncMap{
		"pathescape": url.PathEscape,

		"pwd": func() string {
			pwd, err := os.Getwd()

			if err != nil {
				panic(err)
			}

			return pwd
		},
	}
}
