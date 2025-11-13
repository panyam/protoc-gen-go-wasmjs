import { SecondRequest as SecondRequestInterface, SecondResponse as SecondResponseInterface } from "./interfaces";




/**
 * Request messages
 */
export class SecondRequest implements SecondRequestInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "test_one_package.v1.SecondRequest";

  a: number = 0;
  b: string = "";

  
}



export class SecondResponse implements SecondResponseInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "test_one_package.v1.SecondResponse";

  x: number = 0;

  
}


