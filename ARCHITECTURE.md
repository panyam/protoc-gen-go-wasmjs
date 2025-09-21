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
type Config struct {
    // Generation Control
    GenerateWasm       bool   // Generate WASM wrapper (default: true)
    GenerateTypeScript bool   // Generate TypeScript client (default: true)
    WasmExportPath     string // Path where WASM wrapper should be generated
    TSExportPath       string // Path where TypeScript files should be generated
    
    // Service Selection
    Services         []string // Specific services to generate
    MethodInclude    []string // Glob patterns to include
    MethodExclude    []string // Glob patterns to exclude
    MethodRenames    map[string]string // Method name transformations
    
    // JavaScript API
    JSStructure      string   // namespaced, flat, service_based
    JSNamespace      string   // Global namespace name
    ModuleName       string   // WASM module name
    
    // Build Integration
    WasmPackageSuffix   string // Package suffix for WASM wrapper
    GenerateBuildScript bool   // Generate build script for WASM compilation
}
```

### 4. Template System (`pkg/generator/templates/`)
Embedded templates using Go's `embed` package:
- `wasm.go.tmpl` - Go WASM wrapper generation
- `client_simple.ts.tmpl` - Simplified TypeScript client generation
- `interfaces.ts.tmpl` - TypeScript interface generation
- `models.ts.tmpl` - TypeScript model class generation
- `factory.ts.tmpl` - TypeScript factory generation
- `build.sh.tmpl` - Build script generation
- `main.go.tmpl` - Example usage generation

### 5. TypeScript Generation (`pkg/generator/tsgenerator.go`)
Dedicated TypeScript generation logic:
- Proto message analysis and field type conversion
- Interface and model class generation
- Enhanced factory pattern with cross-package composition
- Schema generation with field metadata and proto field IDs
- Deserializer generation with factory integration
- Package-based nested directory structure
- Cross-package dependency detection and import management
- **Comprehensive Enum Support**: Complete enum collection, generation, and import system
- **wasmjs.v1 Package Filtering**: Intelligent filtering to exclude annotation packages from artifact generation

### 6. Type System (`pkg/generator/types.go`)
Data structures passed to templates and message analysis:
```go
type TemplateData struct {
    Services    []ServiceData
    Config      *Config
    JSNamespace string
    ModuleName  string
    APIStructure string
    Imports     []ImportInfo      // Unique package imports with aliases
    PackageMap  map[string]string // Maps full package path to alias
}

type MessageInfo struct {
    Name         string      // Message name (e.g., "Book")
    TSName       string      // TypeScript interface name (e.g., "Book")
    Fields       []FieldInfo // All fields in the message
    PackageName  string      // Proto package name
    IsNested     bool        // Whether this is a nested message
    Comment      string      // Leading comment from proto
}

type FieldInfo struct {
    Name         string    // Original proto field name
    JSONName     string    // JSON field name (camelCase)
    TSType       string    // TypeScript type
    IsRepeated   bool      // Whether this is a repeated field
    IsOptional   bool      // Whether this is an optional field
    MessageType  string    // For message fields, the message type name
    DefaultValue string    // Default value for the field
}

type EnumInfo struct {
    Name               string           // Enum name (e.g., "GameStatus")
    TSName             string           // TypeScript enum name
    Values             []EnumValueInfo  // All enum values
    PackageName        string           // Proto package name
    ProtoFile          string           // Source proto file
    Comment            string           // Leading comment from proto
    FullyQualifiedName string           // Fully qualified name with package
}

type EnumValueInfo struct {
    Name     string // Enum value name (e.g., "GAME_STATUS_UNSPECIFIED")
    TSName   string // TypeScript enum value name
    Number   int32  // Enum value number
    Comment  string // Leading comment from proto
}
```

## Key Architectural Patterns

### 1. wasmjs.v1 Package Filtering
The system intelligently filters wasmjs annotation packages from artifact generation:
```go
// In cmd/protoc-gen-go-wasmjs/main.go (lines 94-97)
packageName := strings.ReplaceAll(*protoFile.Package, ".", "")
if packageName == "wasmjs.v1" {
    continue  // Skip wasmjs annotation packages
}
```
This ensures annotation packages remain visible to the proto compiler while avoiding unwanted artifact generation.

### 2. Export Pattern for Dependency Injection
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

### 3. Comprehensive Enum Support
Complete enum collection and TypeScript generation system:
```go
// EnumInfo represents proto enums with full metadata
type EnumInfo struct {
    Name               string
    TSName             string
    Values             []EnumValueInfo
    PackageName        string
    FullyQualifiedName string
}

// collectAllEnums gathers enums from all proto files
func collectAllEnums(files []*descriptorpb.FileDescriptorProto, packageFilter string) []EnumInfo
```

Template generation includes enums in interfaces.ts:
```typescript
{{range .Enums}}export enum {{.TSName}} {
{{range .Values}}  {{.TSName}} = {{.Number}},
{{end}}}
{{end}}
```

And proper imports in all TypeScript files:
```typescript
import { MessageInterface, GameStatus } from "./interfaces";
```

### 4. Self-Generated TypeScript Architecture
Generates complete TypeScript structure directly from proto definitions:

```
For each proto package (e.g., library.v2):
├── library_interfaces.ts           // TypeScript interfaces  
├── library_models.ts               // Concrete class implementations
├── factory.ts                      // Enhanced factories with cross-package composition
├── library_schemas.ts              // Schema definitions with field metadata
└── library_deserializer.ts         // Schema-aware deserializers
```

#### Enhanced TypeScript Generation Structure
```typescript
// 1. Interfaces for flexibility and type safety
export interface Book {
  base?: BaseMessage;  // Cross-package reference
  title: string;
  author: string;
  tags?: string[];     // Optional repeated field
  available: boolean;
}

// Generated enums with proper TypeScript syntax
export enum GameStatus {
  GAME_STATUS_UNSPECIFIED = 0,
  GAME_STATUS_WAITING_FOR_PLAYERS = 1,
  GAME_STATUS_IN_PROGRESS = 2,
  GAME_STATUS_FINISHED = 3,
}

// 2. Concrete implementations with proper defaults
export class Book implements BookInterface {
  base?: BaseMessage;
  title: string = "";
  author: string = "";
  tags?: string[];
  available: boolean = false;
}

// 3. Enhanced factories with cross-package composition
export class LibraryV2Factory {
  // Cross-package dependency injection
  private commonFactory = new LibraryCommonFactory();
  
  // Context-aware factory methods
  newBook = (
    parent?: any,
    attributeName?: string, 
    attributeKey?: string | number,
    data?: any
  ): FactoryResult<BookInterface> => {
    const instance = new ConcreteBook();
    return { instance, fullyLoaded: false }; // Delegates to deserializer
  }
  
  // Cross-package factory delegation
  getFactoryMethod(messageType: string): FactoryMethod | undefined {
    const packageName = extractPackage(messageType);
    if (packageName === "library.common") {
      return this.commonFactory[getMethodName(messageType)];
    }
    return this[getMethodName(messageType)];
  }
}

// 4. Schema definitions with complete field metadata
export const BookSchema: MessageSchema = {
  name: "Book",
  fields: [
    {
      name: "base",
      type: FieldType.MESSAGE,
      id: 1,
      messageType: "library.common.BaseMessage"  // Cross-package reference
    },
    { name: "title", type: FieldType.STRING, id: 2 },
    { name: "tags", type: FieldType.REPEATED, id: 8, repeated: true },
    // ... field definitions with proto field IDs
  ]
};

// 5. Schema-aware deserializer with factory integration
export class LibraryV2Deserializer {
  constructor(
    private schemaRegistry: Record<string, MessageSchema>,
    private factory: FactoryInterface
  ) {}
  
  // Type-safe deserialization using schema information
  deserialize<T>(instance: T, data: any, messageType: string): T {
    const schema = this.schemaRegistry[messageType];
    for (const fieldSchema of schema.fields) {
      if (fieldSchema.type === FieldType.MESSAGE) {
        // Cross-package factory delegation
        const factoryMethod = this.factory.getFactoryMethod?.(fieldSchema.messageType!);
        if (factoryMethod) {
          const result = factoryMethod(instance, fieldSchema.name, undefined, data[fieldSchema.name]);
          instance[fieldSchema.name] = result.fullyLoaded ? 
            result.instance : 
            this.deserialize(result.instance, data[fieldSchema.name], fieldSchema.messageType!);
        }
      }
      // ... handle other field types
    }
    return instance;
  }
}
```

### 5. Runtime Package Architecture (@protoc-gen-go-wasmjs/runtime)

#### NPM Package Structure
```typescript
@protoc-gen-go-wasmjs/runtime/
├── browser/
│   └── service-manager.ts      # BrowserServiceManager for WASM service calls
├── client/
│   ├── types.ts               # WASMResponse, WasmError interfaces
│   └── base-client.ts         # WASMServiceClient base class
├── schema/
│   ├── types.ts               # FieldType, FieldSchema, MessageSchema
│   ├── base-deserializer.ts   # BaseDeserializer with all logic methods
│   └── base-registry.ts       # BaseSchemaRegistry with utility methods
└── types/
    ├── factory.ts             # FactoryInterface, FactoryResult
    └── patches.ts             # Patch operation types for stateful proxies
```

#### Template Inheritance Pattern
Generated TypeScript classes extend runtime base classes:

```typescript
// Generated client (simplified)
import { WASMServiceClient } from '@protoc-gen-go-wasmjs/runtime';

export class MyServicesClient extends WASMServiceClient {
  constructor() {
    super();
    this.myService = new MyServiceClientImpl(this);
  }
  
  // Only template-specific methods (API structure, WASM loading)
  protected getWasmMethod(methodPath: string): Function { /* generated */ }
  private async loadWASMModule(wasmPath: string): Promise<void> { /* generated */ }
}

// Generated deserializer (simplified)  
import { BaseDeserializer } from '@protoc-gen-go-wasmjs/runtime';

export class MyDeserializer extends BaseDeserializer {
  constructor() {
    super(mySchemaRegistry, myFactory); // Package-specific config
  }
  
  // Only static factory method (uses package-specific deserializer)
  static from<T>(messageType: string, data: any): T { /* generated */ }
}
```

#### Benefits of Runtime Package Approach
1. **90% bundle size reduction**: Static utilities no longer duplicated
2. **Centralized maintenance**: Runtime fixes benefit all projects immediately
3. **Tree-shakeable imports**: Consumers bundle only needed utilities
4. **Modern TypeScript support**: Full ESM/CJS builds with type definitions
5. **Inheritance-based**: Generated classes focus only on template-specific logic

### 6. Multi-Target Generation
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

### 7. Dual-Target Architecture
WASM and TypeScript artifacts can be generated independently:
```yaml
# Generate only WASM
- local: protoc-gen-go-wasmjs
  out: ./gen/wasm
  opt: [generate_typescript=false]

# Generate only TypeScript (with self-generated classes)
- local: protoc-gen-go-wasmjs  
  out: ./frontend/src/clients
  opt: [generate_wasm=false, ts_export_path=./frontend/src/types]
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