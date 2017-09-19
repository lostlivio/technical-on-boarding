# import configuration variables
include Makefile.env
export $(shell sed -E 's/\??=.*//' Makefile.env)

LDFLAGS=-ldflags "-X main.Version=${PROJECT_VERSION} -X main.Build=${PROJECT_BUILD}"
PATH:="$(PATH):$(GOPATH)/bin/:`pwd`/go/bin/"

all: vet lint test build

build: $(PKG_NAME)

# TODO: use glide to populate vendored dependencies

setup: 
	@go version
	@echo GOPATH IS ${GOPATH}
	@echo app.version is $(PROJECT_VERSION)+$(PROJECT_BUILD)
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
	gofmt -w -s $(APP_PATH)
	gofmt -w -s $(CMD_PATH)
	$(GOPATH)/bin/goimports -w $(APP_PATH)
	$(GOPATH)/bin/goimports -w $(CMD_PATH)
	$(GOPATH)/bin/golint $(APP_PATH) $(APP_PATH)controllers $(APP_PATH)jobs $(APP_PATH)jobs/github
	$(GOPATH)/bin/golint $(CMD_PATH)
	$(GOPATH)/bin/gosimple $(APP_PATH)
	$(GOPATH)/bin/gosimple $(CMD_PATH)

vet: 
	go vet -v -printf=false $(APP_PATH)

clean:
	-rm -vf ./coverage.* ./$(CMD_NAME) ./$(APP_NAME)
	-rm -rf ./test-results/

godoc.txt: $(shell find ./ -name '*.go')
	godoc $(APP_PATH) > $@

docs:  godoc.txt

.PHONY: vet lint test test-cover setup clean docs
