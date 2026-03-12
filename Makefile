APP_NAME := snip
CMD_DIR := ./cmd/snip
BIN_DIR := ./bin
BIN_PATH := $(BIN_DIR)/$(APP_NAME)

PREFIX ?= /usr/local
BINDIR ?= $(PREFIX)/bin
COMPLETION_BASH ?= $(PREFIX)/share/bash-completion/completions
COMPLETION_ZSH ?= $(PREFIX)/share/zsh/site-functions
COMPLETION_FISH ?= $(PREFIX)/share/fish/vendor_completions.d

.PHONY: build install uninstall completions install-completions clean

build:
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN_PATH) $(CMD_DIR)

install: build
	@mkdir -p $(BINDIR)
	install -m 0755 $(BIN_PATH) $(BINDIR)/$(APP_NAME)

uninstall:
	@rm -f $(BINDIR)/$(APP_NAME)
	@rm -f $(COMPLETION_BASH)/$(APP_NAME)
	@rm -f $(COMPLETION_ZSH)/_$(APP_NAME)
	@rm -f $(COMPLETION_FISH)/$(APP_NAME).fish

completions: build
	@mkdir -p $(BIN_DIR)/completions
	$(BIN_PATH) completion bash > $(BIN_DIR)/completions/$(APP_NAME).bash
	$(BIN_PATH) completion zsh > $(BIN_DIR)/completions/_$(APP_NAME)
	$(BIN_PATH) completion fish > $(BIN_DIR)/completions/$(APP_NAME).fish

install-completions: completions
	@mkdir -p $(COMPLETION_BASH)
	@mkdir -p $(COMPLETION_ZSH)
	@mkdir -p $(COMPLETION_FISH)
	install -m 0644 $(BIN_DIR)/completions/$(APP_NAME).bash $(COMPLETION_BASH)/$(APP_NAME)
	install -m 0644 $(BIN_DIR)/completions/_$(APP_NAME) $(COMPLETION_ZSH)/_$(APP_NAME)
	install -m 0644 $(BIN_DIR)/completions/$(APP_NAME).fish $(COMPLETION_FISH)/$(APP_NAME).fish

clean:
	@rm -rf $(BIN_DIR)
