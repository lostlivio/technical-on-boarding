
PATH:="$(PATH):$(GOPATH)/bin/:`pwd`/go/bin/"

all: vet lint test build


build: prepare_workload

# TODO: use glide to populate vendored dependencies


setup: 
	@go version
	@echo GOPATH IS ${GOPATH}
	go get github.com/satori/go.uuid
	go get -u -fix github.com/google/go-github/github
	go get gopkg.in/yaml.v2
	go get golang.org/x/oauth2


prepare_workload: setup ./github/issues/*.go ./cmd/*.go
	go build -v -o $@  github.com/samsung-cnct/technical-on-boarding/onboarding/cmd

test: setup vet lint
	go test -race -v github.com/samsung-cnct/technical-on-boarding/onboarding/github/issues/
	
coverage.html: $(shell find ./ -name '*.go')
	go test -covermode=count -coverprofile=coverage.prof github.com/samsung-cnct/technical-on-boarding/onboarding/github/issues/
	go tool cover -html=coverage.prof -o $@

test-cover: coverage.html

lint: setup
	go get -u github.com/golang/lint/golint
	go get -u golang.org/x/tools/cmd/goimports
	go get -u honnef.co/go/tools/cmd/gosimple
	gofmt -w -s ./github/issues/
	gofmt -w -s ./cmd/
	$(GOPATH)/bin/goimports -w ./github/issues/
	$(GOPATH)/bin/goimports -w ./cmd/
	$(GOPATH)/bin/golint ./github/issues/
	$(GOPATH)/bin/golint ./cmd/
	$(GOPATH)/bin/gosimple ./github/issues/
	$(GOPATH)/bin/gosimple ./cmd/
	
vet: 
	go vet -v github.com/samsung-cnct/technical-on-boarding/onboarding/github/issues/

clean:
	-rm ./coverage.* ./prepare_workload


godoc.txt: $(shell find ./ -name '*.go')
	godoc github.com/samsung-cnct/technical-on-boarding/onboarding/github/issues > $@

docs:  godoc.txt


.PHONY: vet lint test test-cover setup clean docs