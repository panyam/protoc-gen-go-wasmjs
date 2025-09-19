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

package filters

import (
	"strings"

	"google.golang.org/protobuf/compiler/protogen"

	"github.com/panyam/protoc-gen-go-wasmjs/pkg/core"
)

// PackageFilter provides filtering logic for protobuf packages.
// This determines which packages should be processed for code generation
// based on their content and configuration criteria.
type PackageFilter struct {
	analyzer         *core.ProtoAnalyzer
	messageCollector *MessageCollector
	enumCollector    *EnumCollector
}

// NewPackageFilter creates a new package filter with the necessary dependencies.
func NewPackageFilter(analyzer *core.ProtoAnalyzer, messageCollector *MessageCollector, enumCollector *EnumCollector) *PackageFilter {
	return &PackageFilter{
		analyzer:         analyzer,
		messageCollector: messageCollector,
		enumCollector:    enumCollector,
	}
}

// ShouldIncludePackage determines if a package should be processed for generation.
// This applies filtering logic to determine if a package is worth generating code for.
// Returns PackageFilterResult with detailed reasoning about the decision.
func (pf *PackageFilter) ShouldIncludePackage(packageName string, files []*protogen.File, criteria *FilterCriteria) PackageFilterResult {
	// Check for annotation package exclusion (wasmjs.v1, etc.)
	if criteria.ExcludeAnnotationPackages && pf.isAnnotationPackage(packageName) {
		return PackageFilterResult{
			FilterResult: Excluded("package is an annotation package (wasmjs.v1, etc.)"),
		}
	}

	// Analyze package content
	hasServices := pf.hasAnyServices(files)
	hasMessages := pf.messageCollector.HasAnyMessages(files, criteria)
	hasEnums := pf.enumCollector.HasAnyEnums(files, criteria)

	// Check for empty package exclusion
	if criteria.ExcludeEmptyPackages && !hasServices && !hasMessages && !hasEnums {
		return PackageFilterResult{
			FilterResult: Excluded("package has no services, messages, or enums"),
			HasServices:  hasServices,
			HasMessages:  hasMessages,
			HasEnums:     hasEnums,
		}
	}

	// Include packages that have content
	reason := "package has content: "
	var reasons []string
	if hasServices {
		reasons = append(reasons, "services")
	}
	if hasMessages {
		reasons = append(reasons, "messages")
	}
	if hasEnums {
		reasons = append(reasons, "enums")
	}

	if len(reasons) > 0 {
		reason += strings.Join(reasons, ", ")
	} else {
		reason = "package included despite no content (empty package exclusion disabled)"
	}

	return PackageFilterResult{
		FilterResult: Included(reason),
		HasServices:  hasServices,
		HasMessages:  hasMessages,
		HasEnums:     hasEnums,
	}
}

// isAnnotationPackage checks if a package is an annotation package that should be excluded.
// Annotation packages like wasmjs.v1 are library files used for proto annotations,
// not user code that should have generated artifacts.
func (pf *PackageFilter) isAnnotationPackage(packageName string) bool {
	// Known annotation packages that should be excluded from generation
	annotationPackages := []string{
		"wasmjs.v1",
		"google.protobuf", // Standard protobuf types (usually handled specially)
	}

	for _, annoPkg := range annotationPackages {
		if packageName == annoPkg {
			return true
		}
	}

	return false
}

// hasAnyServices checks if any of the files contain gRPC services.
func (pf *PackageFilter) hasAnyServices(files []*protogen.File) bool {
	for _, file := range files {
		if len(file.Services) > 0 {
			return true
		}
	}
	return false
}

// FilterPackages filters a collection of packages (represented as file groups).
// This applies package-level filtering to determine which packages should be processed.
// Returns a map of package names to their files, containing only packages that should be processed.
func (pf *PackageFilter) FilterPackages(allFiles []*protogen.File, criteria *FilterCriteria) (map[string][]*protogen.File, *FilterStats) {
	packageFiles := make(map[string][]*protogen.File)
	stats := NewFilterStats()

	// Group files by package first
	allPackageFiles := make(map[string][]*protogen.File)
	for _, file := range allFiles {
		if !file.Generate {
			continue // Skip files not marked for generation
		}

		packageName := string(file.Desc.Package())
		allPackageFiles[packageName] = append(allPackageFiles[packageName], file)
	}

	// Apply package filtering
	for packageName, files := range allPackageFiles {
		result := pf.ShouldIncludePackage(packageName, files, criteria)
		stats.PackagesTotal++

		if result.Include {
			packageFiles[packageName] = files

			// Update content statistics
			if result.HasServices {
				// Count actual services (will be refined by service filtering)
				for _, file := range files {
					stats.ServicesTotal += len(file.Services)
				}
			}
		}
	}

	return packageFiles, stats
}

// GetPackageNames extracts package names from filtered packages.
// This is useful for template generation and reporting.
func (pf *PackageFilter) GetPackageNames(packageFiles map[string][]*protogen.File) []string {
	var packageNames []string
	for packageName := range packageFiles {
		packageNames = append(packageNames, packageName)
	}
	return packageNames
}

// GetPackageContentSummary analyzes what types of content each package contains.
// This provides insight into the structure of filtered packages.
func (pf *PackageFilter) GetPackageContentSummary(packageFiles map[string][]*protogen.File, criteria *FilterCriteria) map[string]PackageFilterResult {
	summary := make(map[string]PackageFilterResult)

	for packageName, files := range packageFiles {
		result := pf.ShouldIncludePackage(packageName, files, criteria)
		summary[packageName] = result
	}

	return summary
}
