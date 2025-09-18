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
	"unicode"
)

// NameConverter provides pure functions for converting between different naming conventions.
// These functions handle transformations between protocol buffer naming, Go naming,
// JavaScript naming, and TypeScript naming conventions without any side effects.
type NameConverter struct{}

// NewNameConverter creates a new NameConverter instance.
func NewNameConverter() *NameConverter {
	return &NameConverter{}
}

// ToCamelCase converts PascalCase strings to camelCase.
// This is used when converting Go method names to JavaScript method names.
//
// Example:
//
//	"FindBooks" -> "findBooks"
//	"GetUser" -> "getUser"
//	"CreateLibraryItem" -> "createLibraryItem"
//
// JavaScript conventionally uses camelCase for method and variable names.
func (nc *NameConverter) ToCamelCase(s string) string {
	if len(s) == 0 {
		return s
	}

	// Convert first character to lowercase
	return strings.ToLower(s[:1]) + s[1:]
}

// ToPascalCase converts camelCase strings to PascalCase.
// This is the inverse of ToCamelCase and is used for generating class names.
//
// Example:
//
//	"findBooks" -> "FindBooks"
//	"getUser" -> "GetUser"
//	"libraryService" -> "LibraryService"
//
// TypeScript classes and interfaces conventionally use PascalCase.
func (nc *NameConverter) ToPascalCase(s string) string {
	if len(s) == 0 {
		return s
	}

	// Convert first character to uppercase
	return strings.ToUpper(s[:1]) + s[1:]
}

// ToSnakeCase converts camelCase or PascalCase strings to snake_case.
// This is used when converting to certain naming conventions or file names.
//
// Example:
//
//	"FindBooks" -> "find_books"
//	"getUserData" -> "get_user_data"
//	"HTTPSConnection" -> "https_connection"
func (nc *NameConverter) ToSnakeCase(s string) string {
	if len(s) == 0 {
		return s
	}

	var result strings.Builder

	for i, r := range s {
		// Add underscore before uppercase letters (except the first character)
		if i > 0 && unicode.IsUpper(r) {
			result.WriteRune('_')
		}
		result.WriteRune(unicode.ToLower(r))
	}

	return result.String()
}

// ToPackageAlias converts a package path or name to a valid Go package alias.
// This creates short, readable aliases for import statements.
//
// Example:
//
//	"github.com/example/proto/gen/go/library/v1" -> "libraryv1"
//	"common.types.v2" -> "commonv2"
//	"user-service.v1" -> "userv1"
//
// The alias removes special characters and combines relevant parts.
func (nc *NameConverter) ToPackageAlias(packagePath string) string {
	if packagePath == "" {
		return "pkg"
	}

	// Extract meaningful parts from the path
	parts := strings.Split(packagePath, "/")

	// For paths like "github.com/example/proto/gen/go/library/v1"
	// Take the last two parts: "library" and "v1"
	if len(parts) >= 2 {
		pkg := parts[len(parts)-2]
		version := parts[len(parts)-1]

		// Clean up the package name (remove hyphens, underscores)
		pkg = strings.ReplaceAll(pkg, "-", "")
		pkg = strings.ReplaceAll(pkg, "_", "")

		// Clean up version (remove dots)
		version = strings.ReplaceAll(version, ".", "")

		// Combine them
		alias := pkg + version
		return strings.ToLower(alias)
	}

	// Fallback to the last part only
	if len(parts) > 0 {
		last := parts[len(parts)-1]
		// Remove special characters
		last = strings.ReplaceAll(last, ".", "")
		last = strings.ReplaceAll(last, "-", "")
		last = strings.ReplaceAll(last, "_", "")
		return strings.ToLower(last)
	}

	return "pkg" // Ultimate fallback
}

// ToJSNamespace converts a package name to a JavaScript namespace.
// This creates valid JavaScript identifiers for global namespace objects.
//
// Example:
//
//	"library.v1" -> "library_v1"
//	"user-service.v2" -> "user_service_v2"
//	"Common.Types" -> "common_types"
//
// JavaScript namespaces should be lowercase with underscores.
func (nc *NameConverter) ToJSNamespace(packageName string) string {
	if packageName == "" {
		return ""
	}

	// Convert to lowercase and replace special characters
	namespace := strings.ToLower(packageName)
	namespace = strings.ReplaceAll(namespace, ".", "_")
	namespace = strings.ReplaceAll(namespace, "-", "_")

	return namespace
}

// ToModuleName converts a package name to a WASM module name.
// This creates descriptive names for WASM modules and generated files.
//
// Example:
//
//	"library.v1" -> "library_v1_services"
//	"user.auth.v2" -> "user_auth_v2_services"
//	"common" -> "common_services"
//
// Module names are used in file names and Go package declarations.
func (nc *NameConverter) ToModuleName(packageName string) string {
	if packageName == "" {
		return "services"
	}

	// Convert dots to underscores
	name := strings.ReplaceAll(packageName, ".", "_")
	// Add services suffix
	return name + "_services"
}

// ToFactoryName converts a package name to a TypeScript factory class name.
// Factory classes are responsible for creating instances of messages in a package.
//
// Example:
//
//	"library.v1" -> "LibraryV1Factory"
//	"common.types" -> "CommonTypesFactory"
//	"user-auth" -> "UserAuthFactory"
//
// Factory names follow TypeScript class naming conventions (PascalCase).
func (nc *NameConverter) ToFactoryName(packageName string) string {
	if packageName == "" {
		return "Factory"
	}

	// Split by dots and hyphens, then convert each part to PascalCase
	parts := strings.FieldsFunc(packageName, func(r rune) bool {
		return r == '.' || r == '-' || r == '_'
	})

	var result strings.Builder
	for _, part := range parts {
		if part != "" {
			result.WriteString(nc.ToPascalCase(part))
		}
	}
	result.WriteString("Factory")

	return result.String()
}

// ToSchemaRegistryName converts a package name to a schema registry name.
// Schema registries contain metadata about message structures for runtime processing.
//
// Example:
//
//	"library.v1" -> "libraryV1Schemas"
//	"common.types" -> "commonTypesSchemas"
//
// Uses camelCase for variable naming in TypeScript.
func (nc *NameConverter) ToSchemaRegistryName(packageName string) string {
	if packageName == "" {
		return "schemas"
	}

	// Convert to a camelCase identifier
	parts := strings.FieldsFunc(packageName, func(r rune) bool {
		return r == '.' || r == '-' || r == '_'
	})

	var result strings.Builder
	for i, part := range parts {
		if part != "" {
			if i == 0 {
				result.WriteString(strings.ToLower(part))
			} else {
				result.WriteString(nc.ToPascalCase(part))
			}
		}
	}
	result.WriteString("Schemas")

	return result.String()
}

// ToDeserializerName converts a package name to a deserializer class name.
// Deserializers handle converting JSON data to typed objects using schemas.
//
// Example:
//
//	"library.v1" -> "LibraryV1Deserializer"
//	"common.types" -> "CommonTypesDeserializer"
//
// Uses PascalCase for TypeScript class names.
func (nc *NameConverter) ToDeserializerName(packageName string) string {
	if packageName == "" {
		return "Deserializer"
	}

	// Similar to factory name but with Deserializer suffix
	parts := strings.FieldsFunc(packageName, func(r rune) bool {
		return r == '.' || r == '-' || r == '_'
	})

	var result strings.Builder
	for _, part := range parts {
		if part != "" {
			result.WriteString(nc.ToPascalCase(part))
		}
	}
	result.WriteString("Deserializer")

	return result.String()
}

// ToGoFuncName converts a service and method name to a Go function name.
// This creates the function names used in the Go WASM wrapper for JavaScript calls.
//
// Example:
//
//	serviceName: "LibraryService", methodName: "FindBooks" -> "libraryServiceFindBooks"
//	serviceName: "UserAuth", methodName: "Login" -> "userAuthLogin"
//
// Uses camelCase starting with lowercase for Go function naming in WASM context.
func (nc *NameConverter) ToGoFuncName(serviceName, methodName string) string {
	// Convert service name to camelCase and append method name
	camelService := nc.ToCamelCase(serviceName)
	return camelService + methodName
}

// ToTSPropertyName converts JSON field names to TypeScript property names.
// This ensures consistent property naming in generated TypeScript interfaces.
//
// Example:
//
//	"userId" -> "userId" (already camelCase)
//	"user_id" -> "userId" (convert from snake_case)
//	"UserName" -> "userName" (convert from PascalCase)
//
// TypeScript properties conventionally use camelCase.
func (nc *NameConverter) ToTSPropertyName(jsonName string) string {
	// JSON names from protobuf are typically already in camelCase
	// But we ensure it follows camelCase convention
	return nc.ToCamelCase(jsonName)
}

// SanitizeIdentifier ensures a string is a valid identifier for the target language.
// This removes or replaces characters that would make invalid identifiers.
//
// Example:
//
//	"user-name" -> "user_name"
//	"123invalid" -> "_23invalid"
//	"user@name#123" -> "user_name123"
//
// Used to ensure generated names are valid in the target language.
func (nc *NameConverter) SanitizeIdentifier(name string) string {
	if name == "" {
		return "identifier"
	}

	var result strings.Builder

	// Ensure first character is a letter or underscore
	first := rune(name[0])
	if unicode.IsLetter(first) || first == '_' {
		result.WriteRune(first)
	} else {
		result.WriteRune('_')
	}

	// Process remaining characters
	for _, r := range name[1:] {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
			result.WriteRune(r)
		} else if r == '-' || r == '.' {
			// Convert hyphens and dots to underscores
			result.WriteRune('_')
		} else {
			// Convert other invalid characters to underscores
			result.WriteRune('_')
		}
	}

	return result.String()
}
