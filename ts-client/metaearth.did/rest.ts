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

export interface DidDidInfo {
  /**
   * same as the MEID
   * MEIDNFT.umeid
   */
  did?: string;

  /** MEID.account */
  address?: string;

  /**
   * public_key is mapped to the user address
   * the issuer will use public_key to encrypt the user's certificate to ensure
   * the privacy of the off-chain certificate
   */
  pubkey?: string;
  status?: DidDidStatus;

  /** unused! */
  regionId?: string;
  kycLevel?: DidKycLevel;
}

export enum DidDidStatus {
  DID_STATUS_INACTIVE = "DID_STATUS_INACTIVE",
  DID_STATUS_ACTIVE = "DID_STATUS_ACTIVE",
}

export enum DidKycLevel {
  KYC_LEVEL_NONE = "KYC_LEVEL_NONE",
  KYC_LEVEL_ONE = "KYC_LEVEL_ONE",
  KYC_LEVEL_TWO = "KYC_LEVEL_TWO",
  KYC_LEVEL_THREE = "KYC_LEVEL_THREE",
  KYC_LEVEL_FOUR = "KYC_LEVEL_FOUR",
  KYC_LEVEL_FIVE = "KYC_LEVEL_FIVE",
}

export type DidMsgCreateDidResponse = object;

export type DidMsgCreateServiceResponse = object;

export type DidMsgCreateVCResponse = object;

export type DidMsgRemoveVCResponse = object;

export type DidMsgUpdateDidStatusResponse = object;

export type DidMsgUpdateServiceStatusResponse = object;

export type DidMsgUpdateVCResponse = object;

export interface DidQueryCredentialResponse {
  credential?: MetaearthdidCredential;
}

export interface DidQueryCredentialsResponse {
  credentials?: MetaearthdidCredential[];

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

export interface DidQueryDidDocumentResponse {
  doc?: MetaearthdidDidDocument;
}

export interface DidQueryDidInfosResponse {
  infos?: DidDidInfo[];

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

export interface DidQueryDidResponse {
  info?: DidDidInfo;
}

export interface DidQueryServiceResponse {
  service?: MetaearthdidService;
}

export enum DidServiceStatus {
  SERVICE_STATUS_INACTIVE = "SERVICE_STATUS_INACTIVE",
  SERVICE_STATUS_ACTIVE = "SERVICE_STATUS_ACTIVE",
}

export interface MetaearthdidCredential {
  did?: string;
  sid?: string;
  hash?: string;
  uri?: string;

  /**
   * data is an app specific data of the credential. Optional
   * @format byte
   */
  data?: string;
}

export interface MetaearthdidDidDocument {
  info?: DidDidInfo;
  vcs?: MetaearthdidCredential[];
}

export interface MetaearthdidService {
  sid?: string;
  name?: string;
  description?: string;
  issuers?: string[];
  status?: DidServiceStatus;
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
 * @title metaearth/did/credential.proto
 * @version version not set
 */
export class Api<SecurityDataType extends unknown> extends HttpClient<SecurityDataType> {
  /**
   * No description
   *
   * @tags Query
   * @name QueryCredential
   * @request GET:/metaearth/did/credential
   */
  queryCredential = (query?: { did?: string; sid?: string }, params: RequestParams = {}) =>
    this.request<DidQueryCredentialResponse, RpcStatus>({
      path: `/metaearth/did/credential`,
      method: "GET",
      query: query,
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryCredentials
   * @request GET:/metaearth/did/credentials
   */
  queryCredentials = (
    query?: {
      sid?: string;
      filter?: string;
      "pagination.key"?: string;
      "pagination.offset"?: string;
      "pagination.limit"?: string;
      "pagination.count_total"?: boolean;
      "pagination.reverse"?: boolean;
    },
    params: RequestParams = {},
  ) =>
    this.request<DidQueryCredentialsResponse, RpcStatus>({
      path: `/metaearth/did/credentials`,
      method: "GET",
      query: query,
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryDid
   * @summary query user did doc
   * @request GET:/metaearth/did/did
   */
  queryDid = (query?: { address?: string }, params: RequestParams = {}) =>
    this.request<DidQueryDidResponse, RpcStatus>({
      path: `/metaearth/did/did`,
      method: "GET",
      query: query,
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryDidDocument
   * @request GET:/metaearth/did/didDocument
   */
  queryDidDocument = (query?: { did?: string }, params: RequestParams = {}) =>
    this.request<DidQueryDidDocumentResponse, RpcStatus>({
      path: `/metaearth/did/didDocument`,
      method: "GET",
      query: query,
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryDidInfos
   * @summary Queries a list of DidInfos items.
   * @request GET:/metaearth/did/did_infos
   */
  queryDidInfos = (
    query?: {
      "pagination.key"?: string;
      "pagination.offset"?: string;
      "pagination.limit"?: string;
      "pagination.count_total"?: boolean;
      "pagination.reverse"?: boolean;
    },
    params: RequestParams = {},
  ) =>
    this.request<DidQueryDidInfosResponse, RpcStatus>({
      path: `/metaearth/did/did_infos`,
      method: "GET",
      query: query,
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryService
   * @request GET:/metaearth/did/service
   */
  queryService = (query?: { sid?: string }, params: RequestParams = {}) =>
    this.request<DidQueryServiceResponse, RpcStatus>({
      path: `/metaearth/did/service`,
      method: "GET",
      query: query,
      format: "json",
      ...params,
    });
}
