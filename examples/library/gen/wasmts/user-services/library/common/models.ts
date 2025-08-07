import { BaseMessage as BaseMessageInterface, Metadata as MetadataInterface, ErrorInfo as ErrorInfoInterface } from "./interfaces";
import { LibraryCommonDeserializer } from "./deserializer";


/**
 * BaseMessage provides common fields for all library messages
 */
export class BaseMessage implements BaseMessageInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "library.common.BaseMessage";

  id: string = "";
  timestamp: number = 0;
  version: string = "";

  /**
   * Create and deserialize an instance from raw data
   * @param data Raw data to deserialize
   * @returns Deserialized BaseMessage instance or null if creation failed
   */
  static from(data: any) {
    return LibraryCommonDeserializer.from<BaseMessage>(BaseMessage.MESSAGE_TYPE, data);
  }
}


/**
 * Metadata provides additional context for requests/responses
 */
export class Metadata implements MetadataInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "library.common.Metadata";

  requestId: string = "";
  userAgent: string = "";
  headers?: Map<string, string>;

  /**
   * Create and deserialize an instance from raw data
   * @param data Raw data to deserialize
   * @returns Deserialized Metadata instance or null if creation failed
   */
  static from(data: any) {
    return LibraryCommonDeserializer.from<Metadata>(Metadata.MESSAGE_TYPE, data);
  }
}


/**
 * ErrorInfo provides structured error information
 */
export class ErrorInfo implements ErrorInfoInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "library.common.ErrorInfo";

  code: string = "";
  message: string = "";
  details: string[] = [];

  /**
   * Create and deserialize an instance from raw data
   * @param data Raw data to deserialize
   * @returns Deserialized ErrorInfo instance or null if creation failed
   */
  static from(data: any) {
    return LibraryCommonDeserializer.from<ErrorInfo>(ErrorInfo.MESSAGE_TYPE, data);
  }
}


