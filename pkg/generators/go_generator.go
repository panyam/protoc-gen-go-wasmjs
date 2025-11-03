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

package generators

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"

	"github.com/panyam/protoc-gen-go-wasmjs/pkg/builders"
	"github.com/panyam/protoc-gen-go-wasmjs/pkg/filters"
	"github.com/panyam/protoc-gen-go-wasmjs/pkg/renderers"
)

// GoGenerator orchestrates the complete Go WASM generation process.
// This is the top-level generator that coordinates all layers to produce Go WASM artifacts.
type GoGenerator struct {
	// Embed base generator for artifact collection
	*BaseGenerator

	// Builder and renderer specific to Go WASM
	dataBuilder *builders.GoDataBuilder
	renderer    *renderers.GoRenderer
}

// NewGoGenerator creates a new Go generator with all necessary dependencies.
// This sets up the complete processing pipeline for Go WASM generation.
func NewGoGenerator(plugin *protogen.Plugin) *GoGenerator {
	// Create base generator with artifact collection capabilities
	baseGenerator := NewBaseGenerator(plugin)

	// Create Go-specific builder and renderer  
	dataBuilder := builders.NewGoDataBuilder(
		baseGenerator.analyzer,
		baseGenerator.pathCalc,
		baseGenerator.nameConv,
		baseGenerator.serviceFilter,
		baseGenerator.methodFilter,
		baseGenerator.msgCollector,
		baseGenerator.enumCollector,
	)
	renderer := renderers.NewGoRenderer()

	return &GoGenerator{
		BaseGenerator: baseGenerator,
		dataBuilder:   dataBuilder,
		renderer:      renderer,
	}
}

// Generate performs the complete Go WASM generation process.
// This is the main entry point that coordinates all layers to produce Go artifacts.
func (gg *GoGenerator) Generate(config *builders.GenerationConfig, filterCriteria *filters.FilterCriteria) error {
	// Phase 1: Filter packages
	log.Printf("Starting package filtering from %d input files", len(gg.plugin.Files))
	for i, file := range gg.plugin.Files {
		log.Printf("File %d: %s, Generate=%v, Package=%s, Services=%d", 
			i, file.Desc.Path(), file.Generate, file.Desc.Package(), len(file.Services))
	}
	packageFiles, stats := gg.packageFilter.FilterPackages(gg.plugin.Files, filterCriteria)
	
	log.Printf("Package filtering result: %d packages, stats: %s", len(packageFiles), stats.Summary())
	if len(packageFiles) == 0 {
		return nil // No packages to process
	}

	// Phase 2: Collect all browser services (from all packages)
	var allBrowserServices []*protogen.Service
	for _, files := range packageFiles {
		browserServices := gg.serviceFilter.GetBrowserProvidedServices(files, filterCriteria)
		allBrowserServices = append(allBrowserServices, browserServices...)
	}

	// Phase 3: Generate WASM wrapper for each package
	for packageName, files := range packageFiles {
		packageInfo := &builders.PackageInfo{
			Name:      packageName,
			Path:      gg.pathCalc.BuildPackagePath(packageName),
			GoPackage: string(files[0].GoImportPath),
			Files:     files,
		}

		// Build template data
		templateData, err := gg.dataBuilder.BuildTemplateData(packageInfo, allBrowserServices, filterCriteria, config)
		if err != nil {
			return fmt.Errorf("failed to build template data for package %s: %w", packageName, err)
		}

		// Skip packages with no services
		if templateData == nil {
			log.Printf("Skipping package %s: no services to generate", packageName)
			continue
		}
		
		log.Printf("Generated template data for package %s: %d services, %d browser clients", 
			packageName, len(templateData.Services), len(templateData.BrowserClients))

		// Validate template data
		if err := gg.renderer.ValidateGoTemplateData(templateData); err != nil {
			return fmt.Errorf("invalid template data for package %s: %w", packageName, err)
		}

		// Generate files for this package
		if err := gg.generatePackageFiles(templateData, config); err != nil {
			return fmt.Errorf("failed to generate files for package %s: %w", packageName, err)
		}
	}

	// Log generation summary (not to stdout to avoid corrupting protobuf response)
	log.Printf("Go generator processed %s", stats.Summary())

	return nil
}

// generatePackageFiles handles complete file generation for a package using file planning.
// This is the new approach where the generator controls all file creation and naming.
func (gg *GoGenerator) generatePackageFiles(data *builders.GoTemplateData, config *builders.GenerationConfig) error {
	// Phase 1: Plan what files we need to generate
	filePlan := gg.planGoFiles(data, config)

	// Use old generator pattern: create files on-demand when rendering
	if err := gg.renderFilesDirectly(filePlan, data, config); err != nil {
	 return fmt.Errorf("file rendering failed: %w", err)
	}

	return nil
}

// renderFilesDirectly renders files using the old generator pattern - create on demand.
// This matches exactly how the old generator works to avoid wire protocol issues.
func (gg *GoGenerator) renderFilesDirectly(filePlan *builders.FilePlan, data *builders.GoTemplateData, config *builders.GenerationConfig) error {
	for _, spec := range filePlan.Specs {
		log.Printf("Creating and rendering file on-demand: %s", spec.Filename)

		// Create GeneratedFile immediately before rendering (like old generator)
		generatedFile := gg.plugin.NewGeneratedFile(spec.Filename, "")
		log.Printf("Created GeneratedFile: %s -> %p", spec.Filename, generatedFile)

		// Render based on file type
		switch spec.Type {
		case "converters":
			log.Printf("CONVERTERS: Attempting to render converters...")
			if err := gg.renderer.RenderConvertersDirect(generatedFile, data); err != nil {
				log.Printf("CONVERTERS: ERROR rendering converters: %v", err)
				return fmt.Errorf("failed to render converters file %s: %w", spec.Filename, err)
			}
			log.Printf("CONVERTERS: Converters rendered successfully")

		case "exports":
			log.Printf("EXPORTS: Attempting to render exports...")
			if err := gg.renderer.RenderExportsDirect(generatedFile, data); err != nil {
				log.Printf("EXPORTS: ERROR rendering exports: %v", err)
				return fmt.Errorf("failed to render exports file %s: %w", spec.Filename, err)
			}
			log.Printf("EXPORTS: Exports rendered successfully")

		case "browser_clients":
			log.Printf("BROWSER_CLIENTS: Attempting to render browser clients...")
			if err := gg.renderer.RenderBrowserClientsDirect(generatedFile, data); err != nil {
				log.Printf("BROWSER_CLIENTS: ERROR rendering browser clients: %v", err)
				return fmt.Errorf("failed to render browser clients file %s: %w", spec.Filename, err)
			}
			log.Printf("BROWSER_CLIENTS: Browser clients rendered successfully")

		case "example":
			log.Printf("MAIN: Attempting to render main file...")
			if err := gg.renderer.RenderMainExampleDirect(generatedFile, data); err != nil {
				log.Printf("MAIN: ERROR rendering main file: %v", err)
				return fmt.Errorf("failed to render main file %s: %w", spec.Filename, err)
			}
			log.Printf("MAIN: Main file rendered successfully")

		case "script":
			log.Printf("BUILD: Attempting to render build script...")
			if err := gg.renderer.RenderBuildScriptDirect(generatedFile, data); err != nil {
				log.Printf("BUILD: ERROR rendering build script: %v", err)
				return fmt.Errorf("failed to render build script %s: %w", spec.Filename, err)
			}
			log.Printf("BUILD: Build script rendered successfully")

		default:
			log.Printf("Unknown file type %s for %s", spec.Type, spec.Filename)
		}
	}

	return nil
}

// planGoFiles creates a file plan for Go WASM generation based on data and configuration.
// This method decides what files to generate and where to put them.
func (gg *GoGenerator) planGoFiles(data *builders.GoTemplateData, config *builders.GenerationConfig) *builders.FilePlan {
	var specs []builders.FileSpec

	// Split WASM generation into 3 files for better modularity:
	// 1. Converters - syscall/js converters (createJSResponse) and stream wrappers
	// 2. Exports - Exports struct, RegisterAPI, method wrappers
	// 3. Browser clients - Browser service client implementations

	// Use GoPackage to determine output path to avoid collisions when multiple files
	// have the same proto package but different go_package options
	packagePath := gg.calculateOutputPath(data)
	baseName := gg.calculateBaseName(data)

	// Generate converters file (needed whenever we have services)
	// This contains createJSResponse() which is used by all service method wrappers
	if len(data.Services) > 0 {
		convertersFilename := filepath.Join(packagePath, baseName+"_converters.wasm.go")
		log.Printf("Planning converters file: %s", convertersFilename)
		specs = append(specs, builders.FileSpec{
			Name:     "converters",
			Filename: convertersFilename,
			Type:     "converters",
			Required: true,
			ContentHints: builders.ContentHints{
				HasServices: true,
			},
		})
	}

	// Generate exports file (main WASM wrapper)
	exportsFilename := filepath.Join(packagePath, baseName+"_exports.wasm.go")
	log.Printf("Planning exports file: %s", exportsFilename)
	specs = append(specs, builders.FileSpec{
		Name:     "exports",
		Filename: exportsFilename,
		Type:     "exports",
		Required: true,
		ContentHints: builders.ContentHints{
			HasServices:        len(data.Services) > 0,
			HasBrowserServices: data.HasBrowserClients,
		},
	})

	// Generate browser clients file (only if we have browser clients)
	if data.HasBrowserClients {
		browserClientsFilename := filepath.Join(packagePath, baseName+"_browser_clients.wasm.go")
		log.Printf("Planning browser clients file: %s", browserClientsFilename)
		specs = append(specs, builders.FileSpec{
			Name:     "browser_clients",
			Filename: browserClientsFilename,
			Type:     "browser_clients",
			Required: true,
			ContentHints: builders.ContentHints{
				HasBrowserServices: true,
			},
		})
	}

	// Always generate main example (helps users understand integration)
	mainFilename := gg.calculateMainFilename(data.PackageName, config)
	specs = append(specs, builders.FileSpec{
		Name:     "main",
		Filename: mainFilename,
		Type:     "example",
		Required: true,
		ContentHints: builders.ContentHints{
			IsExample: true,
		},
	})

	// Generate build script if enabled
	if config.GenerateBuildScript {
		buildFilename := gg.calculateBuildScriptFilename(config)
		specs = append(specs, builders.FileSpec{
			Name:     "build",
			Filename: buildFilename,
			Type:     "script",
			Required: false,
			ContentHints: builders.ContentHints{
				IsBuildScript: true,
			},
		})
	}

	return &builders.FilePlan{
		PackageName: data.PackageName,
		Specs:       specs,
		Config:      config,
	}
}


// calculateWasmFilename determines the output filename for the WASM wrapper.
func (gg *GoGenerator) calculateWasmFilename(packageName string, config *builders.GenerationConfig) string {
	packagePath := gg.pathCalc.BuildPackagePath(packageName)
	baseName := strings.ReplaceAll(packageName, ".", "_")
	return filepath.Join(packagePath, baseName+".wasm.go")
}

// calculateMainFilename determines the output filename for the main example.
func (gg *GoGenerator) calculateMainFilename(packageName string, config *builders.GenerationConfig) string {
	packagePath := gg.pathCalc.BuildPackagePath(packageName)
	return filepath.Join(packagePath, "main.go.example")
}

// calculateBuildScriptFilename determines the output filename for the build script.
func (gg *GoGenerator) calculateBuildScriptFilename(config *builders.GenerationConfig) string {
	return filepath.Join(config.WasmExportPath, "build.sh")
}

// calculateOutputPath determines the output directory path for generated files.
// It uses the go_package path (if available) to avoid collisions when multiple proto files
// have the same proto package but different go_package options.
func (gg *GoGenerator) calculateOutputPath(data *builders.GoTemplateData) string {
	if data.GoPackage != "" {
		// Extract the relative path from go_package
		// e.g., "github.com/.../gen/go/test/v1/models" -> "test/v1/models"
		parts := strings.Split(data.GoPackage, "/")

		// Find "gen/go" or similar marker and take everything after it
		for i, part := range parts {
			if part == "go" && i > 0 && (parts[i-1] == "gen" || parts[i-1] == "pb") {
				if i+1 < len(parts) {
					remainingPath := strings.Join(parts[i+1:], "/")
					return remainingPath
				}
			}
		}

		// Fallback: use last 2-3 components if no marker found
		if len(parts) >= 2 {
			return strings.Join(parts[len(parts)-2:], "/")
		}
	}

	// Fallback to proto package name
	return gg.pathCalc.BuildPackagePath(data.PackageName)
}

// calculateBaseName determines the base filename for generated files.
// It combines the proto package name with the go_package suffix to ensure uniqueness.
func (gg *GoGenerator) calculateBaseName(data *builders.GoTemplateData) string {
	baseName := strings.ReplaceAll(data.PackageName, ".", "_")

	// If go_package ends with a different suffix than proto package, include it
	if data.GoPackage != "" {
		parts := strings.Split(data.GoPackage, "/")
		if len(parts) > 0 {
			lastPart := parts[len(parts)-1]
			// If the last part is different from the package version, append it
			packageParts := strings.Split(data.PackageName, ".")
			if len(packageParts) > 0 {
				lastPackagePart := packageParts[len(packageParts)-1]
				if lastPart != lastPackagePart && lastPart != "" {
					// Avoid duplication - only append if not already in baseName
					if !strings.HasSuffix(baseName, "_"+lastPart) {
						baseName = baseName + "_" + lastPart
					}
				}
			}
		}
	}

	return baseName
}

// ValidateConfig validates the configuration for Go generation.
func (gg *GoGenerator) ValidateConfig(config *builders.GenerationConfig) error {
	if config.WasmExportPath == "" {
		return fmt.Errorf("WasmExportPath cannot be empty")
	}

	if config.JSStructure == "" {
		config.JSStructure = "namespaced" // Default
	}

	validStructures := map[string]bool{
		"namespaced":    true,
		"flat":          true,
		"service_based": true,
	}

	if !validStructures[config.JSStructure] {
		return fmt.Errorf("invalid JSStructure: %s (supported: namespaced, flat, service_based)", config.JSStructure)
	}

	return nil
}
