/* eslint-disable */
import Long from "long";
import _m0 from "protobufjs/minimal";
import { PageRequest, PageResponse } from "../../cosmos/base/query/v1beta1/pagination";
import { Coin, DecCoin } from "../../cosmos/base/v1beta1/coin";
import { QueryDelegationRequest, QueryDelegationResponse } from "../../cosmos/staking/v1beta1/query";
import { Delegation } from "../../cosmos/staking/v1beta1/staking";
import {
  FixedDeposit,
  FixedDepositCfg,
  FixedDepositState,
  fixedDepositStateFromJSON,
  fixedDepositStateToJSON,
  RegionAllFixedDepositCfg,
} from "./fixed_deposit";
import { Record, ReviewRecord } from "./record";
import { Region } from "./region";
import { Stake } from "./stake";

export const protobufPackage = "metaearth.wstaking";

export interface QueryRegionRequest {
  regionId: string;
}

export interface QueryRegionResponse {
  region: Region | undefined;
}

export interface QueryAllRegionRequest {
  pagination: PageRequest | undefined;
}

export interface QueryAllRegionResponse {
  region: Region[];
  pagination: PageResponse | undefined;
}

export interface QueryDelegationRewardsRequest {
  /** delegator_address defines the delegator address to query for. */
  delegatorAddress: string;
  /** validator_address defines the validator address to query for. */
  validatorAddress: string;
}

export interface QueryDelegationRewardsResponse {
  /** rewards defines the rewards accrued by a delegation. */
  rewards: DecCoin[];
}

export interface QueryFixedDepositByAcctRequest {
  /** cosmos.base.query.v1beta1.PageRequest pagination = 1; */
  account: string;
  queryType: FixedDepositState;
}

export interface QueryFixedDepositByAcctResponse {
  /** cosmos.base.query.v1beta1.PageResponse pagination = 2; */
  FixedDeposit: FixedDeposit[];
}

export interface QueryFixedDepositByRegionRequest {
  regionId: string;
  pagination: PageRequest | undefined;
  queryType: FixedDepositState;
}

export interface QueryFixedDepositByRegionResponse {
  FixedDeposit: FixedDeposit[];
  pagination: PageResponse | undefined;
}

export interface QueryGetFixedDepositRequest {
  address: string;
  id: number;
}

export interface QueryGetFixedDepositResponse {
  FixedDeposit: FixedDeposit | undefined;
}

export interface QueryAllFixedDepositRequest {
  pagination: PageRequest | undefined;
}

export interface QueryAllFixedDepositResponse {
  FixedDeposit: FixedDeposit[];
  pagination: PageResponse | undefined;
}

export interface QueryFixedDepositCfgRequest {
  regionIds: string[];
}

export interface QueryFixedDepositCfgResponse {
  RegionFixedDepositCfgs: RegionAllFixedDepositCfg[];
}

export interface QueryFixedDepositCfgByTermRequest {
  regionId: string;
  term: number;
}

export interface QueryFixedDepositCfgByTermResponse {
  FixedDepositCfg: FixedDepositCfg | undefined;
}

export interface QueryFixedDepositTotalAmountRequest {
}

export interface QueryFixedDepositTotalAmountResponse {
  amount: Coin | undefined;
}

export interface QueryFixedDepositAmountByMeidRequest {
  account: string;
}

export interface QueryFixedDepositAmountByMeidResponse {
  amount: Coin | undefined;
}

export interface QueryAllDelegationsRequest {
  pagination: PageRequest | undefined;
}

export interface QueryAllDelegationsResponse {
  delegations: Delegation[];
  pagination: PageResponse | undefined;
}

export interface QueryStakesRequest {
  pagination: PageRequest | undefined;
}

export interface QueryStakesResponse {
  stakes: Stake[];
  pagination: PageResponse | undefined;
}

export interface QueryAllRecords {
  pagination: PageRequest | undefined;
}

export interface QueryAllRecordsResponse {
  records: Record[];
  pagination: PageResponse | undefined;
}

export interface QueryRecordsByAddress {
  account: string;
}

export interface QueryRecordsByAddressResponse {
  /** cosmos.base.query.v1beta1.PageResponse pagination = 2; */
  records: Record[];
}

export interface QueryReviewRecordByNumber {
  actionNumber: string;
}

export interface QueryReviewRecordByNumberResponse {
  /** cosmos.base.query.v1beta1.PageResponse pagination = 2; */
  reviewRecord: ReviewRecord | undefined;
}

function createBaseQueryRegionRequest(): QueryRegionRequest {
  return { regionId: "" };
}

export const QueryRegionRequest = {
  encode(message: QueryRegionRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.regionId !== "") {
      writer.uint32(10).string(message.regionId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryRegionRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryRegionRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.regionId = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryRegionRequest {
    return { regionId: isSet(object.regionId) ? String(object.regionId) : "" };
  },

  toJSON(message: QueryRegionRequest): unknown {
    const obj: any = {};
    message.regionId !== undefined && (obj.regionId = message.regionId);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryRegionRequest>, I>>(object: I): QueryRegionRequest {
    const message = createBaseQueryRegionRequest();
    message.regionId = object.regionId ?? "";
    return message;
  },
};

function createBaseQueryRegionResponse(): QueryRegionResponse {
  return { region: undefined };
}

export const QueryRegionResponse = {
  encode(message: QueryRegionResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.region !== undefined) {
      Region.encode(message.region, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryRegionResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryRegionResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.region = Region.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryRegionResponse {
    return { region: isSet(object.region) ? Region.fromJSON(object.region) : undefined };
  },

  toJSON(message: QueryRegionResponse): unknown {
    const obj: any = {};
    message.region !== undefined && (obj.region = message.region ? Region.toJSON(message.region) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryRegionResponse>, I>>(object: I): QueryRegionResponse {
    const message = createBaseQueryRegionResponse();
    message.region = (object.region !== undefined && object.region !== null)
      ? Region.fromPartial(object.region)
      : undefined;
    return message;
  },
};

function createBaseQueryAllRegionRequest(): QueryAllRegionRequest {
  return { pagination: undefined };
}

export const QueryAllRegionRequest = {
  encode(message: QueryAllRegionRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryAllRegionRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllRegionRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.pagination = PageRequest.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryAllRegionRequest {
    return { pagination: isSet(object.pagination) ? PageRequest.fromJSON(object.pagination) : undefined };
  },

  toJSON(message: QueryAllRegionRequest): unknown {
    const obj: any = {};
    message.pagination !== undefined
      && (obj.pagination = message.pagination ? PageRequest.toJSON(message.pagination) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryAllRegionRequest>, I>>(object: I): QueryAllRegionRequest {
    const message = createBaseQueryAllRegionRequest();
    message.pagination = (object.pagination !== undefined && object.pagination !== null)
      ? PageRequest.fromPartial(object.pagination)
      : undefined;
    return message;
  },
};

function createBaseQueryAllRegionResponse(): QueryAllRegionResponse {
  return { region: [], pagination: undefined };
}

export const QueryAllRegionResponse = {
  encode(message: QueryAllRegionResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.region) {
      Region.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.pagination !== undefined) {
      PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryAllRegionResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllRegionResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.region.push(Region.decode(reader, reader.uint32()));
          break;
        case 2:
          message.pagination = PageResponse.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryAllRegionResponse {
    return {
      region: Array.isArray(object?.region) ? object.region.map((e: any) => Region.fromJSON(e)) : [],
      pagination: isSet(object.pagination) ? PageResponse.fromJSON(object.pagination) : undefined,
    };
  },

  toJSON(message: QueryAllRegionResponse): unknown {
    const obj: any = {};
    if (message.region) {
      obj.region = message.region.map((e) => e ? Region.toJSON(e) : undefined);
    } else {
      obj.region = [];
    }
    message.pagination !== undefined
      && (obj.pagination = message.pagination ? PageResponse.toJSON(message.pagination) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryAllRegionResponse>, I>>(object: I): QueryAllRegionResponse {
    const message = createBaseQueryAllRegionResponse();
    message.region = object.region?.map((e) => Region.fromPartial(e)) || [];
    message.pagination = (object.pagination !== undefined && object.pagination !== null)
      ? PageResponse.fromPartial(object.pagination)
      : undefined;
    return message;
  },
};

function createBaseQueryDelegationRewardsRequest(): QueryDelegationRewardsRequest {
  return { delegatorAddress: "", validatorAddress: "" };
}

export const QueryDelegationRewardsRequest = {
  encode(message: QueryDelegationRewardsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.delegatorAddress !== "") {
      writer.uint32(10).string(message.delegatorAddress);
    }
    if (message.validatorAddress !== "") {
      writer.uint32(18).string(message.validatorAddress);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryDelegationRewardsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryDelegationRewardsRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.delegatorAddress = reader.string();
          break;
        case 2:
          message.validatorAddress = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryDelegationRewardsRequest {
    return {
      delegatorAddress: isSet(object.delegatorAddress) ? String(object.delegatorAddress) : "",
      validatorAddress: isSet(object.validatorAddress) ? String(object.validatorAddress) : "",
    };
  },

  toJSON(message: QueryDelegationRewardsRequest): unknown {
    const obj: any = {};
    message.delegatorAddress !== undefined && (obj.delegatorAddress = message.delegatorAddress);
    message.validatorAddress !== undefined && (obj.validatorAddress = message.validatorAddress);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryDelegationRewardsRequest>, I>>(
    object: I,
  ): QueryDelegationRewardsRequest {
    const message = createBaseQueryDelegationRewardsRequest();
    message.delegatorAddress = object.delegatorAddress ?? "";
    message.validatorAddress = object.validatorAddress ?? "";
    return message;
  },
};

function createBaseQueryDelegationRewardsResponse(): QueryDelegationRewardsResponse {
  return { rewards: [] };
}

export const QueryDelegationRewardsResponse = {
  encode(message: QueryDelegationRewardsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.rewards) {
      DecCoin.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryDelegationRewardsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryDelegationRewardsResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.rewards.push(DecCoin.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryDelegationRewardsResponse {
    return { rewards: Array.isArray(object?.rewards) ? object.rewards.map((e: any) => DecCoin.fromJSON(e)) : [] };
  },

  toJSON(message: QueryDelegationRewardsResponse): unknown {
    const obj: any = {};
    if (message.rewards) {
      obj.rewards = message.rewards.map((e) => e ? DecCoin.toJSON(e) : undefined);
    } else {
      obj.rewards = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryDelegationRewardsResponse>, I>>(
    object: I,
  ): QueryDelegationRewardsResponse {
    const message = createBaseQueryDelegationRewardsResponse();
    message.rewards = object.rewards?.map((e) => DecCoin.fromPartial(e)) || [];
    return message;
  },
};

function createBaseQueryFixedDepositByAcctRequest(): QueryFixedDepositByAcctRequest {
  return { account: "", queryType: 0 };
}

export const QueryFixedDepositByAcctRequest = {
  encode(message: QueryFixedDepositByAcctRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.account !== "") {
      writer.uint32(18).string(message.account);
    }
    if (message.queryType !== 0) {
      writer.uint32(24).int32(message.queryType);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryFixedDepositByAcctRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryFixedDepositByAcctRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 2:
          message.account = reader.string();
          break;
        case 3:
          message.queryType = reader.int32() as any;
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryFixedDepositByAcctRequest {
    return {
      account: isSet(object.account) ? String(object.account) : "",
      queryType: isSet(object.queryType) ? fixedDepositStateFromJSON(object.queryType) : 0,
    };
  },

  toJSON(message: QueryFixedDepositByAcctRequest): unknown {
    const obj: any = {};
    message.account !== undefined && (obj.account = message.account);
    message.queryType !== undefined && (obj.queryType = fixedDepositStateToJSON(message.queryType));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryFixedDepositByAcctRequest>, I>>(
    object: I,
  ): QueryFixedDepositByAcctRequest {
    const message = createBaseQueryFixedDepositByAcctRequest();
    message.account = object.account ?? "";
    message.queryType = object.queryType ?? 0;
    return message;
  },
};

function createBaseQueryFixedDepositByAcctResponse(): QueryFixedDepositByAcctResponse {
  return { FixedDeposit: [] };
}

export const QueryFixedDepositByAcctResponse = {
  encode(message: QueryFixedDepositByAcctResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.FixedDeposit) {
      FixedDeposit.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryFixedDepositByAcctResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryFixedDepositByAcctResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.FixedDeposit.push(FixedDeposit.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryFixedDepositByAcctResponse {
    return {
      FixedDeposit: Array.isArray(object?.FixedDeposit)
        ? object.FixedDeposit.map((e: any) => FixedDeposit.fromJSON(e))
        : [],
    };
  },

  toJSON(message: QueryFixedDepositByAcctResponse): unknown {
    const obj: any = {};
    if (message.FixedDeposit) {
      obj.FixedDeposit = message.FixedDeposit.map((e) => e ? FixedDeposit.toJSON(e) : undefined);
    } else {
      obj.FixedDeposit = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryFixedDepositByAcctResponse>, I>>(
    object: I,
  ): QueryFixedDepositByAcctResponse {
    const message = createBaseQueryFixedDepositByAcctResponse();
    message.FixedDeposit = object.FixedDeposit?.map((e) => FixedDeposit.fromPartial(e)) || [];
    return message;
  },
};

function createBaseQueryFixedDepositByRegionRequest(): QueryFixedDepositByRegionRequest {
  return { regionId: "", pagination: undefined, queryType: 0 };
}

export const QueryFixedDepositByRegionRequest = {
  encode(message: QueryFixedDepositByRegionRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.regionId !== "") {
      writer.uint32(10).string(message.regionId);
    }
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }
    if (message.queryType !== 0) {
      writer.uint32(24).int32(message.queryType);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryFixedDepositByRegionRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryFixedDepositByRegionRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.regionId = reader.string();
          break;
        case 2:
          message.pagination = PageRequest.decode(reader, reader.uint32());
          break;
        case 3:
          message.queryType = reader.int32() as any;
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryFixedDepositByRegionRequest {
    return {
      regionId: isSet(object.regionId) ? String(object.regionId) : "",
      pagination: isSet(object.pagination) ? PageRequest.fromJSON(object.pagination) : undefined,
      queryType: isSet(object.queryType) ? fixedDepositStateFromJSON(object.queryType) : 0,
    };
  },

  toJSON(message: QueryFixedDepositByRegionRequest): unknown {
    const obj: any = {};
    message.regionId !== undefined && (obj.regionId = message.regionId);
    message.pagination !== undefined
      && (obj.pagination = message.pagination ? PageRequest.toJSON(message.pagination) : undefined);
    message.queryType !== undefined && (obj.queryType = fixedDepositStateToJSON(message.queryType));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryFixedDepositByRegionRequest>, I>>(
    object: I,
  ): QueryFixedDepositByRegionRequest {
    const message = createBaseQueryFixedDepositByRegionRequest();
    message.regionId = object.regionId ?? "";
    message.pagination = (object.pagination !== undefined && object.pagination !== null)
      ? PageRequest.fromPartial(object.pagination)
      : undefined;
    message.queryType = object.queryType ?? 0;
    return message;
  },
};

function createBaseQueryFixedDepositByRegionResponse(): QueryFixedDepositByRegionResponse {
  return { FixedDeposit: [], pagination: undefined };
}

export const QueryFixedDepositByRegionResponse = {
  encode(message: QueryFixedDepositByRegionResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.FixedDeposit) {
      FixedDeposit.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.pagination !== undefined) {
      PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryFixedDepositByRegionResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryFixedDepositByRegionResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.FixedDeposit.push(FixedDeposit.decode(reader, reader.uint32()));
          break;
        case 2:
          message.pagination = PageResponse.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryFixedDepositByRegionResponse {
    return {
      FixedDeposit: Array.isArray(object?.FixedDeposit)
        ? object.FixedDeposit.map((e: any) => FixedDeposit.fromJSON(e))
        : [],
      pagination: isSet(object.pagination) ? PageResponse.fromJSON(object.pagination) : undefined,
    };
  },

  toJSON(message: QueryFixedDepositByRegionResponse): unknown {
    const obj: any = {};
    if (message.FixedDeposit) {
      obj.FixedDeposit = message.FixedDeposit.map((e) => e ? FixedDeposit.toJSON(e) : undefined);
    } else {
      obj.FixedDeposit = [];
    }
    message.pagination !== undefined
      && (obj.pagination = message.pagination ? PageResponse.toJSON(message.pagination) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryFixedDepositByRegionResponse>, I>>(
    object: I,
  ): QueryFixedDepositByRegionResponse {
    const message = createBaseQueryFixedDepositByRegionResponse();
    message.FixedDeposit = object.FixedDeposit?.map((e) => FixedDeposit.fromPartial(e)) || [];
    message.pagination = (object.pagination !== undefined && object.pagination !== null)
      ? PageResponse.fromPartial(object.pagination)
      : undefined;
    return message;
  },
};

function createBaseQueryGetFixedDepositRequest(): QueryGetFixedDepositRequest {
  return { address: "", id: 0 };
}

export const QueryGetFixedDepositRequest = {
  encode(message: QueryGetFixedDepositRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.address !== "") {
      writer.uint32(10).string(message.address);
    }
    if (message.id !== 0) {
      writer.uint32(16).uint64(message.id);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryGetFixedDepositRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetFixedDepositRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.address = reader.string();
          break;
        case 2:
          message.id = longToNumber(reader.uint64() as Long);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryGetFixedDepositRequest {
    return {
      address: isSet(object.address) ? String(object.address) : "",
      id: isSet(object.id) ? Number(object.id) : 0,
    };
  },

  toJSON(message: QueryGetFixedDepositRequest): unknown {
    const obj: any = {};
    message.address !== undefined && (obj.address = message.address);
    message.id !== undefined && (obj.id = Math.round(message.id));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryGetFixedDepositRequest>, I>>(object: I): QueryGetFixedDepositRequest {
    const message = createBaseQueryGetFixedDepositRequest();
    message.address = object.address ?? "";
    message.id = object.id ?? 0;
    return message;
  },
};

function createBaseQueryGetFixedDepositResponse(): QueryGetFixedDepositResponse {
  return { FixedDeposit: undefined };
}

export const QueryGetFixedDepositResponse = {
  encode(message: QueryGetFixedDepositResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.FixedDeposit !== undefined) {
      FixedDeposit.encode(message.FixedDeposit, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryGetFixedDepositResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetFixedDepositResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.FixedDeposit = FixedDeposit.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryGetFixedDepositResponse {
    return { FixedDeposit: isSet(object.FixedDeposit) ? FixedDeposit.fromJSON(object.FixedDeposit) : undefined };
  },

  toJSON(message: QueryGetFixedDepositResponse): unknown {
    const obj: any = {};
    message.FixedDeposit !== undefined
      && (obj.FixedDeposit = message.FixedDeposit ? FixedDeposit.toJSON(message.FixedDeposit) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryGetFixedDepositResponse>, I>>(object: I): QueryGetFixedDepositResponse {
    const message = createBaseQueryGetFixedDepositResponse();
    message.FixedDeposit = (object.FixedDeposit !== undefined && object.FixedDeposit !== null)
      ? FixedDeposit.fromPartial(object.FixedDeposit)
      : undefined;
    return message;
  },
};

function createBaseQueryAllFixedDepositRequest(): QueryAllFixedDepositRequest {
  return { pagination: undefined };
}

export const QueryAllFixedDepositRequest = {
  encode(message: QueryAllFixedDepositRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryAllFixedDepositRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllFixedDepositRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.pagination = PageRequest.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryAllFixedDepositRequest {
    return { pagination: isSet(object.pagination) ? PageRequest.fromJSON(object.pagination) : undefined };
  },

  toJSON(message: QueryAllFixedDepositRequest): unknown {
    const obj: any = {};
    message.pagination !== undefined
      && (obj.pagination = message.pagination ? PageRequest.toJSON(message.pagination) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryAllFixedDepositRequest>, I>>(object: I): QueryAllFixedDepositRequest {
    const message = createBaseQueryAllFixedDepositRequest();
    message.pagination = (object.pagination !== undefined && object.pagination !== null)
      ? PageRequest.fromPartial(object.pagination)
      : undefined;
    return message;
  },
};

function createBaseQueryAllFixedDepositResponse(): QueryAllFixedDepositResponse {
  return { FixedDeposit: [], pagination: undefined };
}

export const QueryAllFixedDepositResponse = {
  encode(message: QueryAllFixedDepositResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.FixedDeposit) {
      FixedDeposit.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.pagination !== undefined) {
      PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryAllFixedDepositResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllFixedDepositResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.FixedDeposit.push(FixedDeposit.decode(reader, reader.uint32()));
          break;
        case 2:
          message.pagination = PageResponse.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryAllFixedDepositResponse {
    return {
      FixedDeposit: Array.isArray(object?.FixedDeposit)
        ? object.FixedDeposit.map((e: any) => FixedDeposit.fromJSON(e))
        : [],
      pagination: isSet(object.pagination) ? PageResponse.fromJSON(object.pagination) : undefined,
    };
  },

  toJSON(message: QueryAllFixedDepositResponse): unknown {
    const obj: any = {};
    if (message.FixedDeposit) {
      obj.FixedDeposit = message.FixedDeposit.map((e) => e ? FixedDeposit.toJSON(e) : undefined);
    } else {
      obj.FixedDeposit = [];
    }
    message.pagination !== undefined
      && (obj.pagination = message.pagination ? PageResponse.toJSON(message.pagination) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryAllFixedDepositResponse>, I>>(object: I): QueryAllFixedDepositResponse {
    const message = createBaseQueryAllFixedDepositResponse();
    message.FixedDeposit = object.FixedDeposit?.map((e) => FixedDeposit.fromPartial(e)) || [];
    message.pagination = (object.pagination !== undefined && object.pagination !== null)
      ? PageResponse.fromPartial(object.pagination)
      : undefined;
    return message;
  },
};

function createBaseQueryFixedDepositCfgRequest(): QueryFixedDepositCfgRequest {
  return { regionIds: [] };
}

export const QueryFixedDepositCfgRequest = {
  encode(message: QueryFixedDepositCfgRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.regionIds) {
      writer.uint32(10).string(v!);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryFixedDepositCfgRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryFixedDepositCfgRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.regionIds.push(reader.string());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryFixedDepositCfgRequest {
    return { regionIds: Array.isArray(object?.regionIds) ? object.regionIds.map((e: any) => String(e)) : [] };
  },

  toJSON(message: QueryFixedDepositCfgRequest): unknown {
    const obj: any = {};
    if (message.regionIds) {
      obj.regionIds = message.regionIds.map((e) => e);
    } else {
      obj.regionIds = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryFixedDepositCfgRequest>, I>>(object: I): QueryFixedDepositCfgRequest {
    const message = createBaseQueryFixedDepositCfgRequest();
    message.regionIds = object.regionIds?.map((e) => e) || [];
    return message;
  },
};

function createBaseQueryFixedDepositCfgResponse(): QueryFixedDepositCfgResponse {
  return { RegionFixedDepositCfgs: [] };
}

export const QueryFixedDepositCfgResponse = {
  encode(message: QueryFixedDepositCfgResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.RegionFixedDepositCfgs) {
      RegionAllFixedDepositCfg.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryFixedDepositCfgResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryFixedDepositCfgResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.RegionFixedDepositCfgs.push(RegionAllFixedDepositCfg.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryFixedDepositCfgResponse {
    return {
      RegionFixedDepositCfgs: Array.isArray(object?.RegionFixedDepositCfgs)
        ? object.RegionFixedDepositCfgs.map((e: any) => RegionAllFixedDepositCfg.fromJSON(e))
        : [],
    };
  },

  toJSON(message: QueryFixedDepositCfgResponse): unknown {
    const obj: any = {};
    if (message.RegionFixedDepositCfgs) {
      obj.RegionFixedDepositCfgs = message.RegionFixedDepositCfgs.map((e) =>
        e ? RegionAllFixedDepositCfg.toJSON(e) : undefined
      );
    } else {
      obj.RegionFixedDepositCfgs = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryFixedDepositCfgResponse>, I>>(object: I): QueryFixedDepositCfgResponse {
    const message = createBaseQueryFixedDepositCfgResponse();
    message.RegionFixedDepositCfgs = object.RegionFixedDepositCfgs?.map((e) => RegionAllFixedDepositCfg.fromPartial(e))
      || [];
    return message;
  },
};

function createBaseQueryFixedDepositCfgByTermRequest(): QueryFixedDepositCfgByTermRequest {
  return { regionId: "", term: 0 };
}

export const QueryFixedDepositCfgByTermRequest = {
  encode(message: QueryFixedDepositCfgByTermRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.regionId !== "") {
      writer.uint32(10).string(message.regionId);
    }
    if (message.term !== 0) {
      writer.uint32(16).int64(message.term);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryFixedDepositCfgByTermRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryFixedDepositCfgByTermRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.regionId = reader.string();
          break;
        case 2:
          message.term = longToNumber(reader.int64() as Long);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryFixedDepositCfgByTermRequest {
    return {
      regionId: isSet(object.regionId) ? String(object.regionId) : "",
      term: isSet(object.term) ? Number(object.term) : 0,
    };
  },

  toJSON(message: QueryFixedDepositCfgByTermRequest): unknown {
    const obj: any = {};
    message.regionId !== undefined && (obj.regionId = message.regionId);
    message.term !== undefined && (obj.term = Math.round(message.term));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryFixedDepositCfgByTermRequest>, I>>(
    object: I,
  ): QueryFixedDepositCfgByTermRequest {
    const message = createBaseQueryFixedDepositCfgByTermRequest();
    message.regionId = object.regionId ?? "";
    message.term = object.term ?? 0;
    return message;
  },
};

function createBaseQueryFixedDepositCfgByTermResponse(): QueryFixedDepositCfgByTermResponse {
  return { FixedDepositCfg: undefined };
}

export const QueryFixedDepositCfgByTermResponse = {
  encode(message: QueryFixedDepositCfgByTermResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.FixedDepositCfg !== undefined) {
      FixedDepositCfg.encode(message.FixedDepositCfg, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryFixedDepositCfgByTermResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryFixedDepositCfgByTermResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.FixedDepositCfg = FixedDepositCfg.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryFixedDepositCfgByTermResponse {
    return {
      FixedDepositCfg: isSet(object.FixedDepositCfg) ? FixedDepositCfg.fromJSON(object.FixedDepositCfg) : undefined,
    };
  },

  toJSON(message: QueryFixedDepositCfgByTermResponse): unknown {
    const obj: any = {};
    message.FixedDepositCfg !== undefined
      && (obj.FixedDepositCfg = message.FixedDepositCfg ? FixedDepositCfg.toJSON(message.FixedDepositCfg) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryFixedDepositCfgByTermResponse>, I>>(
    object: I,
  ): QueryFixedDepositCfgByTermResponse {
    const message = createBaseQueryFixedDepositCfgByTermResponse();
    message.FixedDepositCfg = (object.FixedDepositCfg !== undefined && object.FixedDepositCfg !== null)
      ? FixedDepositCfg.fromPartial(object.FixedDepositCfg)
      : undefined;
    return message;
  },
};

function createBaseQueryFixedDepositTotalAmountRequest(): QueryFixedDepositTotalAmountRequest {
  return {};
}

export const QueryFixedDepositTotalAmountRequest = {
  encode(_: QueryFixedDepositTotalAmountRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryFixedDepositTotalAmountRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryFixedDepositTotalAmountRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): QueryFixedDepositTotalAmountRequest {
    return {};
  },

  toJSON(_: QueryFixedDepositTotalAmountRequest): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryFixedDepositTotalAmountRequest>, I>>(
    _: I,
  ): QueryFixedDepositTotalAmountRequest {
    const message = createBaseQueryFixedDepositTotalAmountRequest();
    return message;
  },
};

function createBaseQueryFixedDepositTotalAmountResponse(): QueryFixedDepositTotalAmountResponse {
  return { amount: undefined };
}

export const QueryFixedDepositTotalAmountResponse = {
  encode(message: QueryFixedDepositTotalAmountResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.amount !== undefined) {
      Coin.encode(message.amount, writer.uint32(58).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryFixedDepositTotalAmountResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryFixedDepositTotalAmountResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 7:
          message.amount = Coin.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryFixedDepositTotalAmountResponse {
    return { amount: isSet(object.amount) ? Coin.fromJSON(object.amount) : undefined };
  },

  toJSON(message: QueryFixedDepositTotalAmountResponse): unknown {
    const obj: any = {};
    message.amount !== undefined && (obj.amount = message.amount ? Coin.toJSON(message.amount) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryFixedDepositTotalAmountResponse>, I>>(
    object: I,
  ): QueryFixedDepositTotalAmountResponse {
    const message = createBaseQueryFixedDepositTotalAmountResponse();
    message.amount = (object.amount !== undefined && object.amount !== null)
      ? Coin.fromPartial(object.amount)
      : undefined;
    return message;
  },
};

function createBaseQueryFixedDepositAmountByMeidRequest(): QueryFixedDepositAmountByMeidRequest {
  return { account: "" };
}

export const QueryFixedDepositAmountByMeidRequest = {
  encode(message: QueryFixedDepositAmountByMeidRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.account !== "") {
      writer.uint32(10).string(message.account);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryFixedDepositAmountByMeidRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryFixedDepositAmountByMeidRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.account = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryFixedDepositAmountByMeidRequest {
    return { account: isSet(object.account) ? String(object.account) : "" };
  },

  toJSON(message: QueryFixedDepositAmountByMeidRequest): unknown {
    const obj: any = {};
    message.account !== undefined && (obj.account = message.account);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryFixedDepositAmountByMeidRequest>, I>>(
    object: I,
  ): QueryFixedDepositAmountByMeidRequest {
    const message = createBaseQueryFixedDepositAmountByMeidRequest();
    message.account = object.account ?? "";
    return message;
  },
};

function createBaseQueryFixedDepositAmountByMeidResponse(): QueryFixedDepositAmountByMeidResponse {
  return { amount: undefined };
}

export const QueryFixedDepositAmountByMeidResponse = {
  encode(message: QueryFixedDepositAmountByMeidResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.amount !== undefined) {
      Coin.encode(message.amount, writer.uint32(58).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryFixedDepositAmountByMeidResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryFixedDepositAmountByMeidResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 7:
          message.amount = Coin.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryFixedDepositAmountByMeidResponse {
    return { amount: isSet(object.amount) ? Coin.fromJSON(object.amount) : undefined };
  },

  toJSON(message: QueryFixedDepositAmountByMeidResponse): unknown {
    const obj: any = {};
    message.amount !== undefined && (obj.amount = message.amount ? Coin.toJSON(message.amount) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryFixedDepositAmountByMeidResponse>, I>>(
    object: I,
  ): QueryFixedDepositAmountByMeidResponse {
    const message = createBaseQueryFixedDepositAmountByMeidResponse();
    message.amount = (object.amount !== undefined && object.amount !== null)
      ? Coin.fromPartial(object.amount)
      : undefined;
    return message;
  },
};

function createBaseQueryAllDelegationsRequest(): QueryAllDelegationsRequest {
  return { pagination: undefined };
}

export const QueryAllDelegationsRequest = {
  encode(message: QueryAllDelegationsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryAllDelegationsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllDelegationsRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.pagination = PageRequest.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryAllDelegationsRequest {
    return { pagination: isSet(object.pagination) ? PageRequest.fromJSON(object.pagination) : undefined };
  },

  toJSON(message: QueryAllDelegationsRequest): unknown {
    const obj: any = {};
    message.pagination !== undefined
      && (obj.pagination = message.pagination ? PageRequest.toJSON(message.pagination) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryAllDelegationsRequest>, I>>(object: I): QueryAllDelegationsRequest {
    const message = createBaseQueryAllDelegationsRequest();
    message.pagination = (object.pagination !== undefined && object.pagination !== null)
      ? PageRequest.fromPartial(object.pagination)
      : undefined;
    return message;
  },
};

function createBaseQueryAllDelegationsResponse(): QueryAllDelegationsResponse {
  return { delegations: [], pagination: undefined };
}

export const QueryAllDelegationsResponse = {
  encode(message: QueryAllDelegationsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.delegations) {
      Delegation.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.pagination !== undefined) {
      PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryAllDelegationsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllDelegationsResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.delegations.push(Delegation.decode(reader, reader.uint32()));
          break;
        case 2:
          message.pagination = PageResponse.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryAllDelegationsResponse {
    return {
      delegations: Array.isArray(object?.delegations) ? object.delegations.map((e: any) => Delegation.fromJSON(e)) : [],
      pagination: isSet(object.pagination) ? PageResponse.fromJSON(object.pagination) : undefined,
    };
  },

  toJSON(message: QueryAllDelegationsResponse): unknown {
    const obj: any = {};
    if (message.delegations) {
      obj.delegations = message.delegations.map((e) => e ? Delegation.toJSON(e) : undefined);
    } else {
      obj.delegations = [];
    }
    message.pagination !== undefined
      && (obj.pagination = message.pagination ? PageResponse.toJSON(message.pagination) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryAllDelegationsResponse>, I>>(object: I): QueryAllDelegationsResponse {
    const message = createBaseQueryAllDelegationsResponse();
    message.delegations = object.delegations?.map((e) => Delegation.fromPartial(e)) || [];
    message.pagination = (object.pagination !== undefined && object.pagination !== null)
      ? PageResponse.fromPartial(object.pagination)
      : undefined;
    return message;
  },
};

function createBaseQueryStakesRequest(): QueryStakesRequest {
  return { pagination: undefined };
}

export const QueryStakesRequest = {
  encode(message: QueryStakesRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryStakesRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryStakesRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.pagination = PageRequest.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryStakesRequest {
    return { pagination: isSet(object.pagination) ? PageRequest.fromJSON(object.pagination) : undefined };
  },

  toJSON(message: QueryStakesRequest): unknown {
    const obj: any = {};
    message.pagination !== undefined
      && (obj.pagination = message.pagination ? PageRequest.toJSON(message.pagination) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryStakesRequest>, I>>(object: I): QueryStakesRequest {
    const message = createBaseQueryStakesRequest();
    message.pagination = (object.pagination !== undefined && object.pagination !== null)
      ? PageRequest.fromPartial(object.pagination)
      : undefined;
    return message;
  },
};

function createBaseQueryStakesResponse(): QueryStakesResponse {
  return { stakes: [], pagination: undefined };
}

export const QueryStakesResponse = {
  encode(message: QueryStakesResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.stakes) {
      Stake.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.pagination !== undefined) {
      PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryStakesResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryStakesResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.stakes.push(Stake.decode(reader, reader.uint32()));
          break;
        case 2:
          message.pagination = PageResponse.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryStakesResponse {
    return {
      stakes: Array.isArray(object?.stakes) ? object.stakes.map((e: any) => Stake.fromJSON(e)) : [],
      pagination: isSet(object.pagination) ? PageResponse.fromJSON(object.pagination) : undefined,
    };
  },

  toJSON(message: QueryStakesResponse): unknown {
    const obj: any = {};
    if (message.stakes) {
      obj.stakes = message.stakes.map((e) => e ? Stake.toJSON(e) : undefined);
    } else {
      obj.stakes = [];
    }
    message.pagination !== undefined
      && (obj.pagination = message.pagination ? PageResponse.toJSON(message.pagination) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryStakesResponse>, I>>(object: I): QueryStakesResponse {
    const message = createBaseQueryStakesResponse();
    message.stakes = object.stakes?.map((e) => Stake.fromPartial(e)) || [];
    message.pagination = (object.pagination !== undefined && object.pagination !== null)
      ? PageResponse.fromPartial(object.pagination)
      : undefined;
    return message;
  },
};

function createBaseQueryAllRecords(): QueryAllRecords {
  return { pagination: undefined };
}

export const QueryAllRecords = {
  encode(message: QueryAllRecords, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryAllRecords {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllRecords();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.pagination = PageRequest.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryAllRecords {
    return { pagination: isSet(object.pagination) ? PageRequest.fromJSON(object.pagination) : undefined };
  },

  toJSON(message: QueryAllRecords): unknown {
    const obj: any = {};
    message.pagination !== undefined
      && (obj.pagination = message.pagination ? PageRequest.toJSON(message.pagination) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryAllRecords>, I>>(object: I): QueryAllRecords {
    const message = createBaseQueryAllRecords();
    message.pagination = (object.pagination !== undefined && object.pagination !== null)
      ? PageRequest.fromPartial(object.pagination)
      : undefined;
    return message;
  },
};

function createBaseQueryAllRecordsResponse(): QueryAllRecordsResponse {
  return { records: [], pagination: undefined };
}

export const QueryAllRecordsResponse = {
  encode(message: QueryAllRecordsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.records) {
      Record.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.pagination !== undefined) {
      PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryAllRecordsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllRecordsResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.records.push(Record.decode(reader, reader.uint32()));
          break;
        case 2:
          message.pagination = PageResponse.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryAllRecordsResponse {
    return {
      records: Array.isArray(object?.records) ? object.records.map((e: any) => Record.fromJSON(e)) : [],
      pagination: isSet(object.pagination) ? PageResponse.fromJSON(object.pagination) : undefined,
    };
  },

  toJSON(message: QueryAllRecordsResponse): unknown {
    const obj: any = {};
    if (message.records) {
      obj.records = message.records.map((e) => e ? Record.toJSON(e) : undefined);
    } else {
      obj.records = [];
    }
    message.pagination !== undefined
      && (obj.pagination = message.pagination ? PageResponse.toJSON(message.pagination) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryAllRecordsResponse>, I>>(object: I): QueryAllRecordsResponse {
    const message = createBaseQueryAllRecordsResponse();
    message.records = object.records?.map((e) => Record.fromPartial(e)) || [];
    message.pagination = (object.pagination !== undefined && object.pagination !== null)
      ? PageResponse.fromPartial(object.pagination)
      : undefined;
    return message;
  },
};

function createBaseQueryRecordsByAddress(): QueryRecordsByAddress {
  return { account: "" };
}

export const QueryRecordsByAddress = {
  encode(message: QueryRecordsByAddress, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.account !== "") {
      writer.uint32(10).string(message.account);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryRecordsByAddress {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryRecordsByAddress();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.account = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryRecordsByAddress {
    return { account: isSet(object.account) ? String(object.account) : "" };
  },

  toJSON(message: QueryRecordsByAddress): unknown {
    const obj: any = {};
    message.account !== undefined && (obj.account = message.account);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryRecordsByAddress>, I>>(object: I): QueryRecordsByAddress {
    const message = createBaseQueryRecordsByAddress();
    message.account = object.account ?? "";
    return message;
  },
};

function createBaseQueryRecordsByAddressResponse(): QueryRecordsByAddressResponse {
  return { records: [] };
}

export const QueryRecordsByAddressResponse = {
  encode(message: QueryRecordsByAddressResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.records) {
      Record.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryRecordsByAddressResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryRecordsByAddressResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.records.push(Record.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryRecordsByAddressResponse {
    return { records: Array.isArray(object?.records) ? object.records.map((e: any) => Record.fromJSON(e)) : [] };
  },

  toJSON(message: QueryRecordsByAddressResponse): unknown {
    const obj: any = {};
    if (message.records) {
      obj.records = message.records.map((e) => e ? Record.toJSON(e) : undefined);
    } else {
      obj.records = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryRecordsByAddressResponse>, I>>(
    object: I,
  ): QueryRecordsByAddressResponse {
    const message = createBaseQueryRecordsByAddressResponse();
    message.records = object.records?.map((e) => Record.fromPartial(e)) || [];
    return message;
  },
};

function createBaseQueryReviewRecordByNumber(): QueryReviewRecordByNumber {
  return { actionNumber: "" };
}

export const QueryReviewRecordByNumber = {
  encode(message: QueryReviewRecordByNumber, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.actionNumber !== "") {
      writer.uint32(10).string(message.actionNumber);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryReviewRecordByNumber {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryReviewRecordByNumber();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.actionNumber = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryReviewRecordByNumber {
    return { actionNumber: isSet(object.actionNumber) ? String(object.actionNumber) : "" };
  },

  toJSON(message: QueryReviewRecordByNumber): unknown {
    const obj: any = {};
    message.actionNumber !== undefined && (obj.actionNumber = message.actionNumber);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryReviewRecordByNumber>, I>>(object: I): QueryReviewRecordByNumber {
    const message = createBaseQueryReviewRecordByNumber();
    message.actionNumber = object.actionNumber ?? "";
    return message;
  },
};

function createBaseQueryReviewRecordByNumberResponse(): QueryReviewRecordByNumberResponse {
  return { reviewRecord: undefined };
}

export const QueryReviewRecordByNumberResponse = {
  encode(message: QueryReviewRecordByNumberResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.reviewRecord !== undefined) {
      ReviewRecord.encode(message.reviewRecord, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryReviewRecordByNumberResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryReviewRecordByNumberResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.reviewRecord = ReviewRecord.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryReviewRecordByNumberResponse {
    return { reviewRecord: isSet(object.reviewRecord) ? ReviewRecord.fromJSON(object.reviewRecord) : undefined };
  },

  toJSON(message: QueryReviewRecordByNumberResponse): unknown {
    const obj: any = {};
    message.reviewRecord !== undefined
      && (obj.reviewRecord = message.reviewRecord ? ReviewRecord.toJSON(message.reviewRecord) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryReviewRecordByNumberResponse>, I>>(
    object: I,
  ): QueryReviewRecordByNumberResponse {
    const message = createBaseQueryReviewRecordByNumberResponse();
    message.reviewRecord = (object.reviewRecord !== undefined && object.reviewRecord !== null)
      ? ReviewRecord.fromPartial(object.reviewRecord)
      : undefined;
    return message;
  },
};

export interface Query {
  /** Queries a list of Region items. */
  Region(request: QueryRegionRequest): Promise<QueryRegionResponse>;
  AllRegion(request: QueryAllRegionRequest): Promise<QueryAllRegionResponse>;
  DelegationRewards(request: QueryDelegationRewardsRequest): Promise<QueryDelegationRewardsResponse>;
  FixedDepositTotalAmount(request: QueryFixedDepositTotalAmountRequest): Promise<QueryFixedDepositTotalAmountResponse>;
  FixedDepositAmountByMeid(
    request: QueryFixedDepositAmountByMeidRequest,
  ): Promise<QueryFixedDepositAmountByMeidResponse>;
  /** Queries a list of FixedDepositByAcct items. */
  FixedDepositByAcct(request: QueryFixedDepositByAcctRequest): Promise<QueryFixedDepositByAcctResponse>;
  /** Queries a list of FixedDepositByRegion items. */
  FixedDepositByRegion(request: QueryFixedDepositByRegionRequest): Promise<QueryFixedDepositByRegionResponse>;
  /** Queries a list of FixedDeposit items. */
  FixedDeposit(request: QueryGetFixedDepositRequest): Promise<QueryGetFixedDepositResponse>;
  FixedDepositAll(request: QueryAllFixedDepositRequest): Promise<QueryAllFixedDepositResponse>;
  FixedDepositCfg(request: QueryFixedDepositCfgRequest): Promise<QueryFixedDepositCfgResponse>;
  FixedDepositCfgByTerm(request: QueryFixedDepositCfgByTermRequest): Promise<QueryFixedDepositCfgByTermResponse>;
  /** Delegation queries delegate info for given validator delegator pair. */
  Delegation(request: QueryDelegationRequest): Promise<QueryDelegationResponse>;
  AllDelegations(request: QueryAllDelegationsRequest): Promise<QueryAllDelegationsResponse>;
  Stakes(request: QueryStakesRequest): Promise<QueryStakesResponse>;
  /** Queries a list of Record */
  QueryAllRecord(request: QueryAllRecords): Promise<QueryAllRecordsResponse>;
  QueryRecordByAddress(request: QueryRecordsByAddress): Promise<QueryRecordsByAddressResponse>;
  QueryReviewRecordByID(request: QueryReviewRecordByNumber): Promise<QueryReviewRecordByNumberResponse>;
}

export class QueryClientImpl implements Query {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.Region = this.Region.bind(this);
    this.AllRegion = this.AllRegion.bind(this);
    this.DelegationRewards = this.DelegationRewards.bind(this);
    this.FixedDepositTotalAmount = this.FixedDepositTotalAmount.bind(this);
    this.FixedDepositAmountByMeid = this.FixedDepositAmountByMeid.bind(this);
    this.FixedDepositByAcct = this.FixedDepositByAcct.bind(this);
    this.FixedDepositByRegion = this.FixedDepositByRegion.bind(this);
    this.FixedDeposit = this.FixedDeposit.bind(this);
    this.FixedDepositAll = this.FixedDepositAll.bind(this);
    this.FixedDepositCfg = this.FixedDepositCfg.bind(this);
    this.FixedDepositCfgByTerm = this.FixedDepositCfgByTerm.bind(this);
    this.Delegation = this.Delegation.bind(this);
    this.AllDelegations = this.AllDelegations.bind(this);
    this.Stakes = this.Stakes.bind(this);
    this.QueryAllRecord = this.QueryAllRecord.bind(this);
    this.QueryRecordByAddress = this.QueryRecordByAddress.bind(this);
    this.QueryReviewRecordByID = this.QueryReviewRecordByID.bind(this);
  }
  Region(request: QueryRegionRequest): Promise<QueryRegionResponse> {
    const data = QueryRegionRequest.encode(request).finish();
    const promise = this.rpc.request("metaearth.wstaking.Query", "Region", data);
    return promise.then((data) => QueryRegionResponse.decode(new _m0.Reader(data)));
  }

  AllRegion(request: QueryAllRegionRequest): Promise<QueryAllRegionResponse> {
    const data = QueryAllRegionRequest.encode(request).finish();
    const promise = this.rpc.request("metaearth.wstaking.Query", "AllRegion", data);
    return promise.then((data) => QueryAllRegionResponse.decode(new _m0.Reader(data)));
  }

  DelegationRewards(request: QueryDelegationRewardsRequest): Promise<QueryDelegationRewardsResponse> {
    const data = QueryDelegationRewardsRequest.encode(request).finish();
    const promise = this.rpc.request("metaearth.wstaking.Query", "DelegationRewards", data);
    return promise.then((data) => QueryDelegationRewardsResponse.decode(new _m0.Reader(data)));
  }

  FixedDepositTotalAmount(request: QueryFixedDepositTotalAmountRequest): Promise<QueryFixedDepositTotalAmountResponse> {
    const data = QueryFixedDepositTotalAmountRequest.encode(request).finish();
    const promise = this.rpc.request("metaearth.wstaking.Query", "FixedDepositTotalAmount", data);
    return promise.then((data) => QueryFixedDepositTotalAmountResponse.decode(new _m0.Reader(data)));
  }

  FixedDepositAmountByMeid(
    request: QueryFixedDepositAmountByMeidRequest,
  ): Promise<QueryFixedDepositAmountByMeidResponse> {
    const data = QueryFixedDepositAmountByMeidRequest.encode(request).finish();
    const promise = this.rpc.request("metaearth.wstaking.Query", "FixedDepositAmountByMeid", data);
    return promise.then((data) => QueryFixedDepositAmountByMeidResponse.decode(new _m0.Reader(data)));
  }

  FixedDepositByAcct(request: QueryFixedDepositByAcctRequest): Promise<QueryFixedDepositByAcctResponse> {
    const data = QueryFixedDepositByAcctRequest.encode(request).finish();
    const promise = this.rpc.request("metaearth.wstaking.Query", "FixedDepositByAcct", data);
    return promise.then((data) => QueryFixedDepositByAcctResponse.decode(new _m0.Reader(data)));
  }

  FixedDepositByRegion(request: QueryFixedDepositByRegionRequest): Promise<QueryFixedDepositByRegionResponse> {
    const data = QueryFixedDepositByRegionRequest.encode(request).finish();
    const promise = this.rpc.request("metaearth.wstaking.Query", "FixedDepositByRegion", data);
    return promise.then((data) => QueryFixedDepositByRegionResponse.decode(new _m0.Reader(data)));
  }

  FixedDeposit(request: QueryGetFixedDepositRequest): Promise<QueryGetFixedDepositResponse> {
    const data = QueryGetFixedDepositRequest.encode(request).finish();
    const promise = this.rpc.request("metaearth.wstaking.Query", "FixedDeposit", data);
    return promise.then((data) => QueryGetFixedDepositResponse.decode(new _m0.Reader(data)));
  }

  FixedDepositAll(request: QueryAllFixedDepositRequest): Promise<QueryAllFixedDepositResponse> {
    const data = QueryAllFixedDepositRequest.encode(request).finish();
    const promise = this.rpc.request("metaearth.wstaking.Query", "FixedDepositAll", data);
    return promise.then((data) => QueryAllFixedDepositResponse.decode(new _m0.Reader(data)));
  }

  FixedDepositCfg(request: QueryFixedDepositCfgRequest): Promise<QueryFixedDepositCfgResponse> {
    const data = QueryFixedDepositCfgRequest.encode(request).finish();
    const promise = this.rpc.request("metaearth.wstaking.Query", "FixedDepositCfg", data);
    return promise.then((data) => QueryFixedDepositCfgResponse.decode(new _m0.Reader(data)));
  }

  FixedDepositCfgByTerm(request: QueryFixedDepositCfgByTermRequest): Promise<QueryFixedDepositCfgByTermResponse> {
    const data = QueryFixedDepositCfgByTermRequest.encode(request).finish();
    const promise = this.rpc.request("metaearth.wstaking.Query", "FixedDepositCfgByTerm", data);
    return promise.then((data) => QueryFixedDepositCfgByTermResponse.decode(new _m0.Reader(data)));
  }

  Delegation(request: QueryDelegationRequest): Promise<QueryDelegationResponse> {
    const data = QueryDelegationRequest.encode(request).finish();
    const promise = this.rpc.request("metaearth.wstaking.Query", "Delegation", data);
    return promise.then((data) => QueryDelegationResponse.decode(new _m0.Reader(data)));
  }

  AllDelegations(request: QueryAllDelegationsRequest): Promise<QueryAllDelegationsResponse> {
    const data = QueryAllDelegationsRequest.encode(request).finish();
    const promise = this.rpc.request("metaearth.wstaking.Query", "AllDelegations", data);
    return promise.then((data) => QueryAllDelegationsResponse.decode(new _m0.Reader(data)));
  }

  Stakes(request: QueryStakesRequest): Promise<QueryStakesResponse> {
    const data = QueryStakesRequest.encode(request).finish();
    const promise = this.rpc.request("metaearth.wstaking.Query", "Stakes", data);
    return promise.then((data) => QueryStakesResponse.decode(new _m0.Reader(data)));
  }

  QueryAllRecord(request: QueryAllRecords): Promise<QueryAllRecordsResponse> {
    const data = QueryAllRecords.encode(request).finish();
    const promise = this.rpc.request("metaearth.wstaking.Query", "QueryAllRecord", data);
    return promise.then((data) => QueryAllRecordsResponse.decode(new _m0.Reader(data)));
  }

  QueryRecordByAddress(request: QueryRecordsByAddress): Promise<QueryRecordsByAddressResponse> {
    const data = QueryRecordsByAddress.encode(request).finish();
    const promise = this.rpc.request("metaearth.wstaking.Query", "QueryRecordByAddress", data);
    return promise.then((data) => QueryRecordsByAddressResponse.decode(new _m0.Reader(data)));
  }

  QueryReviewRecordByID(request: QueryReviewRecordByNumber): Promise<QueryReviewRecordByNumberResponse> {
    const data = QueryReviewRecordByNumber.encode(request).finish();
    const promise = this.rpc.request("metaearth.wstaking.Query", "QueryReviewRecordByID", data);
    return promise.then((data) => QueryReviewRecordByNumberResponse.decode(new _m0.Reader(data)));
  }
}

interface Rpc {
  request(service: string, method: string, data: Uint8Array): Promise<Uint8Array>;
}

declare var self: any | undefined;
declare var window: any | undefined;
declare var global: any | undefined;
var globalThis: any = (() => {
  if (typeof globalThis !== "undefined") {
    return globalThis;
  }
  if (typeof self !== "undefined") {
    return self;
  }
  if (typeof window !== "undefined") {
    return window;
  }
  if (typeof global !== "undefined") {
    return global;
  }
  throw "Unable to locate global object";
})();

type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;

export type DeepPartial<T> = T extends Builtin ? T
  : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
  : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

type KeysOfUnion<T> = T extends T ? keyof T : never;
export type Exact<P, I extends P> = P extends Builtin ? P
  : P & { [K in keyof P]: Exact<P[K], I[K]> } & { [K in Exclude<keyof I, KeysOfUnion<P>>]: never };

function longToNumber(long: Long): number {
  if (long.gt(Number.MAX_SAFE_INTEGER)) {
    throw new globalThis.Error("Value is larger than Number.MAX_SAFE_INTEGER");
  }
  return long.toNumber();
}

if (_m0.util.Long !== Long) {
  _m0.util.Long = Long as any;
  _m0.configure();
}

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}
