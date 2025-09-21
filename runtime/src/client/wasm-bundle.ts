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
 * Configuration for API structure and bundle behavior
 */
export interface WASMBundleConfig {
    moduleName: string;
    apiStructure: 'namespaced' | 'flat' | 'service_based';
    jsNamespace: string;
}

/**
 * WASM Bundle - manages loading and shared access to a WASM module
 * One bundle per WASM file, shared by multiple service clients
 */
export class WASMBundle {
    private wasm: any = null;
    private wasmLoadPromise: Promise<void> | null = null;
    private browserServiceManager: BrowserServiceManager | null = null;
    private config: WASMBundleConfig;

    constructor(config: WASMBundleConfig) {
        this.config = config;
        this.browserServiceManager = new BrowserServiceManager();
    }

    /**
     * Register a browser service implementation
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
     * Load the WASM module asynchronously (singleton pattern)
     */
    public async loadWasm(wasmPath: string): Promise<void> {
        if (this.wasmLoadPromise) {
            return this.wasmLoadPromise;
        }

        this.wasmLoadPromise = this.loadWASMModule(wasmPath);
        return this.wasmLoadPromise;
    }

    /**
     * Get WASM method function by path
     */
    public getWasmMethod(methodPath: string): Function {
        this.ensureWASMLoaded();

        switch (this.config.apiStructure) {
            case 'namespaced':
                // Handle namespaced structure: namespace.service.method
                const parts = methodPath.split('.');
                let current = this.wasm;
                for (const part of parts) {
                    current = current[part];
                    if (!current) {
                        throw new Error(`Method not found: ${methodPath}`);
                    }
                }
                return current;

            case 'flat':
                // Handle flat structure: direct method name
                const method = this.wasm[methodPath];
                if (!method) {
                    throw new Error(`Method not found: ${methodPath}`);
                }
                return method;

            case 'service_based':
                // Handle service-based structure: services.service.method
                const serviceParts = methodPath.split('.');
                let serviceCurrent = this.wasm;
                for (const part of serviceParts) {
                    serviceCurrent = serviceCurrent[part];
                    if (!serviceCurrent) {
                        throw new Error(`Method not found: ${methodPath}`);
                    }
                }
                return serviceCurrent;

            default:
                throw new Error(`Unsupported API structure: ${this.config.apiStructure}`);
        }
    }

    /**
     * Internal method to call WASM functions with JSON conversion
     */
    public callMethod<TRequest, TResponse>(
        methodPath: string,
        request: TRequest
    ): Promise<TResponse> {
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
    private ensureWASMLoaded(): void {
        if (!this.isReady()) {
            throw new Error('WASM module not loaded. Call loadWasm() and waitUntilReady() first.');
        }
    }

    /**
     * Load the WASM module implementation
     */
    private async loadWASMModule(wasmPath: string): Promise<void> {
        console.log(`Loading ${this.config.moduleName} WASM module...`);

        // Check if WASM is already loaded (for testing environments) 
        if (this.checkIfPreLoaded()) {
            console.log('WASM module already loaded (pre-loaded in test environment)');
            return;
        }

        // Load Go's WASM support
        if (!(window as any).Go) {
            const script = document.createElement('script');
            script.src = '/wasm_exec.js';
            document.head.appendChild(script);

            await new Promise<void>((resolve, reject) => {
                script.onload = () => resolve();
                script.onerror = () => reject(new Error('Failed to load wasm_exec.js'));
            });
        }

        // Initialize Go WASM runtime
        const go = new (window as any).Go();
        const wasmModule = await WebAssembly.instantiateStreaming(
            fetch(wasmPath),
            go.importObject
        );

        // Run the WASM module
        go.run(wasmModule.instance);

        // Start browser service manager
        if (this.browserServiceManager) {
            this.browserServiceManager.setWasmModule(window);
            this.browserServiceManager.startProcessing();
        }

        // Verify WASM APIs are available
        this.verifyWASMLoaded();

        console.log(`${this.config.moduleName} WASM module loaded successfully`);
    }

    /**
     * Check if WASM is pre-loaded (for testing)
     */
    private checkIfPreLoaded(): boolean {
        switch (this.config.apiStructure) {
            case 'namespaced':
                if ((window as any)[this.config.jsNamespace]) {
                    this.wasm = (window as any)[this.config.jsNamespace];
                    return true;
                }
                return false;

            case 'flat':
                // For flat structure, we need to check for any method existence
                // This is a simplified check - in reality we'd check for a known method
                if ((window as any)[this.config.jsNamespace + 'LoadUserData']) {
                    this.wasm = window as any;
                    return true;
                }
                return false;

            case 'service_based':
                if ((window as any).services) {
                    this.wasm = (window as any).services;
                    return true;
                }
                return false;

            default:
                return false;
        }
    }

    /**
     * Verify WASM APIs are available after loading
     */
    private verifyWASMLoaded(): void {
        switch (this.config.apiStructure) {
            case 'namespaced':
                if (!(window as any)[this.config.jsNamespace]) {
                    throw new Error('WASM APIs not found - module may not have loaded correctly');
                }
                this.wasm = (window as any)[this.config.jsNamespace];
                break;

            case 'flat':
                // For flat structure, check for a known method
                if (!(window as any)[this.config.jsNamespace + 'LoadUserData']) {
                    throw new Error('WASM APIs not found - module may not have loaded correctly');
                }
                this.wasm = window as any;
                break;

            case 'service_based':
                if (!(window as any).services) {
                    throw new Error('WASM APIs not found - module may not have loaded correctly');
                }
                this.wasm = (window as any).services;
                break;

            default:
                throw new Error(`Unsupported API structure: ${this.config.apiStructure}`);
        }
    }
}
