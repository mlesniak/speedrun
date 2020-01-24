#!/bin/bash

go get github.com/markbates/pkger/cmd/pkger
pkger -include /assets
go build

# go generate?
