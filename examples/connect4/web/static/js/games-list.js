// Games List Application - Main JavaScript Module
// Handles game discovery, creation, and navigation

// Application state
let gamesList = [];

// UI Elements
const gamesListContainer = document.getElementById('gamesList');
const createGameForm = document.getElementById('createGameForm');

// Initialize the page
function initializePage() {
    loadGames();
    startAutoRefresh();
}

// Load games from localStorage
function loadGames() {
    const games = [];
    
    // Scan localStorage for game states
    for (let i = 0; i < localStorage.length; i++) {
        const key = localStorage.key(i);
        if (key && key.startsWith('connect4_game_')) {
            try {
                const gameId = key.replace('connect4_game_', '');
                const gameState = JSON.parse(localStorage.getItem(key));
                
                if (gameState && gameState.gameId) {
                    games.push({
                        id: gameId,
                        state: gameState,
                        lastUpdated: new Date(gameState.lastUpdated || Date.now())
                    });
                }
            } catch (e) {
                console.error('Failed to parse game state for key:', key, e);
            }
        }
    }
    
    // Sort by last updated
    games.sort((a, b) => b.lastUpdated - a.lastUpdated);
    
    gamesList = games;
    renderGamesList();
}

// Render the games list
function renderGamesList() {
    if (!gamesListContainer) return;
    
    gamesListContainer.innerHTML = '';
    
    if (gamesList.length === 0) {
        gamesListContainer.innerHTML = `
            <div style="text-align: center; padding: 40px; color: #6b7280;">
                <h3>No active games</h3>
                <p>Create a new game to get started!</p>
            </div>
        `;
        return;
    }
    
    gamesList.forEach(game => {
        const gameCard = createGameCard(game);
        gamesListContainer.appendChild(gameCard);
    });
}

// Create a game card element
function createGameCard(game) {
    const card = document.createElement('div');
    card.className = 'game-card';
    
    const statusText = getGameStatusText(game.state);
    const playerCount = game.state.players ? game.state.players.length : 0;
    const maxPlayers = game.state.config ? game.state.config.maxPlayers : 2;
    
    card.innerHTML = `
        <div class="game-card-header">
            <h3>Game: ${game.id}</h3>
            <span class="game-status ${game.state.status?.toLowerCase()}">${statusText}</span>
        </div>
        <div class="game-card-body">
            <div class="game-info">
                <div class="info-item">
                    <span class="label">Players:</span>
                    <span class="value">${playerCount}/${maxPlayers}</span>
                </div>
                <div class="info-item">
                    <span class="label">Turn:</span>
                    <span class="value">#${game.state.turnNumber || 0}</span>
                </div>
                <div class="info-item">
                    <span class="label">Last Updated:</span>
                    <span class="value">${formatTime(game.lastUpdated)}</span>
                </div>
            </div>
            ${game.state.players && game.state.players.length > 0 ? `
                <div class="players-preview">
                    ${game.state.players.map(player => `
                        <div class="player-chip">
                            <div class="player-color" style="background-color: ${player.color}"></div>
                            <span>${player.name}</span>
                        </div>
                    `).join('')}
                </div>
            ` : ''}
        </div>
        <div class="game-card-actions">
            <a href="/${game.id}" class="btn btn-primary">
                ${playerCount < maxPlayers ? 'Join Game' : 'View Game'}
            </a>
            <button class="btn btn-secondary" onclick="deleteGame('${game.id}')">Delete</button>
        </div>
    `;
    
    return card;
}

// Create a new game
function createNewGame() {
    const gameNameInput = document.getElementById('gameName');
    const playerNameInput = document.getElementById('playerName');
    
    const gameName = gameNameInput.value.trim();
    const playerName = playerNameInput.value.trim();
    
    if (!gameName) {
        alert('Please enter a game name');
        return;
    }
    
    if (!playerName) {
        alert('Please enter your name');
        return;
    }
    
    // Generate a clean game ID from the name
    const gameId = gameName.toLowerCase()
        .replace(/[^a-z0-9\s]/g, '')
        .replace(/\s+/g, '-')
        .substring(0, 20);
    
    if (!gameId) {
        alert('Please enter a valid game name');
        return;
    }
    
    // Check if game already exists
    const existingGame = localStorage.getItem(`connect4_game_${gameId}`);
    if (existingGame) {
        if (!confirm(`Game "${gameId}" already exists. Do you want to join it instead?`)) {
            return;
        }
    }
    
    // Navigate to the game page where it will be created
    window.location.href = `/${gameId}`;
}

// Delete a game
function deleteGame(gameId) {
    if (confirm(`Are you sure you want to delete game "${gameId}"?`)) {
        localStorage.removeItem(`connect4_game_${gameId}`);
        localStorage.removeItem(`connect4_patches_${gameId}`);
        
        // Trigger event for other windows
        window.dispatchEvent(new CustomEvent('gameStateChanged', { 
            detail: { gameId, state: null } 
        }));
        
        loadGames(); // Refresh the list
        addNotification(`Game "${gameId}" deleted successfully`);
    }
}

// Get human-readable status text
function getGameStatusText(gameState) {
    switch (gameState.status) {
        case 'WAITING_FOR_PLAYERS': return 'Waiting for Players';
        case 'IN_PROGRESS': return 'In Progress';
        case 'FINISHED': return 'Finished';
        default: return 'Unknown';
    }
}

// Format time for display
function formatTime(date) {
    const now = new Date();
    const diff = now - date;
    
    if (diff < 60000) { // Less than 1 minute
        return 'Just now';
    } else if (diff < 3600000) { // Less than 1 hour
        const minutes = Math.floor(diff / 60000);
        return `${minutes} minute${minutes !== 1 ? 's' : ''} ago`;
    } else if (diff < 86400000) { // Less than 1 day
        const hours = Math.floor(diff / 3600000);
        return `${hours} hour${hours !== 1 ? 's' : ''} ago`;
    } else {
        return date.toLocaleDateString();
    }
}

// Add notification message
function addNotification(message) {
    const notification = document.createElement('div');
    notification.className = 'notification';
    notification.textContent = message;
    
    document.body.appendChild(notification);
    
    // Auto-remove after 3 seconds
    setTimeout(() => {
        notification.remove();
    }, 3000);
}

// Listen for game state changes from other windows
window.addEventListener('gameStateChanged', function(e) {
    loadGames(); // Refresh the list when games change
});

// Listen for localStorage changes from other windows
window.addEventListener('storage', function(e) {
    if (e.key && e.key.startsWith('connect4_game_')) {
        loadGames(); // Refresh the list
    }
});

// Auto-refresh games list every 5 seconds
function startAutoRefresh() {
    setInterval(() => {
        loadGames();
    }, 5000);
}

// Handle Enter key in create game form
function handleCreateGameKeyPress(event) {
    if (event.key === 'Enter') {
        createNewGame();
    }
}

// Make functions globally available
window.createNewGame = createNewGame;
window.deleteGame = deleteGame;
window.handleCreateGameKeyPress = handleCreateGameKeyPress;

// Initialize the page when DOM is loaded
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', initializePage);
} else {
    initializePage();
}
