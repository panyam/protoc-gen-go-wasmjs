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
	"fmt"
	"path/filepath"
	
	"google.golang.org/protobuf/compiler/protogen"
	
	"github.com/panyam/protoc-gen-go-wasmjs/pkg/core"
)

// MethodFilter provides filtering logic for gRPC methods.
// This determines which methods within a service should be included in generation
// based on glob patterns, annotations, and streaming support.
type MethodFilter struct {
	analyzer *core.ProtoAnalyzer
}

// NewMethodFilter creates a new method filter with the given analyzer.
func NewMethodFilter(analyzer *core.ProtoAnalyzer) *MethodFilter {
	return &MethodFilter{
		analyzer: analyzer,
	}
}

// ShouldIncludeMethod determines if a method should be included in generation.
// This applies filtering criteria in order of precedence:
// 1. Annotation-based exclusion (highest priority)  
// 2. Streaming method limitations (client streaming not supported)
// 3. Explicit exclude patterns
// 4. Explicit include patterns (if configured)
// 5. Default inclusion
//
// Returns MethodFilterResult with detailed reasoning and method metadata.
func (mf *MethodFilter) ShouldIncludeMethod(method *protogen.Method, criteria *FilterCriteria) MethodFilterResult {
	methodName := string(method.Desc.Name())
	
	// Check annotation-based exclusion first (highest priority)
	if mf.analyzer.IsMethodExcluded(method) {
		return MethodFilterResult{
			FilterResult: Excluded("method marked with wasm_method_exclude annotation"),
		}
	}
	
	// Check streaming limitations - client streaming not supported
	if method.Desc.IsStreamingClient() {
		return MethodFilterResult{
			FilterResult: Excluded("client streaming methods not supported"),
		}
	}
	
	// Get method metadata from annotations
	customJSName := mf.analyzer.GetCustomMethodName(method)
	isAsync := mf.analyzer.IsAsyncMethod(method)
	isServerStreaming := method.Desc.IsStreamingServer()
	
	// Apply exclude patterns
	if criteria.HasMethodExcludes() {
		for _, pattern := range criteria.MethodExcludes {
			if matched, _ := filepath.Match(pattern, methodName); matched {
				return MethodFilterResult{
					FilterResult: Excluded("method matches exclude pattern: " + pattern),
				}
			}
		}
	}
	
	// Apply include patterns (if configured)
	if criteria.HasMethodIncludes() {
		for _, pattern := range criteria.MethodIncludes {
			if matched, _ := filepath.Match(pattern, methodName); matched {
				return MethodFilterResult{
					FilterResult:      Included("method matches include pattern: " + pattern),
					CustomJSName:      customJSName,
					IsAsync:           isAsync,
					IsServerStreaming: isServerStreaming,
				}
			}
		}
		// If includes are specified but no patterns matched, exclude
		return MethodFilterResult{
			FilterResult: Excluded("method doesn't match any include patterns"),
		}
	}
	
	// Default: include methods that aren't explicitly excluded
	return MethodFilterResult{
		FilterResult:      Included("method included by default (no exclusion rules matched)"),
		CustomJSName:      customJSName,
		IsAsync:           isAsync,
		IsServerStreaming: isServerStreaming,
	}
}

// FilterMethods filters all methods from a service based on criteria.
// This is a convenience method for filtering all methods in a service at once.
// Returns methods that should be included along with their metadata.
func (mf *MethodFilter) FilterMethods(service *protogen.Service, criteria *FilterCriteria) ([]MethodFilterResult, *FilterStats) {
	var results []MethodFilterResult
	stats := NewFilterStats()
	
	for _, method := range service.Methods {
		result := mf.ShouldIncludeMethod(method, criteria)
		stats.AddMethodResult(result)
		
		if result.Include {
			results = append(results, result)
		}
	}
	
	return results, stats
}

// GetIncludedMethods returns only the methods that passed filtering.
// This extracts the actual protogen.Method objects from filter results.
func (mf *MethodFilter) GetIncludedMethods(service *protogen.Service, criteria *FilterCriteria) []*protogen.Method {
	var methods []*protogen.Method
	
	for _, method := range service.Methods {
		result := mf.ShouldIncludeMethod(method, criteria)
		if result.Include {
			methods = append(methods, method)
		}
	}
	
	return methods
}

// HasAnyMethods checks if any methods in a service would be included with the given criteria.
// This is useful for determining if a service is worth generating at all.
func (mf *MethodFilter) HasAnyMethods(service *protogen.Service, criteria *FilterCriteria) bool {
	for _, method := range service.Methods {
		result := mf.ShouldIncludeMethod(method, criteria)
		if result.Include {
			return true
		}
	}
	return false
}

// GetMethodJSName returns the JavaScript name for a method, applying renames and annotations.
// This combines configuration-based renames with annotation-based custom names.
// Annotation names take precedence over configuration renames.
func (mf *MethodFilter) GetMethodJSName(method *protogen.Method, criteria *FilterCriteria, nameConverter *core.NameConverter) string {
	methodName := string(method.Desc.Name())
	
	// Check for annotation-based custom name first (highest priority)
	if customName := mf.analyzer.GetCustomMethodName(method); customName != "" {
		return customName
	}
	
	// Check for configuration-based rename
	if criteria.HasMethodRenames() {
		if renamed := criteria.GetMethodRename(methodName); renamed != methodName {
			return renamed
		}
	}
	
	// Default: convert to camelCase for JavaScript
	return nameConverter.ToCamelCase(methodName)
}

// ValidateMethodPatterns validates that all glob patterns in the criteria are valid.
// This should be called during configuration validation to catch invalid patterns early.
func (mf *MethodFilter) ValidateMethodPatterns(criteria *FilterCriteria) error {
	// Validate include patterns
	for _, pattern := range criteria.MethodIncludes {
		if _, err := filepath.Match(pattern, "test"); err != nil {
			return fmt.Errorf("invalid method include pattern '%s': %w", pattern, err)
		}
	}
	
	// Validate exclude patterns
	for _, pattern := range criteria.MethodExcludes {
		if _, err := filepath.Match(pattern, "test"); err != nil {
			return fmt.Errorf("invalid method exclude pattern '%s': %w", pattern, err)
		}
	}
	
	return nil
}
