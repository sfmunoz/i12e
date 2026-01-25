package butane

import (
	"embed"
	"text/template"

	"github.com/sfmunoz/i12e/internal/tplutil"
)

//go:embed templates/*.yaml templates/*.sh
var FS embed.FS

func tplNew(fname string, addFuncs bool) (*template.Template, error) {
	t := template.New(fname)
	if addFuncs {
		t = t.Funcs(tplutil.FuncMap())
	}
	return t.Option("missingkey=error").ParseFS(FS, "templates/"+fname)
}
