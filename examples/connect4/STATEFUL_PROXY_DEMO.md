# Stateful Proxy System - Transport Architecture Demo

## Overview

This Connect4 example demonstrates a **pluggable stateful proxy system** that can seamlessly switch between different transport mechanisms for real-time collaboration.

## Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Browser Tab A │    │   Browser Tab B │    │   Server/WASM   │
│                 │    │                 │    │                 │
│ ┌─────────────┐ │    │ ┌─────────────┐ │    │ ┌─────────────┐ │
│ │StatefulProxy│ │    │ │StatefulProxy│ │    │ │   Service   │ │
│ │             │ │    │ │             │ │    │ │   Logic     │ │
│ │ Transport:  │ │    │ │ Transport:  │ │    │ │             │ │
│ │ BroadCast   │◄────►│ │ BroadCast   │ │    │ │ Generates   │ │
│ │ Channel     │ │    │ │ Channel     │ │    │ │ Patches     │ │
│ └─────────────┘ │    │ └─────────────┘ │    │ └─────────────┘ │
│       ▲         │    │       ▲         │    │       ▲         │
│       │ Patches │    │       │ Patches │    │       │ Calls   │
│       ▼         │    │       ▼         │    │       ▼         │
│ ┌─────────────┐ │    │ ┌─────────────┐ │    │ ┌─────────────┐ │
│ │   Game UI   │ │    │ │   Game UI   │ │    │ │  WASM API   │ │
│ └─────────────┘ │    │ └─────────────┘ │    │ └─────────────┘ │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## Transport Options

### 1. BroadcastChannel (Default)
- **Use Case**: Cross-tab collaboration within same browser
- **URL**: `http://localhost:8000/game.html?gameId=test123`
- **Pros**: Zero-latency, perfect for local multiplayer
- **Cons**: Limited to same browser/origin

### 2. WebSocket
- **Use Case**: Real-time server-based collaboration
- **URL**: `http://localhost:8000/game.html?gameId=test123&transport=websocket`
- **Pros**: True real-time across browsers/devices
- **Cons**: Requires WebSocket server

### 3. Server-Sent Events (SSE)
- **Use Case**: One-way server push with HTTP fallback
- **URL**: `http://localhost:8000/game.html?gameId=test123&transport=sse`
- **Pros**: Simpler than WebSocket, works through proxies
- **Cons**: One-way only (sends via HTTP POST)

## Testing the System

### 1. Basic Cross-Tab Communication
1. Open: `http://localhost:8000/games.html`
2. Create a new game
3. Open the same game URL in multiple tabs
4. Watch real-time updates via BroadcastChannel

### 2. Runtime Transport Switching
```javascript
// In browser console:
switchTransport('websocket');  // Switch to WebSocket
switchTransport('broadcast');  // Switch back to BroadcastChannel
switchTransport('sse');        // Switch to Server-Sent Events
```

### 3. Manual Patch Testing
```javascript
// Test patch application directly:
testPatch(0, 'currentPlayerId', '"player_123"');  // SET operation
testPatch(1, 'winners', '"player_456"');          // INSERT_LIST operation
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

### For Production
```javascript
// Use WebSocket for real-time collaboration
const proxy = new StatefulProxy(gameId, 'websocket');
```

### For Offline-First Apps
```javascript
// Start with local storage, upgrade to WebSocket when online
let proxy = new StatefulProxy(gameId, 'broadcast');
if (navigator.onLine) {
    await proxy.switchTransport('websocket');
}
```

## WASM Integration

The system integrates seamlessly with the WASM service:

1. **User Action** → `dropPiece(column)`
2. **WASM Call** → `connect4Client.connect4Service.dropPiece(...)`
3. **Patch Generation** → Service returns patches array
4. **Proxy Application** → `statefulProxy.applyPatchesFromWasm(patches)`
5. **Transport Broadcast** → Other tabs receive patches automatically
6. **UI Update** → All UIs reflect the change instantly

This creates a seamless real-time collaborative experience with the game logic running in WASM but the state synchronization handled by the pluggable transport layer.
