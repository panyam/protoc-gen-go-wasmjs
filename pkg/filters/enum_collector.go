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
	"strings"
	
	"google.golang.org/protobuf/compiler/protogen"
	
	"github.com/panyam/protoc-gen-go-wasmjs/pkg/core"
)

// EnumValueInfo represents metadata about a protobuf enum value.
type EnumValueInfo struct {
	Name    string // Original proto enum value name (e.g., "WAITING_FOR_PLAYERS")
	Number  int32  // Enum value number
	Comment string // Value comment from proto
}

// EnumInfo represents metadata about a protobuf enum for generation.
type EnumInfo struct {
	Name               string          // Enum name (e.g., "GameStatus")
	FullyQualifiedName string          // Fully qualified name (e.g., "connect4.GameStatus")
	PackageName        string          // Proto package name
	ProtoFile          string          // Source proto file path
	IsNested           bool            // Whether this is nested within a message
	Values             []EnumValueInfo // All values in the enum
	Comment            string          // Leading comment from proto
}

// EnumCollector provides logic for collecting and filtering protobuf enums.
// This handles the complex logic of finding all enums across packages while
// applying appropriate filtering criteria.
type EnumCollector struct {
	analyzer *core.ProtoAnalyzer
}

// NewEnumCollector creates a new enum collector with the given analyzer.
func NewEnumCollector(analyzer *core.ProtoAnalyzer) *EnumCollector {
	return &EnumCollector{
		analyzer: analyzer,
	}
}

// CollectEnums collects all enums from the given files that meet the criteria.
// This handles recursive collection of nested enums and applies filtering rules.
// Returns CollectionResult with detailed statistics about what was found and collected.
func (ec *EnumCollector) CollectEnums(files []*protogen.File, criteria *FilterCriteria) CollectionResult[EnumInfo] {
	var enums []EnumInfo
	totalFound := 0
	
	for _, file := range files {
		// Collect top-level enums
		for _, enum := range file.Enums {
			totalFound++
			enumInfo := ec.buildEnumInfo(enum, file, false)
			enums = append(enums, enumInfo)
		}
		
		// Collect nested enums from messages if not excluded
		if !criteria.ExcludeNestedEnums {
			for _, message := range file.Messages {
				nestedEnums, nestedCount := ec.collectNestedEnums(message, file, criteria)
				enums = append(enums, nestedEnums...)
				totalFound += nestedCount
			}
		}
	}
	
	return NewCollectionResult(enums, totalFound, len(files))
}

// collectNestedEnums recursively collects nested enum definitions from messages.
// Returns the collected enums and the total count (including any that were filtered).
func (ec *EnumCollector) collectNestedEnums(message *protogen.Message, file *protogen.File, criteria *FilterCriteria) ([]EnumInfo, int) {
	var nestedEnums []EnumInfo
	totalCount := 0
	
	// Collect enums directly nested in this message
	for _, enum := range message.Enums {
		totalCount++
		enumInfo := ec.buildEnumInfo(enum, file, true)
		nestedEnums = append(nestedEnums, enumInfo)
	}
	
	// Recursively collect enums from nested messages
	if !criteria.ExcludeNestedEnums {
		for _, nested := range message.Messages {
			// Skip map entry messages
			if criteria.ExcludeMapEntries && nested.Desc.IsMapEntry() {
				continue
			}
			
			deeplyNestedEnums, deepCount := ec.collectNestedEnums(nested, file, criteria)
			nestedEnums = append(nestedEnums, deeplyNestedEnums...)
			totalCount += deepCount
		}
	}
	
	return nestedEnums, totalCount
}

// buildEnumInfo constructs EnumInfo from a protogen.Enum.
// This extracts essential metadata needed for filtering and generation decisions.
func (ec *EnumCollector) buildEnumInfo(enum *protogen.Enum, file *protogen.File, isNested bool) EnumInfo {
	enumName := string(enum.Desc.Name())
	packageName := string(file.Desc.Package())
	
	// Build fully qualified name (e.g., "connect4.GameStatus")
	fullyQualifiedName := packageName + "." + enumName
	
	// Build enum values
	var values []EnumValueInfo
	for _, value := range enum.Values {
		valueInfo := EnumValueInfo{
			Name:    string(value.Desc.Name()),
			Number:  int32(value.Desc.Number()),
			Comment: strings.TrimSpace(string(value.Comments.Leading)),
		}
		values = append(values, valueInfo)
	}
	
	return EnumInfo{
		Name:               enumName,
		FullyQualifiedName: fullyQualifiedName,
		PackageName:        packageName,
		ProtoFile:          file.Desc.Path(),
		IsNested:           isNested,
		Values:             values,
		Comment:            strings.TrimSpace(string(enum.Comments.Leading)),
	}
}

// HasAnyEnums checks if any enums would be collected with the given criteria.
// This is useful for early termination when no enums would be generated.
func (ec *EnumCollector) HasAnyEnums(files []*protogen.File, criteria *FilterCriteria) bool {
	for _, file := range files {
		// Check top-level enums
		if len(file.Enums) > 0 {
			return true
		}
		
		// Check nested enums if not excluded
		if !criteria.ExcludeNestedEnums {
			for _, message := range file.Messages {
				if ec.hasNestedEnums(message) {
					return true
				}
			}
		}
	}
	return false
}

// hasNestedEnums recursively checks if a message contains any nested enums.
func (ec *EnumCollector) hasNestedEnums(message *protogen.Message) bool {
	// Check direct nested enums
	if len(message.Enums) > 0 {
		return true
	}
	
	// Check nested messages recursively
	for _, nested := range message.Messages {
		if ec.hasNestedEnums(nested) {
			return true
		}
	}
	
	return false
}

// CollectTopLevelEnums collects only top-level enums (no nesting).
// This is a specialized method for cases where nested enums should be ignored.
func (ec *EnumCollector) CollectTopLevelEnums(files []*protogen.File, criteria *FilterCriteria) CollectionResult[EnumInfo] {
	var enums []EnumInfo
	totalFound := 0
	
	for _, file := range files {
		for _, enum := range file.Enums {
			totalFound++
			enumInfo := ec.buildEnumInfo(enum, file, false)
			enums = append(enums, enumInfo)
		}
	}
	
	return NewCollectionResult(enums, totalFound, len(files))
}

// CollectEnumsByPackage collects enums grouped by package name.
// This is useful for generating package-specific TypeScript files.
func (ec *EnumCollector) CollectEnumsByPackage(files []*protogen.File, criteria *FilterCriteria) map[string]CollectionResult[EnumInfo] {
	packageEnums := make(map[string]CollectionResult[EnumInfo])
	
	// Group files by package
	packageFiles := make(map[string][]*protogen.File)
	for _, file := range files {
		packageName := string(file.Desc.Package())
		packageFiles[packageName] = append(packageFiles[packageName], file)
	}
	
	// Collect enums for each package
	for packageName, pkgFiles := range packageFiles {
		result := ec.CollectEnums(pkgFiles, criteria)
		packageEnums[packageName] = result
	}
	
	return packageEnums
}

// GetEnumNames extracts just the enum names from a collection result.
// This is useful for template generation where only names are needed.
func (ec *EnumCollector) GetEnumNames(result CollectionResult[EnumInfo]) []string {
	var names []string
	for _, enum := range result.Items {
		names = append(names, enum.Name)
	}
	return names
}

// GetEnumsByFile groups enums by their source proto file.
// This is useful for generating file-specific imports and references.
func (ec *EnumCollector) GetEnumsByFile(result CollectionResult[EnumInfo]) map[string][]EnumInfo {
	enumsByFile := make(map[string][]EnumInfo)
	
	for _, enum := range result.Items {
		enumsByFile[enum.ProtoFile] = append(enumsByFile[enum.ProtoFile], enum)
	}
	
	return enumsByFile
}

// GetEnumValues extracts all enum values from a collection result.
// This is useful for validation and code generation that needs all possible values.
func (ec *EnumCollector) GetEnumValues(result CollectionResult[EnumInfo]) map[string][]EnumValueInfo {
	enumValues := make(map[string][]EnumValueInfo)
	
	for _, enum := range result.Items {
		enumValues[enum.Name] = enum.Values
	}
	
	return enumValues
}
