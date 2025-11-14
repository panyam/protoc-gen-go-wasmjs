import { HelperUtilType, ParentUtilMessage_NestedUtilType } from "../../utils/v1/interfaces";
import { FieldMask, Timestamp } from "@bufbuild/protobuf/wkt";


import { LoadUserDataRequest as LoadUserDataRequestInterface, LoadUserDataResponse as LoadUserDataResponseInterface, StateUpdateRequest as StateUpdateRequestInterface, UIUpdate as UIUpdateInterface, TestRecord as TestRecordInterface, PreferencesRequest as PreferencesRequestInterface, PreferencesResponse as PreferencesResponseInterface, CallbackDemoRequest as CallbackDemoRequestInterface, CallbackDemoResponse as CallbackDemoResponseInterface } from "./interfaces";




/**
 * Request to load user data
 */
export class LoadUserDataRequest implements LoadUserDataRequestInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "presenter.v1.LoadUserDataRequest";
  readonly __MESSAGE_TYPE = LoadUserDataRequest.MESSAGE_TYPE;

  userId: string = "";

  
}


/**
 * Response with user data
 */
export class LoadUserDataResponse implements LoadUserDataResponseInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "presenter.v1.LoadUserDataResponse";
  readonly __MESSAGE_TYPE = LoadUserDataResponse.MESSAGE_TYPE;

  username: string = "";
  email: string = "";
  permissions: string[] = [];
  fromCache: boolean = false;
  createdAt?: Timestamp;

  
}


/**
 * Request to update state
 */
export class StateUpdateRequest implements StateUpdateRequestInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "presenter.v1.StateUpdateRequest";
  readonly __MESSAGE_TYPE = StateUpdateRequest.MESSAGE_TYPE;

  action: string = "";
  params: Record<string, string> = {};
  updateMask?: FieldMask;

  
}


/**
 * UI update message (streamed)
 */
export class UIUpdate implements UIUpdateInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "presenter.v1.UIUpdate";
  readonly __MESSAGE_TYPE = UIUpdate.MESSAGE_TYPE;

  component: string = "";
  action: string = "";
  data: Record<string, string> = {};

  
}



export class TestRecord implements TestRecordInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "presenter.v1.TestRecord";
  readonly __MESSAGE_TYPE = TestRecord.MESSAGE_TYPE;

  helperRecord?: HelperUtilType;
  nestedHelper?: ParentUtilMessage_NestedUtilType;

  
}


/**
 * Request to save preferences
 */
export class PreferencesRequest implements PreferencesRequestInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "presenter.v1.PreferencesRequest";
  readonly __MESSAGE_TYPE = PreferencesRequest.MESSAGE_TYPE;

  preferences: Record<string, string> = {};

  
}


/**
 * Response from preferences save
 */
export class PreferencesResponse implements PreferencesResponseInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "presenter.v1.PreferencesResponse";
  readonly __MESSAGE_TYPE = PreferencesResponse.MESSAGE_TYPE;

  saved: boolean = false;
  itemsSaved: number = 0;

  
}


/**
 * Request to run callback demo
 */
export class CallbackDemoRequest implements CallbackDemoRequestInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "presenter.v1.CallbackDemoRequest";
  readonly __MESSAGE_TYPE = CallbackDemoRequest.MESSAGE_TYPE;

  demoName: string = "";

  
}


/**
 * Response from callback demo
 */
export class CallbackDemoResponse implements CallbackDemoResponseInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "presenter.v1.CallbackDemoResponse";
  readonly __MESSAGE_TYPE = CallbackDemoResponse.MESSAGE_TYPE;

  collectedInputs: string[] = [];
  completed: boolean = false;

  
}


