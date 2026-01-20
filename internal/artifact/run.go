package artifact

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"embed"
	"fmt"
	"os"
	"time"

	"github.com/sfmunoz/i12e/internal/config"
	"github.com/sfmunoz/logit"
)

var log = logit.Logit().WithLevel(logit.LevelInfo)

//go:embed static/* templates/*
var FS embed.FS

type Artifact struct {
	tnow   time.Time
	tarOut *tar.Writer
	gzOut  *gzip.Writer
	obuf   *bytes.Buffer
}

func newArtifact() *Artifact {
	var obuf bytes.Buffer
	gzOut := gzip.NewWriter(&obuf)
	tarOut := tar.NewWriter(gzOut)
	return &Artifact{
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

func (a *Artifact) etcCrictlYaml() error {
	body, err := FS.ReadFile("static/crictl.yaml")
	if err != nil {
		return err
	}
	hdr := &tar.Header{
		Typeflag: tar.TypeReg,
		Name:     "etc/crictl.yaml",
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
	return nil
}

func (a *Artifact) optBinE() error {
	body, err := FS.ReadFile("static/opt-bin-e")
	if err != nil {
		return err
	}
	hdr := &tar.Header{
		Typeflag: tar.TypeReg,
		Name:     "opt/bin/e",
		ModTime:  a.tnow,
		Size:     int64(len(body)),
		Mode:     0755,
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

func Run(cfg *config.Config) error {
	if cfg == nil {
		return fmt.Errorf("artifact.Run(): undefined config")
	}
	a := newArtifact()
	funcList := []func() error{a.folders, a.etcCrictlYaml, a.optBinE}
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
