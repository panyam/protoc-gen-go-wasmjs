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

import { describe, it, expect, beforeEach } from 'vitest';
import { WASMServiceClient, BrowserServiceManager, WasmError } from '../index.js';

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

            await client.callMethodWithCallback('testMethod', { input: 'test' }, (response, error) => {
                callbackResult = response;
                callbackError = error;
            });

            expect(callbackResult).toEqual({ result: 'test-data' });
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
                expect(error.message).toBe('Mock loading failure');
            }
        });
    });
});

describe('BrowserServiceManager Tests', () => {
    let manager: BrowserServiceManager;

    beforeEach(() => {
        manager = new BrowserServiceManager();
    });

    describe('Service Registration', () => {
        it('should register services correctly', () => {
            const mockService = {
                testMethod: async () => ({ success: true })
            };

            expect(() => {
                manager.registerService('TestService', mockService);
            }).not.toThrow();
        });

        it('should handle multiple service registrations', () => {
            const service1 = { method1: async () => ({}) };
            const service2 = { method2: async () => ({}) };

            manager.registerService('Service1', service1);
            manager.registerService('Service2', service2);

            // Should not throw - multiple registrations should be allowed
            expect(true).toBe(true);
        });
    });

    describe('WASM Module Integration', () => {
        it('should set WASM module reference', () => {
            const mockWasmModule = {
                __wasmGetNextBrowserCall: () => null,
                __wasmDeliverBrowserResponse: () => true
            };

            expect(() => {
                manager.setWasmModule(mockWasmModule as any);
            }).not.toThrow();
        });

        it('should handle missing WASM functions gracefully', () => {
            const incompleteMockModule = {};
            
            manager.setWasmModule(incompleteMockModule as any);
            
            // getNextBrowserCall should return null if WASM function missing
            const call = manager['getNextBrowserCall']();
            expect(call).toBeNull();
            
            // deliverBrowserResponse should return false if WASM function missing
            const delivered = manager['deliverBrowserResponse']('test-id', '{}', null);
            expect(delivered).toBe(false);
        });
    });

    describe('Processing Loop Control', () => {
        it('should start and stop processing', () => {
            expect(manager['processing']).toBe(false);
            
            // Start processing (will run until stopped)
            manager.startProcessing();
            expect(manager['processing']).toBe(true);
            
            // Stop processing
            manager.stopProcessing();
            expect(manager['processing']).toBe(false);
        });
    });
});
