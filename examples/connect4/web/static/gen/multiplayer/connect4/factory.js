import { GameState as ConcreteGameState, GameConfig as ConcreteGameConfig, Player as ConcretePlayer, GameBoard as ConcreteGameBoard, BoardRow as ConcreteBoardRow, PlayerStats as ConcretePlayerStats, GetGameRequest as ConcreteGetGameRequest, DropPieceRequest as ConcreteDropPieceRequest, DropPieceResponse as ConcreteDropPieceResponse, PieceDropResult as ConcretePieceDropResult, LineInfo as ConcreteLineInfo, Position as ConcretePosition, JoinGameRequest as ConcreteJoinGameRequest, JoinGameResponse as ConcreteJoinGameResponse, CreateGameRequest as ConcreteCreateGameRequest, CreateGameResponse as ConcreteCreateGameResponse } from "./models";
import { WasmjsV1Factory } from "../wasmjs/v1/factory";
/**
 * Enhanced factory with context-aware object construction
 */
export class Connect4Factory {
    constructor() {
        // Dependency factory for wasmjs.v1 package
        this.v1Factory = new WasmjsV1Factory();
        /**
         * Enhanced factory method for GameState
         * @param parent Parent object containing this field
         * @param attributeName Field name in parent object
         * @param attributeKey Array index, map key, or union tag (for containers)
         * @param data Raw data to potentially populate from
         * @returns Factory result with instance and population status
         */
        this.newGameState = (parent, attributeName, attributeKey, data) => {
            const out = new ConcreteGameState();
            // Factory does not populate by default - let deserializer handle it
            return { instance: out, fullyLoaded: false };
        };
        /**
         * Enhanced factory method for GameConfig
         * @param parent Parent object containing this field
         * @param attributeName Field name in parent object
         * @param attributeKey Array index, map key, or union tag (for containers)
         * @param data Raw data to potentially populate from
         * @returns Factory result with instance and population status
         */
        this.newGameConfig = (parent, attributeName, attributeKey, data) => {
            const out = new ConcreteGameConfig();
            // Factory does not populate by default - let deserializer handle it
            return { instance: out, fullyLoaded: false };
        };
        /**
         * Enhanced factory method for Player
         * @param parent Parent object containing this field
         * @param attributeName Field name in parent object
         * @param attributeKey Array index, map key, or union tag (for containers)
         * @param data Raw data to potentially populate from
         * @returns Factory result with instance and population status
         */
        this.newPlayer = (parent, attributeName, attributeKey, data) => {
            const out = new ConcretePlayer();
            // Factory does not populate by default - let deserializer handle it
            return { instance: out, fullyLoaded: false };
        };
        /**
         * Enhanced factory method for GameBoard
         * @param parent Parent object containing this field
         * @param attributeName Field name in parent object
         * @param attributeKey Array index, map key, or union tag (for containers)
         * @param data Raw data to potentially populate from
         * @returns Factory result with instance and population status
         */
        this.newGameBoard = (parent, attributeName, attributeKey, data) => {
            const out = new ConcreteGameBoard();
            // Factory does not populate by default - let deserializer handle it
            return { instance: out, fullyLoaded: false };
        };
        /**
         * Enhanced factory method for BoardRow
         * @param parent Parent object containing this field
         * @param attributeName Field name in parent object
         * @param attributeKey Array index, map key, or union tag (for containers)
         * @param data Raw data to potentially populate from
         * @returns Factory result with instance and population status
         */
        this.newBoardRow = (parent, attributeName, attributeKey, data) => {
            const out = new ConcreteBoardRow();
            // Factory does not populate by default - let deserializer handle it
            return { instance: out, fullyLoaded: false };
        };
        /**
         * Enhanced factory method for PlayerStats
         * @param parent Parent object containing this field
         * @param attributeName Field name in parent object
         * @param attributeKey Array index, map key, or union tag (for containers)
         * @param data Raw data to potentially populate from
         * @returns Factory result with instance and population status
         */
        this.newPlayerStats = (parent, attributeName, attributeKey, data) => {
            const out = new ConcretePlayerStats();
            // Factory does not populate by default - let deserializer handle it
            return { instance: out, fullyLoaded: false };
        };
        /**
         * Enhanced factory method for GetGameRequest
         * @param parent Parent object containing this field
         * @param attributeName Field name in parent object
         * @param attributeKey Array index, map key, or union tag (for containers)
         * @param data Raw data to potentially populate from
         * @returns Factory result with instance and population status
         */
        this.newGetGameRequest = (parent, attributeName, attributeKey, data) => {
            const out = new ConcreteGetGameRequest();
            // Factory does not populate by default - let deserializer handle it
            return { instance: out, fullyLoaded: false };
        };
        /**
         * Enhanced factory method for DropPieceRequest
         * @param parent Parent object containing this field
         * @param attributeName Field name in parent object
         * @param attributeKey Array index, map key, or union tag (for containers)
         * @param data Raw data to potentially populate from
         * @returns Factory result with instance and population status
         */
        this.newDropPieceRequest = (parent, attributeName, attributeKey, data) => {
            const out = new ConcreteDropPieceRequest();
            // Factory does not populate by default - let deserializer handle it
            return { instance: out, fullyLoaded: false };
        };
        /**
         * Enhanced factory method for DropPieceResponse
         * @param parent Parent object containing this field
         * @param attributeName Field name in parent object
         * @param attributeKey Array index, map key, or union tag (for containers)
         * @param data Raw data to potentially populate from
         * @returns Factory result with instance and population status
         */
        this.newDropPieceResponse = (parent, attributeName, attributeKey, data) => {
            const out = new ConcreteDropPieceResponse();
            // Factory does not populate by default - let deserializer handle it
            return { instance: out, fullyLoaded: false };
        };
        /**
         * Enhanced factory method for PieceDropResult
         * @param parent Parent object containing this field
         * @param attributeName Field name in parent object
         * @param attributeKey Array index, map key, or union tag (for containers)
         * @param data Raw data to potentially populate from
         * @returns Factory result with instance and population status
         */
        this.newPieceDropResult = (parent, attributeName, attributeKey, data) => {
            const out = new ConcretePieceDropResult();
            // Factory does not populate by default - let deserializer handle it
            return { instance: out, fullyLoaded: false };
        };
        /**
         * Enhanced factory method for LineInfo
         * @param parent Parent object containing this field
         * @param attributeName Field name in parent object
         * @param attributeKey Array index, map key, or union tag (for containers)
         * @param data Raw data to potentially populate from
         * @returns Factory result with instance and population status
         */
        this.newLineInfo = (parent, attributeName, attributeKey, data) => {
            const out = new ConcreteLineInfo();
            // Factory does not populate by default - let deserializer handle it
            return { instance: out, fullyLoaded: false };
        };
        /**
         * Enhanced factory method for Position
         * @param parent Parent object containing this field
         * @param attributeName Field name in parent object
         * @param attributeKey Array index, map key, or union tag (for containers)
         * @param data Raw data to potentially populate from
         * @returns Factory result with instance and population status
         */
        this.newPosition = (parent, attributeName, attributeKey, data) => {
            const out = new ConcretePosition();
            // Factory does not populate by default - let deserializer handle it
            return { instance: out, fullyLoaded: false };
        };
        /**
         * Enhanced factory method for JoinGameRequest
         * @param parent Parent object containing this field
         * @param attributeName Field name in parent object
         * @param attributeKey Array index, map key, or union tag (for containers)
         * @param data Raw data to potentially populate from
         * @returns Factory result with instance and population status
         */
        this.newJoinGameRequest = (parent, attributeName, attributeKey, data) => {
            const out = new ConcreteJoinGameRequest();
            // Factory does not populate by default - let deserializer handle it
            return { instance: out, fullyLoaded: false };
        };
        /**
         * Enhanced factory method for JoinGameResponse
         * @param parent Parent object containing this field
         * @param attributeName Field name in parent object
         * @param attributeKey Array index, map key, or union tag (for containers)
         * @param data Raw data to potentially populate from
         * @returns Factory result with instance and population status
         */
        this.newJoinGameResponse = (parent, attributeName, attributeKey, data) => {
            const out = new ConcreteJoinGameResponse();
            // Factory does not populate by default - let deserializer handle it
            return { instance: out, fullyLoaded: false };
        };
        /**
         * Enhanced factory method for CreateGameRequest
         * @param parent Parent object containing this field
         * @param attributeName Field name in parent object
         * @param attributeKey Array index, map key, or union tag (for containers)
         * @param data Raw data to potentially populate from
         * @returns Factory result with instance and population status
         */
        this.newCreateGameRequest = (parent, attributeName, attributeKey, data) => {
            const out = new ConcreteCreateGameRequest();
            // Factory does not populate by default - let deserializer handle it
            return { instance: out, fullyLoaded: false };
        };
        /**
         * Enhanced factory method for CreateGameResponse
         * @param parent Parent object containing this field
         * @param attributeName Field name in parent object
         * @param attributeKey Array index, map key, or union tag (for containers)
         * @param data Raw data to potentially populate from
         * @returns Factory result with instance and population status
         */
        this.newCreateGameResponse = (parent, attributeName, attributeKey, data) => {
            const out = new ConcreteCreateGameResponse();
            // Factory does not populate by default - let deserializer handle it
            return { instance: out, fullyLoaded: false };
        };
        /**
         * Factory method for converting protobuf Timestamp data to native Date
         * @param parent Parent object containing this field
         * @param attributeName Field name in parent object
         * @param attributeKey Array index, map key, or union tag (for containers)
         * @param data Raw protobuf timestamp data
         * @returns Factory result with Date instance
         */
        this.newTimestamp = (parent, attributeName, attributeKey, data) => {
            if (!data) {
                return { instance: new Date(), fullyLoaded: true };
            }
            let date;
            if (typeof data === 'string') {
                // Handle ISO string format
                date = new Date(data);
            }
            else if (data.seconds !== undefined) {
                // Handle protobuf format with seconds/nanos
                const seconds = typeof data.seconds === 'string'
                    ? parseInt(data.seconds, 10)
                    : data.seconds;
                const nanos = data.nanos || 0;
                date = new Date(seconds * 1000 + Math.floor(nanos / 1000000));
            }
            else {
                date = new Date();
            }
            return { instance: date, fullyLoaded: true };
        };
        /**
         * Factory method for converting protobuf FieldMask data to native string array
         * @param parent Parent object containing this field
         * @param attributeName Field name in parent object
         * @param attributeKey Array index, map key, or union tag (for containers)
         * @param data Raw protobuf field mask data
         * @returns Factory result with string array instance
         */
        this.newFieldMask = (parent, attributeName, attributeKey, data) => {
            if (!data) {
                return { instance: [], fullyLoaded: true };
            }
            let paths;
            if (Array.isArray(data)) {
                paths = data;
            }
            else if (data.paths && Array.isArray(data.paths)) {
                paths = data.paths;
            }
            else {
                paths = [];
            }
            return { instance: paths, fullyLoaded: true };
        };
    }
    /**
     * Get factory method for a fully qualified message type
     * Enables cross-package factory delegation
     */
    getFactoryMethod(messageType) {
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
            return this[methodName];
        }
        // Check external type factory mappings
        const externalFactory = this.externalTypeFactories()[messageType];
        if (externalFactory) {
            return externalFactory;
        }
        // Delegate to appropriate dependency factory
        if (packageName === "wasmjs.v1") {
            return this.v1Factory[methodName];
        }
        return undefined;
    }
    /**
     * Generic object deserializer that respects factory decisions
     */
    deserializeObject(instance, data) {
        if (!data || typeof data !== 'object')
            return instance;
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
    externalTypeFactories() {
        return {
            "google.protobuf.Timestamp": this.newTimestamp,
            "google.protobuf.FieldMask": this.newFieldMask,
        };
    }
    ;
    /**
     * Convert native Date to protobuf Timestamp format for serialization
     */
    serializeTimestamp(date) {
        if (!date)
            return null;
        return {
            seconds: Math.floor(date.getTime() / 1000).toString(),
            nanos: (date.getTime() % 1000) * 1000000
        };
    }
    /**
     * Convert native string array to protobuf FieldMask format for serialization
     */
    serializeFieldMask(paths) {
        if (!paths || !Array.isArray(paths))
            return null;
        return { paths };
    }
}
