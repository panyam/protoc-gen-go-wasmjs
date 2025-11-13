# Development Guide

This guide covers development workflows, testing strategies, and contribution guidelines for protoc-gen-go-wasmjs.

## Table of Contents

- [Development Setup](#development-setup)
- [Testing Strategy](#testing-strategy)
- [Running Tests](#running-tests)
- [Test Organization](#test-organization)
- [Adding New Tests](#adding-new-tests)
- [Development Workflow](#development-workflow)
- [Code Quality Standards](#code-quality-standards)

## Development Setup

### Prerequisites

- Go 1.23.5 or later
- Protocol Buffers compiler (`protoc`)
- Buf CLI (recommended for proto management)

### Building the Plugin

```bash
# Build the main plugin
make tool

# Build all targets (includes stateful plugin and WASM binary)
make all

# Install to GOBIN
make install
```

### Project Structure

The project follows a layered architecture designed for testability and maintainability:

```
pkg/
├── core/                    # Layer 1: Pure utility functions
│   ├── proto_analyzer.go    # Proto file parsing & metadata extraction
│   ├── path_calculator.go   # File path & import calculations  
│   ├── name_converter.go    # Naming convention conversions
│   └── *_test.go           # Comprehensive test suites
│
├── filters/                # Layer 2: Business logic ✅
│   ├── filter_config.go    # Filter criteria & configuration parsing
│   ├── filter_result.go    # Result types & statistics
│   ├── service_filter.go   # Service inclusion/exclusion
│   ├── method_filter.go    # Method filtering with patterns
│   ├── message_collector.go # Message collection & filtering
│   ├── enum_collector.go   # Enum collection & filtering
│   ├── package_filter.go   # Package-level filtering
│   └── *_test.go          # Unit & integration tests
│
├── builders/               # Layer 3: Template data building (planned)
├── renderers/              # Layer 4: Template execution (planned)
└── generators/             # Top-level generators (planned)
```

## Testing Strategy

Our testing approach emphasizes **comprehensive coverage with clear documentation** of what each test validates and why it matters.

### Testing Philosophy

1. **Pure Functions First**: Core utilities are extracted as pure functions that are easy to test in isolation
2. **Comprehensive Edge Cases**: Every function tests normal cases, edge cases, and error conditions
3. **Self-Documenting Tests**: Each test includes detailed comments explaining its purpose
4. **Fast Feedback**: Unit tests run quickly and provide immediate feedback
5. **Layered Testing**: Different layers have different testing strategies

### Test Categories

#### 1. Unit Tests (Layer 1 - Core Utilities)

**Purpose**: Test individual pure functions in isolation  
**Coverage**: 100% of core utilities  
**Speed**: Very fast (<100ms total)

**Example Structure**:
```go
// TestProtoAnalyzer_ExtractPackageName tests the extraction of package names from
// fully qualified message types. This is critical for cross-package imports and
// dependency resolution in TypeScript generation.
func TestProtoAnalyzer_ExtractPackageName(t *testing.T) {
    analyzer := NewProtoAnalyzer()
    
    tests := []struct {
        name             string // Test case description
        fullMessageType  string // Input: fully qualified message type
        expectedPackage  string // Expected output: package name
        reason           string // Why this test case is important
    }{
        {
            name:            "standard package with version",
            fullMessageType: "library.v1.Book",
            expectedPackage: "library.v1",
            reason:          "Most common case - package with version needs correct extraction",
        },
        // ... more test cases
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := analyzer.ExtractPackageName(tt.fullMessageType)
            if result != tt.expectedPackage {
                t.Errorf("ExtractPackageName(%s) = %s, want %s\nReason: %s", 
                    tt.fullMessageType, result, tt.expectedPackage, tt.reason)
            }
        })
    }
}
```

#### 2. Integration Tests (Planned)

**Purpose**: Test interactions between layers  
**Coverage**: Cross-layer functionality  
**Speed**: Moderate (1-5 seconds)

#### 3. End-to-End Tests (Existing)

**Purpose**: Test complete generation pipeline  
**Coverage**: Full workflows with example projects  
**Speed**: Slower (10-30 seconds)

## Running Tests

### Core Unit Tests

Run the fast unit tests for core utilities:

```bash
# Run all core tests
go test ./pkg/core/... -v

# Run with coverage
go test ./pkg/core/... -v -cover

# Run specific test
go test ./pkg/core/ -run TestProtoAnalyzer_ExtractPackageName -v

# Run tests in parallel
go test ./pkg/core/... -v -parallel 4
```

### Filter Layer Tests

Run the business logic tests for filtering:

```bash
# Run all filter tests
go test ./pkg/filters/... -v

# Run specific filter component
go test ./pkg/filters/ -run TestServiceFilter -v
go test ./pkg/filters/ -run TestMethodFilter -v

# Run integration tests
go test ./pkg/filters/ -run TestFilterLayer_Integration -v
```

### All Tests

```bash
# Run all tests in the project
go test ./... -v

# Run with race detection
go test ./... -v -race

# Generate coverage report
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

### Example Tests

Test the plugin against real examples:

```bash
# Test library example
cd examples/library && buf generate

# Test connect4 example
cd examples/connect4 && make test

# Test all examples
make test-examples
```

## Test Organization

### Test File Naming

- `*_test.go` for unit tests in the same package
- Test files mirror the structure of source files
- One test file per source file for clarity

### Test Function Naming

```go
// Pattern: Test<Type>_<Method>_<Scenario>
func TestProtoAnalyzer_ExtractPackageName(t *testing.T)          // Basic functionality
func TestPathCalculator_CalculateRelativePath_CrossPlatform(t *testing.T) // Specific scenario
func TestNameConverter_SanitizeIdentifier_EdgeCases(t *testing.T)         // Edge cases
```

### Test Structure

Each test follows a consistent structure:

```go
func TestComponent_Method(t *testing.T) {
    // 1. Setup
    component := NewComponent()
    
    // 2. Test cases with clear documentation
    tests := []struct {
        name     string // What this test case covers
        input    string // Input parameters
        expected string // Expected result
        reason   string // Why this test matters
    }{
        {
            name:     "standard case",
            input:    "typical input",
            expected: "expected output",
            reason:   "Most common usage pattern",
        },
        // ... more cases including edge cases
    }

    // 3. Execute tests
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := component.Method(tt.input)
            if result != tt.expected {
                t.Errorf("Method(%s) = %s, want %s\nReason: %s",
                    tt.input, result, tt.expected, tt.reason)
            }
        })
    }
}
```

## Adding New Tests

When adding new functionality, follow these guidelines:

### 1. Write Tests First (TDD)

```bash
# 1. Write the test
echo "func TestNewFunction(t *testing.T) { /* test */ }" >> pkg/core/component_test.go

# 2. Run test (should fail)
go test ./pkg/core/ -run TestNewFunction

# 3. Implement function
# 4. Run test (should pass)
go test ./pkg/core/ -run TestNewFunction -v
```

### 2. Cover All Cases

For each new function, ensure tests cover:

- **Happy path**: Normal usage with valid inputs
- **Edge cases**: Empty strings, nil values, boundary conditions
- **Error cases**: Invalid inputs, malformed data
- **Cross-platform**: Path separators, case sensitivity (if applicable)

### 3. Document Test Purpose

Each test case should include:

```go
{
    name:     "descriptive test name",           // What scenario
    input:    "test input",                     // What input
    expected: "expected result",                // What output
    reason:   "Why this test is important",     // Why it matters
}
```

### 4. Test Error Messages

Include the reason in error messages:

```go
if result != expected {
    t.Errorf("Function(%v) = %v, want %v\nReason: %s",
        input, result, expected, reason)
}
```

## Development Workflow

### 1. Local Development

```bash
# 1. Make changes
# 2. Run relevant tests
go test ./pkg/core/... -v

# 3. Run full test suite
go test ./... -v

# 4. Build and test with examples
make tool
cd examples/library && buf generate
```

### 2. Before Committing

```bash
# Run all tests
go test ./... -v

# Check code formatting
go fmt ./...

# Run linters (if configured)
golangci-lint run

# Test build
make all
```

### 3. Adding New Features

1. **Plan**: Update architecture docs if needed
2. **Extract**: Add pure functions to appropriate core layer
3. **Test**: Write comprehensive unit tests
4. **Integrate**: Update existing code to use new utilities
5. **Validate**: Run existing examples to ensure backward compatibility

## Code Quality Standards

### Test Quality Requirements

- **High coverage** for testable core utilities (pkg/core/)
  - **Pure functions**: 90%+ coverage (path calculations, name conversions)
  - **Proto-dependent functions**: 40-60% coverage (require complex mocks)
  - **Overall target**: 60%+ for core utilities package
- **Meaningful test names** that describe the scenario
- **Clear error messages** that help debug failures
- **Fast execution** (<1 second for all unit tests)
- **Cross-platform compatibility** for path operations

### Coverage Considerations

Some functions in the core utilities depend on `protogen` types that are complex to mock:
- `IsBrowserProvidedService()`: Requires `protogen.Service` with options
- `GetCustomMethodName()`: Requires `protogen.Method` with annotations
- `IsMapField()`: Requires `protogen.Field` with message types

These functions are integration-tested through the full generator pipeline in the examples.

### Code Documentation

- All public functions must have comprehensive comments
- Test functions must document what they test and why
- Complex logic should have inline comments
- README and architecture docs must be kept up to date

### Error Handling

- Handle edge cases gracefully (empty inputs, invalid data)
- Provide meaningful error messages
- Use fallbacks for non-critical failures
- Test error conditions explicitly

## Continuous Integration

### GitHub Actions (Future)

The project should include CI that runs:

```yaml
# .github/workflows/test.yml
name: Tests
on: [push, pull_request]
jobs:
  test:
    strategy:
      matrix:
        os: [ubuntu-latest, windows-latest, macos-latest]
        go: ['1.23.5']
    runs-on: ${{ matrix.os }}
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v4
        with:
          go-version: ${{ matrix.go }}
      - run: go test ./... -v -race
      - run: make all
      - run: make test-examples
```

### Test Coverage Goals

- **Core utilities**: 100% line and branch coverage
- **Business logic**: 90%+ coverage
- **Integration**: Key workflows covered
- **Examples**: All examples must build and generate valid output

## Debugging Tests

### Common Issues

1. **Cross-platform path failures**:
   ```bash
   # Use filepath.ToSlash() for consistent separators
   result = filepath.ToSlash(result)
   expected = filepath.ToSlash(expected)
   ```

2. **Flaky tests**:
   ```bash
   # Run multiple times to detect
   go test ./pkg/core/ -count=10 -run TestProblematic
   ```

3. **Race conditions**:
   ```bash
   # Enable race detector
   go test ./... -race
   ```

### Useful Test Commands

```bash
# Verbose output with test names
go test ./pkg/core/... -v

# Run only failed tests
go test ./pkg/core/... -v -failfast

# Test with coverage and race detection
go test ./pkg/core/... -v -cover -race

# Profile memory usage
go test ./pkg/core/... -memprofile=mem.prof

# Profile CPU usage
go test ./pkg/core/... -cpuprofile=cpu.prof
```

## Getting Help

- **Architecture**: See [ARCHITECTURE.md](ARCHITECTURE.md)
- **Usage**: See [README.md](README.md)
- **Examples**: Check `examples/` directory
- **Issues**: Use GitHub issues for questions and bug reports

## Contributing

1. **Fork** the repository
2. **Create** a feature branch
3. **Write** tests for new functionality
4. **Implement** the feature
5. **Ensure** all tests pass
6. **Update** documentation
7. **Submit** a pull request

### Pull Request Checklist

- [ ] All tests pass (`go test ./... -v`)
- [ ] Code is formatted (`go fmt ./...`)
- [ ] New functions have comprehensive tests
- [ ] Documentation is updated
- [ ] Examples still work
- [ ] Runtime package builds (`cd runtime && pnpm run build`)
- [ ] TypeScript types are valid (`cd runtime && pnpm run typecheck`)
- [ ] Generated code compiles (test with examples)
- [ ] Backward compatibility is maintained

## Runtime Package Development

The TypeScript runtime package provides shared utilities for generated code:

### **Development Workflow**
```bash
# Install dependencies
cd runtime && pnpm install

# Build the package
pnpm run build

# Run type checking
pnpm run typecheck

# Watch mode for development
pnpm run dev
```

### **Testing Runtime Package**
```bash
# Test with examples
cd example && make buf

# Verify imports resolve
cd example/web && pnpm run typecheck
```

### **Publishing Runtime Package** (Future)
```bash
cd runtime
pnpm run build
pnpm publish --access public
```
