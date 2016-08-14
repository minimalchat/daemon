# Commands
GOCMD = /usr/bin/go

LINT_CMD = $(GOPATH)/bin/golint
TEST_CMD = $(GOCMD) test
COMPILE_CMD = $(GOCMD) build
RUN_CMD = $(GOCMD) run

# Directories
PACKAGE = github.com/mihok/letschat-daemon
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
	$(TEST_CMD)

clean:
	rm -rf $(DIST)/letschat-daemon

compile:
	mkdir -p $(DIST)
	cd $(SRC)
	$(COMPILE_CMD) -o $(DIST)/letschat

go:
	cd $(SRC)
	$(RUN_CMD) main.go