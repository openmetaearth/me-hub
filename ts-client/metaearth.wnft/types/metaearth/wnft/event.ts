/* eslint-disable */
import _m0 from "protobufjs/minimal";

export const protobufPackage = "metaearth.wnft";

/** EventNewClass is emitted on NewClass */
export interface EventNewClass {
  /** class_id associated with the nft */
  classId: string;
  totalSupply: string;
}

function createBaseEventNewClass(): EventNewClass {
  return { classId: "", totalSupply: "" };
}

export const EventNewClass = {
  encode(message: EventNewClass, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.classId !== "") {
      writer.uint32(10).string(message.classId);
    }
    if (message.totalSupply !== "") {
      writer.uint32(18).string(message.totalSupply);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): EventNewClass {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseEventNewClass();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.classId = reader.string();
          break;
        case 2:
          message.totalSupply = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): EventNewClass {
    return {
      classId: isSet(object.classId) ? String(object.classId) : "",
      totalSupply: isSet(object.totalSupply) ? String(object.totalSupply) : "",
    };
  },

  toJSON(message: EventNewClass): unknown {
    const obj: any = {};
    message.classId !== undefined && (obj.classId = message.classId);
    message.totalSupply !== undefined && (obj.totalSupply = message.totalSupply);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<EventNewClass>, I>>(object: I): EventNewClass {
    const message = createBaseEventNewClass();
    message.classId = object.classId ?? "";
    message.totalSupply = object.totalSupply ?? "";
    return message;
  },
};

type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;

export type DeepPartial<T> = T extends Builtin ? T
  : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
  : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

type KeysOfUnion<T> = T extends T ? keyof T : never;
export type Exact<P, I extends P> = P extends Builtin ? P
  : P & { [K in keyof P]: Exact<P[K], I[K]> } & { [K in Exclude<keyof I, KeysOfUnion<P>>]: never };

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}
