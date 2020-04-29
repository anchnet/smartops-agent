@echo off
@echo CloudOps Agent Windows Build Package Script
@echo =====================================================================
@echo.





:: Get Passed Parameters
@echo %0 :: Get Passed Parameters...
@echo ---------------------------------------------------------------------
Set "Version=1.0"
:: First Parameter
if not "%~1"=="" (
    echo.%1 | FIND /I "=" > nul && (
        :: Named Parameter
        set "%~1"
    ) || (
        :: Positional Parameter
        set "Version=%~1"
    )
)

:: If Version not defined, Get the version from Git
if "%Version%"=="" (
    for /f "delims=" %%a in ('git describe') do @set "Version=%%a"
)

Set "CurDir=%~dp0"
Set "BldDir=%CurDir%buildenv"
Set "CnfDir=%CurDir%buildenv\conf"
Set "InsDir=%CurDir%installer"
Set "PreDir=%CurDir%prereqs"
for /f "delims=" %%a in ('git rev-parse --show-toplevel') do @set "SrcDir=%%a"

:: build
rd /S /Q "%BldDir%"
md %BldDir%
go build -o %BldDir%\agent.exe ..\..\..\cmd\main.go
md %CnfDir%
copy  "..\..\..\conf\cloudops.yaml.example" "%CnfDir%\cloudops.yaml"

:: Find the NSIS Installer
If Exist "C:\Program Files\NSIS\" (
    Set "NSIS=C:\Program Files\NSIS\"
) Else (
    Set "NSIS=C:\Program Files (x86)\NSIS\"
)
If not Exist "%NSIS%NSIS.exe" (
    @echo "NSIS not found in %NSIS%"
    exit /b 1
)

:: Add NSIS to the Path
Set "PATH=%NSIS%;%PATH%"
@echo.


@echo Copying SSM to buildenv
@echo ----------------------------------------------------------------------

:: Set the location of the ssm to download
Set Url64="https://repo.saltstack.com/windows/dependencies/64/ssm-2.24-103-gdee49fc.exe"
Set Url32="https://repo.saltstack.com/windows/dependencies/32/ssm-2.24-103-gdee49fc.exe"

:: Check for 64 bit by finding the Program Files (x86) directory
If Defined ProgramFiles(x86) (
    powershell -ExecutionPolicy RemoteSigned -File download_url_file.ps1 -url "%Url64%" -file "%BldDir%\ssm.exe"
) Else (
    powershell -ExecutionPolicy RemoteSigned -File download_url_file.ps1 -url "%Url32%" -file "%BldDir%\ssm.exe"
)
@echo.

:: Make sure the "prereq" directory exists and is empty
If Exist "%PreDir%" rd /s /q "%PreDir%"
mkdir "%PreDir%"

:: Copy down the 32 bit binaries
set Url60=http://repo.saltstack.com/windows/dependencies/32/ucrt/Windows6.0-KB2999226-x86.msu
set Name60=Windows6.0-KB2999226-x86.msu
set Url61=http://repo.saltstack.com/windows/dependencies/32/ucrt/Windows6.1-KB2999226-x86.msu
set Name61=Windows6.1-KB2999226-x86.msu
set Url80=http://repo.saltstack.com/windows/dependencies/32/ucrt/Windows8-RT-KB2999226-x86.msu
set Name80=Windows8-RT-KB2999226-x86.msu
set Url81=http://repo.saltstack.com/windows/dependencies/32/ucrt/Windows8.1-KB2999226-x86.msu
set Name81=Windows8.1-KB2999226-x86.msu
@echo - Downloading %Name60%
powershell -ExecutionPolicy RemoteSigned -File download_url_file.ps1 -url %Url60% -file "%PreDir%\%Name60%"
@echo - Downloading %Name61%
powershell -ExecutionPolicy RemoteSigned -File download_url_file.ps1 -url %Url61% -file "%PreDir%\%Name61%"
@echo - Downloading %Name80%
powershell -ExecutionPolicy RemoteSigned -File download_url_file.ps1 -url %Url80% -file "%PreDir%\%Name80%"
@echo - Downloading %Name81%
powershell -ExecutionPolicy RemoteSigned -File download_url_file.ps1 -url %Url81% -file "%PreDir%\%Name81%"




@echo Building the installer...
@echo ----------------------------------------------------------------------

:: Make the Salt Minion Installer
makensis.exe /DVersion=%Version%  "%InsDir%\Setup.nsi"
@echo.

@echo.
@echo ======================================================================
@echo Script completed...
@echo ======================================================================
@echo Installation file can be found in the following directory:
@echo %InsDir%

:done
if [%Version%] ==ls [] pause
