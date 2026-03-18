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

func buildValuesContent() (string, error) {
	rcloneConf, err := os.ReadFile(rcloneConfPath)
	if err != nil {
		return "", fmt.Errorf("read %s: %w", rcloneConfPath, err)
	}
	lines := strings.Split(strings.TrimRight(string(rcloneConf), "\n"), "\n")
	var indented strings.Builder
	for _, line := range lines {
		indented.WriteString("    ")
		indented.WriteString(line)
		indented.WriteByte('\n')
	}
	return "rclone:\n  conf: |\n" + indented.String() + "  postgres:\n    password: changeme_now", nil
}

func postgresRclone() error {
	tpl, err := template.New("templates/postgres-rclone.yaml").Funcs(tplutil.FuncMap()).Option("missingkey=error").Parse(postgresRcloneYaml)
	if err != nil {
		return err
	}
	valuesContent, err := buildValuesContent()
	if err != nil {
		return err
	}
	data := struct {
		Version       string
		Namespace     string
		ValuesContent *bytes.Buffer
	}{
		Version:       "0.0.4",
		Namespace:     "i12e",
		ValuesContent: bytes.NewBufferString(valuesContent),
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
