package config

import (
	"bytes"
	"text/template"
)

const gunicornServiceTmpl = `[Unit]
Description=Gunicorn for {{.AppName}}
After=network.target

[Service]
User={{.User}}
Group=www-data
WorkingDirectory={{.WorkDir}}
EnvironmentFile={{.EnvFile}}
ExecStart={{.ExecStart}}
Restart=always

[Install]
WantedBy=multi-user.target
`

const uvicornServiceTmpl = `[Unit]
Description=Uvicorn for {{.AppName}}
After=network.target

[Service]
User={{.User}}
Group=www-data
WorkingDirectory={{.WorkDir}}
EnvironmentFile={{.EnvFile}}
ExecStart={{.ExecStart}}
Restart=always

[Install]
WantedBy=multi-user.target
`

const celeryServiceTmpl = `[Unit]
Description=Celery Worker for {{.AppName}}
After=network.target

[Service]
User={{.User}}
Group=www-data
WorkingDirectory={{.WorkDir}}
EnvironmentFile={{.EnvFile}}
ExecStart={{.ExecStart}}
Restart=always

[Install]
WantedBy=multi-user.target
`

type SystemdConfig struct {
	AppName   string
	User      string
	WorkDir   string
	EnvFile   string
	ExecStart string
}

func GenerateGunicorn(cfg SystemdConfig) (string, error) {
	return renderTmpl(gunicornServiceTmpl, cfg)
}

func GenerateUvicorn(cfg SystemdConfig) (string, error) {
	return renderTmpl(uvicornServiceTmpl, cfg)
}

func GenerateCelery(cfg SystemdConfig) (string, error) {
	return renderTmpl(celeryServiceTmpl, cfg)
}

func renderTmpl(tmplStr string, cfg SystemdConfig) (string, error) {
	tmpl, err := template.New("").Parse(tmplStr)
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
