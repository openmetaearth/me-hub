/* eslint-disable */
import Long from "long";
import _m0 from "protobufjs/minimal";
import { Coin } from "../../cosmos/base/v1beta1/coin";
import { Timestamp } from "../../google/protobuf/timestamp";
import { DistrInfo } from "./distr_info";

export const protobufPackage = "metaearth.streamer";

/**
 * Stream is an object that stores and distributes yields to recipients who
 * satisfy certain conditions. Currently streams support conditions around the
 * duration for which a given denom is locked.
 */
export interface Stream {
  /** id is the unique ID of a Stream */
  id: number;
  /** distribute_to is the distr_info. */
  distributeTo:
    | DistrInfo
    | undefined;
  /**
   * coins is the total amount of coins that have been in the stream
   * Can distribute multiple coin denoms
   */
  coins: Coin[];
  /** start_time is the distribution start time */
  startTime:
    | Date
    | undefined;
  /**
   * distr_epoch_identifier is what epoch type di-stribution will be triggered
   * by
   * (day, week, etc.)
   */
  distrEpochIdentifier: string;
  /**
   * num_epochs_paid_over is the number of total epochs distribution will be
   * completed over
   */
  numEpochsPaidOver: number;
  /**
   * filled_epochs is the number of epochs distribution has been completed on
   * already
   */
  filledEpochs: number;
  /** distributed_coins are coins that have been distributed already */
  distributedCoins: Coin[];
}

function createBaseStream(): Stream {
  return {
    id: 0,
    distributeTo: undefined,
    coins: [],
    startTime: undefined,
    distrEpochIdentifier: "",
    numEpochsPaidOver: 0,
    filledEpochs: 0,
    distributedCoins: [],
  };
}

export const Stream = {
  encode(message: Stream, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== 0) {
      writer.uint32(8).uint64(message.id);
    }
    if (message.distributeTo !== undefined) {
      DistrInfo.encode(message.distributeTo, writer.uint32(18).fork()).ldelim();
    }
    for (const v of message.coins) {
      Coin.encode(v!, writer.uint32(26).fork()).ldelim();
    }
    if (message.startTime !== undefined) {
      Timestamp.encode(toTimestamp(message.startTime), writer.uint32(34).fork()).ldelim();
    }
    if (message.distrEpochIdentifier !== "") {
      writer.uint32(42).string(message.distrEpochIdentifier);
    }
    if (message.numEpochsPaidOver !== 0) {
      writer.uint32(48).uint64(message.numEpochsPaidOver);
    }
    if (message.filledEpochs !== 0) {
      writer.uint32(56).uint64(message.filledEpochs);
    }
    for (const v of message.distributedCoins) {
      Coin.encode(v!, writer.uint32(66).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Stream {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseStream();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.id = longToNumber(reader.uint64() as Long);
          break;
        case 2:
          message.distributeTo = DistrInfo.decode(reader, reader.uint32());
          break;
        case 3:
          message.coins.push(Coin.decode(reader, reader.uint32()));
          break;
        case 4:
          message.startTime = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          break;
        case 5:
          message.distrEpochIdentifier = reader.string();
          break;
        case 6:
          message.numEpochsPaidOver = longToNumber(reader.uint64() as Long);
          break;
        case 7:
          message.filledEpochs = longToNumber(reader.uint64() as Long);
          break;
        case 8:
          message.distributedCoins.push(Coin.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): Stream {
    return {
      id: isSet(object.id) ? Number(object.id) : 0,
      distributeTo: isSet(object.distributeTo) ? DistrInfo.fromJSON(object.distributeTo) : undefined,
      coins: Array.isArray(object?.coins) ? object.coins.map((e: any) => Coin.fromJSON(e)) : [],
      startTime: isSet(object.startTime) ? fromJsonTimestamp(object.startTime) : undefined,
      distrEpochIdentifier: isSet(object.distrEpochIdentifier) ? String(object.distrEpochIdentifier) : "",
      numEpochsPaidOver: isSet(object.numEpochsPaidOver) ? Number(object.numEpochsPaidOver) : 0,
      filledEpochs: isSet(object.filledEpochs) ? Number(object.filledEpochs) : 0,
      distributedCoins: Array.isArray(object?.distributedCoins)
        ? object.distributedCoins.map((e: any) => Coin.fromJSON(e))
        : [],
    };
  },

  toJSON(message: Stream): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = Math.round(message.id));
    message.distributeTo !== undefined
      && (obj.distributeTo = message.distributeTo ? DistrInfo.toJSON(message.distributeTo) : undefined);
    if (message.coins) {
      obj.coins = message.coins.map((e) => e ? Coin.toJSON(e) : undefined);
    } else {
      obj.coins = [];
    }
    message.startTime !== undefined && (obj.startTime = message.startTime.toISOString());
    message.distrEpochIdentifier !== undefined && (obj.distrEpochIdentifier = message.distrEpochIdentifier);
    message.numEpochsPaidOver !== undefined && (obj.numEpochsPaidOver = Math.round(message.numEpochsPaidOver));
    message.filledEpochs !== undefined && (obj.filledEpochs = Math.round(message.filledEpochs));
    if (message.distributedCoins) {
      obj.distributedCoins = message.distributedCoins.map((e) => e ? Coin.toJSON(e) : undefined);
    } else {
      obj.distributedCoins = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<Stream>, I>>(object: I): Stream {
    const message = createBaseStream();
    message.id = object.id ?? 0;
    message.distributeTo = (object.distributeTo !== undefined && object.distributeTo !== null)
      ? DistrInfo.fromPartial(object.distributeTo)
      : undefined;
    message.coins = object.coins?.map((e) => Coin.fromPartial(e)) || [];
    message.startTime = object.startTime ?? undefined;
    message.distrEpochIdentifier = object.distrEpochIdentifier ?? "";
    message.numEpochsPaidOver = object.numEpochsPaidOver ?? 0;
    message.filledEpochs = object.filledEpochs ?? 0;
    message.distributedCoins = object.distributedCoins?.map((e) => Coin.fromPartial(e)) || [];
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
