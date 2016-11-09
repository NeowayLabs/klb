all:
	@echo "did you mean 'make test' ?"

deps:
	@echo "Downloading jq..."
	sudo wget "https://github.com/stedolan/jq/releases/download/jq-1.5/jq-linux64" -O /usr/bin/jq
	sudo chmod "+x" /usr/bin/jq

depsdev:
	@echo "Getting dependencies for dev"
	go get -u github.com/Azure/azure-sdk-for-go

test:
	cd aws/tests && make test

testazure:
	cd tests/azure && go test ./...
