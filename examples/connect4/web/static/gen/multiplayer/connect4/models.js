import { Connect4Deserializer } from "./deserializer";
export class GameState {
    constructor() {
        this.gameId = "";
        this.players = [];
        this.currentPlayerId = "";
        this.turnNumber = 0;
        this.status = 0;
        this.winners = [];
        this.lastMoveTime = 0;
        this.moveTimeoutSeconds = 0;
    }
    /**
     * Create and deserialize an instance from raw data
     * @param data Raw data to deserialize
     * @returns Deserialized GameState instance or null if creation failed
     */
    static from(data) {
        return Connect4Deserializer.from(GameState.MESSAGE_TYPE, data);
    }
}
/**
 * Fully qualified message type for schema resolution
 */
GameState.MESSAGE_TYPE = "connect4.GameState";
export class GameConfig {
    constructor() {
        this.boardWidth = 0;
        this.boardHeight = 0;
        this.minPlayers = 0;
        this.maxPlayers = 0;
        this.connectLength = 0;
        this.allowMultipleWinners = false;
        this.moveTimeoutSeconds = 0;
    }
    /**
     * Create and deserialize an instance from raw data
     * @param data Raw data to deserialize
     * @returns Deserialized GameConfig instance or null if creation failed
     */
    static from(data) {
        return Connect4Deserializer.from(GameConfig.MESSAGE_TYPE, data);
    }
}
/**
 * Fully qualified message type for schema resolution
 */
GameConfig.MESSAGE_TYPE = "connect4.GameConfig";
export class Player {
    constructor() {
        this.id = "";
        this.name = "";
        this.color = "";
        this.isConnected = false;
        this.joinOrder = 0;
    }
    /**
     * Create and deserialize an instance from raw data
     * @param data Raw data to deserialize
     * @returns Deserialized Player instance or null if creation failed
     */
    static from(data) {
        return Connect4Deserializer.from(Player.MESSAGE_TYPE, data);
    }
}
/**
 * Fully qualified message type for schema resolution
 */
Player.MESSAGE_TYPE = "connect4.Player";
export class GameBoard {
    constructor() {
        this.width = 0;
        this.height = 0;
        /** Board representation: grid[y][x] = player_id (empty = "") */
        this.rows = [];
        this.columnHeights = [];
    }
    /**
     * Create and deserialize an instance from raw data
     * @param data Raw data to deserialize
     * @returns Deserialized GameBoard instance or null if creation failed
     */
    static from(data) {
        return Connect4Deserializer.from(GameBoard.MESSAGE_TYPE, data);
    }
}
/**
 * Fully qualified message type for schema resolution
 */
GameBoard.MESSAGE_TYPE = "connect4.GameBoard";
export class BoardRow {
    constructor() {
        this.cells = [];
    }
    /**
     * Create and deserialize an instance from raw data
     * @param data Raw data to deserialize
     * @returns Deserialized BoardRow instance or null if creation failed
     */
    static from(data) {
        return Connect4Deserializer.from(BoardRow.MESSAGE_TYPE, data);
    }
}
/**
 * Fully qualified message type for schema resolution
 */
BoardRow.MESSAGE_TYPE = "connect4.BoardRow";
export class PlayerStats {
    constructor() {
        this.piecesPlayed = 0;
        this.winningLines = 0;
        this.hasWon = false;
        this.totalMoveTime = 0;
    }
    /**
     * Create and deserialize an instance from raw data
     * @param data Raw data to deserialize
     * @returns Deserialized PlayerStats instance or null if creation failed
     */
    static from(data) {
        return Connect4Deserializer.from(PlayerStats.MESSAGE_TYPE, data);
    }
}
/**
 * Fully qualified message type for schema resolution
 */
PlayerStats.MESSAGE_TYPE = "connect4.PlayerStats";
/**
 * Request/Response messages
 */
export class GetGameRequest {
    constructor() {
        this.gameId = "";
    }
    /**
     * Create and deserialize an instance from raw data
     * @param data Raw data to deserialize
     * @returns Deserialized GetGameRequest instance or null if creation failed
     */
    static from(data) {
        return Connect4Deserializer.from(GetGameRequest.MESSAGE_TYPE, data);
    }
}
/**
 * Fully qualified message type for schema resolution
 */
GetGameRequest.MESSAGE_TYPE = "connect4.GetGameRequest";
export class DropPieceRequest {
    constructor() {
        this.gameId = "";
        this.playerId = "";
        this.column = 0;
    }
    /**
     * Create and deserialize an instance from raw data
     * @param data Raw data to deserialize
     * @returns Deserialized DropPieceRequest instance or null if creation failed
     */
    static from(data) {
        return Connect4Deserializer.from(DropPieceRequest.MESSAGE_TYPE, data);
    }
}
/**
 * Fully qualified message type for schema resolution
 */
DropPieceRequest.MESSAGE_TYPE = "connect4.DropPieceRequest";
export class DropPieceResponse {
    constructor() {
        this.success = false;
        this.errorMessage = "";
        this.patches = [];
        this.changeNumber = 0;
    }
    /**
     * Create and deserialize an instance from raw data
     * @param data Raw data to deserialize
     * @returns Deserialized DropPieceResponse instance or null if creation failed
     */
    static from(data) {
        return Connect4Deserializer.from(DropPieceResponse.MESSAGE_TYPE, data);
    }
}
/**
 * Fully qualified message type for schema resolution
 */
DropPieceResponse.MESSAGE_TYPE = "connect4.DropPieceResponse";
export class PieceDropResult {
    constructor() {
        this.finalRow = 0;
        this.finalColumn = 0;
        this.formedLine = false;
        this.winningLines = [];
    }
    /**
     * Create and deserialize an instance from raw data
     * @param data Raw data to deserialize
     * @returns Deserialized PieceDropResult instance or null if creation failed
     */
    static from(data) {
        return Connect4Deserializer.from(PieceDropResult.MESSAGE_TYPE, data);
    }
}
/**
 * Fully qualified message type for schema resolution
 */
PieceDropResult.MESSAGE_TYPE = "connect4.PieceDropResult";
export class LineInfo {
    constructor() {
        this.positions = [];
        this.direction = "";
        this.length = 0;
    }
    /**
     * Create and deserialize an instance from raw data
     * @param data Raw data to deserialize
     * @returns Deserialized LineInfo instance or null if creation failed
     */
    static from(data) {
        return Connect4Deserializer.from(LineInfo.MESSAGE_TYPE, data);
    }
}
/**
 * Fully qualified message type for schema resolution
 */
LineInfo.MESSAGE_TYPE = "connect4.LineInfo";
export class Position {
    constructor() {
        this.row = 0;
        this.column = 0;
    }
    /**
     * Create and deserialize an instance from raw data
     * @param data Raw data to deserialize
     * @returns Deserialized Position instance or null if creation failed
     */
    static from(data) {
        return Connect4Deserializer.from(Position.MESSAGE_TYPE, data);
    }
}
/**
 * Fully qualified message type for schema resolution
 */
Position.MESSAGE_TYPE = "connect4.Position";
export class JoinGameRequest {
    constructor() {
        this.gameId = "";
        this.playerName = "";
        this.preferredColor = "";
    }
    /**
     * Create and deserialize an instance from raw data
     * @param data Raw data to deserialize
     * @returns Deserialized JoinGameRequest instance or null if creation failed
     */
    static from(data) {
        return Connect4Deserializer.from(JoinGameRequest.MESSAGE_TYPE, data);
    }
}
/**
 * Fully qualified message type for schema resolution
 */
JoinGameRequest.MESSAGE_TYPE = "connect4.JoinGameRequest";
export class JoinGameResponse {
    constructor() {
        this.success = false;
        this.errorMessage = "";
        this.playerId = "";
        this.assignedColor = "";
    }
    /**
     * Create and deserialize an instance from raw data
     * @param data Raw data to deserialize
     * @returns Deserialized JoinGameResponse instance or null if creation failed
     */
    static from(data) {
        return Connect4Deserializer.from(JoinGameResponse.MESSAGE_TYPE, data);
    }
}
/**
 * Fully qualified message type for schema resolution
 */
JoinGameResponse.MESSAGE_TYPE = "connect4.JoinGameResponse";
export class CreateGameRequest {
    constructor() {
        this.gameId = "";
        this.creatorName = "";
    }
    /**
     * Create and deserialize an instance from raw data
     * @param data Raw data to deserialize
     * @returns Deserialized CreateGameRequest instance or null if creation failed
     */
    static from(data) {
        return Connect4Deserializer.from(CreateGameRequest.MESSAGE_TYPE, data);
    }
}
/**
 * Fully qualified message type for schema resolution
 */
CreateGameRequest.MESSAGE_TYPE = "connect4.CreateGameRequest";
export class CreateGameResponse {
    constructor() {
        this.success = false;
        this.errorMessage = "";
        this.playerId = "";
    }
    /**
     * Create and deserialize an instance from raw data
     * @param data Raw data to deserialize
     * @returns Deserialized CreateGameResponse instance or null if creation failed
     */
    static from(data) {
        return Connect4Deserializer.from(CreateGameResponse.MESSAGE_TYPE, data);
    }
}
/**
 * Fully qualified message type for schema resolution
 */
CreateGameResponse.MESSAGE_TYPE = "connect4.CreateGameResponse";
