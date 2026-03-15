#!/bin/sh
[ -f /etc/i12e/flags/artifact-pulled ] && exit 0
set -x -e -o pipefail
rclone cat rem:artifact.tar.gz | tar -C / -xvz
rm -f /etc/i12e/flags/artifact-tuned
