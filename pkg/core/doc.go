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
Package core provides pure utility functions for protobuf code generation.

# Overview

The core package implements stateless, side-effect-free utility functions
for common operations in protocol buffer code generation. All functions in
this package are deterministic and testable without external dependencies.

# Design Philosophy

The core package follows these principles:

  - Pure Functions: No side effects, deterministic outputs
  - No External State: All required data passed as parameters
  - Comprehensive Testing: 30+ unit tests with 100% coverage
  - Clear Separation: Distinct utilities for different concerns

# Components

The package provides three main utility types:

NameConverter

Handles naming convention transformations between different languages:

  - Proto naming → Go naming (PascalCase)
  - Go naming → JavaScript naming (camelCase)
  - Go naming → TypeScript naming
  - Special cases: acronyms, reserved words, flattening

PathCalculator

Calculates import paths and file paths:

  - Proto package → directory path
  - Cross-package import path calculation
  - Relative path resolution
  - TypeScript import path generation

ProtoAnalyzer

Analyzes protobuf definitions for generation decisions:

  - Service detection (regular vs browser-provided)
  - Method streaming detection (server, client, bidirectional)
  - Async method detection
  - Annotation inspection
  - Well-known type detection

WellKnownTypes

Handles Google's well-known protobuf types:

  - Type detection (google.protobuf.*)
  - TypeScript type mapping
  - Import path resolution
  - Special handling for Timestamp, Duration, Any, etc.

# Usage Example

Name Conversion:

	import "github.com/panyam/protoc-gen-go-wasmjs/pkg/core"

	converter := core.NewNameConverter()

	// Convert Go method name to JavaScript
	jsName := converter.ToCamelCase("FindBooks")  // "findBooks"

	// Flatten nested type names
	flatName := converter.FlattenTypeName("Author.Address")  // "Author_Address"

	// Handle acronyms properly
	name := converter.ToTitleCase("http_client")  // "HTTPClient"

Path Calculation:

	calculator := core.NewPathCalculator()

	// Convert proto package to directory path
	path := calculator.BuildPackagePath("library.v1")  // "library/v1"

	// Calculate relative import path
	rel := calculator.CalculateRelativeImportPath(
	    "presenter/v1",    // from this package
	    "utils/v1"         // to this package
	)  // "../../utils/v1"

Proto Analysis:

	analyzer := core.NewProtoAnalyzer()

	// Check if service has browser_provided annotation
	isBrowser := analyzer.IsBrowserProvidedService(service)

	// Detect streaming type
	streamType := analyzer.GetStreamingType(method)
	switch streamType {
	case core.StreamingTypeServerStreaming:
	    // Generate streaming handler
	case core.StreamingTypeUnary:
	    // Generate unary handler
	}

	// Check if method has async annotation
	isAsync := analyzer.IsAsyncMethod(method)

Well-Known Types:

	// Check if type is well-known
	if core.IsWellKnownType("google.protobuf.Timestamp") {
	    // Use special TypeScript handling
	    tsType := "Date"
	}

	// Get TypeScript type for well-known type
	tsType := core.GetTypeScriptTypeForWellKnown("google.protobuf.Duration")
	// Returns: "string"  (ISO 8601 duration format)

# Architectural Role

The core package sits at the bottom of the dependency hierarchy:

	┌─────────────────┐
	│   generators    │  (uses core utilities)
	└────────┬────────┘
	         │
	┌────────▼────────┐
	│    builders     │  (uses core utilities)
	└────────┬────────┘
	         │
	┌────────▼────────┐
	│    filters      │  (uses core utilities)
	└────────┬────────┘
	         │
	┌────────▼────────┐
	│      core       │  (no dependencies on other pkg/*)
	└─────────────────┘

This ensures:

  - No circular dependencies
  - Easy unit testing
  - Reusable across different generators
  - Clear separation of concerns

# Testing

The core package has comprehensive test coverage:

  - name_converter_test.go: 15+ test cases
  - path_calculator_test.go: 12+ test cases
  - proto_analyzer_test.go: 10+ test cases
  - All edge cases covered
  - Clear test names describing behavior

Run tests:

	go test ./pkg/core/...

# Performance

All core utilities are designed for performance:

  - No allocations in hot paths where possible
  - String builders for concatenation
  - Early returns for common cases
  - Cached computations where appropriate

# Links

Related packages:

  - github.com/panyam/protoc-gen-go-wasmjs/pkg/generators: Uses core for generation
  - github.com/panyam/protoc-gen-go-wasmjs/pkg/builders: Uses core for data building
  - github.com/panyam/protoc-gen-go-wasmjs/pkg/filters: Uses core for filtering
*/
package core
