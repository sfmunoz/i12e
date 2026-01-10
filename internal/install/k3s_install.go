package install

import (
	"github.com/sfmunoz/i12e/internal/cmdutil"
)

const script = `#!/bin/sh
[ -x /opt/libexec/i12e-k3s-install.sh ] && exec /opt/libexec/i12e-k3s-install.sh
exit 0
`

func K3sInstall() {
	cmdutil.RunCmd("/bin/sh", "-c", script)
}
