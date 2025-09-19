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
	"fmt"
	"log"

	"google.golang.org/protobuf/compiler/protogen"
)

// FileSpec defines a file that should be generated, including its logical name,
// output filename, and metadata about what type of content it will contain.
type FileSpec struct {
	// Logical name for this file (used for getting templates and data)
	Name string // e.g., "wasm", "client", "interfaces", "factory"

	// Actual output filename (decided by generator)
	Filename string // e.g., "library/v1/library_v1.wasm.go", "LibraryClient.ts"

	// File type for template selection and data building
	Type string // e.g., "wasm", "client", "interfaces", "models", "factory"

	// Whether this file is required or optional
	Required bool // If true, generation fails if template/data unavailable

	// Content hints for conditional rendering
	ContentHints ContentHints // Metadata about what content this file will have
}

// ContentHints provides metadata about what a file will contain.
// This helps generators make decisions about data building and template selection.
type ContentHints struct {
	HasServices        bool // File will contain service-related content
	HasMessages        bool // File will contain message-related content
	HasEnums           bool // File will contain enum-related content
	HasBrowserServices bool // File will contain browser service content
	IsExample          bool // File is an example/documentation file
	IsBuildScript      bool // File is a build/compilation script
}

// FilePlan represents the complete plan for files to be generated for a package.
// This is created by generators during the planning phase.
type FilePlan struct {
	PackageName string     // Package this plan is for
	Specs       []FileSpec // All files to be generated

	// Generation context
	Config *GenerationConfig // Configuration used for planning
}

// GeneratedFileSet represents a collection of protogen GeneratedFile objects
// that correspond to a FilePlan. This is created by generators and passed to renderers.
type GeneratedFileSet struct {
	Files map[string]*protogen.GeneratedFile // Maps FileSpec.Name to GeneratedFile
	Plan  *FilePlan                          // Original plan this set was created from
}

// NewGeneratedFileSet creates a new file set from a plan using protogen plugin.
func NewGeneratedFileSet(plan *FilePlan, plugin *protogen.Plugin) *GeneratedFileSet {
	files := make(map[string]*protogen.GeneratedFile)

	for _, spec := range plan.Specs {
		file := plugin.NewGeneratedFile(spec.Filename, "")
		log.Printf("Created GeneratedFile: %s -> %p", spec.Filename, file)
		files[spec.Name] = file
	}

	return &GeneratedFileSet{
		Files: files,
		Plan:  plan,
	}
}

// GetFile returns the GeneratedFile for a given logical file name.
// Returns nil if the file wasn't planned.
func (gfs *GeneratedFileSet) GetFile(name string) *protogen.GeneratedFile {
	return gfs.Files[name]
}

// HasFile checks if a file with the given logical name exists in the set.
func (gfs *GeneratedFileSet) HasFile(name string) bool {
	_, exists := gfs.Files[name]
	return exists
}

// GetFilesByType returns all files of a specific type.
// This is useful for batch operations on similar files.
func (gfs *GeneratedFileSet) GetFilesByType(fileType string) map[string]*protogen.GeneratedFile {
	result := make(map[string]*protogen.GeneratedFile)

	for _, spec := range gfs.Plan.Specs {
		if spec.Type == fileType {
			// Include files that were planned for this type, even if GeneratedFile is nil (for testing)
			result[spec.Name] = gfs.Files[spec.Name]
		}
	}

	return result
}

// GetRequiredFiles returns only the files marked as required.
// This helps ensure critical files are always generated.
func (gfs *GeneratedFileSet) GetRequiredFiles() map[string]*protogen.GeneratedFile {
	result := make(map[string]*protogen.GeneratedFile)

	for _, spec := range gfs.Plan.Specs {
		if spec.Required {
			// Include required files that were planned, even if GeneratedFile is nil (for testing)
			result[spec.Name] = gfs.Files[spec.Name]
		}
	}

	return result
}

// ValidateFileSet ensures all required files were created and are ready for rendering.
func (gfs *GeneratedFileSet) ValidateFileSet() error {
	for _, spec := range gfs.Plan.Specs {
		if spec.Required {
			if file := gfs.Files[spec.Name]; file == nil {
				return fmt.Errorf("required file '%s' (%s) was not created", spec.Name, spec.Filename)
			}
		}
	}

	return nil
}

// GetAllFilenames returns all filenames from the file set.
// This is useful for logging and debugging.
func (gfs *GeneratedFileSet) GetAllFilenames() []string {
	var filenames []string
	for _, spec := range gfs.Plan.Specs {
		filenames = append(filenames, spec.Filename)
	}
	return filenames
}
