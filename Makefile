.PHONY: test vendor

version=latest

all:
	@echo "did you mean 'make test' ?"

vendor:
	./hack/vendor.sh

guard-%:
	@ if [ "${${*}}" = "" ]; then \
                echo "Env var '$*' not set"; \
                exit 1; \
        fi

publish:image
	docker push neowaylabs/klb:$(version)

image:
	export TERMINFO=""
	docker build . -t neowaylabs/klb:$(version)

credentials: image guard-sh guard-subscription guard-service-principal guard-service-secret
	docker run -ti --rm -v `pwd`:/credentials -w /credentials neowaylabs/klb:$(version) \
		/credentials/tools/azure/getcredentials.sh \
		$(sh) "$(subscription)" "$(service-principal)" "$(service-secret)"

createsp: image guard-subscription-id guard-service-principal guard-service-secret
	docker run -ti --rm -v `pwd`:/createsp -w /createsp neowaylabs/klb:$(version) \
		/createsp/tools/azure/createsp.sh \
		"$(subscription-id)" "$(service-principal)" "$(service-secret)"

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

integration_timeout=50m
examples_timeout=90m
all_timeout=90m
logger=file
parallel=60 #Explore I/O parallelization
cpu=4
gotest=go test -v ./tests/azure -parallel $(parallel) -cpu $(cpu)
gotestargs=-args -logger $(logger)

test: image
	./hack/run.sh nash ./azure/vm_test.sh

test-integration: image
	./hack/run.sh $(gotest) -timeout $(integration_timeout) -run=$(run) ./... $(gotestargs)

test-examples: image
	./hack/run.sh $(gotest) -timeout $(examples_timeout) -tags=examples -run=TestExamples $(gotestargs)

# It is recommended to use this locally. It takes too much time for the CI
test-all: test
	./hack/run.sh $(gotest) -timeout $(all_timeout) -tags=examples $(gotestargs)

cleanup: image
	./hack/run-tty.sh ./tools/azure/cleanup.sh

testhost:
	$(gotest) -run=$(run) ./... $(gotestargs)
