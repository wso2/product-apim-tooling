@echo off

set ImportExportCLI_INSTALLER_DIR=importExportCLI\bin
del importExportCLI /s /q >nul 2>&1
rmdir msi /s /q >nul 2>&1
rmdir importExportCLI /s /q >nul 2>&1
mkdir %ImportExportCLI_INSTALLER_DIR%

set ImportExportCLI_VERSION=1.0

set GOOS=windows
set GOARCH=amd64

for /f %%x in ('wmic path win32_utctime get /format:list ^| findstr "="') do set %%x
set UTC_TIME=%Year%-%Month%-%Day% %Hour%:%Minute%:%Second% UTC
echo ImportExportCLI-%ImportExportCLI_VERSION% build started at '%UTC_TIME%' for '%GOOS%-%GOARCH%'
go build -gcflags=-trimpath=$GOPATH -asmflags=-trimpath=$GOPATH -ldflags "-X main.importExportCLIVersion=%ImportExportCLI_VERSION% -X 'main.buildDate=%UTC_TIME%'" ..\importExportCLI.go

copy importExportCLI.exe %ImportExportCLI_INSTALLER_DIR% >nul 2>&1
copy ..\resources\LICENSE.txt importExportCLI >nul 2>&1
copy ..\resources\README.txt importExportCLI >nul 2>&1

go run -ldflags "-X main.importExportCLIVersion=%ImportExportCLI_VERSION% -X main.importExportCLIPOS=windows -X main.importExportCLIPArch=x64" release.go
