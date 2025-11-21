#!/bin/bash

set -e -o pipefail

cd "$(dirname "$0")"

[ "$ENV" != "prod" ] && ENV="dev"

HELM="helm template $ENV ."

for FNAME in secrets.yaml ../secrets.yaml
do
  [ -f "$FNAME" ] && HELM="${HELM} -f secrets://${FNAME}"
done

set -x
$HELM | docker run --rm -i quay.io/coreos/butane:latest > os.json

ls -l os.json
