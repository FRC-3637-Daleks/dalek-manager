@echo off
if not exist bin mkdir bin
set GOPATH=%cd%
go build -o %GOPATH%\bin\FRCUserProgram.exe %GOPATH%\src\main.go