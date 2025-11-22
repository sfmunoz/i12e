# i12e os

**i12e** (infrastructure) **os** (Operating System)

- [Usage](#usage)

## Usage

Default (I12E_ENV=dev and I12_DEBUG=0):
```
$ ./os/build.sh
+ helm template dev . -f secrets://secrets.yaml
+ docker run --rm -i quay.io/coreos/butane:latest
+ ls -l os.json
-rw------- 1 sfm sfm 1105 Nov 22 07:28 os.json
```
Production:
```
$ I12E_ENV=prod ./os/build.sh
+ helm template prod . -f secrets://secrets.yaml
+ docker run --rm -i quay.io/coreos/butane:latest
+ ls -l os.json
-rw------- 1 sfm sfm 1106 Nov 22 07:28 os.json
```
Debug:
```
$ I12E_DEBUG=1 ./os/build.sh
+ exec helm template --debug dev . -f secrets://secrets.yaml
level=DEBUG msg="Original chart version" version=""
level=DEBUG msg="Chart path" path=/home/sfm/src/i12e/os
level=DEBUG msg="executing plugin command" pluginName=secrets command="/home/sfm/.local/share/helm/plugins/helm-secrets/scripts/run.sh downloader    secrets://secrets.yaml"
level=DEBUG msg="number of dependencies in the chart" dependencies=0
---
# Source: os/templates/flatcar.yaml
variant: flatcar
version: 1.0.0

systemd:
  units:
    - name: k3s-install.service
(...)
```
