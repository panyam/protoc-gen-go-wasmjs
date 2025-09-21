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

// protoc-gen-go-wasmjs-ts generates TypeScript clients and types for gRPC services.
// This is a focused generator that only produces TypeScript artifacts,
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
