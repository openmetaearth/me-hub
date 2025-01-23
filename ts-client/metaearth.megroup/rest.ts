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
 * GroupInfo represents the high-level on-chain information for a group.
 */
export interface MegroupGroupInfo {
  /**
   * id is the unique ID of the group.
   * @format uint64
   */
  id?: string;

  /** admin is the account address of the group's admin. */
  admin?: string;

  /** metadata is any arbitrary metadata to attached to the group. */
  metadata?: string;

  /**
   * version is used to track changes to a group's membership structure that
   * would break existing proposals. Whenever any members weight is changed,
   * or any member is added or removed this version is incremented and will
   * cause proposals based on older versions of this group to fail
   * @format uint64
   */
  version?: string;

  /** total_weight is the sum of the group members' weights. */
  total_weight?: string;

  /**
   * created_at is a timestamp specifying when a group was created.
   * @format date-time
   */
  created_at?: string;
  regionID?: string;
}

export interface MegroupMember {
  /** address is the member's account address. */
  address?: string;

  /** weight is the member's voting weight that should be greater than 0. */
  weight?: string;

  /** metadata is any arbitrary metadata attached to the member. */
  metadata?: string;

  /**
   * added_at is a timestamp specifying when a member was added.
   * @format date-time
   */
  added_at?: string;
}

export interface MegroupMsgCreateGroupResponse {
  /** @format uint64 */
  id?: string;
}

export type MegroupMsgDeleteGroupResponse = object;

export type MegroupMsgJoinGroupResponse = object;

export type MegroupMsgLeaveGroupResponse = object;

export type MegroupMsgUpdateGroupResponse = object;

export interface MegroupQueryAllGroupResponse {
  Group?: MegroupGroupInfo[];

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

export interface MegroupQueryGetGroupMemberCountResponse {
  /** @format uint64 */
  groupMemberCount?: string;
}

export interface MegroupQueryGetGroupMemberResponse {
  GroupMember?: MetaearthmegroupGroupMember;
}

export interface MegroupQueryGetGroupResponse {
  /** GroupInfo represents the high-level on-chain information for a group. */
  Group?: MegroupGroupInfo;
}

export interface MegroupQueryGroupAllMemberResponse {
  /** @format uint64 */
  groupID?: string;
  GroupMember?: MetaearthmegroupGroupMember[];

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
 * QueryParamsResponse is response type for the Query/Params RPC method.
 */
export interface MegroupQueryParamsResponse {
  /** params holds all the parameters of this module. */
  params?: MetaearthmegroupParams;
}

export interface MetaearthmegroupGroupMember {
  /** @format uint64 */
  group_id?: string;
  member?: MegroupMember;
}

/**
 * Params defines the parameters for the module.
 */
export type MetaearthmegroupParams = object;

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
 * @title metaearth/megroup/genesis.proto
 * @version version not set
 */
export class Api<SecurityDataType extends unknown> extends HttpClient<SecurityDataType> {
  /**
   * No description
   *
   * @tags Query
   * @name QueryGroupAll
   * @request GET:/me-hub/megroup/all_group
   */
  queryGroupAll = (
    query?: {
      "pagination.key"?: string;
      "pagination.offset"?: string;
      "pagination.limit"?: string;
      "pagination.count_total"?: boolean;
      "pagination.reverse"?: boolean;
    },
    params: RequestParams = {},
  ) =>
    this.request<MegroupQueryAllGroupResponse, RpcStatus>({
      path: `/me-hub/megroup/all_group`,
      method: "GET",
      query: query,
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryGroupByMember
   * @request GET:/me-hub/megroup/group/{address}
   */
  queryGroupByMember = (address: string, params: RequestParams = {}) =>
    this.request<MegroupQueryGetGroupResponse, RpcStatus>({
      path: `/me-hub/megroup/group/${address}`,
      method: "GET",
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryGroup
   * @summary Queries a list of Group items.
   * @request GET:/me-hub/megroup/group/{id}
   */
  queryGroup = (id: string, params: RequestParams = {}) =>
    this.request<MegroupQueryGetGroupResponse, RpcStatus>({
      path: `/me-hub/megroup/group/${id}`,
      method: "GET",
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryGroupMemberAll
   * @request GET:/me-hub/megroup/group_all_member/{groupID}
   */
  queryGroupMemberAll = (
    groupId: string,
    query?: {
      "pagination.key"?: string;
      "pagination.offset"?: string;
      "pagination.limit"?: string;
      "pagination.count_total"?: boolean;
      "pagination.reverse"?: boolean;
    },
    params: RequestParams = {},
  ) =>
    this.request<MegroupQueryGroupAllMemberResponse, RpcStatus>({
      path: `/me-hub/megroup/group_all_member/${groupId}`,
      method: "GET",
      query: query,
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryGroupMember
   * @summary Queries a list of GroupMember items.
   * @request GET:/me-hub/megroup/group_member/{address}
   */
  queryGroupMember = (address: string, params: RequestParams = {}) =>
    this.request<MegroupQueryGetGroupMemberResponse, RpcStatus>({
      path: `/me-hub/megroup/group_member/${address}`,
      method: "GET",
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryGroupMemberCount
   * @summary Queries a list of GroupMemberCount items.
   * @request GET:/me-hub/megroup/group_member_count/{groupId}
   */
  queryGroupMemberCount = (groupId: string, params: RequestParams = {}) =>
    this.request<MegroupQueryGetGroupMemberCountResponse, RpcStatus>({
      path: `/me-hub/megroup/group_member_count/${groupId}`,
      method: "GET",
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryParams
   * @summary Parameters queries the parameters of the module.
   * @request GET:/me-hub/megroup/params
   */
  queryParams = (params: RequestParams = {}) =>
    this.request<MegroupQueryParamsResponse, RpcStatus>({
      path: `/me-hub/megroup/params`,
      method: "GET",
      format: "json",
      ...params,
    });
}
