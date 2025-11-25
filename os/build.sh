#!/bin/bash

set -e -o pipefail

cd "$(dirname "$0")"

[ "$I12E_ENV" != "prod" ] && I12E_ENV="dev"

CMD="helm template ${I12E_ENV} ."
[ "$I12E_DEBUG" = "1" ] && CMD="${CMD} --debug"
for FNAME in secrets.yaml ../secrets.yaml
do
  [ -f "$FNAME" ] && CMD="${CMD} -f secrets://${FNAME}"
done

if [ "$I12E_DEBUG" = "1" ]
then
  set -x
  $CMD
  { set +x; } 2> /dev/null
  exit 0
fi

if [ "$I12E_DIST" = "1" ]
then
  REPLY_OK="$(uuidgen | awk -F "-" '{ print $2 }')"
  if [ "$REPLY_OK" = "" ] ; then echo "error creating 'REPLY_OK'" ; exit 1 ; fi
  echo
  read -p "dist is a potentially dangerous operation... type '${REPLY_OK}' if you are sure: "
  if [ "$REPLY" != "$REPLY_OK" ] ; then echo -en "\naborted!!\n\n" ; exit 1 ; fi
  export I12E_DIST
fi

export CMD
BUILD_SH="build/build.sh"

set -x
rm -rf build
mkdir build
$CMD | yq -r .script > $BUILD_SH
ls -l $BUILD_SH
bash $BUILD_SH
{ set +x; } 2> /dev/null
