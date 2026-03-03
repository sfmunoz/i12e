#!/bin/bash

#
# Since I'm using multiple VCS at different levels (git and hg)
#   build failed: exit status 1: error obtaining VCS status: multiple VCS detected: ...
#   Use -buildvcs=false to disable VCS stamping.
#

[ "$GOFLAGS" = "" ] && GOFLAGS="-buildvcs=false"
export GOFLAGS
set -x
goreleaser "$@"
