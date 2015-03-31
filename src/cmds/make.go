package cmds

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
)

var ()

func Make() {
	path, _ := os.Getwd()
	biname := filepath.Base(path)
	template := fmt.Sprintf(`
# cross compile OSs
OSS:=darwin linux windows
# files (cert, config, ...) to copy into dist directory
FILES:=
# test cases
TESTS:=
# benchmark cases
BENCHES:=

BUILD_SUFFIX:=.build
INSTALL_SUFFIX:=.install
BUILDS:=$(addsuffix $(BUILD_SUFFIX),$(OSS))
INSTALLS:=$(addsuffix $(INSTALL_SUFFIX),$(OSS))

SYS_OS=$(shell go version | awk '{print $$NF}' | awk -F/ '{print $$1}')
SYS_ARCH=$(shell go version | awk '{print $$NF}' | awk -F/ '{print $$2}')
ifndef $(GOOS)
	GOOS=%s
endif
ifndef $(GOARCH)
	GOARCH=%s
endif

all: $(INSTALLS)

build:
	mkdir -p dist/$(GOOS)_$(GOARCH)
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=0 go build -o dist/$(GOOS)_$(GOARCH)/%s src/*.go

$(BUILDS):
	$(MAKE) build GOOS=$(basename $@) GOARCH=$(GOARCH)

$(INSTALLS): %%$(INSTALL_SUFFIX):%%$(BUILD_SUFFIX)
	for i in $(FILES); do cp -f src/$$i dist/$(basename $@)_$(GOARCH)/. ; done

test:
	for i in $(TESTS); do go test $$i ; done

bench:
	for i in $(BENCHES); do go test -bench . $$i ; done

clean:
	rm -Rf dist pkg

.PHONY: all build test bench clean
.PHONY: $(BUILDS) $(INSTALLS)`, runtime.GOOS, runtime.GOARCH, biname, biname)
	ioutil.WriteFile("Makefile", []byte(template), 0644)
}
