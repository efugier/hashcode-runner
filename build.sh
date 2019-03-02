#!/bin/sh
go build runner.go utils.go
GOOS=windows GOARCH=amd64 go build runner.go utils.go
