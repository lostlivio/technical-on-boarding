#!/bin/sh


[ -z "$GOPATH" ] && echo "It looks like your \$GOPATH is empty." && exit 1

go get -u github.com/davecgh/go-spew/spew