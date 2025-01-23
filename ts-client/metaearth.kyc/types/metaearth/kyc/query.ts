/* eslint-disable */
import _m0 from "protobufjs/minimal";
import { PageRequest, PageResponse } from "../../cosmos/base/query/v1beta1/pagination";
import { NFT } from "../../cosmos/nft/v1beta1/nft";
import { Credential } from "../did/credential";
import { DidInfo } from "../did/did";
import { Protocol } from "./protocol";

export const protobufPackage = "metaearth.kyc";

export interface QueryProtocol {
}

export interface QueryProtocolResponse {
  protocol: Protocol | undefined;
}

export interface QueryDID {
  address: string;
}

export interface QueryDIDResponse {
  info: DidInfo | undefined;
}

export interface QueryDIDs {
  regionId: string;
  pagination: PageRequest | undefined;
}

export interface QueryDIDsResponse {
  infos: DidInfo[];
  pagination: PageResponse | undefined;
}

export interface QueryKYC {
  did: string;
}

export interface QueryKYCResponse {
  kyc: Credential | undefined;
}

export interface QueryKYCs {
  regionId: string;
  pagination: PageRequest | undefined;
}

export interface QueryKYCsResponse {
  KYCs: Credential[];
  pagination: PageResponse | undefined;
}

export interface QuerySBT {
  did: string;
}

export interface QuerySBTResponse {
  sbt: NFT | undefined;
}

function createBaseQueryProtocol(): QueryProtocol {
  return {};
}

export const QueryProtocol = {
  encode(_: QueryProtocol, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryProtocol {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryProtocol();
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

  fromJSON(_: any): QueryProtocol {
    return {};
  },

  toJSON(_: QueryProtocol): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryProtocol>, I>>(_: I): QueryProtocol {
    const message = createBaseQueryProtocol();
    return message;
  },
};

function createBaseQueryProtocolResponse(): QueryProtocolResponse {
  return { protocol: undefined };
}

export const QueryProtocolResponse = {
  encode(message: QueryProtocolResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.protocol !== undefined) {
      Protocol.encode(message.protocol, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryProtocolResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryProtocolResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.protocol = Protocol.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryProtocolResponse {
    return { protocol: isSet(object.protocol) ? Protocol.fromJSON(object.protocol) : undefined };
  },

  toJSON(message: QueryProtocolResponse): unknown {
    const obj: any = {};
    message.protocol !== undefined && (obj.protocol = message.protocol ? Protocol.toJSON(message.protocol) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryProtocolResponse>, I>>(object: I): QueryProtocolResponse {
    const message = createBaseQueryProtocolResponse();
    message.protocol = (object.protocol !== undefined && object.protocol !== null)
      ? Protocol.fromPartial(object.protocol)
      : undefined;
    return message;
  },
};

function createBaseQueryDID(): QueryDID {
  return { address: "" };
}

export const QueryDID = {
  encode(message: QueryDID, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.address !== "") {
      writer.uint32(10).string(message.address);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryDID {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryDID();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.address = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryDID {
    return { address: isSet(object.address) ? String(object.address) : "" };
  },

  toJSON(message: QueryDID): unknown {
    const obj: any = {};
    message.address !== undefined && (obj.address = message.address);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryDID>, I>>(object: I): QueryDID {
    const message = createBaseQueryDID();
    message.address = object.address ?? "";
    return message;
  },
};

function createBaseQueryDIDResponse(): QueryDIDResponse {
  return { info: undefined };
}

export const QueryDIDResponse = {
  encode(message: QueryDIDResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.info !== undefined) {
      DidInfo.encode(message.info, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryDIDResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryDIDResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.info = DidInfo.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryDIDResponse {
    return { info: isSet(object.info) ? DidInfo.fromJSON(object.info) : undefined };
  },

  toJSON(message: QueryDIDResponse): unknown {
    const obj: any = {};
    message.info !== undefined && (obj.info = message.info ? DidInfo.toJSON(message.info) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryDIDResponse>, I>>(object: I): QueryDIDResponse {
    const message = createBaseQueryDIDResponse();
    message.info = (object.info !== undefined && object.info !== null) ? DidInfo.fromPartial(object.info) : undefined;
    return message;
  },
};

function createBaseQueryDIDs(): QueryDIDs {
  return { regionId: "", pagination: undefined };
}

export const QueryDIDs = {
  encode(message: QueryDIDs, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.regionId !== "") {
      writer.uint32(10).string(message.regionId);
    }
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryDIDs {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryDIDs();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.regionId = reader.string();
          break;
        case 2:
          message.pagination = PageRequest.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryDIDs {
    return {
      regionId: isSet(object.regionId) ? String(object.regionId) : "",
      pagination: isSet(object.pagination) ? PageRequest.fromJSON(object.pagination) : undefined,
    };
  },

  toJSON(message: QueryDIDs): unknown {
    const obj: any = {};
    message.regionId !== undefined && (obj.regionId = message.regionId);
    message.pagination !== undefined
      && (obj.pagination = message.pagination ? PageRequest.toJSON(message.pagination) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryDIDs>, I>>(object: I): QueryDIDs {
    const message = createBaseQueryDIDs();
    message.regionId = object.regionId ?? "";
    message.pagination = (object.pagination !== undefined && object.pagination !== null)
      ? PageRequest.fromPartial(object.pagination)
      : undefined;
    return message;
  },
};

function createBaseQueryDIDsResponse(): QueryDIDsResponse {
  return { infos: [], pagination: undefined };
}

export const QueryDIDsResponse = {
  encode(message: QueryDIDsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.infos) {
      DidInfo.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.pagination !== undefined) {
      PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryDIDsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryDIDsResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.infos.push(DidInfo.decode(reader, reader.uint32()));
          break;
        case 2:
          message.pagination = PageResponse.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryDIDsResponse {
    return {
      infos: Array.isArray(object?.infos) ? object.infos.map((e: any) => DidInfo.fromJSON(e)) : [],
      pagination: isSet(object.pagination) ? PageResponse.fromJSON(object.pagination) : undefined,
    };
  },

  toJSON(message: QueryDIDsResponse): unknown {
    const obj: any = {};
    if (message.infos) {
      obj.infos = message.infos.map((e) => e ? DidInfo.toJSON(e) : undefined);
    } else {
      obj.infos = [];
    }
    message.pagination !== undefined
      && (obj.pagination = message.pagination ? PageResponse.toJSON(message.pagination) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryDIDsResponse>, I>>(object: I): QueryDIDsResponse {
    const message = createBaseQueryDIDsResponse();
    message.infos = object.infos?.map((e) => DidInfo.fromPartial(e)) || [];
    message.pagination = (object.pagination !== undefined && object.pagination !== null)
      ? PageResponse.fromPartial(object.pagination)
      : undefined;
    return message;
  },
};

function createBaseQueryKYC(): QueryKYC {
  return { did: "" };
}

export const QueryKYC = {
  encode(message: QueryKYC, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.did !== "") {
      writer.uint32(10).string(message.did);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryKYC {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryKYC();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.did = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryKYC {
    return { did: isSet(object.did) ? String(object.did) : "" };
  },

  toJSON(message: QueryKYC): unknown {
    const obj: any = {};
    message.did !== undefined && (obj.did = message.did);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryKYC>, I>>(object: I): QueryKYC {
    const message = createBaseQueryKYC();
    message.did = object.did ?? "";
    return message;
  },
};

function createBaseQueryKYCResponse(): QueryKYCResponse {
  return { kyc: undefined };
}

export const QueryKYCResponse = {
  encode(message: QueryKYCResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.kyc !== undefined) {
      Credential.encode(message.kyc, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryKYCResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryKYCResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.kyc = Credential.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryKYCResponse {
    return { kyc: isSet(object.kyc) ? Credential.fromJSON(object.kyc) : undefined };
  },

  toJSON(message: QueryKYCResponse): unknown {
    const obj: any = {};
    message.kyc !== undefined && (obj.kyc = message.kyc ? Credential.toJSON(message.kyc) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryKYCResponse>, I>>(object: I): QueryKYCResponse {
    const message = createBaseQueryKYCResponse();
    message.kyc = (object.kyc !== undefined && object.kyc !== null) ? Credential.fromPartial(object.kyc) : undefined;
    return message;
  },
};

function createBaseQueryKYCs(): QueryKYCs {
  return { regionId: "", pagination: undefined };
}

export const QueryKYCs = {
  encode(message: QueryKYCs, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.regionId !== "") {
      writer.uint32(10).string(message.regionId);
    }
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryKYCs {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryKYCs();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.regionId = reader.string();
          break;
        case 2:
          message.pagination = PageRequest.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryKYCs {
    return {
      regionId: isSet(object.regionId) ? String(object.regionId) : "",
      pagination: isSet(object.pagination) ? PageRequest.fromJSON(object.pagination) : undefined,
    };
  },

  toJSON(message: QueryKYCs): unknown {
    const obj: any = {};
    message.regionId !== undefined && (obj.regionId = message.regionId);
    message.pagination !== undefined
      && (obj.pagination = message.pagination ? PageRequest.toJSON(message.pagination) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryKYCs>, I>>(object: I): QueryKYCs {
    const message = createBaseQueryKYCs();
    message.regionId = object.regionId ?? "";
    message.pagination = (object.pagination !== undefined && object.pagination !== null)
      ? PageRequest.fromPartial(object.pagination)
      : undefined;
    return message;
  },
};

function createBaseQueryKYCsResponse(): QueryKYCsResponse {
  return { KYCs: [], pagination: undefined };
}

export const QueryKYCsResponse = {
  encode(message: QueryKYCsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.KYCs) {
      Credential.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.pagination !== undefined) {
      PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryKYCsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryKYCsResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.KYCs.push(Credential.decode(reader, reader.uint32()));
          break;
        case 2:
          message.pagination = PageResponse.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryKYCsResponse {
    return {
      KYCs: Array.isArray(object?.KYCs) ? object.KYCs.map((e: any) => Credential.fromJSON(e)) : [],
      pagination: isSet(object.pagination) ? PageResponse.fromJSON(object.pagination) : undefined,
    };
  },

  toJSON(message: QueryKYCsResponse): unknown {
    const obj: any = {};
    if (message.KYCs) {
      obj.KYCs = message.KYCs.map((e) => e ? Credential.toJSON(e) : undefined);
    } else {
      obj.KYCs = [];
    }
    message.pagination !== undefined
      && (obj.pagination = message.pagination ? PageResponse.toJSON(message.pagination) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryKYCsResponse>, I>>(object: I): QueryKYCsResponse {
    const message = createBaseQueryKYCsResponse();
    message.KYCs = object.KYCs?.map((e) => Credential.fromPartial(e)) || [];
    message.pagination = (object.pagination !== undefined && object.pagination !== null)
      ? PageResponse.fromPartial(object.pagination)
      : undefined;
    return message;
  },
};

function createBaseQuerySBT(): QuerySBT {
  return { did: "" };
}

export const QuerySBT = {
  encode(message: QuerySBT, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.did !== "") {
      writer.uint32(10).string(message.did);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QuerySBT {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQuerySBT();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.did = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QuerySBT {
    return { did: isSet(object.did) ? String(object.did) : "" };
  },

  toJSON(message: QuerySBT): unknown {
    const obj: any = {};
    message.did !== undefined && (obj.did = message.did);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QuerySBT>, I>>(object: I): QuerySBT {
    const message = createBaseQuerySBT();
    message.did = object.did ?? "";
    return message;
  },
};

function createBaseQuerySBTResponse(): QuerySBTResponse {
  return { sbt: undefined };
}

export const QuerySBTResponse = {
  encode(message: QuerySBTResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.sbt !== undefined) {
      NFT.encode(message.sbt, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QuerySBTResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQuerySBTResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.sbt = NFT.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QuerySBTResponse {
    return { sbt: isSet(object.sbt) ? NFT.fromJSON(object.sbt) : undefined };
  },

  toJSON(message: QuerySBTResponse): unknown {
    const obj: any = {};
    message.sbt !== undefined && (obj.sbt = message.sbt ? NFT.toJSON(message.sbt) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QuerySBTResponse>, I>>(object: I): QuerySBTResponse {
    const message = createBaseQuerySBTResponse();
    message.sbt = (object.sbt !== undefined && object.sbt !== null) ? NFT.fromPartial(object.sbt) : undefined;
    return message;
  },
};

export interface Query {
  /** query protocol */
  Protocol(request: QueryProtocol): Promise<QueryProtocolResponse>;
  DID(request: QueryDID): Promise<QueryDIDResponse>;
  DIDs(request: QueryDIDs): Promise<QueryDIDsResponse>;
  KYC(request: QueryKYC): Promise<QueryKYCResponse>;
  KYCs(request: QueryKYCs): Promise<QueryKYCsResponse>;
  SBT(request: QuerySBT): Promise<QuerySBTResponse>;
}

export class QueryClientImpl implements Query {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.Protocol = this.Protocol.bind(this);
    this.DID = this.DID.bind(this);
    this.DIDs = this.DIDs.bind(this);
    this.KYC = this.KYC.bind(this);
    this.KYCs = this.KYCs.bind(this);
    this.SBT = this.SBT.bind(this);
  }
  Protocol(request: QueryProtocol): Promise<QueryProtocolResponse> {
    const data = QueryProtocol.encode(request).finish();
    const promise = this.rpc.request("metaearth.kyc.Query", "Protocol", data);
    return promise.then((data) => QueryProtocolResponse.decode(new _m0.Reader(data)));
  }

  DID(request: QueryDID): Promise<QueryDIDResponse> {
    const data = QueryDID.encode(request).finish();
    const promise = this.rpc.request("metaearth.kyc.Query", "DID", data);
    return promise.then((data) => QueryDIDResponse.decode(new _m0.Reader(data)));
  }

  DIDs(request: QueryDIDs): Promise<QueryDIDsResponse> {
    const data = QueryDIDs.encode(request).finish();
    const promise = this.rpc.request("metaearth.kyc.Query", "DIDs", data);
    return promise.then((data) => QueryDIDsResponse.decode(new _m0.Reader(data)));
  }

  KYC(request: QueryKYC): Promise<QueryKYCResponse> {
    const data = QueryKYC.encode(request).finish();
    const promise = this.rpc.request("metaearth.kyc.Query", "KYC", data);
    return promise.then((data) => QueryKYCResponse.decode(new _m0.Reader(data)));
  }

  KYCs(request: QueryKYCs): Promise<QueryKYCsResponse> {
    const data = QueryKYCs.encode(request).finish();
    const promise = this.rpc.request("metaearth.kyc.Query", "KYCs", data);
    return promise.then((data) => QueryKYCsResponse.decode(new _m0.Reader(data)));
  }

  SBT(request: QuerySBT): Promise<QuerySBTResponse> {
    const data = QuerySBT.encode(request).finish();
    const promise = this.rpc.request("metaearth.kyc.Query", "SBT", data);
    return promise.then((data) => QuerySBTResponse.decode(new _m0.Reader(data)));
  }
}

interface Rpc {
  request(service: string, method: string, data: Uint8Array): Promise<Uint8Array>;
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
