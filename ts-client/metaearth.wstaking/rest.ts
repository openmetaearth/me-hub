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

export interface MetaearthwstakingFixedDeposit {
  /** @format uint64 */
  id?: string;
  account?: string;

  /**
   * Coin defines a token with a denomination and an amount.
   *
   * NOTE: The amount field is an Int which implements the custom method
   * signatures required by gogoproto.
   */
  principal?: V1Beta1Coin;

  /**
   * Coin defines a token with a denomination and an amount.
   *
   * NOTE: The amount field is an Int which implements the custom method
   * signatures required by gogoproto.
   */
  interest?: V1Beta1Coin;

  /** @format date-time */
  start_time?: string;

  /** @format date-time */
  end_time?: string;

  /** @format int64 */
  term?: string;
  rate?: string;
}

export interface MetaearthwstakingFixedDepositCfg {
  regionId?: string;

  /** @format int64 */
  term?: string;
  rate?: string;
  status?: WstakingFIXEDDEPOSITCFGSTATUS;
}

export interface MetaearthwstakingHeight {
  /**
   * the revision that the client is currently on
   * @format uint64
   */
  revision_number?: string;

  /**
   * the height within the given revision
   * @format uint64
   */
  revision_height?: string;
}

export interface MetaearthwstakingRegion {
  regionId?: string;
  name?: string;
  creator?: string;
  operator_address?: string;
  nft_class_id?: string;
  region_treasure_addr?: string;
  deposit_interest_addr?: string;

  /** tokens define the region tokens share */
  region_share?: string;
  delegate_interest?: string;
  delegate_amount?: string;
  fixed_deposit_amount?: string;
}

export interface MetaearthwstakingReviewRecord {
  recordHash?: string;
  actionNumber?: string;
  recordResult?: string;
  reviewedAddress?: string;
}

/**
* Stake represents the bond with tokens held by an account. It is
owned by one staker, and is associated with the voting power of one
validator.
*/
export interface MetaearthwstakingStake {
  /** staker_address is the bech32-encoded address of the staker. */
  staker_address?: string;

  /** validator_address is the bech32-encoded address of the validator. */
  validator_address?: string;

  /** shares define the stake shares received. */
  shares?: string;

  /** @format int64 */
  startHeight?: string;
  rewards?: string;
  amount?: string;
  unmovable?: string;
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
* Delegation represents the bond with tokens held by an account. It is
owned by one delegator, and is associated with the voting power of one
validator.
*/
export interface Stakingv1Beta1Delegation {
  /** delegator_address is the bech32-encoded address of the delegator. */
  delegator_address?: string;

  /** validator_address is the bech32-encoded address of the validator. */
  validator_address?: string;

  /** shares define the delegation shares received. */
  shares?: string;

  /** @format int64 */
  startHeight?: string;
  amount?: string;
  unmovable?: string;
  unMeidAmount?: string;
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
* DecCoin defines a token with a denomination and a decimal amount.

NOTE: The amount field is an Dec which implements the custom method
signatures required by gogoproto.
*/
export interface V1Beta1DecCoin {
  denom?: string;
  amount?: string;
}

/**
* DelegationResponse is equivalent to Delegation except that it contains a
balance in addition to shares which is more suitable for client responses.
*/
export interface V1Beta1DelegationResponse {
  /**
   * Delegation represents the bond with tokens held by an account. It is
   * owned by one delegator, and is associated with the voting power of one
   * validator.
   */
  delegation?: Stakingv1Beta1Delegation;

  /**
   * Coin defines a token with a denomination and an amount.
   *
   * NOTE: The amount field is an Int which implements the custom method
   * signatures required by gogoproto.
   */
  balance?: V1Beta1Coin;
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

/**
 * QueryDelegationResponse is response type for the Query/Delegation RPC method.
 */
export interface V1Beta1QueryDelegationResponse {
  /** delegation_responses defines the delegation info of a delegation. */
  delegation_response?: V1Beta1DelegationResponse;
}

export enum WstakingFIXEDDEPOSITCFGSTATUS {
  FIXED_DEPOSIT_CFG_ACTIVE = "FIXED_DEPOSIT_CFG_ACTIVE",
  FIXED_DEPOSIT_CFG_INACTIVE = "FIXED_DEPOSIT_CFG_INACTIVE",
}

export enum WstakingFixedDepositState {
  ALL_STATE = "ALL_STATE",
  NOT_EXPIRED = "NOT_EXPIRED",
  EXPIRED = "EXPIRED",
}

export interface WstakingMsgDoFixedDepositResponse {
  /** @format uint64 */
  id?: string;
}

export type WstakingMsgIbcTransferFromRegionTreasureResponse = object;

export interface WstakingMsgNewFixedDepositCfgResp {
  retcode?: string;
}

export type WstakingMsgNewRecordResponse = object;

export interface WstakingMsgNewRegionResponse {
  region_id?: string;
}

export interface WstakingMsgRemoveFixedDepositCfgResp {
  retcode?: string;
}

export interface WstakingMsgRemoveRegionResponse {
  retcode?: string;
}

export type WstakingMsgReviewRecordResponse = object;

export interface WstakingMsgSetFixedDepositCfgRateResp {
  retcode?: string;
}

export interface WstakingMsgSetFixedDepositCfgStatusResp {
  retcode?: string;
}

/**
 * MsgStakeResponse defines the Msg/Stake response type.
 */
export type WstakingMsgStakeResponse = object;

export type WstakingMsgTransferRegionResponse = object;

/**
 * MsgUnstakeResponse defines the Msg/MsgUnstake response type.
 */
export interface WstakingMsgUnstakeResponse {
  /** @format date-time */
  completion_time?: string;
}

/**
* MsgWithdrawDelegatorRewardResponse defines the Msg/WithdrawDelegatorReward
response type.
*/
export interface WstakingMsgWithdrawDelegatorRewardResponse {
  /** Since: cosmos-sdk 0.46 */
  amount?: V1Beta1Coin[];
}

export interface WstakingMsgWithdrawFixedDepositResponse {
  /**
   * Coin defines a token with a denomination and an amount.
   *
   * NOTE: The amount field is an Int which implements the custom method
   * signatures required by gogoproto.
   */
  principal?: V1Beta1Coin;

  /**
   * Coin defines a token with a denomination and an amount.
   *
   * NOTE: The amount field is an Int which implements the custom method
   * signatures required by gogoproto.
   */
  interest?: V1Beta1Coin;

  /** @format int64 */
  term?: string;
  rate?: string;
}

export type WstakingMsgWithdrawFromGlobalDaoFeePoolResp = object;

export type WstakingMsgWithdrawFromRegionResp = object;

export interface WstakingQueryAllDelegationsResponse {
  delegations?: Stakingv1Beta1Delegation[];

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

export interface WstakingQueryAllFixedDepositResponse {
  FixedDeposit?: MetaearthwstakingFixedDeposit[];

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

export interface WstakingQueryAllRecordsResponse {
  records?: WstakingRecord[];

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

export interface WstakingQueryAllRegionResponse {
  region?: MetaearthwstakingRegion[];

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

export interface WstakingQueryDelegationRewardsResponse {
  /** rewards defines the rewards accrued by a delegation. */
  rewards?: V1Beta1DecCoin[];
}

export interface WstakingQueryFixedDepositAmountByMeidResponse {
  /**
   * Coin defines a token with a denomination and an amount.
   *
   * NOTE: The amount field is an Int which implements the custom method
   * signatures required by gogoproto.
   */
  amount?: V1Beta1Coin;
}

export interface WstakingQueryFixedDepositByAcctResponse {
  /** cosmos.base.query.v1beta1.PageResponse pagination = 2; */
  FixedDeposit?: MetaearthwstakingFixedDeposit[];
}

export interface WstakingQueryFixedDepositByRegionResponse {
  FixedDeposit?: MetaearthwstakingFixedDeposit[];

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

export interface WstakingQueryFixedDepositCfgByTermResponse {
  FixedDepositCfg?: MetaearthwstakingFixedDepositCfg;
}

export interface WstakingQueryFixedDepositCfgResponse {
  RegionFixedDepositCfgs?: WstakingRegionAllFixedDepositCfg[];
}

export interface WstakingQueryFixedDepositTotalAmountResponse {
  /**
   * Coin defines a token with a denomination and an amount.
   *
   * NOTE: The amount field is an Int which implements the custom method
   * signatures required by gogoproto.
   */
  amount?: V1Beta1Coin;
}

export interface WstakingQueryGetFixedDepositResponse {
  FixedDeposit?: MetaearthwstakingFixedDeposit;
}

export interface WstakingQueryRecordsByAddressResponse {
  /** cosmos.base.query.v1beta1.PageResponse pagination = 2; */
  records?: WstakingRecord[];
}

export interface WstakingQueryRegionResponse {
  region?: MetaearthwstakingRegion;
}

export interface WstakingQueryReviewRecordByNumberResponse {
  /** cosmos.base.query.v1beta1.PageResponse pagination = 2; */
  reviewRecord?: MetaearthwstakingReviewRecord;
}

export interface WstakingQueryStakesResponse {
  stakes?: MetaearthwstakingStake[];

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

export interface WstakingRecord {
  recordNumber?: string;
  url?: string;
  from?: string;
}

export interface WstakingRegionAllFixedDepositCfg {
  regionId?: string;
  RegionFixedDepositCfg?: WstakingRegionFixedDepositCfg[];
}

export interface WstakingRegionFixedDepositCfg {
  /** @format int64 */
  term?: string;
  rate?: string;
  status?: WstakingFIXEDDEPOSITCFGSTATUS;
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
 * @title metaearth/wstaking/fixed_deposit.proto
 * @version version not set
 */
export class Api<SecurityDataType extends unknown> extends HttpClient<SecurityDataType> {
  /**
   * No description
   *
   * @tags Query
   * @name QueryFixedDepositByRegion
   * @summary Queries a list of FixedDepositByRegion items.
   * @request GET:/cosmos/staking/v1beta1/fixed_deposit_by_region/{region_id}
   */
  queryFixedDepositByRegion = (
    regionId: string,
    query?: {
      "pagination.key"?: string;
      "pagination.offset"?: string;
      "pagination.limit"?: string;
      "pagination.count_total"?: boolean;
      "pagination.reverse"?: boolean;
      query_type?: "ALL_STATE" | "NOT_EXPIRED" | "EXPIRED";
    },
    params: RequestParams = {},
  ) =>
    this.request<WstakingQueryFixedDepositByRegionResponse, RpcStatus>({
      path: `/cosmos/staking/v1beta1/fixed_deposit_by_region/${regionId}`,
      method: "GET",
      query: query,
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryAllDelegations
   * @request GET:/metaearth/wstaking/all-delegations
   */
  queryAllDelegations = (
    query?: {
      "pagination.key"?: string;
      "pagination.offset"?: string;
      "pagination.limit"?: string;
      "pagination.count_total"?: boolean;
      "pagination.reverse"?: boolean;
    },
    params: RequestParams = {},
  ) =>
    this.request<WstakingQueryAllDelegationsResponse, RpcStatus>({
      path: `/metaearth/wstaking/all-delegations`,
      method: "GET",
      query: query,
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryAllRegion
   * @request GET:/metaearth/wstaking/all-region
   */
  queryAllRegion = (
    query?: {
      "pagination.key"?: string;
      "pagination.offset"?: string;
      "pagination.limit"?: string;
      "pagination.count_total"?: boolean;
      "pagination.reverse"?: boolean;
    },
    params: RequestParams = {},
  ) =>
    this.request<WstakingQueryAllRegionResponse, RpcStatus>({
      path: `/metaearth/wstaking/all-region`,
      method: "GET",
      query: query,
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryDelegationRewards
   * @request GET:/metaearth/wstaking/delegation-rewards/{delegator_address}
   */
  queryDelegationRewards = (
    delegatorAddress: string,
    query?: { validator_address?: string },
    params: RequestParams = {},
  ) =>
    this.request<WstakingQueryDelegationRewardsResponse, RpcStatus>({
      path: `/metaearth/wstaking/delegation-rewards/${delegatorAddress}`,
      method: "GET",
      query: query,
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryDelegation
   * @summary Delegation queries delegate info for given validator delegator pair.
   * @request GET:/metaearth/wstaking/delegation/{delegator_addr}
   */
  queryDelegation = (delegatorAddr: string, query?: { validator_addr?: string }, params: RequestParams = {}) =>
    this.request<V1Beta1QueryDelegationResponse, RpcStatus>({
      path: `/metaearth/wstaking/delegation/${delegatorAddr}`,
      method: "GET",
      query: query,
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryFixedDepositAll
   * @request GET:/metaearth/wstaking/fixed_deposit
   */
  queryFixedDepositAll = (
    query?: {
      "pagination.key"?: string;
      "pagination.offset"?: string;
      "pagination.limit"?: string;
      "pagination.count_total"?: boolean;
      "pagination.reverse"?: boolean;
    },
    params: RequestParams = {},
  ) =>
    this.request<WstakingQueryAllFixedDepositResponse, RpcStatus>({
      path: `/metaearth/wstaking/fixed_deposit`,
      method: "GET",
      query: query,
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryFixedDeposit
   * @summary Queries a list of FixedDeposit items.
   * @request GET:/metaearth/wstaking/fixed_deposit/{id}
   */
  queryFixedDeposit = (id: string, query?: { address?: string }, params: RequestParams = {}) =>
    this.request<WstakingQueryGetFixedDepositResponse, RpcStatus>({
      path: `/metaearth/wstaking/fixed_deposit/${id}`,
      method: "GET",
      query: query,
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryFixedDepositAmountByMeid
   * @request GET:/metaearth/wstaking/fixed_deposit_amount_by_meid/{account}
   */
  queryFixedDepositAmountByMeid = (account: string, params: RequestParams = {}) =>
    this.request<WstakingQueryFixedDepositAmountByMeidResponse, RpcStatus>({
      path: `/metaearth/wstaking/fixed_deposit_amount_by_meid/${account}`,
      method: "GET",
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryFixedDepositByAcct
   * @summary Queries a list of FixedDepositByAcct items.
   * @request GET:/metaearth/wstaking/fixed_deposit_by_acct/{account}/{query_type}
   */
  queryFixedDepositByAcct = (
    account: string,
    queryType: "ALL_STATE" | "NOT_EXPIRED" | "EXPIRED",
    params: RequestParams = {},
  ) =>
    this.request<WstakingQueryFixedDepositByAcctResponse, RpcStatus>({
      path: `/metaearth/wstaking/fixed_deposit_by_acct/${account}/${queryType}`,
      method: "GET",
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryFixedDepositCfg
   * @request GET:/metaearth/wstaking/fixed_deposit_cfg
   */
  queryFixedDepositCfg = (query?: { regionIds?: string[] }, params: RequestParams = {}) =>
    this.request<WstakingQueryFixedDepositCfgResponse, RpcStatus>({
      path: `/metaearth/wstaking/fixed_deposit_cfg`,
      method: "GET",
      query: query,
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryFixedDepositCfgByTerm
   * @request GET:/metaearth/wstaking/fixed_deposit_cfg_by_term
   */
  queryFixedDepositCfgByTerm = (query?: { regionId?: string; term?: string }, params: RequestParams = {}) =>
    this.request<WstakingQueryFixedDepositCfgByTermResponse, RpcStatus>({
      path: `/metaearth/wstaking/fixed_deposit_cfg_by_term`,
      method: "GET",
      query: query,
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryFixedDepositTotalAmount
   * @request GET:/metaearth/wstaking/fixed_deposit_total_amount
   */
  queryFixedDepositTotalAmount = (params: RequestParams = {}) =>
    this.request<WstakingQueryFixedDepositTotalAmountResponse, RpcStatus>({
      path: `/metaearth/wstaking/fixed_deposit_total_amount`,
      method: "GET",
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryQueryRecordByAddress
   * @request GET:/metaearth/wstaking/record/{account}
   */
  queryQueryRecordByAddress = (account: string, params: RequestParams = {}) =>
    this.request<WstakingQueryRecordsByAddressResponse, RpcStatus>({
      path: `/metaearth/wstaking/record/${account}`,
      method: "GET",
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryQueryAllRecord
   * @summary Queries a list of Record
   * @request GET:/metaearth/wstaking/records
   */
  queryQueryAllRecord = (
    query?: {
      "pagination.key"?: string;
      "pagination.offset"?: string;
      "pagination.limit"?: string;
      "pagination.count_total"?: boolean;
      "pagination.reverse"?: boolean;
    },
    params: RequestParams = {},
  ) =>
    this.request<WstakingQueryAllRecordsResponse, RpcStatus>({
      path: `/metaearth/wstaking/records`,
      method: "GET",
      query: query,
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryRegion
   * @summary Queries a list of Region items.
   * @request GET:/metaearth/wstaking/region/{regionId}
   */
  queryRegion = (regionId: string, params: RequestParams = {}) =>
    this.request<WstakingQueryRegionResponse, RpcStatus>({
      path: `/metaearth/wstaking/region/${regionId}`,
      method: "GET",
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryQueryReviewRecordById
   * @request GET:/metaearth/wstaking/review_record/{action_number}
   */
  queryQueryReviewRecordByID = (actionNumber: string, params: RequestParams = {}) =>
    this.request<WstakingQueryReviewRecordByNumberResponse, RpcStatus>({
      path: `/metaearth/wstaking/review_record/${actionNumber}`,
      method: "GET",
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryStakes
   * @request GET:/metaearth/wstaking/stakes
   */
  queryStakes = (
    query?: {
      "pagination.key"?: string;
      "pagination.offset"?: string;
      "pagination.limit"?: string;
      "pagination.count_total"?: boolean;
      "pagination.reverse"?: boolean;
    },
    params: RequestParams = {},
  ) =>
    this.request<WstakingQueryStakesResponse, RpcStatus>({
      path: `/metaearth/wstaking/stakes`,
      method: "GET",
      query: query,
      format: "json",
      ...params,
    });
}
