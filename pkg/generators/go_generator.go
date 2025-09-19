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
	"github.com/panyam/protoc-gen-go-wasmjs/pkg/core"
	"github.com/panyam/protoc-gen-go-wasmjs/pkg/filters"
	"github.com/panyam/protoc-gen-go-wasmjs/pkg/renderers"
)

// GoGenerator orchestrates the complete Go WASM generation process.
// This is the top-level generator that coordinates all layers to produce Go WASM artifacts.
type GoGenerator struct {
	// Core dependencies
	analyzer *core.ProtoAnalyzer
	pathCalc *core.PathCalculator
	nameConv *core.NameConverter

	// Filter layer
	packageFilter *filters.PackageFilter
	serviceFilter *filters.ServiceFilter
	methodFilter  *filters.MethodFilter
	msgCollector  *filters.MessageCollector
	enumCollector *filters.EnumCollector

	// Builder and renderer
	dataBuilder *builders.GoDataBuilder
	renderer    *renderers.GoRenderer

	// Generation context
	plugin *protogen.Plugin
}

// NewGoGenerator creates a new Go generator with all necessary dependencies.
// This sets up the complete processing pipeline for Go WASM generation.
func NewGoGenerator(plugin *protogen.Plugin) *GoGenerator {
	// Create core utilities
	analyzer := core.NewProtoAnalyzer()
	pathCalc := core.NewPathCalculator()
	nameConv := core.NewNameConverter()

	// Create filter layer
	msgCollector := filters.NewMessageCollector(analyzer)
	enumCollector := filters.NewEnumCollector(analyzer)
	packageFilter := filters.NewPackageFilter(analyzer, msgCollector, enumCollector)
	serviceFilter := filters.NewServiceFilter(analyzer)
	methodFilter := filters.NewMethodFilter(analyzer)

	// Create builder and renderer
	dataBuilder := builders.NewGoDataBuilder(analyzer, pathCalc, nameConv, serviceFilter, methodFilter, msgCollector, enumCollector)
	renderer := renderers.NewGoRenderer()

	return &GoGenerator{
		analyzer:      analyzer,
		pathCalc:      pathCalc,
		nameConv:      nameConv,
		packageFilter: packageFilter,
		serviceFilter: serviceFilter,
		methodFilter:  methodFilter,
		msgCollector:  msgCollector,
		enumCollector: enumCollector,
		dataBuilder:   dataBuilder,
		renderer:      renderer,
		plugin:        plugin,
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
		
		log.Printf("Generated template data for package %s: %d services, %d browser services", 
			packageName, len(templateData.Services), len(templateData.BrowserServices))

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
		case "wasm":
			log.Printf("WASM: Attempting to render WASM wrapper...")
			if err := gg.renderer.RenderWasmWrapperDirect(generatedFile, data); err != nil {
				log.Printf("WASM: ERROR rendering WASM wrapper: %v", err)
				return fmt.Errorf("failed to render WASM file %s: %w", spec.Filename, err)
			}
			log.Printf("WASM: WASM wrapper rendered successfully")
			
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

	// Always generate WASM wrapper (this is the main artifact)
	wasmFilename := gg.calculateWasmFilename(data.PackageName, config)
	log.Printf("Planning WASM file: %s", wasmFilename)
	specs = append(specs, builders.FileSpec{
		Name:     "wasm",
		Filename: wasmFilename,
		Type:     "wasm",
		Required: true,
		ContentHints: builders.ContentHints{
			HasServices:        len(data.Services) > 0,
			HasBrowserServices: data.HasBrowserServices,
		},
	})

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

// renderFilesFromPlan executes the file plan by rendering all planned files.
func (gg *GoGenerator) renderFilesFromPlan(fileSet *builders.GeneratedFileSet, data *builders.GoTemplateData, config *builders.GenerationConfig) error {
	// Render WASM wrapper using direct approach like old generator
	if wasmFile := fileSet.GetFile("wasm"); wasmFile != nil {
		log.Printf("Rendering WASM wrapper directly to file...")
		if err := gg.renderer.RenderWasmWrapperDirect(wasmFile, data); err != nil {
			return fmt.Errorf("failed to render WASM wrapper: %w", err)
		}
		log.Printf("WASM wrapper rendered directly to file successfully")
	}

	// Render main example
	if mainFile := fileSet.GetFile("main"); mainFile != nil {
		content, err := gg.renderer.RenderMainExample(data)
		if err != nil {
			return fmt.Errorf("failed to render main example: %w", err)
		}
		if content != "" {
			mainFile.P(content)
		}
	}

	// Render build script if planned
	if buildFile := fileSet.GetFile("build"); buildFile != nil {
		content, err := gg.renderer.RenderBuildScript(data)
		if err != nil {
			return fmt.Errorf("failed to render build script: %w", err)
		}
		if content != "" {
			buildFile.P(content)
		}
	}

	return nil
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
