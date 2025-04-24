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

export interface Evmv1Params {
  /**
   * evm_denom represents the token denomination used to run the EVM state
   * transitions.
   */
  evm_denom?: string;

  /** enable_create toggles state transitions that use the vm.Create function */
  enable_create?: boolean;

  /** enable_call toggles state transitions that use the vm.Call function */
  enable_call?: boolean;

  /** extra_eips defines the additional EIPs for the vm.Config */
  extra_eips?: string[];

  /**
   * chain_config defines the EVM chain configuration parameters
   * ChainConfig defines the Ethereum ChainConfig parameters using *sdk.Int values
   * instead of *big.Int.
   */
  chain_config?: V1ChainConfig;

  /**
   * allow_unprotected_txs defines if replay-protected (i.e non EIP155
   * signed) transactions can be executed on the state machine.
   */
  allow_unprotected_txs?: boolean;
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
* ChainConfig defines the Ethereum ChainConfig parameters using *sdk.Int values
instead of *big.Int.
*/
export interface V1ChainConfig {
  /** homestead_block switch (nil no fork, 0 = already homestead) */
  homestead_block?: string;

  /** dao_fork_block corresponds to TheDAO hard-fork switch block (nil no fork) */
  dao_fork_block?: string;

  /** dao_fork_support defines whether the nodes supports or opposes the DAO hard-fork */
  dao_fork_support?: boolean;

  /**
   * eip150_block: EIP150 implements the Gas price changes
   * (https://github.com/ethereum/EIPs/issues/150) EIP150 HF block (nil no fork)
   */
  eip150_block?: string;

  /** eip150_hash: EIP150 HF hash (needed for header only clients as only gas pricing changed) */
  eip150_hash?: string;

  /** eip155_block: EIP155Block HF block */
  eip155_block?: string;

  /** eip158_block: EIP158 HF block */
  eip158_block?: string;

  /** byzantium_block: Byzantium switch block (nil no fork, 0 = already on byzantium) */
  byzantium_block?: string;

  /** constantinople_block: Constantinople switch block (nil no fork, 0 = already activated) */
  constantinople_block?: string;

  /** petersburg_block: Petersburg switch block (nil same as Constantinople) */
  petersburg_block?: string;

  /** istanbul_block: Istanbul switch block (nil no fork, 0 = already on istanbul) */
  istanbul_block?: string;

  /** muir_glacier_block: Eip-2384 (bomb delay) switch block (nil no fork, 0 = already activated) */
  muir_glacier_block?: string;

  /** berlin_block: Berlin switch block (nil = no fork, 0 = already on berlin) */
  berlin_block?: string;

  /** london_block: London switch block (nil = no fork, 0 = already on london) */
  london_block?: string;

  /** arrow_glacier_block: Eip-4345 (bomb delay) switch block (nil = no fork, 0 = already activated) */
  arrow_glacier_block?: string;

  /** gray_glacier_block: EIP-5133 (bomb delay) switch block (nil = no fork, 0 = already activated) */
  gray_glacier_block?: string;

  /** merge_netsplit_block: Virtual fork after The Merge to use as a network splitter */
  merge_netsplit_block?: string;

  /** shanghai_block switch block (nil = no fork, 0 = already on shanghai) */
  shanghai_block?: string;

  /** cancun_block switch block (nil = no fork, 0 = already on cancun) */
  cancun_block?: string;
}

export interface V1EstimateGasResponse {
  /**
   * gas returns the estimated gas
   * @format uint64
   */
  gas?: string;
}

/**
* Log represents an protobuf compatible Ethereum Log that defines a contract
log event. These events are generated by the LOG opcode and stored/indexed by
the node.

NOTE: address, topics and data are consensus fields. The rest of the fields
are derived, i.e. filled in by the nodes, but not secured by consensus.
*/
export interface V1Log {
  /** address of the contract that generated the event */
  address?: string;

  /** topics is a list of topics provided by the contract. */
  topics?: string[];

  /**
   * data which is supplied by the contract, usually ABI-encoded
   * @format byte
   */
  data?: string;

  /**
   * block_number of the block in which the transaction was included
   * @format uint64
   */
  block_number?: string;

  /** tx_hash is the transaction hash */
  tx_hash?: string;

  /**
   * tx_index of the transaction in the block
   * @format uint64
   */
  tx_index?: string;

  /** block_hash of the block in which the transaction was included */
  block_hash?: string;

  /**
   * index of the log in the block
   * @format uint64
   */
  index?: string;

  /**
   * removed is true if this log was reverted due to a chain
   * reorganisation. You must pay attention to this field if you receive logs
   * through a filter query.
   */
  removed?: boolean;
}

/**
 * MsgEthereumTx encapsulates an Ethereum transaction as an SDK message.
 */
export interface V1MsgEthereumTx {
  /**
   * data is inner transaction data of the Ethereum transaction
   * `Any` contains an arbitrary serialized protocol buffer message along with a
   * URL that describes the type of the serialized message.
   *
   * Protobuf library provides support to pack/unpack Any values in the form
   * of utility functions or additional generated methods of the Any type.
   * Example 1: Pack and unpack a message in C++.
   *     Foo foo = ...;
   *     Any any;
   *     any.PackFrom(foo);
   *     ...
   *     if (any.UnpackTo(&foo)) {
   *       ...
   *     }
   * Example 2: Pack and unpack a message in Java.
   *     Any any = Any.pack(foo);
   *     if (any.is(Foo.class)) {
   *       foo = any.unpack(Foo.class);
   *  Example 3: Pack and unpack a message in Python.
   *     foo = Foo(...)
   *     any = Any()
   *     any.Pack(foo)
   *     if any.Is(Foo.DESCRIPTOR):
   *       any.Unpack(foo)
   *  Example 4: Pack and unpack a message in Go
   *      foo := &pb.Foo{...}
   *      any, err := anypb.New(foo)
   *      if err != nil {
   *        ...
   *      }
   *      ...
   *      foo := &pb.Foo{}
   *      if err := any.UnmarshalTo(foo); err != nil {
   * The pack methods provided by protobuf library will by default use
   * 'type.googleapis.com/full.type.name' as the type URL and the unpack
   * methods only use the fully qualified type name after the last '/'
   * in the type URL, for example "foo.bar.com/x/y.z" will yield type
   * name "y.z".
   * JSON
   * ====
   * The JSON representation of an `Any` value uses the regular
   * representation of the deserialized, embedded message, with an
   * additional field `@type` which contains the type URL. Example:
   *     package google.profile;
   *     message Person {
   *       string first_name = 1;
   *       string last_name = 2;
   *     {
   *       "@type": "type.googleapis.com/google.profile.Person",
   *       "firstName": <string>,
   *       "lastName": <string>
   * If the embedded message type is well-known and has a custom JSON
   * representation, that representation will be embedded adding a field
   * `value` which holds the custom JSON in addition to the `@type`
   * field. Example (for message [google.protobuf.Duration][]):
   *       "@type": "type.googleapis.com/google.protobuf.Duration",
   *       "value": "1.212s"
   */
  data?: ProtobufAny;

  /**
   * size is the encoded storage size of the transaction (DEPRECATED)
   * @format double
   */
  size?: number;

  /** hash of the transaction in hex format */
  hash?: string;

  /**
   * from is the ethereum signer address in hex format. This address value is checked
   * against the address derived from the signature (V, R, S) using the
   * secp256k1 elliptic curve
   */
  from?: string;
}

/**
 * MsgEthereumTxResponse defines the Msg/EthereumTx response type.
 */
export interface V1MsgEthereumTxResponse {
  /**
   * hash of the ethereum transaction in hex format. This hash differs from the
   * Tendermint sha256 hash of the transaction bytes. See
   * https://github.com/tendermint/tendermint/issues/6539 for reference
   */
  hash?: string;

  /**
   * logs contains the transaction hash and the proto-compatible ethereum
   * logs.
   */
  logs?: V1Log[];

  /**
   * ret is the returned data from evm function (result or data supplied with revert
   * opcode)
   * @format byte
   */
  ret?: string;

  /** vm_error is the error returned by vm execution */
  vm_error?: string;

  /**
   * gas_used specifies how much gas was consumed by the transaction
   * @format uint64
   */
  gas_used?: string;
}

/**
* MsgUpdateParamsResponse defines the response structure for executing a
MsgUpdateParams message.
*/
export type V1MsgUpdateParamsResponse = object;

/**
 * QueryAccountResponse is the response type for the Query/Account RPC method.
 */
export interface V1QueryAccountResponse {
  /** balance is the balance of the EVM denomination. */
  balance?: string;

  /** code_hash is the hex-formatted code bytes from the EOA. */
  code_hash?: string;

  /**
   * nonce is the account's sequence number.
   * @format uint64
   */
  nonce?: string;
}

/**
 * QueryBalanceResponse is the response type for the Query/Balance RPC method.
 */
export interface V1QueryBalanceResponse {
  /** balance is the balance of the EVM denomination. */
  balance?: string;
}

/**
 * QueryBaseFeeResponse returns the EIP1559 base fee.
 */
export interface V1QueryBaseFeeResponse {
  /** base_fee is the EIP1559 base fee */
  base_fee?: string;
}

/**
* QueryCodeResponse is the response type for the Query/Code RPC
method.
*/
export interface V1QueryCodeResponse {
  /**
   * code represents the code bytes from an ethereum address.
   * @format byte
   */
  code?: string;
}

/**
* QueryCosmosAccountResponse is the response type for the Query/CosmosAccount
RPC method.
*/
export interface V1QueryCosmosAccountResponse {
  /** cosmos_address is the cosmos address of the account. */
  cosmos_address?: string;

  /**
   * sequence is the account's sequence number.
   * @format uint64
   */
  sequence?: string;

  /**
   * account_number is the account number
   * @format uint64
   */
  account_number?: string;
}

/**
 * QueryParamsResponse defines the response type for querying x/evm parameters.
 */
export interface V1QueryParamsResponse {
  /** params define the evm module parameters. */
  params?: Evmv1Params;
}

/**
* QueryStorageResponse is the response type for the Query/Storage RPC
method.
*/
export interface V1QueryStorageResponse {
  /** value defines the storage state value hash associated with the given key. */
  value?: string;
}

export interface V1QueryTraceBlockResponse {
  /**
   * data is the response serialized in bytes
   * @format byte
   */
  data?: string;
}

export interface V1QueryTraceTxResponse {
  /**
   * data is the response serialized in bytes
   * @format byte
   */
  data?: string;
}

/**
* QueryValidatorAccountResponse is the response type for the
Query/ValidatorAccount RPC method.
*/
export interface V1QueryValidatorAccountResponse {
  /** account_address is the cosmos address of the account in bech32 format. */
  account_address?: string;

  /**
   * sequence is the account's sequence number.
   * @format uint64
   */
  sequence?: string;

  /**
   * account_number is the account number
   * @format uint64
   */
  account_number?: string;
}

export interface V1QueryVirtualFrontierBankContractByDenomResponse {
  /** VFBCPair is a pair of a native Cosmos token (base) with an ERC20 virtual frontier bank contract address. */
  pair?: V1VFBCPair;
}

export interface V1QueryVirtualFrontierBankContractsResponse {
  /** pairs is list of available virtual frontier bank contracts */
  pairs?: V1VFBCPair[];

  /** pagination defines the pagination in the response. */
  pagination?: V1Beta1PageResponse;
}

export interface V1QueryVirtualFrontierContractByAddressResponse {
  virtual_frontier_contract_json?: string;
}

export interface V1QueryVirtualFrontierContractsResponse {
  virtual_frontier_contracts_json?: string[];

  /** pagination defines the pagination in the response. */
  pagination?: V1Beta1PageResponse;
}

/**
 * TraceConfig holds extra parameters to trace functions.
 */
export interface V1TraceConfig {
  /** tracer is a custom javascript tracer */
  tracer?: string;

  /**
   * timeout overrides the default timeout of 5 seconds for JavaScript-based tracing
   * calls
   */
  timeout?: string;

  /**
   * reexec defines the number of blocks the tracer is willing to go back
   * @format uint64
   */
  reexec?: string;

  /** disable_stack switches stack capture */
  disable_stack?: boolean;

  /** disable_storage switches storage capture */
  disable_storage?: boolean;

  /** debug can be used to print output during capture end */
  debug?: boolean;

  /**
   * limit defines the maximum length of output, but zero means unlimited
   * @format int32
   */
  limit?: number;

  /**
   * overrides can be used to execute a trace using future fork rules
   * ChainConfig defines the Ethereum ChainConfig parameters using *sdk.Int values
   * instead of *big.Int.
   */
  overrides?: V1ChainConfig;

  /** enable_memory switches memory capture */
  enable_memory?: boolean;

  /** enable_return_data switches the capture of return data */
  enable_return_data?: boolean;

  /** tracer_json_config configures the tracer using a JSON string */
  tracer_json_config?: string;
}

/**
 * VFBCPair is a pair of a native Cosmos token (base) with an ERC20 virtual frontier bank contract address.
 */
export interface V1VFBCPair {
  /** contract_address is the hex address of the contract */
  contract_address?: string;

  /** min_denom defines the cosmos base denomination */
  min_denom?: string;

  /** enabled is status of the contract, whether it is enabled or not */
  enabled?: boolean;
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
 * @title ethermint/evm/v1/events.proto
 * @version version not set
 */
export class Api<SecurityDataType extends unknown> extends HttpClient<SecurityDataType> {
  /**
   * No description
   *
   * @tags Query
   * @name QueryAccount
   * @summary Account queries an Ethereum account.
   * @request GET:/ethermint/evm/v1/account/{address}
   */
  queryAccount = (address: string, params: RequestParams = {}) =>
    this.request<V1QueryAccountResponse, RpcStatus>({
      path: `/ethermint/evm/v1/account/${address}`,
      method: "GET",
      format: "json",
      ...params,
    });

  /**
 * No description
 * 
 * @tags Query
 * @name QueryBalance
 * @summary Balance queries the balance of a the EVM denomination for a single
EthAccount.
 * @request GET:/ethermint/evm/v1/balances/{address}
 */
  queryBalance = (address: string, params: RequestParams = {}) =>
    this.request<V1QueryBalanceResponse, RpcStatus>({
      path: `/ethermint/evm/v1/balances/${address}`,
      method: "GET",
      format: "json",
      ...params,
    });

  /**
 * No description
 * 
 * @tags Query
 * @name QueryBaseFee
 * @summary BaseFee queries the base fee of the parent block of the current block,
it's similar to feemarket module's method, but also checks london hardfork status.
 * @request GET:/ethermint/evm/v1/base_fee
 */
  queryBaseFee = (params: RequestParams = {}) =>
    this.request<V1QueryBaseFeeResponse, RpcStatus>({
      path: `/ethermint/evm/v1/base_fee`,
      method: "GET",
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryCode
   * @summary Code queries the balance of all coins for a single account.
   * @request GET:/ethermint/evm/v1/codes/{address}
   */
  queryCode = (address: string, params: RequestParams = {}) =>
    this.request<V1QueryCodeResponse, RpcStatus>({
      path: `/ethermint/evm/v1/codes/${address}`,
      method: "GET",
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryCosmosAccount
   * @summary CosmosAccount queries an Ethereum account's Cosmos Address.
   * @request GET:/ethermint/evm/v1/cosmos_account/{address}
   */
  queryCosmosAccount = (address: string, params: RequestParams = {}) =>
    this.request<V1QueryCosmosAccountResponse, RpcStatus>({
      path: `/ethermint/evm/v1/cosmos_account/${address}`,
      method: "GET",
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryEstimateGas
   * @summary EstimateGas implements the `eth_estimateGas` rpc api
   * @request GET:/ethermint/evm/v1/estimate_gas
   */
  queryEstimateGas = (
    query?: { args?: string; gas_cap?: string; proposer_address?: string; chain_id?: string },
    params: RequestParams = {},
  ) =>
    this.request<V1EstimateGasResponse, RpcStatus>({
      path: `/ethermint/evm/v1/estimate_gas`,
      method: "GET",
      query: query,
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryEthCall
   * @summary EthCall implements the `eth_call` rpc api
   * @request GET:/ethermint/evm/v1/eth_call
   */
  queryEthCall = (
    query?: { args?: string; gas_cap?: string; proposer_address?: string; chain_id?: string },
    params: RequestParams = {},
  ) =>
    this.request<V1MsgEthereumTxResponse, RpcStatus>({
      path: `/ethermint/evm/v1/eth_call`,
      method: "GET",
      query: query,
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Msg
   * @name MsgEthereumTx
   * @summary EthereumTx defines a method submitting Ethereum transactions.
   * @request POST:/ethermint/evm/v1/ethereum_tx
   */
  msgEthereumTx = (
    query?: { "data.type_url"?: string; "data.value"?: string; size?: number; hash?: string; from?: string },
    params: RequestParams = {},
  ) =>
    this.request<V1MsgEthereumTxResponse, RpcStatus>({
      path: `/ethermint/evm/v1/ethereum_tx`,
      method: "POST",
      query: query,
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryParams
   * @summary Params queries the parameters of x/evm module.
   * @request GET:/ethermint/evm/v1/params
   */
  queryParams = (params: RequestParams = {}) =>
    this.request<V1QueryParamsResponse, RpcStatus>({
      path: `/ethermint/evm/v1/params`,
      method: "GET",
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryStorage
   * @summary Storage queries the balance of all coins for a single account.
   * @request GET:/ethermint/evm/v1/storage/{address}/{key}
   */
  queryStorage = (address: string, key: string, params: RequestParams = {}) =>
    this.request<V1QueryStorageResponse, RpcStatus>({
      path: `/ethermint/evm/v1/storage/${address}/${key}`,
      method: "GET",
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryTraceBlock
   * @summary TraceBlock implements the `debug_traceBlockByNumber` and `debug_traceBlockByHash` rpc api
   * @request GET:/ethermint/evm/v1/trace_block
   */
  queryTraceBlock = (
    query?: {
      "trace_config.tracer"?: string;
      "trace_config.timeout"?: string;
      "trace_config.reexec"?: string;
      "trace_config.disable_stack"?: boolean;
      "trace_config.disable_storage"?: boolean;
      "trace_config.debug"?: boolean;
      "trace_config.limit"?: number;
      "trace_config.overrides.homestead_block"?: string;
      "trace_config.overrides.dao_fork_block"?: string;
      "trace_config.overrides.dao_fork_support"?: boolean;
      "trace_config.overrides.eip150_block"?: string;
      "trace_config.overrides.eip150_hash"?: string;
      "trace_config.overrides.eip155_block"?: string;
      "trace_config.overrides.eip158_block"?: string;
      "trace_config.overrides.byzantium_block"?: string;
      "trace_config.overrides.constantinople_block"?: string;
      "trace_config.overrides.petersburg_block"?: string;
      "trace_config.overrides.istanbul_block"?: string;
      "trace_config.overrides.muir_glacier_block"?: string;
      "trace_config.overrides.berlin_block"?: string;
      "trace_config.overrides.london_block"?: string;
      "trace_config.overrides.arrow_glacier_block"?: string;
      "trace_config.overrides.gray_glacier_block"?: string;
      "trace_config.overrides.merge_netsplit_block"?: string;
      "trace_config.overrides.shanghai_block"?: string;
      "trace_config.overrides.cancun_block"?: string;
      "trace_config.enable_memory"?: boolean;
      "trace_config.enable_return_data"?: boolean;
      "trace_config.tracer_json_config"?: string;
      block_number?: string;
      block_hash?: string;
      block_time?: string;
      proposer_address?: string;
      chain_id?: string;
    },
    params: RequestParams = {},
  ) =>
    this.request<V1QueryTraceBlockResponse, RpcStatus>({
      path: `/ethermint/evm/v1/trace_block`,
      method: "GET",
      query: query,
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryTraceTx
   * @summary TraceTx implements the `debug_traceTransaction` rpc api
   * @request GET:/ethermint/evm/v1/trace_tx
   */
  queryTraceTx = (
    query?: {
      "msg.data.type_url"?: string;
      "msg.data.value"?: string;
      "msg.size"?: number;
      "msg.hash"?: string;
      "msg.from"?: string;
      "trace_config.tracer"?: string;
      "trace_config.timeout"?: string;
      "trace_config.reexec"?: string;
      "trace_config.disable_stack"?: boolean;
      "trace_config.disable_storage"?: boolean;
      "trace_config.debug"?: boolean;
      "trace_config.limit"?: number;
      "trace_config.overrides.homestead_block"?: string;
      "trace_config.overrides.dao_fork_block"?: string;
      "trace_config.overrides.dao_fork_support"?: boolean;
      "trace_config.overrides.eip150_block"?: string;
      "trace_config.overrides.eip150_hash"?: string;
      "trace_config.overrides.eip155_block"?: string;
      "trace_config.overrides.eip158_block"?: string;
      "trace_config.overrides.byzantium_block"?: string;
      "trace_config.overrides.constantinople_block"?: string;
      "trace_config.overrides.petersburg_block"?: string;
      "trace_config.overrides.istanbul_block"?: string;
      "trace_config.overrides.muir_glacier_block"?: string;
      "trace_config.overrides.berlin_block"?: string;
      "trace_config.overrides.london_block"?: string;
      "trace_config.overrides.arrow_glacier_block"?: string;
      "trace_config.overrides.gray_glacier_block"?: string;
      "trace_config.overrides.merge_netsplit_block"?: string;
      "trace_config.overrides.shanghai_block"?: string;
      "trace_config.overrides.cancun_block"?: string;
      "trace_config.enable_memory"?: boolean;
      "trace_config.enable_return_data"?: boolean;
      "trace_config.tracer_json_config"?: string;
      block_number?: string;
      block_hash?: string;
      block_time?: string;
      proposer_address?: string;
      chain_id?: string;
    },
    params: RequestParams = {},
  ) =>
    this.request<V1QueryTraceTxResponse, RpcStatus>({
      path: `/ethermint/evm/v1/trace_tx`,
      method: "GET",
      query: query,
      format: "json",
      ...params,
    });

  /**
 * No description
 * 
 * @tags Query
 * @name QueryValidatorAccount
 * @summary ValidatorAccount queries an Ethereum account's from a validator consensus
Address.
 * @request GET:/ethermint/evm/v1/validator_account/{cons_address}
 */
  queryValidatorAccount = (consAddress: string, params: RequestParams = {}) =>
    this.request<V1QueryValidatorAccountResponse, RpcStatus>({
      path: `/ethermint/evm/v1/validator_account/${consAddress}`,
      method: "GET",
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryVirtualFrontierBankContractByDenom
   * @summary VirtualFrontierBankContractByDenom returns the bank contract of the corresponding min denom
   * @request GET:/ethermint/evm/v1/virtual_frontier_bank_contract/by_denom/{min_denom}
   */
  queryVirtualFrontierBankContractByDenom = (minDenom: string, params: RequestParams = {}) =>
    this.request<V1QueryVirtualFrontierBankContractByDenomResponse, RpcStatus>({
      path: `/ethermint/evm/v1/virtual_frontier_bank_contract/by_denom/${minDenom}`,
      method: "GET",
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryListVirtualFrontierBankContracts
   * @summary ListVirtualFrontierBankContracts returns the list of bank contract
   * @request GET:/ethermint/evm/v1/virtual_frontier_bank_contracts
   */
  queryListVirtualFrontierBankContracts = (
    query?: {
      "pagination.key"?: string;
      "pagination.offset"?: string;
      "pagination.limit"?: string;
      "pagination.count_total"?: boolean;
      "pagination.reverse"?: boolean;
    },
    params: RequestParams = {},
  ) =>
    this.request<V1QueryVirtualFrontierBankContractsResponse, RpcStatus>({
      path: `/ethermint/evm/v1/virtual_frontier_bank_contracts`,
      method: "GET",
      query: query,
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryListVirtualFrontierContracts
   * @summary ListVirtualFrontierContracts returns Json array of information of virtual frontier contracts.
   * @request GET:/ethermint/evm/v1/virtual_frontier_contracts
   */
  queryListVirtualFrontierContracts = (
    query?: {
      "pagination.key"?: string;
      "pagination.offset"?: string;
      "pagination.limit"?: string;
      "pagination.count_total"?: boolean;
      "pagination.reverse"?: boolean;
    },
    params: RequestParams = {},
  ) =>
    this.request<V1QueryVirtualFrontierContractsResponse, RpcStatus>({
      path: `/ethermint/evm/v1/virtual_frontier_contracts`,
      method: "GET",
      query: query,
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryVirtualFrontierContractByAddress
   * @summary VirtualFrontierContractByAddress returns Json of information of the corresponding virtual frontier contract by address
   * @request GET:/ethermint/evm/v1/virtual_frontier_contracts/{address}
   */
  queryVirtualFrontierContractByAddress = (address: string, params: RequestParams = {}) =>
    this.request<V1QueryVirtualFrontierContractByAddressResponse, RpcStatus>({
      path: `/ethermint/evm/v1/virtual_frontier_contracts/${address}`,
      method: "GET",
      format: "json",
      ...params,
    });
}
