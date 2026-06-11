PKG := github.com/hagatasdelus/agysession
CURRENT_REVISION := $(shell git rev-parse --short HEAD 2>/dev/null || echo "HEAD")
BUILD_LDFLAGS := -s -w -X $(PKG)/version.Revision=$(CURRENT_REVISION)
GOBIN := $(shell go env GOPATH)/bin
VERSION := $(shell godzil show-version)

.PHONY: default ci test build install lint verify clean credits release crossbuild upload

default: test

ci: $(GOBIN)/godzil $(GOBIN)/gostyle lint test

test:
	go test ./... -coverprofile=coverage.out -covermode=count -count=1

build:
	go build -ldflags="$(BUILD_LDFLAGS)" -trimpath -o agysession .

install:
	go install -ldflags="$(BUILD_LDFLAGS)" -trimpath .

lint:
	golangci-lint run ./...
	go vet -vettool=`which gostyle` -gostyle.config=$(PWD)/.gostyle.yml ./...

verify: lint test

clean:
	rm -rf $(BIN)
	go clean

$(GOBIN)/gostyle:
	go install github.com/k1LoW/gostyle@latest

$(GOBIN)/godzil:
	go install github.com/Songmu/godzil/cmd/godzil@latest

$(GOBIN)/ghr:
	go install github.com/tcnksm/ghr@latest

credits: $(GOBIN)/godzil
	go mod download
	godzil credits -w .

release: $(GOBIN)/godzil
	godzil release

DIST_DIR = dist
crossbuild: $(GOBIN)/godzil
	rm -rf $(DIST_DIR)
	godzil crossbuild -pv=v$(VERSION) -build-ldflags=$(BUILD_LDFLAGS) -os=linux,darwin \
		-static -d=$(DIST_DIR) ./cmd/* 
	cd $(DIST_DIR) && shasum -a 256 $$(find * -type f -maxdepth 0) > SHA256SUMS

upload:
	ghr -body="$$(./godzil changelog --latest -F markdown)" v$(VERSION) $(DIST_DIR)
