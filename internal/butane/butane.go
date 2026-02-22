package butane

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/sfmunoz/i12e/internal/cmdutil"
	"github.com/sfmunoz/i12e/internal/config"
	"github.com/sfmunoz/logit"
)

var log = logit.Logit().WithLevel(logit.LevelInfo)

type Butane struct {
	cfg *config.Config
}

func newButane(cfg *config.Config) (*Butane, error) {
	if cfg == nil {
		return nil, fmt.Errorf("newButane(): undefined config")
	}
	return &Butane{cfg: cfg}, nil
}

func (b *Butane) butaneCmd() *exec.Cmd {
	if b.cfg.Butane.Output == config.OutputDebug {
		return exec.Command("butane", "-s", "-p")
	}
	return exec.Command("butane", "-s")
}

func (b *Butane) ignitionConfigMergeSource() (*bytes.Buffer, error) {
	fname := b.cfg.Env.ButaneEncYaml()
	_, err := os.Stat(fname)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	bo1, be1, err := cmdutil.RunSimple(exec.Command("sops", "decrypt", fname))
	if err != nil {
		return nil, fmt.Errorf("'sops decrypt %s' failed: err=%s; buf_err=%s", fname, err, be1)
	}
	cmd := b.butaneCmd()
	cmd.Stdin = bo1
	bo2, be2, err := cmdutil.RunSimple(cmd)
	if err != nil {
		return nil, fmt.Errorf("'butane' failed: err=%s; buf_err=%s", err, be2)
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

func (b *Butane) butaneRender() (*bytes.Buffer, error) {
	tpl, err := tplNew("flatcar.yaml", true)
	if err != nil {
		return nil, err
	}
	icms, err := b.ignitionConfigMergeSource()
	if err != nil {
		return nil, err
	}
	rcloneConf, err := rcloneConfig(b.cfg.Rclone.Remote)
	if err != nil {
		return nil, err
	}
	data := struct {
		IgnitionConfigMergeSource *bytes.Buffer
		I12eVersion               string
		I12eSha256sum             string
		Mode                      string
		SshAuthorizedKeys         []string
		RcloneConf                *bytes.Buffer
	}{
		IgnitionConfigMergeSource: icms,
		I12eVersion:               b.cfg.I12e.Version,
		I12eSha256sum:             b.cfg.I12e.Sha256sum,
		Mode:                      b.cfg.Butane.Mode.String(),
		SshAuthorizedKeys:         b.cfg.Butane.SshAuthorizedKeys,
		RcloneConf:                rcloneConf,
	}
	var ret bytes.Buffer
	err = tpl.Execute(&ret, &data)
	if err != nil {
		return nil, err
	}
	if b.cfg.Butane.Output == config.OutputDebug {
		log.Info("======== butane begin ========")
		for _, line := range strings.Split(ret.String(), "\n") {
			log.Info(line)
		}
		log.Info("-------- butane end --------")
	}
	return &ret, nil
}

func (b *Butane) ignitionRender(buf *bytes.Buffer) (*bytes.Buffer, error) {
	cmd := b.butaneCmd()
	cmd.Stdin = buf
	bo, be, err := cmdutil.RunSimple(cmd)
	if err != nil {
		return nil, fmt.Errorf("'butane' failed: err=%s; buf_err=%s", err, be)
	}
	if b.cfg.Butane.Output == config.OutputDebug {
		log.Info("======== ignition begin ========")
		for _, line := range strings.Split(bo.String(), "\n") {
			log.Info(line)
		}
		log.Info("-------- ignition end --------")
	}
	return bo, nil
}

func (b *Butane) bashRaw(ibuf *bytes.Buffer) (*bytes.Buffer, error) {
	gzBuf, err := b64Gzip(ibuf)
	if err != nil {
		return nil, err
	}
	tpl, err := tplNew("bash_raw.sh", false)
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
	return &ret, nil
}

func (b *Butane) bashB64(ibuf *bytes.Buffer) (*bytes.Buffer, error) {
	gzBuf, err := b64Gzip(ibuf)
	if err != nil {
		return nil, err
	}
	tpl, err := tplNew("bash_b64.sh", false)
	if err != nil {
		return nil, err
	}
	var ret bytes.Buffer
	data := struct{ Buf *bytes.Buffer }{Buf: gzBuf}
	err = tpl.Execute(&ret, &data)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func (b *Butane) run() error {
	buf, err := b.butaneRender()
	if err != nil {
		return err
	}
	ignitionBuf, err := b.ignitionRender(buf)
	if err != nil {
		return err
	}
	if b.cfg.Butane.Output == config.OutputDebug {
		return nil
	}
	if b.cfg.Butane.Output == config.OutputIgnition {
		fmt.Fprint(os.Stdout, ignitionBuf.String())
		return nil
	}
	bufRaw, err := b.bashRaw(ignitionBuf)
	if err != nil {
		return err
	}
	if b.cfg.Butane.Output == config.OutputBashRaw {
		fmt.Fprint(os.Stdout, bufRaw.String())
		return nil
	}
	bufB64, err := b.bashB64(bufRaw)
	if err != nil {
		return err
	}
	if b.cfg.Butane.Output == config.OutputBashB64 {
		fmt.Fprint(os.Stdout, bufB64.String())
		return nil
	}
	return nil
}

func Run(cfg *config.Config) error {
	b, err := newButane(cfg)
	if err != nil {
		return err
	}
	return b.run()
}
