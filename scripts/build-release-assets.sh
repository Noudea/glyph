#!/usr/bin/env bash
set -euo pipefail

if [[ $# -lt 1 || $# -gt 2 ]]; then
  echo "Usage: $0 <version-tag> [output-dir]" >&2
  echo "Example: $0 v0.1.0 dist" >&2
  exit 1
fi

version_tag="$1"
out_dir="${2:-dist}"
version_no_v="${version_tag#v}"

project_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
case "${out_dir}" in
  /*) out_dir_abs="${out_dir}" ;;
  *) out_dir_abs="${project_root}/${out_dir}" ;;
esac
mkdir -p "${out_dir_abs}"

# Keep target matrix small and practical for now.
targets=(
  "darwin/amd64"
  "darwin/arm64"
  "linux/amd64"
  "linux/arm64"
  "windows/amd64"
  "windows/arm64"
)

tmp_dirs=()
cleanup() {
  for dir in "${tmp_dirs[@]}"; do
    rm -rf "${dir}"
  done
}
trap cleanup EXIT

for target in "${targets[@]}"; do
  IFS=/ read -r goos goarch <<<"${target}"

  tmp_dir="$(mktemp -d)"
  tmp_dirs+=("${tmp_dir}")

  bin_name="glyph"
  ext=""
  if [[ "${goos}" == "windows" ]]; then
    ext=".exe"
  fi

  out_bin="${tmp_dir}/${bin_name}${ext}"

  echo "Building ${goos}/${goarch}..."
  (
    cd "${project_root}"
    CGO_ENABLED=0 GOOS="${goos}" GOARCH="${goarch}" go build -trimpath -o "${out_bin}" ./cmd/glyph
  )

  if [[ "${goos}" == "windows" ]]; then
    asset="glyph_${version_no_v}_${goos}_${goarch}.zip"
    (
      cd "${tmp_dir}"
      zip -q "${out_dir_abs}/${asset}" "${bin_name}${ext}"
    )
  else
    asset="glyph_${version_no_v}_${goos}_${goarch}.tar.gz"
    (
      cd "${tmp_dir}"
      tar -czf "${out_dir_abs}/${asset}" "${bin_name}${ext}"
    )
  fi

  echo "Created ${out_dir_abs}/${asset}"
done

echo "All release assets are ready in ${out_dir_abs}/"
