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
	docker build . -t neowaylabs/klb

shell: image
	docker run -ti neowaylabs/klb /usr/bin/nash

libdir=$(NASHPATH)/lib/klb
bindir=$(NASHPATH)/bin
install: guard-NASHPATH
	rm -rf $(libdir)
	mkdir -p $(libdir)
	mkdir -p $(bindir)
	cp -pr ./aws $(libdir)
	cp -pr ./azure $(libdir)
	cp -pr ./tools/azure/getcredentials.sh $(bindir)/azure-credentials.sh
	cp -pr ./tools/azure/createsp.sh $(bindir)/createsp.sh

timeout=60m
logger=file
parallel=30 #Explore I/O parallelization
gotest=cd tests/azure && go test -parallel $(parallel) -timeout $(timeout) -race
gotestargs=-args -logger $(logger)

test:
	$(gotest) -run=$(run) ./... $(gotestargs)
