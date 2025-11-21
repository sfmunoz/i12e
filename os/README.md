# i12e os

**i12e** (infrastructure) **os** (Operating System)

- [Usage](#usage)

## Usage

Data from **values.yaml**:
```
$ helm template os .
```
Data from **values.yaml** overwritten with **secrets.yaml**:
```
$ helm template os . -f secrets://secrets.yaml
```
