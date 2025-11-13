import { FetchRequest as FetchRequestInterface, FetchResponse as FetchResponseInterface, StorageKeyRequest as StorageKeyRequestInterface, StorageValueResponse as StorageValueResponseInterface, StorageSetRequest as StorageSetRequestInterface, StorageSetResponse as StorageSetResponseInterface, CookieRequest as CookieRequestInterface, CookieResponse as CookieResponseInterface, AlertRequest as AlertRequestInterface, AlertResponse as AlertResponseInterface, PromptRequest as PromptRequestInterface, PromptResponse as PromptResponseInterface, LogRequest as LogRequestInterface, LogResponse as LogResponseInterface } from "./interfaces";




/**
 * Request to fetch data from a URL
 */
export class FetchRequest implements FetchRequestInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "browser.v1.FetchRequest";

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

  logged: boolean = false;

  
}


