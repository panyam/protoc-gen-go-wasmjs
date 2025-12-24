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

/**
 * Browser Service Manager
 *
 * Manages browser service implementations that can be called from WASM.
 *
 * There are two calling patterns:
 *
 * 1. Synchronous: WASM calls __browserServices.ServiceName.methodName(request)
 *    and receives the result directly. The JS function must return a value (not a Promise).
 *
 * 2. Asynchronous: WASM calls __wasmCallBrowserMethodAsync(service, method, request, callbackId)
 *    and JS calls __wasmDeliverBrowserResponse(callbackId, response, error) when done.
 *    This is for browser APIs that return Promises (fetch, IndexedDB, etc.)
 */
export class BrowserServiceManager {
    private serviceImplementations = new Map<string, any>();
    private initialized = false;

    constructor() {}

    /**
     * Register a browser service implementation.
     * The service object should have methods matching the proto service definition.
     *
     * Methods can be either:
     * - Synchronous: Return a value directly (for sync browser APIs)
     * - Asynchronous: Return a Promise (for async browser APIs like fetch)
     */
    registerService(name: string, implementation: any): void {
        this.serviceImplementations.set(name, implementation);
        this.updateBrowserServices();
    }

    /**
     * Initialize the browser service bridge.
     * This sets up the global functions that WASM uses to call browser services.
     */
    initialize(): void {
        if (this.initialized) return;
        this.initialized = true;

        // Set up the async method handler
        // WASM calls this for async browser methods (ones that return Promises)
        (window as any).__wasmCallBrowserMethodAsync = async (
            service: string,
            method: string,
            request: string,
            callbackId: string
        ) => {
            try {
                const impl = this.serviceImplementations.get(service);
                if (!impl) {
                    this.deliverResponse(callbackId, null, `Service '${service}' not registered`);
                    return;
                }

                // Convert method name to camelCase (e.g., "getLocalStorage" -> "getLocalStorage")
                const methodName = method.charAt(0).toLowerCase() + method.slice(1);
                const methodFunc = impl[methodName];
                if (!methodFunc) {
                    this.deliverResponse(callbackId, null, `Method '${methodName}' not found on service '${service}'`);
                    return;
                }

                // Parse request and call the method
                const requestObj = JSON.parse(request);
                const response = await Promise.resolve(methodFunc.call(impl, requestObj));

                // Deliver the response back to WASM
                this.deliverResponse(callbackId, JSON.stringify(response), null);
            } catch (error: any) {
                this.deliverResponse(callbackId, null, error.message || String(error));
            }
        };

        // Update the browser services object
        this.updateBrowserServices();
    }

    /**
     * Update the __browserServices global object with current registrations.
     * This is called automatically when services are registered.
     */
    private updateBrowserServices(): void {
        const browserServices: Record<string, any> = {};

        for (const [name, impl] of this.serviceImplementations) {
            const serviceProxy: Record<string, any> = {};

            // Get all methods from the implementation
            // Check both own properties and prototype chain
            const methodNames = this.getAllMethodNames(impl);

            for (const methodName of methodNames) {
                // Create a wrapper that handles sync calls
                serviceProxy[methodName] = (request: string): string => {
                    try {
                        const requestObj = JSON.parse(request);
                        const result = impl[methodName](requestObj);

                        // For sync methods, result should not be a Promise
                        // If it is, this is a programming error - async methods should use
                        // __wasmCallBrowserMethodAsync instead
                        if (result instanceof Promise) {
                            console.warn(
                                `Method ${name}.${methodName} returned a Promise but was called synchronously. ` +
                                `Mark this method as async in the proto definition.`
                            );
                            return JSON.stringify({ error: 'Method is async but called synchronously' });
                        }

                        return JSON.stringify(result);
                    } catch (error: any) {
                        return JSON.stringify({ error: error.message || String(error) });
                    }
                };
            }

            browserServices[name] = serviceProxy;
        }

        (window as any).__browserServices = browserServices;
    }

    /**
     * Get all method names from an object, including prototype chain.
     */
    private getAllMethodNames(obj: any): string[] {
        const methods = new Set<string>();

        // Get own property names (for plain objects)
        for (const name of Object.keys(obj)) {
            if (typeof obj[name] === 'function') {
                methods.add(name);
            }
        }

        // Also check prototype chain (for class instances)
        let proto = Object.getPrototypeOf(obj);
        while (proto && proto !== Object.prototype) {
            for (const name of Object.getOwnPropertyNames(proto)) {
                if (name !== 'constructor' && typeof obj[name] === 'function') {
                    methods.add(name);
                }
            }
            proto = Object.getPrototypeOf(proto);
        }

        return Array.from(methods);
    }

    /**
     * Deliver an async response back to WASM.
     */
    private deliverResponse(callbackId: string, response: string | null, error: string | null): void {
        if (typeof (window as any).__wasmDeliverBrowserResponse === 'function') {
            (window as any).__wasmDeliverBrowserResponse(callbackId, response, error);
        } else {
            console.error('__wasmDeliverBrowserResponse not available - WASM not initialized?');
        }
    }

    /**
     * Check if a service is registered.
     */
    hasService(name: string): boolean {
        return this.serviceImplementations.has(name);
    }

    /**
     * Get the number of registered services.
     */
    getServiceCount(): number {
        return this.serviceImplementations.size;
    }
}
