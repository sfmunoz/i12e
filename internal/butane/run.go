package butane

import (
	"bytes"
	"embed"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"text/template"

	"github.com/sfmunoz/i12e/internal/cmdutil"
	"github.com/sfmunoz/i12e/internal/config"
	"github.com/sfmunoz/logit"
)

var log = logit.Logit().
	WithLevel(logit.LevelInfo).
	With("mod", "i12e").
	With("pkg", "butane")

type FlatcarYaml struct {
	IgnitionConfigMergeSource *bytes.Buffer
	I12eVersion               string
	I12eSha256sum             string
	Mode                      string
	SshAuthorizedKeys         []string
	RcloneConf                string
}

var i12eVersion = "v0.0.18"
var i12eSha256Sum = "cfe8d33bc00805344dbe4008d87b896ea0c3bb0618cc69bcf5bc0462af4a2709"

//go:embed templates/*.yaml
var FS embed.FS

func ignitionConfigMergeSource(cfg *config.Config) (*bytes.Buffer, error) {
	fname := "butane-dev.yaml"
	if cfg.Prod {
		fname = "butane-prod.yaml"
	}
	_, err := os.Stat(fname)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	bo1, be1, err := cmdutil.RunSimple(exec.Command("sops", "decrypt", fname))
	if err != nil {
		return nil, fmt.Errorf("'sops decrypt' failed: err=%s; buf_err=%s; prod=%t", err, be1, cfg.Prod)
	}
	cmd := exec.Command("butane")
	cmd.Stdin = bo1
	bo2, be2, err := cmdutil.RunSimple(cmd)
	if err != nil {
		return nil, fmt.Errorf("'butane' failed: err=%s; buf_err=%s; prod=%t", err, be2, cfg.Prod)
	}
	var ret bytes.Buffer
	_, err = ret.WriteString("data:;base64,")
	if err != nil {
		return nil, err
	}
	enc := base64.NewEncoder(base64.StdEncoding, &ret)
	_, err = enc.Write(bo2.Bytes())
	if err != nil {
		return nil, err
	}
	err = enc.Close()
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func flatcarYamlRender(cfg *config.Config) error {
	tpl := template.New("flatcar.yaml") // must match basename of the file
	tpl, err := tpl.Option("missingkey=error").ParseFS(FS, "templates/flatcar.yaml")
	if err != nil {
		return err
	}
	i, err := ignitionConfigMergeSource(cfg)
	if err != nil {
		return err
	}
	log.Info("butane.Run()", "ignition", i)
	f := FlatcarYaml{
		IgnitionConfigMergeSource: i,
		I12eVersion:               i12eVersion,
		I12eSha256sum:             i12eSha256Sum,
		Mode:                      cfg.Mode.String(),
		SshAuthorizedKeys:         cfg.SshAuthorizedKeys,
		RcloneConf:                "**** RcloneConf ****",
	}
	err = tpl.Execute(os.Stdout, &f)
	if err != nil {
		return err
	}
	return nil
}

func Run(cfg *config.Config) error {
	if cfg == nil {
		return fmt.Errorf("butane.Run(): undefined config")
	}
	log.Info("butane.Run()", "cfg", cfg)
	return flatcarYamlRender(cfg)
}
