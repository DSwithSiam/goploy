package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/dswithsiam/goploy/pkg/config"
	"github.com/dswithsiam/goploy/pkg/ssh"
	"github.com/spf13/cobra"
)

var (
	domain      string
	serverIP    string
	projectPath string
	framework   string
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy an app to your VPS",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Deploy called with:")
		fmt.Println("  Domain:", domain)
		fmt.Println("  Server IP:", serverIP)
		fmt.Println("  Project Path/Repo:", projectPath)
		fmt.Println("  Framework:", framework)

		// Prompt for SSH user and password
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter SSH username: ")
		user, _ := reader.ReadString('\n')
		user = strings.TrimSpace(user)
		fmt.Print("Enter SSH password: ")
		pw, _ := reader.ReadString('\n')
		pw = strings.TrimSpace(pw)

		addr := serverIP + ":22"
		client, err := ssh.NewSSHClient(user, addr, pw)
		if err != nil {
			fmt.Println("SSH connection failed:", err)
			return
		}
		defer client.Close()

		// Install nginx, redis, python3-pip
		fmt.Println("Installing nginx, redis, python3-pip...")
		client.Run("sudo apt-get update -y")
		client.Run("sudo apt-get install -y nginx redis-server python3-pip")

		// Install python dependencies
		if framework == "django" {
			fmt.Println("Installing Gunicorn for Django...")
			client.Run("sudo pip3 install gunicorn")
		} else if framework == "fastapi" {
			fmt.Println("Installing Uvicorn for FastAPI...")
			client.Run("sudo pip3 install uvicorn")
		}

		// Ask if Celery worker is needed before .env and systemd steps
		celeryNeeded := false
		fmt.Print("Deploy Celery worker? (y/N): ")
		celeryAns, _ := reader.ReadString('\n')
		celeryAns = strings.ToLower(strings.TrimSpace(celeryAns))
		if celeryAns == "y" || celeryAns == "yes" {
			celeryNeeded = true
			fmt.Println("Installing Celery and Redis Python dependencies...")
			client.Run("sudo pip3 install celery redis")
		}

		// Prompt for .env variables
		envVars := map[string]string{}
		fmt.Println("Enter environment variables (key=value), blank line to finish:")
		for {
			fmt.Print("> ")
			line, _ := reader.ReadString('\n')
			line = strings.TrimSpace(line)
			if line == "" {
				break
			}
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				envVars[parts[0]] = parts[1]
			}
		}
		envContent, err := config.GenerateEnv(envVars)
		if err != nil {
			fmt.Println("Failed to generate .env:", err)
			return
		}
		envTmp, err := os.CreateTemp("", ".env-*")
		if err != nil {
			fmt.Println("Temp file error:", err)
			return
		}
		defer os.Remove(envTmp.Name())
		envTmp.WriteString(envContent)
		envTmp.Seek(0, 0)
		err = client.Upload(envTmp, "/tmp/goploy.env")
		if err != nil {
			fmt.Println("Upload .env failed:", err)
			return
		}
		client.Run("sudo mv /tmp/goploy.env " + projectPath + "/.env")

		// Generate nginx config
		nginxConf, err := config.GenerateNginx(config.NginxConfig{
			Domain:  domain,
			AppPort: 8000,
		})
		if err != nil {
			fmt.Println("Failed to generate nginx config:", err)
			return
		}
		// Upload nginx config
		fmt.Println("Uploading nginx config...")
		tmpFile, err := os.CreateTemp("", "nginx-*.conf")
		if err != nil {
			fmt.Println("Temp file error:", err)
			return
		}
		defer os.Remove(tmpFile.Name())
		tmpFile.WriteString(nginxConf)
		tmpFile.Seek(0, 0)
		err = client.Upload(tmpFile, "/tmp/goploy_nginx.conf")
		if err != nil {
			fmt.Println("Upload failed:", err)
			return
		}
		client.Run("sudo mv /tmp/goploy_nginx.conf /etc/nginx/sites-available/goploy")
		client.Run("sudo ln -sf /etc/nginx/sites-available/goploy /etc/nginx/sites-enabled/goploy")
		client.Run("sudo nginx -t && sudo systemctl restart nginx")

		// Generate and upload systemd service
		var svcContent string
		svcCfg := config.SystemdConfig{
			AppName: "goploy-app",
			User:    user,
			WorkDir: projectPath,
			EnvFile: projectPath + "/.env",
		}
		if framework == "django" {
			svcCfg.ExecStart = "gunicorn app.wsgi:application --bind 0.0.0.0:8000"
			svcContent, err = config.GenerateGunicorn(svcCfg)
		} else if framework == "fastapi" {
			svcCfg.ExecStart = "uvicorn app:app --host 0.0.0.0 --port 8000"
			svcContent, err = config.GenerateUvicorn(svcCfg)
		}
		if err != nil {
			fmt.Println("Failed to generate systemd service:", err)
			return
		}
		svcTmp, err := os.CreateTemp("", "goploy-app.service")
		if err != nil {
			fmt.Println("Temp file error:", err)
			return
		}
		defer os.Remove(svcTmp.Name())
		svcTmp.WriteString(svcContent)
		svcTmp.Seek(0, 0)
		err = client.Upload(svcTmp, "/tmp/goploy-app.service")
		if err != nil {
			fmt.Println("Upload service failed:", err)
			return
		}
		client.Run("sudo mv /tmp/goploy-app.service /etc/systemd/system/goploy-app.service")
		client.Run("sudo systemctl daemon-reload && sudo systemctl enable --now goploy-app.service")

		// Celery worker (optional)
		fmt.Print("Deploy Celery worker? (y/N): ")
		celeryAns, _ := reader.ReadString('\n')
		celeryAns = strings.ToLower(strings.TrimSpace(celeryAns))
		if celeryAns == "y" || celeryAns == "yes" {
			svcCfg.ExecStart = "celery -A app worker --loglevel=info"
			celeryContent, err := config.GenerateCelery(svcCfg)
			if err != nil {
				fmt.Println("Failed to generate celery service:", err)
			} else {
				celeryTmp, err := os.CreateTemp("", "goploy-celery.service")
				if err != nil {
					fmt.Println("Temp file error:", err)
				} else {
					defer os.Remove(celeryTmp.Name())
					celeryTmp.WriteString(celeryContent)
					celeryTmp.Seek(0, 0)
					err = client.Upload(celeryTmp, "/tmp/goploy-celery.service")
					if err != nil {
						fmt.Println("Upload celery service failed:", err)
					} else {
						client.Run("sudo mv /tmp/goploy-celery.service /etc/systemd/system/goploy-celery.service")
						client.Run("sudo systemctl daemon-reload && sudo systemctl enable --now goploy-celery.service")
					}
				}
			}
		}

		fmt.Println("Deployment steps complete.")
	},
}

func init() {
	deployCmd.Flags().StringVarP(&domain, "domain", "d", "", "Domain name for the app (optional)")
	deployCmd.Flags().StringVarP(&serverIP, "server", "s", "", "Server IP address")
	deployCmd.Flags().StringVarP(&projectPath, "project", "p", "", "Project path or git repo")
	deployCmd.Flags().StringVarP(&framework, "framework", "f", "", "Framework type (Django, FastAPI, Node, static)")
	// domain is now optional
	deployCmd.MarkFlagRequired("server")
	deployCmd.MarkFlagRequired("project")
	deployCmd.MarkFlagRequired("framework")
	rootCmd.AddCommand(deployCmd)
}
