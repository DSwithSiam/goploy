package config

import (
	"bytes"
	"text/template"
)

const nginxTemplate = `server {
    listen 80;
    server_name {{.Domain}};
    location / {
        proxy_pass http://127.0.0.1:{{.AppPort}};
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}`

type NginxConfig struct {
	Domain  string
	AppPort int
}

func GenerateNginx(cfg NginxConfig) (string, error) {
	tmpl, err := template.New("nginx").Parse(nginxTemplate)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, cfg)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
