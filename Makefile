
# cross compile OSs
OSS:=darwin linux windows
# files to copy into dist/{OSS}/
FILES:=
# test cases
TESTS:=
# benchmark cases
BENCHES:=

BUILD_SUFFIX:=.build
INSTALL_SUFFIX:=.install
BUILDS:=$(addsuffix $(BUILD_SUFFIX),$(OSS))
INSTALLS:=$(addsuffix $(INSTALL_SUFFIX),$(OSS))

GIT_HASH=$(shell git rev-parse --short HEAD)
SYS_OS=$(shell go version | awk '{print $$NF}' | awk -F/ '{print $$1}')
SYS_ARCH=$(shell go version | awk '{print $$NF}' | awk -F/ '{print $$2}')
ifndef $(GOOS)
	GOOS=darwin
endif
ifndef $(GOARCH)
	GOARCH=amd64
endif

all: $(INSTALLS)

build:
	mkdir -p dist/$(GOOS)_$(GOARCH)
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=0 go build -o dist/$(GOOS)_$(GOARCH)/gg src/*.go

$(BUILDS):
	$(MAKE) build GOOS=$(basename $@) GOARCH=$(GOARCH)

$(INSTALLS): %$(INSTALL_SUFFIX):%$(BUILD_SUFFIX)
	for i in $(FILES); do cp -f src/$$i dist/$(basename $@)_$(GOARCH)/. ; done
	tar -czf dist/$(basename $@)_$(GOARCH)_$(GIT_HASH).tar.gz dist/$(basename $@)_$(GOARCH)/*

test:
	for i in $(TESTS); do go test $$i ; done

bench:
	for i in $(BENCHES); do go test -bench . $$i ; done

clean:
	rm -Rf dist

.PHONY: all build test bench clean
.PHONY: $(BUILDS) $(INSTALLS)
