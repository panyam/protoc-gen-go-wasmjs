// src/client/types.ts
var WasmError = class extends Error {
  constructor(message, methodPath) {
    super(message);
    this.methodPath = methodPath;
    this.name = "WasmError";
  }
};

export { WasmError };
//# sourceMappingURL=index.mjs.map
//# sourceMappingURL=index.mjs.map