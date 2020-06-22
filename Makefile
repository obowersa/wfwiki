MODULE   = $(shell env GO111MODULE=on $(GO) list -m)
DATE    ?= $(shell date +%FT%T%z)
VERSION ?= $(shell git describe --tags --always --dirty --match=v* 2> /dev/null || \
		   cat $(CURDIR)/.version 2> /dev/null || echo v0)

PKGS     = $(or $(PKG),$(shell env GO111MODULE=on $(GO) list ./...))
BIN      = $(CURDIR)/bin

GO      = go
GOLINT  = golangci-lint
BUILD   = $(CURDIR)/cmd/wfwiki/main.go

TARGET_WINDOWS_EXTENSION = "exe"
TARGET_WINDOWS_AMD64 = GOOS=windows GOARCH=amd64
TARGET_LINUX_AMD64 = GOOS=linux GOARCH=amd64

V = 0
Q = $(if $(filter 1,$V),,@)
M = $(shell printf "\033[34;1mâ–¶\033[0m")

export GO111MODULE=on

## Add support for release tags once we're a bit further along
all: clean lint fmt test build_windows_amd64 build_linux_amd64 | $(BIN) ; $(info $(M) b) @

## Go build
build_windows_amd64: ; $(info $(M) building windows_amd64 executable...) @ ## Build for windows amd64
	$(Q) $(TARGET_WINDOWS_AMD64) $(GO) build \
		-ldflags '-X main.Version=$(VERSION) -X main.BuildDate=$(DATE)' \
		-o $(BIN)/$(basename $(MODULE)).$(TARGET_WINDOWS_EXTENSION) $(BUILD)

build_linux_amd64: ; $(info $(M) building linux_amd64 executable...) @ ## Build for linux amd64
	$(Q) $(TARGET_LINUX_AMD64) $(GO) build \
		-ldflags '-X main.Version=$(VERSION) -X main.BuildDate=$(DATE)' \
		-o $(BIN)/$(basename $(MODULE)) $(BUILD)

##Linting
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

