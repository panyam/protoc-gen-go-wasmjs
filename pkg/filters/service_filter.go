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
	"google.golang.org/protobuf/compiler/protogen"

	"github.com/panyam/protoc-gen-go-wasmjs/pkg/core"
)

// ServiceFilter provides filtering logic for gRPC services.
// This determines which services should be included in code generation
// based on configuration criteria and protobuf annotations.
type ServiceFilter struct {
	analyzer *core.ProtoAnalyzer
}

// NewServiceFilter creates a new service filter with the given analyzer.
func NewServiceFilter(analyzer *core.ProtoAnalyzer) *ServiceFilter {
	return &ServiceFilter{
		analyzer: analyzer,
	}
}

// ShouldIncludeService determines if a service should be included in generation.
// This applies multiple filtering criteria in order of precedence:
// 1. Annotation-based exclusion (highest priority)
// 2. Explicit service list filtering
// 3. Default inclusion
//
// Returns ServiceFilterResult with detailed reasoning for the decision.
func (sf *ServiceFilter) ShouldIncludeService(service *protogen.Service, criteria *FilterCriteria) ServiceFilterResult {
	serviceName := string(service.Desc.Name())

	// Check annotation-based exclusion first (highest priority)
	if sf.analyzer.IsServiceExcluded(service) {
		return ServiceFilterResult{
			FilterResult: Excluded("service marked with wasm_service_exclude annotation"),
		}
	}

	// Check if service is browser-provided (affects generation but doesn't exclude)
	isBrowserProvided := sf.analyzer.IsBrowserProvidedService(service)

	// Get custom service name if specified
	customName := sf.analyzer.GetCustomServiceName(service)

	// Apply service list filtering if configured
	if criteria.HasServiceFilter() {
		if criteria.ServicesSet[serviceName] {
			return ServiceFilterResult{
				FilterResult:      Included("service explicitly included in services list"),
				IsBrowserProvided: isBrowserProvided,
				CustomName:        customName,
			}
		} else {
			return ServiceFilterResult{
				FilterResult: Excluded("service not in configured services list"),
			}
		}
	}

	// Default: include all services that aren't explicitly excluded
	return ServiceFilterResult{
		FilterResult:      Included("service included by default (no exclusion rules matched)"),
		IsBrowserProvided: isBrowserProvided,
		CustomName:        customName,
	}
}

// FilterServices filters a collection of services from multiple files.
// This is a convenience method for batch filtering operations.
// Returns all services that should be included along with filtering statistics.
func (sf *ServiceFilter) FilterServices(files []*protogen.File, criteria *FilterCriteria) ([]ServiceFilterResult, *FilterStats) {
	var results []ServiceFilterResult
	stats := NewFilterStats()

	for _, file := range files {
		for _, service := range file.Services {
			result := sf.ShouldIncludeService(service, criteria)
			stats.AddServiceResult(result)

			if result.Include {
				results = append(results, result)
			}
		}
	}

	return results, stats
}

// GetIncludedServices returns only the services that passed filtering.
// This extracts the actual protogen.Service objects from filter results.
func (sf *ServiceFilter) GetIncludedServices(files []*protogen.File, criteria *FilterCriteria) []*protogen.Service {
	var services []*protogen.Service

	for _, file := range files {
		for _, service := range file.Services {
			result := sf.ShouldIncludeService(service, criteria)
			if result.Include {
				services = append(services, service)
			}
		}
	}

	return services
}

// HasAnyServices checks if any services would be included with the given criteria.
// This is useful for early termination when no services would be generated.
func (sf *ServiceFilter) HasAnyServices(files []*protogen.File, criteria *FilterCriteria) bool {
	for _, file := range files {
		for _, service := range file.Services {
			result := sf.ShouldIncludeService(service, criteria)
			if result.Include {
				return true
			}
		}
	}
	return false
}

// GetBrowserProvidedServices returns only browser-provided services that passed filtering.
// Browser-provided services require special handling in the generation process.
func (sf *ServiceFilter) GetBrowserProvidedServices(files []*protogen.File, criteria *FilterCriteria) []*protogen.Service {
	var browserServices []*protogen.Service

	for _, file := range files {
		for _, service := range file.Services {
			result := sf.ShouldIncludeService(service, criteria)
			if result.Include && result.IsBrowserProvided {
				browserServices = append(browserServices, service)
			}
		}
	}

	return browserServices
}

// GetRegularServices returns only regular (non-browser-provided) services that passed filtering.
// Regular services are implemented in Go WASM and exposed to JavaScript.
func (sf *ServiceFilter) GetRegularServices(files []*protogen.File, criteria *FilterCriteria) []*protogen.Service {
	var regularServices []*protogen.Service

	for _, file := range files {
		for _, service := range file.Services {
			result := sf.ShouldIncludeService(service, criteria)
			if result.Include && !result.IsBrowserProvided {
				regularServices = append(regularServices, service)
			}
		}
	}

	return regularServices
}
