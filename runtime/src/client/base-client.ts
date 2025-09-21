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

import { BrowserServiceManager } from '../browser/service-manager.js';
import { WasmError } from './types.js';

/**
 * Base WASM service client containing all non-template-dependent logic
 */
export abstract class WASMServiceClient {
    protected wasm: any;
    protected wasmLoadPromise: Promise<void> | null = null;
    protected browserServiceManager: BrowserServiceManager | null = null;

    constructor() {
        this.browserServiceManager = new BrowserServiceManager();
    }

    /**
     * Register a browser service implementation
     * Can be used to register browser services from any package
     */
    public registerBrowserService(name: string, implementation: any): void {
        if (!this.browserServiceManager) {
            throw new Error('Browser service manager not initialized');
        }
        this.browserServiceManager.registerService(name, implementation);
    }

    /**
     * Check if WASM is ready for operations
     */
    public isReady(): boolean {
        return this.wasm !== null && this.wasm !== undefined;
    }

    /**
     * Wait for WASM to be ready (use during initialization)
     */
    public async waitUntilReady(): Promise<void> {
        if (!this.wasmLoadPromise) {
            throw new Error('WASM loading not started. Call loadWasm() first.');
        }
        await this.wasmLoadPromise;
    }

    /**
     * Internal method to call WASM functions with JSON conversion
     */
    public callMethod<TRequest, TResponse>(
        methodPath: string,
        request: TRequest
    ): Promise<TResponse> {
        this.ensureWASMLoaded();

        try {
            // Convert request to JSON
            const jsonReq = JSON.parse(JSON.stringify(request));
            const wasmMethod = this.getWasmMethod(methodPath);
            const wasmResponse = wasmMethod(JSON.stringify(jsonReq));

            if (!wasmResponse.success) {
                throw new WasmError(wasmResponse.message, methodPath);
            }

            // Return response data directly
            return wasmResponse.data;
        } catch (error) {
            if (error instanceof WasmError) {
                throw error;
            }
            throw new WasmError(
                `Call error: ${error instanceof Error ? error.message : String(error)}`,
                methodPath
            );
        }
    }

    /**
     * Internal method to call async WASM functions with callback
     */
    public callMethodWithCallback<TRequest>(
        methodPath: string,
        request: TRequest,
        callback: (response: any, error?: string) => void
    ): Promise<void> {
        this.ensureWASMLoaded();

        try {
            // Convert request to JSON
            const jsonReq = JSON.parse(JSON.stringify(request));
            const wasmMethod = this.getWasmMethod(methodPath);
            
            // Call WASM method with callback function
            const wasmResponse = wasmMethod(JSON.stringify(jsonReq), callback);

            if (!wasmResponse.success) {
                throw new WasmError(wasmResponse.message, methodPath);
            }

            // Async methods return immediately
            return Promise.resolve();
        } catch (error) {
            if (error instanceof WasmError) {
                throw error;
            }
            throw new WasmError(
                `Call error: ${error instanceof Error ? error.message : String(error)}`,
                methodPath
            );
        }
    }

    /**
     * Internal method to call server streaming WASM functions
     */
    public callStreamingMethod<TRequest, TResponse>(
        methodPath: string,
        request: TRequest,
        callback: (response: TResponse | null, error: string | null, done: boolean) => boolean
    ): void {
        this.ensureWASMLoaded();

        try {
            // Convert request to JSON
            const jsonReq = JSON.parse(JSON.stringify(request));
            const wasmMethod = this.getWasmMethod(methodPath);

            // Wrap the callback to parse JSON responses
            const wrappedCallback = (responseStr: string | null, error: string | null, done: boolean): boolean => {
                let response: TResponse | null = null;
                if (responseStr && !error) {
                    try {
                        response = JSON.parse(responseStr);
                    } catch (e) {
                        // If parsing fails, pass the raw string
                        response = responseStr as any;
                    }
                }
                return callback(response, error, done);
            };

            // Call WASM streaming method with wrapped callback
            const wasmResponse = wasmMethod(JSON.stringify(jsonReq), wrappedCallback);

            if (!wasmResponse.success) {
                throw new WasmError(wasmResponse.message, methodPath);
            }

            // Streaming methods return immediately
        } catch (error) {
            if (error instanceof WasmError) {
                throw error;
            }
            throw new WasmError(
                `Streaming call error: ${error instanceof Error ? error.message : String(error)}`,
                methodPath
            );
        }
    }

    /**
     * Ensure WASM module is loaded (synchronous version for service calls)
     */
    protected ensureWASMLoaded(): void {
        if (!this.isReady()) {
            throw new Error('WASM module not loaded. Call loadWasm() and waitUntilReady() first.');
        }
    }

    /**
     * Abstract method for getting WASM method function by path
     * Implementation depends on API structure (namespaced, flat, service_based)
     */
    protected abstract getWasmMethod(methodPath: string): Function;
}
