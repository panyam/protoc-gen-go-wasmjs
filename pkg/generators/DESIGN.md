# Generator Layer Design

## Purpose

The generator layer provides **top-level orchestration** of the complete code generation process. This is Layer 5 of the architecture - it coordinates all lower layers (core, filters, builders, renderers) to produce the final generated code from protobuf definitions.

## Design Principles

### 1. **Single Responsibility Generators**
- **GoGenerator**: Only generates Go WASM artifacts
- **TSGenerator**: Only generates TypeScript artifacts  
- Each generator is focused and optimized for its target language

### 2. **Complete File Control**
- Generators make ALL file creation and naming decisions
- File planning separates logical organization from physical paths
- Flexible file organization strategies (platform-specific, content-based, etc.)

### 3. **Layered Orchestration**
- Coordinates all lower layers through clean interfaces
- Handles error propagation and validation
- Provides generation statistics and reporting

## Architecture Overview

```
Generator (Layer 5) - Top-level orchestration
    ↓
Builder (Layer 3) - Template data transformation  
    ↓
Filter (Layer 2) - Business logic decisions
    ↓  
Core (Layer 1) - Pure utility functions

Generator also coordinates:
    → Renderer (Layer 4) - Template execution
    → protogen.Plugin - File creation and management
```

## Key Components

### GoGenerator (`go_generator.go`)

**Purpose**: Orchestrate complete Go WASM generation process

**Key Interface**:
```go
type GoGenerator struct {
    // All layer dependencies injected
    analyzer      *core.ProtoAnalyzer
    pathCalc      *core.PathCalculator  
    nameConv      *core.NameConverter
    packageFilter *filters.PackageFilter
    serviceFilter *filters.ServiceFilter
    methodFilter  *filters.MethodFilter
    dataBuilder   *builders.GoDataBuilder
    renderer      *renderers.GoRenderer
    plugin        *protogen.Plugin
}

// Main generation entry point
func (gg *GoGenerator) Generate(config *GenerationConfig, filterCriteria *FilterCriteria) error
```

**Generation Pipeline**:
1. **Package Filtering** - Determine which packages to process
2. **Browser Service Collection** - Gather browser services from all packages
3. **Per-Package Generation**:
   - Build template data using filtered proto data
   - Plan file generation strategy
   - Create protogen GeneratedFile objects
   - Execute templates through renderer

**File Planning Strategy**:
```go
// GoGenerator decides what files to generate
func (gg *GoGenerator) planGoFiles(data *GoTemplateData, config *GenerationConfig) *FilePlan {
    var specs []FileSpec
    
    // Always: WASM wrapper (main artifact)
    specs = append(specs, FileSpec{
        Name: "wasm", 
        Filename: "library/v1/library_v1.wasm.go",
        Type: "wasm",
        Required: true,
    })
    
    // Always: Usage example
    specs = append(specs, FileSpec{
        Name: "main",
        Filename: "library/v1/main.go.example", 
        Type: "example",
        Required: true,
    })
    
    // Conditional: Build script
    if config.GenerateBuildScript {
        specs = append(specs, FileSpec{
            Name: "build",
            Filename: "build.sh",
            Type: "script", 
            Required: false,
        })
    }
    
    return &FilePlan{PackageName: data.PackageName, Specs: specs}
}
```

### TSGenerator (`ts_generator.go`)

**Purpose**: Orchestrate complete TypeScript generation process

**Key Interface**:
```go
type TSGenerator struct {
    // All layer dependencies injected
    // Plus message and enum collectors for type generation
    msgCollector  *filters.MessageCollector
    enumCollector *filters.EnumCollector
    dataBuilder   *builders.TSDataBuilder
    renderer      *renderers.TSRenderer
}

// Main generation entry point  
func (tg *TSGenerator) Generate(config *GenerationConfig, filterCriteria *FilterCriteria) error
```

**Generation Pipeline**:
1. **Package Filtering** - Determine which packages to process
2. **Per-Package Generation**:
   - **Client Generation** (if package has services)
   - **Type Generation** (if package has messages/enums)
   - File planning for TypeScript artifacts
   - Template execution through renderer

**File Planning Strategy**:
```go
// TSGenerator plans multiple file types
func (tg *TSGenerator) planTSFiles(packageInfo *PackageInfo, criteria *FilterCriteria, config *GenerationConfig) *FilePlan {
    var specs []FileSpec
    
    // Client file (if services exist)
    if hasServices {
        specs = append(specs, FileSpec{
            Name: "client",
            Filename: "LibraryClient.ts",
            Type: "client",
            Required: true,
        })
    }
    
    // Type files (if messages/enums exist)
    if hasTypes {
        specs = append(specs, []FileSpec{
            {Name: "interfaces", Filename: "library/v1/interfaces.ts", Type: "interfaces"},
            {Name: "models", Filename: "library/v1/models.ts", Type: "models"},
            {Name: "factory", Filename: "library/v1/factory.ts", Type: "factory"},
            {Name: "schemas", Filename: "library/v1/schemas.ts", Type: "schemas"},
            {Name: "deserializer", Filename: "library/v1/deserializer.ts", Type: "deserializer"},
        }...)
    }
    
    return &FilePlan{PackageName: packageInfo.Name, Specs: specs}
}
```

## File Management Architecture

### Generator-Controlled File Creation
```go
// Generator creates ALL protogen files upfront
func (g *Generator) generatePackageFiles(data *TemplateData, config *Config) error {
    // 1. Plan files based on content and configuration
    filePlan := g.planFiles(data, config)
    
    // 2. Create all GeneratedFile objects
    fileSet := NewGeneratedFileSet(filePlan, g.plugin)
    
    // 3. Validate file plan
    if err := fileSet.ValidateFileSet(); err != nil {
        return err
    }
    
    // 4. Render each file
    for _, spec := range filePlan.Specs {
        file := fileSet.GetFile(spec.Name)
        template := g.getTemplate(spec.Type)
        renderData := g.buildDataForFile(spec, data)
        
        renderer.RenderToFile(file, template, renderData)
    }
}
```

### Benefits of This Approach
- **No filename conflicts** - protogen files created once
- **Flexible file strategies** - generator can implement any organization
- **Reusable files** - multiple templates can write to same file
- **Conditional generation** - files created only when needed
- **Complete control** - generator decides everything about files

## Configuration Integration

### Command-Line to Generator Bridge
```go
// Binary (cmd/) parses flags
var services = flagSet.String("services", "", "Services to generate")

// Generator receives structured configuration
config := &GenerationConfig{
    WasmExportPath: *wasmExportPath,
    JSStructure:    *jsStructure,
}

filterCriteria, err := filters.ParseFromConfig(*services, *includes, *excludes, *renames)

// Generator orchestrates generation
generator.Generate(config, filterCriteria)
```

### Configuration Validation
```go
// Generators validate their specific configuration needs
func (gg *GoGenerator) ValidateConfig(config *GenerationConfig) error {
    if config.WasmExportPath == "" {
        return fmt.Errorf("WasmExportPath cannot be empty")
    }
    
    validStructures := map[string]bool{"namespaced": true, "flat": true, "service_based": true}
    if !validStructures[config.JSStructure] {
        return fmt.Errorf("invalid JSStructure: %s", config.JSStructure)
    }
    
    return nil
}
```

## Usage Examples

### Basic Generation
```go
// Create generator
generator := generators.NewGoGenerator(plugin)

// Configure generation
config := &builders.GenerationConfig{
    WasmExportPath: "./gen/wasm",
    JSStructure:    "namespaced",
}

criteria, _ := filters.ParseFromConfig("UserService", "", "*Internal", "")

// Generate
err := generator.Generate(config, criteria)
```

### Advanced File Organization
```go
// Custom filename strategy
func (gg *GoGenerator) calculateWasmFilename(packageName string, config *GenerationConfig) string {
    switch config.Target {
    case "browser":
        return "js/" + packageName + "_wasm.js"
    case "node": 
        return "lib/" + packageName + ".wasm.go"
    default:
        return gg.pathCalc.BuildPackagePath(packageName) + "/" + packageName + ".wasm.go"
    }
}
```

### Multiple Templates Per File
```go
// Combine multiple template sections
func (gg *GoGenerator) renderCombinedWasmFile(file *protogen.GeneratedFile, data *GoTemplateData) error {
    // Header section
    headerTemplate := gg.getHeaderTemplate()
    renderer.RenderToFile(file, headerTemplate, data.HeaderData)
    
    // Service methods section
    methodsTemplate := gg.getMethodsTemplate() 
    renderer.RenderToFile(file, methodsTemplate, data.MethodsData)
    
    // Footer section
    footerTemplate := gg.getFooterTemplate()
    renderer.RenderToFile(file, footerTemplate, data.FooterData)
}
```

## Testing Strategy

### Generator Testing
- **Component initialization** ensures all dependencies are wired correctly
- **Configuration validation** tests various config scenarios
- **File planning** tests verify correct file specifications
- **Integration tests** validate layer coordination

### Error Handling Testing
- **Invalid configuration** scenarios
- **Missing template data** scenarios  
- **File creation failures** (though limited without real protogen)

## Extension Points

### Adding New Generators
```go
// Add generator for new target
type WebGenerator struct {
    // Similar structure to GoGenerator/TSGenerator
    analyzer    *core.ProtoAnalyzer
    dataBuilder *builders.WebDataBuilder
    renderer    *renderers.WebRenderer
}

func (wg *WebGenerator) Generate(config *WebConfig, criteria *FilterCriteria) error {
    // Implement web-specific generation pipeline
}
```

### Adding New File Types
```go
// Extend file planning with new types
FileSpec{
    Name: "middleware",
    Type: "middleware",
    ContentHints: ContentHints{
        HasMiddleware: true,
    },
}
```

### Adding New Template Sources
```go
// Support different template sources
func (g *Generator) getTemplate(fileType string) string {
    switch g.config.TemplateSource {
    case "embedded":
        return g.getEmbeddedTemplate(fileType)
    case "files":
        return g.loadTemplateFile(fileType + ".tmpl")
    case "remote":
        return g.fetchRemoteTemplate(fileType)
    }
}
```

## Dependencies

### Internal Dependencies
- `pkg/core` - Proto analysis, path calculations, name conversions
- `pkg/filters` - Business logic filtering and collection
- `pkg/builders` - Template data building and file planning
- `pkg/renderers` - Template execution

### External Dependencies
- `google.golang.org/protobuf/compiler/protogen` - Protobuf compiler integration
- Standard library only

## Future Considerations

### Performance
- **Parallel package processing** for independent packages
- **Lazy template loading** for large template sets
- **Incremental generation** for large proto files

### Features
- **Custom template sources** (files, remote, etc.)
- **Generation middleware** for custom processing steps
- **Plugin composition** for combining generators

### Reliability
- **Graceful degradation** when optional files fail
- **Partial generation** support for large projects
- **Recovery strategies** for template execution failures
