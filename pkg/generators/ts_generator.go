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

	"google.golang.org/protobuf/compiler/protogen"

	"github.com/panyam/protoc-gen-go-wasmjs/pkg/builders"
	"github.com/panyam/protoc-gen-go-wasmjs/pkg/core"
	"github.com/panyam/protoc-gen-go-wasmjs/pkg/filters"
	"github.com/panyam/protoc-gen-go-wasmjs/pkg/renderers"
)

// TSGenerator orchestrates the complete TypeScript generation process.
// This is the top-level generator that coordinates all layers to produce TypeScript artifacts.
type TSGenerator struct {
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
	dataBuilder *builders.TSDataBuilder
	renderer    *renderers.TSRenderer

	// Generation context
	plugin *protogen.Plugin
}

// NewTSGenerator creates a new TypeScript generator with all necessary dependencies.
// This sets up the complete processing pipeline for TypeScript generation.
func NewTSGenerator(plugin *protogen.Plugin) *TSGenerator {
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
	dataBuilder := builders.NewTSDataBuilder(analyzer, pathCalc, nameConv, serviceFilter, methodFilter, msgCollector, enumCollector)
	renderer := renderers.NewTSRenderer()

	return &TSGenerator{
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

// Generate performs the complete TypeScript generation process.
// This is the main entry point that coordinates all layers to produce TypeScript artifacts.
func (tg *TSGenerator) Generate(config *builders.GenerationConfig, filterCriteria *filters.FilterCriteria) error {
	// Phase 1: Filter packages
	packageFiles, stats := tg.packageFilter.FilterPackages(tg.plugin.Files, filterCriteria)

	if len(packageFiles) == 0 {
		return nil // No packages to process
	}

	// Phase 2: Generate TypeScript artifacts for each package
	for packageName, files := range packageFiles {
		packageInfo := &builders.PackageInfo{
			Name:  packageName,
			Path:  tg.pathCalc.BuildPackagePath(packageName),
			Files: files,
		}

		// Generate files for this package
		if err := tg.generatePackageFiles(packageInfo, filterCriteria, config); err != nil {
			return fmt.Errorf("failed to generate files for package %s: %w", packageName, err)
		}
	}

	// Log generation summary (not to stdout to avoid corrupting protobuf response)
	log.Printf("TypeScript generator processed %s", stats.Summary())

	return nil
}

// generatePackageFiles handles complete file generation for a package using file planning.
// This is the new approach where the generator controls all file creation and naming.
func (tg *TSGenerator) generatePackageFiles(
	packageInfo *builders.PackageInfo,
	criteria *filters.FilterCriteria,
	config *builders.GenerationConfig,
) error {

	// Phase 1: Plan what files we need to generate
	filePlan := tg.planTSFiles(packageInfo, criteria, config)

	if len(filePlan.Specs) == 0 {
		return nil // No files to generate
	}

	// Phase 2: Create all GeneratedFile objects
	fileSet := builders.NewGeneratedFileSet(filePlan, tg.plugin)

	// Phase 3: Validate file set
	if err := fileSet.ValidateFileSet(); err != nil {
		return fmt.Errorf("file planning validation failed: %w", err)
	}

	// Phase 4: Render each file
	if err := tg.renderFilesFromPlan(fileSet, packageInfo, criteria, config); err != nil {
		return fmt.Errorf("file rendering failed: %w", err)
	}

	return nil
}

// planTSFiles creates a file plan for TypeScript generation based on package content and configuration.
// This method decides what TypeScript files to generate and where to put them.
func (tg *TSGenerator) planTSFiles(
	packageInfo *builders.PackageInfo,
	criteria *filters.FilterCriteria,
	config *builders.GenerationConfig,
) *builders.FilePlan {

	var specs []builders.FileSpec
	hasServices := tg.serviceFilter.HasAnyServices(packageInfo.Files, criteria)
	hasTypes := tg.hasTypesToGenerate(packageInfo.Files, criteria)

	// Plan client file if package has services
	if hasServices {
		clientFilename := tg.calculateClientFilename(packageInfo, config)
		specs = append(specs, builders.FileSpec{
			Name:     "client",
			Filename: clientFilename,
			Type:     "client",
			Required: true,
			ContentHints: builders.ContentHints{
				HasServices: true,
			},
		})
	}

	// Plan BrowserServiceManager if we have any services
	// TODO: This should be moved to an npm package in the future
	// For now, generate it alongside the first client file to avoid duplicates
	if hasServices {
		specs = append(specs, builders.FileSpec{
			Name:     "browserServiceManager",
			Filename: "browserServiceManager.ts",
			Type:     "browserServiceManager",
			Required: false, // Not required if no browser services
			ContentHints: builders.ContentHints{
				HasServices: true,
			},
		})
	}

	// Plan type files if package has messages/enums
	if hasTypes {
		// Plan interfaces file
		interfacesFilename := tg.calculateInterfacesFilename(packageInfo, config)
		specs = append(specs, builders.FileSpec{
			Name:     "interfaces",
			Filename: interfacesFilename,
			Type:     "interfaces",
			Required: true,
			ContentHints: builders.ContentHints{
				HasMessages: true,
				HasEnums:    true,
			},
		})

		// Plan models file
		modelsFilename := tg.calculateModelsFilename(packageInfo, config)
		specs = append(specs, builders.FileSpec{
			Name:     "models",
			Filename: modelsFilename,
			Type:     "models",
			Required: false,
			ContentHints: builders.ContentHints{
				HasMessages: true,
			},
		})

		// Plan factory file (optional)
		if config.GenerateFactories {
			factoryFilename := tg.calculateFactoryFilename(packageInfo, config)
			specs = append(specs, builders.FileSpec{
				Name:     "factory",
				Filename: factoryFilename,
				Type:     "factory",
				Required: false,
				ContentHints: builders.ContentHints{
					HasMessages: true,
				},
			})
		}

		// Plan schemas file
		schemasFilename := tg.calculateSchemasFilename(packageInfo, config)
		specs = append(specs, builders.FileSpec{
			Name:     "schemas",
			Filename: schemasFilename,
			Type:     "schemas",
			Required: false,
			ContentHints: builders.ContentHints{
				HasMessages: true,
			},
		})

		// Plan deserializer file
		deserializerFilename := tg.calculateDeserializerFilename(packageInfo, config)
		specs = append(specs, builders.FileSpec{
			Name:     "deserializer",
			Filename: deserializerFilename,
			Type:     "deserializer",
			Required: false,
			ContentHints: builders.ContentHints{
				HasMessages: true,
			},
		})
	}

	return &builders.FilePlan{
		PackageName: packageInfo.Name,
		Specs:       specs,
		Config:      config,
	}
}

// renderFilesFromPlan executes the file plan by rendering all planned files.
func (tg *TSGenerator) renderFilesFromPlan(
	fileSet *builders.GeneratedFileSet,
	packageInfo *builders.PackageInfo,
	criteria *filters.FilterCriteria,
	config *builders.GenerationConfig,
) error {

	// Render client file if planned
	if clientFile := fileSet.GetFile("client"); clientFile != nil {
		clientData, err := tg.dataBuilder.BuildClientData(packageInfo, criteria, config)
		if err != nil {
			return fmt.Errorf("failed to build client data: %w", err)
		}

		if clientData != nil {
			if err := tg.renderer.RenderClient(clientFile, clientData); err != nil {
				return fmt.Errorf("failed to render client: %w", err)
			}
		}
	}

	// Render BrowserServiceManager if planned
	// TODO: Move this to an npm package in the future
	if bsmFile := fileSet.GetFile("browserServiceManager"); bsmFile != nil {
		// BrowserServiceManager doesn't need specific data, it's a generic utility
		if err := tg.renderer.RenderBrowserServiceManager(bsmFile); err != nil {
			return fmt.Errorf("failed to render BrowserServiceManager: %w", err)
		}
	}

	// Render type files if planned
	typeFiles := fileSet.GetFilesByType("interfaces")
	if len(typeFiles) > 0 || fileSet.HasFile("models") || fileSet.HasFile("factory") {
		typeData, err := tg.dataBuilder.BuildTypeData(packageInfo, criteria, config)
		if err != nil {
			return fmt.Errorf("failed to build type data: %w", err)
		}

		if typeData != nil {
			// Render interfaces
			if interfacesFile := fileSet.GetFile("interfaces"); interfacesFile != nil {
				if err := tg.renderer.RenderInterfaces(interfacesFile, typeData); err != nil {
					return fmt.Errorf("failed to render interfaces: %w", err)
				}
			}

			// Render models
			if modelsFile := fileSet.GetFile("models"); modelsFile != nil {
				if err := tg.renderer.RenderModels(modelsFile, typeData); err != nil {
					return fmt.Errorf("failed to render models: %w", err)
				}
			}

			// Render factory
			if factoryFile := fileSet.GetFile("factory"); factoryFile != nil {
				if err := tg.renderer.RenderFactory(factoryFile, typeData); err != nil {
					return fmt.Errorf("failed to render factory: %w", err)
				}
			}

			// Render schemas
			if schemasFile := fileSet.GetFile("schemas"); schemasFile != nil {
				if err := tg.renderer.RenderSchemas(schemasFile, typeData); err != nil {
					return fmt.Errorf("failed to render schemas: %w", err)
				}
			}

			// Render deserializer
			if deserializerFile := fileSet.GetFile("deserializer"); deserializerFile != nil {
				if err := tg.renderer.RenderDeserializer(deserializerFile, typeData); err != nil {
					return fmt.Errorf("failed to render deserializer: %w", err)
				}
			}
		}
	}

	return nil
}

// hasTypesToGenerate checks if a package has messages or enums that need TypeScript types.
func (tg *TSGenerator) hasTypesToGenerate(files []*protogen.File, criteria *filters.FilterCriteria) bool {
	return tg.msgCollector.HasAnyMessages(files, criteria) || tg.enumCollector.HasAnyEnums(files, criteria)
}

// File naming calculation methods - Generator controls all file naming decisions

// calculateClientFilename determines the output filename for the TypeScript client.
func (tg *TSGenerator) calculateClientFilename(packageInfo *builders.PackageInfo, config *builders.GenerationConfig) string {
	moduleName := tg.getModuleName(packageInfo.Name, config)
	return moduleName + "Client.ts"
}

// calculateInterfacesFilename determines the output filename for TypeScript interfaces.
func (tg *TSGenerator) calculateInterfacesFilename(packageInfo *builders.PackageInfo, config *builders.GenerationConfig) string {
	return filepath.Join(packageInfo.Path, "interfaces.ts")
}

// calculateModelsFilename determines the output filename for TypeScript model classes.
func (tg *TSGenerator) calculateModelsFilename(packageInfo *builders.PackageInfo, config *builders.GenerationConfig) string {
	return filepath.Join(packageInfo.Path, "models.ts")
}

// calculateFactoryFilename determines the output filename for TypeScript factory classes.
func (tg *TSGenerator) calculateFactoryFilename(packageInfo *builders.PackageInfo, config *builders.GenerationConfig) string {
	return filepath.Join(packageInfo.Path, "factory.ts")
}

// calculateSchemasFilename determines the output filename for TypeScript schemas.
func (tg *TSGenerator) calculateSchemasFilename(packageInfo *builders.PackageInfo, config *builders.GenerationConfig) string {
	return filepath.Join(packageInfo.Path, "schemas.ts")
}

// calculateDeserializerFilename determines the output filename for TypeScript deserializers.
func (tg *TSGenerator) calculateDeserializerFilename(packageInfo *builders.PackageInfo, config *builders.GenerationConfig) string {
	return filepath.Join(packageInfo.Path, "deserializer.ts")
}

// getModuleName determines the TypeScript module name.
func (tg *TSGenerator) getModuleName(packageName string, config *builders.GenerationConfig) string {
	if config.ModuleName != "" {
		return config.ModuleName
	}
	return tg.nameConv.ToModuleName(packageName)
}

// ValidateConfig validates the configuration for TypeScript generation.
func (tg *TSGenerator) ValidateConfig(config *builders.GenerationConfig) error {
	if config.TSExportPath == "" {
		return fmt.Errorf("TSExportPath cannot be empty")
	}

	return nil
}
