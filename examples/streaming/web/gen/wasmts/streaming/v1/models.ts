import { TickRequest as TickRequestInterface, TickResponse as TickResponseInterface } from "./interfaces";
import { StreamingV1Deserializer } from "./deserializer";


/**
 * Test message for streaming
 */
export class TickRequest implements TickRequestInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "streaming.v1.TickRequest";

  count: number = 0;
  intervalMs: number = 0;
  message: string = "";

  /**
   * Create and deserialize an instance from raw data
   * @param data Raw data to deserialize
   * @returns Deserialized TickRequest instance or null if creation failed
   */
  static from(data: any) {
    return StreamingV1Deserializer.from<TickRequest>(TickRequest.MESSAGE_TYPE, data);
  }
}



export class TickResponse implements TickResponseInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "streaming.v1.TickResponse";

  tickNumber: number = 0;
  timestamp: number = 0;
  message: string = "";
  isFinal: boolean = false;

  /**
   * Create and deserialize an instance from raw data
   * @param data Raw data to deserialize
   * @returns Deserialized TickResponse instance or null if creation failed
   */
  static from(data: any) {
    return StreamingV1Deserializer.from<TickResponse>(TickResponse.MESSAGE_TYPE, data);
  }
}


