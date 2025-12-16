# genesis

- [References](#references)
- [Install](#install)
- [Uninstall](#uninstall)

## References

- https://github.com/sfmunoz/k8s-playground/blob/main/crd-operator/README.md
- https://kopf.readthedocs.io/en/stable/install/
- https://github.com/nolar/kopf
  - https://github.com/nolar/kopf/tree/main/examples/01-minimal
- https://helm.sh/docs/chart_best_practices/custom_resource_definitions/

## Install
```
helm upgrade --install -n genesis --create-namespace genesis genesis
```

## Uninstall

[https://helm.sh/docs/chart_best_practices/custom_resource_definitions/](https://helm.sh/docs/chart_best_practices/custom_resource_definitions/): CRDs are not deleted

```
helm uninstall -n genesis genesis
kubectl delete namespaces genesis
kubectl delete crd kopfpeerings.kopf.dev
kubectl delete crd clusterkopfpeerings.kopf.dev
```
