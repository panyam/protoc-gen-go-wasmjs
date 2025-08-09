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
        this.loadExistingGames();
        this.initializeWasmClient();
    }

    private async initializeWasmClient(): Promise<void> {
        try {
            console.log('Initializing WASM client for games list...');
            this.connect4Client = new Connect4Client();
            
            // Initialize storage transport for global game state management
            this.storageTransport = new IndexedDBTransport('global');
            await this.storageTransport.init();
            
            // Set up storage callbacks before loading WASM
            await this.setupStorageCallbacks();
            
            await this.connect4Client.loadWasm('/static/wasm/multiplayer_connect4.wasm');
            await this.connect4Client.waitUntilReady();
            console.log('WASM client ready for game operations');
        } catch (error) {
            console.error('Failed to initialize WASM client:', error);
        }
    }

    private async setupStorageCallbacks(): Promise<void> {
        if (!this.storageTransport || !(window as any).setWasmStorageCallbacks) {
            console.warn('Storage transport or WASM callbacks not available');
            return;
        }

        // Create callback functions for WASM to use
        const saveCallback = async (gameId: string, gameStateJson: string) => {
            try {
                const gameState = JSON.parse(gameStateJson);
                await this.storageTransport!.saveGameState(gameId, gameState);
                console.log(`Saved game state for ${gameId} to IndexedDB`);
            } catch (error) {
                console.error('Failed to save game state:', error);
            }
        };

        const loadCallback = async (gameId: string): Promise<string | null> => {
            try {
                const gameState = await this.storageTransport!.loadGameState(gameId);
                if (gameState) {
                    console.log(`Loaded game state for ${gameId} from IndexedDB`);
                    return JSON.stringify(gameState);
                }
                return null;
            } catch (error) {
                console.error('Failed to load game state:', error);
                return null;
            }
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

    private loadExistingGames(): void {
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

        this.gamesContainer.innerHTML = games.map(game => `
            <div class="game-item" data-game-id="${game.gameId}">
                <div class="game-info">
                    <h3>${game.gameId}</h3>
                    <p>Player: ${game.playerName}</p>
                    <p>Status: ${game.gameStatus || 'Unknown'}</p>
                    <small>Last played: ${new Date(game.lastPlayed).toLocaleString()}</small>
                </div>
                <div class="game-actions">
                    <a href="/${game.gameId}" class="btn">Continue Game</a>
                    <button class="btn btn-secondary" onclick="gamesManager.removeGame('${game.gameId}')">Remove</button>
                </div>
            </div>
        `).join('');
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

        try {
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
                // Store the game locally
                this.storeGame({
                    gameId,
                    playerName: '', // Will be set when joining a slot
                    lastPlayed: Date.now(),
                    gameStatus: 'Created'
                });

                // Navigate to the game
                window.location.href = `/${gameId}`;
            } else {
                alert(`Failed to create game: ${response.message}`);
            }
        } catch (error) {
            console.error('Error creating game:', error);
            alert('Failed to create game. Please try again.');
        }
    }

    private getStoredGames(): StoredGame[] {
        try {
            const gamesData = localStorage.getItem('connect4Games');
            return gamesData ? JSON.parse(gamesData) : [];
        } catch (error) {
            console.error('Error loading stored games:', error);
            return [];
        }
    }

    private storeGame(game: StoredGame): void {
        try {
            const games = this.getStoredGames();
            const existingIndex = games.findIndex(g => g.gameId === game.gameId);
            
            if (existingIndex >= 0) {
                games[existingIndex] = game;
            } else {
                games.push(game);
            }
            
            localStorage.setItem('connect4Games', JSON.stringify(games));
            this.loadExistingGames(); // Refresh the display
        } catch (error) {
            console.error('Error storing game:', error);
        }
    }

    public removeGame(gameId: string): void {
        try {
            const games = this.getStoredGames();
            const filteredGames = games.filter(g => g.gameId !== gameId);
            localStorage.setItem('connect4Games', JSON.stringify(filteredGames));
            this.loadExistingGames(); // Refresh the display
        } catch (error) {
            console.error('Error removing game:', error);
        }
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
