#!/usr/bin/env sh
set -eu

REPO="${GLYPH_REPO:-Noudea/glyph}"
BINARY_NAME="glyph"
INPUT_VERSION="${1:-latest}"

require_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "error: required command not found: $1" >&2
    exit 1
  fi
}

resolve_os() {
  case "$(uname -s)" in
    Linux) echo "linux" ;;
    Darwin) echo "darwin" ;;
    *)
      echo "error: unsupported operating system: $(uname -s)" >&2
      exit 1
      ;;
  esac
}

resolve_arch() {
  case "$(uname -m)" in
    x86_64|amd64) echo "amd64" ;;
    arm64|aarch64) echo "arm64" ;;
    *)
      echo "error: unsupported architecture: $(uname -m)" >&2
      exit 1
      ;;
  esac
}

resolve_latest_tag() {
  api="https://api.github.com/repos/${REPO}/releases/latest"
  tag="$(curl -fsSL "${api}" | sed -n 's/.*"tag_name"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p' | head -n 1)"
  if [ -z "${tag}" ]; then
    echo "error: could not resolve latest release tag from ${api}" >&2
    exit 1
  fi
  printf '%s' "${tag}"
}

require_cmd curl
require_cmd tar
require_cmd uname
require_cmd mktemp

os="$(resolve_os)"
arch="$(resolve_arch)"

if [ "${INPUT_VERSION}" = "latest" ]; then
  tag="$(resolve_latest_tag)"
else
  case "${INPUT_VERSION}" in
    v*) tag="${INPUT_VERSION}" ;;
    *) tag="v${INPUT_VERSION}" ;;
  esac
fi

version_no_v="${tag#v}"
asset="glyph_${version_no_v}_${os}_${arch}.tar.gz"
url="https://github.com/${REPO}/releases/download/${tag}/${asset}"

install_dir="${GLYPH_INSTALL_DIR:-}"
if [ -z "${install_dir}" ]; then
  if [ -w /usr/local/bin ]; then
    install_dir="/usr/local/bin"
  else
    install_dir="${HOME}/.local/bin"
  fi
fi

mkdir -p "${install_dir}"

tmp_dir="$(mktemp -d)"
trap 'rm -rf "${tmp_dir}"' EXIT INT TERM

echo "Downloading ${asset} from ${tag}..."
curl -fL "${url}" -o "${tmp_dir}/${asset}"

tar -xzf "${tmp_dir}/${asset}" -C "${tmp_dir}"
if [ ! -f "${tmp_dir}/${BINARY_NAME}" ]; then
  echo "error: archive does not contain ${BINARY_NAME}" >&2
  exit 1
fi

chmod +x "${tmp_dir}/${BINARY_NAME}"

destination="${install_dir}/${BINARY_NAME}"
if command -v install >/dev/null 2>&1; then
  install -m 0755 "${tmp_dir}/${BINARY_NAME}" "${destination}"
else
  cp "${tmp_dir}/${BINARY_NAME}" "${destination}"
  chmod 0755 "${destination}"
fi

echo "Installed ${BINARY_NAME} ${tag} to ${destination}"
case ":${PATH}:" in
  *":${install_dir}:"*) ;;
  *)
    echo "Note: ${install_dir} is not in your PATH."
    ;;
esac
