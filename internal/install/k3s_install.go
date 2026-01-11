package install

import (
	"github.com/sfmunoz/i12e/internal/cmdutil"
)

func K3sInstall() {
	cmdutil.RunCmd("/opt/libexec/i12e/k3s-install.sh")
}
