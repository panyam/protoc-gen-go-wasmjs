'use strict';

// src/client/types.ts
var WasmError = class extends Error {
  constructor(message, methodPath) {
    super(message);
    this.methodPath = methodPath;
    this.name = "WasmError";
  }
};

exports.WasmError = WasmError;
//# sourceMappingURL=index.js.map
//# sourceMappingURL=index.js.map