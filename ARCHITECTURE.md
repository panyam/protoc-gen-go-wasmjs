# protoc-gen-go-wasmjs Architecture

## Overview

protoc-gen-go-wasmjs follows a layered plugin architecture with BaseGenerator artifact collection that integrates with the Protocol Buffers compiler toolchain. It implements a 4-step artifact processing approach that separates protogen dependency from file mapping decisions, enabling flexible artifact grouping and reliable file generation order.

## High-Level Architecture

```
┌─────────────────┐     ┌──────────────────┐     ┌─────────────────┐
│   .proto files  │────▶│     protoc       │────▶│  Generated Code │
└─────────────────┘     └──────────────────┘     └─────────────────┘
                               │                           │
                               ▼                           ▼
                     ┌──────────────────┐         ┌─────────────────┐
                     │  BaseGenerator   │         │ Service clients │
                     │  + TSGenerator   │         │ Base bundle     │
                     │  + GoGenerator   │         │ TypeScript types│
                     └──────────────────┘         │ WASM wrapper    │
                                                  └─────────────────┘
```

## 4-Step Artifact Processing Approach

The architecture implements a clean separation between artifact collection and file generation:

```
1. COLLECT ALL ARTIFACTS    → BaseGenerator.CollectAllArtifacts()
   ├─ Get map of all artifacts from protogen
   └─ Available regardless of protoc's Generate flags

2. CLASSIFY ARTIFACTS       → ArtifactCatalog
   ├─ Services (regular + browser)
   ├─ Messages by package
   └─ Enums by package

3. MAP ARTIFACTS TO FILES   → planFilesFromCatalog()
   ├─ Generator-specific slice/dice/group logic
   ├─ N artifacts → 1 file (bundle with multiple services)
   └─ 1 artifact → 1 file (per-service clients)

4. CREATE PROTOGEN FILES    → fileSet.CreateFiles(plugin)
   ├─ Send final mapping to protogen
   └─ Only after all artifact mapping decisions are complete
```

## Core Components

### 1. BaseGenerator (`pkg/generators/base_generator.go`)
The foundation component that provides artifact collection for all generators:
- Collects complete artifact catalog from ALL proto files (ignores Generate flags)
- Classifies artifacts into services, messages, enums by package
- Provides shared utilities (ProtoAnalyzer, PathCalculator, NameConverter)
- Embedded by both GoGenerator and TSGenerator for consistency

```go
type ArtifactCatalog struct {
    Services        []ServiceArtifact  // Regular services
    BrowserServices []ServiceArtifact  // Browser-provided services
    Messages        []MessageArtifact  // Messages by package
    Enums           []EnumArtifact     // Enums by package
    Packages        map[string]*PackageInfo // Complete package map
}
```

### 2. TSGenerator (`pkg/generators/ts_generator.go`)
TypeScript-specific generator that embeds BaseGenerator:
- Collects all artifacts using BaseGenerator.CollectAllArtifacts()
- Maps artifacts to files with TypeScript-specific logic
- Generates service clients at package level (presenter/v1/presenterServiceClient.ts)
- Generates simple base bundle at module level (index.ts)
- Renders complete TypeScript artifact set per package:
  - `interfaces.ts` - Type definitions (always generated)
  - `models.ts` - Concrete implementations (always generated)
  - `factory.ts` - Object factories (when `generate_factories=true`)
  - `schemas.ts` - Field metadata (always generated)
  - `deserializer.ts` - Data population (always generated)

### 3. GoGenerator (`pkg/generators/go_generator.go`)
Go WASM-specific generator that embeds BaseGenerator:
- Uses BaseGenerator for artifact collection and utilities
- Maps artifacts to Go WASM files with Go-specific logic
- Generates WASM wrappers, examples, and build scripts
- Maintains direct file creation for Go artifacts

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

### 4. Template System (`pkg/renderers/templates/`)
Embedded templates using Go's `embed` package with inheritance-based architecture:

**Go Templates:**
- `wasm.go.tmpl` - Go WASM wrapper generation
- `main.go.tmpl` - Example usage generation
- `build.sh.tmpl` - Build script generation

**TypeScript Templates:**
- `client_simple.ts.tmpl` - Service client generation (cleaned of bundle code)
- `bundle.ts.tmpl` - Simple base bundle class extending WASMBundle
- `browser_service.ts.tmpl` - Browser service interfaces
- `interfaces.ts.tmpl` - TypeScript interface generation (type definitions only)
- `models.ts.tmpl` - TypeScript model class generation (concrete implementations)
- `factory.ts.tmpl` - TypeScript factory generation (object construction)
- `schemas.ts.tmpl` - Schema definitions (field metadata)
- `deserializer.ts.tmpl` - Schema-aware deserializers (data population)

**TypeScript Generation Model:**
The generator follows a clean separation between interfaces and implementations:
- **`interfaces.ts`**: Pure TypeScript interfaces for type safety and flexibility
- **`models.ts`**: Concrete classes implementing the interfaces with default values
- **`factory.ts`**: Factory methods for object construction with context awareness
- **`schemas.ts`**: Field metadata and protobuf schema information
- **`deserializer.ts`**: Schema-driven deserialization with factory integration

This model allows users to work with interfaces for type definitions while having concrete implementations available when needed. The factory/deserializer system handles proper default values and recursive object construction.

### 5. Simplified Bundle Architecture
The bundle architecture uses composition and inheritance patterns for maximum flexibility:

**Generated Base Bundle:**
```typescript
// generated/index.ts - Simple base class with module configuration
export class ExampleBundle extends WASMBundle {
    constructor() {
        super({
            moduleName: 'example',
            apiStructure: 'namespaced',
            jsNamespace: 'example'
        });
    }
}
```

**User Composition Pattern:**
```typescript
// User creates their own bundle with needed services
const wasmBundle = new ExampleBundle();
const presenterService = new PresenterServiceClient(wasmBundle);
const browserAPI = new BrowserAPIServiceClient(wasmBundle);

// Users can also extend for convenience
class MyAppBundle extends ExampleBundle {
    public readonly presenter: PresenterServiceClient;
    constructor() {
        super();
        this.presenter = new PresenterServiceClient(this);
    }
}
```

**Benefits:**
- No cross-package coordination complexity
- Users include only needed services
- Clean separation between WASM management and service usage
- Eliminates duplicate file generation issues
- Maximum flexibility for different use cases

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
// Generated bundle client (simplified)
import { WASMBundle, WASMBundleConfig, ServiceClient } from '@protoc-gen-go-wasmjs/runtime';

export class My_servicesBundle {
  private wasmBundle: WASMBundle;
  public readonly myService: MyServiceClient;
  
  constructor() {
    const config: WASMBundleConfig = {
      moduleName: 'my_services',
      apiStructure: 'namespaced',
      jsNamespace: 'myApp'
    };
    this.wasmBundle = new WASMBundle(config);
    this.myService = new MyServiceClient(this.wasmBundle);
  }
  
  // Only template-specific methods (WASM loading, service registration)
  async loadWasm(wasmPath: string): Promise<void> { /* generated */ }
  registerBrowserService(name: string, implementation: any): void { /* generated */ }
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

### 6. JavaScript API Structure Options

The generator supports three different API structures for how WASM methods are exposed to JavaScript:

#### **Namespaced Structure** (`js_structure=namespaced`) - **Default & Recommended**

Creates a clean hierarchical API structure:

```javascript
// Generated WASM API
window.myNamespace = {
  userService: {
    getUser: function(request) { /* ... */ },
    createUser: function(request) { /* ... */ }
  },
  orderService: {
    createOrder: function(request) { /* ... */ },
    updateOrder: function(request) { /* ... */ }
  }
};

// TypeScript client usage
const response = await client.userService.getUser({id: "123"});
```

**Benefits:**
- **🎯 Clean organization**: Services grouped logically
- **🔍 Easy discovery**: IDE autocomplete shows service structure
- **📦 Namespace isolation**: No method name conflicts between services
- **🧹 Readable code**: `client.userService.getUser()` is self-documenting

#### **Flat Structure** (`js_structure=flat`)

Creates flat function names with prefixes:

```javascript
// Generated WASM API
window.myNamespaceUserServiceGetUser = function(request) { /* ... */ };
window.myNamespaceUserServiceCreateUser = function(request) { /* ... */ };
window.myNamespaceOrderServiceCreateOrder = function(request) { /* ... */ };

// Internal client usage (less readable)
const response = await client.callMethod('myNamespaceUserServiceGetUser', request);
```

**When to use:**
- **Legacy compatibility** with existing flat API expectations
- **Minimal bundle size** for single-service projects
- **Simple debugging** with predictable global function names

#### **Service-Based Structure** (`js_structure=service_based`)

Creates service-oriented grouping:

```javascript
// Generated WASM API
window.services = {
  user: {
    getUser: function(request) { /* ... */ },
    createUser: function(request) { /* ... */ }
  },
  order: {
    createOrder: function(request) { /* ... */ },
    updateOrder: function(request) { /* ... */ }
  }
};

// TypeScript client usage
const response = await client.services.user.getUser({id: "123"});
```

**When to use:**
- **Micro-frontend architecture** where services are primary concept
- **Multiple independent service modules** loaded dynamically
- **Framework integration** where `services` is a standard pattern

#### **API Structure Impact on Generation**

The chosen structure affects multiple aspects:

**1. WASM Method Resolution:**
```typescript
// Namespaced: namespace.service.method
protected getWasmMethod(methodPath: string): Function {
  const parts = methodPath.split('.');
  let current = this.wasm; // = window.myNamespace
  for (const part of parts) {
    current = current[part]; // Navigate: service → method
  }
  return current;
}

// Flat: direct function name
protected getWasmMethod(methodPath: string): Function {
  return this.wasm[methodPath]; // Direct: window.myNamespaceServiceMethod
}
```

**2. Client Interface Generation:**
```typescript
// Namespaced generates clean service clients
export class MyClient extends WASMServiceClient {
  public readonly userService: UserServiceClientImpl;
  public readonly orderService: OrderServiceClientImpl;
}

// Flat generates method-based calls
export class MyClient extends WASMServiceClient {
  async getUserData() { return this.callMethod('myNamespaceUserServiceGetUser', ...); }
}
```

**3. Method Call Paths:**
- **Namespaced**: `userService.getUser` → calls `namespace.userService.getUser`
- **Flat**: `getUser` → calls `namespaceUserServiceGetUser`  
- **Service-based**: `userService.getUser` → calls `services.user.getUser`

#### **Configuration Examples**

```yaml
# Namespaced (recommended for most projects)
- local: protoc-gen-go-wasmjs
  opt:
    - js_structure=namespaced
    - js_namespace=myApp
    - module_name=my_services

# Flat (for legacy compatibility)  
- local: protoc-gen-go-wasmjs
  opt:
    - js_structure=flat
    - js_namespace=MyApp
    - module_name=my_services

# Service-based (for micro-frontends)
- local: protoc-gen-go-wasmjs
  opt:
    - js_structure=service_based
    - js_namespace=unused  # Not used in service_based mode
    - module_name=my_services
```

### 7. Multi-Target Generation
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
