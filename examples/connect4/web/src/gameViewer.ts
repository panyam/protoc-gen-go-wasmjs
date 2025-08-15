// Entry point for the game viewer page (/{gameId})
// Handles individual Connect4 game interface and real-time multiplayer

import Connect4Client from '../gen/wasmts/multiplayer_connect4Client.client';
import { GameState, GameConfig, Player } from '../gen/wasmts/connect4/models';
import { StatefulTransport, TransportFactory, IndexedDBTransport } from './transport';

// Game interface state
interface GameUI {
    gameId: string;
    playerId: string;
    gameState: GameState | null;
    transport: StatefulTransport | null;
    connect4Client: Connect4Client | null;
}

class GameViewer {
    private ui: GameUI;
    private elements: { [key: string]: HTMLElement | null } = {};
    private selectedSlot: number = -1;
    private globalStorageTransport: IndexedDBTransport | null = null; // For storage callbacks
    private gameTransport: IndexedDBTransport | null = null; // For real-time patches

    constructor() {
        // Extract game ID and optional player index from URL path
        // Supports: /gameId and /gameId/players/0 (or 1, 2, etc.)
        const pathParts = window.location.pathname.split('/').filter(p => p);
        const gameId = pathParts[0] || '';
        
        // Check if this is a player-specific URL: /gameId/players/0
        let urlPlayerIndex = -1;
        let urlPlayerId = '';
        if (pathParts.length >= 3 && pathParts[1] === 'players') {
            const indexStr = pathParts[2];
            const index = parseInt(indexStr);
            if (!isNaN(index) && index >= 0) {
                urlPlayerIndex = index;
                console.log('🔗 Player-specific URL detected:', { gameId, playerIndex: urlPlayerIndex });
            }
        }
        
        this.ui = {
            gameId,
            playerId: urlPlayerId, // Will be set later when we load game state
            gameState: null,
            transport: null,
            connect4Client: null
        };
        
        // Store the URL player index for later use
        (this as any).urlPlayerIndex = urlPlayerIndex;

        this.init();
    }

    private async init(): Promise<void>{
        // Wait for DOM to be ready
        if (document.readyState === 'loading') {
            document.addEventListener('DOMContentLoaded', () => this.initializeUI());
        } else {
            await this.initializeUI();
        }
    }

    private async initializeUI(): Promise<void> {
        // Get UI elements
        this.elements = {
            playerSlots: document.getElementById('playerSlots'),
            slotsContainer: document.getElementById('slotsContainer'),
            joinSlotModal: document.getElementById('joinSlotModal'),
            joinSlotForm: document.getElementById('joinSlotForm'),
            gameInterface: document.getElementById('gameInterface'),
            gameControls: document.getElementById('gameControls'),
            errorState: document.getElementById('errorState'),
            gameBoard: document.getElementById('gameBoard'),
            gameStatus: document.getElementById('gameStatus'),
            currentPlayerName: document.getElementById('currentPlayerName'),
            currentPlayerColor: document.getElementById('currentPlayerColor'),
            turnNumber: document.getElementById('turnNumber'),
            playersList: document.getElementById('playersList'),
            gameLog: document.getElementById('gameLog'),
            currentGameId: document.getElementById('currentGameId'),
            gameUrl: document.getElementById('gameUrl')
        };

        // Set game ID and URL in UI
        if (this.elements.currentGameId) {
            this.elements.currentGameId.textContent = this.ui.gameId;
        }
        if (this.elements.gameUrl) {
            this.elements.gameUrl.textContent = window.location.href;
        }

        // Check if we have a valid game ID
        if (!this.ui.gameId || !this.isValidGameId(this.ui.gameId)) {
            this.showError('Invalid game ID');
            return;
        }

        // Set up form handlers
        if (this.elements.joinSlotForm) {
            this.elements.joinSlotForm.addEventListener('submit', (e) => this.handleJoinSlot(e));
        }

        // Initialize game components
        await this.initializeWasmClient();
        await this.initializeStatefulProxy();
        await this.loadGameState();
    }

    private async initializeWasmClient(): Promise<void> {
        console.log('Starting initializeWasmClient in gameViewer...');
        console.log('Loading WASM module...');
        this.ui.connect4Client = new Connect4Client();
        
        // Initialize global storage transport for storage callbacks (same as index.ts)
        console.log('Initializing global IndexedDB transport...');
        this.globalStorageTransport = new IndexedDBTransport('global');
        await this.globalStorageTransport.init();
        console.log('Global IndexedDB transport initialized');
        
        // Initialize game-specific transport for real-time patches
        console.log('Initializing game-specific transport...');
        this.gameTransport = new IndexedDBTransport(this.ui.gameId);
        await this.gameTransport.init();
        console.log('Game-specific transport initialized');
        
        // Load WASM first so that setWasmStorageCallbacks is available
        console.log('Loading WASM...');
        await this.ui.connect4Client.loadWasm('/static/wasm/multiplayer_connect4.wasm');
        await this.ui.connect4Client.waitUntilReady();
        console.log('WASM loaded and ready');
        
        // Set up storage callbacks after WASM is ready
        console.log('Setting up storage callbacks...');
        await this.setupStorageCallbacks();
        console.log('Storage callbacks setup complete');
        
        console.log('WASM module loaded successfully!');
    }

    private async setupStorageCallbacks(): Promise<void> {
        if (!this.globalStorageTransport || !(window as any).setWasmStorageCallbacks) {
            console.error('Global storage transport or WASM callbacks not available:', {
                globalStorageTransport: !!this.globalStorageTransport,
                setWasmStorageCallbacks: !!(window as any).setWasmStorageCallbacks
            });
            return;
        }

        // Simple, clean async callbacks using global transport - let failures bubble up to WASM
        const saveCallback = async (gameId: string, gameStateJson: string) => {
            const gameState = JSON.parse(gameStateJson);
            
            console.log('💾 SaveCallback triggered for game:', gameId);
            console.log('💾 Game state being saved players:', gameState.players);
            console.log('💾 Number of players being saved:', gameState.players?.length || 0);
            
            // Save to global DB for cross-page persistence
            await this.globalStorageTransport!.saveGameState(gameId, gameState);
            // Also save to game-specific DB for real-time sync (optional)
            await this.gameTransport!.saveGameState(gameId, gameState);
            console.log(`💾 Saved game state for ${gameId} to both global and game-specific IndexedDB`);
        };

        const loadCallback = async (gameId: string) => {
          console.log("Did loadCallback get called?: ", gameId)
          
          // Load from global DB for cross-page persistence
            const gameState = await this.globalStorageTransport!.loadGameState(gameId);
            if (gameState) {
                console.log(`Loaded game state for ${gameId} from global IndexedDB`);
                return JSON.stringify(gameState);
            } else {
                console.log(`No game state found for ${gameId}, checking all stored games...`);
                // Debug: list all games to see what's actually stored in global DB
                await this.globalStorageTransport!.debugListAllGames();
            }
            return null; // Game not found
        };

        const pollCallback = (gameId: string) => {
            // Set up polling for external changes using game-specific transport
            this.gameTransport!.onGameStateChanged((gameState: any) => {
                // Update UI when external changes are detected
                this.handleExternalGameStateChange(gameState);
                
                // Notify WASM of external change
                if ((window as any).wasmOnExternalStorageChange) {
                    (window as any).wasmOnExternalStorageChange(gameId, JSON.stringify(gameState));
                }
            });
        };

        // Configure WASM with these callbacks
        console.log('About to set WASM storage callbacks...');
        const result = (window as any).setWasmStorageCallbacks(
            saveCallback,
            loadCallback, 
            pollCallback
        );
        
        console.log('setWasmStorageCallbacks result:', result);
        if (result && result.success) {
            console.log('Storage callbacks configured successfully in gameViewer');
        } else {
            console.error('Failed to configure storage callbacks in gameViewer:', result);
        }
    }

    private handleExternalGameStateChange(gameState: any): void {
        if (!this.ui.gameState || JSON.stringify(this.ui.gameState) !== JSON.stringify(gameState)) {
            this.ui.gameState = GameState.from(gameState);
            this.updateGameDisplay();
            console.log('Updated UI from external game state change');
        }
    }

    private async initializeStatefulProxy(): Promise<void> {
        // Use the already-initialized game-specific transport for real-time patches
        this.ui.transport = this.gameTransport;
        
        // Set up state change listener for real-time patches
        if (this.ui.transport) {
            this.ui.transport.subscribe((patches: any[]) => {
                console.log('Received state patches:', patches);
                this.applyPatches(patches);
            });
        }

        console.log('Stateful proxy initialized using game-specific transport');
    }

    private async loadGameState(): Promise<void> {
        // Use the new async callback-based API to avoid deadlocks
        await this.ui.connect4Client!.connect4Service.getGame({
            gameId: this.ui.gameId
        }, (response, error) => {
            if (error) {
                console.error('Failed to load game state:', error);
                this.showError('Game not found - please check the game ID or create a new game');
                return;
            }
            
            if (response) {
                this.ui.gameState = GameState.from(JSON.parse(response));
                console.log('Loaded game state via callback:', this.ui.gameState);
                
                        // Debug: Log the loaded game state
        console.log('🔍 Loaded game state players:', this.ui.gameState.players);
        console.log('🔍 Number of players loaded:', this.ui.gameState.players?.length || 0);
        
        // Show player selection if players exist
        this.checkPlayerSelection();
                
                // Always show the game interface - users can participate as players or viewers
                console.log('Game state players:', this.ui.gameState.players);
                console.log('Game state config:', this.ui.gameState.config);
                console.log('Game state board:', this.ui.gameState.board);
                
                this.showGameInterface();
                this.updateGameDisplay();
            } else {
                this.showError('Game not found - please check the game ID or create a new game');
            }
        });
    }

    private checkPlayerSelection(): void {
        if (!this.ui.gameState) return;
        
        console.log('🔍 checkPlayerSelection called:', {
            urlPlayerIndex: (this as any).urlPlayerIndex,
            currentPlayerId: this.ui.playerId,
            gameStatePlayers: this.ui.gameState.players?.map((p, i) => ({ index: i, id: p.id, name: p.name })),
            url: window.location.pathname
        });
        
        // If URL contains player index, validate and auto-select
        const urlPlayerIndex = (this as any).urlPlayerIndex;
        if (urlPlayerIndex >= 0) {
            const player = this.ui.gameState.players?.[urlPlayerIndex];
            if (player) {
                console.log('🔗 Auto-selecting player from URL index:', { playerIndex: urlPlayerIndex, playerId: player.id, playerName: player.name });
                this.ui.playerId = player.id;
                this.savePlayerIdentity(player.id, player.name);
                return; // Skip showing modal
            } else {
                console.warn('🚨 Player index from URL not found in game:', {
                    urlPlayerIndex,
                    availablePlayers: this.ui.gameState.players?.length || 0,
                    players: this.ui.gameState.players
                });
                // Player doesn't exist at that index, fall through to normal selection
            }
        }
        
        // If there are existing players, show player selection
        if (this.ui.gameState.players && this.ui.gameState.players.length > 0) {
            console.log('🎭 Showing player selection modal - no URL player index or player not found');
            this.showPlayerSelectionModal();
        } else {
            // No players yet - user needs to join first
            this.ui.playerId = '';
            console.log('🎮 No players in game yet - user needs to join');
        }
    }

    private showPlayerSelectionModal(): void {
        if (!this.ui.gameState?.players) return;
        
        // Create player selection modal
        const modal = document.createElement('div');
        modal.className = 'modal';
        modal.id = 'playerSelectionModal';
        
        const playerOptions = this.ui.gameState.players.map((player, index) => {
            const isCurrentPlayer = player.id === this.ui.gameState?.currentPlayerId;
            const playerColors = this.getPlayerColors(this.ui.gameState?.config?.maxPlayers || 2);
            const playerColor = playerColors[index];
            
            return `
                <div class="player-option" data-player-id="${player.id}" onclick="gameViewer.selectPlayer('${player.id}')">
                    <div class="player-info">
                        <span class="player-color" style="color: ${playerColor}">●</span>
                        <span class="player-name">${player.name}</span>
                        ${isCurrentPlayer ? '<span class="current-indicator">(Current Turn)</span>' : ''}
                    </div>
                </div>
            `;
        }).join('');
        
        modal.innerHTML = `
            <div class="modal-content">
                <h3>Select Player</h3>
                <p>Choose which player you want to play as in this tab:</p>
                <div class="player-selection">
                    ${playerOptions}
                    <div class="player-option spectator" onclick="gameViewer.selectPlayer('')">
                        <div class="player-info">
                            <span class="player-color">👁️</span>
                            <span class="player-name">Spectate Only</span>
                        </div>
                    </div>
                </div>
                <div class="modal-note">
                    <small>You can play as any existing player or just watch the game.</small>
                </div>
            </div>
        `;
        
        document.body.appendChild(modal);
    }

    public selectPlayer(playerId: string): void {
        this.ui.playerId = playerId;
        
        // Close the selection modal
        const modal = document.getElementById('playerSelectionModal');
        if (modal) {
            modal.remove();
        }
        
        if (playerId) {
            const player = this.ui.gameState?.players?.find(p => p.id === playerId);
            console.log('🎮 Selected player:', player?.name, 'ID:', playerId);
            this.addLogEntry(`You are now playing as ${player?.name}`);
        } else {
            console.log('🎮 Selected spectator mode');
            this.addLogEntry('You are now spectating');
        }
        
        // Update the display to reflect selection
        this.updateGameDisplay();
    }

    private savePlayerIdentity(playerId: string, playerName: string): void {
        const identity = {
            playerId: playerId,
            playerName: playerName,
            joinedAt: Date.now()
        };
        
        localStorage.setItem(`connect4_player_${this.ui.gameId}`, JSON.stringify(identity));
        console.log('🎮 Saved player identity:', playerName, 'ID:', playerId);
    }


    private async handleJoinGame(event: Event): Promise<void> {
        event.preventDefault();
        
        const formData = new FormData(event.target as HTMLFormElement);
        const playerName = formData.get('playerName') as string;

        if (!playerName.trim()) {
            alert('Please enter your name');
            return;
        }

        this.ui.playerId = `player_${Date.now()}`;
        
        // Try to join existing game first
        const joinResponse = await this.joinGame(playerName);
        
        if (joinResponse.success) {
            this.ui.gameState = GameState.from(joinResponse.data);
            this.showGameInterface();
            this.updateGameDisplay();
            this.addLogEntry(`${playerName} joined the game`);
        } else {
            // If join fails, try to create a new game
            console.log('Join failed, attempting to create new game...');
            const createResponse = await this.createGame(playerName);
            
            if (createResponse.success) {
                this.ui.gameState = GameState.from(createResponse.gameState);
                this.showGameInterface();
                this.updateGameDisplay();
                this.addLogEntry(`Game created by ${playerName}`);
            } else {
                throw new Error(createResponse.errorMessage || 'Failed to create game');
            }
        }
    }

    private async joinGame(playerName: string): Promise<any> {
        if (!this.ui.connect4Client) {
            throw new Error('WASM client not initialized');
        }

        return new Promise((resolve, reject) => {
            this.ui.connect4Client!.connect4Service.joinGame({
                gameId: this.ui.gameId,
                playerId: this.ui.playerId,
                playerName: playerName
            }, (response, error) => {
                if (error) {
                    reject(new Error(error));
                    return;
                }
                
                if (response) {
                    const parsedResponse = JSON.parse(response);
                    resolve(parsedResponse);
                } else {
                    reject(new Error('No response received'));
                }
            });
        });
    }

    private async createGame(playerName: string): Promise<any> {
        if (!this.ui.connect4Client) {
            throw new Error('WASM client not initialized');
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

        return new Promise((resolve, reject) => {
            this.ui.connect4Client!.connect4Service.createGame({
                gameId: this.ui.gameId,
                playerId: this.ui.playerId,
                playerName: playerName,
                config: gameConfig
            }, (response, error) => {
                if (error) {
                    reject(new Error(error));
                    return;
                }
                
                if (response) {
                    const parsedResponse = JSON.parse(response);
                    resolve(parsedResponse);
                } else {
                    reject(new Error('Failed to create game'));
                }
            });
        });
    }

    public async dropPiece(column: number): Promise<void> {
        if (!this.ui.connect4Client || !this.ui.gameState) {
            console.error('Game not properly initialized');
            return;
        }

        // Check if user has selected a player
        if (!this.ui.playerId) {
            alert('You need to select a player first! Choose which player you want to play as.');
            return;
        }

        // Check if it's the selected player's turn
        if (this.ui.gameState.currentPlayerId !== this.ui.playerId) {
            const currentPlayer = this.ui.gameState.players?.find(p => p.id === this.ui.gameState?.currentPlayerId);
            const selectedPlayer = this.ui.gameState.players?.find(p => p.id === this.ui.playerId);
            alert(`It's ${currentPlayer?.name}'s turn! You're playing as ${selectedPlayer?.name}.`);
            return;
        }

        console.log('🎮 Dropping piece:', {
            gameId: this.ui.gameId,
            playerId: this.ui.playerId,
            column: column
        });

        await this.ui.connect4Client.connect4Service.dropPiece({
            gameId: this.ui.gameId,
            playerId: this.ui.playerId,
            column: column
        }, (response, error) => {
            console.log('🎮 Drop piece callback - Response:', response);
            console.log('🎮 Drop piece callback - Error:', error);

            if (error) {
                console.error('Failed to drop piece:', error);
                alert(`Move failed: ${error}`);
                return;
            }

            if (response) {
                const parsedResponse = JSON.parse(response);
                console.log('🎮 Parsed drop piece response:', parsedResponse);

                if (parsedResponse.success) {
                    // Update game state from response
                    this.ui.gameState = GameState.from(parsedResponse.gameState || parsedResponse.data);
                    this.updateGameDisplay();
                    
                    // Send state update through stateful transport
                    if (this.ui.transport) {
                        this.ui.transport.sendPatches([{
                            operation: 'update',
                            path: '',
                            value: this.ui.gameState,
                            timestamp: Date.now(),
                            source: this.ui.playerId
                        }]);
                    }

                    // Add move to log
                    const player = this.ui.gameState.players?.find(p => p.id === this.ui.playerId);
                    this.addLogEntry(`${player?.name || 'You'} dropped piece in column ${column + 1}`);
                } else {
                    console.error('Failed to drop piece:', parsedResponse.errorMessage);
                    alert(`Move failed: ${parsedResponse.errorMessage || 'Unknown error'}`);
                }
            } else {
                console.error('🎮 No response received from dropPiece');
                alert('Move failed: No response received');
            }
        });
    }

    private applyPatches(patches: any[]): void {
        console.log('🔄 Applying patches:', {
            patchCount: patches.length,
            currentPlayerId: this.ui.playerId,
            urlPlayerIndex: (this as any).urlPlayerIndex,
            beforeState: this.ui.gameState?.players?.map((p, i) => ({ index: i, id: p.id, name: p.name }))
        });
        
        for (const patch of patches) {
            if (patch.operation === 'update' && patch.value) {
                const newState = GameState.from(patch.value);
                if (newState) {
                    // Validate patch: reject if it would make the game state worse
                    const currentPlayerCount = this.ui.gameState?.players?.length || 0;
                    const newPlayerCount = newState.players?.length || 0;
                    
                    // Don't apply patches that remove all players (likely corrupted)
                    if (currentPlayerCount > 0 && newPlayerCount === 0) {
                        console.warn('🚨 Rejecting patch with empty players array:', {
                            currentPlayerCount,
                            newPlayerCount,
                            patchValue: patch.value
                        });
                        continue; // Skip this patch
                    }
                    
                    const previousPlayerId = this.ui.playerId;
                    this.ui.gameState = newState;
                    
                    console.log('🔄 Game state updated from patch:', {
                        previousPlayerId,
                        currentPlayerCount,
                        newPlayerCount,
                        newPlayers: newState.players?.map((p, i) => ({ index: i, id: p.id, name: p.name })),
                        willMaintainPlayerSelection: !!previousPlayerId
                    });
                    
                    // If we had a player selected from URL, maintain that selection
                    if (previousPlayerId && (this as any).urlPlayerIndex >= 0) {
                        const urlPlayerIndex = (this as any).urlPlayerIndex;
                        const urlPlayer = newState.players?.[urlPlayerIndex];
                        if (urlPlayer && urlPlayer.id !== previousPlayerId) {
                            console.log('🔗 Re-applying URL player selection after patch:', { 
                                urlPlayerIndex, 
                                newPlayerId: urlPlayer.id, 
                                previousPlayerId 
                            });
                            this.ui.playerId = urlPlayer.id;
                        }
                    }
                    
                    this.updateGameDisplay();
                    this.addLogEntry('Game state updated from another player');
                }
            }
        }
    }

    private showPlayerSlots(): void {
        if (this.elements.playerSlots) {
            this.elements.playerSlots.classList.remove('hidden');
        }
        if (this.elements.gameInterface) {
            this.elements.gameInterface.classList.add('hidden');
        }
        if (this.elements.gameControls) {
            this.elements.gameControls.classList.add('hidden');
        }
        if (this.elements.errorState) {
            this.elements.errorState.classList.add('hidden');
        }
        if (this.elements.joinSlotModal) {
            this.elements.joinSlotModal.classList.add('hidden');
        }
        
        this.generatePlayerSlots();
    }

    private showGameInterface(): void {
        console.log('showGameInterface called');
        console.log('Elements:', {
            playerSlots: !!this.elements.playerSlots,
            gameInterface: !!this.elements.gameInterface,
            gameControls: !!this.elements.gameControls,
            gameBoard: !!this.elements.gameBoard,
            slotsContainer: !!this.elements.slotsContainer
        });
        
        // Show both left panel and center panel
        if (this.elements.playerSlots) {
            this.elements.playerSlots.classList.remove('hidden');
            console.log('Showed player slots');
        }
        if (this.elements.gameInterface) {
            this.elements.gameInterface.classList.remove('hidden');
            console.log('Showed game interface');
        }
        if (this.elements.gameControls) {
            this.elements.gameControls.classList.remove('hidden');
            console.log('Showed game controls');
        }
        if (this.elements.errorState) {
            this.elements.errorState.classList.add('hidden');
        }
        if (this.elements.joinSlotModal) {
            this.elements.joinSlotModal.classList.add('hidden');
        }
        
        console.log('Calling generatePlayerSlots...');
        this.generatePlayerSlots();
        console.log('Calling initializeGameBoard...');
        this.initializeGameBoard();
        console.log('showGameInterface complete');
    }

    private showError(message: string): void {
        if (this.elements.playerSlots) {
            this.elements.playerSlots.classList.add('hidden');
        }
        if (this.elements.gameInterface) {
            this.elements.gameInterface.classList.add('hidden');
        }
        if (this.elements.gameControls) {
            this.elements.gameControls.classList.add('hidden');
        }
        if (this.elements.errorState) {
            this.elements.errorState.classList.remove('hidden');
            const errorMessage = this.elements.errorState.querySelector('p');
            if (errorMessage) {
                errorMessage.textContent = message;
            }
        }
        if (this.elements.joinSlotModal) {
            this.elements.joinSlotModal.classList.add('hidden');
        }
    }

    private initializeGameBoard(): void {
        if (!this.elements.gameBoard || !this.ui.gameState?.board) return;

        const rows = this.ui.gameState.config?.boardHeight || 6;
        const cols = this.ui.gameState.config?.boardWidth || 7;
        const playerColors = this.getPlayerColors(this.ui.gameState.config?.maxPlayers || 2);

        let boardHTML = '';
        for (let row = 0; row < rows; row++) {
            boardHTML += '<div class="board-row">';
            for (let col = 0; col < cols; col++) {
                const cellValue = this.ui.gameState.board.rows[row]?.cells[col] || '';
                const player = this.ui.gameState.players.find(p => p.id === cellValue);
                const playerIndex = player ? this.ui.gameState.players.indexOf(player) : -1;
                const pieceColor = playerIndex >= 0 ? playerColors[playerIndex] : '';
                
                boardHTML += `
                    <div class="board-cell ${cellValue ? 'occupied' : 'empty'}" 
                         data-row="${row}" 
                         data-col="${col}"
                         onclick="gameViewer.dropPiece(${col})"
                         style="cursor: pointer;">
                        ${cellValue ? `<span class="piece" style="color: ${pieceColor}">●</span>` : ''}
                    </div>
                `;
            }
            boardHTML += '</div>';
        }

        this.elements.gameBoard.innerHTML = boardHTML;
    }

    private updateBoardOnly(): void {
        if (!this.elements.gameBoard || !this.ui.gameState?.board) return;

        const cells = this.elements.gameBoard.querySelectorAll('.board-cell');
        const playerColors = this.getPlayerColors(this.ui.gameState.config?.maxPlayers || 2);
        
        cells.forEach((cell: Element) => {
            const htmlCell = cell as HTMLElement;
            const row = parseInt(htmlCell.dataset.row || '0');
            const col = parseInt(htmlCell.dataset.col || '0');
            const cellValue = this.ui.gameState!.board!.rows[row]?.cells[col] || '';
            const player = this.ui.gameState!.players.find(p => p.id === cellValue);
            const playerIndex = player ? this.ui.gameState!.players.indexOf(player) : -1;
            const pieceColor = playerIndex >= 0 ? playerColors[playerIndex] : '';
            
            if (cellValue) {
                htmlCell.className = 'board-cell occupied';
                htmlCell.innerHTML = `<span class="piece" style="color: ${pieceColor}">●</span>`;
            } else {
                htmlCell.className = 'board-cell empty';
                htmlCell.innerHTML = '';
            }
        });
    }

    private updateGameDisplay(): void {
        if (!this.ui.gameState) return;

        console.log('🎨 updateGameDisplay called:', {
            currentPlayerId: this.ui.playerId,
            urlPlayerIndex: (this as any).urlPlayerIndex,
            gameStatePlayers: this.ui.gameState.players?.map((p, i) => ({ index: i, id: p.id, name: p.name }))
        });

        // Update turn information
        if (this.elements.turnNumber) {
            this.elements.turnNumber.textContent = this.ui.gameState.turnNumber.toString();
        }

        // Update current player
        const currentPlayer = this.ui.gameState.players.find(p => p.id === this.ui.gameState!.currentPlayerId);
        const selectedPlayer = this.ui.gameState.players.find(p => p.id === this.ui.playerId);
    
        if (this.elements.currentPlayerName && currentPlayer) {
            let displayText = currentPlayer.name;
            if (this.ui.playerId) {
                if (currentPlayer.id === this.ui.playerId) {
                    displayText += ' (Your Turn!)';
                } else {
                    displayText += ` | You: ${selectedPlayer?.name || 'Spectating'}`;
                }
            } else {
                displayText += ' | You: Spectating';
            }
            
            this.elements.currentPlayerName.textContent = displayText;
            
            // Update current player color
            if (this.elements.currentPlayerColor) {
                const playerIndex = this.ui.gameState.players.indexOf(currentPlayer);
                const playerColors = this.getPlayerColors(this.ui.gameState.config?.maxPlayers || 2);
                const playerColor = playerColors[playerIndex] || '#e74c3c';
                this.elements.currentPlayerColor.style.backgroundColor = playerColor;
                this.elements.currentPlayerColor.innerHTML = '●';
                this.elements.currentPlayerColor.style.color = playerColor;
            }
        }

        // Update game status
        if (this.elements.gameStatus) {
            let statusText = 'In Progress';
            if (this.ui.gameState.status === 2) { // Assuming 2 is game over
                statusText = this.ui.gameState.winners.length > 0 ? 
                    `Game Over - Winner: ${this.ui.gameState.winners.join(', ')}` : 
                    'Game Over - Draw';
            } else if (this.ui.gameState.players.length < 2) {
                statusText = 'Waiting for players...';
            }
            this.elements.gameStatus.textContent = statusText;
        }

        // Update players list
        if (this.elements.playersList) {
            const playerColors = this.getPlayerColors(this.ui.gameState.config?.maxPlayers || 2);
            this.elements.playersList.innerHTML = this.ui.gameState.players.map((player, index) => {
                const playerColor = playerColors[index] || '#e74c3c';
                return `
                    <div class="player-item ${player.id === this.ui.gameState!.currentPlayerId ? 'current' : ''}">
                        <span class="player-name">${player.name}</span>
                        <span class="player-color" style="color: ${playerColor}">●</span>
                    </div>
                `;
            }).join('');
        }

        // Update player slots and board separately to avoid cascading updates
        this.generatePlayerSlots();
        this.updateBoardOnly();
    }


    private addLogEntry(message: string): void {
        if (!this.elements.gameLog) return;
        
        const timestamp = new Date().toLocaleTimeString();
        const logEntry = document.createElement('div');
        logEntry.className = 'log-entry';
        logEntry.innerHTML = `<span class="timestamp">[${timestamp}]</span> ${message}`;
        
        this.elements.gameLog.appendChild(logEntry);
        this.elements.gameLog.scrollTop = this.elements.gameLog.scrollHeight;
    }

    private isValidGameId(gameId: string): boolean {
        if (!gameId || gameId.length === 0 || gameId.length > 50) {
            return false;
        }
        return /^[a-zA-Z0-9-]+$/.test(gameId);
    }

    private getPlayerColors(maxPlayers: number): string[] {
        const defaultColors = [
            '#e74c3c', // Red
            '#3498db', // Blue  
            '#2ecc71', // Green
            '#f39c12', // Orange
            '#9b59b6', // Purple
            '#1abc9c'  // Teal
        ];
        
        return defaultColors.slice(0, maxPlayers);
    }

    private generatePlayerSlots(): void {
        console.log('generatePlayerSlots called');
        console.log('slotsContainer exists:', !!this.elements.slotsContainer);
        console.log('gameState exists:', !!this.ui.gameState);
        
        if (!this.elements.slotsContainer || !this.ui.gameState) {
            console.log('Early return from generatePlayerSlots - missing elements');
            return;
        }

        const maxPlayers = this.ui.gameState.config?.maxPlayers || 2;
        const currentPlayers = this.ui.gameState.players || [];
        const playerColors = this.getPlayerColors(maxPlayers);
        
        console.log('Player slot generation data:', {
            maxPlayers,
            currentPlayersCount: currentPlayers.length,
            currentPlayers,
            playerColors
        });

        let slotsHTML = '';
        for (let slotIndex = 0; slotIndex < maxPlayers; slotIndex++) {
            const player = currentPlayers.find((p, index) => index === slotIndex);
            const slotColor = playerColors[slotIndex];
            
            if (player) {
                // Check different states
                const isSelectedPlayer = player.id === this.ui.playerId;
                const isCurrentPlayer = player.id === this.ui.gameState?.currentPlayerId;
                
                let extraClasses = '';
                let indicators = [];
                let statusText = 'Joined';
                
                if (isSelectedPlayer) {
                    extraClasses += ' selected';
                    indicators.push('Playing as');
                    statusText = 'You';
                }
                
                if (isCurrentPlayer) {
                    extraClasses += ' current-turn';
                    indicators.push('Current Turn');
                    if (!isSelectedPlayer) {
                        statusText = 'Turn Active';
                    }
                }
                
                const indicatorText = indicators.length > 0 ? ` (${indicators.join(', ')})` : '';
                
                // Occupied slot
                const isGeneralGamePage = !window.location.pathname.includes('/players/');
                
                slotsHTML += `
                    <div class="player-slot occupied${extraClasses}" data-slot="${slotIndex}" style="border-left: 4px solid ${slotColor}">
                        <div class="slot-header">Player ${slotIndex + 1}${indicatorText}</div>
                        <div class="player-info">
                            <span class="player-name">${player.name}</span>
                            <span class="player-color" style="color: ${slotColor}">●</span>
                        </div>
                        <div class="slot-status">${statusText}</div>
                        ${!isSelectedPlayer ? `<button class="switch-btn" onclick="gameViewer.selectPlayer('${player.id}')">Play as ${player.name}</button>` : ''}
                        ${isGeneralGamePage ? `<button class="direct-link-btn" onclick="gameViewer.goToPlayerIndex(${slotIndex})" title="Go to player-specific page">🔗 Direct Link</button>` : ''}
                    </div>
                `;
            } else {
                // Empty slot
                slotsHTML += `
                    <div class="player-slot empty" data-slot="${slotIndex}" onclick="gameViewer.openJoinModal(${slotIndex})" style="border-left: 4px solid ${slotColor}">
                        <div class="slot-header">Player ${slotIndex + 1}</div>
                        <div class="empty-message">
                            <span class="player-color" style="color: ${slotColor}">●</span>
                            Click to join
                        </div>
                        <div class="slot-status">Available</div>
                    </div>
                `;
            }
        }

        console.log('Generated slots HTML:', slotsHTML);
        this.elements.slotsContainer.innerHTML = slotsHTML;
        console.log('Set slotsContainer innerHTML complete');
    }

    public openJoinModal(slotIndex: number): void {
        this.selectedSlot = slotIndex;
        
        if (this.elements.joinSlotModal) {
            this.elements.joinSlotModal.classList.remove('hidden');
            
            // Update slot number in modal
            const slotNumberElement = this.elements.joinSlotModal.querySelector('#slotNumber');
            if (slotNumberElement) {
                slotNumberElement.textContent = (slotIndex + 1).toString();
            }
            
            // Focus on the name input
            const nameInput = this.elements.joinSlotModal.querySelector('#playerName') as HTMLInputElement;
            if (nameInput) {
                nameInput.focus();
            }
        }
    }

    private async handleJoinSlot(event: Event): Promise<void> {
        event.preventDefault();
        
        const formData = new FormData(event.target as HTMLFormElement);
        const playerName = formData.get('playerName') as string;

        if (!playerName.trim()) {
            alert('Please enter your name');
            return;
        }

        if (this.selectedSlot < 0) {
            alert('No slot selected');
            return;
        }

        // Debug: Log game state before joining
        console.log('🎮 Attempting to join slot:', {
            selectedSlot: this.selectedSlot,
            playerName: playerName,
            gameId: this.ui.gameId,
            currentGameState: this.ui.gameState,
            currentPlayers: this.ui.gameState?.players,
            maxPlayers: this.ui.gameState?.config?.maxPlayers,
            playerCount: this.ui.gameState?.players?.length || 0
        });

        this.ui.playerId = `player_${Date.now()}`;
    
        await this.ui.connect4Client!.connect4Service.joinGame({
            gameId: this.ui.gameId,
            playerId: this.ui.playerId,
            playerName: playerName
        }, (response, error) => {
            console.log('🧪 JoinGame callback - Raw response:', response);
            console.log('🧪 JoinGame callback - Error:', error);
            
            if (error) {
                console.error('Error joining slot:', error);
                alert(`Failed to join slot: ${error}`);
                return;
            }
            
            if (response) {
                console.log('🧪 Join response received:', response);
                const parsedResponse = JSON.parse(response);
                console.log('🧪 Parsed join response:', parsedResponse);
                
                if (!parsedResponse.success) {
                    console.error('🚨 Join failed - Server response:', {
                        success: parsedResponse.success,
                        errorMessage: parsedResponse.errorMessage,
                        fullResponse: parsedResponse
                    });
                }
                
                if (parsedResponse.success) {
                    this.ui.playerId = parsedResponse.playerId;
                    this.ui.gameState = GameState.from(parsedResponse.gameState);
                    
                    console.log('🎉 Player joined successfully!');
                    console.log('🎉 New game state players:', this.ui.gameState.players);
                    console.log('🎉 Number of players after join:', this.ui.gameState.players?.length || 0);
                    
                    // Save player identity for persistence across refreshes
                    this.savePlayerIdentity(parsedResponse.playerId, playerName);
                    
                    this.showGameInterface();
                    this.updateGameDisplay();
                    this.addLogEntry(`${playerName} joined as Player ${this.selectedSlot + 1}`);
                    
                    // Close the modal
                    if (this.elements.joinSlotModal) {
                        this.elements.joinSlotModal.classList.add('hidden');
                    }
                    
                    // Redirect to player-specific URL if we're on the general game page
                    const currentPath = window.location.pathname;
                    const isGeneralGamePage = !currentPath.includes('/players/');
                    if (isGeneralGamePage) {
                        // Find the player index in the game state
                        const playerIndex = this.ui.gameState.players?.findIndex(p => p.id === parsedResponse.playerId) ?? -1;
                        if (playerIndex >= 0) {
                            const playerSpecificUrl = `/${this.ui.gameId}/players/${playerIndex}`;
                            console.log('🔗 Redirecting to player-specific URL:', playerSpecificUrl);
                            window.history.pushState({}, '', playerSpecificUrl);
                        }
                    }
                } else {
                    console.error('🧪 Join failed with error:', parsedResponse.errorMessage);
                    alert(`Failed to join slot: ${parsedResponse.errorMessage || 'Unknown error'}`);
                }
            } else {
                console.error('🧪 No response received from joinGame');
                alert('Failed to join slot: No response received');
            }
        });
    }

    // Public methods for HTML onclick handlers
    public resetGame(): void {
        if (confirm('Are you sure you want to start a new game?')) {
            window.location.reload();
        }
    }
    
    public goToPlayerIndex(playerIndex: number): void {
        const playerSpecificUrl = `/${this.ui.gameId}/players/${playerIndex}`;
        console.log('🔗 Navigating to player-specific URL:', playerSpecificUrl);
        window.location.href = playerSpecificUrl;
    }

    public leaveGame(): void {
        if (confirm('Are you sure you want to leave this game?')) {
            window.location.href = '/';
        }
    }

    public closeJoinModal(): void {
        if (this.elements.joinSlotModal) {
            this.elements.joinSlotModal.classList.add('hidden');
        }
        this.selectedSlot = -1;
        
        // Clear the form
        const form = this.elements.joinSlotForm as HTMLFormElement;
        if (form) {
            form.reset();
        }
    }
}

// Initialize the game viewer
const gameViewer = new GameViewer();

// Make it globally available for HTML onclick handlers
(window as any).gameViewer = gameViewer;
(window as any).resetGame = () => gameViewer.resetGame();
(window as any).leaveGame = () => gameViewer.leaveGame();

export default gameViewer;
