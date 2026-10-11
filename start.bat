@echo off
rem English text UTF-8，English text GBK English text。
chcp 65001 >nul 2>&1
rem ARTEX English text（Windows）
rem
rem English text：
rem   start.bat                  English text（Ctrl-C English text）
rem   start.bat -addr :9000      English text artex
rem
rem English text：English text artex.exe English text，English text。
rem
rem   0      English text     -> English text
rem   75     English text     -> English text（English text"English text"English text"English text"）
rem   English text   English text             -> English text（1->2->4…English text 60 English text）
rem
rem English text、SHA256 English text、English text，English text artex English text
rem （selfupdate English text）。English text，English text start.sh English text。

setlocal enabledelayedexpansion
cd /d "%~dp0"

set "BIN=artex.exe"
if not exist "%BIN%" (
	echo [artex] English text %BIN% 1>&2
	exit /b 1
)

set "RESTART_CODE=75"
set "MAX_DELAY=60"
set /a delay=1

:loop
"%BIN%" %*
set "code=!ERRORLEVEL!"

if "!code!"=="0" (
	echo [artex] English text
	exit /b 0
)

if "!code!"=="%RESTART_CODE%" (
	rem English text/English text：English text artex English text。
	echo [artex] English text（English text）…
	set /a delay=1
	goto loop
)

echo [artex] English text ^(code=!code!^)，!delay!s English text 1>&2
rem timeout English text，English text ping English text（English text N English text N+1 English text）。
set /a pings=!delay!+1
ping -n !pings! 127.0.0.1 >nul 2>&1
set /a delay=!delay!*2
if !delay! gtr %MAX_DELAY% set /a delay=%MAX_DELAY%
goto loop
