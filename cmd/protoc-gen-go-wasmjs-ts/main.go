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
protoc-gen-go-wasmjs-ts is a Protocol Buffers compiler plugin that generates TypeScript clients and types for gRPC services.

# Overview

This generator is part of the protoc-gen-go-wasmjs split architecture. It focuses exclusively
on generating TypeScript client code and type definitions, while the Go generator
(protoc-gen-go-wasmjs-go) handles WASM wrapper generation.

# Installation

	go install github.com/panyam/protoc-gen-go-wasmjs/cmd/protoc-gen-go-wasmjs-ts@latest

Verify installation:

	which protoc-gen-go-wasmjs-ts
	# Should output: /path/to/go/bin/protoc-gen-go-wasmjs-ts

# Runtime Package

Generated TypeScript code requires the runtime utilities package:

	npm install @protoc-gen-go-wasmjs/runtime
	# or
	pnpm add @protoc-gen-go-wasmjs/runtime

# Usage with buf

Add to your buf.gen.yaml:

	version: v2
	plugins:
	  # Generate TypeScript clients and types
	  - local: protoc-gen-go-wasmjs-ts
	    out: ./web/src/generated
	    opt:
	      - js_structure=namespaced
	      - js_namespace=myApp
	      - generate_clients=true
	      - generate_types=true
	      - generate_factories=true

Generate code:

	buf generate

# Generated Files

The generator produces a comprehensive TypeScript ecosystem per proto package:

Bundle Level (Module):

  - index.ts: Base bundle class extending WASMBundle

Package Level:

  - {package}/interfaces.ts: Pure TypeScript interfaces
  - {package}/models.ts: Concrete implementations with defaults
  - {package}/schemas.ts: Field metadata for runtime introspection
  - {package}/deserializer.ts: Schema-driven data population
  - {package}/factory.ts: Object factories (when generate_factories=true)

Service Level:

  - {package}/{serviceName}Client.ts: Per-service typed clients

Example generated structure:

	web/src/generated/
	├── index.ts                              # Base bundle class
	├── presenter/v1/
	│   ├── presenterServiceClient.ts         # Service client
	│   ├── interfaces.ts                     # Type definitions
	│   ├── models.ts                         # Concrete classes
	│   ├── schemas.ts                        # Field metadata
	│   ├── deserializer.ts                   # Data population
	│   └── factory.ts                        # Object factories
	└── browser/v1/
	    ├── browserAPIClient.ts
	    ├── interfaces.ts
	    ├── models.ts
	    ├── schemas.ts
	    ├── deserializer.ts
	    └── factory.ts

# Configuration Options

Core Generation:

  - ts_export_path: Path where TypeScript files should be generated (default: ".")
  - module_name: Module name for generated code (default: package_services)

JavaScript API Structure:

  - js_structure: API structure - namespaced|flat|service_based (default: "namespaced")
  - js_namespace: Global JavaScript namespace (default: lowercase package name)

Content Filtering:

  - generate_clients: Generate TypeScript client classes (default: true)
  - generate_types: Generate TypeScript interfaces/models (default: true)
  - generate_factories: Generate TypeScript factory classes (default: true)

Service & Method Selection:

  - services: Comma-separated list of services to generate clients for (default: all)
  - method_include: Comma-separated glob patterns for methods to include
  - method_exclude: Comma-separated glob patterns for methods to exclude
  - method_rename: Comma-separated method renames (e.g., "OldName:NewName")

# Usage Example

Define your service:

	syntax = "proto3";
	package library.v1;

	service LibraryService {
	  rpc FindBooks(FindBooksRequest) returns (FindBooksResponse);
	  rpc CreateBook(CreateBookRequest) returns (CreateBookResponse);
	}

	message FindBooksRequest {
	  string query = 1;
	  int32 max_results = 2;
	}

	message FindBooksResponse {
	  repeated Book books = 1;
	}

	message Book {
	  string id = 1;
	  string title = 2;
	  string author = 3;
	}

Generate TypeScript clients:

	buf generate

Use in your application:

	import { Library_servicesBundle } from './generated';
	import type { FindBooksRequest } from './generated/library/v1/interfaces';

	// Create and load WASM bundle
	const bundle = new Library_servicesBundle();
	await bundle.loadWasm('./library.wasm');

	// Fully typed method call with IntelliSense
	const request: FindBooksRequest = {
	  query: "golang",
	  maxResults: 10
	};
	const response = await bundle.libraryService.findBooks(request);

	// TypeScript knows response.books is Book[]
	console.log(`Found ${response.books.length} books`);
	response.books.forEach(book => {
	  console.log(`${book.title} by ${book.author}`);
	});

# TypeScript Generation Model

The generator creates a clean separation between interfaces and implementations:

Interfaces (Pure Type Definitions):

	// interfaces.ts
	export interface Book {
	  id: string;
	  title: string;
	  author: string;
	}

Models (Concrete Implementations):

	// models.ts
	export class Book implements BookInterface {
	  id: string = "";
	  title: string = "";
	  author: string = "";
	}

Factories (Object Construction):

	// factory.ts
	export class LibraryV1Factory {
	  newBook(parent?: any, attributeName?: string, data?: any): FactoryResult<Book> {
	    const instance = new ConcreteBook();
	    return { instance, fullyLoaded: false };
	  }
	}

Schemas (Runtime Metadata):

	// schemas.ts
	export const BookSchema: MessageSchema = {
	  name: "Book",
	  fields: [
	    { name: "id", type: FieldType.STRING, id: 1 },
	    { name: "title", type: FieldType.STRING, id: 2 },
	    { name: "author", type: FieldType.STRING, id: 3 },
	  ]
	};

Deserializers (Data Population):

	// deserializer.ts
	export class LibraryV1Deserializer extends BaseDeserializer {
	  static from<T>(messageType: string, data: any): T {
	    const deserializer = new LibraryV1Deserializer();
	    return deserializer.createAndDeserialize<T>(messageType, data);
	  }
	}

# Service Client Generation

Each service gets its own typed client file:

	// presenter/v1/presenterServiceClient.ts
	import { ServiceClient } from '@protoc-gen-go-wasmjs/runtime';
	import type { LoadUserDataRequest, LoadUserDataResponse } from './interfaces';

	export class PresenterServiceServiceClient extends ServiceClient {
	  async loadUserData(request: LoadUserDataRequest): Promise<LoadUserDataResponse> {
	    return this.callMethod('presenterService.loadUserData', request);
	  }

	  // Streaming method with typed callback
	  updateUIState(
	    request: StateUpdateRequest,
	    callback: (response: UIUpdate | null, error: string | null, done: boolean) => boolean
	  ): void {
	    return this.callStreamingMethod('presenterService.updateUIState', request, callback);
	  }
	}

# Bundle-Based Architecture

The generator creates a simple base bundle class that users compose with service clients:

	// index.ts - Generated base bundle
	export class Library_servicesBundle extends WASMBundle {
	  constructor() {
	    super({
	      moduleName: 'library_services',
	      apiStructure: 'namespaced',
	      jsNamespace: 'myApp'
	    });
	  }
	}

	// User code - Composition pattern
	const bundle = new Library_servicesBundle();
	const libraryService = new LibraryServiceServiceClient(bundle);
	await bundle.loadWasm('./library.wasm');

# Browser Service Integration

For services that call browser APIs:

	service BrowserAPI {
	    option (wasmjs.v1.browser_provided) = true;

	    rpc GetLocalStorage(StorageKeyRequest) returns (StorageValueResponse);
	    rpc Fetch(FetchRequest) returns (FetchResponse);
	}

TypeScript implementation:

	const bundle = new My_servicesBundle();

	// Register browser service implementation
	bundle.registerBrowserService('BrowserAPI', {
	  async getLocalStorage(request) {
	    return {
	      value: localStorage.getItem(request.key) || '',
	      exists: true
	    };
	  },
	  async fetch(request) {
	    const response = await fetch(request.url);
	    return { body: await response.text() };
	  }
	});

	await bundle.loadWasm('./my_services.wasm');

# Architecture

The generator uses a layered architecture:

  1. TSGenerator (pkg/generators/ts_generator.go): Top-level orchestrator
  2. BaseGenerator: Artifact collection and classification
  3. TSDataBuilder: Template data construction
  4. TSRenderer: Template rendering
  5. Filters: Service/method filtering

This separation enables:

  - Clear separation of concerns
  - Comprehensive unit testing
  - Package-based generation (follows proto structure)
  - Cross-package import resolution

# Advanced Features

Selective Generation:

	# Only generate clients, skip types
	opt:
	  - generate_clients=true
	  - generate_types=false

	# Only generate types, skip clients
	opt:
	  - generate_clients=false
	  - generate_types=true

	# Skip factories to reduce bundle size
	opt:
	  - generate_factories=false

Method Filtering:

	# Include only specific methods
	opt:
	  - method_include=Find*,Get*,Create*

	# Exclude internal methods
	opt:
	  - method_exclude=*Internal,*Debug

Service Selection:

	# Generate clients only for specific services
	opt:
	  - services=LibraryService,UserService

# Full Type Safety

All generated code is fully typed for TypeScript:

  - Method parameters use generated interfaces
  - Return types are properly inferred
  - Streaming callbacks have typed signatures
  - IntelliSense works throughout

Example with full type checking:

	const request: FindBooksRequest = {
	  query: "golang",
	  maxResults: 10,
	  // TypeScript error if you add invalid fields
	};

	const response = await service.findBooks(request);
	// TypeScript knows response is FindBooksResponse
	// response.books is Book[]
	// Each book has id, title, author properties

# Error Handling

The generator validates configuration and provides detailed error messages:

  - Invalid configuration: Reports specific issues
  - Missing dependencies: Checks for required proto definitions
  - Template errors: Shows template execution failures
  - Import resolution: Validates cross-package references

# Links

Related Tools:

  - protoc-gen-go-wasmjs-go: Go WASM wrapper generator
  - @protoc-gen-go-wasmjs/runtime: NPM runtime package
  - protoc: Protocol Buffers compiler
  - buf: Modern protobuf tooling

Documentation:

  - GitHub: https://github.com/panyam/protoc-gen-go-wasmjs
  - Examples: https://github.com/panyam/protoc-gen-go-wasmjs/tree/main/examples
  - Runtime Package: https://www.npmjs.com/package/@protoc-gen-go-wasmjs/runtime
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
	tsExportPath := flagSet.String("ts_export_path", ".", "Path where TypeScript files should be generated")

	// Service & method selection (for client generation)
	services := flagSet.String("services", "", "Comma-separated list of services to generate clients for (default: all)")
	methodInclude := flagSet.String("method_include", "", "Comma-separated glob patterns for methods to include")
	methodExclude := flagSet.String("method_exclude", "", "Comma-separated glob patterns for methods to exclude")
	methodRename := flagSet.String("method_rename", "", "Comma-separated method renames (e.g., OldName:NewName)")

	// JavaScript API structure
	jsStructure := flagSet.String("js_structure", "namespaced", "JavaScript API structure (namespaced|flat|service_based)")
	jsNamespace := flagSet.String("js_namespace", "", "Global JavaScript namespace (default: lowercase package name)")
	
	// TypeScript-specific options
	moduleName := flagSet.String("module_name", "", "TypeScript module name (default: package_services)")

	// Content filtering
	generateClients := flagSet.Bool("generate_clients", true, "Generate TypeScript client classes for services")
	generateTypes := flagSet.Bool("generate_types", true, "Generate TypeScript interfaces and models for messages/enums")
	generateFactories := flagSet.Bool("generate_factories", true, "Generate TypeScript factory classes for creating message objects")

	protogen.Options{
		ParamFunc: flagSet.Set,
	}.Run(func(gen *protogen.Plugin) error {
		log.Printf("TS GENERATOR: Received parameters: js_structure=%s, js_namespace=%s", *jsStructure, *jsNamespace)
		
		// Create generation configuration
		config := &builders.GenerationConfig{
			TSExportPath:      *tsExportPath,
			JSStructure:       *jsStructure,
			JSNamespace:       *jsNamespace,
			ModuleName:        *moduleName,
			GenerateClients:   *generateClients,
			GenerateTypes:     *generateTypes,
			GenerateFactories: *generateFactories,
		}

		// Create filter criteria from configuration
		filterCriteria, err := filters.ParseFromConfig(*services, *methodInclude, *methodExclude, *methodRename)
		if err != nil {
			return fmt.Errorf("invalid filter configuration: %w", err)
		}

		// Adjust filter criteria based on generation options
		if !config.GenerateClients {
			// If not generating clients, exclude all services
			filterCriteria.ServicesSet = make(map[string]bool) // Empty = exclude all when HasServiceFilter() is false
			// We need a way to indicate "exclude all services"
			filterCriteria.ServicesSet["__EXCLUDE_ALL__"] = false
		}

		if !config.GenerateTypes {
			// If not generating types, exclude nested content
			filterCriteria.ExcludeNestedMessages = true
			filterCriteria.ExcludeNestedEnums = true
			// TODO: Add option to exclude all messages/enums
		}

		// Create TypeScript generator
		tsGenerator := generators.NewTSGenerator(gen)

		// Validate configuration
		if err := tsGenerator.ValidateConfig(config); err != nil {
			return fmt.Errorf("invalid configuration: %w", err)
		}

		// Perform generation
		if err := tsGenerator.Generate(config, filterCriteria); err != nil {
			return fmt.Errorf("TypeScript generation failed: %w", err)
		}

		log.Printf("protoc-gen-go-wasmjs-ts: Generation completed successfully")
		return nil
	})
}
