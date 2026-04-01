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

	"github.com/panyam/protokit/packages"
)

// PathCalculator provides functions for calculating file paths and import paths.
// Generic methods delegate to protokit/packages. TS-specific methods are defined here.
type PathCalculator struct{}

// NewPathCalculator creates a new PathCalculator instance.
func NewPathCalculator() *PathCalculator {
	return &PathCalculator{}
}

// --- Generic methods (delegated to protokit) ---

// CalculateRelativePath calculates the relative path from one directory to another.
// Delegates to protokit/packages.CalculateRelativePath.
func (pc *PathCalculator) CalculateRelativePath(fromPath, toPath string) string {
	return packages.CalculateRelativePath(fromPath, toPath)
}

// BuildPackagePath converts a dot-separated package name to a directory path.
// Delegates to protokit/packages.BuildPackagePath.
func (pc *PathCalculator) BuildPackagePath(packageName string) string {
	return packages.BuildPackagePath(packageName)
}

// NormalizePath normalizes a file path for consistent processing.
// Delegates to protokit/packages.NormalizePath.
func (pc *PathCalculator) NormalizePath(path string) string {
	return packages.NormalizePath(path)
}

// IsAbsolutePath checks if a path is absolute.
// Delegates to protokit/packages.IsAbsolutePath.
func (pc *PathCalculator) IsAbsolutePath(path string) bool {
	return packages.IsAbsolutePath(path)
}

// JoinPaths joins path components with normalization.
// Delegates to protokit/packages.JoinPaths.
func (pc *PathCalculator) JoinPaths(components ...string) string {
	return packages.JoinPaths(components...)
}

// --- TS/WASM-specific methods (stay in this package) ---

// BuildCrossPackageImportPath calculates the relative import path between two packages.
func (pc *PathCalculator) BuildCrossPackageImportPath(currentPackage, targetPackage string) string {
	if currentPackage == targetPackage {
		return "."
	}
	currentPath := packages.BuildPackagePath(currentPackage)
	currentLevels := strings.Count(currentPath, "/") + 1
	upLevels := strings.Repeat("../", currentLevels)
	targetImportName := strings.ReplaceAll(targetPackage, ".", "_")
	return upLevels + targetImportName
}

// GetFactoryImportPath generates the relative import path for a factory dependency.
func (pc *PathCalculator) GetFactoryImportPath(dependencyPackage, currentPackage string) string {
	if dependencyPackage == currentPackage {
		return "./factory"
	}
	return pc.BuildCrossPackageImportPath(currentPackage, dependencyPackage) + "/factory"
}

// GenerateOutputFilePath creates the complete file path for a generated file.
func (pc *PathCalculator) GenerateOutputFilePath(baseOutputPath, packageName, fileName string) string {
	packagePath := packages.BuildPackagePath(packageName)
	return filepath.Join(baseOutputPath, packagePath, fileName)
}

// GetGoPackageAlias creates a Go package alias from an import path.
// Uses last two path segments for unique identifiers.
func (pc *PathCalculator) GetGoPackageAlias(packagePath string) string {
	parts := strings.Split(packagePath, "/")
	if len(parts) >= 2 {
		pkg := parts[len(parts)-2]
		version := parts[len(parts)-1]
		pkg = strings.ReplaceAll(pkg, ".", "")
		pkg = strings.ReplaceAll(pkg, "-", "")
		pkg = strings.ReplaceAll(pkg, "_", "")
		version = strings.ReplaceAll(version, ".", "")
		version = strings.ReplaceAll(version, "-", "")
		version = strings.ReplaceAll(version, "_", "")
		return strings.ToLower(pkg + version)
	}
	last := parts[len(parts)-1]
	last = strings.ReplaceAll(last, ".", "")
	last = strings.ReplaceAll(last, "-", "")
	last = strings.ReplaceAll(last, "_", "")
	return strings.ToLower(last)
}
