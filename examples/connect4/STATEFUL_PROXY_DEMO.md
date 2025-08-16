# Stateful Proxy System - Transport Architecture Demo

## Overview

This Connect4 example demonstrates a **pluggable stateful proxy system** with **enhanced player management** that can seamlessly switch between different transport mechanisms for real-time collaboration. Recent updates include player-specific URLs, simple player IDs, and improved cross-tab synchronization.

## Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Browser Tab A │    │   Browser Tab B │    │   HTTP Server   │
│                 │    │                 │    │                 │
│ ┌─────────────┐ │    │ ┌─────────────┐ │    │ ┌─────────────┐ │
│ │ Transport   │ │    │ │ Transport   │ │    │ │   Static    │ │
│ │ Layer       │ │    │ │ Layer       │ │    │ │   Files +   │ │
│ │             │ │    │ │             │ │    │ │   Templates │ │
│ │ IndexedDB + │◄────►│ │ IndexedDB + │ │    │ │             │ │
│ │ Polling     │ │    │ │ Polling     │ │    │ │ (Go server) │ │
│ └─────────────┘ │    │ └─────────────┘ │    │ └─────────────┘ │
│       ▲         │    │       ▲         │    │                 │
│       │ State   │    │       │ State   │    │                 │
│       ▼         │    │       ▼         │    │                 │
│ ┌─────────────┐ │    │ ┌─────────────┐ │    │                 │
│ │  WASM Game  │ │    │ │  WASM Game  │ │    │                 │
│ │  Service    │ │    │ │  Service    │ │    │                 │
│ └─────────────┘ │    │ └─────────────┘ │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## Transport Options

### 1. IndexedDB + Polling (Default)
- **Use Case**: Cross-tab collaboration with persistence
- **URL**: `http://localhost:8080/{gameId}` or `http://localhost:8080/{gameId}/players/{playerIndex}`
- **Pros**: Persistent across browser sessions, works offline, enhanced schema (v4) prevents constraint errors
- **Cons**: ~1 second polling delay

### 2. BroadcastChannel
- **Use Case**: Fast cross-tab collaboration within same browser
- **Switch**: `switchTransport('broadcast')` in console
- **Pros**: Zero-latency, instant updates
- **Cons**: Limited to same browser/origin, no persistence

### 3. WebSocket (Ready, needs server)
- **Use Case**: Real-time server-based collaboration
- **Switch**: `switchTransport('websocket')` in console
- **Pros**: True real-time across browsers/devices
- **Cons**: Requires WebSocket server endpoint

### 4. Server-Sent Events (Ready, needs server)
- **Use Case**: One-way server push with HTTP fallback
- **Switch**: `switchTransport('sse')` in console
- **Pros**: Simpler than WebSocket, works through proxies
- **Cons**: One-way only (sends via HTTP POST)

## Testing the System

### 1. Basic Cross-Tab Communication
1. **Start server**: `make all` or `make web`
2. **Open**: `http://localhost:8080/`
3. **Create game**: Enter Game ID "TestGame", Player Name "Creator"
   - Automatically redirected to `http://localhost:8080/TestGame/players/0`
4. **Open player-specific URL in new tab**: Navigate to `http://localhost:8080/TestGame/players/1`
   - Or use general URL `http://localhost:8080/TestGame` and select Player 2 from modal
5. **Join as Player2**: Enter Player Name "Player2"
6. **Test moves**: Drop pieces and watch real-time sync via enhanced IndexedDB + BroadcastChannel

### 2. Runtime Transport Switching
```javascript
// In browser console (F12):
showTransportStatus();           // Show current transport
switchTransport('broadcast');    // Switch to BroadcastChannel (faster)
switchTransport('indexeddb');    // Switch back to IndexedDB (persistent)
switchTransport('websocket');    // Switch to WebSocket (needs server)
```

### 3. Testing Game Persistence
```javascript
// Create a game, make moves, then:
localStorage.getItem('connect4_game_test-game'); // View stored state
// Refresh page - game should resume automatically
```

## Key Features

### Transport Abstraction
- **Pluggable**: Switch transports at runtime
- **Consistent API**: Same interface for all transport types
- **Fallback Support**: Graceful degradation when transports fail

### Patch-Based Updates
- **Differential**: Only changes are transmitted
- **Ordered**: Change numbers ensure correct sequencing
- **Conflict Resolution**: Timestamps and user IDs for merge conflicts

### State Management
- **Persistent**: LocalStorage for state persistence
- **Reactive**: UI automatically updates on state changes
- **Synchronized**: All clients see same state instantly

## Real-World Usage

### For Local Development
```javascript
// Use BroadcastChannel for rapid prototyping
const proxy = new StatefulProxy(gameId, 'broadcast');
```

### For Production (Future)
```javascript
// Use WebSocket for real-time collaboration across browsers
const transport = TransportFactory.create(gameId, 'websocket');
```

### For Offline-First Apps
```javascript
// Start with IndexedDB, upgrade to WebSocket when online
let transport = TransportFactory.create(gameId, 'indexeddb');
if (navigator.onLine) {
    await transport.switchTransport('websocket');
}
```

## WASM Integration Flow

The system integrates seamlessly with the WASM service using async callbacks:

1. **User Action** → Click column in UI
2. **TypeScript** → `gameViewer.dropPiece(column)` with player validation
3. **WASM Call** → `connect4Client.connect4Service.dropPiece({gameId, playerId: '0', column}, callback)`
4. **Go Service** → Validates move, updates game state with simple player ID, returns async response
5. **Async Callback** → Receives game state update without blocking browser
6. **State Update** → Update local game state from WASM callback
7. **Transport Broadcast** → `transport.sendPatches([...patches])` with uniqueId for reliability
8. **Cross-Tab Sync** → Other tabs receive patches via enhanced IndexedDB + BroadcastChannel
9. **UI Update** → All game UIs reflect the change instantly (BroadcastChannel) or within ~1 second (IndexedDB)

## Current vs Future Architecture

**Current (Working)**:
- **Enhanced local multiplayer**: Cross-tab via IndexedDB + polling and BroadcastChannel
- **Player-specific URLs**: Direct links like `/TestGame/players/0` with auto-selection
- **Simple player IDs**: Clean numeric indices (0, 1) instead of complex timestamps
- **Enhanced state persistence**: Survives browser restarts with player identity restoration
- **WASM integration**: Full game logic in WebAssembly with async callback support
- **Robust transport layer**: Enhanced IndexedDB schema (v4) with constraint error fixes
- **Transport switching**: Runtime pluggable transports with improved reliability

**Future (Server Required)**:
- **Cross-browser multiplayer**: WebSocket server coordination
- **Real-time updates**: <100ms latency via WebSocket
- **Conflict resolution**: Server-side authoritative state
- **Spectator mode**: Watch games without joining
