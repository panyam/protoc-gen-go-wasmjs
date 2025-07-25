# protoc-gen-go-wasmjs Implementation Plan

## Overview
Transform existing MCP generator into a WASM generator that produces browser-compatible gRPC services with extensive template customization. Generate Go WASM wrappers and TypeScript clients that integrate seamlessly with buf.build toolchain.

## Architecture Philosophy
- **Ecosystem Integration**: Work with existing protoc generators (protoc-gen-go-grpc, protoc-gen-es, etc.)
- **Template Flexibility**: Extensive customization via configuration and template overrides
- **Multi-Service Support**: Bundle related services in single WASM modules
- **Clean API Structure**: Namespaced APIs instead of flat function names
- **Local-First Enablement**: Same types work for WASM and HTTP clients

## Progress Tracking

### ✅ Phase 1: Project Transformation & Configuration System (COMPLETED)
- [x] **1.1 Core Project Updates**
  - [x] Update `go.mod` from `protoc-gen-go-mcp` → `protoc-gen-go-wasmjs`
  - [x] Rename `cmd/protoc-gen-go-mcp/` → `cmd/protoc-gen-go-wasmjs/`
  - [x] Update all import paths and remove MCP dependencies
  - [x] Remove MCP-specific code from `pkg/generator/generator.go`

- [x] **1.2 Advanced Configuration System**
  - [x] Create `pkg/generator/config.go` with comprehensive option parsing
  - [x] Implement glob pattern matching for method filtering
  - [x] Add configuration validation with helpful error messages
  - [x] Support template directory discovery and validation

### ✅ Phase 2: Template System Architecture (COMPLETED)
- [x] **2.1 Template Data Structure**
  - [x] Design core `TemplateData` struct with multi-service support
  - [x] Implement `ServiceData` and `MethodData` structures
  - [x] Add customization fields (JSNamespace, ModuleName, APIStructure)

- [x] **2.2 Template Override System**
  - [x] Default templates embedded in binary using `go:embed`
  - [x] Template discovery from `template_dir` option
  - [x] Clean separation between Go WASM, TypeScript client, and build script templates
  - [x] Template helper functions for customization

### ✅ Phase 3: Multi-Service WASM Generation (COMPLETED)
- [x] **3.1 Namespaced API Structure (Default)**
  - [x] Generate clean namespaced JavaScript structure
  - [x] Multi-service WASM wrapper with global service registry
  - [x] Service initialization and injection functions

- [x] **3.2 Alternative API Structures**
  - [x] Flat structure: `bookstoreLibraryFindBooks()` (backward compatibility)
  - [x] Service-based: `services.library.findBooks()` (enterprise style)
  - [x] Configurable via `js_structure` option

- [x] **3.3 Method Generation with Filtering**
  - [x] Implement include/exclude glob pattern matching
  - [x] Method renaming functionality
  - [x] Template logic for conditional method generation

### ✅ Phase 4: Enhanced TypeScript Client Generation (COMPLETED)
- [x] **4.1 Multi-Service Client Structure**
  - [x] Generate main client class with service-specific sub-clients
  - [x] WASM loading and initialization logic
  - [x] Clean API separation between services

- [x] **4.2 Generator-Specific Integration**
  - [x] protoc-gen-es integration with `.toJson()` and `.fromJson()`
  - [x] protoc-gen-ts integration with `.toJSON()` and `fromJSON()`
  - [x] Generic fallback with JSON.stringify/parse
  - [x] Auto-detection of conversion methods

- [x] **4.3 Advanced Type Conversion**
  - [x] Convention-based auto-detection system
  - [x] Method path resolution for namespaced APIs
  - [x] Structured error handling with WasmError class

### ✅ Phase 5: Build Pipeline & Developer Experience (COMPLETED)
- [x] **5.1 Generated Build Integration**
  - [x] Generate `build.sh` for WASM compilation
  - [x] Include `wasm_exec.js` copying and versioning
  - [x] Support for different Go versions and build flags
  - [x] Integration with existing buf workflows

- [x] **5.2 Output Structure**
  - [x] Organize generated files with configurable export paths
  - [x] Separate TS and WASM output directories
  - [x] Include necessary runtime files

### 🔄 Phase 6: Advanced Features & Error Handling (IN PROGRESS)
- [x] **6.1 Error Handling & Debugging**
  - [x] Structured error responses with error codes
  - [ ] Source maps for WASM debugging
  - [ ] Performance monitoring hooks
  - [ ] Graceful degradation patterns

- [x] **6.2 Method Customization Examples**
  - [x] Real-world configuration examples in documentation
  - [x] Complex filtering scenarios
  - [x] Custom naming conventions

### ✅ Phase 7: Example & Documentation (COMPLETED)
- [x] **7.1 Complete LibraryService Example**
  - [x] Multi-service proto definition (library.proto with LibraryService and UserService)
  - [x] Full buf.gen.yaml with all options demonstrated
  - [x] WASM annotation examples
  - [x] Generated code demonstrates all features

- [x] **7.2 Migration & Best Practices Guide**
  - [x] Complete README with usage examples
  - [x] Local-first architecture patterns
  - [x] TypeScript integration guide
  - [x] Configuration reference

## Current Checkpoint: January 25, 2025

### What We've Accomplished
**Complete, production-ready WASM generator** with advanced multi-target capabilities:

1. **Core Architecture**: Successfully migrated from MCP tool generation to WASM binding generation
2. **Template System**: Implemented robust `go:embed` template system with Go WASM, TypeScript client, and build script generation
3. **Configuration**: Comprehensive option parsing supporting all major customization scenarios
4. **Multi-Service Support**: Can bundle related services in single WASM modules with clean APIs
5. **TypeScript Integration**: Auto-detects and works with popular protobuf TypeScript generators (protoc-gen-es, protoc-gen-ts)
6. **WASM Annotations**: Custom protobuf annotations for fine-grained control
7. **Multi-File Package Support**: Correctly handles packages with multiple proto files, avoiding duplicate generation
8. **Relative Path Resolution**: Proper TypeScript import path calculation across complex directory structures
9. **Real-World Testing**: Successfully generates WASM modules for complex projects with multiple services

### Key Learnings
1. **Template Approach**: Using `go:embed` with separate template files is much cleaner than string constants
2. **Export Path Separation**: Important to separate import paths (where we read types) from export paths (where we write files)
3. **API Structure Flexibility**: Different projects need different JavaScript API styles (namespaced vs flat vs service-based)
4. **Type Conversion**: Auto-detection of TypeScript conversion methods enables broad generator compatibility
5. **Local-First Value**: The LibraryService example clearly demonstrates the power of identical logic in server and browser
6. **Package Grouping**: Proto files should be grouped by package to generate one WASM module per package
7. **Service Detection**: Early return logic must check all files in a package, not just the primary file
8. **Multi-Target Need**: Different pages/use cases need different service combinations for optimal bundle sizes
9. **Dependency Injection**: Generated `main()` functions prevent dependency injection - export pattern is more flexible

### Technical Achievements
- **Zero breaking changes** during transformation
- **Embedded templates** eliminate string escaping issues
- **Comprehensive configuration** with validation and helpful error messages
- **Multi-service bundling** reduces WASM overhead
- **Clean separation** between generated code and user implementations
- **Duplicate file prevention** via package-based generation
- **Correct relative imports** across complex directory structures
- **Working end-to-end** with real proto definitions and multiple services

## Next Steps

### Immediate (Current Sprint)
- [x] **End-to-End Testing**: Test complete generation pipeline with LibraryService example
- [x] **WASM Compilation**: Verify generated Go code compiles to WASM successfully
- [x] **TypeScript Integration**: Test with real protoc-gen-es generated types
- [x] **Multi-Package Support**: Fix duplicate file generation and service detection
- [ ] **Multi-Target Implementation**: Add support for multiple targets per project
- [ ] **Export Pattern**: Replace main() generation with flexible Export pattern
- [ ] **Documentation Update**: Update README for new multi-target workflow

### Short Term (Next Month)
- [ ] **Advanced Multi-Target**: Support complex service combinations and custom naming
- [ ] **Dependency Injection Examples**: Show real-world service implementations with DB, auth, etc.
- [ ] **Bundle Optimization**: Analyze and optimize WASM bundle sizes per target
- [ ] **Browser Demo**: Create complete browser demo with multiple targets
- [ ] **Community Feedback**: Gather feedback from early adopters

### Medium Term (Next Quarter)
- [ ] **Streaming Support**: Research and implement streaming RPC support for WASM
- [ ] **Advanced Templates**: Template inheritance and partial overrides
- [ ] **Monitoring Integration**: Performance monitoring and analytics hooks
- [ ] **IDE Support**: Language server and IDE plugin support

### Long Term (6+ Months)
- [ ] **Multi-Language**: Explore Rust/C++ service implementations with same TypeScript clients
- [ ] **Edge Computing**: Optimize for edge/CDN deployment scenarios
- [ ] **Enterprise Features**: Advanced security, authentication, and authorization patterns
- [ ] **Ecosystem Integration**: Integration with popular frameworks (React, Vue, Angular)

## Upcoming Architecture: Multi-Target + Export Pattern

### The Problem We're Solving
Current single-module approach generates one WASM binary with all services, but real applications need:
- **Page-specific bundles**: User page only needs UsersService, Game page needs GamesService + WorldsService
- **Dependency injection**: Service implementations need database connections, config, auth clients
- **Custom initialization**: Each target may need different middleware, endpoints, logging

### The Solution: Multi-Target Export Pattern
```yaml
# Multiple targets in buf.gen.yaml
plugins:
  # User page target (minimal bundle)
  - local: protoc-gen-go-wasmjs
    out: ./gen/wasm/user-page
    opt: [services=UsersService, export_pattern=true]
    
  # Game page target (optimized bundle)  
  - local: protoc-gen-go-wasmjs
    out: ./gen/wasm/game-page
    opt: [services=GamesService,WorldsService, export_pattern=true]
```

Generated exports allow full dependency injection:
```go
// User creates cmd/user-page-wasm/main.go
exports := &UserPageServicesExports{
    UsersService: &services.UsersService{
        DB: postgresDB,
        Auth: authService,
    },
}
exports.RegisterAPI()
```

### Benefits
- **🎯 Targeted bundles**: Each page gets exactly what it needs
- **💉 Full DI control**: Inject any dependencies (DB, auth, config)
- **📦 Smaller bundles**: No unused services in production
- **🔧 Customizable**: Add middleware, custom endpoints, logging
- **🧪 Testable**: Easy to mock services for testing

## Success Metrics
- [x] **Generator compiles successfully** without errors
- [x] **Example generates complete WASM module** 
- [x] **TypeScript client works with generated types**
- [x] **Real-world multi-service project works**
- [x] **Multi-file packages handled correctly**
- [ ] **Multi-target generation implemented**
- [ ] **Export pattern with dependency injection**
- [ ] **Browser demo demonstrates local-first capability**
- [ ] **Performance comparable to manual WASM implementation**

## Configuration Schema

```yaml
# Full configuration schema
plugins:
  - plugin: go-wasmjs
    out: gen/wasm
    opt:
      # Core integration
      - ts_generator=protoc-gen-es        # protoc-gen-es, protoc-gen-ts, etc.
      - ts_import_path=./gen/ts           # where TS types are generated (for imports)
      - ts_export_path=./gen/wasm         # where TS client should be generated
      - wasm_export_path=./gen/wasm       # where WASM wrapper should be generated
      
      # Service & method selection
      - services=LibraryService,UserService  # specific services (default: all)
      - method_include=Find*,Get*,Create*     # glob patterns for methods
      - method_exclude=*Internal,*Debug       # exclude patterns
      - method_rename=FindBooks:searchBooks   # rename methods
      
      # JS API structure
      - js_structure=namespaced           # namespaced|flat|service_based
      - js_namespace=bookstore            # global namespace
      - module_name=bookstore_services    # WASM module name
      
      # Customization
      - template_dir=./custom-templates   # override templates
      - wasm_template=custom.wasm.go.tmpl # specific templates
      - ts_template=custom.client.ts.tmpl
      
      # Build integration
      - wasm_package_suffix=wasm
      - generate_build_script=true
```

## Template Data Structures

```go
type TemplateData struct {
    // Core data
    Services    []ServiceData
    Config      GeneratorConfig
    
    // Customization
    JSNamespace string
    ModuleName  string
    APIStructure string // namespaced|flat|service_based
    
    // Build info
    GeneratedImports []string
    BuildScript      string
}

type ServiceData struct {
    Name        string
    GoType      string
    JSName      string
    Methods     []MethodData
    PackagePath string
}

type MethodData struct {
    Name           string
    JSName         string     // customizable via config
    GoFuncName     string     // internal Go function name
    ShouldGenerate bool       // based on include/exclude filters
    RequestType    string
    ResponseType   string
    Comment        string
}
```

## Generated Output Structure

```
gen/
├── go/                              # from protoc-gen-go-grpc
│   ├── library/v1/library_grpc.pb.go
│   └── user/v1/user_grpc.pb.go
├── ts/                              # from user's TS generator
│   ├── library/v1/library_pb.ts
│   └── user/v1/user_pb.ts
└── wasm/                            # from protoc-gen-go-wasmjs
    ├── bookstore_services.wasm.go   # Go WASM wrapper
    ├── BookstoreServicesClient.ts   # TS client
    ├── build.sh                     # Build script
    └── compiled/
        ├── bookstore_services.wasm  # Compiled WASM
        └── wasm_exec.js            # Go WASM runtime
```

## Success Criteria
- [ ] Generate multi-service WASM with clean namespaced APIs
- [ ] TypeScript clients work seamlessly with protoc-gen-es types
- [ ] Template customization works for real-world scenarios
- [ ] Build pipeline integrates smoothly with buf workflows
- [ ] Performance is competitive with manual WASM implementations

## Key Deliverables
1. **PLAN.md** - ✅ This comprehensive plan with progress tracking
2. **Complete generator rewrite** with template flexibility
3. **Multi-service example** demonstrating all features
4. **Documentation** for configuration and customization
5. **Migration tools** from existing solutions

---

## Notes & Decisions

### Implementation Decisions
- **Namespaced APIs**: Default to clean `namespace.service.method()` structure vs flat function names
- **Multi-service bundling**: Enable related services in single WASM module for efficiency
- **Template flexibility**: Support both configuration-driven and template override customization
- **Type system integration**: Auto-detect and work with multiple TypeScript protobuf generators

### Performance Considerations
- **Bundle size**: Multi-service modules reduce overhead vs one WASM per service
- **Method filtering**: Allow selective generation to reduce bundle size
- **Type conversion**: Minimize JSON marshaling overhead with direct conversion methods

### Developer Experience
- **buf.build integration**: Work seamlessly with existing buf workflows
- **Build automation**: Generate build scripts and runtime file management
- **Error handling**: Provide clear error messages and debugging capabilities