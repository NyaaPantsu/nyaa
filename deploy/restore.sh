#!/bin/bash
# Restore the database from a postgres dump


pg_restore --username "${POSTGRES_USER}" -d ${POSTGRES_DB} "/${PANTSU_POSTGRES_DBFILE}"

last_nyaa_torrent_id=923001
# Setting the sequence to start from after the latest nyaa torrent.
psql -U postgres "${POSTGRES_DB}" <<EOF
alter sequence torrents_torrent_id_seq start
                    ${last_nyaa_torrent_id} restart
EOF
