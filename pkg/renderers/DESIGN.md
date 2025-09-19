# Renderer Layer Design

## Purpose

The renderer layer provides **pure template execution** without file management concerns. This is Layer 4 of the architecture - it takes template data from builders and template content from generators, then executes templates to produce the final generated code.

## Design Principles

### 1. **Pure Template Execution**
- Renderers only execute templates - no file management
- Generators control all file creation and naming decisions
- Clean separation between template logic and file organization

### 2. **GeneratedFile-Based Interface**
- Accept `protogen.GeneratedFile` objects instead of filenames
- Enable multiple templates per file and file reuse
- Support advanced file composition patterns

### 3. **Shared Template System**
- Common template helper functions across all renderers
- Consistent template execution with error handling
- Template validation and debugging support

## Key Components

### Template Helpers (`template_helpers.go`)

**Purpose**: Shared template functions and execution utilities

**Key Functions**:
```go
// Template function map with common helpers
func getTemplateFuncMap() template.FuncMap {
    return template.FuncMap{
        // String manipulation
        "title":      func(s string) string // PascalCase
        "lower":      func(s string) string // lowercase
        "replaceAll": func(old, new, s string) string
        
        // Conditionals
        "eq":  func(a, b interface{}) bool
        "and": func(a, b bool) bool
        "or":  func(a, b bool) bool
        
        // Array operations
        "join":  func(sep string, elems []string) string
        "split": func(sep, s string) []string
    }
}

// Consistent template execution
func ExecuteTemplate(templateName, templateContent string, data interface{}, output Writer) error

// Template validation
func ValidateTemplate(templateName, templateContent string) error
```

**Design Decisions**:
- **Rich function map** covers common template operations
- **Error handling** with descriptive messages
- **Template validation** for development and testing
- **Reusable across languages** (Go and TypeScript templates)

### GoRenderer (`go_renderer.go`)

**Purpose**: Execute Go templates with validation and error handling

**Core Interface**:
```go
type GoRenderer struct {
    // No dependencies - pure template executor
}

// Core method - all others use this
func (gr *GoRenderer) RenderToFile(
    file *protogen.GeneratedFile, 
    templateContent string, 
    data interface{}
) error

// Convenience methods with validation
func (gr *GoRenderer) RenderWasmWrapper(file *protogen.GeneratedFile, data *GoTemplateData, templateContent string) error
func (gr *GoRenderer) RenderMainExample(file *protogen.GeneratedFile, data *GoTemplateData, templateContent string) error
func (gr *GoRenderer) RenderBuildScript(file *protogen.GeneratedFile, data *GoTemplateData, templateContent string) error
```

**Design Decisions**:
- **Single core method** (`RenderToFile`) handles all template execution
- **Convenience methods** add data validation for specific template types
- **No file management** - accepts `GeneratedFile` objects from generator
- **Type-specific validation** ensures data integrity before rendering

### TSRenderer (`ts_renderer.go`)

**Purpose**: Execute TypeScript templates with validation and error handling

**Core Interface**:
```go
type TSRenderer struct {
    // No dependencies - pure template executor
}

// Core method
func (tr *TSRenderer) RenderToFile(file *protogen.GeneratedFile, templateContent string, data interface{}) error

// TypeScript-specific convenience methods
func (tr *TSRenderer) RenderClient(file *protogen.GeneratedFile, data *TSTemplateData, templateContent string) error
func (tr *TSRenderer) RenderInterfaces(file *protogen.GeneratedFile, data *TSTemplateData, templateContent string) error
func (tr *TSRenderer) RenderModels(file *protogen.GeneratedFile, data *TSTemplateData, templateContent string) error
func (tr *TSRenderer) RenderFactory(file *protogen.GeneratedFile, data *TSTemplateData, templateContent string) error
func (tr *TSRenderer) RenderSchemas(file *protogen.GeneratedFile, data *TSTemplateData, templateContent string) error
func (tr *TSRenderer) RenderDeserializer(file *protogen.GeneratedFile, data *TSTemplateData, templateContent string) error
```

**Design Decisions**:
- **Multiple TypeScript artifacts** supported through different methods
- **Same core interface** as GoRenderer for consistency
- **Rich validation** ensures TypeScript-specific data requirements
- **Flexible template support** for different TypeScript file types

## Architecture Benefits

### Clean Separation of Concerns

**Before (Old Architecture)**:
```go
// Mixed concerns - bad
func generateWasmWrapper(data *TemplateData) error {
    filename := calculateFileName(data)        // File management
    file := plugin.NewGeneratedFile(filename)  // Protogen integration
    template := parseTemplate(content)         // Template parsing
    return template.Execute(file, data)        // Template execution
}
```

**After (New Architecture)**:
```go
// Generator: Controls file decisions
wasmFile := plugin.NewGeneratedFile("library_v1.wasm.go", "")

// Renderer: Pure template execution
renderer.RenderToFile(wasmFile, templateContent, data)
```

### Advanced File Patterns Enabled

#### 1. **Multiple Templates → Single File**
```go
// Generator creates one file, renders multiple template sections
combinedFile := plugin.NewGeneratedFile("service.go", "")

// Build up file with multiple template executions
renderer.RenderToFile(combinedFile, headerTemplate, headerData)
renderer.RenderToFile(combinedFile, methodsTemplate, methodsData)
renderer.RenderToFile(combinedFile, footerTemplate, footerData)
```

#### 2. **Conditional File Creation**
```go
// Generator decides what files to create based on content
var files []*protogen.GeneratedFile

if hasServices {
    files = append(files, plugin.NewGeneratedFile("services.ts", ""))
}
if hasMessages {
    files = append(files, plugin.NewGeneratedFile("types.ts", ""))
}

// Render only the files we need
for _, file := range files {
    renderer.RenderToFile(file, template, data)
}
```

#### 3. **Platform-Specific File Organization**
```go
// Generator can implement different file strategies
switch config.Target {
case "browser":
    clientFile := plugin.NewGeneratedFile("client/library.js", "")
case "node":
    serverFile := plugin.NewGeneratedFile("lib/library.mjs", "")
case "wasm":
    wasmFile := plugin.NewGeneratedFile("wasm/library.wasm.go", "")
}
```

## Template Integration

### Current Template System Integration
The renderer layer is designed to integrate with the existing template system:

```go
// Templates come from various sources
func (gg *GoGenerator) getWasmTemplate() string {
    // Option 1: Embedded templates (current approach)
    return embeddedWasmTemplate
    
    // Option 2: Custom templates (user-provided)
    if customTemplate := loadCustomTemplate(config.WasmTemplate); customTemplate != "" {
        return customTemplate
    }
    
    // Option 3: External template files
    return loadTemplateFile("templates/wasm.go.tmpl")
}

// Renderer executes whatever template content it receives
renderer.RenderWasmWrapper(wasmFile, data, template)
```

### Template Function Integration
```go
// Template content can use all helper functions
const wasmTemplate = `
//go:build js && wasm
package {{ .ModuleName }}

// {{ .PackageName | replaceAll "." "_" | title }}ServicesExports
type {{ .PackageName | replaceAll "." "_" | title }}ServicesExports struct {
{{- range .Services }}
    {{ .Name }} {{ .GoType }}
{{- end }}
}
`
```

## Error Handling

### Template Execution Errors
```go
// Detailed error messages with context
func (gr *GoRenderer) RenderToFile(file *protogen.GeneratedFile, templateContent string, data interface{}) error {
    if file == nil {
        return fmt.Errorf("GeneratedFile cannot be nil")
    }
    
    if templateContent == "" {
        return fmt.Errorf("template content cannot be empty") 
    }
    
    // ExecuteTemplate provides detailed parsing and execution errors
    return ExecuteTemplate("go", templateContent, data, file)
}
```

### Data Validation
```go
// Pre-rendering validation catches issues early
func (gr *GoRenderer) ValidateGoTemplateData(data *GoTemplateData) error {
    if data.PackageName == "" {
        return fmt.Errorf("PackageName cannot be empty")
    }
    
    for i, service := range data.Services {
        if len(service.Methods) == 0 {
            return fmt.Errorf("service %s at index %d has no methods", service.Name, i)
        }
    }
    
    return nil
}
```

## Testing Strategy

### Unit Testing
- **Template execution** with mock data structures
- **Data validation** with valid and invalid inputs
- **Error handling** for various failure scenarios

### Integration Testing
- **File planning** integration with real protogen objects
- **Template rendering** with actual template content
- **Multi-file scenarios** testing file reuse and composition

## Dependencies

### Internal Dependencies
- `pkg/builders` - For template data types
- None for core renderer functionality

### External Dependencies
- `google.golang.org/protobuf/compiler/protogen` - For GeneratedFile objects
- `text/template` - For template execution
- Standard library only

## Usage Guidelines

### For Template Authors
```go
// Templates receive rich data structures
{{- range .Services }}
Service: {{ .Name }}
JS Name: {{ .JSName }}
{{- range .Methods }}
  Method: {{ .Name }} → {{ .JSName }}
  Async: {{ .IsAsync }}
{{- end }}
{{- end }}
```

### For Generator Authors
```go
// Generators control all file decisions
func (g *Generator) Generate() {
    // 1. Plan files
    plan := g.planFiles(data, config)
    
    // 2. Create protogen files
    fileSet := NewGeneratedFileSet(plan, plugin)
    
    // 3. Render with complete control
    for _, spec := range plan.Specs {
        file := fileSet.GetFile(spec.Name)
        template := g.getTemplate(spec.Type)
        data := g.buildData(spec, packageData)
        
        renderer.RenderToFile(file, template, data)
    }
}
```

### For Testing
```go
// Easy to test without complex protogen setup
func TestRenderer_RenderToFile(t *testing.T) {
    var buf bytes.Buffer
    data := &TestData{Name: "test"}
    template := "Hello {{ .Name }}"
    
    err := ExecuteTemplate("test", template, data, &buf)
    
    assert.NoError(t, err)
    assert.Equal(t, "Hello test", buf.String())
}
```
