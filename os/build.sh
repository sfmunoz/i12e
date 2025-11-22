#!/bin/bash

set -e -o pipefail

cd "$(dirname "$0")"

[ "$I12E_ENV" != "prod" ] && I12E_ENV="dev"

HELM="helm template $I12E_ENV ."

for FNAME in secrets.yaml ../secrets.yaml
do
  [ -f "$FNAME" ] && HELM="${HELM} -f secrets://${FNAME}"
done

set -x
$HELM | docker run --rm -i quay.io/coreos/butane:latest > os.json

ls -l os.json
