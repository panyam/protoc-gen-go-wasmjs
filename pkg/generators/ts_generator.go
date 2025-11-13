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

// TSGenerator orchestrates the complete TypeScript generation process.
// This is the top-level generator that coordinates all layers to produce TypeScript artifacts.
type TSGenerator struct {
	// Embed base generator for artifact collection
	*BaseGenerator

	// Builder and renderer specific to TypeScript
	dataBuilder *builders.TSDataBuilder
	renderer    *renderers.TSRenderer
}

// NewTSGenerator creates a new TypeScript generator with all necessary dependencies.
// This sets up the complete processing pipeline for TypeScript generation.
func NewTSGenerator(plugin *protogen.Plugin) *TSGenerator {
	// Create base generator with artifact collection capabilities
	baseGenerator := NewBaseGenerator(plugin)

	// Create TypeScript-specific builder and renderer
	dataBuilder := builders.NewTSDataBuilder(
		baseGenerator.analyzer, 
		baseGenerator.pathCalc, 
		baseGenerator.nameConv, 
		baseGenerator.serviceFilter, 
		baseGenerator.methodFilter, 
		baseGenerator.msgCollector, 
		baseGenerator.enumCollector,
	)
	renderer := renderers.NewTSRenderer()

	return &TSGenerator{
		BaseGenerator: baseGenerator,
		dataBuilder:   dataBuilder,
		renderer:      renderer,
	}
}

// Generate performs the complete TypeScript generation process.
// This uses BaseGenerator to collect all artifacts first, then maps them to files.
func (tg *TSGenerator) Generate(config *builders.GenerationConfig, filterCriteria *filters.FilterCriteria) error {
	// Phase 1: Collect all artifacts from all packages
	catalog, err := tg.CollectAllArtifacts(config, filterCriteria)
	if err != nil {
		return fmt.Errorf("artifact collection failed: %w", err)
	}

	// Phase 2: Plan files based on artifacts (TypeScript-specific slice/dice/group logic)
	filePlan := tg.planFilesFromCatalog(catalog, config)
	
	if len(filePlan.Specs) == 0 {
		return nil // No files to generate
	}

	// Phase 3: Create file set structure (without protogen files yet)
	fileSet := builders.NewGeneratedFileSet(filePlan)

	// Phase 4: Create actual protogen files after all mapping decisions are made
	if err := fileSet.CreateFiles(tg.plugin); err != nil {
		return fmt.Errorf("file creation failed: %w", err)
	}

	// Phase 5: Validate file set
	if err := fileSet.ValidateFileSet(); err != nil {
		return fmt.Errorf("file planning validation failed: %w", err)
	}

	// Phase 6: Render all files
	if err := tg.renderFilesFromCatalog(fileSet, catalog, config, filterCriteria); err != nil {
		return fmt.Errorf("file rendering failed: %w", err)
	}

	log.Printf("TypeScript generator processed %d services, %d browser services across %d packages",
		len(catalog.Services), len(catalog.BrowserServices), len(catalog.Packages))

	return nil
}

// planFilesFromCatalog creates a file plan based on the complete artifact catalog.
// This is where TypeScript-specific slice/dice/group logic happens.
func (tg *TSGenerator) planFilesFromCatalog(catalog *ArtifactCatalog, config *builders.GenerationConfig) *builders.FilePlan {
	var specs []builders.FileSpec

	// Plan service client files (one per service, organized by package)
	for _, svcArtifact := range catalog.Services {
		serviceClientFilename := tg.calculateServiceClientFilename(svcArtifact.Package, svcArtifact.Service, config)
		specs = append(specs, builders.FileSpec{
			Name:     fmt.Sprintf("client_%s_%s", svcArtifact.Package.Name, svcArtifact.Service.GoName),
			Filename: serviceClientFilename,
			Type:     "service_client",
			Required: true,
			ContentHints: builders.ContentHints{
				HasServices: true,
			},
			Metadata: map[string]interface{}{
				"service":     svcArtifact.Service,
				"packageInfo": svcArtifact.Package,
			},
		})
	}

	// Plan browser service client files  
	for _, browserSvcArtifact := range catalog.BrowserServices {
		serviceClientFilename := tg.calculateServiceClientFilename(browserSvcArtifact.Package, browserSvcArtifact.Service, config)
		specs = append(specs, builders.FileSpec{
			Name:     fmt.Sprintf("client_%s_%s", browserSvcArtifact.Package.Name, browserSvcArtifact.Service.GoName),
			Filename: serviceClientFilename,
			Type:     "service_client",
			Required: true,
			ContentHints: builders.ContentHints{
				HasServices:        true,
				HasBrowserServices: true,
			},
			Metadata: map[string]interface{}{
				"service":     browserSvcArtifact.Service,
				"packageInfo": browserSvcArtifact.Package,
			},
		})
	}

	// Always plan module-level bundle file (simple base class with module config)
	// Generate bundle once per module - protoc will deduplicate automatically
	specs = append(specs, builders.FileSpec{
		Name:     "bundle",
		Filename: "index.ts", // Module-level bundle
		Type:     "bundle",
		Required: true,
		ContentHints: builders.ContentHints{
			HasServices: false, // Simple bundle doesn't contain services
		},
	})

	// Plan type files per package
	// Track which packages we've already planned to avoid duplicates
	processedPackages := make(map[string]bool)

	for _, msgArtifact := range catalog.Messages {
		packageInfo := msgArtifact.Package

		// Skip if we've already processed this package
		if processedPackages[packageInfo.Name] {
			continue
		}
		processedPackages[packageInfo.Name] = true

		// Interfaces file
		interfacesFilename := tg.calculateInterfacesFilename(packageInfo, config)
		specs = append(specs, builders.FileSpec{
			Name:     fmt.Sprintf("interfaces_%s", packageInfo.Name),
			Filename: interfacesFilename,
			Type:     "interfaces",
			Required: true,
			ContentHints: builders.ContentHints{
				HasMessages: true,
			},
			Metadata: map[string]interface{}{
				"packageInfo": packageInfo,
			},
		})

		// Models file
		modelsFilename := tg.calculateModelsFilename(packageInfo, config)
		specs = append(specs, builders.FileSpec{
			Name:     fmt.Sprintf("models_%s", packageInfo.Name),
			Filename: modelsFilename,
			Type:     "models",
			Required: false,
			ContentHints: builders.ContentHints{
				HasMessages: true,
			},
			Metadata: map[string]interface{}{
				"packageInfo": packageInfo,
			},
		})

		// Schemas file (keep this - it's a basic unit)
		schemasFilename := tg.calculateSchemasFilename(packageInfo, config)
		specs = append(specs, builders.FileSpec{
			Name:     fmt.Sprintf("schemas_%s", packageInfo.Name),
			Filename: schemasFilename,
			Type:     "schemas",
			Required: false,
			ContentHints: builders.ContentHints{
				HasMessages: true,
			},
			Metadata: map[string]interface{}{
				"packageInfo": packageInfo,
			},
		})

		// TODO: Factory and Deserializer are disabled for now since they need
		// package-wide aggregation across multiple buf invocations
		// We'll implement these later with a different strategy

		// Factory file (DISABLED - requires package-wide aggregation)
		// if config.GenerateFactories {
		// 	factoryFilename := tg.calculateFactoryFilename(packageInfo, config)
		// 	specs = append(specs, builders.FileSpec{
		// 		Name:     fmt.Sprintf("factory_%s", packageInfo.Name),
		// 		Filename: factoryFilename,
		// 		Type:     "factory",
		// 		Required: false,
		// 		ContentHints: builders.ContentHints{
		// 			HasMessages: true,
		// 		},
		// 		Metadata: map[string]interface{}{
		// 			"packageInfo": packageInfo,
		// 		},
		// 	})
		// }

		// Deserializer file (DISABLED - requires package-wide aggregation)
		// deserializerFilename := tg.calculateDeserializerFilename(packageInfo, config)
		// specs = append(specs, builders.FileSpec{
		// 	Name:     fmt.Sprintf("deserializer_%s", packageInfo.Name),
		// 	Filename: deserializerFilename,
		// 	Type:     "deserializer",
		// 	Required: false,
		// 	ContentHints: builders.ContentHints{
		// 		HasMessages: true,
		// 	},
		// 	Metadata: map[string]interface{}{
		// 		"packageInfo": packageInfo,
		// 	},
		// })
	}

	return &builders.FilePlan{
		PackageName: "module", // This is module-level, not package-level
		Specs:       specs,
		Config:      config,
	}
}

// renderFilesFromCatalog renders all files using the artifact catalog.
func (tg *TSGenerator) renderFilesFromCatalog(
	fileSet *builders.GeneratedFileSet,
	catalog *ArtifactCatalog,
	config *builders.GenerationConfig,
	criteria *filters.FilterCriteria,
) error {
	// Render service client files
	serviceClientFiles := fileSet.GetFilesByType("service_client")
	for fileName, serviceFile := range serviceClientFiles {
		spec := fileSet.GetFileSpec(fileName)
		if spec != nil && spec.Metadata != nil {
			service := spec.Metadata["service"].(*protogen.Service)
			packageInfo := spec.Metadata["packageInfo"].(*builders.PackageInfo)

			// Build service client data for this specific service
			serviceClientData, err := tg.dataBuilder.BuildServiceClientData(packageInfo, service, criteria, config)
			if err != nil {
				return fmt.Errorf("failed to build service client data for %s: %w", service.GoName, err)
			}

			if serviceClientData != nil {
				if err := tg.renderer.RenderServiceClient(serviceFile, serviceClientData); err != nil {
					return fmt.Errorf("failed to render service client %s: %w", service.GoName, err)
				}
			}
		}
	}

	// Render module-level bundle file
	if bundleFile := fileSet.GetFile("bundle"); bundleFile != nil {
		bundleData, err := tg.buildBundleDataFromCatalog(catalog, config)
		if err != nil {
			return fmt.Errorf("failed to build bundle data: %w", err)
		}

		if bundleData != nil {
			if err := tg.renderer.RenderBundle(bundleFile, bundleData); err != nil {
				return fmt.Errorf("failed to render bundle: %w", err)
			}
		}
	}

	// Render type files
	// Track which packages we've already rendered type files for to avoid duplicates
	renderedPackages := make(map[string]*builders.TSTemplateData)

	interfaceFiles := fileSet.GetFilesByType("interfaces")
	for fileName, interfaceFile := range interfaceFiles {
		spec := fileSet.GetFileSpec(fileName)
		if spec != nil && spec.Metadata != nil {
			packageInfo := spec.Metadata["packageInfo"].(*builders.PackageInfo)

			// Build type data for this package (cache it for reuse)
			typeData, err := tg.dataBuilder.BuildTypeData(packageInfo, criteria, config)
			if err != nil {
				return fmt.Errorf("failed to build type data for %s: %w", packageInfo.Name, err)
			}

			if typeData != nil {
				renderedPackages[packageInfo.Name] = typeData

				if err := tg.renderer.RenderInterfaces(interfaceFile, typeData); err != nil {
					return fmt.Errorf("failed to render interfaces for %s: %w", packageInfo.Name, err)
				}
			}
		}
	}

	// Render models files
	modelsFiles := fileSet.GetFilesByType("models")
	for fileName, modelsFile := range modelsFiles {
		spec := fileSet.GetFileSpec(fileName)
		if spec != nil && spec.Metadata != nil {
			packageInfo := spec.Metadata["packageInfo"].(*builders.PackageInfo)

			// Reuse type data if we already built it
			typeData := renderedPackages[packageInfo.Name]
			if typeData == nil {
				var err error
				typeData, err = tg.dataBuilder.BuildTypeData(packageInfo, criteria, config)
				if err != nil {
					return fmt.Errorf("failed to build type data for %s: %w", packageInfo.Name, err)
				}
				renderedPackages[packageInfo.Name] = typeData
			}

			if typeData != nil {
				if err := tg.renderer.RenderModels(modelsFile, typeData); err != nil {
					return fmt.Errorf("failed to render models for %s: %w", packageInfo.Name, err)
				}
			}
		}
	}

	// Render factory files
	factoryFiles := fileSet.GetFilesByType("factory")
	for fileName, factoryFile := range factoryFiles {
		spec := fileSet.GetFileSpec(fileName)
		if spec != nil && spec.Metadata != nil {
			packageInfo := spec.Metadata["packageInfo"].(*builders.PackageInfo)

			// Reuse type data if we already built it
			typeData := renderedPackages[packageInfo.Name]
			if typeData == nil {
				var err error
				typeData, err = tg.dataBuilder.BuildTypeData(packageInfo, criteria, config)
				if err != nil {
					return fmt.Errorf("failed to build type data for %s: %w", packageInfo.Name, err)
				}
				renderedPackages[packageInfo.Name] = typeData
			}

			if typeData != nil {
				if err := tg.renderer.RenderFactory(factoryFile, typeData); err != nil {
					return fmt.Errorf("failed to render factory for %s: %w", packageInfo.Name, err)
				}
			}
		}
	}

	// Render schemas files
	schemasFiles := fileSet.GetFilesByType("schemas")
	for fileName, schemasFile := range schemasFiles {
		spec := fileSet.GetFileSpec(fileName)
		if spec != nil && spec.Metadata != nil {
			packageInfo := spec.Metadata["packageInfo"].(*builders.PackageInfo)

			// Reuse type data if we already built it
			typeData := renderedPackages[packageInfo.Name]
			if typeData == nil {
				var err error
				typeData, err = tg.dataBuilder.BuildTypeData(packageInfo, criteria, config)
				if err != nil {
					return fmt.Errorf("failed to build type data for %s: %w", packageInfo.Name, err)
				}
				renderedPackages[packageInfo.Name] = typeData
			}

			if typeData != nil {
				if err := tg.renderer.RenderSchemas(schemasFile, typeData); err != nil {
					return fmt.Errorf("failed to render schemas for %s: %w", packageInfo.Name, err)
				}
			}
		}
	}

	// Render deserializer files
	deserializerFiles := fileSet.GetFilesByType("deserializer")
	for fileName, deserializerFile := range deserializerFiles {
		spec := fileSet.GetFileSpec(fileName)
		if spec != nil && spec.Metadata != nil {
			packageInfo := spec.Metadata["packageInfo"].(*builders.PackageInfo)

			// Reuse type data if we already built it
			typeData := renderedPackages[packageInfo.Name]
			if typeData == nil {
				var err error
				typeData, err = tg.dataBuilder.BuildTypeData(packageInfo, criteria, config)
				if err != nil {
					return fmt.Errorf("failed to build type data for %s: %w", packageInfo.Name, err)
				}
				renderedPackages[packageInfo.Name] = typeData
			}

			if typeData != nil {
				if err := tg.renderer.RenderDeserializer(deserializerFile, typeData); err != nil {
					return fmt.Errorf("failed to render deserializer for %s: %w", packageInfo.Name, err)
				}
			}
		}
	}

	return nil
}

// buildBundleDataFromCatalog creates bundle template data with just module configuration.
// The simplified bundle only needs module config - no service information needed.
func (tg *TSGenerator) buildBundleDataFromCatalog(catalog *ArtifactCatalog, config *builders.GenerationConfig) (*builders.TSTemplateData, error) {
	// Build minimal bundle template data - just module configuration
	return &builders.TSTemplateData{
		PackageName:  "module",                           // Module-level bundle
		PackagePath:  ".",                                // Root level path
		ModuleName:   tg.getModuleName("", config),       // Module-level name
		APIStructure: config.JSStructure,                 // Pass-through configuration
		JSNamespace:  config.JSNamespace,                 // Pass-through configuration
		Services:     []builders.ServiceData{},           // No services needed for simple bundle
		Messages:     []builders.TSMessageInfo{},         // No messages needed
		Enums:        []builders.TSEnumInfo{},             // No enums needed
		// Minimal flags to satisfy validation
		HasBrowserServices: false,
		HasBrowserClients:  false,
		HasMessages:        false,
		HasEnums:           false,
	}, nil
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

	// Phase 2: Create file set structure (without protogen files yet)
	fileSet := builders.NewGeneratedFileSet(filePlan)

	// Phase 3: Create actual protogen files after all mapping decisions are made
	if err := fileSet.CreateFiles(tg.plugin); err != nil {
		return fmt.Errorf("file creation failed: %w", err)
	}

	// Phase 4: Validate file set
	if err := fileSet.ValidateFileSet(); err != nil {
		return fmt.Errorf("file planning validation failed: %w", err)
	}

	// Phase 5: Render each file
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

	// Plan separate client files per service
	if hasServices {
		// Get filtered services for this package
		filteredServices := tg.serviceFilter.GetIncludedServices(packageInfo.Files, criteria)
		
		for _, service := range filteredServices {
			serviceClientFilename := tg.calculateServiceClientFilename(packageInfo, service, config)
			specs = append(specs, builders.FileSpec{
				Name:     fmt.Sprintf("client_%s", service.GoName),
				Filename: serviceClientFilename,
				Type:     "service_client",
				Required: true,
				ContentHints: builders.ContentHints{
					HasServices: true,
				},
				// Store service info for template data building
				Metadata: map[string]interface{}{
					"service": service,
				},
			})
		}

		// Plan a single bundle file per module (shared by all services)
		// This should be generated at the root level (e.g., generated/index.ts)
		bundleFilename := tg.calculateBundleFilename(packageInfo, config)
		specs = append(specs, builders.FileSpec{
			Name:     "bundle",
			Filename: bundleFilename,
			Type:     "bundle",
			Required: true,
			ContentHints: builders.ContentHints{
				HasServices: true,
			},
		})
	}

	// BrowserServiceManager is now imported from @protoc-gen-go-wasmjs/runtime package
	// No longer generating it as a file

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

		// Plan schemas file (keep this - it's a basic unit)
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

		// TODO: Factory and Deserializer disabled - require package-wide aggregation
		// Plan factory file (optional) - DISABLED
		// if config.GenerateFactories {
		// 	factoryFilename := tg.calculateFactoryFilename(packageInfo, config)
		// 	specs = append(specs, builders.FileSpec{
		// 		Name:     "factory",
		// 		Filename: factoryFilename,
		// 		Type:     "factory",
		// 		Required: false,
		// 		ContentHints: builders.ContentHints{
		// 			HasMessages: true,
		// 		},
		// 	})
		// }

		// Plan deserializer file - DISABLED
		// deserializerFilename := tg.calculateDeserializerFilename(packageInfo, config)
		// specs = append(specs, builders.FileSpec{
		// 	Name:     "deserializer",
		// 	Filename: deserializerFilename,
		// 	Type:     "deserializer",
		// 	Required: false,
		// 	ContentHints: builders.ContentHints{
		// 		HasMessages: true,
		// 	},
		// })
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

	// Render service client files (one per service)
	serviceClientFiles := fileSet.GetFilesByType("service_client")
	for fileName, serviceFile := range serviceClientFiles {
		// Get the service info from metadata
		spec := fileSet.GetFileSpec(fileName)
		if spec != nil && spec.Metadata != nil {
			if service, ok := spec.Metadata["service"].(*protogen.Service); ok {
				// Build client data for this specific service
				serviceClientData, err := tg.dataBuilder.BuildServiceClientData(packageInfo, service, criteria, config)
				if err != nil {
					return fmt.Errorf("failed to build service client data for %s: %w", service.GoName, err)
				}

				if serviceClientData != nil {
					if err := tg.renderer.RenderServiceClient(serviceFile, serviceClientData); err != nil {
						return fmt.Errorf("failed to render service client %s: %w", service.GoName, err)
					}
				}
			}
		}
	}

	// Render bundle file if planned
	if bundleFile := fileSet.GetFile("bundle"); bundleFile != nil {
		// Build bundle data (includes all services in this package)
		bundleData, err := tg.dataBuilder.BuildClientData(packageInfo, criteria, config)
		if err != nil {
			return fmt.Errorf("failed to build bundle data: %w", err)
		}

		if bundleData != nil {
			if err := tg.renderer.RenderBundle(bundleFile, bundleData); err != nil {
				return fmt.Errorf("failed to render bundle: %w", err)
			}
		}
	}

	// BrowserServiceManager is now imported from @protoc-gen-go-wasmjs/runtime package
	// No longer rendering BrowserServiceManager - it's imported from runtime package

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



// calculateServiceClientFilename determines the output filename for a specific service client.
func (tg *TSGenerator) calculateServiceClientFilename(packageInfo *builders.PackageInfo, service *protogen.Service, config *builders.GenerationConfig) string {
	// Generate file in the proto file's directory structure
	// e.g., presenter/v1/services/presenterServiceClient.ts (if service is in services/ folder)
	serviceFileName := tg.convertToFileName(service.GoName) + "Client.ts"
	dir := tg.getProtoFileDirectory(packageInfo)
	return filepath.Join(dir, serviceFileName)
}

// convertToFileName converts a service name to a filename-friendly format
func (tg *TSGenerator) convertToFileName(serviceName string) string {
	// Convert PascalCase to camelCase for filenames
	// e.g., "PresenterService" -> "presenterService"
	if len(serviceName) == 0 {
		return serviceName
	}
	return strings.ToLower(serviceName[:1]) + serviceName[1:]
}

// getProtoFileDirectory extracts the directory path from proto files.
// Uses the proto source file path (not the go_package path) to determine output directory.
func (tg *TSGenerator) getProtoFileDirectory(packageInfo *builders.PackageInfo) string {
	// Find files with Generate=true in this invocation
	var generateFiles []*protogen.File
	for _, file := range packageInfo.Files {
		if file.Generate {
			generateFiles = append(generateFiles, file)
		}
	}

	if len(generateFiles) == 0 {
		// Fallback: use first file
		if len(packageInfo.Files) > 0 {
			protoPath := string(packageInfo.Files[0].Desc.Path())
			return filepath.Dir(protoPath)
		}
		// Last resort: use package path
		return packageInfo.Path
	}

	// Use the directory of the first Generate=true file's proto source path
	// file.Desc.Path() returns the proto file path like "test_one_package/v1/models/test_service.proto"
	// We want the directory: "test_one_package/v1/models"
	protoPath := string(generateFiles[0].Desc.Path())
	dir := filepath.Dir(protoPath)

	return dir
}

// calculateInterfacesFilename determines the output filename for TypeScript interfaces.
// Uses proto file directory structure to avoid collisions when buf invokes the plugin
// multiple times for the same package.
func (tg *TSGenerator) calculateInterfacesFilename(packageInfo *builders.PackageInfo, config *builders.GenerationConfig) string {
	dir := tg.getProtoFileDirectory(packageInfo)
	return filepath.Join(dir, "interfaces.ts")
}

// calculateModelsFilename determines the output filename for TypeScript model classes.
func (tg *TSGenerator) calculateModelsFilename(packageInfo *builders.PackageInfo, config *builders.GenerationConfig) string {
	dir := tg.getProtoFileDirectory(packageInfo)
	return filepath.Join(dir, "models.ts")
}

// calculateFactoryFilename determines the output filename for TypeScript factory classes.
func (tg *TSGenerator) calculateFactoryFilename(packageInfo *builders.PackageInfo, config *builders.GenerationConfig) string {
	dir := tg.getProtoFileDirectory(packageInfo)
	return filepath.Join(dir, "factory.ts")
}

// calculateSchemasFilename determines the output filename for TypeScript schemas.
func (tg *TSGenerator) calculateSchemasFilename(packageInfo *builders.PackageInfo, config *builders.GenerationConfig) string {
	dir := tg.getProtoFileDirectory(packageInfo)
	return filepath.Join(dir, "schemas.ts")
}

// calculateDeserializerFilename determines the output filename for TypeScript deserializers.
func (tg *TSGenerator) calculateDeserializerFilename(packageInfo *builders.PackageInfo, config *builders.GenerationConfig) string {
	dir := tg.getProtoFileDirectory(packageInfo)
	return filepath.Join(dir, "deserializer.ts")
}

// calculateBundleFilename determines the output filename for the TypeScript bundle.
// For now, generate one bundle per package alongside the service files.
func (tg *TSGenerator) calculateBundleFilename(packageInfo *builders.PackageInfo, config *builders.GenerationConfig) string {
	// Place bundle alongside service files in the package directory
	return filepath.Join(packageInfo.Path, "index.ts")
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

	// Set default JSStructure if not specified
	if config.JSStructure == "" {
		config.JSStructure = "namespaced" // Default
	}

	// Validate JSStructure
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
