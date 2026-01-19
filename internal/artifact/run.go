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

func tarGz() (*bytes.Buffer, error) {
	var buf bytes.Buffer
	gzOut := gzip.NewWriter(&buf)
	tarOut := tar.NewWriter(gzOut)
	var files = []struct {
		Name, Body string
	}{
		{"file1.txt", "1st file"},
		{"file2.txt", "2nd file"},
		{"file3.txt", "3rd and last file"},
	}
	tnow := time.Now().UTC()
	for _, file := range files {
		hdr := &tar.Header{
			Name:    file.Name,
			ModTime: tnow,
			Size:    int64(len(file.Body)),
			Mode:    0644,
			Uid:     0,
			Gid:     0,
			Uname:   "root",
			Gname:   "root",
		}
		if err := tarOut.WriteHeader(hdr); err != nil {
			return nil, err
		}
		if _, err := tarOut.Write([]byte(file.Body)); err != nil {
			return nil, err
		}
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
