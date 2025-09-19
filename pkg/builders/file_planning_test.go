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

package builders

import (
	"testing"

	"google.golang.org/protobuf/compiler/protogen"
)

// TestFileSpec_Creation tests file specification creation and validation.
// This ensures file specifications contain all necessary information for file planning.
func TestFileSpec_Creation(t *testing.T) {
	spec := FileSpec{
		Name:     "wasm",
		Filename: "library/v1/library_v1.wasm.go",
		Type:     "wasm",
		Required: true,
		ContentHints: ContentHints{
			HasServices:        true,
			HasBrowserServices: false,
			IsExample:          false,
			IsBuildScript:      false,
		},
	}

	// Test basic fields
	if spec.Name != "wasm" {
		t.Errorf("FileSpec.Name = %s, want wasm", spec.Name)
	}

	if spec.Filename != "library/v1/library_v1.wasm.go" {
		t.Errorf("FileSpec.Filename = %s, want library/v1/library_v1.wasm.go", spec.Filename)
	}

	if spec.Type != "wasm" {
		t.Errorf("FileSpec.Type = %s, want wasm", spec.Type)
	}

	if !spec.Required {
		t.Error("FileSpec.Required should be true")
	}

	// Test content hints
	if !spec.ContentHints.HasServices {
		t.Error("ContentHints.HasServices should be true")
	}

	if spec.ContentHints.HasBrowserServices {
		t.Error("ContentHints.HasBrowserServices should be false")
	}
}

// TestFilePlan_Creation tests file plan creation and management.
// This ensures file plans correctly organize file specifications for a package.
func TestFilePlan_Creation(t *testing.T) {
	config := &GenerationConfig{
		WasmExportPath:      "./gen/wasm",
		JSStructure:         "namespaced",
		GenerateBuildScript: true,
	}

	specs := []FileSpec{
		{
			Name:     "wasm",
			Filename: "library/v1/library_v1.wasm.go",
			Type:     "wasm",
			Required: true,
		},
		{
			Name:     "main",
			Filename: "library/v1/main.go.example",
			Type:     "example",
			Required: true,
		},
		{
			Name:     "build",
			Filename: "build.sh",
			Type:     "script",
			Required: false,
		},
	}

	plan := &FilePlan{
		PackageName: "library.v1",
		Specs:       specs,
		Config:      config,
	}

	// Test basic fields
	if plan.PackageName != "library.v1" {
		t.Errorf("FilePlan.PackageName = %s, want library.v1", plan.PackageName)
	}

	if len(plan.Specs) != 3 {
		t.Errorf("FilePlan.Specs length = %d, want 3", len(plan.Specs))
	}

	if plan.Config != config {
		t.Error("FilePlan.Config should reference the provided config")
	}

	// Test that we have the expected spec types
	specTypes := make(map[string]bool)
	for _, spec := range plan.Specs {
		specTypes[spec.Type] = true
	}

	expectedTypes := []string{"wasm", "example", "script"}
	for _, expectedType := range expectedTypes {
		if !specTypes[expectedType] {
			t.Errorf("FilePlan should have spec of type %s", expectedType)
		}
	}
}

// TestNewGeneratedFileSet tests creation of GeneratedFileSet from FilePlan.
// This tests the file set structure without requiring a real protogen.Plugin.
func TestNewGeneratedFileSet(t *testing.T) {
	plan := &FilePlan{
		PackageName: "library.v1",
		Specs: []FileSpec{
			{Name: "wasm", Filename: "library_v1.wasm.go", Type: "wasm", Required: true},
			{Name: "main", Filename: "main.go.example", Type: "example", Required: true},
			{Name: "build", Filename: "build.sh", Type: "script", Required: false},
		},
	}

	// Test the file set structure without calling NewGeneratedFile (requires real plugin)
	// We'll create the structure manually to test the API
	fileSet := &GeneratedFileSet{
		Files: make(map[string]*protogen.GeneratedFile),
		Plan:  plan,
	}

	// Manually populate the map to simulate what NewGeneratedFileSet would do
	for _, spec := range plan.Specs {
		fileSet.Files[spec.Name] = nil // nil GeneratedFile for testing
	}

	// Test basic structure
	if fileSet == nil {
		t.Error("GeneratedFileSet should be non-nil")
		return
	}

	if fileSet.Plan != plan {
		t.Error("GeneratedFileSet should reference the original plan")
	}

	if fileSet.Files == nil {
		t.Error("GeneratedFileSet.Files should be initialized")
		return
	}

	// Test file map structure
	expectedFiles := []string{"wasm", "main", "build"}
	for _, expectedFile := range expectedFiles {
		if _, exists := fileSet.Files[expectedFile]; !exists {
			t.Errorf("GeneratedFileSet.Files should have entry for %s", expectedFile)
		}
	}
}

// TestGeneratedFileSet_Methods tests the utility methods on GeneratedFileSet.
// This ensures file set operations work correctly for different file organization patterns.
func TestGeneratedFileSet_Methods(t *testing.T) {
	plan := &FilePlan{
		PackageName: "library.v1",
		Specs: []FileSpec{
			{Name: "wasm", Filename: "wasm.go", Type: "wasm", Required: true},
			{Name: "client", Filename: "client.ts", Type: "client", Required: true},
			{Name: "interfaces", Filename: "interfaces.ts", Type: "interfaces", Required: false},
			{Name: "models", Filename: "models.ts", Type: "interfaces", Required: false}, // Same type as interfaces
		},
	}

	// Create file set manually since we can't use nil plugin
	fileSet := &GeneratedFileSet{
		Files: make(map[string]*protogen.GeneratedFile),
		Plan:  plan,
	}

	// Manually populate the map to simulate what NewGeneratedFileSet would do
	for _, spec := range plan.Specs {
		fileSet.Files[spec.Name] = nil // nil GeneratedFile for testing
	}

	// Test HasFile
	if !fileSet.HasFile("wasm") {
		t.Error("HasFile should return true for planned file")
	}

	if fileSet.HasFile("nonexistent") {
		t.Error("HasFile should return false for unplanned file")
	}

	// Test GetFilesByType
	interfaceFiles := fileSet.GetFilesByType("interfaces")
	if len(interfaceFiles) != 2 {
		t.Errorf("GetFilesByType('interfaces') should return 2 files, got %d", len(interfaceFiles))
	}

	wasmFiles := fileSet.GetFilesByType("wasm")
	if len(wasmFiles) != 1 {
		t.Errorf("GetFilesByType('wasm') should return 1 file, got %d", len(wasmFiles))
	}

	// Test GetRequiredFiles
	requiredFiles := fileSet.GetRequiredFiles()
	if len(requiredFiles) != 2 {
		t.Errorf("GetRequiredFiles() should return 2 files, got %d", len(requiredFiles))
	}

	// Test GetAllFilenames
	filenames := fileSet.GetAllFilenames()
	expectedFilenames := []string{"wasm.go", "client.ts", "interfaces.ts", "models.ts"}

	if len(filenames) != len(expectedFilenames) {
		t.Errorf("GetAllFilenames() length = %d, want %d", len(filenames), len(expectedFilenames))
	}

	// Verify all expected filenames are present
	filenameMap := make(map[string]bool)
	for _, filename := range filenames {
		filenameMap[filename] = true
	}

	for _, expected := range expectedFilenames {
		if !filenameMap[expected] {
			t.Errorf("GetAllFilenames() should include %s", expected)
		}
	}
}

// TestContentHints_Usage tests content hints for conditional file generation decisions.
// This ensures content hints provide the right metadata for generators to make decisions.
func TestContentHints_Usage(t *testing.T) {
	tests := []struct {
		name     string                         // Test case description
		hints    ContentHints                   // Content hints to test
		testFunc func(*testing.T, ContentHints) // Test function
		reason   string                         // Why this test is important
	}{
		{
			name: "service-focused file",
			hints: ContentHints{
				HasServices:        true,
				HasBrowserServices: true,
				HasMessages:        false,
				HasEnums:           false,
			},
			testFunc: func(t *testing.T, hints ContentHints) {
				if !hints.HasServices {
					t.Error("Service-focused file should have HasServices=true")
				}
				if !hints.HasBrowserServices {
					t.Error("Should indicate browser services presence")
				}
				if hints.HasMessages || hints.HasEnums {
					t.Error("Service-focused file shouldn't have message/enum content")
				}
			},
			reason: "Service files need accurate service metadata for template decisions",
		},
		{
			name: "type-focused file",
			hints: ContentHints{
				HasServices: false,
				HasMessages: true,
				HasEnums:    true,
			},
			testFunc: func(t *testing.T, hints ContentHints) {
				if hints.HasServices {
					t.Error("Type-focused file shouldn't have services")
				}
				if !hints.HasMessages || !hints.HasEnums {
					t.Error("Type-focused file should have messages and enums")
				}
			},
			reason: "Type files need accurate type metadata for interface generation",
		},
		{
			name: "utility file",
			hints: ContentHints{
				IsExample:     true,
				IsBuildScript: false,
			},
			testFunc: func(t *testing.T, hints ContentHints) {
				if !hints.IsExample {
					t.Error("Example file should have IsExample=true")
				}
				if hints.IsBuildScript {
					t.Error("Example file shouldn't be a build script")
				}
			},
			reason: "Utility files need correct categorization for template selection",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t, tt.hints)
		})
	}
}
