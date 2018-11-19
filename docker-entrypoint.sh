#!/bin/bash
set -e

if [ $# -eq 0 ]; then
    /usr/bin/app-service
else
    exec "$@"
fi
