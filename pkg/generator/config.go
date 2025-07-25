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

package generator

import (
	"fmt"
	"path/filepath"
	"strings"
)

// Config holds all configuration options for the WASM generator
type Config struct {
	// Core integration
	TSGenerator  string // protoc-gen-es, protoc-gen-ts, etc.
	TSImportPath string // Path where TypeScript types are generated (for imports)
	GenerateWasm       bool   // Generate WASM wrapper (default: true)
	GenerateTypeScript bool   // Generate TypeScript client (default: true)
	WasmExportPath string // Path where WASM wrapper should be generated
	
	// Service & method selection
	Services      string // Comma-separated list of services
	MethodInclude string // Comma-separated glob patterns for methods to include
	MethodExclude string // Comma-separated glob patterns for methods to exclude
	MethodRename  string // Comma-separated method renames (OldName:NewName)
	
	// JS API structure
	JSStructure string // namespaced|flat|service_based
	JSNamespace string // Global JavaScript namespace
	ModuleName  string // WASM module name
	
	// Customization
	TemplateDir  string // Directory containing custom templates
	WasmTemplate string // Custom WASM template file
	TSTemplate   string // Custom TypeScript template file
	
	// Build integration
	WasmPackageSuffix   string // Package suffix for WASM wrapper
	GenerateBuildScript bool   // Generate build script for WASM compilation
	
	// Parsed configuration (populated by Validate)
	ServicesSet      map[string]bool
	MethodIncludes   []string
	MethodExcludes   []string
	MethodRenames    map[string]string
}

// Validate validates and normalizes the configuration
func (c *Config) Validate() error {
	// Validate TypeScript generator
	if c.TSGenerator == "" {
		c.TSGenerator = "protoc-gen-es"
	}
	
	supportedTSGenerators := map[string]bool{
		"protoc-gen-es": true,
		"protoc-gen-ts": true,
	}
	if !supportedTSGenerators[c.TSGenerator] {
		return fmt.Errorf("unsupported ts_generator: %s (supported: protoc-gen-es, protoc-gen-ts)", c.TSGenerator)
	}
	
	// Validate TypeScript import path (where we read existing types)
	if c.TSImportPath == "" {
		c.TSImportPath = "./gen/ts"
	}
	
	// Set defaults for generation flags
	if c.GenerateWasm == false && c.GenerateTypeScript == false {
		// If both are explicitly false, enable both (default behavior)
		c.GenerateWasm = true
		c.GenerateTypeScript = true
	}
	// If neither is explicitly set, default to generating both
	// This will be handled by the flags parsing in main.go
	
	// Validate WASM export path (where we write WASM wrapper)
	if c.WasmExportPath == "" {
		c.WasmExportPath = "." // Default to current directory (protoc out directory)
	}
	
	// Parse services list
	c.ServicesSet = make(map[string]bool)
	if c.Services != "" {
		for _, service := range strings.Split(c.Services, ",") {
			service = strings.TrimSpace(service)
			if service != "" {
				c.ServicesSet[service] = true
			}
		}
	}
	
	// Parse method includes
	if c.MethodInclude != "" {
		for _, pattern := range strings.Split(c.MethodInclude, ",") {
			pattern = strings.TrimSpace(pattern)
			if pattern != "" {
				c.MethodIncludes = append(c.MethodIncludes, pattern)
			}
		}
	}
	
	// Parse method excludes
	if c.MethodExclude != "" {
		for _, pattern := range strings.Split(c.MethodExclude, ",") {
			pattern = strings.TrimSpace(pattern)
			if pattern != "" {
				c.MethodExcludes = append(c.MethodExcludes, pattern)
			}
		}
	}
	
	// Parse method renames
	c.MethodRenames = make(map[string]string)
	if c.MethodRename != "" {
		for _, rename := range strings.Split(c.MethodRename, ",") {
			rename = strings.TrimSpace(rename)
			if rename != "" {
				parts := strings.SplitN(rename, ":", 2)
				if len(parts) != 2 {
					return fmt.Errorf("invalid method rename format: %s (expected OldName:NewName)", rename)
				}
				oldName := strings.TrimSpace(parts[0])
				newName := strings.TrimSpace(parts[1])
				if oldName == "" || newName == "" {
					return fmt.Errorf("invalid method rename: empty old or new name in %s", rename)
				}
				c.MethodRenames[oldName] = newName
			}
		}
	}
	
	// Validate JS structure
	if c.JSStructure == "" {
		c.JSStructure = "namespaced"
	}
	
	validStructures := map[string]bool{
		"namespaced":     true,
		"flat":           true,
		"service_based":  true,
	}
	if !validStructures[c.JSStructure] {
		return fmt.Errorf("invalid js_structure: %s (supported: namespaced, flat, service_based)", c.JSStructure)
	}
	
	// Validate template directory if specified
	if c.TemplateDir != "" {
		absPath, err := filepath.Abs(c.TemplateDir)
		if err != nil {
			return fmt.Errorf("invalid template_dir: %w", err)
		}
		c.TemplateDir = absPath
	}
	
	// Set defaults for WASM package suffix
	if c.WasmPackageSuffix == "" {
		c.WasmPackageSuffix = "wasm"
	}
	
	return nil
}

// ShouldGenerateService returns true if the given service should be generated
func (c *Config) ShouldGenerateService(serviceName string) bool {
	if len(c.ServicesSet) == 0 {
		return true // Generate all services if none specified
	}
	return c.ServicesSet[serviceName]
}

// ShouldGenerateMethod returns true if the given method should be generated
func (c *Config) ShouldGenerateMethod(methodName string) bool {
	// Check excludes first
	for _, pattern := range c.MethodExcludes {
		if matched, _ := filepath.Match(pattern, methodName); matched {
			return false
		}
	}
	
	// If no includes specified, include by default (unless excluded above)
	if len(c.MethodIncludes) == 0 {
		return true
	}
	
	// Check includes
	for _, pattern := range c.MethodIncludes {
		if matched, _ := filepath.Match(pattern, methodName); matched {
			return true
		}
	}
	
	return false
}

// GetMethodJSName returns the JavaScript name for a method, applying renames if configured
func (c *Config) GetMethodJSName(methodName string) string {
	if renamed, exists := c.MethodRenames[methodName]; exists {
		return renamed
	}
	
	// Convert to camelCase for JavaScript
	return toCamelCase(methodName)
}

// GetDefaultJSNamespace returns the default JavaScript namespace based on package name
func (c *Config) GetDefaultJSNamespace(packageName string) string {
	if c.JSNamespace != "" {
		return c.JSNamespace
	}
	
	// Convert package name to lowercase and replace dots with underscores
	namespace := strings.ToLower(packageName)
	namespace = strings.ReplaceAll(namespace, ".", "_")
	return namespace
}

// GetDefaultModuleName returns the default WASM module name
func (c *Config) GetDefaultModuleName(packageName string) string {
	if c.ModuleName != "" {
		return c.ModuleName
	}
	
	// Convert package name and add suffix
	name := strings.ReplaceAll(packageName, ".", "_")
	return name + "_services"
}

// GetTSImportPathForProto returns the TypeScript import path for a given proto file
func (c *Config) GetTSImportPathForProto(protoFile string) string {
	// Remove .proto extension and construct path based on TS generator
	baseName := strings.TrimSuffix(protoFile, ".proto")
	
	switch c.TSGenerator {
	case "protoc-gen-es":
		return c.TSImportPath + "/" + baseName + "_pb.js"
	case "protoc-gen-ts":
		return c.TSImportPath + "/" + baseName + "_pb"
	default:
		return c.TSImportPath + "/" + baseName + "_pb"
	}
}

// GetRelativeTSImportPathForProto returns the relative TypeScript import path for a given proto file
// relative to where the TypeScript client is generated (co-located with WASM artifacts)
func (c *Config) GetRelativeTSImportPathForProto(protoFile string) string {
	// Remove .proto extension and construct path based on TS generator
	baseName := strings.TrimSuffix(protoFile, ".proto")
	
	// Calculate relative path from WasmExportPath (where TS client is generated) to TSImportPath
	relativePath := c.calculateRelativePath(c.WasmExportPath, c.TSImportPath)
	
	// Construct the full relative import path with proper filename
	var filename string
	switch c.TSGenerator {
	case "protoc-gen-es":
		filename = baseName + "_pb.js"
	case "protoc-gen-ts":
		filename = baseName + "_pb"
	default:
		filename = baseName + "_pb"
	}
	
	return relativePath + "/" + filename
}

// calculateRelativePath calculates the relative path from fromPath to toPath
// Both paths should be relative to the protoc working directory (where buf.gen.yaml is)
func (c *Config) calculateRelativePath(fromPath, toPath string) string {
	// Clean the paths to handle . and .. properly
	fromPath = filepath.Clean(fromPath)
	toPath = filepath.Clean(toPath)
	
	// Calculate relative path
	relPath, err := filepath.Rel(fromPath, toPath)
	if err != nil {
		// Fallback to absolute path if relative calculation fails
		return toPath
	}
	
	// Convert to forward slashes for TypeScript imports
	relPath = filepath.ToSlash(relPath)
	
	// Ensure path starts with ./ for relative imports
	if !strings.HasPrefix(relPath, "./") && !strings.HasPrefix(relPath, "../") {
		relPath = "./" + relPath
	}
	
	return relPath
}

// Helper function to convert PascalCase to camelCase
func toCamelCase(s string) string {
	if len(s) == 0 {
		return s
	}
	
	// Convert first character to lowercase
	return strings.ToLower(s[:1]) + s[1:]
}