/* eslint-disable */
import _m0 from "protobufjs/minimal";

export const protobufPackage = "ethermint.evm.v1";

/** VFContractType is the type of the virtual frontier contract */
export enum VFContractType {
  VFC_TYPE_UNSPECIFIED = 0,
  /** VFC_TYPE_BANK - VFC_TYPE_BANK indicates the VFC is a Virtual Frontier Bank Contract */
  VFC_TYPE_BANK = 1,
  UNRECOGNIZED = -1,
}

export function vFContractTypeFromJSON(object: any): VFContractType {
  switch (object) {
    case 0:
    case "VFC_TYPE_UNSPECIFIED":
      return VFContractType.VFC_TYPE_UNSPECIFIED;
    case 1:
    case "VFC_TYPE_BANK":
      return VFContractType.VFC_TYPE_BANK;
    case -1:
    case "UNRECOGNIZED":
    default:
      return VFContractType.UNRECOGNIZED;
  }
}

export function vFContractTypeToJSON(object: VFContractType): string {
  switch (object) {
    case VFContractType.VFC_TYPE_UNSPECIFIED:
      return "VFC_TYPE_UNSPECIFIED";
    case VFContractType.VFC_TYPE_BANK:
      return "VFC_TYPE_BANK";
    case VFContractType.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

/**
 * VirtualFrontierContract are virtual contracts that able to call directly without actual invoking EVM.
 * Prohibited to call directly from EVM.
 */
export interface VirtualFrontierContract {
  /** address is the contract address */
  address: string;
  /** active indicate the activation of the contract. If not active, invoking methods is disabled. */
  active: boolean;
  /** type of the virtual frontier contract */
  type: VFContractType;
  /** metadata contains marshalled binary of the metadata of the virtual frontier contract */
  metadata: Uint8Array;
}

/**
 * VFBankContractMetadata contains metadata of the Bank Contract, which represents ERC-20 contracts as proxy of the user bank assets.
 * Each bank contract is represented as a single asset of user bank.
 */
export interface VFBankContractMetadata {
  /** min_denom is the base denomination of the asset */
  minDenom: string;
}

/** UpdateVirtualFrontierBankContractsProposal is a gov Content type to update the virtual frontier bank contracts. */
export interface UpdateVirtualFrontierBankContractsProposal {
  /** title of the proposal */
  title: string;
  /** description of the proposal */
  description: string;
  /** contracts is slice of the update content for each contract */
  contracts: VirtualFrontierBankContractProposalContent[];
}

/**
 * VirtualFrontierBankContractProposalContent contains the adjustment to the virtual frontier bank contract.
 * Currently only used to toggle activation of the contract.
 */
export interface VirtualFrontierBankContractProposalContent {
  /** contract_address that want to be updated */
  contractAddress: string;
  /** active indicate the activation of the contract. If not active, invoking methods is disabled. */
  active: boolean;
}

function createBaseVirtualFrontierContract(): VirtualFrontierContract {
  return { address: "", active: false, type: 0, metadata: new Uint8Array() };
}

export const VirtualFrontierContract = {
  encode(message: VirtualFrontierContract, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.address !== "") {
      writer.uint32(10).string(message.address);
    }
    if (message.active === true) {
      writer.uint32(16).bool(message.active);
    }
    if (message.type !== 0) {
      writer.uint32(24).int32(message.type);
    }
    if (message.metadata.length !== 0) {
      writer.uint32(34).bytes(message.metadata);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): VirtualFrontierContract {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseVirtualFrontierContract();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.address = reader.string();
          break;
        case 2:
          message.active = reader.bool();
          break;
        case 3:
          message.type = reader.int32() as any;
          break;
        case 4:
          message.metadata = reader.bytes();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): VirtualFrontierContract {
    return {
      address: isSet(object.address) ? String(object.address) : "",
      active: isSet(object.active) ? Boolean(object.active) : false,
      type: isSet(object.type) ? vFContractTypeFromJSON(object.type) : 0,
      metadata: isSet(object.metadata) ? bytesFromBase64(object.metadata) : new Uint8Array(),
    };
  },

  toJSON(message: VirtualFrontierContract): unknown {
    const obj: any = {};
    message.address !== undefined && (obj.address = message.address);
    message.active !== undefined && (obj.active = message.active);
    message.type !== undefined && (obj.type = vFContractTypeToJSON(message.type));
    message.metadata !== undefined
      && (obj.metadata = base64FromBytes(message.metadata !== undefined ? message.metadata : new Uint8Array()));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<VirtualFrontierContract>, I>>(object: I): VirtualFrontierContract {
    const message = createBaseVirtualFrontierContract();
    message.address = object.address ?? "";
    message.active = object.active ?? false;
    message.type = object.type ?? 0;
    message.metadata = object.metadata ?? new Uint8Array();
    return message;
  },
};

function createBaseVFBankContractMetadata(): VFBankContractMetadata {
  return { minDenom: "" };
}

export const VFBankContractMetadata = {
  encode(message: VFBankContractMetadata, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.minDenom !== "") {
      writer.uint32(10).string(message.minDenom);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): VFBankContractMetadata {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseVFBankContractMetadata();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.minDenom = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): VFBankContractMetadata {
    return { minDenom: isSet(object.minDenom) ? String(object.minDenom) : "" };
  },

  toJSON(message: VFBankContractMetadata): unknown {
    const obj: any = {};
    message.minDenom !== undefined && (obj.minDenom = message.minDenom);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<VFBankContractMetadata>, I>>(object: I): VFBankContractMetadata {
    const message = createBaseVFBankContractMetadata();
    message.minDenom = object.minDenom ?? "";
    return message;
  },
};

function createBaseUpdateVirtualFrontierBankContractsProposal(): UpdateVirtualFrontierBankContractsProposal {
  return { title: "", description: "", contracts: [] };
}

export const UpdateVirtualFrontierBankContractsProposal = {
  encode(message: UpdateVirtualFrontierBankContractsProposal, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.title !== "") {
      writer.uint32(10).string(message.title);
    }
    if (message.description !== "") {
      writer.uint32(18).string(message.description);
    }
    for (const v of message.contracts) {
      VirtualFrontierBankContractProposalContent.encode(v!, writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateVirtualFrontierBankContractsProposal {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateVirtualFrontierBankContractsProposal();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.title = reader.string();
          break;
        case 2:
          message.description = reader.string();
          break;
        case 3:
          message.contracts.push(VirtualFrontierBankContractProposalContent.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UpdateVirtualFrontierBankContractsProposal {
    return {
      title: isSet(object.title) ? String(object.title) : "",
      description: isSet(object.description) ? String(object.description) : "",
      contracts: Array.isArray(object?.contracts)
        ? object.contracts.map((e: any) => VirtualFrontierBankContractProposalContent.fromJSON(e))
        : [],
    };
  },

  toJSON(message: UpdateVirtualFrontierBankContractsProposal): unknown {
    const obj: any = {};
    message.title !== undefined && (obj.title = message.title);
    message.description !== undefined && (obj.description = message.description);
    if (message.contracts) {
      obj.contracts = message.contracts.map((e) =>
        e ? VirtualFrontierBankContractProposalContent.toJSON(e) : undefined
      );
    } else {
      obj.contracts = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<UpdateVirtualFrontierBankContractsProposal>, I>>(
    object: I,
  ): UpdateVirtualFrontierBankContractsProposal {
    const message = createBaseUpdateVirtualFrontierBankContractsProposal();
    message.title = object.title ?? "";
    message.description = object.description ?? "";
    message.contracts = object.contracts?.map((e) => VirtualFrontierBankContractProposalContent.fromPartial(e)) || [];
    return message;
  },
};

function createBaseVirtualFrontierBankContractProposalContent(): VirtualFrontierBankContractProposalContent {
  return { contractAddress: "", active: false };
}

export const VirtualFrontierBankContractProposalContent = {
  encode(message: VirtualFrontierBankContractProposalContent, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.contractAddress !== "") {
      writer.uint32(10).string(message.contractAddress);
    }
    if (message.active === true) {
      writer.uint32(16).bool(message.active);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): VirtualFrontierBankContractProposalContent {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseVirtualFrontierBankContractProposalContent();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.contractAddress = reader.string();
          break;
        case 2:
          message.active = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): VirtualFrontierBankContractProposalContent {
    return {
      contractAddress: isSet(object.contractAddress) ? String(object.contractAddress) : "",
      active: isSet(object.active) ? Boolean(object.active) : false,
    };
  },

  toJSON(message: VirtualFrontierBankContractProposalContent): unknown {
    const obj: any = {};
    message.contractAddress !== undefined && (obj.contractAddress = message.contractAddress);
    message.active !== undefined && (obj.active = message.active);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<VirtualFrontierBankContractProposalContent>, I>>(
    object: I,
  ): VirtualFrontierBankContractProposalContent {
    const message = createBaseVirtualFrontierBankContractProposalContent();
    message.contractAddress = object.contractAddress ?? "";
    message.active = object.active ?? false;
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

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}
