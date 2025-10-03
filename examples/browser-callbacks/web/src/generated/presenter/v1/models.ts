import { HelperUtilType, ParentUtilMessage_NestedUtilType } from "../../utils/v1/interfaces";
import { FieldMask, Timestamp } from "@bufbuild/protobuf/wkt";


import { LoadUserDataRequest as LoadUserDataRequestInterface, LoadUserDataResponse as LoadUserDataResponseInterface, StateUpdateRequest as StateUpdateRequestInterface, UIUpdate as UIUpdateInterface, TestRecord as TestRecordInterface, PreferencesRequest as PreferencesRequestInterface, PreferencesResponse as PreferencesResponseInterface, CallbackDemoRequest as CallbackDemoRequestInterface, CallbackDemoResponse as CallbackDemoResponseInterface } from "./interfaces";
import { Presenter_v1Deserializer } from "./deserializer";


/**
 * Request to load user data
 */
export class LoadUserDataRequest implements LoadUserDataRequestInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "presenter.v1.LoadUserDataRequest";

  userId: string = "";

  /**
   * Create and deserialize an instance from raw data
   * @param data Raw data to deserialize
   * @returns Deserialized LoadUserDataRequest instance or null if creation failed
   */
  static from(data: any) {
    return Presenter_v1Deserializer.from<LoadUserDataRequest>(LoadUserDataRequest.MESSAGE_TYPE, data);
  }
}


/**
 * Response with user data
 */
export class LoadUserDataResponse implements LoadUserDataResponseInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "presenter.v1.LoadUserDataResponse";

  username: string = "";
  email: string = "";
  permissions: string[] = [];
  fromCache: boolean = false;
  createdAt?: Timestamp;

  /**
   * Create and deserialize an instance from raw data
   * @param data Raw data to deserialize
   * @returns Deserialized LoadUserDataResponse instance or null if creation failed
   */
  static from(data: any) {
    return Presenter_v1Deserializer.from<LoadUserDataResponse>(LoadUserDataResponse.MESSAGE_TYPE, data);
  }
}


/**
 * Request to update state
 */
export class StateUpdateRequest implements StateUpdateRequestInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "presenter.v1.StateUpdateRequest";

  action: string = "";
  params: Record<string, string> = {};
  updateMask?: FieldMask;

  /**
   * Create and deserialize an instance from raw data
   * @param data Raw data to deserialize
   * @returns Deserialized StateUpdateRequest instance or null if creation failed
   */
  static from(data: any) {
    return Presenter_v1Deserializer.from<StateUpdateRequest>(StateUpdateRequest.MESSAGE_TYPE, data);
  }
}


/**
 * UI update message (streamed)
 */
export class UIUpdate implements UIUpdateInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "presenter.v1.UIUpdate";

  component: string = "";
  action: string = "";
  data: Record<string, string> = {};

  /**
   * Create and deserialize an instance from raw data
   * @param data Raw data to deserialize
   * @returns Deserialized UIUpdate instance or null if creation failed
   */
  static from(data: any) {
    return Presenter_v1Deserializer.from<UIUpdate>(UIUpdate.MESSAGE_TYPE, data);
  }
}



export class TestRecord implements TestRecordInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "presenter.v1.TestRecord";

  helperRecord?: HelperUtilType;
  nestedHelper?: ParentUtilMessage_NestedUtilType;

  /**
   * Create and deserialize an instance from raw data
   * @param data Raw data to deserialize
   * @returns Deserialized TestRecord instance or null if creation failed
   */
  static from(data: any) {
    return Presenter_v1Deserializer.from<TestRecord>(TestRecord.MESSAGE_TYPE, data);
  }
}


/**
 * Request to save preferences
 */
export class PreferencesRequest implements PreferencesRequestInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "presenter.v1.PreferencesRequest";

  preferences: Record<string, string> = {};

  /**
   * Create and deserialize an instance from raw data
   * @param data Raw data to deserialize
   * @returns Deserialized PreferencesRequest instance or null if creation failed
   */
  static from(data: any) {
    return Presenter_v1Deserializer.from<PreferencesRequest>(PreferencesRequest.MESSAGE_TYPE, data);
  }
}


/**
 * Response from preferences save
 */
export class PreferencesResponse implements PreferencesResponseInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "presenter.v1.PreferencesResponse";

  saved: boolean = false;
  itemsSaved: number = 0;

  /**
   * Create and deserialize an instance from raw data
   * @param data Raw data to deserialize
   * @returns Deserialized PreferencesResponse instance or null if creation failed
   */
  static from(data: any) {
    return Presenter_v1Deserializer.from<PreferencesResponse>(PreferencesResponse.MESSAGE_TYPE, data);
  }
}


/**
 * Request to run callback demo
 */
export class CallbackDemoRequest implements CallbackDemoRequestInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "presenter.v1.CallbackDemoRequest";

  demoName: string = "";

  /**
   * Create and deserialize an instance from raw data
   * @param data Raw data to deserialize
   * @returns Deserialized CallbackDemoRequest instance or null if creation failed
   */
  static from(data: any) {
    return Presenter_v1Deserializer.from<CallbackDemoRequest>(CallbackDemoRequest.MESSAGE_TYPE, data);
  }
}


/**
 * Response from callback demo
 */
export class CallbackDemoResponse implements CallbackDemoResponseInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "presenter.v1.CallbackDemoResponse";

  collectedInputs: string[] = [];
  completed: boolean = false;

  /**
   * Create and deserialize an instance from raw data
   * @param data Raw data to deserialize
   * @returns Deserialized CallbackDemoResponse instance or null if creation failed
   */
  static from(data: any) {
    return Presenter_v1Deserializer.from<CallbackDemoResponse>(CallbackDemoResponse.MESSAGE_TYPE, data);
  }
}


