#!/usr/bin/env bash

echo "Updating Vendor .."
dep ensure -v
echo "Installing Linter .."
go get github.com/alecthomas/gometalinter
gometalinter --install --update
echo "Installing CLI Tools .."
go install ./cmd/goreportcard-cli
echo "Done .."