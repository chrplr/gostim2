# Simple installation script for gostim2 on Windows
# Installs to the user's Local AppData folder.

$InstallDir = Join-Path $env:LOCALAPPDATA "gostim2"
$Executables = @("gostim2.exe", "gostim2-gui.exe")

# 1. Create the installation directory if it doesn't exist
if (-not (Test-Path $InstallDir)) {
    Write-Host "Creating installation directory: $InstallDir" -ForegroundColor Cyan
    New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null
}

# 2. Copy the executables and examples
foreach ($exe in $Executables) {
    if (Test-Path $exe) {
        Write-Host "Installing $exe..." -ForegroundColor Green
        Copy-Item -Path $exe -Destination $InstallDir -Force
    } else {
        Write-Warning "Could not find $exe in the current folder. Skipping."
    }
}

if (Test-Path "examples") {
    Write-Host "Installing examples..." -ForegroundColor Green
    Copy-Item -Path "examples" -Destination $InstallDir -Recurse -Force
}

# 3. Add to User PATH
$CurrentPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($CurrentPath -notlike "*$InstallDir*") {
    Write-Host "Adding $InstallDir to the user PATH..." -ForegroundColor Cyan
    $NewPath = "$CurrentPath;$InstallDir"
    [Environment]::SetEnvironmentVariable("Path", $NewPath, "User")
    Write-Host "Success! You may need to restart your terminal for changes to take effect." -ForegroundColor Green
} else {
    Write-Host "$InstallDir is already in the PATH." -ForegroundColor Yellow
}

Write-Host "`nInstallation complete. You can now run 'gostim2' or 'gostim2-gui' from any terminal." -ForegroundColor White
pause
