/* eslint-disable */
import _m0 from "protobufjs/minimal";
import { KycLevel, kycLevelFromJSON, kycLevelToJSON } from "../did/did";

export const protobufPackage = "metaearth.kyc";

export interface MsgApprove {
  issuer: string;
  /** me-id */
  did: string;
  regionId: string;
  address: string;
  pubkey: string;
  level: KycLevel;
  uri: string;
  hash: string;
  inviter: string;
}

export interface MsgApproveResponse {
}

export interface MsgUpdate {
  issuer: string;
  did: string;
  regionId: string;
  level: KycLevel;
  uri: string;
  hash: string;
  inviter: string;
}

export interface MsgUpdateResponse {
}

export interface MsgRemove {
  issuer: string;
  did: string;
}

export interface MsgRemoveResponse {
}

export interface MsgCreateSBT {
  issuer: string;
  did: string;
  uri: string;
  uriHash: string;
  data: Uint8Array;
}

export interface MsgCreateSBTResponse {
}

export interface MsgUpdateSBT {
  issuer: string;
  did: string;
  uri: string;
  uriHash: string;
  data: Uint8Array;
}

export interface MsgUpdateSBTResponse {
}

export interface MsgDeleteSBT {
  issuer: string;
  did: string;
}

export interface MsgDeleteSBTResponse {
}

function createBaseMsgApprove(): MsgApprove {
  return { issuer: "", did: "", regionId: "", address: "", pubkey: "", level: 0, uri: "", hash: "", inviter: "" };
}

export const MsgApprove = {
  encode(message: MsgApprove, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.issuer !== "") {
      writer.uint32(10).string(message.issuer);
    }
    if (message.did !== "") {
      writer.uint32(18).string(message.did);
    }
    if (message.regionId !== "") {
      writer.uint32(26).string(message.regionId);
    }
    if (message.address !== "") {
      writer.uint32(34).string(message.address);
    }
    if (message.pubkey !== "") {
      writer.uint32(42).string(message.pubkey);
    }
    if (message.level !== 0) {
      writer.uint32(48).int32(message.level);
    }
    if (message.uri !== "") {
      writer.uint32(58).string(message.uri);
    }
    if (message.hash !== "") {
      writer.uint32(66).string(message.hash);
    }
    if (message.inviter !== "") {
      writer.uint32(74).string(message.inviter);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgApprove {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgApprove();
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
          message.regionId = reader.string();
          break;
        case 4:
          message.address = reader.string();
          break;
        case 5:
          message.pubkey = reader.string();
          break;
        case 6:
          message.level = reader.int32() as any;
          break;
        case 7:
          message.uri = reader.string();
          break;
        case 8:
          message.hash = reader.string();
          break;
        case 9:
          message.inviter = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgApprove {
    return {
      issuer: isSet(object.issuer) ? String(object.issuer) : "",
      did: isSet(object.did) ? String(object.did) : "",
      regionId: isSet(object.regionId) ? String(object.regionId) : "",
      address: isSet(object.address) ? String(object.address) : "",
      pubkey: isSet(object.pubkey) ? String(object.pubkey) : "",
      level: isSet(object.level) ? kycLevelFromJSON(object.level) : 0,
      uri: isSet(object.uri) ? String(object.uri) : "",
      hash: isSet(object.hash) ? String(object.hash) : "",
      inviter: isSet(object.inviter) ? String(object.inviter) : "",
    };
  },

  toJSON(message: MsgApprove): unknown {
    const obj: any = {};
    message.issuer !== undefined && (obj.issuer = message.issuer);
    message.did !== undefined && (obj.did = message.did);
    message.regionId !== undefined && (obj.regionId = message.regionId);
    message.address !== undefined && (obj.address = message.address);
    message.pubkey !== undefined && (obj.pubkey = message.pubkey);
    message.level !== undefined && (obj.level = kycLevelToJSON(message.level));
    message.uri !== undefined && (obj.uri = message.uri);
    message.hash !== undefined && (obj.hash = message.hash);
    message.inviter !== undefined && (obj.inviter = message.inviter);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgApprove>, I>>(object: I): MsgApprove {
    const message = createBaseMsgApprove();
    message.issuer = object.issuer ?? "";
    message.did = object.did ?? "";
    message.regionId = object.regionId ?? "";
    message.address = object.address ?? "";
    message.pubkey = object.pubkey ?? "";
    message.level = object.level ?? 0;
    message.uri = object.uri ?? "";
    message.hash = object.hash ?? "";
    message.inviter = object.inviter ?? "";
    return message;
  },
};

function createBaseMsgApproveResponse(): MsgApproveResponse {
  return {};
}

export const MsgApproveResponse = {
  encode(_: MsgApproveResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgApproveResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgApproveResponse();
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

  fromJSON(_: any): MsgApproveResponse {
    return {};
  },

  toJSON(_: MsgApproveResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgApproveResponse>, I>>(_: I): MsgApproveResponse {
    const message = createBaseMsgApproveResponse();
    return message;
  },
};

function createBaseMsgUpdate(): MsgUpdate {
  return { issuer: "", did: "", regionId: "", level: 0, uri: "", hash: "", inviter: "" };
}

export const MsgUpdate = {
  encode(message: MsgUpdate, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.issuer !== "") {
      writer.uint32(10).string(message.issuer);
    }
    if (message.did !== "") {
      writer.uint32(18).string(message.did);
    }
    if (message.regionId !== "") {
      writer.uint32(26).string(message.regionId);
    }
    if (message.level !== 0) {
      writer.uint32(32).int32(message.level);
    }
    if (message.uri !== "") {
      writer.uint32(42).string(message.uri);
    }
    if (message.hash !== "") {
      writer.uint32(50).string(message.hash);
    }
    if (message.inviter !== "") {
      writer.uint32(58).string(message.inviter);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdate {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdate();
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
          message.regionId = reader.string();
          break;
        case 4:
          message.level = reader.int32() as any;
          break;
        case 5:
          message.uri = reader.string();
          break;
        case 6:
          message.hash = reader.string();
          break;
        case 7:
          message.inviter = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgUpdate {
    return {
      issuer: isSet(object.issuer) ? String(object.issuer) : "",
      did: isSet(object.did) ? String(object.did) : "",
      regionId: isSet(object.regionId) ? String(object.regionId) : "",
      level: isSet(object.level) ? kycLevelFromJSON(object.level) : 0,
      uri: isSet(object.uri) ? String(object.uri) : "",
      hash: isSet(object.hash) ? String(object.hash) : "",
      inviter: isSet(object.inviter) ? String(object.inviter) : "",
    };
  },

  toJSON(message: MsgUpdate): unknown {
    const obj: any = {};
    message.issuer !== undefined && (obj.issuer = message.issuer);
    message.did !== undefined && (obj.did = message.did);
    message.regionId !== undefined && (obj.regionId = message.regionId);
    message.level !== undefined && (obj.level = kycLevelToJSON(message.level));
    message.uri !== undefined && (obj.uri = message.uri);
    message.hash !== undefined && (obj.hash = message.hash);
    message.inviter !== undefined && (obj.inviter = message.inviter);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgUpdate>, I>>(object: I): MsgUpdate {
    const message = createBaseMsgUpdate();
    message.issuer = object.issuer ?? "";
    message.did = object.did ?? "";
    message.regionId = object.regionId ?? "";
    message.level = object.level ?? 0;
    message.uri = object.uri ?? "";
    message.hash = object.hash ?? "";
    message.inviter = object.inviter ?? "";
    return message;
  },
};

function createBaseMsgUpdateResponse(): MsgUpdateResponse {
  return {};
}

export const MsgUpdateResponse = {
  encode(_: MsgUpdateResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateResponse();
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

  fromJSON(_: any): MsgUpdateResponse {
    return {};
  },

  toJSON(_: MsgUpdateResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgUpdateResponse>, I>>(_: I): MsgUpdateResponse {
    const message = createBaseMsgUpdateResponse();
    return message;
  },
};

function createBaseMsgRemove(): MsgRemove {
  return { issuer: "", did: "" };
}

export const MsgRemove = {
  encode(message: MsgRemove, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.issuer !== "") {
      writer.uint32(10).string(message.issuer);
    }
    if (message.did !== "") {
      writer.uint32(18).string(message.did);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgRemove {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgRemove();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.issuer = reader.string();
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

  fromJSON(object: any): MsgRemove {
    return {
      issuer: isSet(object.issuer) ? String(object.issuer) : "",
      did: isSet(object.did) ? String(object.did) : "",
    };
  },

  toJSON(message: MsgRemove): unknown {
    const obj: any = {};
    message.issuer !== undefined && (obj.issuer = message.issuer);
    message.did !== undefined && (obj.did = message.did);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgRemove>, I>>(object: I): MsgRemove {
    const message = createBaseMsgRemove();
    message.issuer = object.issuer ?? "";
    message.did = object.did ?? "";
    return message;
  },
};

function createBaseMsgRemoveResponse(): MsgRemoveResponse {
  return {};
}

export const MsgRemoveResponse = {
  encode(_: MsgRemoveResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgRemoveResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgRemoveResponse();
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

  fromJSON(_: any): MsgRemoveResponse {
    return {};
  },

  toJSON(_: MsgRemoveResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgRemoveResponse>, I>>(_: I): MsgRemoveResponse {
    const message = createBaseMsgRemoveResponse();
    return message;
  },
};

function createBaseMsgCreateSBT(): MsgCreateSBT {
  return { issuer: "", did: "", uri: "", uriHash: "", data: new Uint8Array() };
}

export const MsgCreateSBT = {
  encode(message: MsgCreateSBT, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.issuer !== "") {
      writer.uint32(10).string(message.issuer);
    }
    if (message.did !== "") {
      writer.uint32(18).string(message.did);
    }
    if (message.uri !== "") {
      writer.uint32(26).string(message.uri);
    }
    if (message.uriHash !== "") {
      writer.uint32(34).string(message.uriHash);
    }
    if (message.data.length !== 0) {
      writer.uint32(42).bytes(message.data);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgCreateSBT {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgCreateSBT();
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
          message.uri = reader.string();
          break;
        case 4:
          message.uriHash = reader.string();
          break;
        case 5:
          message.data = reader.bytes();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgCreateSBT {
    return {
      issuer: isSet(object.issuer) ? String(object.issuer) : "",
      did: isSet(object.did) ? String(object.did) : "",
      uri: isSet(object.uri) ? String(object.uri) : "",
      uriHash: isSet(object.uriHash) ? String(object.uriHash) : "",
      data: isSet(object.data) ? bytesFromBase64(object.data) : new Uint8Array(),
    };
  },

  toJSON(message: MsgCreateSBT): unknown {
    const obj: any = {};
    message.issuer !== undefined && (obj.issuer = message.issuer);
    message.did !== undefined && (obj.did = message.did);
    message.uri !== undefined && (obj.uri = message.uri);
    message.uriHash !== undefined && (obj.uriHash = message.uriHash);
    message.data !== undefined
      && (obj.data = base64FromBytes(message.data !== undefined ? message.data : new Uint8Array()));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgCreateSBT>, I>>(object: I): MsgCreateSBT {
    const message = createBaseMsgCreateSBT();
    message.issuer = object.issuer ?? "";
    message.did = object.did ?? "";
    message.uri = object.uri ?? "";
    message.uriHash = object.uriHash ?? "";
    message.data = object.data ?? new Uint8Array();
    return message;
  },
};

function createBaseMsgCreateSBTResponse(): MsgCreateSBTResponse {
  return {};
}

export const MsgCreateSBTResponse = {
  encode(_: MsgCreateSBTResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgCreateSBTResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgCreateSBTResponse();
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

  fromJSON(_: any): MsgCreateSBTResponse {
    return {};
  },

  toJSON(_: MsgCreateSBTResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgCreateSBTResponse>, I>>(_: I): MsgCreateSBTResponse {
    const message = createBaseMsgCreateSBTResponse();
    return message;
  },
};

function createBaseMsgUpdateSBT(): MsgUpdateSBT {
  return { issuer: "", did: "", uri: "", uriHash: "", data: new Uint8Array() };
}

export const MsgUpdateSBT = {
  encode(message: MsgUpdateSBT, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.issuer !== "") {
      writer.uint32(10).string(message.issuer);
    }
    if (message.did !== "") {
      writer.uint32(18).string(message.did);
    }
    if (message.uri !== "") {
      writer.uint32(26).string(message.uri);
    }
    if (message.uriHash !== "") {
      writer.uint32(34).string(message.uriHash);
    }
    if (message.data.length !== 0) {
      writer.uint32(42).bytes(message.data);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateSBT {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateSBT();
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
          message.uri = reader.string();
          break;
        case 4:
          message.uriHash = reader.string();
          break;
        case 5:
          message.data = reader.bytes();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgUpdateSBT {
    return {
      issuer: isSet(object.issuer) ? String(object.issuer) : "",
      did: isSet(object.did) ? String(object.did) : "",
      uri: isSet(object.uri) ? String(object.uri) : "",
      uriHash: isSet(object.uriHash) ? String(object.uriHash) : "",
      data: isSet(object.data) ? bytesFromBase64(object.data) : new Uint8Array(),
    };
  },

  toJSON(message: MsgUpdateSBT): unknown {
    const obj: any = {};
    message.issuer !== undefined && (obj.issuer = message.issuer);
    message.did !== undefined && (obj.did = message.did);
    message.uri !== undefined && (obj.uri = message.uri);
    message.uriHash !== undefined && (obj.uriHash = message.uriHash);
    message.data !== undefined
      && (obj.data = base64FromBytes(message.data !== undefined ? message.data : new Uint8Array()));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgUpdateSBT>, I>>(object: I): MsgUpdateSBT {
    const message = createBaseMsgUpdateSBT();
    message.issuer = object.issuer ?? "";
    message.did = object.did ?? "";
    message.uri = object.uri ?? "";
    message.uriHash = object.uriHash ?? "";
    message.data = object.data ?? new Uint8Array();
    return message;
  },
};

function createBaseMsgUpdateSBTResponse(): MsgUpdateSBTResponse {
  return {};
}

export const MsgUpdateSBTResponse = {
  encode(_: MsgUpdateSBTResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateSBTResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateSBTResponse();
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

  fromJSON(_: any): MsgUpdateSBTResponse {
    return {};
  },

  toJSON(_: MsgUpdateSBTResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgUpdateSBTResponse>, I>>(_: I): MsgUpdateSBTResponse {
    const message = createBaseMsgUpdateSBTResponse();
    return message;
  },
};

function createBaseMsgDeleteSBT(): MsgDeleteSBT {
  return { issuer: "", did: "" };
}

export const MsgDeleteSBT = {
  encode(message: MsgDeleteSBT, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.issuer !== "") {
      writer.uint32(10).string(message.issuer);
    }
    if (message.did !== "") {
      writer.uint32(18).string(message.did);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgDeleteSBT {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgDeleteSBT();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.issuer = reader.string();
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

  fromJSON(object: any): MsgDeleteSBT {
    return {
      issuer: isSet(object.issuer) ? String(object.issuer) : "",
      did: isSet(object.did) ? String(object.did) : "",
    };
  },

  toJSON(message: MsgDeleteSBT): unknown {
    const obj: any = {};
    message.issuer !== undefined && (obj.issuer = message.issuer);
    message.did !== undefined && (obj.did = message.did);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgDeleteSBT>, I>>(object: I): MsgDeleteSBT {
    const message = createBaseMsgDeleteSBT();
    message.issuer = object.issuer ?? "";
    message.did = object.did ?? "";
    return message;
  },
};

function createBaseMsgDeleteSBTResponse(): MsgDeleteSBTResponse {
  return {};
}

export const MsgDeleteSBTResponse = {
  encode(_: MsgDeleteSBTResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgDeleteSBTResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgDeleteSBTResponse();
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

  fromJSON(_: any): MsgDeleteSBTResponse {
    return {};
  },

  toJSON(_: MsgDeleteSBTResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgDeleteSBTResponse>, I>>(_: I): MsgDeleteSBTResponse {
    const message = createBaseMsgDeleteSBTResponse();
    return message;
  },
};

/** Msg defines the Msg service. */
export interface Msg {
  Approve(request: MsgApprove): Promise<MsgApproveResponse>;
  Update(request: MsgUpdate): Promise<MsgUpdateResponse>;
  Remove(request: MsgRemove): Promise<MsgRemoveResponse>;
  CreateSBT(request: MsgCreateSBT): Promise<MsgCreateSBTResponse>;
  UpdateSBT(request: MsgUpdateSBT): Promise<MsgUpdateSBTResponse>;
  DeleteSBT(request: MsgDeleteSBT): Promise<MsgDeleteSBTResponse>;
}

export class MsgClientImpl implements Msg {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.Approve = this.Approve.bind(this);
    this.Update = this.Update.bind(this);
    this.Remove = this.Remove.bind(this);
    this.CreateSBT = this.CreateSBT.bind(this);
    this.UpdateSBT = this.UpdateSBT.bind(this);
    this.DeleteSBT = this.DeleteSBT.bind(this);
  }
  Approve(request: MsgApprove): Promise<MsgApproveResponse> {
    const data = MsgApprove.encode(request).finish();
    const promise = this.rpc.request("metaearth.kyc.Msg", "Approve", data);
    return promise.then((data) => MsgApproveResponse.decode(new _m0.Reader(data)));
  }

  Update(request: MsgUpdate): Promise<MsgUpdateResponse> {
    const data = MsgUpdate.encode(request).finish();
    const promise = this.rpc.request("metaearth.kyc.Msg", "Update", data);
    return promise.then((data) => MsgUpdateResponse.decode(new _m0.Reader(data)));
  }

  Remove(request: MsgRemove): Promise<MsgRemoveResponse> {
    const data = MsgRemove.encode(request).finish();
    const promise = this.rpc.request("metaearth.kyc.Msg", "Remove", data);
    return promise.then((data) => MsgRemoveResponse.decode(new _m0.Reader(data)));
  }

  CreateSBT(request: MsgCreateSBT): Promise<MsgCreateSBTResponse> {
    const data = MsgCreateSBT.encode(request).finish();
    const promise = this.rpc.request("metaearth.kyc.Msg", "CreateSBT", data);
    return promise.then((data) => MsgCreateSBTResponse.decode(new _m0.Reader(data)));
  }

  UpdateSBT(request: MsgUpdateSBT): Promise<MsgUpdateSBTResponse> {
    const data = MsgUpdateSBT.encode(request).finish();
    const promise = this.rpc.request("metaearth.kyc.Msg", "UpdateSBT", data);
    return promise.then((data) => MsgUpdateSBTResponse.decode(new _m0.Reader(data)));
  }

  DeleteSBT(request: MsgDeleteSBT): Promise<MsgDeleteSBTResponse> {
    const data = MsgDeleteSBT.encode(request).finish();
    const promise = this.rpc.request("metaearth.kyc.Msg", "DeleteSBT", data);
    return promise.then((data) => MsgDeleteSBTResponse.decode(new _m0.Reader(data)));
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
