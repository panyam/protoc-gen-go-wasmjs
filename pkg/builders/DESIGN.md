# Builder Layer Design

## Purpose

The builder layer transforms **filtered protobuf data into template-ready structures**. This is Layer 3 of the architecture - it takes the output from filters and creates rich data structures that templates can consume to generate code.

## Design Principles

### 1. **Template-Focused Data Structures**
- Build data structures optimized for template consumption
- Include all metadata templates need (names, types, imports, etc.)
- Separate data building from template execution

### 2. **Generator-Specific Builders**
- **GoDataBuilder**: Focuses on Go WASM generation needs
- **TSDataBuilder**: Focuses on TypeScript generation needs
- **Shared types** for common data structures

### 3. **File Planning Architecture**
- Generators control all file creation decisions
- File specifications separate logical names from physical filenames
- Content hints guide conditional file generation

## Key Components

### Shared Types (`shared_types.go`)

**Purpose**: Common data structures used by both Go and TypeScript builders

**Core Types**:
```go
// Configuration for both generators
type GenerationConfig struct {
    WasmExportPath      string // Where to write WASM files
    TSExportPath        string // Where to write TypeScript files
    JSStructure         string // namespaced|flat|service_based
    JSNamespace         string // Global JavaScript namespace
    ModuleName          string // WASM module name
    GenerateBuildScript bool   // Whether to generate build scripts
}

// Import information for templates
type ImportInfo struct {
    Path  string // Full import path
    Alias string // Package alias
}

// Service data for templates
type ServiceData struct {
    Name              string      // Service name
    GoType            string      // Go interface type (for WASM)
    JSName            string      // JavaScript name
    IsBrowserProvided bool        // Service type
    Methods           []MethodData // Filtered methods
}

// Method data for templates
type MethodData struct {
    Name              string // Original method name
    JSName            string // JavaScript method name
    GoFuncName        string // Go function name (for WASM)
    RequestType       string // Go request type
    ResponseType      string // Go response type
    RequestTSType     string // TypeScript request type
    ResponseTSType    string // TypeScript response type
    IsAsync           bool   // Async handling required
    IsServerStreaming bool   // Streaming metadata
}
```

**Design Decisions**:
- **Separate Go and TypeScript type fields** (GoType vs JSName)
- **Rich metadata** for different generation contexts
- **Comment preservation** from original protobuf definitions
- **Cross-package support** with import management

### GoDataBuilder (`go_data_builder.go`)

**Purpose**: Build template data specifically for Go WASM generation

**Key Interface**:
```go
type GoDataBuilder struct {
    analyzer      *core.ProtoAnalyzer
    pathCalc      *core.PathCalculator
    nameConv      *core.NameConverter
    serviceFilter *filters.ServiceFilter
    methodFilter  *filters.MethodFilter
}

// Main building method
func BuildTemplateData(
    packageInfo *PackageInfo,
    allBrowserServices []*protogen.Service,
    criteria *filters.FilterCriteria,
    config *GenerationConfig,
) (*GoTemplateData, error)
```

**Output Structure**:
```go
type GoTemplateData struct {
    // Package metadata
    PackageName string
    GoPackage   string
    ModuleName  string
    
    // Service data
    Services        []ServiceData // Regular services
    BrowserServices []ServiceData // Browser-provided services
    
    // JavaScript API configuration
    JSNamespace     string
    APIStructure    string
    
    // Import management
    Imports         []ImportInfo
    PackageMap      map[string]string
}
```

**Design Decisions**:
- **Separate regular and browser services** for different template handling
- **Import management** with automatic alias generation
- **JavaScript API configuration** for different structuring options
- **Cross-package browser service support** (services from other packages)

### TSDataBuilder (`ts_data_builder.go`)

**Purpose**: Build template data specifically for TypeScript generation

**Key Interface**:
```go
type TSDataBuilder struct {
    // Similar dependencies to GoDataBuilder
    // Plus message and enum collectors
    messageCollector *filters.MessageCollector
    enumCollector    *filters.EnumCollector
}

// Separate building methods for different TypeScript artifacts
func BuildClientData(...) (*TSTemplateData, error)  // For client generation
func BuildTypeData(...) (*TSTemplateData, error)    // For type generation
```

**Output Structure**:
```go
type TSTemplateData struct {
    // Package metadata
    PackageName string
    PackagePath string
    ModuleName  string
    
    // Content data (varies by generation type)
    Services []ServiceData        // For client generation
    Messages []filters.MessageInfo // For type generation
    Enums    []filters.EnumInfo    // For type generation
    
    // Import management
    ImportBasePath  string
    ExternalImports []ExternalImport
}
```

**Design Decisions**:
- **Separate methods** for different TypeScript artifacts (clients vs types)
- **Flexible content** - same data structure used for different purposes
- **Cross-package imports** for type dependencies
- **Package-based organization** for directory structure

### File Planning (`file_planning.go`)

**Purpose**: Plan what files to generate and manage protogen GeneratedFile objects

**Key Types**:
```go
// File specification
type FileSpec struct {
    Name         string       // Logical name ("wasm", "client", "interfaces")
    Filename     string       // Physical filename ("library/v1/library_v1.wasm.go")
    Type         string       // File type for template selection
    Required     bool         // Whether file is required
    ContentHints ContentHints // Metadata about file content
}

// Complete file plan for a package
type FilePlan struct {
    PackageName string
    Specs       []FileSpec
    Config      *GenerationConfig
}

// Collection of protogen files ready for rendering
type GeneratedFileSet struct {
    Files map[string]*protogen.GeneratedFile // Maps logical name to protogen file
    Plan  *FilePlan                          // Original plan
}
```

**Design Decisions**:
- **Logical vs physical names** - separates template concerns from file paths
- **Content hints** guide conditional generation decisions
- **Required vs optional** files for validation
- **Batch file creation** - all protogen files created at once

## Data Flow

### Go WASM Generation Flow
```
Proto Files → ServiceFilter → GoDataBuilder → GoTemplateData
                    ↓
            FilePlan → GeneratedFileSet → GoRenderer → .wasm.go files
```

### TypeScript Generation Flow
```
Proto Files → MessageCollector → TSDataBuilder → TSTemplateData
                    ↓                    ↓
            EnumCollector →        FilePlan → GeneratedFileSet → TSRenderer → .ts files
```

## Usage Examples

### Go Data Building
```go
// Create builder with dependencies
goBuilder := builders.NewGoDataBuilder(analyzer, pathCalc, nameConv, serviceFilter, methodFilter)

// Build template data
templateData, err := goBuilder.BuildTemplateData(packageInfo, browserServices, criteria, config)

// Use in templates
// templateData.Services contains filtered services
// templateData.Imports contains required imports
// templateData.JSNamespace contains calculated namespace
```

### TypeScript Data Building
```go
// Create builder
tsBuilder := builders.NewTSDataBuilder(analyzer, pathCalc, nameConv, serviceFilter, methodFilter, msgCollector, enumCollector)

// Build client data
clientData, err := tsBuilder.BuildClientData(packageInfo, criteria, config)

// Build type data
typeData, err := tsBuilder.BuildTypeData(packageInfo, criteria, config)
```

### File Planning
```go
// Plan files
plan := &builders.FilePlan{
    PackageName: "library.v1",
    Specs: []builders.FileSpec{
        {Name: "wasm", Filename: "library_v1.wasm.go", Type: "wasm", Required: true},
        {Name: "main", Filename: "main.go.example", Type: "example", Required: true},
    },
}

// Create protogen files
fileSet := builders.NewGeneratedFileSet(plan, plugin)

// Use in rendering
wasmFile := fileSet.GetFile("wasm")
renderer.RenderToFile(wasmFile, template, data)
```

## Testing Strategy

### Builder Testing
- **Template data validation** ensures all required fields are present
- **Import management** tests verify correct package alias generation
- **Cross-package support** tests browser service handling

### File Planning Testing  
- **File specification** creation and validation
- **File set operations** (GetFilesByType, GetRequiredFiles, etc.)
- **Content hints** for conditional generation logic

### Integration Testing
- **End-to-end workflows** from proto files to template data
- **Complex filtering scenarios** with multiple criteria types
- **Error handling** for invalid configurations

## Extension Points

### Adding New Data Builders
```go
// Add specialized builder for new generation target
type WebDataBuilder struct {
    // Dependencies for web-specific generation
}

func (wb *WebDataBuilder) BuildWebData(...) (*WebTemplateData, error) {
    // Build data for web-specific templates
}
```

### Adding New Template Data
```go
// Extend template data with new fields
type GoTemplateData struct {
    // ... existing fields
    
    // New fields for extended functionality
    CustomDirectives []string    // Custom template directives
    FeatureFlags     FeatureFlags // Feature toggle support
}
```

### Adding New File Types
```go
// Add new file specifications
FileSpec{
    Name: "middleware",
    Type: "middleware", 
    ContentHints: ContentHints{
        HasMiddleware: true,
    },
}
```

## Architecture Benefits

### Separation of Concerns
- **Data building** separate from template execution
- **File planning** separate from file creation
- **Language-specific logic** isolated in dedicated builders

### Testability
- **Rich data structures** easy to validate and test
- **Pure transformation logic** with predictable outputs
- **Mock-free testing** for most builder logic

### Flexibility
- **Multiple template data types** for different generation needs
- **Configurable file organization** through file planning
- **Cross-package support** for complex proto structures

### Maintainability  
- **Clear interfaces** between filtering and rendering
- **Comprehensive validation** catches issues early
- **Consistent patterns** across all builders

## Future Considerations

### Performance
- **Lazy evaluation** for expensive data building operations
- **Caching** for repeated calculations (import paths, etc.)
- **Parallel building** for independent packages

### Features
- **Template data versioning** for backward compatibility
- **Conditional data building** based on target platform
- **Enhanced cross-package dependency tracking**

### Testing
- **Golden file testing** for template data structures
- **Property-based testing** for data transformation logic
- **Performance benchmarks** for large proto files
