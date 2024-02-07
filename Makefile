PKG     = github.com/ArcticOJ/blizzard/v0
BIN		= arctic
HASH    = $(shell git rev-parse --short HEAD)
DATE    = $(shell date +%s)
TAG     = $(shell git describe --tags --always --abbrev=0 --match="v[0-9]*.[0-9]*.[0-9]*" 2> /dev/null)
VERSION = $(shell echo "${TAG}" | sed 's/^.//')

DEV_FLAGS = -ldflags "-X '${PKG}/build.Hash=${HASH}' -X '${PKG}/build._date=${DATE}'"
REL_FLAGS = -ldflags "-X '${PKG}/build.Version=${VERSION}' -X '${PKG}/build.Hash=${HASH}' -X '${PKG}/build._date=${DATE}' -s -w"

# TODO: test before releasing
release: main.go main_headless.go main_nothing.go main_orca.go main_nothing.go
	go build ${REL_FLAGS} -tags ui,headless,orca -o ${OUT}

dev: main.go main_headless.go main_nothing.go main_orca.go main_nothing.go
	go build ${DEV_FLAGS} -tags headless,orca -o ${OUT}

