# Next Steps for protoc-gen-go-wasmjs

## ✅ Completed Work (Phase 4)

### Split Architecture Successfully Implemented
- **Go Generator (`protoc-gen-go-wasmjs-go`)**: Fully functional, generates WASM wrappers
- **TypeScript Generator (`protoc-gen-go-wasmjs-ts`)**: Fully functional, generates TypeScript clients
- **Both generators tested** with simple and complex examples

### Critical Bug Fixes Applied

#### Go Generator Issues Fixed:
1. **Missing wasm package import** - Now always included for WASM generation
2. **Empty request/response types** - Properly qualified with package aliases
3. **Protobuf wire protocol corruption** - Removed all stdout writes
4. **Template execution verification** - Added proper error handling

#### TypeScript Generator Issues Fixed:
1. **Template data structure mismatches** - Created TypeScript-specific data structures
2. **Missing fields** - Added all required fields (APIStructure, JSNamespace, etc.)
3. **Boolean logic errors** - Fixed template conditionals with HasMessages/HasEnums flags
4. **Stdout corruption** - Removed fmt.Printf calls

## ✅ **CRITICAL ISSUES RESOLVED** (September 2025)

### ✅ **1. Template Inheritance Issues (CRITICAL) - RESOLVED**
- **✅ Base class properties**: `wasmLoadPromise`, `browserServiceManager` properly accessible
- **✅ Base class methods**: `loadWasm`, `registerBrowserService`, `callMethod` all working
- **✅ Inheritance chain**: Generated clients properly extend `WASMServiceClient`
- **✅ Runtime package integration**: `@protoc-gen-go-wasmjs/runtime` imports working correctly
- **Issue**: Was a Vite dev server caching problem - resolved by restarting dev server

### ✅ **2. Per-Service Client Generation (ARCHITECTURAL IMPROVEMENT) - IMPLEMENTED**
- **✅ Separate client files**: Each service generates its own client file
- **✅ Directory structure**: Follows proto package hierarchy (`presenter/v1/presenterServiceClient.ts`)
- **✅ No file conflicts**: Eliminates overwriting issues from multiple services
- **✅ Clean organization**: Browser services and WASM services properly separated
- **✅ Comprehensive tests**: Unit tests and integration tests for new functionality

### ✅ **3. Browser Service Communication (CRITICAL) - FIXED**
- **✅ Main thread blocking**: Fixed with `async_method` annotations preventing deadlocks
- **✅ Protobuf deserialization**: Fixed pointer instantiation in `CallBrowserService`
- **✅ JSON → JS object conversion**: Go now passes proper JavaScript objects to callbacks
- **✅ End-to-end functionality**: Example callbacks working with prompts, localStorage, etc.

### ✅ **4. TypeScript Development Environment (HIGH) - WORKING**
- **✅ Vite integration**: Modern TypeScript compilation and hot reload
- **✅ pnpm workspace**: Runtime package properly linked as workspace dependency
- **✅ TypeScript compilation**: All generated code compiles without errors
- **✅ Runtime package**: Clean inheritance-based architecture working

### ✅ **5. Bundle Naming Issue (CRITICAL) - RESOLVED** (September 2025)
- **✅ Root cause identified**: Line 223 in `TSDataBuilder.BuildServiceClientData` used package names instead of configured module_name
- **✅ Fix implemented**: Updated to use `tb.getModuleName(packageInfo.Name, config)` method
- **✅ Correct behavior**: Both `presenter.v1` and `browser.v1` packages now generate `exampleBundle`
- **✅ Configuration usage**: Proper usage of `module_name=example` parameter from buf.gen.yaml
- **✅ Tests updated**: Integration tests and examples updated to reflect correct naming
- **✅ Regression prevention**: Added debug tests to prevent future occurrences

### ✅ **6. Cross-Package Type Imports (CRITICAL) - RESOLVED** (October 2025)
- **✅ Issue**: Missing imports for types from other proto packages in same project
- **✅ Root cause**: Import collection logic only handled well-known types, not cross-package message types
- **✅ Fix implemented**:
  - Uses protobuf descriptor API (`field.Message.Desc.FullName()`, `ParentFile().Package()`) instead of string parsing
  - Added `MessagePackage` and `IsNestedType` fields to `TSFieldInfo` for accurate metadata
  - Fixed `MessageCollector` to use `Desc.FullName()` for correct fully qualified names
- **✅ Nested type support**: Properly flattens nested types (e.g., `ParentMessage_NestedType`) to avoid name collisions
- **✅ Relative imports**: Correctly calculates relative import paths (e.g., `../../utils/v1/interfaces`)
- **✅ Tests added**: Comprehensive unit tests for package extraction and type name flattening
- **✅ Example verification**: `example` example now correctly imports `HelperUtilType` and `ParentUtilMessage_NestedUtilType`

### ✅ **7. Factory/Deserializer/Models Generation (ARCHITECTURAL) - COMPLETED** (October 2025)
- **✅ Issue**: TODO at line 179 in `ts_generator.go` - factory/models/schemas/deserializer files not being generated
- **✅ Root cause**: New catalog-based `planFilesFromCatalog()` method incomplete - had TODO comment instead of implementation
- **✅ Fix implemented**:
  - Completed file planning for models, factory, schemas, and deserializer in `planFilesFromCatalog()`
  - Added rendering logic in `renderFilesFromCatalog()` for all type artifact files
  - Factory generation respects `generate_factories=true` configuration flag
  - Added package deduplication to avoid generating same files multiple times
  - Implemented caching of `TSTemplateData` to avoid rebuilding for each file type
- **✅ Generated files per package**:
  - `interfaces.ts` - Pure TypeScript interfaces (always generated)
  - `models.ts` - Concrete class implementations with defaults (always generated)
  - `factory.ts` - Object construction with context awareness (when `generate_factories=true`)
  - `schemas.ts` - Field metadata for runtime introspection (always generated)
  - `deserializer.ts` - Schema-driven data population (always generated)
- **✅ Architecture benefits**:
  - Clean separation: interfaces for types, models for implementations
  - Factory system handles proper default values and recursive construction
  - Deserializer uses schema metadata for type-safe field resolution
  - All existing tests pass with new generation logic

### ✅ **8. pkg.go.dev Documentation (DOCUMENTATION) - COMPLETED** (October 2025)
- **✅ Comprehensive godoc added**: All packages now have complete documentation
- **✅ Files created**:
  - Root `doc.go` - Complete project overview with architecture diagrams
  - `pkg/generators/doc.go` - Generation architecture and artifact processing
  - `pkg/core/doc.go` - Pure utilities (NameConverter, PathCalculator, ProtoAnalyzer)
  - `pkg/wasm/doc.go` - WASM runtime and browser service channels
  - `pkg/builders/doc.go` - Template data building and file planning
- **✅ Detailed type documentation**:
  - All exported types in `pkg/generators/base_generator.go`
  - All exported types in `pkg/wasm/browser_channel.go`
  - Comprehensive examples throughout
- **✅ Command documentation**:
  - `cmd/protoc-gen-go-wasmjs-go/main.go` - Complete Go generator guide
  - `cmd/protoc-gen-go-wasmjs-ts/main.go` - Complete TypeScript generator guide
- **✅ Documentation features**:
  - Rich examples with Library service
  - ASCII architecture diagrams
  - Complete configuration reference
  - Cross-package linking
  - Usage patterns and best practices
  - Error handling documentation

## ✅ **Recent Improvements** (January 2025)

### **Go WASM File Splitting - COMPLETED**
Split monolithic Go WASM generation into 3 modular files:

**Generated Files:**
1. `<package>_converters.wasm.go` - JS/Go type converters and stream wrappers
   - `createJSResponse()` helper used by all service methods
   - Server stream wrapper implementations
   - Generated whenever services exist

2. `<package>_exports.wasm.go` - Service exports and registration
   - Exports struct for dependency injection
   - `RegisterAPI()` method for JS namespace registration
   - Service method wrappers

3. `<package>_browser_clients.wasm.go` - Browser service clients
   - Browser service client implementations
   - Only generated when browser-provided services exist

**Benefits:**
- **Better organization** - Clear separation of concerns
- **Easier maintenance** - Each file has focused responsibility
- **Selective generation** - Only generate what's needed
- **Foundation for future improvements** - Easier to add non-wasm utilities later

### ✅ **9. go_package Output Path Collision (CRITICAL) - RESOLVED** (November 2025)
- **✅ Issue**: When multiple proto files had same proto package but different `go_package` options, generator was invoked multiple times and files collided/overwrote each other
- **✅ Root cause**: Output path calculation only used proto package name, ignoring `go_package` differences
- **✅ Fix implemented**:
  - Added `calculateOutputPath()` method that extracts path from `go_package` option
  - Added `calculateBaseName()` method that incorporates go_package suffix for uniqueness
  - Updated `planGoFiles()` to use go_package-aware path calculation
- **✅ Test cases added**:
  - `test_one_package/` - Same proto package, different go_package (models vs services)
  - `test_multi_packages/` - Different proto packages (test_multi_packages.v1.models vs test_multi_packages.v1.services)
  - `test_broken/` - Original failing case now working
- **✅ Benefits**:
  - Supports standard Go pattern: separate models and services packages to avoid gRPC dependencies
  - Proper file organization: `.../v1/models/...exports.wasm.go` and `.../v1/services/...exports.wasm.go`
  - No more file collisions when using multiple go_package options
  - Works with both same and different proto package patterns

### ✅ **10. Annotation-Based Factory Generation (ARCHITECTURAL) - COMPLETED** (November 2025)
- **✅ Issue**: Automatic package-level factory generation was too complex with multiple buf invocations and caused duplicate schema files
- **✅ Root cause**: Trying to aggregate across multiple protoc invocations with automatic detection
- **✅ New approach - Explicit annotation-based generation**:
  - Added `(wasmjs.v1.ts_factory) = true` file-level annotation in `annotations.proto`
  - User creates dedicated `factory.proto` file with imports of desired messages
  - Factory file location determines output location (e.g., `x/y/z/factory.proto` → `x/y/z/factory.ts`)
  - Only imports from same package are included (different packages filtered out)
- **✅ Fix implemented**:
  - Added `IsTypeScriptFactoryFile()` detection in `ProtoAnalyzer`
  - Implemented `collectFactoryArtifacts()` and `collectImportedMessages()` in `TSGenerator`
  - Created `BuildFactoryData()` in `TSDataBuilder` with correct relative import paths
  - Fixed Generate flag checks to prevent external dependencies (google/api) from being generated
- **✅ Generated output**:
  - **factory.ts** - Combined factory + deserializer in single file
  - **schemas.ts** - Aggregated schemas from all imported message directories
  - Proper import aliasing: interfaces as `TypeInterface`, models as `ConcreteType`
  - Relative imports correctly calculated from factory location to message directories
- **✅ Benefits**:
  - **Simple and predictable**: User explicitly controls what goes in factory
  - **Works across buf invocations**: No dependency on invocation boundaries
  - **Clean imports**: Automatically calculates correct relative paths
  - **Schema consolidation**: Aggregates directory-level schemas with unique aliases
  - **Single file**: Combined factory + deserializer reduces import complexity

### ✅ **11. InvocationStyle System (December 2025) - COMPLETED**
- **✅ Issue**: Ambiguous `async_method` annotation didn't clearly distinguish between different async patterns
- **✅ Solution**: Explicit `invocation_style` annotation with SYNC, PROMISE, and CALLBACK options
- **✅ Files created/modified**:
  - `pkg/wasm/promise.go` (NEW) - Go helpers for creating JavaScript Promises
  - `proto/wasmjs/v1/annotations.proto` - InvocationStyle enum
  - `pkg/core/proto_analyzer.go` - InvocationStyle type and GetInvocationStyle method
  - `pkg/builders/shared_types.go` - MethodData with InvocationStyle field
  - `pkg/renderers/templates/wasm_exports.go.tmpl` - PROMISE branch using CreateJSPromise
  - `pkg/renderers/templates/client_simple.ts.tmpl` - Updated method signatures
  - `runtime/src/client/service-client.ts` - callMethodSync/callMethodAsync
  - `runtime/src/client/wasm-bundle.ts` - callMethodSync/callMethodAsync
- **✅ Dead code removed**:
  - `pkg/generator/generator.go`, `templates.go`, `tsgenerator.go`, `types.go` - Deleted
  - `pkg/renderers/templates/wasm.go.tmpl` - Deleted
- **✅ Benefits**:
  - Clear semantics for each invocation style
  - PROMISE style allows async/await in TypeScript
  - Prevents deadlocks when calling browser APIs from Go WASM
  - Fully typed method signatures in generated TypeScript

## 🚀 **NEXT PHASE: Production Polish**

### **Phase 2: Streaming Improvements (Priority: MEDIUM)**
Enhance streaming support with better error handling:
- Typed streaming callbacks with proper response types
- Better cancellation support
- Progress reporting for long-running streams

### **Phase 3: Performance Optimization (Priority: LOW)**
- Consider SharedArrayBuffer for better WASM/JS data transfer
- Template caching for faster regeneration
- Parallel generation for large projects

## 📋 Medium-Term Goals

### 1. Feature Parity
- Ensure all features from monolithic generator work in split version
- Streaming support verification
- Custom template support
- All filtering options working

### 2. Migration Guide
- Document migration from old to new generators
- Create compatibility wrapper if needed
- Update all examples to use new generators

### 3. Documentation
- Complete API documentation for all layers
- Usage guide for split generators
- Template customization guide
- Architecture decision records (ADRs)

### 4. CI/CD Integration
- Add GitHub Actions for testing
- Automated release process
- Version compatibility matrix

## 🔮 Long-Term Vision

### 1. Additional Language Support
- Consider Python WASM generator
- Rust WASM generator
- Other target languages as needed

### 2. Advanced Features
- Hot reload support for development
- Source map generation
- Advanced debugging capabilities
- Performance profiling tools

### 3. Ecosystem Integration
- Buf Schema Registry integration
- VS Code extension
- Build tool plugins (webpack, vite, etc.)
- Framework integrations (React, Vue, Angular)

## 📊 Technical Debt to Address

1. **TODO Comments**: Several TODO items in code need addressing:
   - Field analysis in TSDataBuilder
   - External imports in TSDataBuilder
   - Oneof group analysis
   - Browser service detection

2. **Template Consolidation**:
   - Review if templates can be simplified
   - Consider template inheritance/composition

3. **Error Messages**:
   - Improve error messages with more context
   - Add suggestions for common issues

4. **Performance**:
   - Consider parallel generation for large projects
   - Template caching for repeated use
   - Memory usage optimization

## ✅ Success Criteria for Production Release

- [ ] All examples generate and run successfully
- [ ] Performance within 10% of old generator
- [ ] Comprehensive test coverage (>80%)
- [ ] Documentation complete
- [ ] Migration guide available
- [ ] No known critical bugs
- [ ] Community feedback incorporated

## 📝 Notes

The refactoring to a split architecture has been successful. The new design provides:
- Better testability through layer separation
- Cleaner code organization
- Easier maintenance and extension
- Language-specific optimizations
- Improved error handling

The generators are now ready for broader testing and community feedback before the production release.
