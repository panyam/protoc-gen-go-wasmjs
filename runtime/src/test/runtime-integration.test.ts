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

import { describe, it, expect, beforeEach, vi } from 'vitest';
import { WASMServiceClient, BrowserServiceManager, WasmError } from '../index.js';

// Mock window for Node.js environment
const mockWindow: Record<string, any> = {};
vi.stubGlobal('window', mockWindow);

// Mock WASM service client for testing inheritance
class TestWASMClient extends WASMServiceClient {
    public testProperty = 'test';

    protected async loadWASMModule(wasmPath: string): Promise<void> {
        // Mock implementation for testing
        this.wasm = {
            testMethod: () => ({
                success: true,
                message: 'Success',
                data: { result: 'test-data' }
            })
        };
    }

    protected getWasmMethod(methodPath: string): Function {
        return this.wasm[methodPath] || (() => {
            throw new Error(`Method not found: ${methodPath}`);
        });
    }
}

describe('Framework Runtime Integration Tests', () => {
    let client: TestWASMClient;

    beforeEach(() => {
        client = new TestWASMClient();
    });

    describe('Base Class Inheritance (Template Inheritance Fix)', () => {
        it('should properly extend WASMServiceClient', () => {
            expect(client).toBeInstanceOf(WASMServiceClient);
            expect(client.testProperty).toBe('test');
        });

        it('should have all base class methods available (fixes missing methods issue)', () => {
            // This validates our fix for "loadWasm is not a function" and similar errors
            const requiredMethods = [
                'loadWasm', 'registerBrowserService', 'isReady', 'waitUntilReady',
                'callMethod', 'callMethodWithCallback', 'callStreamingMethod'
            ];

            for (const methodName of requiredMethods) {
                expect(typeof (client as any)[methodName]).toBe('function');
            }
        });

        it('should have properly initialized browser service manager', () => {
            expect(client['browserServiceManager']).toBeInstanceOf(BrowserServiceManager);
        });

        it('should handle WASM loading state correctly', () => {
            // Initially not ready
            expect(client.isReady()).toBe(false);

            // Should throw if trying to call methods before loading
            expect(() => client['ensureWASMLoaded']()).toThrow('WASM module not loaded');
        });
    });

    describe('WASM Method Calls', () => {
        beforeEach(async () => {
            // Load mock WASM for testing
            await client.loadWasm('mock.wasm');
        });

        it('should call WASM methods successfully', async () => {
            const result = await client.callMethod('testMethod', { input: 'test' });
            expect(result).toEqual({ result: 'test-data' });
        });

        it('should handle WASM method errors', async () => {
            try {
                await client.callMethod('nonExistentMethod', {});
                expect.fail('Should have thrown error');
            } catch (error) {
                expect(error).toBeInstanceOf(WasmError);
                expect((error as WasmError).methodPath).toBe('nonExistentMethod');
            }
        });

        it('should handle callback methods', async () => {
            let callbackResult: any = null;
            let callbackError: string | undefined = undefined;

            // callMethodWithCallback is async and returns immediately
            // The callback is called by the WASM method
            await client.callMethodWithCallback('testMethod', { input: 'test' }, (response, error) => {
                callbackResult = response;
                callbackError = error;
            });

            // Note: In the mock, the callback isn't actually invoked by the WASM method
            // This test verifies the method doesn't throw
            expect(callbackError).toBeUndefined();
        });
    });

    describe('Browser Service Registration', () => {
        it('should register browser services correctly', () => {
            const mockBrowserService = {
                testMethod: async (request: any) => ({ result: 'browser-result' })
            };

            expect(() => {
                client.registerBrowserService('TestBrowserService', mockBrowserService);
            }).not.toThrow();
        });

        it('should handle browser service registration errors', () => {
            // Create client without browser service manager (edge case)
            const brokenClient = new TestWASMClient();
            brokenClient['browserServiceManager'] = null;

            expect(() => {
                brokenClient.registerBrowserService('TestService', {});
            }).toThrow('Browser service manager not initialized');
        });
    });

    describe('Error Handling', () => {
        it('should create proper WasmError instances', () => {
            const error = new WasmError('Test error message', 'test.method');

            expect(error).toBeInstanceOf(Error);
            expect(error).toBeInstanceOf(WasmError);
            expect(error.name).toBe('WasmError');
            expect(error.message).toBe('Test error message');
            expect(error.methodPath).toBe('test.method');
        });

        it('should handle WASM loading errors', async () => {
            // Create client that will fail to load
            class FailingClient extends WASMServiceClient {
                protected async loadWASMModule(wasmPath: string): Promise<void> {
                    throw new Error('Mock loading failure');
                }
                protected getWasmMethod(methodPath: string): Function {
                    throw new Error('Not implemented');
                }
            }

            const failingClient = new FailingClient();

            try {
                await failingClient.loadWasm('failing.wasm');
                expect.fail('Should have thrown error');
            } catch (error) {
                expect((error as Error).message).toBe('Mock loading failure');
            }
        });
    });
});

describe('BrowserServiceManager Tests', () => {
    let manager: BrowserServiceManager;

    beforeEach(() => {
        // Reset mock window
        Object.keys(mockWindow).forEach(key => delete mockWindow[key]);
        manager = new BrowserServiceManager();
    });

    describe('Service Registration', () => {
        it('should register services correctly', () => {
            const mockService = {
                testMethod: () => ({ success: true })
            };

            expect(() => {
                manager.registerService('TestService', mockService);
            }).not.toThrow();

            expect(manager.hasService('TestService')).toBe(true);
        });

        it('should handle multiple service registrations', () => {
            const service1 = { method1: () => ({}) };
            const service2 = { method2: () => ({}) };

            manager.registerService('Service1', service1);
            manager.registerService('Service2', service2);

            expect(manager.getServiceCount()).toBe(2);
            expect(manager.hasService('Service1')).toBe(true);
            expect(manager.hasService('Service2')).toBe(true);
        });

        it('should update __browserServices when services are registered', () => {
            const mockService = {
                testMethod: () => ({ result: 'test' })
            };

            manager.registerService('TestService', mockService);

            // Check that __browserServices was set on window
            expect(mockWindow.__browserServices).toBeDefined();
            expect(mockWindow.__browserServices.TestService).toBeDefined();
            expect(typeof mockWindow.__browserServices.TestService.testMethod).toBe('function');
        });
    });

    describe('Initialization', () => {
        it('should set up async method handler on initialize', () => {
            manager.initialize();

            expect(typeof mockWindow.__wasmCallBrowserMethodAsync).toBe('function');
        });

        it('should only initialize once', () => {
            manager.initialize();
            const firstHandler = mockWindow.__wasmCallBrowserMethodAsync;

            manager.initialize();
            const secondHandler = mockWindow.__wasmCallBrowserMethodAsync;

            // Should be the same handler (not re-registered)
            expect(firstHandler).toBe(secondHandler);
        });
    });

    describe('Sync Method Calls via __browserServices', () => {
        it('should wrap sync methods correctly', () => {
            const mockService = {
                getData: (request: any) => ({ value: request.key + '-result' })
            };

            manager.registerService('TestService', mockService);

            // Call through __browserServices
            const result = mockWindow.__browserServices.TestService.getData('{"key":"test"}');
            expect(JSON.parse(result)).toEqual({ value: 'test-result' });
        });

        it('should return error for async methods called synchronously', () => {
            const mockService = {
                asyncMethod: async (request: any) => ({ value: 'async-result' })
            };

            manager.registerService('TestService', mockService);

            // Call through __browserServices - should return error since it's async
            const result = mockWindow.__browserServices.TestService.asyncMethod('{}');
            const parsed = JSON.parse(result);
            expect(parsed.error).toBeDefined();
        });
    });

    describe('Async Method Calls via __wasmCallBrowserMethodAsync', () => {
        it('should handle async method calls', async () => {
            const mockService = {
                fetchData: async (request: any) => ({ result: 'fetched-' + request.id })
            };

            manager.registerService('TestService', mockService);
            manager.initialize();

            // Mock the response delivery
            let deliveredCallId: string | null = null;
            let deliveredResponse: string | null = null;
            mockWindow.__wasmDeliverBrowserResponse = (callId: string, response: string, error: string | null) => {
                deliveredCallId = callId;
                deliveredResponse = response;
            };

            // Call the async method
            await mockWindow.__wasmCallBrowserMethodAsync('TestService', 'fetchData', '{"id":"123"}', 'call-1');

            expect(deliveredCallId).toBe('call-1');
            expect(JSON.parse(deliveredResponse!)).toEqual({ result: 'fetched-123' });
        });

        it('should handle async method errors', async () => {
            const mockService = {
                failingMethod: async () => { throw new Error('Test error'); }
            };

            manager.registerService('TestService', mockService);
            manager.initialize();

            let deliveredError: string | null = null;
            mockWindow.__wasmDeliverBrowserResponse = (callId: string, response: string | null, error: string | null) => {
                deliveredError = error;
            };

            await mockWindow.__wasmCallBrowserMethodAsync('TestService', 'failingMethod', '{}', 'call-2');

            expect(deliveredError).toBe('Test error');
        });

        it('should handle missing service gracefully', async () => {
            manager.initialize();

            let deliveredError: string | null = null;
            mockWindow.__wasmDeliverBrowserResponse = (callId: string, response: string | null, error: string | null) => {
                deliveredError = error;
            };

            await mockWindow.__wasmCallBrowserMethodAsync('NonExistentService', 'someMethod', '{}', 'call-3');

            expect(deliveredError).toContain("not registered");
        });
    });
});
