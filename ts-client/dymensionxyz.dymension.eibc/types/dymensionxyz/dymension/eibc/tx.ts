/* eslint-disable */
import _m0 from "protobufjs/minimal";

export const protobufPackage = "dymensionxyz.dymension.eibc";

/** MsgFulfillOrder defines the FulfillOrder request type. */
export interface MsgFulfillOrder {
  /**
   * fulfiller_address is the bech32-encoded address of the account which the
   * message was sent from.
   */
  fulfillerAddress: string;
  /** order_id is the unique identifier of the order to be fulfilled. */
  orderId: string;
  /** expected_fee is the nominal fee set in the order. */
  expectedFee: string;
}

/** MsgFulfillOrderResponse defines the FulfillOrder response type. */
export interface MsgFulfillOrderResponse {
}

export interface MsgUpdateDemandOrder {
  /**
   * owner_address is the bech32-encoded address of the account owns the order.
   * This is expected to be the address of the order recipient.
   */
  ownerAddress: string;
  /** order_id is the unique identifier of the order to be updated. */
  orderId: string;
  /** new_fee is the new fee amount to be set in the order. */
  newFee: string;
}

export interface MsgUpdateDemandOrderResponse {
}

function createBaseMsgFulfillOrder(): MsgFulfillOrder {
  return { fulfillerAddress: "", orderId: "", expectedFee: "" };
}

export const MsgFulfillOrder = {
  encode(message: MsgFulfillOrder, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.fulfillerAddress !== "") {
      writer.uint32(10).string(message.fulfillerAddress);
    }
    if (message.orderId !== "") {
      writer.uint32(18).string(message.orderId);
    }
    if (message.expectedFee !== "") {
      writer.uint32(26).string(message.expectedFee);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgFulfillOrder {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgFulfillOrder();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.fulfillerAddress = reader.string();
          break;
        case 2:
          message.orderId = reader.string();
          break;
        case 3:
          message.expectedFee = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgFulfillOrder {
    return {
      fulfillerAddress: isSet(object.fulfillerAddress) ? String(object.fulfillerAddress) : "",
      orderId: isSet(object.orderId) ? String(object.orderId) : "",
      expectedFee: isSet(object.expectedFee) ? String(object.expectedFee) : "",
    };
  },

  toJSON(message: MsgFulfillOrder): unknown {
    const obj: any = {};
    message.fulfillerAddress !== undefined && (obj.fulfillerAddress = message.fulfillerAddress);
    message.orderId !== undefined && (obj.orderId = message.orderId);
    message.expectedFee !== undefined && (obj.expectedFee = message.expectedFee);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgFulfillOrder>, I>>(object: I): MsgFulfillOrder {
    const message = createBaseMsgFulfillOrder();
    message.fulfillerAddress = object.fulfillerAddress ?? "";
    message.orderId = object.orderId ?? "";
    message.expectedFee = object.expectedFee ?? "";
    return message;
  },
};

function createBaseMsgFulfillOrderResponse(): MsgFulfillOrderResponse {
  return {};
}

export const MsgFulfillOrderResponse = {
  encode(_: MsgFulfillOrderResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgFulfillOrderResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgFulfillOrderResponse();
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

  fromJSON(_: any): MsgFulfillOrderResponse {
    return {};
  },

  toJSON(_: MsgFulfillOrderResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgFulfillOrderResponse>, I>>(_: I): MsgFulfillOrderResponse {
    const message = createBaseMsgFulfillOrderResponse();
    return message;
  },
};

function createBaseMsgUpdateDemandOrder(): MsgUpdateDemandOrder {
  return { ownerAddress: "", orderId: "", newFee: "" };
}

export const MsgUpdateDemandOrder = {
  encode(message: MsgUpdateDemandOrder, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.ownerAddress !== "") {
      writer.uint32(10).string(message.ownerAddress);
    }
    if (message.orderId !== "") {
      writer.uint32(18).string(message.orderId);
    }
    if (message.newFee !== "") {
      writer.uint32(26).string(message.newFee);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateDemandOrder {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateDemandOrder();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.ownerAddress = reader.string();
          break;
        case 2:
          message.orderId = reader.string();
          break;
        case 3:
          message.newFee = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgUpdateDemandOrder {
    return {
      ownerAddress: isSet(object.ownerAddress) ? String(object.ownerAddress) : "",
      orderId: isSet(object.orderId) ? String(object.orderId) : "",
      newFee: isSet(object.newFee) ? String(object.newFee) : "",
    };
  },

  toJSON(message: MsgUpdateDemandOrder): unknown {
    const obj: any = {};
    message.ownerAddress !== undefined && (obj.ownerAddress = message.ownerAddress);
    message.orderId !== undefined && (obj.orderId = message.orderId);
    message.newFee !== undefined && (obj.newFee = message.newFee);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgUpdateDemandOrder>, I>>(object: I): MsgUpdateDemandOrder {
    const message = createBaseMsgUpdateDemandOrder();
    message.ownerAddress = object.ownerAddress ?? "";
    message.orderId = object.orderId ?? "";
    message.newFee = object.newFee ?? "";
    return message;
  },
};

function createBaseMsgUpdateDemandOrderResponse(): MsgUpdateDemandOrderResponse {
  return {};
}

export const MsgUpdateDemandOrderResponse = {
  encode(_: MsgUpdateDemandOrderResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateDemandOrderResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateDemandOrderResponse();
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

  fromJSON(_: any): MsgUpdateDemandOrderResponse {
    return {};
  },

  toJSON(_: MsgUpdateDemandOrderResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgUpdateDemandOrderResponse>, I>>(_: I): MsgUpdateDemandOrderResponse {
    const message = createBaseMsgUpdateDemandOrderResponse();
    return message;
  },
};

/** Msg defines the Msg service. */
export interface Msg {
  FulfillOrder(request: MsgFulfillOrder): Promise<MsgFulfillOrderResponse>;
  UpdateDemandOrder(request: MsgUpdateDemandOrder): Promise<MsgUpdateDemandOrderResponse>;
}

export class MsgClientImpl implements Msg {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.FulfillOrder = this.FulfillOrder.bind(this);
    this.UpdateDemandOrder = this.UpdateDemandOrder.bind(this);
  }
  FulfillOrder(request: MsgFulfillOrder): Promise<MsgFulfillOrderResponse> {
    const data = MsgFulfillOrder.encode(request).finish();
    const promise = this.rpc.request("dymensionxyz.dymension.eibc.Msg", "FulfillOrder", data);
    return promise.then((data) => MsgFulfillOrderResponse.decode(new _m0.Reader(data)));
  }

  UpdateDemandOrder(request: MsgUpdateDemandOrder): Promise<MsgUpdateDemandOrderResponse> {
    const data = MsgUpdateDemandOrder.encode(request).finish();
    const promise = this.rpc.request("dymensionxyz.dymension.eibc.Msg", "UpdateDemandOrder", data);
    return promise.then((data) => MsgUpdateDemandOrderResponse.decode(new _m0.Reader(data)));
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
