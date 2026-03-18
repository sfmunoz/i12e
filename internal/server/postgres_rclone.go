package server

import (
	"bytes"
	"crypto/sha256"
	_ "embed"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/sfmunoz/i12e/internal/tplutil"
)

//go:embed templates/postgres-rclone.yaml
var postgresRcloneYaml string

const destPath = "/var/lib/rancher/k3s/server/manifests/postgres-rclone.yaml"

const rcloneConfPath = "/root/.config/rclone/rclone.conf"

func removeBlankLines(s string) string {
	lines := strings.Split(strings.TrimRight(s, "\n"), "\n")
	filtered := lines[:0]
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			filtered = append(filtered, line)
		}
	}
	return strings.Join(filtered, "\n")
}

func postgresRclone() error {
	tpl, err := template.New("templates/postgres-rclone.yaml").Funcs(tplutil.FuncMap()).Option("missingkey=error").Parse(postgresRcloneYaml)
	if err != nil {
		return err
	}
	rcloneConf, err := os.ReadFile(rcloneConfPath)
	if err != nil {
		return fmt.Errorf("read %s: %w", rcloneConfPath, err)
	}
	data := struct {
		Version          string
		Namespace        string
		RcloneConf       *bytes.Buffer
		PostgresPassword string
	}{
		Version:          "0.0.4",
		Namespace:        "i12e",
		RcloneConf:       bytes.NewBufferString(removeBlankLines(string(rcloneConf))),
		PostgresPassword: "changeme_now",
	}
	var body bytes.Buffer
	err = tpl.Execute(&body, data)
	if err != nil {
		return err
	}
	bufNew := body.Bytes()
	sumNew := sha256.Sum256(bufNew)
	bufOld, err := os.ReadFile(destPath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if err == nil && sha256.Sum256(bufOld) == sumNew {
		return nil
	}
	log.Info("updating manifest", "destPath", destPath)
	if err := os.WriteFile(destPath, bufNew, 0600); err != nil {
		return err
	}
	if err := os.Chown(destPath, 0, 0); err != nil {
		return fmt.Errorf("chown %s: %w", destPath, err)
	}
	return nil
}
