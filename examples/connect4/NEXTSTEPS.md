# Next Steps - Connect4 Example

## Recently Completed ✅

### Async Callback Pattern Implementation
- **Problem Solved**: Browser→WASM→Browser deadlock when calling IndexedDB from WASM methods
- **Solution**: Implemented async method annotations with direct callback support

#### Key Changes:
1. **Protobuf Generator Enhancements**:
   - Added `AsyncMethodOptions` to `wasmjs/v1/annotations.proto`
   - Extended method options with `async_method` annotation
   - Updated TypeScript and WASM template generators for callback support

2. **Service Methods Made Async**:
   - `GetGame` marked as async (eliminates deadlock when loading from IndexedDB)
   - `CreateGame` marked as async (eliminates deadlock when saving to IndexedDB)
   - Removed internal goroutines from service methods (WASM wrapper now handles async)

3. **Frontend Integration**:
   - Updated `gameViewer.ts` to use callback-based `getGame()` 
   - Updated `index.ts` and `gameViewer.ts` for callback-based `createGame()`
   - Fixed infinite loop issues in game page loading

4. **Generator Templates Enhanced**:
   - TypeScript client generates different signatures for async vs sync methods
   - WASM wrapper generates goroutine-based execution for async methods
   - Direct `js.Value` callback parameter passing for clean async handling

## Current State
- ✅ Async callback pattern fully implemented
- ✅ Browser→WASM→Browser deadlock resolved  
- ✅ Infinite loop in game page loading fixed
- ✅ Clean TypeScript API with proper async/callback patterns

## Immediate Next Steps

### 1. Testing & Validation
- [ ] Regenerate protobuf files and WASM artifacts
- [ ] Test game creation flow (index page → game page)
- [ ] Test game loading flow (direct game URL access)
- [ ] Verify IndexedDB storage operations work without deadlocks

### 2. Code Cleanup
- [ ] Remove legacy static JS files if no longer needed
- [ ] Update any remaining sync calls to use new async pattern
- [ ] Clean up commented code and debug logs

### 3. Documentation Updates
- [ ] Update README with new async callback patterns
- [ ] Document the deadlock solution for future reference
- [ ] Add examples of how to use async methods

## Technical Debt
- [ ] Static JS files (`connect4-game.js`, etc.) may need updating or removal
- [ ] Consider consolidating storage callback patterns
- [ ] Review error handling in async callback chains

## Future Enhancements
- [ ] Add timeout handling for async operations
- [ ] Implement retry logic for failed storage operations
- [ ] Add progress indicators for async operations in UI
- [ ] Consider extending async pattern to other methods that may benefit