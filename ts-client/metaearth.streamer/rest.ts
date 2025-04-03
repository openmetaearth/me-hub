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

export interface ProtobufAny {
  "@type"?: string;
}

export interface RpcStatus {
  /** @format int32 */
  code?: number;
  message?: string;
  details?: ProtobufAny[];
}

export interface StreamerActiveStreamsResponse {
  /** Active gagues only */
  data?: StreamerStream[];

  /**
   * Pagination defines pagination for the response
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

export interface StreamerDistrInfo {
  total_weight?: string;
  records?: StreamerDistrRecord[];
}

export interface StreamerDistrRecord {
  /** @format uint64 */
  gauge_id?: string;
  weight?: string;
}

export interface StreamerModuleToDistributeCoinsResponse {
  /** Coins that have yet to be distributed */
  coins?: V1Beta1Coin[];
}

/**
* Stream is an object that stores and distributes yields to recipients who
satisfy certain conditions. Currently streams support conditions around the
duration for which a given denom is locked.
*/
export interface StreamerStream {
  /**
   * id is the unique ID of a Stream
   * @format uint64
   */
  id?: string;

  /** distribute_to is the distr_info. */
  distribute_to?: StreamerDistrInfo;

  /**
   * coins is the total amount of coins that have been in the stream
   * Can distribute multiple coin denoms
   */
  coins?: V1Beta1Coin[];

  /**
   * start_time is the distribution start time
   * @format date-time
   */
  start_time?: string;

  /**
   * distr_epoch_identifier is what epoch type di-stribution will be triggered
   * by
   * (day, week, etc.)
   */
  distr_epoch_identifier?: string;

  /**
   * num_epochs_paid_over is the number of total epochs distribution will be
   * completed over
   * @format uint64
   */
  num_epochs_paid_over?: string;

  /**
   * filled_epochs is the number of epochs distribution has been completed on
   * already
   * @format uint64
   */
  filled_epochs?: string;

  /** distributed_coins are coins that have been distributed already */
  distributed_coins?: V1Beta1Coin[];
}

export interface StreamerStreamByIDResponse {
  /**
   * Stream that corresponds to provided gague ID
   * Stream is an object that stores and distributes yields to recipients who
   * satisfy certain conditions. Currently streams support conditions around the
   * duration for which a given denom is locked.
   */
  stream?: StreamerStream;
}

export interface StreamerStreamsResponse {
  /** Upcoming and active streams */
  data?: StreamerStream[];

  /**
   * Pagination defines pagination for the response
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

export interface StreamerUpcomingStreamsResponse {
  /** Streams whose distribution is upcoming */
  data?: StreamerStream[];

  /**
   * Pagination defines pagination for the response
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
 * @title metaearth/streamer/distr_info.proto
 * @version version not set
 */
export class Api<SecurityDataType extends unknown> extends HttpClient<SecurityDataType> {
  /**
   * No description
   *
   * @tags Query
   * @name QueryActiveStreams
   * @summary ActiveStreams returns active streams
   * @request GET:/metaearth/streamer/active_streams
   */
  queryActiveStreams = (
    query?: {
      "pagination.key"?: string;
      "pagination.offset"?: string;
      "pagination.limit"?: string;
      "pagination.count_total"?: boolean;
      "pagination.reverse"?: boolean;
    },
    params: RequestParams = {},
  ) =>
    this.request<StreamerActiveStreamsResponse, RpcStatus>({
      path: `/metaearth/streamer/active_streams`,
      method: "GET",
      query: query,
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryModuleToDistributeCoins
   * @summary ModuleToDistributeCoins returns coins that are going to be distributed
   * @request GET:/metaearth/streamer/module_to_distribute_coins
   */
  queryModuleToDistributeCoins = (params: RequestParams = {}) =>
    this.request<StreamerModuleToDistributeCoinsResponse, RpcStatus>({
      path: `/metaearth/streamer/module_to_distribute_coins`,
      method: "GET",
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryStreamById
   * @summary StreamByID returns streams by their respective ID
   * @request GET:/metaearth/streamer/stream_by_id/{id}
   */
  queryStreamByID = (id: string, params: RequestParams = {}) =>
    this.request<StreamerStreamByIDResponse, RpcStatus>({
      path: `/metaearth/streamer/stream_by_id/${id}`,
      method: "GET",
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryStreams
   * @summary Streams returns both upcoming and active streams
   * @request GET:/metaearth/streamer/streams
   */
  queryStreams = (
    query?: {
      "pagination.key"?: string;
      "pagination.offset"?: string;
      "pagination.limit"?: string;
      "pagination.count_total"?: boolean;
      "pagination.reverse"?: boolean;
    },
    params: RequestParams = {},
  ) =>
    this.request<StreamerStreamsResponse, RpcStatus>({
      path: `/metaearth/streamer/streams`,
      method: "GET",
      query: query,
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryUpcomingStreams
   * @summary Returns scheduled streams that have not yet occurred
   * @request GET:/metaearth/streamer/upcoming_streams
   */
  queryUpcomingStreams = (
    query?: {
      "pagination.key"?: string;
      "pagination.offset"?: string;
      "pagination.limit"?: string;
      "pagination.count_total"?: boolean;
      "pagination.reverse"?: boolean;
    },
    params: RequestParams = {},
  ) =>
    this.request<StreamerUpcomingStreamsResponse, RpcStatus>({
      path: `/metaearth/streamer/upcoming_streams`,
      method: "GET",
      query: query,
      format: "json",
      ...params,
    });
}
