@echo off

REM 设置版本号
set version=2.0.0

REM windows amd64
go build
rename gokins.exe gokins-%version%-windows-amd64.exe

REM linux amd64
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build
rename gokins gokins-%version%-linux-amd64

REM linux arm
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=arm
go build
rename gokins gokins-%version%-linux-arm

REM linux arm64
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=arm64
go build
rename gokins gokins-%version%-linux-arm64

REM macOS amd64
SET CGO_ENABLED=0
SET GOOS=darwin
SET GOARCH=amd64
go build
rename gokins gokins-%version%-darwin-amd64

REM macOS arm64
SET CGO_ENABLED=0
SET GOOS=darwin
SET GOARCH=arm64
go build
rename gokins gokins-%version%-darwin-arm64
