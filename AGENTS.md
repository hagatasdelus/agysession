# agysession 

Answer the user in Japanese.

`agysession` is a CLI tool that allows you to list, preview, and resume your Antigravity CLI sessions across all your projects. It provides a fuzzy-finder interface powered by `fzf` to quickly find the session you want, and it faithfully resumes the session in its original working directory.

## Build & Test

Requires Go 1.26+.

```sh
make test       # run tests (run `go test ./... -coverprofile=coverage.out -covermode=count -count=1`)
make build      # build the binary (run `go build -ldflags="$(BUILD_LDFLAGS)" -trimpath -o agysession .`)
make lint       # run linters and vet (run `golangci-lint run ./...` and `go vet -vettool=$(which gostyle) -gostyle.config=$(PWD)/.gostyle.yml ./...`)
make verify     # run formatting, checks
make clean      # clean build artifacts (run `rm -rf $(BIN)` and `go clean`)
make install    # install the binary to $GOPATH/bin (run `go install -ldflags="$(BUILD_LDFLAGS)" -trimpath .`)
```

## Quick Commands

```bash
make build      # build
make test       # tests
make lint       # lint and vet
make check      # run linters and type checks
make verify     # run formatting, checks, and tests
```

## Development Style

- Basically develop using TDD. Proceed in the order of exploration, Red, Green, and Refactoring.
- If there are KPI or coverage targets (such as number of tests or execution time), attempt to achieve them.
- If instructions or requirements are unclear, ask questions to clarify them before taking action.
- Express specifications in both documents and test cases; if there is a discrepancy, prioritize the tests for review.

## Code Design

- Maintain separation of concerns.
- Separate state and logic.
- Emphasize readability and maintainability.
- Maintain a structure that makes internal implementation easy to refactor or reimplement.

## TDD Development

Perform development and modifications while cycling through the Red/Green/Refactor stages of Kent Beck's TDD (tdd skill).

## Read external source code (e.g., from GitHub, npm, or PyPI)

Use the opensrc skill.
