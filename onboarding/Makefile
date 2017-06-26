

prepare_workload: ./github/issues/*.go ./cmd/*.go
	go build -v -o $@  github.com/samsung-cnct/technical-on-boarding/onboarding/cmd

build: prepare_workload

test:
	go test -v ./github/issues