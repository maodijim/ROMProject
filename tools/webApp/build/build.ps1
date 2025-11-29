# powershell
# File: `tools/webApp/build/build.ps1`
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

if (-not (Test-Path $OutDir)) {
    New-Item -ItemType Directory -Path $OutDir | Out-Null
}

# Optional frontend build (run once)
$frontendOut = $null
$pkgJson = Join-Path $WebAppDir "package.json"
if (Test-Path $pkgJson) {
    Write-Host "Detected frontend package.json; attempting frontend build..."
    Push-Location $WebAppDir
    if (Get-Command npm -ErrorAction SilentlyContinue) {
        try { npm ci } catch { npm install }
        try { npm run build } catch { Write-Host "npm build failed or no build script" }
    } elseif (Get-Command yarn -ErrorAction SilentlyContinue) {
        try { yarn install --frozen-lockfile } catch { yarn }
        try { yarn build } catch { Write-Host "yarn build failed or no build script" }
    } else {
        Write-Host "No npm/yarn found; skipping frontend build."
    }

    if (Test-Path (Join-Path $WebAppDir "dist")) { $frontendOut = Join-Path $WebAppDir "dist" }
    elseif (Test-Path (Join-Path $WebAppDir "build")) { $frontendOut = Join-Path $WebAppDir "build" }
    elseif (Test-Path (Join-Path $WebAppDir "public")) { $frontendOut = Join-Path $WebAppDir "public" }
    Pop-Location
}

# Target list
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

    # copy static folder from repo
    $staticSrc = Join-Path $WebAppDir "static"
    if (Test-Path $staticSrc) {
        Write-Host "Copying `static` -> $binDir\static"
        $dest = Join-Path $binDir "static"
        if (-not (Test-Path $dest)) { New-Item -ItemType Directory -Path $dest | Out-Null }
        robocopy $staticSrc $dest /MIR | Out-Null
    }
# File: `tools/webApp/build/build.ps1`
#!/usr/bin/env pwsh
<#
Build multiple binaries and package static assets.
Outputs to: tools/webApp/build/dist/<os>-<arch>/ and creates per-target zip files
plus a combined tools/webApp/build/dist/webapp-all-platforms.zip
#>
param()

$ErrorActionPreference = 'Stop'
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$RepoRoot = Resolve-Path (Join-Path $ScriptDir "..\..") | Select-Object -ExpandProperty Path
$WebAppDir = Join-Path $RepoRoot "tools\webApp"
$OutDir = Join-Path $WebAppDir "build\dist"

if (-not (Test-Path $OutDir)) {
    New-Item -ItemType Directory -Path $OutDir | Out-Null
}

# Optional frontend build (run once)
$frontendOut = $null
$pkgJson = Join-Path $WebAppDir "package.json"
if (Test-Path $pkgJson) {
    Write-Host "Detected frontend package.json; attempting frontend build..."
    Push-Location $WebAppDir
    if (Get-Command npm -ErrorAction SilentlyContinue) {
        try { npm ci } catch { npm install }
        try { npm run build } catch { Write-Host "npm build failed or no build script" }
    } elseif (Get-Command yarn -ErrorAction SilentlyContinue) {
        try { yarn install --frozen-lockfile } catch { yarn }
        try { yarn build } catch { Write-Host "yarn build failed or no build script" }
    } else {
        Write-Host "No npm/yarn found; skipping frontend build."
    }

    if (Test-Path (Join-Path $WebAppDir "dist")) { $frontendOut = Join-Path $WebAppDir "dist" }
    elseif (Test-Path (Join-Path $WebAppDir "build")) { $frontendOut = Join-Path $WebAppDir "build" }
    elseif (Test-Path (Join-Path $WebAppDir "public")) { $frontendOut = Join-Path $WebAppDir "public" }
    Pop-Location
}

# Target list
$targets = @(
    "darwin:amd64",
    "darwin:arm64",
    "linux:amd64",
    "linux:arm64",
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

    # copy static folder from repo
    $staticSrc = Join-Path $WebAppDir "static"
    if (Test-Path $staticSrc) {
        Write-Host "Copying static -> $binDir\static"
        $dest = Join-Path $binDir "static"
        if (-not (Test-Path $dest)) { New-Item -ItemType Directory -Path $dest | Out-Null }
        robocopy $staticSrc $dest /MIR | Out-Null
    }

    # copy frontend build output (if any) into static (overlay)
    if ($frontendOut -and (Test-Path $frontendOut)) {
        Write-Host "Copying frontend build ($frontendOut) -> $binDir\static"
        $dest = Join-Path $binDir "static"
        if (-not (Test-Path $dest)) { New-Item -ItemType Directory -Path $dest | Out-Null }
        robocopy $frontendOut $dest /MIR | Out-Null
    }

    # create per-target zip archive
    Push-Location $OutDir
    $srcPattern = "$($goos)-$($goarch)\*"
    $zipPath = Join-Path $OutDir "$($goos)-$($goarch).zip"
    if (Test-Path $zipPath) { Remove-Item $zipPath -Force }
    Compress-Archive -Path $srcPattern -DestinationPath $zipPath -Force
    Pop-Location
}

# create a combined zip archive containing all built directories
Push-Location $OutDir
$dirs = Get-ChildItem -Directory | Select-Object -ExpandProperty Name
if ($dirs.Count -gt 0) {
    $combined = Join-Path $OutDir "webapp-all-platforms.zip"
    if (Test-Path $combined) { Remove-Item $combined -Force }
    Compress-Archive -Path $dirs -DestinationPath $combined -Force
    Write-Host "Created combined archive: $combined"
}
Pop-Location

Write-Host "Build & packaging complete: $OutDir"

    # copy frontend build output (if any) into static (overlay)
    if ($frontendOut -and (Test-Path $frontendOut)) {
        Write-Host "Copying frontend build ($frontendOut) -> $binDir\static"
        $dest = Join-Path $binDir "static"
        if (-not (Test-Path $dest)) { New-Item -ItemType Directory -Path $dest | Out-Null }
        robocopy $frontendOut $dest /MIR | Out-Null
    }
}

Write-Host "Build complete: $OutDir"
