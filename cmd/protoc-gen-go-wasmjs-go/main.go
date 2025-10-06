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
protoc-gen-go-wasmjs-go is a Protocol Buffers compiler plugin that generates Go WASM wrappers for gRPC services.

# Overview

This generator is part of the protoc-gen-go-wasmjs split architecture. It focuses exclusively
on generating Go WASM wrapper code, while the TypeScript generator (protoc-gen-go-wasmjs-ts)
handles client generation.

# Installation

	go install github.com/panyam/protoc-gen-go-wasmjs/cmd/protoc-gen-go-wasmjs-go@latest

Verify installation:

	which protoc-gen-go-wasmjs-go
	# Should output: /path/to/go/bin/protoc-gen-go-wasmjs-go

# Usage with buf

Add to your buf.gen.yaml:

	version: v2
	plugins:
	  # Generate standard Go protobuf types
	  - remote: buf.build/protocolbuffers/go
	    out: ./gen/go
	    opt: paths=source_relative

	  # Generate gRPC service interfaces
	  - remote: buf.build/grpc/go
	    out: ./gen/go
	    opt: paths=source_relative

	  # Generate Go WASM wrappers
	  - local: protoc-gen-go-wasmjs-go
	    out: ./gen/wasm/go
	    opt:
	      - js_structure=namespaced
	      - js_namespace=myApp
	      - module_name=my_services
	      - generate_build_script=true

Generate code:

	buf generate

# Generated Files

The generator produces the following files per proto package:

  - {module_name}.wasm.go: WASM wrapper with exports struct and RegisterAPI
  - main.go.example: Example usage showing dependency injection pattern
  - build.sh: Build script for WASM compilation (if generate_build_script=true)

Example generated structure:

	gen/wasm/go/
	├── library/v1/
	│   ├── library_v1.wasm.go       # WASM wrapper
	│   ├── main.go.example          # Usage example
	│   └── build.sh                 # Build script
	└── user/v1/
	    ├── user_v1.wasm.go
	    ├── main.go.example
	    └── build.sh

# Configuration Options

Core Generation:

  - wasm_export_path: Path where WASM wrapper should be generated (default: ".")
  - module_name: WASM module name (default: package_services)

JavaScript API Structure:

  - js_structure: API structure - namespaced|flat|service_based (default: "namespaced")
  - js_namespace: Global JavaScript namespace (default: lowercase package name)

Service & Method Selection:

  - services: Comma-separated list of services to generate (default: all)
  - method_include: Comma-separated glob patterns for methods to include
  - method_exclude: Comma-separated glob patterns for methods to exclude
  - method_rename: Comma-separated method renames (e.g., "OldName:NewName")

Build Integration:

  - wasm_package_suffix: Package suffix for WASM wrapper (default: "wasm")
  - generate_build_script: Generate build.sh script (default: true)

# Usage Example

Define your service:

	syntax = "proto3";
	package library.v1;

	service LibraryService {
	  rpc FindBooks(FindBooksRequest) returns (FindBooksResponse);
	  rpc CreateBook(CreateBookRequest) returns (CreateBookResponse);
	}

Generate WASM wrapper:

	buf generate

Implement your service with dependency injection:

	package main

	import (
	    "your-project/gen/wasm/library_v1"
	    libraryv1 "your-project/gen/go/library/v1"
	)

	type LibraryServiceImpl struct {
	    db    *sql.DB
	    cache *redis.Client
	}

	func (s *LibraryServiceImpl) FindBooks(ctx context.Context, req *libraryv1.FindBooksRequest) (*libraryv1.FindBooksResponse, error) {
	    // Your implementation
	    return &libraryv1.FindBooksResponse{Books: books}, nil
	}

	func main() {
	    // Initialize with dependency injection
	    exports := &library_v1.Library_v1_servicesServicesExports{
	        LibraryService: &LibraryServiceImpl{
	            db:    database,
	            cache: redisClient,
	        },
	    }

	    // Register JavaScript API
	    exports.RegisterAPI()

	    // Keep WASM running
	    select {}
	}

Build the WASM binary:

	cd gen/wasm/library/v1
	GOOS=js GOARCH=wasm go build -o library.wasm

	# Or use the generated build script
	./build.sh

# JavaScript API Structures

The generator supports three different API structures:

Namespaced (Recommended):

	# buf.gen.yaml
	opt:
	  - js_structure=namespaced
	  - js_namespace=myApp

	# JavaScript:
	window.myApp.libraryService.findBooks(request)

Flat:

	# buf.gen.yaml
	opt:
	  - js_structure=flat
	  - js_namespace=MyApp

	# JavaScript:
	window.MyAppLibraryServiceFindBooks(request)

Service-Based:

	# buf.gen.yaml
	opt:
	  - js_structure=service_based

	# JavaScript:
	window.services.library.findBooks(request)

# Architecture

The generator uses a layered architecture:

  1. GoGenerator (pkg/generators/go_generator.go): Top-level orchestrator
  2. BaseGenerator: Artifact collection and classification
  3. GoDataBuilder: Template data construction
  4. GoRenderer: Template rendering
  5. Filters: Service/method filtering

This separation enables:

  - Clear separation of concerns
  - Comprehensive unit testing
  - Easy customization through configuration

# Browser Service Support

For services that need to call browser APIs:

	service BrowserAPI {
	    option (wasmjs.v1.browser_provided) = true;

	    rpc GetLocalStorage(StorageKeyRequest) returns (StorageValueResponse);
	}

The generator automatically:

  - Detects browser-provided services
  - Generates client stubs that call JavaScript
  - Handles async communication to prevent deadlocks

# Advanced Features

Method Filtering:

	# Include only specific methods
	opt:
	  - method_include=Find*,Get*,Create*

	# Exclude internal methods
	opt:
	  - method_exclude=*Internal,*Debug

Service Selection:

	# Generate only specific services
	opt:
	  - services=LibraryService,UserService

Method Renaming:

	# Rename methods in generated code
	opt:
	  - method_rename=FindBooks:searchBooks,GetUser:fetchUser

# Error Handling

The generator validates configuration and provides detailed error messages:

  - Invalid configuration: Reports specific issues
  - Missing dependencies: Checks for required proto definitions
  - Template errors: Shows template execution failures

# Links

Related Tools:

  - protoc-gen-go-wasmjs-ts: TypeScript client generator
  - protoc: Protocol Buffers compiler
  - buf: Modern protobuf tooling

Documentation:

  - GitHub: https://github.com/panyam/protoc-gen-go-wasmjs
  - Examples: https://github.com/panyam/protoc-gen-go-wasmjs/tree/main/examples
*/
package main

import (
	"flag"
	"fmt"
	"log"

	"google.golang.org/protobuf/compiler/protogen"

	"github.com/panyam/protoc-gen-go-wasmjs/pkg/builders"
	"github.com/panyam/protoc-gen-go-wasmjs/pkg/filters"
	"github.com/panyam/protoc-gen-go-wasmjs/pkg/generators"
)

func main() {
	var flagSet flag.FlagSet

	// Core generation options
	wasmExportPath := flagSet.String("wasm_export_path", ".", "Path where WASM wrapper should be generated")

	// Service & method selection
	services := flagSet.String("services", "", "Comma-separated list of services to generate (default: all)")
	methodInclude := flagSet.String("method_include", "", "Comma-separated glob patterns for methods to include")
	methodExclude := flagSet.String("method_exclude", "", "Comma-separated glob patterns for methods to exclude")
	methodRename := flagSet.String("method_rename", "", "Comma-separated method renames (e.g., OldName:NewName)")

	// JavaScript API structure
	jsStructure := flagSet.String("js_structure", "namespaced", "JavaScript API structure (namespaced|flat|service_based)")
	jsNamespace := flagSet.String("js_namespace", "", "Global JavaScript namespace (default: lowercase package name)")
	moduleName := flagSet.String("module_name", "", "WASM module name (default: package_services)")

	// Build integration
	wasmPackageSuffix := flagSet.String("wasm_package_suffix", "wasm", "Package suffix for WASM wrapper")
	generateBuildScript := flagSet.Bool("generate_build_script", true, "Generate build script for WASM compilation")

	protogen.Options{
		ParamFunc: flagSet.Set,
	}.Run(func(gen *protogen.Plugin) error {
		log.Printf("NEW GENERATOR: Plugin callback started")
		log.Printf("NEW GENERATOR: Request has %d files", len(gen.Files))
		log.Printf("NEW GENERATOR: Request parameters: %+v", gen.Request.GetParameter())
		
		defer func() {
			log.Printf("NEW GENERATOR: Plugin callback ending")
			log.Printf("NEW GENERATOR: Response has %d files", len(gen.Response().File))
			for i, file := range gen.Response().File {
				log.Printf("NEW GENERATOR: Response file %d: %s (%d bytes)", 
					i, file.GetName(), len(file.GetContent()))
			}
		}()
		// Create generation configuration
		config := &builders.GenerationConfig{
			WasmExportPath:      *wasmExportPath,
			JSStructure:         *jsStructure,
			JSNamespace:         *jsNamespace,
			ModuleName:          *moduleName,
			WasmPackageSuffix:   *wasmPackageSuffix,
			GenerateBuildScript: *generateBuildScript,
		}

		// Create filter criteria from configuration
		filterCriteria, err := filters.ParseFromConfig(*services, *methodInclude, *methodExclude, *methodRename)
		if err != nil {
			return fmt.Errorf("invalid filter configuration: %w", err)
		}

		// Create Go generator
		goGenerator := generators.NewGoGenerator(gen)

		// Validate configuration
		if err := goGenerator.ValidateConfig(config); err != nil {
			return fmt.Errorf("invalid configuration: %w", err)
		}

		// Perform generation with detailed error handling
		if err := goGenerator.Generate(config, filterCriteria); err != nil {
			log.Printf("protoc-gen-go-wasmjs-go: Generation failed: %v", err)
			return fmt.Errorf("Go generation failed: %w", err)
		}

		log.Printf("protoc-gen-go-wasmjs-go: Generation completed successfully")
		return nil
	})
}
