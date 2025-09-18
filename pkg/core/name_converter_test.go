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

package core

import (
	"testing"
)

// TestNameConverter_ToCamelCase tests conversion from PascalCase to camelCase.
// This is critical for generating JavaScript method names from Go method names,
// ensuring the generated API follows JavaScript naming conventions.
func TestNameConverter_ToCamelCase(t *testing.T) {
	converter := NewNameConverter()

	tests := []struct {
		name     string // Test case description
		input    string // Input PascalCase string
		expected string // Expected camelCase output
		reason   string // Why this test case is important
	}{
		{
			name:     "standard method name",
			input:    "FindBooks",
			expected: "findBooks",
			reason:   "Most common case - Go method names converted to JS method names",
		},
		{
			name:     "single word",
			input:    "Login",
			expected: "login",
			reason:   "Simple method names should work correctly",
		},
		{
			name:     "multi-word method",
			input:    "CreateLibraryItem",
			expected: "createLibraryItem",
			reason:   "Complex method names preserve internal capitalization",
		},
		{
			name:     "already camelCase",
			input:    "getUserData",
			expected: "getUserData",
			reason:   "Already camelCase strings should remain unchanged",
		},
		{
			name:     "single letter",
			input:    "A",
			expected: "a",
			reason:   "Edge case - single letter should be lowercase",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
			reason:   "Empty strings should be handled gracefully",
		},
		{
			name:     "all uppercase abbreviation",
			input:    "HTTPSConnection",
			expected: "hTTPSConnection",
			reason:   "Abbreviations should only lowercase the first letter",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := converter.ToCamelCase(tt.input)
			if result != tt.expected {
				t.Errorf("ToCamelCase(%s) = %s, want %s\nReason: %s",
					tt.input, result, tt.expected, tt.reason)
			}
		})
	}
}

// TestNameConverter_ToPascalCase tests conversion from camelCase to PascalCase.
// This is used for generating TypeScript class names and interface names from
// variable names or other camelCase identifiers.
func TestNameConverter_ToPascalCase(t *testing.T) {
	converter := NewNameConverter()

	tests := []struct {
		name     string // Test case description
		input    string // Input camelCase string
		expected string // Expected PascalCase output
		reason   string // Why this test case is important
	}{
		{
			name:     "standard camelCase",
			input:    "findBooks",
			expected: "FindBooks",
			reason:   "Common case - JS method names to TS interface method names",
		},
		{
			name:     "service name",
			input:    "libraryService",
			expected: "LibraryService",
			reason:   "Service names need PascalCase for TypeScript classes",
		},
		{
			name:     "single word",
			input:    "user",
			expected: "User",
			reason:   "Simple names should capitalize first letter",
		},
		{
			name:     "already PascalCase",
			input:    "CreateUser",
			expected: "CreateUser",
			reason:   "Already PascalCase strings should remain unchanged",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
			reason:   "Empty strings should be handled gracefully",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := converter.ToPascalCase(tt.input)
			if result != tt.expected {
				t.Errorf("ToPascalCase(%s) = %s, want %s\nReason: %s",
					tt.input, result, tt.expected, tt.reason)
			}
		})
	}
}

// TestNameConverter_ToSnakeCase tests conversion to snake_case.
// This is used for generating certain file names or when interfacing with
// systems that expect snake_case naming conventions.
func TestNameConverter_ToSnakeCase(t *testing.T) {
	converter := NewNameConverter()

	tests := []struct {
		name     string // Test case description
		input    string // Input string (camelCase or PascalCase)
		expected string // Expected snake_case output
		reason   string // Why this test case is important
	}{
		{
			name:     "PascalCase method",
			input:    "FindBooks",
			expected: "find_books",
			reason:   "Method names may need snake_case for certain contexts",
		},
		{
			name:     "camelCase property",
			input:    "getUserData",
			expected: "get_user_data",
			reason:   "Property names may need snake_case conversion",
		},
		{
			name:     "single word",
			input:    "User",
			expected: "user",
			reason:   "Single words should just be lowercased",
		},
		{
			name:     "consecutive capitals",
			input:    "HTTPSConnection",
			expected: "h_t_t_p_s_connection",
			reason:   "Consecutive capitals should each get underscores",
		},
		{
			name:     "already lowercase",
			input:    "simple",
			expected: "simple",
			reason:   "Already lowercase strings should remain unchanged",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
			reason:   "Empty strings should be handled gracefully",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := converter.ToSnakeCase(tt.input)
			if result != tt.expected {
				t.Errorf("ToSnakeCase(%s) = %s, want %s\nReason: %s",
					tt.input, result, tt.expected, tt.reason)
			}
		})
	}
}

// TestNameConverter_ToPackageAlias tests Go package alias generation.
// Package aliases must be valid Go identifiers and should be short but descriptive
// to avoid conflicts in import statements.
func TestNameConverter_ToPackageAlias(t *testing.T) {
	converter := NewNameConverter()

	tests := []struct {
		name          string // Test case description
		packagePath   string // Input package path
		expectedAlias string // Expected package alias
		reason        string // Why this test case is important
	}{
		{
			name:          "standard versioned package",
			packagePath:   "github.com/example/proto/gen/go/library/v1",
			expectedAlias: "libraryv1",
			reason:        "Most common case - versioned packages need readable aliases",
		},
		{
			name:          "hyphenated package name",
			packagePath:   "github.com/example/user-service/v2",
			expectedAlias: "userservicev2",
			reason:        "Hyphens should be removed to create valid Go identifiers",
		},
		{
			name:          "dotted package name",
			packagePath:   "common.types.v1",
			expectedAlias: "commontypesv1",
			reason:        "Should handle dotted package names as single segment",
		},
		{
			name:          "single segment",
			packagePath:   "utils",
			expectedAlias: "utils",
			reason:        "Simple package names should be used as-is",
		},
		{
			name:          "empty package path",
			packagePath:   "",
			expectedAlias: "pkg",
			reason:        "Empty paths should have a fallback alias",
		},
		{
			name:          "complex nested path",
			packagePath:   "github.com/company/services/user-auth/types/v3",
			expectedAlias: "typesv3",
			reason:        "Complex paths should use the most specific parts",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := converter.ToPackageAlias(tt.packagePath)
			if result != tt.expectedAlias {
				t.Errorf("ToPackageAlias(%s) = %s, want %s\nReason: %s",
					tt.packagePath, result, tt.expectedAlias, tt.reason)
			}
		})
	}
}

// TestNameConverter_ToJSNamespace tests JavaScript namespace generation.
// Namespaces are used for global objects in the browser environment and must
// follow JavaScript identifier rules while being descriptive.
func TestNameConverter_ToJSNamespace(t *testing.T) {
	converter := NewNameConverter()

	tests := []struct {
		name        string // Test case description
		packageName string // Input package name
		expectedNS  string // Expected JavaScript namespace
		reason      string // Why this test case is important
	}{
		{
			name:        "standard versioned package",
			packageName: "library.v1",
			expectedNS:  "library_v1",
			reason:      "Dots should be converted to underscores for JS identifiers",
		},
		{
			name:        "hyphenated package",
			packageName: "user-service.v2",
			expectedNS:  "user_service_v2",
			reason:      "Hyphens should be converted to underscores",
		},
		{
			name:        "mixed case package",
			packageName: "Library.Types",
			expectedNS:  "library_types",
			reason:      "Should be lowercase for consistent JavaScript style",
		},
		{
			name:        "simple package",
			packageName: "common",
			expectedNS:  "common",
			reason:      "Simple names should work as-is",
		},
		{
			name:        "empty package",
			packageName: "",
			expectedNS:  "",
			reason:      "Empty input should return empty string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := converter.ToJSNamespace(tt.packageName)
			if result != tt.expectedNS {
				t.Errorf("ToJSNamespace(%s) = %s, want %s\nReason: %s",
					tt.packageName, result, tt.expectedNS, tt.reason)
			}
		})
	}
}

// TestNameConverter_ToModuleName tests WASM module name generation.
// Module names are used for WASM modules and Go packages, so they need to be
// descriptive and follow Go package naming conventions.
func TestNameConverter_ToModuleName(t *testing.T) {
	converter := NewNameConverter()

	tests := []struct {
		name         string // Test case description
		packageName  string // Input package name
		expectedName string // Expected module name
		reason       string // Why this test case is important
	}{
		{
			name:         "standard package",
			packageName:  "library.v1",
			expectedName: "library_v1_services",
			reason:       "Standard packages get descriptive module names",
		},
		{
			name:         "nested package",
			packageName:  "user.auth.v2",
			expectedName: "user_auth_v2_services",
			reason:       "Nested packages should flatten with underscores",
		},
		{
			name:         "simple package",
			packageName:  "common",
			expectedName: "common_services",
			reason:       "Simple packages get the services suffix",
		},
		{
			name:         "empty package",
			packageName:  "",
			expectedName: "services",
			reason:       "Empty package should fallback to just 'services'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := converter.ToModuleName(tt.packageName)
			if result != tt.expectedName {
				t.Errorf("ToModuleName(%s) = %s, want %s\nReason: %s",
					tt.packageName, result, tt.expectedName, tt.reason)
			}
		})
	}
}

// TestNameConverter_ToFactoryName tests TypeScript factory class name generation.
// Factory names must be valid TypeScript class names (PascalCase) and should
// be descriptive enough to identify the package they serve.
func TestNameConverter_ToFactoryName(t *testing.T) {
	converter := NewNameConverter()

	tests := []struct {
		name         string // Test case description
		packageName  string // Input package name
		expectedName string // Expected factory class name
		reason       string // Why this test case is important
	}{
		{
			name:         "standard versioned package",
			packageName:  "library.v1",
			expectedName: "LibraryV1Factory",
			reason:       "Standard case - PascalCase with version and Factory suffix",
		},
		{
			name:         "multi-part package",
			packageName:  "common.types",
			expectedName: "CommonTypesFactory",
			reason:       "Multi-part packages should combine all parts",
		},
		{
			name:         "hyphenated package",
			packageName:  "user-auth",
			expectedName: "UserAuthFactory",
			reason:       "Hyphens should be converted to PascalCase parts",
		},
		{
			name:         "underscore package",
			packageName:  "data_store.v1",
			expectedName: "DataStoreV1Factory",
			reason:       "Underscores should be converted to PascalCase parts",
		},
		{
			name:         "empty package",
			packageName:  "",
			expectedName: "Factory",
			reason:       "Empty package should fallback to just 'Factory'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := converter.ToFactoryName(tt.packageName)
			if result != tt.expectedName {
				t.Errorf("ToFactoryName(%s) = %s, want %s\nReason: %s",
					tt.packageName, result, tt.expectedName, tt.reason)
			}
		})
	}
}

// TestNameConverter_ToGoFuncName tests Go WASM function name generation.
// These function names are exposed to JavaScript and must be unique across
// all services while being readable and following Go conventions.
func TestNameConverter_ToGoFuncName(t *testing.T) {
	converter := NewNameConverter()

	tests := []struct {
		name         string // Test case description
		serviceName  string // Service name
		methodName   string // Method name
		expectedName string // Expected Go function name
		reason       string // Why this test case is important
	}{
		{
			name:         "standard service method",
			serviceName:  "LibraryService",
			methodName:   "FindBooks",
			expectedName: "libraryServiceFindBooks",
			reason:       "Standard case - camelCase service + PascalCase method",
		},
		{
			name:         "simple service method",
			serviceName:  "UserAuth",
			methodName:   "Login",
			expectedName: "userAuthLogin",
			reason:       "Shorter names should work consistently",
		},
		{
			name:         "complex method name",
			serviceName:  "DataProcessor",
			methodName:   "ProcessLargeDataSet",
			expectedName: "dataProcessorProcessLargeDataSet",
			reason:       "Complex method names should be preserved",
		},
		{
			name:         "single word service",
			serviceName:  "Auth",
			methodName:   "Validate",
			expectedName: "authValidate",
			reason:       "Single word services should work correctly",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := converter.ToGoFuncName(tt.serviceName, tt.methodName)
			if result != tt.expectedName {
				t.Errorf("ToGoFuncName(%s, %s) = %s, want %s\nReason: %s",
					tt.serviceName, tt.methodName, result, tt.expectedName, tt.reason)
			}
		})
	}
}

// TestNameConverter_SanitizeIdentifier tests identifier sanitization for target languages.
// This ensures generated names are always valid identifiers regardless of proto content.
func TestNameConverter_SanitizeIdentifier(t *testing.T) {
	converter := NewNameConverter()

	tests := []struct {
		name         string // Test case description
		input        string // Input string (potentially invalid identifier)
		expectedName string // Expected sanitized identifier
		reason       string // Why this test case is important
	}{
		{
			name:         "valid identifier",
			input:        "userName",
			expectedName: "userName",
			reason:       "Valid identifiers should pass through unchanged",
		},
		{
			name:         "hyphenated name",
			input:        "user-name",
			expectedName: "user_name",
			reason:       "Hyphens should be converted to underscores",
		},
		{
			name:         "starts with number",
			input:        "123invalid",
			expectedName: "_23invalid",
			reason:       "Identifiers starting with numbers should be prefixed",
		},
		{
			name:         "dotted name",
			input:        "user.name",
			expectedName: "user_name",
			reason:       "Dots should be converted to underscores",
		},
		{
			name:         "mixed invalid characters",
			input:        "user@name#123",
			expectedName: "user_name_123",
			reason:       "Invalid characters should be converted to underscores",
		},
		{
			name:         "empty string",
			input:        "",
			expectedName: "identifier",
			reason:       "Empty input should get a fallback identifier",
		},
		{
			name:         "all invalid characters",
			input:        "@#$%",
			expectedName: "____",
			reason:       "All invalid characters should be handled",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := converter.SanitizeIdentifier(tt.input)
			if result != tt.expectedName {
				t.Errorf("SanitizeIdentifier(%s) = %s, want %s\nReason: %s",
					tt.input, result, tt.expectedName, tt.reason)
			}
		})
	}
}
