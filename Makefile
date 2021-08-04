include env.mk

release: build bump-version tag-kd6rmx

tag-kd6rmx:
	git tag $(VERSION)
	git push origin $(VERSION)

build:
	env GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o build/kd6ctl ./cmd/kd6ctl

install: build
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