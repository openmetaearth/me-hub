/* eslint-disable */
import Long from "long";
import _m0 from "protobufjs/minimal";
import { StateInfoIndex } from "./state_info";

export const protobufPackage = "dymensionxyz.dymension.rollapp";

/**
 * RollappGenesisState is a partial repr of the state the hub can expect the
 * rollapp to be in upon genesis
 */
export interface RollappGenesisState {
  /**
   * If true, then full usage of the canonical ibc transfer channel is enabled.
   * Note: in v3.1.0 and prior this field marked the completion of the 'genesis
   * event' Keeping and renaming the field enables a seamless upgrade
   * https://www.notion.so/dymension/ADR-x-Genesis-Bridge-Phase-2-89769aa551b5440b9ed403a101775ce1?pvs=4#89698384d815435b87393dbe45bc5a74
   * to the new genesis transfer protocol
   * Note: if this field is false, ibc transfers may still be allowed in one or
   * either direction.
   */
  transfersEnabled: boolean;
}

/**
 * Rollapp defines a rollapp object. First the RollApp is created and then
 * sequencers can be created and attached. The RollApp is identified by
 * rollappId
 */
export interface Rollapp {
  /**
   * The unique identifier of the rollapp chain.
   * The rollappId follows the same standard as cosmos chain_id.
   */
  rollappId: string;
  /** creator is the bech32-encoded address of the rollapp creator. */
  creator: string;
  /**
   * version is the software and configuration version.
   * starts from 1 and increases by one on every MsgUpdateState
   */
  version: number;
  /** maxSequencers is the maximum number of sequencers. */
  maxSequencers: number;
  /**
   * permissionedAddresses is a bech32-encoded address list of the sequencers
   * that are allowed to serve this rollappId. In the case of an empty list, the
   * rollapp is considered permissionless.
   */
  permissionedAddresses: string[];
  /**
   * genesis_state is a partial repr of the state the hub can expect the rollapp
   * to be in upon genesis
   */
  genesisState:
    | RollappGenesisState
    | undefined;
  /** channel_id will be set to the canonical IBC channel of the rollapp. */
  channelId: string;
  /** frozen is a boolean that indicates if the rollapp is frozen. */
  frozen: boolean;
  /** registeredDenoms is a list of registered denom bases on this rollapp */
  registeredDenoms: string[];
}

/** Rollapp summary is a compact representation of Rollapp */
export interface RollappSummary {
  /**
   * The unique identifier of the rollapp chain.
   * The rollappId follows the same standard as cosmos chain_id.
   */
  rollappId: string;
  /** Defines the index of the last rollapp UpdateState. */
  latestStateIndex:
    | StateInfoIndex
    | undefined;
  /** Defines the index of the last rollapp UpdateState that was finalized. */
  latestFinalizedStateIndex: StateInfoIndex | undefined;
}

function createBaseRollappGenesisState(): RollappGenesisState {
  return { transfersEnabled: false };
}

export const RollappGenesisState = {
  encode(message: RollappGenesisState, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.transfersEnabled === true) {
      writer.uint32(16).bool(message.transfersEnabled);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RollappGenesisState {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRollappGenesisState();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 2:
          message.transfersEnabled = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): RollappGenesisState {
    return { transfersEnabled: isSet(object.transfersEnabled) ? Boolean(object.transfersEnabled) : false };
  },

  toJSON(message: RollappGenesisState): unknown {
    const obj: any = {};
    message.transfersEnabled !== undefined && (obj.transfersEnabled = message.transfersEnabled);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<RollappGenesisState>, I>>(object: I): RollappGenesisState {
    const message = createBaseRollappGenesisState();
    message.transfersEnabled = object.transfersEnabled ?? false;
    return message;
  },
};

function createBaseRollapp(): Rollapp {
  return {
    rollappId: "",
    creator: "",
    version: 0,
    maxSequencers: 0,
    permissionedAddresses: [],
    genesisState: undefined,
    channelId: "",
    frozen: false,
    registeredDenoms: [],
  };
}

export const Rollapp = {
  encode(message: Rollapp, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.rollappId !== "") {
      writer.uint32(10).string(message.rollappId);
    }
    if (message.creator !== "") {
      writer.uint32(18).string(message.creator);
    }
    if (message.version !== 0) {
      writer.uint32(24).uint64(message.version);
    }
    if (message.maxSequencers !== 0) {
      writer.uint32(32).uint64(message.maxSequencers);
    }
    for (const v of message.permissionedAddresses) {
      writer.uint32(42).string(v!);
    }
    if (message.genesisState !== undefined) {
      RollappGenesisState.encode(message.genesisState, writer.uint32(58).fork()).ldelim();
    }
    if (message.channelId !== "") {
      writer.uint32(66).string(message.channelId);
    }
    if (message.frozen === true) {
      writer.uint32(72).bool(message.frozen);
    }
    for (const v of message.registeredDenoms) {
      writer.uint32(82).string(v!);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Rollapp {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRollapp();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.rollappId = reader.string();
          break;
        case 2:
          message.creator = reader.string();
          break;
        case 3:
          message.version = longToNumber(reader.uint64() as Long);
          break;
        case 4:
          message.maxSequencers = longToNumber(reader.uint64() as Long);
          break;
        case 5:
          message.permissionedAddresses.push(reader.string());
          break;
        case 7:
          message.genesisState = RollappGenesisState.decode(reader, reader.uint32());
          break;
        case 8:
          message.channelId = reader.string();
          break;
        case 9:
          message.frozen = reader.bool();
          break;
        case 10:
          message.registeredDenoms.push(reader.string());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): Rollapp {
    return {
      rollappId: isSet(object.rollappId) ? String(object.rollappId) : "",
      creator: isSet(object.creator) ? String(object.creator) : "",
      version: isSet(object.version) ? Number(object.version) : 0,
      maxSequencers: isSet(object.maxSequencers) ? Number(object.maxSequencers) : 0,
      permissionedAddresses: Array.isArray(object?.permissionedAddresses)
        ? object.permissionedAddresses.map((e: any) => String(e))
        : [],
      genesisState: isSet(object.genesisState) ? RollappGenesisState.fromJSON(object.genesisState) : undefined,
      channelId: isSet(object.channelId) ? String(object.channelId) : "",
      frozen: isSet(object.frozen) ? Boolean(object.frozen) : false,
      registeredDenoms: Array.isArray(object?.registeredDenoms)
        ? object.registeredDenoms.map((e: any) => String(e))
        : [],
    };
  },

  toJSON(message: Rollapp): unknown {
    const obj: any = {};
    message.rollappId !== undefined && (obj.rollappId = message.rollappId);
    message.creator !== undefined && (obj.creator = message.creator);
    message.version !== undefined && (obj.version = Math.round(message.version));
    message.maxSequencers !== undefined && (obj.maxSequencers = Math.round(message.maxSequencers));
    if (message.permissionedAddresses) {
      obj.permissionedAddresses = message.permissionedAddresses.map((e) => e);
    } else {
      obj.permissionedAddresses = [];
    }
    message.genesisState !== undefined
      && (obj.genesisState = message.genesisState ? RollappGenesisState.toJSON(message.genesisState) : undefined);
    message.channelId !== undefined && (obj.channelId = message.channelId);
    message.frozen !== undefined && (obj.frozen = message.frozen);
    if (message.registeredDenoms) {
      obj.registeredDenoms = message.registeredDenoms.map((e) => e);
    } else {
      obj.registeredDenoms = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<Rollapp>, I>>(object: I): Rollapp {
    const message = createBaseRollapp();
    message.rollappId = object.rollappId ?? "";
    message.creator = object.creator ?? "";
    message.version = object.version ?? 0;
    message.maxSequencers = object.maxSequencers ?? 0;
    message.permissionedAddresses = object.permissionedAddresses?.map((e) => e) || [];
    message.genesisState = (object.genesisState !== undefined && object.genesisState !== null)
      ? RollappGenesisState.fromPartial(object.genesisState)
      : undefined;
    message.channelId = object.channelId ?? "";
    message.frozen = object.frozen ?? false;
    message.registeredDenoms = object.registeredDenoms?.map((e) => e) || [];
    return message;
  },
};

function createBaseRollappSummary(): RollappSummary {
  return { rollappId: "", latestStateIndex: undefined, latestFinalizedStateIndex: undefined };
}

export const RollappSummary = {
  encode(message: RollappSummary, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.rollappId !== "") {
      writer.uint32(10).string(message.rollappId);
    }
    if (message.latestStateIndex !== undefined) {
      StateInfoIndex.encode(message.latestStateIndex, writer.uint32(18).fork()).ldelim();
    }
    if (message.latestFinalizedStateIndex !== undefined) {
      StateInfoIndex.encode(message.latestFinalizedStateIndex, writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RollappSummary {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRollappSummary();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.rollappId = reader.string();
          break;
        case 2:
          message.latestStateIndex = StateInfoIndex.decode(reader, reader.uint32());
          break;
        case 3:
          message.latestFinalizedStateIndex = StateInfoIndex.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): RollappSummary {
    return {
      rollappId: isSet(object.rollappId) ? String(object.rollappId) : "",
      latestStateIndex: isSet(object.latestStateIndex) ? StateInfoIndex.fromJSON(object.latestStateIndex) : undefined,
      latestFinalizedStateIndex: isSet(object.latestFinalizedStateIndex)
        ? StateInfoIndex.fromJSON(object.latestFinalizedStateIndex)
        : undefined,
    };
  },

  toJSON(message: RollappSummary): unknown {
    const obj: any = {};
    message.rollappId !== undefined && (obj.rollappId = message.rollappId);
    message.latestStateIndex !== undefined
      && (obj.latestStateIndex = message.latestStateIndex
        ? StateInfoIndex.toJSON(message.latestStateIndex)
        : undefined);
    message.latestFinalizedStateIndex !== undefined
      && (obj.latestFinalizedStateIndex = message.latestFinalizedStateIndex
        ? StateInfoIndex.toJSON(message.latestFinalizedStateIndex)
        : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<RollappSummary>, I>>(object: I): RollappSummary {
    const message = createBaseRollappSummary();
    message.rollappId = object.rollappId ?? "";
    message.latestStateIndex = (object.latestStateIndex !== undefined && object.latestStateIndex !== null)
      ? StateInfoIndex.fromPartial(object.latestStateIndex)
      : undefined;
    message.latestFinalizedStateIndex =
      (object.latestFinalizedStateIndex !== undefined && object.latestFinalizedStateIndex !== null)
        ? StateInfoIndex.fromPartial(object.latestFinalizedStateIndex)
        : undefined;
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
