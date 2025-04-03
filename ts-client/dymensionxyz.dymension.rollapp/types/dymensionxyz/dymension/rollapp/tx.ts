/* eslint-disable */
import Long from "long";
import _m0 from "protobufjs/minimal";
import { BlockDescriptors } from "./block_descriptor";

export const protobufPackage = "dymensionxyz.dymension.rollapp";

/** MsgCreateRollapp creates a new rollapp chain on the hub. */
export interface MsgCreateRollapp {
  /** creator is the bech32-encoded address of the rollapp creator */
  creator: string;
  /**
   * rollappId is the unique identifier of the rollapp chain.
   * The rollappId follows the same standard as cosmos chain_id
   */
  rollappId: string;
  /** maxSequencers is the maximum number of sequencers */
  maxSequencers: number;
  /**
   * permissionedAddresses is a bech32-encoded address list of the
   * sequencers that are allowed to serve this rollappId.
   * In the case of an empty list, the rollapp is considered permissionless
   */
  permissionedAddresses: string[];
}

export interface MsgCreateRollappResponse {
}

/**
 * MsgUpdateState updates a rollapp state with a block batch.
 * a block batch is a list of ordered blocks (by height)
 */
export interface MsgUpdateState {
  /** creator is the bech32-encoded address of the sequencer sending the update */
  creator: string;
  /**
   * rollappId is the rollapp that the sequencer belongs to and asking to update
   * The rollappId follows the same standard as cosmos chain_id
   */
  rollappId: string;
  /** startHeight is the block height of the first block in the batch */
  startHeight: number;
  /** numBlocks is the number of blocks included in this batch update */
  numBlocks: number;
  /** DAPath is the description of the location on the DA layer */
  DAPath: string;
  /** version is the version of the rollapp */
  version: number;
  /**
   * BDs is a list of block description objects (one per block)
   * the list must be ordered by height, starting from startHeight to
   * startHeight+numBlocks-1
   */
  BDs: BlockDescriptors | undefined;
}

export interface MsgUpdateStateResponse {
}

export interface MsgSkipDelayRollapp {
  creator: string;
  rollappId: string;
  skip: boolean;
}

export interface MsgSkipDelayRollappResponse {
}

function createBaseMsgCreateRollapp(): MsgCreateRollapp {
  return { creator: "", rollappId: "", maxSequencers: 0, permissionedAddresses: [] };
}

export const MsgCreateRollapp = {
  encode(message: MsgCreateRollapp, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.rollappId !== "") {
      writer.uint32(18).string(message.rollappId);
    }
    if (message.maxSequencers !== 0) {
      writer.uint32(48).uint64(message.maxSequencers);
    }
    for (const v of message.permissionedAddresses) {
      writer.uint32(58).string(v!);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgCreateRollapp {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgCreateRollapp();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.rollappId = reader.string();
          break;
        case 6:
          message.maxSequencers = longToNumber(reader.uint64() as Long);
          break;
        case 7:
          message.permissionedAddresses.push(reader.string());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgCreateRollapp {
    return {
      creator: isSet(object.creator) ? String(object.creator) : "",
      rollappId: isSet(object.rollappId) ? String(object.rollappId) : "",
      maxSequencers: isSet(object.maxSequencers) ? Number(object.maxSequencers) : 0,
      permissionedAddresses: Array.isArray(object?.permissionedAddresses)
        ? object.permissionedAddresses.map((e: any) => String(e))
        : [],
    };
  },

  toJSON(message: MsgCreateRollapp): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    message.rollappId !== undefined && (obj.rollappId = message.rollappId);
    message.maxSequencers !== undefined && (obj.maxSequencers = Math.round(message.maxSequencers));
    if (message.permissionedAddresses) {
      obj.permissionedAddresses = message.permissionedAddresses.map((e) => e);
    } else {
      obj.permissionedAddresses = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgCreateRollapp>, I>>(object: I): MsgCreateRollapp {
    const message = createBaseMsgCreateRollapp();
    message.creator = object.creator ?? "";
    message.rollappId = object.rollappId ?? "";
    message.maxSequencers = object.maxSequencers ?? 0;
    message.permissionedAddresses = object.permissionedAddresses?.map((e) => e) || [];
    return message;
  },
};

function createBaseMsgCreateRollappResponse(): MsgCreateRollappResponse {
  return {};
}

export const MsgCreateRollappResponse = {
  encode(_: MsgCreateRollappResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgCreateRollappResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgCreateRollappResponse();
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

  fromJSON(_: any): MsgCreateRollappResponse {
    return {};
  },

  toJSON(_: MsgCreateRollappResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgCreateRollappResponse>, I>>(_: I): MsgCreateRollappResponse {
    const message = createBaseMsgCreateRollappResponse();
    return message;
  },
};

function createBaseMsgUpdateState(): MsgUpdateState {
  return { creator: "", rollappId: "", startHeight: 0, numBlocks: 0, DAPath: "", version: 0, BDs: undefined };
}

export const MsgUpdateState = {
  encode(message: MsgUpdateState, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.rollappId !== "") {
      writer.uint32(18).string(message.rollappId);
    }
    if (message.startHeight !== 0) {
      writer.uint32(24).uint64(message.startHeight);
    }
    if (message.numBlocks !== 0) {
      writer.uint32(32).uint64(message.numBlocks);
    }
    if (message.DAPath !== "") {
      writer.uint32(42).string(message.DAPath);
    }
    if (message.version !== 0) {
      writer.uint32(48).uint64(message.version);
    }
    if (message.BDs !== undefined) {
      BlockDescriptors.encode(message.BDs, writer.uint32(58).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateState {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateState();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.rollappId = reader.string();
          break;
        case 3:
          message.startHeight = longToNumber(reader.uint64() as Long);
          break;
        case 4:
          message.numBlocks = longToNumber(reader.uint64() as Long);
          break;
        case 5:
          message.DAPath = reader.string();
          break;
        case 6:
          message.version = longToNumber(reader.uint64() as Long);
          break;
        case 7:
          message.BDs = BlockDescriptors.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgUpdateState {
    return {
      creator: isSet(object.creator) ? String(object.creator) : "",
      rollappId: isSet(object.rollappId) ? String(object.rollappId) : "",
      startHeight: isSet(object.startHeight) ? Number(object.startHeight) : 0,
      numBlocks: isSet(object.numBlocks) ? Number(object.numBlocks) : 0,
      DAPath: isSet(object.DAPath) ? String(object.DAPath) : "",
      version: isSet(object.version) ? Number(object.version) : 0,
      BDs: isSet(object.BDs) ? BlockDescriptors.fromJSON(object.BDs) : undefined,
    };
  },

  toJSON(message: MsgUpdateState): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    message.rollappId !== undefined && (obj.rollappId = message.rollappId);
    message.startHeight !== undefined && (obj.startHeight = Math.round(message.startHeight));
    message.numBlocks !== undefined && (obj.numBlocks = Math.round(message.numBlocks));
    message.DAPath !== undefined && (obj.DAPath = message.DAPath);
    message.version !== undefined && (obj.version = Math.round(message.version));
    message.BDs !== undefined && (obj.BDs = message.BDs ? BlockDescriptors.toJSON(message.BDs) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgUpdateState>, I>>(object: I): MsgUpdateState {
    const message = createBaseMsgUpdateState();
    message.creator = object.creator ?? "";
    message.rollappId = object.rollappId ?? "";
    message.startHeight = object.startHeight ?? 0;
    message.numBlocks = object.numBlocks ?? 0;
    message.DAPath = object.DAPath ?? "";
    message.version = object.version ?? 0;
    message.BDs = (object.BDs !== undefined && object.BDs !== null)
      ? BlockDescriptors.fromPartial(object.BDs)
      : undefined;
    return message;
  },
};

function createBaseMsgUpdateStateResponse(): MsgUpdateStateResponse {
  return {};
}

export const MsgUpdateStateResponse = {
  encode(_: MsgUpdateStateResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateStateResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateStateResponse();
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

  fromJSON(_: any): MsgUpdateStateResponse {
    return {};
  },

  toJSON(_: MsgUpdateStateResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgUpdateStateResponse>, I>>(_: I): MsgUpdateStateResponse {
    const message = createBaseMsgUpdateStateResponse();
    return message;
  },
};

function createBaseMsgSkipDelayRollapp(): MsgSkipDelayRollapp {
  return { creator: "", rollappId: "", skip: false };
}

export const MsgSkipDelayRollapp = {
  encode(message: MsgSkipDelayRollapp, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.rollappId !== "") {
      writer.uint32(18).string(message.rollappId);
    }
    if (message.skip === true) {
      writer.uint32(24).bool(message.skip);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgSkipDelayRollapp {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgSkipDelayRollapp();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.rollappId = reader.string();
          break;
        case 3:
          message.skip = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgSkipDelayRollapp {
    return {
      creator: isSet(object.creator) ? String(object.creator) : "",
      rollappId: isSet(object.rollappId) ? String(object.rollappId) : "",
      skip: isSet(object.skip) ? Boolean(object.skip) : false,
    };
  },

  toJSON(message: MsgSkipDelayRollapp): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    message.rollappId !== undefined && (obj.rollappId = message.rollappId);
    message.skip !== undefined && (obj.skip = message.skip);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgSkipDelayRollapp>, I>>(object: I): MsgSkipDelayRollapp {
    const message = createBaseMsgSkipDelayRollapp();
    message.creator = object.creator ?? "";
    message.rollappId = object.rollappId ?? "";
    message.skip = object.skip ?? false;
    return message;
  },
};

function createBaseMsgSkipDelayRollappResponse(): MsgSkipDelayRollappResponse {
  return {};
}

export const MsgSkipDelayRollappResponse = {
  encode(_: MsgSkipDelayRollappResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgSkipDelayRollappResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgSkipDelayRollappResponse();
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

  fromJSON(_: any): MsgSkipDelayRollappResponse {
    return {};
  },

  toJSON(_: MsgSkipDelayRollappResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgSkipDelayRollappResponse>, I>>(_: I): MsgSkipDelayRollappResponse {
    const message = createBaseMsgSkipDelayRollappResponse();
    return message;
  },
};

/** Msg defines the Msg service. */
export interface Msg {
  CreateRollapp(request: MsgCreateRollapp): Promise<MsgCreateRollappResponse>;
  SkipDelayRollapp(request: MsgSkipDelayRollapp): Promise<MsgSkipDelayRollappResponse>;
  UpdateState(request: MsgUpdateState): Promise<MsgUpdateStateResponse>;
}

export class MsgClientImpl implements Msg {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.CreateRollapp = this.CreateRollapp.bind(this);
    this.SkipDelayRollapp = this.SkipDelayRollapp.bind(this);
    this.UpdateState = this.UpdateState.bind(this);
  }
  CreateRollapp(request: MsgCreateRollapp): Promise<MsgCreateRollappResponse> {
    const data = MsgCreateRollapp.encode(request).finish();
    const promise = this.rpc.request("dymensionxyz.dymension.rollapp.Msg", "CreateRollapp", data);
    return promise.then((data) => MsgCreateRollappResponse.decode(new _m0.Reader(data)));
  }

  SkipDelayRollapp(request: MsgSkipDelayRollapp): Promise<MsgSkipDelayRollappResponse> {
    const data = MsgSkipDelayRollapp.encode(request).finish();
    const promise = this.rpc.request("dymensionxyz.dymension.rollapp.Msg", "SkipDelayRollapp", data);
    return promise.then((data) => MsgSkipDelayRollappResponse.decode(new _m0.Reader(data)));
  }

  UpdateState(request: MsgUpdateState): Promise<MsgUpdateStateResponse> {
    const data = MsgUpdateState.encode(request).finish();
    const promise = this.rpc.request("dymensionxyz.dymension.rollapp.Msg", "UpdateState", data);
    return promise.then((data) => MsgUpdateStateResponse.decode(new _m0.Reader(data)));
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
