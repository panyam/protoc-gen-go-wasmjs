import { SampleRequest as SampleRequestInterface, SampleResponse as SampleResponseInterface } from "./interfaces";




/**
 * Request messages
 */
export class SampleRequest implements SampleRequestInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "test_multi_packages.v1.models.SampleRequest";

  a: number = 0;
  b: string = "";

  
}



export class SampleResponse implements SampleResponseInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "test_multi_packages.v1.models.SampleResponse";

  x: number = 0;

  
}


