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
	"github.com/panyam/protokit/fields"
	"github.com/panyam/protokit/messages"

	wasmjsv1 "github.com/panyam/protoc-gen-go-wasmjs/proto/gen/go/wasmjs/v1"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
)

// InvocationStyle represents how a method is invoked across the JS/WASM boundary.
type InvocationStyle int

// InvocationStyle constants for method invocation patterns.
const (
	InvocationStyleSync     InvocationStyle = iota // Returns value directly, executes on JS thread
	InvocationStylePromise                         // Returns JS Promise, executes in goroutine
	InvocationStyleCallback                        // Takes callback parameter, executes in goroutine
)

// String returns the string representation of the invocation style.
func (s InvocationStyle) String() string {
	switch s {
	case InvocationStyleSync:
		return "sync"
	case InvocationStylePromise:
		return "promise"
	case InvocationStyleCallback:
		return "callback"
	default:
		return "sync"
	}
}

// ProtoAnalyzer provides functions for analyzing protobuf definitions.
// Generic methods delegate to protokit. WASM-specific annotation methods are defined here.
type ProtoAnalyzer struct{}

// NewProtoAnalyzer creates a new ProtoAnalyzer instance.
func NewProtoAnalyzer() *ProtoAnalyzer {
	return &ProtoAnalyzer{}
}

// --- Generic methods (delegated to protokit) ---

// ExtractPackageName extracts the package name from a fully qualified message type.
// Delegates to protokit/messages.ExtractPackageName.
func (pa *ProtoAnalyzer) ExtractPackageName(fullMessageType string) string {
	return messages.ExtractPackageName(fullMessageType)
}

// ExtractMessageName extracts the message name from a fully qualified message type.
// Delegates to protokit/messages.ExtractMessageName.
func (pa *ProtoAnalyzer) ExtractMessageName(fullMessageType string) string {
	return messages.ExtractMessageName(fullMessageType)
}

// GetProtoFieldType returns the proto field type as a string.
// For scalars returns the kind ("string", "int32", etc.).
// For messages returns the message name. For enums returns the enum name.
func (pa *ProtoAnalyzer) GetProtoFieldType(field *protogen.Field) string {
	kind := fields.GetFieldKind(field)
	switch kind {
	case "message":
		return string(field.Message.Desc.Name())
	case "enum":
		return string(field.Enum.Desc.Name())
	default:
		return kind
	}
}

// GetFullyQualifiedMessageType returns the fully qualified type name for a message field.
// Delegates to protokit/messages.GetFullyQualifiedType.
func (pa *ProtoAnalyzer) GetFullyQualifiedMessageType(field *protogen.Field) string {
	return messages.GetFullyQualifiedType(field)
}

// IsMapField checks if a protobuf field represents a map type.
// Delegates to protokit/fields.IsMapField.
func (pa *ProtoAnalyzer) IsMapField(field *protogen.Field) bool {
	return fields.IsMapField(field)
}

// GetMapKeyValueTypes extracts the key and value types from a map field.
// Returns ("any", "any") if the field is not a valid map field.
func (pa *ProtoAnalyzer) GetMapKeyValueTypes(field *protogen.Field) (keyType, valueType string) {
	keyField, valueField := fields.GetMapKeyValueFields(field)
	if keyField == nil || valueField == nil {
		return "any", "any"
	}
	return pa.GetProtoFieldType(keyField), pa.GetProtoFieldType(valueField)
}

// IsNestedMessage checks if a message is nested within another message.
// Delegates to protokit/messages.IsNestedMessage.
func (pa *ProtoAnalyzer) IsNestedMessage(message *protogen.Message) bool {
	return messages.IsNestedMessage(message)
}

// GetOneofGroups extracts all oneof group names from a message.
// Delegates to protokit/messages.GetOneofGroups.
func (pa *ProtoAnalyzer) GetOneofGroups(message *protogen.Message) []string {
	return messages.GetOneofGroups(message)
}

// GetBaseFileName extracts the filename without extension from a proto file path.
// Delegates to protokit/messages.GetBaseFileName.
func (pa *ProtoAnalyzer) GetBaseFileName(protoFile string) string {
	return messages.GetBaseFileName(protoFile)
}

// --- WASM-specific annotation methods (stay in this package) ---

// IsBrowserProvidedService checks if a service is marked as browser-provided.
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

// IsTypeScriptFactoryFile checks if a proto file is marked as a TypeScript factory definition.
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

// GetInvocationStyle returns the invocation style for a method.
func (pa *ProtoAnalyzer) GetInvocationStyle(method *protogen.Method) InvocationStyle {
	if method.Desc.Options() != nil {
		if styleOpt := proto.GetExtension(method.Desc.Options(), wasmjsv1.E_InvocationStyle); styleOpt != nil {
			if style, ok := styleOpt.(wasmjsv1.InvocationStyle); ok {
				switch style {
				case wasmjsv1.InvocationStyle_INVOCATION_STYLE_PROMISE:
					return InvocationStylePromise
				case wasmjsv1.InvocationStyle_INVOCATION_STYLE_CALLBACK:
					return InvocationStyleCallback
				case wasmjsv1.InvocationStyle_INVOCATION_STYLE_SYNC:
					return InvocationStyleSync
				}
			}
		}
		// Backward compatibility: check deprecated async_method annotation
		if asyncOpts := proto.GetExtension(method.Desc.Options(), wasmjsv1.E_AsyncMethod); asyncOpts != nil {
			if opts, ok := asyncOpts.(*wasmjsv1.AsyncMethodOptions); ok && opts != nil && opts.GetIsAsync() {
				return InvocationStyleCallback
			}
		}
	}
	return InvocationStyleSync
}

// IsMethodExcluded checks if a method is marked for exclusion from WASM generation.
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
