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
            
            request.onerror = () => reject(request.error);
            request.onsuccess = () => {
                this.db = request.result;
                this.startPolling();
                resolve();
            };
            
            request.onupgradeneeded = (event) => {
                const db = (event.target as IDBOpenDBRequest).result;
                
                // Create patches store (existing functionality)
                if (!db.objectStoreNames.contains(this.patchesStoreName)) {
                    const patchesStore = db.createObjectStore(this.patchesStoreName, { keyPath: 'id', autoIncrement: true });
                    patchesStore.createIndex('timestamp', 'timestamp', { unique: false });
                    patchesStore.createIndex('gameId', 'gameId', { unique: false });
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
    async saveGameState(gameId: string, gameState: any): Promise<void> {
        if (!this.db) throw new Error('Database not initialized');
        
        const transaction = this.db.transaction([this.gameStatesStoreName], 'readwrite');
        const store = transaction.objectStore(this.gameStatesStoreName);
        
        const gameStateData = {
            gameId,
            gameState,
            timestamp: Date.now(),
            lastModified: Date.now(),
            tabId: this.getTabId()
        };
        
        return new Promise((resolve, reject) => {
            const request = store.put(gameStateData); // Use put to update existing
            request.onsuccess = () => resolve();
            request.onerror = () => reject(request.error);
        });
    }

    // Load complete game state from IndexedDB
    async loadGameState(gameId: string): Promise<any | null> {
        if (!this.db) throw new Error('Database not initialized');
        
        const transaction = this.db.transaction([this.gameStatesStoreName], 'readonly');
        const store = transaction.objectStore(this.gameStatesStoreName);
        
        return new Promise((resolve, reject) => {
            const request = store.get(gameId);
            request.onsuccess = () => {
                const result = request.result;
                resolve(result ? result.gameState : null);
            };
            request.onerror = () => reject(request.error);
        });
    }

    // Set callback for game state changes
    onGameStateChanged(callback: (gameState: any) => void): void {
        this.onGameStateChange = callback;
        this.startGameStatePolling();
    }

    private startPolling(): void {
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
        const gameIndex = store.index('gameId');
        
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
