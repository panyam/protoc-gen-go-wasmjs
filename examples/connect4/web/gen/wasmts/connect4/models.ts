import { GameState as GameStateInterface, GameConfig as GameConfigInterface, Player as PlayerInterface, GameBoard as GameBoardInterface, BoardRow as BoardRowInterface, PlayerStats as PlayerStatsInterface, GetGameRequest as GetGameRequestInterface, DropPieceRequest as DropPieceRequestInterface, DropPieceResponse as DropPieceResponseInterface, PieceDropResult as PieceDropResultInterface, LineInfo as LineInfoInterface, Position as PositionInterface, JoinGameRequest as JoinGameRequestInterface, JoinGameResponse as JoinGameResponseInterface, CreateGameRequest as CreateGameRequestInterface, CreateGameResponse as CreateGameResponseInterface } from "./interfaces";
import { Connect4Deserializer } from "./deserializer";



export class GameState implements GameStateInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "connect4.GameState";

  gameId: string = "";
  config?: GameConfig;
  players: Player[] = [];
  board?: GameBoard;
  currentPlayerId: string = "";
  turnNumber: number = 0;
  status: GameStatus = 0;
  winners: string[] = [];
  playerStats?: Map<string, PlayerStats>;
  lastMoveTime: number = 0;
  moveTimeoutSeconds: number = 0;

  /**
   * Create and deserialize an instance from raw data
   * @param data Raw data to deserialize
   * @returns Deserialized GameState instance or null if creation failed
   */
  static from(data: any) {
    return Connect4Deserializer.from<GameState>(GameState.MESSAGE_TYPE, data);
  }
}



export class GameConfig implements GameConfigInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "connect4.GameConfig";

  boardWidth: number = 0;
  boardHeight: number = 0;
  minPlayers: number = 0;
  maxPlayers: number = 0;
  connectLength: number = 0;
  allowMultipleWinners: boolean = false;
  moveTimeoutSeconds: number = 0;

  /**
   * Create and deserialize an instance from raw data
   * @param data Raw data to deserialize
   * @returns Deserialized GameConfig instance or null if creation failed
   */
  static from(data: any) {
    return Connect4Deserializer.from<GameConfig>(GameConfig.MESSAGE_TYPE, data);
  }
}



export class Player implements PlayerInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "connect4.Player";

  id: string = "";
  name: string = "";
  color: string = "";
  isConnected: boolean = false;
  joinOrder: number = 0;

  /**
   * Create and deserialize an instance from raw data
   * @param data Raw data to deserialize
   * @returns Deserialized Player instance or null if creation failed
   */
  static from(data: any) {
    return Connect4Deserializer.from<Player>(Player.MESSAGE_TYPE, data);
  }
}



export class GameBoard implements GameBoardInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "connect4.GameBoard";

  width: number = 0;
  height: number = 0;
  /** Board representation: grid[y][x] = player_id (empty = "") */
  rows: BoardRow[] = [];
  columnHeights: number[] = [];

  /**
   * Create and deserialize an instance from raw data
   * @param data Raw data to deserialize
   * @returns Deserialized GameBoard instance or null if creation failed
   */
  static from(data: any) {
    return Connect4Deserializer.from<GameBoard>(GameBoard.MESSAGE_TYPE, data);
  }
}



export class BoardRow implements BoardRowInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "connect4.BoardRow";

  cells: string[] = [];

  /**
   * Create and deserialize an instance from raw data
   * @param data Raw data to deserialize
   * @returns Deserialized BoardRow instance or null if creation failed
   */
  static from(data: any) {
    return Connect4Deserializer.from<BoardRow>(BoardRow.MESSAGE_TYPE, data);
  }
}



export class PlayerStats implements PlayerStatsInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "connect4.PlayerStats";

  piecesPlayed: number = 0;
  winningLines: number = 0;
  hasWon: boolean = false;
  totalMoveTime: number = 0;

  /**
   * Create and deserialize an instance from raw data
   * @param data Raw data to deserialize
   * @returns Deserialized PlayerStats instance or null if creation failed
   */
  static from(data: any) {
    return Connect4Deserializer.from<PlayerStats>(PlayerStats.MESSAGE_TYPE, data);
  }
}


/**
 * Request/Response messages
 */
export class GetGameRequest implements GetGameRequestInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "connect4.GetGameRequest";

  gameId: string = "";

  /**
   * Create and deserialize an instance from raw data
   * @param data Raw data to deserialize
   * @returns Deserialized GetGameRequest instance or null if creation failed
   */
  static from(data: any) {
    return Connect4Deserializer.from<GetGameRequest>(GetGameRequest.MESSAGE_TYPE, data);
  }
}



export class DropPieceRequest implements DropPieceRequestInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "connect4.DropPieceRequest";

  gameId: string = "";
  playerId: string = "";
  column: number = 0;

  /**
   * Create and deserialize an instance from raw data
   * @param data Raw data to deserialize
   * @returns Deserialized DropPieceRequest instance or null if creation failed
   */
  static from(data: any) {
    return Connect4Deserializer.from<DropPieceRequest>(DropPieceRequest.MESSAGE_TYPE, data);
  }
}



export class DropPieceResponse implements DropPieceResponseInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "connect4.DropPieceResponse";

  success: boolean = false;
  errorMessage: string = "";
  patches: MessagePatch[] = [];
  changeNumber: number = 0;
  result?: PieceDropResult;

  /**
   * Create and deserialize an instance from raw data
   * @param data Raw data to deserialize
   * @returns Deserialized DropPieceResponse instance or null if creation failed
   */
  static from(data: any) {
    return Connect4Deserializer.from<DropPieceResponse>(DropPieceResponse.MESSAGE_TYPE, data);
  }
}



export class PieceDropResult implements PieceDropResultInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "connect4.PieceDropResult";

  finalRow: number = 0;
  finalColumn: number = 0;
  formedLine: boolean = false;
  winningLines: LineInfo[] = [];

  /**
   * Create and deserialize an instance from raw data
   * @param data Raw data to deserialize
   * @returns Deserialized PieceDropResult instance or null if creation failed
   */
  static from(data: any) {
    return Connect4Deserializer.from<PieceDropResult>(PieceDropResult.MESSAGE_TYPE, data);
  }
}



export class LineInfo implements LineInfoInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "connect4.LineInfo";

  positions: Position[] = [];
  direction: string = "";
  length: number = 0;

  /**
   * Create and deserialize an instance from raw data
   * @param data Raw data to deserialize
   * @returns Deserialized LineInfo instance or null if creation failed
   */
  static from(data: any) {
    return Connect4Deserializer.from<LineInfo>(LineInfo.MESSAGE_TYPE, data);
  }
}



export class Position implements PositionInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "connect4.Position";

  row: number = 0;
  column: number = 0;

  /**
   * Create and deserialize an instance from raw data
   * @param data Raw data to deserialize
   * @returns Deserialized Position instance or null if creation failed
   */
  static from(data: any) {
    return Connect4Deserializer.from<Position>(Position.MESSAGE_TYPE, data);
  }
}



export class JoinGameRequest implements JoinGameRequestInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "connect4.JoinGameRequest";

  gameId: string = "";
  playerName: string = "";
  preferredColor: string = "";

  /**
   * Create and deserialize an instance from raw data
   * @param data Raw data to deserialize
   * @returns Deserialized JoinGameRequest instance or null if creation failed
   */
  static from(data: any) {
    return Connect4Deserializer.from<JoinGameRequest>(JoinGameRequest.MESSAGE_TYPE, data);
  }
}



export class JoinGameResponse implements JoinGameResponseInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "connect4.JoinGameResponse";

  success: boolean = false;
  errorMessage: string = "";
  playerId: string = "";
  assignedColor: string = "";
  gameState?: GameState;

  /**
   * Create and deserialize an instance from raw data
   * @param data Raw data to deserialize
   * @returns Deserialized JoinGameResponse instance or null if creation failed
   */
  static from(data: any) {
    return Connect4Deserializer.from<JoinGameResponse>(JoinGameResponse.MESSAGE_TYPE, data);
  }
}



export class CreateGameRequest implements CreateGameRequestInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "connect4.CreateGameRequest";

  gameId: string = "";
  config?: GameConfig;
  creatorName: string = "";

  /**
   * Create and deserialize an instance from raw data
   * @param data Raw data to deserialize
   * @returns Deserialized CreateGameRequest instance or null if creation failed
   */
  static from(data: any) {
    return Connect4Deserializer.from<CreateGameRequest>(CreateGameRequest.MESSAGE_TYPE, data);
  }
}



export class CreateGameResponse implements CreateGameResponseInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "connect4.CreateGameResponse";

  success: boolean = false;
  errorMessage: string = "";
  playerId: string = "";
  gameState?: GameState;

  /**
   * Create and deserialize an instance from raw data
   * @param data Raw data to deserialize
   * @returns Deserialized CreateGameResponse instance or null if creation failed
   */
  static from(data: any) {
    return Connect4Deserializer.from<CreateGameResponse>(CreateGameResponse.MESSAGE_TYPE, data);
  }
}


