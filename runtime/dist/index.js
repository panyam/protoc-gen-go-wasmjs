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

// src/schema/types.ts
var FieldType = /* @__PURE__ */ ((FieldType2) => {
  FieldType2["STRING"] = "string";
  FieldType2["NUMBER"] = "number";
  FieldType2["BOOLEAN"] = "boolean";
  FieldType2["MESSAGE"] = "message";
  FieldType2["REPEATED"] = "repeated";
  FieldType2["MAP"] = "map";
  FieldType2["ONEOF"] = "oneof";
  return FieldType2;
})(FieldType || {});

// src/client/types.ts
var WasmError = class extends Error {
  constructor(message, methodPath) {
    super(message);
    this.methodPath = methodPath;
    this.name = "WasmError";
  }
};

// src/types/patches.ts
var PatchOperation = /* @__PURE__ */ ((PatchOperation2) => {
  PatchOperation2["SET"] = "SET";
  PatchOperation2["INSERT_LIST"] = "INSERT_LIST";
  PatchOperation2["REMOVE_LIST"] = "REMOVE_LIST";
  PatchOperation2["MOVE_LIST"] = "MOVE_LIST";
  PatchOperation2["INSERT_MAP"] = "INSERT_MAP";
  PatchOperation2["REMOVE_MAP"] = "REMOVE_MAP";
  PatchOperation2["CLEAR_LIST"] = "CLEAR_LIST";
  PatchOperation2["CLEAR_MAP"] = "CLEAR_MAP";
  return PatchOperation2;
})(PatchOperation || {});
var PatchSource = /* @__PURE__ */ ((PatchSource2) => {
  PatchSource2["LOCAL"] = "LOCAL";
  PatchSource2["REMOTE"] = "REMOTE";
  PatchSource2["SERVER"] = "SERVER";
  PatchSource2["STORAGE"] = "STORAGE";
  return PatchSource2;
})(PatchSource || {});

exports.BrowserServiceManager = BrowserServiceManager;
exports.FieldType = FieldType;
exports.PatchOperation = PatchOperation;
exports.PatchSource = PatchSource;
exports.WasmError = WasmError;
//# sourceMappingURL=index.js.map
//# sourceMappingURL=index.js.map