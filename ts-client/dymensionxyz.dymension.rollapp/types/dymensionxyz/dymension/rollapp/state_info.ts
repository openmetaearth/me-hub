/* eslint-disable */
import Long from "long";
import _m0 from "protobufjs/minimal";
import { Status, statusFromJSON, statusToJSON } from "../common/status";
import { BlockDescriptors } from "./block_descriptor";

export const protobufPackage = "dymensionxyz.dymension.rollapp";

/**
 * StateInfoIndex is the data used for indexing and retrieving a StateInfo
 * it updated and saved with every UpdateState in StateInfo.
 * We use the this structure also for:
 * 1. LatestStateInfoIndex which defines the rollapps' current (latest) index of
 * the last UpdateState
 * 2. LatestFinalizedStateIndex which defines the rollapps' current (latest)
 * index of the latest StateInfo that was finalized
 */
export interface StateInfoIndex {
  /**
   * rollappId is the rollapp that the sequencer belongs to and asking to update
   * it used to identify the what rollapp a StateInfo belongs
   * The rollappId follows the same standard as cosmos chain_id
   */
  rollappId: string;
  /**
   * index is a sequential increasing number, updating on each
   * state update used for indexing to a specific state info, the first index is
   * 1
   */
  index: number;
}

/** StateInfo defines a rollapps' state. */
export interface StateInfo {
  /**
   * stateInfoIndex defines what rollapp the state belongs to
   * and in which index it can be referenced
   */
  stateInfoIndex:
    | StateInfoIndex
    | undefined;
  /** sequencer is the bech32-encoded address of the sequencer sent the update */
  sequencer: string;
  /** startHeight is the block height of the first block in the batch */
  startHeight: number;
  /** numBlocks is the number of blocks included in this batch update */
  numBlocks: number;
  /** DAPath is the description of the location on the DA layer */
  DAPath: string;
  /** version is the version of the rollapp */
  version: number;
  /** creationHeight is the height at which the UpdateState took place */
  creationHeight: number;
  /** status is the status of the state update */
  status: Status;
  /**
   * BDs is a list of block description objects (one per block)
   * the list must be ordered by height, starting from startHeight to
   * startHeight+numBlocks-1
   */
  BDs: BlockDescriptors | undefined;
}

/** StateInfoSummary is a compact representation of StateInfo */
export interface StateInfoSummary {
  /**
   * stateInfoIndex defines what rollapp the state belongs to
   * and in which index it can be referenced
   */
  stateInfoIndex:
    | StateInfoIndex
    | undefined;
  /** status is the status of the state update */
  status: Status;
  /** creationHeight is the height at which the UpdateState took place */
  creationHeight: number;
}

/**
 * BlockHeightToFinalizationQueue defines a map from block height to list of
 * states to finalized
 */
export interface BlockHeightToFinalizationQueue {
  /** creationHeight is the block height that the state should be finalized */
  creationHeight: number;
  /**
   * finalizationQueue is a list of states that are waiting to be finalized
   * when the block height becomes creationHeight
   */
  finalizationQueue: StateInfoIndex[];
}

function createBaseStateInfoIndex(): StateInfoIndex {
  return { rollappId: "", index: 0 };
}

export const StateInfoIndex = {
  encode(message: StateInfoIndex, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.rollappId !== "") {
      writer.uint32(10).string(message.rollappId);
    }
    if (message.index !== 0) {
      writer.uint32(16).uint64(message.index);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): StateInfoIndex {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseStateInfoIndex();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.rollappId = reader.string();
          break;
        case 2:
          message.index = longToNumber(reader.uint64() as Long);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): StateInfoIndex {
    return {
      rollappId: isSet(object.rollappId) ? String(object.rollappId) : "",
      index: isSet(object.index) ? Number(object.index) : 0,
    };
  },

  toJSON(message: StateInfoIndex): unknown {
    const obj: any = {};
    message.rollappId !== undefined && (obj.rollappId = message.rollappId);
    message.index !== undefined && (obj.index = Math.round(message.index));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<StateInfoIndex>, I>>(object: I): StateInfoIndex {
    const message = createBaseStateInfoIndex();
    message.rollappId = object.rollappId ?? "";
    message.index = object.index ?? 0;
    return message;
  },
};

function createBaseStateInfo(): StateInfo {
  return {
    stateInfoIndex: undefined,
    sequencer: "",
    startHeight: 0,
    numBlocks: 0,
    DAPath: "",
    version: 0,
    creationHeight: 0,
    status: 0,
    BDs: undefined,
  };
}

export const StateInfo = {
  encode(message: StateInfo, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.stateInfoIndex !== undefined) {
      StateInfoIndex.encode(message.stateInfoIndex, writer.uint32(10).fork()).ldelim();
    }
    if (message.sequencer !== "") {
      writer.uint32(18).string(message.sequencer);
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
    if (message.creationHeight !== 0) {
      writer.uint32(56).uint64(message.creationHeight);
    }
    if (message.status !== 0) {
      writer.uint32(64).int32(message.status);
    }
    if (message.BDs !== undefined) {
      BlockDescriptors.encode(message.BDs, writer.uint32(74).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): StateInfo {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseStateInfo();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.stateInfoIndex = StateInfoIndex.decode(reader, reader.uint32());
          break;
        case 2:
          message.sequencer = reader.string();
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
          message.creationHeight = longToNumber(reader.uint64() as Long);
          break;
        case 8:
          message.status = reader.int32() as any;
          break;
        case 9:
          message.BDs = BlockDescriptors.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): StateInfo {
    return {
      stateInfoIndex: isSet(object.stateInfoIndex) ? StateInfoIndex.fromJSON(object.stateInfoIndex) : undefined,
      sequencer: isSet(object.sequencer) ? String(object.sequencer) : "",
      startHeight: isSet(object.startHeight) ? Number(object.startHeight) : 0,
      numBlocks: isSet(object.numBlocks) ? Number(object.numBlocks) : 0,
      DAPath: isSet(object.DAPath) ? String(object.DAPath) : "",
      version: isSet(object.version) ? Number(object.version) : 0,
      creationHeight: isSet(object.creationHeight) ? Number(object.creationHeight) : 0,
      status: isSet(object.status) ? statusFromJSON(object.status) : 0,
      BDs: isSet(object.BDs) ? BlockDescriptors.fromJSON(object.BDs) : undefined,
    };
  },

  toJSON(message: StateInfo): unknown {
    const obj: any = {};
    message.stateInfoIndex !== undefined
      && (obj.stateInfoIndex = message.stateInfoIndex ? StateInfoIndex.toJSON(message.stateInfoIndex) : undefined);
    message.sequencer !== undefined && (obj.sequencer = message.sequencer);
    message.startHeight !== undefined && (obj.startHeight = Math.round(message.startHeight));
    message.numBlocks !== undefined && (obj.numBlocks = Math.round(message.numBlocks));
    message.DAPath !== undefined && (obj.DAPath = message.DAPath);
    message.version !== undefined && (obj.version = Math.round(message.version));
    message.creationHeight !== undefined && (obj.creationHeight = Math.round(message.creationHeight));
    message.status !== undefined && (obj.status = statusToJSON(message.status));
    message.BDs !== undefined && (obj.BDs = message.BDs ? BlockDescriptors.toJSON(message.BDs) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<StateInfo>, I>>(object: I): StateInfo {
    const message = createBaseStateInfo();
    message.stateInfoIndex = (object.stateInfoIndex !== undefined && object.stateInfoIndex !== null)
      ? StateInfoIndex.fromPartial(object.stateInfoIndex)
      : undefined;
    message.sequencer = object.sequencer ?? "";
    message.startHeight = object.startHeight ?? 0;
    message.numBlocks = object.numBlocks ?? 0;
    message.DAPath = object.DAPath ?? "";
    message.version = object.version ?? 0;
    message.creationHeight = object.creationHeight ?? 0;
    message.status = object.status ?? 0;
    message.BDs = (object.BDs !== undefined && object.BDs !== null)
      ? BlockDescriptors.fromPartial(object.BDs)
      : undefined;
    return message;
  },
};

function createBaseStateInfoSummary(): StateInfoSummary {
  return { stateInfoIndex: undefined, status: 0, creationHeight: 0 };
}

export const StateInfoSummary = {
  encode(message: StateInfoSummary, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.stateInfoIndex !== undefined) {
      StateInfoIndex.encode(message.stateInfoIndex, writer.uint32(10).fork()).ldelim();
    }
    if (message.status !== 0) {
      writer.uint32(16).int32(message.status);
    }
    if (message.creationHeight !== 0) {
      writer.uint32(24).uint64(message.creationHeight);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): StateInfoSummary {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseStateInfoSummary();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.stateInfoIndex = StateInfoIndex.decode(reader, reader.uint32());
          break;
        case 2:
          message.status = reader.int32() as any;
          break;
        case 3:
          message.creationHeight = longToNumber(reader.uint64() as Long);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): StateInfoSummary {
    return {
      stateInfoIndex: isSet(object.stateInfoIndex) ? StateInfoIndex.fromJSON(object.stateInfoIndex) : undefined,
      status: isSet(object.status) ? statusFromJSON(object.status) : 0,
      creationHeight: isSet(object.creationHeight) ? Number(object.creationHeight) : 0,
    };
  },

  toJSON(message: StateInfoSummary): unknown {
    const obj: any = {};
    message.stateInfoIndex !== undefined
      && (obj.stateInfoIndex = message.stateInfoIndex ? StateInfoIndex.toJSON(message.stateInfoIndex) : undefined);
    message.status !== undefined && (obj.status = statusToJSON(message.status));
    message.creationHeight !== undefined && (obj.creationHeight = Math.round(message.creationHeight));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<StateInfoSummary>, I>>(object: I): StateInfoSummary {
    const message = createBaseStateInfoSummary();
    message.stateInfoIndex = (object.stateInfoIndex !== undefined && object.stateInfoIndex !== null)
      ? StateInfoIndex.fromPartial(object.stateInfoIndex)
      : undefined;
    message.status = object.status ?? 0;
    message.creationHeight = object.creationHeight ?? 0;
    return message;
  },
};

function createBaseBlockHeightToFinalizationQueue(): BlockHeightToFinalizationQueue {
  return { creationHeight: 0, finalizationQueue: [] };
}

export const BlockHeightToFinalizationQueue = {
  encode(message: BlockHeightToFinalizationQueue, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.creationHeight !== 0) {
      writer.uint32(8).uint64(message.creationHeight);
    }
    for (const v of message.finalizationQueue) {
      StateInfoIndex.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): BlockHeightToFinalizationQueue {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseBlockHeightToFinalizationQueue();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creationHeight = longToNumber(reader.uint64() as Long);
          break;
        case 2:
          message.finalizationQueue.push(StateInfoIndex.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): BlockHeightToFinalizationQueue {
    return {
      creationHeight: isSet(object.creationHeight) ? Number(object.creationHeight) : 0,
      finalizationQueue: Array.isArray(object?.finalizationQueue)
        ? object.finalizationQueue.map((e: any) => StateInfoIndex.fromJSON(e))
        : [],
    };
  },

  toJSON(message: BlockHeightToFinalizationQueue): unknown {
    const obj: any = {};
    message.creationHeight !== undefined && (obj.creationHeight = Math.round(message.creationHeight));
    if (message.finalizationQueue) {
      obj.finalizationQueue = message.finalizationQueue.map((e) => e ? StateInfoIndex.toJSON(e) : undefined);
    } else {
      obj.finalizationQueue = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<BlockHeightToFinalizationQueue>, I>>(
    object: I,
  ): BlockHeightToFinalizationQueue {
    const message = createBaseBlockHeightToFinalizationQueue();
    message.creationHeight = object.creationHeight ?? 0;
    message.finalizationQueue = object.finalizationQueue?.map((e) => StateInfoIndex.fromPartial(e)) || [];
    return message;
  },
};

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
