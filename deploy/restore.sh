#!/bin/bash
# Restore the database from a postgres dump

pg_restore --username "${POSTGRES_USER}" -d ${POSTGRES_DB} "/${PANTSU_POSTGRES_DBFILE}"
