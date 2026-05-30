set shell := ["bash", "-cu"]

pkg := "github.com/hagatasdelus/agysession"
commit := `git rev-parse --short HEAD 2>/dev/null || echo "HEAD"`
build_ldflags := "-s -w -X " + pkg + "/version.Revision=" + commit

# Run tests
default: test

# Run tests with coverage
test:
	go test ./... -coverprofile=coverage.out -covermode=count -count=1

# Build the binary
build:
	go build -ldflags="{{build_ldflags}}" -trimpath -o agysession .

# Run linter
lint:
	golangci-lint run ./...

# Install dev dependencies
depsdev:
	go install github.com/Songmu/gocredits/cmd/gocredits@latest
	go install github.com/k1LoW/gostyle@latest

# Update credits
credits: depsdev
	mise exec -- go mod download
	mise exec -- gocredits -w .
