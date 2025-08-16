# Connect4 Architecture

## Overview

This Connect4 implementation demonstrates **stateful proxy patterns** with **async callback support** and **player management systems** to enable real-time collaborative gaming without Browser→WASM→Browser deadlocks. Recent enhancements include simple player IDs, player-specific URLs, and enhanced cross-tab synchronization.

## Key Architectural Components

### 1. Async Method Pattern (Deadlock Prevention)

**Problem**: When WASM methods need to call browser APIs (like IndexedDB), a deadlock occurs:
```
Browser Thread → WASM Method → IndexedDB Call → BLOCKED (main thread busy)
```

**Solution**: Async method annotations with callback parameters:
```
Browser Thread → WASM Method (with callback) → Returns Immediately
  └→ Goroutine → IndexedDB Call → Callback Invocation → Browser
```

#### Implementation:
- **Protobuf Annotation**: `option (wasmjs.v1.async_method) = { is_async: true };`
- **Generated TypeScript**: Methods accept `(request, callback)` parameters
- **WASM Wrapper**: Executes service calls in goroutines with `js.Value` callbacks
- **Service Layer**: Direct synchronous calls to storage (no internal goroutines)

### 2. Stateful Proxy System

Enables real-time collaboration through differential updates:

```typescript
// State changes generate patches
const patches = [
  { op: 'replace', path: '/board/rows/2/cells/3', value: 'player_123' },
  { op: 'replace', path: '/currentPlayerId', value: 'player_456' }
];

// Broadcast to other clients
transport.broadcastPatches(patches);
```

### 3. Transport Layer Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                     Browser Tabs                           │
├─────────────────┬─────────────────┬─────────────────────────┤
│   Tab A         │      Tab B      │       Tab C             │
│                 │                 │                         │
│ ┌─────────────┐ │ ┌─────────────┐ │ ┌─────────────────────┐ │
│ │ GameViewer  │ │ │ GameViewer  │ │ │ GameViewer          │ │
│ │  (UI Layer) │ │ │  (UI Layer) │ │ │  (UI Layer)         │ │
│ └─────────────┘ │ └─────────────┘ │ └─────────────────────┘ │
│        │        │        │        │           │             │
│ ┌─────────────┐ │ ┌─────────────┐ │ ┌─────────────────────┐ │
│ │ WASM Client │ │ │ WASM Client │ │ │ WASM Client         │ │
│ │ (Async API) │ │ │ (Async API) │ │ │ (Async API)         │ │
│ └─────────────┘ │ └─────────────┘ │ └─────────────────────┘ │
│        │        │        │        │           │             │
│ ┌─────────────┐ │ ┌─────────────┐ │ ┌─────────────────────┐ │
│ │ Connect4    │ │ │ Connect4    │ │ │ Connect4            │ │
│ │ Service     │ │ │ Service     │ │ │ Service             │ │
│ │ (Go/WASM)   │ │ │ (Go/WASM)   │ │ │ (Go/WASM)           │ │
│ └─────────────┘ │ └─────────────┘ │ └─────────────────────┘ │
└─────────────────┴─────────────────┴─────────────────────────┘
         │                  │                     │
         └──────────────────┼─────────────────────┘
                            │
              ┌─────────────────────────┐
              │     Storage Layer       │
              │                         │
              │ ┌─────────┐ ┌─────────┐ │
              │ │IndexedDB│ │LocalStor│ │
              │ │(Persist)│ │ (Cache) │ │
              │ └─────────┘ └─────────┘ │
              │                         │
              │ ┌─────────────────────┐ │
              │ │ BroadcastChannel    │ │
              │ │ (Real-time Sync)    │ │
              │ └─────────────────────┘ │
              └─────────────────────────┘
```

### 4. Player Management System

The Connect4 implementation includes a comprehensive player management system with simple IDs and URL-based selection:

#### Player ID System
- **Simple Numeric IDs**: Players use clean indices (0, 1, 2) instead of complex timestamps
- **Frontend-Backend Consistency**: Both Go service and TypeScript client use the same ID format
- **URL Integration**: Player IDs map directly to URL paths (`/GameName/players/0`)

#### URL-Based Player Selection
```typescript
// URL patterns supported:
// /GameName              - General view with player selection modal
// /GameName/players/0    - Auto-select Player 1 (index 0)
// /GameName/players/1    - Auto-select Player 2 (index 1)

// URL parsing in gameViewer.ts
const urlParts = window.location.pathname.split('/');
if (urlParts.length >= 4 && urlParts[2] === 'players') {
  const urlPlayerIndex = parseInt(urlParts[3], 10);
  this.ui.playerId = urlPlayerIndex.toString();
}
```

#### Player Identity Persistence
- **LocalStorage**: Maintains player identity across sessions
- **Cross-Tab Sync**: Player selection persists across browser tabs
- **Session Recovery**: Automatic player restoration on page reload

#### Player Selection Flow
1. **Direct URL**: `/GameName/players/0` automatically selects Player 1
2. **General URL**: `/GameName` shows player selection modal
3. **Join Redirect**: After joining, redirects to player-specific URL
4. **Cross-Tab Consistency**: Player identity synced via localStorage

## Data Flow

### 1. Game Creation (Async Pattern)
```typescript
// Browser initiates game creation
await client.connect4Service.createGame(gameConfig, (response, error) => {
  if (error) {
    handleError(error);
    return;
  }
  
  const gameState = JSON.parse(response);
  // Auto-redirect to player-specific URL for creator (always player 0)
  window.location.href = `/${gameState.gameId}/players/0`;
});

// WASM Flow:
// 1. WASM wrapper receives request + callback
// 2. Spawns goroutine immediately, returns success
// 3. Goroutine calls Connect4Service.CreateGame() 
// 4. Service creates player with simple ID "0" and saves to IndexedDB
// 5. Goroutine invokes callback with result containing player ID "0"
```

### 2. Game Loading (Async Pattern)
```typescript
// Browser requests game state
await client.connect4Service.getGame({ gameId }, (response, error) => {
  if (response) {
    updateGameUI(JSON.parse(response));
  }
});

// WASM Flow:
// 1. WASM wrapper spawns goroutine for storage operation
// 2. Service loads from IndexedDB synchronously 
// 3. Callback receives parsed game state
// 4. UI updates without deadlock
```

### 3. Real-time Moves (Async Pattern with Patches)
```typescript
// Real-time piece drop with async callback
await client.connect4Service.dropPiece({
  gameId,
  playerId: '0',  // Simple numeric player ID
  column
}, (response, error) => {
  if (response) {
    const result = JSON.parse(response);
    // Patches automatically sent via transport layer
    updateGameUI(result.gameState);
  }
});

// Enhanced patch generation with simple player IDs
const patches = [{
  operation: 'update',
  path: 'board.rows[2].cells[3]',
  value: '0',  // Simple player ID as cell value
  source: '0', // Who made the move
  timestamp: Date.now(),
  uniqueId: `${gameId}_${Date.now()}_${Math.random()}`
}];
```

## Storage Strategy

### Primary Storage: IndexedDB
- **Persistent** across browser sessions
- **Structured** game state storage
- **Async API** (handled by callback pattern)
- **Enhanced Schema** (v4): Fixed constraint errors with auto-increment patches store
- **Unique Patch IDs**: Prevents duplicate key conflicts during cross-tab synchronization

### Secondary Storage: LocalStorage  
- **Fast access** for UI state caching
- **Synchronous** API for immediate reads
- **Smaller capacity** for metadata only
- **Player Identity**: Stores selected player ID and name for session persistence

### Real-time Sync: BroadcastChannel
- **Tab-to-tab** communication
- **Instant** state synchronization
- **Same origin** only

## Key Design Principles

### 1. Deadlock-Free Architecture
- **Async annotations** prevent Browser→WASM→Browser blocking
- **Goroutine isolation** keeps main thread responsive  
- **Direct callback invocation** for clean async handling

### 2. Local-First Design
- **Immediate UI feedback** with local state
- **Background synchronization** with other clients
- **Offline capability** with persistent storage

### 3. Type-Safe Code Generation
- **Protobuf definitions** drive all interfaces
- **Generated TypeScript** with full type checking
- **Consistent APIs** across sync/async patterns

### 4. Transport Flexibility
- **Pluggable backends** (IndexedDB, WebSocket, SSE)
- **Configurable conflict resolution**
- **Extensible patch formats**
- **Enhanced Reliability**: Fixed IndexedDB constraints and added patch validation
- **Race Condition Prevention**: Validates patches to prevent empty state overwrites

This architecture enables **real-time collaborative gaming** while maintaining **responsive UIs** and **preventing deadlocks** in browser-WASM interactions.