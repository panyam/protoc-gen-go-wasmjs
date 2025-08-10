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
        // Extract game ID from URL path
        const pathParts = window.location.pathname.split('/').filter(p => p);
        const gameId = pathParts[0] || '';
        
        this.ui = {
            gameId,
            playerId: '',
            gameState: null,
            transport: null,
            connect4Client: null
        };

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
            // Save to global DB for cross-page persistence
            await this.globalStorageTransport!.saveGameState(gameId, gameState);
            // Also save to game-specific DB for real-time sync (optional)
            await this.gameTransport!.saveGameState(gameId, gameState);
            console.log(`Saved game state for ${gameId} to both global and game-specific IndexedDB`);
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
        this.ui.gameState = GameState.from(gameState);
        this.updateGameDisplay();
        console.log('Updated UI from external game state change');
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
        // WASM will now handle loading from storage automatically via callbacks
        // Just try to get the game state - WASM will load from IndexedDB if not in memory
        // Ensure WASM is ready before making calls
        await this.ui.connect4Client!.waitUntilReady();
        
        const response = await this.ui.connect4Client!.connect4Service.getGame({
            gameId: this.ui.gameId
        });

        if (response.success) {
            this.ui.gameState = GameState.from(response.data);
            console.log('Loaded game state (WASM handled storage):', this.ui.gameState);
            
            // Always show the game interface - users can participate as players or viewers
            console.log('Game state players:', this.ui.gameState.players);
            console.log('Game state config:', this.ui.gameState.config);
            console.log('Game state board:', this.ui.gameState.board);
            
            this.showGameInterface();
            this.updateGameDisplay();
            return;
        }

        // If WASM couldn't find/load the game, show error
        this.showError('Game not found - please check the game ID or create a new game');
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
        
        try {
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
                    this.ui.gameState = GameState.from(createResponse.data);
                    this.showGameInterface();
                    this.updateGameDisplay();
                    this.addLogEntry(`Game created by ${playerName}`);
                } else {
                    throw new Error(createResponse.message || 'Failed to create game');
                }
            }
        } catch (error) {
            console.error('Error joining game:', error);
            this.showError('Failed to join or create game. Please try again.');
        }
    }

    private async joinGame(playerName: string): Promise<any> {
        if (!this.ui.connect4Client) {
            throw new Error('WASM client not initialized');
        }

        return await this.ui.connect4Client.connect4Service.joinGame({
            gameId: this.ui.gameId,
            playerId: this.ui.playerId,
            playerName: playerName
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

        const response = await this.ui.connect4Client.connect4Service.createGame({
            gameId: this.ui.gameId,
            playerId: this.ui.playerId,
            playerName: playerName,
            config: gameConfig
        }) as any;

        if (!response.success) {
            throw new Error(response.message || 'Failed to create game');
        }

        return response;
    }

    public async dropPiece(column: number): Promise<void> {
        if (!this.ui.connect4Client || !this.ui.gameState) {
            console.error('Game not properly initialized');
            return;
        }

        try {
            const response = await this.ui.connect4Client.connect4Service.dropPiece({
                gameId: this.ui.gameId,
                playerId: this.ui.playerId,
                column: column
            }) as any;

            if (response.success) {
                this.ui.gameState = GameState.from(response.data);
                this.updateGameDisplay();
                
                // Send state update through stateful transport
                if (this.ui.transport) {
                    await this.ui.transport.sendPatches([{
                        operation: 'update',
                        path: '',
                        value: this.ui.gameState,
                        timestamp: Date.now(),
                        source: this.ui.playerId
                    }]);
                }
            } else {
                console.error('Failed to drop piece:', response.message);
            }
        } catch (error) {
            console.error('Error dropping piece:', error);
        }
    }

    private applyPatches(patches: any[]): void {
        for (const patch of patches) {
            if (patch.operation === 'update' && patch.value) {
                try {
                    const newState = GameState.from(patch.value);
                    if (newState) {
                        this.ui.gameState = newState;
                        this.updateGameDisplay();
                        this.addLogEntry('Game state updated from another player');
                    }
                } catch (error) {
                    console.error('Failed to apply patch:', error);
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

    private updateGameDisplay(): void {
        if (!this.ui.gameState) return;

        // Update turn information
        if (this.elements.turnNumber) {
            this.elements.turnNumber.textContent = this.ui.gameState.turnNumber.toString();
        }

        // Update current player
        const currentPlayer = this.ui.gameState.players.find(p => p.id === this.ui.gameState!.currentPlayerId);
        if (this.elements.currentPlayerName && currentPlayer) {
            this.elements.currentPlayerName.textContent = currentPlayer.name;
            
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

        // Update board
        this.initializeGameBoard();
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
                // Occupied slot
                slotsHTML += `
                    <div class="player-slot occupied" data-slot="${slotIndex}" style="border-left: 4px solid ${slotColor}">
                        <div class="slot-header">Player ${slotIndex + 1}</div>
                        <div class="player-info">
                            <span class="player-name">${player.name}</span>
                            <span class="player-color" style="color: ${slotColor}">●</span>
                        </div>
                        <div class="slot-status">Joined</div>
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

        this.ui.playerId = `player_${Date.now()}`;
        
        try {
            const response = await this.ui.connect4Client!.connect4Service.joinGame({
                gameId: this.ui.gameId,
                playerId: this.ui.playerId,
                playerName: playerName
            });

            if (response.success) {
                this.ui.gameState = GameState.from(response.data);
                this.showGameInterface();
                this.updateGameDisplay();
                this.addLogEntry(`${playerName} joined as Player ${this.selectedSlot + 1}`);
            } else {
                alert(`Failed to join slot: ${response.message}`);
            }
        } catch (error) {
            console.error('Error joining slot:', error);
            alert('Failed to join slot. Please try again.');
        }
    }

    // Public methods for HTML onclick handlers
    public resetGame(): void {
        if (confirm('Are you sure you want to start a new game?')) {
            window.location.reload();
        }
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
