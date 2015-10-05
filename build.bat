@echo off
if not exist bin mkdir bin
set GOPATH=%cd%
go build -o %GOPATH%\bin\robot-manager.exe %GOPATH%\src\main.go