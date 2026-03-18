package server

import (
	"bytes"
	_ "embed"
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
	println(body.String())
	return nil
}
