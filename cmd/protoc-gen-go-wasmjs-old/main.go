// Copyright 2025 Sri Panyam
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main

import (
	"flag"
	"log"

	"github.com/panyam/protoc-gen-go-wasmjs/pkg/generator"
	"google.golang.org/protobuf/compiler/protogen"
)

func main() {
	var flagSet flag.FlagSet

	// Core integration options
	generateWasm := flagSet.Bool("generate_wasm", true, "Generate WASM wrapper (default: true)")
	generateTypeScript := flagSet.Bool("generate_typescript", true, "Generate TypeScript client (default: true)")
	wasmExportPath := flagSet.String("wasm_export_path", ".", "Path where WASM wrapper should be generated")
	tsExportPath := flagSet.String("ts_export_path", ".", "Path where TypeScript files should be generated")

	// Service & method selection
	services := flagSet.String("services", "", "Comma-separated list of services to generate (default: all)")
	methodInclude := flagSet.String("method_include", "", "Comma-separated glob patterns for methods to include")
	methodExclude := flagSet.String("method_exclude", "", "Comma-separated glob patterns for methods to exclude")
	methodRename := flagSet.String("method_rename", "", "Comma-separated method renames (e.g., OldName:NewName)")

	// JS API structure
	jsStructure := flagSet.String("js_structure", "namespaced", "JavaScript API structure (namespaced|flat|service_based)")
	jsNamespace := flagSet.String("js_namespace", "", "Global JavaScript namespace (default: lowercase package name)")
	moduleName := flagSet.String("module_name", "", "WASM module name (default: package_services)")

	// Customization
	templateDir := flagSet.String("template_dir", "", "Directory containing custom templates")
	wasmTemplate := flagSet.String("wasm_template", "", "Custom WASM template file")
	tsTemplate := flagSet.String("ts_template", "", "Custom TypeScript template file")

	// Build integration
	wasmPackageSuffix := flagSet.String("wasm_package_suffix", "wasm", "Package suffix for WASM wrapper")
	generateBuildScript := flagSet.Bool("generate_build_script", true, "Generate build script for WASM compilation")

	protogen.Options{
		ParamFunc: flagSet.Set,
	}.Run(func(gen *protogen.Plugin) error {
		log.Printf("OLD GENERATOR: Plugin callback started")
		log.Printf("OLD GENERATOR: Request has %d files", len(gen.Files))
		log.Printf("OLD GENERATOR: Request parameters: %+v", gen.Request.GetParameter())
		
		defer func() {
			log.Printf("OLD GENERATOR: Plugin callback ending")
			log.Printf("OLD GENERATOR: Response has %d files", len(gen.Response().File))
			for i, file := range gen.Response().File {
				log.Printf("OLD GENERATOR: Response file %d: %s (%d bytes)", 
					i, file.GetName(), len(file.GetContent()))
			}
		}()
		// Create configuration from parsed flags
		config := &generator.Config{
			// Core integration
			GenerateWasm:       *generateWasm,
			GenerateTypeScript: *generateTypeScript,
			WasmExportPath:     *wasmExportPath,
			TSExportPath:       *tsExportPath,

			// Service & method selection
			Services:      *services,
			MethodInclude: *methodInclude,
			MethodExclude: *methodExclude,
			MethodRename:  *methodRename,

			// JS API structure
			JSStructure: *jsStructure,
			JSNamespace: *jsNamespace,
			ModuleName:  *moduleName,

			// Customization
			TemplateDir:  *templateDir,
			WasmTemplate: *wasmTemplate,
			TSTemplate:   *tsTemplate,

			// Build integration
			WasmPackageSuffix:   *wasmPackageSuffix,
			GenerateBuildScript: *generateBuildScript,
		}

		// Group files by package to generate one WASM module per package
		packageFiles := make(map[string][]*protogen.File)

		// Collect all browser-provided services across all packages
		allBrowserServices := make([]*generator.BrowserServiceInfo, 0)

		// Group files by package and collect browser services
		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}
			packageName := string(f.Desc.Package())

			// Skip wasmjs annotation packages - they are library files, not user code
			if packageName == "wasmjs.v1" {
				continue
			}

			packageFiles[packageName] = append(packageFiles[packageName], f)

			// Collect browser-provided services from this file
			for _, service := range f.Services {
				if generator.IsBrowserProvidedService(service) {
					allBrowserServices = append(allBrowserServices, &generator.BrowserServiceInfo{
						File:    f,
						Service: service,
					})
				}
			}
		}

		// Generate one WASM module per package
		for _, files := range packageFiles {
			if len(files) == 0 {
				continue
			}

			// Use the first file as the primary file, but collect services from all files
			primaryFile := files[0]
			fileGen := generator.NewFileGenerator(primaryFile, gen, config)

			// Set the additional files for this package
			fileGen.SetPackageFiles(files)

			// Set all browser services so they can be included in generation
			fileGen.SetAllBrowserServices(allBrowserServices)

			if err := fileGen.Generate(); err != nil {
				return err
			}
		}

		return nil
	})
}
