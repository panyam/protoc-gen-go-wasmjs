# Connect4 Architecture

## Overview

This Connect4 implementation demonstrates **stateful proxy patterns** with **async callback support** to enable real-time collaborative gaming without Browser→WASM→Browser deadlocks.

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

## Data Flow

### 1. Game Creation (Async Pattern)
```typescript
// Browser initiates
await client.connect4Service.createGame(gameConfig, (response, error) => {
  if (error) {
    handleError(error);
    return;
  }
  
  const gameState = JSON.parse(response);
  navigateToGame(gameState.gameId);
});

// WASM Flow:
// 1. WASM wrapper receives request + callback
// 2. Spawns goroutine immediately, returns success
// 3. Goroutine calls Connect4Service.CreateGame() 
// 4. Service saves to IndexedDB (no goroutine - direct call)
// 5. Goroutine invokes callback with result
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

### 3. Real-time Moves (Sync Pattern)
```typescript
// Immediate response for game moves
const response = await client.connect4Service.dropPiece({
  gameId,
  playerId,
  column
});

// Generates patches for other clients
const patches = response.patches;
transport.broadcastPatches(patches);
```

## Storage Strategy

### Primary Storage: IndexedDB
- **Persistent** across browser sessions
- **Structured** game state storage
- **Async API** (handled by callback pattern)

### Secondary Storage: LocalStorage  
- **Fast access** for UI state caching
- **Synchronous** API for immediate reads
- **Smaller capacity** for metadata only

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

This architecture enables **real-time collaborative gaming** while maintaining **responsive UIs** and **preventing deadlocks** in browser-WASM interactions.