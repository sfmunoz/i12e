package butane

import (
	"embed"
	"fmt"

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

func flatcarYamlDump() error {
	buf, err := FS.ReadFile("templates/flatcar.yaml")
	if err != nil {
		return err
	}
	log.Info("flatcarYamlDump()", "buf", string(buf))
	return nil
}

func Run(prod bool) error {
	log.Info("butane.Run()", "prod", prod)
	fname := "secrets-dev.yaml"
	if prod {
		fname = "secrets-prod.yaml"
	}
	buf, err := cmdutil.SopsDecrypt(fname)
	if err != nil {
		return fmt.Errorf("cmdutil.SopsDecrypt() failed: err=%s; prod=%t", err, prod)
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
	return flatcarYamlDump()
}
