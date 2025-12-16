# i12e: infrastructure

**i12e** is an **infrastructure** [numeronym](https://en.wikipedia.org/wiki/Numeronym)

- [Architecture](#architecture)
- [Requirements](#requirements)
  - [helm](#helm)
  - [age](#age)
  - [sops](#sops)
  - [helm-secrets](#helm-secrets)
  - [yq](#yq)
- Modules (not in this page)
  - [os](os/README.md)
  - [dhcpd](dhcpd/README.md)

## Architecture

```mermaid
flowchart LR
    u(["User"])
    ks["K3S<br/>SQLite3"]
    ke1["K3S etcd<br/>master<br/>server-1"]
    ke2["K3S etcd<br/>server-2"]
    ke3["K3S etcd<br/>server-3"]
    ke4["K3S etcd<br/>agent-1"]
    ke5["K3S etcd<br/>agent-2"]
    fs[("fileserver<br/>s3, gcs, ...")] 
    d{"Deploy"}
    u -->|"kubectl<br/>helm<br/>..."| d
    d -->|"spark<br/>+<br/>restore"| ks
    ks -->|"ignition<br/>+<br/>reboot"| ke1
    d -->|"os<br/>deploy"| ke1
    ke1 -.->|cloud API<br/>+<br/>ignition| ke2
    ke1 -.->|cloud API<br/>+<br/>ignition| ke3
    ke1 -.->|cloud API<br/>+<br/>ignition| ke4
    ke1 -.->|cloud API<br/>+<br/>ignition| ke5
    ks -.-|"restore<br/>(rclone)"| fs -.-|"backup<br/>(rclone)"| ke1
```

## Requirements

### helm

- https://helm.sh/
- https://github.com/helm/helm

```
$ brew install helm
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

### yq

- https://mikefarah.gitbook.io/yq/
- https://github.com/kislyuk/yq
```
# apt install yq
```
