#!/usr/bin/env bash
export LISTEN_PORT=8080
go get -v -d && go run /app/main.go
