#!/usr/bin/env bash
set -euo pipefail
# Builds multiple binaries and packages static assets.
# Outputs to: tools/webApp/build/dist/<os>-<arch>/ and creates per-target archives
# plus a combined archive tools/webApp/build/dist/webapp-all-platforms.{tar.gz,zip}

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../../.." && pwd)"
WEBAPP_DIR="$REPO_ROOT/tools/webApp"
OUTDIR="$WEBAPP_DIR/build/dist"

mkdir -p "$OUTDIR"

# Optional frontend build (run once)
FRONTEND_OUT=""
if [ -f "$WEBAPP_DIR/package.json" ]; then
  echo "Detected frontend package.json; attempting frontend build..."
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
  # prefer common output dirs
  if [ -d "dist" ]; then
    FRONTEND_OUT="$WEBAPP_DIR/dist"
  elif [ -d "build" ]; then
    FRONTEND_OUT="$WEBAPP_DIR/build"
  elif [ -d "public" ]; then
    FRONTEND_OUT="$WEBAPP_DIR/public"
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
  env GOOS="$GOOS" GOARCH="$GOARCH" go build -o "$BIN_DIR/$BINARY_NAME" "$WEBAPP_DIR"

  # copy static folder from repo
  if [ -d "$WEBAPP_DIR/static" ]; then
    echo "Copying static -> $BIN_DIR/static"
    mkdir -p "$BIN_DIR/static"
    rsync -a --delete "$WEBAPP_DIR/static/" "$BIN_DIR/static/"
  fi

  # copy frontend build output (if any) into static (overlay)
  if [ -n "$FRONTEND_OUT" ] && [ -d "$FRONTEND_OUT" ]; then
    echo "Copying frontend build ($FRONTEND_OUT) -> $BIN_DIR/static/"
    mkdir -p "$BIN_DIR/static"
    rsync -a --delete "$FRONTEND_OUT/" "$BIN_DIR/static/"
  fi

  # create per-target compressed archive (zip for windows if zip available, tar.gz otherwise)
  if [ "$GOOS" = "windows" ]; then
    if command -v zip >/dev/null 2>&1; then
      echo "Creating zip archive for $GOOS-$GOARCH"
      (cd "$OUTDIR" && zip -r -q "${GOOS}-${GOARCH}.zip" "${GOOS}-${GOARCH}")
    else
      echo "zip not found; falling back to tar.gz for $GOOS-$GOARCH"
      tar -C "$OUTDIR" -czf "$OUTDIR/${GOOS}-${GOARCH}.tar.gz" "${GOOS}-${GOARCH}"
    fi
  else
    echo "Creating tar.gz archive for $GOOS-$GOARCH"
    tar -C "$OUTDIR" -czf "$OUTDIR/${GOOS}-${GOARCH}.tar.gz" "${GOOS}-${GOARCH}"
  fi
done

# create a combined archive containing all built directories
dirs=()
for d in "$OUTDIR"/*; do
  [ -d "$d" ] || continue
  dirs+=("$(basename "$d")")
done

if [ "${#dirs[@]}" -gt 0 ]; then
  echo "Creating combined tar.gz archive for all platforms"
  tar -C "$OUTDIR" -czf "$OUTDIR/webapp-all-platforms.tar.gz" "${dirs[@]}"

  if command -v zip >/dev/null 2>&1; then
    echo "Also creating combined zip archive for all platforms"
    (cd "$OUTDIR" && zip -r -q "webapp-all-platforms.zip" "${dirs[@]}")
  fi
fi

echo "Build & packaging complete: $OUTDIR"
