package artifact

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/sfmunoz/i12e/internal/cmdutil"
	"github.com/sfmunoz/i12e/internal/config"
	"github.com/sfmunoz/logit"
)

var log = logit.Logit().WithLevel(logit.LevelInfo)

type Artifact struct {
	cfg    *config.ArtifactConfig
	tnow   time.Time
	tarOut *tar.Writer
	gzOut  *gzip.Writer
	obuf   *bytes.Buffer
}

func newArtifact(cfg *config.ArtifactConfig) (*Artifact, error) {
	if cfg == nil {
		return nil, fmt.Errorf("newArtifact(): undefined config")
	}
	var obuf bytes.Buffer
	gzOut := gzip.NewWriter(&obuf)
	tarOut := tar.NewWriter(gzOut)
	return &Artifact{
		cfg:    cfg,
		tnow:   time.Now().UTC().Truncate(time.Second), // Truncate(): to prevent 'timestamp in the future' warnings
		tarOut: tarOut,
		gzOut:  gzOut,
		obuf:   &obuf,
	}, nil
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

func (a *Artifact) folders() error {
	flist := []struct {
		Name string
		Mode int64
	}{
		{Name: "etc/i12e", Mode: 0700},
		{Name: "etc/i12e/flags", Mode: 0700},
		{Name: "etc/i12e/k3s", Mode: 0700},
		{Name: "etc/i12e/flux", Mode: 0700},
		{Name: "etc/systemd/system.conf.d", Mode: 0755},
		{Name: "etc/wireguard", Mode: 0700}, // it's ok and harmless: it already exists like this
		{Name: "opt/libexec", Mode: 0755},
		{Name: "opt/libexec/i12e", Mode: 0755},
		{Name: "opt/libexec/i12e/plugins", Mode: 0755},
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
	if len(a.cfg.Mesh.EndpointInterface) < 1 {
		log.Info("skipping: undefined 'mesh.endpoint_interface'", "targetName", targetName)
		return nil
	}
	body := []byte(a.cfg.Mesh.EndpointInterface + "\n")
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
		K3sToken:      a.cfg.Artifact.K3sToken,
		K3sAgentToken: a.cfg.Artifact.K3sAgentToken,
		K3sUrl:        fmt.Sprintf("https://%s:6443", a.cfg.K3s.TlsSan),
		TlsSan:        a.cfg.K3s.TlsSan,
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

func (a *Artifact) etcI12eFluxCfg() error {
	data := struct {
		Cluster     string
		GitHubToken string
		AgeKey      string
	}{
		Cluster:     a.cfg.Env.String(),
		GitHubToken: a.cfg.Flux.GitHubToken,
		AgeKey:      a.cfg.Flux.AgeKey,
	}
	return a.addTemplate("flux.cfg", "etc/i12e/flux/flux.cfg", 0600, &data)
}

func (a *Artifact) etcNftablesConf() error {
	data := struct {
		PortKnocking []int
		EndpointPort int
	}{
		PortKnocking: a.cfg.Artifact.PortKnocking,
		EndpointPort: a.cfg.Mesh.EndpointPort,
	}
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

func (a *Artifact) etcSystemdSystemNftablesService() error {
	if err := a.addStatic("static/nftables.service", "etc/systemd/system/nftables.service", 0644); err != nil {
		return err
	}
	return nil
}

func (a *Artifact) etcWireguard() error {
	if len(a.cfg.WgConf) < 1 {
		log.Info("skipping: undefined 'wg_conf' array")
		return nil
	}
	for i, v := range a.cfg.WgConf {
		targetName := fmt.Sprintf("etc/wireguard/wg%d.conf", i)
		body := []byte(v)
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

func (a *Artifact) optLibexecI12eK3sInstallSh() error {
	if err := a.addStatic("static/k3s-install.sh", "opt/libexec/i12e/k3s-install.sh", 0755); err != nil {
		return err
	}
	return nil
}

func (a *Artifact) optLibexecI12ePlugins() error {
	const pluginsDir = "plugins"
	entries, err := os.ReadDir(pluginsDir)
	if err != nil {
		if os.IsNotExist(err) {
			log.Info("skipping: 'plugins/' folder not found")
			return nil
		}
		return err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		srcPath := filepath.Join(pluginsDir, entry.Name())
		body, err := os.ReadFile(srcPath)
		if err != nil {
			return err
		}
		targetName := "opt/libexec/i12e/plugins/" + entry.Name()
		hdr := &tar.Header{
			Typeflag: tar.TypeReg,
			Name:     targetName,
			ModTime:  a.tnow,
			Size:     int64(len(body)),
			Mode:     0644,
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
	}
	return nil
}

func (a *Artifact) etcI12eFlagsArtifactPulled() error {
	if err := a.addEmpty("etc/i12e/flags/artifact-pulled", 0600); err != nil {
		return err
	}
	return nil
}

func (a *Artifact) rclonePush() error {
	remFile := fmt.Sprintf("%s:artifact.tar.gz", a.cfg.Rclone.Remote)
	log.Info("rclonePush()", "remFile", remFile)
	// --------
	sha256_1 := fmt.Sprintf("%x", sha256.Sum256(a.obuf.Bytes())) // keep it before buffer read
	// --------
	cmd1 := exec.Command("rclone", "rcat", remFile)
	cmd1.Stdin = a.obuf
	bo1, be1, err1 := cmdutil.RunSimple(cmd1)
	if err1 != nil {
		return fmt.Errorf("'rclone rcat' failed': %s (stdout=%s, stderr=%s)", err1, bo1.String(), be1.String())
	}
	// --------
	cmd2 := exec.Command("rclone", "cat", remFile)
	bo2, be2, err2 := cmdutil.RunSimple(cmd2)
	if err2 != nil {
		return fmt.Errorf("'rclone cat' failed': %s (stdout=%s, stderr=%s)", err2, bo2.String(), be2.String())
	}
	// --------
	sha256_2 := fmt.Sprintf("%x", sha256.Sum256(bo2.Bytes()))
	log.Info("sha256(bef)", "sha256", sha256_1)
	log.Info("sha256(aft)", "sha256", sha256_2)
	if sha256_1 != sha256_2 {
		return fmt.Errorf("sha256 checksum mismatch")
	}
	// --------
	cmd3 := exec.Command("tar", "tvz")
	cmd3.Stdin = bo2
	bo3, be3, err3 := cmdutil.RunSimple(cmd3)
	if err3 != nil {
		return fmt.Errorf("'tar tvz' failed': %s (stdout=%s, stderr=%s)", err3, bo3.String(), be3.String())
	}
	for _, line := range strings.Split(bo3.String(), "\n") {
		log.Info(">", "tgz", line)
	}
	return nil
}

func (a *Artifact) run() error {
	funcList := []func() error{
		a.folders,
		a.etcCrictlYaml,
		a.etcFlatcarUpdateConf,
		a.etcI12eIfaceTxt,
		a.etcI12eK3sConfigYaml,
		a.etcI12eK3sOverrideConf,
		a.etcI12eFluxCfg,
		a.etcNftablesConf,
		a.etcSystemdSystemConfDI12eConf,
		a.etcSystemdSystemNftablesService,
		a.etcWireguard,
		a.optBinE,
		a.optLibexecI12eArtifactTuneSh,
		a.optLibexecI12eK3sInstallSh,
		a.optLibexecI12ePlugins,
		a.etcI12eFlagsArtifactPulled,
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
	if err := a.rclonePush(); err != nil {
		return err
	}
	return nil
}

func Run(cfg *config.ArtifactConfig) error {
	a, err := newArtifact(cfg)
	if err != nil {
		return err
	}
	return a.run()
}
