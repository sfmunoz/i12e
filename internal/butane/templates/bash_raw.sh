set -x -e -o pipefail
sudo rm -fv "{{ .ConfigIgn }}"
base64 -d <<< "{{ .Buf }}" | \
  gunzip | \
  sudo flatcar-reset \
    --keep-machine-id \
    --keep-paths \
    '/etc/ssh/ssh_host_.*' \
    /var/log \
    /var/lib/rancher/k3s/agent/containerd \
    /var/lib/rancher/k3s/storage \
    /home/core/.bash_history \
    /root/.bash_history \
    -F \
    /dev/stdin
sudo test -s "{{ .ConfigIgn }}"
sudo jq . "{{ .ConfigIgn }}"
sudo systemd-run bash -c 'sleep 1 ; systemctl reboot'
