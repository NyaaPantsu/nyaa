# Create a backup from the database
#!/bin/bash

set -eu

NYAAPANTSU_USERNAME="$1"
NYAAPANTSU_PASSWORD="$2"
NYAAPANTSU_DB="$3"
NYAAPANTSU_PASSPHRASE_FILE="$4"

dump_file="${NYAAPANTSU_DB}_$(date +'%Y_%m_%d_%H_%M').backup"

pg_dump -U "${NYAAPANTSU_USERNAME}" -f "${dump_file}"

xz -z "${dump_file}"

compressed_dump_file="${dump_file}.xz"

gpg2 --batch --yes --passphrase-fd 0 \
     --output "${compressed_dump_file}.sig" \
     --detach-sig "${compressed_dump_file}" < "${NYAAPANTSU_PASSPHRASE_FILE}"
