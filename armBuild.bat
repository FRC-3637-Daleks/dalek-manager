@echo off
if not exist cross-bin mkdir cross-bin
set GOPATH=%cd%
set GOOS="linux"
set GOARCH="arm"
set GOARM=7
go build -o %GOPATH%\bin\robot-manager.exe %GOPATH%\src\main.go