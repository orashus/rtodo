#!/usr/bin/env bash
# Install rtodo from GitHub Releases. No Go toolchain required.
#
#   curl -fsSL https://raw.githubusercontent.com/orashus/rtodo/main/install.sh | bash
#
# Optional env vars:
#   VERSION      release tag (default: latest), e.g. v1.0.0
#   INSTALL_DIR  directory to place the binary (default: /usr/local/bin or ~/.local/bin)

set -euo pipefail

REPO="${REPO:-orashus/rtodo}"
BINARY="rtodo"

need_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "error: '$1' is required" >&2
    exit 1
  fi
}

need_cmd curl
need_cmd tar
need_cmd uname
need_cmd mktemp

os="$(uname -s | tr '[:upper:]' '[:lower:]')"
arch="$(uname -m)"

case "$os" in
  linux | darwin) ;;
  *)
    echo "error: unsupported OS '$os' (supported: linux, darwin)" >&2
    exit 1
    ;;
esac

case "$arch" in
  x86_64 | amd64) arch="amd64" ;;
  aarch64 | arm64) arch="arm64" ;;
  *)
    echo "error: unsupported architecture '$arch' (supported: amd64, arm64)" >&2
    exit 1
    ;;
esac

asset="${BINARY}-${os}-${arch}.tar.gz"

if [ -n "${VERSION:-}" ]; then
  tag="$VERSION"
  base="https://github.com/${REPO}/releases/download/${tag}"
else
  tag="latest"
  base="https://github.com/${REPO}/releases/latest/download"
fi

url="${base}/${asset}"
sums_url="${base}/checksums.txt"

if [ -z "${INSTALL_DIR:-}" ]; then
  if [ -d /usr/local/bin ] && [ -w /usr/local/bin ]; then
    INSTALL_DIR="/usr/local/bin"
  else
    INSTALL_DIR="${HOME}/.local/bin"
  fi
fi

mkdir -p "$INSTALL_DIR"

tmp="$(mktemp -d)"
cleanup() { rm -rf "$tmp"; }
trap cleanup EXIT

echo "Downloading ${asset} (${tag})..."
if ! curl -fsSL "$url" -o "${tmp}/${asset}"; then
  echo "error: failed to download ${url}" >&2
  echo "hint: Contact Rash for publishing a GitHub Release with that asset, or set VERSION=vX.Y.Z" >&2
  exit 1
fi

if curl -fsSL "$sums_url" -o "${tmp}/checksums.txt" 2>/dev/null; then
  expected="$(grep -E "[[:space:]]${asset}\$" "${tmp}/checksums.txt" | awk '{print $1}')"
  if [ -n "$expected" ]; then
    if command -v sha256sum >/dev/null 2>&1; then
      actual="$(sha256sum "${tmp}/${asset}" | awk '{print $1}')"
    elif command -v shasum >/dev/null 2>&1; then
      actual="$(shasum -a 256 "${tmp}/${asset}" | awk '{print $1}')"
    else
      actual=""
    fi
    if [ -n "$actual" ] && [ "$actual" != "$expected" ]; then
      echo "error: checksum mismatch for ${asset}" >&2
      exit 1
    fi
  fi
fi

tar -xzf "${tmp}/${asset}" -C "$tmp"
if [ ! -f "${tmp}/${BINARY}" ]; then
  echo "error: archive did not contain '${BINARY}'" >&2
  exit 1
fi

chmod +x "${tmp}/${BINARY}"
dest="${INSTALL_DIR}/${BINARY}"
mv "${tmp}/${BINARY}" "$dest"

echo "Installed ${BINARY} to ${dest}"
if ! command -v "$BINARY" >/dev/null 2>&1; then
  echo "Add ${INSTALL_DIR} to your PATH, then run: ${BINARY} --version"
else
  "$BINARY" --version || true
fi
