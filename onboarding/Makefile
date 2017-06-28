
PATH:="$(PATH):$(GOPATH)/bin/:`pwd`/go/bin/"

all: lint test build


build: prepare_workload


deps: 
	go version
	@echo GOPATH IS ${GOPATH}
	go get github.com/satori/go.uuid
	go get -u -fix github.com/google/go-github/github
	go get gopkg.in/yaml.v2
	go get golang.org/x/oauth2


prepare_workload: ./github/issues/*.go ./cmd/*.go
	go build -v -o $@  github.com/samsung-cnct/technical-on-boarding/onboarding/cmd

test:
	go test -v -cover github.com/samsung-cnct/technical-on-boarding/onboarding/github/issues/
	go vet -v github.com/samsung-cnct/technical-on-boarding/onboarding/github/issues/


lint: deps
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
	

