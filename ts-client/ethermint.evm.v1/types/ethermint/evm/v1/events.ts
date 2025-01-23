/* eslint-disable */
import _m0 from "protobufjs/minimal";

export const protobufPackage = "ethermint.evm.v1";

/** EventEthereumTx defines the event for an Ethereum transaction */
export interface EventEthereumTx {
  /** amount */
  amount: string;
  /** eth_hash is the Ethereum hash of the transaction */
  ethHash: string;
  /** index of the transaction in the block */
  index: string;
  /** gas_used is the amount of gas used by the transaction */
  gasUsed: string;
  /** hash is the Tendermint hash of the transaction */
  hash: string;
  /** recipient of the transaction */
  recipient: string;
  /** eth_tx_failed contains a VM error should it occur */
  ethTxFailed: string;
}

/** EventTxLog defines the event for an Ethereum transaction log */
export interface EventTxLog {
  /** tx_logs is an array of transaction logs */
  txLogs: string[];
}

/** EventMessage */
export interface EventMessage {
  /** module which emits the event */
  module: string;
  /** sender of the message */
  sender: string;
  /** tx_type is the type of the message */
  txType: string;
}

/** EventBlockBloom defines an Ethereum block bloom filter event */
export interface EventBlockBloom {
  /** bloom is the bloom filter of the block */
  bloom: string;
}

function createBaseEventEthereumTx(): EventEthereumTx {
  return { amount: "", ethHash: "", index: "", gasUsed: "", hash: "", recipient: "", ethTxFailed: "" };
}

export const EventEthereumTx = {
  encode(message: EventEthereumTx, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.amount !== "") {
      writer.uint32(10).string(message.amount);
    }
    if (message.ethHash !== "") {
      writer.uint32(18).string(message.ethHash);
    }
    if (message.index !== "") {
      writer.uint32(26).string(message.index);
    }
    if (message.gasUsed !== "") {
      writer.uint32(34).string(message.gasUsed);
    }
    if (message.hash !== "") {
      writer.uint32(42).string(message.hash);
    }
    if (message.recipient !== "") {
      writer.uint32(50).string(message.recipient);
    }
    if (message.ethTxFailed !== "") {
      writer.uint32(58).string(message.ethTxFailed);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): EventEthereumTx {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseEventEthereumTx();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.amount = reader.string();
          break;
        case 2:
          message.ethHash = reader.string();
          break;
        case 3:
          message.index = reader.string();
          break;
        case 4:
          message.gasUsed = reader.string();
          break;
        case 5:
          message.hash = reader.string();
          break;
        case 6:
          message.recipient = reader.string();
          break;
        case 7:
          message.ethTxFailed = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): EventEthereumTx {
    return {
      amount: isSet(object.amount) ? String(object.amount) : "",
      ethHash: isSet(object.ethHash) ? String(object.ethHash) : "",
      index: isSet(object.index) ? String(object.index) : "",
      gasUsed: isSet(object.gasUsed) ? String(object.gasUsed) : "",
      hash: isSet(object.hash) ? String(object.hash) : "",
      recipient: isSet(object.recipient) ? String(object.recipient) : "",
      ethTxFailed: isSet(object.ethTxFailed) ? String(object.ethTxFailed) : "",
    };
  },

  toJSON(message: EventEthereumTx): unknown {
    const obj: any = {};
    message.amount !== undefined && (obj.amount = message.amount);
    message.ethHash !== undefined && (obj.ethHash = message.ethHash);
    message.index !== undefined && (obj.index = message.index);
    message.gasUsed !== undefined && (obj.gasUsed = message.gasUsed);
    message.hash !== undefined && (obj.hash = message.hash);
    message.recipient !== undefined && (obj.recipient = message.recipient);
    message.ethTxFailed !== undefined && (obj.ethTxFailed = message.ethTxFailed);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<EventEthereumTx>, I>>(object: I): EventEthereumTx {
    const message = createBaseEventEthereumTx();
    message.amount = object.amount ?? "";
    message.ethHash = object.ethHash ?? "";
    message.index = object.index ?? "";
    message.gasUsed = object.gasUsed ?? "";
    message.hash = object.hash ?? "";
    message.recipient = object.recipient ?? "";
    message.ethTxFailed = object.ethTxFailed ?? "";
    return message;
  },
};

function createBaseEventTxLog(): EventTxLog {
  return { txLogs: [] };
}

export const EventTxLog = {
  encode(message: EventTxLog, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.txLogs) {
      writer.uint32(10).string(v!);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): EventTxLog {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseEventTxLog();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.txLogs.push(reader.string());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): EventTxLog {
    return { txLogs: Array.isArray(object?.txLogs) ? object.txLogs.map((e: any) => String(e)) : [] };
  },

  toJSON(message: EventTxLog): unknown {
    const obj: any = {};
    if (message.txLogs) {
      obj.txLogs = message.txLogs.map((e) => e);
    } else {
      obj.txLogs = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<EventTxLog>, I>>(object: I): EventTxLog {
    const message = createBaseEventTxLog();
    message.txLogs = object.txLogs?.map((e) => e) || [];
    return message;
  },
};

function createBaseEventMessage(): EventMessage {
  return { module: "", sender: "", txType: "" };
}

export const EventMessage = {
  encode(message: EventMessage, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.module !== "") {
      writer.uint32(10).string(message.module);
    }
    if (message.sender !== "") {
      writer.uint32(18).string(message.sender);
    }
    if (message.txType !== "") {
      writer.uint32(26).string(message.txType);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): EventMessage {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseEventMessage();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.module = reader.string();
          break;
        case 2:
          message.sender = reader.string();
          break;
        case 3:
          message.txType = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): EventMessage {
    return {
      module: isSet(object.module) ? String(object.module) : "",
      sender: isSet(object.sender) ? String(object.sender) : "",
      txType: isSet(object.txType) ? String(object.txType) : "",
    };
  },

  toJSON(message: EventMessage): unknown {
    const obj: any = {};
    message.module !== undefined && (obj.module = message.module);
    message.sender !== undefined && (obj.sender = message.sender);
    message.txType !== undefined && (obj.txType = message.txType);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<EventMessage>, I>>(object: I): EventMessage {
    const message = createBaseEventMessage();
    message.module = object.module ?? "";
    message.sender = object.sender ?? "";
    message.txType = object.txType ?? "";
    return message;
  },
};

function createBaseEventBlockBloom(): EventBlockBloom {
  return { bloom: "" };
}

export const EventBlockBloom = {
  encode(message: EventBlockBloom, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.bloom !== "") {
      writer.uint32(10).string(message.bloom);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): EventBlockBloom {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseEventBlockBloom();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.bloom = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): EventBlockBloom {
    return { bloom: isSet(object.bloom) ? String(object.bloom) : "" };
  },

  toJSON(message: EventBlockBloom): unknown {
    const obj: any = {};
    message.bloom !== undefined && (obj.bloom = message.bloom);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<EventBlockBloom>, I>>(object: I): EventBlockBloom {
    const message = createBaseEventBlockBloom();
    message.bloom = object.bloom ?? "";
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
