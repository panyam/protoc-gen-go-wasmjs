import { SecondRequest as SecondRequestInterface, SecondResponse as SecondResponseInterface } from "./interfaces";




/**
 * Request messages
 */
export class SecondRequest implements SecondRequestInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "test_one_package.v1.SecondRequest";
  readonly __MESSAGE_TYPE = SecondRequest.MESSAGE_TYPE;

  a: number = 0;
  b: string = "";

  
}



export class SecondResponse implements SecondResponseInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "test_one_package.v1.SecondResponse";
  readonly __MESSAGE_TYPE = SecondResponse.MESSAGE_TYPE;

  x: number = 0;

  
}


