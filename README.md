# i12e: infrastructure

**i12e** is an **infrastructure** [numeronym](https://en.wikipedia.org/wiki/Numeronym)

- [Guidelines](#guidelines)
- [Architecture](#architecture)
- [Requirements](#requirements)
  - [butane](#butane)
  - [age](#age)
  - [sops](#sops)
  - [helm](#helm)
  - [helm-secrets](#helm-secrets)
- [Genesis](#genesis)
- Modules (not in this page)
  - [os](os/README.md)
  - [dhcpd](dhcpd/README.md)

## Guidelines

Simple outline:

- [KISS](https://en.wikipedia.org/wiki/KISS_principle)
- Flexible
- Full control
- Cluster API inspired but...
  - Opinionated
  - Smaller
  - Flatcar OS
  - k3s for k8s
- Volatile: always pulls recent data from stable storage
- rclone is the tool to push-to/pull-from storage
- Nomad: easily move from one cloud provider to another
- backup-centric: it's the pillar of the system

## Architecture

**Notice**: be aware that this is an over-simplified architecture. Details will be provided as they are defined

```mermaid
flowchart LR
    github_repo["github.com/sfmunoz/i12e<br/>(repository)"]
    github_rel["github.com/sfmunoz/i12e<br/>(releases)"]
    rclone_conf["~/.config/rclone/rclone.conf<br/>(encrypted)"]
    local(["local"])
    fs[("fileserver<br/>s3, gcs, ...")] 
    host(["host (target)"])
    github_repo -->|"(1) git clone/pull"| local
    rclone_conf -->|"(2) config pull"| local
    local -->|"(3) config push<br/>(rclone)"| fs
    local -->|"(4) ignition push (ssh)<br/>i12e + rclone.conf"| host
    github_rel -->|"(5) i12e-flatcar.raw pull"| host
    fs -->|"(6) config pull"| host
```

Details:

- **(1)** Code is pulled from **github.com/sfmunoz/i12e**
- **(2)** Rclone config is read from **~/.config/rclone/rclone.conf** (encrypted)
- **(3)** Configuration is pushed to rclone-compatible storage
- **(4)** Ignition configuration is pushed to target host (**rclone.conf** included)
- **(5)** **i12e-flatcar.raw** is pulled from github
- **(6)** Target host pulls whatever is required from rclone-compatible storage

## Requirements

### butane

- https://coreos.github.io/butane/
- https://github.com/coreos/butane

```
$ brew install butane
```

### age

- https://age-encryption.org/
- https://github.com/FiloSottile/age

```
# apt install age
```

### sops

- https://getsops.io/
- https://github.com/getsops/sops

```
$ brew install sops
```

### helm

- https://helm.sh/
- https://github.com/helm/helm

```
$ brew install helm
```

### helm-secrets

[https://github.com/jkroepke/helm-secrets](https://github.com/jkroepke/helm-secrets)

```
$ helm plugin list
NAME    VERSION TYPE    APIVERSION      PROVENANCE      SOURCE
(... nothing ...)

$ helm plugin install --verify=false https://github.com/jkroepke/helm-secrets
WARNING: Skipping plugin signature verification
Installed plugin: secrets

$ helm plugin list
NAME    VERSION         TYPE            APIVERSION      PROVENANCE      SOURCE
secrets 4.8.0-dev       getter/v1       legacy          unknown         unknown
```

## Genesis

Help (artifact):
```
$ go run main.go artifact -h
Artifact management:

  - generation: tar+gz artifact
  - push to remote using rclone

Usage:
  i12e artifact [flags]

Flags:
  -h, --help   help for artifact

Global Flags:
  -p, --prod   Environment: 'prod' if set (default: 'dev')
```
Help (butane):
```
$ go run main.go butane -h
Run butane to generate ignition code

Examples:
  Reset flatcar host over ssh (default: '-o bash_b64'):
    $ i12e butane | ssh core@192.168.56.51 bash

  Generate ignition file:
    $ i12e butane -o ignition

Usage:
  i12e butane [flags]

Flags:
  -h, --help            help for butane
  -m, --mode string     Set target mode: ["main" "server" "agent"] (default "main")
  -o, --output string   Set output format: ["bash_b64" "bash_raw" "ignition" "debug"] (default "bash_b64")

Global Flags:
  -p, --prod   Environment: 'prod' if set (default: 'dev')
```
Artifact generation:
```
$ go run main.go artifact
2026-01-21T19:19:30.815Z 0d00h00m00.193s [I] rclonePush() remFile=d00:artifact.tar.gz
2026-01-21T19:19:32.660Z 0d00h00m02.037s [I] sha256(bef) sha256=7c3889fd21ccfb89c049f9300b1b3060fa3e23ffdd610b1b2635b23d38890965
2026-01-21T19:19:32.660Z 0d00h00m02.037s [I] sha256(aft) sha256=7c3889fd21ccfb89c049f9300b1b3060fa3e23ffdd610b1b2635b23d38890965
2026-01-21T19:19:32.667Z 0d00h00m02.044s [I] > tgz=drwx------ root/root         0 2026-01-21 19:19 etc/i12e
2026-01-21T19:19:32.667Z 0d00h00m02.044s [I] > tgz=drwx------ root/root         0 2026-01-21 19:19 etc/i12e/flags
2026-01-21T19:19:32.667Z 0d00h00m02.044s [I] > tgz=drwx------ root/root         0 2026-01-21 19:19 etc/i12e/k3s
2026-01-21T19:19:32.667Z 0d00h00m02.044s [I] > tgz=drwxr-xr-x root/root         0 2026-01-21 19:19 etc/systemd/system/k3s.service.d
2026-01-21T19:19:32.667Z 0d00h00m02.044s [I] > tgz=drwxr-xr-x root/root         0 2026-01-21 19:19 etc/systemd/system.conf.d
2026-01-21T19:19:32.667Z 0d00h00m02.044s [I] > tgz=drwxr-xr-x root/root         0 2026-01-21 19:19 opt/libexec
2026-01-21T19:19:32.667Z 0d00h00m02.044s [I] > tgz=drwxr-xr-x root/root         0 2026-01-21 19:19 opt/libexec/i12e
2026-01-21T19:19:32.667Z 0d00h00m02.044s [I] > tgz=-rw-r--r-- root/root       145 2026-01-21 19:19 etc/crictl.yaml
2026-01-21T19:19:32.667Z 0d00h00m02.044s [I] > tgz=-rw-r--r-- root/root        73 2026-01-21 19:19 etc/flatcar/update.conf
2026-01-21T19:19:32.667Z 0d00h00m02.044s [I] > tgz=-rw------- root/root         7 2026-01-21 19:19 etc/i12e/iface.txt
2026-01-21T19:19:32.667Z 0d00h00m02.044s [I] > tgz=-rw------- root/root       242 2026-01-21 19:19 etc/i12e/k3s/config-main.yaml
2026-01-21T19:19:32.667Z 0d00h00m02.044s [I] > tgz=-rw------- root/root       260 2026-01-21 19:19 etc/i12e/k3s/config-server.yaml
2026-01-21T19:19:32.667Z 0d00h00m02.044s [I] > tgz=-rw------- root/root        83 2026-01-21 19:19 etc/i12e/k3s/config-agent.yaml
2026-01-21T19:19:32.667Z 0d00h00m02.044s [I] > tgz=-rw-r--r-- root/root       129 2026-01-21 19:19 etc/i12e/k3s/override-main.conf
2026-01-21T19:19:32.667Z 0d00h00m02.044s [I] > tgz=-rw-r--r-- root/root       129 2026-01-21 19:19 etc/i12e/k3s/override-server.conf
2026-01-21T19:19:32.667Z 0d00h00m02.044s [I] > tgz=-rw-r--r-- root/root       128 2026-01-21 19:19 etc/i12e/k3s/override-agent.conf
2026-01-21T19:19:32.667Z 0d00h00m02.044s [I] > tgz=-rw------- root/root      2113 2026-01-21 19:19 etc/nftables.conf
2026-01-21T19:19:32.667Z 0d00h00m02.044s [I] > tgz=-rw-r--r-- root/root        68 2026-01-21 19:19 etc/systemd/system.conf.d/i12e.conf
2026-01-21T19:19:32.667Z 0d00h00m02.044s [I] > tgz=-rw-r--r-- root/root         0 2026-01-21 19:19 etc/systemd/system/k3s.service.d/override.conf
2026-01-21T19:19:32.667Z 0d00h00m02.044s [I] > tgz=-rw-r--r-- root/root       304 2026-01-21 19:19 etc/systemd/system/nftables.service
2026-01-21T19:19:32.667Z 0d00h00m02.044s [I] > tgz=-rw------- root/root         0 2026-01-21 19:19 etc/rancher/k3s/config.yaml
2026-01-21T19:19:32.667Z 0d00h00m02.044s [I] > tgz=-rwxr-xr-x root/root       178 2026-01-21 19:19 opt/bin/e
2026-01-21T19:19:32.667Z 0d00h00m02.044s [I] > tgz=-rwxr-xr-x root/root      1081 2026-01-21 19:19 opt/libexec/i12e/artifact-tune.sh
2026-01-21T19:19:32.667Z 0d00h00m02.044s [I] > tgz=-rw------- root/root         0 2026-01-21 19:19 etc/i12e/flags/artifact-pulled
2026-01-21T19:19:32.667Z 0d00h00m02.044s [I] > tgz=
```
Butane generation:
```
$ go run main.go butane
base64 -d <<< "H4sIA...(quite long base64 encoded gzipped script)...oIAAA=" | gunzip | bash
```
Butane injection (over ssh):
```
$ go run main.go butane -o bash_b64 | ssh core@192.168.56.51 bash
+ sudo rm -fv /oem/config.ign
removed '/oem/config.ign'
+ base64 -d
+ gunzip
+ sudo flatcar-reset --keep-machine-id --keep-paths '/etc/ssh/ssh_host_.*' /var/log /var/lib/rancher/k3s/agent/containerd -F /dev/stdin
WARNING: Running without --backup can cause data loss if the keep paths don't work as expected.
Also check whether your regex works as wanted with --preview-delete and --preview-keep.

Wrote machine ID as kernel cmdline parameter to /oem/grub.cfg
Removed any ignition.config.url kernel cmdline parameter in /oem/grub.cfg
Wrote Ignition file /oem/config.ign
Prepared /selective-os-reset and /boot/flatcar/first_boot
Staged OS reset, you can reboot now
+ sudo test -s /oem/config.ign
+ sudo jq . /oem/config.ign
{
  "ignition": {
    "version": "3.3.0"
  },
  (... ignition config ...)
}
+ sudo systemd-run bash -c 'sleep 1 ; systemctl reboot'
Running as unit: run-rb96ef8572bb2485e9ba0e96db33005c0.service; invocation ID: 314d5d8b3f144e8a923ff6ba0ba8b353
```
