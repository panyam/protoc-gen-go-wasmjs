// Entry point for the games list page (/)
// Handles game discovery, creation, and navigation

import Connect4Client from '../gen/wasmts/multiplayer_connect4Client.client';
import { IndexedDBTransport } from './transport';

// Types for game storage
interface StoredGame {
    gameId: string;
    playerName: string;
    lastPlayed: number;
    gameStatus?: string;
}

class GamesListManager {
    private gamesContainer: HTMLElement | null = null;
    private createForm: HTMLFormElement | null = null;
    private connect4Client: Connect4Client | null = null;
    private numPlayersSelect: HTMLSelectElement | null = null;
    private boardWidthInput: HTMLInputElement | null = null;
    private boardHeightInput: HTMLInputElement | null = null;
    private dimensionsDisplay: HTMLElement | null = null;
    private storageTransport: IndexedDBTransport | null = null;

    constructor() {
        this.init();
    }

    private async init(): Promise<void> {
        // Wait for DOM to be ready
        if (document.readyState === 'loading') {
            document.addEventListener('DOMContentLoaded', () => this.initializeUI());
        } else {
            this.initializeUI();
        }
    }

    private initializeUI(): void {
        this.gamesContainer = document.getElementById('gamesList');
        this.createForm = document.getElementById('createGameForm') as HTMLFormElement;
        this.numPlayersSelect = document.getElementById('numPlayers') as HTMLSelectElement;
        this.boardWidthInput = document.getElementById('boardWidth') as HTMLInputElement;
        this.boardHeightInput = document.getElementById('boardHeight') as HTMLInputElement;
        this.dimensionsDisplay = document.getElementById('dimensionsDisplay');

        if (this.createForm) {
            this.createForm.addEventListener('submit', (e) => this.handleCreateGame(e));
        }

        if (this.numPlayersSelect) {
            this.numPlayersSelect.addEventListener('change', () => this.updateBoardDimensions());
        }

        if (this.boardWidthInput) {
            this.boardWidthInput.addEventListener('input', () => this.updateDimensionsDisplay());
        }

        if (this.boardHeightInput) {
            this.boardHeightInput.addEventListener('input', () => this.updateDimensionsDisplay());
        }

        this.updateBoardDimensions();
        this.initializeWasmClient().then(async () => {
            // Load and display games list from IndexedDB
            await this.loadExistingGames();
        });
    }

    private async initializeWasmClient(): Promise<void> {
        console.log('Initializing WASM client for games list...');
        this.connect4Client = new Connect4Client();
        
        // Initialize storage transport for global game state management
        this.storageTransport = new IndexedDBTransport('global');
        await this.storageTransport.init();
        
        // Load WASM first so that setWasmStorageCallbacks is available
        await this.connect4Client.loadWasm('/static/wasm/multiplayer_connect4.wasm');
        await this.connect4Client.waitUntilReady();
        
        // Set up storage callbacks after WASM is ready
        await this.setupStorageCallbacks();
        
        console.log('WASM client ready for game operations');
    }

    private async setupStorageCallbacks(): Promise<void> {
        if (!this.storageTransport || !(window as any).setWasmStorageCallbacks) {
            console.error('Storage transport or WASM callbacks not available:', {
                storageTransport: !!this.storageTransport,
                setWasmStorageCallbacks: !!(window as any).setWasmStorageCallbacks
            });
            return;
        }

        // Simple, clean async callbacks - let failures bubble up to WASM
        const saveCallback = async (gameId: string, gameStateJson: string) => {
            const gameState = JSON.parse(gameStateJson);
            console.log("Before Saving to IndexDB")
            await this.storageTransport!.saveGameState(gameId, gameState);
            console.log(`After Saved game ${gameId} to IndexedDB`);
        };

        const loadCallback = async (gameId: string) => {
            const gameState = await this.storageTransport!.loadGameState(gameId);
            if (gameState) {
                console.log(`Loaded game ${gameId} from IndexedDB`);
                return JSON.stringify(gameState);
            }
            return null; // Game not found
        };

        const pollCallback = (gameId: string) => {
            // Set up polling for external changes to this game
            this.storageTransport!.onGameStateChanged((gameState: any) => {
                if ((window as any).wasmOnExternalStorageChange) {
                    (window as any).wasmOnExternalStorageChange(gameId, JSON.stringify(gameState));
                }
            });
        };

        // Configure WASM with these callbacks
        const result = (window as any).setWasmStorageCallbacks(
            saveCallback,
            loadCallback, 
            pollCallback
        );
        
        if (result && result.success) {
            console.log('Storage callbacks configured successfully');
        } else {
            console.error('Failed to configure storage callbacks:', result);
        }
    }

    private async loadExistingGames(): Promise<void> {
        if (!this.gamesContainer) return;

        // Load games directly from IndexedDB
        const games = await this.getAllGamesFromIndexedDB();
        
        if (games.length === 0) {
            this.gamesContainer.innerHTML = `
                <div class="no-games">
                    <p>No games found. Create your first game!</p>
                </div>
            `;
            return;
        }

        console.log("Listed Games: ", games)

        this.gamesContainer.innerHTML = games.map(game => `
            <div class="game-item" data-game-id="${game.game_id}">
                <div class="game-info">
                    <h3>${game.game_id}</h3>
                    <p>Dimensions: ${game.config?.board_width || 7}×${game.config?.board_height || 6}</p>
                    <p>Players: ${game.players?.length || 0}/${game.config?.max_players || 2}</p>
                    <p>Status: ${this.getGameStatusText(game)}</p>
                </div>
                <div class="game-actions">
                    <a href="/${game.game_id}" class="btn">Continue Game</a>
                    <button class="btn btn-secondary" onclick="gamesManager.removeGame('${game.game_id}')">Remove</button>
                </div>
            </div>
        `).join('');
    }


    private async getAllGamesFromIndexedDB(): Promise<any[]> {
        return new Promise((resolve, reject) => {
            if (!this.storageTransport || !this.storageTransport['db']) {
                resolve([]);
                return;
            }

            const db = this.storageTransport['db'];
            const transaction = db.transaction(['gameStates'], 'readonly');
            const store = transaction.objectStore('gameStates');
            
            const request = store.getAll();
            request.onsuccess = () => {
                const results = request.result || [];
                const games = results.map(item => item.gameState);
                resolve(games);
            };
            request.onerror = () => reject(request.error);
        });
    }

    private getGameStatusText(game: any): string {
        if (!game.status) return 'Unknown';
        
        switch (game.status) {
            case 0: return 'Waiting for players';
            case 1: return 'In progress';
            case 2: return 'Finished';
            default: return 'Unknown';
        }
    }

    private async handleCreateGame(event: Event): Promise<void> {
        event.preventDefault();
        
        const formData = new FormData(event.target as HTMLFormElement);
        const gameId = formData.get('gameId') as string;
        const numPlayers = parseInt(formData.get('numPlayers') as string) || 2;
        const boardWidth = parseInt(formData.get('boardWidth') as string) || 7;
        const boardHeight = parseInt(formData.get('boardHeight') as string) || 6;

        if (!gameId) {
            alert('Please enter a game name');
            return;
        }

        // Validate game ID format
        if (!this.isValidGameId(gameId)) {
            alert('Game ID can only contain letters, numbers, and hyphens');
            return;
        }

        // Create the game using WASM client
        if (!this.connect4Client) {
            alert('Game engine not loaded yet, please wait and try again');
            return;
        }

        const response = await this.connect4Client.connect4Service.createGame({
            gameId: gameId,
            creatorName: 'Creator', // Placeholder, will be set when joining slot
            config: {
                boardWidth: boardWidth,
                boardHeight: boardHeight,
                connectLength: 4,
                maxPlayers: numPlayers,
                minPlayers: 2,
                allowMultipleWinners: false,
                moveTimeoutSeconds: 30
            }
        });

        if (response.success) {
            console.log('Game created successfully in WASM (saved to IndexedDB via callback)');
            
            // WASM has already awaited storage completion before returning success
            // window.location.href = `/${gameId}`;
        } else {
            alert(`Failed to create game: ${response.message}`);
        }
    }

    private getStoredGames(): StoredGame[] {
        const gamesData = localStorage.getItem('connect4Games');
        return gamesData ? JSON.parse(gamesData) : [];
    }

    private storeGame(game: StoredGame): void {
        const games = this.getStoredGames();
        const existingIndex = games.findIndex(g => g.gameId === game.gameId);
        
        if (existingIndex >= 0) {
            games[existingIndex] = game;
        } else {
            games.push(game);
        }
        
        localStorage.setItem('connect4Games', JSON.stringify(games));
        this.loadExistingGames(); // Refresh the display
    }

    public handleFormSubmit(): void {
        if (this.createForm) {
            const event = new Event('submit', { bubbles: true, cancelable: true });
            this.createForm.dispatchEvent(event);
        }
    }

    public async removeGame(gameId: string): Promise<void> {
        // Remove from IndexedDB
        if (this.storageTransport) {
            await this.removeGameFromIndexedDB(gameId);
            console.log(`Removed game ${gameId} from IndexedDB`);
        }
        
        // Refresh the display
        await this.loadExistingGames();
    }

    private async removeGameFromIndexedDB(gameId: string): Promise<void> {
        return new Promise((resolve, reject) => {
            if (!this.storageTransport || !this.storageTransport['db']) {
                resolve();
                return;
            }

            const db = this.storageTransport['db'];
            const transaction = db.transaction(['gameStates'], 'readwrite');
            const store = transaction.objectStore('gameStates');
            
            const request = store.delete(gameId);
            request.onsuccess = () => resolve();
            request.onerror = () => reject(request.error);
        });
    }

    private isValidGameId(gameId: string): boolean {
        if (!gameId || gameId.length === 0 || gameId.length > 50) {
            return false;
        }
        
        // Allow alphanumeric characters and hyphens
        return /^[a-zA-Z0-9-]+$/.test(gameId);
    }

    private calculateDefaultDimensions(numPlayers: number): { width: number, height: number } {
        // Base dimensions for 2 players: 7x6
        // Scale proportionally for more players
        switch (numPlayers) {
            case 2:
                return { width: 7, height: 6 };
            case 3:
                return { width: 8, height: 7 };
            case 4:
                return { width: 9, height: 8 };
            case 5:
                return { width: 10, height: 9 };
            case 6:
                return { width: 11, height: 10 };
            default:
                return { width: 7, height: 6 };
        }
    }

    private updateBoardDimensions(): void {
        if (!this.numPlayersSelect || !this.boardWidthInput || !this.boardHeightInput) return;

        const numPlayers = parseInt(this.numPlayersSelect.value);
        const defaultDimensions = this.calculateDefaultDimensions(numPlayers);
        
        this.boardWidthInput.value = defaultDimensions.width.toString();
        this.boardHeightInput.value = defaultDimensions.height.toString();
        
        this.updateDimensionsDisplay();
    }

    private updateDimensionsDisplay(): void {
        if (!this.boardWidthInput || !this.boardHeightInput || !this.dimensionsDisplay) return;

        const width = this.boardWidthInput.value || '7';
        const height = this.boardHeightInput.value || '6';
        
        this.dimensionsDisplay.textContent = `${width}×${height}`;
    }
}

// Initialize the games list manager
const gamesManager = new GamesListManager();

// Make it globally available for HTML onclick handlers
(window as any).gamesManager = gamesManager;

export default gamesManager;
