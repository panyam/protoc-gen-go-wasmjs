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
├── pkg/generator/                # Core generation logic
│   ├── templates/                # Embedded templates
│   │   ├── wasm.go.tmpl         # Go WASM wrapper
│   │   ├── client_simple.ts.tmpl # Simplified TypeScript client
│   │   ├── interfaces.ts.tmpl   # TypeScript interfaces
│   │   ├── models.ts.tmpl       # TypeScript model classes
│   │   ├── factory.ts.tmpl      # Enhanced TypeScript factories
│   │   ├── schemas.ts.tmpl      # Schema definitions for type-safe deserialization
│   │   ├── deserializer.ts.tmpl # Schema-aware deserializers
│   │   └── build.sh.tmpl        # Build script
│   ├── config.go                # Configuration parsing
│   ├── generator.go             # Main generation logic
│   ├── tsgenerator.go          # TypeScript-specific generation
│   └── types.go                 # Template data structures
├── proto/wasmjs/v1/             # WASM annotations
├── examples/
│   ├── connect4/                # Working multiplayer Connect4 demo
│   └── library/                 # Library management examples
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
```typescript
// Generated per proto package
export interface Book { id: string; title: string; base?: BaseMessage; }
export class Book implements BookInterface { /* ... */ }

// Enhanced factories with context-aware construction
export class LibraryV2Factory {
  private commonFactory = new LibraryCommonFactory(); // Cross-package dependency
  
  newBook = (parent?: any, attributeName?: string, attributeKey?: string | number, data?: any): FactoryResult<Book> => {
    const instance = new ConcreteBook();
    return { instance, fullyLoaded: false }; // Delegates to deserializer
  }
  
  getFactoryMethod(messageType: string) { /* Cross-package delegation */ }
}

// Schema-aware deserializer with factory composition
export class LibraryV2Deserializer {
  constructor(private schemaRegistry: Record<string, MessageSchema>, private factory: FactoryInterface) {}
  
  deserialize<T>(instance: T, data: any, messageType: string): T { /* Schema-based field processing */ }
  createAndDeserialize<T>(messageType: string, data: any): T | null { /* Factory integration */ }
}

// Generated schemas with field metadata
export const BookSchema: MessageSchema = {
  name: "Book",
  fields: [
    { name: "base", type: FieldType.MESSAGE, id: 1, messageType: "library.common.BaseMessage" },
    { name: "title", type: FieldType.STRING, id: 2 },
    // ... other fields with proto field IDs and types
  ]
};
```

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
The project has completed a comprehensive refactoring and achieved **production-ready status** with:
- ✅ **Enhanced Factory Method Design** with context-aware construction and parent object tracking
- ✅ **Schema-Aware Deserialization** with type-safe field resolution and proto field ID support
- ✅ **Cross-Package Factory Composition** with automatic dependency detection and delegation
- ✅ **Package-Scoped Schema Registries** for conflict-free multi-version support
- ✅ **Self-contained TypeScript generation** eliminating external generator dependencies
- ✅ **Simplified client architecture** with direct JSON serialization
- ✅ **Multi-target generation support** for flexible deployment patterns
- ✅ **Template-based generation system** with full customization support
- ✅ **Production-ready code generation** with comprehensive testing
- ✅ **Per-service client generation** following proto directory structure
- ✅ **Runtime package architecture** with inheritance-based TypeScript clients
- ✅ **Browser service communication** with full WASM ↔ JavaScript integration

### Recent Quality & TypeScript Improvements (Latest)
- ✅ **Native Map Type Support**: Fixed proto `map<K,V>` fields to generate TypeScript `Map<K,V>` instead of synthetic interfaces
- ✅ **Framework Schema Separation**: Separated framework types (`FieldType`, `FieldSchema`) into `deserializer_schemas.ts` for cleaner architecture  
- ✅ **Package-Based Generation**: Transitioned from file-based to package-based TypeScript generation eliminating import issues
- ✅ **TypeScript Type Safety**: Fixed factory method subscripting and interface compatibility issues for full type safety
- ✅ **External Type Mapping System**: Comprehensive support for external protobuf types with configurable mappings, factory integration, and proper import handling
- ✅ **Developer Experience Enhancements**: Ergonomic API improvements with MESSAGE_TYPE constants, static deserializer methods, and performance-optimized shared instances

### Latest Bug Fixes & Enum Support (January 2025)
- ✅ **wasmjs.v1 Package Filtering**: Fixed artifact generation for wasmjs annotation packages - they are now correctly excluded from generation while remaining visible for proto compilation
- ✅ **Comprehensive Enum Support**: Implemented complete enum collection, generation, and import system for TypeScript
  - Enums are generated in interfaces.ts with proper TypeScript enum syntax
  - All generated TypeScript files (models.ts, factory.ts) now correctly import and reference enums
  - Cross-package enum references work seamlessly with the import resolution system
  - Fixed template data structures to include enums in all generation contexts
- ✅ **Enhanced Cross-Package Import Detection**: Improved import resolution to filter out wasmjs.v1 dependencies in factory composition

### Split Generator Architecture & Per-Service Generation (September 2025)
- ✅ **Phase 1: Split Architecture**: Separate Go and TypeScript generators with layered architecture
- ✅ **Phase 2: Runtime Package Migration**: Extracted common utilities to `@protoc-gen-go-wasmjs/runtime`
- ✅ **Phase 3: Template Inheritance**: Generated clients extend base classes from runtime package
- ✅ **Phase 4: Per-Service Generation**: Each service generates to separate file following proto structure
- ✅ **Phase 5: Browser Service Integration**: Full WASM ↔ browser communication with proper serialization
- ✅ **Phase 6: Production Validation**: End-to-end testing with working browser-callbacks example

### Server Streaming Support (August 2025)
- ✅ **Phase 1: Server Streaming Implementation**: Complete server-side streaming from WASM to JavaScript
  - Callback-based streaming API: `method(request, (response, error, done) => boolean)`
  - Goroutine-based WASM implementation with proper `stream.Recv()` handling
  - User-controlled stream cancellation via callback return values
  - Proper error handling and EOF detection
  - TypeScript interface generation with correct streaming signatures
- 🔄 **Phase 2: Client Streaming** (Planned): Connection objects with `send()`, `close()`, `isOpen()` methods
- 🔄 **Phase 3: Bidirectional Streaming** (Planned): Combined server and client streaming capabilities

**Major Architecture Achievements**: 
1. **Factory Composition System**: Implemented sophisticated cross-package factory delegation enabling seamless object creation across package boundaries with automatic dependency injection
2. **Schema-Aware Architecture**: Built complete schema generation and deserialization system with field metadata, proto field IDs, and oneof support for type-safe runtime processing
3. **Self-Generated TypeScript**: Successfully transitioned from complex conversion-based architecture to streamlined self-generated TypeScript classes that match Go's protojson format exactly
4. **Type-Safe Map Handling**: Proper conversion of protobuf map fields to native TypeScript Map types with synthetic message filtering
5. **External Type Integration**: Complete external type mapping system with configurable mappings, table-driven factory methods, and seamless conversion between protobuf and TypeScript types
6. **Developer Experience Excellence**: Type-safe MESSAGE_TYPE constants, ergonomic static deserialization methods, and performance-optimized shared factory instances for production-ready usage

**Production Readiness**: System handles complex nested object hierarchies, cross-package dependencies, real-world proto features (maps, external types), per-service client generation, browser service integration, and maintains full TypeScript type safety with comprehensive test validation.