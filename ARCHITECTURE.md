# protoc-gen-go-wasmjs Architecture

## Overview

protoc-gen-go-wasmjs follows a plugin architecture that integrates with the Protocol Buffers compiler toolchain. It generates two primary artifacts: Go WASM wrappers and TypeScript clients, enabling seamless communication between JavaScript environments and Go service implementations compiled to WebAssembly.

## High-Level Architecture

```
┌─────────────────┐     ┌──────────────────┐     ┌─────────────────┐
│   .proto files  │────▶│     protoc       │────▶│  Generated Code │
└─────────────────┘     └──────────────────┘     └─────────────────┘
                               │                           │
                               ▼                           ▼
                     ┌──────────────────┐         ┌─────────────────┐
                     │ protoc-gen-go-   │         │ • WASM wrapper  │
                     │     wasmjs        │         │ • TS client     │
                     └──────────────────┘         │ • Build script  │
                                                  └─────────────────┘
```

## Core Components

### 1. Plugin Entry Point (`cmd/protoc-gen-go-wasmjs/main.go`)
- Receives CodeGeneratorRequest from protoc
- Delegates to generator package
- Returns CodeGeneratorResponse with generated files

### 2. Generator (`pkg/generator/generator.go`)
The main orchestrator that:
- Parses configuration options
- Analyzes proto files and services
- Groups files by package
- Applies filtering and transformation rules
- Executes templates

### 3. Configuration System (`pkg/generator/config.go`)
Comprehensive configuration parsing:
```go
type GeneratorConfig struct {
    // TypeScript Integration
    TSGenerator      string   // protoc-gen-es, protoc-gen-ts
    TSImportPath     string   // Where to import TS types from
    TSImportExtension string  // js, ts, none, or auto-detect
    
    // Service Selection
    Services         []string // Specific services to generate
    MethodInclude    []string // Glob patterns to include
    MethodExclude    []string // Glob patterns to exclude
    
    // JavaScript API
    JSStructure      string   // namespaced, flat, service_based
    JSNamespace      string   // Global namespace name
    ModuleName       string   // WASM module name
    
    // Generation Control
    GenerateWASM     bool     // Generate Go WASM wrapper
    GenerateTypeScript bool   // Generate TS client
}
```

### 4. Template System (`pkg/generator/templates/`)
Embedded templates using Go's `embed` package:
- `wasm.go.tmpl` - Go WASM wrapper generation
- `client.ts.tmpl` - TypeScript client generation
- `build.sh.tmpl` - Build script generation
- `main.go.tmpl` - Example usage generation

### 5. Type System (`pkg/generator/types.go`)
Data structures passed to templates:
```go
type TemplateData struct {
    Services    []ServiceData
    Config      GeneratorConfig
    JSNamespace string
    ModuleName  string
    APIStructure string
    TSImports   []TSImportGroup  // Smart import grouping
}
```

## Key Architectural Patterns

### 1. Export Pattern for Dependency Injection
Instead of generating fixed `main()` functions, we generate export structs:
```go
type ServicesExports struct {
    UserService  UserServiceServer
    OrderService OrderServiceServer
}

func (exports *ServicesExports) RegisterAPI() {
    // Registers JavaScript APIs
}
```

This allows users to inject their own implementations with full control over dependencies.

### 2. Smart Import Detection
Analyzes proto file sources to generate accurate TypeScript imports:
```go
// For each method, determine source proto file
sourceFile := method.Input.Desc.ParentFile().Path()
// Group imports by source file
imports[sourceFile] = append(imports[sourceFile], typeName)
```

### 3. Proto to JSON Conversion System

#### Client-Side Architecture
The TypeScript client implements a sophisticated conversion pipeline:

```typescript
class Client {
    private conversionOptions: ConversionOptions = {
        handleOneofs: true,
        emitDefaults: false,
        fieldTransformer?: (field: string) => string,
        bigIntHandler?: (value: bigint) => string
    };
    
    callMethod(method: string, request: any): Promise<any> {
        // 1. Convert request to JSON with custom handling
        const json = this.convertToJson(request);
        
        // 2. Call WASM method
        const response = wasmMethod(JSON.stringify(json));
        
        // 3. Convert response from JSON
        return this.convertFromJson(response.data);
    }
}
```

#### Conversion Pipeline
1. **Native Method Detection**: Check for toJson/toJSON methods
2. **Custom Conversions**: Apply field transformations
3. **Oneof Handling**: Flatten oneof structures if configured
4. **BigInt Serialization**: Handle BigInt values properly
5. **Default Value Management**: Control emission of defaults

#### WASM-Side Configuration
Go WASM wrapper uses protojson with specific options:
```go
// Unmarshal with flexibility
protojson.UnmarshalOptions{
    DiscardUnknown: true,
    AllowPartial:   true,
}

// Marshal for TypeScript compatibility
protojson.MarshalOptions{
    UseProtoNames:   false,  // Use JSON names
    EmitUnpopulated: false,  // Skip defaults
    UseEnumNumbers:  false,  // Use strings
}
```

### 4. Multi-Target Generation
Supports generating different combinations of services for different use cases:
```yaml
# User page - only UserService
- local: protoc-gen-go-wasmjs
  out: ./gen/wasm/user-page
  opt: [services=UserService]

# Admin page - all services  
- local: protoc-gen-go-wasmjs
  out: ./gen/wasm/admin-page
  opt: []  # All services
```

### 5. Dual-Target Architecture
WASM and TypeScript artifacts can be generated independently:
```yaml
# Generate only WASM
- local: protoc-gen-go-wasmjs
  out: ./gen/wasm
  opt: [generate_typescript=false]

# Generate only TypeScript
- local: protoc-gen-go-wasmjs  
  out: ./frontend/src/clients
  opt: [generate_wasm=false]
```

## Data Flow

### 1. Generation Flow
```
Proto Files → protoc → Plugin → Templates → Generated Files
                ↓
         Configuration
```

### 2. Runtime Flow
```
TypeScript Client → JSON Request → WASM Method → Go Service
                                        ↓
                                   JSON Response
                                        ↓
                                 TypeScript Client
```

### 3. Conversion Flow
```
TS Object → toJson() → Custom Conversions → JSON → WASM
                              ↓
                    • Oneof flattening
                    • Field transformation  
                    • BigInt handling
                    • Default filtering
```

## Design Principles

### 1. **Zero Configuration**
Works out of the box with sensible defaults while allowing extensive customization.

### 2. **Generator Agnostic**
Detects and works with any TypeScript protobuf generator through convention-based detection.

### 3. **Flexible Deployment**
Supports various deployment patterns from monolithic to micro-frontends.

### 4. **Type Safety**
Maintains full type safety from proto definitions through to TypeScript usage.

### 5. **Performance First**
Minimal overhead in the conversion layer with optimized JSON handling.

## Extension Points

### 1. Template Customization
Users can provide custom templates via `template_dir` option.

### 2. Service Filtering
Fine-grained control over which services and methods to generate.

### 3. Conversion Middleware
The conversion system can be extended with custom transformers.

### 4. Build Integration
Generated build scripts can be customized for specific environments.

## Security Considerations

### 1. Input Validation
All JSON parsing includes error handling and validation.

### 2. Timeout Protection
WASM methods include context timeouts to prevent hanging.

### 3. Error Isolation
Errors in one service call don't affect others.

### 4. No Direct Memory Access
JavaScript and WASM communicate only through JSON, ensuring memory safety.

## Performance Characteristics

### 1. Startup Cost
- One-time WASM module loading (~1-5MB depending on services)
- TypeScript client initialization is negligible

### 2. Per-Call Overhead
- JSON serialization/deserialization
- Proto conversion transformations
- Minimal compared to network calls

### 3. Memory Usage
- WASM module remains in memory
- No per-call memory leaks
- Garbage collection handled by browser

## Future Architecture Considerations

### 1. Streaming Support
Current architecture uses request/response pattern. Streaming would require:
- WebSocket or EventSource integration
- Stream-aware TypeScript clients
- Go-side streaming handlers

### 2. Shared Memory
Could leverage SharedArrayBuffer for better performance in supported browsers.

### 3. Module Federation
Architecture supports splitting into multiple WASM modules for micro-frontend patterns.

### 4. Web Workers
WASM execution could be moved to Web Workers for non-blocking operations.