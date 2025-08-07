"use strict";
var Connect4Index = (() => {
  var __defProp = Object.defineProperty;
  var __getOwnPropDesc = Object.getOwnPropertyDescriptor;
  var __getOwnPropNames = Object.getOwnPropertyNames;
  var __hasOwnProp = Object.prototype.hasOwnProperty;
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

  // src/index.ts
  var index_exports = {};
  __export(index_exports, {
    default: () => index_default
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

  // src/index.ts
  var GamesListManager = class {
    constructor() {
      this.gamesContainer = null;
      this.createForm = null;
      this.connect4Client = null;
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
      this.gamesContainer = document.getElementById("gamesList");
      this.createForm = document.getElementById("createGameForm");
      if (this.createForm) {
        this.createForm.addEventListener("submit", (e) => this.handleCreateGame(e));
      }
      this.loadExistingGames();
      this.initializeWasmClient();
    }
    initializeWasmClient() {
      return __async(this, null, function* () {
        try {
          console.log("Initializing WASM client for games list...");
          this.connect4Client = new multiplayer_connect4Client_client_default();
          yield this.connect4Client.loadWasm("/static/wasm/multiplayer_connect4.wasm");
          yield this.connect4Client.waitUntilReady();
          console.log("WASM client ready for game operations");
        } catch (error) {
          console.error("Failed to initialize WASM client:", error);
        }
      });
    }
    loadExistingGames() {
      if (!this.gamesContainer) return;
      const games = this.getStoredGames();
      if (games.length === 0) {
        this.gamesContainer.innerHTML = `
                <div class="no-games">
                    <p>No games found. Create your first game!</p>
                </div>
            `;
        return;
      }
      this.gamesContainer.innerHTML = games.map((game) => `
            <div class="game-item" data-game-id="${game.gameId}">
                <div class="game-info">
                    <h3>${game.gameId}</h3>
                    <p>Player: ${game.playerName}</p>
                    <p>Status: ${game.gameStatus || "Unknown"}</p>
                    <small>Last played: ${new Date(game.lastPlayed).toLocaleString()}</small>
                </div>
                <div class="game-actions">
                    <a href="/${game.gameId}" class="btn">Continue Game</a>
                    <button class="btn btn-secondary" onclick="gamesManager.removeGame('${game.gameId}')">Remove</button>
                </div>
            </div>
        `).join("");
    }
    handleCreateGame(event) {
      return __async(this, null, function* () {
        event.preventDefault();
        const formData = new FormData(event.target);
        const gameId = formData.get("gameId");
        const playerName = formData.get("playerName");
        if (!gameId || !playerName) {
          alert("Please fill in all fields");
          return;
        }
        try {
          if (!this.isValidGameId(gameId)) {
            alert("Game ID can only contain letters, numbers, and hyphens");
            return;
          }
          this.storeGame({
            gameId,
            playerName,
            lastPlayed: Date.now(),
            gameStatus: "Created"
          });
          window.location.href = `/${gameId}`;
        } catch (error) {
          console.error("Error creating game:", error);
          alert("Failed to create game. Please try again.");
        }
      });
    }
    getStoredGames() {
      try {
        const gamesData = localStorage.getItem("connect4Games");
        return gamesData ? JSON.parse(gamesData) : [];
      } catch (error) {
        console.error("Error loading stored games:", error);
        return [];
      }
    }
    storeGame(game) {
      try {
        const games = this.getStoredGames();
        const existingIndex = games.findIndex((g) => g.gameId === game.gameId);
        if (existingIndex >= 0) {
          games[existingIndex] = game;
        } else {
          games.push(game);
        }
        localStorage.setItem("connect4Games", JSON.stringify(games));
        this.loadExistingGames();
      } catch (error) {
        console.error("Error storing game:", error);
      }
    }
    removeGame(gameId) {
      try {
        const games = this.getStoredGames();
        const filteredGames = games.filter((g) => g.gameId !== gameId);
        localStorage.setItem("connect4Games", JSON.stringify(filteredGames));
        this.loadExistingGames();
      } catch (error) {
        console.error("Error removing game:", error);
      }
    }
    isValidGameId(gameId) {
      if (!gameId || gameId.length === 0 || gameId.length > 50) {
        return false;
      }
      return /^[a-zA-Z0-9-]+$/.test(gameId);
    }
  };
  var gamesManager = new GamesListManager();
  window.gamesManager = gamesManager;
  var index_default = gamesManager;
  return __toCommonJS(index_exports);
})();
