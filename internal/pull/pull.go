package pull

import (
	"github.com/sfmunoz/i12e/internal/cmdutil"
	"github.com/sfmunoz/logit"
)

var log = logit.Logit().
	WithLevel(logit.LevelInfo).
	With("mod", "i12e").
	With("pkg", "pull")

const script = `#!/bin/sh
FLAG_FILE="/etc/i12e/z.flag"
REBOOT_FILE="/etc/i12e/reboot-required"
set -e -o pipefail
function reboot_if_required {
  [ -f "$REBOOT_FILE" ] || return 0
  set -x
  systemctl daemon-reload
  [ -f /etc/systemd/system/k3s.service ] || touch /etc/i12e/k3s-install-required
  rm -f "$REBOOT_FILE"
  { set +x; } 2> /dev/null
  if [ -f "$REBOOT_FILE" ]
  then
    echo "error: reboot aborted: cannot delete '$REBOOT_FILE' before 'systemctl reboot' execution"
    return 1
  fi
  set -x
  systemctl reboot
}
function pull_if_needed {
  if [ -f "$FLAG_FILE" ]
  then
    echo "pull not needed: '${FLAG_FILE}' already exists"
  else
    echo "pulling artifact.tar.gz provided that '${FLAG_FILE}' doesn't exist..."
    set -x
    rclone cat rem:artifact.tar.gz | tar -C / -xvz
    touch "$REBOOT_FILE"
    { set +x; } 2>/dev/null
  fi
}
pull_if_needed
reboot_if_required
`

func Pull() {
	log.Info("Pull()...")
	cmdutil.RunCmd("/bin/sh", "-c", script)
}
