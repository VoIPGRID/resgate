#!/bin/bash -e
# Run from directory above via ./scripts/cover.sh

go test -v -covermode=atomic -coverprofile=./cover.out -coverpkg=./server/... ./...
go tool cover -html=cover.out
