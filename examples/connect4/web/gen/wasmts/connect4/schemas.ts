
// Generated TypeScript schemas from proto file
// DO NOT EDIT - This file is auto-generated

import { FieldType, FieldSchema, MessageSchema } from "./deserializer_schemas";


/**
 * Schema for GameState message
 */
export const GameStateSchema: MessageSchema = {
  name: "GameState",
  fields: [
    {
      name: "gameId",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "config",
      type: FieldType.MESSAGE,
      id: 2,
      messageType: "connect4.GameConfig",
    },
    {
      name: "players",
      type: FieldType.MESSAGE,
      id: 3,
      messageType: "connect4.Player",
      repeated: true,
    },
    {
      name: "board",
      type: FieldType.MESSAGE,
      id: 4,
      messageType: "connect4.GameBoard",
    },
    {
      name: "currentPlayerId",
      type: FieldType.STRING,
      id: 5,
    },
    {
      name: "turnNumber",
      type: FieldType.NUMBER,
      id: 6,
    },
    {
      name: "status",
      type: FieldType.STRING,
      id: 7,
    },
    {
      name: "winners",
      type: FieldType.REPEATED,
      id: 8,
      repeated: true,
    },
    {
      name: "playerStats",
      type: FieldType.MESSAGE,
      id: 9,
      messageType: "connect4.PlayerStatsEntry",
    },
    {
      name: "lastMoveTime",
      type: FieldType.NUMBER,
      id: 10,
    },
    {
      name: "moveTimeoutSeconds",
      type: FieldType.NUMBER,
      id: 11,
    },
  ],
};


/**
 * Schema for GameConfig message
 */
export const GameConfigSchema: MessageSchema = {
  name: "GameConfig",
  fields: [
    {
      name: "boardWidth",
      type: FieldType.NUMBER,
      id: 1,
    },
    {
      name: "boardHeight",
      type: FieldType.NUMBER,
      id: 2,
    },
    {
      name: "minPlayers",
      type: FieldType.NUMBER,
      id: 3,
    },
    {
      name: "maxPlayers",
      type: FieldType.NUMBER,
      id: 4,
    },
    {
      name: "connectLength",
      type: FieldType.NUMBER,
      id: 5,
    },
    {
      name: "allowMultipleWinners",
      type: FieldType.BOOLEAN,
      id: 6,
    },
    {
      name: "moveTimeoutSeconds",
      type: FieldType.NUMBER,
      id: 7,
    },
  ],
};


/**
 * Schema for Player message
 */
export const PlayerSchema: MessageSchema = {
  name: "Player",
  fields: [
    {
      name: "id",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "name",
      type: FieldType.STRING,
      id: 2,
    },
    {
      name: "color",
      type: FieldType.STRING,
      id: 3,
    },
    {
      name: "isConnected",
      type: FieldType.BOOLEAN,
      id: 4,
    },
    {
      name: "joinOrder",
      type: FieldType.NUMBER,
      id: 5,
    },
  ],
};


/**
 * Schema for GameBoard message
 */
export const GameBoardSchema: MessageSchema = {
  name: "GameBoard",
  fields: [
    {
      name: "width",
      type: FieldType.NUMBER,
      id: 1,
    },
    {
      name: "height",
      type: FieldType.NUMBER,
      id: 2,
    },
    {
      name: "rows",
      type: FieldType.MESSAGE,
      id: 3,
      messageType: "connect4.BoardRow",
      repeated: true,
    },
    {
      name: "columnHeights",
      type: FieldType.REPEATED,
      id: 4,
      repeated: true,
    },
  ],
};


/**
 * Schema for BoardRow message
 */
export const BoardRowSchema: MessageSchema = {
  name: "BoardRow",
  fields: [
    {
      name: "cells",
      type: FieldType.REPEATED,
      id: 1,
      repeated: true,
    },
  ],
};


/**
 * Schema for PlayerStats message
 */
export const PlayerStatsSchema: MessageSchema = {
  name: "PlayerStats",
  fields: [
    {
      name: "piecesPlayed",
      type: FieldType.NUMBER,
      id: 1,
    },
    {
      name: "winningLines",
      type: FieldType.NUMBER,
      id: 2,
    },
    {
      name: "hasWon",
      type: FieldType.BOOLEAN,
      id: 3,
    },
    {
      name: "totalMoveTime",
      type: FieldType.NUMBER,
      id: 4,
    },
  ],
};


/**
 * Schema for GetGameRequest message
 */
export const GetGameRequestSchema: MessageSchema = {
  name: "GetGameRequest",
  fields: [
    {
      name: "gameId",
      type: FieldType.STRING,
      id: 1,
    },
  ],
};


/**
 * Schema for DropPieceRequest message
 */
export const DropPieceRequestSchema: MessageSchema = {
  name: "DropPieceRequest",
  fields: [
    {
      name: "gameId",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "playerId",
      type: FieldType.STRING,
      id: 2,
    },
    {
      name: "column",
      type: FieldType.NUMBER,
      id: 3,
    },
  ],
};


/**
 * Schema for DropPieceResponse message
 */
export const DropPieceResponseSchema: MessageSchema = {
  name: "DropPieceResponse",
  fields: [
    {
      name: "success",
      type: FieldType.BOOLEAN,
      id: 1,
    },
    {
      name: "errorMessage",
      type: FieldType.STRING,
      id: 2,
    },
    {
      name: "changeNumber",
      type: FieldType.NUMBER,
      id: 4,
    },
    {
      name: "result",
      type: FieldType.MESSAGE,
      id: 5,
      messageType: "connect4.PieceDropResult",
    },
  ],
};


/**
 * Schema for PieceDropResult message
 */
export const PieceDropResultSchema: MessageSchema = {
  name: "PieceDropResult",
  fields: [
    {
      name: "finalRow",
      type: FieldType.NUMBER,
      id: 1,
    },
    {
      name: "finalColumn",
      type: FieldType.NUMBER,
      id: 2,
    },
    {
      name: "formedLine",
      type: FieldType.BOOLEAN,
      id: 3,
    },
    {
      name: "winningLines",
      type: FieldType.MESSAGE,
      id: 4,
      messageType: "connect4.LineInfo",
      repeated: true,
    },
  ],
};


/**
 * Schema for LineInfo message
 */
export const LineInfoSchema: MessageSchema = {
  name: "LineInfo",
  fields: [
    {
      name: "positions",
      type: FieldType.MESSAGE,
      id: 1,
      messageType: "connect4.Position",
      repeated: true,
    },
    {
      name: "direction",
      type: FieldType.STRING,
      id: 2,
    },
    {
      name: "length",
      type: FieldType.NUMBER,
      id: 3,
    },
  ],
};


/**
 * Schema for Position message
 */
export const PositionSchema: MessageSchema = {
  name: "Position",
  fields: [
    {
      name: "row",
      type: FieldType.NUMBER,
      id: 1,
    },
    {
      name: "column",
      type: FieldType.NUMBER,
      id: 2,
    },
  ],
};


/**
 * Schema for JoinGameRequest message
 */
export const JoinGameRequestSchema: MessageSchema = {
  name: "JoinGameRequest",
  fields: [
    {
      name: "gameId",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "playerName",
      type: FieldType.STRING,
      id: 2,
    },
    {
      name: "preferredColor",
      type: FieldType.STRING,
      id: 3,
    },
  ],
};


/**
 * Schema for JoinGameResponse message
 */
export const JoinGameResponseSchema: MessageSchema = {
  name: "JoinGameResponse",
  fields: [
    {
      name: "success",
      type: FieldType.BOOLEAN,
      id: 1,
    },
    {
      name: "errorMessage",
      type: FieldType.STRING,
      id: 2,
    },
    {
      name: "playerId",
      type: FieldType.STRING,
      id: 3,
    },
    {
      name: "assignedColor",
      type: FieldType.STRING,
      id: 4,
    },
    {
      name: "gameState",
      type: FieldType.MESSAGE,
      id: 5,
      messageType: "connect4.GameState",
    },
  ],
};


/**
 * Schema for CreateGameRequest message
 */
export const CreateGameRequestSchema: MessageSchema = {
  name: "CreateGameRequest",
  fields: [
    {
      name: "gameId",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "config",
      type: FieldType.MESSAGE,
      id: 2,
      messageType: "connect4.GameConfig",
    },
    {
      name: "creatorName",
      type: FieldType.STRING,
      id: 3,
    },
  ],
};


/**
 * Schema for CreateGameResponse message
 */
export const CreateGameResponseSchema: MessageSchema = {
  name: "CreateGameResponse",
  fields: [
    {
      name: "success",
      type: FieldType.BOOLEAN,
      id: 1,
    },
    {
      name: "errorMessage",
      type: FieldType.STRING,
      id: 2,
    },
    {
      name: "playerId",
      type: FieldType.STRING,
      id: 3,
    },
    {
      name: "gameState",
      type: FieldType.MESSAGE,
      id: 4,
      messageType: "connect4.GameState",
    },
  ],
};



/**
 * Package-scoped schema registry for connect4
 */
export const Connect4SchemaRegistry: Record<string, MessageSchema> = {
  "connect4.GameState": GameStateSchema,
  "connect4.GameConfig": GameConfigSchema,
  "connect4.Player": PlayerSchema,
  "connect4.GameBoard": GameBoardSchema,
  "connect4.BoardRow": BoardRowSchema,
  "connect4.PlayerStats": PlayerStatsSchema,
  "connect4.GetGameRequest": GetGameRequestSchema,
  "connect4.DropPieceRequest": DropPieceRequestSchema,
  "connect4.DropPieceResponse": DropPieceResponseSchema,
  "connect4.PieceDropResult": PieceDropResultSchema,
  "connect4.LineInfo": LineInfoSchema,
  "connect4.Position": PositionSchema,
  "connect4.JoinGameRequest": JoinGameRequestSchema,
  "connect4.JoinGameResponse": JoinGameResponseSchema,
  "connect4.CreateGameRequest": CreateGameRequestSchema,
  "connect4.CreateGameResponse": CreateGameResponseSchema,
};

/**
 * Get schema for a message type from connect4 package
 */
export function getSchema(messageType: string): MessageSchema | undefined {
  return Connect4SchemaRegistry[messageType];
}

/**
 * Get field schema by name from connect4 package
 */
export function getFieldSchema(messageType: string, fieldName: string): FieldSchema | undefined {
  const schema = getSchema(messageType);
  return schema?.fields.find(field => field.name === fieldName);
}

/**
 * Get field schema by proto field ID from connect4 package
 */
export function getFieldSchemaById(messageType: string, fieldId: number): FieldSchema | undefined {
  const schema = getSchema(messageType);
  return schema?.fields.find(field => field.id === fieldId);
}

/**
 * Check if field is part of a oneof group in connect4 package
 */
export function isOneofField(messageType: string, fieldName: string): boolean {
  const fieldSchema = getFieldSchema(messageType, fieldName);
  return fieldSchema?.oneofGroup !== undefined;
}

/**
 * Get all fields in a oneof group from connect4 package
 */
export function getOneofFields(messageType: string, oneofGroup: string): FieldSchema[] {
  const schema = getSchema(messageType);
  return schema?.fields.filter(field => field.oneofGroup === oneofGroup) || [];
}