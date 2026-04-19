#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SOURCE_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
MANIFEST_PATH="${SOURCE_ROOT}/dist/template-files/manifest.tsv"
TARGET_ROOT="${1:-}"

if [[ -z "${TARGET_ROOT}" ]]; then
  echo "Usage: ${0} <target-repo-root>" >&2
  exit 1
fi

if [[ ! -f "${MANIFEST_PATH}" ]]; then
  echo "Error: template file manifest not found at ${MANIFEST_PATH}" >&2
  exit 1
fi

while IFS='|' read -r source_rel target_rel mode; do
  if [[ -z "${source_rel}" ]] || [[ "${source_rel}" == \#* ]]; then
    continue
  fi

  source_path="${SOURCE_ROOT}/${source_rel}"
  target_path="${TARGET_ROOT}/${target_rel}"
  target_dir="$(dirname "${target_path}")"

  if [[ ! -f "${source_path}" ]]; then
    echo "Error: expected managed file not found at ${source_path}" >&2
    exit 1
  fi

  mkdir -p "${target_dir}"
  install -m "${mode}" "${source_path}" "${target_path}"
  echo "Synced ${target_rel} from ${source_rel}"
done < "${MANIFEST_PATH}"
