#!/bin/bash

set -e -o pipefail

function helm_cmd {
  PROF="$1"
  REL="$I12E_ENV"
  [ "$PROF" != "" ] && REL="${REL}-${PROF}"
  CMD="helm template"
  [ "$I12E_DEBUG" = "1" ] && CMD="${CMD} --debug"
  CMD="${CMD} ${REL} ."
  for FNAME in secrets.yaml ../secrets.yaml
  do
    [ -f "$FNAME" ] && CMD="${CMD} -f secrets://${FNAME}"
  done
  echo $CMD
}

cd "$(dirname "$0")"

[ "$I12E_ENV" != "prod" ] && I12E_ENV="dev"

HELM="$(helm_cmd)"
HELM_DIST="$(helm_cmd "dist-sh")"

if [ "$I12E_DEBUG" = "1" ]
then
  set -x
  $HELM
  $HELM_DIST
  { set +x; } 2> /dev/null
  exit 0
fi

set -x
mkdir -p build
$HELM | docker run --rm -i quay.io/coreos/butane:latest > build/os.json
ls -l build/os.json
{ set +x; } 2> /dev/null

[ "$I12E_DIST" = "1" ] || exit 0

set -x
$HELM_DIST | yq -r .script > build/dist.sh
ls -l build/dist.sh
{ set +x; } 2> /dev/null
