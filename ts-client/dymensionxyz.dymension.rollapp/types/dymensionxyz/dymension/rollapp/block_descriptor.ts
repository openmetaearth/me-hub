/* eslint-disable */
import Long from "long";
import _m0 from "protobufjs/minimal";

export const protobufPackage = "dymensionxyz.dymension.rollapp";

/** BlockDescriptor defines a single rollapp chain block description. */
export interface BlockDescriptor {
  /** height is the height of the block */
  height: number;
  /**
   * stateRoot is a 32 byte array of the hash of the block (state root of the
   * block)
   */
  stateRoot: Uint8Array;
}

/** BlockDescriptors defines list of BlockDescriptor. */
export interface BlockDescriptors {
  BD: BlockDescriptor[];
}

function createBaseBlockDescriptor(): BlockDescriptor {
  return { height: 0, stateRoot: new Uint8Array() };
}

export const BlockDescriptor = {
  encode(message: BlockDescriptor, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.height !== 0) {
      writer.uint32(8).uint64(message.height);
    }
    if (message.stateRoot.length !== 0) {
      writer.uint32(18).bytes(message.stateRoot);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): BlockDescriptor {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseBlockDescriptor();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.height = longToNumber(reader.uint64() as Long);
          break;
        case 2:
          message.stateRoot = reader.bytes();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): BlockDescriptor {
    return {
      height: isSet(object.height) ? Number(object.height) : 0,
      stateRoot: isSet(object.stateRoot) ? bytesFromBase64(object.stateRoot) : new Uint8Array(),
    };
  },

  toJSON(message: BlockDescriptor): unknown {
    const obj: any = {};
    message.height !== undefined && (obj.height = Math.round(message.height));
    message.stateRoot !== undefined
      && (obj.stateRoot = base64FromBytes(message.stateRoot !== undefined ? message.stateRoot : new Uint8Array()));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<BlockDescriptor>, I>>(object: I): BlockDescriptor {
    const message = createBaseBlockDescriptor();
    message.height = object.height ?? 0;
    message.stateRoot = object.stateRoot ?? new Uint8Array();
    return message;
  },
};

function createBaseBlockDescriptors(): BlockDescriptors {
  return { BD: [] };
}

export const BlockDescriptors = {
  encode(message: BlockDescriptors, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.BD) {
      BlockDescriptor.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): BlockDescriptors {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseBlockDescriptors();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.BD.push(BlockDescriptor.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): BlockDescriptors {
    return { BD: Array.isArray(object?.BD) ? object.BD.map((e: any) => BlockDescriptor.fromJSON(e)) : [] };
  },

  toJSON(message: BlockDescriptors): unknown {
    const obj: any = {};
    if (message.BD) {
      obj.BD = message.BD.map((e) => e ? BlockDescriptor.toJSON(e) : undefined);
    } else {
      obj.BD = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<BlockDescriptors>, I>>(object: I): BlockDescriptors {
    const message = createBaseBlockDescriptors();
    message.BD = object.BD?.map((e) => BlockDescriptor.fromPartial(e)) || [];
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
