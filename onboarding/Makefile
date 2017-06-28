
PATH:="$(PATH):$(GOPATH)/bin/"


prepare_workload: ./github/issues/*.go ./cmd/*.go
	go build -v -o $@  github.com/samsung-cnct/technical-on-boarding/onboarding/cmd

build: prepare_workload

test:
	go test -v github.com/samsung-cnct/technical-on-boarding/onboarding/github/issues/
	go vet -v github.com/samsung-cnct/technical-on-boarding/onboarding/github/issues/


lint:
	gofmt -w -s ./github/issues/
	gofmt -w -s ./cmd/
	goimports -w ./github/issues/
	goimports -w ./cmd/
	go get -u github.com/golang/lint/golint
	go get honnef.co/go/tools/cmd/gosimple
	golint ./github/issues/
	golint ./cmd/
	gosimple ./github/issues/
	gosimple ./cmd/
	

