GOCMD=go
GOBUILD=$(GOCMD) build
GOGET=$(GOCMD) get
CURRENT := $(shell pwd)
OUTPUTFILENAME := "main"
BUILDDIR=./build
PKGDIR=$(BUILDDIR)/pkg
GOXOS := "linux"
GOXARCH := "amd64"
GOXOUTPUT := "$(PKGDIR)/$(OUTPUTFILENAME)_{{.OS}}_{{.Arch}}/$(OUTPUTFILENAME)"

.PHONY: setup
## Install dependencies
setup:
	$(GOGET) github.com/mitchellh/gox
	$(GOGET) -d -t ./...

.PHONY: cross-build
## Cross build binaries
cross-build:
	rm -rf $(PKGDIR)
	gox -os=$(GOXOS) -arch=$(GOXARCH) -output=$(GOXOUTPUT)