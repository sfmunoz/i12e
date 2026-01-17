package butane

import (
	"embed"

	"github.com/sfmunoz/i12e/internal/cmdutil"
	"github.com/sfmunoz/logit"
	"github.com/spf13/viper"
)

var log = logit.Logit().
	WithLevel(logit.LevelInfo).
	With("mod", "i12e").
	With("pkg", "butane")

//go:embed templates/*.yaml
var FS embed.FS

func flatcarYamlDump() {
	buf, err := FS.ReadFile("templates/flatcar.yaml")
	if err != nil {
		log.Fatal("FS.ReadFile() failed", "err", err)
	}
	log.Info("flatcarYamlDump()", "buf", string(buf))
}

func Run(prod bool) {
	log.Info("butane.Run()", "prod", prod)
	fname := "secrets-dev.yaml"
	if prod {
		fname = "secrets-prod.yaml"
	}
	buf, err := cmdutil.SopsDecrypt(fname)
	if err != nil {
		log.Fatal("cmdutil.SopsDecrypt() failed", "err", err, "prod", prod)
	}
	log.Info("cmdutil.SopsDecrypt() OK", "buf", buf)
	viper.SetConfigType("yaml")
	viper.ReadConfig(buf)
	klist := viper.GetStringSlice("ssh_authorized_keys")
	for i, sshKey := range klist {
		log.Info("butane.Run()", "i", i, "sshKey", sshKey)
	}
	rcloneRemote := viper.Get("rclone_remote")
	log.Info("butane.Run()", "rcloneRemote", rcloneRemote)
	kubeVip := viper.GetStringMapString("kube_vip")
	for k, v := range kubeVip {
		log.Info("butane.Run() kubeVip", "k", k, "v", v)
	}
	flatcarYamlDump()
}
