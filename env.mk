# host os, arch and version
GOOS    := $(shell go env GOOS)
GOARCH  := $(shell go env GOARCH)

# current git sha of the repo
GITSHA := $(shell git rev-parse --short HEAD)

# vision-scan version
VERSION := $(GITSHA)