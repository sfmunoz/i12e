#!/bin/bash
cd "$(dirname "$0")"
exec bash ../plugin-run.sh csi-rclone
