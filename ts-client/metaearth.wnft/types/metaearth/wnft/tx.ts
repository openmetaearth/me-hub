/* eslint-disable */
import Long from "long";
import _m0 from "protobufjs/minimal";

export const protobufPackage = "metaearth.wnft";

export interface MsgNewClass {
  /**
   * class_id defines the unique identifier of the nft classification, similar
   * to the contract address of ERC721
   */
  classId: string;
  totalSupply: number;
  /** sender is the address of the owner of nft */
  sender: string;
  name: string;
  symbol: string;
  /** description is a brief description of nft classification. Optional */
  description: string;
  /**
   * uri for the class metadata stored off chain. It can define schema for Class
   * and NFT `Data` attributes. Optional
   */
  uri: string;
  /** uri_hash is a hash of the document pointed by uri. Optional */
  uriHash: string;
}

/** MsgNewClassResponse defines the Msg/NewClass response type. */
export interface MsgNewClassResponse {
}

export interface MsgMintNFT {
  /**
   * class_id defines the unique identifier of the nft classification, similar
   * to the contract address of ERC721
   */
  classId: string;
  tokenId: string;
  uri: string;
  /** uri_hash is a hash of the document pointed by uri */
  uriHash: string;
  /** sender is the address of the owner of nft */
  creator: string;
  receiver: string;
}

/** MsgMintResponse defines the Msg/Mint response type. */
export interface MsgMintNFTResponse {
}

/**
 * MsgSend represents a message to send a nft from one account to another
 * account.
 */
export interface MsgSend {
  /**
   * class_id defines the unique identifier of the nft classification, similar
   * to the contract address of ERC721
   */
  classId: string;
  /** id defines the unique identification of nft */
  id: string;
  /** sender is the address of the owner of nft */
  sender: string;
  /** receiver is the receiver address of nft */
  receiver: string;
}

/** MsgSendResponse defines the Msg/Send response type. */
export interface MsgSendResponse {
}

function createBaseMsgNewClass(): MsgNewClass {
  return { classId: "", totalSupply: 0, sender: "", name: "", symbol: "", description: "", uri: "", uriHash: "" };
}

export const MsgNewClass = {
  encode(message: MsgNewClass, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.classId !== "") {
      writer.uint32(10).string(message.classId);
    }
    if (message.totalSupply !== 0) {
      writer.uint32(16).uint64(message.totalSupply);
    }
    if (message.sender !== "") {
      writer.uint32(26).string(message.sender);
    }
    if (message.name !== "") {
      writer.uint32(34).string(message.name);
    }
    if (message.symbol !== "") {
      writer.uint32(42).string(message.symbol);
    }
    if (message.description !== "") {
      writer.uint32(50).string(message.description);
    }
    if (message.uri !== "") {
      writer.uint32(58).string(message.uri);
    }
    if (message.uriHash !== "") {
      writer.uint32(66).string(message.uriHash);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgNewClass {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgNewClass();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.classId = reader.string();
          break;
        case 2:
          message.totalSupply = longToNumber(reader.uint64() as Long);
          break;
        case 3:
          message.sender = reader.string();
          break;
        case 4:
          message.name = reader.string();
          break;
        case 5:
          message.symbol = reader.string();
          break;
        case 6:
          message.description = reader.string();
          break;
        case 7:
          message.uri = reader.string();
          break;
        case 8:
          message.uriHash = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgNewClass {
    return {
      classId: isSet(object.classId) ? String(object.classId) : "",
      totalSupply: isSet(object.totalSupply) ? Number(object.totalSupply) : 0,
      sender: isSet(object.sender) ? String(object.sender) : "",
      name: isSet(object.name) ? String(object.name) : "",
      symbol: isSet(object.symbol) ? String(object.symbol) : "",
      description: isSet(object.description) ? String(object.description) : "",
      uri: isSet(object.uri) ? String(object.uri) : "",
      uriHash: isSet(object.uriHash) ? String(object.uriHash) : "",
    };
  },

  toJSON(message: MsgNewClass): unknown {
    const obj: any = {};
    message.classId !== undefined && (obj.classId = message.classId);
    message.totalSupply !== undefined && (obj.totalSupply = Math.round(message.totalSupply));
    message.sender !== undefined && (obj.sender = message.sender);
    message.name !== undefined && (obj.name = message.name);
    message.symbol !== undefined && (obj.symbol = message.symbol);
    message.description !== undefined && (obj.description = message.description);
    message.uri !== undefined && (obj.uri = message.uri);
    message.uriHash !== undefined && (obj.uriHash = message.uriHash);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgNewClass>, I>>(object: I): MsgNewClass {
    const message = createBaseMsgNewClass();
    message.classId = object.classId ?? "";
    message.totalSupply = object.totalSupply ?? 0;
    message.sender = object.sender ?? "";
    message.name = object.name ?? "";
    message.symbol = object.symbol ?? "";
    message.description = object.description ?? "";
    message.uri = object.uri ?? "";
    message.uriHash = object.uriHash ?? "";
    return message;
  },
};

function createBaseMsgNewClassResponse(): MsgNewClassResponse {
  return {};
}

export const MsgNewClassResponse = {
  encode(_: MsgNewClassResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgNewClassResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgNewClassResponse();
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

  fromJSON(_: any): MsgNewClassResponse {
    return {};
  },

  toJSON(_: MsgNewClassResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgNewClassResponse>, I>>(_: I): MsgNewClassResponse {
    const message = createBaseMsgNewClassResponse();
    return message;
  },
};

function createBaseMsgMintNFT(): MsgMintNFT {
  return { classId: "", tokenId: "", uri: "", uriHash: "", creator: "", receiver: "" };
}

export const MsgMintNFT = {
  encode(message: MsgMintNFT, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.classId !== "") {
      writer.uint32(10).string(message.classId);
    }
    if (message.tokenId !== "") {
      writer.uint32(18).string(message.tokenId);
    }
    if (message.uri !== "") {
      writer.uint32(26).string(message.uri);
    }
    if (message.uriHash !== "") {
      writer.uint32(34).string(message.uriHash);
    }
    if (message.creator !== "") {
      writer.uint32(42).string(message.creator);
    }
    if (message.receiver !== "") {
      writer.uint32(50).string(message.receiver);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgMintNFT {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgMintNFT();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.classId = reader.string();
          break;
        case 2:
          message.tokenId = reader.string();
          break;
        case 3:
          message.uri = reader.string();
          break;
        case 4:
          message.uriHash = reader.string();
          break;
        case 5:
          message.creator = reader.string();
          break;
        case 6:
          message.receiver = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgMintNFT {
    return {
      classId: isSet(object.classId) ? String(object.classId) : "",
      tokenId: isSet(object.tokenId) ? String(object.tokenId) : "",
      uri: isSet(object.uri) ? String(object.uri) : "",
      uriHash: isSet(object.uriHash) ? String(object.uriHash) : "",
      creator: isSet(object.creator) ? String(object.creator) : "",
      receiver: isSet(object.receiver) ? String(object.receiver) : "",
    };
  },

  toJSON(message: MsgMintNFT): unknown {
    const obj: any = {};
    message.classId !== undefined && (obj.classId = message.classId);
    message.tokenId !== undefined && (obj.tokenId = message.tokenId);
    message.uri !== undefined && (obj.uri = message.uri);
    message.uriHash !== undefined && (obj.uriHash = message.uriHash);
    message.creator !== undefined && (obj.creator = message.creator);
    message.receiver !== undefined && (obj.receiver = message.receiver);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgMintNFT>, I>>(object: I): MsgMintNFT {
    const message = createBaseMsgMintNFT();
    message.classId = object.classId ?? "";
    message.tokenId = object.tokenId ?? "";
    message.uri = object.uri ?? "";
    message.uriHash = object.uriHash ?? "";
    message.creator = object.creator ?? "";
    message.receiver = object.receiver ?? "";
    return message;
  },
};

function createBaseMsgMintNFTResponse(): MsgMintNFTResponse {
  return {};
}

export const MsgMintNFTResponse = {
  encode(_: MsgMintNFTResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgMintNFTResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgMintNFTResponse();
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

  fromJSON(_: any): MsgMintNFTResponse {
    return {};
  },

  toJSON(_: MsgMintNFTResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgMintNFTResponse>, I>>(_: I): MsgMintNFTResponse {
    const message = createBaseMsgMintNFTResponse();
    return message;
  },
};

function createBaseMsgSend(): MsgSend {
  return { classId: "", id: "", sender: "", receiver: "" };
}

export const MsgSend = {
  encode(message: MsgSend, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.classId !== "") {
      writer.uint32(10).string(message.classId);
    }
    if (message.id !== "") {
      writer.uint32(18).string(message.id);
    }
    if (message.sender !== "") {
      writer.uint32(26).string(message.sender);
    }
    if (message.receiver !== "") {
      writer.uint32(34).string(message.receiver);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgSend {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgSend();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.classId = reader.string();
          break;
        case 2:
          message.id = reader.string();
          break;
        case 3:
          message.sender = reader.string();
          break;
        case 4:
          message.receiver = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgSend {
    return {
      classId: isSet(object.classId) ? String(object.classId) : "",
      id: isSet(object.id) ? String(object.id) : "",
      sender: isSet(object.sender) ? String(object.sender) : "",
      receiver: isSet(object.receiver) ? String(object.receiver) : "",
    };
  },

  toJSON(message: MsgSend): unknown {
    const obj: any = {};
    message.classId !== undefined && (obj.classId = message.classId);
    message.id !== undefined && (obj.id = message.id);
    message.sender !== undefined && (obj.sender = message.sender);
    message.receiver !== undefined && (obj.receiver = message.receiver);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgSend>, I>>(object: I): MsgSend {
    const message = createBaseMsgSend();
    message.classId = object.classId ?? "";
    message.id = object.id ?? "";
    message.sender = object.sender ?? "";
    message.receiver = object.receiver ?? "";
    return message;
  },
};

function createBaseMsgSendResponse(): MsgSendResponse {
  return {};
}

export const MsgSendResponse = {
  encode(_: MsgSendResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgSendResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgSendResponse();
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

  fromJSON(_: any): MsgSendResponse {
    return {};
  },

  toJSON(_: MsgSendResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgSendResponse>, I>>(_: I): MsgSendResponse {
    const message = createBaseMsgSendResponse();
    return message;
  },
};

/** Msg defines the nft Msg service. */
export interface Msg {
  /** NewClass defines a method to new class. */
  NewClass(request: MsgNewClass): Promise<MsgNewClassResponse>;
  /** MintNFT defines a method to new nft. */
  MintNFT(request: MsgMintNFT): Promise<MsgMintNFTResponse>;
  /** Send defines a method to send a nft from one account to another account. */
  Send(request: MsgSend): Promise<MsgSendResponse>;
}

export class MsgClientImpl implements Msg {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.NewClass = this.NewClass.bind(this);
    this.MintNFT = this.MintNFT.bind(this);
    this.Send = this.Send.bind(this);
  }
  NewClass(request: MsgNewClass): Promise<MsgNewClassResponse> {
    const data = MsgNewClass.encode(request).finish();
    const promise = this.rpc.request("metaearth.wnft.Msg", "NewClass", data);
    return promise.then((data) => MsgNewClassResponse.decode(new _m0.Reader(data)));
  }

  MintNFT(request: MsgMintNFT): Promise<MsgMintNFTResponse> {
    const data = MsgMintNFT.encode(request).finish();
    const promise = this.rpc.request("metaearth.wnft.Msg", "MintNFT", data);
    return promise.then((data) => MsgMintNFTResponse.decode(new _m0.Reader(data)));
  }

  Send(request: MsgSend): Promise<MsgSendResponse> {
    const data = MsgSend.encode(request).finish();
    const promise = this.rpc.request("metaearth.wnft.Msg", "Send", data);
    return promise.then((data) => MsgSendResponse.decode(new _m0.Reader(data)));
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
