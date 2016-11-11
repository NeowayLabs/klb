all:
	@echo "did you mean 'make test' ?"

deps: aws-deps azure-deps jq-dep

aws-deps:
	pip install --user awscli

azure-deps: jq-dep
	npm install -g azure-cli

jq-dep:
	@echo "Downloading jq..."
	mkdir -p $(GOPATH)/bin
	wget "https://github.com/stedolan/jq/releases/download/jq-1.5/jq-linux64" -O $(GOPATH)/bin/jq
	chmod "+x" $(GOPATH)/bin/jq
depsdev:
	@echo "Getting dependencies for dev"
	go get -d ./tests/...

testazure: depsdev
	cd tests/azure && go test ./...
