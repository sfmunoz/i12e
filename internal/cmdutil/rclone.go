package cmdutil

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/sfmunoz/i12e/internal/config"
	"go.yaml.in/yaml/v3"
)

func RcloneConfig(env config.Env) (*bytes.Buffer, error) {
	if env != config.EnvDev && env != config.EnvProd {
		return nil, fmt.Errorf("unsupported '%s' environment (valid: 'dev' or 'prod')", env.String())
	}
	secretsPath := strings.TrimSpace(os.Getenv("I12E_SECRETS"))
	if len(secretsPath) < 1 {
		secretsPath = path.Join("..", "i12e-secrets")
	}
	rcloneConfYaml := path.Join(secretsPath, "clusters", env.String(), "kube-system", "rclone-conf.yaml")
	bo, be, err := RunSimple(exec.Command("sops", "decrypt", rcloneConfYaml))
	if err != nil {
		return nil, fmt.Errorf("'sops decrypt %s' failed: err=%s; buf_err=%s", rcloneConfYaml, err, be)
	}
	var secret struct {
		StringData struct {
			ConfigData string `yaml:"configData"`
		} `yaml:"stringData"`
	}
	if err := yaml.Unmarshal(bo.Bytes(), &secret); err != nil {
		return nil, err
	}
	ret := secret.StringData.ConfigData
	if len(ret) < 1 {
		return nil, fmt.Errorf("cannot get 'stringData.configData' from secret")
	}
	return bytes.NewBufferString(ret), nil
}
