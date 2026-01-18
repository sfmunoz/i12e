package butane

import (
	"embed"
	"os"
	"text/template"

	"github.com/sfmunoz/i12e/internal/config"
	"github.com/sfmunoz/logit"
)

var log = logit.Logit().
	WithLevel(logit.LevelInfo).
	With("mod", "i12e").
	With("pkg", "butane")

var i12eVersion = "v0.0.18"
var i12eSha256Sum = "cfe8d33bc00805344dbe4008d87b896ea0c3bb0618cc69bcf5bc0462af4a2709"

//go:embed templates/*.yaml
var FS embed.FS

func flatcarYamlRender() error {
	tpl := template.New("flatcar.yaml") // must much basename of the file
	tpl, err := tpl.Option("missingkey=error").ParseFS(FS, "templates/flatcar.yaml")
	if err != nil {
		return err
	}
	type FlatcarYaml struct {
		IgnitionConfigMergeLocal string
		I12eSha256sum            string
		Mode                     string
		I12eVersion              string
		RcloneConf               string
		SshAuthorizedKeys        []string
	}
	v := FlatcarYaml{
		IgnitionConfigMergeLocal: "**** IgnitionConfigMergeLocal ****",
		I12eVersion:              i12eVersion,
		I12eSha256sum:            i12eSha256Sum,
		Mode:                     "**** Mode ****",
		SshAuthorizedKeys: []string{
			"**** SshAuthorizedKey1 ****",
			"**** SshAuthorizedKey2 ****",
			"**** SshAuthorizedKey3 ****",
		},
		RcloneConf: "**** RcloneConf ****",
	}
	log.Info("**** flatcarYamlRender", "tpl-name", tpl.Name())
	err = tpl.Execute(os.Stdout, &v)
	if err != nil {
		return err
	}
	return nil
}

func Run(cfg *config.Config) error {
	log.Info("butane.Run()", "cfg", cfg)
	return flatcarYamlRender()
}
