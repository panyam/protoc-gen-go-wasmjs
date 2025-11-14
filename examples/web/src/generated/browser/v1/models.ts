import { FetchRequest as FetchRequestInterface, FetchResponse as FetchResponseInterface, StorageKeyRequest as StorageKeyRequestInterface, StorageValueResponse as StorageValueResponseInterface, StorageSetRequest as StorageSetRequestInterface, StorageSetResponse as StorageSetResponseInterface, CookieRequest as CookieRequestInterface, CookieResponse as CookieResponseInterface, AlertRequest as AlertRequestInterface, AlertResponse as AlertResponseInterface, PromptRequest as PromptRequestInterface, PromptResponse as PromptResponseInterface, LogRequest as LogRequestInterface, LogResponse as LogResponseInterface } from "./interfaces";




/**
 * Request to fetch data from a URL
 */
export class FetchRequest implements FetchRequestInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "browser.v1.FetchRequest";
  readonly __MESSAGE_TYPE = FetchRequest.MESSAGE_TYPE;

  url: string = "";
  method: string = "";
  headers: Record<string, string> = {};
  body: string = "";

  
}


/**
 * Response from fetch
 */
export class FetchResponse implements FetchResponseInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "browser.v1.FetchResponse";
  readonly __MESSAGE_TYPE = FetchResponse.MESSAGE_TYPE;

  status: number = 0;
  statusText: string = "";
  headers: Record<string, string> = {};
  body: string = "";

  
}


/**
 * Request for localStorage key
 */
export class StorageKeyRequest implements StorageKeyRequestInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "browser.v1.StorageKeyRequest";
  readonly __MESSAGE_TYPE = StorageKeyRequest.MESSAGE_TYPE;

  key: string = "";

  
}


/**
 * Response with localStorage value
 */
export class StorageValueResponse implements StorageValueResponseInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "browser.v1.StorageValueResponse";
  readonly __MESSAGE_TYPE = StorageValueResponse.MESSAGE_TYPE;

  value: string = "";
  exists: boolean = false;

  
}


/**
 * Request to set localStorage
 */
export class StorageSetRequest implements StorageSetRequestInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "browser.v1.StorageSetRequest";
  readonly __MESSAGE_TYPE = StorageSetRequest.MESSAGE_TYPE;

  key: string = "";
  value: string = "";

  
}


/**
 * Response from localStorage set
 */
export class StorageSetResponse implements StorageSetResponseInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "browser.v1.StorageSetResponse";
  readonly __MESSAGE_TYPE = StorageSetResponse.MESSAGE_TYPE;

  success: boolean = false;

  
}


/**
 * Request for cookie
 */
export class CookieRequest implements CookieRequestInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "browser.v1.CookieRequest";
  readonly __MESSAGE_TYPE = CookieRequest.MESSAGE_TYPE;

  name: string = "";

  
}


/**
 * Response with cookie value
 */
export class CookieResponse implements CookieResponseInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "browser.v1.CookieResponse";
  readonly __MESSAGE_TYPE = CookieResponse.MESSAGE_TYPE;

  value: string = "";
  exists: boolean = false;

  
}


/**
 * Request to show alert
 */
export class AlertRequest implements AlertRequestInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "browser.v1.AlertRequest";
  readonly __MESSAGE_TYPE = AlertRequest.MESSAGE_TYPE;

  message: string = "";

  
}


/**
 * Response from alert
 */
export class AlertResponse implements AlertResponseInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "browser.v1.AlertResponse";
  readonly __MESSAGE_TYPE = AlertResponse.MESSAGE_TYPE;

  shown: boolean = false;

  
}


/**
 * Request for user prompt
 */
export class PromptRequest implements PromptRequestInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "browser.v1.PromptRequest";
  readonly __MESSAGE_TYPE = PromptRequest.MESSAGE_TYPE;

  message: string = "";
  defaultValue: string = "";

  
}


/**
 * Response from user prompt
 */
export class PromptResponse implements PromptResponseInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "browser.v1.PromptResponse";
  readonly __MESSAGE_TYPE = PromptResponse.MESSAGE_TYPE;

  value: string = "";
  cancelled: boolean = false;

  
}


/**
 * Request to log to window
 */
export class LogRequest implements LogRequestInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "browser.v1.LogRequest";
  readonly __MESSAGE_TYPE = LogRequest.MESSAGE_TYPE;

  message: string = "";
  level: string = "";

  
}


/**
 * Response from log to window
 */
export class LogResponse implements LogResponseInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "browser.v1.LogResponse";
  readonly __MESSAGE_TYPE = LogResponse.MESSAGE_TYPE;

  logged: boolean = false;

  
}


