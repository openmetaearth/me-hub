/* eslint-disable */
import _m0 from "protobufjs/minimal";
import { Service } from "../did/issuer";

export const protobufPackage = "metaearth.kyc";

export interface Region {
  id: string;
  name: string;
  authType: string[];
}

export interface Protocol {
  service: Service | undefined;
  regions: Region[];
}

function createBaseRegion(): Region {
  return { id: "", name: "", authType: [] };
}

export const Region = {
  encode(message: Region, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== "") {
      writer.uint32(10).string(message.id);
    }
    if (message.name !== "") {
      writer.uint32(18).string(message.name);
    }
    for (const v of message.authType) {
      writer.uint32(26).string(v!);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Region {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRegion();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.id = reader.string();
          break;
        case 2:
          message.name = reader.string();
          break;
        case 3:
          message.authType.push(reader.string());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): Region {
    return {
      id: isSet(object.id) ? String(object.id) : "",
      name: isSet(object.name) ? String(object.name) : "",
      authType: Array.isArray(object?.authType) ? object.authType.map((e: any) => String(e)) : [],
    };
  },

  toJSON(message: Region): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = message.id);
    message.name !== undefined && (obj.name = message.name);
    if (message.authType) {
      obj.authType = message.authType.map((e) => e);
    } else {
      obj.authType = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<Region>, I>>(object: I): Region {
    const message = createBaseRegion();
    message.id = object.id ?? "";
    message.name = object.name ?? "";
    message.authType = object.authType?.map((e) => e) || [];
    return message;
  },
};

function createBaseProtocol(): Protocol {
  return { service: undefined, regions: [] };
}

export const Protocol = {
  encode(message: Protocol, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.service !== undefined) {
      Service.encode(message.service, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.regions) {
      Region.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Protocol {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseProtocol();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.service = Service.decode(reader, reader.uint32());
          break;
        case 2:
          message.regions.push(Region.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): Protocol {
    return {
      service: isSet(object.service) ? Service.fromJSON(object.service) : undefined,
      regions: Array.isArray(object?.regions) ? object.regions.map((e: any) => Region.fromJSON(e)) : [],
    };
  },

  toJSON(message: Protocol): unknown {
    const obj: any = {};
    message.service !== undefined && (obj.service = message.service ? Service.toJSON(message.service) : undefined);
    if (message.regions) {
      obj.regions = message.regions.map((e) => e ? Region.toJSON(e) : undefined);
    } else {
      obj.regions = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<Protocol>, I>>(object: I): Protocol {
    const message = createBaseProtocol();
    message.service = (object.service !== undefined && object.service !== null)
      ? Service.fromPartial(object.service)
      : undefined;
    message.regions = object.regions?.map((e) => Region.fromPartial(e)) || [];
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
