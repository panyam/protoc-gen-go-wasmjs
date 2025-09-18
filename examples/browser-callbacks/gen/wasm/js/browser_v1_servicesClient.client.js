var __defProp = Object.defineProperty;
var __defNormalProp = (obj, key, value) => key in obj ? __defProp(obj, key, { enumerable: true, configurable: true, writable: true, value }) : obj[key] = value;
var __publicField = (obj, key, value) => __defNormalProp(obj, typeof key !== "symbol" ? key + "" : key, value);
import { BrowserServiceManager } from "./browserServiceManager.js";
class WasmError extends Error {
  constructor(message, methodPath) {
    super(message);
    this.methodPath = methodPath;
    this.name = "WasmError";
  }
}
class Browser_v1_servicesClient {
  // Service-specific clients
  constructor() {
    __publicField(this, "wasm");
    __publicField(this, "wasmLoadPromise", null);
    __publicField(this, "browserServiceManager", null);
    this.browserServiceManager = new BrowserServiceManager();
  }
  /**
   * Register a browser service implementation
   * Can be used to register browser services from any package
   */
  registerBrowserService(name, implementation) {
    if (!this.browserServiceManager) {
      throw new Error("Browser service manager not initialized");
    }
    this.browserServiceManager.registerService(name, implementation);
  }
  /**
   * Load the WASM module asynchronously
   */
  async loadWasm(wasmPath = "./browser_v1_services.wasm") {
    if (this.wasmLoadPromise) {
      return this.wasmLoadPromise;
    }
    this.wasmLoadPromise = this.loadWASMModule(wasmPath);
    return this.wasmLoadPromise;
  }
  /**
   * Check if WASM is ready for operations
   */
  isReady() {
    return this.wasm !== null && this.wasm !== void 0;
  }
  /**
   * Wait for WASM to be ready (use during initialization)
   */
  async waitUntilReady() {
    if (!this.wasmLoadPromise) {
      throw new Error("WASM loading not started. Call loadWasm() first.");
    }
    await this.wasmLoadPromise;
  }
  /**
   * Internal method to call WASM functions with JSON conversion
   */
  callMethod(methodPath, request) {
    this.ensureWASMLoaded();
    try {
      const jsonReq = JSON.parse(JSON.stringify(request));
      const wasmMethod = this.getWasmMethod(methodPath);
      const wasmResponse = wasmMethod(JSON.stringify(jsonReq));
      if (!wasmResponse.success) {
        throw new WasmError(wasmResponse.message, methodPath);
      }
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
  callMethodWithCallback(methodPath, request, callback) {
    this.ensureWASMLoaded();
    try {
      const jsonReq = JSON.parse(JSON.stringify(request));
      const wasmMethod = this.getWasmMethod(methodPath);
      const wasmResponse = wasmMethod(JSON.stringify(jsonReq), callback);
      if (!wasmResponse.success) {
        throw new WasmError(wasmResponse.message, methodPath);
      }
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
  callStreamingMethod(methodPath, request, callback) {
    this.ensureWASMLoaded();
    try {
      const jsonReq = JSON.parse(JSON.stringify(request));
      const wasmMethod = this.getWasmMethod(methodPath);
      const wrappedCallback = (responseStr, error, done) => {
        let response = null;
        if (responseStr && !error) {
          try {
            response = JSON.parse(responseStr);
          } catch (e) {
            response = responseStr;
          }
        }
        return callback(response, error, done);
      };
      const wasmResponse = wasmMethod(JSON.stringify(jsonReq), wrappedCallback);
      if (!wasmResponse.success) {
        throw new WasmError(wasmResponse.message, methodPath);
      }
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
   * Load the WASM module implementation
   */
  async loadWASMModule(wasmPath) {
    console.log("Loading browser_v1_services WASM module...");
    if (window.browser_v1) {
      console.log("WASM module already loaded (pre-loaded in test environment)");
      this.wasm = window.browser_v1;
      return;
    }
    if (!window.Go) {
      const script = document.createElement("script");
      script.src = "/static/wasm/wasm_exec.js";
      document.head.appendChild(script);
      await new Promise((resolve, reject) => {
        script.onload = () => resolve();
        script.onerror = () => reject(new Error("Failed to load wasm_exec.js"));
      });
    }
    const go = new window.Go();
    const wasmModule = await WebAssembly.instantiateStreaming(
      fetch(wasmPath),
      go.importObject
    );
    go.run(wasmModule.instance);
    if (this.browserServiceManager) {
      this.browserServiceManager.setWasmModule(window);
      this.browserServiceManager.startProcessing();
    }
    if (!window.browser_v1) {
      throw new Error("WASM APIs not found - module may not have loaded correctly");
    }
    this.wasm = window.browser_v1;
    console.log("browser_v1_services WASM module loaded successfully");
  }
  /**
   * Ensure WASM module is loaded (synchronous version for service calls)
   */
  ensureWASMLoaded() {
    if (!this.isReady()) {
      throw new Error("WASM module not loaded. Call loadWasm() and waitUntilReady() first.");
    }
  }
  /**
   * Get WASM method function by path
   */
  getWasmMethod(methodPath) {
    const parts = methodPath.split(".");
    let current = this.wasm;
    for (const part of parts) {
      current = current[part];
      if (!current) {
        throw new Error(`Method not found: ${methodPath}`);
      }
    }
    return current;
  }
}
var browser_v1_servicesClient_client_default = Browser_v1_servicesClient;
export {
  BrowserServiceManager,
  Browser_v1_servicesClient,
  WasmError,
  browser_v1_servicesClient_client_default as default
};
