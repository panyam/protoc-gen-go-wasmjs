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
Package protoc-gen-go-wasmjs generates WASM bindings and TypeScript clients for Protocol Buffer gRPC services.

# Overview

protoc-gen-go-wasmjs is a Protocol Buffers compiler plugin that enables local-first applications
by generating WebAssembly (WASM) wrappers and TypeScript clients from your gRPC service definitions.
This allows you to run the same service logic in both server and browser environments with full
type safety and seamless integration.

# Key Features

  - Dual-target code generation: WASM wrappers and TypeScript clients
  - Full TypeScript type safety with automatic interface generation
  - Local-first architecture: same service logic runs server-side or in the browser
  - Browser service communication: WASM can call browser APIs (localStorage, fetch, etc.)
  - Dependency injection support through export pattern
  - Flexible API structures: namespaced, flat, or service-based
  - Comprehensive streaming support (server, client, bidirectional)
  - Runtime package with inheritance-based architecture

# Installation

Install the plugin generators:

	# Install Go WASM generator
	go install github.com/panyam/protoc-gen-go-wasmjs/cmd/protoc-gen-go-wasmjs-go@latest

	# Install TypeScript generator
	go install github.com/panyam/protoc-gen-go-wasmjs/cmd/protoc-gen-go-wasmjs-ts@latest

For TypeScript runtime support:

	npm install @protoc-gen-go-wasmjs/runtime

# Quick Start

Configure your buf.gen.yaml:

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

	  # Generate TypeScript clients
	  - local: protoc-gen-go-wasmjs-ts
	    out: ./web/src/generated
	    opt:
	      - js_structure=namespaced
	      - js_namespace=myApp

Run generation:

	buf generate

# Architecture

The plugin follows a layered architecture with distinct responsibilities:

		┌─────────────────┐
		│   .proto files  │
		└────────┬────────┘
		         │
		         ▼
		┌─────────────────┐     ┌──────────────────┐
		│  BaseGenerator  │────▶│   TSGenerator    │──┐
		│  (Artifact      │     │   (TypeScript)   │  │
		│   Collection)   │     └──────────────────┘  │
		└────────┬────────┘                           │
		         │                                    │
		         │         ┌──────────────────┐       │
		         └────────▶│   GoGenerator    │───────┤
		                   │   (WASM)         │       │
		                   └──────────────────┘       │
		                                              │
		                                              ▼
	 	                                    ┌──────────────────┐
	 	                                    │ Generated Files  │
	 	                                    │ • WASM wrappers  │
		                                    │ • TS clients     │
		                                    │ • TS interfaces  │
		                                    │ • Factories      │
		                                    └──────────────────┘

# 4-Step Artifact Processing

The BaseGenerator implements a clean separation between artifact collection and file generation:

 1. COLLECT ALL ARTIFACTS → BaseGenerator.CollectAllArtifacts()
    ├─ Get map of all artifacts from protogen
    └─ Available regardless of protoc's Generate flags

 2. CLASSIFY ARTIFACTS → ArtifactCatalog
    ├─ Services (regular + browser)
    ├─ Messages by package
    └─ Enums by package

 3. MAP ARTIFACTS TO FILES → planFilesFromCatalog()
    ├─ Generator-specific slice/dice/group logic
    ├─ N artifacts → 1 file (bundle with multiple services)
    └─ 1 artifact → 1 file (per-service clients)

 4. CREATE PROTOGEN FILES → fileSet.CreateFiles(plugin)
    ├─ Send final mapping to protogen
    └─ Only after all artifact mapping decisions are complete

# Usage Example

Define your service in protobuf:

	syntax = "proto3";
	package library.v1;

	import "wasmjs/v1/annotations.proto";

	service LibraryService {
	  rpc FindBooks(FindBooksRequest) returns (FindBooksResponse);

	  // Async method with callback support
	  rpc LoadData(LoadDataRequest) returns (LoadDataResponse) {
	    option (wasmjs.v1.async_method) = { is_async: true };
	  };
	}

	message FindBooksRequest {
	  string query = 1;
	}

	message FindBooksResponse {
	  repeated Book books = 1;
	}

	message Book {
	  string id = 1;
	  string title = 2;
	  string author = 3;
	}

Implement the service in Go with dependency injection:

	package main

	import (
	    "context"
	    "your-project/gen/wasm/library_services"
	    libraryv1 "your-project/gen/go/library/v1"
	)

	type LibraryServiceImpl struct {
	    db    *sql.DB
	    cache *redis.Client
	}

	func (s *LibraryServiceImpl) FindBooks(
	    ctx context.Context,
	    req *libraryv1.FindBooksRequest,
	) (*libraryv1.FindBooksResponse, error) {
	    // Your implementation
	    return &libraryv1.FindBooksResponse{Books: books}, nil
	}

	func main() {
	    // Initialize with dependency injection
	    exports := &library_services.Library_servicesServicesExports{
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

	GOOS=js GOARCH=wasm go build -o library.wasm

Use in TypeScript with full type safety:

	import { Library_servicesBundle } from './generated';
	import type { FindBooksRequest } from './generated/library/v1/interfaces';

	// Create and load WASM bundle
	const bundle = new Library_servicesBundle();
	await bundle.loadWasm('./library.wasm');

	// Fully typed method call
	const request: FindBooksRequest = { query: "golang" };
	const response = await bundle.libraryService.findBooks(request);

	console.log(`Found ${response.books.length} books`);

# Browser Service Integration

For services that need to call browser APIs:

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

# Configuration Options

Core Integration:

  - ts_export_path: Path where TypeScript files should be generated
  - wasm_export_path: Path where WASM wrapper should be generated
  - module_name: WASM module name (default: package_services)

Service & Method Selection:

  - services: Specific services to generate (comma-separated)
  - method_include: Include methods by glob pattern (e.g., "Find*,Get*")
  - method_exclude: Exclude methods by glob pattern (e.g., "*Internal")
  - method_rename: Rename methods (e.g., "FindBooks:searchBooks")

JavaScript API Structure:

  - js_structure: API structure (namespaced|flat|service_based)
  - js_namespace: Global namespace name
  - generate_clients: Generate TypeScript client classes (default: true)
  - generate_types: Generate TypeScript interfaces/models (default: true)
  - generate_factories: Generate TypeScript factory classes (default: true)

Build Integration:

  - generate_build_script: Generate build.sh script (default: true)
  - wasm_package_suffix: Package suffix for WASM wrapper (default: "wasm")

# Package Organization

The plugin is organized into focused packages:

  - cmd/protoc-gen-go-wasmjs-go: Go WASM generator entry point
  - cmd/protoc-gen-go-wasmjs-ts: TypeScript generator entry point
  - pkg/generators: Core generation logic (BaseGenerator, TSGenerator, GoGenerator)
  - pkg/builders: Template data building and file planning
  - pkg/renderers: Template rendering with proper imports
  - pkg/filters: Service/method filtering and artifact collection
  - pkg/core: Pure utility functions (name conversion, path calculation, proto analysis)
  - pkg/wasm: WASM runtime utilities (browser channel, helpers)
  - proto/wasmjs/v1: WASM annotation definitions

# TypeScript Generation

The generator creates a complete TypeScript ecosystem per proto package:

  - interfaces.ts: Pure TypeScript interfaces for type safety
  - models.ts: Concrete implementations with default values
  - factory.ts: Object factories with context awareness (when generate_factories=true)
  - schemas.ts: Field metadata for runtime introspection
  - deserializer.ts: Schema-driven data population
  - {serviceName}Client.ts: Per-service typed clients

# Local-First Use Case

The primary use case enables applications where the same service logic runs in different environments:

Server Environment (Full Dataset):

	type LibraryService struct {
	    db *sql.DB // Access to millions of books
	}

	func (s *LibraryService) FindBooks(ctx context.Context, req *FindBooksRequest) (*FindBooksResponse, error) {
	    // Query full database
	    return s.searchDatabase(req.Query)
	}

Browser Environment (Local Subset):

	type LibraryService struct {
	    books []Book // Local subset from localStorage
	}

	func (s *LibraryService) FindBooks(ctx context.Context, req *FindBooksRequest) (*FindBooksResponse, error) {
	    // Search local books only
	    return s.searchLocalBooks(req.Query)
	}

Both use the same gRPC interface definition, maintaining API consistency across deployments.

# Links

  - GitHub: https://github.com/panyam/protoc-gen-go-wasmjs
  - Examples: https://github.com/panyam/protoc-gen-go-wasmjs/tree/main/examples
  - Runtime Package: https://www.npmjs.com/package/@protoc-gen-go-wasmjs/runtime
*/
package protoc_gen_go_wasmjs
