Unicode true
!include "MUI2.nsh"
!include "LogicLib.nsh"

Name "Sify PDV"
OutFile "dist\SifyPDV-Setup.exe"
InstallDir "$LOCALAPPDATA\Programs\Sify PDV"
RequestExecutionLevel user
VIProductVersion "1.0.0.0"
VIAddVersionKey /LANG=1046 "ProductName" "Sify PDV"
VIAddVersionKey /LANG=1046 "FileDescription" "Instalador do Sify PDV"
VIAddVersionKey /LANG=1046 "CompanyName" "Sify Sistemas"
VIAddVersionKey /LANG=1046 "FileVersion" "1.0.0"
!define APPKEY "Software\Microsoft\Windows\CurrentVersion\Uninstall\Sify PDV"
!define MUI_FINISHPAGE_RUN "$INSTDIR\SifyPDV.exe"
!define MUI_FINISHPAGE_RUN_TEXT "Abrir o Sify PDV"
!define MUI_FINISHPAGE_TEXT "O Sify PDV está pronto. Abra pela primeira vez com internet. Se você já usava o PDV no navegador, exporte o backup no navegador antigo e importe na nova janela antes de registrar vendas."
!insertmacro MUI_PAGE_WELCOME
!insertmacro MUI_PAGE_INSTFILES
!insertmacro MUI_PAGE_FINISH
!insertmacro MUI_UNPAGE_CONFIRM
!insertmacro MUI_UNPAGE_INSTFILES
!insertmacro MUI_LANGUAGE "PortugueseBR"

Section "Sify PDV"
  SetShellVarContext current
  SetOutPath "$INSTDIR"
  SetOverwrite on
  File "dist\SifyPDV.exe"
  WriteUninstaller "$INSTDIR\Desinstalar.exe"
  CreateShortcut "$SMPROGRAMS\Sify PDV.lnk" "$INSTDIR\SifyPDV.exe"
  CreateShortcut "$DESKTOP\Sify PDV.lnk" "$INSTDIR\SifyPDV.exe"
  WriteRegStr HKCU "${APPKEY}" "DisplayName" "Sify PDV"
  WriteRegStr HKCU "${APPKEY}" "DisplayVersion" "1.0.0"
  WriteRegStr HKCU "${APPKEY}" "Publisher" "Sify Sistemas"
  WriteRegStr HKCU "${APPKEY}" "DisplayIcon" "$INSTDIR\SifyPDV.exe"
  WriteRegStr HKCU "${APPKEY}" "UninstallString" '$\"$INSTDIR\Desinstalar.exe$\"'
  WriteRegDWORD HKCU "${APPKEY}" "NoModify" 1
  WriteRegDWORD HKCU "${APPKEY}" "NoRepair" 1
SectionEnd

Section "Uninstall"
  SetShellVarContext current
  Delete "$DESKTOP\Sify PDV.lnk"
  Delete "$SMPROGRAMS\Sify PDV.lnk"
  Delete "$INSTDIR\SifyPDV.exe"
  Delete "$INSTDIR\Desinstalar.exe"
  DeleteRegKey HKCU "${APPKEY}"
  RMDir "$INSTDIR"
SectionEnd
