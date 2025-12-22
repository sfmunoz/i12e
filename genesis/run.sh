#!/bin/bash

set -e -o pipefail
set -x

cd "$(dirname "$0")"

exec docker run -it --rm \
  -v ./app:/app/genesis:ro \
  -e PYTHONUNBUFFERED=1 \
  -e PYTHONPATH=/app \
  -e GENESIS_LOCAL=1 \
  ghcr.io/sfmunoz/k8s-bulk:v1.6.0 \
  python3 -m genesis
