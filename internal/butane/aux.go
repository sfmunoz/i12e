package butane

import (
	"bytes"
	"compress/gzip"
	"embed"
	"encoding/base64"
	"os/exec"
	"text/template"

	"github.com/sfmunoz/i12e/internal/config"
	"github.com/sfmunoz/i12e/internal/tplutil"
)

//go:embed templates/*.yaml templates/*.sh
var FS embed.FS

func butaneCmd(cfg *config.Config) *exec.Cmd {
	if cfg.Bout == config.BoutDebug {
		return exec.Command("butane", "-s", "-p")
	}
	return exec.Command("butane", "-s")
}

func tplNew(fname string, addFuncs bool) (*template.Template, error) {
	t := template.New(fname)
	if addFuncs {
		t = t.Funcs(tplutil.FuncMap())
	}
	return t.Option("missingkey=error").ParseFS(FS, "templates/"+fname)
}

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
