@echo off
set BIN=.\bin\symo.exe

for /F %%i in ('date /t') do set BUILD_DATE=%%i
for /F %%i in ('time /t') do set BUILD_TIME=%%i

git log --pretty="format:%%h" -n 1 > temp.txt
set /p GIT_HASH=<temp.txt
del temp.txt

set "LDFLAGS=-X main.buildDate=%BUILD_DATE%-%BUILD_TIME% -X main.gitHash=%GIT_HASH%"
echo %LDFLAGS%

go build -v -o %BIN% -ldflags "%LDFLAGS%" -tags pprof .\cmd\symo

set LOAD_BIN=.\bin\load.exe
go build -v -o %LOAD_BIN% .\cmd\load

call %LOAD_BIN%
