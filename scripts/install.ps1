param(
  [Parameter(Position = 0)]
  [string]$Version = "latest"
)

$ErrorActionPreference = "Stop"

$repo = if ($env:GLYPH_REPO) { $env:GLYPH_REPO } else { "Noudea/glyph" }
$binaryName = "glyph.exe"

if ($env:OS -ne "Windows_NT") {
  throw "This installer is for Windows only."
}

function Resolve-Architecture {
  $arch = [System.Runtime.InteropServices.RuntimeInformation]::OSArchitecture.ToString().ToLowerInvariant()
  switch ($arch) {
    "x64" { return "amd64" }
    "arm64" { return "arm64" }
    default { throw "Unsupported architecture: $arch" }
  }
}

function Resolve-Tag([string]$inputVersion, [string]$repository) {
  if ($inputVersion -eq "latest") {
    $api = "https://api.github.com/repos/$repository/releases/latest"
    $resp = Invoke-RestMethod -Uri $api -Headers @{ "User-Agent" = "glyph-installer" }
    if (-not $resp.tag_name) {
      throw "Could not resolve latest release tag from $api"
    }
    return [string]$resp.tag_name
  }

  if ($inputVersion.StartsWith("v")) {
    return $inputVersion
  }

  return "v$inputVersion"
}

function Ensure-PathContains([string]$pathToAdd) {
  $normalize = {
    param([string]$p)
    if (-not $p) { return "" }
    return $p.Trim().TrimEnd('\\').ToLowerInvariant()
  }

  $userPath = [Environment]::GetEnvironmentVariable("Path", "User")
  $entries = @()
  if ($userPath) {
    $entries = $userPath -split ';' | Where-Object { $_ -and $_.Trim() -ne "" }
  }

  $normalizedEntries = $entries | ForEach-Object { & $normalize $_ }
  $normalizedTarget = & $normalize $pathToAdd

  if ($normalizedEntries -notcontains $normalizedTarget) {
    $newPath = if ($userPath) { "$userPath;$pathToAdd" } else { $pathToAdd }
    [Environment]::SetEnvironmentVariable("Path", $newPath, "User")
    Write-Host "Updated user PATH with $pathToAdd"
  }

  $runtimeEntries = ($env:Path -split ';' | ForEach-Object { & $normalize $_ })
  if ($runtimeEntries -notcontains $normalizedTarget) {
    $env:Path = if ($env:Path) { "$env:Path;$pathToAdd" } else { $pathToAdd }
  }
}

$arch = Resolve-Architecture
$tag = Resolve-Tag -inputVersion $Version -repository $repo
$versionNoV = $tag.TrimStart("v")
$asset = "glyph_${versionNoV}_windows_${arch}.zip"
$url = "https://github.com/$repo/releases/download/$tag/$asset"

$installDir = if ($env:GLYPH_INSTALL_DIR) {
  $env:GLYPH_INSTALL_DIR
} else {
  Join-Path $env:LOCALAPPDATA "Programs\\glyph\\bin"
}

New-Item -ItemType Directory -Path $installDir -Force | Out-Null

$tempDir = Join-Path ([System.IO.Path]::GetTempPath()) ("glyph-install-" + [Guid]::NewGuid().ToString("N"))
New-Item -ItemType Directory -Path $tempDir -Force | Out-Null

try {
  $archivePath = Join-Path $tempDir $asset
  Write-Host "Downloading $asset from $tag..."
  Invoke-WebRequest -Uri $url -OutFile $archivePath -Headers @{ "User-Agent" = "glyph-installer" }

  Expand-Archive -Path $archivePath -DestinationPath $tempDir -Force

  $sourceBinary = Join-Path $tempDir $binaryName
  if (-not (Test-Path $sourceBinary)) {
    throw "Archive does not contain $binaryName"
  }

  $destination = Join-Path $installDir $binaryName
  Copy-Item -Path $sourceBinary -Destination $destination -Force

  Ensure-PathContains -pathToAdd $installDir

  Write-Host "Installed glyph $tag to $destination"
  Write-Host "Run: glyph"
}
finally {
  Remove-Item -Path $tempDir -Recurse -Force -ErrorAction SilentlyContinue
}
