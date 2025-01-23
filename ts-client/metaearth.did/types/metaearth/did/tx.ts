/* eslint-disable */
import _m0 from "protobufjs/minimal";
import { DidStatus, didStatusFromJSON, didStatusToJSON } from "./did";
import { ServiceStatus, serviceStatusFromJSON, serviceStatusToJSON } from "./issuer";

export const protobufPackage = "metaearth.did";

export interface MsgCreateDid {
  creator: string;
  did: string;
  address: string;
  pubkey: string;
}

export interface MsgCreateDidResponse {
}

export interface MsgUpdateDidStatus {
  creator: string;
  did: string;
  status: DidStatus;
}

export interface MsgUpdateDidStatusResponse {
}

export interface MsgRemoveDid {
  creator: string;
  did: string;
}

export interface MsgRemoveDidResponse {
}

export interface MsgCreateService {
  creator: string;
  sid: string;
  name: string;
  description: string;
  issuers: string[];
}

export interface MsgCreateServiceResponse {
}

export interface MsgUpdateServiceStatus {
  creator: string;
  sid: string;
  status: ServiceStatus;
}

export interface MsgUpdateServiceStatusResponse {
}

export interface MsgRemoveService {
  creator: string;
  sid: string;
}

export interface MsgRemoveServiceResponse {
}

export interface MsgCreateVC {
  issuer: string;
  did: string;
  sid: string;
  hash: string;
  uri: string;
  data: Uint8Array;
  filters: Uint8Array[];
}

export interface MsgCreateVCResponse {
}

export interface MsgUpdateVC {
  issuer: string;
  did: string;
  sid: string;
  hash: string;
  uri: string;
  data: Uint8Array;
  filters: Uint8Array[];
}

export interface MsgUpdateVCResponse {
}

export interface MsgRemoveVC {
  issuer: string;
  did: string;
  sid: string;
}

export interface MsgRemoveVCResponse {
}

function createBaseMsgCreateDid(): MsgCreateDid {
  return { creator: "", did: "", address: "", pubkey: "" };
}

export const MsgCreateDid = {
  encode(message: MsgCreateDid, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.did !== "") {
      writer.uint32(18).string(message.did);
    }
    if (message.address !== "") {
      writer.uint32(26).string(message.address);
    }
    if (message.pubkey !== "") {
      writer.uint32(34).string(message.pubkey);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgCreateDid {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgCreateDid();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.did = reader.string();
          break;
        case 3:
          message.address = reader.string();
          break;
        case 4:
          message.pubkey = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgCreateDid {
    return {
      creator: isSet(object.creator) ? String(object.creator) : "",
      did: isSet(object.did) ? String(object.did) : "",
      address: isSet(object.address) ? String(object.address) : "",
      pubkey: isSet(object.pubkey) ? String(object.pubkey) : "",
    };
  },

  toJSON(message: MsgCreateDid): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    message.did !== undefined && (obj.did = message.did);
    message.address !== undefined && (obj.address = message.address);
    message.pubkey !== undefined && (obj.pubkey = message.pubkey);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgCreateDid>, I>>(object: I): MsgCreateDid {
    const message = createBaseMsgCreateDid();
    message.creator = object.creator ?? "";
    message.did = object.did ?? "";
    message.address = object.address ?? "";
    message.pubkey = object.pubkey ?? "";
    return message;
  },
};

function createBaseMsgCreateDidResponse(): MsgCreateDidResponse {
  return {};
}

export const MsgCreateDidResponse = {
  encode(_: MsgCreateDidResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgCreateDidResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgCreateDidResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): MsgCreateDidResponse {
    return {};
  },

  toJSON(_: MsgCreateDidResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgCreateDidResponse>, I>>(_: I): MsgCreateDidResponse {
    const message = createBaseMsgCreateDidResponse();
    return message;
  },
};

function createBaseMsgUpdateDidStatus(): MsgUpdateDidStatus {
  return { creator: "", did: "", status: 0 };
}

export const MsgUpdateDidStatus = {
  encode(message: MsgUpdateDidStatus, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.did !== "") {
      writer.uint32(18).string(message.did);
    }
    if (message.status !== 0) {
      writer.uint32(24).int32(message.status);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateDidStatus {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateDidStatus();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.did = reader.string();
          break;
        case 3:
          message.status = reader.int32() as any;
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgUpdateDidStatus {
    return {
      creator: isSet(object.creator) ? String(object.creator) : "",
      did: isSet(object.did) ? String(object.did) : "",
      status: isSet(object.status) ? didStatusFromJSON(object.status) : 0,
    };
  },

  toJSON(message: MsgUpdateDidStatus): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    message.did !== undefined && (obj.did = message.did);
    message.status !== undefined && (obj.status = didStatusToJSON(message.status));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgUpdateDidStatus>, I>>(object: I): MsgUpdateDidStatus {
    const message = createBaseMsgUpdateDidStatus();
    message.creator = object.creator ?? "";
    message.did = object.did ?? "";
    message.status = object.status ?? 0;
    return message;
  },
};

function createBaseMsgUpdateDidStatusResponse(): MsgUpdateDidStatusResponse {
  return {};
}

export const MsgUpdateDidStatusResponse = {
  encode(_: MsgUpdateDidStatusResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateDidStatusResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateDidStatusResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): MsgUpdateDidStatusResponse {
    return {};
  },

  toJSON(_: MsgUpdateDidStatusResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgUpdateDidStatusResponse>, I>>(_: I): MsgUpdateDidStatusResponse {
    const message = createBaseMsgUpdateDidStatusResponse();
    return message;
  },
};

function createBaseMsgRemoveDid(): MsgRemoveDid {
  return { creator: "", did: "" };
}

export const MsgRemoveDid = {
  encode(message: MsgRemoveDid, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.did !== "") {
      writer.uint32(18).string(message.did);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgRemoveDid {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgRemoveDid();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.did = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgRemoveDid {
    return {
      creator: isSet(object.creator) ? String(object.creator) : "",
      did: isSet(object.did) ? String(object.did) : "",
    };
  },

  toJSON(message: MsgRemoveDid): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    message.did !== undefined && (obj.did = message.did);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgRemoveDid>, I>>(object: I): MsgRemoveDid {
    const message = createBaseMsgRemoveDid();
    message.creator = object.creator ?? "";
    message.did = object.did ?? "";
    return message;
  },
};

function createBaseMsgRemoveDidResponse(): MsgRemoveDidResponse {
  return {};
}

export const MsgRemoveDidResponse = {
  encode(_: MsgRemoveDidResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgRemoveDidResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgRemoveDidResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): MsgRemoveDidResponse {
    return {};
  },

  toJSON(_: MsgRemoveDidResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgRemoveDidResponse>, I>>(_: I): MsgRemoveDidResponse {
    const message = createBaseMsgRemoveDidResponse();
    return message;
  },
};

function createBaseMsgCreateService(): MsgCreateService {
  return { creator: "", sid: "", name: "", description: "", issuers: [] };
}

export const MsgCreateService = {
  encode(message: MsgCreateService, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.sid !== "") {
      writer.uint32(18).string(message.sid);
    }
    if (message.name !== "") {
      writer.uint32(26).string(message.name);
    }
    if (message.description !== "") {
      writer.uint32(34).string(message.description);
    }
    for (const v of message.issuers) {
      writer.uint32(42).string(v!);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgCreateService {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgCreateService();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.sid = reader.string();
          break;
        case 3:
          message.name = reader.string();
          break;
        case 4:
          message.description = reader.string();
          break;
        case 5:
          message.issuers.push(reader.string());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgCreateService {
    return {
      creator: isSet(object.creator) ? String(object.creator) : "",
      sid: isSet(object.sid) ? String(object.sid) : "",
      name: isSet(object.name) ? String(object.name) : "",
      description: isSet(object.description) ? String(object.description) : "",
      issuers: Array.isArray(object?.issuers) ? object.issuers.map((e: any) => String(e)) : [],
    };
  },

  toJSON(message: MsgCreateService): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    message.sid !== undefined && (obj.sid = message.sid);
    message.name !== undefined && (obj.name = message.name);
    message.description !== undefined && (obj.description = message.description);
    if (message.issuers) {
      obj.issuers = message.issuers.map((e) => e);
    } else {
      obj.issuers = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgCreateService>, I>>(object: I): MsgCreateService {
    const message = createBaseMsgCreateService();
    message.creator = object.creator ?? "";
    message.sid = object.sid ?? "";
    message.name = object.name ?? "";
    message.description = object.description ?? "";
    message.issuers = object.issuers?.map((e) => e) || [];
    return message;
  },
};

function createBaseMsgCreateServiceResponse(): MsgCreateServiceResponse {
  return {};
}

export const MsgCreateServiceResponse = {
  encode(_: MsgCreateServiceResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgCreateServiceResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgCreateServiceResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): MsgCreateServiceResponse {
    return {};
  },

  toJSON(_: MsgCreateServiceResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgCreateServiceResponse>, I>>(_: I): MsgCreateServiceResponse {
    const message = createBaseMsgCreateServiceResponse();
    return message;
  },
};

function createBaseMsgUpdateServiceStatus(): MsgUpdateServiceStatus {
  return { creator: "", sid: "", status: 0 };
}

export const MsgUpdateServiceStatus = {
  encode(message: MsgUpdateServiceStatus, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.sid !== "") {
      writer.uint32(18).string(message.sid);
    }
    if (message.status !== 0) {
      writer.uint32(24).int32(message.status);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateServiceStatus {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateServiceStatus();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.sid = reader.string();
          break;
        case 3:
          message.status = reader.int32() as any;
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgUpdateServiceStatus {
    return {
      creator: isSet(object.creator) ? String(object.creator) : "",
      sid: isSet(object.sid) ? String(object.sid) : "",
      status: isSet(object.status) ? serviceStatusFromJSON(object.status) : 0,
    };
  },

  toJSON(message: MsgUpdateServiceStatus): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    message.sid !== undefined && (obj.sid = message.sid);
    message.status !== undefined && (obj.status = serviceStatusToJSON(message.status));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgUpdateServiceStatus>, I>>(object: I): MsgUpdateServiceStatus {
    const message = createBaseMsgUpdateServiceStatus();
    message.creator = object.creator ?? "";
    message.sid = object.sid ?? "";
    message.status = object.status ?? 0;
    return message;
  },
};

function createBaseMsgUpdateServiceStatusResponse(): MsgUpdateServiceStatusResponse {
  return {};
}

export const MsgUpdateServiceStatusResponse = {
  encode(_: MsgUpdateServiceStatusResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateServiceStatusResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateServiceStatusResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): MsgUpdateServiceStatusResponse {
    return {};
  },

  toJSON(_: MsgUpdateServiceStatusResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgUpdateServiceStatusResponse>, I>>(_: I): MsgUpdateServiceStatusResponse {
    const message = createBaseMsgUpdateServiceStatusResponse();
    return message;
  },
};

function createBaseMsgRemoveService(): MsgRemoveService {
  return { creator: "", sid: "" };
}

export const MsgRemoveService = {
  encode(message: MsgRemoveService, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.sid !== "") {
      writer.uint32(18).string(message.sid);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgRemoveService {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgRemoveService();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.sid = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgRemoveService {
    return {
      creator: isSet(object.creator) ? String(object.creator) : "",
      sid: isSet(object.sid) ? String(object.sid) : "",
    };
  },

  toJSON(message: MsgRemoveService): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    message.sid !== undefined && (obj.sid = message.sid);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgRemoveService>, I>>(object: I): MsgRemoveService {
    const message = createBaseMsgRemoveService();
    message.creator = object.creator ?? "";
    message.sid = object.sid ?? "";
    return message;
  },
};

function createBaseMsgRemoveServiceResponse(): MsgRemoveServiceResponse {
  return {};
}

export const MsgRemoveServiceResponse = {
  encode(_: MsgRemoveServiceResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgRemoveServiceResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgRemoveServiceResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): MsgRemoveServiceResponse {
    return {};
  },

  toJSON(_: MsgRemoveServiceResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgRemoveServiceResponse>, I>>(_: I): MsgRemoveServiceResponse {
    const message = createBaseMsgRemoveServiceResponse();
    return message;
  },
};

function createBaseMsgCreateVC(): MsgCreateVC {
  return { issuer: "", did: "", sid: "", hash: "", uri: "", data: new Uint8Array(), filters: [] };
}

export const MsgCreateVC = {
  encode(message: MsgCreateVC, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.issuer !== "") {
      writer.uint32(10).string(message.issuer);
    }
    if (message.did !== "") {
      writer.uint32(18).string(message.did);
    }
    if (message.sid !== "") {
      writer.uint32(26).string(message.sid);
    }
    if (message.hash !== "") {
      writer.uint32(34).string(message.hash);
    }
    if (message.uri !== "") {
      writer.uint32(42).string(message.uri);
    }
    if (message.data.length !== 0) {
      writer.uint32(50).bytes(message.data);
    }
    for (const v of message.filters) {
      writer.uint32(58).bytes(v!);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgCreateVC {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgCreateVC();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.issuer = reader.string();
          break;
        case 2:
          message.did = reader.string();
          break;
        case 3:
          message.sid = reader.string();
          break;
        case 4:
          message.hash = reader.string();
          break;
        case 5:
          message.uri = reader.string();
          break;
        case 6:
          message.data = reader.bytes();
          break;
        case 7:
          message.filters.push(reader.bytes());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgCreateVC {
    return {
      issuer: isSet(object.issuer) ? String(object.issuer) : "",
      did: isSet(object.did) ? String(object.did) : "",
      sid: isSet(object.sid) ? String(object.sid) : "",
      hash: isSet(object.hash) ? String(object.hash) : "",
      uri: isSet(object.uri) ? String(object.uri) : "",
      data: isSet(object.data) ? bytesFromBase64(object.data) : new Uint8Array(),
      filters: Array.isArray(object?.filters) ? object.filters.map((e: any) => bytesFromBase64(e)) : [],
    };
  },

  toJSON(message: MsgCreateVC): unknown {
    const obj: any = {};
    message.issuer !== undefined && (obj.issuer = message.issuer);
    message.did !== undefined && (obj.did = message.did);
    message.sid !== undefined && (obj.sid = message.sid);
    message.hash !== undefined && (obj.hash = message.hash);
    message.uri !== undefined && (obj.uri = message.uri);
    message.data !== undefined
      && (obj.data = base64FromBytes(message.data !== undefined ? message.data : new Uint8Array()));
    if (message.filters) {
      obj.filters = message.filters.map((e) => base64FromBytes(e !== undefined ? e : new Uint8Array()));
    } else {
      obj.filters = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgCreateVC>, I>>(object: I): MsgCreateVC {
    const message = createBaseMsgCreateVC();
    message.issuer = object.issuer ?? "";
    message.did = object.did ?? "";
    message.sid = object.sid ?? "";
    message.hash = object.hash ?? "";
    message.uri = object.uri ?? "";
    message.data = object.data ?? new Uint8Array();
    message.filters = object.filters?.map((e) => e) || [];
    return message;
  },
};

function createBaseMsgCreateVCResponse(): MsgCreateVCResponse {
  return {};
}

export const MsgCreateVCResponse = {
  encode(_: MsgCreateVCResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgCreateVCResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgCreateVCResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): MsgCreateVCResponse {
    return {};
  },

  toJSON(_: MsgCreateVCResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgCreateVCResponse>, I>>(_: I): MsgCreateVCResponse {
    const message = createBaseMsgCreateVCResponse();
    return message;
  },
};

function createBaseMsgUpdateVC(): MsgUpdateVC {
  return { issuer: "", did: "", sid: "", hash: "", uri: "", data: new Uint8Array(), filters: [] };
}

export const MsgUpdateVC = {
  encode(message: MsgUpdateVC, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.issuer !== "") {
      writer.uint32(10).string(message.issuer);
    }
    if (message.did !== "") {
      writer.uint32(18).string(message.did);
    }
    if (message.sid !== "") {
      writer.uint32(26).string(message.sid);
    }
    if (message.hash !== "") {
      writer.uint32(34).string(message.hash);
    }
    if (message.uri !== "") {
      writer.uint32(42).string(message.uri);
    }
    if (message.data.length !== 0) {
      writer.uint32(50).bytes(message.data);
    }
    for (const v of message.filters) {
      writer.uint32(58).bytes(v!);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateVC {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateVC();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.issuer = reader.string();
          break;
        case 2:
          message.did = reader.string();
          break;
        case 3:
          message.sid = reader.string();
          break;
        case 4:
          message.hash = reader.string();
          break;
        case 5:
          message.uri = reader.string();
          break;
        case 6:
          message.data = reader.bytes();
          break;
        case 7:
          message.filters.push(reader.bytes());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgUpdateVC {
    return {
      issuer: isSet(object.issuer) ? String(object.issuer) : "",
      did: isSet(object.did) ? String(object.did) : "",
      sid: isSet(object.sid) ? String(object.sid) : "",
      hash: isSet(object.hash) ? String(object.hash) : "",
      uri: isSet(object.uri) ? String(object.uri) : "",
      data: isSet(object.data) ? bytesFromBase64(object.data) : new Uint8Array(),
      filters: Array.isArray(object?.filters) ? object.filters.map((e: any) => bytesFromBase64(e)) : [],
    };
  },

  toJSON(message: MsgUpdateVC): unknown {
    const obj: any = {};
    message.issuer !== undefined && (obj.issuer = message.issuer);
    message.did !== undefined && (obj.did = message.did);
    message.sid !== undefined && (obj.sid = message.sid);
    message.hash !== undefined && (obj.hash = message.hash);
    message.uri !== undefined && (obj.uri = message.uri);
    message.data !== undefined
      && (obj.data = base64FromBytes(message.data !== undefined ? message.data : new Uint8Array()));
    if (message.filters) {
      obj.filters = message.filters.map((e) => base64FromBytes(e !== undefined ? e : new Uint8Array()));
    } else {
      obj.filters = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgUpdateVC>, I>>(object: I): MsgUpdateVC {
    const message = createBaseMsgUpdateVC();
    message.issuer = object.issuer ?? "";
    message.did = object.did ?? "";
    message.sid = object.sid ?? "";
    message.hash = object.hash ?? "";
    message.uri = object.uri ?? "";
    message.data = object.data ?? new Uint8Array();
    message.filters = object.filters?.map((e) => e) || [];
    return message;
  },
};

function createBaseMsgUpdateVCResponse(): MsgUpdateVCResponse {
  return {};
}

export const MsgUpdateVCResponse = {
  encode(_: MsgUpdateVCResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateVCResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateVCResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): MsgUpdateVCResponse {
    return {};
  },

  toJSON(_: MsgUpdateVCResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgUpdateVCResponse>, I>>(_: I): MsgUpdateVCResponse {
    const message = createBaseMsgUpdateVCResponse();
    return message;
  },
};

function createBaseMsgRemoveVC(): MsgRemoveVC {
  return { issuer: "", did: "", sid: "" };
}

export const MsgRemoveVC = {
  encode(message: MsgRemoveVC, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.issuer !== "") {
      writer.uint32(10).string(message.issuer);
    }
    if (message.did !== "") {
      writer.uint32(18).string(message.did);
    }
    if (message.sid !== "") {
      writer.uint32(26).string(message.sid);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgRemoveVC {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgRemoveVC();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.issuer = reader.string();
          break;
        case 2:
          message.did = reader.string();
          break;
        case 3:
          message.sid = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgRemoveVC {
    return {
      issuer: isSet(object.issuer) ? String(object.issuer) : "",
      did: isSet(object.did) ? String(object.did) : "",
      sid: isSet(object.sid) ? String(object.sid) : "",
    };
  },

  toJSON(message: MsgRemoveVC): unknown {
    const obj: any = {};
    message.issuer !== undefined && (obj.issuer = message.issuer);
    message.did !== undefined && (obj.did = message.did);
    message.sid !== undefined && (obj.sid = message.sid);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgRemoveVC>, I>>(object: I): MsgRemoveVC {
    const message = createBaseMsgRemoveVC();
    message.issuer = object.issuer ?? "";
    message.did = object.did ?? "";
    message.sid = object.sid ?? "";
    return message;
  },
};

function createBaseMsgRemoveVCResponse(): MsgRemoveVCResponse {
  return {};
}

export const MsgRemoveVCResponse = {
  encode(_: MsgRemoveVCResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgRemoveVCResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgRemoveVCResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): MsgRemoveVCResponse {
    return {};
  },

  toJSON(_: MsgRemoveVCResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgRemoveVCResponse>, I>>(_: I): MsgRemoveVCResponse {
    const message = createBaseMsgRemoveVCResponse();
    return message;
  },
};

/** Msg defines the Msg service. */
export interface Msg {
  CreateDid(request: MsgCreateDid): Promise<MsgCreateDidResponse>;
  /** rpc RemoveDid(MsgRemoveDid) returns (MsgRemoveDidResponse); */
  UpdateDidStatus(request: MsgUpdateDidStatus): Promise<MsgUpdateDidStatusResponse>;
  CreateService(request: MsgCreateService): Promise<MsgCreateServiceResponse>;
  /** rpc RemoveService(MsgRemoveService) returns (MsgRemoveServiceResponse); */
  UpdateServiceStatus(request: MsgUpdateServiceStatus): Promise<MsgUpdateServiceStatusResponse>;
  CreateVC(request: MsgCreateVC): Promise<MsgCreateVCResponse>;
  UpdateVC(request: MsgUpdateVC): Promise<MsgUpdateVCResponse>;
  RemoveVC(request: MsgRemoveVC): Promise<MsgRemoveVCResponse>;
}

export class MsgClientImpl implements Msg {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.CreateDid = this.CreateDid.bind(this);
    this.UpdateDidStatus = this.UpdateDidStatus.bind(this);
    this.CreateService = this.CreateService.bind(this);
    this.UpdateServiceStatus = this.UpdateServiceStatus.bind(this);
    this.CreateVC = this.CreateVC.bind(this);
    this.UpdateVC = this.UpdateVC.bind(this);
    this.RemoveVC = this.RemoveVC.bind(this);
  }
  CreateDid(request: MsgCreateDid): Promise<MsgCreateDidResponse> {
    const data = MsgCreateDid.encode(request).finish();
    const promise = this.rpc.request("metaearth.did.Msg", "CreateDid", data);
    return promise.then((data) => MsgCreateDidResponse.decode(new _m0.Reader(data)));
  }

  UpdateDidStatus(request: MsgUpdateDidStatus): Promise<MsgUpdateDidStatusResponse> {
    const data = MsgUpdateDidStatus.encode(request).finish();
    const promise = this.rpc.request("metaearth.did.Msg", "UpdateDidStatus", data);
    return promise.then((data) => MsgUpdateDidStatusResponse.decode(new _m0.Reader(data)));
  }

  CreateService(request: MsgCreateService): Promise<MsgCreateServiceResponse> {
    const data = MsgCreateService.encode(request).finish();
    const promise = this.rpc.request("metaearth.did.Msg", "CreateService", data);
    return promise.then((data) => MsgCreateServiceResponse.decode(new _m0.Reader(data)));
  }

  UpdateServiceStatus(request: MsgUpdateServiceStatus): Promise<MsgUpdateServiceStatusResponse> {
    const data = MsgUpdateServiceStatus.encode(request).finish();
    const promise = this.rpc.request("metaearth.did.Msg", "UpdateServiceStatus", data);
    return promise.then((data) => MsgUpdateServiceStatusResponse.decode(new _m0.Reader(data)));
  }

  CreateVC(request: MsgCreateVC): Promise<MsgCreateVCResponse> {
    const data = MsgCreateVC.encode(request).finish();
    const promise = this.rpc.request("metaearth.did.Msg", "CreateVC", data);
    return promise.then((data) => MsgCreateVCResponse.decode(new _m0.Reader(data)));
  }

  UpdateVC(request: MsgUpdateVC): Promise<MsgUpdateVCResponse> {
    const data = MsgUpdateVC.encode(request).finish();
    const promise = this.rpc.request("metaearth.did.Msg", "UpdateVC", data);
    return promise.then((data) => MsgUpdateVCResponse.decode(new _m0.Reader(data)));
  }

  RemoveVC(request: MsgRemoveVC): Promise<MsgRemoveVCResponse> {
    const data = MsgRemoveVC.encode(request).finish();
    const promise = this.rpc.request("metaearth.did.Msg", "RemoveVC", data);
    return promise.then((data) => MsgRemoveVCResponse.decode(new _m0.Reader(data)));
  }
}

interface Rpc {
  request(service: string, method: string, data: Uint8Array): Promise<Uint8Array>;
}

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

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}
