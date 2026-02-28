#!/bin/bash
export RCLONE_CONFIG="$(dirname "$0")/rclone.conf"
exec rclone "$@"
