!include "MUI2.nsh"

Name "gostim2"
OutFile "gostim2-setup.exe"
InstallDir "$LOCALAPPDATA\gostim2"
RequestExecutionLevel user

!define MUI_ABORTWARNING

!insertmacro MUI_PAGE_WELCOME
!insertmacro MUI_PAGE_DIRECTORY
!insertmacro MUI_PAGE_INSTFILES
!insertmacro MUI_PAGE_FINISH

!insertmacro MUI_UNPAGE_CONFIRM
!insertmacro MUI_UNPAGE_INSTFILES

!insertmacro MUI_LANGUAGE "English"

Section "Install"
    SetOutPath "$INSTDIR"
    
    # Files to include
    File "gostim2.exe"
    File "gostim2-gui.exe"
    File "README.txt"
    File /r "examples"

    # Create uninstaller
    WriteUninstaller "$INSTDIR\Uninstall.exe"

    # Create shortcuts
    CreateDirectory "$SMPROGRAMS\gostim2"
    CreateShortcut "$SMPROGRAMS\gostim2\gostim2-gui.lnk" "$INSTDIR\gostim2-gui.exe"
    CreateShortcut "$SMPROGRAMS\gostim2\Uninstall.lnk" "$INSTDIR\Uninstall.exe"
SectionEnd

Section "Uninstall"
    Delete "$INSTDIR\gostim2.exe"
    Delete "$INSTDIR\gostim2-gui.exe"
    Delete "$INSTDIR\README.txt"
    RMDir /r "$INSTDIR\examples"
    Delete "$INSTDIR\Uninstall.exe"

    RMDir "$INSTDIR"
    Delete "$SMPROGRAMS\gostim2\*.lnk"
    RMDir "$SMPROGRAMS\gostim2"
SectionEnd
