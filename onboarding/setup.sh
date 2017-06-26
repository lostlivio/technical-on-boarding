#!/bin/sh


[ -z "$GOPATH" ] && echo "It looks like your \$GOPATH is empty." && exit 1


go get github.com/satori/go.uuid
go get github.com/google/go-github/github
go get gopkg.in/yaml.v2

