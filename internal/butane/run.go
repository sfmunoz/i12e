package butane

import (
	"bytes"
	"compress/gzip"
	"embed"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/sfmunoz/i12e/internal/cmdutil"
	"github.com/sfmunoz/i12e/internal/config"
	"github.com/sfmunoz/logit"
)

var log = logit.Logit().WithLevel(logit.LevelInfo)

type FlatcarYaml struct {
	IgnitionConfigMergeSource *bytes.Buffer
	I12eVersion               string
	I12eSha256sum             string
	Mode                      string
	SshAuthorizedKeys         []string
	RcloneConf                *bytes.Buffer
}

var i12eVersion = "v0.0.18"
var i12eSha256Sum = "cfe8d33bc00805344dbe4008d87b896ea0c3bb0618cc69bcf5bc0462af4a2709"

//go:embed templates/*.yaml templates/*.sh
var FS embed.FS

func b64Gzip(ibuf *bytes.Buffer) (*bytes.Buffer, error) {
	var ret bytes.Buffer
	b64Enc := base64.NewEncoder(base64.StdEncoding, &ret)
	gzEnc := gzip.NewWriter(b64Enc)
	_, err := gzEnc.Write(ibuf.Bytes())
	if err != nil {
		return nil, err
	}
	if err := gzEnc.Close(); err != nil {
		return nil, err
	}
	if err := b64Enc.Close(); err != nil {
		return nil, err
	}
	return &ret, nil
}

func bashRaw(cfg *config.Config, ibuf *bytes.Buffer) (*bytes.Buffer, error) {
	gzBuf, err := b64Gzip(ibuf)
	if err != nil {
		return nil, err
	}
	tpl := template.New("bash_raw.sh")
	tpl, err = tpl.Option("missingkey=error").ParseFS(FS, "templates/bash_raw.sh")
	if err != nil {
		return nil, err
	}
	var ret bytes.Buffer
	data := struct {
		ConfigIgn string
		Buf       *bytes.Buffer
	}{
		ConfigIgn: "/oem/config.ign",
		Buf:       gzBuf,
	}
	err = tpl.Execute(&ret, &data)
	if err != nil {
		return nil, err
	}
	log.Info("======== bashRaw begin ========")
	for _, line := range strings.Split(ret.String(), "\n") {
		log.Info(line)
	}
	log.Info("-------- bashRaw end --------")
	return &ret, nil
}

func bashB64(cfg *config.Config, ibuf *bytes.Buffer) (*bytes.Buffer, error) {
	gzBuf, err := b64Gzip(ibuf)
	if err != nil {
		return nil, err
	}
	tpl := template.New("bash_b64.sh")
	tpl, err = tpl.Option("missingkey=error").ParseFS(FS, "templates/bash_b64.sh")
	if err != nil {
		return nil, err
	}
	var ret bytes.Buffer
	data := struct{ Buf *bytes.Buffer }{Buf: gzBuf}
	err = tpl.Execute(&ret, &data)
	if err != nil {
		return nil, err
	}
	log.Info("======== bashB64 begin ========")
	for _, line := range strings.Split(ret.String(), "\n") {
		log.Info(line)
	}
	log.Info("-------- bashB64 end --------")
	return &ret, nil
}

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
	cmd := exec.Command("butane", "-s")
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

func indent(spaces int, v fmt.Stringer) string {
	pad := strings.Repeat(" ", spaces)
	return pad + strings.ReplaceAll(v.String(), "\n", "\n"+pad)
}

func nindent(spaces int, v fmt.Stringer) string {
	return "\n" + indent(spaces, v)
}

func butaneRender(cfg *config.Config) (*bytes.Buffer, error) {
	var funcMap = template.FuncMap{
		"indent":  indent,
		"nindent": nindent,
	}
	tpl := template.New("flatcar.yaml").Funcs(funcMap)
	tpl, err := tpl.Option("missingkey=error").ParseFS(FS, "templates/flatcar.yaml")
	if err != nil {
		return nil, err
	}
	icms, err := ignitionConfigMergeSource(cfg)
	if err != nil {
		return nil, err
	}
	rcloneConfig, err := RcloneConfig(cfg)
	if err != nil {
		return nil, err
	}
	f := FlatcarYaml{
		IgnitionConfigMergeSource: icms,
		I12eVersion:               i12eVersion,
		I12eSha256sum:             i12eSha256Sum,
		Mode:                      cfg.Mode.String(),
		SshAuthorizedKeys:         cfg.SshAuthorizedKeys,
		RcloneConf:                rcloneConfig,
	}
	var ret bytes.Buffer
	err = tpl.Execute(&ret, &f)
	if err != nil {
		return nil, err
	}
	log.Info("======== butane begin ========")
	for _, line := range strings.Split(ret.String(), "\n") {
		log.Info(line)
	}
	log.Info("-------- butane end --------")
	return &ret, nil
}

func ignitionRender(cfg *config.Config, buf *bytes.Buffer, pretty bool) (*bytes.Buffer, error) {
	cmd := func() *exec.Cmd {
		if pretty {
			return exec.Command("butane", "-s", "-p")
		}
		return exec.Command("butane", "-s")
	}()
	cmd.Stdin = buf
	bo, be, err := cmdutil.RunSimple(cmd)
	if err != nil {
		return nil, fmt.Errorf("'butane' failed: err=%s; buf_err=%s; prod=%t", err, be, cfg.Prod)
	}
	log.Info("======== ignition begin ========")
	for _, line := range strings.Split(bo.String(), "\n") {
		log.Info(line)
	}
	log.Info("-------- ignition end --------")
	return bo, nil
}

func Run(cfg *config.Config) error {
	if cfg == nil {
		return fmt.Errorf("butane.Run(): undefined config")
	}
	if cfg.Bout != config.BoutDebug {
		return fmt.Errorf("butane.Run(): output mode '%s' not implemented yet: use '%s'", cfg.Bout, config.BoutDebug)
	}
	buf, err := butaneRender(cfg)
	if err != nil {
		return err
	}
	ignitionBuf, err := ignitionRender(cfg, buf, true)
	if err != nil {
		return err
	}
	log.Info("ignitionBuf", "ignitionBuf", ignitionBuf)
	bufRaw, err := bashRaw(cfg, ignitionBuf)
	if err != nil {
		return err
	}
	_, err = bashB64(cfg, bufRaw)
	if err != nil {
		return err
	}
	return nil
}
