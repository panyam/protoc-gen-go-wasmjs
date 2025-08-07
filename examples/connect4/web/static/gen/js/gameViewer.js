"use strict";
var Connect4GameViewer = (() => {
  var __defProp = Object.defineProperty;
  var __getOwnPropDesc = Object.getOwnPropertyDescriptor;
  var __getOwnPropNames = Object.getOwnPropertyNames;
  var __getOwnPropSymbols = Object.getOwnPropertySymbols;
  var __hasOwnProp = Object.prototype.hasOwnProperty;
  var __propIsEnum = Object.prototype.propertyIsEnumerable;
  var __defNormalProp = (obj, key, value) => key in obj ? __defProp(obj, key, { enumerable: true, configurable: true, writable: true, value }) : obj[key] = value;
  var __spreadValues = (a, b) => {
    for (var prop in b || (b = {}))
      if (__hasOwnProp.call(b, prop))
        __defNormalProp(a, prop, b[prop]);
    if (__getOwnPropSymbols)
      for (var prop of __getOwnPropSymbols(b)) {
        if (__propIsEnum.call(b, prop))
          __defNormalProp(a, prop, b[prop]);
      }
    return a;
  };
  var __export = (target, all) => {
    for (var name in all)
      __defProp(target, name, { get: all[name], enumerable: true });
  };
  var __copyProps = (to, from, except, desc) => {
    if (from && typeof from === "object" || typeof from === "function") {
      for (let key of __getOwnPropNames(from))
        if (!__hasOwnProp.call(to, key) && key !== except)
          __defProp(to, key, { get: () => from[key], enumerable: !(desc = __getOwnPropDesc(from, key)) || desc.enumerable });
    }
    return to;
  };
  var __toCommonJS = (mod) => __copyProps(__defProp({}, "__esModule", { value: true }), mod);
  var __async = (__this, __arguments, generator) => {
    return new Promise((resolve, reject) => {
      var fulfilled = (value) => {
        try {
          step(generator.next(value));
        } catch (e) {
          reject(e);
        }
      };
      var rejected = (value) => {
        try {
          step(generator.throw(value));
        } catch (e) {
          reject(e);
        }
      };
      var step = (x) => x.done ? resolve(x.value) : Promise.resolve(x.value).then(fulfilled, rejected);
      step((generator = generator.apply(__this, __arguments)).next());
    });
  };

  // src/gameViewer.ts
  var gameViewer_exports = {};
  __export(gameViewer_exports, {
    default: () => gameViewer_default
  });

  // gen/wasmts/multiplayer_connect4Client.client.ts
  var WasmError = class extends Error {
    constructor(message, methodPath) {
      super(message);
      this.methodPath = methodPath;
      this.name = "WasmError";
    }
  };
  var Multiplayer_connect4Client = class {
    constructor() {
      this.wasmLoadPromise = null;
      this.connect4Service = new Connect4ServiceClientImpl(this);
    }
    /**
     * Load the WASM module asynchronously
     */
    loadWasm(wasmPath = "./multiplayer_connect4.wasm") {
      return __async(this, null, function* () {
        if (this.wasmLoadPromise) {
          return this.wasmLoadPromise;
        }
        this.wasmLoadPromise = this.loadWASMModule(wasmPath);
        return this.wasmLoadPromise;
      });
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
    waitUntilReady() {
      return __async(this, null, function* () {
        if (!this.wasmLoadPromise) {
          throw new Error("WASM loading not started. Call loadWasm() first.");
        }
        yield this.wasmLoadPromise;
      });
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
     * Load the WASM module implementation
     */
    loadWASMModule(wasmPath) {
      return __async(this, null, function* () {
        console.log("Loading multiplayer_connect4 WASM module...");
        if (window.multiplayerConnect4) {
          console.log("WASM module already loaded (pre-loaded in test environment)");
          this.wasm = window.multiplayerConnect4;
          return;
        }
        if (!window.Go) {
          const script = document.createElement("script");
          script.src = "/static/wasm/wasm_exec.js";
          document.head.appendChild(script);
          yield new Promise((resolve, reject) => {
            script.onload = () => resolve();
            script.onerror = () => reject(new Error("Failed to load wasm_exec.js"));
          });
        }
        const go = new window.Go();
        const wasmModule = yield WebAssembly.instantiateStreaming(
          fetch(wasmPath),
          go.importObject
        );
        go.run(wasmModule.instance);
        if (!window.multiplayerConnect4) {
          throw new Error("WASM APIs not found - module may not have loaded correctly");
        }
        this.wasm = window.multiplayerConnect4;
        console.log("multiplayer_connect4 WASM module loaded successfully");
      });
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
  };
  var Connect4ServiceClientImpl = class {
    constructor(parent) {
      this.parent = parent;
    }
    getGame(request) {
      return __async(this, null, function* () {
        return this.parent.callMethod("connect4Service.getGame", request);
      });
    }
    dropPiece(request) {
      return __async(this, null, function* () {
        return this.parent.callMethod("connect4Service.dropPiece", request);
      });
    }
    joinGame(request) {
      return __async(this, null, function* () {
        return this.parent.callMethod("connect4Service.joinGame", request);
      });
    }
    createGame(request) {
      return __async(this, null, function* () {
        return this.parent.callMethod("connect4Service.createGame", request);
      });
    }
  };
  var multiplayer_connect4Client_client_default = Multiplayer_connect4Client;

  // gen/wasmts/wasmjs/v1/schemas.ts
  var StatefulOptionsSchema = {
    name: "StatefulOptions",
    fields: [
      {
        name: "enabled",
        type: "boolean" /* BOOLEAN */,
        id: 1
      },
      {
        name: "stateMessageType",
        type: "string" /* STRING */,
        id: 2
      },
      {
        name: "conflictResolution",
        type: "string" /* STRING */,
        id: 3
      }
    ]
  };
  var StatefulMethodOptionsSchema = {
    name: "StatefulMethodOptions",
    fields: [
      {
        name: "returnsPatches",
        type: "boolean" /* BOOLEAN */,
        id: 1
      },
      {
        name: "broadcasts",
        type: "boolean" /* BOOLEAN */,
        id: 2
      }
    ]
  };
  var MessagePatchSchema = {
    name: "MessagePatch",
    fields: [
      {
        name: "operation",
        type: "string" /* STRING */,
        id: 1
      },
      {
        name: "fieldPath",
        type: "string" /* STRING */,
        id: 2
      },
      {
        name: "valueJson",
        type: "string" /* STRING */,
        id: 3
      },
      {
        name: "index",
        type: "number" /* NUMBER */,
        id: 4
      },
      {
        name: "key",
        type: "string" /* STRING */,
        id: 5
      },
      {
        name: "oldIndex",
        type: "number" /* NUMBER */,
        id: 6
      },
      {
        name: "changeNumber",
        type: "number" /* NUMBER */,
        id: 7
      },
      {
        name: "timestamp",
        type: "number" /* NUMBER */,
        id: 8
      },
      {
        name: "userId",
        type: "string" /* STRING */,
        id: 9
      },
      {
        name: "transactionId",
        type: "string" /* STRING */,
        id: 10
      }
    ]
  };
  var PatchBatchSchema = {
    name: "PatchBatch",
    fields: [
      {
        name: "messageType",
        type: "string" /* STRING */,
        id: 1
      },
      {
        name: "entityId",
        type: "string" /* STRING */,
        id: 2
      },
      {
        name: "patches",
        type: "message" /* MESSAGE */,
        id: 3,
        messageType: "wasmjs.v1.MessagePatch",
        repeated: true
      },
      {
        name: "changeNumber",
        type: "number" /* NUMBER */,
        id: 4
      },
      {
        name: "source",
        type: "string" /* STRING */,
        id: 5
      },
      {
        name: "metadata",
        type: "message" /* MESSAGE */,
        id: 6,
        messageType: "wasmjs.v1.MetadataEntry"
      }
    ]
  };
  var PatchResponseSchema = {
    name: "PatchResponse",
    fields: [
      {
        name: "patchBatches",
        type: "message" /* MESSAGE */,
        id: 1,
        messageType: "wasmjs.v1.PatchBatch",
        repeated: true
      },
      {
        name: "success",
        type: "boolean" /* BOOLEAN */,
        id: 2
      },
      {
        name: "errorMessage",
        type: "string" /* STRING */,
        id: 3
      },
      {
        name: "newChangeNumber",
        type: "number" /* NUMBER */,
        id: 4
      }
    ]
  };
  var WasmjsV1SchemaRegistry = {
    "wasmjs.v1.StatefulOptions": StatefulOptionsSchema,
    "wasmjs.v1.StatefulMethodOptions": StatefulMethodOptionsSchema,
    "wasmjs.v1.MessagePatch": MessagePatchSchema,
    "wasmjs.v1.PatchBatch": PatchBatchSchema,
    "wasmjs.v1.PatchResponse": PatchResponseSchema
  };

  // gen/wasmts/wasmjs/v1/deserializer.ts
  var DEFAULT_FACTORY = new WasmjsV1Factory();
  var WasmjsV1Deserializer = class _WasmjsV1Deserializer {
    constructor(schemaRegistry = WasmjsV1SchemaRegistry, factory = DEFAULT_FACTORY) {
      this.schemaRegistry = schemaRegistry;
      this.factory = factory;
    }
    /**
     * Deserialize an object using schema information
     * @param instance The target instance to populate
     * @param data The source data to deserialize from
     * @param messageType The fully qualified message type (e.g., "library.v1.Book")
     * @returns The populated instance
     */
    deserialize(instance, data, messageType) {
      if (!data || typeof data !== "object") {
        return instance;
      }
      const schema = this.schemaRegistry[messageType];
      if (!schema) {
        return this.fallbackDeserialize(instance, data);
      }
      for (const fieldSchema of schema.fields) {
        const fieldValue = data[fieldSchema.name];
        if (fieldValue === null || fieldValue === void 0) {
          continue;
        }
        this.deserializeField(instance, fieldSchema, fieldValue);
      }
      return instance;
    }
    /**
     * Deserialize a single field based on its schema
     */
    deserializeField(instance, fieldSchema, fieldValue) {
      const fieldName = fieldSchema.name;
      switch (fieldSchema.type) {
        case "string" /* STRING */:
        case "number" /* NUMBER */:
        case "boolean" /* BOOLEAN */:
          instance[fieldName] = fieldValue;
          break;
        case "message" /* MESSAGE */:
          if (fieldSchema.repeated) {
            instance[fieldName] = this.deserializeMessageArray(
              fieldValue,
              fieldSchema.messageType,
              instance,
              fieldName
            );
          } else {
            instance[fieldName] = this.deserializeMessageField(
              fieldValue,
              fieldSchema.messageType,
              instance,
              fieldName
            );
          }
          break;
        case "repeated" /* REPEATED */:
          if (Array.isArray(fieldValue)) {
            instance[fieldName] = [...fieldValue];
          }
          break;
        case "oneof" /* ONEOF */:
          instance[fieldName] = fieldValue;
          break;
        case "map" /* MAP */:
          instance[fieldName] = __spreadValues({}, fieldValue);
          break;
        default:
          instance[fieldName] = fieldValue;
          break;
      }
    }
    /**
     * Deserialize a single message field
     */
    deserializeMessageField(fieldValue, messageType, parent, attributeName) {
      let factoryMethod;
      if (this.factory.getFactoryMethod) {
        factoryMethod = this.factory.getFactoryMethod(messageType);
      } else {
        const factoryMethodName = this.getFactoryMethodName(messageType);
        factoryMethod = this.factory[factoryMethodName];
      }
      if (factoryMethod) {
        const result = factoryMethod(parent, attributeName, void 0, fieldValue);
        if (result.fullyLoaded) {
          return result.instance;
        } else {
          return this.deserialize(result.instance, fieldValue, messageType);
        }
      }
      return this.fallbackDeserialize({}, fieldValue);
    }
    /**
     * Deserialize an array of message objects
     */
    deserializeMessageArray(fieldValue, messageType, parent, attributeName) {
      if (!Array.isArray(fieldValue)) {
        return [];
      }
      let factoryMethod;
      if (this.factory.getFactoryMethod) {
        factoryMethod = this.factory.getFactoryMethod(messageType);
      } else {
        const factoryMethodName = this.getFactoryMethodName(messageType);
        factoryMethod = this.factory[factoryMethodName];
      }
      return fieldValue.map((item, index) => {
        if (factoryMethod) {
          const result = factoryMethod(parent, attributeName, index, item);
          if (result.fullyLoaded) {
            return result.instance;
          } else {
            return this.deserialize(result.instance, item, messageType);
          }
        }
        return this.fallbackDeserialize({}, item);
      });
    }
    /**
     * Convert message type to factory method name
     * "library.v1.Book" -> "newBook"
     */
    getFactoryMethodName(messageType) {
      const parts = messageType.split(".");
      const typeName = parts[parts.length - 1];
      return "new" + typeName;
    }
    /**
     * Fallback deserializer for when no schema is available
     */
    fallbackDeserialize(instance, data) {
      if (!data || typeof data !== "object") {
        return instance;
      }
      for (const [key, value] of Object.entries(data)) {
        if (value !== null && value !== void 0) {
          instance[key] = value;
        }
      }
      return instance;
    }
    /**
     * Create and deserialize a new instance of a message type
     */
    createAndDeserialize(messageType, data) {
      let factoryMethod;
      if (this.factory.getFactoryMethod) {
        factoryMethod = this.factory.getFactoryMethod(messageType);
      } else {
        const factoryMethodName = this.getFactoryMethodName(messageType);
        factoryMethod = this.factory[factoryMethodName];
      }
      if (!factoryMethod) {
        throw new Error(`Could not find factory method to deserialize: ${messageType}`);
      }
      const result = factoryMethod(void 0, void 0, void 0, data);
      if (result.fullyLoaded) {
        return result.instance;
      } else {
        return this.deserialize(result.instance, data, messageType);
      }
    }
    /**
     * Static utility method to create and deserialize a message without needing a deserializer instance
     * @param messageType Fully qualified message type (use Class.MESSAGE_TYPE)
     * @param data Raw data to deserialize
     * @returns Deserialized instance or null if creation failed
     */
    static from(messageType, data) {
      const deserializer = new _WasmjsV1Deserializer();
      return deserializer.createAndDeserialize(messageType, data);
    }
  };

  // gen/wasmts/wasmjs/v1/models.ts
  var _StatefulOptions = class _StatefulOptions {
    constructor() {
      /** Whether stateful proxy generation is enabled */
      this.enabled = false;
      /** The fully qualified name of the message type that represents the state
      (e.g., "example.Game", "library.Map") */
      this.stateMessageType = "";
      /** Strategy for resolving conflicts when multiple changes occur */
      this.conflictResolution = 0;
    }
    /**
     * Create and deserialize an instance from raw data
     * @param data Raw data to deserialize
     * @returns Deserialized StatefulOptions instance or null if creation failed
     */
    static from(data) {
      return WasmjsV1Deserializer.from(_StatefulOptions.MESSAGE_TYPE, data);
    }
  };
  /**
   * Fully qualified message type for schema resolution
   */
  _StatefulOptions.MESSAGE_TYPE = "wasmjs.v1.StatefulOptions";
  var StatefulOptions = _StatefulOptions;
  var _StatefulMethodOptions = class _StatefulMethodOptions {
    constructor() {
      /** Whether this method returns patch operations instead of full objects */
      this.returnsPatches = false;
      /** Whether changes from this method should be broadcast to other clients */
      this.broadcasts = false;
    }
    /**
     * Create and deserialize an instance from raw data
     * @param data Raw data to deserialize
     * @returns Deserialized StatefulMethodOptions instance or null if creation failed
     */
    static from(data) {
      return WasmjsV1Deserializer.from(_StatefulMethodOptions.MESSAGE_TYPE, data);
    }
  };
  /**
   * Fully qualified message type for schema resolution
   */
  _StatefulMethodOptions.MESSAGE_TYPE = "wasmjs.v1.StatefulMethodOptions";
  var StatefulMethodOptions = _StatefulMethodOptions;
  var _MessagePatch = class _MessagePatch {
    constructor() {
      /** The type of operation to perform */
      this.operation = 0;
      /** Path to the field being modified (e.g., "players[2].name", "places['tile_123'].latitude") */
      this.fieldPath = "";
      /** The new value to set (for SET, INSERT_LIST, INSERT_MAP operations)
      Encoded as JSON for type flexibility */
      this.valueJson = "";
      /** Index for list operations (INSERT_LIST, REMOVE_LIST, MOVE_LIST) */
      this.index = 0;
      /** Map key for map operations (INSERT_MAP, REMOVE_MAP) */
      this.key = "";
      /** Source index for MOVE_LIST operations */
      this.oldIndex = 0;
      /** Monotonically increasing change number for ordering */
      this.changeNumber = 0;
      /** Timestamp when the change was created (microseconds since epoch) */
      this.timestamp = 0;
      /** Optional user ID who made the change (for conflict resolution) */
      this.userId = "";
      /** Optional transaction ID to group related patches */
      this.transactionId = "";
    }
    /**
     * Create and deserialize an instance from raw data
     * @param data Raw data to deserialize
     * @returns Deserialized MessagePatch instance or null if creation failed
     */
    static from(data) {
      return WasmjsV1Deserializer.from(_MessagePatch.MESSAGE_TYPE, data);
    }
  };
  /**
   * Fully qualified message type for schema resolution
   */
  _MessagePatch.MESSAGE_TYPE = "wasmjs.v1.MessagePatch";
  var MessagePatch = _MessagePatch;
  var _PatchBatch = class _PatchBatch {
    constructor() {
      /** The fully qualified protobuf message type (e.g., "example.Game") */
      this.messageType = "";
      /** The unique identifier of the entity being modified */
      this.entityId = "";
      /** List of patches to apply in order */
      this.patches = [];
      /** The highest change number in this batch */
      this.changeNumber = 0;
      /** Source of the changes */
      this.source = 0;
    }
    /**
     * Create and deserialize an instance from raw data
     * @param data Raw data to deserialize
     * @returns Deserialized PatchBatch instance or null if creation failed
     */
    static from(data) {
      return WasmjsV1Deserializer.from(_PatchBatch.MESSAGE_TYPE, data);
    }
  };
  /**
   * Fully qualified message type for schema resolution
   */
  _PatchBatch.MESSAGE_TYPE = "wasmjs.v1.PatchBatch";
  var PatchBatch = _PatchBatch;
  var _PatchResponse = class _PatchResponse {
    constructor() {
      /** The patches to apply */
      this.patchBatches = [];
      /** Success status */
      this.success = false;
      /** Error message if success is false */
      this.errorMessage = "";
      /** The new change number after applying these patches */
      this.newChangeNumber = 0;
    }
    /**
     * Create and deserialize an instance from raw data
     * @param data Raw data to deserialize
     * @returns Deserialized PatchResponse instance or null if creation failed
     */
    static from(data) {
      return WasmjsV1Deserializer.from(_PatchResponse.MESSAGE_TYPE, data);
    }
  };
  /**
   * Fully qualified message type for schema resolution
   */
  _PatchResponse.MESSAGE_TYPE = "wasmjs.v1.PatchResponse";
  var PatchResponse = _PatchResponse;

  // gen/wasmts/wasmjs/v1/factory.ts
  var WasmjsV1Factory = class {
    constructor() {
      /**
       * Enhanced factory method for StatefulOptions
       * @param parent Parent object containing this field
       * @param attributeName Field name in parent object
       * @param attributeKey Array index, map key, or union tag (for containers)
       * @param data Raw data to potentially populate from
       * @returns Factory result with instance and population status
       */
      this.newStatefulOptions = (parent, attributeName, attributeKey, data) => {
        const out = new StatefulOptions();
        return { instance: out, fullyLoaded: false };
      };
      /**
       * Enhanced factory method for StatefulMethodOptions
       * @param parent Parent object containing this field
       * @param attributeName Field name in parent object
       * @param attributeKey Array index, map key, or union tag (for containers)
       * @param data Raw data to potentially populate from
       * @returns Factory result with instance and population status
       */
      this.newStatefulMethodOptions = (parent, attributeName, attributeKey, data) => {
        const out = new StatefulMethodOptions();
        return { instance: out, fullyLoaded: false };
      };
      /**
       * Enhanced factory method for MessagePatch
       * @param parent Parent object containing this field
       * @param attributeName Field name in parent object
       * @param attributeKey Array index, map key, or union tag (for containers)
       * @param data Raw data to potentially populate from
       * @returns Factory result with instance and population status
       */
      this.newMessagePatch = (parent, attributeName, attributeKey, data) => {
        const out = new MessagePatch();
        return { instance: out, fullyLoaded: false };
      };
      /**
       * Enhanced factory method for PatchBatch
       * @param parent Parent object containing this field
       * @param attributeName Field name in parent object
       * @param attributeKey Array index, map key, or union tag (for containers)
       * @param data Raw data to potentially populate from
       * @returns Factory result with instance and population status
       */
      this.newPatchBatch = (parent, attributeName, attributeKey, data) => {
        const out = new PatchBatch();
        return { instance: out, fullyLoaded: false };
      };
      /**
       * Enhanced factory method for PatchResponse
       * @param parent Parent object containing this field
       * @param attributeName Field name in parent object
       * @param attributeKey Array index, map key, or union tag (for containers)
       * @param data Raw data to potentially populate from
       * @returns Factory result with instance and population status
       */
      this.newPatchResponse = (parent, attributeName, attributeKey, data) => {
        const out = new PatchResponse();
        return { instance: out, fullyLoaded: false };
      };
      /**
       * Factory method for converting protobuf Timestamp data to native Date
       * @param parent Parent object containing this field
       * @param attributeName Field name in parent object  
       * @param attributeKey Array index, map key, or union tag (for containers)
       * @param data Raw protobuf timestamp data
       * @returns Factory result with Date instance
       */
      this.newTimestamp = (parent, attributeName, attributeKey, data) => {
        if (!data) {
          return { instance: /* @__PURE__ */ new Date(), fullyLoaded: true };
        }
        let date;
        if (typeof data === "string") {
          date = new Date(data);
        } else if (data.seconds !== void 0) {
          const seconds = typeof data.seconds === "string" ? parseInt(data.seconds, 10) : data.seconds;
          const nanos = data.nanos || 0;
          date = new Date(seconds * 1e3 + Math.floor(nanos / 1e6));
        } else {
          date = /* @__PURE__ */ new Date();
        }
        return { instance: date, fullyLoaded: true };
      };
      /**
       * Factory method for converting protobuf FieldMask data to native string array
       * @param parent Parent object containing this field
       * @param attributeName Field name in parent object
       * @param attributeKey Array index, map key, or union tag (for containers)
       * @param data Raw protobuf field mask data
       * @returns Factory result with string array instance
       */
      this.newFieldMask = (parent, attributeName, attributeKey, data) => {
        if (!data) {
          return { instance: [], fullyLoaded: true };
        }
        let paths;
        if (Array.isArray(data)) {
          paths = data;
        } else if (data.paths && Array.isArray(data.paths)) {
          paths = data.paths;
        } else {
          paths = [];
        }
        return { instance: paths, fullyLoaded: true };
      };
    }
    /**
     * Get factory method for a fully qualified message type
     * Enables cross-package factory delegation
     */
    getFactoryMethod(messageType) {
      const parts = messageType.split(".");
      if (parts.length < 2) {
        return void 0;
      }
      const packageName = parts.slice(0, -1).join(".");
      const typeName = parts[parts.length - 1];
      const methodName = "new" + typeName;
      const currentPackage = "wasmjs.v1";
      if (packageName === currentPackage) {
        return this[methodName];
      }
      const externalFactory = this.externalTypeFactories()[messageType];
      if (externalFactory) {
        return externalFactory;
      }
      return void 0;
    }
    /**
     * Generic object deserializer that respects factory decisions
     */
    deserializeObject(instance, data) {
      if (!data || typeof data !== "object") return instance;
      for (const [key, value] of Object.entries(data)) {
        if (value !== null && value !== void 0) {
          instance[key] = value;
        }
      }
      return instance;
    }
    // External type conversion methods
    /**
     * Mapping of external types to their factory methods
     */
    externalTypeFactories() {
      return {
        "google.protobuf.Timestamp": this.newTimestamp,
        "google.protobuf.FieldMask": this.newFieldMask
      };
    }
    /**
     * Convert native Date to protobuf Timestamp format for serialization
     */
    serializeTimestamp(date) {
      if (!date) return null;
      return {
        seconds: Math.floor(date.getTime() / 1e3).toString(),
        nanos: date.getTime() % 1e3 * 1e6
      };
    }
    /**
     * Convert native string array to protobuf FieldMask format for serialization
     */
    serializeFieldMask(paths) {
      if (!paths || !Array.isArray(paths)) return null;
      return { paths };
    }
  };

  // gen/wasmts/connect4/factory.ts
  var Connect4Factory = class {
    constructor() {
      // Dependency factory for wasmjs.v1 package
      this.v1Factory = new WasmjsV1Factory();
      /**
       * Enhanced factory method for GameState
       * @param parent Parent object containing this field
       * @param attributeName Field name in parent object
       * @param attributeKey Array index, map key, or union tag (for containers)
       * @param data Raw data to potentially populate from
       * @returns Factory result with instance and population status
       */
      this.newGameState = (parent, attributeName, attributeKey, data) => {
        const out = new GameState();
        return { instance: out, fullyLoaded: false };
      };
      /**
       * Enhanced factory method for GameConfig
       * @param parent Parent object containing this field
       * @param attributeName Field name in parent object
       * @param attributeKey Array index, map key, or union tag (for containers)
       * @param data Raw data to potentially populate from
       * @returns Factory result with instance and population status
       */
      this.newGameConfig = (parent, attributeName, attributeKey, data) => {
        const out = new GameConfig();
        return { instance: out, fullyLoaded: false };
      };
      /**
       * Enhanced factory method for Player
       * @param parent Parent object containing this field
       * @param attributeName Field name in parent object
       * @param attributeKey Array index, map key, or union tag (for containers)
       * @param data Raw data to potentially populate from
       * @returns Factory result with instance and population status
       */
      this.newPlayer = (parent, attributeName, attributeKey, data) => {
        const out = new Player();
        return { instance: out, fullyLoaded: false };
      };
      /**
       * Enhanced factory method for GameBoard
       * @param parent Parent object containing this field
       * @param attributeName Field name in parent object
       * @param attributeKey Array index, map key, or union tag (for containers)
       * @param data Raw data to potentially populate from
       * @returns Factory result with instance and population status
       */
      this.newGameBoard = (parent, attributeName, attributeKey, data) => {
        const out = new GameBoard();
        return { instance: out, fullyLoaded: false };
      };
      /**
       * Enhanced factory method for BoardRow
       * @param parent Parent object containing this field
       * @param attributeName Field name in parent object
       * @param attributeKey Array index, map key, or union tag (for containers)
       * @param data Raw data to potentially populate from
       * @returns Factory result with instance and population status
       */
      this.newBoardRow = (parent, attributeName, attributeKey, data) => {
        const out = new BoardRow();
        return { instance: out, fullyLoaded: false };
      };
      /**
       * Enhanced factory method for PlayerStats
       * @param parent Parent object containing this field
       * @param attributeName Field name in parent object
       * @param attributeKey Array index, map key, or union tag (for containers)
       * @param data Raw data to potentially populate from
       * @returns Factory result with instance and population status
       */
      this.newPlayerStats = (parent, attributeName, attributeKey, data) => {
        const out = new PlayerStats();
        return { instance: out, fullyLoaded: false };
      };
      /**
       * Enhanced factory method for GetGameRequest
       * @param parent Parent object containing this field
       * @param attributeName Field name in parent object
       * @param attributeKey Array index, map key, or union tag (for containers)
       * @param data Raw data to potentially populate from
       * @returns Factory result with instance and population status
       */
      this.newGetGameRequest = (parent, attributeName, attributeKey, data) => {
        const out = new GetGameRequest();
        return { instance: out, fullyLoaded: false };
      };
      /**
       * Enhanced factory method for DropPieceRequest
       * @param parent Parent object containing this field
       * @param attributeName Field name in parent object
       * @param attributeKey Array index, map key, or union tag (for containers)
       * @param data Raw data to potentially populate from
       * @returns Factory result with instance and population status
       */
      this.newDropPieceRequest = (parent, attributeName, attributeKey, data) => {
        const out = new DropPieceRequest();
        return { instance: out, fullyLoaded: false };
      };
      /**
       * Enhanced factory method for DropPieceResponse
       * @param parent Parent object containing this field
       * @param attributeName Field name in parent object
       * @param attributeKey Array index, map key, or union tag (for containers)
       * @param data Raw data to potentially populate from
       * @returns Factory result with instance and population status
       */
      this.newDropPieceResponse = (parent, attributeName, attributeKey, data) => {
        const out = new DropPieceResponse();
        return { instance: out, fullyLoaded: false };
      };
      /**
       * Enhanced factory method for PieceDropResult
       * @param parent Parent object containing this field
       * @param attributeName Field name in parent object
       * @param attributeKey Array index, map key, or union tag (for containers)
       * @param data Raw data to potentially populate from
       * @returns Factory result with instance and population status
       */
      this.newPieceDropResult = (parent, attributeName, attributeKey, data) => {
        const out = new PieceDropResult();
        return { instance: out, fullyLoaded: false };
      };
      /**
       * Enhanced factory method for LineInfo
       * @param parent Parent object containing this field
       * @param attributeName Field name in parent object
       * @param attributeKey Array index, map key, or union tag (for containers)
       * @param data Raw data to potentially populate from
       * @returns Factory result with instance and population status
       */
      this.newLineInfo = (parent, attributeName, attributeKey, data) => {
        const out = new LineInfo();
        return { instance: out, fullyLoaded: false };
      };
      /**
       * Enhanced factory method for Position
       * @param parent Parent object containing this field
       * @param attributeName Field name in parent object
       * @param attributeKey Array index, map key, or union tag (for containers)
       * @param data Raw data to potentially populate from
       * @returns Factory result with instance and population status
       */
      this.newPosition = (parent, attributeName, attributeKey, data) => {
        const out = new Position();
        return { instance: out, fullyLoaded: false };
      };
      /**
       * Enhanced factory method for JoinGameRequest
       * @param parent Parent object containing this field
       * @param attributeName Field name in parent object
       * @param attributeKey Array index, map key, or union tag (for containers)
       * @param data Raw data to potentially populate from
       * @returns Factory result with instance and population status
       */
      this.newJoinGameRequest = (parent, attributeName, attributeKey, data) => {
        const out = new JoinGameRequest();
        return { instance: out, fullyLoaded: false };
      };
      /**
       * Enhanced factory method for JoinGameResponse
       * @param parent Parent object containing this field
       * @param attributeName Field name in parent object
       * @param attributeKey Array index, map key, or union tag (for containers)
       * @param data Raw data to potentially populate from
       * @returns Factory result with instance and population status
       */
      this.newJoinGameResponse = (parent, attributeName, attributeKey, data) => {
        const out = new JoinGameResponse();
        return { instance: out, fullyLoaded: false };
      };
      /**
       * Enhanced factory method for CreateGameRequest
       * @param parent Parent object containing this field
       * @param attributeName Field name in parent object
       * @param attributeKey Array index, map key, or union tag (for containers)
       * @param data Raw data to potentially populate from
       * @returns Factory result with instance and population status
       */
      this.newCreateGameRequest = (parent, attributeName, attributeKey, data) => {
        const out = new CreateGameRequest();
        return { instance: out, fullyLoaded: false };
      };
      /**
       * Enhanced factory method for CreateGameResponse
       * @param parent Parent object containing this field
       * @param attributeName Field name in parent object
       * @param attributeKey Array index, map key, or union tag (for containers)
       * @param data Raw data to potentially populate from
       * @returns Factory result with instance and population status
       */
      this.newCreateGameResponse = (parent, attributeName, attributeKey, data) => {
        const out = new CreateGameResponse();
        return { instance: out, fullyLoaded: false };
      };
      /**
       * Factory method for converting protobuf Timestamp data to native Date
       * @param parent Parent object containing this field
       * @param attributeName Field name in parent object  
       * @param attributeKey Array index, map key, or union tag (for containers)
       * @param data Raw protobuf timestamp data
       * @returns Factory result with Date instance
       */
      this.newTimestamp = (parent, attributeName, attributeKey, data) => {
        if (!data) {
          return { instance: /* @__PURE__ */ new Date(), fullyLoaded: true };
        }
        let date;
        if (typeof data === "string") {
          date = new Date(data);
        } else if (data.seconds !== void 0) {
          const seconds = typeof data.seconds === "string" ? parseInt(data.seconds, 10) : data.seconds;
          const nanos = data.nanos || 0;
          date = new Date(seconds * 1e3 + Math.floor(nanos / 1e6));
        } else {
          date = /* @__PURE__ */ new Date();
        }
        return { instance: date, fullyLoaded: true };
      };
      /**
       * Factory method for converting protobuf FieldMask data to native string array
       * @param parent Parent object containing this field
       * @param attributeName Field name in parent object
       * @param attributeKey Array index, map key, or union tag (for containers)
       * @param data Raw protobuf field mask data
       * @returns Factory result with string array instance
       */
      this.newFieldMask = (parent, attributeName, attributeKey, data) => {
        if (!data) {
          return { instance: [], fullyLoaded: true };
        }
        let paths;
        if (Array.isArray(data)) {
          paths = data;
        } else if (data.paths && Array.isArray(data.paths)) {
          paths = data.paths;
        } else {
          paths = [];
        }
        return { instance: paths, fullyLoaded: true };
      };
    }
    /**
     * Get factory method for a fully qualified message type
     * Enables cross-package factory delegation
     */
    getFactoryMethod(messageType) {
      const parts = messageType.split(".");
      if (parts.length < 2) {
        return void 0;
      }
      const packageName = parts.slice(0, -1).join(".");
      const typeName = parts[parts.length - 1];
      const methodName = "new" + typeName;
      const currentPackage = "connect4";
      if (packageName === currentPackage) {
        return this[methodName];
      }
      const externalFactory = this.externalTypeFactories()[messageType];
      if (externalFactory) {
        return externalFactory;
      }
      if (packageName === "wasmjs.v1") {
        return this.v1Factory[methodName];
      }
      return void 0;
    }
    /**
     * Generic object deserializer that respects factory decisions
     */
    deserializeObject(instance, data) {
      if (!data || typeof data !== "object") return instance;
      for (const [key, value] of Object.entries(data)) {
        if (value !== null && value !== void 0) {
          instance[key] = value;
        }
      }
      return instance;
    }
    // External type conversion methods
    /**
     * Mapping of external types to their factory methods
     */
    externalTypeFactories() {
      return {
        "google.protobuf.Timestamp": this.newTimestamp,
        "google.protobuf.FieldMask": this.newFieldMask
      };
    }
    /**
     * Convert native Date to protobuf Timestamp format for serialization
     */
    serializeTimestamp(date) {
      if (!date) return null;
      return {
        seconds: Math.floor(date.getTime() / 1e3).toString(),
        nanos: date.getTime() % 1e3 * 1e6
      };
    }
    /**
     * Convert native string array to protobuf FieldMask format for serialization
     */
    serializeFieldMask(paths) {
      if (!paths || !Array.isArray(paths)) return null;
      return { paths };
    }
  };

  // gen/wasmts/connect4/schemas.ts
  var GameStateSchema = {
    name: "GameState",
    fields: [
      {
        name: "gameId",
        type: "string" /* STRING */,
        id: 1
      },
      {
        name: "config",
        type: "message" /* MESSAGE */,
        id: 2,
        messageType: "connect4.GameConfig"
      },
      {
        name: "players",
        type: "message" /* MESSAGE */,
        id: 3,
        messageType: "connect4.Player",
        repeated: true
      },
      {
        name: "board",
        type: "message" /* MESSAGE */,
        id: 4,
        messageType: "connect4.GameBoard"
      },
      {
        name: "currentPlayerId",
        type: "string" /* STRING */,
        id: 5
      },
      {
        name: "turnNumber",
        type: "number" /* NUMBER */,
        id: 6
      },
      {
        name: "status",
        type: "string" /* STRING */,
        id: 7
      },
      {
        name: "winners",
        type: "repeated" /* REPEATED */,
        id: 8,
        repeated: true
      },
      {
        name: "playerStats",
        type: "message" /* MESSAGE */,
        id: 9,
        messageType: "connect4.PlayerStatsEntry"
      },
      {
        name: "lastMoveTime",
        type: "number" /* NUMBER */,
        id: 10
      },
      {
        name: "moveTimeoutSeconds",
        type: "number" /* NUMBER */,
        id: 11
      }
    ]
  };
  var GameConfigSchema = {
    name: "GameConfig",
    fields: [
      {
        name: "boardWidth",
        type: "number" /* NUMBER */,
        id: 1
      },
      {
        name: "boardHeight",
        type: "number" /* NUMBER */,
        id: 2
      },
      {
        name: "minPlayers",
        type: "number" /* NUMBER */,
        id: 3
      },
      {
        name: "maxPlayers",
        type: "number" /* NUMBER */,
        id: 4
      },
      {
        name: "connectLength",
        type: "number" /* NUMBER */,
        id: 5
      },
      {
        name: "allowMultipleWinners",
        type: "boolean" /* BOOLEAN */,
        id: 6
      },
      {
        name: "moveTimeoutSeconds",
        type: "number" /* NUMBER */,
        id: 7
      }
    ]
  };
  var PlayerSchema = {
    name: "Player",
    fields: [
      {
        name: "id",
        type: "string" /* STRING */,
        id: 1
      },
      {
        name: "name",
        type: "string" /* STRING */,
        id: 2
      },
      {
        name: "color",
        type: "string" /* STRING */,
        id: 3
      },
      {
        name: "isConnected",
        type: "boolean" /* BOOLEAN */,
        id: 4
      },
      {
        name: "joinOrder",
        type: "number" /* NUMBER */,
        id: 5
      }
    ]
  };
  var GameBoardSchema = {
    name: "GameBoard",
    fields: [
      {
        name: "width",
        type: "number" /* NUMBER */,
        id: 1
      },
      {
        name: "height",
        type: "number" /* NUMBER */,
        id: 2
      },
      {
        name: "rows",
        type: "message" /* MESSAGE */,
        id: 3,
        messageType: "connect4.BoardRow",
        repeated: true
      },
      {
        name: "columnHeights",
        type: "repeated" /* REPEATED */,
        id: 4,
        repeated: true
      }
    ]
  };
  var BoardRowSchema = {
    name: "BoardRow",
    fields: [
      {
        name: "cells",
        type: "repeated" /* REPEATED */,
        id: 1,
        repeated: true
      }
    ]
  };
  var PlayerStatsSchema = {
    name: "PlayerStats",
    fields: [
      {
        name: "piecesPlayed",
        type: "number" /* NUMBER */,
        id: 1
      },
      {
        name: "winningLines",
        type: "number" /* NUMBER */,
        id: 2
      },
      {
        name: "hasWon",
        type: "boolean" /* BOOLEAN */,
        id: 3
      },
      {
        name: "totalMoveTime",
        type: "number" /* NUMBER */,
        id: 4
      }
    ]
  };
  var GetGameRequestSchema = {
    name: "GetGameRequest",
    fields: [
      {
        name: "gameId",
        type: "string" /* STRING */,
        id: 1
      }
    ]
  };
  var DropPieceRequestSchema = {
    name: "DropPieceRequest",
    fields: [
      {
        name: "gameId",
        type: "string" /* STRING */,
        id: 1
      },
      {
        name: "playerId",
        type: "string" /* STRING */,
        id: 2
      },
      {
        name: "column",
        type: "number" /* NUMBER */,
        id: 3
      }
    ]
  };
  var DropPieceResponseSchema = {
    name: "DropPieceResponse",
    fields: [
      {
        name: "success",
        type: "boolean" /* BOOLEAN */,
        id: 1
      },
      {
        name: "errorMessage",
        type: "string" /* STRING */,
        id: 2
      },
      {
        name: "patches",
        type: "message" /* MESSAGE */,
        id: 3,
        messageType: "wasmjs.v1.MessagePatch",
        repeated: true
      },
      {
        name: "changeNumber",
        type: "number" /* NUMBER */,
        id: 4
      },
      {
        name: "result",
        type: "message" /* MESSAGE */,
        id: 5,
        messageType: "connect4.PieceDropResult"
      }
    ]
  };
  var PieceDropResultSchema = {
    name: "PieceDropResult",
    fields: [
      {
        name: "finalRow",
        type: "number" /* NUMBER */,
        id: 1
      },
      {
        name: "finalColumn",
        type: "number" /* NUMBER */,
        id: 2
      },
      {
        name: "formedLine",
        type: "boolean" /* BOOLEAN */,
        id: 3
      },
      {
        name: "winningLines",
        type: "message" /* MESSAGE */,
        id: 4,
        messageType: "connect4.LineInfo",
        repeated: true
      }
    ]
  };
  var LineInfoSchema = {
    name: "LineInfo",
    fields: [
      {
        name: "positions",
        type: "message" /* MESSAGE */,
        id: 1,
        messageType: "connect4.Position",
        repeated: true
      },
      {
        name: "direction",
        type: "string" /* STRING */,
        id: 2
      },
      {
        name: "length",
        type: "number" /* NUMBER */,
        id: 3
      }
    ]
  };
  var PositionSchema = {
    name: "Position",
    fields: [
      {
        name: "row",
        type: "number" /* NUMBER */,
        id: 1
      },
      {
        name: "column",
        type: "number" /* NUMBER */,
        id: 2
      }
    ]
  };
  var JoinGameRequestSchema = {
    name: "JoinGameRequest",
    fields: [
      {
        name: "gameId",
        type: "string" /* STRING */,
        id: 1
      },
      {
        name: "playerName",
        type: "string" /* STRING */,
        id: 2
      },
      {
        name: "preferredColor",
        type: "string" /* STRING */,
        id: 3
      }
    ]
  };
  var JoinGameResponseSchema = {
    name: "JoinGameResponse",
    fields: [
      {
        name: "success",
        type: "boolean" /* BOOLEAN */,
        id: 1
      },
      {
        name: "errorMessage",
        type: "string" /* STRING */,
        id: 2
      },
      {
        name: "playerId",
        type: "string" /* STRING */,
        id: 3
      },
      {
        name: "assignedColor",
        type: "string" /* STRING */,
        id: 4
      },
      {
        name: "gameState",
        type: "message" /* MESSAGE */,
        id: 5,
        messageType: "connect4.GameState"
      }
    ]
  };
  var CreateGameRequestSchema = {
    name: "CreateGameRequest",
    fields: [
      {
        name: "gameId",
        type: "string" /* STRING */,
        id: 1
      },
      {
        name: "config",
        type: "message" /* MESSAGE */,
        id: 2,
        messageType: "connect4.GameConfig"
      },
      {
        name: "creatorName",
        type: "string" /* STRING */,
        id: 3
      }
    ]
  };
  var CreateGameResponseSchema = {
    name: "CreateGameResponse",
    fields: [
      {
        name: "success",
        type: "boolean" /* BOOLEAN */,
        id: 1
      },
      {
        name: "errorMessage",
        type: "string" /* STRING */,
        id: 2
      },
      {
        name: "playerId",
        type: "string" /* STRING */,
        id: 3
      },
      {
        name: "gameState",
        type: "message" /* MESSAGE */,
        id: 4,
        messageType: "connect4.GameState"
      }
    ]
  };
  var Connect4SchemaRegistry = {
    "connect4.GameState": GameStateSchema,
    "connect4.GameConfig": GameConfigSchema,
    "connect4.Player": PlayerSchema,
    "connect4.GameBoard": GameBoardSchema,
    "connect4.BoardRow": BoardRowSchema,
    "connect4.PlayerStats": PlayerStatsSchema,
    "connect4.GetGameRequest": GetGameRequestSchema,
    "connect4.DropPieceRequest": DropPieceRequestSchema,
    "connect4.DropPieceResponse": DropPieceResponseSchema,
    "connect4.PieceDropResult": PieceDropResultSchema,
    "connect4.LineInfo": LineInfoSchema,
    "connect4.Position": PositionSchema,
    "connect4.JoinGameRequest": JoinGameRequestSchema,
    "connect4.JoinGameResponse": JoinGameResponseSchema,
    "connect4.CreateGameRequest": CreateGameRequestSchema,
    "connect4.CreateGameResponse": CreateGameResponseSchema
  };

  // gen/wasmts/connect4/deserializer.ts
  var DEFAULT_FACTORY2 = new Connect4Factory();
  var Connect4Deserializer = class _Connect4Deserializer {
    constructor(schemaRegistry = Connect4SchemaRegistry, factory = DEFAULT_FACTORY2) {
      this.schemaRegistry = schemaRegistry;
      this.factory = factory;
    }
    /**
     * Deserialize an object using schema information
     * @param instance The target instance to populate
     * @param data The source data to deserialize from
     * @param messageType The fully qualified message type (e.g., "library.v1.Book")
     * @returns The populated instance
     */
    deserialize(instance, data, messageType) {
      if (!data || typeof data !== "object") {
        return instance;
      }
      const schema = this.schemaRegistry[messageType];
      if (!schema) {
        return this.fallbackDeserialize(instance, data);
      }
      for (const fieldSchema of schema.fields) {
        const fieldValue = data[fieldSchema.name];
        if (fieldValue === null || fieldValue === void 0) {
          continue;
        }
        this.deserializeField(instance, fieldSchema, fieldValue);
      }
      return instance;
    }
    /**
     * Deserialize a single field based on its schema
     */
    deserializeField(instance, fieldSchema, fieldValue) {
      const fieldName = fieldSchema.name;
      switch (fieldSchema.type) {
        case "string" /* STRING */:
        case "number" /* NUMBER */:
        case "boolean" /* BOOLEAN */:
          instance[fieldName] = fieldValue;
          break;
        case "message" /* MESSAGE */:
          if (fieldSchema.repeated) {
            instance[fieldName] = this.deserializeMessageArray(
              fieldValue,
              fieldSchema.messageType,
              instance,
              fieldName
            );
          } else {
            instance[fieldName] = this.deserializeMessageField(
              fieldValue,
              fieldSchema.messageType,
              instance,
              fieldName
            );
          }
          break;
        case "repeated" /* REPEATED */:
          if (Array.isArray(fieldValue)) {
            instance[fieldName] = [...fieldValue];
          }
          break;
        case "oneof" /* ONEOF */:
          instance[fieldName] = fieldValue;
          break;
        case "map" /* MAP */:
          instance[fieldName] = __spreadValues({}, fieldValue);
          break;
        default:
          instance[fieldName] = fieldValue;
          break;
      }
    }
    /**
     * Deserialize a single message field
     */
    deserializeMessageField(fieldValue, messageType, parent, attributeName) {
      let factoryMethod;
      if (this.factory.getFactoryMethod) {
        factoryMethod = this.factory.getFactoryMethod(messageType);
      } else {
        const factoryMethodName = this.getFactoryMethodName(messageType);
        factoryMethod = this.factory[factoryMethodName];
      }
      if (factoryMethod) {
        const result = factoryMethod(parent, attributeName, void 0, fieldValue);
        if (result.fullyLoaded) {
          return result.instance;
        } else {
          return this.deserialize(result.instance, fieldValue, messageType);
        }
      }
      return this.fallbackDeserialize({}, fieldValue);
    }
    /**
     * Deserialize an array of message objects
     */
    deserializeMessageArray(fieldValue, messageType, parent, attributeName) {
      if (!Array.isArray(fieldValue)) {
        return [];
      }
      let factoryMethod;
      if (this.factory.getFactoryMethod) {
        factoryMethod = this.factory.getFactoryMethod(messageType);
      } else {
        const factoryMethodName = this.getFactoryMethodName(messageType);
        factoryMethod = this.factory[factoryMethodName];
      }
      return fieldValue.map((item, index) => {
        if (factoryMethod) {
          const result = factoryMethod(parent, attributeName, index, item);
          if (result.fullyLoaded) {
            return result.instance;
          } else {
            return this.deserialize(result.instance, item, messageType);
          }
        }
        return this.fallbackDeserialize({}, item);
      });
    }
    /**
     * Convert message type to factory method name
     * "library.v1.Book" -> "newBook"
     */
    getFactoryMethodName(messageType) {
      const parts = messageType.split(".");
      const typeName = parts[parts.length - 1];
      return "new" + typeName;
    }
    /**
     * Fallback deserializer for when no schema is available
     */
    fallbackDeserialize(instance, data) {
      if (!data || typeof data !== "object") {
        return instance;
      }
      for (const [key, value] of Object.entries(data)) {
        if (value !== null && value !== void 0) {
          instance[key] = value;
        }
      }
      return instance;
    }
    /**
     * Create and deserialize a new instance of a message type
     */
    createAndDeserialize(messageType, data) {
      let factoryMethod;
      if (this.factory.getFactoryMethod) {
        factoryMethod = this.factory.getFactoryMethod(messageType);
      } else {
        const factoryMethodName = this.getFactoryMethodName(messageType);
        factoryMethod = this.factory[factoryMethodName];
      }
      if (!factoryMethod) {
        throw new Error(`Could not find factory method to deserialize: ${messageType}`);
      }
      const result = factoryMethod(void 0, void 0, void 0, data);
      if (result.fullyLoaded) {
        return result.instance;
      } else {
        return this.deserialize(result.instance, data, messageType);
      }
    }
    /**
     * Static utility method to create and deserialize a message without needing a deserializer instance
     * @param messageType Fully qualified message type (use Class.MESSAGE_TYPE)
     * @param data Raw data to deserialize
     * @returns Deserialized instance or null if creation failed
     */
    static from(messageType, data) {
      const deserializer = new _Connect4Deserializer();
      return deserializer.createAndDeserialize(messageType, data);
    }
  };

  // gen/wasmts/connect4/models.ts
  var _GameState = class _GameState {
    constructor() {
      this.gameId = "";
      this.players = [];
      this.currentPlayerId = "";
      this.turnNumber = 0;
      this.status = 0;
      this.winners = [];
      this.lastMoveTime = 0;
      this.moveTimeoutSeconds = 0;
    }
    /**
     * Create and deserialize an instance from raw data
     * @param data Raw data to deserialize
     * @returns Deserialized GameState instance or null if creation failed
     */
    static from(data) {
      return Connect4Deserializer.from(_GameState.MESSAGE_TYPE, data);
    }
  };
  /**
   * Fully qualified message type for schema resolution
   */
  _GameState.MESSAGE_TYPE = "connect4.GameState";
  var GameState = _GameState;
  var _GameConfig = class _GameConfig {
    constructor() {
      this.boardWidth = 0;
      this.boardHeight = 0;
      this.minPlayers = 0;
      this.maxPlayers = 0;
      this.connectLength = 0;
      this.allowMultipleWinners = false;
      this.moveTimeoutSeconds = 0;
    }
    /**
     * Create and deserialize an instance from raw data
     * @param data Raw data to deserialize
     * @returns Deserialized GameConfig instance or null if creation failed
     */
    static from(data) {
      return Connect4Deserializer.from(_GameConfig.MESSAGE_TYPE, data);
    }
  };
  /**
   * Fully qualified message type for schema resolution
   */
  _GameConfig.MESSAGE_TYPE = "connect4.GameConfig";
  var GameConfig = _GameConfig;
  var _Player = class _Player {
    constructor() {
      this.id = "";
      this.name = "";
      this.color = "";
      this.isConnected = false;
      this.joinOrder = 0;
    }
    /**
     * Create and deserialize an instance from raw data
     * @param data Raw data to deserialize
     * @returns Deserialized Player instance or null if creation failed
     */
    static from(data) {
      return Connect4Deserializer.from(_Player.MESSAGE_TYPE, data);
    }
  };
  /**
   * Fully qualified message type for schema resolution
   */
  _Player.MESSAGE_TYPE = "connect4.Player";
  var Player = _Player;
  var _GameBoard = class _GameBoard {
    constructor() {
      this.width = 0;
      this.height = 0;
      /** Board representation: grid[y][x] = player_id (empty = "") */
      this.rows = [];
      this.columnHeights = [];
    }
    /**
     * Create and deserialize an instance from raw data
     * @param data Raw data to deserialize
     * @returns Deserialized GameBoard instance or null if creation failed
     */
    static from(data) {
      return Connect4Deserializer.from(_GameBoard.MESSAGE_TYPE, data);
    }
  };
  /**
   * Fully qualified message type for schema resolution
   */
  _GameBoard.MESSAGE_TYPE = "connect4.GameBoard";
  var GameBoard = _GameBoard;
  var _BoardRow = class _BoardRow {
    constructor() {
      this.cells = [];
    }
    /**
     * Create and deserialize an instance from raw data
     * @param data Raw data to deserialize
     * @returns Deserialized BoardRow instance or null if creation failed
     */
    static from(data) {
      return Connect4Deserializer.from(_BoardRow.MESSAGE_TYPE, data);
    }
  };
  /**
   * Fully qualified message type for schema resolution
   */
  _BoardRow.MESSAGE_TYPE = "connect4.BoardRow";
  var BoardRow = _BoardRow;
  var _PlayerStats = class _PlayerStats {
    constructor() {
      this.piecesPlayed = 0;
      this.winningLines = 0;
      this.hasWon = false;
      this.totalMoveTime = 0;
    }
    /**
     * Create and deserialize an instance from raw data
     * @param data Raw data to deserialize
     * @returns Deserialized PlayerStats instance or null if creation failed
     */
    static from(data) {
      return Connect4Deserializer.from(_PlayerStats.MESSAGE_TYPE, data);
    }
  };
  /**
   * Fully qualified message type for schema resolution
   */
  _PlayerStats.MESSAGE_TYPE = "connect4.PlayerStats";
  var PlayerStats = _PlayerStats;
  var _GetGameRequest = class _GetGameRequest {
    constructor() {
      this.gameId = "";
    }
    /**
     * Create and deserialize an instance from raw data
     * @param data Raw data to deserialize
     * @returns Deserialized GetGameRequest instance or null if creation failed
     */
    static from(data) {
      return Connect4Deserializer.from(_GetGameRequest.MESSAGE_TYPE, data);
    }
  };
  /**
   * Fully qualified message type for schema resolution
   */
  _GetGameRequest.MESSAGE_TYPE = "connect4.GetGameRequest";
  var GetGameRequest = _GetGameRequest;
  var _DropPieceRequest = class _DropPieceRequest {
    constructor() {
      this.gameId = "";
      this.playerId = "";
      this.column = 0;
    }
    /**
     * Create and deserialize an instance from raw data
     * @param data Raw data to deserialize
     * @returns Deserialized DropPieceRequest instance or null if creation failed
     */
    static from(data) {
      return Connect4Deserializer.from(_DropPieceRequest.MESSAGE_TYPE, data);
    }
  };
  /**
   * Fully qualified message type for schema resolution
   */
  _DropPieceRequest.MESSAGE_TYPE = "connect4.DropPieceRequest";
  var DropPieceRequest = _DropPieceRequest;
  var _DropPieceResponse = class _DropPieceResponse {
    constructor() {
      this.success = false;
      this.errorMessage = "";
      this.patches = [];
      this.changeNumber = 0;
    }
    /**
     * Create and deserialize an instance from raw data
     * @param data Raw data to deserialize
     * @returns Deserialized DropPieceResponse instance or null if creation failed
     */
    static from(data) {
      return Connect4Deserializer.from(_DropPieceResponse.MESSAGE_TYPE, data);
    }
  };
  /**
   * Fully qualified message type for schema resolution
   */
  _DropPieceResponse.MESSAGE_TYPE = "connect4.DropPieceResponse";
  var DropPieceResponse = _DropPieceResponse;
  var _PieceDropResult = class _PieceDropResult {
    constructor() {
      this.finalRow = 0;
      this.finalColumn = 0;
      this.formedLine = false;
      this.winningLines = [];
    }
    /**
     * Create and deserialize an instance from raw data
     * @param data Raw data to deserialize
     * @returns Deserialized PieceDropResult instance or null if creation failed
     */
    static from(data) {
      return Connect4Deserializer.from(_PieceDropResult.MESSAGE_TYPE, data);
    }
  };
  /**
   * Fully qualified message type for schema resolution
   */
  _PieceDropResult.MESSAGE_TYPE = "connect4.PieceDropResult";
  var PieceDropResult = _PieceDropResult;
  var _LineInfo = class _LineInfo {
    constructor() {
      this.positions = [];
      this.direction = "";
      this.length = 0;
    }
    /**
     * Create and deserialize an instance from raw data
     * @param data Raw data to deserialize
     * @returns Deserialized LineInfo instance or null if creation failed
     */
    static from(data) {
      return Connect4Deserializer.from(_LineInfo.MESSAGE_TYPE, data);
    }
  };
  /**
   * Fully qualified message type for schema resolution
   */
  _LineInfo.MESSAGE_TYPE = "connect4.LineInfo";
  var LineInfo = _LineInfo;
  var _Position = class _Position {
    constructor() {
      this.row = 0;
      this.column = 0;
    }
    /**
     * Create and deserialize an instance from raw data
     * @param data Raw data to deserialize
     * @returns Deserialized Position instance or null if creation failed
     */
    static from(data) {
      return Connect4Deserializer.from(_Position.MESSAGE_TYPE, data);
    }
  };
  /**
   * Fully qualified message type for schema resolution
   */
  _Position.MESSAGE_TYPE = "connect4.Position";
  var Position = _Position;
  var _JoinGameRequest = class _JoinGameRequest {
    constructor() {
      this.gameId = "";
      this.playerName = "";
      this.preferredColor = "";
    }
    /**
     * Create and deserialize an instance from raw data
     * @param data Raw data to deserialize
     * @returns Deserialized JoinGameRequest instance or null if creation failed
     */
    static from(data) {
      return Connect4Deserializer.from(_JoinGameRequest.MESSAGE_TYPE, data);
    }
  };
  /**
   * Fully qualified message type for schema resolution
   */
  _JoinGameRequest.MESSAGE_TYPE = "connect4.JoinGameRequest";
  var JoinGameRequest = _JoinGameRequest;
  var _JoinGameResponse = class _JoinGameResponse {
    constructor() {
      this.success = false;
      this.errorMessage = "";
      this.playerId = "";
      this.assignedColor = "";
    }
    /**
     * Create and deserialize an instance from raw data
     * @param data Raw data to deserialize
     * @returns Deserialized JoinGameResponse instance or null if creation failed
     */
    static from(data) {
      return Connect4Deserializer.from(_JoinGameResponse.MESSAGE_TYPE, data);
    }
  };
  /**
   * Fully qualified message type for schema resolution
   */
  _JoinGameResponse.MESSAGE_TYPE = "connect4.JoinGameResponse";
  var JoinGameResponse = _JoinGameResponse;
  var _CreateGameRequest = class _CreateGameRequest {
    constructor() {
      this.gameId = "";
      this.creatorName = "";
    }
    /**
     * Create and deserialize an instance from raw data
     * @param data Raw data to deserialize
     * @returns Deserialized CreateGameRequest instance or null if creation failed
     */
    static from(data) {
      return Connect4Deserializer.from(_CreateGameRequest.MESSAGE_TYPE, data);
    }
  };
  /**
   * Fully qualified message type for schema resolution
   */
  _CreateGameRequest.MESSAGE_TYPE = "connect4.CreateGameRequest";
  var CreateGameRequest = _CreateGameRequest;
  var _CreateGameResponse = class _CreateGameResponse {
    constructor() {
      this.success = false;
      this.errorMessage = "";
      this.playerId = "";
    }
    /**
     * Create and deserialize an instance from raw data
     * @param data Raw data to deserialize
     * @returns Deserialized CreateGameResponse instance or null if creation failed
     */
    static from(data) {
      return Connect4Deserializer.from(_CreateGameResponse.MESSAGE_TYPE, data);
    }
  };
  /**
   * Fully qualified message type for schema resolution
   */
  _CreateGameResponse.MESSAGE_TYPE = "connect4.CreateGameResponse";
  var CreateGameResponse = _CreateGameResponse;

  // src/gameViewer.ts
  var GameViewer = class {
    constructor() {
      this.elements = {};
      const pathParts = window.location.pathname.split("/").filter((p) => p);
      const gameId = pathParts[0] || "";
      this.ui = {
        gameId,
        playerId: "",
        gameState: null,
        statefulProxy: null,
        connect4Client: null
      };
      this.init();
    }
    init() {
      return __async(this, null, function* () {
        if (document.readyState === "loading") {
          document.addEventListener("DOMContentLoaded", () => this.initializeUI());
        } else {
          this.initializeUI();
        }
      });
    }
    initializeUI() {
      this.elements = {
        joinGameForm: document.getElementById("joinGameForm"),
        gameInterface: document.getElementById("gameInterface"),
        errorState: document.getElementById("errorState"),
        gameBoard: document.getElementById("gameBoard"),
        gameStatus: document.getElementById("gameStatus"),
        currentPlayerName: document.getElementById("currentPlayerName"),
        currentPlayerColor: document.getElementById("currentPlayerColor"),
        turnNumber: document.getElementById("turnNumber"),
        playersList: document.getElementById("playersList"),
        gameLog: document.getElementById("gameLog"),
        currentGameId: document.getElementById("currentGameId"),
        gameUrl: document.getElementById("gameUrl")
      };
      if (this.elements.currentGameId) {
        this.elements.currentGameId.textContent = this.ui.gameId;
      }
      if (this.elements.gameUrl) {
        this.elements.gameUrl.textContent = window.location.href;
      }
      if (!this.ui.gameId || !this.isValidGameId(this.ui.gameId)) {
        this.showError("Invalid game ID");
        return;
      }
      const joinForm = document.getElementById("joinGameForm");
      if (joinForm) {
        joinForm.addEventListener("submit", (e) => this.handleJoinGame(e));
      }
      this.initializeWasmClient();
      this.initializeStatefulProxy();
      this.loadStoredGameState();
    }
    initializeWasmClient() {
      return __async(this, null, function* () {
        try {
          console.log("Loading WASM module...");
          this.ui.connect4Client = new multiplayer_connect4Client_client_default();
          yield this.ui.connect4Client.loadWasm("/static/wasm/multiplayer_connect4.wasm");
          yield this.ui.connect4Client.waitUntilReady();
          console.log("WASM module loaded successfully!");
        } catch (error) {
          console.error("Failed to load WASM:", error);
          this.showError("Failed to load game engine");
        }
      });
    }
    initializeStatefulProxy() {
      try {
        this.ui.statefulProxy = new StatefulProxy(this.ui.gameId, "indexeddb");
        this.ui.statefulProxy.onStateChange((patches) => {
          console.log("Received state patches:", patches);
          this.applyPatches(patches);
        });
        console.log("Stateful proxy initialized");
      } catch (error) {
        console.error("Failed to initialize stateful proxy:", error);
      }
    }
    loadStoredGameState() {
      try {
        const stored = localStorage.getItem(`connect4_game_${this.ui.gameId}`);
        if (stored) {
          const parsedState = JSON.parse(stored);
          this.ui.gameState = GameState.from(parsedState);
          console.log("Loaded stored game state:", this.ui.gameState);
          if (this.ui.gameState && this.ui.gameState.players.length > 0) {
            const storedPlayerId = localStorage.getItem(`connect4_player_${this.ui.gameId}`);
            if (storedPlayerId) {
              this.ui.playerId = storedPlayerId;
              this.showGameInterface();
              this.updateGameDisplay();
              return;
            }
          }
        }
      } catch (error) {
        console.error("Failed to parse stored game state:", error);
      }
      this.showJoinForm();
    }
    handleJoinGame(event) {
      return __async(this, null, function* () {
        event.preventDefault();
        const formData = new FormData(event.target);
        const playerName = formData.get("playerName");
        if (!playerName.trim()) {
          alert("Please enter your name");
          return;
        }
        this.ui.playerId = `player_${Date.now()}`;
        try {
          const joinResponse = yield this.joinGame(playerName);
          if (joinResponse.success) {
            this.ui.gameState = GameState.from(joinResponse.data);
            this.storeGameState();
            this.showGameInterface();
            this.updateGameDisplay();
            this.addLogEntry(`${playerName} joined the game`);
          } else {
            console.log("Join failed, attempting to create new game...");
            const createResponse = yield this.createGame(playerName);
            if (createResponse.success) {
              this.ui.gameState = GameState.from(createResponse.data);
              this.storeGameState();
              this.showGameInterface();
              this.updateGameDisplay();
              this.addLogEntry(`Game created by ${playerName}`);
            } else {
              throw new Error(createResponse.message || "Failed to create game");
            }
          }
        } catch (error) {
          console.error("Error joining game:", error);
          this.showError("Failed to join or create game. Please try again.");
        }
      });
    }
    joinGame(playerName) {
      return __async(this, null, function* () {
        if (!this.ui.connect4Client) {
          throw new Error("WASM client not initialized");
        }
        return yield this.ui.connect4Client.callMethod("connect4Service.joinGame", {
          gameId: this.ui.gameId,
          playerId: this.ui.playerId,
          playerName
        });
      });
    }
    createGame(playerName) {
      return __async(this, null, function* () {
        if (!this.ui.connect4Client) {
          throw new Error("WASM client not initialized");
        }
        const gameConfig = {
          boardWidth: 7,
          boardHeight: 6,
          connectLength: 4,
          maxPlayers: 2,
          minPlayers: 2,
          allowMultipleWinners: false,
          moveTimeoutSeconds: 30
        };
        const response = yield this.ui.connect4Client.callMethod("connect4Service.createGame", {
          gameId: this.ui.gameId,
          playerId: this.ui.playerId,
          playerName,
          config: gameConfig
        });
        if (!response.success) {
          throw new Error(response.message || "Failed to create game");
        }
        return response;
      });
    }
    dropPiece(column) {
      return __async(this, null, function* () {
        if (!this.ui.connect4Client || !this.ui.gameState) {
          console.error("Game not properly initialized");
          return;
        }
        try {
          const response = yield this.ui.connect4Client.callMethod("connect4Service.dropPiece", {
            gameId: this.ui.gameId,
            playerId: this.ui.playerId,
            column
          });
          if (response.success) {
            this.ui.gameState = GameState.from(response.data);
            this.storeGameState();
            this.updateGameDisplay();
            if (this.ui.statefulProxy) {
              this.ui.statefulProxy.sendPatches([{
                operation: "update",
                path: "",
                value: this.ui.gameState,
                timestamp: Date.now(),
                source: this.ui.playerId
              }]);
            }
          } else {
            console.error("Failed to drop piece:", response.message);
          }
        } catch (error) {
          console.error("Error dropping piece:", error);
        }
      });
    }
    applyPatches(patches) {
      for (const patch of patches) {
        if (patch.operation === "update" && patch.value) {
          try {
            const newState = GameState.from(patch.value);
            if (newState) {
              this.ui.gameState = newState;
              this.storeGameState();
              this.updateGameDisplay();
              this.addLogEntry("Game state updated from another player");
            }
          } catch (error) {
            console.error("Failed to apply patch:", error);
          }
        }
      }
    }
    showJoinForm() {
      if (this.elements.joinGameForm) {
        this.elements.joinGameForm.classList.remove("hidden");
      }
      if (this.elements.gameInterface) {
        this.elements.gameInterface.classList.add("hidden");
      }
      if (this.elements.errorState) {
        this.elements.errorState.classList.add("hidden");
      }
    }
    showGameInterface() {
      if (this.elements.joinGameForm) {
        this.elements.joinGameForm.classList.add("hidden");
      }
      if (this.elements.gameInterface) {
        this.elements.gameInterface.classList.remove("hidden");
      }
      if (this.elements.errorState) {
        this.elements.errorState.classList.add("hidden");
      }
      this.initializeGameBoard();
    }
    showError(message) {
      if (this.elements.joinGameForm) {
        this.elements.joinGameForm.classList.add("hidden");
      }
      if (this.elements.gameInterface) {
        this.elements.gameInterface.classList.add("hidden");
      }
      if (this.elements.errorState) {
        this.elements.errorState.classList.remove("hidden");
        const errorMessage = this.elements.errorState.querySelector("p");
        if (errorMessage) {
          errorMessage.textContent = message;
        }
      }
    }
    initializeGameBoard() {
      var _a, _b, _c, _d;
      if (!this.elements.gameBoard || !((_a = this.ui.gameState) == null ? void 0 : _a.board)) return;
      const rows = ((_b = this.ui.gameState.config) == null ? void 0 : _b.boardHeight) || 6;
      const cols = ((_c = this.ui.gameState.config) == null ? void 0 : _c.boardWidth) || 7;
      let boardHTML = "";
      for (let row = 0; row < rows; row++) {
        boardHTML += '<div class="board-row">';
        for (let col = 0; col < cols; col++) {
          const cellValue = ((_d = this.ui.gameState.board.rows[row]) == null ? void 0 : _d.cells[col]) || "";
          const pieceClass = cellValue ? `piece-${cellValue}` : "";
          boardHTML += `
                    <div class="board-cell ${pieceClass}" 
                         data-row="${row}" 
                         data-col="${col}"
                         onclick="gameViewer.dropPiece(${col})">
                        ${cellValue ? "\u25CF" : ""}
                    </div>
                `;
        }
        boardHTML += "</div>";
      }
      this.elements.gameBoard.innerHTML = boardHTML;
    }
    updateGameDisplay() {
      if (!this.ui.gameState) return;
      if (this.elements.turnNumber) {
        this.elements.turnNumber.textContent = this.ui.gameState.turnNumber.toString();
      }
      const currentPlayer = this.ui.gameState.players.find((p) => p.id === this.ui.gameState.currentPlayerId);
      if (this.elements.currentPlayerName && currentPlayer) {
        this.elements.currentPlayerName.textContent = currentPlayer.name;
      }
      if (this.elements.gameStatus) {
        let statusText = "In Progress";
        if (this.ui.gameState.status === 2) {
          statusText = this.ui.gameState.winners.length > 0 ? `Game Over - Winner: ${this.ui.gameState.winners.join(", ")}` : "Game Over - Draw";
        } else if (this.ui.gameState.players.length < 2) {
          statusText = "Waiting for players...";
        }
        this.elements.gameStatus.textContent = statusText;
      }
      if (this.elements.playersList) {
        this.elements.playersList.innerHTML = this.ui.gameState.players.map((player) => `
                <div class="player-item ${player.id === this.ui.gameState.currentPlayerId ? "current" : ""}">
                    <span class="player-name">${player.name}</span>
                    <span class="player-color player-color-${player.color || "red"}">\u25CF</span>
                </div>
            `).join("");
      }
      this.initializeGameBoard();
    }
    storeGameState() {
      if (this.ui.gameState) {
        localStorage.setItem(`connect4_game_${this.ui.gameId}`, JSON.stringify(this.ui.gameState));
        localStorage.setItem(`connect4_player_${this.ui.gameId}`, this.ui.playerId);
      }
    }
    addLogEntry(message) {
      if (!this.elements.gameLog) return;
      const timestamp = (/* @__PURE__ */ new Date()).toLocaleTimeString();
      const logEntry = document.createElement("div");
      logEntry.className = "log-entry";
      logEntry.innerHTML = `<span class="timestamp">[${timestamp}]</span> ${message}`;
      this.elements.gameLog.appendChild(logEntry);
      this.elements.gameLog.scrollTop = this.elements.gameLog.scrollHeight;
    }
    isValidGameId(gameId) {
      if (!gameId || gameId.length === 0 || gameId.length > 50) {
        return false;
      }
      return /^[a-zA-Z0-9-]+$/.test(gameId);
    }
    // Public methods for HTML onclick handlers
    resetGame() {
      if (confirm("Are you sure you want to start a new game?")) {
        localStorage.removeItem(`connect4_game_${this.ui.gameId}`);
        localStorage.removeItem(`connect4_player_${this.ui.gameId}`);
        window.location.reload();
      }
    }
    leaveGame() {
      if (confirm("Are you sure you want to leave this game?")) {
        localStorage.removeItem(`connect4_game_${this.ui.gameId}`);
        localStorage.removeItem(`connect4_player_${this.ui.gameId}`);
        window.location.href = "/";
      }
    }
  };
  var gameViewer = new GameViewer();
  window.gameViewer = gameViewer;
  window.joinCurrentGame = () => {
    const form = document.getElementById("joinGameForm");
    if (form) {
      form.dispatchEvent(new Event("submit"));
    }
  };
  window.resetGame = () => gameViewer.resetGame();
  window.leaveGame = () => gameViewer.leaveGame();
  var gameViewer_default = gameViewer;
  return __toCommonJS(gameViewer_exports);
})();
