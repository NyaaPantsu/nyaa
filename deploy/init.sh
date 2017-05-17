#!/bin/bash

set -eux

# TODO Doesn't scale, find another way to wait until db is ready
if [[ "${PANTSU_DBTYPE}" = "postgres" ]]; then
  echo 'Waiting for the database to be ready...'
  sleep 40
fi

go get github.com/NyaaPantsu/nyaa
go build
./nyaa -host 0.0.0.0 \
       -port "${PANTSU_INTERNAL_PORT}" \
       -dbtype "${PANTSU_DBTYPE}" \
       -dbparams "${PANTSU_DBPARAMS}"
