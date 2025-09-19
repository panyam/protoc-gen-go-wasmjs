'use strict';

// src/browser/service-manager.ts
var BrowserServiceManager = class {
  constructor() {
    this.processing = false;
    this.serviceImplementations = /* @__PURE__ */ new Map();
  }
  /**
   * Register a browser service implementation
   */
  registerService(name, implementation) {
    this.serviceImplementations.set(name, implementation);
  }
  /**
   * Set the WASM module reference
   */
  setWasmModule(wasmModule) {
    this.wasmModule = wasmModule;
  }
  /**
   * Start processing browser service calls
   */
  async startProcessing() {
    if (this.processing) return;
    this.processing = true;
    while (this.processing) {
      const call = this.getNextBrowserCall();
      if (!call) {
        await new Promise((resolve) => setTimeout(resolve, 10));
        continue;
      }
      this.processCall(call);
    }
  }
  /**
   * Process a single browser service call asynchronously
   */
  async processCall(call) {
    try {
      const service = this.serviceImplementations.get(call.service);
      if (!service) {
        throw new Error(`No implementation registered for service: ${call.service}`);
      }
      const methodName = call.method.charAt(0).toLowerCase() + call.method.slice(1);
      const method = service[methodName];
      if (!method) {
        throw new Error(`Method ${methodName} not found on service ${call.service}`);
      }
      const request = JSON.parse(call.request);
      const response = await Promise.resolve(method.call(service, request));
      this.deliverBrowserResponse(call.id, JSON.stringify(response), null);
    } catch (error) {
      this.deliverBrowserResponse(call.id, null, error.message || String(error));
    }
  }
  /**
   * Stop processing browser service calls
   */
  stopProcessing() {
    this.processing = false;
  }
  /**
   * Get the next browser call from WASM
   */
  getNextBrowserCall() {
    if (typeof window.__wasmGetNextBrowserCall === "function") {
      return window.__wasmGetNextBrowserCall();
    }
    return null;
  }
  /**
   * Deliver a response back to WASM (called internally)
   */
  deliverBrowserResponse(callId, response, error) {
    if (!window.__wasmDeliverBrowserResponse) {
      return false;
    }
    return window.__wasmDeliverBrowserResponse(callId, response, error);
  }
};

exports.BrowserServiceManager = BrowserServiceManager;
//# sourceMappingURL=index.js.map
//# sourceMappingURL=index.js.map