all:
	@echo "did you mean 'make test' ?"

deps: aws-deps azure-deps jq-dep

aws-deps:
	sudo pip install awscli

azure-deps: jq-dep
	sudo npm install -g azure-cli

jq-dep:
	@echo "Downloading jq..."
	sudo wget "https://github.com/stedolan/jq/releases/download/jq-1.5/jq-linux64" -O /usr/bin/jq
	sudo chmod "+x" /usr/bin/jq
depsdev:
	@echo "Getting dependencies for dev"
	go get -d ./tests/...

testazure: depsdev
	cd tests/azure && go test ./...
