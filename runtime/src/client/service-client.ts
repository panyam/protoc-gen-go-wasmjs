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

import { WASMBundle } from './wasm-bundle.js';

/**
 * Base service client that references a shared WASM bundle
 * Lightweight facade for service-specific method calls
 */
export abstract class ServiceClient {
    protected bundle: WASMBundle;

    constructor(bundle: WASMBundle) {
        this.bundle = bundle;
    }

    /**
     * Check if the underlying WASM bundle is ready
     */
    public isReady(): boolean {
        return this.bundle.isReady();
    }

    /**
     * Wait for the underlying WASM bundle to be ready
     */
    public async waitUntilReady(): Promise<void> {
        return this.bundle.waitUntilReady();
    }

    /**
     * Call a synchronous WASM method
     */
    protected callMethod<TRequest, TResponse>(
        methodPath: string,
        request: TRequest
    ): Promise<TResponse> {
        return this.bundle.callMethod(methodPath, request);
    }

    /**
     * Call an asynchronous WASM method with callback
     */
    protected callMethodWithCallback<TRequest>(
        methodPath: string,
        request: TRequest,
        callback: (response: any, error?: string) => void
    ): Promise<void> {
        return this.bundle.callMethodWithCallback(methodPath, request, callback);
    }

    /**
     * Call a server streaming WASM method
     */
    protected callStreamingMethod<TRequest, TResponse>(
        methodPath: string,
        request: TRequest,
        callback: (response: TResponse | null, error: string | null, done: boolean) => boolean
    ): void {
        return this.bundle.callStreamingMethod(methodPath, request, callback);
    }
}
