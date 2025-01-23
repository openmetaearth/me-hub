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

export interface Feemarketv1Params {
  /** no_base_fee forces the EIP-1559 base fee to 0 (needed for 0 price calls) */
  no_base_fee?: boolean;

  /**
   * base_fee_change_denominator bounds the amount the base fee can change
   * between blocks.
   * @format int64
   */
  base_fee_change_denominator?: number;

  /**
   * elasticity_multiplier bounds the maximum gas limit an EIP-1559 block may
   * have.
   * @format int64
   */
  elasticity_multiplier?: number;

  /**
   * enable_height defines at which block height the base fee calculation is enabled.
   * @format int64
   */
  enable_height?: string;

  /** base_fee for EIP-1559 blocks. */
  base_fee?: string;

  /** min_gas_price defines the minimum gas price value for cosmos and eth transactions */
  min_gas_price?: string;

  /**
   * min_gas_multiplier bounds the minimum gas used to be charged
   * to senders based on gas limit
   */
  min_gas_multiplier?: string;
}

export interface ProtobufAny {
  "@type"?: string;
}

export interface RpcStatus {
  /** @format int32 */
  code?: number;
  message?: string;
  details?: ProtobufAny[];
}

/**
* MsgUpdateParamsResponse defines the response structure for executing a
MsgUpdateParams message.
*/
export type V1MsgUpdateParamsResponse = object;

/**
 * QueryBaseFeeResponse returns the EIP1559 base fee.
 */
export interface V1QueryBaseFeeResponse {
  /** base_fee is the EIP1559 base fee */
  base_fee?: string;
}

/**
 * QueryBlockGasResponse returns block gas used for a given height.
 */
export interface V1QueryBlockGasResponse {
  /**
   * gas is the returned block gas
   * @format int64
   */
  gas?: string;
}

/**
 * QueryParamsResponse defines the response type for querying x/evm parameters.
 */
export interface V1QueryParamsResponse {
  /** params define the evm module parameters. */
  params?: Feemarketv1Params;
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
 * @title ethermint/feemarket/v1/events.proto
 * @version version not set
 */
export class Api<SecurityDataType extends unknown> extends HttpClient<SecurityDataType> {
  /**
   * No description
   *
   * @tags Query
   * @name QueryBaseFee
   * @summary BaseFee queries the base fee of the parent block of the current block.
   * @request GET:/ethermint/feemarket/v1/base_fee
   */
  queryBaseFee = (params: RequestParams = {}) =>
    this.request<V1QueryBaseFeeResponse, RpcStatus>({
      path: `/ethermint/feemarket/v1/base_fee`,
      method: "GET",
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryBlockGas
   * @summary BlockGas queries the gas used at a given block height
   * @request GET:/ethermint/feemarket/v1/block_gas
   */
  queryBlockGas = (params: RequestParams = {}) =>
    this.request<V1QueryBlockGasResponse, RpcStatus>({
      path: `/ethermint/feemarket/v1/block_gas`,
      method: "GET",
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryParams
   * @summary Params queries the parameters of x/feemarket module.
   * @request GET:/ethermint/feemarket/v1/params
   */
  queryParams = (params: RequestParams = {}) =>
    this.request<V1QueryParamsResponse, RpcStatus>({
      path: `/ethermint/feemarket/v1/params`,
      method: "GET",
      format: "json",
      ...params,
    });
}
