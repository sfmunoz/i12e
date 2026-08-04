package server

import (
	"bytes"
	"crypto/sha256"
	_ "embed"
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/sfmunoz/i12e/internal/tplutil"
)

//go:embed templates/rclone-conf.yaml
var rcloneConfYaml string

const manifestsFolder = "/var/lib/rancher/k3s/server/manifests"
const rcloneConfYamlPath = manifestsFolder + "/rclone-conf.yaml"

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

func rcloneConfYamlRender() (*bytes.Buffer, error) {
	tpl, err := template.New("templates/rclone-conf.yaml").Funcs(tplutil.FuncMap()).Option("missingkey=error").Parse(rcloneConfYaml)
	if err != nil {
		return nil, err
	}
	rcloneConf, err := os.ReadFile(rcloneConfPath)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", rcloneConfPath, err)
	}
	data := struct {
		Namespace  string
		RcloneConf string
	}{
		Namespace:  "i12e",
		RcloneConf: base64.StdEncoding.EncodeToString([]byte(removeBlankLines(string(rcloneConf)))),
	}
	var body bytes.Buffer
	if err = tpl.Execute(&body, data); err != nil {
		return nil, err
	}
	return &body, nil
}

func rcloneConfYamlWrite(bufNew *bytes.Buffer) error {
	sumNew := sha256.Sum256(bufNew.Bytes())
	bufOld, err := os.ReadFile(rcloneConfYamlPath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if err == nil && sha256.Sum256(bufOld) == sumNew {
		return nil
	}
	if err != nil {
		log.Info("creating manifest", "destPath", rcloneConfYamlPath)
	} else {
		log.Info("updating manifest", "destPath", rcloneConfYamlPath)
	}
	if err := os.WriteFile(rcloneConfYamlPath, bufNew.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Chown(rcloneConfYamlPath, 0, 0); err != nil {
		return fmt.Errorf("chown %s: %w", rcloneConfYamlPath, err)
	}
	return nil
}

func rcloneConf() error {
	if _, err := os.Stat(manifestsFolder); os.IsNotExist(err) {
		return nil
	}
	bufNew, err := rcloneConfYamlRender()
	if err != nil {
		return err
	}
	return rcloneConfYamlWrite(bufNew)
}
