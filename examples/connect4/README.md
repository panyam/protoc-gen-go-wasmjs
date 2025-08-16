# Connect4 with Stateful Proxies

This example demonstrates the **stateful proxy extension** for `protoc-gen-go-wasmjs`, which enables real-time collaborative applications with pluggable transport mechanisms.

## What are Stateful Proxies?

Stateful proxies provide:

1. **Differential Updates**: Instead of sending entire game states, only changes (patches) are transmitted
2. **Pluggable Transports**: IndexedDB + polling, BroadcastChannel, WebSocket, or SSE
3. **Conflict Resolution**: Built-in strategies for handling concurrent modifications  
4. **Type Safety**: Generated TypeScript clients with full type checking
5. **Persistent State**: LocalStorage + IndexedDB for cross-session continuity

## Game Features

- **2 players** with unique colors (configurable for more)
- **Simple player IDs** - Clean numeric indices (0, 1, 2) instead of complex timestamps
- **Player-specific URLs** - Direct links like `/GameName/players/0` for seamless sharing
- **Auto-player selection** - URLs automatically select the correct player
- **Standard 7x6 board** (configurable sizes)
- **Real-time collaboration** via pluggable transport system
- **Enhanced cross-tab sync** - Fixed IndexedDB constraint errors with proper schema management
- **Gravity-based piece placement**
- **Turn-based gameplay** with validation
- **Line detection** (horizontal, vertical, diagonal)
- **Game persistence** across browser sessions
- **Player selection modal** - Choose your player identity when accessing general game URLs

## Architecture

```
┌─────────────────┐    IndexedDB     ┌─────────────────┐
│   Browser Tab A │    + Polling     │   Browser Tab B │
│                 │ ◄──────────────► │                 │
│ ┌─────────────┐ │                  │ ┌─────────────┐ │
│ │ Transport   │ │                  │ │ Transport   │ │
│ │ Layer       │ │   BroadcastCh    │ │ Layer       │ │
│ │ (Pluggable) │ │ ◄──────────────► │ │ (Pluggable) │ │
│ └─────────────┘ │                  │ └─────────────┘ │
│        │        │                  │        │        │
│ ┌─────────────┐ │                  │ ┌─────────────┐ │
│ │ WASM        │ │                  │ │ WASM        │ │
│ │ Service     │ │                  │ │ Service     │ │
│ │ (Go)        │ │                  │ │ (Go)        │ │
│ └─────────────┘ │                  │ └─────────────┘ │
└─────────────────┘                  └─────────────────┘
                            │
                   ┌─────────────────┐
                   │ HTTP Server     │
                   │ (Go templates)  │
                   └─────────────────┘
```

## Getting Started

### Prerequisites
- Go 1.23+
- Node.js 18+ and pnpm
- Buf CLI (`brew install bufbuild/buf/buf`)

### Build & Run

1. **Build everything and start server:**
   ```bash
   make all
   # Opens http://localhost:8080
   ```

2. **Or step by step:**
   ```bash
   make generate  # Generate protobuf + WASM bindings
   make wasm      # Build WASM binary
   make web       # Build frontend + start server
   ```

### Test Stateful Proxy Generation

```bash
make test
```

## Project Structure

```
examples/connect4/
├── proto/connect4/
│   └── game.proto              # Service definitions with stateful annotations
├── cmd/
│   ├── wasm/main.go           # WASM binary entry point
│   └── serve/main.go          # HTTP server
├── services/
│   └── connect4.go            # Go service implementation
├── web/                       # Frontend (TypeScript + Webpack)
│   ├── src/                   # TypeScript source files
│   ├── static/gen/js/         # Compiled JavaScript bundles
│   ├── static/wasm/           # WASM binaries and loader
│   └── gen/wasmts/            # Generated TypeScript clients
├── gen/                       # Generated protobuf Go code
├── buf.yaml                   # Buf configuration
├── buf.gen.yaml              # Code generation config
└── Makefile                  # Build automation
```

## Key Stateful Features

### Service Annotations
```protobuf
service Connect4Service {
  option (wasmjs.v1.stateful) = {
    enabled: true
    state_message_type: "connect4.GameState"
    conflict_resolution: CHANGE_NUMBER_BASED
  };
  
  // Async method annotation for IndexedDB operations (prevents deadlocks)
  rpc GetGame(GetGameRequest) returns (GameState) {
    option (wasmjs.v1.async_method) = { is_async: true };
  };
  
  rpc CreateGame(CreateGameRequest) returns (CreateGameResponse) {
    option (wasmjs.v1.async_method) = { is_async: true };
  };
  
  // Stateful method for real-time updates
  rpc DropPiece(DropPieceRequest) returns (DropPieceResponse) {
    option (wasmjs.v1.stateful_method) = {
      returns_patches: true
      broadcasts: true
    };
  };
}
```

### Player-Specific URLs

The Connect4 example now supports clean, shareable player-specific URLs:

```
General game view:           /GameName
Player-specific URLs:        /GameName/players/0
                            /GameName/players/1
Direct creation and join:    Create game -> auto-redirect to /GameName/players/0
```

**Features:**
- **Simple Player IDs**: Uses clean numeric indices (0, 1) instead of complex timestamps
- **Auto-Selection**: Player-specific URLs automatically select the correct player
- **Direct Sharing**: Send `/TestGame/players/1` to let someone join as Player 2
- **Player Selection Modal**: General URLs show a modal to choose your player identity
- **Cross-Tab Consistency**: Player identity persists across tabs and browser sessions

### Generated TypeScript Client
```typescript
// Auto-generated WASM client with async method support
import Connect4Client from './gen/wasmts/multiplayer_connect4Client.client';

const client = new Connect4Client();
await client.loadWasm('/static/wasm/multiplayer_connect4.wasm');

// Async method calls with callbacks (prevents IndexedDB deadlocks)
await client.connect4Service.joinGame({ 
  gameId: 'my-game', 
  playerName: 'Player1' 
}, (response, error) => {
  if (error) {
    console.error('Failed to join game:', error);
    return;
  }
  
  const gameState = JSON.parse(response);
  console.log('Joined as player:', gameState.playerId); // Simple ID like "0" or "1"
});

// Synchronous method calls for game moves
const response = await client.connect4Service.dropPiece({
  gameId: 'my-game',
  playerId: '0',  // Simple numeric player ID
  column: 3
});
```

### Real-time Updates
```typescript
// Enhanced transport system with fixed IndexedDB schema
const transport = TransportFactory.create(gameId, 'indexeddb');

transport.subscribe((patches) => {
  // Apply incoming state changes with validation
  statefulProxy.applyPatches(patches);
  updateGameUI();
});

// Send state changes to other clients with proper patch structure
await transport.sendPatches([{
  operation: 'update',
  path: 'board.rows[2].cells[3]',
  value: '0',  // Simple player ID
  source: '0', // Who made the change
  timestamp: Date.now(),
  uniqueId: `${gameId}_${Date.now()}_${Math.random()}`
}]);
```

## Game Rules & Validation

### Valid Moves
- **Turn-based** - Only current player can move  
- **Column bounds** - Must be within board width  
- **Available space** - Column not full  
- **Gravity** - Pieces fall to lowest position  

### Invalid Moves
- **Out of turn** - Not your turn  
- **Full column** - No space available  
- **Invalid column** - Outside board bounds  

### Winning Conditions
- **Connect 4** - Horizontal, vertical, or diagonal  
- **Multiple winners** - Game continues (configurable)  
- **Score tracking** - Pieces played, lines formed  

## Performance Benefits

- **Local-First** - Instant UI feedback with localStorage persistence
- **Differential Patches** - Only send changes, not full state  
- **Transport Flexibility** - IndexedDB polling, BroadcastChannel, or WebSocket ready
- **Cross-Session Persistence** - Resume games across browser restarts
- **WASM Performance** - Game logic runs in compiled WebAssembly

## Development

### Clean Generated Files
```bash
make clean
```

### Development Mode (with local plugin changes)
```bash
make generate-dev  # Uses local wasmjs proto dependencies
```

### Frontend Development
```bash
cd web
pnpm install
npm run dev        # Webpack dev server with hot reload
```

### Customize Board Size
Edit `game.proto` and modify `GameConfig` defaults, then run `make generate`.

This example showcases how the **protoc-gen-go-wasmjs** plugin enables **real-time collaborative applications** with pluggable transport mechanisms!
