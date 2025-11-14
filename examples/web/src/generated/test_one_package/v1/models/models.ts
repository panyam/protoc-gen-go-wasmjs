import { Any, FieldMask, Timestamp } from "@bufbuild/protobuf/wkt";


import { SampleRequest as SampleRequestInterface, SampleResponse as SampleResponseInterface } from "./interfaces";




/**
 * Request messages
 */
export class SampleRequest implements SampleRequestInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "test_one_package.v1.SampleRequest";
  readonly __MESSAGE_TYPE = SampleRequest.MESSAGE_TYPE;

  a: number = 0;
  b: string = "";
  entityData?: Any;
  updateMask?: FieldMask;

  
}



export class SampleResponse implements SampleResponseInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "test_one_package.v1.SampleResponse";
  readonly __MESSAGE_TYPE = SampleResponse.MESSAGE_TYPE;

  x: number = 0;
  createdAt?: Timestamp;
  updatedAt?: Timestamp;

  
}


