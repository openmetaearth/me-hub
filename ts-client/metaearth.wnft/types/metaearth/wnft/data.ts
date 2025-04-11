/* eslint-disable */
import _m0 from "protobufjs/minimal";

export const protobufPackage = "metaearth.wnft";

export interface Extension {
  data: string;
}

export interface ClassMetadata {
  creator: string;
}

function createBaseExtension(): Extension {
  return { data: "" };
}

export const Extension = {
  encode(message: Extension, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.data !== "") {
      writer.uint32(10).string(message.data);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Extension {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseExtension();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.data = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): Extension {
    return { data: isSet(object.data) ? String(object.data) : "" };
  },

  toJSON(message: Extension): unknown {
    const obj: any = {};
    message.data !== undefined && (obj.data = message.data);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<Extension>, I>>(object: I): Extension {
    const message = createBaseExtension();
    message.data = object.data ?? "";
    return message;
  },
};

function createBaseClassMetadata(): ClassMetadata {
  return { creator: "" };
}

export const ClassMetadata = {
  encode(message: ClassMetadata, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ClassMetadata {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseClassMetadata();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ClassMetadata {
    return { creator: isSet(object.creator) ? String(object.creator) : "" };
  },

  toJSON(message: ClassMetadata): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<ClassMetadata>, I>>(object: I): ClassMetadata {
    const message = createBaseClassMetadata();
    message.creator = object.creator ?? "";
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
