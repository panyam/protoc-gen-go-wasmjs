// =============================================================================
// STATEFUL PROXY TRANSPORT SYSTEM - Pluggable Architecture
// =============================================================================

// Abstract base class for transport implementations
class StatefulTransport {
    constructor(gameId) {
        this.gameId = gameId;
        this.onPatchReceived = null; // Callback for incoming patches
    }
    
    // Send patches to other clients/tabs
    async sendPatches(patches) {
        throw new Error('sendPatches must be implemented by transport');
    }
    
    // Subscribe to incoming patches
    subscribe(callback) {
        this.onPatchReceived = callback;
    }
    
    // Clean up resources
    destroy() {
        // Override in implementations
    }
}

// IndexedDB + Polling transport for persistent cross-page communication
class IndexedDBTransport extends StatefulTransport {
    constructor(gameId, pollInterval = 1000) {
        super(gameId);
        this.pollInterval = pollInterval;
        this.dbName = `connect4-stateful-${gameId}`;
        this.storeName = 'patches';
        this.db = null;
        this.lastPatchId = 0;
        this.pollTimer = null;
        this.isPolling = false;
        
        this.initDB();
    }
    
    async initDB() {
        return new Promise((resolve, reject) => {
            const request = indexedDB.open(this.dbName, 1);
            
            request.onerror = () => reject(request.error);
            request.onsuccess = () => {
                this.db = request.result;
                this.startPolling();
                resolve();
            };
            
            request.onupgradeneeded = (event) => {
                const db = event.target.result;
                if (!db.objectStoreNames.contains(this.storeName)) {
                    const store = db.createObjectStore(this.storeName, { 
                        keyPath: 'id', 
                        autoIncrement: true 
                    });
                    store.createIndex('timestamp', 'timestamp', { unique: false });
                    store.createIndex('gameId', 'gameId', { unique: false });
                }
            };
        });
    }
    
    async sendPatches(patches) {
        if (!this.db || !patches || patches.length === 0) return;
        
        const transaction = this.db.transaction([this.storeName], 'readwrite');
        const store = transaction.objectStore(this.storeName);
        
        const patchRecord = {
            gameId: this.gameId,
            patches: patches,
            timestamp: Date.now(),
            source: 'indexeddb'
        };
        
        try {
            await new Promise((resolve, reject) => {
                const request = store.add(patchRecord);
                request.onsuccess = () => resolve();
                request.onerror = () => reject(request.error);
            });
            
            console.log(`Stored ${patches.length} patches to IndexedDB`);
        } catch (error) {
            console.error('Failed to store patches:', error);
        }
    }
    
    async pollForPatches() {
        if (!this.db || this.isPolling) return;
        
        this.isPolling = true;
        
        try {
            const transaction = this.db.transaction([this.storeName], 'readonly');
            const store = transaction.objectStore(this.storeName);
            const index = store.index('gameId');
            
            const patches = await new Promise((resolve, reject) => {
                const patches = [];
                const request = index.openCursor(IDBKeyRange.only(this.gameId));
                
                request.onsuccess = (event) => {
                    const cursor = event.target.result;
                    if (cursor) {
                        const record = cursor.value;
                        if (record.id > this.lastPatchId) {
                            patches.push(...record.patches);
                            this.lastPatchId = record.id;
                        }
                        cursor.continue();
                    } else {
                        resolve(patches);
                    }
                };
                request.onerror = () => reject(request.error);
            });
            
            if (patches.length > 0 && this.onPatchReceived) {
                console.log(`Polled ${patches.length} new patches from IndexedDB`);
                this.onPatchReceived(patches);
            }
        } catch (error) {
            console.error('Failed to poll patches:', error);
        } finally {
            this.isPolling = false;
        }
    }
    
    startPolling() {
        if (this.pollTimer) return;
        
        console.log(`Started IndexedDB polling every ${this.pollInterval}ms for game ${this.gameId}`);
        this.pollTimer = setInterval(() => {
            this.pollForPatches();
        }, this.pollInterval);
        
        // Do an initial poll
        this.pollForPatches();
    }
    
    stopPolling() {
        if (this.pollTimer) {
            clearInterval(this.pollTimer);
            this.pollTimer = null;
            console.log(`Stopped IndexedDB polling for game ${this.gameId}`);
        }
    }
    
    destroy() {
        this.stopPolling();
        if (this.db) {
            this.db.close();
        }
    }
}

// BroadcastChannel transport for cross-tab communication (kept for comparison)
class BroadcastChannelTransport extends StatefulTransport {
    constructor(gameId) {
        super(gameId);
        this.channel = new BroadcastChannel(`connect4-stateful-${gameId}`);
        this.channel.onmessage = (event) => this.handleMessage(event);
    }
    
    async sendPatches(patches) {
        this.channel.postMessage({
            type: 'patches',
            gameId: this.gameId,
            patches: patches,
            timestamp: Date.now(),
            source: 'broadcast'
        });
    }
    
    handleMessage(event) {
        const { type, gameId, patches } = event.data;
        if (type === 'patches' && gameId === this.gameId && this.onPatchReceived) {
            this.onPatchReceived(patches);
        }
    }
    
    destroy() {
        this.channel.close();
    }
}

// WebSocket transport for server-based collaboration
class WebSocketTransport extends StatefulTransport {
    constructor(gameId, wsUrl = null) {
        super(gameId);
        this.wsUrl = wsUrl || `ws://localhost:8080/game/${gameId}`;
        this.ws = null;
        this.reconnectAttempts = 0;
        this.maxReconnectAttempts = 5;
        this.connect();
    }
    
    async connect() {
        try {
            this.ws = new WebSocket(this.wsUrl);
            
            this.ws.onopen = () => {
                console.log(`WebSocket connected for game ${this.gameId}`);
                this.reconnectAttempts = 0;
            };
            
            this.ws.onmessage = (event) => {
                const data = JSON.parse(event.data);
                if (data.type === 'patches' && this.onPatchReceived) {
                    this.onPatchReceived(data.patches);
                }
            };
            
            this.ws.onclose = () => {
                this.handleReconnect();
            };
            
            this.ws.onerror = (error) => {
                console.error('WebSocket error:', error);
            };
        } catch (error) {
            console.error('Failed to connect WebSocket:', error);
            this.handleReconnect();
        }
    }
    
    async sendPatches(patches) {
        if (this.ws && this.ws.readyState === WebSocket.OPEN) {
            this.ws.send(JSON.stringify({
                type: 'patches',
                gameId: this.gameId,
                patches: patches,
                timestamp: Date.now(),
                source: 'websocket'
            }));
        } else {
            console.warn('WebSocket not ready, patches not sent');
        }
    }
    
    handleReconnect() {
        if (this.reconnectAttempts < this.maxReconnectAttempts) {
            this.reconnectAttempts++;
            const delay = Math.pow(2, this.reconnectAttempts) * 1000; // Exponential backoff
            console.log(`Reconnecting WebSocket in ${delay}ms (attempt ${this.reconnectAttempts})`);
            setTimeout(() => this.connect(), delay);
        }
    }
    
    destroy() {
        if (this.ws) {
            this.ws.close();
        }
    }
}

// Server-Sent Events transport for real-time updates from server
class SSETransport extends StatefulTransport {
    constructor(gameId, sseUrl = null) {
        super(gameId);
        this.sseUrl = sseUrl || `/events/game/${gameId}`;
        this.eventSource = null;
        this.connect();
    }
    
    connect() {
        this.eventSource = new EventSource(this.sseUrl);
        
        this.eventSource.onmessage = (event) => {
            const data = JSON.parse(event.data);
            if (data.type === 'patches' && this.onPatchReceived) {
                this.onPatchReceived(data.patches);
            }
        };
        
        this.eventSource.onerror = (error) => {
            console.error('SSE error:', error);
        };
    }
    
    async sendPatches(patches) {
        // SSE is receive-only, send via regular HTTP POST
        try {
            await fetch(`/api/game/${this.gameId}/patches`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ patches })
            });
        } catch (error) {
            console.error('Failed to send patches via HTTP:', error);
        }
    }
    
    destroy() {
        if (this.eventSource) {
            this.eventSource.close();
        }
    }
}

// =============================================================================
// STATEFUL PROXY - Transport-Agnostic Implementation
// =============================================================================

class StatefulProxy {
    constructor(gameId, transportType = 'indexeddb') {
        this.gameId = gameId;
        this.gameState = null;
        this.patches = [];
        this.lastChangeNumber = 0;
        
        // Initialize pluggable transport
        this.transport = this.createTransport(transportType);
        this.transport.subscribe((patches) => this.handleIncomingPatches(patches));
        
        // Storage keys for localStorage backup
        this.storageKey = `connect4_game_${gameId}`;
        this.patchesKey = `connect4_patches_${gameId}`;
        
        // UI update callback
        this.onStateChanged = null;
    }
    
    // Factory method for creating transport instances
    createTransport(type) {
        switch (type) {
            case 'indexeddb':
                return new IndexedDBTransport(this.gameId, 1000); // Poll every 1 second
            case 'broadcast':
                return new BroadcastChannelTransport(this.gameId);
            case 'websocket':
                return new WebSocketTransport(this.gameId);
            case 'sse':
                return new SSETransport(this.gameId);
            default:
                console.warn(`Unknown transport type: ${type}, falling back to IndexedDB`);
                return new IndexedDBTransport(this.gameId, 1000);
        }
    }
    
    // Switch transport at runtime
    async switchTransport(newType) {
        const oldTransport = this.transport;
        this.transport = this.createTransport(newType);
        this.transport.subscribe((patches) => this.handleIncomingPatches(patches));
        
        // Clean up old transport
        oldTransport.destroy();
        
        console.log(`Switched transport from ${oldTransport.constructor.name} to ${this.transport.constructor.name}`);
    }
    
    // Set the game state and initialize proxy
    setGameState(gameState) {
        this.gameState = gameState;
        this.saveState();
        this.notifyStateChanged();
    }
    
    // Apply patches from WASM service calls
    async applyPatchesFromWasm(patches) {
        if (!patches || patches.length === 0) return;
        
        // Apply patches locally
        this.applyPatches(patches);
        
        // Broadcast to other clients
        await this.transport.sendPatches(patches);
        
        // Store patches for conflict resolution
        this.patches.push(...patches);
        this.saveState();
    }
    
    // Handle incoming patches from transport
    handleIncomingPatches(patches) {
        console.log(`Received ${patches.length} patches from transport`);
        this.applyPatches(patches);
        this.saveState();
    }
    
    // Apply patches to local game state
    applyPatches(patches) {
        if (!this.gameState || !patches || patches.length === 0) return;
        
        // Sort patches by change number for correct ordering
        patches.sort((a, b) => (a.changeNumber || 0) - (b.changeNumber || 0));
        
        patches.forEach(patch => {
            if (patch.changeNumber <= this.lastChangeNumber) {
                console.log(`Skipping old patch ${patch.changeNumber} (current: ${this.lastChangeNumber})`);
                return;
            }
            
            this.applyPatch(patch);
            this.lastChangeNumber = Math.max(this.lastChangeNumber, patch.changeNumber || 0);
        });
        
        this.notifyStateChanged();
    }
    
    // Apply a single patch to the game state
    applyPatch(patch) {
        const fieldPath = patch.fieldPath || patch.field_path;
        const value = this.parseJsonValue(patch.valueJson || patch.value_json);
        
        console.log(`Applying patch: ${patch.operation} ${fieldPath} = ${patch.valueJson || patch.value_json}`);
        
        switch (patch.operation) {
            case 0: // SET
                this.setNestedValue(this.gameState, fieldPath, value);
                break;
            case 1: // INSERT_LIST
                this.insertIntoArray(this.gameState, fieldPath, value, patch.index);
                break;
            case 2: // REMOVE_LIST
                this.removeFromArray(this.gameState, fieldPath, patch.index);
                break;
            default:
                console.warn(`Unknown patch operation: ${patch.operation}`);
        }
    }
    
    // Parse JSON value from patch
    parseJsonValue(jsonValue) {
        try {
            return JSON.parse(jsonValue);
        } catch (e) {
            return jsonValue;
        }
    }
    
    // Set nested object value by path (e.g., "board.rows[0].cells[1]")
    setNestedValue(obj, path, value) {
        const parts = path.split(/[\.\[\]]+/).filter(Boolean);
        let current = obj;
        
        for (let i = 0; i < parts.length - 1; i++) {
            const part = parts[i];
            const nextPart = parts[i + 1];
            
            if (!current[part]) {
                current[part] = isNaN(nextPart) ? {} : [];
            }
            current = current[part];
        }
        
        const lastPart = parts[parts.length - 1];
        current[lastPart] = value;
    }
    
    // Insert into array by path
    insertIntoArray(obj, path, value, index = -1) {
        const parts = path.split(/[\.\[\]]+/).filter(Boolean);
        let current = obj;
        
        for (let i = 0; i < parts.length; i++) {
            if (i === parts.length - 1) {
                if (!Array.isArray(current[parts[i]])) {
                    current[parts[i]] = [];
                }
                if (index >= 0) {
                    current[parts[i]].splice(index, 0, value);
                } else {
                    current[parts[i]].push(value);
                }
            } else {
                current = current[parts[i]];
            }
        }
    }
    
    // Remove from array by path
    removeFromArray(obj, path, index) {
        const parts = path.split(/[\.\[\]]+/).filter(Boolean);
        let current = obj;
        
        for (let i = 0; i < parts.length; i++) {
            if (i === parts.length - 1) {
                if (Array.isArray(current[parts[i]]) && index >= 0) {
                    current[parts[i]].splice(index, 1);
                }
            } else {
                current = current[parts[i]];
            }
        }
    }
    
    // Save current state to localStorage
    saveState() {
        if (this.gameState) {
            localStorage.setItem(this.storageKey, JSON.stringify(this.gameState));
            localStorage.setItem(this.patchesKey, JSON.stringify(this.patches));
        }
    }
    
    // Load state from localStorage
    loadState() {
        const savedState = localStorage.getItem(this.storageKey);
        const savedPatches = localStorage.getItem(this.patchesKey);
        
        if (savedState) {
            this.gameState = JSON.parse(savedState);
        }
        if (savedPatches) {
            this.patches = JSON.parse(savedPatches);
            this.lastChangeNumber = Math.max(...this.patches.map(p => p.changeNumber || 0), 0);
        }
        
        return this.gameState;
    }
    
    // Set callback for state changes
    onStateChange(callback) {
        this.onStateChanged = callback;
    }
    
    // Notify that state has changed
    notifyStateChanged() {
        if (this.onStateChanged) {
            this.onStateChanged(this.gameState);
        }
    }
    
    // Clean up resources
    destroy() {
        if (this.transport) {
            this.transport.destroy();
        }
    }
}

// Global function to switch transports for debugging
window.switchTransport = function(transportType) {
    if (window.statefulProxy) {
        console.log(`Switching from ${window.statefulProxy.transport.constructor.name} to ${transportType}`);
        window.statefulProxy.switchTransport(transportType);
    } else {
        console.warn('No active stateful proxy found');
    }
};

// Global function to show transport status
window.showTransportStatus = function() {
    if (window.statefulProxy) {
        const transport = window.statefulProxy.transport;
        console.log(`Current transport: ${transport.constructor.name}`);
        if (transport instanceof IndexedDBTransport) {
            console.log(`- Polling interval: ${transport.pollInterval}ms`);
            console.log(`- Last patch ID: ${transport.lastPatchId}`);
            console.log(`- DB name: ${transport.dbName}`);
        }
    } else {
        console.warn('No active stateful proxy found');
    }
};

// Export classes for module usage
if (typeof module !== 'undefined' && module.exports) {
    module.exports = {
        StatefulProxy,
        StatefulTransport,
        IndexedDBTransport,
        BroadcastChannelTransport,
        WebSocketTransport,
        SSETransport
    };
}
