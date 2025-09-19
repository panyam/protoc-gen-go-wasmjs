# Filter Layer Design

## Purpose

The filter layer implements **business logic for determining what to generate**. This is Layer 2 of the architecture - it takes protobuf definitions and configuration criteria, then decides which services, methods, messages, and enums should be included in code generation.

## Design Principles

### 1. **Rich Result Types**
- Filtering operations return detailed results explaining decisions
- Human-readable reasons for inclusion/exclusion
- Additional metadata for template generation

### 2. **Centralized Configuration**
- Single `FilterCriteria` type handles all filtering options
- Consistent parsing and validation across all filter types
- Easy to extend with new filtering options

### 3. **Composable Filters**
- Each filter handles one concern (services, methods, messages, etc.)
- Filters can be combined for complex scenarios
- Statistics collection for debugging and reporting

## Key Components

### FilterCriteria (`filter_config.go`)

**Purpose**: Centralized configuration for all filtering decisions

**Key Types**:
```go
type FilterCriteria struct {
    // Service filtering
    ServicesSet map[string]bool // Specific services to include
    
    // Method filtering with glob patterns
    MethodIncludes []string      // Patterns to include
    MethodExcludes []string      // Patterns to exclude
    MethodRenames  map[string]string // Name transformations
    
    // Content filtering flags
    ExcludeAnnotationPackages bool // Skip wasmjs.v1 packages
    ExcludeEmptyPackages     bool // Skip packages with no content
    ExcludeMapEntries        bool // Skip synthetic map messages
}
```

**Design Decisions**:
- **String-based configuration** for easy parsing from command line
- **Boolean flags** for common filtering scenarios
- **Map structures** for efficient lookup during filtering
- **Validation during parsing** to catch errors early

### ServiceFilter (`service_filter.go`)

**Purpose**: Determine which gRPC services should be included in generation

**Key Interface**:
```go
type ServiceFilter struct {
    analyzer *core.ProtoAnalyzer
}

// Main filtering logic
func ShouldIncludeService(service *protogen.Service, criteria *FilterCriteria) ServiceFilterResult

// Batch operations
func FilterServices(files []*protogen.File, criteria *FilterCriteria) ([]ServiceFilterResult, *FilterStats)
func GetBrowserProvidedServices(files []*protogen.File, criteria *FilterCriteria) []*protogen.Service
```

**Filtering Priority Order**:
1. **Annotation exclusion** (highest priority) - `wasm_service_exclude`
2. **Service list filtering** - explicit services configuration
3. **Default inclusion** - include if not explicitly excluded

**Design Decisions**:
- **Annotation takes precedence** over configuration
- **Browser service detection** affects generation but doesn't exclude
- **Rich result types** include metadata for template generation
- **Batch operations** for efficiency when processing multiple services

### MethodFilter (`method_filter.go`)

**Purpose**: Determine which methods within services should be included

**Key Interface**:
```go
type MethodFilter struct {
    analyzer *core.ProtoAnalyzer
}

// Main filtering logic
func ShouldIncludeMethod(method *protogen.Method, criteria *FilterCriteria) MethodFilterResult

// Utilities
func ValidateMethodPatterns(criteria *FilterCriteria) error
func GetMethodJSName(method *protogen.Method, criteria *FilterCriteria, nameConverter *core.NameConverter) string
```

**Filtering Priority Order**:
1. **Annotation exclusion** - `wasm_method_exclude`
2. **Streaming limitations** - client streaming not supported
3. **Exclude patterns** - glob pattern matching
4. **Include patterns** (if configured) - must match to be included
5. **Default inclusion** - include if not excluded

**Design Decisions**:
- **Glob pattern support** for flexible method selection (`Get*`, `*Internal`)
- **Streaming validation** prevents unsupported method types
- **Custom naming** from annotations takes precedence over config renames
- **Pattern validation** catches invalid glob syntax early

### MessageCollector (`message_collector.go`)

**Purpose**: Collect and filter protobuf messages across packages

**Key Interface**:
```go
type MessageCollector struct {
    analyzer *core.ProtoAnalyzer
}

// Collection operations
func CollectMessages(files []*protogen.File, criteria *FilterCriteria) CollectionResult[MessageInfo]
func CollectMessagesByPackage(files []*protogen.File, criteria *FilterCriteria) map[string]CollectionResult[MessageInfo]

// Utilities
func HasAnyMessages(files []*protogen.File, criteria *FilterCriteria) bool
func GetMessagesByFile(result CollectionResult[MessageInfo]) map[string][]MessageInfo
```

**Design Decisions**:
- **Recursive collection** handles nested messages correctly
- **Map entry filtering** excludes synthetic messages for map fields
- **Metadata preservation** (package names, proto files, comments)
- **Generic result types** with statistics for debugging

### EnumCollector (`enum_collector.go`)

**Purpose**: Collect and filter protobuf enums across packages

**Similar interface to MessageCollector but for enums**:
```go
func CollectEnums(files []*protogen.File, criteria *FilterCriteria) CollectionResult[EnumInfo]
```

### PackageFilter (`package_filter.go`)

**Purpose**: Determine which packages should be processed for generation

**Key Interface**:
```go
type PackageFilter struct {
    analyzer          *core.ProtoAnalyzer
    messageCollector  *MessageCollector
    enumCollector     *EnumCollector
}

func ShouldIncludePackage(packageName string, files []*protogen.File, criteria *FilterCriteria) PackageFilterResult
func FilterPackages(allFiles []*protogen.File, criteria *FilterCriteria) (map[string][]*protogen.File, *FilterStats)
```

**Filtering Logic**:
1. **Annotation package exclusion** - skip `wasmjs.v1` and similar
2. **Empty package detection** - skip packages with no services/messages/enums
3. **Content analysis** - determine what types of content packages have

## Result Types

### FilterResult
```go
type FilterResult struct {
    Include bool   // Whether to include the item
    Reason  string // Human-readable explanation
}
```

### Specialized Results
```go
type ServiceFilterResult struct {
    FilterResult
    IsBrowserProvided bool   // Service metadata
    CustomName        string // Custom JavaScript name
}

type MethodFilterResult struct {
    FilterResult
    CustomJSName      string // Custom JavaScript method name
    IsAsync           bool   // Async handling required
    IsServerStreaming bool   // Streaming metadata
}
```

### Statistics
```go
type FilterStats struct {
    ServicesTotal, ServicesIncluded, ServicesExcluded int
    MethodsTotal, MethodsIncluded, MethodsExcluded    int
    MessagesTotal, EnumsTotal, PackagesTotal          int
}
```

## Usage Patterns

### Basic Service Filtering
```go
criteria := filters.NewFilterCriteria()
criteria.ServicesSet = map[string]bool{"UserService": true}

filter := filters.NewServiceFilter(analyzer)
result := filter.ShouldIncludeService(service, criteria)
if result.Include {
    // Generate this service
}
```

### Complex Method Filtering
```go
criteria, _ := filters.ParseFromConfig(
    "UserService,LibraryService",  // services
    "Get*,Find*",                  // includes  
    "*Internal,*Debug",            // excludes
    "FindBooks:searchBooks",       // renames
)

methodFilter := filters.NewMethodFilter(analyzer)
result := methodFilter.ShouldIncludeMethod(method, criteria)
```

### Batch Collection
```go
collector := filters.NewMessageCollector(analyzer)
messages := collector.CollectMessages(files, criteria)
fmt.Printf("Found %d messages from %d files", len(messages.Items), messages.FilesScanned)
```

## Testing Strategy

### Unit Tests (25+ test cases)
- **Configuration parsing** with edge cases and error handling
- **Filter result creation** and metadata validation
- **Pattern validation** for glob syntax checking
- **Complex scenarios** combining multiple filter types

### Integration Tests
- **Component coordination** between different filters
- **Real-world scenarios** from actual usage patterns
- **Statistics accuracy** for debugging support

### End-to-End Testing
- Full filtering workflows tested through example proto files
- Annotation-based filtering tested with real protobuf definitions

## Dependencies

### Internal Dependencies
- `pkg/core` - For proto analysis, path calculations, name conversions

### External Dependencies
- `google.golang.org/protobuf/compiler/protogen` - For protobuf definitions
- `path/filepath` - For glob pattern matching
- Standard library only

## Extension Points

### Adding New Filter Types
```go
// Add new specialized filter
type FieldFilter struct {
    analyzer *core.ProtoAnalyzer
}

func (ff *FieldFilter) ShouldIncludeField(field *protogen.Field, criteria *FilterCriteria) FieldFilterResult
```

### Adding New Criteria
```go
// Extend FilterCriteria with new options
type FilterCriteria struct {
    // ... existing fields
    
    // New filtering options
    FieldIncludes []string      // Field-level filtering
    TypeMappings  map[string]string // Custom type mappings
}
```

### Adding New Result Metadata
```go
// Extend result types with new metadata
type ServiceFilterResult struct {
    FilterResult
    // ... existing fields
    
    // New metadata
    Complexity    int    // Method count or complexity score
    Dependencies  []string // Cross-package dependencies
}
```

## Common Patterns

### Filter Composition
```go
// Combine multiple filters for complex logic
packageFilter := NewPackageFilter(analyzer, msgCollector, enumCollector)
serviceFilter := NewServiceFilter(analyzer)

packages, _ := packageFilter.FilterPackages(files, criteria)
for _, files := range packages {
    services := serviceFilter.GetIncludedServices(files, criteria)
    // Process filtered services
}
```

### Statistics Collection
```go
stats := filters.NewFilterStats()
for _, service := range services {
    result := filter.ShouldIncludeService(service, criteria)
    stats.AddServiceResult(result)
}
fmt.Println(stats.Summary()) // "3/5 services, 15/20 methods, ..."
```

### Error Handling
```go
// Validate configuration before using
if err := methodFilter.ValidateMethodPatterns(criteria); err != nil {
    return fmt.Errorf("invalid method patterns: %w", err)
}
```
