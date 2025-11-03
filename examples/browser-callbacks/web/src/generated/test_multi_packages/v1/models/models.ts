import { SampleRequest as SampleRequestInterface, SampleResponse as SampleResponseInterface } from "./interfaces";
import { Test_multi_packages_v1_modelsDeserializer } from "./deserializer";


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

  /**
   * Create and deserialize an instance from raw data
   * @param data Raw data to deserialize
   * @returns Deserialized SampleRequest instance or null if creation failed
   */
  static from(data: any) {
    return Test_multi_packages_v1_modelsDeserializer.from<SampleRequest>(SampleRequest.MESSAGE_TYPE, data);
  }
}



export class SampleResponse implements SampleResponseInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "test_multi_packages.v1.models.SampleResponse";

  x: number = 0;

  /**
   * Create and deserialize an instance from raw data
   * @param data Raw data to deserialize
   * @returns Deserialized SampleResponse instance or null if creation failed
   */
  static from(data: any) {
    return Test_multi_packages_v1_modelsDeserializer.from<SampleResponse>(SampleResponse.MESSAGE_TYPE, data);
  }
}


