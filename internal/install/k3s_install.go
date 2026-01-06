package install

import (
	"github.com/sfmunoz/i12e/internal/cmdutil"
)

const script = `#!/bin/sh
set -e -o pipefail
TRIGGER_FILE="/etc/i12e/k3s-install-required"
[ -f "${TRIGGER_FILE}" ] || exit 0
set -x
export INSTALL_K3S_VERSION="v1.34.3+k3s1"
curl -sfL https://get.k3s.io | sh -s
rm -f "${TRIGGER_FILE}"
`

func K3sInstall() {
	log.Info("K3sInstall()...")
	cmdutil.RunCmd("/bin/sh", "-c", script)
}
