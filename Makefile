ifeq ($(origin VERSION), undefined)
	VERSION != git rev-parse --short HEAD
endif
HOST_GOOS=$(shell go env GOOS)
HOST_GOARCH=$(shell go env GOARCH)
REPOPATH = github.com/JayH5/proxy

VERBOSE_1 := -v
VERBOSE_2 := -v -x

build: vendor
	$(BUILD_ENV_FLAGS) go build $(VERBOSE_$(V)) -ldflags "-X $(REPOPATH).Version=$(VERSION)"

test: tools/glide
	go test -race $(VERBOSE_$(V)) $(shell ./tools/glide novendor)

vet: tools/glide
	go vet $(VERBOSE_$(V)) $(shell ./tools/glide novendor)

fmt: tools/glide
	go fmt $(VERBOSE_$(V)) $(shell ./tools/glide novendor)

clean:
	# Nothing to do for now

cleanall: clean
	rm -rf bin tools vendor

release:
	mkdir -p release
	goxc -d ./release -tasks-=go-vet,go-test -os="linux darwin" -pv=$(VERSION) -build-ldflags="-X $(REPOPATH).Version=$(VERSION)" -resources-include="README.md,LICENSE" -main-dirs-exclude="vendor"

vendor: tools/glide
	./tools/glide install

tools/glide:
	@echo "Downloading glide"
	mkdir -p tools
	curl -L https://github.com/Masterminds/glide/releases/download/v0.11.0/glide-v0.11.0-$(HOST_GOOS)-$(HOST_GOARCH).tar.gz | tar -xz -C tools
	mv tools/$(HOST_GOOS)-$(HOST_GOARCH)/glide tools/glide
	rm -r tools/$(HOST_GOOS)-$(HOST_GOARCH)

help:
	@echo "Influential make variables"
	@echo "  V                 - Build verbosity {0,1,2}."
	@echo "  BUILD_ENV_FLAGS   - Environment added to 'go build'."

.PHONY: build test vet fmt clean cleanall help
