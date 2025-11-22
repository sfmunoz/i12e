#!/bin/bash

set -e -o pipefail

cd "$(dirname "$0")"

[ "$I12E_ENV" != "prod" ] && I12E_ENV="dev"

HELM="helm template"
[ "$I12E_DEBUG" = "1" ] && HELM="${HELM} --debug"
HELM="${HELM} $I12E_ENV ."

for FNAME in secrets.yaml ../secrets.yaml
do
  [ -f "$FNAME" ] && HELM="${HELM} -f secrets://${FNAME}"
done

if [ "$I12E_DEBUG" = "1" ]
then
  set -x
  exec $HELM
fi

set -x
mkdir -p build
$HELM | docker run --rm -i quay.io/coreos/butane:latest > build/os.json
ls -l build/os.json
