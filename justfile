# GoForGo Development Justfile
# See https://github.com/casey/just for more information

set shell := ["bash", "-c"]

# Default recipe - show available commands
default:
    @just --list

# Build the GoForGo CLI binary
build:
    @echo "🔨 Building GoForGo CLI..."
    mkdir -p bin
    go build -ldflags="-X 'github.com/stonecharioteer/goforgo/internal/cli.version={{version}}' -X 'github.com/stonecharioteer/goforgo/internal/cli.commit={{commit}}' -X 'github.com/stonecharioteer/goforgo/internal/cli.date={{date}}'" -o bin/goforgo ./cmd/goforgo
    @echo "✅ Binary built: bin/goforgo"

# Build with race detection (for development)
build-race:
    @echo "🔨 Building GoForGo CLI with race detection..."
    mkdir -p bin
    go build -race -ldflags="-X 'github.com/stonecharioteer/goforgo/internal/cli.version={{version}}' -X 'github.com/stonecharioteer/goforgo/internal/cli.commit={{commit}}' -X 'github.com/stonecharioteer/goforgo/internal/cli.date={{date}}'" -o bin/goforgo-race ./cmd/goforgo
    @echo "✅ Race-enabled binary built: bin/goforgo-race"

# Run all tests
test:
    @echo "🧪 Running tests..."
    go test -v ./...

# Run tests with coverage
test-coverage:
    @echo "🧪 Running tests with coverage..."
    go test -v -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out -o coverage.html
    @echo "📊 Coverage report generated: coverage.html"

# Run benchmarks
bench:
    @echo "⚡ Running benchmarks..."
    go test -bench=. -benchmem ./...

# Lint the code
lint:
    @echo "🔍 Linting code..."
    golangci-lint run ./...

# Format the code
fmt:
    @echo "📝 Formatting code..."
    go fmt ./...
    goimports -w .

# Tidy dependencies
tidy:
    @echo "🧹 Tidying dependencies..."
    go mod tidy

# Clean build artifacts
clean:
    @echo "🧽 Cleaning build artifacts..."
    rm -rf bin/
    rm -f coverage.out coverage.html
    rm -rf test-init/

# Install development dependencies
install-deps:
    @echo "📦 Installing development dependencies..."
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    go install golang.org/x/tools/cmd/goimports@latest

# Development build (fast, includes dev-only commands like solve)
dev-build:
    @echo "🚀 Building for development..."
    mkdir -p bin
    go build -tags dev -o bin/goforgo-dev ./cmd/goforgo
    @echo "✅ Development binary built: bin/goforgo-dev"

# Test the CLI in a temporary directory
test-cli: build
    @echo "🧪 Testing CLI functionality..."
    rm -rf test-init/
    mkdir -p test-init
    cd test-init && ../bin/goforgo init
    @echo "✅ CLI test completed"

# Run goforgo in development mode with a test setup
dev-run: dev-build test-cli
    @echo "🎮 Starting GoForGo in development mode..."
    cd test-init && ../bin/goforgo-dev

# Build release binaries for multiple platforms
build-release:
    @echo "🏗️  Building release binaries..."
    mkdir -p dist
    # Linux AMD64
    GOOS=linux GOARCH=amd64 go build -ldflags="-s -w -X 'github.com/stonecharioteer/goforgo/internal/cli.version={{version}}' -X 'github.com/stonecharioteer/goforgo/internal/cli.commit={{commit}}' -X 'github.com/stonecharioteer/goforgo/internal/cli.date={{date}}'" -o dist/goforgo-linux-amd64 ./cmd/goforgo
    # Linux ARM64
    GOOS=linux GOARCH=arm64 go build -ldflags="-s -w -X 'github.com/stonecharioteer/goforgo/internal/cli.version={{version}}' -X 'github.com/stonecharioteer/goforgo/internal/cli.commit={{commit}}' -X 'github.com/stonecharioteer/goforgo/internal/cli.date={{date}}'" -o dist/goforgo-linux-arm64 ./cmd/goforgo
    # macOS AMD64
    GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w -X 'github.com/stonecharioteer/goforgo/internal/cli.version={{version}}' -X 'github.com/stonecharioteer/goforgo/internal/cli.commit={{commit}}' -X 'github.com/stonecharioteer/goforgo/internal/cli.date={{date}}'" -o dist/goforgo-darwin-amd64 ./cmd/goforgo
    # macOS ARM64 (Apple Silicon)
    GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w -X 'github.com/stonecharioteer/goforgo/internal/cli.version={{version}}' -X 'github.com/stonecharioteer/goforgo/internal/cli.commit={{commit}}' -X 'github.com/stonecharioteer/goforgo/internal/cli.date={{date}}'" -o dist/goforgo-darwin-arm64 ./cmd/goforgo
    # Windows AMD64
    GOOS=windows GOARCH=amd64 go build -ldflags="-s -w -X 'github.com/stonecharioteer/goforgo/internal/cli.version={{version}}' -X 'github.com/stonecharioteer/goforgo/internal/cli.commit={{commit}}' -X 'github.com/stonecharioteer/goforgo/internal/cli.date={{date}}'" -o dist/goforgo-windows-amd64.exe ./cmd/goforgo
    @echo "✅ Release binaries built in dist/"

# Check if code is ready for commit
pre-commit: fmt lint test
    @echo "✅ Code is ready for commit!"

# Get version from git
version := `git describe --tags --always --dirty 2>/dev/null || echo "dev"`

# Get commit hash
commit := `git rev-parse --short HEAD 2>/dev/null || echo "unknown"`

# Get build date
date := `date -u +%Y-%m-%dT%H:%M:%SZ`

# Show project info
info:
    @echo "📋 GoForGo Project Information"
    @echo "=============================="
    @echo "Version: {{version}}"
    @echo "Commit:  {{commit}}"
    @echo "Date:    {{date}}"
    @echo "Go:      $(go version)"
    @echo ""
    @echo "📊 Project Stats:"
    @echo "Lines of Go code: $(find . -name '*.go' -not -path './vendor/*' | xargs wc -l | tail -1 | awk '{print $1}')"
    @echo "Number of packages: $(go list ./... | wc -l)"

# Watch for changes and rebuild (requires entr)
watch:
    @echo "👀 Watching for changes... (requires 'entr' to be installed)"
    find . -name '*.go' | entr -r just dev-build

# Generate documentation
docs:
    @echo "📚 Generating documentation..."
    go doc -all ./... > docs/api.txt
    @echo "✅ Documentation generated in docs/"