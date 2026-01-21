package artifact

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"os"
	"time"

	"github.com/sfmunoz/i12e/internal/config"
	"github.com/sfmunoz/logit"
)

var log = logit.Logit().WithLevel(logit.LevelInfo)

type Artifact struct {
	cfg    *config.Config
	tnow   time.Time
	tarOut *tar.Writer
	gzOut  *gzip.Writer
	obuf   *bytes.Buffer
}

func newArtifact(cfg *config.Config) *Artifact {
	var obuf bytes.Buffer
	gzOut := gzip.NewWriter(&obuf)
	tarOut := tar.NewWriter(gzOut)
	return &Artifact{
		cfg:    cfg,
		tnow:   time.Now().UTC(),
		tarOut: tarOut,
		gzOut:  gzOut,
		obuf:   &obuf,
	}
}

func (a *Artifact) Close() error {
	if err := a.tarOut.Close(); err != nil {
		return err
	}
	if err := a.gzOut.Close(); err != nil {
		return err
	}
	return nil
}

func (a *Artifact) Obuf() *bytes.Buffer {
	return a.obuf
}

func (a *Artifact) folders() error {
	flist := []struct {
		Name string
		Mode int64
	}{
		{Name: "etc/i12e", Mode: 0700},
		{Name: "etc/i12e/flags", Mode: 0700},
		{Name: "etc/i12e/k3s", Mode: 0700},
		{Name: "etc/systemd/system/k3s.service.d", Mode: 0755},
		{Name: "etc/systemd/system.conf.d", Mode: 0755},
		{Name: "opt/libexec", Mode: 0755},
		{Name: "opt/libexec/i12e", Mode: 0755},
	}
	for _, folder := range flist {
		hdr := &tar.Header{
			Typeflag: tar.TypeDir,
			Name:     folder.Name,
			ModTime:  a.tnow,
			Mode:     folder.Mode,
			Uid:      0,
			Gid:      0,
			Uname:    "root",
			Gname:    "root",
		}
		if err := a.tarOut.WriteHeader(hdr); err != nil {
			return err
		}
	}
	return nil
}

func (a *Artifact) addStatic(staticFname, targetFname string, mode int64) error {
	body, err := FS.ReadFile(staticFname)
	if err != nil {
		return err
	}
	hdr := &tar.Header{
		Typeflag: tar.TypeReg,
		Name:     targetFname,
		ModTime:  a.tnow,
		Size:     int64(len(body)),
		Mode:     mode,
		Uid:      0,
		Gid:      0,
		Uname:    "root",
		Gname:    "root",
	}
	if err := a.tarOut.WriteHeader(hdr); err != nil {
		return err
	}
	if _, err := a.tarOut.Write(body); err != nil {
		return err
	}
	return nil
}

func (a *Artifact) addEmpty(targetFname string, mode int64) error {
	hdr := &tar.Header{
		Typeflag: tar.TypeReg,
		Name:     targetFname,
		ModTime:  a.tnow,
		Size:     0,
		Mode:     mode,
		Uid:      0,
		Gid:      0,
		Uname:    "root",
		Gname:    "root",
	}
	if err := a.tarOut.WriteHeader(hdr); err != nil {
		return err
	}
	return nil
}

func (a *Artifact) addTemplate(templateFname, targetFname string, mode int64, data any) error {
	tpl, err := tplNew(templateFname, true)
	if err != nil {
		return err
	}
	var body bytes.Buffer
	err = tpl.Execute(&body, data)
	if err != nil {
		return err
	}
	hdr := &tar.Header{
		Typeflag: tar.TypeReg,
		Name:     targetFname,
		ModTime:  a.tnow,
		Size:     int64(body.Len()),
		Mode:     mode,
		Uid:      0,
		Gid:      0,
		Uname:    "root",
		Gname:    "root",
	}
	if err := a.tarOut.WriteHeader(hdr); err != nil {
		return err
	}
	if _, err := a.tarOut.Write(body.Bytes()); err != nil {
		return err
	}
	return nil
}

func (a *Artifact) etcCrictlYaml() error {
	if err := a.addStatic("static/crictl.yaml", "etc/crictl.yaml", 0644); err != nil {
		return err
	}
	return nil
}

func (a *Artifact) etcFlatcarUpdateConf() error {
	if err := a.addStatic("static/flatcar-update.conf", "etc/flatcar/update.conf", 0644); err != nil {
		return err
	}
	return nil
}

func (a *Artifact) etcI12eIfaceTxt() error {
	targetName := "etc/i12e/iface.txt"
	if a.cfg.Flannel == nil {
		log.Info("skipping: undefined 'flannel'", "targetName", targetName)
		return nil
	}
	if len(a.cfg.Flannel.Interface) < 1 {
		log.Info("skipping: undefined 'flannel.interface'", "targetName", targetName)
	}
	body := []byte(a.cfg.Flannel.Interface)
	hdr := &tar.Header{
		Typeflag: tar.TypeReg,
		Name:     targetName,
		ModTime:  a.tnow,
		Size:     int64(len(body)),
		Mode:     0600,
		Uid:      0,
		Gid:      0,
		Uname:    "root",
		Gname:    "root",
	}
	if err := a.tarOut.WriteHeader(hdr); err != nil {
		return err
	}
	if _, err := a.tarOut.Write(body); err != nil {
		return err
	}
	return nil
}

func (a *Artifact) etcI12eK3sConfigYaml() error {
	data := struct {
		I12eMode      string
		K3sToken      string
		K3sAgentToken string
		K3sUrl        string
		TlsSan        string
	}{
		I12eMode:      "",
		K3sToken:      a.cfg.K3sToken,
		K3sAgentToken: a.cfg.K3sAgentToken,
		K3sUrl:        a.cfg.K3sUrl,
		TlsSan:        a.cfg.TlsSan,
	}
	for _, m := range config.ValidModes() {
		data.I12eMode = m
		targetName := fmt.Sprintf("etc/i12e/k3s/config-%s.yaml", m)
		if err := a.addTemplate("k3s-config.yaml", targetName, 0600, &data); err != nil {
			return err
		}
	}
	return nil
}

func (a *Artifact) etcI12eK3sOverrideConf() error {
	for _, m := range config.ValidModes() {
		data := struct{ I12eMode string }{I12eMode: m}
		targetName := fmt.Sprintf("etc/i12e/k3s/override-%s.conf", m)
		if err := a.addTemplate("k3s-override.conf", targetName, 0644, &data); err != nil {
			return err
		}
	}
	return nil
}

func (a *Artifact) etcNftablesConf() error {
	data := struct{ PortKnocking []int }{PortKnocking: a.cfg.PortKnocking}
	if err := a.addTemplate("nftables.conf", "etc/nftables.conf", 0600, &data); err != nil {
		return err
	}
	return nil
}

func (a *Artifact) etcSystemdSystemConfDI12eConf() error {
	if err := a.addStatic("static/systemd-i12e.conf", "etc/systemd/system.conf.d/i12e.conf", 0644); err != nil {
		return err
	}
	return nil
}

func (a *Artifact) etcSystemdSystemK3sServiceDOverrideConf() error {
	if err := a.addEmpty("etc/systemd/system/k3s.service.d/override.conf", 0644); err != nil {
		return err
	}
	return nil
}

func (a *Artifact) etcSystemdSystemNftablesService() error {
	if err := a.addStatic("static/nftables.service", "etc/systemd/system/nftables.service", 0644); err != nil {
		return err
	}
	return nil
}

func (a *Artifact) etcRancherK3sConfigYaml() error {
	if err := a.addEmpty("etc/rancher/k3s/config.yaml", 0600); err != nil {
		return err
	}
	return nil
}

func (a *Artifact) optBinE() error {
	if err := a.addStatic("static/opt-bin-e", "opt/bin/e", 0755); err != nil {
		return err
	}
	return nil
}

func (a *Artifact) optLibexecI12eArtifactTuneSh() error {
	if err := a.addStatic("static/artifact-tune.sh", "opt/libexec/i12e/artifact-tune.sh", 0755); err != nil {
		return err
	}
	return nil
}

func Run(cfg *config.Config) error {
	if cfg == nil {
		return fmt.Errorf("artifact.Run(): undefined config")
	}
	a := newArtifact(cfg)
	funcList := []func() error{
		a.folders,
		a.etcCrictlYaml,
		a.etcFlatcarUpdateConf,
		a.etcI12eIfaceTxt,
		a.etcI12eK3sConfigYaml,
		a.etcI12eK3sOverrideConf,
		a.etcNftablesConf,
		a.etcSystemdSystemConfDI12eConf,
		a.etcSystemdSystemK3sServiceDOverrideConf,
		a.etcSystemdSystemNftablesService,
		a.etcRancherK3sConfigYaml,
		a.optBinE,
		a.optLibexecI12eArtifactTuneSh,
	}
	for _, f := range funcList {
		if err := f(); err != nil {
			a.Close()
			return err
		}
	}
	if err := a.Close(); err != nil {
		return err
	}
	fmt.Fprint(os.Stdout, a.Obuf().String())
	return nil
}
