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

//go:build js && wasm

package wasm

import (
	"encoding/json"
	"syscall/js"
)

// CreateJSResponse creates a JavaScript-compatible response object
// This is used by all generated WASM service methods to return standardized responses
func CreateJSResponse(success bool, message string, data any) any {
	response := map[string]any{
		"success": success,
		"message": message,
	}

	if data != nil {
		response["data"] = data
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		// Fallback error response if marshaling fails
		errorResponse := map[string]any{
			"success": false,
			"message": "Failed to marshal response: " + err.Error(),
		}
		errorBytes, _ := json.Marshal(errorResponse)
		return js.Global().Get("JSON").Call("parse", string(errorBytes))
	}

	return js.Global().Get("JSON").Call("parse", string(responseBytes))
}