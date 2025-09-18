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
	"strings"
)

// PathCalculator provides pure functions for calculating file paths, import paths,
// and directory structures. All functions are stateless and perform path manipulations
// without any external dependencies or side effects.
type PathCalculator struct{}

// NewPathCalculator creates a new PathCalculator instance.
func NewPathCalculator() *PathCalculator {
	return &PathCalculator{}
}

// CalculateRelativePath calculates the relative path from one directory to another.
// Both paths should be relative to a common root (e.g., protoc working directory).
//
// Example:
//
//	fromPath: "./gen/wasm"
//	toPath: "./gen/ts"
//	returns: "../ts"
//
// This is essential for generating correct import statements in TypeScript files.
func (pc *PathCalculator) CalculateRelativePath(fromPath, toPath string) string {
	// Clean the paths to handle . and .. properly
	fromPath = filepath.Clean(fromPath)
	toPath = filepath.Clean(toPath)

	// Handle same directory case
	if fromPath == toPath {
		return "."
	}

	// Calculate relative path
	relPath, err := filepath.Rel(fromPath, toPath)
	if err != nil {
		// Fallback to absolute path if relative calculation fails
		return toPath
	}

	// Convert to forward slashes for TypeScript imports
	relPath = filepath.ToSlash(relPath)

	// Ensure path starts with ./ for relative imports (but not if it's already ../ or .)
	if relPath != "." && !strings.HasPrefix(relPath, "./") && !strings.HasPrefix(relPath, "../") {
		relPath = "./" + relPath
	}

	return relPath
}

// BuildPackagePath converts a proto package name to a directory path structure.
// This creates nested directories from dot-separated package names.
//
// Example:
//
//	packageName: "library.v1"
//	returns: "library/v1"
//
// Used for organizing generated TypeScript files in package-based directory structures.
func (pc *PathCalculator) BuildPackagePath(packageName string) string {
	// Convert "library.v1" to "library/v1" structure
	packagePath := strings.ReplaceAll(packageName, ".", "/")
	return packagePath
}

// BuildCrossPackageImportPath calculates the relative import path between two packages.
// This is used when TypeScript files in one package need to import from another package.
//
// Example:
//
//	currentPackage: "library.v1"
//	targetPackage: "common.v1"
//	returns: "../common_v1"
//
// The calculation assumes both packages are in the same root directory structure.
func (pc *PathCalculator) BuildCrossPackageImportPath(currentPackage, targetPackage string) string {
	if currentPackage == targetPackage {
		return "." // Same package, relative to current directory
	}

	// Convert package names to directory paths
	currentPath := pc.BuildPackagePath(currentPackage)

	// Calculate how many levels up we need to go from current package
	currentLevels := strings.Count(currentPath, "/") + 1
	upLevels := strings.Repeat("../", currentLevels)

	// Convert target package to underscore format for import
	targetImportName := strings.ReplaceAll(targetPackage, ".", "_")

	// Build the import path
	importPath := upLevels + targetImportName

	return importPath
}

// GetFactoryImportPath generates the relative import path for a factory dependency.
// Factories are used for creating instances across package boundaries in TypeScript.
//
// Example:
//
//	dependencyPackage: "common.v1"
//	currentPackage: "library.v1"
//	returns: "../common_v1/factory"
//
// This ensures factory imports use the correct relative paths.
func (pc *PathCalculator) GetFactoryImportPath(dependencyPackage, currentPackage string) string {
	if dependencyPackage == currentPackage {
		return "./factory" // Same package
	}

	crossPackagePath := pc.BuildCrossPackageImportPath(currentPackage, dependencyPackage)
	return crossPackagePath + "/factory"
}

// GenerateOutputFilePath creates the complete file path for a generated file.
// This combines the base output directory, package structure, and filename.
//
// Example:
//
//	baseOutputPath: "./gen/ts"
//	packageName: "library.v1"
//	fileName: "interfaces.ts"
//	returns: "./gen/ts/library/v1/interfaces.ts"
//
// Used by generators to determine where to write generated files.
func (pc *PathCalculator) GenerateOutputFilePath(baseOutputPath, packageName, fileName string) string {
	packagePath := pc.BuildPackagePath(packageName)
	return filepath.Join(baseOutputPath, packagePath, fileName)
}

// GetGoPackageAlias creates a Go package alias from an import path.
// This generates valid Go identifiers for use in import statements.
//
// Example:
//
//	packagePath: "github.com/example/proto/gen/go/library/v1"
//	returns: "libraryv1"
//
// The alias combines the last two path segments to create unique, readable identifiers.
func (pc *PathCalculator) GetGoPackageAlias(packagePath string) string {
	// Extract package name from the path (e.g., "library/v1" -> "libraryv1")
	parts := strings.Split(packagePath, "/")
	if len(parts) >= 2 {
		// Take last two parts and combine them
		pkg := parts[len(parts)-2]
		version := parts[len(parts)-1]

		// Remove special characters to create valid Go identifier
		pkg = strings.ReplaceAll(pkg, ".", "")
		pkg = strings.ReplaceAll(pkg, "-", "")
		pkg = strings.ReplaceAll(pkg, "_", "")

		version = strings.ReplaceAll(version, ".", "")
		version = strings.ReplaceAll(version, "-", "")
		version = strings.ReplaceAll(version, "_", "")

		alias := pkg + version
		return strings.ToLower(alias)
	}
	// Fallback to last part only
	last := parts[len(parts)-1]
	last = strings.ReplaceAll(last, ".", "")
	last = strings.ReplaceAll(last, "-", "")
	last = strings.ReplaceAll(last, "_", "")
	return strings.ToLower(last)
}

// NormalizePath normalizes a file path for consistent processing.
// This handles different path separators and resolves relative path components.
//
// Example:
//
//	path: "./gen\\ts/../wasm"
//	returns: "./gen/wasm" (on all platforms)
//
// Ensures paths are consistent regardless of operating system or input format.
func (pc *PathCalculator) NormalizePath(path string) string {
	// Preserve whether the path started with "./"
	startsWithDot := strings.HasPrefix(path, "./")

	// Clean the path to resolve . and .. components
	cleanPath := filepath.Clean(path)

	// Convert to forward slashes for consistency
	normalizedPath := filepath.ToSlash(cleanPath)

	// If original path started with "./" and clean removed it, add it back
	// But don't add it if the path now starts with ../ or is just "."
	if startsWithDot &&
		!strings.HasPrefix(normalizedPath, "./") &&
		!strings.HasPrefix(normalizedPath, "../") &&
		normalizedPath != "." {
		normalizedPath = "./" + normalizedPath
	}

	return normalizedPath
}

// IsAbsolutePath checks if a path is absolute or relative.
// This is used to determine if paths need to be made relative for imports.
//
// Returns true for absolute paths like "/usr/local" or "C:\Users"
// Returns false for relative paths like "./gen" or "../types"
func (pc *PathCalculator) IsAbsolutePath(path string) bool {
	return filepath.IsAbs(path)
}

// JoinPaths joins multiple path components into a single path.
// This is a wrapper around filepath.Join with additional normalization.
//
// Example:
//
//	components: ["./gen", "ts", "library/v1", "interfaces.ts"]
//	returns: "./gen/ts/library/v1/interfaces.ts"
//
// Ensures all path joins use consistent separators and normalization.
func (pc *PathCalculator) JoinPaths(components ...string) string {
	if len(components) == 0 {
		return ""
	}

	// Track if first component starts with "./"
	startsWithDot := len(components) > 0 && strings.HasPrefix(components[0], "./")

	// Filter out empty components
	var nonEmpty []string
	for _, comp := range components {
		if comp != "" {
			nonEmpty = append(nonEmpty, comp)
		}
	}

	if len(nonEmpty) == 0 {
		return ""
	}

	joined := filepath.Join(nonEmpty...)
	normalized := pc.NormalizePath(joined)

	// If original first component started with "./" and result doesn't, add it back
	if startsWithDot && !strings.HasPrefix(normalized, "./") && !strings.HasPrefix(normalized, "../") {
		normalized = "./" + normalized
	}

	return normalized
}
