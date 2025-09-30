// Generated TypeScript interfaces from proto file
// DO NOT EDIT - This file is auto-generated

/**
 * Scheme describes the schemes supported by the OpenAPI Swagger
 and Operation objects.
 */
export enum Scheme {
  UNKNOWN = 0,
  HTTP = 1,
  HTTPS = 2,
  WS = 3,
  WSS = 4,
}

/**
 * `Type` is a supported HTTP header type.
 See https://swagger.io/specification/v2/#parameterType.
 */
export enum Type {
  UNKNOWN = 0,
  STRING = 1,
  NUMBER = 2,
  INTEGER = 3,
  BOOLEAN = 4,
}


export enum JSONSchemaSimpleTypes {
  UNKNOWN = 0,
  ARRAY = 1,
  BOOLEAN = 2,
  INTEGER = 3,
  NULL = 4,
  NUMBER = 5,
  OBJECT = 6,
  STRING = 7,
}

/**
 * The type of the security scheme. Valid values are "basic",
 "apiKey" or "oauth2".
 */
export enum Type {
  TYPE_INVALID = 0,
  TYPE_BASIC = 1,
  TYPE_API_KEY = 2,
  TYPE_OAUTH2 = 3,
}

/**
 * The location of the API key. Valid values are "query" or "header".
 */
export enum In {
  IN_INVALID = 0,
  IN_QUERY = 1,
  IN_HEADER = 2,
}

/**
 * The flow used by the OAuth2 security scheme. Valid values are
 "implicit", "password", "application" or "accessCode".
 */
export enum Flow {
  FLOW_INVALID = 0,
  FLOW_IMPLICIT = 1,
  FLOW_PASSWORD = 2,
  FLOW_APPLICATION = 3,
  FLOW_ACCESS_CODE = 4,
}


/**
 * `Swagger` is a representation of OpenAPI v2 specification's Swagger object.

 See: https://github.com/OAI/OpenAPI-Specification/blob/3.0.0/versions/2.0.md#swaggerObject

 Example:

  option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_swagger) = {
    info: {
      title: "Echo API";
      version: "1.0";
      description: "";
      contact: {
        name: "gRPC-Gateway project";
        url: "https://github.com/grpc-ecosystem/grpc-gateway";
        email: "none@example.com";
      };
      license: {
        name: "BSD 3-Clause License";
        url: "https://github.com/grpc-ecosystem/grpc-gateway/blob/main/LICENSE";
      };
    };
    schemes: HTTPS;
    consumes: "application/json";
    produces: "application/json";
  };
 */
export interface Swagger {
  /** Specifies the OpenAPI Specification version being used. It can be
 used by the OpenAPI UI and other clients to interpret the API listing. The
 value MUST be "2.0". */
  swagger: string;
  /** Provides metadata about the API. The metadata can be used by the
 clients if needed. */
  info?: Info;
  /** The host (name or ip) serving the API. This MUST be the host only and does
 not include the scheme nor sub-paths. It MAY include a port. If the host is
 not included, the host serving the documentation is to be used (including
 the port). The host does not support path templating. */
  host: string;
  /** The base path on which the API is served, which is relative to the host. If
 it is not included, the API is served directly under the host. The value
 MUST start with a leading slash (/). The basePath does not support path
 templating.
 Note that using `base_path` does not change the endpoint paths that are
 generated in the resulting OpenAPI file. If you wish to use `base_path`
 with relatively generated OpenAPI paths, the `base_path` prefix must be
 manually removed from your `google.api.http` paths and your code changed to
 serve the API from the `base_path`. */
  basePath: string;
  /** The transfer protocol of the API. Values MUST be from the list: "http",
 "https", "ws", "wss". If the schemes is not included, the default scheme to
 be used is the one used to access the OpenAPI definition itself. */
  schemes: Scheme[];
  /** A list of MIME types the APIs can consume. This is global to all APIs but
 can be overridden on specific API calls. Value MUST be as described under
 Mime Types. */
  consumes: string[];
  /** A list of MIME types the APIs can produce. This is global to all APIs but
 can be overridden on specific API calls. Value MUST be as described under
 Mime Types. */
  produces: string[];
  /** An object to hold responses that can be used across operations. This
 property does not define global responses for all operations. */
  responses: Record<string, Response>;
  /** Security scheme definitions that can be used across the specification. */
  securityDefinitions?: SecurityDefinitions;
  /** A declaration of which security schemes are applied for the API as a whole.
 The list of values describes alternative security schemes that can be used
 (that is, there is a logical OR between the security requirements).
 Individual operations can override this definition. */
  security?: SecurityRequirement[];
  /** A list of tags for API documentation control. Tags can be used for logical
 grouping of operations by resources or any other qualifier. */
  tags?: Tag[];
  /** Additional external documentation. */
  externalDocs?: ExternalDocumentation;
  /** Custom properties that start with "x-" such as "x-foo" used to describe
 extra functionality that is not covered by the standard OpenAPI Specification.
 See: https://swagger.io/docs/specification/2-0/swagger-extensions/ */
  extensions: Record<string, Value>;
}


/**
 * `Operation` is a representation of OpenAPI v2 specification's Operation object.

 See: https://github.com/OAI/OpenAPI-Specification/blob/3.0.0/versions/2.0.md#operationObject

 Example:

  service EchoService {
    rpc Echo(SimpleMessage) returns (SimpleMessage) {
      option (google.api.http) = {
        get: "/v1/example/echo/{id}"
      };

      option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_operation) = {
        summary: "Get a message.";
        operation_id: "getMessage";
        tags: "echo";
        responses: {
          key: "200"
            value: {
            description: "OK";
          }
        }
      };
    }
  }
 */
export interface Operation {
  /** A list of tags for API documentation control. Tags can be used for logical
 grouping of operations by resources or any other qualifier. */
  tags: string[];
  /** A short summary of what the operation does. For maximum readability in the
 swagger-ui, this field SHOULD be less than 120 characters. */
  summary: string;
  /** A verbose explanation of the operation behavior. GFM syntax can be used for
 rich text representation. */
  description: string;
  /** Additional external documentation for this operation. */
  externalDocs?: ExternalDocumentation;
  /** Unique string used to identify the operation. The id MUST be unique among
 all operations described in the API. Tools and libraries MAY use the
 operationId to uniquely identify an operation, therefore, it is recommended
 to follow common programming naming conventions. */
  operationId: string;
  /** A list of MIME types the operation can consume. This overrides the consumes
 definition at the OpenAPI Object. An empty value MAY be used to clear the
 global definition. Value MUST be as described under Mime Types. */
  consumes: string[];
  /** A list of MIME types the operation can produce. This overrides the produces
 definition at the OpenAPI Object. An empty value MAY be used to clear the
 global definition. Value MUST be as described under Mime Types. */
  produces: string[];
  /** The list of possible responses as they are returned from executing this
 operation. */
  responses: Record<string, Response>;
  /** The transfer protocol for the operation. Values MUST be from the list:
 "http", "https", "ws", "wss". The value overrides the OpenAPI Object
 schemes definition. */
  schemes: Scheme[];
  /** Declares this operation to be deprecated. Usage of the declared operation
 should be refrained. Default value is false. */
  deprecated: boolean;
  /** A declaration of which security schemes are applied for this operation. The
 list of values describes alternative security schemes that can be used
 (that is, there is a logical OR between the security requirements). This
 definition overrides any declared top-level security. To remove a top-level
 security declaration, an empty array can be used. */
  security?: SecurityRequirement[];
  /** Custom properties that start with "x-" such as "x-foo" used to describe
 extra functionality that is not covered by the standard OpenAPI Specification.
 See: https://swagger.io/docs/specification/2-0/swagger-extensions/ */
  extensions: Record<string, Value>;
  /** Custom parameters such as HTTP request headers.
 See: https://swagger.io/docs/specification/2-0/describing-parameters/
 and https://swagger.io/specification/v2/#parameter-object. */
  parameters?: Parameters;
}


/**
 * `Parameters` is a representation of OpenAPI v2 specification's parameters object.
 Note: This technically breaks compatibility with the OpenAPI 2 definition structure as we only
 allow header parameters to be set here since we do not want users specifying custom non-header
 parameters beyond those inferred from the Protobuf schema.
 See: https://swagger.io/specification/v2/#parameter-object
 */
export interface Parameters {
  /** `Headers` is one or more HTTP header parameter.
 See: https://swagger.io/docs/specification/2-0/describing-parameters/#header-parameters */
  headers?: HeaderParameter[];
}


/**
 * `HeaderParameter` a HTTP header parameter.
 See: https://swagger.io/specification/v2/#parameter-object
 */
export interface HeaderParameter {
  /** `Name` is the header name. */
  name: string;
  /** `Description` is a short description of the header. */
  description: string;
  /** `Type` is the type of the object. The value MUST be one of "string", "number", "integer", or "boolean". The "array" type is not supported.
 See: https://swagger.io/specification/v2/#parameterType. */
  type: Type;
  /** `Format` The extending format for the previously mentioned type. */
  format: string;
  /** `Required` indicates if the header is optional */
  required: boolean;
}


/**
 * `Header` is a representation of OpenAPI v2 specification's Header object.

 See: https://github.com/OAI/OpenAPI-Specification/blob/3.0.0/versions/2.0.md#headerObject
 */
export interface Header {
  /** `Description` is a short description of the header. */
  description: string;
  /** The type of the object. The value MUST be one of "string", "number", "integer", or "boolean". The "array" type is not supported. */
  type: string;
  /** `Format` The extending format for the previously mentioned type. */
  format: string;
  /** `Default` Declares the value of the header that the server will use if none is provided.
 See: https://tools.ietf.org/html/draft-fge-json-schema-validation-00#section-6.2.
 Unlike JSON Schema this value MUST conform to the defined type for the header. */
  default: string;
  /** 'Pattern' See https://tools.ietf.org/html/draft-fge-json-schema-validation-00#section-5.2.3. */
  pattern: string;
}


/**
 * `Response` is a representation of OpenAPI v2 specification's Response object.

 See: https://github.com/OAI/OpenAPI-Specification/blob/3.0.0/versions/2.0.md#responseObject
 */
export interface Response {
  /** `Description` is a short description of the response.
 GFM syntax can be used for rich text representation. */
  description: string;
  /** `Schema` optionally defines the structure of the response.
 If `Schema` is not provided, it means there is no content to the response. */
  schema?: Schema;
  /** `Headers` A list of headers that are sent with the response.
 `Header` name is expected to be a string in the canonical format of the MIME header key
 See: https://golang.org/pkg/net/textproto/#CanonicalMIMEHeaderKey */
  headers: Record<string, Header>;
  /** `Examples` gives per-mimetype response examples.
 See: https://github.com/OAI/OpenAPI-Specification/blob/3.0.0/versions/2.0.md#example-object */
  examples: Record<string, string>;
  /** Custom properties that start with "x-" such as "x-foo" used to describe
 extra functionality that is not covered by the standard OpenAPI Specification.
 See: https://swagger.io/docs/specification/2-0/swagger-extensions/ */
  extensions: Record<string, Value>;
}


/**
 * `Info` is a representation of OpenAPI v2 specification's Info object.

 See: https://github.com/OAI/OpenAPI-Specification/blob/3.0.0/versions/2.0.md#infoObject

 Example:

  option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_swagger) = {
    info: {
      title: "Echo API";
      version: "1.0";
      description: "";
      contact: {
        name: "gRPC-Gateway project";
        url: "https://github.com/grpc-ecosystem/grpc-gateway";
        email: "none@example.com";
      };
      license: {
        name: "BSD 3-Clause License";
        url: "https://github.com/grpc-ecosystem/grpc-gateway/blob/main/LICENSE";
      };
    };
    ...
  };
 */
export interface Info {
  /** The title of the application. */
  title: string;
  /** A short description of the application. GFM syntax can be used for rich
 text representation. */
  description: string;
  /** The Terms of Service for the API. */
  termsOfService: string;
  /** The contact information for the exposed API. */
  contact?: Contact;
  /** The license information for the exposed API. */
  license?: License;
  /** Provides the version of the application API (not to be confused
 with the specification version). */
  version: string;
  /** Custom properties that start with "x-" such as "x-foo" used to describe
 extra functionality that is not covered by the standard OpenAPI Specification.
 See: https://swagger.io/docs/specification/2-0/swagger-extensions/ */
  extensions: Record<string, Value>;
}


/**
 * `Contact` is a representation of OpenAPI v2 specification's Contact object.

 See: https://github.com/OAI/OpenAPI-Specification/blob/3.0.0/versions/2.0.md#contactObject

 Example:

  option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_swagger) = {
    info: {
      ...
      contact: {
        name: "gRPC-Gateway project";
        url: "https://github.com/grpc-ecosystem/grpc-gateway";
        email: "none@example.com";
      };
      ...
    };
    ...
  };
 */
export interface Contact {
  /** The identifying name of the contact person/organization. */
  name: string;
  /** The URL pointing to the contact information. MUST be in the format of a
 URL. */
  url: string;
  /** The email address of the contact person/organization. MUST be in the format
 of an email address. */
  email: string;
}


/**
 * `License` is a representation of OpenAPI v2 specification's License object.

 See: https://github.com/OAI/OpenAPI-Specification/blob/3.0.0/versions/2.0.md#licenseObject

 Example:

  option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_swagger) = {
    info: {
      ...
      license: {
        name: "BSD 3-Clause License";
        url: "https://github.com/grpc-ecosystem/grpc-gateway/blob/main/LICENSE";
      };
      ...
    };
    ...
  };
 */
export interface License {
  /** The license name used for the API. */
  name: string;
  /** A URL to the license used for the API. MUST be in the format of a URL. */
  url: string;
}


/**
 * `ExternalDocumentation` is a representation of OpenAPI v2 specification's
 ExternalDocumentation object.

 See: https://github.com/OAI/OpenAPI-Specification/blob/3.0.0/versions/2.0.md#externalDocumentationObject

 Example:

  option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_swagger) = {
    ...
    external_docs: {
      description: "More about gRPC-Gateway";
      url: "https://github.com/grpc-ecosystem/grpc-gateway";
    }
    ...
  };
 */
export interface ExternalDocumentation {
  /** A short description of the target documentation. GFM syntax can be used for
 rich text representation. */
  description: string;
  /** The URL for the target documentation. Value MUST be in the format
 of a URL. */
  url: string;
}


/**
 * `Schema` is a representation of OpenAPI v2 specification's Schema object.

 See: https://github.com/OAI/OpenAPI-Specification/blob/3.0.0/versions/2.0.md#schemaObject
 */
export interface Schema {
  jsonSchema?: JSONSchema;
  /** Adds support for polymorphism. The discriminator is the schema property
 name that is used to differentiate between other schema that inherit this
 schema. The property name used MUST be defined at this schema and it MUST
 be in the required property list. When used, the value MUST be the name of
 this schema or any schema that inherits it. */
  discriminator: string;
  /** Relevant only for Schema "properties" definitions. Declares the property as
 "read only". This means that it MAY be sent as part of a response but MUST
 NOT be sent as part of the request. Properties marked as readOnly being
 true SHOULD NOT be in the required list of the defined schema. Default
 value is false. */
  readOnly: boolean;
  /** Additional external documentation for this schema. */
  externalDocs?: ExternalDocumentation;
  /** A free-form property to include an example of an instance for this schema in JSON.
 This is copied verbatim to the output. */
  example: string;
}


/**
 * `EnumSchema` is subset of fields from the OpenAPI v2 specification's Schema object.
 Only fields that are applicable to Enums are included
 See: https://github.com/OAI/OpenAPI-Specification/blob/3.0.0/versions/2.0.md#schemaObject

 Example:

  option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_enum) = {
    ...
    title: "MyEnum";
    description:"This is my nice enum";
    example: "ZERO";
    required: true;
    ...
  };
 */
export interface EnumSchema {
  /** A short description of the schema. */
  description: string;
  default: string;
  /** The title of the schema. */
  title: string;
  required: boolean;
  readOnly: boolean;
  /** Additional external documentation for this schema. */
  externalDocs?: ExternalDocumentation;
  example: string;
  /** Ref is used to define an external reference to include in the message.
 This could be a fully qualified proto message reference, and that type must
 be imported into the protofile. If no message is identified, the Ref will
 be used verbatim in the output.
 For example:
  `ref: ".google.protobuf.Timestamp"`. */
  ref: string;
  /** Custom properties that start with "x-" such as "x-foo" used to describe
 extra functionality that is not covered by the standard OpenAPI Specification.
 See: https://swagger.io/docs/specification/2-0/swagger-extensions/ */
  extensions: Record<string, Value>;
}


/**
 * `JSONSchema` represents properties from JSON Schema taken, and as used, in
 the OpenAPI v2 spec.

 This includes changes made by OpenAPI v2.

 See: https://github.com/OAI/OpenAPI-Specification/blob/3.0.0/versions/2.0.md#schemaObject

 See also: https://cswr.github.io/JsonSchema/spec/basic_types/,
 https://github.com/json-schema-org/json-schema-spec/blob/master/schema.json

 Example:

  message SimpleMessage {
    option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_schema) = {
      json_schema: {
        title: "SimpleMessage"
        description: "A simple message."
        required: ["id"]
      }
    };

    // Id represents the message identifier.
    string id = 1; [
        (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_field) = {
          description: "The unique identifier of the simple message."
        }];
  }
 */
export interface JSONSchema {
  /** Ref is used to define an external reference to include in the message.
 This could be a fully qualified proto message reference, and that type must
 be imported into the protofile. If no message is identified, the Ref will
 be used verbatim in the output.
 For example:
  `ref: ".google.protobuf.Timestamp"`. */
  ref: string;
  /** The title of the schema. */
  title: string;
  /** A short description of the schema. */
  description: string;
  default: string;
  readOnly: boolean;
  /** A free-form property to include a JSON example of this field. This is copied
 verbatim to the output swagger.json. Quotes must be escaped.
 This property is the same for 2.0 and 3.0.0 https://github.com/OAI/OpenAPI-Specification/blob/3.0.0/versions/3.0.0.md#schemaObject  https://github.com/OAI/OpenAPI-Specification/blob/3.0.0/versions/2.0.md#schemaObject */
  example: string;
  multipleOf: number;
  /** Maximum represents an inclusive upper limit for a numeric instance. The
 value of MUST be a number, */
  maximum: number;
  exclusiveMaximum: boolean;
  /** minimum represents an inclusive lower limit for a numeric instance. The
 value of MUST be a number, */
  minimum: number;
  exclusiveMinimum: boolean;
  maxLength: number;
  minLength: number;
  pattern: string;
  maxItems: number;
  minItems: number;
  uniqueItems: boolean;
  maxProperties: number;
  minProperties: number;
  required: string[];
  /** Items in 'array' must be unique. */
  array: string[];
  type: JSONSchemaSimpleTypes[];
  /** `Format` */
  format: string;
  /** Items in `enum` must be unique https://tools.ietf.org/html/draft-fge-json-schema-validation-00#section-5.5.1 */
  enum: string[];
  /** Additional field level properties used when generating the OpenAPI v2 file. */
  fieldConfiguration?: FieldConfiguration;
  /** Custom properties that start with "x-" such as "x-foo" used to describe
 extra functionality that is not covered by the standard OpenAPI Specification.
 See: https://swagger.io/docs/specification/2-0/swagger-extensions/ */
  extensions: Record<string, Value>;
}


/**
 * 'FieldConfiguration' provides additional field level properties used when generating the OpenAPI v2 file.
 These properties are not defined by OpenAPIv2, but they are used to control the generation.
 */
export interface FieldConfiguration {
}


/**
 * `Tag` is a representation of OpenAPI v2 specification's Tag object.

 See: https://github.com/OAI/OpenAPI-Specification/blob/3.0.0/versions/2.0.md#tagObject
 */
export interface Tag {
  /** The name of the tag. Use it to allow override of the name of a
 global Tag object, then use that name to reference the tag throughout the
 OpenAPI file. */
  name: string;
  /** A short description for the tag. GFM syntax can be used for rich text
 representation. */
  description: string;
  /** Additional external documentation for this tag. */
  externalDocs?: ExternalDocumentation;
  /** Custom properties that start with "x-" such as "x-foo" used to describe
 extra functionality that is not covered by the standard OpenAPI Specification.
 See: https://swagger.io/docs/specification/2-0/swagger-extensions/ */
  extensions: Record<string, Value>;
}


/**
 * `SecurityDefinitions` is a representation of OpenAPI v2 specification's
 Security Definitions object.

 See: https://github.com/OAI/OpenAPI-Specification/blob/3.0.0/versions/2.0.md#securityDefinitionsObject

 A declaration of the security schemes available to be used in the
 specification. This does not enforce the security schemes on the operations
 and only serves to provide the relevant details for each scheme.
 */
export interface SecurityDefinitions {
  /** A single security scheme definition, mapping a "name" to the scheme it
 defines. */
  security: Record<string, SecurityScheme>;
}


/**
 * `SecurityScheme` is a representation of OpenAPI v2 specification's
 Security Scheme object.

 See: https://github.com/OAI/OpenAPI-Specification/blob/3.0.0/versions/2.0.md#securitySchemeObject

 Allows the definition of a security scheme that can be used by the
 operations. Supported schemes are basic authentication, an API key (either as
 a header or as a query parameter) and OAuth2's common flows (implicit,
 password, application and access code).
 */
export interface SecurityScheme {
  /** The type of the security scheme. Valid values are "basic",
 "apiKey" or "oauth2". */
  type: Type;
  /** A short description for security scheme. */
  description: string;
  /** The name of the header or query parameter to be used.
 Valid for apiKey. */
  name: string;
  /** The location of the API key. Valid values are "query" or
 "header".
 Valid for apiKey. */
  in: In;
  /** The flow used by the OAuth2 security scheme. Valid values are
 "implicit", "password", "application" or "accessCode".
 Valid for oauth2. */
  flow: Flow;
  /** The authorization URL to be used for this flow. This SHOULD be in
 the form of a URL.
 Valid for oauth2/implicit and oauth2/accessCode. */
  authorizationUrl: string;
  /** The token URL to be used for this flow. This SHOULD be in the
 form of a URL.
 Valid for oauth2/password, oauth2/application and oauth2/accessCode. */
  tokenUrl: string;
  /** The available scopes for the OAuth2 security scheme.
 Valid for oauth2. */
  scopes?: Scopes;
  /** Custom properties that start with "x-" such as "x-foo" used to describe
 extra functionality that is not covered by the standard OpenAPI Specification.
 See: https://swagger.io/docs/specification/2-0/swagger-extensions/ */
  extensions: Record<string, Value>;
}


/**
 * `SecurityRequirement` is a representation of OpenAPI v2 specification's
 Security Requirement object.

 See: https://github.com/OAI/OpenAPI-Specification/blob/3.0.0/versions/2.0.md#securityRequirementObject

 Lists the required security schemes to execute this operation. The object can
 have multiple security schemes declared in it which are all required (that
 is, there is a logical AND between the schemes).

 The name used for each property MUST correspond to a security scheme
 declared in the Security Definitions.
 */
export interface SecurityRequirement {
  /** Each name must correspond to a security scheme which is declared in
 the Security Definitions. If the security scheme is of type "oauth2",
 then the value is a list of scope names required for the execution.
 For other security scheme types, the array MUST be empty. */
  securityRequirement: Record<string, SecurityRequirementValue>;
}


/**
 * If the security scheme is of type "oauth2", then the value is a list of
 scope names required for the execution. For other security scheme types,
 the array MUST be empty.
 */
export interface SecurityRequirementValue {
}


/**
 * `Scopes` is a representation of OpenAPI v2 specification's Scopes object.

 See: https://github.com/OAI/OpenAPI-Specification/blob/3.0.0/versions/2.0.md#scopesObject

 Lists the available scopes for an OAuth2 security scheme.
 */
export interface Scopes {
  /** Maps between a name of a scope to a short description of it (as the value
 of the property). */
  scope: Record<string, string>;
}

