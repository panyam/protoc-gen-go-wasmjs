package services

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"syscall/js"
	"time"

	pb "github.com/panyam/protoc-gen-go-wasmjs/examples/connect4/gen/go/connect4"
	wasmjs "github.com/panyam/protoc-gen-go-wasmjs/examples/connect4/gen/go/wasmjs/v1"
)

// Connect4Service implements the Connect4 game logic with stateful patch generation
type Connect4Service struct {
	pb.UnimplementedConnect4ServiceServer

	games         map[string]*pb.GameState
	changeCounter int64
	mu            sync.RWMutex

	// Storage callbacks from browser
	saveCallback js.Value
	loadCallback js.Value
	pollCallback js.Value
	callbacksSet bool
}

// NewConnect4Service creates a new Connect4 service instance
func NewConnect4Service() *Connect4Service {
	return &Connect4Service{
		games: make(map[string]*pb.GameState),
	}
}

// CreateGame creates a new Connect4 game
func (s *Connect4Service) CreateGame(ctx context.Context, req *pb.CreateGameRequest) (*pb.CreateGameResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	gameId := req.GameId
	if gameId == "" {
		gameId = fmt.Sprintf("game_%d", time.Now().Unix())
	}

	// Create game configuration
	config := &pb.GameConfig{
		BoardWidth:    7,
		BoardHeight:   6,
		ConnectLength: 4,
		MaxPlayers:    2,
		MinPlayers:    2,
	}
	if req.Config != nil {
		config = req.Config
	}

	// Initialize empty board
	board := &pb.GameBoard{
		Width:         config.BoardWidth,
		Height:        config.BoardHeight,
		Rows:          make([]*pb.BoardRow, config.BoardHeight),
		ColumnHeights: make([]int32, config.BoardWidth),
	}

	for i := int32(0); i < config.BoardHeight; i++ {
		board.Rows[i] = &pb.BoardRow{
			Cells: make([]string, config.BoardWidth),
		}
	}

	// Create initial player with simple index-based ID
	player := &pb.Player{
		Id:          "0", // First player always gets index 0
		Name:        req.CreatorName,
		Color:       "#ff0000", // Red
		IsConnected: true,
		JoinOrder:   0,
	}

	// Create initial game state
	game := &pb.GameState{
		GameId:          gameId,
		Config:          config,
		Board:           board,
		Players:         []*pb.Player{player},
		CurrentPlayerId: player.Id,
		Status:          pb.GameStatus_GAME_STATUS_WAITING_FOR_PLAYERS,
		TurnNumber:      0,
		PlayerStats:     make(map[string]*pb.PlayerStats),
	}

	// Initialize player stats
	game.PlayerStats[player.Id] = &pb.PlayerStats{
		PiecesPlayed:  0,
		WinningLines:  0,
		HasWon:        false,
		TotalMoveTime: 0,
	}

	s.games[gameId] = game

	// Save to storage via callback (async - fire and forget for now)
	s.saveGameToStorage(gameId, game)

	return &pb.CreateGameResponse{
		Success:   true,
		PlayerId:  player.Id,
		GameState: game,
	}, nil
}

// JoinGame allows a player to join an existing game
func (s *Connect4Service) JoinGame(ctx context.Context, req *pb.JoinGameRequest) (*pb.JoinGameResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	game, exists := s.games[req.GameId]
	if !exists {
		return &pb.JoinGameResponse{
			Success:      false,
			ErrorMessage: fmt.Sprintf("Game not found: %s", req.GameId),
		}, nil
	}

	// Check if player is already in game (by name)
	for _, player := range game.Players {
		if player.Name == req.PlayerName {
			return &pb.JoinGameResponse{
				Success:       true,
				PlayerId:      player.Id,
				GameState:     game,
				AssignedColor: player.Color,
			}, nil
		}
	}

	// Check if game is full
	if len(game.Players) >= int(game.Config.MaxPlayers) {
		return &pb.JoinGameResponse{
			Success:      false,
			ErrorMessage: "Game is full",
		}, nil
	}

	// Create new player with simple index-based ID
	playerIndex := len(game.Players)
	player := &pb.Player{
		Id:          fmt.Sprintf("%d", playerIndex), // Simple index: "0", "1", "2"...
		Name:        req.PlayerName,
		Color:       "#0000ff", // Blue (for second player)
		IsConnected: true,
		JoinOrder:   int32(playerIndex),
	}

	// Add player to game
	game.Players = append(game.Players, player)

	// Initialize player stats
	game.PlayerStats[player.Id] = &pb.PlayerStats{
		PiecesPlayed:  0,
		WinningLines:  0,
		HasWon:        false,
		TotalMoveTime: 0,
	}

	// Start game if we have enough players
	if len(game.Players) >= int(game.Config.MinPlayers) {
		game.Status = pb.GameStatus_GAME_STATUS_IN_PROGRESS
	}

	// Save updated game to storage (async)
	s.saveGameToStorage(req.GameId, game)

	return &pb.JoinGameResponse{
		Success:       true,
		PlayerId:      player.Id,
		AssignedColor: player.Color,
		GameState:     game,
	}, nil
}

// GetGame retrieves the current state of a game
func (s *Connect4Service) GetGame(ctx context.Context, req *pb.GetGameRequest) (*pb.GameState, error) {
	s.mu.Lock() // Use write lock since we might load from storage
	defer s.mu.Unlock()

	game, exists := s.games[req.GameId]
	if !exists {
		// Try to load from storage
		loadedGame := s.loadGameFromStorage(req.GameId)
		if loadedGame != nil {
			s.games[req.GameId] = loadedGame
			return loadedGame, nil
		}
	}

	return game, nil
}

// DropPiece handles a player dropping a piece in a column
func (s *Connect4Service) DropPiece(ctx context.Context, req *pb.DropPieceRequest) (*pb.DropPieceResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	game, exists := s.games[req.GameId]
	if !exists {
		return &pb.DropPieceResponse{
			Success:      false,
			ErrorMessage: fmt.Sprintf("Game not found: %s", req.GameId),
		}, nil
	}

	// Validate move
	if game.Status != pb.GameStatus_GAME_STATUS_IN_PROGRESS {
		return &pb.DropPieceResponse{
			Success:      false,
			ErrorMessage: "Game is not in progress",
		}, nil
	}

	if game.CurrentPlayerId != req.PlayerId {
		return &pb.DropPieceResponse{
			Success:      false,
			ErrorMessage: "Not your turn",
		}, nil
	}

	if req.Column < 0 || req.Column >= game.Config.BoardWidth {
		return &pb.DropPieceResponse{
			Success:      false,
			ErrorMessage: fmt.Sprintf("Invalid column: %d", req.Column),
		}, nil
	}

	if game.Board.ColumnHeights[req.Column] >= game.Config.BoardHeight {
		return &pb.DropPieceResponse{
			Success:      false,
			ErrorMessage: fmt.Sprintf("Column is full: %d", req.Column),
		}, nil
	}

	// Calculate the row where the piece will land
	row := game.Config.BoardHeight - 1 - game.Board.ColumnHeights[req.Column]

	// Place the piece
	game.Board.Rows[row].Cells[req.Column] = req.PlayerId
	game.Board.ColumnHeights[req.Column]++

	// Generate patches for the state changes
	s.changeCounter++
	var patches []*wasmjs.MessagePatch

	// Patch for placing the piece
	patches = append(patches, &wasmjs.MessagePatch{
		Operation:    wasmjs.PatchOperation_SET,
		FieldPath:    fmt.Sprintf("board.rows[%d].cells[%d]", row, req.Column),
		ValueJson:    fmt.Sprintf(`"%s"`, req.PlayerId),
		ChangeNumber: s.changeCounter,
		Timestamp:    time.Now().UnixMicro(),
	})

	// Patch for updating column height
	patches = append(patches, &wasmjs.MessagePatch{
		Operation:    wasmjs.PatchOperation_SET,
		FieldPath:    fmt.Sprintf("board.columnHeights[%d]", req.Column),
		ValueJson:    fmt.Sprintf("%d", game.Board.ColumnHeights[req.Column]),
		ChangeNumber: s.changeCounter,
		Timestamp:    time.Now().UnixMicro(),
	})

	// Check for winning lines
	winningLines := s.checkForWinningLines(game, row, req.Column, req.PlayerId)

	result := &pb.PieceDropResult{
		FinalRow:     row,
		FinalColumn:  req.Column,
		FormedLine:   len(winningLines) > 0,
		WinningLines: winningLines,
	}

	if len(winningLines) > 0 {
		// Player wins!
		game.Status = pb.GameStatus_GAME_STATUS_FINISHED
		game.Winners = append(game.Winners, req.PlayerId)

		// Update player stats
		if stats, exists := game.PlayerStats[req.PlayerId]; exists {
			stats.HasWon = true
			stats.WinningLines = int32(len(winningLines))
		}

		patches = append(patches, &wasmjs.MessagePatch{
			Operation:    wasmjs.PatchOperation_SET,
			FieldPath:    "status",
			ValueJson:    fmt.Sprintf("%d", int(pb.GameStatus_GAME_STATUS_FINISHED)),
			ChangeNumber: s.changeCounter,
			Timestamp:    time.Now().UnixMicro(),
		})

		patches = append(patches, &wasmjs.MessagePatch{
			Operation:    wasmjs.PatchOperation_INSERT_LIST,
			FieldPath:    "winners",
			ValueJson:    fmt.Sprintf(`"%s"`, req.PlayerId),
			ChangeNumber: s.changeCounter,
			Timestamp:    time.Now().UnixMicro(),
		})
	} else {
		// Switch to next player
		currentPlayerIndex := -1
		for i, player := range game.Players {
			if player.Id == req.PlayerId {
				currentPlayerIndex = i
				break
			}
		}
		if currentPlayerIndex >= 0 {
			nextPlayerIndex := (currentPlayerIndex + 1) % len(game.Players)
			game.CurrentPlayerId = game.Players[nextPlayerIndex].Id
		}
		game.TurnNumber++

		patches = append(patches, &wasmjs.MessagePatch{
			Operation:    wasmjs.PatchOperation_SET,
			FieldPath:    "currentPlayerId",
			ValueJson:    fmt.Sprintf(`"%s"`, game.CurrentPlayerId),
			ChangeNumber: s.changeCounter,
			Timestamp:    time.Now().UnixMicro(),
		})

		patches = append(patches, &wasmjs.MessagePatch{
			Operation:    wasmjs.PatchOperation_SET,
			FieldPath:    "turnNumber",
			ValueJson:    fmt.Sprintf("%d", game.TurnNumber),
			ChangeNumber: s.changeCounter,
			Timestamp:    time.Now().UnixMicro(),
		})
	}

	// Update player stats
	if stats, exists := game.PlayerStats[req.PlayerId]; exists {
		stats.PiecesPlayed++
	}

	// Save updated game state to storage (async)
	s.saveGameToStorage(req.GameId, game)

	return &pb.DropPieceResponse{
		Success:      true,
		ChangeNumber: s.changeCounter,
		Result:       result,
	}, nil
}

// checkForWinningLines checks if placing a piece creates any winning lines
func (s *Connect4Service) checkForWinningLines(game *pb.GameState, row, col int32, playerId string) []*pb.LineInfo {
	var winningLines []*pb.LineInfo

	directions := []struct {
		dr, dc    int32
		direction string
	}{
		{0, 1, "horizontal"},
		{1, 0, "vertical"},
		{1, 1, "diagonal_up"},
		{1, -1, "diagonal_down"},
	}

	for _, dir := range directions {
		positions := s.findLineInDirection(game, row, col, dir.dr, dir.dc, playerId)
		if len(positions) >= int(game.Config.ConnectLength) {
			winningLines = append(winningLines, &pb.LineInfo{
				Positions: positions,
				Direction: dir.direction,
				Length:    int32(len(positions)),
			})
		}
	}

	return winningLines
}

// findLineInDirection finds the longest continuous line in a given direction
func (s *Connect4Service) findLineInDirection(game *pb.GameState, row, col, dr, dc int32, playerId string) []*pb.Position {
	var positions []*pb.Position

	// Add the current cell
	positions = append(positions, &pb.Position{Row: row, Column: col})

	// Check in positive direction
	r, c := row+dr, col+dc
	for r >= 0 && r < game.Config.BoardHeight && c >= 0 && c < game.Config.BoardWidth &&
		game.Board.Rows[r].Cells[c] == playerId {
		positions = append(positions, &pb.Position{Row: r, Column: c})
		r, c = r+dr, c+dc
	}

	// Check in negative direction
	r, c = row-dr, col-dc
	for r >= 0 && r < game.Config.BoardHeight && c >= 0 && c < game.Config.BoardWidth &&
		game.Board.Rows[r].Cells[c] == playerId {
		// Insert at beginning to maintain order
		positions = append([]*pb.Position{{Row: r, Column: c}}, positions...)
		r, c = r-dr, c-dc
	}

	return positions
}

// SetStorageCallbacks configures the browser storage callbacks
func (s *Connect4Service) SetStorageCallbacks(saveCallback, loadCallback, pollCallback js.Value) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.saveCallback = saveCallback
	s.loadCallback = loadCallback
	s.pollCallback = pollCallback
	s.callbacksSet = true

	fmt.Println("Storage callbacks configured successfully")
}

// HandleExternalStorageChange processes external changes to game state
func (s *Connect4Service) HandleExternalStorageChange(gameId, gameStateJson string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Parse the updated game state
	var gameState pb.GameState
	if err := json.Unmarshal([]byte(gameStateJson), &gameState); err != nil {
		fmt.Printf("Failed to unmarshal external game state for %s: %v\n", gameId, err)
		return
	}

	// Update our internal state
	s.games[gameId] = &gameState
	fmt.Printf("Updated game %s from external storage change\n", gameId)
}

// saveGameToStorage saves game state to browser storage via callback
func (s *Connect4Service) saveGameToStorage(gameId string, game *pb.GameState) {
	if !s.callbacksSet || s.saveCallback.IsUndefined() {
		return // No callback configured
	}

	// Convert game state to JSON
	gameStateJson, err := json.Marshal(game)
	if err != nil {
		fmt.Printf("Failed to marshal game state for %s: %v\n", gameId, err)
		return
	}

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Panic in saveGameToStorage for %s: %v\n", gameId, r)
		}
	}()

	// Call save callback synchronously (WASM wrapper now handles async)
	promise := s.saveCallback.Invoke(gameId, string(gameStateJson))

	// Await the result to log success/errors
	result, errValues := await(promise)
	if errValues != nil && len(errValues) > 0 {
		fmt.Printf("Failed to save game state for %s: %v\n", gameId, errValues[0])
	} else {
		fmt.Printf("Successfully saved game state for %s\n", gameId)
		if len(result) > 0 && !result[0].IsNull() {
			// Browser can return additional info if needed
		}
	}
}

// loadGameFromStorage loads game state from browser storage via callback
func (s *Connect4Service) loadGameFromStorage(gameId string) *pb.GameState {
	fmt.Println("Load callback called?: ", s.callbacksSet, s.loadCallback)
	if !s.callbacksSet || s.loadCallback.IsUndefined() {
		return nil // No callback configured
	}

	// Call browser load callback synchronously (WASM wrapper now handles async)
	promise := s.loadCallback.Invoke(gameId)

	// Await the Promise result synchronously
	result, errValues := await(promise)
	if errValues != nil && len(errValues) > 0 {
		fmt.Printf("Failed to load game state for %s: %v\n", gameId, errValues[0])
		return nil
	}

	fmt.Println("Load result: ", result)
	if len(result) == 0 || result[0].IsNull() || result[0].IsUndefined() {
		fmt.Printf("No game state found for %s\n", gameId)
		return nil // Game not found
	}

	gameStateJson := result[0].String()
	if gameStateJson != "" {
		// Parse the JSON
		var gameState pb.GameState
		if err := json.Unmarshal([]byte(gameStateJson), &gameState); err != nil {
			fmt.Printf("Failed to unmarshal loaded game state for %s: %v\n", gameId, err)
			return nil
		}

		fmt.Printf("Successfully loaded game %s from storage\n", gameId)
		return &gameState
	}

	return nil
}

// await helper for handling JavaScript Promises in Go WASM
func await(awaitable js.Value) ([]js.Value, []js.Value) {
	then := make(chan []js.Value, 1)
	defer close(then)
	thenFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		fmt.Println("Did we come to THEN?")
		then <- args
		// go func() { }()
		return nil
	})
	defer thenFunc.Release()

	catch := make(chan []js.Value, 1)
	defer close(catch)
	catchFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		fmt.Println("Did we come to catch?")
		catch <- args
		return nil
	})
	defer catchFunc.Release()

	res := awaitable.Call("then", thenFunc)
	res = res.Call("catch", catchFunc)

	t := time.NewTicker(1000 * time.Millisecond)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			fmt.Println("Here 7??")
		case result := <-then:
			return result, nil
		case err := <-catch:
			return nil, err
		}
	}
}
