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
	"strings"

	wasmjsv1 "github.com/panyam/protoc-gen-go-wasmjs/proto/gen/go/wasmjs/v1"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
)

// ProtoAnalyzer provides pure functions for analyzing protobuf definitions.
// These functions extract metadata and properties from proto files, services,
// methods, messages, and fields without any side effects or external dependencies.
type ProtoAnalyzer struct{}

// NewProtoAnalyzer creates a new ProtoAnalyzer instance.
func NewProtoAnalyzer() *ProtoAnalyzer {
	return &ProtoAnalyzer{}
}

// ExtractPackageName extracts the package name from a fully qualified message type.
// For example, "library.v1.Book" returns "library.v1".
// This is used when analyzing message dependencies across packages.
func (pa *ProtoAnalyzer) ExtractPackageName(fullMessageType string) string {
	parts := strings.Split(fullMessageType, ".")
	if len(parts) <= 1 {
		return ""
	}
	// Take all parts except the last one (which is the message name)
	return strings.Join(parts[:len(parts)-1], ".")
}

// ExtractMessageName extracts the message name from a fully qualified message type.
// For example, "library.v1.Book" returns "Book".
// This is used when generating TypeScript interfaces and class names.
func (pa *ProtoAnalyzer) ExtractMessageName(fullMessageType string) string {
	parts := strings.Split(fullMessageType, ".")
	if len(parts) == 0 {
		return ""
	}
	// Return the last part (which is the message name)
	return parts[len(parts)-1]
}

// GetProtoFieldType returns the proto field type as a string representation.
// This analyzes the field's protobuf kind and returns a standardized string
// that can be used for type conversion and template generation.
func (pa *ProtoAnalyzer) GetProtoFieldType(field *protogen.Field) string {
	kind := field.Desc.Kind()
	switch kind.String() {
	case "double", "float", "int32", "int64", "uint32", "uint64", "sint32", "sint64", "fixed32", "fixed64", "sfixed32", "sfixed64":
		return kind.String()
	case "bool":
		return "bool"
	case "string":
		return "string"
	case "bytes":
		return "bytes"
	case "message":
		return string(field.Message.Desc.Name())
	case "enum":
		return string(field.Enum.Desc.Name())
	default:
		return kind.String()
	}
}

// GetFullyQualifiedMessageType returns the fully qualified type name for a message field.
// This includes the package name and message name, used for cross-package references.
// For example, a field of type Book in package library.v1 returns "library.v1.Book".
func (pa *ProtoAnalyzer) GetFullyQualifiedMessageType(field *protogen.Field) string {
	if field.Message == nil {
		return ""
	}

	// Get the full package name + message name
	packageName := string(field.Message.Desc.ParentFile().Package())
	messageName := string(field.Message.Desc.Name())

	if packageName != "" {
		return packageName + "." + messageName
	}
	return messageName
}

// IsMapField checks if a protobuf field represents a map type.
// Map fields are implemented as repeated message fields with special map entry messages.
func (pa *ProtoAnalyzer) IsMapField(field *protogen.Field) bool {
	if field.Message != nil {
		// For map fields, protogen might not set IsList but the message will be a map entry
		return field.Message.Desc.IsMapEntry()
	}
	return false
}

// GetMapKeyValueTypes extracts the key and value types from a map field.
// Map entry messages have exactly 2 fields: key (field 1) and value (field 2).
// Returns ("any", "any") if the field is not a valid map field.
func (pa *ProtoAnalyzer) GetMapKeyValueTypes(field *protogen.Field) (keyType, valueType string) {
	if !pa.IsMapField(field) {
		return "any", "any"
	}

	mapEntry := field.Message

	// Map entry messages have exactly 2 fields: key (field 1) and value (field 2)
	var keyField, valueField *protogen.Field
	for _, f := range mapEntry.Fields {
		if f.Desc.Number() == 1 {
			keyField = f
		} else if f.Desc.Number() == 2 {
			valueField = f
		}
	}

	if keyField == nil || valueField == nil {
		return "any", "any"
	}

	// Get the proto types for key and value
	keyType = pa.GetProtoFieldType(keyField)
	valueType = pa.GetProtoFieldType(valueField)

	return keyType, valueType
}

// IsBrowserProvidedService checks if a service is marked as browser-provided using annotations.
// Browser-provided services are implemented by JavaScript code rather than Go WASM,
// allowing access to browser APIs like fetch, localStorage, etc.
func (pa *ProtoAnalyzer) IsBrowserProvidedService(service *protogen.Service) bool {
	if service.Desc.Options() != nil {
		if browserOpts := proto.GetExtension(service.Desc.Options(), wasmjsv1.E_BrowserProvided); browserOpts != nil {
			if provided, ok := browserOpts.(bool); ok {
				return provided
			}
		}
	}
	return false
}

// IsTypeScriptFactoryFile checks if a proto file is marked as a TypeScript factory definition file.
// Factory files generate combined factory + deserializer TypeScript code for all messages
// imported from the same package.
func (pa *ProtoAnalyzer) IsTypeScriptFactoryFile(file *protogen.File) bool {
	if file.Desc.Options() != nil {
		if factoryOpt := proto.GetExtension(file.Desc.Options(), wasmjsv1.E_TsFactory); factoryOpt != nil {
			if isFactory, ok := factoryOpt.(bool); ok {
				return isFactory
			}
		}
	}
	return false
}

// GetCustomMethodName retrieves the custom JavaScript method name from wasmjs annotations.
// Returns empty string if no custom name is specified.
// This allows proto methods to have different names in the generated JavaScript API.
func (pa *ProtoAnalyzer) GetCustomMethodName(method *protogen.Method) string {
	if method.Desc.Options() != nil {
		if nameOpt := proto.GetExtension(method.Desc.Options(), wasmjsv1.E_WasmMethodName); nameOpt != nil {
			if name, ok := nameOpt.(string); ok && name != "" {
				return name
			}
		}
	}
	return ""
}

// IsAsyncMethod checks if a method is marked as async using wasmjs annotations.
// Async methods require callback parameters to prevent WASM deadlocks when calling browser APIs.
func (pa *ProtoAnalyzer) IsAsyncMethod(method *protogen.Method) bool {
	if method.Desc.Options() != nil {
		if asyncOpts := proto.GetExtension(method.Desc.Options(), wasmjsv1.E_AsyncMethod); asyncOpts != nil {
			if opts, ok := asyncOpts.(*wasmjsv1.AsyncMethodOptions); ok && opts != nil {
				return opts.GetIsAsync()
			}
		}
	}
	return false
}

// IsMethodExcluded checks if a method is marked for exclusion from WASM generation.
// Excluded methods won't appear in the generated JavaScript API.
func (pa *ProtoAnalyzer) IsMethodExcluded(method *protogen.Method) bool {
	if method.Desc.Options() != nil {
		if excludeOpt := proto.GetExtension(method.Desc.Options(), wasmjsv1.E_WasmMethodExclude); excludeOpt != nil {
			if exclude, ok := excludeOpt.(bool); ok {
				return exclude
			}
		}
	}
	return false
}

// IsServiceExcluded checks if a service is marked for exclusion from WASM generation.
// Excluded services won't appear in the generated JavaScript API.
func (pa *ProtoAnalyzer) IsServiceExcluded(service *protogen.Service) bool {
	if service.Desc.Options() != nil {
		if excludeOpt := proto.GetExtension(service.Desc.Options(), wasmjsv1.E_WasmServiceExclude); excludeOpt != nil {
			if exclude, ok := excludeOpt.(bool); ok {
				return exclude
			}
		}
	}
	return false
}

// GetCustomServiceName retrieves the custom JavaScript service name from wasmjs annotations.
// Returns empty string if no custom name is specified.
// This affects the namespaced API structure (e.g., namespace.customName.method).
func (pa *ProtoAnalyzer) GetCustomServiceName(service *protogen.Service) string {
	if service.Desc.Options() != nil {
		if nameOpt := proto.GetExtension(service.Desc.Options(), wasmjsv1.E_WasmServiceName); nameOpt != nil {
			if name, ok := nameOpt.(string); ok && name != "" {
				return name
			}
		}
	}
	return ""
}

// GetBaseFileName extracts the filename without extension from a proto file path.
// For example, "proto/library/v1/library.proto" returns "library".
// This is used for generating TypeScript file names and import paths.
func (pa *ProtoAnalyzer) GetBaseFileName(protoFile string) string {
	// Extract filename from path
	parts := strings.Split(protoFile, "/")
	filename := parts[len(parts)-1]

	// Remove .proto extension
	if strings.HasSuffix(filename, ".proto") {
		filename = filename[:len(filename)-6]
	}

	return filename
}

// GetOneofGroups extracts all oneof group names from a message.
// Returns a slice of unique oneof group names present in the message.
// This is used for TypeScript interface generation with proper oneof handling.
func (pa *ProtoAnalyzer) GetOneofGroups(message *protogen.Message) []string {
	var oneofGroups []string
	oneofMap := make(map[string]bool)

	for _, oneof := range message.Oneofs {
		oneofName := string(oneof.Desc.Name())
		if !oneofMap[oneofName] {
			oneofGroups = append(oneofGroups, oneofName)
			oneofMap[oneofName] = true
		}
	}

	return oneofGroups
}

// IsNestedMessage checks if a message is nested within another message.
// This is determined by checking if the message's parent is another message
// rather than a file descriptor.
func (pa *ProtoAnalyzer) IsNestedMessage(message *protogen.Message) bool {
	// In protogen, nested messages have a parent message
	// We can check this by looking at the parent descriptor
	parent := message.Desc.Parent()

	// If parent is not a file descriptor, it's nested
	switch parent.(type) {
	case interface{ IsFile() bool }:
		return false // Parent is a file, so not nested
	default:
		return true // Parent is something else (likely a message), so nested
	}
}
