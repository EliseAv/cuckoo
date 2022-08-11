@echo off
rem consider migrating to syso:
rem https://pkg.go.dev/github.com/hallazzang/syso

echo on
rsrc -manifest windows.manifest -ico check.ico
go build -ldflags="-H windowsgui"
