# current git sha of the repo
GITSHA := $(shell git rev-parse --short HEAD)

# version
VERSION := $(GITSHA)

GOOS    := $(shell go env GOOS)
GOARCH  := $(shell go env GOARCH)

release: bump-version build-all tag

tag:
	git tag $(VERSION)
	git push origin $(VERSION)

build-windows: 
	mkdir -p build/kd6ctl-$(VERSION)-windows-amd64
	env GOOS=windows GOARCH=amd64 go build -o build/kd6ctl-$(VERSION)-windows-amd64/kd6ctl.exe ./cmd/kd6ctl

build-linux:
	mkdir -p build/kd6ctl-$(VERSION)-linux-amd64
	env GOOS=linux GOARCH=amd64 go build -o build/kd6ctl-$(VERSION)-linux-amd64/kd6ctl ./cmd/kd6ctl

build-macos:
	mkdir -p build/kd6ctl-$(VERSION)-macos-amd64
	env GOOS=darwin GOARCH=amd64 go build -o build/kd6ctl-$(VERSION)-macos-amd64/kd6ctl ./cmd/kd6ctl

build-macos-m1:
	mkdir -p build/kd6ctl-$(VERSION)-macos-arm64
	env GOOS=darwin GOARCH=arm64 go build -o build/kd6ctl-$(VERSION)-macos-arm64/kd6ctl ./cmd/kd6ctl

build-all: clean build-windows build-linux build-macos build-macos-m1

.PHONY: build
build:
ifeq ($(GOOS), windows)
	target_name="kd6ctl.exe"
else
	target_name="kd6ctl"
endif
	mkdir -p build/kd6ctl-$(VERSION)-$(GOOS)-$(GOARCH)
	env GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o build/kd6ctl-$(VERSION)-$(GOOS)-$(GOARCH)/$(target_name) ./cmd/kd6ctl

clean:
	rm -rf build/*

install:
	go install ./cmd/kd6ctl 

bump-version:
	$(eval v := $(shell git describe --tags --abbrev=0 2>/dev/null | xargs git rev-parse | xargs git tag --points-at | tail -1 | sed -Ee 's/^v|-.*//'))
ifeq ($(bump), major)
	$(eval VERSION := v$(shell echo $v | awk -F'.' '{printf("%d.0.0", $$1+1, 0, 0)}'))
else ifeq ($(bump), minor)
	$(eval VERSION := v$(shell echo $v | awk -F'.' '{printf("%d.%d.0", $$1, $$2+1, 0)}'))
else ifeq ($(bump), patch)
	$(eval VERSION := v$(shell echo $v | awk -F'.' '{printf("%d.%d.%d", $$1, $$2, $$3+1)}'))
else ifeq ($(bump), )
	$(eval VERSION := v$(shell echo $v | awk -F'.' '{printf("%d.%d.%d", $$1, $$2, $$3+1)}'))
else ifeq ($(bump), which)
	$(eval VERSION := v$(shell echo $v | awk -F'.' '{printf("%d.%d.%d", $$1, $$2, $$3)}'))
else
	$(error $(bump) is not supported. available bump values are: major, minor, patch or leave empty for patch as default)
endif
	@echo bump version: $(VERSION)
