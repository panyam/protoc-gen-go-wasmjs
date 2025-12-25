# Implementation Plan: Invocation Style Refactoring

## Overview

Refactor the `async_method` annotation into a clearer `invocation_style` system that properly distinguishes between:
- **SYNC**: Returns value directly, executes on JS thread
- **PROMISE**: Returns JS Promise, executes in goroutine (safe for browser calls)
- **CALLBACK**: Takes callback parameter, executes in goroutine (legacy/streaming-like)

## Background

### The Problem
Go WASM runs on the same thread as JavaScript's event loop. When Go blocks (e.g., waiting for a browser API that returns a Promise), it deadlocks because JS can't resolve the Promise while Go is blocking.

### The Solution
- **SYNC**: Only for methods that don't call browser APIs returning Promises
- **PROMISE/CALLBACK**: Spawn a goroutine, return immediately to JS, deliver result later

### Key Insight
Go can create and return JavaScript Promises directly using `js.Global().Get("Promise").New(executor)`. This gives JS clean async/await syntax.

---

## Implementation Phases

### Phase 1: Add Promise Helpers to Runtime
**Status**: [ ] Not Started

**Files to create/modify:**
- [ ] `pkg/wasm/promise.go` (NEW)

**Tasks:**
1. Create `CreateJSPromise(executor func(resolve, reject func(any))) js.Value`
2. Create `CreateResolvedPromise(value any) js.Value`
3. Create `CreateRejectedPromise(errorMsg string) js.Value`

**Code:**
```go
//go:build js && wasm

package wasm

import (
    "fmt"
    "syscall/js"
)

// CreateJSPromise creates a JavaScript Promise that can be resolved/rejected from Go.
func CreateJSPromise(executor func(resolve, reject func(any))) js.Value {
    promiseConstructor := js.Global().Get("Promise")

    return promiseConstructor.New(js.FuncOf(func(this js.Value, args []js.Value) any {
        resolveJS := args[0]
        rejectJS := args[1]

        resolve := func(value any) {
            resolveJS.Invoke(js.ValueOf(value))
        }
        reject := func(err any) {
            errorConstructor := js.Global().Get("Error")
            jsError := errorConstructor.New(fmt.Sprintf("%v", err))
            rejectJS.Invoke(jsError)
        }

        executor(resolve, reject)
        return nil
    }))
}

// CreateResolvedPromise creates an already-resolved Promise.
func CreateResolvedPromise(value any) js.Value {
    return js.Global().Get("Promise").Call("resolve", js.ValueOf(value))
}

// CreateRejectedPromise creates an already-rejected Promise.
func CreateRejectedPromise(errorMsg string) js.Value {
    errorConstructor := js.Global().Get("Error")
    jsError := errorConstructor.New(errorMsg)
    return js.Global().Get("Promise").Call("reject", jsError)
}
```

---

### Phase 2: Update Proto Annotations
**Status**: [ ] Not Started

**Files to modify:**
- [ ] `proto/wasmjs/v1/annotations.proto`

**Tasks:**
1. Add `InvocationStyle` enum
2. Add `invocation_style` method extension
3. Keep `async_method` for backward compatibility (deprecated)

**Code:**
```protobuf
// InvocationStyle specifies how a method should be invoked across the JS/WASM boundary.
enum InvocationStyle {
  // SYNC: Method returns value directly. Executes on JS thread.
  // WARNING: Do not use for methods that call browser APIs returning Promises.
  INVOCATION_STYLE_SYNC = 0;

  // PROMISE: Method returns a JavaScript Promise. Executes in goroutine.
  // Safe for methods that call browser APIs. JS can use async/await.
  INVOCATION_STYLE_PROMISE = 1;

  // CALLBACK: Method takes a callback parameter. Executes in goroutine.
  // Safe for methods that call browser APIs. Legacy pattern.
  INVOCATION_STYLE_CALLBACK = 2;
}

// invocation_style specifies how a method should be invoked.
extend google.protobuf.MethodOptions {
  InvocationStyle invocation_style = 50010;
}
```

4. Regenerate Go code: `buf generate` or `protoc`

---

### Phase 3: Update Builders & Analyzers
**Status**: [ ] Not Started

**Files to modify:**
- [ ] `pkg/builders/shared_types.go`
- [ ] `pkg/core/proto_analyzer.go`
- [ ] `pkg/builders/go_data_builder.go`
- [ ] `pkg/builders/ts_data_builder.go`

**Tasks:**

#### 3.1 Update shared_types.go
Replace `IsAsync bool` with proper invocation style:

```go
// MethodData represents a gRPC method prepared for template generation.
type MethodData struct {
    // ... existing fields ...

    // InvocationStyle determines how the method is called across JS/WASM boundary
    // Values: "sync", "promise", "callback"
    InvocationStyle string

    // Deprecated: Use InvocationStyle instead
    IsAsync bool
}

// Helper methods
func (m *MethodData) IsSync() bool {
    return m.InvocationStyle == "" || m.InvocationStyle == "sync"
}

func (m *MethodData) IsPromise() bool {
    return m.InvocationStyle == "promise"
}

func (m *MethodData) IsCallback() bool {
    return m.InvocationStyle == "callback"
}
```

#### 3.2 Update proto_analyzer.go
Parse the new annotation:

```go
func (a *ProtoAnalyzer) extractInvocationStyle(method *protogen.Method) string {
    // Check new annotation first
    if style := proto.GetExtension(method.Desc.Options(), wasmjsv1.E_InvocationStyle); style != nil {
        switch style.(wasmjsv1.InvocationStyle) {
        case wasmjsv1.InvocationStyle_INVOCATION_STYLE_PROMISE:
            return "promise"
        case wasmjsv1.InvocationStyle_INVOCATION_STYLE_CALLBACK:
            return "callback"
        default:
            return "sync"
        }
    }

    // Backward compatibility: check old async_method annotation
    if asyncOpts := proto.GetExtension(method.Desc.Options(), wasmjsv1.E_AsyncMethod); asyncOpts != nil {
        if opts, ok := asyncOpts.(*wasmjsv1.AsyncMethodOptions); ok && opts.GetIsAsync() {
            return "callback"  // Old async_method maps to callback
        }
    }

    return "sync"
}
```

#### 3.3 Update data builders
Populate `InvocationStyle` field from analyzed methods.

---

### Phase 4: Update WASM Export Templates
**Status**: [ ] Not Started

**Files to modify:**
- [ ] `pkg/renderers/templates/wasm_exports.go.tmpl`

**Tasks:**
1. Add PROMISE branch using `CreateJSPromise`
2. Keep CALLBACK branch (rename from IsAsync)
3. SYNC branch unchanged

**Template changes:**
```go
{{- if .IsServerStreaming }}
// ... streaming code unchanged ...
{{- else if .IsPromise }}
// Promise method: returns JS Promise, executes in goroutine
func (exports *...) {{ .GoFuncName }}(this js.Value, args []js.Value) any {
    // Parse request...

    return wasm.CreateJSPromise(func(resolve, reject func(any)) {
        go func() {
            ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
            defer cancel()

            resp, err := exports.{{ $serviceName }}.{{ .Name }}(ctx, req)
            if err != nil {
                reject(err.Error())
                return
            }

            // Marshal and resolve
            responseJSON, _ := marshaller.Marshal(resp, wasm.MarshalOptions{...})
            var jsObject interface{}
            json.Unmarshal(responseJSON, &jsObject)
            resolve(js.ValueOf(jsObject))
        }()
    })
}
{{- else if .IsCallback }}
// Callback method: takes callback, executes in goroutine
// ... existing IsAsync code ...
{{- else }}
// Sync method: returns value directly
// ... existing sync code ...
{{- end }}
```

---

### Phase 5: Update TypeScript Client Templates
**Status**: [ ] Not Started

**Files to modify:**
- [ ] `pkg/renderers/templates/client_simple.ts.tmpl`
- [ ] `runtime/src/client/service-client.ts`
- [ ] `runtime/src/client/wasm-bundle.ts`

**Tasks:**

#### 5.1 Update client_simple.ts.tmpl

**Interface signatures:**
```typescript
{{- range .Methods }}
    {{- if .ShouldGenerate }}
        {{- if .IsServerStreaming }}
    {{ .JSName }}(request: {{ .RequestTSType }}, callback: (response: {{ .ResponseTSType }} | null, error: string | null, done: boolean) => boolean): void;
        {{- else if .IsPromise }}
    {{ .JSName }}(request: {{ .RequestTSType }}): Promise<{{ .ResponseTSType }}>;
        {{- else if .IsCallback }}
    {{ .JSName }}(request: {{ .RequestTSType }}, callback: (response: {{ .ResponseTSType }}, error?: string) => void): void;
        {{- else }}
    {{ .JSName }}(request: {{ .RequestTSType }}): {{ .ResponseTSType }};
        {{- end }}
    {{- end }}
{{- end }}
```

**Implementation:**
```typescript
{{- if .IsPromise }}
    async {{ .JSName }}(request: {{ .RequestTSType }}): Promise<{{ .ResponseTSType }}> {
        return this.callMethodAsync('{{ $serviceJSName }}.{{ .JSName }}', request);
    }
{{- else if .IsCallback }}
    {{ .JSName }}(request: {{ .RequestTSType }}, callback: (response: {{ .ResponseTSType }}, error?: string) => void): void {
        this.callMethodWithCallback('{{ $serviceJSName }}.{{ .JSName }}', request, callback);
    }
{{- else }}
    {{ .JSName }}(request: {{ .RequestTSType }}): {{ .ResponseTSType }} {
        return this.callMethodSync('{{ $serviceJSName }}.{{ .JSName }}', request);
    }
{{- end }}
```

#### 5.2 Update service-client.ts

```typescript
export abstract class ServiceClient {
    // Sync call - returns value directly
    protected callMethodSync<TRequest, TResponse>(
        methodPath: string,
        request: TRequest
    ): TResponse {
        return this.bundle.callMethodSync(methodPath, request);
    }

    // Promise call - WASM returns a Promise
    protected callMethodAsync<TRequest, TResponse>(
        methodPath: string,
        request: TRequest
    ): Promise<TResponse> {
        return this.bundle.callMethodAsync(methodPath, request);
    }

    // Callback call - WASM invokes callback later
    protected callMethodWithCallback<TRequest>(
        methodPath: string,
        request: TRequest,
        callback: (response: any, error?: string) => void
    ): void {
        this.bundle.callMethodWithCallback(methodPath, request, callback);
    }
}
```

#### 5.3 Update wasm-bundle.ts

```typescript
// Sync call - WASM returns value directly
public callMethodSync<TRequest, TResponse>(
    methodPath: string,
    request: TRequest
): TResponse {
    const wasmMethod = this.getWasmMethod(methodPath);
    const wasmResponse = wasmMethod(JSON.stringify(request));

    if (!wasmResponse.success) {
        throw new WasmError(wasmResponse.message, methodPath);
    }
    return wasmResponse.data;
}

// Promise call - WASM returns a Promise
public callMethodAsync<TRequest, TResponse>(
    methodPath: string,
    request: TRequest
): Promise<TResponse> {
    const wasmMethod = this.getWasmMethod(methodPath);
    // WASM method returns a JS Promise directly
    return wasmMethod(JSON.stringify(request));
}

// Callback call - WASM invokes callback later
public callMethodWithCallback<TRequest>(
    methodPath: string,
    request: TRequest,
    callback: (response: any, error?: string) => void
): void {
    const wasmMethod = this.getWasmMethod(methodPath);
    wasmMethod(JSON.stringify(request), callback);
}
```

---

### Phase 6: Update Browser Client Templates
**Status**: [ ] Not Started

**Files to modify:**
- [ ] `pkg/renderers/templates/wasm_browser_clients.go.tmpl`
- [ ] `pkg/renderers/templates/client_simple.ts.tmpl` (BrowserServer section)

**Tasks:**

#### 6.1 Update wasm_browser_clients.go.tmpl

```go
{{- range .Methods }}
{{- if .IsPromise }}
// {{ .Name }} calls the browser-provided {{ .Name }} method (returns Promise).
func (c *{{ $service.Name }}Client) {{ .Name }}(ctx context.Context, req *{{ .RequestType }}) *wasm.Future[*{{ .ResponseType }}] {
    return wasm.CallBrowserServiceAsync[*{{ .RequestType }}, *{{ .ResponseType }}](
        c.channel, ctx, "{{ $service.Name }}", "{{ .JSName }}", req,
    )
}
{{- else }}
// {{ .Name }} calls the browser-provided {{ .Name }} method (sync).
func (c *{{ $service.Name }}Client) {{ .Name }}(ctx context.Context, req *{{ .RequestType }}) (*{{ .ResponseType }}, error) {
    return wasm.CallBrowserService[*{{ .RequestType }}, *{{ .ResponseType }}](
        c.channel, ctx, "{{ $service.Name }}", "{{ .JSName }}", req,
    )
}
{{- end }}
{{- end }}
```

#### 6.2 Update BrowserServer interface in client_simple.ts.tmpl

```typescript
{{- range .BrowserClients }}
export interface {{ .Name }}Server {
    {{- range .Methods }}
        {{- if .ShouldGenerate }}
            {{- if .IsPromise }}
    {{ .JSName }}(request: {{ .RequestTSType }}): Promise<{{ .ResponseTSType }}>;
            {{- else }}
    {{ .JSName }}(request: {{ .RequestTSType }}): {{ .ResponseTSType }};
            {{- end }}
        {{- end }}
    {{- end }}
}
{{- end }}
```

---

### Phase 7: Testing & Documentation
**Status**: [ ] Not Started

**Files to modify:**
- [ ] `examples/proto/` - Add example methods with different invocation styles
- [ ] `README.md` - Update documentation
- [ ] `pkg/wasm/doc.go` - Update package docs

**Tasks:**
1. Add test proto with all three invocation styles
2. Test SYNC method (no browser calls)
3. Test PROMISE method (with browser calls, verify no deadlock)
4. Test CALLBACK method (legacy pattern)
5. Update documentation with examples

---

## Migration Guide

### For existing users

Old annotation:
```protobuf
rpc GetGame(GetGameRequest) returns (GameState) {
    option (wasmjs.v1.async_method) = { is_async: true };
}
```

New annotation:
```protobuf
rpc GetGame(GetGameRequest) returns (GameState) {
    option (wasmjs.v1.invocation_style) = INVOCATION_STYLE_PROMISE;
}
```

The old `async_method` annotation will continue to work (maps to CALLBACK style).

---

## Summary Table

| Style | JS Signature | Go Behavior | Use When |
|-------|--------------|-------------|----------|
| SYNC | `method(req): Resp` | Execute on JS thread | Simple methods, no browser calls |
| PROMISE | `method(req): Promise<Resp>` | Return Promise, goroutine | Methods that call browser APIs |
| CALLBACK | `method(req, cb): void` | Return void, goroutine | Legacy, prefer PROMISE |
| STREAMING | `method(req, cb): void` | Return void, goroutine | Server streaming RPCs |

---

## Progress Tracking

- [x] Phase 1: Promise Helpers - COMPLETED (pkg/wasm/promise.go)
- [x] Phase 2: Proto Annotations - COMPLETED (annotations.proto)
- [x] Phase 3: Builders & Analyzers - COMPLETED (proto_analyzer.go, shared_types.go)
- [x] Phase 4: WASM Export Templates - COMPLETED (wasm_exports.go.tmpl)
- [x] Phase 5: TypeScript Client Templates - COMPLETED (client_simple.ts.tmpl, runtime)
- [x] Phase 6: Browser Client Templates - COMPLETED (wasm_browser_clients.go.tmpl)
- [x] Phase 7: Testing & Documentation - COMPLETED (examples regenerated)

**Implementation Complete: December 2025**

Additional cleanup performed:
- Removed dead code: pkg/generator/ directory, wasm.go.tmpl
- Renamed callMethod to callMethodAsync for clarity
- Updated examples with PROMISE and CALLBACK invocation styles
