#!/bin/bash
set -e

echo "🧪 Running protoc-gen-go-wasmjs test suite..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print status
print_status() {
    echo -e "${BLUE}▶${NC} $1"
}

print_success() {
    echo -e "${GREEN}✅${NC} $1"
}

print_error() {
    echo -e "${RED}❌${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠️${NC} $1"
}

# Check if Go is installed
if ! command -v go &> /dev/null; then
    print_error "Go is not installed or not in PATH"
    exit 1
fi

print_success "Go $(go version | cut -d' ' -f3) detected"

# 1. Run core unit tests
print_status "Running core utility tests..."
if go test ./pkg/core/... -v -timeout=30s; then
    print_success "Core utility tests passed"
else
    print_error "Core utility tests failed"
    exit 1
fi

# 1.5. Run filter layer tests
print_status "Running filter layer tests..."
if go test ./pkg/filters/... -v -timeout=30s; then
    print_success "Filter layer tests passed"
else
    print_error "Filter layer tests failed"
    exit 1
fi

# 1.6. Run builder layer tests
print_status "Running builder layer tests..."
if go test ./pkg/builders/... -v -timeout=30s; then
    print_success "Builder layer tests passed"
else
    print_error "Builder layer tests failed"
    exit 1
fi

# 1.7. Run new generator layer tests
print_status "Running new generator layer tests..."
if go test ./pkg/generators/... -v -timeout=30s; then
    print_success "New generator layer tests passed"
else
    print_error "New generator layer tests failed"
    exit 1
fi

# 2. Run all tests with race detection
print_status "Running all tests with race detection..."
if go test ./... -race -timeout=60s; then
    print_success "All tests passed with race detection"
else
    print_error "Tests failed with race detection"
    exit 1
fi

# 3. Check test coverage for core utilities
print_status "Checking test coverage for core utilities..."
COVERAGE=$(go test ./pkg/core/... -cover -covermode=count -coverprofile=/tmp/core_coverage.out 2>&1 | grep "coverage:" | tail -1 | grep -oE '[0-9]+\.[0-9]+%' || echo "0.0%")
echo "Core utilities coverage: $COVERAGE"

# Parse coverage percentage
COVERAGE_NUM=$(echo $COVERAGE | sed 's/%//')
if (( $(echo "$COVERAGE_NUM >= 80.0" | bc -l) )); then
    print_success "Excellent test coverage: $COVERAGE"
elif (( $(echo "$COVERAGE_NUM >= 60.0" | bc -l) )); then
    print_success "Good test coverage: $COVERAGE (some functions require protogen mocks)"
elif (( $(echo "$COVERAGE_NUM >= 40.0" | bc -l) )); then
    print_warning "Moderate test coverage: $COVERAGE (aim for 60%+ for pure functions)"
else
    print_error "Low test coverage: $COVERAGE (should be 40%+)"
    exit 1
fi

# 4. Build all plugins
print_status "Building plugins..."
if make all > /dev/null 2>&1; then
    print_success "All plugins built successfully"
else
    print_error "Plugin build failed"
    exit 1
fi

# 5. Check if plugin binaries exist
EXPECTED_BINARIES=(
    "./bin/protoc-gen-go-wasmjs"
    "./bin/protoc-gen-go-wasmjs-go"
    "./bin/protoc-gen-go-wasmjs-ts"
    "./bin/protoc-gen-go-wasmjs-stateful"
)

for binary in "${EXPECTED_BINARIES[@]}"; do
    if [ -f "$binary" ]; then
        print_success "Binary created: $binary"
    else
        print_error "Binary not found: $binary"
        exit 1
    fi
done

# 6. Test with library example (if buf is available)
if command -v buf &> /dev/null; then
    print_status "Testing with library example..."
    cd examples/library
    
    # Test the original generator (new split generators need template integration)
    if make bufdev > /dev/null 2>&1; then
        print_success "Library example generated successfully"
        
        # Check if expected files were generated  
        EXPECTED_FILES=(
            "gen/go/library/v1/library.pb.go"
            "gen/go/library/v1/library_grpc.pb.go"
        )
        
        for file in "${EXPECTED_FILES[@]}"; do
            if [ -f "$file" ]; then
                print_success "Generated: $file"
            else
                print_warning "Missing: $file"
            fi
        done
        
        # Check for WASM artifacts
        WASM_FILES=$(find gen -name "*.wasm.go" 2>/dev/null || true)
        if [ -n "$WASM_FILES" ]; then
            print_success "WASM artifacts generated"
        else
            print_warning "No WASM artifacts found"
        fi
        
        # Check for TypeScript artifacts
        TS_FILES=$(find gen -name "*.ts" 2>/dev/null || true)
        if [ -n "$TS_FILES" ]; then
            print_success "TypeScript artifacts generated"
        else
            print_warning "No TypeScript artifacts found"
        fi
        
    else
        print_error "Library example generation failed"
        cd - > /dev/null
        exit 1
    fi
    
    cd - > /dev/null
else
    print_warning "Buf not found, skipping example test"
fi

# 7. Code formatting check
print_status "Checking code formatting..."
UNFORMATTED=$(go fmt ./... 2>&1 || true)
if [ -z "$UNFORMATTED" ]; then
    print_success "Code is properly formatted"
else
    print_error "Code formatting issues found:"
    echo "$UNFORMATTED"
    exit 1
fi

# 8. Check for common Go issues (if available)
if command -v go vet &> /dev/null; then
    print_status "Running go vet..."
    if go vet ./... > /dev/null 2>&1; then
        print_success "No issues found with go vet"
    else
        print_error "Issues found with go vet:"
        go vet ./...
        exit 1
    fi
fi

# Summary
echo ""
echo "🎉 All tests passed! Summary:"
print_success "✅ Core utility tests: PASSED"
print_success "✅ Filter layer tests: PASSED"
print_success "✅ Builder layer tests: PASSED"
print_success "✅ New generator layer tests: PASSED"
print_success "✅ All tests with race detection: PASSED"
print_success "✅ Test coverage: $COVERAGE"
print_success "✅ All plugins build: SUCCESS"
print_success "✅ Code formatting: CLEAN"
if command -v buf &> /dev/null; then
    print_success "✅ Example generation: SUCCESS"
fi

echo ""
echo "Ready for development! 🚀"
echo ""
echo "Quick commands:"
echo "  go test ./pkg/core/... -v        # Run core tests"
echo "  go test ./pkg/filters/... -v     # Run filter tests"
echo "  go test ./pkg/builders/... -v    # Run builder tests"
echo "  go test ./pkg/generators/... -v  # Run generator tests"
echo "  go test ./pkg/... -v             # Run all new layer tests"
echo "  make split                       # Build new split generators"
echo "  make tool                        # Build original generator"
echo "  cd examples/library && make bufdev  # Test with example"
echo ""
echo "See DEVELOPMENT.md for detailed development guide."
