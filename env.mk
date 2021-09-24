# current git sha of the repo
GITSHA := $(shell git rev-parse --short HEAD)

# vision-scan version
VERSION := $(GITSHA)
