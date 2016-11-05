# Commands
GO_CMD = go
LINT_CMD = $(GOPATH)/bin/golint

# Directories
PACKAGE = github.com/minimalchat/mnml-daemon
SRC = $(GOPATH)/src/$(PACKAGE)
DIST = $(GOPATH)/bin

.PHONY: build lint

build: lint test clean compile

run: lint test go

lint:
	$(LINT_CMD) $(SRC)
	# $(LINT_CMD) $(SRC) $(TEST)

test:
	cd $(SRC)
	$(GO_CMD) test

clean:
	rm -rf $(DIST)/mnml-daemon

compile:
	mkdir -p $(DIST)
	cd $(SRC)
	$(GO_CMD) build -o $(DIST)/letschat

go:
	cd $(SRC)
	$(GO_CMD) run main.go
