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

const rcloneConfPath = "/root/.config/rclone/rclone.conf"

const manifestsFolder = "/var/lib/rancher/k3s/server/manifests"
const postgresRcloneYamlPath = manifestsFolder + "/postgres-rclone.yaml"

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

func postgresRcloneYamlRender() (*bytes.Buffer, error) {
	tpl, err := template.New("templates/postgres-rclone.yaml").Funcs(tplutil.FuncMap()).Option("missingkey=error").Parse(postgresRcloneYaml)
	if err != nil {
		return nil, err
	}
	rcloneConf, err := os.ReadFile(rcloneConfPath)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", rcloneConfPath, err)
	}
	data := struct {
		Version    string
		Namespace  string
		RcloneConf *bytes.Buffer
	}{
		Version:    "0.0.8",
		Namespace:  "i12e",
		RcloneConf: bytes.NewBufferString(removeBlankLines(string(rcloneConf))),
	}
	var body bytes.Buffer
	if err = tpl.Execute(&body, data); err != nil {
		return nil, err
	}
	return &body, nil
}

func postgresRcloneYamlWrite(bufNew *bytes.Buffer) error {
	sumNew := sha256.Sum256(bufNew.Bytes())
	bufOld, err := os.ReadFile(postgresRcloneYamlPath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if err == nil && sha256.Sum256(bufOld) == sumNew {
		return nil
	}
	if err != nil {
		log.Info("creating manifest", "destPath", postgresRcloneYamlPath)
	} else {
		log.Info("updating manifest", "destPath", postgresRcloneYamlPath)
	}
	if err := os.WriteFile(postgresRcloneYamlPath, bufNew.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Chown(postgresRcloneYamlPath, 0, 0); err != nil {
		return fmt.Errorf("chown %s: %w", postgresRcloneYamlPath, err)
	}
	return nil
}

func postgresRclone() error {
	if _, err := os.Stat(manifestsFolder); os.IsNotExist(err) {
		return nil
	}
	bufNew, err := postgresRcloneYamlRender()
	if err != nil {
		return err
	}
	return postgresRcloneYamlWrite(bufNew)
}
