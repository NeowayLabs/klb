.PHONY: deps aws-deps azure-deps testazure test

ifndef TESTRUN
TESTRUN=".*"
endif

ifndef GOPATH
$(error $$GOPATH is not set)
endif

all:
	@echo "did you mean 'make test' ?"

deps: aws-deps azure-deps jq-dep

aws-deps:
	pip install --user awscli

azure-deps: jq-dep
	npm install -g azure-cli

jq-dep: $(GOPATH)/bin/jq

$(GOPATH)/bin/jq:
	@echo "Downloading jq..."
	mkdir -p $(GOPATH)/bin
	wget "https://github.com/stedolan/jq/releases/download/jq-1.5/jq-linux64" -O $(GOPATH)/bin/jq
	chmod "+x" $(GOPATH)/bin/jq

depsdev:
	@echo "updating dev dependencies"
	go get -u -d ./tests/...

timeout=10m
logger=file
parallel=30 #Explore I/O parallelization
gotest=cd tests/azure && go test -parallel $(parallel) -timeout $(timeout) -race
gotestargs=-args -logger $(logger)

testall: depsdev
	$(gotest) ./... $(gotestargs)

test: depsdev
	$(gotest) -run=$(run) ./... $(gotestargs)
