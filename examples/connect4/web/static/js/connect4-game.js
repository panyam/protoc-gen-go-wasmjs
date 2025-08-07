// Connect4 Game Application - Main JavaScript Module
// Handles game initialization, UI updates, and WASM integration

// Application state
let gameState = null;
let currentPlayerId = null;
let currentGameId = null;
let connect4Client = null;
let statefulProxy = null;

// UI Elements
const joinGameForm = document.getElementById('joinGameForm');
const gameInterface = document.getElementById('gameInterface');
const errorState = document.getElementById('errorState');
const gameBoard = document.getElementById('gameBoard');
const gameStatus = document.getElementById('gameStatus');
const currentPlayerName = document.getElementById('currentPlayerName');
const currentPlayerColor = document.getElementById('currentPlayerColor');
const turnNumber = document.getElementById('turnNumber');
const playersList = document.getElementById('playersList');
const gameLog = document.getElementById('gameLog');

// Initialize WASM client
async function initializeWasmClient() {
    try {
        addLogEntry('Loading WASM module...');
        connect4Client = new window.Connect4Client();
        await connect4Client.loadWasm('./static/wasm/multiplayer_connect4.wasm');
        await connect4Client.waitUntilReady();
        addLogEntry('WASM module loaded successfully!');
        return true;
    } catch (error) {
        console.error('Failed to load WASM:', error);
        addLogEntry(`Failed to load WASM module: ${error.message}`);
        return false;
    }
}

// Get URL parameters
function getUrlParameter(name) {
    const urlParams = new URLSearchParams(window.location.search);
    return urlParams.get(name);
}

// Game state management
function loadGameState(gameId) {
    const stored = localStorage.getItem(`connect4_game_${gameId}`);
    if (stored) {
        try {
            return JSON.parse(stored);
        } catch (e) {
            console.error('Failed to parse stored game state:', e);
        }
    }
    return null;
}

function saveGameState(gameId, state) {
    localStorage.setItem(`connect4_game_${gameId}`, JSON.stringify(state));
    // Trigger storage event for other windows
    window.dispatchEvent(new CustomEvent('gameStateChanged', { detail: { gameId, state } }));
}

function deleteGameState(gameId) {
    localStorage.removeItem(`connect4_game_${gameId}`);
    window.dispatchEvent(new CustomEvent('gameStateChanged', { detail: { gameId, state: null } }));
}

// Listen for game state changes from other windows
window.addEventListener('gameStateChanged', function(e) {
    const { gameId, state } = e.detail;
    if (gameId === currentGameId) {
        if (!state) {
            // Game was deleted
            addLogEntry(`Game "${gameId}" was deleted - all players have left`);
            setTimeout(() => {
                window.location.href = '/';
            }, 2000);
        } else if (gameState) {
            // Game state updated
            gameState = state;
            updateUI();
            addLogEntry('Game state updated from another window');
        }
    }
});

// Also listen for localStorage changes from other windows
window.addEventListener('storage', function(e) {
    if (e.key === `connect4_game_${currentGameId}`) {
        const newState = e.newValue ? JSON.parse(e.newValue) : null;
        if (!newState) {
            addLogEntry(`Game "${currentGameId}" was deleted - all players have left`);
            setTimeout(() => {
                window.location.href = '/';
            }, 2000);
        } else if (gameState) {
            gameState = newState;
            updateUI();
            addLogEntry('Game state updated from another window');
        }
    }
});

// Page initialization
async function initializePage() {
    // Check if we're on a game page (has gameId in URL path)
    const pathParts = window.location.pathname.split('/').filter(Boolean);
    if (pathParts.length > 0) {
        currentGameId = pathParts[0];
        const playerId = getUrlParameter('playerId');

        if (!currentGameId) {
            showErrorState();
            return;
        }

        // Initialize WASM client first
        const wasmReady = await initializeWasmClient();
        if (!wasmReady) {
            showErrorState();
            return;
        }

        document.getElementById('currentGameId').textContent = currentGameId;
        document.getElementById('gameUrl').textContent = `${window.location.origin}/${currentGameId}`;

        // Try to get game state from WASM service
        try {
            const wasmGameState = await connect4Client.connect4Service.getGame({ gameId: currentGameId });
            if (wasmGameState) {
                gameState = wasmGameState;
                
                // Initialize stateful proxy for cross-page synchronization
                statefulProxy = new StatefulProxy(currentGameId, 'indexeddb');
                statefulProxy.setGameState(gameState);
                statefulProxy.onStateChange((newState) => {
                    gameState = newState;
                    updateUI();
                    addLogEntry('Game state synchronized from other tab/page');
                });
                window.statefulProxy = statefulProxy; // For debugging
                
                // Check if user was already in this game
                if (playerId && gameState.players.some(p => p.id === playerId)) {
                    currentPlayerId = playerId;
                    showGameInterface();
                    addLogEntry(`Reconnected as ${gameState.players.find(p => p.id === playerId).name}`);
                } else {
                    showJoinForm();
                }
            } else {
                // Game doesn't exist in WASM service yet
                showJoinForm();
            }
        } catch (error) {
            console.log('Game not found in WASM service, will create new one');
            showJoinForm();
        }

        // Set up periodic sync with localStorage for cross-tab communication
        startLocalStorageSync();
    }
}

function showErrorState() {
    errorState.classList.remove('hidden');
    joinGameForm.classList.add('hidden');
    gameInterface.classList.add('hidden');
}

function showJoinForm() {
    joinGameForm.classList.remove('hidden');
    gameInterface.classList.add('hidden');
    errorState.classList.add('hidden');
}

function showGameInterface() {
    joinGameForm.classList.add('hidden');
    gameInterface.classList.remove('hidden');
    errorState.classList.add('hidden');
    initializeBoard(gameState.config.boardWidth, gameState.config.boardHeight);
    updateUI();
}

// Join current game - now calls WASM service
async function joinCurrentGame() {
    const playerName = document.getElementById('playerName').value;
    
    if (!playerName) {
        alert('Please enter your name');
        return;
    }

    if (!connect4Client) {
        alert('WASM module not loaded. Please refresh the page.');
        return;
    }

    try {
        addLogEntry(`Joining game "${currentGameId}" as ${playerName}...`);

        // Try to join existing game first
        let response = await connect4Client.connect4Service.joinGame({
            gameId: currentGameId,
            playerName: playerName,
            preferredColor: '' // Let service assign color
        });

        if (response.success) {
            currentPlayerId = response.playerId;
            gameState = response.gameState;
            
            // Initialize stateful proxy for cross-page synchronization
            if (!statefulProxy) {
                statefulProxy = new StatefulProxy(currentGameId, 'indexeddb');
                statefulProxy.onStateChange((newState) => {
                    gameState = newState;
                    updateUI();
                    addLogEntry('Game state synchronized from other tab/page');
                });
                window.statefulProxy = statefulProxy; // For debugging
            }
            statefulProxy.setGameState(gameState);
            
            // Update URL with player ID for reconnection
            const newUrl = `/${currentGameId}?playerId=${currentPlayerId}`;
            window.history.replaceState({}, '', newUrl);
            
            // Save to localStorage for cross-tab sync
            saveGameState(currentGameId, gameState);
            
            showGameInterface();
            addLogEntry(`Successfully joined game! You are ${playerName} (${response.assignedColor})`);

            if (gameState.status === 'IN_PROGRESS') {
                addLogEntry('Game started! Enough players have joined.');
            }
        } else {
            // Join failed, try to create the game
            addLogEntry('Game not found, creating new game...');
            await createNewGameAndJoin(currentGameId, playerName);
        }

    } catch (error) {
        console.error('Error joining game:', error);
        addLogEntry(`Failed to join game: ${error.message}`);
        
        // Fallback: try to create new game
        try {
            await createNewGameAndJoin(currentGameId, playerName);
        } catch (createError) {
            console.error('Error creating game:', createError);
            addLogEntry(`Failed to create game: ${createError.message}`);
            alert('Failed to join or create game. Please try again.');
        }
    }
}

async function createNewGameAndJoin(gameId, playerName) {
    try {
        addLogEntry(`Creating new game "${gameId}"...`);

        const response = await connect4Client.connect4Service.createGame({
            gameId: gameId,
            config: {
                boardWidth: 7,
                boardHeight: 6,
                minPlayers: 2,
                maxPlayers: 4,
                connectLength: 4,
                allowMultipleWinners: false,
                moveTimeoutSeconds: 30
            },
            creatorName: playerName
        });

        if (response.success) {
            currentPlayerId = response.playerId;
            gameState = response.gameState;
            
            // Initialize stateful proxy for cross-page synchronization
            if (!statefulProxy) {
                statefulProxy = new StatefulProxy(gameId, 'indexeddb');
                statefulProxy.onStateChange((newState) => {
                    gameState = newState;
                    updateUI();
                    addLogEntry('Game state synchronized from other tab/page');
                });
                window.statefulProxy = statefulProxy; // For debugging
            }
            statefulProxy.setGameState(gameState);
            
            // Update URL with player ID
            const newUrl = `/${gameId}?playerId=${currentPlayerId}`;
            window.history.replaceState({}, '', newUrl);
            
            // Save to localStorage for cross-tab sync
            saveGameState(gameId, gameState);
            
            showGameInterface();
            addLogEntry(`Game "${gameId}" created successfully! You are the first player.`);
            addLogEntry('Waiting for more players to start...');
        } else {
            throw new Error(response.errorMessage || 'Failed to create game');
        }
    } catch (error) {
        console.error('Error creating game:', error);
        throw error;
    }
}

// Initialize the board
function initializeBoard(width = 7, height = 6) {
    gameBoard.innerHTML = '';
    gameBoard.style.gridTemplateColumns = `repeat(${width}, 1fr)`;
    
    for (let row = 0; row < height; row++) {
        for (let col = 0; col < width; col++) {
            const cell = document.createElement('div');
            cell.className = 'cell';
            cell.dataset.row = row;
            cell.dataset.col = col;
            cell.onclick = () => dropPiece(col);
            gameBoard.appendChild(cell);
        }
    }
}

// Update the UI with current game state
function updateUI() {
    if (!gameState) return;

    // Update game status
    gameStatus.textContent = getStatusText(gameState.status);
    turnNumber.textContent = gameState.turnNumber;

    // Update current player info
    if (gameState.currentPlayerId) {
        const currentPlayer = gameState.players.find(p => p.id === gameState.currentPlayerId);
        if (currentPlayer) {
            currentPlayerName.textContent = currentPlayer.name;
            currentPlayerColor.style.backgroundColor = currentPlayer.color;
        }
    }

    // Update board
    updateBoard();

    // Update players list
    updatePlayersList();
}

// Update the board display
function updateBoard() {
    const cells = gameBoard.querySelectorAll('.cell');
    
    cells.forEach(cell => {
        const row = parseInt(cell.dataset.row);
        const col = parseInt(cell.dataset.col);
        const playerId = gameState.board.rows[row].cells[col];
        
        if (playerId) {
            const player = gameState.players.find(p => p.id === playerId);
            if (player) {
                cell.style.backgroundColor = player.color;
                cell.classList.add('occupied');
            }
        } else {
            cell.style.backgroundColor = '#3b82f6';
            cell.classList.remove('occupied');
        }
    });
}

// Update players list
function updatePlayersList() {
    playersList.innerHTML = '';
    
    gameState.players.forEach(player => {
        const playerDiv = document.createElement('div');
        playerDiv.className = 'player-item';
        
        const colorDiv = document.createElement('div');
        colorDiv.className = 'player-color';
        colorDiv.style.backgroundColor = player.color;
        
        const nameSpan = document.createElement('span');
        nameSpan.textContent = player.name;
        if (player.id === currentPlayerId) {
            nameSpan.textContent += ' (You)';
        }
        
        const statusSpan = document.createElement('span');
        statusSpan.textContent = player.isConnected ? ' • Online' : ' • Offline';
        statusSpan.style.color = player.isConnected ? '#10b981' : '#ef4444';
        
        playerDiv.appendChild(colorDiv);
        playerDiv.appendChild(nameSpan);
        playerDiv.appendChild(statusSpan);
        playersList.appendChild(playerDiv);
    });
}

// Drop a piece in the specified column
async function dropPiece(column) {
    if (!gameState || gameState.status !== 'IN_PROGRESS') {
        addLogEntry('Game is not in progress');
        return;
    }

    if (gameState.currentPlayerId !== currentPlayerId) {
        addLogEntry('It\'s not your turn');
        return;
    }

    if (gameState.board.columnHeights[column] >= gameState.config.boardHeight) {
        addLogEntry('Column is full');
        return;
    }

    try {
        // Simulate piece drop (in real implementation, this would call WASM service)
        const row = gameState.config.boardHeight - 1 - gameState.board.columnHeights[column];
        
        // Place the piece
        gameState.board.rows[row].cells[column] = currentPlayerId;
        gameState.board.columnHeights[column]++;
        
        // Check for winning lines
        const winningLines = checkForWinningLines(row, column, currentPlayerId);
        
        if (winningLines.length > 0) {
            gameState.winners.push(currentPlayerId);
            gameState.status = 'FINISHED';
            highlightWinningLines(winningLines);
            
            const player = gameState.players.find(p => p.id === currentPlayerId);
            addLogEntry(`🎉 ${player.name} wins!`);
        } else {
            // Advance to next player
            const currentIndex = gameState.players.findIndex(p => p.id === gameState.currentPlayerId);
            const nextIndex = (currentIndex + 1) % gameState.players.length;
            gameState.currentPlayerId = gameState.players[nextIndex].id;
            gameState.turnNumber++;
            
            const player = gameState.players.find(p => p.id === currentPlayerId);
            const nextPlayer = gameState.players.find(p => p.id === gameState.currentPlayerId);
            addLogEntry(`${player.name} dropped piece in column ${column + 1}. ${nextPlayer.name}'s turn.`);
        }

        // Save updated game state via stateful proxy
        if (statefulProxy) {
            statefulProxy.setGameState(gameState);
        } else {
            saveGameState(currentGameId, gameState);
        }
        updateUI();

    } catch (error) {
        console.error('Error dropping piece:', error);
        addLogEntry('Failed to drop piece');
    }
}

// Check for winning lines (simplified version)
function checkForWinningLines(row, col, playerId) {
    const lines = [];
    const directions = [
        { dx: 1, dy: 0, name: 'horizontal' },
        { dx: 0, dy: 1, name: 'vertical' },
        { dx: 1, dy: 1, name: 'diagonal_down' },
        { dx: 1, dy: -1, name: 'diagonal_up' }
    ];

    directions.forEach(dir => {
        const positions = [];
        
        // Check both directions from the placed piece
        for (let direction = -1; direction <= 1; direction += 2) {
            let r = row;
            let c = col;
            
            while (r >= 0 && r < gameState.config.boardHeight && 
                   c >= 0 && c < gameState.config.boardWidth &&
                   gameState.board.rows[r].cells[c] === playerId) {
                positions.push({ row: r, column: c });
                r += direction * dir.dy;
                c += direction * dir.dx;
            }
        }

        if (positions.length >= gameState.config.connectLength) {
            lines.push({ positions, direction: dir.name });
        }
    });

    return lines;
}

// Highlight winning lines on the board
function highlightWinningLines(lines) {
    lines.forEach(line => {
        line.positions.forEach(pos => {
            const cell = document.querySelector(`[data-row="${pos.row}"][data-col="${pos.column}"]`);
            if (cell) {
                cell.classList.add('winning-line');
            }
        });
    });
}

// Utility functions
function getStatusText(status) {
    switch (status) {
        case 'WAITING_FOR_PLAYERS': return 'Waiting for players...';
        case 'IN_PROGRESS': return 'Game in progress';
        case 'FINISHED': return 'Game finished';
        default: return 'Unknown status';
    }
}

function addLogEntry(message) {
    const entry = document.createElement('div');
    entry.className = 'log-entry';
    entry.textContent = `[${new Date().toLocaleTimeString()}] ${message}`;
    gameLog.appendChild(entry);
    gameLog.scrollTop = gameLog.scrollHeight;
}

// Game management functions
function resetGame() {
    if (confirm('Are you sure you want to start a new game?')) {
        if (gameState) {
            gameState.board.rows.forEach(row => row.cells.fill(''));
            gameState.board.columnHeights.fill(0);
            gameState.status = 'IN_PROGRESS';
            gameState.turnNumber = 0;
            gameState.winners = [];
            gameState.currentPlayerId = gameState.players[0].id;
            
            // Save via stateful proxy
            if (statefulProxy) {
                statefulProxy.setGameState(gameState);
            } else {
                saveGameState(currentGameId, gameState);
            }
            updateUI();
            addLogEntry('New game started!');
            
            // Remove winning line highlights
            document.querySelectorAll('.winning-line').forEach(cell => {
                cell.classList.remove('winning-line');
            });
        }
    }
}

function leaveGame() {
    if (confirm('Are you sure you want to leave the game?')) {
        // Clean up stateful proxy
        if (statefulProxy) {
            statefulProxy.destroy();
            statefulProxy = null;
        }
        
        // Remove current player from game state
        if (gameState && currentPlayerId) {
            gameState.players = gameState.players.filter(p => p.id !== currentPlayerId);
            
            // If no players left, delete the game entirely
            if (gameState.players.length === 0) {
                deleteGameState(currentGameId);
                addLogEntry(`Game "${currentGameId}" deleted - all players have left`);
            } else {
                // Still have players, just save the updated state
                saveGameState(currentGameId, gameState);
                addLogEntry(`You left the game. ${gameState.players.length} player(s) remaining.`);
            }
        }
        
        // Navigate back to games list
        window.location.href = '/';
    }
}

function startLocalStorageSync() {
    // Auto-refresh game state every 2 seconds for real-time updates
    setInterval(() => {
        if (currentGameId && gameState) {
            const latestState = loadGameState(currentGameId);
            if (latestState && JSON.stringify(latestState) !== JSON.stringify(gameState)) {
                gameState = latestState;
                updateUI();
            }
        }
    }, 2000);
}

// Make functions globally available
window.joinCurrentGame = joinCurrentGame;
window.resetGame = resetGame;
window.leaveGame = leaveGame;

// Initialize the page when DOM is loaded
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', initializePage);
} else {
    initializePage();
}
