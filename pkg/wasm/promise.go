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
	"fmt"
	"syscall/js"
)

// CreateJSPromise creates a JavaScript Promise that can be resolved or rejected from Go.
//
// The executor function is called synchronously when CreateJSPromise is invoked.
// For async work, spawn a goroutine inside the executor:
//
//	return wasm.CreateJSPromise(func(resolve, reject func(any)) {
//	    go func() {
//	        result, err := doAsyncWork()
//	        if err != nil {
//	            reject(err.Error())
//	            return
//	        }
//	        resolve(result)
//	    }()
//	})
//
// The Promise is returned immediately to JavaScript, allowing the JS event loop
// to continue while the goroutine does its work.
func CreateJSPromise(executor func(resolve, reject func(any))) js.Value {
	promiseConstructor := js.Global().Get("Promise")

	return promiseConstructor.New(js.FuncOf(func(this js.Value, args []js.Value) any {
		resolveJS := args[0]
		rejectJS := args[1]

		resolve := func(value any) {
			resolveJS.Invoke(js.ValueOf(value))
		}

		reject := func(err any) {
			// Create a proper JS Error object for better stack traces in JS
			errorConstructor := js.Global().Get("Error")
			jsError := errorConstructor.New(fmt.Sprintf("%v", err))
			rejectJS.Invoke(jsError)
		}

		executor(resolve, reject)
		return nil
	}))
}

// CreateResolvedPromise creates an already-resolved Promise with the given value.
// This is useful for returning immediate success without spawning a goroutine.
//
//	if cachedResult != nil {
//	    return wasm.CreateResolvedPromise(cachedResult)
//	}
func CreateResolvedPromise(value any) js.Value {
	return js.Global().Get("Promise").Call("resolve", js.ValueOf(value))
}

// CreateRejectedPromise creates an already-rejected Promise with the given error message.
// This is useful for returning immediate errors without spawning a goroutine.
//
//	if req == nil {
//	    return wasm.CreateRejectedPromise("request is required")
//	}
func CreateRejectedPromise(errorMsg string) js.Value {
	errorConstructor := js.Global().Get("Error")
	jsError := errorConstructor.New(errorMsg)
	return js.Global().Get("Promise").Call("reject", jsError)
}

// CreateRejectedPromiseFromError creates an already-rejected Promise from a Go error.
func CreateRejectedPromiseFromError(err error) js.Value {
	if err == nil {
		return CreateRejectedPromise("unknown error")
	}
	return CreateRejectedPromise(err.Error())
}
