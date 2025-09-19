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
	"testing"

	"github.com/panyam/protoc-gen-go-wasmjs/pkg/core"
)

// TestFilterLayer_Integration demonstrates how the filter layer components work together.
// This integration test shows the complete filtering workflow without requiring
// complex protogen mocks, focusing on the coordination between filter components.
func TestFilterLayer_Integration(t *testing.T) {
	// Create filter layer components
	analyzer := core.NewProtoAnalyzer()
	messageCollector := NewMessageCollector(analyzer)
	enumCollector := NewEnumCollector(analyzer)
	packageFilter := NewPackageFilter(analyzer, messageCollector, enumCollector)
	serviceFilter := NewServiceFilter(analyzer)
	methodFilter := NewMethodFilter(analyzer)

	// Test component initialization
	if packageFilter == nil || serviceFilter == nil || methodFilter == nil {
		t.Error("Filter components should initialize successfully")
		return
	}

	// Test filter criteria creation and parsing
	criteria, err := ParseFromConfig(
		"UserService,LibraryService",              // services
		"Get*,Find*",                              // includes
		"*Internal,*Debug",                        // excludes
		"FindBooks:searchBooks,GetUser:fetchUser", // renames
	)

	if err != nil {
		t.Errorf("ParseFromConfig() failed: %v", err)
		return
	}

	// Validate parsed criteria
	if !criteria.HasServiceFilter() {
		t.Error("Should have service filter with configured services")
	}

	if !criteria.HasMethodIncludes() {
		t.Error("Should have method includes")
	}

	if !criteria.HasMethodExcludes() {
		t.Error("Should have method excludes")
	}

	if !criteria.HasMethodRenames() {
		t.Error("Should have method renames")
	}

	// Test method rename lookup
	if criteria.GetMethodRename("FindBooks") != "searchBooks" {
		t.Error("FindBooks should rename to searchBooks")
	}

	if criteria.GetMethodRename("GetUser") != "fetchUser" {
		t.Error("GetUser should rename to fetchUser")
	}

	if criteria.GetMethodRename("UnknownMethod") != "UnknownMethod" {
		t.Error("Unknown methods should return original name")
	}

	// Test filter stats
	stats := NewFilterStats()
	stats.ServicesTotal = 5
	stats.ServicesIncluded = 3
	stats.ServicesExcluded = 2
	stats.MethodsTotal = 20
	stats.MethodsIncluded = 15
	stats.MethodsExcluded = 5
	stats.MessagesTotal = 10
	stats.EnumsTotal = 3
	stats.PackagesTotal = 2

	summary := stats.Summary()
	expected := "Filtering Summary: 3/5 services, 15/20 methods, 10 messages, 3 enums from 2 packages"
	if summary != expected {
		t.Errorf("FilterStats.Summary() = %s, want %s", summary, expected)
	}

	// Test pattern validation
	if err := methodFilter.ValidateMethodPatterns(criteria); err != nil {
		t.Errorf("ValidateMethodPatterns() failed: %v", err)
	}

	// Test filter result creation
	includedResult := Included("test inclusion")
	if !includedResult.Include {
		t.Error("Included() should create result with Include=true")
	}

	excludedResult := Excluded("test exclusion")
	if excludedResult.Include {
		t.Error("Excluded() should create result with Include=false")
	}
}

// TestFilterCriteria_ComplexScenarios tests complex filtering scenarios that combine
// multiple criteria types. This ensures the filter system handles realistic use cases.
func TestFilterCriteria_ComplexScenarios(t *testing.T) {
	tests := []struct {
		name     string                            // Test scenario description
		setup    func() *FilterCriteria            // Setup function
		testFunc func(*testing.T, *FilterCriteria) // Test function
		reason   string                            // Why this scenario is important
	}{
		{
			name: "strict service filtering with method patterns",
			setup: func() *FilterCriteria {
				criteria, _ := ParseFromConfig(
					"UserService",       // Only UserService
					"Get*,Update*",      // Only Get and Update methods
					"*Internal",         // Exclude internal methods
					"GetUser:fetchUser", // Rename GetUser
				)
				return criteria
			},
			testFunc: func(t *testing.T, criteria *FilterCriteria) {
				// Should have strict service filtering
				if !criteria.HasServiceFilter() {
					t.Error("Should have service filter configured")
				}

				// Should have both include and exclude patterns
				if !criteria.HasMethodIncludes() || !criteria.HasMethodExcludes() {
					t.Error("Should have both method includes and excludes")
				}

				// Should have method renames
				if criteria.GetMethodRename("GetUser") != "fetchUser" {
					t.Error("GetUser should be renamed to fetchUser")
				}
			},
			reason: "Real applications often combine service filtering with method patterns",
		},
		{
			name: "open service filtering with strict method control",
			setup: func() *FilterCriteria {
				criteria, _ := ParseFromConfig(
					"",                       // All services
					"",                       // No specific includes
					"*Debug,*Test,*Internal", // Exclude development methods
					"",                       // No renames
				)
				return criteria
			},
			testFunc: func(t *testing.T, criteria *FilterCriteria) {
				// Should not have service filtering (include all services)
				if criteria.HasServiceFilter() {
					t.Error("Should not have service filter (include all)")
				}

				// Should have excludes but not includes
				if criteria.HasMethodIncludes() {
					t.Error("Should not have method includes (include all by default)")
				}

				if !criteria.HasMethodExcludes() {
					t.Error("Should have method excludes for development methods")
				}
			},
			reason: "Production deployments often exclude debug/test methods from all services",
		},
		{
			name: "development configuration",
			setup: func() *FilterCriteria {
				criteria := NewFilterCriteria()
				// Development: include everything, no filtering
				return criteria
			},
			testFunc: func(t *testing.T, criteria *FilterCriteria) {
				// Should have no filtering configured
				if criteria.HasServiceFilter() || criteria.HasMethodIncludes() || criteria.HasMethodExcludes() {
					t.Error("Development config should have no filtering")
				}

				// Should still exclude annotation packages and map entries
				if !criteria.ExcludeAnnotationPackages || !criteria.ExcludeMapEntries {
					t.Error("Should always exclude annotation packages and map entries")
				}
			},
			reason: "Development environments need all functionality available for testing",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			criteria := tt.setup()
			tt.testFunc(t, criteria)
		})
	}
}

// TestFilterResult_Types tests the various filter result types and their usage.
// This ensures result types provide the right metadata for different filtering operations.
func TestFilterResult_Types(t *testing.T) {
	// Test basic FilterResult
	basicIncluded := Included("basic inclusion test")
	if !basicIncluded.Include || basicIncluded.Reason == "" {
		t.Error("Basic included result should have Include=true and non-empty reason")
	}

	basicExcluded := Excluded("basic exclusion test")
	if basicExcluded.Include || basicExcluded.Reason == "" {
		t.Error("Basic excluded result should have Include=false and non-empty reason")
	}

	// Test ServiceFilterResult
	serviceResult := ServiceFilterResult{
		FilterResult:      Included("service test"),
		IsBrowserProvided: true,
		CustomName:        "customService",
	}

	if !serviceResult.Include {
		t.Error("ServiceFilterResult should inherit Include from FilterResult")
	}

	if !serviceResult.IsBrowserProvided {
		t.Error("ServiceFilterResult should preserve IsBrowserProvided metadata")
	}

	if serviceResult.CustomName != "customService" {
		t.Error("ServiceFilterResult should preserve CustomName metadata")
	}

	// Test MethodFilterResult
	methodResult := MethodFilterResult{
		FilterResult:      Included("method test"),
		CustomJSName:      "customMethod",
		IsAsync:           true,
		IsServerStreaming: false,
	}

	if !methodResult.Include {
		t.Error("MethodFilterResult should inherit Include from FilterResult")
	}

	if methodResult.CustomJSName != "customMethod" {
		t.Error("MethodFilterResult should preserve CustomJSName metadata")
	}

	if !methodResult.IsAsync {
		t.Error("MethodFilterResult should preserve IsAsync metadata")
	}

	// Test PackageFilterResult
	packageResult := PackageFilterResult{
		FilterResult: Included("package test"),
		HasServices:  true,
		HasMessages:  true,
		HasEnums:     false,
	}

	if !packageResult.Include {
		t.Error("PackageFilterResult should inherit Include from FilterResult")
	}

	if !packageResult.HasServices || !packageResult.HasMessages {
		t.Error("PackageFilterResult should preserve content metadata")
	}

	if packageResult.HasEnums {
		t.Error("PackageFilterResult should accurately reflect HasEnums=false")
	}
}
