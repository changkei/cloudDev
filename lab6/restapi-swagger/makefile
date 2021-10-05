#-------------------------
# Build artifacts
#-------------------------
GO	:= go

pkgs	= $(shell GOFLAGS=-mod=vendor $(GO) list ./... | grep -vE -e /vendor/ -e /pkg/swagger/)
pkgDirs	=	$(shell GOFLAGS=-mod=vendor $(GO) list -f {{.Dir}} ./... | grep -vE -e /vendor/ -e /pkg/swagger/)
DIR_OUT:=/tmp


build:
	$(GO) build -o bin/http-go-server internal/main.go
:

#-------------------------
# Target: swagger.validate
#-------------------------
.PHONY:	swagger.validate

swagger.validate:
	swagger validate pkg/swagger/swagger.yml

#-------------------------
# Code generation
#-------------------------
.PHONY:	generate

## Generate go code
generate:
	@echo "==> generating go code"
	GOFLAGS=-mod=vendor $(GO) generate $(pkgs)

