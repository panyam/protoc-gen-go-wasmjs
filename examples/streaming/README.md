# Server Streaming Example

This example demonstrates **Phase 1** of the streaming implementation: **Server-side streaming** from WASM to JavaScript.

## Overview

Server streaming allows a WASM service to send multiple responses for a single request using a callback pattern. This is perfect for:
- Live data feeds (ticks, prices, events)
- Progress updates
- Real-time notifications
- Incremental data loading

## Proto Definition

```proto
service TickService {
    // Server streaming: sends multiple responses for one request
    rpc GetTicks(TickRequest) returns (stream TickResponse);
    
    // Regular unary method for comparison
    rpc GetSingleTick(TickRequest) returns (TickResponse);
}
```

## Generated API

### Server Streaming Method
```typescript
client.tickService.getTicks(
    request,
    (response, error, done) => {
        if (error) {
            console.error('Stream error:', error);
            return false; // Stop stream
        }
        if (done) {
            console.log('Stream complete');
            return false; 
        }
        
        console.log('Got tick:', response);
        return true; // Continue streaming
    }
);
```

### Regular Unary Method
```typescript
const singleTick = await client.tickService.getSingleTick(request);
```

## Usage

1. **Generate code**: `make generate-dev`
2. **Build WASM**: `make wasm` 
3. **Implement service**: Provide `TickServiceServer` implementation
4. **Load in browser**: `client.loadWasm()` then call streaming methods

## Implementation Details

### WASM Side (Generated)
- Runs in goroutine to avoid blocking
- Calls `stream.Recv()` in loop until `io.EOF`
- Invokes JavaScript callback: `callback(responseJSON, errorString, isDone)`
- Respects user cancellation via callback return value

### TypeScript Side (Generated)  
- Method signature: `(request, callback) => void`
- Callback type: `(response | null, error | null, done: boolean) => boolean`
- Returns `true` to continue, `false` to stop stream

## Next Phases

- **Phase 2**: Client streaming with connection objects (`conn.send()`, `conn.close()`)
- **Phase 3**: Bidirectional streaming (combination of both)

## Testing

```bash
cd examples/streaming
make test  # Generate and validate streaming code
```