# version info
VERSION   ?= `cat VERSION`
BUILD     ?= `git rev-parse HEAD`
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD}"

# paths
PATH:="$(PATH):$(GOPATH)/bin/:`pwd`/go/bin/"
APP_PATH:="./app/"
CMD_PATH:="./cmd/"

# naming
PROJECT_NAME:=technical-on-boarding
APP_NAME:="$(PROJECT_NAME)-app"
CMD_NAME:="$(PROJECT_NAME)-cmd"

PKG_NAME:="$(PROJECT_NAME)-v$(VERSION)"
PKG_PATH:="github.com/samsung-cnct/technical-on-boarding/onboarding"

all: vet lint test build

build: $(PKG_NAME)

# TODO: use glide to populate vendored dependencies

setup: 
	@go version
	@echo GOPATH IS ${GOPATH}
	@echo app.version is $(VERSION)+$(BUILD)
	go get github.com/satori/go.uuid
	go get -u -fix github.com/google/go-github/github
	go get gopkg.in/yaml.v2
	go get golang.org/x/oauth2
	go get github.com/revel/cmd/revel
	go get github.com/revel/revel
	go get github.com/revel/cron
	go get github.com/masterminds/semver

$(PKG_NAME): $(APP_NAME) $(CMD_NAME)

$(APP_NAME): #setup $(shell find $(APP_PATH) -name '*.go')
	go build -v $(LDFLAGS) $(APP_PATH)

$(CMD_NAME): setup $(shell find $(CMD_PATH) -name '*.go')
	go build -v $(LDFLAGS) -o $@ $(CMD_PATH)

test: setup vet lint
	go test -race -v $(APP_PATH)jobs/github
	
coverage.html: $(shell find $(APP_PATH) $(CMD_PATH) -name '*.go')
	go test -covermode=count -coverprofile=coverage.prof $(APP_PATH)
	go tool cover -html=coverage.prof -o $@

test-cover: coverage.html

lint: setup
	go get -u github.com/golang/lint/golint
	go get -u golang.org/x/tools/cmd/goimports
	go get -u honnef.co/go/tools/cmd/gosimple
	gofmt -w -s ./app/jobs/github/
	gofmt -w -s ./cmd/
	$(GOPATH)/bin/goimports -w ./app/jobs/github/
	$(GOPATH)/bin/goimports -w ./cmd/
	$(GOPATH)/bin/golint ./app/jobs/github/
	$(GOPATH)/bin/golint ./cmd/
	$(GOPATH)/bin/gosimple ./app/jobs/github/
	$(GOPATH)/bin/gosimple ./cmd/

vet: 
	go vet -v -printf=false $(APP_PATH)

clean:
	-rm -vf ./coverage.* ./$(CMD_NAME) ./$(APP_NAME) ./$(PKG_NAME)

godoc.txt: $(shell find ./ -name '*.go')
	godoc $(APP_PATH) > $@

docs:  godoc.txt

.PHONY: vet lint test test-cover setup clean docs
