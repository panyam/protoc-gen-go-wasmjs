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

	"google.golang.org/protobuf/compiler/protogen"

	"github.com/panyam/protoc-gen-go-wasmjs/pkg/builders"
	"github.com/panyam/protoc-gen-go-wasmjs/pkg/core"
	"github.com/panyam/protoc-gen-go-wasmjs/pkg/filters"
)

// ArtifactCatalog contains all artifacts discovered from proto files
// This provides a complete view of everything available for generation
type ArtifactCatalog struct {
	Services        []ServiceArtifact
	BrowserServices []ServiceArtifact
	Messages        []MessageArtifact
	Enums           []EnumArtifact
	Packages        map[string]*builders.PackageInfo
}

// ServiceArtifact represents a service ready for generation
type ServiceArtifact struct {
	Service   *protogen.Service
	Package   *builders.PackageInfo
	IsBrowser bool
}

// MessageArtifact represents messages ready for generation
type MessageArtifact struct {
	Messages []filters.MessageInfo
	Package  *builders.PackageInfo
}

// EnumArtifact represents enums ready for generation
type EnumArtifact struct {
	Enums   []filters.EnumInfo
	Package *builders.PackageInfo
}

// BaseGenerator provides shared artifact collection and file planning utilities.
// Each target-specific generator (TS, Go) embeds this to get access to complete artifacts.
type BaseGenerator struct {
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

	// Generation context
	plugin *protogen.Plugin
}

// NewBaseGenerator creates a new base generator with all required dependencies.
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
// This gives each generator a complete view of everything available for generation.
func (bg *BaseGenerator) CollectAllArtifacts(config *builders.GenerationConfig, criteria *filters.FilterCriteria) (*ArtifactCatalog, error) {
	log.Printf("BaseGenerator: Collecting artifacts from %d proto files", len(bg.plugin.Files))

	catalog := &ArtifactCatalog{
		Packages: make(map[string]*builders.PackageInfo),
	}

	// Phase 1: Build complete package map from ALL files (ignore Generate flag for artifact collection)
	allPackageFiles := make(map[string][]*protogen.File)
	for _, file := range bg.plugin.Files {
		// Include ALL files for artifact collection, not just those marked for generation
		packageName := string(file.Desc.Package())
		// Skip only system packages that we definitely don't want
		if packageName == "google.protobuf" {
			continue
		}
		allPackageFiles[packageName] = append(allPackageFiles[packageName], file)
	}

	// Convert to PackageInfo and add to catalog
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

		// Collect messages if needed
		if config.GenerateTypes {
			messageResult := bg.msgCollector.CollectMessages(packageInfo.Files, criteria)
			if len(messageResult.Items) > 0 {
				catalog.Messages = append(catalog.Messages, MessageArtifact{
					Messages: messageResult.Items,
					Package:  packageInfo,
				})
				log.Printf("    Found %d messages in %s", len(messageResult.Items), packageName)
			}
		}

		// Collect enums if needed
		if config.GenerateTypes {
			enumResult := bg.enumCollector.CollectEnums(packageInfo.Files, criteria)
			if len(enumResult.Items) > 0 {
				catalog.Enums = append(catalog.Enums, EnumArtifact{
					Enums:   enumResult.Items,
					Package: packageInfo,
				})
				log.Printf("    Found %d enums in %s", len(enumResult.Items), packageName)
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
