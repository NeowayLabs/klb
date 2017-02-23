.PHONY: deps aws-deps azure-deps testazure test vendor

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
	sudo npm install -g azure-cli

jq-dep: $(GOPATH)/bin/jq

$(GOPATH)/bin/jq:
	@echo "Downloading jq..."
	mkdir -p $(GOPATH)/bin
	wget "https://github.com/stedolan/jq/releases/download/jq-1.5/jq-linux64" -O $(GOPATH)/bin/jq
	chmod "+x" $(GOPATH)/bin/jq

vendor:
	./hack/vendor.sh


guard-%:
	@ if [ "${${*}}" = "" ]; then \
                echo "Env var '$*' not set"; \
                exit 1; \
        fi

libdir=$(NASHPATH)/lib/klb
bindir=$(NASHPATH)/bin
install: guard-NASHPATH
	rm -rf $(libdir)
	mkdir -p $(libdir)
	mkdir -p $(bindir)
	cp -pr ./aws $(libdir)
	cp -pr ./azure $(libdir)
	cp -pr ./tools/azure/getcredentials.sh $(bindir)/azure-credentials.sh

timeout=30m
logger=file
parallel=30 #Explore I/O parallelization
#FIXME ADD RACE DETECTOR AGAIN BEFORE i4k GETS PISSED OFF
#gotest=cd tests/azure && go test -parallel $(parallel) -timeout $(timeout) -race
gotest=cd tests/azure && go test -parallel $(parallel) -timeout $(timeout)
gotestargs=-args -logger $(logger)

test:
	$(gotest) -run=$(run) ./... $(gotestargs)
