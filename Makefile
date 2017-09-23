# Commands
GO_CMD = `which go`
LINT_CMD = $(GOPATH)/bin/golint
DOCKER_CMD = `which docker`

# Directories
PACKAGE = github.com/minimalchat/daemon
SRC = $(GOPATH)/src/$(PACKAGE)
DIST = $(SRC)/dist

default: lint test coverage clean compile

build: lint test clean compile

run: lint test go

dependencies:
	cat $(SRC)/requirements.txt | xargs -I \\# go get -u github.com/\\#

lint:
	$(LINT_CMD) $(SRC)
	# $(LINT_CMD) $(SRC) $(TEST)

test:
	cd $(SRC)
	$(GO_CMD) test -v ./...
	# $(GOPATH)/bin/goveralls -coverprofile=coverage.out -service=travis-ci -repotoken $(COVERALLS_TOKEN)

coverage:
	cd $(SRC)
	$(GOPATH)/bin/overalls -project=$(PACKAGE) -covermode=count
	$(GOPATH)/bin/goveralls -coverprofile=overalls.coverprofile -service=travis-ci -repotoken $(COVERALLS_TOKEN)

clean:
	rm -rf $(DIST)/mnml-daemon

protob-gen:
	protoc --plugin=protoc-gen-go=$(GOPATH)bin/protoc-gen-go --go_out=Mclient/client.proto=github.com/minimalchat/daemon/client:. client/*.proto
	protoc --plugin=protoc-gen-go=$(GOPATH)bin/protoc-gen-go --go_out=Mclient/client.proto=github.com/minimalchat/daemon/client:. chat/*.proto
	protoc --plugin=protoc-gen-go=$(GOPATH)bin/protoc-gen-go --go_out=Mclient/client.proto=github.com/minimalchat/daemon/client:. operator/*.proto

compile:
	mkdir -p $(DIST)
	cd $(SRC)
	$(GO_CMD) build -o $(DIST)/daemon

docker: compile
	$(DOCKER_CMD) build -t minimalchat/daemon $(SRC)

go:
	cd $(SRC)
	$(GO_CMD) run main.go -cors
