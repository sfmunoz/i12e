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

OS_JSON="build/os.json"

set -x
mkdir -p build
$HELM | docker run --rm -i quay.io/coreos/butane:latest > $OS_JSON
ls -l $OS_JSON
{ set +x; } 2> /dev/null

[ "$I12E_DIST" = "1" ] || exit 0

REPLY_OK="$(uuidgen | awk -F "-" '{ print $2 }')"

if [ "$REPLY_OK" = "" ]
then
  echo "error creating 'REPLY_OK'"
  exit 1
fi

echo
read -p "dist is a potentially dangerous operation... type '${REPLY_OK}' if you are sure: "

if [ "$REPLY" != "$REPLY_OK" ]
then
  echo
  echo "aborted!!"
  echo
  exit 1
fi

DIST_SH="build/dist.sh"

set -x
$HELM_DIST | yq -r .script > $DIST_SH
ls -l $DIST_SH
bash $DIST_SH $OS_JSON
{ set +x; } 2> /dev/null
