#! /bin/bash

[ ! -f .env ] || export $(grep -v '^#' .env | xargs)

go run ./cmd/main.go
