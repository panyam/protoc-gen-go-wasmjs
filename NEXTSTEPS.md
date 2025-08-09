# Next Steps

## Recently Completed (January 2025)

- ✅ **Connect4 Example Restoration & Documentation Overhaul**: Fixed corrupted demo and aligned docs with reality
  - **Independent Module Setup**: Added go.mod for standalone Connect4 example with proper parent module replacement
  - **WASM Integration Fixes**: Corrected protobuf enum constants, import paths, and struct references for working WASM compilation
  - **Documentation Accuracy**: Major rewrite of all .md files to reflect actual working implementation vs outdated claims
  - **Transport Architecture Reality**: Updated docs to show working IndexedDB+polling and BroadcastChannel vs non-existent WebSocket server
  - **Build Process Alignment**: Fixed all Makefiles, build instructions, and file path references to match webpack+TypeScript reality
  - **Working Demo**: Cross-tab multiplayer Connect4 with state persistence and pluggable transports now fully functional

- ✅ **Enhanced Factory & Deserialization System**: Completed comprehensive factory composition and schema-aware deserialization
  - **Context-Aware Factory Methods**: Implemented parent object tracking, attribute names, and container keys for granular control
  - **Cross-Package Factory Composition**: Automatic dependency detection, import generation, and factory delegation across package boundaries
  - **Schema-Aware Deserialization**: Generated schema files with field metadata, proto field IDs, and oneof support for type-safe runtime processing
  - **Package-Scoped Registries**: Conflict-free multi-version support with fully qualified messageType names
  - **Factory-Deserializer Integration**: Seamless delegation between factory creation and deserializer population
  - **Production Testing**: 100% validation success with complex nested objects and real-world scenarios
  - **Enhanced Client Integration**: Demonstration client using the new deserializer system

- ✅ **Major Architecture Simplification**: Successfully completed TypeScript architecture transformation
  - **Self-Contained TypeScript Generation**: Eliminated dependencies on external TypeScript generators (protoc-gen-es, protoc-gen-ts)
  - **Simplified Client Architecture**: Replaced complex conversion system with direct JSON serialization
  - **Template-Based Generation**: Implemented interfaces, models, and factory generation using Go templates
  - **Configuration Cleanup**: Removed obsolete fields and streamlined configuration options
  - **Default Value Handling**: Fixed array and message type defaults with proper optional field support
  - **End-to-End Testing**: Validated complete architecture with working examples
  - **Performance Improvement**: Eliminated ~200 lines of complex conversion logic

- ✅ **Quality & TypeScript Refinements**: Latest round of critical bug fixes and improvements
  - **Native Map Type Support**: Fixed proto `map<K,V>` fields to generate native TypeScript `Map<K,V>` types instead of synthetic interfaces
  - **Framework Schema Architecture**: Separated framework types (`FieldType`, `FieldSchema`, `MessageSchema`) into dedicated `deserializer_schemas.ts` files
  - **Package-Based Generation**: Completed transition from file-based to package-based TypeScript generation eliminating import conflicts
  - **Type Safety Improvements**: Fixed factory method subscripting (`this[methodName]`) and `FactoryInterface` compatibility issues
  - **Build System Stability**: Resolved all TypeScript compilation errors and maintained 100% build success rate

## Recently Completed (August 2025)
- ✅ **Symlink Elimination & BSR Integration**: Successfully resolved protobuf dependency management for both examples
  - **Published wasmjs Protos**: Published `wasmjs/v1/annotations.proto` to buf.build registry at `buf.build/panyam/protoc-gen-go-wasmjs`
  - **Production Mode**: Both examples now use published wasmjs proto dependencies + local plugin installation
  - **Development Mode**: Clean development workflow with local symlinks and buf.lock management
  - **Dual Configuration**: Separate buf.yaml/buf.gen.yaml files for production vs development workflows
  - **User Experience**: End users no longer need symlink management - just add one dependency line and install plugin
  - **Library Example**: Complete elimination of symlink requirement with comprehensive Makefile targets
  - **Connect4 Example**: Applied same symlink elimination pattern with real-time multiplayer focus
  - **Documentation**: Created SETUP.md guides for both examples explaining production vs development workflows

## Immediate Tasks (Current Priority)
- ✅ **External Package Import Support**: Implemented comprehensive external type mapping system
  - **Google Protobuf Types**: Full support for `google.protobuf.Timestamp` → `Date`, `google.protobuf.FieldMask` → `string[]`
  - **Import Generation**: Automatic detection and exclusion of external types from factory dependencies
  - **Type Mapping**: Configurable mapping system with default mappings for well-known types
  - **Factory Integration**: Table-driven `externalTypeFactories` with `newXYZ`/`serializeXYZ` methods
  - **Package Resolution**: Proper handling of well-known types without generating non-existent imports

- ✅ **Enhanced Developer Experience**: Implemented ergonomic API improvements for type-safe deserialization
  - **MESSAGE_TYPE Constants**: Each message class has `static readonly MESSAGE_TYPE` with fully qualified name
  - **Static Deserializer Method**: `MyDeserializer.from<T>(messageType, data)` for convenient deserialization
  - **Optional Constructor Parameters**: Deserializer constructor with default factory and schema registry
  - **Shared Factory Instance**: Performance-optimized singleton factory to avoid unnecessary instantiation

- ✅ **Bug Fixes & Enum Support**: Latest critical fixes for production stability
  - **wasmjs.v1 Package Filtering**: Fixed artifact generation for wasmjs annotation packages using package name detection in main.go (lines 94-97)
  - **Comprehensive Enum Support**: Implemented complete enum collection and generation system
    - Added EnumInfo and EnumValueInfo types to represent proto enums
    - Added collectAllEnums() function to gather enums from all proto files
    - Updated generation logic to handle packages with enums but no messages
    - Enhanced TypeScript templates to generate and import enums correctly
  - **Cross-Package Import Filtering**: Enhanced import detection to exclude wasmjs.v1 from factory dependencies
  - **Template Import Resolution**: Fixed enum imports in models.ts, factory.ts, and all generated TypeScript files
  
- [ ] **Enhanced Browser Demo**: Create a complete browser demo showcasing the enhanced factory and deserialization system
  - Demonstrate cross-package factory composition with real dependencies
  - Show schema-aware deserialization with complex nested objects
  - Include examples of context-aware factory methods with external types
  - Performance comparisons with previous systems
  - Real-world scenario simulation (library management system)
  
- [ ] **Performance Analysis**: Benchmark the enhanced factory and deserialization system
  - Measure factory composition overhead vs direct creation
  - Analyze schema-aware deserialization performance
  - Compare cross-package delegation efficiency
  - Document complex object creation performance characteristics

- [x] **Documentation Refresh**: Updated documentation to reflect current architecture and capabilities
  - ✅ Connect4 example docs completely rewritten to show actual working implementation
  - ✅ Fixed all architecture diagrams to match IndexedDB+polling transport reality  
  - ✅ Corrected all build instructions and file path references
  - ✅ Updated README examples to show actual TypeScript client usage patterns
  - [ ] Create comprehensive examples showing enhanced factory patterns and cross-package composition
  - [ ] Document the complete interface/model/factory/schema/deserializer ecosystem with real examples

## Short Term (Next Month)
- [ ] **Streaming Support**: Research and implement streaming RPC support for WASM
  - Server-streaming RPCs
  - Client-streaming RPCs
  - Bidirectional streaming
  
- [ ] **Advanced Templates**: Template inheritance and partial overrides
  - Allow extending base templates
  - Support for custom template functions
  - Template composition for complex scenarios

- [ ] **Error Recovery**: Implement graceful error recovery
  - Retry mechanisms for transient failures
  - Circuit breaker pattern
  - Fallback to HTTP when WASM fails

## Medium Term (Next Quarter)
- [ ] **Monitoring Integration**: Performance monitoring and analytics
  - OpenTelemetry integration
  - Custom metrics for WASM performance
  - Distributed tracing support
  
- [ ] **IDE Support**: Language server and IDE plugin support
  - VS Code extension for proto to WASM development
  - IntelliJ IDEA plugin
  - Syntax highlighting for WASM annotations

- [ ] **Advanced Generation Features**:
  - Custom template functions and helpers
  - Template inheritance and composition
  - Proto extension support for custom annotations

## Long Term (6+ Months)
- [ ] **Multi-Language Support**: 
  - Rust service implementations with same TypeScript clients
  - C++ WASM generation
  - Shared TypeScript client generation

- [ ] **Edge Computing Optimization**:
  - CDN-friendly WASM deployment
  - Edge worker compatibility
  - Lazy loading strategies

- [ ] **Enterprise Features**:
  - Advanced authentication patterns
  - Authorization middleware
  - Audit logging
  - Compliance features (GDPR, HIPAA)

## Community & Ecosystem
- [ ] **Documentation Improvements**:
  - Video tutorials
  - Example repository with common patterns
  - Migration guides from other solutions
  
- [ ] **Framework Integration**:
  - React hooks for WASM services
  - Vue composables
  - Angular services
  - Svelte stores

- [ ] **Testing Tools**:
  - Mock generation for WASM services
  - Integration testing framework
  - Performance testing suite

## Technical Debt
- [ ] **Code Cleanup**:
  - Add more unit tests for new TypeScript generation logic
  - Improve error messages throughout the generation pipeline
  - Refactor template data structures for better extensibility

- [ ] **Performance Optimizations**:
  - Lazy WASM loading strategies
  - Optimize generated TypeScript bundle size
  - Implement caching for repeated generation tasks

## Research Topics
- [ ] **WebAssembly Component Model**: Investigate integration with WASM Component Model
- [ ] **SharedArrayBuffer**: Explore using SharedArrayBuffer for better performance
- [ ] **WASM SIMD**: Leverage WASM SIMD for data processing
- [ ] **Module Federation**: Integration with Webpack Module Federation