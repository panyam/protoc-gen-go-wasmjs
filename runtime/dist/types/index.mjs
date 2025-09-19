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

export { PatchOperation, PatchSource };
//# sourceMappingURL=index.mjs.map
//# sourceMappingURL=index.mjs.map