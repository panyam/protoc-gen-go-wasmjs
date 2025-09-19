# Core Layer Design

## Purpose

The core layer provides **pure utility functions** for analyzing protobuf definitions, calculating file paths, and converting between naming conventions. This is Layer 1 of the architecture - it has no external dependencies and consists entirely of stateless, side-effect-free functions.

## Design Principles

### 1. **Pure Functions Only**
- All functions are stateless with no side effects
- Same input always produces same output
- No external dependencies (files, network, etc.)
- Easy to test in isolation

### 2. **Single Responsibility** 
- Each component has one focused purpose
- **ProtoAnalyzer**: Protobuf metadata extraction
- **PathCalculator**: File path manipulation
- **NameConverter**: Naming convention transformations

### 3. **Cross-Platform Compatibility**
- Path operations work on Windows, macOS, and Linux
- Consistent behavior across different environments
- Proper handling of path separators and edge cases

## Key Components

### ProtoAnalyzer (`proto_analyzer.go`)

**Purpose**: Extract metadata and properties from protobuf definitions

**Key Functions**:
```go
// Package and type analysis
ExtractPackageName(fullMessageType string) string
ExtractMessageName(fullMessageType string) string
GetBaseFileName(protoFile string) string

// Field analysis
GetProtoFieldType(field *protogen.Field) string
GetFullyQualifiedMessageType(field *protogen.Field) string
IsMapField(field *protogen.Field) bool
GetMapKeyValueTypes(field *protogen.Field) (string, string)

// Annotation analysis (requires protogen objects)
IsBrowserProvidedService(service *protogen.Service) bool
GetCustomMethodName(method *protogen.Method) string
IsAsyncMethod(method *protogen.Method) bool
IsMethodExcluded(method *protogen.Method) bool
IsServiceExcluded(service *protogen.Service) bool
```

**Design Decisions**:
- **String-based analysis** for functions that can be pure (package names, file paths)
- **Protogen-dependent functions** for complex annotation analysis (tested via integration)
- **Consistent return types** (empty string for not found, boolean for existence checks)

### PathCalculator (`path_calculator.go`)

**Purpose**: Handle all file path calculations and import path resolution

**Key Functions**:
```go
// Relative path calculations
CalculateRelativePath(fromPath, toPath string) string
BuildCrossPackageImportPath(currentPackage, targetPackage string) string
GetFactoryImportPath(dependencyPackage, currentPackage string) string

// Package structure
BuildPackagePath(packageName string) string
GenerateOutputFilePath(baseOutputPath, packageName, fileName string) string

// Path utilities
NormalizePath(path string) string
IsAbsolutePath(path string) bool
JoinPaths(components ...string) string

// Go package handling
GetGoPackageAlias(packagePath string) string
```

**Design Decisions**:
- **Forward slashes** always used for TypeScript imports (cross-platform)
- **Relative path preservation** with proper `./` and `../` handling  
- **Package aliasing** combines last two path segments for readable Go imports
- **Graceful fallbacks** for edge cases (empty paths, malformed input)

### NameConverter (`name_converter.go`)

**Purpose**: Convert between different naming conventions (Go, JavaScript, TypeScript)

**Key Functions**:
```go
// Case conversions
ToCamelCase(s string) string       // PascalCase → camelCase
ToPascalCase(s string) string      // camelCase → PascalCase  
ToSnakeCase(s string) string       // Any case → snake_case

// Specialized naming
ToPackageAlias(packagePath string) string        // Go package aliases
ToJSNamespace(packageName string) string         // JavaScript namespaces
ToModuleName(packageName string) string          // WASM module names
ToFactoryName(packageName string) string         // TypeScript factory classes
ToGoFuncName(serviceName, methodName string) string // Go WASM function names

// Validation
SanitizeIdentifier(name string) string           // Ensure valid identifiers
```

**Design Decisions**:
- **Language-specific conventions** (camelCase for JS, PascalCase for TS classes)
- **Special character handling** (hyphens → underscores, dots → removed)
- **Descriptive naming** (factory names include "Factory" suffix)
- **Fallback safety** (invalid identifiers get underscore prefix/replacement)

## Testing Strategy

### 100% Pure Function Coverage
- **30+ unit tests** covering all string manipulation functions
- **Cross-platform tests** for path operations  
- **Edge case coverage** (empty strings, invalid input, boundary conditions)
- **Self-documenting tests** with detailed explanations of what and why

### Test Structure Example
```go
func TestProtoAnalyzer_ExtractPackageName(t *testing.T) {
    tests := []struct {
        name             string // Test case description
        fullMessageType  string // Input
        expectedPackage  string // Expected output
        reason           string // Why this test matters
    }{
        {
            name:            "standard package with version",
            fullMessageType: "library.v1.Book",
            expectedPackage: "library.v1", 
            reason:          "Most common case - correct extraction is critical",
        },
        // ... more test cases
    }
}
```

### Integration Testing
Functions requiring `protogen` objects are tested through the complete generator pipeline in `examples/`.

## Usage Examples

### Typical Usage in Higher Layers
```go
// In filter layer
analyzer := core.NewProtoAnalyzer()
if analyzer.IsBrowserProvidedService(service) {
    // Handle browser service
}

// In builder layer  
pathCalc := core.NewPathCalculator()
importPath := pathCalc.BuildCrossPackageImportPath("library.v1", "common.v1")

// In generator layer
nameConv := core.NewNameConverter()
jsName := nameConv.ToCamelCase(methodName)
```

### Chain Operations
```go
// Complex path and naming operations
packagePath := pathCalc.BuildPackagePath(packageName)
alias := pathCalc.GetGoPackageAlias(importPath)
moduleName := nameConv.ToModuleName(packageName)
filename := pathCalc.JoinPaths(outputPath, packagePath, baseName+".wasm.go")
```

## Extension Points

### Adding New Analyzers
```go
// Add new proto analysis functions
func (pa *ProtoAnalyzer) GetCustomAnnotation(method *protogen.Method) string {
    // Analyze new annotation types
}
```

### Adding New Naming Conventions
```go
// Add new naming patterns
func (nc *NameConverter) ToKebabCase(s string) string {
    // Convert to kebab-case for certain contexts
}
```

### Adding New Path Strategies
```go
// Add new path calculation strategies
func (pc *PathCalculator) BuildFlatImportPath(packages ...string) string {
    // Alternative import path strategy
}
```

## Architecture Benefits

### For Testing
- **Fast tests** (<100ms for all core tests)
- **No mocking required** for pure functions
- **Easy to isolate** and test individual functions
- **Clear test failures** with descriptive error messages

### For Maintenance
- **Pure functions** are easy to understand and modify
- **No hidden dependencies** or global state
- **Predictable behavior** across all use cases
- **Safe refactoring** - breaking changes are caught by tests

### For Extension
- **Composable functions** can be combined in new ways
- **Consistent patterns** for adding new functionality
- **Clear interfaces** between components
- **Backward compatibility** easy to maintain

## Dependencies

### Internal Dependencies
- None - this is the foundation layer

### External Dependencies  
- `google.golang.org/protobuf/compiler/protogen` (for protogen-dependent functions only)
- Standard library only (strings, path/filepath, unicode)

## Files

- `proto_analyzer.go` - Protobuf definition analysis (15 functions)
- `path_calculator.go` - File path calculations (9 functions)
- `name_converter.go` - Naming conversions (10 functions)
- `*_test.go` - Comprehensive test suites (30+ test cases)

## Future Considerations

### Performance
- All functions are already optimized for common cases
- Could add caching for expensive operations if needed
- String operations are efficient for typical proto sizes

### Functionality
- Could add more protobuf analysis functions as needed
- Could add more naming convention support
- Could add more path manipulation utilities

### Testing
- Consider adding property-based testing for path operations
- Could add benchmarks for performance-critical functions
- Consider adding fuzzing for string manipulation functions
