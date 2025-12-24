# genesis

- [References](#references)
- [Install](#install)
- [Uninstall](#uninstall)
- [Helm OCI package build](#helm-oci-package-build)

## References

- https://github.com/sfmunoz/k8s-playground/blob/main/crd-operator/README.md
- https://kopf.readthedocs.io/en/stable/install/
- https://github.com/nolar/kopf
  - https://github.com/nolar/kopf/tree/main/examples/01-minimal
- https://helm.sh/docs/chart_best_practices/custom_resource_definitions/

## Install
```
$ helm upgrade --install -n genesis --create-namespace -f secrets://secrets.yaml genesis genesis
```

## Uninstall

[https://helm.sh/docs/chart_best_practices/custom_resource_definitions/](https://helm.sh/docs/chart_best_practices/custom_resource_definitions/): CRDs are not deleted by Helm

```
$ helm uninstall -n genesis genesis
$ kubectl delete namespaces genesis
$ kubectl delete crd kopfpeerings.kopf.dev
$ kubectl delete crd clusterkopfpeerings.kopf.dev
$ kubectl delete crd gdeployments.sfmunoz.com
```

## Helm OCI package build

**(1)** Build the package:
```
$ helm package genesis
Successfully packaged chart and saved it to: /home/sfm/src/i12e/genesis-0.1.0.tgz
```
**(2)** Generate TOKEN with `write:packages` permissions (**Settings > Developer settings > Personal access tokens**)

**(3)** Login to **ghcr.io** using that token (it's saved to **~/.config/helm/registry/config.json**):
```
$ helm registry login ghcr.io --username sfmunoz
Password: 
Login Succeeded
```
**(4)** Push:
```
$ helm push genesis-0.1.0.tgz oci://ghcr.io/sfmunoz
Pushed: ghcr.io/sfmunoz/genesis:0.1.0
Digest: sha256:93b32f63dd2d7d13ed4762344f1d9314da9e8a8f66b6c75276f5590c2f73a16b
```
**(5)** (optional) Logout:
```
$ helm registry logout ghcr.io
Removing login credentials for ghcr.io
```
**(6a)** Install
```
$ helm upgrade --install -f secrets://secrets.yaml -n genesis --create-namespace genesis oci://ghcr.io/sfmunoz/genesis --version 0.1.0
Release "genesis" does not exist. Installing it now.
Pulled: ghcr.io/sfmunoz/genesis:0.1.0
Digest: sha256:93b32f63dd2d7d13ed4762344f1d9314da9e8a8f66b6c75276f5590c2f73a16b
NAME: genesis
LAST DEPLOYED: Wed Dec 24 11:15:42 2025
NAMESPACE: genesis
STATUS: deployed
REVISION: 1
DESCRIPTION: Install complete
TEST SUITE: None
```
**(6b)** Install (without secrets):
```
$ helm upgrade --install -n genesis --set-json '{"env":{"dev":{"ssh_authorized_keys":["...ssh-public-key here..."]}}}' --create-namespace genesis oci://ghcr.io/sfmunoz/genesis --version 0.1.0
...
```

**(7)** (first time) Connect package to repository: **Packages > genesis**

**(8)** (first time) Make the package public: **Packages > genesis > Package settings > Change package visibility**
