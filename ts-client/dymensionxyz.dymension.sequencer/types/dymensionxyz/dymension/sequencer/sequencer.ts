/* eslint-disable */
import Long from "long";
import _m0 from "protobufjs/minimal";
import { Coin } from "../../../cosmos/base/v1beta1/coin";
import { Any } from "../../../google/protobuf/any";
import { Timestamp } from "../../../google/protobuf/timestamp";
import { Description } from "./description";
import { OperatingStatus, operatingStatusFromJSON, operatingStatusToJSON } from "./operating_status";

export const protobufPackage = "dymensionxyz.dymension.sequencer";

/**
 * Sequencer defines a sequencer identified by its' address (sequencerAddress).
 * The sequencer could be attached to only one rollapp (rollappId).
 */
export interface Sequencer {
  /**
   * sequencerAddress is the bech32-encoded address of the sequencer account
   * which is the account that the message was sent from.
   */
  sequencerAddress: string;
  /**
   * pubkey is the public key of the sequencers' dymint client, as a Protobuf
   * Any.
   */
  dymintPubKey:
    | Any
    | undefined;
  /** rollappId defines the rollapp to which the sequencer belongs. */
  rollappId: string;
  /** description defines the descriptive terms for the sequencer. */
  description:
    | Description
    | undefined;
  /**
   * jailed defined whether the sequencer has been jailed from bonded status or
   * not.
   */
  jailed: boolean;
  /** proposer defines whether the sequencer is a proposer or not. */
  proposer: boolean;
  /** status is the sequencer status (bonded/unbonding/unbonded). */
  status: OperatingStatus;
  /** tokens define the delegated tokens (incl. self-delegation). */
  tokens: Coin[];
  /**
   * unbonding_height defines, if unbonding, the height at which this sequencer
   * has begun unbonding.
   */
  unbondingHeight: number;
  /**
   * unbond_time defines, if unbonding, the min time for the sequencer to
   * complete unbonding.
   */
  unbondTime: Date | undefined;
}

function createBaseSequencer(): Sequencer {
  return {
    sequencerAddress: "",
    dymintPubKey: undefined,
    rollappId: "",
    description: undefined,
    jailed: false,
    proposer: false,
    status: 0,
    tokens: [],
    unbondingHeight: 0,
    unbondTime: undefined,
  };
}

export const Sequencer = {
  encode(message: Sequencer, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.sequencerAddress !== "") {
      writer.uint32(10).string(message.sequencerAddress);
    }
    if (message.dymintPubKey !== undefined) {
      Any.encode(message.dymintPubKey, writer.uint32(18).fork()).ldelim();
    }
    if (message.rollappId !== "") {
      writer.uint32(26).string(message.rollappId);
    }
    if (message.description !== undefined) {
      Description.encode(message.description, writer.uint32(34).fork()).ldelim();
    }
    if (message.jailed === true) {
      writer.uint32(40).bool(message.jailed);
    }
    if (message.proposer === true) {
      writer.uint32(48).bool(message.proposer);
    }
    if (message.status !== 0) {
      writer.uint32(56).int32(message.status);
    }
    for (const v of message.tokens) {
      Coin.encode(v!, writer.uint32(66).fork()).ldelim();
    }
    if (message.unbondingHeight !== 0) {
      writer.uint32(72).int64(message.unbondingHeight);
    }
    if (message.unbondTime !== undefined) {
      Timestamp.encode(toTimestamp(message.unbondTime), writer.uint32(82).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Sequencer {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSequencer();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.sequencerAddress = reader.string();
          break;
        case 2:
          message.dymintPubKey = Any.decode(reader, reader.uint32());
          break;
        case 3:
          message.rollappId = reader.string();
          break;
        case 4:
          message.description = Description.decode(reader, reader.uint32());
          break;
        case 5:
          message.jailed = reader.bool();
          break;
        case 6:
          message.proposer = reader.bool();
          break;
        case 7:
          message.status = reader.int32() as any;
          break;
        case 8:
          message.tokens.push(Coin.decode(reader, reader.uint32()));
          break;
        case 9:
          message.unbondingHeight = longToNumber(reader.int64() as Long);
          break;
        case 10:
          message.unbondTime = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): Sequencer {
    return {
      sequencerAddress: isSet(object.sequencerAddress) ? String(object.sequencerAddress) : "",
      dymintPubKey: isSet(object.dymintPubKey) ? Any.fromJSON(object.dymintPubKey) : undefined,
      rollappId: isSet(object.rollappId) ? String(object.rollappId) : "",
      description: isSet(object.description) ? Description.fromJSON(object.description) : undefined,
      jailed: isSet(object.jailed) ? Boolean(object.jailed) : false,
      proposer: isSet(object.proposer) ? Boolean(object.proposer) : false,
      status: isSet(object.status) ? operatingStatusFromJSON(object.status) : 0,
      tokens: Array.isArray(object?.tokens) ? object.tokens.map((e: any) => Coin.fromJSON(e)) : [],
      unbondingHeight: isSet(object.unbondingHeight) ? Number(object.unbondingHeight) : 0,
      unbondTime: isSet(object.unbondTime) ? fromJsonTimestamp(object.unbondTime) : undefined,
    };
  },

  toJSON(message: Sequencer): unknown {
    const obj: any = {};
    message.sequencerAddress !== undefined && (obj.sequencerAddress = message.sequencerAddress);
    message.dymintPubKey !== undefined
      && (obj.dymintPubKey = message.dymintPubKey ? Any.toJSON(message.dymintPubKey) : undefined);
    message.rollappId !== undefined && (obj.rollappId = message.rollappId);
    message.description !== undefined
      && (obj.description = message.description ? Description.toJSON(message.description) : undefined);
    message.jailed !== undefined && (obj.jailed = message.jailed);
    message.proposer !== undefined && (obj.proposer = message.proposer);
    message.status !== undefined && (obj.status = operatingStatusToJSON(message.status));
    if (message.tokens) {
      obj.tokens = message.tokens.map((e) => e ? Coin.toJSON(e) : undefined);
    } else {
      obj.tokens = [];
    }
    message.unbondingHeight !== undefined && (obj.unbondingHeight = Math.round(message.unbondingHeight));
    message.unbondTime !== undefined && (obj.unbondTime = message.unbondTime.toISOString());
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<Sequencer>, I>>(object: I): Sequencer {
    const message = createBaseSequencer();
    message.sequencerAddress = object.sequencerAddress ?? "";
    message.dymintPubKey = (object.dymintPubKey !== undefined && object.dymintPubKey !== null)
      ? Any.fromPartial(object.dymintPubKey)
      : undefined;
    message.rollappId = object.rollappId ?? "";
    message.description = (object.description !== undefined && object.description !== null)
      ? Description.fromPartial(object.description)
      : undefined;
    message.jailed = object.jailed ?? false;
    message.proposer = object.proposer ?? false;
    message.status = object.status ?? 0;
    message.tokens = object.tokens?.map((e) => Coin.fromPartial(e)) || [];
    message.unbondingHeight = object.unbondingHeight ?? 0;
    message.unbondTime = object.unbondTime ?? undefined;
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
