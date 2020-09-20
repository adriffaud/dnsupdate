#!/bin/bash

GOOS=linux GOARCH=arm GOARM=6 go build -o dnsupdate-armv6
GOOS=linux GOARCH=arm GOARM=7 go build -o dnsupdate-armv7
GOOS=linux GOARCH=amd64 go build -o dnsupdate-linux64
GOOS=darwin GOARCH=amd64 go build -o dnsupdate-darwin64
GOOS=windows GOARCH=amd64 go build -o dnsupdate-windows64.exe