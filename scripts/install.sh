#!/bin/sh
# Instalador de gitnomad para macOS y Linux / gitnomad installer for macOS and Linux
#
#   curl -fsSL https://rodrigoquezadab.github.io/gitnomad/install.sh | sh
#   curl -fsSL https://rodrigoquezadab.github.io/gitnomad/install.sh | sh -s -- --dir ~/dev/github
#
# Opciones / options:
#   --dir RUTA           carpeta donde se clonarán los repos / folder to clone repos into
#   --version vX.Y.Z     versión a instalar (por defecto la última) / version (default: latest)
#   --install-dir RUTA   dónde dejar el ejecutable (por defecto ~/.local/bin)
#
# Variables: GITNOMAD_VERSION, GITNOMAD_INSTALL_DIR, GITNOMAD_DIR, GITNOMAD_REPO,
#            GITNOMAD_DOWNLOAD_URL (espejo / mirror), GITNOMAD_NO_MODIFY_PATH=1

set -eu

REPO="${GITNOMAD_REPO:-rodrigoquezadab/gitnomad}"
VERSION="${GITNOMAD_VERSION:-latest}"
INSTALL_DIR="${GITNOMAD_INSTALL_DIR:-$HOME/.local/bin}"
CLONE_DIR="${GITNOMAD_DIR:-}"

case "${LC_ALL:-${LC_MESSAGES:-${LANG:-}}}" in
  es*) ES=1 ;;
  *) ES=0 ;;
esac

# say "texto en español" "english text"
say() { if [ "$ES" = 1 ]; then printf '%s\n' "$1"; else printf '%s\n' "$2"; fi; }
fail() { if [ "$ES" = 1 ]; then printf 'error: %s\n' "$1" >&2; else printf 'error: %s\n' "$2" >&2; fi; exit 1; }

need_arg() { [ $# -ge 2 ] || fail "falta el valor de $1" "missing value for $1"; }

while [ $# -gt 0 ]; do
  case "$1" in
    --dir) need_arg "$@"; CLONE_DIR="$2"; shift 2 ;;
    --dir=*) CLONE_DIR="${1#--dir=}"; shift ;;
    --version) need_arg "$@"; VERSION="$2"; shift 2 ;;
    --version=*) VERSION="${1#--version=}"; shift ;;
    --install-dir) need_arg "$@"; INSTALL_DIR="$2"; shift 2 ;;
    --install-dir=*) INSTALL_DIR="${1#--install-dir=}"; shift ;;
    -h|--help) sed -n '2,15p' "$0" 2>/dev/null || true; exit 0 ;;
    *) fail "opción desconocida: $1" "unknown option: $1" ;;
  esac
done

# --- Sistema y arquitectura / OS and architecture ---
case "$(uname -s)" in
  Darwin) OS=darwin ;;
  Linux) OS=linux ;;
  MINGW*|MSYS*|CYGWIN*) fail "en Windows usa PowerShell: irm https://rodrigoquezadab.github.io/gitnomad/install.ps1 | iex" \
                             "on Windows use PowerShell: irm https://rodrigoquezadab.github.io/gitnomad/install.ps1 | iex" ;;
  *) fail "sistema no soportado: $(uname -s)" "unsupported OS: $(uname -s)" ;;
esac

case "$(uname -m)" in
  x86_64|amd64) ARCH=amd64 ;;
  arm64|aarch64) ARCH=arm64 ;;
  *) fail "arquitectura no soportada: $(uname -m)" "unsupported architecture: $(uname -m)" ;;
esac
# Mac con Apple Silicon corriendo la terminal bajo Rosetta: usar el binario nativo.
if [ "$OS" = darwin ] && [ "$ARCH" = amd64 ] && [ "$(sysctl -n hw.optional.arm64 2>/dev/null || echo 0)" = 1 ]; then
  ARCH=arm64
fi

ASSET="gitnomad-$OS-$ARCH"
if [ -n "${GITNOMAD_DOWNLOAD_URL:-}" ]; then
  BASE="${GITNOMAD_DOWNLOAD_URL%/}"
elif [ "$VERSION" = latest ]; then
  BASE="https://github.com/$REPO/releases/latest/download"
else
  BASE="https://github.com/$REPO/releases/download/$VERSION"
fi

# --- Descarga / download ---
if command -v curl >/dev/null 2>&1; then
  download() { curl -fsSL --retry 3 -o "$2" "$1"; }
elif command -v wget >/dev/null 2>&1; then
  download() { wget -q -O "$2" "$1"; }
else
  fail "se necesita curl o wget" "curl or wget is required"
fi

TMP="$(mktemp -d 2>/dev/null || mktemp -d -t gitnomad)"
trap 'rm -rf "$TMP"' EXIT INT TERM

say "Descargando $ASSET ($VERSION)..." "Downloading $ASSET ($VERSION)..."
download "$BASE/$ASSET" "$TMP/gitnomad" ||
  fail "no se pudo descargar $BASE/$ASSET" "could not download $BASE/$ASSET"

# --- Verificar SHA-256 / verify checksum ---
if download "$BASE/checksums.txt" "$TMP/checksums.txt" 2>/dev/null; then
  EXPECTED="$(awk -v f="$ASSET" '$2 == f || $2 == "*" f { print $1 }' "$TMP/checksums.txt")"
  if command -v sha256sum >/dev/null 2>&1; then
    ACTUAL="$(sha256sum "$TMP/gitnomad" | awk '{print $1}')"
  elif command -v shasum >/dev/null 2>&1; then
    ACTUAL="$(shasum -a 256 "$TMP/gitnomad" | awk '{print $1}')"
  else
    ACTUAL=""
  fi
  if [ -z "$EXPECTED" ]; then
    fail "checksums.txt no incluye $ASSET" "checksums.txt does not list $ASSET"
  elif [ -z "$ACTUAL" ]; then
    say "aviso: no hay sha256sum ni shasum; se omite la verificación" "warning: no sha256sum or shasum; skipping verification"
  elif [ "$EXPECTED" != "$ACTUAL" ]; then
    fail "el archivo descargado no coincide con su checksum (¿descarga corrupta?)" \
         "downloaded file does not match its checksum (corrupted download?)"
  else
    say "Checksum verificado." "Checksum verified."
  fi
else
  say "aviso: no se encontró checksums.txt; se omite la verificación" "warning: checksums.txt not found; skipping verification"
fi

# --- Instalar / install ---
mkdir -p "$INSTALL_DIR"
chmod +x "$TMP/gitnomad"
mv -f "$TMP/gitnomad" "$INSTALL_DIR/gitnomad"
BIN="$INSTALL_DIR/gitnomad"
say "Instalado en $BIN" "Installed to $BIN"
"$BIN" version || true

# --- PATH ---
case ":$PATH:" in
  *":$INSTALL_DIR:"*) ON_PATH=1 ;;
  *) ON_PATH=0 ;;
esac
if [ "$ON_PATH" = 0 ]; then
  case "$(basename "${SHELL:-sh}")" in
    zsh) RC="$HOME/.zshrc"; LINE="export PATH=\"$INSTALL_DIR:\$PATH\"" ;;
    bash) if [ "$OS" = darwin ]; then RC="$HOME/.bash_profile"; else RC="$HOME/.bashrc"; fi
          LINE="export PATH=\"$INSTALL_DIR:\$PATH\"" ;;
    fish) RC="$HOME/.config/fish/config.fish"; LINE="fish_add_path \"$INSTALL_DIR\"" ;;
    *) RC="$HOME/.profile"; LINE="export PATH=\"$INSTALL_DIR:\$PATH\"" ;;
  esac
  if [ "${GITNOMAD_NO_MODIFY_PATH:-0}" != 1 ]; then
    if ! grep -qsF "$LINE" "$RC"; then
      mkdir -p "$(dirname "$RC")"
      printf '\n# gitnomad\n%s\n' "$LINE" >> "$RC"
      say "Agregué $INSTALL_DIR al PATH en $RC" "Added $INSTALL_DIR to PATH in $RC"
    fi
    say "Abre una terminal nueva (o corre: $LINE) para usar 'gitnomad'." \
        "Open a new terminal (or run: $LINE) to use 'gitnomad'."
  else
    say "Agrega esto a tu $RC: $LINE" "Add this to your $RC: $LINE"
  fi
fi

# --- Carpeta de repos / repos folder ---
if [ -n "$CLONE_DIR" ]; then
  "$BIN" config set base_dir "$CLONE_DIR" >/dev/null
  say "Carpeta de repos: $CLONE_DIR" "Repos folder: $CLONE_DIR"
fi

# --- Requisitos / requirements ---
MISSING=0
if ! command -v git >/dev/null 2>&1; then
  MISSING=1
  if [ "$OS" = darwin ]; then
    say "Falta git: xcode-select --install" "git is missing: xcode-select --install"
  else
    say "Falta git: instálalo con tu gestor de paquetes (ej. sudo apt install git)" \
        "git is missing: install it with your package manager (e.g. sudo apt install git)"
  fi
fi
if ! command -v gh >/dev/null 2>&1; then
  MISSING=1
  if [ "$OS" = darwin ]; then
    say "Falta GitHub CLI: brew install gh" "GitHub CLI is missing: brew install gh"
  else
    say "Falta GitHub CLI: https://github.com/cli/cli#installation" "GitHub CLI is missing: https://github.com/cli/cli#installation"
  fi
fi

echo
say "Siguientes pasos:" "Next steps:"
[ "$MISSING" = 1 ] && say "  0. Instala lo que falta (arriba)" "  0. Install what's missing (above)"
N=1
step() { echo "  $N. $1"; N=$((N + 1)); }
step "gh auth login"
if [ -z "$CLONE_DIR" ]; then step "gitnomad config set base_dir ~/dev/github"; fi
step "gitnomad doctor --fix"
step "gitnomad clone"
