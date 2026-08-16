#!/bin/bash
set -x
exit 0
cd "$(dirname "$0")"
exec bash ../plugin-run.sh wikijs
