MODULE   = $(shell env GO111MODULE=on $(GO) list -m)
DATE    ?= $(shell date +%FT%T%z)
VERSION ?= $(shell git describe --tags --always --dirty --match=v* 2> /dev/null || \
		   cat $(CURDIR)/.version 2> /dev/null || echo v0)
PKGS     = $(or $(PKG),$(shell env GO111MODULE=on $(GO) list ./...))

GO      = go
GOLINT  = golangci-lint

TARGET_WINDOWS_EXTENSION = "exe"

V = 0
Q = $(if $(filter 1,$V),,@)
M = $(shell printf "\033[34;1mâ–¶\033[0m")

export GO111MODULE=on

## Add support for release tags once we're a bit further along
all: clean lint fmt test| $(BIN) ; $(info $(M) b) @ 

lint: |  ; $(info $(M) running golangci-lint...) @ ## Run golint on all packages
	$Q $(GOLINT) run ./...

fmt: |  ; $(info $(M) running go fmt) @ ## Run go fmt on all packages
	$Q $(GO) fmt $(PKGS)

test: |  ; $(info $(M) running go test) @ ## Run go test on all packages
	$Q $(GO) test $(PKGS)

# Misc
clean: ; $(info $(M) cleaning...)	@ ## Cleanup everything
	@rm -rf $(BIN)

help:
	@awk '/^[a-zA-Z\-\_0-9]+:/ {                    \
	  nb = sub( /^## /, "", helpMsg );              \
	  if(nb == 0) {                                 \
		helpMsg = $$0;                              \
		nb = sub( /^[^:]*:.* ## /, "", helpMsg );   \
	  }                                             \
	  if (nb)                                       \
		print  $$1 "\t" helpMsg;                    \
	}                                               \
	{ helpMsg = $$0 }'                              \
	$(MAKEFILE_LIST) | column -ts $$'\t' |          \
	grep '^[^ ]*'

version:
	@echo $(VERSION)

