package config

import (
	"bytes"
	"text/template"
)

const envTemplate = `{{range $k, $v := .}}{{$k}}={{$v}}
{{end}}`

func GenerateEnv(vars map[string]string) (string, error) {
	tmpl, err := template.New("env").Parse(envTemplate)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, vars)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
