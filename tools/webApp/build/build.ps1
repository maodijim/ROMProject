<#
Build multiple binaries and package static assets.
Outputs to: tools/webApp/build/dist/<os>-<arch>/
#>
param()

$ErrorActionPreference = 'Stop'
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$RepoRoot = Resolve-Path (Join-Path $ScriptDir "..\..\..") | Select-Object -ExpandProperty Path
$WebAppDir = Join-Path $RepoRoot "tools\webApp"
$OutDir = Join-Path $WebAppDir "build\dist"

if (-not (Test-Path $OutDir)) { New-Item -ItemType Directory -Path $OutDir | Out-Null }

# Optional frontend build (run if package.json exists) but do not bundle its output
$pkgJson = Join-Path $WebAppDir "package.json"
if (Test-Path $pkgJson) {
    Write-Host "Detected frontend package.json; attempting frontend build (output will not be bundled)..."
    Push-Location $WebAppDir
    if (Get-Command npm -ErrorAction SilentlyContinue) {
        try { npm ci } catch { try { npm install } catch { Write-Host "npm install failed" } }
        try { npm run build } catch { Write-Host "npm build failed or no build script" }
    } elseif (Get-Command yarn -ErrorAction SilentlyContinue) {
        try { yarn install --frozen-lockfile } catch { try { yarn } catch { Write-Host "yarn install failed" } }
        try { yarn build } catch { Write-Host "yarn build failed or no build script" }
    } else {
        Write-Host "No npm/yarn found; skipping frontend build."
    }
    Pop-Location
}

# Targets (match existing build.sh)
$targets = @(
    "darwin:amd64",
    "darwin:arm64",
    "linux:amd64",
    "windows:amd64"
)

foreach ($t in $targets) {
    $parts = $t -split ":"
    $goos = $parts[0]
    $goarch = $parts[1]

    $binDir = Join-Path $OutDir ("$goos-$goarch")
    if (-not (Test-Path $binDir)) { New-Item -ItemType Directory -Path $binDir | Out-Null }

    $binaryName = "webapp"
    if ($goos -eq "windows") { $binaryName = "$binaryName.exe" }

    Write-Host "Building Go webapp (GOOS=$goos GOARCH=$goarch) -> $binDir\$binaryName"
    $env:GOOS = $goos
    $env:GOARCH = $goarch
    & go build -o (Join-Path $binDir $binaryName) (Join-Path $WebAppDir)

    # create per-target compressed archive (binary only)
    if ($goos -eq "windows") {
        $zipPath = Join-Path $OutDir "$goos-$goarch.zip"
        if (Test-Path $zipPath) { Remove-Item $zipPath -Force }
        Compress-Archive -Path (Join-Path $binDir $binaryName) -DestinationPath $zipPath -Force
        Write-Host "Created $zipPath"
    } else {
        $tarPath = Join-Path $OutDir "$goos-$goarch.tar.gz"
        if (Test-Path $tarPath) { Remove-Item $tarPath -Force }
        if (Get-Command tar -ErrorAction SilentlyContinue) {
            # use tar to create gzipped archive containing only the binary
            & tar -C $binDir -czf $tarPath (Split-Path $binaryName -Leaf)
            Write-Host "Created $tarPath"
        } else {
            # fallback to Compress-Archive into zip if tar not available
            $zipFallback = Join-Path $OutDir "$goos-$goarch.zip"
            if (Test-Path $zipFallback) { Remove-Item $zipFallback -Force }
            Compress-Archive -Path (Join-Path $binDir $binaryName) -DestinationPath $zipFallback -Force
            Write-Host "tar not found; created $zipFallback instead"
        }
    }
}

# create a combined archive containing all per-target directories (no static)
$dirs = Get-ChildItem -Path $OutDir -Directory | Select-Object -ExpandProperty Name
if ($dirs.Count -gt 0) {
    # combined tar.gz
    $combinedTar = Join-Path $OutDir "webapp-all-platforms.tar.gz"
    if (Test-Path $combinedTar) { Remove-Item $combinedTar -Force }
    if (Get-Command tar -ErrorAction SilentlyContinue) {
        & tar -C $OutDir -czf $combinedTar $dirs
        Write-Host "Created $combinedTar"
    } else {
        Write-Host "tar not found; skipping combined tar.gz creation"
    }

    # combined zip (uses Compress-Archive)
    $combinedZip = Join-Path $OutDir "webapp-all-platforms.zip"
    if (Test-Path $combinedZip) { Remove-Item $combinedZip -Force }
    Push-Location $OutDir
    try {
        Compress-Archive -Path $dirs -DestinationPath $combinedZip -Force
        Write-Host "Created $combinedZip"
    } catch {
        Write-Host "Compress-Archive failed for combined zip: $_"
    } finally {
        Pop-Location
    }
}

Write-Host "Build & packaging complete: $OutDir"
