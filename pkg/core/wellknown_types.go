// Copyright 2025 Sri Panyam
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package core

// WellKnownTypeMapping represents a mapping from protobuf well-known type to TypeScript
type WellKnownTypeMapping struct {
	ProtoType    string // Full proto type name (e.g., "google.protobuf.Timestamp")
	TSType       string // TypeScript type name to use (e.g., "Timestamp")
	ImportSource string // Where to import from (e.g., "@bufbuild/protobuf")
	IsNative     bool   // Whether this maps to a native TS type (e.g., Date)
}

// WellKnownTypesMapper handles mapping of protobuf well-known types to TypeScript types
type WellKnownTypesMapper struct {
	mappings map[string]WellKnownTypeMapping
}

// NewWellKnownTypesMapper creates a new mapper with default @bufbuild/protobuf mappings
func NewWellKnownTypesMapper() *WellKnownTypesMapper {
	mapper := &WellKnownTypesMapper{
		mappings: make(map[string]WellKnownTypeMapping),
	}
	mapper.initializeDefaults()
	return mapper
}

// initializeDefaults sets up the default mappings for well-known types using @bufbuild/protobuf
func (m *WellKnownTypesMapper) initializeDefaults() {
	// Core wrapper types from google/protobuf/wrappers.proto
	m.addMapping("google.protobuf.DoubleValue", "DoubleValue", "@bufbuild/protobuf")
	m.addMapping("google.protobuf.FloatValue", "FloatValue", "@bufbuild/protobuf")
	m.addMapping("google.protobuf.Int64Value", "Int64Value", "@bufbuild/protobuf")
	m.addMapping("google.protobuf.UInt64Value", "UInt64Value", "@bufbuild/protobuf")
	m.addMapping("google.protobuf.Int32Value", "Int32Value", "@bufbuild/protobuf")
	m.addMapping("google.protobuf.UInt32Value", "UInt32Value", "@bufbuild/protobuf")
	m.addMapping("google.protobuf.BoolValue", "BoolValue", "@bufbuild/protobuf")
	m.addMapping("google.protobuf.StringValue", "StringValue", "@bufbuild/protobuf")
	m.addMapping("google.protobuf.BytesValue", "BytesValue", "@bufbuild/protobuf")

	// Time-related types
	m.addMapping("google.protobuf.Timestamp", "Timestamp", "@bufbuild/protobuf/wkt")
	m.addMapping("google.protobuf.Duration", "Duration", "@bufbuild/protobuf")

	// Structural types
	m.addMapping("google.protobuf.Any", "Any", "@bufbuild/protobuf")
	m.addMapping("google.protobuf.Empty", "Empty", "@bufbuild/protobuf")
	m.addMapping("google.protobuf.Struct", "Struct", "@bufbuild/protobuf")
	m.addMapping("google.protobuf.Value", "Value", "@bufbuild/protobuf")
	m.addMapping("google.protobuf.ListValue", "ListValue", "@bufbuild/protobuf")
	m.addMapping("google.protobuf.NullValue", "NullValue", "@bufbuild/protobuf")

	// Field mask for partial updates
	m.addMapping("google.protobuf.FieldMask", "FieldMask", "@bufbuild/protobuf/wkt")

	// Type definitions (less commonly used in application code)
	m.addMapping("google.protobuf.Type", "Type", "@bufbuild/protobuf")
	m.addMapping("google.protobuf.Field", "Field", "@bufbuild/protobuf")
	m.addMapping("google.protobuf.Enum", "Enum", "@bufbuild/protobuf")
	m.addMapping("google.protobuf.EnumValue", "EnumValue", "@bufbuild/protobuf")
	m.addMapping("google.protobuf.Option", "Option", "@bufbuild/protobuf")

	// Source context
	m.addMapping("google.protobuf.SourceContext", "SourceContext", "@bufbuild/protobuf")

	// Commonly used Google API types (if needed in the future)
	// These would come from a different package, but we can add support later
	// m.addMapping("google.type.Date", "Date", "@bufbuild/protobuf/google/type")
	// m.addMapping("google.type.Money", "Money", "@bufbuild/protobuf/google/type")
}

// addMapping adds a type mapping to the mapper
func (m *WellKnownTypesMapper) addMapping(protoType, tsType, importSource string) {
	m.mappings[protoType] = WellKnownTypeMapping{
		ProtoType:    protoType,
		TSType:       tsType,
		ImportSource: importSource,
		IsNative:     false,
	}
}

// addNativeMapping adds a mapping to a native TypeScript type
func (m *WellKnownTypesMapper) addNativeMapping(protoType, tsType string) {
	m.mappings[protoType] = WellKnownTypeMapping{
		ProtoType:    protoType,
		TSType:       tsType,
		ImportSource: "",
		IsNative:     true,
	}
}

// GetMapping returns the TypeScript mapping for a protobuf type, if it exists
func (m *WellKnownTypesMapper) GetMapping(protoType string) (WellKnownTypeMapping, bool) {
	mapping, exists := m.mappings[protoType]
	return mapping, exists
}

// IsWellKnownType checks if a type is a well-known protobuf type
func (m *WellKnownTypesMapper) IsWellKnownType(protoType string) bool {
	_, exists := m.mappings[protoType]
	return exists
}

// GetAllMappings returns all configured mappings
func (m *WellKnownTypesMapper) GetAllMappings() map[string]WellKnownTypeMapping {
	// Return a copy to prevent external modification
	result := make(map[string]WellKnownTypeMapping)
	for k, v := range m.mappings {
		result[k] = v
	}
	return result
}

// OverrideMapping allows customization of a specific type mapping
func (m *WellKnownTypesMapper) OverrideMapping(protoType, tsType, importSource string, isNative bool) {
	m.mappings[protoType] = WellKnownTypeMapping{
		ProtoType:    protoType,
		TSType:       tsType,
		ImportSource: importSource,
		IsNative:     isNative,
	}
}

// Example customizations that could be applied:
// mapper.OverrideMapping("google.protobuf.Timestamp", "Date", "", true)  // Use native Date
// mapper.OverrideMapping("google.protobuf.Duration", "number", "", true) // Use number for milliseconds
