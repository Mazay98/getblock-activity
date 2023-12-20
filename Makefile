# suppress output, run `make XXX V=` to be verbose
V := @

# Common
NAME = get-block-activity
VCS = gitlab.com
VERSION := $(shell git describe --always --tags)
CURRENT_TIME := $(shell TZ="Europe/Moscow" date +"%d-%m-%y %T")

# Build
OUT_DIR = ./bin
MAIN_PKG = ./cmd/${NAME}
ACTION ?= build
GC_FLAGS = -gcflags 'all=-N -l'
LD_FLAGS = -ldflags "-s -v -w -X 'main.version=${VERSION}' -X 'main.buildTime=${CURRENT_TIME}'"
BUILD_CMD = CGO_ENABLED=1 go build -o ${OUT_DIR}/${NAME} ${LD_FLAGS} ${MAIN_PKG}
DEBUG_CMD = CGO_ENABLED=1 go build -o ${OUT_DIR}/${NAME} ${GC_FLAGS} ${MAIN_PKG}

# Other
.DEFAULT_GOAL = build

.PHONY: build
build:
	@echo BUILDING PRODUCTION $(NAME)
	$(V)${BUILD_CMD}
	@echo DONE

.PHONY: build-debug
build-debug:
	@echo BUILDING DEBUG $(NAME)
	$(V)${DEBUG_CMD}
	@echo DONE

.PHONY: clean
clean:
	@echo "Removing $(OUT_DIR)"
	$(V)rm -rf $(OUT_DIR)

.PHONY: vendor
vendor:
	$(V)GOPRIVATE=${VCS}/* go mod tidy -compat=1.21.0
	$(V)GOPRIVATE=${VCS}/* go mod vendor
	$(V)git add vendor go.mod go.sum
