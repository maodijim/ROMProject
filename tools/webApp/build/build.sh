#!/usr/bin/env bash
set -euo pipefail
# Builds multiple binaries and packages them per-target.
# Does NOT bundle the `static` folder in any output.
# Outputs to: tools/webApp/build/dist/<os>-<arch>/ and creates per-target archives
# plus a combined archive tools/webApp/build/dist/webapp-all-platforms.{tar.gz,zip}

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../../.." && pwd)"
WEBAPP_DIR="$REPO_ROOT/tools/webApp"
OUTDIR="$WEBAPP_DIR/build/dist"

mkdir -p "$OUTDIR"

# Optional frontend build (run if package.json exists) but do not bundle its output
if [ -f "$WEBAPP_DIR/package.json" ]; then
  echo "Detected frontend package.json; attempting frontend build (output will not be bundled)..."
  pushd "$WEBAPP_DIR" > /dev/null
  if command -v npm >/dev/null 2>&1; then
    npm ci || npm install
    (npm run build || npm run build:prod) || true
  elif command -v yarn >/dev/null 2>&1; then
    yarn install --frozen-lockfile || yarn
    (yarn build || yarn build:prod) || true
  else
    echo "No npm/yarn found; skipping frontend build."
  fi
  popd > /dev/null
fi

# Target matrix: darwin/linux -> amd64 & arm64, windows -> amd64
targets=(
  "darwin:amd64"
  "darwin:arm64"
  "linux:amd64"
  "windows:amd64"
)

for t in "${targets[@]}"; do
  IFS=':' read -r GOOS GOARCH <<< "$t"
  BIN_DIR="$OUTDIR/${GOOS}-${GOARCH}"
  mkdir -p "$BIN_DIR"
  BINARY_NAME="webapp"
  if [ "$GOOS" = "windows" ]; then
    BINARY_NAME="${BINARY_NAME}.exe"
  fi

  echo "Building Go webapp (GOOS=$GOOS GOARCH=$GOARCH) -> $BIN_DIR/$BINARY_NAME"
  env GOOS="$GOOS" GOARCH="$GOARCH" go build -trimpath -ldflags="-s -w" -o "$BIN_DIR/$BINARY_NAME" "$WEBAPP_DIR"

  # create per-target compressed archive (binary only)
  if [ "$GOOS" = "windows" ]; then
    if command -v zip >/dev/null 2>&1; then
      echo "Creating zip archive for $GOOS-$GOARCH (binary only)"
      (cd "$OUTDIR" && zip -j -q "${GOOS}-${GOARCH}.zip" "${GOOS}-${GOARCH}/${BINARY_NAME}")
    else
      echo "zip not found; creating tar.gz with binary only"
      tar -C "$BIN_DIR" -czf "$OUTDIR/${GOOS}-${GOARCH}.tar.gz" "$(basename "$BINARY_NAME")"
    fi
  else
    echo "Creating tar.gz archive for $GOOS-$GOARCH (binary only)"
    tar -C "$BIN_DIR" -czf "$OUTDIR/${GOOS}-${GOARCH}.tar.gz" "$(basename "$BINARY_NAME")"
  fi
done

# create a combined archive containing all per-target directories (no static)
dirs=()
for d in "$OUTDIR"/*; do
  [ -d "$d" ] || continue
  name=$(basename "$d")
  # include only per-target dirs
  dirs+=("$name")
done

if [ "${#dirs[@]}" -gt 0 ]; then
  echo "Creating combined tar.gz archive for all platforms (no static)"
  tar -C "$OUTDIR" -czf "$OUTDIR/webapp-all-platforms.tar.gz" "${dirs[@]}"

  if command -v zip >/dev/null 2>&1; then
    echo "Also creating combined zip archive for all platforms (no static)"
    (cd "$OUTDIR" && zip -r -q "webapp-all-platforms.zip" "${dirs[@]}")
  fi
fi

echo "Build & packaging complete: $OUTDIR"
