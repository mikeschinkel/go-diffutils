linter := "github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.6.2"

# Every directory holding a go.mod, found rather than listed, so a new module
# needs no edit here. Worktrees under .endless are other checkouts, not modules.
modules := `find . -path ./.endless -prune -o -name go.mod -print | sed 's|/go.mod$||' | sort | tr '\n' ' '`

# List recipes
help:
    @just --list

# Write a go.work over every module, so local builds use the code on disk
work:
    rm -f go.work go.work.sum
    go work init {{modules}}

# Tidy, test, lint and vet
ensure-valid: tidy test lint vet

# Run tests in every module
test:
    #!/usr/bin/env bash
    set -euo pipefail
    for module in {{modules}}; do
        echo "Running tests for $module..."
        (cd "$module" && go test -v -race ./...)
    done

# Run golangci-lint in every module
lint:
    #!/usr/bin/env bash
    set -euo pipefail
    for module in {{modules}}; do
        echo "Running linter for $module..."
        (cd "$module" && go run {{linter}} run ./... --timeout=5m)
    done

# Format code with gofmt
fmt:
    gofmt -s -w .

# Run go vet in every module
vet:
    #!/usr/bin/env bash
    set -euo pipefail
    for module in {{modules}}; do
        echo "Running go vet for $module..."
        (cd "$module" && go vet ./...)
    done

# Run go mod tidy in every module, as CI sees it (without go.work)
tidy:
    #!/usr/bin/env bash
    set -euo pipefail
    for module in {{modules}}; do
        echo "Running go mod tidy for $module..."
        (cd "$module" && GOWORK=off go mod tidy)
    done

# Build every module
build:
    #!/usr/bin/env bash
    set -euo pipefail
    for module in {{modules}}; do
        echo "Building $module..."
        (cd "$module" && go build ./...)
    done

# Clean build artifacts
clean:
    #!/usr/bin/env bash
    set -euo pipefail
    for module in {{modules}}; do
        (cd "$module" && go clean)
    done
    rm -f coverage.txt

# Run all CI checks locally (fmt, vet, lint, test)
ci: fmt vet lint test
    @echo "All CI checks passed!"
