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

export enum DymensioncommonStatus {
  PENDING = "PENDING",
  FINALIZED = "FINALIZED",
  REVERTED = "REVERTED",
}

/**
 * Params defines the parameters for the module.
 */
export interface DymensionrollappParams {
  /**
   * dispute_period_in_blocks the number of blocks it takes
   * to change a status of a state from received to finalized.
   * during that period, any user could submit fraud proof
   * @format uint64
   */
  dispute_period_in_blocks?: string;

  /**
   * deployer_whitelist is a list of the
   * accounts that are allowed to create a rollapp and maximum number of
   * rollapps. In the case of an empty list, there are no restrictions
   */
  deployer_whitelist?: RollappDeployerParams[];
  rollapps_enabled?: boolean;
}

export interface DymensionrollappRollapp {
  /**
   * The unique identifier of the rollapp chain.
   * The rollappId follows the same standard as cosmos chain_id.
   */
  rollappId?: string;

  /** creator is the bech32-encoded address of the rollapp creator. */
  creator?: string;

  /**
   * version is the software and configuration version.
   * starts from 1 and increases by one on every MsgUpdateState
   * @format uint64
   */
  version?: string;

  /**
   * maxSequencers is the maximum number of sequencers.
   * @format uint64
   */
  maxSequencers?: string;

  /**
   * permissionedAddresses is a bech32-encoded address list of the sequencers
   * that are allowed to serve this rollappId. In the case of an empty list, the
   * rollapp is considered permissionless.
   */
  permissionedAddresses?: string[];

  /**
   * genesis_state is a partial repr of the state the hub can expect the rollapp
   * to be in upon genesis
   */
  genesis_state?: RollappRollappGenesisState;

  /** channel_id will be set to the canonical IBC channel of the rollapp. */
  channel_id?: string;

  /** frozen is a boolean that indicates if the rollapp is frozen. */
  frozen?: boolean;

  /** registeredDenoms is a list of registered denom bases on this rollapp */
  registeredDenoms?: string[];
}

/**
 * StateInfo defines a rollapps' state.
 */
export interface DymensionrollappStateInfo {
  /**
   * stateInfoIndex defines what rollapp the state belongs to
   * and in which index it can be referenced
   */
  stateInfoIndex?: RollappStateInfoIndex;

  /** sequencer is the bech32-encoded address of the sequencer sent the update */
  sequencer?: string;

  /**
   * startHeight is the block height of the first block in the batch
   * @format uint64
   */
  startHeight?: string;

  /**
   * numBlocks is the number of blocks included in this batch update
   * @format uint64
   */
  numBlocks?: string;

  /** DAPath is the description of the location on the DA layer */
  DAPath?: string;

  /**
   * version is the version of the rollapp
   * @format uint64
   */
  version?: string;

  /**
   * creationHeight is the height at which the UpdateState took place
   * @format uint64
   */
  creationHeight?: string;

  /** status is the status of the state update */
  status?: DymensioncommonStatus;

  /**
   * BDs is a list of block description objects (one per block)
   * the list must be ordered by height, starting from startHeight to
   * startHeight+numBlocks-1
   * BlockDescriptors defines list of BlockDescriptor.
   */
  BDs?: RollappBlockDescriptors;
}

export interface GooglerpcStatus {
  /** @format int32 */
  code?: number;
  message?: string;
  details?: ProtobufAny[];
}

export interface ProtobufAny {
  "@type"?: string;
}

/**
 * BlockDescriptor defines a single rollapp chain block description.
 */
export interface RollappBlockDescriptor {
  /**
   * height is the height of the block
   * @format uint64
   */
  height?: string;

  /**
   * stateRoot is a 32 byte array of the hash of the block (state root of the
   * block)
   * @format byte
   */
  stateRoot?: string;
}

/**
 * BlockDescriptors defines list of BlockDescriptor.
 */
export interface RollappBlockDescriptors {
  BD?: RollappBlockDescriptor[];
}

export interface RollappDeployerParams {
  /**
   * address is a bech32-encoded address of the
   * accounts that are allowed to create a rollapp.
   */
  address?: string;
}

export type RollappMsgCreateRollappResponse = object;

export type RollappMsgSkipDelayRollappResponse = object;

export type RollappMsgUpdateStateResponse = object;

export interface RollappQueryAllRollappResponse {
  rollapp?: RollappRollappSummary[];

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

export interface RollappQueryGetLatestHeightResponse {
  /** @format uint64 */
  height?: string;
}

export interface RollappQueryGetLatestStateIndexResponse {
  stateIndex?: RollappStateInfoIndex;
}

export interface RollappQueryGetRollappResponse {
  rollapp?: DymensionrollappRollapp;

  /** Defines the index of the last rollapp UpdateState. */
  latestStateIndex?: RollappStateInfoIndex;

  /** Defines the index of the last rollapp UpdateState that was finalized. */
  latestFinalizedStateIndex?: RollappStateInfoIndex;

  /** @format uint64 */
  latestHeight?: string;

  /** @format uint64 */
  latestFinalizedHeight?: string;
}

export interface RollappQueryGetStateInfoResponse {
  /** StateInfo defines a rollapps' state. */
  stateInfo?: DymensionrollappStateInfo;
}

/**
 * QueryParamsResponse is response type for the Query/Params RPC method.
 */
export interface RollappQueryParamsResponse {
  /** params holds all the parameters of this module. */
  params?: DymensionrollappParams;
}

export interface RollappQuerySkipDelayRollappResponse {
  rollapps?: string[];
}

export interface RollappRollappGenesisState {
  /**
   * If true, then full usage of the canonical ibc transfer channel is enabled.
   * Note: in v3.1.0 and prior this field marked the completion of the 'genesis
   * event' Keeping and renaming the field enables a seamless upgrade
   * https://www.notion.so/dymension/ADR-x-Genesis-Bridge-Phase-2-89769aa551b5440b9ed403a101775ce1?pvs=4#89698384d815435b87393dbe45bc5a74
   * to the new genesis transfer protocol
   * Note: if this field is false, ibc transfers may still be allowed in one or
   * either direction.
   */
  transfers_enabled?: boolean;
}

export interface RollappRollappSummary {
  /**
   * The unique identifier of the rollapp chain.
   * The rollappId follows the same standard as cosmos chain_id.
   */
  rollappId?: string;

  /** Defines the index of the last rollapp UpdateState. */
  latestStateIndex?: RollappStateInfoIndex;

  /** Defines the index of the last rollapp UpdateState that was finalized. */
  latestFinalizedStateIndex?: RollappStateInfoIndex;
}

export interface RollappStateInfoIndex {
  /**
   * rollappId is the rollapp that the sequencer belongs to and asking to update
   * it used to identify the what rollapp a StateInfo belongs
   * The rollappId follows the same standard as cosmos chain_id
   */
  rollappId?: string;

  /**
   * index is a sequential increasing number, updating on each
   * state update used for indexing to a specific state info, the first index is
   * 1
   * @format uint64
   */
  index?: string;
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
 * @title dymensionxyz/dymension/rollapp/block_descriptor.proto
 * @version version not set
 */
export class Api<SecurityDataType extends unknown> extends HttpClient<SecurityDataType> {
  /**
   * No description
   *
   * @tags Query
   * @name QueryRollappByEip155
   * @summary Queries a Rollapp by index.
   * @request GET:/dymensionxyz/dymension/rollapp/eip155/{eip155}
   */
  queryRollappByEIP155 = (eip155: string, params: RequestParams = {}) =>
    this.request<RollappQueryGetRollappResponse, GooglerpcStatus>({
      path: `/dymensionxyz/dymension/rollapp/eip155/${eip155}`,
      method: "GET",
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryLatestHeight
   * @summary Queries a LatestHeight by rollapp-id.
   * @request GET:/dymensionxyz/dymension/rollapp/latest_height/{rollappId}
   */
  queryLatestHeight = (rollappId: string, query?: { finalized?: boolean }, params: RequestParams = {}) =>
    this.request<RollappQueryGetLatestHeightResponse, GooglerpcStatus>({
      path: `/dymensionxyz/dymension/rollapp/latest_height/${rollappId}`,
      method: "GET",
      query: query,
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryLatestStateIndex
   * @summary Queries a LatestStateIndex by rollapp-id.
   * @request GET:/dymensionxyz/dymension/rollapp/latest_state_index/{rollappId}
   */
  queryLatestStateIndex = (rollappId: string, query?: { finalized?: boolean }, params: RequestParams = {}) =>
    this.request<RollappQueryGetLatestStateIndexResponse, GooglerpcStatus>({
      path: `/dymensionxyz/dymension/rollapp/latest_state_index/${rollappId}`,
      method: "GET",
      query: query,
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryParams
   * @summary Parameters queries the parameters of the module.
   * @request GET:/dymensionxyz/dymension/rollapp/params
   */
  queryParams = (params: RequestParams = {}) =>
    this.request<RollappQueryParamsResponse, GooglerpcStatus>({
      path: `/dymensionxyz/dymension/rollapp/params`,
      method: "GET",
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryRollappAll
   * @summary Queries a list of Rollapp items.
   * @request GET:/dymensionxyz/dymension/rollapp/rollapp
   */
  queryRollappAll = (
    query?: {
      "pagination.key"?: string;
      "pagination.offset"?: string;
      "pagination.limit"?: string;
      "pagination.count_total"?: boolean;
      "pagination.reverse"?: boolean;
    },
    params: RequestParams = {},
  ) =>
    this.request<RollappQueryAllRollappResponse, GooglerpcStatus>({
      path: `/dymensionxyz/dymension/rollapp/rollapp`,
      method: "GET",
      query: query,
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryRollapp
   * @summary Queries a Rollapp by index.
   * @request GET:/dymensionxyz/dymension/rollapp/rollapp/{rollappId}
   */
  queryRollapp = (rollappId: string, params: RequestParams = {}) =>
    this.request<RollappQueryGetRollappResponse, GooglerpcStatus>({
      path: `/dymensionxyz/dymension/rollapp/rollapp/${rollappId}`,
      method: "GET",
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QuerySkipDelayRollapp
   * @request GET:/dymensionxyz/dymension/rollapp/skip_delay_rollapp
   */
  querySkipDelayRollapp = (params: RequestParams = {}) =>
    this.request<RollappQuerySkipDelayRollappResponse, GooglerpcStatus>({
      path: `/dymensionxyz/dymension/rollapp/skip_delay_rollapp`,
      method: "GET",
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryStateInfo
   * @summary Queries a StateInfo by index.
   * @request GET:/dymensionxyz/dymension/rollapp/state_info/{rollappId}/{index}
   */
  queryStateInfo = (
    rollappId: string,
    index: string,
    query?: { height?: string; finalized?: boolean },
    params: RequestParams = {},
  ) =>
    this.request<RollappQueryGetStateInfoResponse, GooglerpcStatus>({
      path: `/dymensionxyz/dymension/rollapp/state_info/${rollappId}/${index}`,
      method: "GET",
      query: query,
      format: "json",
      ...params,
    });
}
