# Instalador de gitnomad para Windows / gitnomad installer for Windows
#
#   irm https://rodrigoquezadab.github.io/gitnomad/install.ps1 | iex
#
# Con carpeta de repos / with a repos folder:
#   & ([scriptblock]::Create((irm https://rodrigoquezadab.github.io/gitnomad/install.ps1))) -Dir 'C:\Dev\github'
#
# Variables: GITNOMAD_VERSION, GITNOMAD_INSTALL_DIR, GITNOMAD_DIR, GITNOMAD_REPO,
#            GITNOMAD_DOWNLOAD_URL (espejo / mirror), GITNOMAD_NO_MODIFY_PATH=1
#
# Nota: no usa 'exit' para no cerrar tu ventana de PowerShell si algo falla.
[CmdletBinding()]
param(
    [string]$Dir = $env:GITNOMAD_DIR,
    [string]$Version = $(if ($env:GITNOMAD_VERSION) { $env:GITNOMAD_VERSION } else { 'latest' }),
    [string]$InstallDir = $(if ($env:GITNOMAD_INSTALL_DIR) { $env:GITNOMAD_INSTALL_DIR } else { Join-Path $env:LOCALAPPDATA 'Programs\gitnomad' })
)

function Install-Gitnomad {
    param([string]$Dir, [string]$Version, [string]$InstallDir)

    $ErrorActionPreference = 'Stop'
    $ProgressPreference = 'SilentlyContinue'   # Invoke-WebRequest es mucho más rápido sin barra
    $es = (Get-Culture).TwoLetterISOLanguageName -eq 'es'
    function Say([string]$esText, [string]$enText) { if ($es) { Write-Host $esText } else { Write-Host $enText } }
    function Fail([string]$esText, [string]$enText) { if ($es) { throw $esText } else { throw $enText } }

    $repo = if ($env:GITNOMAD_REPO) { $env:GITNOMAD_REPO } else { 'rodrigoquezadab/gitnomad' }

    if ($env:OS -ne 'Windows_NT' -and $env:GITNOMAD_ALLOW_NON_WINDOWS -ne '1') {
        Fail 'Este instalador es para Windows. En macOS/Linux usa install.sh.' 'This installer is for Windows. On macOS/Linux use install.sh.'
    }
    if (-not $InstallDir) { Fail 'No se pudo determinar la carpeta de instalación.' 'Could not determine the install folder.' }

    # Windows PowerShell 5.1 a veces no tiene TLS 1.2 activado por defecto.
    try { [Net.ServicePointManager]::SecurityProtocol = [Net.ServicePointManager]::SecurityProtocol -bor [Net.SecurityProtocolType]::Tls12 } catch { }

    $cpu = if ($env:PROCESSOR_ARCHITEW6432) { $env:PROCESSOR_ARCHITEW6432 } else { $env:PROCESSOR_ARCHITECTURE }
    switch ($cpu) {
        'AMD64' { $arch = 'amd64' }
        'ARM64' { $arch = 'arm64' }
        default { Fail "Arquitectura no soportada: $cpu" "Unsupported architecture: $cpu" }
    }

    $asset = "gitnomad-windows-$arch.exe"
    if ($env:GITNOMAD_DOWNLOAD_URL) { $base = $env:GITNOMAD_DOWNLOAD_URL.TrimEnd('/') }
    elseif ($Version -eq 'latest') { $base = "https://github.com/$repo/releases/latest/download" }
    else { $base = "https://github.com/$repo/releases/download/$Version" }

    $tmp = Join-Path ([IO.Path]::GetTempPath()) ("gitnomad-" + [Guid]::NewGuid().ToString('N'))
    New-Item -ItemType Directory -Path $tmp -Force | Out-Null
    try {
        $file = Join-Path $tmp 'gitnomad.exe'
        Say "Descargando $asset ($Version)..." "Downloading $asset ($Version)..."
        try { Invoke-WebRequest -UseBasicParsing -Uri "$base/$asset" -OutFile $file }
        catch { Fail "No se pudo descargar $base/$asset : $($_.Exception.Message)" "Could not download $base/$asset : $($_.Exception.Message)" }

        # Verificar SHA-256
        $sums = Join-Path $tmp 'checksums.txt'
        $haveSums = $true
        try { Invoke-WebRequest -UseBasicParsing -Uri "$base/checksums.txt" -OutFile $sums } catch { $haveSums = $false }
        if ($haveSums) {
            $expected = $null
            foreach ($line in Get-Content $sums) {
                $parts = $line -split '\s+', 2
                if ($parts.Count -eq 2 -and $parts[1].TrimStart('*') -eq $asset) { $expected = $parts[0].ToLower() }
            }
            if (-not $expected) { Fail "checksums.txt no incluye $asset" "checksums.txt does not list $asset" }
            $actual = (Get-FileHash -Algorithm SHA256 -Path $file).Hash.ToLower()
            if ($actual -ne $expected) {
                Fail 'El archivo descargado no coincide con su checksum (¿descarga corrupta?).' 'Downloaded file does not match its checksum (corrupted download?).'
            }
            Say 'Checksum verificado.' 'Checksum verified.'
        } else {
            Say 'Aviso: no se encontró checksums.txt; se omite la verificación.' 'Warning: checksums.txt not found; skipping verification.'
        }

        New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
        $bin = Join-Path $InstallDir 'gitnomad.exe'
        try { Move-Item -Path $file -Destination $bin -Force }
        catch { Fail "No se pudo escribir $bin (¿gitnomad está abierto?): $($_.Exception.Message)" "Could not write $bin (is gitnomad running?): $($_.Exception.Message)" }
        Say "Instalado en $bin" "Installed to $bin"
    }
    finally {
        Remove-Item -Recurse -Force $tmp -ErrorAction SilentlyContinue
    }

    # PATH del usuario (permanente) y de esta sesión
    $sep = [IO.Path]::PathSeparator
    $sessionParts = $env:Path -split [regex]::Escape($sep) | ForEach-Object { $_.TrimEnd('\') }
    if ($sessionParts -notcontains $InstallDir.TrimEnd('\')) { $env:Path = "$InstallDir$sep$env:Path" }
    if ($env:OS -eq 'Windows_NT' -and $env:GITNOMAD_NO_MODIFY_PATH -ne '1') {
        $userPath = [Environment]::GetEnvironmentVariable('Path', 'User')
        $userParts = @()
        if ($userPath) { $userParts = $userPath -split ';' | Where-Object { $_ } | ForEach-Object { $_.TrimEnd('\') } }
        if ($userParts -notcontains $InstallDir.TrimEnd('\')) {
            $newPath = if ($userPath) { "$($userPath.TrimEnd(';'));$InstallDir" } else { $InstallDir }
            [Environment]::SetEnvironmentVariable('Path', $newPath, 'User')
            Say "Agregué $InstallDir a tu PATH." "Added $InstallDir to your PATH."
        }
    }

    & $bin version

    if ($Dir) {
        & $bin config set base_dir $Dir | Out-Null
        if ($LASTEXITCODE -ne 0) { Fail "No se pudo guardar la carpeta $Dir" "Could not save folder $Dir" }
        Say "Carpeta de repos: $Dir" "Repos folder: $Dir"
    }

    $missing = $false
    if (-not (Get-Command git -ErrorAction SilentlyContinue)) {
        $missing = $true
        Say 'Falta git:        winget install --id Git.Git -e' 'git is missing:        winget install --id Git.Git -e'
    }
    if (-not (Get-Command gh -ErrorAction SilentlyContinue)) {
        $missing = $true
        Say 'Falta GitHub CLI: winget install --id GitHub.cli -e' 'GitHub CLI is missing: winget install --id GitHub.cli -e'
    }

    Write-Host ''
    Say 'Siguientes pasos:' 'Next steps:'
    if ($missing) { Say '  0. Instala lo que falta (arriba) y abre una ventana nueva de PowerShell' "  0. Install what's missing (above) and open a new PowerShell window" }
    $steps = @('gh auth login')
    if (-not $Dir) { $steps += 'gitnomad config set base_dir C:\Dev\github' }
    $steps += 'gitnomad doctor --fix', 'gitnomad clone'
    for ($i = 0; $i -lt $steps.Count; $i++) { Write-Host ("  {0}. {1}" -f ($i + 1), $steps[$i]) }
}

Install-Gitnomad -Dir $Dir -Version $Version -InstallDir $InstallDir
