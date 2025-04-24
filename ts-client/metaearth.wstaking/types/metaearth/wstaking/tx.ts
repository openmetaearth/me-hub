/* eslint-disable */
import Long from "long";
import _m0 from "protobufjs/minimal";
import { Coin } from "../../cosmos/base/v1beta1/coin";
import { Timestamp } from "../../google/protobuf/timestamp";
import { fixedDepositCfgStatus, fixedDepositCfgStatusFromJSON, fixedDepositCfgStatusToJSON } from "./fixed_deposit";

export const protobufPackage = "metaearth.wstaking";

/**
 * MsgStake defines a SDK message for performing a stake of coins
 * from a staker to a validator.
 */
export interface MsgStake {
  stakerAddress: string;
  validatorAddress: string;
  amount: Coin | undefined;
}

/** MsgStakeResponse defines the Msg/Stake response type. */
export interface MsgStakeResponse {
}

/**
 * MsgUnstake defines a SDK message for performing an unstake from a
 * stake and a validator.
 */
export interface MsgUnstake {
  stakerAddress: string;
  validatorAddress: string;
  amount: Coin | undefined;
}

/** MsgUnstakeResponse defines the Msg/MsgUnstake response type. */
export interface MsgUnstakeResponse {
  completionTime: Date | undefined;
}

export interface MsgFixedDepositCfg {
  admin: string;
  regionId: string;
  term: number;
  rate: string;
}

export interface MsgFixedDepositCfgResp {
  retcode: string;
}

export interface MsgRemoveFixedDepositCfg {
  admin: string;
  regionId: string;
  term: number;
}

export interface MsgRemoveFixedDepositCfgResp {
  retcode: string;
}

export interface MsgSetFixedDepositCfgStatus {
  admin: string;
  regionId: string;
  term: number;
  status: fixedDepositCfgStatus;
}

export interface MsgSetFixedDepositCfgStatusResp {
  retcode: string;
}

export interface MsgSetFixedDepositCfgRate {
  admin: string;
  regionId: string;
  term: number;
  rate: string;
}

export interface MsgSetFixedDepositCfgRateResp {
  retcode: string;
}

export interface MsgNewRegion {
  creator: string;
  name: string;
  operatorAddress: string;
}

export interface MsgNewRegionResponse {
  regionId: string;
}

export interface MsgRemoveRegion {
  creator: string;
  regionId: string;
}

export interface MsgRemoveRegionResponse {
  retcode: string;
}

export interface MsgWithdrawFromRegion {
  withdrawer: string;
  regionId: string;
  receiver: string;
  amount: Coin[];
}

export interface MsgWithdrawFromRegionResp {
}

export interface MsgWithdrawFromGlobalDaoFeePool {
  withdrawer: string;
  amount: Coin[];
}

export interface MsgWithdrawFromGlobalDaoFeePoolResp {
}

/**
 * MsgWithdrawDelegatorReward represents delegation withdrawal to a delegator
 * from a single validator.
 */
export interface MsgWithdrawDelegatorReward {
  delegatorAddress: string;
  validatorAddress: string;
}

/**
 * MsgWithdrawDelegatorRewardResponse defines the Msg/WithdrawDelegatorReward
 * response type.
 */
export interface MsgWithdrawDelegatorRewardResponse {
  /** Since: cosmos-sdk 0.46 */
  amount: Coin[];
}

export interface MsgNewFixedDepositCfg {
  dao: string;
  regionId: string;
  term: number;
  rate: string;
}

export interface MsgNewFixedDepositCfgResp {
  retcode: string;
}

export interface MsgDoFixedDeposit {
  account: string;
  principal: Coin | undefined;
  term: number;
}

export interface MsgDoFixedDepositResponse {
  id: number;
}

export interface MsgWithdrawFixedDeposit {
  account: string;
  id: number;
}

export interface MsgWithdrawFixedDepositResponse {
  principal: Coin | undefined;
  interest: Coin | undefined;
  term: number;
  rate: string;
}

export interface MsgNewRecord {
  actionNumber: string;
  actionUrl: string;
  from: string;
}

export interface MsgNewRecordResponse {
}

export interface MsgReviewRecord {
  recordHash: string;
  reviewResult: string;
  from: string;
  ActionNumber: string;
  reviewedAddress: string;
}

export interface MsgReviewRecordResponse {
}

export interface MsgTransferRegion {
  fromRegion: string;
  toRegion: string;
  address: string[];
  creator: string;
}

export interface MsgTransferRegionResponse {
}

export interface MsgIbcTransferFromRegionTreasure {
  /** the port on which the packet will be sent */
  sourcePort: string;
  /** the channel by which the packet will be sent */
  sourceChannel: string;
  regionId: string;
  /** the tokens to be transferred */
  token:
    | Coin
    | undefined;
  /**
   * the recipient address on the destination chain
   * string receiver = 5;
   * Timeout height relative to the current block height.
   * The timeout is disabled when set to 0.
   */
  timeoutHeight:
    | Height
    | undefined;
  /**
   * Timeout timestamp in absolute nanoseconds since unix epoch.
   * The timeout is disabled when set to 0.
   */
  timeoutTimestamp: number;
  /** optional memo */
  memo: string;
  /** tx creator address */
  creator: string;
}

export interface Height {
  /** the revision that the client is currently on */
  revisionNumber: number;
  /** the height within the given revision */
  revisionHeight: number;
}

export interface MsgIbcTransferFromRegionTreasureResponse {
}

function createBaseMsgStake(): MsgStake {
  return { stakerAddress: "", validatorAddress: "", amount: undefined };
}

export const MsgStake = {
  encode(message: MsgStake, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.stakerAddress !== "") {
      writer.uint32(10).string(message.stakerAddress);
    }
    if (message.validatorAddress !== "") {
      writer.uint32(18).string(message.validatorAddress);
    }
    if (message.amount !== undefined) {
      Coin.encode(message.amount, writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgStake {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgStake();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.stakerAddress = reader.string();
          break;
        case 2:
          message.validatorAddress = reader.string();
          break;
        case 3:
          message.amount = Coin.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgStake {
    return {
      stakerAddress: isSet(object.stakerAddress) ? String(object.stakerAddress) : "",
      validatorAddress: isSet(object.validatorAddress) ? String(object.validatorAddress) : "",
      amount: isSet(object.amount) ? Coin.fromJSON(object.amount) : undefined,
    };
  },

  toJSON(message: MsgStake): unknown {
    const obj: any = {};
    message.stakerAddress !== undefined && (obj.stakerAddress = message.stakerAddress);
    message.validatorAddress !== undefined && (obj.validatorAddress = message.validatorAddress);
    message.amount !== undefined && (obj.amount = message.amount ? Coin.toJSON(message.amount) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgStake>, I>>(object: I): MsgStake {
    const message = createBaseMsgStake();
    message.stakerAddress = object.stakerAddress ?? "";
    message.validatorAddress = object.validatorAddress ?? "";
    message.amount = (object.amount !== undefined && object.amount !== null)
      ? Coin.fromPartial(object.amount)
      : undefined;
    return message;
  },
};

function createBaseMsgStakeResponse(): MsgStakeResponse {
  return {};
}

export const MsgStakeResponse = {
  encode(_: MsgStakeResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgStakeResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgStakeResponse();
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

  fromJSON(_: any): MsgStakeResponse {
    return {};
  },

  toJSON(_: MsgStakeResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgStakeResponse>, I>>(_: I): MsgStakeResponse {
    const message = createBaseMsgStakeResponse();
    return message;
  },
};

function createBaseMsgUnstake(): MsgUnstake {
  return { stakerAddress: "", validatorAddress: "", amount: undefined };
}

export const MsgUnstake = {
  encode(message: MsgUnstake, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.stakerAddress !== "") {
      writer.uint32(10).string(message.stakerAddress);
    }
    if (message.validatorAddress !== "") {
      writer.uint32(18).string(message.validatorAddress);
    }
    if (message.amount !== undefined) {
      Coin.encode(message.amount, writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUnstake {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUnstake();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.stakerAddress = reader.string();
          break;
        case 2:
          message.validatorAddress = reader.string();
          break;
        case 3:
          message.amount = Coin.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgUnstake {
    return {
      stakerAddress: isSet(object.stakerAddress) ? String(object.stakerAddress) : "",
      validatorAddress: isSet(object.validatorAddress) ? String(object.validatorAddress) : "",
      amount: isSet(object.amount) ? Coin.fromJSON(object.amount) : undefined,
    };
  },

  toJSON(message: MsgUnstake): unknown {
    const obj: any = {};
    message.stakerAddress !== undefined && (obj.stakerAddress = message.stakerAddress);
    message.validatorAddress !== undefined && (obj.validatorAddress = message.validatorAddress);
    message.amount !== undefined && (obj.amount = message.amount ? Coin.toJSON(message.amount) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgUnstake>, I>>(object: I): MsgUnstake {
    const message = createBaseMsgUnstake();
    message.stakerAddress = object.stakerAddress ?? "";
    message.validatorAddress = object.validatorAddress ?? "";
    message.amount = (object.amount !== undefined && object.amount !== null)
      ? Coin.fromPartial(object.amount)
      : undefined;
    return message;
  },
};

function createBaseMsgUnstakeResponse(): MsgUnstakeResponse {
  return { completionTime: undefined };
}

export const MsgUnstakeResponse = {
  encode(message: MsgUnstakeResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.completionTime !== undefined) {
      Timestamp.encode(toTimestamp(message.completionTime), writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUnstakeResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUnstakeResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.completionTime = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgUnstakeResponse {
    return { completionTime: isSet(object.completionTime) ? fromJsonTimestamp(object.completionTime) : undefined };
  },

  toJSON(message: MsgUnstakeResponse): unknown {
    const obj: any = {};
    message.completionTime !== undefined && (obj.completionTime = message.completionTime.toISOString());
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgUnstakeResponse>, I>>(object: I): MsgUnstakeResponse {
    const message = createBaseMsgUnstakeResponse();
    message.completionTime = object.completionTime ?? undefined;
    return message;
  },
};

function createBaseMsgFixedDepositCfg(): MsgFixedDepositCfg {
  return { admin: "", regionId: "", term: 0, rate: "" };
}

export const MsgFixedDepositCfg = {
  encode(message: MsgFixedDepositCfg, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.admin !== "") {
      writer.uint32(10).string(message.admin);
    }
    if (message.regionId !== "") {
      writer.uint32(18).string(message.regionId);
    }
    if (message.term !== 0) {
      writer.uint32(24).int64(message.term);
    }
    if (message.rate !== "") {
      writer.uint32(34).string(message.rate);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgFixedDepositCfg {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgFixedDepositCfg();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.admin = reader.string();
          break;
        case 2:
          message.regionId = reader.string();
          break;
        case 3:
          message.term = longToNumber(reader.int64() as Long);
          break;
        case 4:
          message.rate = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgFixedDepositCfg {
    return {
      admin: isSet(object.admin) ? String(object.admin) : "",
      regionId: isSet(object.regionId) ? String(object.regionId) : "",
      term: isSet(object.term) ? Number(object.term) : 0,
      rate: isSet(object.rate) ? String(object.rate) : "",
    };
  },

  toJSON(message: MsgFixedDepositCfg): unknown {
    const obj: any = {};
    message.admin !== undefined && (obj.admin = message.admin);
    message.regionId !== undefined && (obj.regionId = message.regionId);
    message.term !== undefined && (obj.term = Math.round(message.term));
    message.rate !== undefined && (obj.rate = message.rate);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgFixedDepositCfg>, I>>(object: I): MsgFixedDepositCfg {
    const message = createBaseMsgFixedDepositCfg();
    message.admin = object.admin ?? "";
    message.regionId = object.regionId ?? "";
    message.term = object.term ?? 0;
    message.rate = object.rate ?? "";
    return message;
  },
};

function createBaseMsgFixedDepositCfgResp(): MsgFixedDepositCfgResp {
  return { retcode: "" };
}

export const MsgFixedDepositCfgResp = {
  encode(message: MsgFixedDepositCfgResp, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.retcode !== "") {
      writer.uint32(10).string(message.retcode);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgFixedDepositCfgResp {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgFixedDepositCfgResp();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.retcode = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgFixedDepositCfgResp {
    return { retcode: isSet(object.retcode) ? String(object.retcode) : "" };
  },

  toJSON(message: MsgFixedDepositCfgResp): unknown {
    const obj: any = {};
    message.retcode !== undefined && (obj.retcode = message.retcode);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgFixedDepositCfgResp>, I>>(object: I): MsgFixedDepositCfgResp {
    const message = createBaseMsgFixedDepositCfgResp();
    message.retcode = object.retcode ?? "";
    return message;
  },
};

function createBaseMsgRemoveFixedDepositCfg(): MsgRemoveFixedDepositCfg {
  return { admin: "", regionId: "", term: 0 };
}

export const MsgRemoveFixedDepositCfg = {
  encode(message: MsgRemoveFixedDepositCfg, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.admin !== "") {
      writer.uint32(10).string(message.admin);
    }
    if (message.regionId !== "") {
      writer.uint32(18).string(message.regionId);
    }
    if (message.term !== 0) {
      writer.uint32(24).int64(message.term);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgRemoveFixedDepositCfg {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgRemoveFixedDepositCfg();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.admin = reader.string();
          break;
        case 2:
          message.regionId = reader.string();
          break;
        case 3:
          message.term = longToNumber(reader.int64() as Long);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgRemoveFixedDepositCfg {
    return {
      admin: isSet(object.admin) ? String(object.admin) : "",
      regionId: isSet(object.regionId) ? String(object.regionId) : "",
      term: isSet(object.term) ? Number(object.term) : 0,
    };
  },

  toJSON(message: MsgRemoveFixedDepositCfg): unknown {
    const obj: any = {};
    message.admin !== undefined && (obj.admin = message.admin);
    message.regionId !== undefined && (obj.regionId = message.regionId);
    message.term !== undefined && (obj.term = Math.round(message.term));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgRemoveFixedDepositCfg>, I>>(object: I): MsgRemoveFixedDepositCfg {
    const message = createBaseMsgRemoveFixedDepositCfg();
    message.admin = object.admin ?? "";
    message.regionId = object.regionId ?? "";
    message.term = object.term ?? 0;
    return message;
  },
};

function createBaseMsgRemoveFixedDepositCfgResp(): MsgRemoveFixedDepositCfgResp {
  return { retcode: "" };
}

export const MsgRemoveFixedDepositCfgResp = {
  encode(message: MsgRemoveFixedDepositCfgResp, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.retcode !== "") {
      writer.uint32(10).string(message.retcode);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgRemoveFixedDepositCfgResp {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgRemoveFixedDepositCfgResp();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.retcode = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgRemoveFixedDepositCfgResp {
    return { retcode: isSet(object.retcode) ? String(object.retcode) : "" };
  },

  toJSON(message: MsgRemoveFixedDepositCfgResp): unknown {
    const obj: any = {};
    message.retcode !== undefined && (obj.retcode = message.retcode);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgRemoveFixedDepositCfgResp>, I>>(object: I): MsgRemoveFixedDepositCfgResp {
    const message = createBaseMsgRemoveFixedDepositCfgResp();
    message.retcode = object.retcode ?? "";
    return message;
  },
};

function createBaseMsgSetFixedDepositCfgStatus(): MsgSetFixedDepositCfgStatus {
  return { admin: "", regionId: "", term: 0, status: 0 };
}

export const MsgSetFixedDepositCfgStatus = {
  encode(message: MsgSetFixedDepositCfgStatus, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.admin !== "") {
      writer.uint32(10).string(message.admin);
    }
    if (message.regionId !== "") {
      writer.uint32(18).string(message.regionId);
    }
    if (message.term !== 0) {
      writer.uint32(24).int64(message.term);
    }
    if (message.status !== 0) {
      writer.uint32(32).int32(message.status);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgSetFixedDepositCfgStatus {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgSetFixedDepositCfgStatus();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.admin = reader.string();
          break;
        case 2:
          message.regionId = reader.string();
          break;
        case 3:
          message.term = longToNumber(reader.int64() as Long);
          break;
        case 4:
          message.status = reader.int32() as any;
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgSetFixedDepositCfgStatus {
    return {
      admin: isSet(object.admin) ? String(object.admin) : "",
      regionId: isSet(object.regionId) ? String(object.regionId) : "",
      term: isSet(object.term) ? Number(object.term) : 0,
      status: isSet(object.status) ? fixedDepositCfgStatusFromJSON(object.status) : 0,
    };
  },

  toJSON(message: MsgSetFixedDepositCfgStatus): unknown {
    const obj: any = {};
    message.admin !== undefined && (obj.admin = message.admin);
    message.regionId !== undefined && (obj.regionId = message.regionId);
    message.term !== undefined && (obj.term = Math.round(message.term));
    message.status !== undefined && (obj.status = fixedDepositCfgStatusToJSON(message.status));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgSetFixedDepositCfgStatus>, I>>(object: I): MsgSetFixedDepositCfgStatus {
    const message = createBaseMsgSetFixedDepositCfgStatus();
    message.admin = object.admin ?? "";
    message.regionId = object.regionId ?? "";
    message.term = object.term ?? 0;
    message.status = object.status ?? 0;
    return message;
  },
};

function createBaseMsgSetFixedDepositCfgStatusResp(): MsgSetFixedDepositCfgStatusResp {
  return { retcode: "" };
}

export const MsgSetFixedDepositCfgStatusResp = {
  encode(message: MsgSetFixedDepositCfgStatusResp, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.retcode !== "") {
      writer.uint32(10).string(message.retcode);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgSetFixedDepositCfgStatusResp {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgSetFixedDepositCfgStatusResp();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.retcode = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgSetFixedDepositCfgStatusResp {
    return { retcode: isSet(object.retcode) ? String(object.retcode) : "" };
  },

  toJSON(message: MsgSetFixedDepositCfgStatusResp): unknown {
    const obj: any = {};
    message.retcode !== undefined && (obj.retcode = message.retcode);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgSetFixedDepositCfgStatusResp>, I>>(
    object: I,
  ): MsgSetFixedDepositCfgStatusResp {
    const message = createBaseMsgSetFixedDepositCfgStatusResp();
    message.retcode = object.retcode ?? "";
    return message;
  },
};

function createBaseMsgSetFixedDepositCfgRate(): MsgSetFixedDepositCfgRate {
  return { admin: "", regionId: "", term: 0, rate: "" };
}

export const MsgSetFixedDepositCfgRate = {
  encode(message: MsgSetFixedDepositCfgRate, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.admin !== "") {
      writer.uint32(10).string(message.admin);
    }
    if (message.regionId !== "") {
      writer.uint32(18).string(message.regionId);
    }
    if (message.term !== 0) {
      writer.uint32(24).int64(message.term);
    }
    if (message.rate !== "") {
      writer.uint32(34).string(message.rate);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgSetFixedDepositCfgRate {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgSetFixedDepositCfgRate();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.admin = reader.string();
          break;
        case 2:
          message.regionId = reader.string();
          break;
        case 3:
          message.term = longToNumber(reader.int64() as Long);
          break;
        case 4:
          message.rate = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgSetFixedDepositCfgRate {
    return {
      admin: isSet(object.admin) ? String(object.admin) : "",
      regionId: isSet(object.regionId) ? String(object.regionId) : "",
      term: isSet(object.term) ? Number(object.term) : 0,
      rate: isSet(object.rate) ? String(object.rate) : "",
    };
  },

  toJSON(message: MsgSetFixedDepositCfgRate): unknown {
    const obj: any = {};
    message.admin !== undefined && (obj.admin = message.admin);
    message.regionId !== undefined && (obj.regionId = message.regionId);
    message.term !== undefined && (obj.term = Math.round(message.term));
    message.rate !== undefined && (obj.rate = message.rate);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgSetFixedDepositCfgRate>, I>>(object: I): MsgSetFixedDepositCfgRate {
    const message = createBaseMsgSetFixedDepositCfgRate();
    message.admin = object.admin ?? "";
    message.regionId = object.regionId ?? "";
    message.term = object.term ?? 0;
    message.rate = object.rate ?? "";
    return message;
  },
};

function createBaseMsgSetFixedDepositCfgRateResp(): MsgSetFixedDepositCfgRateResp {
  return { retcode: "" };
}

export const MsgSetFixedDepositCfgRateResp = {
  encode(message: MsgSetFixedDepositCfgRateResp, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.retcode !== "") {
      writer.uint32(10).string(message.retcode);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgSetFixedDepositCfgRateResp {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgSetFixedDepositCfgRateResp();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.retcode = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgSetFixedDepositCfgRateResp {
    return { retcode: isSet(object.retcode) ? String(object.retcode) : "" };
  },

  toJSON(message: MsgSetFixedDepositCfgRateResp): unknown {
    const obj: any = {};
    message.retcode !== undefined && (obj.retcode = message.retcode);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgSetFixedDepositCfgRateResp>, I>>(
    object: I,
  ): MsgSetFixedDepositCfgRateResp {
    const message = createBaseMsgSetFixedDepositCfgRateResp();
    message.retcode = object.retcode ?? "";
    return message;
  },
};

function createBaseMsgNewRegion(): MsgNewRegion {
  return { creator: "", name: "", operatorAddress: "" };
}

export const MsgNewRegion = {
  encode(message: MsgNewRegion, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.name !== "") {
      writer.uint32(18).string(message.name);
    }
    if (message.operatorAddress !== "") {
      writer.uint32(26).string(message.operatorAddress);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgNewRegion {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgNewRegion();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.name = reader.string();
          break;
        case 3:
          message.operatorAddress = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgNewRegion {
    return {
      creator: isSet(object.creator) ? String(object.creator) : "",
      name: isSet(object.name) ? String(object.name) : "",
      operatorAddress: isSet(object.operatorAddress) ? String(object.operatorAddress) : "",
    };
  },

  toJSON(message: MsgNewRegion): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    message.name !== undefined && (obj.name = message.name);
    message.operatorAddress !== undefined && (obj.operatorAddress = message.operatorAddress);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgNewRegion>, I>>(object: I): MsgNewRegion {
    const message = createBaseMsgNewRegion();
    message.creator = object.creator ?? "";
    message.name = object.name ?? "";
    message.operatorAddress = object.operatorAddress ?? "";
    return message;
  },
};

function createBaseMsgNewRegionResponse(): MsgNewRegionResponse {
  return { regionId: "" };
}

export const MsgNewRegionResponse = {
  encode(message: MsgNewRegionResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.regionId !== "") {
      writer.uint32(10).string(message.regionId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgNewRegionResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgNewRegionResponse();
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

  fromJSON(object: any): MsgNewRegionResponse {
    return { regionId: isSet(object.regionId) ? String(object.regionId) : "" };
  },

  toJSON(message: MsgNewRegionResponse): unknown {
    const obj: any = {};
    message.regionId !== undefined && (obj.regionId = message.regionId);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgNewRegionResponse>, I>>(object: I): MsgNewRegionResponse {
    const message = createBaseMsgNewRegionResponse();
    message.regionId = object.regionId ?? "";
    return message;
  },
};

function createBaseMsgRemoveRegion(): MsgRemoveRegion {
  return { creator: "", regionId: "" };
}

export const MsgRemoveRegion = {
  encode(message: MsgRemoveRegion, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.regionId !== "") {
      writer.uint32(18).string(message.regionId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgRemoveRegion {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgRemoveRegion();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.regionId = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgRemoveRegion {
    return {
      creator: isSet(object.creator) ? String(object.creator) : "",
      regionId: isSet(object.regionId) ? String(object.regionId) : "",
    };
  },

  toJSON(message: MsgRemoveRegion): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    message.regionId !== undefined && (obj.regionId = message.regionId);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgRemoveRegion>, I>>(object: I): MsgRemoveRegion {
    const message = createBaseMsgRemoveRegion();
    message.creator = object.creator ?? "";
    message.regionId = object.regionId ?? "";
    return message;
  },
};

function createBaseMsgRemoveRegionResponse(): MsgRemoveRegionResponse {
  return { retcode: "" };
}

export const MsgRemoveRegionResponse = {
  encode(message: MsgRemoveRegionResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.retcode !== "") {
      writer.uint32(10).string(message.retcode);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgRemoveRegionResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgRemoveRegionResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.retcode = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgRemoveRegionResponse {
    return { retcode: isSet(object.retcode) ? String(object.retcode) : "" };
  },

  toJSON(message: MsgRemoveRegionResponse): unknown {
    const obj: any = {};
    message.retcode !== undefined && (obj.retcode = message.retcode);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgRemoveRegionResponse>, I>>(object: I): MsgRemoveRegionResponse {
    const message = createBaseMsgRemoveRegionResponse();
    message.retcode = object.retcode ?? "";
    return message;
  },
};

function createBaseMsgWithdrawFromRegion(): MsgWithdrawFromRegion {
  return { withdrawer: "", regionId: "", receiver: "", amount: [] };
}

export const MsgWithdrawFromRegion = {
  encode(message: MsgWithdrawFromRegion, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.withdrawer !== "") {
      writer.uint32(10).string(message.withdrawer);
    }
    if (message.regionId !== "") {
      writer.uint32(18).string(message.regionId);
    }
    if (message.receiver !== "") {
      writer.uint32(26).string(message.receiver);
    }
    for (const v of message.amount) {
      Coin.encode(v!, writer.uint32(34).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgWithdrawFromRegion {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgWithdrawFromRegion();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.withdrawer = reader.string();
          break;
        case 2:
          message.regionId = reader.string();
          break;
        case 3:
          message.receiver = reader.string();
          break;
        case 4:
          message.amount.push(Coin.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgWithdrawFromRegion {
    return {
      withdrawer: isSet(object.withdrawer) ? String(object.withdrawer) : "",
      regionId: isSet(object.regionId) ? String(object.regionId) : "",
      receiver: isSet(object.receiver) ? String(object.receiver) : "",
      amount: Array.isArray(object?.amount) ? object.amount.map((e: any) => Coin.fromJSON(e)) : [],
    };
  },

  toJSON(message: MsgWithdrawFromRegion): unknown {
    const obj: any = {};
    message.withdrawer !== undefined && (obj.withdrawer = message.withdrawer);
    message.regionId !== undefined && (obj.regionId = message.regionId);
    message.receiver !== undefined && (obj.receiver = message.receiver);
    if (message.amount) {
      obj.amount = message.amount.map((e) => e ? Coin.toJSON(e) : undefined);
    } else {
      obj.amount = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgWithdrawFromRegion>, I>>(object: I): MsgWithdrawFromRegion {
    const message = createBaseMsgWithdrawFromRegion();
    message.withdrawer = object.withdrawer ?? "";
    message.regionId = object.regionId ?? "";
    message.receiver = object.receiver ?? "";
    message.amount = object.amount?.map((e) => Coin.fromPartial(e)) || [];
    return message;
  },
};

function createBaseMsgWithdrawFromRegionResp(): MsgWithdrawFromRegionResp {
  return {};
}

export const MsgWithdrawFromRegionResp = {
  encode(_: MsgWithdrawFromRegionResp, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgWithdrawFromRegionResp {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgWithdrawFromRegionResp();
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

  fromJSON(_: any): MsgWithdrawFromRegionResp {
    return {};
  },

  toJSON(_: MsgWithdrawFromRegionResp): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgWithdrawFromRegionResp>, I>>(_: I): MsgWithdrawFromRegionResp {
    const message = createBaseMsgWithdrawFromRegionResp();
    return message;
  },
};

function createBaseMsgWithdrawFromGlobalDaoFeePool(): MsgWithdrawFromGlobalDaoFeePool {
  return { withdrawer: "", amount: [] };
}

export const MsgWithdrawFromGlobalDaoFeePool = {
  encode(message: MsgWithdrawFromGlobalDaoFeePool, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.withdrawer !== "") {
      writer.uint32(10).string(message.withdrawer);
    }
    for (const v of message.amount) {
      Coin.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgWithdrawFromGlobalDaoFeePool {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgWithdrawFromGlobalDaoFeePool();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.withdrawer = reader.string();
          break;
        case 2:
          message.amount.push(Coin.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgWithdrawFromGlobalDaoFeePool {
    return {
      withdrawer: isSet(object.withdrawer) ? String(object.withdrawer) : "",
      amount: Array.isArray(object?.amount) ? object.amount.map((e: any) => Coin.fromJSON(e)) : [],
    };
  },

  toJSON(message: MsgWithdrawFromGlobalDaoFeePool): unknown {
    const obj: any = {};
    message.withdrawer !== undefined && (obj.withdrawer = message.withdrawer);
    if (message.amount) {
      obj.amount = message.amount.map((e) => e ? Coin.toJSON(e) : undefined);
    } else {
      obj.amount = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgWithdrawFromGlobalDaoFeePool>, I>>(
    object: I,
  ): MsgWithdrawFromGlobalDaoFeePool {
    const message = createBaseMsgWithdrawFromGlobalDaoFeePool();
    message.withdrawer = object.withdrawer ?? "";
    message.amount = object.amount?.map((e) => Coin.fromPartial(e)) || [];
    return message;
  },
};

function createBaseMsgWithdrawFromGlobalDaoFeePoolResp(): MsgWithdrawFromGlobalDaoFeePoolResp {
  return {};
}

export const MsgWithdrawFromGlobalDaoFeePoolResp = {
  encode(_: MsgWithdrawFromGlobalDaoFeePoolResp, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgWithdrawFromGlobalDaoFeePoolResp {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgWithdrawFromGlobalDaoFeePoolResp();
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

  fromJSON(_: any): MsgWithdrawFromGlobalDaoFeePoolResp {
    return {};
  },

  toJSON(_: MsgWithdrawFromGlobalDaoFeePoolResp): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgWithdrawFromGlobalDaoFeePoolResp>, I>>(
    _: I,
  ): MsgWithdrawFromGlobalDaoFeePoolResp {
    const message = createBaseMsgWithdrawFromGlobalDaoFeePoolResp();
    return message;
  },
};

function createBaseMsgWithdrawDelegatorReward(): MsgWithdrawDelegatorReward {
  return { delegatorAddress: "", validatorAddress: "" };
}

export const MsgWithdrawDelegatorReward = {
  encode(message: MsgWithdrawDelegatorReward, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.delegatorAddress !== "") {
      writer.uint32(10).string(message.delegatorAddress);
    }
    if (message.validatorAddress !== "") {
      writer.uint32(18).string(message.validatorAddress);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgWithdrawDelegatorReward {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgWithdrawDelegatorReward();
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

  fromJSON(object: any): MsgWithdrawDelegatorReward {
    return {
      delegatorAddress: isSet(object.delegatorAddress) ? String(object.delegatorAddress) : "",
      validatorAddress: isSet(object.validatorAddress) ? String(object.validatorAddress) : "",
    };
  },

  toJSON(message: MsgWithdrawDelegatorReward): unknown {
    const obj: any = {};
    message.delegatorAddress !== undefined && (obj.delegatorAddress = message.delegatorAddress);
    message.validatorAddress !== undefined && (obj.validatorAddress = message.validatorAddress);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgWithdrawDelegatorReward>, I>>(object: I): MsgWithdrawDelegatorReward {
    const message = createBaseMsgWithdrawDelegatorReward();
    message.delegatorAddress = object.delegatorAddress ?? "";
    message.validatorAddress = object.validatorAddress ?? "";
    return message;
  },
};

function createBaseMsgWithdrawDelegatorRewardResponse(): MsgWithdrawDelegatorRewardResponse {
  return { amount: [] };
}

export const MsgWithdrawDelegatorRewardResponse = {
  encode(message: MsgWithdrawDelegatorRewardResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.amount) {
      Coin.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgWithdrawDelegatorRewardResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgWithdrawDelegatorRewardResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.amount.push(Coin.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgWithdrawDelegatorRewardResponse {
    return { amount: Array.isArray(object?.amount) ? object.amount.map((e: any) => Coin.fromJSON(e)) : [] };
  },

  toJSON(message: MsgWithdrawDelegatorRewardResponse): unknown {
    const obj: any = {};
    if (message.amount) {
      obj.amount = message.amount.map((e) => e ? Coin.toJSON(e) : undefined);
    } else {
      obj.amount = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgWithdrawDelegatorRewardResponse>, I>>(
    object: I,
  ): MsgWithdrawDelegatorRewardResponse {
    const message = createBaseMsgWithdrawDelegatorRewardResponse();
    message.amount = object.amount?.map((e) => Coin.fromPartial(e)) || [];
    return message;
  },
};

function createBaseMsgNewFixedDepositCfg(): MsgNewFixedDepositCfg {
  return { dao: "", regionId: "", term: 0, rate: "" };
}

export const MsgNewFixedDepositCfg = {
  encode(message: MsgNewFixedDepositCfg, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.dao !== "") {
      writer.uint32(10).string(message.dao);
    }
    if (message.regionId !== "") {
      writer.uint32(18).string(message.regionId);
    }
    if (message.term !== 0) {
      writer.uint32(24).int64(message.term);
    }
    if (message.rate !== "") {
      writer.uint32(34).string(message.rate);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgNewFixedDepositCfg {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgNewFixedDepositCfg();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.dao = reader.string();
          break;
        case 2:
          message.regionId = reader.string();
          break;
        case 3:
          message.term = longToNumber(reader.int64() as Long);
          break;
        case 4:
          message.rate = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgNewFixedDepositCfg {
    return {
      dao: isSet(object.dao) ? String(object.dao) : "",
      regionId: isSet(object.regionId) ? String(object.regionId) : "",
      term: isSet(object.term) ? Number(object.term) : 0,
      rate: isSet(object.rate) ? String(object.rate) : "",
    };
  },

  toJSON(message: MsgNewFixedDepositCfg): unknown {
    const obj: any = {};
    message.dao !== undefined && (obj.dao = message.dao);
    message.regionId !== undefined && (obj.regionId = message.regionId);
    message.term !== undefined && (obj.term = Math.round(message.term));
    message.rate !== undefined && (obj.rate = message.rate);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgNewFixedDepositCfg>, I>>(object: I): MsgNewFixedDepositCfg {
    const message = createBaseMsgNewFixedDepositCfg();
    message.dao = object.dao ?? "";
    message.regionId = object.regionId ?? "";
    message.term = object.term ?? 0;
    message.rate = object.rate ?? "";
    return message;
  },
};

function createBaseMsgNewFixedDepositCfgResp(): MsgNewFixedDepositCfgResp {
  return { retcode: "" };
}

export const MsgNewFixedDepositCfgResp = {
  encode(message: MsgNewFixedDepositCfgResp, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.retcode !== "") {
      writer.uint32(10).string(message.retcode);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgNewFixedDepositCfgResp {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgNewFixedDepositCfgResp();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.retcode = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgNewFixedDepositCfgResp {
    return { retcode: isSet(object.retcode) ? String(object.retcode) : "" };
  },

  toJSON(message: MsgNewFixedDepositCfgResp): unknown {
    const obj: any = {};
    message.retcode !== undefined && (obj.retcode = message.retcode);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgNewFixedDepositCfgResp>, I>>(object: I): MsgNewFixedDepositCfgResp {
    const message = createBaseMsgNewFixedDepositCfgResp();
    message.retcode = object.retcode ?? "";
    return message;
  },
};

function createBaseMsgDoFixedDeposit(): MsgDoFixedDeposit {
  return { account: "", principal: undefined, term: 0 };
}

export const MsgDoFixedDeposit = {
  encode(message: MsgDoFixedDeposit, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.account !== "") {
      writer.uint32(10).string(message.account);
    }
    if (message.principal !== undefined) {
      Coin.encode(message.principal, writer.uint32(18).fork()).ldelim();
    }
    if (message.term !== 0) {
      writer.uint32(32).int64(message.term);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgDoFixedDeposit {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgDoFixedDeposit();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.account = reader.string();
          break;
        case 2:
          message.principal = Coin.decode(reader, reader.uint32());
          break;
        case 4:
          message.term = longToNumber(reader.int64() as Long);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgDoFixedDeposit {
    return {
      account: isSet(object.account) ? String(object.account) : "",
      principal: isSet(object.principal) ? Coin.fromJSON(object.principal) : undefined,
      term: isSet(object.term) ? Number(object.term) : 0,
    };
  },

  toJSON(message: MsgDoFixedDeposit): unknown {
    const obj: any = {};
    message.account !== undefined && (obj.account = message.account);
    message.principal !== undefined && (obj.principal = message.principal ? Coin.toJSON(message.principal) : undefined);
    message.term !== undefined && (obj.term = Math.round(message.term));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgDoFixedDeposit>, I>>(object: I): MsgDoFixedDeposit {
    const message = createBaseMsgDoFixedDeposit();
    message.account = object.account ?? "";
    message.principal = (object.principal !== undefined && object.principal !== null)
      ? Coin.fromPartial(object.principal)
      : undefined;
    message.term = object.term ?? 0;
    return message;
  },
};

function createBaseMsgDoFixedDepositResponse(): MsgDoFixedDepositResponse {
  return { id: 0 };
}

export const MsgDoFixedDepositResponse = {
  encode(message: MsgDoFixedDepositResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== 0) {
      writer.uint32(8).uint64(message.id);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgDoFixedDepositResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgDoFixedDepositResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.id = longToNumber(reader.uint64() as Long);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgDoFixedDepositResponse {
    return { id: isSet(object.id) ? Number(object.id) : 0 };
  },

  toJSON(message: MsgDoFixedDepositResponse): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = Math.round(message.id));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgDoFixedDepositResponse>, I>>(object: I): MsgDoFixedDepositResponse {
    const message = createBaseMsgDoFixedDepositResponse();
    message.id = object.id ?? 0;
    return message;
  },
};

function createBaseMsgWithdrawFixedDeposit(): MsgWithdrawFixedDeposit {
  return { account: "", id: 0 };
}

export const MsgWithdrawFixedDeposit = {
  encode(message: MsgWithdrawFixedDeposit, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.account !== "") {
      writer.uint32(10).string(message.account);
    }
    if (message.id !== 0) {
      writer.uint32(16).uint64(message.id);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgWithdrawFixedDeposit {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgWithdrawFixedDeposit();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.account = reader.string();
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

  fromJSON(object: any): MsgWithdrawFixedDeposit {
    return {
      account: isSet(object.account) ? String(object.account) : "",
      id: isSet(object.id) ? Number(object.id) : 0,
    };
  },

  toJSON(message: MsgWithdrawFixedDeposit): unknown {
    const obj: any = {};
    message.account !== undefined && (obj.account = message.account);
    message.id !== undefined && (obj.id = Math.round(message.id));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgWithdrawFixedDeposit>, I>>(object: I): MsgWithdrawFixedDeposit {
    const message = createBaseMsgWithdrawFixedDeposit();
    message.account = object.account ?? "";
    message.id = object.id ?? 0;
    return message;
  },
};

function createBaseMsgWithdrawFixedDepositResponse(): MsgWithdrawFixedDepositResponse {
  return { principal: undefined, interest: undefined, term: 0, rate: "" };
}

export const MsgWithdrawFixedDepositResponse = {
  encode(message: MsgWithdrawFixedDepositResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.principal !== undefined) {
      Coin.encode(message.principal, writer.uint32(10).fork()).ldelim();
    }
    if (message.interest !== undefined) {
      Coin.encode(message.interest, writer.uint32(18).fork()).ldelim();
    }
    if (message.term !== 0) {
      writer.uint32(24).int64(message.term);
    }
    if (message.rate !== "") {
      writer.uint32(34).string(message.rate);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgWithdrawFixedDepositResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgWithdrawFixedDepositResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.principal = Coin.decode(reader, reader.uint32());
          break;
        case 2:
          message.interest = Coin.decode(reader, reader.uint32());
          break;
        case 3:
          message.term = longToNumber(reader.int64() as Long);
          break;
        case 4:
          message.rate = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgWithdrawFixedDepositResponse {
    return {
      principal: isSet(object.principal) ? Coin.fromJSON(object.principal) : undefined,
      interest: isSet(object.interest) ? Coin.fromJSON(object.interest) : undefined,
      term: isSet(object.term) ? Number(object.term) : 0,
      rate: isSet(object.rate) ? String(object.rate) : "",
    };
  },

  toJSON(message: MsgWithdrawFixedDepositResponse): unknown {
    const obj: any = {};
    message.principal !== undefined && (obj.principal = message.principal ? Coin.toJSON(message.principal) : undefined);
    message.interest !== undefined && (obj.interest = message.interest ? Coin.toJSON(message.interest) : undefined);
    message.term !== undefined && (obj.term = Math.round(message.term));
    message.rate !== undefined && (obj.rate = message.rate);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgWithdrawFixedDepositResponse>, I>>(
    object: I,
  ): MsgWithdrawFixedDepositResponse {
    const message = createBaseMsgWithdrawFixedDepositResponse();
    message.principal = (object.principal !== undefined && object.principal !== null)
      ? Coin.fromPartial(object.principal)
      : undefined;
    message.interest = (object.interest !== undefined && object.interest !== null)
      ? Coin.fromPartial(object.interest)
      : undefined;
    message.term = object.term ?? 0;
    message.rate = object.rate ?? "";
    return message;
  },
};

function createBaseMsgNewRecord(): MsgNewRecord {
  return { actionNumber: "", actionUrl: "", from: "" };
}

export const MsgNewRecord = {
  encode(message: MsgNewRecord, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.actionNumber !== "") {
      writer.uint32(10).string(message.actionNumber);
    }
    if (message.actionUrl !== "") {
      writer.uint32(18).string(message.actionUrl);
    }
    if (message.from !== "") {
      writer.uint32(26).string(message.from);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgNewRecord {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgNewRecord();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.actionNumber = reader.string();
          break;
        case 2:
          message.actionUrl = reader.string();
          break;
        case 3:
          message.from = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgNewRecord {
    return {
      actionNumber: isSet(object.actionNumber) ? String(object.actionNumber) : "",
      actionUrl: isSet(object.actionUrl) ? String(object.actionUrl) : "",
      from: isSet(object.from) ? String(object.from) : "",
    };
  },

  toJSON(message: MsgNewRecord): unknown {
    const obj: any = {};
    message.actionNumber !== undefined && (obj.actionNumber = message.actionNumber);
    message.actionUrl !== undefined && (obj.actionUrl = message.actionUrl);
    message.from !== undefined && (obj.from = message.from);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgNewRecord>, I>>(object: I): MsgNewRecord {
    const message = createBaseMsgNewRecord();
    message.actionNumber = object.actionNumber ?? "";
    message.actionUrl = object.actionUrl ?? "";
    message.from = object.from ?? "";
    return message;
  },
};

function createBaseMsgNewRecordResponse(): MsgNewRecordResponse {
  return {};
}

export const MsgNewRecordResponse = {
  encode(_: MsgNewRecordResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgNewRecordResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgNewRecordResponse();
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

  fromJSON(_: any): MsgNewRecordResponse {
    return {};
  },

  toJSON(_: MsgNewRecordResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgNewRecordResponse>, I>>(_: I): MsgNewRecordResponse {
    const message = createBaseMsgNewRecordResponse();
    return message;
  },
};

function createBaseMsgReviewRecord(): MsgReviewRecord {
  return { recordHash: "", reviewResult: "", from: "", ActionNumber: "", reviewedAddress: "" };
}

export const MsgReviewRecord = {
  encode(message: MsgReviewRecord, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.recordHash !== "") {
      writer.uint32(10).string(message.recordHash);
    }
    if (message.reviewResult !== "") {
      writer.uint32(18).string(message.reviewResult);
    }
    if (message.from !== "") {
      writer.uint32(26).string(message.from);
    }
    if (message.ActionNumber !== "") {
      writer.uint32(34).string(message.ActionNumber);
    }
    if (message.reviewedAddress !== "") {
      writer.uint32(42).string(message.reviewedAddress);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgReviewRecord {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgReviewRecord();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.recordHash = reader.string();
          break;
        case 2:
          message.reviewResult = reader.string();
          break;
        case 3:
          message.from = reader.string();
          break;
        case 4:
          message.ActionNumber = reader.string();
          break;
        case 5:
          message.reviewedAddress = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgReviewRecord {
    return {
      recordHash: isSet(object.recordHash) ? String(object.recordHash) : "",
      reviewResult: isSet(object.reviewResult) ? String(object.reviewResult) : "",
      from: isSet(object.from) ? String(object.from) : "",
      ActionNumber: isSet(object.ActionNumber) ? String(object.ActionNumber) : "",
      reviewedAddress: isSet(object.reviewedAddress) ? String(object.reviewedAddress) : "",
    };
  },

  toJSON(message: MsgReviewRecord): unknown {
    const obj: any = {};
    message.recordHash !== undefined && (obj.recordHash = message.recordHash);
    message.reviewResult !== undefined && (obj.reviewResult = message.reviewResult);
    message.from !== undefined && (obj.from = message.from);
    message.ActionNumber !== undefined && (obj.ActionNumber = message.ActionNumber);
    message.reviewedAddress !== undefined && (obj.reviewedAddress = message.reviewedAddress);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgReviewRecord>, I>>(object: I): MsgReviewRecord {
    const message = createBaseMsgReviewRecord();
    message.recordHash = object.recordHash ?? "";
    message.reviewResult = object.reviewResult ?? "";
    message.from = object.from ?? "";
    message.ActionNumber = object.ActionNumber ?? "";
    message.reviewedAddress = object.reviewedAddress ?? "";
    return message;
  },
};

function createBaseMsgReviewRecordResponse(): MsgReviewRecordResponse {
  return {};
}

export const MsgReviewRecordResponse = {
  encode(_: MsgReviewRecordResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgReviewRecordResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgReviewRecordResponse();
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

  fromJSON(_: any): MsgReviewRecordResponse {
    return {};
  },

  toJSON(_: MsgReviewRecordResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgReviewRecordResponse>, I>>(_: I): MsgReviewRecordResponse {
    const message = createBaseMsgReviewRecordResponse();
    return message;
  },
};

function createBaseMsgTransferRegion(): MsgTransferRegion {
  return { fromRegion: "", toRegion: "", address: [], creator: "" };
}

export const MsgTransferRegion = {
  encode(message: MsgTransferRegion, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.fromRegion !== "") {
      writer.uint32(10).string(message.fromRegion);
    }
    if (message.toRegion !== "") {
      writer.uint32(18).string(message.toRegion);
    }
    for (const v of message.address) {
      writer.uint32(26).string(v!);
    }
    if (message.creator !== "") {
      writer.uint32(34).string(message.creator);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgTransferRegion {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgTransferRegion();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.fromRegion = reader.string();
          break;
        case 2:
          message.toRegion = reader.string();
          break;
        case 3:
          message.address.push(reader.string());
          break;
        case 4:
          message.creator = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgTransferRegion {
    return {
      fromRegion: isSet(object.fromRegion) ? String(object.fromRegion) : "",
      toRegion: isSet(object.toRegion) ? String(object.toRegion) : "",
      address: Array.isArray(object?.address) ? object.address.map((e: any) => String(e)) : [],
      creator: isSet(object.creator) ? String(object.creator) : "",
    };
  },

  toJSON(message: MsgTransferRegion): unknown {
    const obj: any = {};
    message.fromRegion !== undefined && (obj.fromRegion = message.fromRegion);
    message.toRegion !== undefined && (obj.toRegion = message.toRegion);
    if (message.address) {
      obj.address = message.address.map((e) => e);
    } else {
      obj.address = [];
    }
    message.creator !== undefined && (obj.creator = message.creator);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgTransferRegion>, I>>(object: I): MsgTransferRegion {
    const message = createBaseMsgTransferRegion();
    message.fromRegion = object.fromRegion ?? "";
    message.toRegion = object.toRegion ?? "";
    message.address = object.address?.map((e) => e) || [];
    message.creator = object.creator ?? "";
    return message;
  },
};

function createBaseMsgTransferRegionResponse(): MsgTransferRegionResponse {
  return {};
}

export const MsgTransferRegionResponse = {
  encode(_: MsgTransferRegionResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgTransferRegionResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgTransferRegionResponse();
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

  fromJSON(_: any): MsgTransferRegionResponse {
    return {};
  },

  toJSON(_: MsgTransferRegionResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgTransferRegionResponse>, I>>(_: I): MsgTransferRegionResponse {
    const message = createBaseMsgTransferRegionResponse();
    return message;
  },
};

function createBaseMsgIbcTransferFromRegionTreasure(): MsgIbcTransferFromRegionTreasure {
  return {
    sourcePort: "",
    sourceChannel: "",
    regionId: "",
    token: undefined,
    timeoutHeight: undefined,
    timeoutTimestamp: 0,
    memo: "",
    creator: "",
  };
}

export const MsgIbcTransferFromRegionTreasure = {
  encode(message: MsgIbcTransferFromRegionTreasure, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.sourcePort !== "") {
      writer.uint32(10).string(message.sourcePort);
    }
    if (message.sourceChannel !== "") {
      writer.uint32(18).string(message.sourceChannel);
    }
    if (message.regionId !== "") {
      writer.uint32(26).string(message.regionId);
    }
    if (message.token !== undefined) {
      Coin.encode(message.token, writer.uint32(34).fork()).ldelim();
    }
    if (message.timeoutHeight !== undefined) {
      Height.encode(message.timeoutHeight, writer.uint32(50).fork()).ldelim();
    }
    if (message.timeoutTimestamp !== 0) {
      writer.uint32(56).uint64(message.timeoutTimestamp);
    }
    if (message.memo !== "") {
      writer.uint32(66).string(message.memo);
    }
    if (message.creator !== "") {
      writer.uint32(74).string(message.creator);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgIbcTransferFromRegionTreasure {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgIbcTransferFromRegionTreasure();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.sourcePort = reader.string();
          break;
        case 2:
          message.sourceChannel = reader.string();
          break;
        case 3:
          message.regionId = reader.string();
          break;
        case 4:
          message.token = Coin.decode(reader, reader.uint32());
          break;
        case 6:
          message.timeoutHeight = Height.decode(reader, reader.uint32());
          break;
        case 7:
          message.timeoutTimestamp = longToNumber(reader.uint64() as Long);
          break;
        case 8:
          message.memo = reader.string();
          break;
        case 9:
          message.creator = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgIbcTransferFromRegionTreasure {
    return {
      sourcePort: isSet(object.sourcePort) ? String(object.sourcePort) : "",
      sourceChannel: isSet(object.sourceChannel) ? String(object.sourceChannel) : "",
      regionId: isSet(object.regionId) ? String(object.regionId) : "",
      token: isSet(object.token) ? Coin.fromJSON(object.token) : undefined,
      timeoutHeight: isSet(object.timeoutHeight) ? Height.fromJSON(object.timeoutHeight) : undefined,
      timeoutTimestamp: isSet(object.timeoutTimestamp) ? Number(object.timeoutTimestamp) : 0,
      memo: isSet(object.memo) ? String(object.memo) : "",
      creator: isSet(object.creator) ? String(object.creator) : "",
    };
  },

  toJSON(message: MsgIbcTransferFromRegionTreasure): unknown {
    const obj: any = {};
    message.sourcePort !== undefined && (obj.sourcePort = message.sourcePort);
    message.sourceChannel !== undefined && (obj.sourceChannel = message.sourceChannel);
    message.regionId !== undefined && (obj.regionId = message.regionId);
    message.token !== undefined && (obj.token = message.token ? Coin.toJSON(message.token) : undefined);
    message.timeoutHeight !== undefined
      && (obj.timeoutHeight = message.timeoutHeight ? Height.toJSON(message.timeoutHeight) : undefined);
    message.timeoutTimestamp !== undefined && (obj.timeoutTimestamp = Math.round(message.timeoutTimestamp));
    message.memo !== undefined && (obj.memo = message.memo);
    message.creator !== undefined && (obj.creator = message.creator);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgIbcTransferFromRegionTreasure>, I>>(
    object: I,
  ): MsgIbcTransferFromRegionTreasure {
    const message = createBaseMsgIbcTransferFromRegionTreasure();
    message.sourcePort = object.sourcePort ?? "";
    message.sourceChannel = object.sourceChannel ?? "";
    message.regionId = object.regionId ?? "";
    message.token = (object.token !== undefined && object.token !== null) ? Coin.fromPartial(object.token) : undefined;
    message.timeoutHeight = (object.timeoutHeight !== undefined && object.timeoutHeight !== null)
      ? Height.fromPartial(object.timeoutHeight)
      : undefined;
    message.timeoutTimestamp = object.timeoutTimestamp ?? 0;
    message.memo = object.memo ?? "";
    message.creator = object.creator ?? "";
    return message;
  },
};

function createBaseHeight(): Height {
  return { revisionNumber: 0, revisionHeight: 0 };
}

export const Height = {
  encode(message: Height, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.revisionNumber !== 0) {
      writer.uint32(8).uint64(message.revisionNumber);
    }
    if (message.revisionHeight !== 0) {
      writer.uint32(16).uint64(message.revisionHeight);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Height {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseHeight();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.revisionNumber = longToNumber(reader.uint64() as Long);
          break;
        case 2:
          message.revisionHeight = longToNumber(reader.uint64() as Long);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): Height {
    return {
      revisionNumber: isSet(object.revisionNumber) ? Number(object.revisionNumber) : 0,
      revisionHeight: isSet(object.revisionHeight) ? Number(object.revisionHeight) : 0,
    };
  },

  toJSON(message: Height): unknown {
    const obj: any = {};
    message.revisionNumber !== undefined && (obj.revisionNumber = Math.round(message.revisionNumber));
    message.revisionHeight !== undefined && (obj.revisionHeight = Math.round(message.revisionHeight));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<Height>, I>>(object: I): Height {
    const message = createBaseHeight();
    message.revisionNumber = object.revisionNumber ?? 0;
    message.revisionHeight = object.revisionHeight ?? 0;
    return message;
  },
};

function createBaseMsgIbcTransferFromRegionTreasureResponse(): MsgIbcTransferFromRegionTreasureResponse {
  return {};
}

export const MsgIbcTransferFromRegionTreasureResponse = {
  encode(_: MsgIbcTransferFromRegionTreasureResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgIbcTransferFromRegionTreasureResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgIbcTransferFromRegionTreasureResponse();
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

  fromJSON(_: any): MsgIbcTransferFromRegionTreasureResponse {
    return {};
  },

  toJSON(_: MsgIbcTransferFromRegionTreasureResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgIbcTransferFromRegionTreasureResponse>, I>>(
    _: I,
  ): MsgIbcTransferFromRegionTreasureResponse {
    const message = createBaseMsgIbcTransferFromRegionTreasureResponse();
    return message;
  },
};

/** Msg defines the staking Msg service. */
export interface Msg {
  Stake(request: MsgStake): Promise<MsgStakeResponse>;
  Unstake(request: MsgUnstake): Promise<MsgUnstakeResponse>;
  NewRegion(request: MsgNewRegion): Promise<MsgNewRegionResponse>;
  RemoveRegion(request: MsgRemoveRegion): Promise<MsgRemoveRegionResponse>;
  /**
   * rpc UnmeidWithdrawDelegatorReward(MsgUnmeidWithdrawDelegatorReward)
   *  returns (MsgUnmeidWithdrawDelegatorRewardResponse);
   */
  WithdrawDelegatorReward(request: MsgWithdrawDelegatorReward): Promise<MsgWithdrawDelegatorRewardResponse>;
  /**
   * rpc RetrieveCoinsFromRegion(MsgRetrieveCoinsFromRegion) returns
   *  (MsgRetrieveCoinsFromRegionResp);
   *  rpc TransferRegion(MsgTransferRegion) returns (MsgTransferRegionResponse);
   *  rpc
   *  RetrieveFeeFromGlobalAdminFeePool(MsgRetrieveFeeFromGlobalAdminFeePool)
   *  returns (MsgRetrieveFeeFromGlobalAdminFeePoolResp);
   */
  NewFixedDepositCfg(request: MsgNewFixedDepositCfg): Promise<MsgNewFixedDepositCfgResp>;
  SetFixedDepositCfgStatus(request: MsgSetFixedDepositCfgStatus): Promise<MsgSetFixedDepositCfgStatusResp>;
  SetFixedDepositCfgRate(request: MsgSetFixedDepositCfgRate): Promise<MsgSetFixedDepositCfgRateResp>;
  DoFixedDeposit(request: MsgDoFixedDeposit): Promise<MsgDoFixedDepositResponse>;
  WithdrawFixedDeposit(request: MsgWithdrawFixedDeposit): Promise<MsgWithdrawFixedDepositResponse>;
  RemoveFixedDepositCfg(request: MsgRemoveFixedDepositCfg): Promise<MsgRemoveFixedDepositCfgResp>;
  WithdrawFromRegion(request: MsgWithdrawFromRegion): Promise<MsgWithdrawFromRegionResp>;
  WithdrawFromGlobalDaoFeePool(request: MsgWithdrawFromGlobalDaoFeePool): Promise<MsgWithdrawFromGlobalDaoFeePoolResp>;
  NewRecord(request: MsgNewRecord): Promise<MsgNewRecordResponse>;
  ReviewRecord(request: MsgReviewRecord): Promise<MsgReviewRecordResponse>;
  TransferRegion(request: MsgTransferRegion): Promise<MsgTransferRegionResponse>;
  IbcTransferFromRegionTreasure(
    request: MsgIbcTransferFromRegionTreasure,
  ): Promise<MsgIbcTransferFromRegionTreasureResponse>;
}

export class MsgClientImpl implements Msg {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.Stake = this.Stake.bind(this);
    this.Unstake = this.Unstake.bind(this);
    this.NewRegion = this.NewRegion.bind(this);
    this.RemoveRegion = this.RemoveRegion.bind(this);
    this.WithdrawDelegatorReward = this.WithdrawDelegatorReward.bind(this);
    this.NewFixedDepositCfg = this.NewFixedDepositCfg.bind(this);
    this.SetFixedDepositCfgStatus = this.SetFixedDepositCfgStatus.bind(this);
    this.SetFixedDepositCfgRate = this.SetFixedDepositCfgRate.bind(this);
    this.DoFixedDeposit = this.DoFixedDeposit.bind(this);
    this.WithdrawFixedDeposit = this.WithdrawFixedDeposit.bind(this);
    this.RemoveFixedDepositCfg = this.RemoveFixedDepositCfg.bind(this);
    this.WithdrawFromRegion = this.WithdrawFromRegion.bind(this);
    this.WithdrawFromGlobalDaoFeePool = this.WithdrawFromGlobalDaoFeePool.bind(this);
    this.NewRecord = this.NewRecord.bind(this);
    this.ReviewRecord = this.ReviewRecord.bind(this);
    this.TransferRegion = this.TransferRegion.bind(this);
    this.IbcTransferFromRegionTreasure = this.IbcTransferFromRegionTreasure.bind(this);
  }
  Stake(request: MsgStake): Promise<MsgStakeResponse> {
    const data = MsgStake.encode(request).finish();
    const promise = this.rpc.request("metaearth.wstaking.Msg", "Stake", data);
    return promise.then((data) => MsgStakeResponse.decode(new _m0.Reader(data)));
  }

  Unstake(request: MsgUnstake): Promise<MsgUnstakeResponse> {
    const data = MsgUnstake.encode(request).finish();
    const promise = this.rpc.request("metaearth.wstaking.Msg", "Unstake", data);
    return promise.then((data) => MsgUnstakeResponse.decode(new _m0.Reader(data)));
  }

  NewRegion(request: MsgNewRegion): Promise<MsgNewRegionResponse> {
    const data = MsgNewRegion.encode(request).finish();
    const promise = this.rpc.request("metaearth.wstaking.Msg", "NewRegion", data);
    return promise.then((data) => MsgNewRegionResponse.decode(new _m0.Reader(data)));
  }

  RemoveRegion(request: MsgRemoveRegion): Promise<MsgRemoveRegionResponse> {
    const data = MsgRemoveRegion.encode(request).finish();
    const promise = this.rpc.request("metaearth.wstaking.Msg", "RemoveRegion", data);
    return promise.then((data) => MsgRemoveRegionResponse.decode(new _m0.Reader(data)));
  }

  WithdrawDelegatorReward(request: MsgWithdrawDelegatorReward): Promise<MsgWithdrawDelegatorRewardResponse> {
    const data = MsgWithdrawDelegatorReward.encode(request).finish();
    const promise = this.rpc.request("metaearth.wstaking.Msg", "WithdrawDelegatorReward", data);
    return promise.then((data) => MsgWithdrawDelegatorRewardResponse.decode(new _m0.Reader(data)));
  }

  NewFixedDepositCfg(request: MsgNewFixedDepositCfg): Promise<MsgNewFixedDepositCfgResp> {
    const data = MsgNewFixedDepositCfg.encode(request).finish();
    const promise = this.rpc.request("metaearth.wstaking.Msg", "NewFixedDepositCfg", data);
    return promise.then((data) => MsgNewFixedDepositCfgResp.decode(new _m0.Reader(data)));
  }

  SetFixedDepositCfgStatus(request: MsgSetFixedDepositCfgStatus): Promise<MsgSetFixedDepositCfgStatusResp> {
    const data = MsgSetFixedDepositCfgStatus.encode(request).finish();
    const promise = this.rpc.request("metaearth.wstaking.Msg", "SetFixedDepositCfgStatus", data);
    return promise.then((data) => MsgSetFixedDepositCfgStatusResp.decode(new _m0.Reader(data)));
  }

  SetFixedDepositCfgRate(request: MsgSetFixedDepositCfgRate): Promise<MsgSetFixedDepositCfgRateResp> {
    const data = MsgSetFixedDepositCfgRate.encode(request).finish();
    const promise = this.rpc.request("metaearth.wstaking.Msg", "SetFixedDepositCfgRate", data);
    return promise.then((data) => MsgSetFixedDepositCfgRateResp.decode(new _m0.Reader(data)));
  }

  DoFixedDeposit(request: MsgDoFixedDeposit): Promise<MsgDoFixedDepositResponse> {
    const data = MsgDoFixedDeposit.encode(request).finish();
    const promise = this.rpc.request("metaearth.wstaking.Msg", "DoFixedDeposit", data);
    return promise.then((data) => MsgDoFixedDepositResponse.decode(new _m0.Reader(data)));
  }

  WithdrawFixedDeposit(request: MsgWithdrawFixedDeposit): Promise<MsgWithdrawFixedDepositResponse> {
    const data = MsgWithdrawFixedDeposit.encode(request).finish();
    const promise = this.rpc.request("metaearth.wstaking.Msg", "WithdrawFixedDeposit", data);
    return promise.then((data) => MsgWithdrawFixedDepositResponse.decode(new _m0.Reader(data)));
  }

  RemoveFixedDepositCfg(request: MsgRemoveFixedDepositCfg): Promise<MsgRemoveFixedDepositCfgResp> {
    const data = MsgRemoveFixedDepositCfg.encode(request).finish();
    const promise = this.rpc.request("metaearth.wstaking.Msg", "RemoveFixedDepositCfg", data);
    return promise.then((data) => MsgRemoveFixedDepositCfgResp.decode(new _m0.Reader(data)));
  }

  WithdrawFromRegion(request: MsgWithdrawFromRegion): Promise<MsgWithdrawFromRegionResp> {
    const data = MsgWithdrawFromRegion.encode(request).finish();
    const promise = this.rpc.request("metaearth.wstaking.Msg", "WithdrawFromRegion", data);
    return promise.then((data) => MsgWithdrawFromRegionResp.decode(new _m0.Reader(data)));
  }

  WithdrawFromGlobalDaoFeePool(request: MsgWithdrawFromGlobalDaoFeePool): Promise<MsgWithdrawFromGlobalDaoFeePoolResp> {
    const data = MsgWithdrawFromGlobalDaoFeePool.encode(request).finish();
    const promise = this.rpc.request("metaearth.wstaking.Msg", "WithdrawFromGlobalDaoFeePool", data);
    return promise.then((data) => MsgWithdrawFromGlobalDaoFeePoolResp.decode(new _m0.Reader(data)));
  }

  NewRecord(request: MsgNewRecord): Promise<MsgNewRecordResponse> {
    const data = MsgNewRecord.encode(request).finish();
    const promise = this.rpc.request("metaearth.wstaking.Msg", "NewRecord", data);
    return promise.then((data) => MsgNewRecordResponse.decode(new _m0.Reader(data)));
  }

  ReviewRecord(request: MsgReviewRecord): Promise<MsgReviewRecordResponse> {
    const data = MsgReviewRecord.encode(request).finish();
    const promise = this.rpc.request("metaearth.wstaking.Msg", "ReviewRecord", data);
    return promise.then((data) => MsgReviewRecordResponse.decode(new _m0.Reader(data)));
  }

  TransferRegion(request: MsgTransferRegion): Promise<MsgTransferRegionResponse> {
    const data = MsgTransferRegion.encode(request).finish();
    const promise = this.rpc.request("metaearth.wstaking.Msg", "TransferRegion", data);
    return promise.then((data) => MsgTransferRegionResponse.decode(new _m0.Reader(data)));
  }

  IbcTransferFromRegionTreasure(
    request: MsgIbcTransferFromRegionTreasure,
  ): Promise<MsgIbcTransferFromRegionTreasureResponse> {
    const data = MsgIbcTransferFromRegionTreasure.encode(request).finish();
    const promise = this.rpc.request("metaearth.wstaking.Msg", "IbcTransferFromRegionTreasure", data);
    return promise.then((data) => MsgIbcTransferFromRegionTreasureResponse.decode(new _m0.Reader(data)));
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

function toTimestamp(date: Date): Timestamp {
  const seconds = date.getTime() / 1_000;
  const nanos = (date.getTime() % 1_000) * 1_000_000;
  return { seconds, nanos };
}

function fromTimestamp(t: Timestamp): Date {
  let millis = t.seconds * 1_000;
  millis += t.nanos / 1_000_000;
  return new Date(millis);
}

function fromJsonTimestamp(o: any): Date {
  if (o instanceof Date) {
    return o;
  } else if (typeof o === "string") {
    return new Date(o);
  } else {
    return fromTimestamp(Timestamp.fromJSON(o));
  }
}

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
