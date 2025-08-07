// Generated TypeScript interfaces from proto file
// DO NOT EDIT - This file is auto-generated


export enum GameStatus {
  GAME_STATUS_UNSPECIFIED = 0,
  GAME_STATUS_WAITING_FOR_PLAYERS = 1,
  GAME_STATUS_IN_PROGRESS = 2,
  GAME_STATUS_FINISHED = 3,
  GAME_STATUS_PAUSED = 4,
}



export interface GameState {
  gameId: string;
  config?: GameConfig;
  players?: Player[];
  board?: GameBoard;
  currentPlayerId: string;
  turnNumber: number;
  status: GameStatus;
  winners: string[];
  playerStats?: Map<string, PlayerStats>;
  lastMoveTime: number;
  moveTimeoutSeconds: number;
}



export interface GameConfig {
  boardWidth: number;
  boardHeight: number;
  minPlayers: number;
  maxPlayers: number;
  connectLength: number;
  allowMultipleWinners: boolean;
  moveTimeoutSeconds: number;
}



export interface Player {
  id: string;
  name: string;
  color: string;
  isConnected: boolean;
  joinOrder: number;
}



export interface GameBoard {
  width: number;
  height: number;
  /** Board representation: grid[y][x] = player_id (empty = "") */
  rows?: BoardRow[];
  columnHeights: number[];
}



export interface BoardRow {
  cells: string[];
}



export interface PlayerStats {
  piecesPlayed: number;
  winningLines: number;
  hasWon: boolean;
  totalMoveTime: number;
}


/**
 * Request/Response messages
 */
export interface GetGameRequest {
  gameId: string;
}



export interface DropPieceRequest {
  gameId: string;
  playerId: string;
  column: number;
}



export interface DropPieceResponse {
  success: boolean;
  errorMessage: string;
  /** repeated wasmjs.v1.MessagePatch patches = 3; */
  changeNumber: number;
  result?: PieceDropResult;
}



export interface PieceDropResult {
  finalRow: number;
  finalColumn: number;
  formedLine: boolean;
  winningLines?: LineInfo[];
}



export interface LineInfo {
  positions?: Position[];
  direction: string;
  length: number;
}



export interface Position {
  row: number;
  column: number;
}



export interface JoinGameRequest {
  gameId: string;
  playerName: string;
  preferredColor: string;
}



export interface JoinGameResponse {
  success: boolean;
  errorMessage: string;
  playerId: string;
  assignedColor: string;
  gameState?: GameState;
}



export interface CreateGameRequest {
  gameId: string;
  config?: GameConfig;
  creatorName: string;
}



export interface CreateGameResponse {
  success: boolean;
  errorMessage: string;
  playerId: string;
  gameState?: GameState;
}

