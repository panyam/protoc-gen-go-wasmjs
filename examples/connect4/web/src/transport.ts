// =============================================================================
// STATEFUL PROXY TRANSPORT SYSTEM - Pluggable Architecture
// =============================================================================

// Abstract base class for transport implementations
export abstract class StatefulTransport {
    protected gameId: string;
    protected onPatchReceived: ((patches: any[]) => void) | null = null;

    constructor(gameId: string) {
        this.gameId = gameId;
    }
    
    // Send patches to other clients/tabs
    abstract sendPatches(patches: any[]): Promise<void>;
    
    // Subscribe to incoming patches
    subscribe(callback: (patches: any[]) => void): void {
        this.onPatchReceived = callback;
    }
    
    // Clean up resources
    destroy(): void {
        // Override in implementations
    }
}

// IndexedDB + Polling transport for persistent cross-page communication
export class IndexedDBTransport extends StatefulTransport {
    private pollInterval: number;
    private dbName: string;
    private patchesStoreName: string;
    private gameStatesStoreName: string;
    private db: IDBDatabase | null = null;
    private pollTimer: number | null = null;
    private gameStatesPollTimer: number | null = null;
    private lastProcessedId: number = 0;
    private lastGameStateTimestamp: number = 0;
    private onGameStateChange: ((gameState: any) => void) | null = null;

    constructor(gameId: string, pollInterval: number = 1000) {
        super(gameId);
        this.pollInterval = pollInterval;
        this.dbName = `connect4_transport_${gameId}`;
        this.patchesStoreName = 'patches';
        this.gameStatesStoreName = 'gameStates';
    }

    async init(): Promise<void> {
        return new Promise((resolve, reject) => {
            const request = indexedDB.open(this.dbName, 1);
            
            request.onerror = () => {
              console.log("Error opening DB: ", request.error)
              reject(request.error);
            }
            request.onsuccess = () => {
              console.log("Successfully opened IndexedDB: ", request)
                this.db = request.result;
                this.startPolling();
                resolve();
            };
            
            request.onupgradeneeded = (event) => {
                console.log("Upgrading....")
                const db = (event.target as IDBOpenDBRequest).result;
                
                // Create patches store (existing functionality)
                if (!db.objectStoreNames.contains(this.patchesStoreName)) {
                    const patchesStore = db.createObjectStore(this.patchesStoreName, { keyPath: 'game_id', autoIncrement: true });
                    patchesStore.createIndex('timestamp', 'timestamp', { unique: false });
                    patchesStore.createIndex('game_id', 'game_id', { unique: false });
                }
                
                // Create game states store (new functionality)
                if (!db.objectStoreNames.contains(this.gameStatesStoreName)) {
                    const gameStatesStore = db.createObjectStore(this.gameStatesStoreName, { keyPath: 'gameId' });
                    gameStatesStore.createIndex('timestamp', 'timestamp', { unique: false });
                    gameStatesStore.createIndex('lastModified', 'lastModified', { unique: false });
                }
            };
        });
    }

    async sendPatches(patches: any[]): Promise<void> {
        if (!this.db) throw new Error('Database not initialized');
        
        const transaction = this.db.transaction([this.patchesStoreName], 'readwrite');
        const store = transaction.objectStore(this.patchesStoreName);
        
        const patchData = {
            gameId: this.gameId,
            patches,
            timestamp: Date.now(),
            tabId: this.getTabId()
        };
        
        return new Promise((resolve, reject) => {
            const request = store.add(patchData);
            request.onsuccess = () => resolve();
            request.onerror = () => reject(request.error);
        });
    }

    // Save complete game state to IndexedDB
    async saveGameState(gameId: string, gameState: any): Promise<IDBValidKey> {
        if (!this.db) throw new Error('Database not initialized');
        console.log("Saving game state to IndexedDB:", { gameId, gameState });
        
        const transaction = this.db!.transaction([this.gameStatesStoreName], 'readwrite');
        const store = transaction.objectStore(this.gameStatesStoreName);
        
        const gameStateData = {
            gameId,
            gameState,
            timestamp: Date.now(),
            lastModified: Date.now(),
            tabId: this.getTabId()
        };
        
        console.log("Game state data to save:", gameStateData);
        const request = store.put(gameStateData); // Use put to update existing
        console.log("Payload: ", gameStateData)
        return new Promise((resolve, reject) => {
            request.onerror = () => {
              console.log("Error saving DB: ", this.db, request.error)
              reject(request.error);
            }
            request.onsuccess = () => {
              console.log("Successfully saved: ", request)
                resolve(request.result);
            };
        })
    }

    // Load complete game state from IndexedDB
    async loadGameState(gameId: string): Promise<any | null> {
        if (!this.db) throw new Error('Database not initialized');
        console.log("Loading game state from IndexedDB for gameId:", gameId);
        
        const transaction = this.db.transaction([this.gameStatesStoreName], 'readonly');
        const store = transaction.objectStore(this.gameStatesStoreName);
        const request = store.get(gameId);
        
        return new Promise((resolve, reject) => {
            request.onerror = () => {
              console.log("Error loading game from DB: ", this.db, request.error)
              reject(request.error);
            }
            request.onsuccess = () => {
              console.log("Load request result:", request.result);
              const result = request.result;
              if (result) {
                console.log("Found game state:", result.gameState);
                resolve(result.gameState);
              } else {
                console.log("No game state found for gameId:", gameId);
                resolve(null);
              }
            };
        })
    }

    // Debug method to see all stored games
    async debugListAllGames(): Promise<void> {
        if (!this.db) throw new Error('Database not initialized');
        console.log("=== DEBUG: Listing all games in IndexedDB ===");
        
        const transaction = this.db.transaction([this.gameStatesStoreName], 'readonly');
        const store = transaction.objectStore(this.gameStatesStoreName);
        const request = store.getAll();
        
        return new Promise((resolve, reject) => {
            request.onerror = () => {
                console.log("Error listing all games:", request.error);
                reject(request.error);
            };
            request.onsuccess = () => {
                console.log("All stored games:", request.result);
                request.result.forEach((item: any, index: number) => {
                    console.log(`Game ${index}:`, {
                        gameId: item.gameId,
                        hasGameState: !!item.gameState,
                        timestamp: item.timestamp,
                        gameState: item.gameState
                    });
                });
                resolve();
            };
        });
    }

    // Set callback for game state changes
    onGameStateChanged(callback: (gameState: any) => void): void {
        this.onGameStateChange = callback;
        this.startGameStatePolling();
    }

    private startPolling(): void {
        if (this.pollTimer) {
            clearInterval(this.pollTimer);
        }
        this.pollTimer = window.setInterval(() => {
            this.checkForNewPatches();
        }, this.pollInterval);
    }

    private startGameStatePolling(): void {
        if (this.gameStatesPollTimer) {
            clearInterval(this.gameStatesPollTimer);
        }
        
        this.gameStatesPollTimer = window.setInterval(() => {
            this.checkForGameStateChanges();
        }, this.pollInterval);
    }

    private async checkForNewPatches(): Promise<void> {
        if (!this.db || !this.onPatchReceived) return;
        
        const transaction = this.db.transaction([this.patchesStoreName], 'readonly');
        const store = transaction.objectStore(this.patchesStoreName);
        const gameIndex = store.index('game_id');
        
        const request = gameIndex.getAll(this.gameId);
        request.onsuccess = () => {
            const results = request.result;
            const newPatches = results.filter(item => 
                item.id > this.lastProcessedId && 
                item.tabId !== this.getTabId()
            );
            
            if (newPatches.length > 0) {
                // Sort by ID to maintain order
                newPatches.sort((a, b) => a.id - b.id);
                
                for (const patchData of newPatches) {
                    this.onPatchReceived!(patchData.patches);
                    this.lastProcessedId = Math.max(this.lastProcessedId, patchData.id);
                }
            }
        };
    }

    private async checkForGameStateChanges(): Promise<void> {
        if (!this.db || !this.onGameStateChange) return;
        
        const transaction = this.db.transaction([this.gameStatesStoreName], 'readonly');
        const store = transaction.objectStore(this.gameStatesStoreName);
        
        const request = store.get(this.gameId);
        request.onsuccess = () => {
            const result = request.result;
            if (result && 
                result.lastModified > this.lastGameStateTimestamp && 
                result.tabId !== this.getTabId()) {
                
                this.lastGameStateTimestamp = result.lastModified;
                this.onGameStateChange!(result.gameState);
            }
        };
    }

    private getTabId(): string {
        if (!window.tabId) {
            window.tabId = Math.random().toString(36).substr(2, 9);
        }
        return window.tabId;
    }

    destroy(): void {
        if (this.pollTimer) {
            clearInterval(this.pollTimer);
            this.pollTimer = null;
        }
        if (this.gameStatesPollTimer) {
            clearInterval(this.gameStatesPollTimer);
            this.gameStatesPollTimer = null;
        }
        if (this.db) {
            this.db.close();
            this.db = null;
        }
    }
}

// BroadcastChannel transport for same-origin tab communication
export class BroadcastChannelTransport extends StatefulTransport {
    private channel: BroadcastChannel | null = null;

    constructor(gameId: string) {
        super(gameId);
        this.channel = new BroadcastChannel(`connect4_${gameId}`);
        this.channel.onmessage = (event) => {
            if (this.onPatchReceived && event.data.type === 'patches') {
                this.onPatchReceived(event.data.patches);
            }
        };
    }

    async sendPatches(patches: any[]): Promise<void> {
        if (this.channel) {
            this.channel.postMessage({
                type: 'patches',
                patches,
                timestamp: Date.now()
            });
        }
    }

    destroy(): void {
        if (this.channel) {
            this.channel.close();
            this.channel = null;
        }
    }
}

// Transport factory
export class TransportFactory {
    static create(gameId: string, type: 'indexeddb' | 'broadcast' = 'indexeddb'): StatefulTransport {
        switch (type) {
            case 'indexeddb':
                return new IndexedDBTransport(gameId);
            case 'broadcast':
                return new BroadcastChannelTransport(gameId);
            default:
                throw new Error(`Unknown transport type: ${type}`);
        }
    }
}

// Extend Window interface for tabId
declare global {
    interface Window {
        tabId: string;
    }
}
