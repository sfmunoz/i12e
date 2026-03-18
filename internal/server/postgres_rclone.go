package server

import (
	"bytes"
	"crypto/sha256"
	_ "embed"
	"fmt"
	"os"
	"text/template"

	"github.com/sfmunoz/i12e/internal/tplutil"
)

//go:embed templates/postgres-rclone.yaml
var postgresRcloneYaml string

const valuesContent = `rclone:
  conf: |
    [rem]
    remote = rustfs:d01
    type = alias
    [rustfs]
    acl = private
    endpoint = http://192.168.56.1:9000
    provider = Other
    secret_access_key = rustfsadmin
    type = s3
    access_key_id = rustfsadmin
  postgres:
    password: changeme_now`

func postgresRclone() error {
	tpl, err := template.New("templates/postgres-rclone.yaml").Funcs(tplutil.FuncMap()).Option("missingkey=error").Parse(postgresRcloneYaml)
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
	const destPath = "/var/lib/rancher/k3s/server/manifests/postgres-rclone.yaml"
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
