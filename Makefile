# Commands
GO_CMD = go
LINT_CMD = $(GOPATH)/bin/golint

# Directories
PACKAGE = github.com/minimalchat/mnml-daemon
SRC = $(GOPATH)/src/$(PACKAGE)
DIST = $(GOPATH)/bin

default: lint test coverage clean compile

build: lint test clean compile

run: lint test go

lint:
	$(LINT_CMD) $(SRC)
	# $(LINT_CMD) $(SRC) $(TEST)

test:
	cd $(SRC)
	$(GO_CMD) test -v ./...
	# $(GOPATH)/bin/goveralls -coverprofile=coverage.out -service=travis-ci -repotoken $(COVERALLS_TOKEN)

coverage:
	cd $(SRC)
	$(DIST)/overalls -project=$(PACKAGE) -covermode=count
	$(GOPATH)/bin/goveralls -coverprofile=overalls.coverprofile -service=travis-ci -repotoken $(COVERALLS_TOKEN)

clean:
	rm -rf $(DIST)/mnml-daemon

compile:
	mkdir -p $(DIST)
	cd $(SRC)
	$(GO_CMD) build -o $(DIST)/mnml-daemon

go:
	cd $(SRC)
	$(GO_CMD) run main.go
