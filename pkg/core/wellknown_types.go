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

import "github.com/panyam/protokit/wellknown"

// WellKnownTypeMapping represents a mapping from protobuf well-known type to TypeScript.
type WellKnownTypeMapping struct {
	ProtoType    string // Full proto type name (e.g., "google.protobuf.Timestamp")
	TSType       string // TypeScript type name to use (e.g., "Timestamp")
	ImportSource string // Where to import from (e.g., "@bufbuild/protobuf")
	IsNative     bool   // Whether this maps to a native TS type (e.g., Date)
}

// WellKnownTypesMapper handles mapping of protobuf well-known types to TypeScript types.
// Wraps protokit/wellknown.Registry with TS-specific default mappings.
type WellKnownTypesMapper struct {
	registry *wellknown.Registry
}

// NewWellKnownTypesMapper creates a new mapper with default @bufbuild/protobuf mappings.
func NewWellKnownTypesMapper() *WellKnownTypesMapper {
	r := wellknown.NewRegistry()

	// Core wrapper types
	r.Register("google.protobuf.DoubleValue", "DoubleValue", "@bufbuild/protobuf", false)
	r.Register("google.protobuf.FloatValue", "FloatValue", "@bufbuild/protobuf", false)
	r.Register("google.protobuf.Int64Value", "Int64Value", "@bufbuild/protobuf", false)
	r.Register("google.protobuf.UInt64Value", "UInt64Value", "@bufbuild/protobuf", false)
	r.Register("google.protobuf.Int32Value", "Int32Value", "@bufbuild/protobuf", false)
	r.Register("google.protobuf.UInt32Value", "UInt32Value", "@bufbuild/protobuf", false)
	r.Register("google.protobuf.BoolValue", "BoolValue", "@bufbuild/protobuf", false)
	r.Register("google.protobuf.StringValue", "StringValue", "@bufbuild/protobuf", false)
	r.Register("google.protobuf.BytesValue", "BytesValue", "@bufbuild/protobuf", false)

	// Time-related types
	r.Register("google.protobuf.Timestamp", "Timestamp", "@bufbuild/protobuf/wkt", false)
	r.Register("google.protobuf.Duration", "Duration", "@bufbuild/protobuf", false)

	// Structural types
	r.Register("google.protobuf.Any", "Any", "@bufbuild/protobuf/wkt", false)
	r.Register("google.protobuf.Empty", "Empty", "@bufbuild/protobuf", false)
	r.Register("google.protobuf.Struct", "Struct", "@bufbuild/protobuf", false)
	r.Register("google.protobuf.Value", "Value", "@bufbuild/protobuf", false)
	r.Register("google.protobuf.ListValue", "ListValue", "@bufbuild/protobuf", false)
	r.Register("google.protobuf.NullValue", "NullValue", "@bufbuild/protobuf", false)

	// Field mask
	r.Register("google.protobuf.FieldMask", "FieldMask", "@bufbuild/protobuf/wkt", false)

	// Type definitions
	r.Register("google.protobuf.Type", "Type", "@bufbuild/protobuf", false)
	r.Register("google.protobuf.Field", "Field", "@bufbuild/protobuf", false)
	r.Register("google.protobuf.Enum", "Enum", "@bufbuild/protobuf", false)
	r.Register("google.protobuf.EnumValue", "EnumValue", "@bufbuild/protobuf", false)
	r.Register("google.protobuf.Option", "Option", "@bufbuild/protobuf", false)
	r.Register("google.protobuf.SourceContext", "SourceContext", "@bufbuild/protobuf", false)

	return &WellKnownTypesMapper{registry: r}
}

// GetMapping returns the TypeScript mapping for a protobuf type, if it exists.
func (m *WellKnownTypesMapper) GetMapping(protoType string) (WellKnownTypeMapping, bool) {
	tm, ok := m.registry.Get(protoType)
	if !ok {
		return WellKnownTypeMapping{}, false
	}
	return WellKnownTypeMapping{
		ProtoType:    tm.ProtoFullName,
		TSType:       tm.TargetType,
		ImportSource: tm.ImportPath,
		IsNative:     tm.IsNative,
	}, true
}

// IsWellKnownType checks if a type is a well-known protobuf type.
func (m *WellKnownTypesMapper) IsWellKnownType(protoType string) bool {
	return m.registry.IsWellKnown(protoType)
}

// GetAllMappings returns all configured mappings.
func (m *WellKnownTypesMapper) GetAllMappings() map[string]WellKnownTypeMapping {
	result := make(map[string]WellKnownTypeMapping)
	for k, tm := range m.registry.AllMappings() {
		result[k] = WellKnownTypeMapping{
			ProtoType:    tm.ProtoFullName,
			TSType:       tm.TargetType,
			ImportSource: tm.ImportPath,
			IsNative:     tm.IsNative,
		}
	}
	return result
}

// OverrideMapping allows customization of a specific type mapping.
func (m *WellKnownTypesMapper) OverrideMapping(protoType, tsType, importSource string, isNative bool) {
	m.registry.Register(protoType, tsType, importSource, isNative)
}
