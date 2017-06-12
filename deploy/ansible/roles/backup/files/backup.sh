# Create a backup from the database
#!/bin/bash

set -eu

NYAAPANTSU_USERNAME="$1"
NYAAPANTSU_PASSWORD="$2"
NYAAPANTSU_DB="$3"
NYAAPANTSU_PASSPHRASE_FILE="$4"
NYAAPANTSU_TRACKER="$5"
NYAAPANTSU_DOWNLOADED_DIR="$6"
NYAAPANTSU_WATCH_DIR="$7"

dump_file="${NYAAPANTSU_DB}_$(date +'%Y_%m_%d_%H_%M').backup"

pg_dump -U "${NYAAPANTSU_USERNAME}" -Fc --exclude-table-data=users -f "${dump_file}"

xz -z "${dump_file}"

compressed_dump_file="${dump_file}.xz"
signature_file="${compressed_dump_file}.sig"

gpg2 --batch --yes --passphrase-fd 0 \
     --output "${signature_file}" \
     --detach-sig "${compressed_dump_file}" < "${NYAAPANTSU_PASSPHRASE_FILE}"

mktorrent -a "${NYAAPANTSU_TRACKER}" \
          -c "Official nyaapantsu database release ($(date +'%Y-%m-%d'))" \
          "${compressed_dump_file}" "${signature_file}"

mv "${compressed_dump_file}" "${signature_file}" "${NYAAPANTSU_DOWNLOADED_DIR}"
mv "${compressed_dump_file}.torrent" "${NYAAPANTSU_WATCH_DIR}/"
