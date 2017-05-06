#!/bin/bash

psql -v ON_ERROR_STOP=1 --username "${POSTGRES_USER}" < "${PANTSU_POSTGRES_DBFILE}"
