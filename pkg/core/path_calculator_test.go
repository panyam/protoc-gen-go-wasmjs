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
	"path/filepath"
	"runtime"
	"testing"
)

// TestPathCalculator_CalculateRelativePath tests relative path calculation between directories.
// This is critical for generating correct TypeScript import statements that work across
// different package structures and output directory configurations.
func TestPathCalculator_CalculateRelativePath(t *testing.T) {
	calculator := NewPathCalculator()

	tests := []struct {
		name         string // Test case description
		fromPath     string // Source directory path
		toPath       string // Target directory path
		expectedPath string // Expected relative path
		reason       string // Why this test case is important
	}{
		{
			name:         "sibling directories",
			fromPath:     "./gen/wasm",
			toPath:       "./gen/ts",
			expectedPath: "../ts",
			reason:       "Most common case - WASM and TS output in sibling directories",
		},
		{
			name:         "nested to parent",
			fromPath:     "./gen/ts/library/v1",
			toPath:       "./gen/ts",
			expectedPath: "../..",
			reason:       "TypeScript files importing from parent package directories",
		},
		{
			name:         "parent to nested",
			fromPath:     "./gen/ts",
			toPath:       "./gen/ts/library/v1",
			expectedPath: "./library/v1",
			reason:       "Parent directories importing nested package files",
		},
		{
			name:         "cross-package imports",
			fromPath:     "./gen/ts/library/v1",
			toPath:       "./gen/ts/common/v1",
			expectedPath: "../../common/v1",
			reason:       "Cross-package dependencies require correct relative paths",
		},
		{
			name:         "same directory",
			fromPath:     "./gen/ts/library/v1",
			toPath:       "./gen/ts/library/v1",
			expectedPath: ".",
			reason:       "Imports within the same directory should use current directory",
		},
		{
			name:         "deeply nested cross-reference",
			fromPath:     "./gen/ts/company/product/module/v2",
			toPath:       "./gen/ts/common/types/v1",
			expectedPath: "../../../../common/types/v1",
			reason:       "Complex package hierarchies need correct path traversal",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculator.CalculateRelativePath(tt.fromPath, tt.toPath)

			// Normalize paths for comparison (handle Windows vs Unix)
			expected := filepath.ToSlash(tt.expectedPath)
			result = filepath.ToSlash(result)

			if result != expected {
				t.Errorf("CalculateRelativePath(%s, %s) = %s, want %s\nReason: %s",
					tt.fromPath, tt.toPath, result, expected, tt.reason)
			}
		})
	}
}

// TestPathCalculator_BuildPackagePath tests conversion of proto package names to directory paths.
// This ensures TypeScript files are organized in the correct directory structure
// matching the proto package hierarchy.
func TestPathCalculator_BuildPackagePath(t *testing.T) {
	calculator := NewPathCalculator()

	tests := []struct {
		name         string // Test case description
		packageName  string // Input proto package name
		expectedPath string // Expected directory path
		reason       string // Why this test case is important
	}{
		{
			name:         "standard versioned package",
			packageName:  "library.v1",
			expectedPath: "library/v1",
			reason:       "Most common pattern - package with version creates nested dirs",
		},
		{
			name:         "deeply nested package",
			packageName:  "company.product.module.v2",
			expectedPath: "company/product/module/v2",
			reason:       "Complex package hierarchies should create proper nesting",
		},
		{
			name:         "simple package",
			packageName:  "common",
			expectedPath: "common",
			reason:       "Single-level packages should create single directory",
		},
		{
			name:         "empty package",
			packageName:  "",
			expectedPath: "",
			reason:       "Empty packages should return empty path (edge case)",
		},
		{
			name:         "single dot",
			packageName:  ".",
			expectedPath: "/",
			reason:       "Malformed package names should be handled gracefully",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculator.BuildPackagePath(tt.packageName)
			if result != tt.expectedPath {
				t.Errorf("BuildPackagePath(%s) = %s, want %s\nReason: %s",
					tt.packageName, result, tt.expectedPath, tt.reason)
			}
		})
	}
}

// TestPathCalculator_BuildCrossPackageImportPath tests import path calculation between packages.
// This is essential for TypeScript files to correctly import types and functions from other packages.
func TestPathCalculator_BuildCrossPackageImportPath(t *testing.T) {
	calculator := NewPathCalculator()

	tests := []struct {
		name           string // Test case description
		currentPackage string // Package doing the importing
		targetPackage  string // Package being imported
		expectedPath   string // Expected import path
		reason         string // Why this test case is important
	}{
		{
			name:           "same package import",
			currentPackage: "library.v1",
			targetPackage:  "library.v1",
			expectedPath:   ".",
			reason:         "Imports within same package should use current directory",
		},
		{
			name:           "sibling package import",
			currentPackage: "library.v1",
			targetPackage:  "common.v1",
			expectedPath:   "../../common_v1",
			reason:         "Sibling packages need to go up one level then down",
		},
		{
			name:           "nested to parent import",
			currentPackage: "company.product.v1",
			targetPackage:  "company.v1",
			expectedPath:   "../../../company_v1",
			reason:         "Nested packages importing from higher levels",
		},
		{
			name:           "parent to nested import",
			currentPackage: "company.v1",
			targetPackage:  "company.product.v1",
			expectedPath:   "../../company_product_v1",
			reason:         "Parent packages importing from nested packages",
		},
		{
			name:           "cross-tree import",
			currentPackage: "library.services.v1",
			targetPackage:  "user.auth.v2",
			expectedPath:   "../../../user_auth_v2",
			reason:         "Complex cross-package dependencies in different trees",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculator.BuildCrossPackageImportPath(tt.currentPackage, tt.targetPackage)
			if result != tt.expectedPath {
				t.Errorf("BuildCrossPackageImportPath(%s, %s) = %s, want %s\nReason: %s",
					tt.currentPackage, tt.targetPackage, result, tt.expectedPath, tt.reason)
			}
		})
	}
}

// TestPathCalculator_GetFactoryImportPath tests factory import path generation.
// Factory imports are used for cross-package object creation in TypeScript,
// so correct paths are essential for compilation and runtime functionality.
func TestPathCalculator_GetFactoryImportPath(t *testing.T) {
	calculator := NewPathCalculator()

	tests := []struct {
		name              string // Test case description
		dependencyPackage string // Package containing the factory
		currentPackage    string // Package importing the factory
		expectedPath      string // Expected import path
		reason            string // Why this test case is important
	}{
		{
			name:              "same package factory",
			dependencyPackage: "library.v1",
			currentPackage:    "library.v1",
			expectedPath:      "./factory",
			reason:            "Same-package factory imports should use relative path",
		},
		{
			name:              "cross-package factory",
			dependencyPackage: "common.v1",
			currentPackage:    "library.v1",
			expectedPath:      "../../common_v1/factory",
			reason:            "Cross-package factories need correct relative path to factory file",
		},
		{
			name:              "nested dependency factory",
			dependencyPackage: "company.product.types.v1",
			currentPackage:    "company.services.v1",
			expectedPath:      "../../../company_product_types_v1/factory",
			reason:            "Complex nested factory dependencies must resolve correctly",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculator.GetFactoryImportPath(tt.dependencyPackage, tt.currentPackage)
			if result != tt.expectedPath {
				t.Errorf("GetFactoryImportPath(%s, %s) = %s, want %s\nReason: %s",
					tt.dependencyPackage, tt.currentPackage, result, tt.expectedPath, tt.reason)
			}
		})
	}
}

// TestPathCalculator_GenerateOutputFilePath tests complete file path generation.
// This ensures generated files are written to the correct locations in the output structure.
func TestPathCalculator_GenerateOutputFilePath(t *testing.T) {
	calculator := NewPathCalculator()

	tests := []struct {
		name           string // Test case description
		baseOutputPath string // Base output directory
		packageName    string // Proto package name
		fileName       string // Generated file name
		expectedPath   string // Expected complete file path
		reason         string // Why this test case is important
	}{
		{
			name:           "standard TypeScript interface file",
			baseOutputPath: "./gen/ts",
			packageName:    "library.v1",
			fileName:       "interfaces.ts",
			expectedPath:   "gen/ts/library/v1/interfaces.ts",
			reason:         "Most common case - TypeScript files in package directories",
		},
		{
			name:           "nested package structure",
			baseOutputPath: "./gen/typescript",
			packageName:    "company.product.types.v2",
			fileName:       "models.ts",
			expectedPath:   "gen/typescript/company/product/types/v2/models.ts",
			reason:         "Deep package hierarchies should create proper nested structure",
		},
		{
			name:           "root level output",
			baseOutputPath: ".",
			packageName:    "common",
			fileName:       "factory.ts",
			expectedPath:   "common/factory.ts",
			reason:         "Root-level output should work without extra nesting",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculator.GenerateOutputFilePath(tt.baseOutputPath, tt.packageName, tt.fileName)

			// Normalize for cross-platform comparison
			expected := filepath.ToSlash(tt.expectedPath)
			result = filepath.ToSlash(result)

			if result != expected {
				t.Errorf("GenerateOutputFilePath(%s, %s, %s) = %s, want %s\nReason: %s",
					tt.baseOutputPath, tt.packageName, tt.fileName, result, expected, tt.reason)
			}
		})
	}
}

// TestPathCalculator_GetGoPackageAlias tests Go package alias generation.
// Package aliases must be valid Go identifiers and unique to avoid import conflicts.
func TestPathCalculator_GetGoPackageAlias(t *testing.T) {
	calculator := NewPathCalculator()

	tests := []struct {
		name          string // Test case description
		packagePath   string // Go package import path
		expectedAlias string // Expected package alias
		reason        string // Why this test case is important
	}{
		{
			name:          "standard versioned package",
			packagePath:   "github.com/example/proto/gen/go/library/v1",
			expectedAlias: "libraryv1",
			reason:        "Standard case - combine last two path segments",
		},
		{
			name:          "complex nested path",
			packagePath:   "github.com/company/protos/gen/go/user/auth/service/v2",
			expectedAlias: "servicev2",
			reason:        "Deep paths should use last two meaningful segments",
		},
		{
			name:          "single segment path",
			packagePath:   "common",
			expectedAlias: "common",
			reason:        "Short paths should use the available segment",
		},
		{
			name:          "path with dots and special chars",
			packagePath:   "github.com/example/user-service.v1",
			expectedAlias: "exampleuserservicev1",
			reason:        "Should take last two path segments and clean them",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculator.GetGoPackageAlias(tt.packagePath)
			if result != tt.expectedAlias {
				t.Errorf("GetGoPackageAlias(%s) = %s, want %s\nReason: %s",
					tt.packagePath, result, tt.expectedAlias, tt.reason)
			}
		})
	}
}

// TestPathCalculator_NormalizePath tests path normalization across platforms.
// This ensures paths work consistently on Windows, macOS, and Linux systems.
func TestPathCalculator_NormalizePath(t *testing.T) {
	calculator := NewPathCalculator()

	tests := []struct {
		name         string // Test case description
		inputPath    string // Input path (possibly with mixed separators)
		expectedPath string // Expected normalized path
		reason       string // Why this test case is important
	}{
		{
			name:         "mixed separators",
			inputPath:    "./gen\\ts/../wasm",
			expectedPath: "./wasm",
			reason:       "Mixed path separators should be normalized consistently",
		},
		{
			name:         "redundant current directory",
			inputPath:    "./gen/./ts/./interfaces.ts",
			expectedPath: "./gen/ts/interfaces.ts",
			reason:       "Redundant ./ components should be cleaned up",
		},
		{
			name:         "parent directory resolution",
			inputPath:    "./gen/ts/../js/../wasm/output.go",
			expectedPath: "./gen/wasm/output.go",
			reason:       "Parent directory traversal should be resolved correctly",
		},
		{
			name:         "already normalized",
			inputPath:    "./gen/ts/library/v1/interfaces.ts",
			expectedPath: "./gen/ts/library/v1/interfaces.ts",
			reason:       "Already normalized paths should remain unchanged",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculator.NormalizePath(tt.inputPath)
			if result != tt.expectedPath {
				t.Errorf("NormalizePath(%s) = %s, want %s\nReason: %s",
					tt.inputPath, result, tt.expectedPath, tt.reason)
			}
		})
	}
}

// TestPathCalculator_IsAbsolutePath tests absolute vs relative path detection.
// This is important for determining whether paths need to be made relative for imports.
func TestPathCalculator_IsAbsolutePath(t *testing.T) {
	calculator := NewPathCalculator()

	tests := []struct {
		name       string // Test case description
		path       string // Input path to test
		isAbsolute bool   // Expected result
		reason     string // Why this test case is important
	}{
		{
			name:       "relative current directory",
			path:       "./gen/ts",
			isAbsolute: false,
			reason:     "Relative paths starting with ./ should be identified correctly",
		},
		{
			name:       "relative parent directory",
			path:       "../gen/wasm",
			isAbsolute: false,
			reason:     "Relative paths with parent traversal should be identified",
		},
		{
			name:       "relative without prefix",
			path:       "gen/output",
			isAbsolute: false,
			reason:     "Relative paths without ./ prefix should be identified",
		},
	}

	// Add platform-specific absolute path tests
	if runtime.GOOS == "windows" {
		tests = append(tests, []struct {
			name       string
			path       string
			isAbsolute bool
			reason     string
		}{
			{
				name:       "windows absolute path",
				path:       "C:\\Users\\developer\\project",
				isAbsolute: true,
				reason:     "Windows absolute paths should be identified on Windows",
			},
			{
				name:       "windows UNC path",
				path:       "\\\\server\\share\\file",
				isAbsolute: true,
				reason:     "Windows UNC paths should be identified as absolute",
			},
		}...)
	} else {
		tests = append(tests, struct {
			name       string
			path       string
			isAbsolute bool
			reason     string
		}{
			name:       "unix absolute path",
			path:       "/usr/local/bin",
			isAbsolute: true,
			reason:     "Unix absolute paths should be identified on Unix systems",
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculator.IsAbsolutePath(tt.path)
			if result != tt.isAbsolute {
				t.Errorf("IsAbsolutePath(%s) = %t, want %t\nReason: %s",
					tt.path, result, tt.isAbsolute, tt.reason)
			}
		})
	}
}

// TestPathCalculator_JoinPaths tests path component joining with normalization.
// This ensures consistent path construction regardless of input format.
func TestPathCalculator_JoinPaths(t *testing.T) {
	calculator := NewPathCalculator()

	tests := []struct {
		name         string   // Test case description
		components   []string // Path components to join
		expectedPath string   // Expected joined and normalized path
		reason       string   // Why this test case is important
	}{
		{
			name:         "standard path join",
			components:   []string{"./gen", "ts", "library", "v1", "interfaces.ts"},
			expectedPath: "./gen/ts/library/v1/interfaces.ts",
			reason:       "Standard path construction should work consistently",
		},
		{
			name:         "empty components handling",
			components:   []string{"./gen", "", "ts", "output.js"},
			expectedPath: "./gen/ts/output.js",
			reason:       "Empty components should be handled gracefully",
		},
		{
			name:         "single component",
			components:   []string{"interfaces.ts"},
			expectedPath: "interfaces.ts",
			reason:       "Single components without ./ prefix shouldn't get one added",
		},
		{
			name:         "no components",
			components:   []string{},
			expectedPath: "",
			reason:       "Empty input should return empty string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculator.JoinPaths(tt.components...)
			if result != tt.expectedPath {
				t.Errorf("JoinPaths(%v) = %s, want %s\nReason: %s",
					tt.components, result, tt.expectedPath, tt.reason)
			}
		})
	}
}
