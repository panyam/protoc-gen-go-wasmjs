# protoc-gen-go-wasmjs Summary

## Overview
protoc-gen-go-wasmjs is a Protocol Buffers compiler plugin that generates WASM bindings and TypeScript clients for gRPC services. It enables local-first applications where the same service logic can run in both server environments (with full database access) and browser environments (with local storage).

## Core Features

### 1. **Dual-Target Architecture**
- Generate WASM wrappers and TypeScript clients separately or together
- Flexible deployment with TypeScript clients directly in frontend directories
- Clean separation of concerns with independent generation targets

### 2. **Self-Contained TypeScript Generation**
- Generates lightweight TypeScript interfaces, models, and factories directly from proto
- Eliminates dependencies on external TypeScript generators (protoc-gen-es, protoc-gen-ts)
- Perfect compatibility with Go's protojson format without conversion layers

### 3. **Simplified TypeScript Architecture**
- Direct JSON serialization/deserialization without complex conversion layers
- Generated TypeScript classes match Go's protojson format exactly
- Optional field handling for message types and arrays
- Lightweight client implementation with minimal overhead

### 4. **Enhanced Factory & Deserialization System**
- Context-aware factory methods with parent object tracking
- Schema-aware deserialization with type-safe field resolution
- Cross-package factory composition with automatic dependency injection
- Package-scoped schema registries for conflict-free multi-version support

### 5. **Multi-Service & Multi-Target Support**
- Bundle multiple services in a single WASM module
- Generate optimized bundles per page/use case
- Service filtering for targeted deployments
- Dependency injection through export pattern

### 6. **Developer Experience**
- Seamless buf.build integration
- Generated build scripts and examples
- Self-contained generation with no external dependencies
- Template-based customization system

## Examples & Demonstrations

### Connect4 Multiplayer Game (`examples/connect4/`)
A fully working real-time multiplayer Connect4 game demonstrating:
- **WASM Integration**: Complete Go service compiled to WebAssembly
- **Pluggable Transports**: IndexedDB+polling and BroadcastChannel for cross-tab multiplayer
- **State Persistence**: Games survive browser restarts via localStorage + IndexedDB  
- **TypeScript Client**: Generated client with full type safety calling WASM methods
- **Independent Module**: Standalone go.mod with parent module replacement for development
- **Production-Ready**: Working demo with proper build system (webpack + TypeScript)

### Library Management (`examples/library/`)  
Complex multi-package examples showing:
- **Cross-Package Dependencies**: Services spanning multiple proto packages
- **Enhanced Factory Patterns**: Context-aware object creation with parent tracking
- **Schema-Aware Deserialization**: Type-safe runtime processing with field metadata
- **Multi-Service Bundling**: Different WASM bundles for different use cases

## Project Structure

```
├── cmd/protoc-gen-go-wasmjs/     # Plugin entry point
├── pkg/
│   ├── generators/               # Core generation logic
│   │   ├── base_generator.go   # Artifact collection & catalog
│   │   ├── ts_generator.go     # TypeScript-specific generation
│   │   └── go_generator.go     # Go WASM-specific generation
│   ├── renderers/
│   │   └── templates/           # Embedded templates
│   │       ├── wasm_converters.go.tmpl     # JS/Go type converters & stream wrappers
│   │       ├── wasm_exports.go.tmpl        # Exports struct & RegisterAPI
│   │       ├── wasm_browser_clients.go.tmpl # Browser service clients
│   │       ├── client_simple.ts.tmpl # TypeScript service clients
│   │       ├── bundle.ts.tmpl       # Bundle base class
│   │       ├── interfaces.ts.tmpl   # TypeScript interfaces
│   │       ├── models.ts.tmpl       # Concrete implementations
│   │       ├── factory.ts.tmpl      # Object factories
│   │       ├── schemas.ts.tmpl      # Field metadata
│   │       ├── deserializer.ts.tmpl # Data deserialization
│   │       └── build.sh.tmpl        # Build script
│   ├── builders/                # Template data building
│   ├── filters/                 # Service/method filtering
│   └── collectors/              # Message/enum collection
├── proto/wasmjs/v1/             # WASM annotations
├── runtime/                     # @protoc-gen-go-wasmjs/runtime package
├── example/
└── PROTO_CONVERSION.md          # Conversion documentation
```

## Key Design Decisions

### 1. **Export Pattern over Main Functions**
Instead of generating `main()` functions, we generate export structs that allow full dependency injection:
```go
exports := &ServicesExports{
    UserService: &myUserService{db: database, auth: authService}
}
exports.RegisterAPI()
```

### 2. **Template-Based Generation**
Using Go's `embed` package for templates provides:
- Clean separation of generation logic and output format
- Easy customization through template overrides
- Maintainable code generation

### 3. **Enhanced TypeScript Generation System**
Generates complete TypeScript ecosystem with advanced features:

**Generated Files Per Package:**
```typescript
// interfaces.ts - Pure TypeScript interfaces for type safety
export interface Book {
  id: string;
  title: string;
  base?: BaseMessage;
}

// models.ts - Concrete implementations with default values
export class Book implements BookInterface {
  id: string = "";
  title: string = "";
  base?: BaseMessage;
}

// factory.ts - Context-aware object construction (when generate_factories=true)
export class LibraryV2Factory {
  private commonFactory = new LibraryCommonFactory(); // Cross-package dependency

  newBook = (parent?: any, attributeName?: string, attributeKey?: string | number, data?: any): FactoryResult<Book> => {
    const instance = new ConcreteBook();
    return { instance, fullyLoaded: false }; // Delegates to deserializer
  }

  getFactoryMethod(messageType: string) { /* Cross-package delegation */ }
}

// schemas.ts - Field metadata for runtime processing
export const BookSchema: MessageSchema = {
  name: "Book",
  fields: [
    { name: "base", type: FieldType.MESSAGE, id: 1, messageType: "library.common.BaseMessage" },
    { name: "title", type: FieldType.STRING, id: 2 },
    // ... other fields with proto field IDs and types
  ]
};

// deserializer.ts - Schema-driven deserialization
export class LibraryV2Deserializer {
  constructor(private schemaRegistry: Record<string, MessageSchema>, private factory: FactoryInterface) {}

  deserialize<T>(instance: T, data: any, messageType: string): T { /* Schema-based field processing */ }
  static from<T>(messageType: string, data: any): T { /* Static convenience method */ }
}
```

**Clean Architecture Benefits:**
- **Interfaces** provide type safety without implementation overhead
- **Models** offer concrete classes when needed
- **Factories** handle object construction with proper defaults
- **Schemas** enable runtime type introspection
- **Deserializers** populate objects with schema awareness

### 4. **Flexible API Structures**
Supports multiple JavaScript API patterns:
- Namespaced: `myapp.service.method()`
- Flat: `myappServiceMethod()`
- Service-based: `services.library.findBooks()`

## Technical Achievements

- **Zero Configuration**: Works out of the box with smart defaults
- **Backward Compatible**: Existing configurations continue to work
- **Production Ready**: Used in real projects with complex proto structures
- **Performance Optimized**: Minimal overhead in proto conversions
- **Type Safe**: Full TypeScript support with proper type inference

## Use Cases

### Local-First Applications
Same business logic runs in both environments:
- **Server**: Full database access, complete dataset
- **Browser**: Local storage, user's subset of data

### Progressive Web Apps
- Offline functionality with WASM services
- Sync when online, work offline
- Reduced server load

### Edge Computing
- Deploy services to CDN edge workers
- Run logic closer to users
- Maintain consistent API across deployments

## Current Status (September 2025)
The project has completed a comprehensive refactoring and achieved **production-ready status** with BaseGenerator architecture:

### **Core Architecture (September 2025)**
- **BaseGenerator Artifact Collection**: Complete 4-step approach separating artifact collection from protogen file creation
- **Cross-Package Visibility**: BaseGenerator collects ALL artifacts regardless of protoc's Generate flags
- **Flexible File Mapping**: N artifacts to 1 file or 1 artifact to N files based on generator logic
- **Simplified Bundle Generation**: Base bundle class with user composition patterns
- **Per-Service Client Files**: Clean package-level service client generation
- **TypeScript Type Safety**: Full compilation with proper interface generation

### **Completed Major Features**
- **Enhanced Factory Method Design** with context-aware construction and parent object tracking
- **Schema-Aware Deserialization** with type-safe field resolution and proto field ID support
- **Cross-Package Factory Composition** with automatic dependency detection and delegation
- **Package-Scoped Schema Registries** for conflict-free multi-version support
- **Self-contained TypeScript generation** eliminating external generator dependencies
- **Runtime package architecture** with inheritance-based TypeScript clients
- **Browser service communication** with full WASM ↔ JavaScript integration

### Recent Quality & TypeScript Improvements (Latest)
- ✅ **Native Map Type Support**: Fixed proto `map<K,V>` fields to generate TypeScript `Map<K,V>` instead of synthetic interfaces
- ✅ **Framework Schema Separation**: Separated framework types (`FieldType`, `FieldSchema`) into `deserializer_schemas.ts` for cleaner architecture
- ✅ **Package-Based Generation**: Transitioned from file-based to package-based TypeScript generation eliminating import issues
- ✅ **TypeScript Type Safety**: Fixed factory method subscripting and interface compatibility issues for full type safety
- ✅ **External Type Mapping System**: Comprehensive support for external protobuf types with configurable mappings, factory integration, and proper import handling
- ✅ **Developer Experience Enhancements**: Ergonomic API improvements with MESSAGE_TYPE constants, static deserializer methods, and performance-optimized shared instances
- ✅ **Factory/Deserializer Generation**: Completed wiring for models, factory, schemas, and deserializer file generation in new catalog-based architecture (October 2025)
- ✅ **pkg.go.dev Documentation**: Comprehensive godoc documentation added to all packages with examples, architecture diagrams, and complete API reference (October 2025)

### Latest Updates (January 2025)

#### Go WASM File Splitting (January 2025)
- **Modular Go Generation**: Split monolithic WASM file into 3 focused files
  - `<package>_converters.wasm.go` - JS/Go type converters (`createJSResponse`) and stream wrappers
  - `<package>_exports.wasm.go` - Exports struct, `RegisterAPI()`, and service method wrappers
  - `<package>_browser_clients.wasm.go` - Browser service client implementations
- **Better Organization**: Clear separation of concerns for maintainability
- **Selective Generation**: Converters only when services exist, browser clients only when needed
- **All files**: Still use `//go:build js && wasm` build tags

#### go_package Output Path Collision Fix (November 2025)
- **Issue Resolved**: Multiple proto files with same proto package but different `go_package` options no longer collide
- **Root Cause**: Output path calculation only used proto package name, causing file overwrites
- **Solution Implemented**:
  - `calculateOutputPath()` method extracts path from `go_package` option (e.g., `.../v1/models` → `test/v1/models`)
  - `calculateBaseName()` method incorporates go_package suffix for unique filenames
  - Files now correctly segregated: `.../v1/models/...exports.wasm.go` and `.../v1/services/...exports.wasm.go`
- **Standard Pattern Support**: Enables separating models and services into different Go packages to avoid gRPC dependencies
- **Test Cases**: Added `test_one_package/`, `test_multi_packages/`, and `test_broken/` examples validating all patterns work correctly

#### Bug Fixes & Enum Support
- **wasmjs.v1 Package Filtering**: Fixed artifact generation for wasmjs annotation packages - they are now correctly excluded from generation while remaining visible for proto compilation
- **Comprehensive Enum Support**: Implemented complete enum collection, generation, and import system for TypeScript
  - Enums are generated in interfaces.ts with proper TypeScript enum syntax
  - All generated TypeScript files (models.ts, factory.ts) now correctly import and reference enums
  - Cross-package enum references work seamlessly with the import resolution system
  - Fixed template data structures to include enums in all generation contexts
- **Enhanced Cross-Package Import Detection**: Improved import resolution to filter out wasmjs.v1 dependencies in factory composition

### Cross-Package Import System (October 2025)
- **Issue Resolved**: Missing TypeScript imports for types from other proto packages in same project
- **Protobuf Descriptor API**: Uses descriptor methods (`FullName()`, `ParentFile().Package()`, `Parent()`) for accurate type information instead of string parsing
- **Nested Type Handling**: Properly flattens nested message types (e.g., `ParentMessage.NestedType` becomes `ParentMessage_NestedType`) to prevent name collisions
- **Metadata Fields**: Added `MessagePackage` and `IsNestedType` to `TSFieldInfo` for storing descriptor-derived information
- **Import Path Calculation**: Correctly generates relative import paths between packages (e.g., `../../utils/v1/interfaces`)
- **MessageCollector Fix**: Fixed to use `message.Desc.FullName()` for accurate fully qualified names including parent messages
- **Type Safety**: Field types use flattened names for nested types, simple names for top-level types
- **Comprehensive Testing**: Unit tests for package extraction, type name flattening, and import path calculation

### Split Generator Architecture & Per-Service Generation (September 2025)
- ✅ **Phase 1: Split Architecture**: Separate Go and TypeScript generators with layered architecture
- ✅ **Phase 2: Runtime Package Migration**: Extracted common utilities to `@protoc-gen-go-wasmjs/runtime`
- ✅ **Phase 3: Template Inheritance**: Generated clients extend base classes from runtime package
- ✅ **Phase 4: Per-Service Generation**: Each service generates to separate file following proto structure
- ✅ **Phase 5: Browser Service Integration**: Full WASM ↔ browser communication with proper serialization
- ✅ **Phase 6: Production Validation**: End-to-end testing with working example

### Server Streaming Support (August 2025)
- ✅ **Phase 1: Server Streaming Implementation**: Complete server-side streaming from WASM to JavaScript
  - Callback-based streaming API: `method(request, (response, error, done) => boolean)`
  - Goroutine-based WASM implementation with proper `stream.Recv()` handling
  - User-controlled stream cancellation via callback return values
  - Proper error handling and EOF detection
  - TypeScript interface generation with correct streaming signatures
- 🔄 **Phase 2: Client Streaming** (Planned): Connection objects with `send()`, `close()`, `isOpen()` methods
- 🔄 **Phase 3: Bidirectional Streaming** (Planned): Combined server and client streaming capabilities

### **BaseGenerator Architecture Implementation (September 2025)**
**Complete artifact-centric approach implemented:**

1. **BaseGenerator Foundation**: Both TSGenerator and GoGenerator embed BaseGenerator for unified artifact collection
2. **4-Step Processing**: CollectAllArtifacts() → Classify → Map → CreateFiles() separates concerns cleanly  
3. **Cross-Package Visibility**: Artifact collection sees ALL packages/services regardless of protoc Generate flags
4. **Flexible File Mapping**: Generator-specific slice/dice/group logic with N:1 and 1:N file mapping
5. **Delayed File Creation**: protogen.NewGeneratedFile() calls delayed until after all artifact mapping decisions
6. **Simplified Bundle**: Base bundle extends WASMBundle with module config, users add services via composition

**Major Architecture Achievements**: 
1. **Artifact-Driven Generation**: Complete separation of artifact discovery from file generation decisions
2. **Generator Independence**: Each generator controls its own file mapping logic without protogen constraints
3. **Bundle Architecture Simplification**: Eliminated complex cross-package coordination in favor of user composition
4. **File Visit Order Resolution**: Protogen only involved in final step after all mapping decisions complete
5. **Template Architecture**: Clean separation between service clients, bundle, browser services, and type files
6. **User Experience**: Composition pattern gives users complete control over service inclusion

**Production Readiness**: System provides complete artifact visibility for mapping decisions, eliminates file generation ordering issues, maintains TypeScript type safety, and enables flexible user patterns for service composition.
