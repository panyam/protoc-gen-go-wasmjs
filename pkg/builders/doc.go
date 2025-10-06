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
Package builders provides template data construction for code generation.

# Overview

The builders package sits between the generators layer and the renderers layer,
transforming filtered artifacts into structured data that templates can consume.
This package implements the "build" phase of the generation pipeline, where raw
protobuf definitions are converted into language-specific template data.

# Architecture Role

The builders package is part of the 3-layer generation architecture:

	┌─────────────────┐
	│   Generators    │  Orchestrate generation flow
	└────────┬────────┘
	         │
	┌────────▼────────┐
	│    Builders     │  Transform artifacts → template data
	└────────┬────────┘
	         │
	┌────────▼────────┐
	│   Renderers     │  Execute templates with data
	└─────────────────┘

This separation enables:

  - Clear data transformation logic
  - Reusable builders across generators
  - Easy testing of template data
  - Generator-specific customization

# Core Components

GenerationConfig

Shared configuration structure used by all generators:

  - Export paths (WASM, TypeScript)
  - JavaScript API structure (namespaced, flat, service-based)
  - Generation control flags (clients, types, factories)
  - Build integration options

GoDataBuilder

Builds template data for Go WASM wrapper generation:

  - Processes services and methods
  - Calculates Go imports and package aliases
  - Builds ServiceData with Go-specific naming
  - Handles browser service detection

TSDataBuilder

Builds template data for TypeScript client generation:

  - Processes services, messages, enums
  - Calculates TypeScript imports and type names
  - Builds per-service client data
  - Handles cross-package imports
  - Generates factory and schema data

PackageInfo

Metadata about a proto package:

  - Package name (e.g., "library.v1")
  - Directory path (e.g., "library/v1")
  - Proto files in this package
  - Used for organizing generated files

# Data Structures

Template data structures represent processed artifacts ready for rendering:

ServiceData

Represents a service ready for code generation:

  - Service names (Go, JavaScript, custom)
  - Package metadata (path, alias)
  - Browser service flag
  - Processed method list

MethodData

Represents a method ready for code generation:

  - Method names (original, JavaScript, Go function)
  - Request/response types with proper imports
  - Streaming detection
  - Async flag
  - Type information

MessageInfo / EnumInfo

Represent types for TypeScript generation:

  - Type names (proto, TypeScript, Go)
  - Package information
  - Field details (for messages)
  - Value details (for enums)

ImportInfo

Represents an import statement:

  - Full import path
  - Package alias
  - Used in Go and TypeScript generation

# Usage Example

Building Go template data:

	import (
	    "github.com/panyam/protoc-gen-go-wasmjs/pkg/builders"
	    "github.com/panyam/protoc-gen-go-wasmjs/pkg/generators"
	)

	// From generator
	config := &builders.GenerationConfig{
	    JSStructure: "namespaced",
	    JSNamespace: "myApp",
	    ModuleName:  "my_services",
	}

	// Create builder
	goBuilder := builders.NewGoDataBuilder(analyzer, pathCalc, nameConv)

	// Build template data for a service
	serviceData := goBuilder.BuildServiceData(
	    serviceArtifact.Service,
	    config,
	    packageInfo,
	)

	// serviceData is ready for template rendering
	// Contains: Name, GoType, JSName, Methods, etc.

Building TypeScript template data:

	tsBuilder := builders.NewTSDataBuilder(analyzer, pathCalc, nameConv)

	// Build service client data
	clientData := tsBuilder.BuildServiceClientData(
	    serviceArtifact,
	    catalog,
	    config,
	)

	// Build type artifact data (messages, enums)
	typeData := tsBuilder.BuildTypeArtifactData(
	    messageArtifact,
	    catalog,
	    config,
	)

	// Both ready for TypeScript template rendering

# File Planning

The builders package also handles file planning - deciding how artifacts
map to output files:

File Metadata:

  - FilePath: Where to write the file
  - TemplateName: Which template to use
  - Data: Template data to render

File Planning Process:

	1. Classify artifacts (services, messages, enums)
	2. Decide file mapping strategy:
	   - Go: One file per package (bundle all services)
	   - TypeScript: Multiple files per package (clients, types, schemas)
	3. Create FilePlanning with metadata for each file
	4. Return to generator for rendering

Example file planning:

	planning := &builders.FilePlanning{
	    Files: []builders.FileMetadata{
	        {
	            FilePath:     "library/v1/library_v1.wasm.go",
	            TemplateName: "wasm.go.tmpl",
	            Data:         goTemplateData,
	        },
	        {
	            FilePath:     "library/v1/interfaces.ts",
	            TemplateName: "interfaces.ts.tmpl",
	            Data:         tsTypeData,
	        },
	    },
	}

# Data Transformation Pipeline

The transformation from artifacts to template data follows this flow:

	Artifact (from generators)
	    │
	    ├─ Extract basic info (name, type, package)
	    ├─ Apply naming conventions (camelCase, PascalCase)
	    ├─ Calculate imports (cross-package references)
	    ├─ Detect special cases (streaming, async, browser)
	    ├─ Build request/response metadata
	    └─ Assemble into ServiceData/MethodData
	    │
	    ▼
	Template Data (to renderers)

# Type Name Resolution

The builders handle complex type name resolution:

Go Types:

  - Fully qualified: "libraryv1.FindBooksRequest"
  - With imports: Import{Path: "...", Alias: "libraryv1"}
  - Handles nested messages and well-known types

TypeScript Types:

  - Interface names: "FindBooksRequest"
  - Import paths: "../../library/v1/interfaces"
  - Cross-package references
  - Well-known type mapping (Timestamp → Date)

# Browser Service Handling

The builders detect and handle browser-provided services:

Detection:

  - Check for (wasmjs.v1.browser_provided) = true
  - Set IsBrowserProvided flag in ServiceData
  - Generate different code paths for browser services

Client Generation:

  - Browser services: Generate TypeScript-to-WASM stubs
  - Regular services: Generate WASM-to-TypeScript clients

# Streaming Support

The builders detect and annotate streaming methods:

  - Server streaming: stream in response
  - Client streaming: stream in request
  - Bidirectional: streams in both
  - Generates appropriate TypeScript signatures

# Testing

The package has comprehensive test coverage:

  - file_planning_test.go: File mapping strategies
  - ts_data_builder_test.go: TypeScript data building
  - typed_callbacks_test.go: Callback type generation
  - bundle_naming_*.go: Module naming logic

Run tests:

	go test ./pkg/builders/...

# Configuration

GenerationConfig provides fine-grained control:

	config := &builders.GenerationConfig{
	    // Where to write files
	    WasmExportPath: "./gen/wasm",
	    TSExportPath:   "./web/src/generated",

	    // JavaScript API structure
	    JSStructure: "namespaced",  // or "flat", "service_based"
	    JSNamespace: "myApp",
	    ModuleName:  "my_services",

	    // What to generate
	    GenerateClients:   true,  // TypeScript clients
	    GenerateTypes:     true,  // Interfaces/models
	    GenerateFactories: true,  // Factory classes

	    // Build integration
	    GenerateBuildScript: true,
	    WasmPackageSuffix:   "wasm",
	}

# Links

Related packages:

  - github.com/panyam/protoc-gen-go-wasmjs/pkg/generators: Uses builders for data
  - github.com/panyam/protoc-gen-go-wasmjs/pkg/renderers: Consumes builder output
  - github.com/panyam/protoc-gen-go-wasmjs/pkg/core: Provides naming utilities
  - github.com/panyam/protoc-gen-go-wasmjs/pkg/filters: Provides filtered artifacts
*/
package builders
