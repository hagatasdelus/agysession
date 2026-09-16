PKG := github.com/hagatasdelus/agysession
CURRENT_REVISION := $(shell git rev-parse --short HEAD 2>/dev/null || echo "HEAD")
BUILD_LDFLAGS := -s -w -X $(PKG)/version.Revision=$(CURRENT_REVISION)
GOBIN := $(shell go env GOPATH)/bin

.PHONY: default ci test build install lint verify clean credits release-dry-run prerelease_for_tagpr

default: test

ci: $(GOBIN)/gostyle lint test

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
	go clean

$(GOBIN)/gostyle:
	go install github.com/k1LoW/gostyle@latest

$(GOBIN)/gocredits:
	go install github.com/Songmu/gocredits/cmd/gocredits@latest

credits: $(GOBIN)/gocredits
	go mod download
	$(GOBIN)/gocredits -w .
	@if [ -f .github/assets/CREDITS.extra ]; then \
		cat .github/assets/CREDITS.extra >> CREDITS; \
	fi

release-dry-run:
	goreleaser release --snapshot --clean

prerelease_for_tagpr: credits
	git add CHANGELOG.md CREDITS go.mod go.sum
