package main

import (
	"context"
	"fmt"
	"time"

	pb "github.com/panyam/protoc-gen-go-wasmjs/examples/connect4/gen/go/connect4"
	wasmjs "github.com/panyam/protoc-gen-go-wasmjs/examples/connect4/gen/go/wasmjs/v1"
)

type Connect4Service struct {
	games         map[string]*pb.GameState
	changeCounter int64
}

func NewConnect4Service() *Connect4Service {
	return &Connect4Service{
		games: make(map[string]*pb.GameState),
	}
}

func (s *Connect4Service) GetGame(ctx context.Context, req *pb.GetGameRequest) (*pb.GameState, error) {
	game, exists := s.games[req.GameId]
	if !exists {
		return nil, fmt.Errorf("game not found")
	}
	return game, nil
}

func (s *Connect4Service) CreateGame(ctx context.Context, req *pb.CreateGameRequest) (*pb.CreateGameResponse, error) {
	// Validate config
	config := req.Config
	if config.BoardWidth < 7 || config.BoardWidth > 20 {
		return &pb.CreateGameResponse{Success: false, ErrorMessage: "Board width must be 7-20"}, nil
	}
	if config.BoardHeight < 6 || config.BoardHeight > 20 {
		return &pb.CreateGameResponse{Success: false, ErrorMessage: "Board height must be 6-20"}, nil
	}

	// Initialize empty board
	board := &pb.GameBoard{
		Width:         config.BoardWidth,
		Height:        config.BoardHeight,
		Rows:          make([]*pb.BoardRow, config.BoardHeight),
		ColumnHeights: make([]int32, config.BoardWidth),
	}

	// Create empty rows
	for i := int32(0); i < config.BoardHeight; i++ {
		row := &pb.BoardRow{
			Cells: make([]string, config.BoardWidth),
		}
		board.Rows[i] = row
	}

	game := &pb.GameState{
		GameId:             req.GameId,
		Config:             config,
		Players:            []*pb.Player{},
		Board:              board,
		Status:             pb.GameStatus_WAITING_FOR_PLAYERS,
		TurnNumber:         0,
		Winners:            []string{},
		PlayerStats:        make(map[string]*pb.PlayerStats),
		MoveTimeoutSeconds: config.MoveTimeoutSeconds,
	}

	s.games[req.GameId] = game

	// Add creator as first player
	joinResp, _ := s.JoinGame(ctx, &pb.JoinGameRequest{
		GameId:     req.GameId,
		PlayerName: req.CreatorName,
	})

	return &pb.CreateGameResponse{
		Success:   true,
		PlayerId:  joinResp.PlayerId,
		GameState: game,
	}, nil
}

func (s *Connect4Service) JoinGame(ctx context.Context, req *pb.JoinGameRequest) (*pb.JoinGameResponse, error) {
	game, exists := s.games[req.GameId]
	if !exists {
		return &pb.JoinGameResponse{Success: false, ErrorMessage: "Game not found"}, nil
	}

	if int32(len(game.Players)) >= game.Config.MaxPlayers {
		return &pb.JoinGameResponse{Success: false, ErrorMessage: "Game is full"}, nil
	}

	// Assign color
	colors := []string{"#FF0000", "#0000FF", "#00FF00", "#FFFF00", "#FF00FF", "#00FFFF", "#FFA500", "#800080", "#FFC0CB", "#A52A2A"}
	playerColor := colors[len(game.Players)%len(colors)]

	playerId := fmt.Sprintf("player_%d_%d", len(game.Players)+1, time.Now().Unix())

	player := &pb.Player{
		Id:          playerId,
		Name:        req.PlayerName,
		Color:       playerColor,
		IsConnected: true,
		JoinOrder:   int32(len(game.Players)),
	}

	game.Players = append(game.Players, player)
	game.PlayerStats[playerId] = &pb.PlayerStats{}

	// Start game if we have minimum players
	if int32(len(game.Players)) >= game.Config.MinPlayers && game.Status == pb.GameStatus_WAITING_FOR_PLAYERS {
		game.Status = pb.GameStatus_IN_PROGRESS
		game.CurrentPlayerId = game.Players[0].Id
		game.LastMoveTime = time.Now().Unix()
	}

	return &pb.JoinGameResponse{
		Success:       true,
		PlayerId:      playerId,
		AssignedColor: playerColor,
		GameState:     game,
	}, nil
}

func (s *Connect4Service) DropPiece(ctx context.Context, req *pb.DropPieceRequest) (*pb.DropPieceResponse, error) {
	game, exists := s.games[req.GameId]
	if !exists {
		return &pb.DropPieceResponse{Success: false, ErrorMessage: "Game not found"}, nil
	}

	// Validate move
	if err := s.validateMove(game, req.PlayerId, req.Column); err != nil {
		return &pb.DropPieceResponse{Success: false, ErrorMessage: err.Error()}, nil
	}

	s.changeCounter++

	// Calculate where piece lands (gravity)
	row := s.findLowestAvailableRow(game.Board, req.Column)

	// Place the piece
	game.Board.Rows[row].Cells[req.Column] = req.PlayerId
	game.Board.ColumnHeights[req.Column]++

	// Update player stats
	if game.PlayerStats[req.PlayerId] == nil {
		game.PlayerStats[req.PlayerId] = &pb.PlayerStats{}
	}
	game.PlayerStats[req.PlayerId].PiecesPlayed++

	// Check for winning lines
	winningLines := s.checkForWinningLines(game.Board, row, req.Column, req.PlayerId, game.Config.ConnectLength)

	var patches []*wasmjs.MessagePatch

	// Generate patch for piece placement
	patches = append(patches, &wasmjs.MessagePatch{
		Operation:    wasmjs.PatchOperation_SET,
		FieldPath:    fmt.Sprintf("board.rows[%d].cells[%d]", row, req.Column),
		ValueJson:    req.PlayerId,
		ChangeNumber: s.changeCounter,
		Timestamp:    time.Now().UnixMicro(),
	})

	// Update column height
	patches = append(patches, &wasmjs.MessagePatch{
		Operation:    wasmjs.PatchOperation_SET,
		FieldPath:    fmt.Sprintf("board.column_heights[%d]", req.Column),
		ValueJson:    fmt.Sprintf("%d", game.Board.ColumnHeights[req.Column]),
		ChangeNumber: s.changeCounter,
		Timestamp:    time.Now().UnixMicro(),
	})

	result := &pb.PieceDropResult{
		FinalRow:     row,
		FinalColumn:  req.Column,
		FormedLine:   len(winningLines) > 0,
		WinningLines: winningLines,
	}

	// Handle winning lines
	if len(winningLines) > 0 {
		if !contains(game.Winners, req.PlayerId) {
			game.Winners = append(game.Winners, req.PlayerId)
			game.PlayerStats[req.PlayerId].HasWon = true

			// Add winner patch
			winnerIndex := int32(len(game.Winners) - 1)
			winnerValue := req.PlayerId
			patches = append(patches, &wasmjs.MessagePatch{
				Operation:    wasmjs.PatchOperation_INSERT_LIST,
				FieldPath:    "winners",
				Index:        winnerIndex,
				ValueJson:    winnerValue,
				ChangeNumber: s.changeCounter,
				Timestamp:    time.Now().UnixMicro(),
			})
		}

		game.PlayerStats[req.PlayerId].WinningLines += int32(len(winningLines))

		// Check if game should end
		if !game.Config.AllowMultipleWinners || s.isBoardFull(game.Board) {
			game.Status = pb.GameStatus_FINISHED
			statusValue := fmt.Sprintf("%d", int32(pb.GameStatus_FINISHED))
			patches = append(patches, &wasmjs.MessagePatch{
				Operation:    wasmjs.PatchOperation_SET,
				FieldPath:    "status",
				ValueJson:    statusValue,
				ChangeNumber: s.changeCounter,
				Timestamp:    time.Now().UnixMicro(),
			})
		}
	}

	// Advance turn (if game not finished)
	if game.Status == pb.GameStatus_IN_PROGRESS {
		s.advanceToNextPlayer(game)

		currentPlayerValue := game.CurrentPlayerId
		patches = append(patches, &wasmjs.MessagePatch{
			Operation:    wasmjs.PatchOperation_SET,
			FieldPath:    "current_player_id",
			ValueJson:    currentPlayerValue,
			ChangeNumber: s.changeCounter,
			Timestamp:    time.Now().UnixMicro(),
		})

		turnNumberValue := fmt.Sprintf("%d", game.TurnNumber)
		patches = append(patches, &wasmjs.MessagePatch{
			Operation:    wasmjs.PatchOperation_SET,
			FieldPath:    "turn_number",
			ValueJson:    turnNumberValue,
			ChangeNumber: s.changeCounter,
			Timestamp:    time.Now().UnixMicro(),
		})

		game.LastMoveTime = time.Now().Unix()
		lastMoveTimeValue := fmt.Sprintf("%d", game.LastMoveTime)
		patches = append(patches, &wasmjs.MessagePatch{
			Operation:    wasmjs.PatchOperation_SET,
			FieldPath:    "last_move_time",
			ValueJson:    lastMoveTimeValue,
			ChangeNumber: s.changeCounter,
			Timestamp:    time.Now().UnixMicro(),
		})
	}

	return &pb.DropPieceResponse{
		Success:      true,
		Patches:      patches,
		ChangeNumber: s.changeCounter,
		Result:       result,
	}, nil
}

func (s *Connect4Service) validateMove(game *pb.GameState, playerId string, column int32) error {
	if game.Status != pb.GameStatus_IN_PROGRESS {
		return fmt.Errorf("game not in progress")
	}

	if game.CurrentPlayerId != playerId {
		return fmt.Errorf("not your turn")
	}

	if column < 0 || column >= game.Board.Width {
		return fmt.Errorf("invalid column")
	}

	if game.Board.ColumnHeights[column] >= game.Board.Height {
		return fmt.Errorf("column is full")
	}

	return nil
}

func (s *Connect4Service) findLowestAvailableRow(board *pb.GameBoard, column int32) int32 {
	for row := board.Height - 1; row >= 0; row-- {
		if board.Rows[row].Cells[column] == "" {
			return row
		}
	}
	return -1 // Should never happen if validation passed
}

func (s *Connect4Service) checkForWinningLines(board *pb.GameBoard, row, col int32, playerId string, connectLength int32) []*pb.LineInfo {
	var lines []*pb.LineInfo

	directions := []struct {
		name   string
		dx, dy int32
	}{
		{"horizontal", 1, 0},
		{"vertical", 0, 1},
		{"diagonal_up", 1, -1},
		{"diagonal_down", 1, 1},
	}

	for _, dir := range directions {
		positions := s.findLineInDirection(board, row, col, dir.dx, dir.dy, playerId, connectLength)
		if len(positions) >= int(connectLength) {
			lines = append(lines, &pb.LineInfo{
				Positions: positions,
				Direction: dir.name,
				Length:    int32(len(positions)),
			})
		}
	}

	return lines
}

func (s *Connect4Service) findLineInDirection(board *pb.GameBoard, startRow, startCol, dx, dy int32, playerId string, minLength int32) []*pb.Position {
	var positions []*pb.Position

	// Check backward and forward from starting position
	for direction := -1; direction <= 1; direction += 2 {
		stepX, stepY := dx*int32(direction), dy*int32(direction)

		for step := int32(0); step < minLength; step++ {
			row := startRow + stepY*step
			col := startCol + stepX*step

			if row < 0 || row >= board.Height || col < 0 || col >= board.Width {
				break
			}

			if board.Rows[row].Cells[col] != playerId {
				break
			}

			positions = append(positions, &pb.Position{Row: row, Column: col})
		}
	}

	return s.removeDuplicatePositions(positions)
}

// Helper functions
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func (s *Connect4Service) isBoardFull(board *pb.GameBoard) bool {
	for _, height := range board.ColumnHeights {
		if height < board.Height {
			return false
		}
	}
	return true
}

func (s *Connect4Service) advanceToNextPlayer(game *pb.GameState) {
	if len(game.Players) == 0 {
		return
	}

	currentIdx := 0
	for i, player := range game.Players {
		if player.Id == game.CurrentPlayerId {
			currentIdx = i
			break
		}
	}

	nextIdx := (currentIdx + 1) % len(game.Players)
	game.CurrentPlayerId = game.Players[nextIdx].Id
	game.TurnNumber++
}

func (s *Connect4Service) removeDuplicatePositions(positions []*pb.Position) []*pb.Position {
	seen := make(map[string]bool)
	var result []*pb.Position

	for _, pos := range positions {
		key := fmt.Sprintf("%d,%d", pos.Row, pos.Column)
		if !seen[key] {
			seen[key] = true
			result = append(result, pos)
		}
	}

	return result
}

// Export for WASM
func main() {
	// WASM export will be handled by the main protoc-gen-go-wasmjs generator
}
