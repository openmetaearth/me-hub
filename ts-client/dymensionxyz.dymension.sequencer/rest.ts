/* eslint-disable */
/* tslint:disable */
/*
 * ---------------------------------------------------------------
 * ## THIS FILE WAS GENERATED VIA SWAGGER-TYPESCRIPT-API        ##
 * ##                                                           ##
 * ## AUTHOR: acacode                                           ##
 * ## SOURCE: https://github.com/acacode/swagger-typescript-api ##
 * ---------------------------------------------------------------
 */

/**
 * Params defines the parameters for the module.
 */
export interface DymensionsequencerParams {
  /**
   * Coin defines a token with a denomination and an amount.
   *
   * NOTE: The amount field is an Int which implements the custom method
   * signatures required by gogoproto.
   */
  min_bond?: V1Beta1Coin;

  /** unbonding_time is the time duration of unbonding. */
  unbonding_time?: string;
}

/**
* Sequencer defines a sequencer identified by its' address (sequencerAddress).
The sequencer could be attached to only one rollapp (rollappId).
*/
export interface DymensionsequencerSequencer {
  /**
   * sequencerAddress is the bech32-encoded address of the sequencer account
   * which is the account that the message was sent from.
   */
  sequencerAddress?: string;

  /**
   * pubkey is the public key of the sequencers' dymint client, as a Protobuf
   * Any.
   */
  dymintPubKey?: ProtobufAny;

  /** rollappId defines the rollapp to which the sequencer belongs. */
  rollappId?: string;

  /** description defines the descriptive terms for the sequencer. */
  description?: SequencerDescription;

  /**
   * jailed defined whether the sequencer has been jailed from bonded status or
   * not.
   */
  jailed?: boolean;

  /** proposer defines whether the sequencer is a proposer or not. */
  proposer?: boolean;

  /** status is the sequencer status (bonded/unbonding/unbonded). */
  status?: SequencerOperatingStatus;

  /** tokens define the delegated tokens (incl. self-delegation). */
  tokens?: V1Beta1Coin[];

  /**
   * unbonding_height defines, if unbonding, the height at which this sequencer
   * has begun unbonding.
   * @format int64
   */
  unbonding_height?: string;

  /**
   * unbond_time defines, if unbonding, the min time for the sequencer to
   * complete unbonding.
   * @format date-time
   */
  unbond_time?: string;
}

/**
* `Any` contains an arbitrary serialized protocol buffer message along with a
URL that describes the type of the serialized message.

Protobuf library provides support to pack/unpack Any values in the form
of utility functions or additional generated methods of the Any type.

Example 1: Pack and unpack a message in C++.

    Foo foo = ...;
    Any any;
    any.PackFrom(foo);
    ...
    if (any.UnpackTo(&foo)) {
      ...
    }

Example 2: Pack and unpack a message in Java.

    Foo foo = ...;
    Any any = Any.pack(foo);
    ...
    if (any.is(Foo.class)) {
      foo = any.unpack(Foo.class);
    }

 Example 3: Pack and unpack a message in Python.

    foo = Foo(...)
    any = Any()
    any.Pack(foo)
    ...
    if any.Is(Foo.DESCRIPTOR):
      any.Unpack(foo)
      ...

 Example 4: Pack and unpack a message in Go

     foo := &pb.Foo{...}
     any, err := anypb.New(foo)
     if err != nil {
       ...
     }
     ...
     foo := &pb.Foo{}
     if err := any.UnmarshalTo(foo); err != nil {
       ...
     }

The pack methods provided by protobuf library will by default use
'type.googleapis.com/full.type.name' as the type URL and the unpack
methods only use the fully qualified type name after the last '/'
in the type URL, for example "foo.bar.com/x/y.z" will yield type
name "y.z".


JSON
====
The JSON representation of an `Any` value uses the regular
representation of the deserialized, embedded message, with an
additional field `@type` which contains the type URL. Example:

    package google.profile;
    message Person {
      string first_name = 1;
      string last_name = 2;
    }

    {
      "@type": "type.googleapis.com/google.profile.Person",
      "firstName": <string>,
      "lastName": <string>
    }

If the embedded message type is well-known and has a custom JSON
representation, that representation will be embedded adding a field
`value` which holds the custom JSON in addition to the `@type`
field. Example (for message [google.protobuf.Duration][]):

    {
      "@type": "type.googleapis.com/google.protobuf.Duration",
      "value": "1.212s"
    }
*/
export interface ProtobufAny {
  /**
   * A URL/resource name that uniquely identifies the type of the serialized
   * protocol buffer message. This string must contain at least
   * one "/" character. The last segment of the URL's path must represent
   * the fully qualified name of the type (as in
   * `path/google.protobuf.Duration`). The name should be in a canonical form
   * (e.g., leading "." is not accepted).
   *
   * In practice, teams usually precompile into the binary all types that they
   * expect it to use in the context of Any. However, for URLs which use the
   * scheme `http`, `https`, or no scheme, one can optionally set up a type
   * server that maps type URLs to message definitions as follows:
   * * If no scheme is provided, `https` is assumed.
   * * An HTTP GET on the URL must yield a [google.protobuf.Type][]
   *   value in binary format, or produce an error.
   * * Applications are allowed to cache lookup results based on the
   *   URL, or have them precompiled into a binary to avoid any
   *   lookup. Therefore, binary compatibility needs to be preserved
   *   on changes to types. (Use versioned type names to manage
   *   breaking changes.)
   * Note: this functionality is not currently available in the official
   * protobuf release, and it is not used for type URLs beginning with
   * type.googleapis.com.
   * Schemes other than `http`, `https` (or the empty scheme) might be
   * used with implementation specific semantics.
   */
  "@type"?: string;
}

export interface RpcStatus {
  /** @format int32 */
  code?: number;
  message?: string;
  details?: ProtobufAny[];
}

/**
 * Description defines a sequencer description.
 */
export interface SequencerDescription {
  /** moniker defines a human-readable name for the sequencer. */
  moniker?: string;

  /** identity defines an optional identity signature (ex. UPort or Keybase). */
  identity?: string;

  /** website defines an optional website link. */
  website?: string;

  /** securityContact defines an optional email for security contact. */
  securityContact?: string;

  /** details define other optional details. */
  details?: string;
}

export type SequencerMsgCreateSequencerResponse = object;

/**
 * MsgUnbondResponse defines the Msg/Unbond response type.
 */
export interface SequencerMsgUnbondResponse {
  /** @format date-time */
  completion_time?: string;
}

/**
* - OPERATING_STATUS_UNBONDED: OPERATING_STATUS_UNBONDED defines a sequencer that is not active and won't
be scheduled
 - OPERATING_STATUS_UNBONDING: UNBONDING defines a sequencer that is currently unbonding.
 - OPERATING_STATUS_BONDED: OPERATING_STATUS_BONDED defines a sequencer that is bonded and can be
scheduled
*/
export enum SequencerOperatingStatus {
  OPERATING_STATUS_UNBONDED = "OPERATING_STATUS_UNBONDED",
  OPERATING_STATUS_UNBONDING = "OPERATING_STATUS_UNBONDING",
  OPERATING_STATUS_BONDED = "OPERATING_STATUS_BONDED",
}

export interface SequencerQueryGetSequencerResponse {
  /**
   * Sequencer defines a sequencer identified by its' address (sequencerAddress).
   * The sequencer could be attached to only one rollapp (rollappId).
   */
  sequencer?: DymensionsequencerSequencer;
}

export interface SequencerQueryGetSequencersByRollappByStatusResponse {
  sequencers?: DymensionsequencerSequencer[];
}

export interface SequencerQueryGetSequencersByRollappResponse {
  sequencers?: DymensionsequencerSequencer[];
}

/**
 * QueryParamsResponse is response type for the Query/Params RPC method.
 */
export interface SequencerQueryParamsResponse {
  /** params holds all the parameters of this module. */
  params?: DymensionsequencerParams;
}

export interface SequencerQuerySequencersResponse {
  sequencers?: DymensionsequencerSequencer[];

  /**
   * PageResponse is to be embedded in gRPC response messages where the
   * corresponding request message has used PageRequest.
   *
   *  message SomeResponse {
   *          repeated Bar results = 1;
   *          PageResponse page = 2;
   *  }
   */
  pagination?: V1Beta1PageResponse;
}

/**
* Coin defines a token with a denomination and an amount.

NOTE: The amount field is an Int which implements the custom method
signatures required by gogoproto.
*/
export interface V1Beta1Coin {
  denom?: string;
  amount?: string;
}

/**
* message SomeRequest {
         Foo some_parameter = 1;
         PageRequest pagination = 2;
 }
*/
export interface V1Beta1PageRequest {
  /**
   * key is a value returned in PageResponse.next_key to begin
   * querying the next page most efficiently. Only one of offset or key
   * should be set.
   * @format byte
   */
  key?: string;

  /**
   * offset is a numeric offset that can be used when key is unavailable.
   * It is less efficient than using key. Only one of offset or key should
   * be set.
   * @format uint64
   */
  offset?: string;

  /**
   * limit is the total number of results to be returned in the result page.
   * If left empty it will default to a value to be set by each app.
   * @format uint64
   */
  limit?: string;

  /**
   * count_total is set to true  to indicate that the result set should include
   * a count of the total number of items available for pagination in UIs.
   * count_total is only respected when offset is used. It is ignored when key
   * is set.
   */
  count_total?: boolean;

  /**
   * reverse is set to true if results are to be returned in the descending order.
   *
   * Since: cosmos-sdk 0.43
   */
  reverse?: boolean;
}

/**
* PageResponse is to be embedded in gRPC response messages where the
corresponding request message has used PageRequest.

 message SomeResponse {
         repeated Bar results = 1;
         PageResponse page = 2;
 }
*/
export interface V1Beta1PageResponse {
  /**
   * next_key is the key to be passed to PageRequest.key to
   * query the next page most efficiently. It will be empty if
   * there are no more results.
   * @format byte
   */
  next_key?: string;

  /**
   * total is total number of results available if PageRequest.count_total
   * was set, its value is undefined otherwise
   * @format uint64
   */
  total?: string;
}

import axios, { AxiosInstance, AxiosRequestConfig, AxiosResponse, ResponseType } from "axios";

export type QueryParamsType = Record<string | number, any>;

export interface FullRequestParams extends Omit<AxiosRequestConfig, "data" | "params" | "url" | "responseType"> {
  /** set parameter to `true` for call `securityWorker` for this request */
  secure?: boolean;
  /** request path */
  path: string;
  /** content type of request body */
  type?: ContentType;
  /** query params */
  query?: QueryParamsType;
  /** format of response (i.e. response.json() -> format: "json") */
  format?: ResponseType;
  /** request body */
  body?: unknown;
}

export type RequestParams = Omit<FullRequestParams, "body" | "method" | "query" | "path">;

export interface ApiConfig<SecurityDataType = unknown> extends Omit<AxiosRequestConfig, "data" | "cancelToken"> {
  securityWorker?: (
    securityData: SecurityDataType | null,
  ) => Promise<AxiosRequestConfig | void> | AxiosRequestConfig | void;
  secure?: boolean;
  format?: ResponseType;
}

export enum ContentType {
  Json = "application/json",
  FormData = "multipart/form-data",
  UrlEncoded = "application/x-www-form-urlencoded",
}

export class HttpClient<SecurityDataType = unknown> {
  public instance: AxiosInstance;
  private securityData: SecurityDataType | null = null;
  private securityWorker?: ApiConfig<SecurityDataType>["securityWorker"];
  private secure?: boolean;
  private format?: ResponseType;

  constructor({ securityWorker, secure, format, ...axiosConfig }: ApiConfig<SecurityDataType> = {}) {
    this.instance = axios.create({ ...axiosConfig, baseURL: axiosConfig.baseURL || "" });
    this.secure = secure;
    this.format = format;
    this.securityWorker = securityWorker;
  }

  public setSecurityData = (data: SecurityDataType | null) => {
    this.securityData = data;
  };

  private mergeRequestParams(params1: AxiosRequestConfig, params2?: AxiosRequestConfig): AxiosRequestConfig {
    return {
      ...this.instance.defaults,
      ...params1,
      ...(params2 || {}),
      headers: {
        ...(this.instance.defaults.headers || {}),
        ...(params1.headers || {}),
        ...((params2 && params2.headers) || {}),
      },
    };
  }

  private createFormData(input: Record<string, unknown>): FormData {
    return Object.keys(input || {}).reduce((formData, key) => {
      const property = input[key];
      formData.append(
        key,
        property instanceof Blob
          ? property
          : typeof property === "object" && property !== null
          ? JSON.stringify(property)
          : `${property}`,
      );
      return formData;
    }, new FormData());
  }

  public request = async <T = any, _E = any>({
    secure,
    path,
    type,
    query,
    format,
    body,
    ...params
  }: FullRequestParams): Promise<AxiosResponse<T>> => {
    const secureParams =
      ((typeof secure === "boolean" ? secure : this.secure) &&
        this.securityWorker &&
        (await this.securityWorker(this.securityData))) ||
      {};
    const requestParams = this.mergeRequestParams(params, secureParams);
    const responseFormat = (format && this.format) || void 0;

    if (type === ContentType.FormData && body && body !== null && typeof body === "object") {
      requestParams.headers.common = { Accept: "*/*" };
      requestParams.headers.post = {};
      requestParams.headers.put = {};

      body = this.createFormData(body as Record<string, unknown>);
    }

    return this.instance.request({
      ...requestParams,
      headers: {
        ...(type && type !== ContentType.FormData ? { "Content-Type": type } : {}),
        ...(requestParams.headers || {}),
      },
      params: query,
      responseType: responseFormat,
      data: body,
      url: path,
    });
  };
}

/**
 * @title dymensionxyz/dymension/sequencer/description.proto
 * @version version not set
 */
export class Api<SecurityDataType extends unknown> extends HttpClient<SecurityDataType> {
  /**
   * No description
   *
   * @tags Query
   * @name QueryParams
   * @summary Parameters queries the parameters of the module.
   * @request GET:/dymensionxyz/dymension/sequencer/params
   */
  queryParams = (params: RequestParams = {}) =>
    this.request<SequencerQueryParamsResponse, RpcStatus>({
      path: `/dymensionxyz/dymension/sequencer/params`,
      method: "GET",
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QuerySequencers
   * @summary Queries a list of Sequencer items.
   * @request GET:/dymensionxyz/dymension/sequencer/sequencer
   */
  querySequencers = (
    query?: {
      "pagination.key"?: string;
      "pagination.offset"?: string;
      "pagination.limit"?: string;
      "pagination.count_total"?: boolean;
      "pagination.reverse"?: boolean;
    },
    params: RequestParams = {},
  ) =>
    this.request<SequencerQuerySequencersResponse, RpcStatus>({
      path: `/dymensionxyz/dymension/sequencer/sequencer`,
      method: "GET",
      query: query,
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QuerySequencer
   * @summary Queries a Sequencer by address.
   * @request GET:/dymensionxyz/dymension/sequencer/sequencer/{sequencerAddress}
   */
  querySequencer = (sequencerAddress: string, params: RequestParams = {}) =>
    this.request<SequencerQueryGetSequencerResponse, RpcStatus>({
      path: `/dymensionxyz/dymension/sequencer/sequencer/${sequencerAddress}`,
      method: "GET",
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QuerySequencersByRollapp
   * @summary Queries a SequencersByRollapp by rollappId.
   * @request GET:/dymensionxyz/dymension/sequencer/sequencers_by_rollapp/{rollappId}
   */
  querySequencersByRollapp = (rollappId: string, params: RequestParams = {}) =>
    this.request<SequencerQueryGetSequencersByRollappResponse, RpcStatus>({
      path: `/dymensionxyz/dymension/sequencer/sequencers_by_rollapp/${rollappId}`,
      method: "GET",
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QuerySequencersByRollappByStatus
   * @summary Queries a SequencersByRollappByStatus
   * @request GET:/dymensionxyz/dymension/sequencer/sequencers_by_rollapp/{rollappId}/{status}
   */
  querySequencersByRollappByStatus = (
    rollappId: string,
    status: "OPERATING_STATUS_UNBONDED" | "OPERATING_STATUS_UNBONDING" | "OPERATING_STATUS_BONDED",
    params: RequestParams = {},
  ) =>
    this.request<SequencerQueryGetSequencersByRollappByStatusResponse, RpcStatus>({
      path: `/dymensionxyz/dymension/sequencer/sequencers_by_rollapp/${rollappId}/${status}`,
      method: "GET",
      format: "json",
      ...params,
    });
}
