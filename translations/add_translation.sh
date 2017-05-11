#!/bin/bash

set -euo pipefail

id="$1"
translation="$2"

for file in *.json; do
  head -n -2 "${file}" > "${file}.tmp"
  echo -e "  },\n  {\n    \"id\": \"${id}\",\n    \"translation\": \"${translation}\"\n  }\n]" >> "${file}.tmp"
  mv "${file}.tmp" "${file}"
done
