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

// protoc-gen-go-wasmjs-go generates Go WASM wrappers for gRPC services.
// This is a focused generator that only produces Go WASM artifacts,
// using the new layered architecture for better maintainability and testing.
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
