/* eslint-disable */
import Long from "long";
import _m0 from "protobufjs/minimal";
import { Coin } from "../../cosmos/base/v1beta1/coin";
import { Timestamp } from "../../google/protobuf/timestamp";
import { DistrRecord } from "./distr_info";

export const protobufPackage = "metaearth.streamer";

export interface CreateStreamProposal {
  title: string;
  description: string;
  distributeToRecords: DistrRecord[];
  /** coins are coin(s) to be distributed by the stream */
  coins: Coin[];
  /** start_time is the distribution start time */
  startTime: Date | undefined;
  distrEpochIdentifier: string;
  /**
   * num_epochs_paid_over is the number of epochs distribution will be completed
   * over
   */
  numEpochsPaidOver: number;
}

export interface TerminateStreamProposal {
  title: string;
  description: string;
  streamId: number;
}

function createBaseCreateStreamProposal(): CreateStreamProposal {
  return {
    title: "",
    description: "",
    distributeToRecords: [],
    coins: [],
    startTime: undefined,
    distrEpochIdentifier: "",
    numEpochsPaidOver: 0,
  };
}

export const CreateStreamProposal = {
  encode(message: CreateStreamProposal, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.title !== "") {
      writer.uint32(10).string(message.title);
    }
    if (message.description !== "") {
      writer.uint32(18).string(message.description);
    }
    for (const v of message.distributeToRecords) {
      DistrRecord.encode(v!, writer.uint32(26).fork()).ldelim();
    }
    for (const v of message.coins) {
      Coin.encode(v!, writer.uint32(34).fork()).ldelim();
    }
    if (message.startTime !== undefined) {
      Timestamp.encode(toTimestamp(message.startTime), writer.uint32(42).fork()).ldelim();
    }
    if (message.distrEpochIdentifier !== "") {
      writer.uint32(50).string(message.distrEpochIdentifier);
    }
    if (message.numEpochsPaidOver !== 0) {
      writer.uint32(56).uint64(message.numEpochsPaidOver);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CreateStreamProposal {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCreateStreamProposal();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.title = reader.string();
          break;
        case 2:
          message.description = reader.string();
          break;
        case 3:
          message.distributeToRecords.push(DistrRecord.decode(reader, reader.uint32()));
          break;
        case 4:
          message.coins.push(Coin.decode(reader, reader.uint32()));
          break;
        case 5:
          message.startTime = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          break;
        case 6:
          message.distrEpochIdentifier = reader.string();
          break;
        case 7:
          message.numEpochsPaidOver = longToNumber(reader.uint64() as Long);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): CreateStreamProposal {
    return {
      title: isSet(object.title) ? String(object.title) : "",
      description: isSet(object.description) ? String(object.description) : "",
      distributeToRecords: Array.isArray(object?.distributeToRecords)
        ? object.distributeToRecords.map((e: any) => DistrRecord.fromJSON(e))
        : [],
      coins: Array.isArray(object?.coins) ? object.coins.map((e: any) => Coin.fromJSON(e)) : [],
      startTime: isSet(object.startTime) ? fromJsonTimestamp(object.startTime) : undefined,
      distrEpochIdentifier: isSet(object.distrEpochIdentifier) ? String(object.distrEpochIdentifier) : "",
      numEpochsPaidOver: isSet(object.numEpochsPaidOver) ? Number(object.numEpochsPaidOver) : 0,
    };
  },

  toJSON(message: CreateStreamProposal): unknown {
    const obj: any = {};
    message.title !== undefined && (obj.title = message.title);
    message.description !== undefined && (obj.description = message.description);
    if (message.distributeToRecords) {
      obj.distributeToRecords = message.distributeToRecords.map((e) => e ? DistrRecord.toJSON(e) : undefined);
    } else {
      obj.distributeToRecords = [];
    }
    if (message.coins) {
      obj.coins = message.coins.map((e) => e ? Coin.toJSON(e) : undefined);
    } else {
      obj.coins = [];
    }
    message.startTime !== undefined && (obj.startTime = message.startTime.toISOString());
    message.distrEpochIdentifier !== undefined && (obj.distrEpochIdentifier = message.distrEpochIdentifier);
    message.numEpochsPaidOver !== undefined && (obj.numEpochsPaidOver = Math.round(message.numEpochsPaidOver));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<CreateStreamProposal>, I>>(object: I): CreateStreamProposal {
    const message = createBaseCreateStreamProposal();
    message.title = object.title ?? "";
    message.description = object.description ?? "";
    message.distributeToRecords = object.distributeToRecords?.map((e) => DistrRecord.fromPartial(e)) || [];
    message.coins = object.coins?.map((e) => Coin.fromPartial(e)) || [];
    message.startTime = object.startTime ?? undefined;
    message.distrEpochIdentifier = object.distrEpochIdentifier ?? "";
    message.numEpochsPaidOver = object.numEpochsPaidOver ?? 0;
    return message;
  },
};

function createBaseTerminateStreamProposal(): TerminateStreamProposal {
  return { title: "", description: "", streamId: 0 };
}

export const TerminateStreamProposal = {
  encode(message: TerminateStreamProposal, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.title !== "") {
      writer.uint32(10).string(message.title);
    }
    if (message.description !== "") {
      writer.uint32(18).string(message.description);
    }
    if (message.streamId !== 0) {
      writer.uint32(32).uint64(message.streamId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): TerminateStreamProposal {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseTerminateStreamProposal();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.title = reader.string();
          break;
        case 2:
          message.description = reader.string();
          break;
        case 4:
          message.streamId = longToNumber(reader.uint64() as Long);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): TerminateStreamProposal {
    return {
      title: isSet(object.title) ? String(object.title) : "",
      description: isSet(object.description) ? String(object.description) : "",
      streamId: isSet(object.streamId) ? Number(object.streamId) : 0,
    };
  },

  toJSON(message: TerminateStreamProposal): unknown {
    const obj: any = {};
    message.title !== undefined && (obj.title = message.title);
    message.description !== undefined && (obj.description = message.description);
    message.streamId !== undefined && (obj.streamId = Math.round(message.streamId));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<TerminateStreamProposal>, I>>(object: I): TerminateStreamProposal {
    const message = createBaseTerminateStreamProposal();
    message.title = object.title ?? "";
    message.description = object.description ?? "";
    message.streamId = object.streamId ?? 0;
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
