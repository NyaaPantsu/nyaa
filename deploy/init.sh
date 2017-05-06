#!/bin/bash

set -ex

# TODO Doesn't scale, find another way to wait until db is ready
sleep 40
go get github.com/ewhal/nyaa
go build
./nyaa -host 0.0.0.0 \
       -port "${PANTSU_INTERNAL_PORT}" \
       -dbtype "${PANTSU_DBTYPE}" \
       -dbparams "${PANTSU_DBPARAMS}"
