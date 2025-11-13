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

package generators

import (
	"log"
	"path/filepath"

	"google.golang.org/protobuf/compiler/protogen"

	"github.com/panyam/protoc-gen-go-wasmjs/pkg/builders"
	"github.com/panyam/protoc-gen-go-wasmjs/pkg/core"
	"github.com/panyam/protoc-gen-go-wasmjs/pkg/filters"
)

// ArtifactCatalog contains all artifacts discovered from proto files.
// It provides a complete view of everything available for generation, regardless of
// protoc's Generate flags. This enables cross-package artifact visibility for
// import resolution, factory composition, and schema registries.
//
// The catalog is populated by BaseGenerator.CollectAllArtifacts and then used by
// target-specific generators (TSGenerator, GoGenerator) to make file mapping decisions.
//
// Example:
//
//	catalog, err := baseGen.CollectAllArtifacts(config, criteria)
//	if err != nil {
//	    return err
//	}
//	// catalog.Services contains all regular services
//	// catalog.BrowserServices contains browser-provided services
//	// catalog.Messages contains messages grouped by package
//	// catalog.Enums contains enums grouped by package
//	// catalog.Packages maps package names to PackageInfo
type ArtifactCatalog struct {
	// Services contains all regular gRPC services to be generated.
	// These are services that will have WASM implementations.
	Services []ServiceArtifact

	// BrowserServices contains services marked with (wasmjs.v1.browser_provided) = true.
	// These services are implemented in JavaScript/TypeScript and called from WASM.
	BrowserServices []ServiceArtifact

	// Messages contains messages grouped by proto package.
	// Each MessageArtifact includes all messages for one package.
	Messages []MessageArtifact

	// Enums contains enums grouped by proto package.
	// Each EnumArtifact includes all enums for one package.
	Enums []EnumArtifact

	// Packages maps proto package names (e.g., "library.v1") to PackageInfo.
	// This provides metadata about each package including its files and path.
	Packages map[string]*builders.PackageInfo
}

// ServiceArtifact represents a gRPC service ready for code generation.
// It combines the service definition with package metadata and browser service detection.
//
// ServiceArtifacts are collected by BaseGenerator and used by target-specific generators
// to produce WASM wrappers (GoGenerator) or TypeScript clients (TSGenerator).
type ServiceArtifact struct {
	// Service is the protogen service definition from the proto file.
	Service *protogen.Service

	// Package provides metadata about the proto package containing this service.
	Package *builders.PackageInfo

	// IsBrowser indicates if this service has (wasmjs.v1.browser_provided) = true.
	// Browser services are implemented in JavaScript and called from WASM.
	IsBrowser bool
}

// MessageArtifact represents a collection of protobuf messages from one package.
// Messages are grouped by package to enable package-level file generation
// (e.g., one interfaces.ts file per package).
//
// MessageArtifacts are used by TSGenerator to create TypeScript interfaces, models,
// factories, schemas, and deserializers.
type MessageArtifact struct {
	// Messages contains all messages discovered in this package.
	Messages []filters.MessageInfo

	// Package provides metadata about the proto package.
	Package *builders.PackageInfo
}

// EnumArtifact represents a collection of protobuf enums from one package.
// Enums are grouped by package to enable package-level file generation.
//
// EnumArtifacts are used by TSGenerator to create TypeScript enums in interfaces.ts files.
type EnumArtifact struct {
	// Enums contains all enums discovered in this package.
	Enums []filters.EnumInfo

	// Package provides metadata about the proto package.
	Package *builders.PackageInfo
}

// BaseGenerator provides shared artifact collection and file planning utilities.
// It implements the foundation of the 4-step artifact processing approach and is
// embedded by target-specific generators (TSGenerator, GoGenerator) to access
// complete artifact catalogs.
//
// Architecture:
//
//	BaseGenerator embeds:
//	  - Core utilities (ProtoAnalyzer, PathCalculator, NameConverter)
//	  - Filter layer (PackageFilter, ServiceFilter, MethodFilter)
//	  - Collectors (MessageCollector, EnumCollector)
//
//	Target generators embed BaseGenerator:
//	  - TSGenerator: Generates TypeScript clients and types
//	  - GoGenerator: Generates Go WASM wrappers
//
// The BaseGenerator ensures all generators have access to the complete artifact
// catalog, enabling cross-package visibility for imports, factories, and schemas.
type BaseGenerator struct {
	// Core dependencies - pure utility functions without side effects
	analyzer *core.ProtoAnalyzer  // Proto file analysis (service detection, etc.)
	pathCalc *core.PathCalculator // Import path calculation
	nameConv *core.NameConverter  // Naming convention conversion

	// Filter layer - business logic for filtering artifacts
	packageFilter *filters.PackageFilter // Package-level filtering
	serviceFilter *filters.ServiceFilter // Service-level filtering
	methodFilter  *filters.MethodFilter  // Method-level filtering
	msgCollector  *filters.MessageCollector // Message collection
	enumCollector *filters.EnumCollector    // Enum collection

	// Generation context - protogen plugin instance
	plugin *protogen.Plugin
}

// NewBaseGenerator creates a new BaseGenerator with all required dependencies initialized.
// This sets up the complete filter and utility pipeline used for artifact collection.
//
// The returned BaseGenerator is ready to be embedded by target-specific generators
// (TSGenerator or GoGenerator).
//
// Example:
//
//	func NewTSGenerator(plugin *protogen.Plugin) *TSGenerator {
//	    return &TSGenerator{
//	        BaseGenerator: NewBaseGenerator(plugin),
//	    }
//	}
func NewBaseGenerator(plugin *protogen.Plugin) *BaseGenerator {
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

	return &BaseGenerator{
		analyzer:      analyzer,
		pathCalc:      pathCalc,
		nameConv:      nameConv,
		packageFilter: packageFilter,
		serviceFilter: serviceFilter,
		methodFilter:  methodFilter,
		msgCollector:  msgCollector,
		enumCollector: enumCollector,
		plugin:        plugin,
	}
}

// CollectAllArtifacts performs complete artifact discovery from all proto files.
// This is Step 1 of the 4-step artifact processing approach.
//
// The method collects artifacts from ALL proto files, regardless of protoc's Generate flags,
// providing generators with complete visibility for cross-package imports, factory composition,
// and schema registries. This is critical for TypeScript generation where imports may reference
// types from packages not being directly generated.
//
// Collection Process:
//
//  1. Build complete package map from all files
//  2. Collect all services (regular and browser-provided)
//  3. Collect all messages by package (if config.GenerateTypes)
//  4. Collect all enums by package (if config.GenerateTypes)
//
// The returned ArtifactCatalog is then passed to target-specific generators for
// file mapping decisions (Step 3) and rendering (Step 4).
//
// Parameters:
//
//	config: Generation configuration controlling what artifacts to collect
//	criteria: Filter criteria for service/method/message selection
//
// Returns:
//
//	ArtifactCatalog: Complete catalog of all discovered artifacts
//	error: Any error encountered during collection
//
// Example:
//
//	config := &builders.GenerationConfig{
//	    GenerateTypes: true,
//	    JSStructure: "namespaced",
//	}
//	criteria := &filters.FilterCriteria{
//	    // Service filtering configuration
//	}
//	catalog, err := baseGen.CollectAllArtifacts(config, criteria)
//	if err != nil {
//	    return fmt.Errorf("artifact collection failed: %w", err)
//	}
//	// Use catalog.Services, catalog.Messages, etc. for generation
func (bg *BaseGenerator) CollectAllArtifacts(config *builders.GenerationConfig, criteria *filters.FilterCriteria) (*ArtifactCatalog, error) {
	log.Printf("BaseGenerator: Collecting artifacts from %d proto files", len(bg.plugin.Files))

	catalog := &ArtifactCatalog{
		Packages: make(map[string]*builders.PackageInfo),
	}

	// Phase 1: Build complete package map from ALL files (ignore Generate flag for artifact collection)
	allPackageFiles := make(map[string][]*protogen.File)

	for _, file := range bg.plugin.Files {
		packageName := string(file.Desc.Package())

		// Collect ALL files regardless of Generate flag
		allPackageFiles[packageName] = append(allPackageFiles[packageName], file)
	}

	// Convert to PackageInfo and add to catalog
	// Use ALL files for complete package visibility, but track which have Generate=true
	for packageName, files := range allPackageFiles {
		if len(files) > 0 {
			packageInfo := &builders.PackageInfo{
				Name:  packageName,
				Path:  bg.pathCalc.BuildPackagePath(packageName),
				Files: files,
			}
			catalog.Packages[packageName] = packageInfo
		}
	}

	log.Printf("Collected packages from ALL files: %v", func() []string {
		var names []string
		for name := range catalog.Packages {
			names = append(names, name)
		}
		return names
	}())

	// Phase 2: Collect all services across all packages
	for packageName, packageInfo := range catalog.Packages {
		log.Printf("  Processing package: %s (%d files)", packageName, len(packageInfo.Files))

		for _, file := range packageInfo.Files {
			// IMPORTANT: Only collect services from files marked for generation
			if !file.Generate {
				log.Printf("    Skipping file %s (Generate=false)", file.Desc.Path())
				continue
			}

			for _, service := range file.Services {
				// Check if service should be included
				serviceResult := bg.serviceFilter.ShouldIncludeService(service, criteria)
				if !serviceResult.Include {
					continue
				}

				serviceArtifact := ServiceArtifact{
					Service:   service,
					Package:   packageInfo,
					IsBrowser: bg.analyzer.IsBrowserProvidedService(service),
				}

				if serviceArtifact.IsBrowser {
					catalog.BrowserServices = append(catalog.BrowserServices, serviceArtifact)
					log.Printf("    Found browser service: %s", service.GoName)
				} else {
					catalog.Services = append(catalog.Services, serviceArtifact)
					log.Printf("    Found regular service: %s", service.GoName)
				}
			}
		}

		// Collect messages and enums per directory (not per package)
		// This ensures each directory gets its own interfaces/models/schemas files
		if config.GenerateTypes {
			// Group files by directory - ONLY include files with Generate=true
			filesByDir := make(map[string][]*protogen.File)
			for _, file := range packageInfo.Files {
				// IMPORTANT: Only collect types from files marked for generation
				if !file.Generate {
					continue
				}
				protoPath := string(file.Desc.Path())
				dir := filepath.Dir(protoPath)
				filesByDir[dir] = append(filesByDir[dir], file)
			}

			// Create message/enum artifacts per directory
			for dir, filesInDir := range filesByDir {
				// Collect messages from files in this directory
				messageResult := bg.msgCollector.CollectMessages(filesInDir, criteria)
				if len(messageResult.Items) > 0 {
					// Create a directory-specific PackageInfo
					dirPackageInfo := &builders.PackageInfo{
						Name:  packageInfo.Name,
						Path:  packageInfo.Path,
						Files: filesInDir,
					}
					catalog.Messages = append(catalog.Messages, MessageArtifact{
						Messages: messageResult.Items,
						Package:  dirPackageInfo,
					})
					log.Printf("    Found %d messages in %s (directory: %s)", len(messageResult.Items), packageName, dir)
				}

				// Collect enums from files in this directory
				enumResult := bg.enumCollector.CollectEnums(filesInDir, criteria)
				if len(enumResult.Items) > 0 {
					// Create a directory-specific PackageInfo
					dirPackageInfo := &builders.PackageInfo{
						Name:  packageInfo.Name,
						Path:  packageInfo.Path,
						Files: filesInDir,
					}
					catalog.Enums = append(catalog.Enums, EnumArtifact{
						Enums:   enumResult.Items,
						Package: dirPackageInfo,
					})
					log.Printf("    Found %d enums in %s (directory: %s)", len(enumResult.Items), packageName, dir)
				}
			}
		}
	}

	log.Printf("Artifact collection complete: %d services, %d browser services, %d message groups, %d enum groups across %d packages",
		len(catalog.Services), len(catalog.BrowserServices), len(catalog.Messages), len(catalog.Enums), len(catalog.Packages))

	return catalog, nil
}

// GetPackageInfo retrieves package information for a given package name
func (catalog *ArtifactCatalog) GetPackageInfo(packageName string) *builders.PackageInfo {
	return catalog.Packages[packageName]
}

// GetServicesForPackage returns all services (regular + browser) for a specific package
func (catalog *ArtifactCatalog) GetServicesForPackage(packageName string) ([]ServiceArtifact, []ServiceArtifact) {
	var services []ServiceArtifact
	var browserServices []ServiceArtifact

	for _, svc := range catalog.Services {
		if svc.Package.Name == packageName {
			services = append(services, svc)
		}
	}

	for _, bsvc := range catalog.BrowserServices {
		if bsvc.Package.Name == packageName {
			browserServices = append(browserServices, bsvc)
		}
	}

	return services, browserServices
}

// HasServicesForModule checks if any services exist for the given module
func (catalog *ArtifactCatalog) HasServicesForModule() bool {
	return len(catalog.Services) > 0 || len(catalog.BrowserServices) > 0
}
