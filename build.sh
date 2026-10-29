#!/usr/bin/env bash
# ============================================================================
#  build.sh - Active Stack
#  Cross-compila binarios estaticos (CGO desactivado) a la raiz del repo.
#  Uso:  ./build.sh        (chmod +x build.sh la primera vez)
#  Requiere: Go 1.26+ en el PATH.
#  Corre en Linux y WSL. Plataformas soportadas: Windows, Linux y WSL (amd64).
# ============================================================================
set -euo pipefail
cd "$(dirname "$0")"

OUT=.
export CGO_ENABLED=0

# build <goos> <goarch> [ext]
build() {
  local goos=$1 goarch=$2 ext=${3:-}
  echo "== ${goos} ${goarch} =="
  GOOS="$goos" GOARCH="$goarch" \
    go build -trimpath -o "${OUT}/active-stack_${goos}_${goarch}${ext}" ./cmd/active-stack
}

# Matriz de targets: las plataformas soportadas. macOS y Termux no se publican.
build windows amd64 .exe
build linux   amd64

echo
echo "Listo. Binarios en la raiz del repo:"
ls -lh active-stack_*

# ----------------------------------------------------------------------------
# Acceso directo en el escritorio (best-effort, idempotente).
# Solo Linux: WSL no tiene escritorio estandar. Falla blanda: nunca rompe el build.
# ----------------------------------------------------------------------------
desktop_shortcut() {
  local os; os="$(uname -s 2>/dev/null || echo unknown)"
  case "$os" in
    Linux)
      local bin="$PWD/active-stack_linux_amd64"
      local desk="${XDG_DESKTOP_DIR:-$HOME/Desktop}"
      [ -x "$bin" ]  || { echo "  (sin binario Linux nativo; omito acceso directo)"; return 0; }
      [ -d "$desk" ] || { echo "  (sin carpeta de escritorio en $desk; omito acceso directo)"; return 0; }
      local lnk="$desk/active-stack.desktop"
      printf '%s\n' \
        '[Desktop Entry]' \
        'Type=Application' \
        'Name=Active Stack' \
        "Exec=$bin" \
        "Path=$PWD" \
        'Terminal=true' \
        'Categories=Development;' > "$lnk"
      chmod +x "$lnk"
      echo "  -> $lnk"
      ;;
    *)
      echo "  (acceso directo automatico solo en Linux/Windows por ahora; SO: $os)"
      ;;
  esac
}

echo
echo "== Acceso directo en el escritorio =="
desktop_shortcut || true

echo
echo "Siguiente paso: ejecuta el binario de tu plataforma, por ejemplo:"
echo "  ./active-stack_linux_amd64 install"
echo "Instala los harnesses y copia el binario a un directorio de tu PATH: no hace falta copiarlo a mano."
