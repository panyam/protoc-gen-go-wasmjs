/**
 * WASM Response interface for all service calls
 */
interface WASMResponse<T = any> {
    success: boolean;
    message: string;
    data: T;
}
/**
 * Error class for WASM-specific errors
 */
declare class WasmError extends Error {
    readonly methodPath?: string | undefined;
    constructor(message: string, methodPath?: string | undefined);
}

export { type WASMResponse, WasmError };
