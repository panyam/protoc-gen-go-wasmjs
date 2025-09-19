/**
 * Browser Service Manager
 * Handles FIFO processing of browser service calls from WASM
 * This is a shared component used by all WASM clients that interact with browser services
 */
declare class BrowserServiceManager {
    private processing;
    private serviceImplementations;
    private wasmModule;
    constructor();
    /**
     * Register a browser service implementation
     */
    registerService(name: string, implementation: any): void;
    /**
     * Set the WASM module reference
     */
    setWasmModule(wasmModule: any): void;
    /**
     * Start processing browser service calls
     */
    startProcessing(): Promise<void>;
    /**
     * Process a single browser service call asynchronously
     */
    private processCall;
    /**
     * Stop processing browser service calls
     */
    stopProcessing(): void;
    /**
     * Get the next browser call from WASM
     */
    private getNextBrowserCall;
    /**
     * Deliver a response back to WASM (called internally)
     */
    private deliverBrowserResponse;
}

export { BrowserServiceManager };
