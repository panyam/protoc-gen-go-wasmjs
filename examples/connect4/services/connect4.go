package services

import (
	"context"
	"fmt"
	"sync"
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

	// Create initial player
	player := &pb.Player{
		Id:          fmt.Sprintf("player_%d", time.Now().UnixNano()),
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
		Status:          pb.GameStatus_WAITING_FOR_PLAYERS,
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

	// Create new player
	player := &pb.Player{
		Id:          fmt.Sprintf("player_%d", time.Now().UnixNano()),
		Name:        req.PlayerName,
		Color:       "#0000ff", // Blue (for second player)
		IsConnected: true,
		JoinOrder:   int32(len(game.Players)),
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
		game.Status = pb.GameStatus_IN_PROGRESS
	}

	return &pb.JoinGameResponse{
		Success:       true,
		PlayerId:      player.Id,
		AssignedColor: player.Color,
		GameState:     game,
	}, nil
}

// GetGame retrieves the current state of a game
func (s *Connect4Service) GetGame(ctx context.Context, req *pb.GetGameRequest) (*pb.GameState, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	game, exists := s.games[req.GameId]
	if !exists {
		return nil, fmt.Errorf("game not found: %s", req.GameId)
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
	if game.Status != pb.GameStatus_IN_PROGRESS {
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
		game.Status = pb.GameStatus_FINISHED
		game.Winners = append(game.Winners, req.PlayerId)

		// Update player stats
		if stats, exists := game.PlayerStats[req.PlayerId]; exists {
			stats.HasWon = true
			stats.WinningLines = int32(len(winningLines))
		}

		patches = append(patches, &wasmjs.MessagePatch{
			Operation:    wasmjs.PatchOperation_SET,
			FieldPath:    "status",
			ValueJson:    fmt.Sprintf("%d", int(pb.GameStatus_FINISHED)),
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

	return &pb.DropPieceResponse{
		Success:      true,
		Patches:      patches,
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
