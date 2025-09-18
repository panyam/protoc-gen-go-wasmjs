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
	"strings"
)

// FilterCriteria defines the configuration and criteria used for filtering 
// services, methods, messages, and enums during code generation.
// This centralizes all filtering logic in a single, testable structure.
type FilterCriteria struct {
	// Service filtering
	ServicesSet map[string]bool // Specific services to include (empty = all)
	
	// Method filtering with glob patterns
	MethodIncludes []string // Glob patterns for methods to include
	MethodExcludes []string // Glob patterns for methods to exclude
	
	// Method name transformations
	MethodRenames map[string]string // Original name -> Custom name mappings
	
	// Package filtering
	ExcludeAnnotationPackages bool     // Skip wasmjs.v1 and similar annotation packages
	ExcludeEmptyPackages     bool     // Skip packages with no services/messages/enums
	
	// Message filtering  
	ExcludeMapEntries    bool // Skip synthetic map entry messages
	ExcludeNestedMessages bool // Skip nested messages (collect only top-level)
	
	// Enum filtering
	ExcludeNestedEnums bool // Skip nested enums (collect only top-level)
}

// NewFilterCriteria creates FilterCriteria with sensible defaults for typical usage.
// These defaults match the current generator behavior to ensure backward compatibility.
func NewFilterCriteria() *FilterCriteria {
	return &FilterCriteria{
		ServicesSet:               make(map[string]bool),
		MethodIncludes:           []string{},
		MethodExcludes:           []string{},
		MethodRenames:            make(map[string]string),
		ExcludeAnnotationPackages: true,  // Always skip wasmjs.v1 packages
		ExcludeEmptyPackages:     true,  // Skip packages with no content
		ExcludeMapEntries:        true,  // Skip synthetic map entry messages
		ExcludeNestedMessages:    false, // Include nested messages by default
		ExcludeNestedEnums:       false, // Include nested enums by default
	}
}

// ServiceFilterCriteria creates criteria specifically for service filtering.
// This is used when you only need to filter services and don't care about messages/enums.
func ServiceFilterCriteria(servicesSet map[string]bool) *FilterCriteria {
	criteria := NewFilterCriteria()
	criteria.ServicesSet = servicesSet
	return criteria
}

// MethodFilterCriteria creates criteria specifically for method filtering.
// This is used when you need fine-grained control over which methods are generated.
func MethodFilterCriteria(includes, excludes []string, renames map[string]string) *FilterCriteria {
	criteria := NewFilterCriteria()
	criteria.MethodIncludes = includes
	criteria.MethodExcludes = excludes
	criteria.MethodRenames = renames
	return criteria
}

// HasServiceFilter returns true if specific services are configured for filtering.
// When false, all services should be included (subject to other filters).
func (fc *FilterCriteria) HasServiceFilter() bool {
	return len(fc.ServicesSet) > 0
}

// HasMethodIncludes returns true if specific method include patterns are configured.
// When true, only methods matching include patterns should be generated.
func (fc *FilterCriteria) HasMethodIncludes() bool {
	return len(fc.MethodIncludes) > 0
}

// HasMethodExcludes returns true if method exclude patterns are configured.
// When true, methods matching exclude patterns should be skipped.
func (fc *FilterCriteria) HasMethodExcludes() bool {
	return len(fc.MethodExcludes) > 0
}

// HasMethodRenames returns true if method name transformations are configured.
// When true, method names should be checked for custom renames.
func (fc *FilterCriteria) HasMethodRenames() bool {
	return len(fc.MethodRenames) > 0
}

// GetMethodRename returns the custom name for a method if configured.
// Returns the original name if no rename is configured.
func (fc *FilterCriteria) GetMethodRename(originalName string) string {
	if renamed, exists := fc.MethodRenames[originalName]; exists {
		return renamed
	}
	return originalName
}

// ParseFromConfig creates FilterCriteria from the generator configuration.
// This bridges the gap between the existing config system and the new filter layer.
func ParseFromConfig(servicesStr, methodIncludeStr, methodExcludeStr, methodRenameStr string) (*FilterCriteria, error) {
	criteria := NewFilterCriteria()
	
	// Parse services list
	if servicesStr != "" {
		for _, service := range parseCommaSeparated(servicesStr) {
			criteria.ServicesSet[service] = true
		}
	}
	
	// Parse method includes
	criteria.MethodIncludes = parseCommaSeparated(methodIncludeStr)
	
	// Parse method excludes  
	criteria.MethodExcludes = parseCommaSeparated(methodExcludeStr)
	
	// Parse method renames
	if methodRenameStr != "" {
		renames, err := parseMethodRenames(methodRenameStr)
		if err != nil {
			return nil, err
		}
		criteria.MethodRenames = renames
	}
	
	return criteria, nil
}

// parseCommaSeparated splits a comma-separated string and trims whitespace.
// Returns empty slice for empty input.
func parseCommaSeparated(str string) []string {
	if str == "" {
		return []string{}
	}
	
	var result []string
	for _, item := range strings.Split(str, ",") {
		item = strings.TrimSpace(item)
		if item != "" {
			result = append(result, item)
		}
	}
	return result
}

// parseMethodRenames parses method rename specifications.
// Format: "OldName1:NewName1,OldName2:NewName2"
// Returns error for invalid format.
func parseMethodRenames(renameStr string) (map[string]string, error) {
	renames := make(map[string]string)
	
	for _, rename := range parseCommaSeparated(renameStr) {
		parts := strings.SplitN(rename, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid method rename format: %s (expected OldName:NewName)", rename)
		}
		
		oldName := strings.TrimSpace(parts[0])
		newName := strings.TrimSpace(parts[1])
		
		if oldName == "" || newName == "" {
			return nil, fmt.Errorf("invalid method rename: empty old or new name in %s", rename)
		}
		
		renames[oldName] = newName
	}
	
	return renames, nil
}
