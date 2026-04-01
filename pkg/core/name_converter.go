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

	"github.com/panyam/protokit/names"
)

// NameConverter provides functions for converting between different naming conventions.
// Generic methods delegate to protokit/names. WASM/TS-specific methods are defined here.
type NameConverter struct{}

// NewNameConverter creates a new NameConverter instance.
func NewNameConverter() *NameConverter {
	return &NameConverter{}
}

// ToCamelCase converts PascalCase strings to camelCase.
// Delegates to protokit/names.ToCamelCase.
func (nc *NameConverter) ToCamelCase(s string) string {
	return names.ToCamelCase(s)
}

// ToPascalCase converts camelCase strings to PascalCase.
// Delegates to protokit/names.ToPascalCase.
func (nc *NameConverter) ToPascalCase(s string) string {
	return names.ToPascalCase(s)
}

// ToSnakeCase converts camelCase or PascalCase strings to snake_case.
// Delegates to protokit/names.ToSnakeCase.
func (nc *NameConverter) ToSnakeCase(s string) string {
	return names.ToSnakeCase(s)
}

// SanitizeIdentifier ensures a string is a valid identifier for the target language.
// Delegates to protokit/names.SanitizeIdentifier.
func (nc *NameConverter) SanitizeIdentifier(name string) string {
	return names.SanitizeIdentifier(name)
}

// ToPackageAlias converts a package path to a valid Go package alias.
// Uses the last two path segments for unique, readable identifiers.
func (nc *NameConverter) ToPackageAlias(packagePath string) string {
	if packagePath == "" {
		return "pkg"
	}
	parts := strings.Split(packagePath, "/")
	if len(parts) >= 2 {
		pkg := parts[len(parts)-2]
		version := parts[len(parts)-1]
		pkg = strings.ReplaceAll(pkg, "-", "")
		pkg = strings.ReplaceAll(pkg, "_", "")
		version = strings.ReplaceAll(version, ".", "")
		return strings.ToLower(pkg + version)
	}
	if len(parts) > 0 {
		last := parts[len(parts)-1]
		last = strings.ReplaceAll(last, ".", "")
		last = strings.ReplaceAll(last, "-", "")
		last = strings.ReplaceAll(last, "_", "")
		return strings.ToLower(last)
	}
	return "pkg"
}

// ToJSNamespace converts a package name to a JavaScript namespace.
func (nc *NameConverter) ToJSNamespace(packageName string) string {
	if packageName == "" {
		return ""
	}
	namespace := strings.ToLower(packageName)
	namespace = strings.ReplaceAll(namespace, ".", "_")
	namespace = strings.ReplaceAll(namespace, "-", "_")
	return namespace
}

// ToModuleName converts a package name to a WASM module name.
func (nc *NameConverter) ToModuleName(packageName string) string {
	if packageName == "" {
		return "services"
	}
	return strings.ReplaceAll(packageName, ".", "_") + "_services"
}

// ToFactoryName converts a package name to a TypeScript factory class name.
func (nc *NameConverter) ToFactoryName(packageName string) string {
	if packageName == "" {
		return "Factory"
	}
	parts := strings.FieldsFunc(packageName, func(r rune) bool {
		return r == '.' || r == '-' || r == '_'
	})
	var result strings.Builder
	for _, part := range parts {
		if part != "" {
			result.WriteString(names.ToPascalCase(part))
		}
	}
	result.WriteString("Factory")
	return result.String()
}

// ToSchemaRegistryName converts a package name to a schema registry name.
func (nc *NameConverter) ToSchemaRegistryName(packageName string) string {
	if packageName == "" {
		return "schemas"
	}
	parts := strings.FieldsFunc(packageName, func(r rune) bool {
		return r == '.' || r == '-' || r == '_'
	})
	var result strings.Builder
	for i, part := range parts {
		if part != "" {
			if i == 0 {
				result.WriteString(strings.ToLower(part))
			} else {
				result.WriteString(names.ToPascalCase(part))
			}
		}
	}
	result.WriteString("Schemas")
	return result.String()
}

// ToDeserializerName converts a package name to a deserializer class name.
func (nc *NameConverter) ToDeserializerName(packageName string) string {
	if packageName == "" {
		return "Deserializer"
	}
	parts := strings.FieldsFunc(packageName, func(r rune) bool {
		return r == '.' || r == '-' || r == '_'
	})
	var result strings.Builder
	for _, part := range parts {
		if part != "" {
			result.WriteString(names.ToPascalCase(part))
		}
	}
	result.WriteString("Deserializer")
	return result.String()
}

// ToGoFuncName converts a service and method name to a Go function name.
func (nc *NameConverter) ToGoFuncName(serviceName, methodName string) string {
	return names.ToCamelCase(serviceName) + methodName
}

// ToTSPropertyName converts JSON field names to TypeScript property names.
func (nc *NameConverter) ToTSPropertyName(jsonName string) string {
	return names.ToCamelCase(jsonName)
}
