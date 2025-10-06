// Copyright 2025 Sri Panyam
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

/*
Package generators provides the core generation logic for protoc-gen-go-wasmjs.

# Overview

The generators package implements a layered architecture for code generation with
a BaseGenerator foundation that separates artifact collection from file generation.
This design enables flexible, target-specific generation while maintaining a consistent
approach to artifact discovery and classification.

# Architecture

The package implements a 4-step artifact processing approach:

	1. COLLECT    → BaseGenerator.CollectAllArtifacts()
	2. CLASSIFY   → ArtifactCatalog (Services, Messages, Enums by package)
	3. MAP        → planFilesFromCatalog() (generator-specific logic)
	4. CREATE     → fileSet.CreateFiles() (send to protogen)

This separation allows:

  - Cross-package artifact visibility (regardless of protoc Generate flags)
  - Flexible file mapping (N artifacts → 1 file, or 1 artifact → N files)
  - Delayed file creation (after all mapping decisions are complete)
  - Generator-specific slice/dice/group logic

# Generator Types

BaseGenerator

The foundation component providing artifact collection for all generators.
It embeds core utilities (ProtoAnalyzer, PathCalculator, NameConverter) and
filter layers (PackageFilter, ServiceFilter, MethodFilter, collectors).

Both TSGenerator and GoGenerator embed BaseGenerator to access the complete
artifact catalog.

TSGenerator

TypeScript-specific generator that produces:

  - Service clients (per-service, package-level files)
  - Base bundle class (module-level)
  - TypeScript interfaces (type definitions)
  - Concrete models (implementations with defaults)
  - Factories (object construction, when generate_factories=true)
  - Schemas (field metadata)
  - Deserializers (schema-driven data population)

GoGenerator

Go WASM-specific generator that produces:

  - WASM wrappers (exports struct with RegisterAPI)
  - Example main.go files
  - Build scripts

# Artifact Catalog

The ArtifactCatalog provides a complete view of everything available for generation:

	type ArtifactCatalog struct {
	    Services        []ServiceArtifact   // Regular services
	    BrowserServices []ServiceArtifact   // Browser-provided services
	    Messages        []MessageArtifact   // Messages by package
	    Enums           []EnumArtifact      // Enums by package
	    Packages        map[string]*PackageInfo
	}

All generators work with this catalog to make file mapping decisions.

# Usage Example

Creating a generator:

	import (
	    "google.golang.org/protobuf/compiler/protogen"
	    "github.com/panyam/protoc-gen-go-wasmjs/pkg/generators"
	    "github.com/panyam/protoc-gen-go-wasmjs/pkg/builders"
	    "github.com/panyam/protoc-gen-go-wasmjs/pkg/filters"
	)

	func main() {
	    protogen.Options{}.Run(func(gen *protogen.Plugin) error {
	        // Create configuration
	        config := &builders.GenerationConfig{
	            JSStructure:  "namespaced",
	            JSNamespace:  "myApp",
	            ModuleName:   "my_services",
	        }

	        // Create filter criteria
	        filterCriteria := &filters.FilterCriteria{
	            // Configure service/method filtering
	        }

	        // Create TypeScript generator
	        tsGen := generators.NewTSGenerator(gen)

	        // Validate configuration
	        if err := tsGen.ValidateConfig(config); err != nil {
	            return err
	        }

	        // Generate TypeScript artifacts
	        if err := tsGen.Generate(config, filterCriteria); err != nil {
	            return err
	        }

	        return nil
	    })
	}

# Generator Workflow

Each generator follows this workflow:

	1. Create generator instance (embeds BaseGenerator)
	2. Validate configuration
	3. Call Generate(config, criteria):
	   a. Collect all artifacts via BaseGenerator
	   b. Classify into ArtifactCatalog
	   c. Plan files from catalog (generator-specific)
	   d. Render files from catalog
	   e. Create protogen files

# File Planning

Generator-specific file planning allows flexible artifact grouping:

TypeScript File Planning:

  - One bundle file per module (index.ts)
  - One client file per service (presenter/v1/presenterServiceClient.ts)
  - One interfaces file per package (presenter/v1/interfaces.ts)
  - One models file per package (presenter/v1/models.ts)
  - Optional factory file per package (presenter/v1/factory.ts)
  - One schemas file per package (presenter/v1/schemas.ts)
  - One deserializer file per package (presenter/v1/deserializer.ts)

Go File Planning:

  - One WASM wrapper per package
  - One example main.go per package
  - One build script per package

# Cross-Package Visibility

The BaseGenerator collects artifacts from ALL proto files, regardless of
protoc's Generate flags. This ensures generators have complete visibility
for making file mapping decisions, especially important for:

  - Cross-package type imports
  - Factory composition
  - Schema registries
  - Import path calculation

# Configuration Options

Generators accept GenerationConfig with these key options:

  - TSExportPath/WasmExportPath: Output directories
  - JSStructure: API structure (namespaced|flat|service_based)
  - JSNamespace: Global JavaScript namespace
  - ModuleName: Module name for generated code
  - GenerateClients: Generate service clients (TS only)
  - GenerateTypes: Generate interfaces/models (TS only)
  - GenerateFactories: Generate factory classes (TS only)
  - GenerateBuildScript: Generate build scripts (Go only)

# Filtering

Generators use FilterCriteria to control what gets generated:

  - Services: Specific services to include/exclude
  - Methods: Glob patterns for method inclusion/exclusion
  - Packages: Package-level filtering
  - Browser services: Automatic detection and separation

# Template Rendering

Both generators use the renderers package for template-based code generation:

  - TSRenderer: Renders TypeScript templates
  - GoRenderer: Renders Go WASM templates

Templates are embedded in the binary and can be overridden via configuration.

# Links

Related packages:

  - github.com/panyam/protoc-gen-go-wasmjs/pkg/builders: Template data building
  - github.com/panyam/protoc-gen-go-wasmjs/pkg/renderers: Template rendering
  - github.com/panyam/protoc-gen-go-wasmjs/pkg/filters: Artifact filtering
  - github.com/panyam/protoc-gen-go-wasmjs/pkg/core: Core utilities
*/
package generators
