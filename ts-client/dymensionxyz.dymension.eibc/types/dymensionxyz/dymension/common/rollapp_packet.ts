/* eslint-disable */
import Long from "long";
import _m0 from "protobufjs/minimal";
import { Packet } from "../../../ibc/core/channel/v1/channel";
import { Status, statusFromJSON, statusToJSON } from "./status";

export const protobufPackage = "dymensionxyz.dymension.common";

export interface RollappPacket {
  rollappId: string;
  packet: Packet | undefined;
  acknowledgement: Uint8Array;
  status: Status;
  ProofHeight: number;
  relayer: Uint8Array;
  type: RollappPacket_Type;
  /** stores the result of onAck, onTimeout or onRecv/writeAck */
  error: string;
  /**
   * who was the original person who gets the money (recipient of ics20
   * transfer) of the packet?
   */
  originalTransferTarget: string;
}

export enum RollappPacket_Type {
  ON_RECV = 0,
  ON_ACK = 1,
  ON_TIMEOUT = 2,
  UNDEFINED = -1,
  UNRECOGNIZED = -1,
}

export function rollappPacket_TypeFromJSON(object: any): RollappPacket_Type {
  switch (object) {
    case 0:
    case "ON_RECV":
      return RollappPacket_Type.ON_RECV;
    case 1:
    case "ON_ACK":
      return RollappPacket_Type.ON_ACK;
    case 2:
    case "ON_TIMEOUT":
      return RollappPacket_Type.ON_TIMEOUT;
    case -1:
    case "UNDEFINED":
      return RollappPacket_Type.UNDEFINED;
    case -1:
    case "UNRECOGNIZED":
    default:
      return RollappPacket_Type.UNRECOGNIZED;
  }
}

export function rollappPacket_TypeToJSON(object: RollappPacket_Type): string {
  switch (object) {
    case RollappPacket_Type.ON_RECV:
      return "ON_RECV";
    case RollappPacket_Type.ON_ACK:
      return "ON_ACK";
    case RollappPacket_Type.ON_TIMEOUT:
      return "ON_TIMEOUT";
    case RollappPacket_Type.UNDEFINED:
      return "UNDEFINED";
    case RollappPacket_Type.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

function createBaseRollappPacket(): RollappPacket {
  return {
    rollappId: "",
    packet: undefined,
    acknowledgement: new Uint8Array(),
    status: 0,
    ProofHeight: 0,
    relayer: new Uint8Array(),
    type: 0,
    error: "",
    originalTransferTarget: "",
  };
}

export const RollappPacket = {
  encode(message: RollappPacket, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.rollappId !== "") {
      writer.uint32(10).string(message.rollappId);
    }
    if (message.packet !== undefined) {
      Packet.encode(message.packet, writer.uint32(18).fork()).ldelim();
    }
    if (message.acknowledgement.length !== 0) {
      writer.uint32(26).bytes(message.acknowledgement);
    }
    if (message.status !== 0) {
      writer.uint32(32).int32(message.status);
    }
    if (message.ProofHeight !== 0) {
      writer.uint32(40).uint64(message.ProofHeight);
    }
    if (message.relayer.length !== 0) {
      writer.uint32(50).bytes(message.relayer);
    }
    if (message.type !== 0) {
      writer.uint32(56).int32(message.type);
    }
    if (message.error !== "") {
      writer.uint32(66).string(message.error);
    }
    if (message.originalTransferTarget !== "") {
      writer.uint32(74).string(message.originalTransferTarget);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RollappPacket {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRollappPacket();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.rollappId = reader.string();
          break;
        case 2:
          message.packet = Packet.decode(reader, reader.uint32());
          break;
        case 3:
          message.acknowledgement = reader.bytes();
          break;
        case 4:
          message.status = reader.int32() as any;
          break;
        case 5:
          message.ProofHeight = longToNumber(reader.uint64() as Long);
          break;
        case 6:
          message.relayer = reader.bytes();
          break;
        case 7:
          message.type = reader.int32() as any;
          break;
        case 8:
          message.error = reader.string();
          break;
        case 9:
          message.originalTransferTarget = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): RollappPacket {
    return {
      rollappId: isSet(object.rollappId) ? String(object.rollappId) : "",
      packet: isSet(object.packet) ? Packet.fromJSON(object.packet) : undefined,
      acknowledgement: isSet(object.acknowledgement) ? bytesFromBase64(object.acknowledgement) : new Uint8Array(),
      status: isSet(object.status) ? statusFromJSON(object.status) : 0,
      ProofHeight: isSet(object.ProofHeight) ? Number(object.ProofHeight) : 0,
      relayer: isSet(object.relayer) ? bytesFromBase64(object.relayer) : new Uint8Array(),
      type: isSet(object.type) ? rollappPacket_TypeFromJSON(object.type) : 0,
      error: isSet(object.error) ? String(object.error) : "",
      originalTransferTarget: isSet(object.originalTransferTarget) ? String(object.originalTransferTarget) : "",
    };
  },

  toJSON(message: RollappPacket): unknown {
    const obj: any = {};
    message.rollappId !== undefined && (obj.rollappId = message.rollappId);
    message.packet !== undefined && (obj.packet = message.packet ? Packet.toJSON(message.packet) : undefined);
    message.acknowledgement !== undefined
      && (obj.acknowledgement = base64FromBytes(
        message.acknowledgement !== undefined ? message.acknowledgement : new Uint8Array(),
      ));
    message.status !== undefined && (obj.status = statusToJSON(message.status));
    message.ProofHeight !== undefined && (obj.ProofHeight = Math.round(message.ProofHeight));
    message.relayer !== undefined
      && (obj.relayer = base64FromBytes(message.relayer !== undefined ? message.relayer : new Uint8Array()));
    message.type !== undefined && (obj.type = rollappPacket_TypeToJSON(message.type));
    message.error !== undefined && (obj.error = message.error);
    message.originalTransferTarget !== undefined && (obj.originalTransferTarget = message.originalTransferTarget);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<RollappPacket>, I>>(object: I): RollappPacket {
    const message = createBaseRollappPacket();
    message.rollappId = object.rollappId ?? "";
    message.packet = (object.packet !== undefined && object.packet !== null)
      ? Packet.fromPartial(object.packet)
      : undefined;
    message.acknowledgement = object.acknowledgement ?? new Uint8Array();
    message.status = object.status ?? 0;
    message.ProofHeight = object.ProofHeight ?? 0;
    message.relayer = object.relayer ?? new Uint8Array();
    message.type = object.type ?? 0;
    message.error = object.error ?? "";
    message.originalTransferTarget = object.originalTransferTarget ?? "";
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

function bytesFromBase64(b64: string): Uint8Array {
  if (globalThis.Buffer) {
    return Uint8Array.from(globalThis.Buffer.from(b64, "base64"));
  } else {
    const bin = globalThis.atob(b64);
    const arr = new Uint8Array(bin.length);
    for (let i = 0; i < bin.length; ++i) {
      arr[i] = bin.charCodeAt(i);
    }
    return arr;
  }
}

function base64FromBytes(arr: Uint8Array): string {
  if (globalThis.Buffer) {
    return globalThis.Buffer.from(arr).toString("base64");
  } else {
    const bin: string[] = [];
    arr.forEach((byte) => {
      bin.push(String.fromCharCode(byte));
    });
    return globalThis.btoa(bin.join(""));
  }
}

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
