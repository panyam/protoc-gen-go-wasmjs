// Generated TypeScript interfaces from proto file
// DO NOT EDIT - This file is auto-generated

import { Any, FieldMask, Timestamp } from "@bufbuild/protobuf/wkt";



/**
 * Request messages
 */
export interface SampleRequest {
  a: number;
  b: string;
  entityData?: Any;
  updateMask?: FieldMask;
}



export interface SampleResponse {
  x: number;
  createdAt?: Timestamp;
  updatedAt?: Timestamp;
}

