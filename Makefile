# qk Makefile
#
# Builds the qk command resolver, runs tests, and installs shell wrappers that
# resolve shorthand invocations into shell commands and execute them.
#
# Variables:
#   BINARY        Name of the compiled binary (default: qk)
#   GO            Go toolchain command (default: go)
#   QK_ROOT       Absolute path to this repository
#   QK_ZSHRC      zsh config file updated by install/uninstall (default: ~/.zshrc)
#   QK_BASHRC     bash config file updated by install/uninstall (default: ~/.bashrc)
#   QK_ZSH_WRAPPER  Path to the zsh function sourced by ~/.zshrc
#   QK_BASH_WRAPPER Path to the bash function sourced by ~/.bashrc
#
# Targets:
#   all         Alias for build
#   build       Untrack work files, run tests, and compile the qk binary
#   test        Run all Go tests
#   install     Build qk and append shell wrapper source lines to ~/.zshrc and ~/.bashrc
#   uninstall   Remove shell wrapper source lines from ~/.zshrc and ~/.bashrc
#   clean       Delete the built binary
#   diagrams    Render Mermaid sources in docs/diagrams/ to SVG
#   help        Print target usage
#   untrack-work  Mark local work command files as assume-unchanged in git

BINARY := qk
GO := go
SHELL := /bin/bash

QK_ROOT := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))
QK_ZSHRC := $(HOME)/.zshrc
QK_BASHRC := $(HOME)/.bashrc
QK_ZSH_WRAPPER := $(QK_ROOT)/install/qk.zsh
QK_BASH_WRAPPER := $(QK_ROOT)/install/qk.bash
QK_ZSH_MARKER := qk/install/qk.zsh
QK_BASH_MARKER := qk/install/qk.bash

.PHONY: all build test clean install uninstall help untrack-work install-rc uninstall-rc diagrams

MMDC := $(QK_ROOT)/node_modules/.bin/mmdc
DIAGRAMS := docs/diagrams
DIAGRAM_SOURCES := $(wildcard $(DIAGRAMS)/*.mmd)
DIAGRAM_SVGS := $(DIAGRAM_SOURCES:.mmd=.svg)

all: build

# help prints available make targets and post-install shell usage.
help:
	@echo "Usage:"
	@echo "  make build      Build $(BINARY) (runs tests first)"
	@echo "  make test       Run tests"
	@echo "  make install    Build and add the qk shell function to ~/.zshrc and ~/.bashrc"
	@echo "  make uninstall  Remove the qk shell function from ~/.zshrc and ~/.bashrc"
	@echo "  make clean      Remove built binary"
	@echo "  make diagrams   Render docs/diagrams/*.mmd to SVG via local mermaid-cli"
	@echo ""
	@echo "After install, reload your shell and run qk from anywhere:"
	@echo "  source ~/.zshrc   # zsh"
	@echo "  source ~/.bashrc  # bash"
	@echo "  qk reload"

# untrack-work prevents accidentally committing local work-specific commands.
untrack-work:
	@work_path="commands/work/work.go"; \
	work_test_path="commands/work/work_test.go"; \
	for path in $$work_path $$work_test_path; do \
		if [ "$$(git ls-files -v $$path 2>/dev/null)" = "H $$path" ]; then \
			echo "$$path should not be tracked by git. Running git update-index --assume-unchanged $$path" >&2; \
			git update-index --assume-unchanged $$path; \
		fi; \
	done

# build produces the qk binary after tests pass.
build: untrack-work test
	$(GO) build -o $(BINARY) .

# test runs the full Go test suite.
test:
	$(GO) test ./...

# clean removes the compiled binary from the repository root.
clean:
	rm -f $(BINARY)

# diagrams renders Mermaid source files to SVG using the local mermaid-cli install.
diagrams: node_modules $(DIAGRAM_SVGS)

node_modules: package.json package-lock.json
	npm install

$(DIAGRAMS)/%.svg: $(DIAGRAMS)/%.mmd node_modules
	$(MMDC) -i $< -o $@

# install-rc adds the qk wrapper to a single shell rc file.
# Usage: make install-rc RC=~/.bashrc WRAPPER=.../qk.bash MARKER=qk/install/qk.bash
install-rc:
	@if [ -z "$(RC)" ] || [ -z "$(WRAPPER)" ] || [ -z "$(MARKER)" ]; then \
		echo "install-rc requires RC, WRAPPER, and MARKER" >&2; \
		exit 1; \
	fi
	@if [ ! -f "$(RC)" ]; then \
		touch "$(RC)"; \
	fi
	@if grep -q "$(MARKER)" "$(RC)"; then \
		echo "qk is already installed in $(RC)"; \
	else \
		printf '\n# qk command wrapper\nsource %s\n' "$(WRAPPER)" >> "$(RC)"; \
		echo "Added qk to $(RC)"; \
	fi

# uninstall-rc removes the qk wrapper from a single shell rc file.
uninstall-rc:
	@if [ -z "$(RC)" ] || [ -z "$(MARKER)" ]; then \
		echo "uninstall-rc requires RC and MARKER" >&2; \
		exit 1; \
	fi
	@if [ ! -f "$(RC)" ]; then \
		echo "Nothing to uninstall in $(RC)"; \
		exit 0; \
	fi
	@if grep -q "$(MARKER)" "$(RC)"; then \
		sed -i '/# qk command wrapper/d;/source .*qk\/install\/qk\.\(zsh\|bash\)/d' "$(RC)"; \
		echo "Removed qk from $(RC)"; \
	else \
		echo "qk is not installed in $(RC)"; \
	fi

# install builds qk and registers shell wrappers in ~/.zshrc and ~/.bashrc.
install: build
	@$(MAKE) --no-print-directory install-rc RC="$(QK_ZSHRC)" WRAPPER="$(QK_ZSH_WRAPPER)" MARKER="$(QK_ZSH_MARKER)"
	@$(MAKE) --no-print-directory install-rc RC="$(QK_BASHRC)" WRAPPER="$(QK_BASH_WRAPPER)" MARKER="$(QK_BASH_MARKER)"
	@echo "Run: source ~/.zshrc  # zsh"
	@echo "Run: source ~/.bashrc # bash"

# uninstall removes shell wrappers from ~/.zshrc and ~/.bashrc.
uninstall:
	@$(MAKE) --no-print-directory uninstall-rc RC="$(QK_ZSHRC)" MARKER="$(QK_ZSH_MARKER)"
	@$(MAKE) --no-print-directory uninstall-rc RC="$(QK_BASHRC)" MARKER="$(QK_BASH_MARKER)"
