

import { GameState as GameStateInterface, GameConfig as GameConfigInterface, Player as PlayerInterface, GameBoard as GameBoardInterface, BoardRow as BoardRowInterface, PlayerStats as PlayerStatsInterface, GetGameRequest as GetGameRequestInterface, DropPieceRequest as DropPieceRequestInterface, DropPieceResponse as DropPieceResponseInterface, PieceDropResult as PieceDropResultInterface, LineInfo as LineInfoInterface, Position as PositionInterface, JoinGameRequest as JoinGameRequestInterface, JoinGameResponse as JoinGameResponseInterface, CreateGameRequest as CreateGameRequestInterface, CreateGameResponse as CreateGameResponseInterface, GameStatus } from "./interfaces";


import { GameState as ConcreteGameState, GameConfig as ConcreteGameConfig, Player as ConcretePlayer, GameBoard as ConcreteGameBoard, BoardRow as ConcreteBoardRow, PlayerStats as ConcretePlayerStats, GetGameRequest as ConcreteGetGameRequest, DropPieceRequest as ConcreteDropPieceRequest, DropPieceResponse as ConcreteDropPieceResponse, PieceDropResult as ConcretePieceDropResult, LineInfo as ConcreteLineInfo, Position as ConcretePosition, JoinGameRequest as ConcreteJoinGameRequest, JoinGameResponse as ConcreteJoinGameResponse, CreateGameRequest as ConcreteCreateGameRequest, CreateGameResponse as ConcreteCreateGameResponse } from "./models";



/**
 * Factory result interface for enhanced factory methods
 */
export interface FactoryResult<T> {
  instance: T;
  fullyLoaded: boolean;
}

/**
 * Enhanced factory with context-aware object construction
 */
export class Connect4Factory {


  /**
   * Enhanced factory method for GameState
   * @param parent Parent object containing this field
   * @param attributeName Field name in parent object
   * @param attributeKey Array index, map key, or union tag (for containers)
   * @param data Raw data to potentially populate from
   * @returns Factory result with instance and population status
   */
  newGameState = (
    parent?: any,
    attributeName?: string,
    attributeKey?: string | number,
    data?: any
  ): FactoryResult<GameStateInterface> => {
    const out = new ConcreteGameState();
    
    // Factory does not populate by default - let deserializer handle it
    return { instance: out, fullyLoaded: false };
  }

  /**
   * Enhanced factory method for GameConfig
   * @param parent Parent object containing this field
   * @param attributeName Field name in parent object
   * @param attributeKey Array index, map key, or union tag (for containers)
   * @param data Raw data to potentially populate from
   * @returns Factory result with instance and population status
   */
  newGameConfig = (
    parent?: any,
    attributeName?: string,
    attributeKey?: string | number,
    data?: any
  ): FactoryResult<GameConfigInterface> => {
    const out = new ConcreteGameConfig();
    
    // Factory does not populate by default - let deserializer handle it
    return { instance: out, fullyLoaded: false };
  }

  /**
   * Enhanced factory method for Player
   * @param parent Parent object containing this field
   * @param attributeName Field name in parent object
   * @param attributeKey Array index, map key, or union tag (for containers)
   * @param data Raw data to potentially populate from
   * @returns Factory result with instance and population status
   */
  newPlayer = (
    parent?: any,
    attributeName?: string,
    attributeKey?: string | number,
    data?: any
  ): FactoryResult<PlayerInterface> => {
    const out = new ConcretePlayer();
    
    // Factory does not populate by default - let deserializer handle it
    return { instance: out, fullyLoaded: false };
  }

  /**
   * Enhanced factory method for GameBoard
   * @param parent Parent object containing this field
   * @param attributeName Field name in parent object
   * @param attributeKey Array index, map key, or union tag (for containers)
   * @param data Raw data to potentially populate from
   * @returns Factory result with instance and population status
   */
  newGameBoard = (
    parent?: any,
    attributeName?: string,
    attributeKey?: string | number,
    data?: any
  ): FactoryResult<GameBoardInterface> => {
    const out = new ConcreteGameBoard();
    
    // Factory does not populate by default - let deserializer handle it
    return { instance: out, fullyLoaded: false };
  }

  /**
   * Enhanced factory method for BoardRow
   * @param parent Parent object containing this field
   * @param attributeName Field name in parent object
   * @param attributeKey Array index, map key, or union tag (for containers)
   * @param data Raw data to potentially populate from
   * @returns Factory result with instance and population status
   */
  newBoardRow = (
    parent?: any,
    attributeName?: string,
    attributeKey?: string | number,
    data?: any
  ): FactoryResult<BoardRowInterface> => {
    const out = new ConcreteBoardRow();
    
    // Factory does not populate by default - let deserializer handle it
    return { instance: out, fullyLoaded: false };
  }

  /**
   * Enhanced factory method for PlayerStats
   * @param parent Parent object containing this field
   * @param attributeName Field name in parent object
   * @param attributeKey Array index, map key, or union tag (for containers)
   * @param data Raw data to potentially populate from
   * @returns Factory result with instance and population status
   */
  newPlayerStats = (
    parent?: any,
    attributeName?: string,
    attributeKey?: string | number,
    data?: any
  ): FactoryResult<PlayerStatsInterface> => {
    const out = new ConcretePlayerStats();
    
    // Factory does not populate by default - let deserializer handle it
    return { instance: out, fullyLoaded: false };
  }

  /**
   * Enhanced factory method for GetGameRequest
   * @param parent Parent object containing this field
   * @param attributeName Field name in parent object
   * @param attributeKey Array index, map key, or union tag (for containers)
   * @param data Raw data to potentially populate from
   * @returns Factory result with instance and population status
   */
  newGetGameRequest = (
    parent?: any,
    attributeName?: string,
    attributeKey?: string | number,
    data?: any
  ): FactoryResult<GetGameRequestInterface> => {
    const out = new ConcreteGetGameRequest();
    
    // Factory does not populate by default - let deserializer handle it
    return { instance: out, fullyLoaded: false };
  }

  /**
   * Enhanced factory method for DropPieceRequest
   * @param parent Parent object containing this field
   * @param attributeName Field name in parent object
   * @param attributeKey Array index, map key, or union tag (for containers)
   * @param data Raw data to potentially populate from
   * @returns Factory result with instance and population status
   */
  newDropPieceRequest = (
    parent?: any,
    attributeName?: string,
    attributeKey?: string | number,
    data?: any
  ): FactoryResult<DropPieceRequestInterface> => {
    const out = new ConcreteDropPieceRequest();
    
    // Factory does not populate by default - let deserializer handle it
    return { instance: out, fullyLoaded: false };
  }

  /**
   * Enhanced factory method for DropPieceResponse
   * @param parent Parent object containing this field
   * @param attributeName Field name in parent object
   * @param attributeKey Array index, map key, or union tag (for containers)
   * @param data Raw data to potentially populate from
   * @returns Factory result with instance and population status
   */
  newDropPieceResponse = (
    parent?: any,
    attributeName?: string,
    attributeKey?: string | number,
    data?: any
  ): FactoryResult<DropPieceResponseInterface> => {
    const out = new ConcreteDropPieceResponse();
    
    // Factory does not populate by default - let deserializer handle it
    return { instance: out, fullyLoaded: false };
  }

  /**
   * Enhanced factory method for PieceDropResult
   * @param parent Parent object containing this field
   * @param attributeName Field name in parent object
   * @param attributeKey Array index, map key, or union tag (for containers)
   * @param data Raw data to potentially populate from
   * @returns Factory result with instance and population status
   */
  newPieceDropResult = (
    parent?: any,
    attributeName?: string,
    attributeKey?: string | number,
    data?: any
  ): FactoryResult<PieceDropResultInterface> => {
    const out = new ConcretePieceDropResult();
    
    // Factory does not populate by default - let deserializer handle it
    return { instance: out, fullyLoaded: false };
  }

  /**
   * Enhanced factory method for LineInfo
   * @param parent Parent object containing this field
   * @param attributeName Field name in parent object
   * @param attributeKey Array index, map key, or union tag (for containers)
   * @param data Raw data to potentially populate from
   * @returns Factory result with instance and population status
   */
  newLineInfo = (
    parent?: any,
    attributeName?: string,
    attributeKey?: string | number,
    data?: any
  ): FactoryResult<LineInfoInterface> => {
    const out = new ConcreteLineInfo();
    
    // Factory does not populate by default - let deserializer handle it
    return { instance: out, fullyLoaded: false };
  }

  /**
   * Enhanced factory method for Position
   * @param parent Parent object containing this field
   * @param attributeName Field name in parent object
   * @param attributeKey Array index, map key, or union tag (for containers)
   * @param data Raw data to potentially populate from
   * @returns Factory result with instance and population status
   */
  newPosition = (
    parent?: any,
    attributeName?: string,
    attributeKey?: string | number,
    data?: any
  ): FactoryResult<PositionInterface> => {
    const out = new ConcretePosition();
    
    // Factory does not populate by default - let deserializer handle it
    return { instance: out, fullyLoaded: false };
  }

  /**
   * Enhanced factory method for JoinGameRequest
   * @param parent Parent object containing this field
   * @param attributeName Field name in parent object
   * @param attributeKey Array index, map key, or union tag (for containers)
   * @param data Raw data to potentially populate from
   * @returns Factory result with instance and population status
   */
  newJoinGameRequest = (
    parent?: any,
    attributeName?: string,
    attributeKey?: string | number,
    data?: any
  ): FactoryResult<JoinGameRequestInterface> => {
    const out = new ConcreteJoinGameRequest();
    
    // Factory does not populate by default - let deserializer handle it
    return { instance: out, fullyLoaded: false };
  }

  /**
   * Enhanced factory method for JoinGameResponse
   * @param parent Parent object containing this field
   * @param attributeName Field name in parent object
   * @param attributeKey Array index, map key, or union tag (for containers)
   * @param data Raw data to potentially populate from
   * @returns Factory result with instance and population status
   */
  newJoinGameResponse = (
    parent?: any,
    attributeName?: string,
    attributeKey?: string | number,
    data?: any
  ): FactoryResult<JoinGameResponseInterface> => {
    const out = new ConcreteJoinGameResponse();
    
    // Factory does not populate by default - let deserializer handle it
    return { instance: out, fullyLoaded: false };
  }

  /**
   * Enhanced factory method for CreateGameRequest
   * @param parent Parent object containing this field
   * @param attributeName Field name in parent object
   * @param attributeKey Array index, map key, or union tag (for containers)
   * @param data Raw data to potentially populate from
   * @returns Factory result with instance and population status
   */
  newCreateGameRequest = (
    parent?: any,
    attributeName?: string,
    attributeKey?: string | number,
    data?: any
  ): FactoryResult<CreateGameRequestInterface> => {
    const out = new ConcreteCreateGameRequest();
    
    // Factory does not populate by default - let deserializer handle it
    return { instance: out, fullyLoaded: false };
  }

  /**
   * Enhanced factory method for CreateGameResponse
   * @param parent Parent object containing this field
   * @param attributeName Field name in parent object
   * @param attributeKey Array index, map key, or union tag (for containers)
   * @param data Raw data to potentially populate from
   * @returns Factory result with instance and population status
   */
  newCreateGameResponse = (
    parent?: any,
    attributeName?: string,
    attributeKey?: string | number,
    data?: any
  ): FactoryResult<CreateGameResponseInterface> => {
    const out = new ConcreteCreateGameResponse();
    
    // Factory does not populate by default - let deserializer handle it
    return { instance: out, fullyLoaded: false };
  }



  /**
   * Get factory method for a fully qualified message type
   * Enables cross-package factory delegation
   */
  getFactoryMethod(messageType: string): ((parent?: any, attributeName?: string, attributeKey?: string | number, data?: any) => FactoryResult<any>) | undefined {
    // Extract package from message type (e.g., "library.common.BaseMessage" -> "library.common")
    const parts = messageType.split('.');
    if (parts.length < 2) {
      return undefined;
    }
    
    const packageName = parts.slice(0, -1).join('.');
    const typeName = parts[parts.length - 1];
    const methodName = 'new' + typeName;
    
    // Check if this is our own package first
    const currentPackage = "connect4";
    if (packageName === currentPackage) {
      return (this as any)[methodName];
    }
    
    // Check external type factory mappings
    const externalFactory = this.externalTypeFactories()[messageType];
    if (externalFactory) {
      return externalFactory;
    }
    
    // Delegate to appropriate dependency factory

    
    return undefined;
  }



  /**
   * Generic object deserializer that respects factory decisions
   */
  protected deserializeObject(instance: any, data: any): any {
    if (!data || typeof data !== 'object') return instance;
    
    for (const [key, value] of Object.entries(data)) {
      if (value !== null && value !== undefined) {
        instance[key] = value;
      }
    }
    return instance;
  }

  // External type conversion methods

  /**
   * Mapping of external types to their factory methods
   */
  private externalTypeFactories(): Record<string, (parent?: any, attributeName?: string, attributeKey?: string | number, data?: any) => FactoryResult<any>> { 
      return {
          "google.protobuf.Timestamp": this.newTimestamp,
          "google.protobuf.FieldMask": this.newFieldMask,
      }
  };

  /**
   * Convert native Date to protobuf Timestamp format for serialization
   */
  serializeTimestamp(date: Date): any {
    if (!date) return null;
    return {
      seconds: Math.floor(date.getTime() / 1000).toString(),
      nanos: (date.getTime() % 1000) * 1000000
    };
  }

  /**
   * Factory method for converting protobuf Timestamp data to native Date
   * @param parent Parent object containing this field
   * @param attributeName Field name in parent object  
   * @param attributeKey Array index, map key, or union tag (for containers)
   * @param data Raw protobuf timestamp data
   * @returns Factory result with Date instance
   */
  newTimestamp = (
    parent?: any,
    attributeName?: string,
    attributeKey?: string | number,
    data?: any
  ): FactoryResult<Date> => {
    if (!data) {
      return { instance: new Date(), fullyLoaded: true };
    }
    
    let date: Date;
    if (typeof data === 'string') {
      // Handle ISO string format
      date = new Date(data);
    } else if (data.seconds !== undefined) {
      // Handle protobuf format with seconds/nanos
      const seconds = typeof data.seconds === 'string' 
        ? parseInt(data.seconds, 10) 
        : data.seconds;
      const nanos = data.nanos || 0;
      date = new Date(seconds * 1000 + Math.floor(nanos / 1000000));
    } else {
      date = new Date();
    }
    
    return { instance: date, fullyLoaded: true };
  }

  /**
   * Convert native string array to protobuf FieldMask format for serialization
   */
  serializeFieldMask(paths: string[]): any {
    if (!paths || !Array.isArray(paths)) return null;
    return { paths };
  }

  /**
   * Factory method for converting protobuf FieldMask data to native string array
   * @param parent Parent object containing this field
   * @param attributeName Field name in parent object
   * @param attributeKey Array index, map key, or union tag (for containers)
   * @param data Raw protobuf field mask data
   * @returns Factory result with string array instance
   */
  newFieldMask = (
    parent?: any,
    attributeName?: string,
    attributeKey?: string | number,
    data?: any
  ): FactoryResult<string[]> => {
    if (!data) {
      return { instance: [], fullyLoaded: true };
    }
    
    let paths: string[];
    if (Array.isArray(data)) {
      paths = data;
    } else if (data.paths && Array.isArray(data.paths)) {
      paths = data.paths;
    } else {
      paths = [];
    }
    
    return { instance: paths, fullyLoaded: true };
  }
}
