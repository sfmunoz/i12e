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

func folders(tarOut *tar.Writer, tnow *time.Time) error {
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
			ModTime:  *tnow,
			Mode:     folder.Mode,
			Uid:      0,
			Gid:      0,
			Uname:    "root",
			Gname:    "root",
		}
		if err := tarOut.WriteHeader(hdr); err != nil {
			return err
		}
	}
	return nil
}

func tarGz() (*bytes.Buffer, error) {
	var buf bytes.Buffer
	gzOut := gzip.NewWriter(&buf)
	tarOut := tar.NewWriter(gzOut)
	tnow := time.Now().UTC()
	folders(tarOut, &tnow)
	var files = []struct {
		Name, Body string
	}{
		{"file1.txt", "1st file"},
		{"file2.txt", "2nd file"},
		{"file3.txt", "3rd and last file"},
	}
	for _, file := range files {
		hdr := &tar.Header{
			Typeflag: tar.TypeReg,
			Name:     file.Name,
			ModTime:  tnow,
			Size:     int64(len(file.Body)),
			Mode:     0644,
			Uid:      0,
			Gid:      0,
			Uname:    "root",
			Gname:    "root",
		}
		if err := tarOut.WriteHeader(hdr); err != nil {
			return nil, err
		}
		if _, err := tarOut.Write([]byte(file.Body)); err != nil {
			return nil, err
		}
	}
	body, err := FS.ReadFile("static/opt-bin-e")
	if err != nil {
		return nil, err
	}
	hdr := &tar.Header{
		Typeflag: tar.TypeReg,
		Name:     "opt/bin/e",
		ModTime:  tnow,
		Size:     int64(len(body)),
		Mode:     0755,
		Uid:      0,
		Gid:      0,
		Uname:    "root",
		Gname:    "root",
	}
	if err := tarOut.WriteHeader(hdr); err != nil {
		return nil, err
	}
	if _, err := tarOut.Write(body); err != nil {
		return nil, err
	}
	if err := tarOut.Close(); err != nil {
		return nil, err
	}
	if err := gzOut.Close(); err != nil {
		return nil, err
	}
	return &buf, nil
}

func Run(cfg *config.Config) error {
	if cfg == nil {
		return fmt.Errorf("artifact.Run(): undefined config")
	}
	buf, err := tarGz()
	if err != nil {
		return err
	}
	fmt.Fprint(os.Stdout, buf.String())
	return nil
}
