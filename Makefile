.PHONY: test vendor

ifndef TESTRUN
TESTRUN=".*"
endif

all:
	@echo "did you mean 'make test' ?"

vendor:
	./hack/vendor.sh

guard-%:
	@ if [ "${${*}}" = "" ]; then \
                echo "Env var '$*' not set"; \
                exit 1; \
        fi

image:
	export TERMINFO=""
	docker build . -t neowaylabs/klb

shell: image
	./hack/run-tty.sh /usr/bin/nash

example-%: image
	./hack/run-tty.sh ./examples/azure/$*/build.sh

example-%-cleanup: image
	./hack/run-tty.sh ./examples/azure/$*/cleanup.sh

libdir=$(NASHPATH)/lib/klb
bindir=$(NASHPATH)/bin
install: guard-NASHPATH
	@rm -rf $(libdir)
	@mkdir -p $(libdir)
	@mkdir -p $(bindir)
	@cp -pr ./aws $(libdir)
	@cp -pr ./azure $(libdir)
	@cp -pr ./tools/azure/getcredentials.sh $(bindir)/azure-credentials.sh
	@cp -pr ./tools/azure/createsp.sh $(bindir)/createsp.sh

timeout=60m
logger=file
parallel=30 #Explore I/O parallelization
cpu=4
gotest=go test ./tests/azure -parallel $(parallel) -timeout $(timeout) -cpu $(cpu)
gotestargs=-args -logger $(logger)

test: image
	./hack/run.sh nash ./azure/vm_test.sh

test-integration: image
	./hack/run.sh $(gotest) -run=$(run) ./... $(gotestargs)

test-examples: image
	./hack/run.sh $(gotest) -tags=examples -run=TestExamples $(gotestargs)

# It is recommended to use this locally. It takes too much time for the CI
test-all: test test-integration test-examples

cleanup: image
	./hack/run-tty.sh ./tools/azure/cleanup.sh

testhost:
	$(gotest) -run=$(run) ./... $(gotestargs)
