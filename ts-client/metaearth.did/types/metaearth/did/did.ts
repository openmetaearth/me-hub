/* eslint-disable */
import _m0 from "protobufjs/minimal";
import { Credential } from "./credential";

export const protobufPackage = "metaearth.did";

/** kyc level enum */
export enum KycLevel {
  KYC_LEVEL_NONE = 0,
  KYC_LEVEL_ONE = 1,
  KYC_LEVEL_TWO = 2,
  KYC_LEVEL_THREE = 3,
  KYC_LEVEL_FOUR = 4,
  KYC_LEVEL_FIVE = 5,
  UNRECOGNIZED = -1,
}

export function kycLevelFromJSON(object: any): KycLevel {
  switch (object) {
    case 0:
    case "KYC_LEVEL_NONE":
      return KycLevel.KYC_LEVEL_NONE;
    case 1:
    case "KYC_LEVEL_ONE":
      return KycLevel.KYC_LEVEL_ONE;
    case 2:
    case "KYC_LEVEL_TWO":
      return KycLevel.KYC_LEVEL_TWO;
    case 3:
    case "KYC_LEVEL_THREE":
      return KycLevel.KYC_LEVEL_THREE;
    case 4:
    case "KYC_LEVEL_FOUR":
      return KycLevel.KYC_LEVEL_FOUR;
    case 5:
    case "KYC_LEVEL_FIVE":
      return KycLevel.KYC_LEVEL_FIVE;
    case -1:
    case "UNRECOGNIZED":
    default:
      return KycLevel.UNRECOGNIZED;
  }
}

export function kycLevelToJSON(object: KycLevel): string {
  switch (object) {
    case KycLevel.KYC_LEVEL_NONE:
      return "KYC_LEVEL_NONE";
    case KycLevel.KYC_LEVEL_ONE:
      return "KYC_LEVEL_ONE";
    case KycLevel.KYC_LEVEL_TWO:
      return "KYC_LEVEL_TWO";
    case KycLevel.KYC_LEVEL_THREE:
      return "KYC_LEVEL_THREE";
    case KycLevel.KYC_LEVEL_FOUR:
      return "KYC_LEVEL_FOUR";
    case KycLevel.KYC_LEVEL_FIVE:
      return "KYC_LEVEL_FIVE";
    case KycLevel.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

/** did status enum */
export enum DidStatus {
  DID_STATUS_INACTIVE = 0,
  DID_STATUS_ACTIVE = 1,
  UNRECOGNIZED = -1,
}

export function didStatusFromJSON(object: any): DidStatus {
  switch (object) {
    case 0:
    case "DID_STATUS_INACTIVE":
      return DidStatus.DID_STATUS_INACTIVE;
    case 1:
    case "DID_STATUS_ACTIVE":
      return DidStatus.DID_STATUS_ACTIVE;
    case -1:
    case "UNRECOGNIZED":
    default:
      return DidStatus.UNRECOGNIZED;
  }
}

export function didStatusToJSON(object: DidStatus): string {
  switch (object) {
    case DidStatus.DID_STATUS_INACTIVE:
      return "DID_STATUS_INACTIVE";
    case DidStatus.DID_STATUS_ACTIVE:
      return "DID_STATUS_ACTIVE";
    case DidStatus.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export interface DidInfo {
  /** same as the MEID */
  did: string;
  /** MEID.account */
  address: string;
  /**
   * public_key is mapped to the user address
   * the issuer will use public_key to encrypt the user's certificate to ensure
   * the privacy of the off-chain certificate
   */
  pubkey: string;
  status: DidStatus;
  /** unused! */
  regionId: string;
  kycLevel: KycLevel;
}

/** did document */
export interface DidDocument {
  info: DidInfo | undefined;
  vcs: Credential[];
}

function createBaseDidInfo(): DidInfo {
  return { did: "", address: "", pubkey: "", status: 0, regionId: "", kycLevel: 0 };
}

export const DidInfo = {
  encode(message: DidInfo, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.did !== "") {
      writer.uint32(10).string(message.did);
    }
    if (message.address !== "") {
      writer.uint32(18).string(message.address);
    }
    if (message.pubkey !== "") {
      writer.uint32(26).string(message.pubkey);
    }
    if (message.status !== 0) {
      writer.uint32(32).int32(message.status);
    }
    if (message.regionId !== "") {
      writer.uint32(42).string(message.regionId);
    }
    if (message.kycLevel !== 0) {
      writer.uint32(48).int32(message.kycLevel);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DidInfo {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDidInfo();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.did = reader.string();
          break;
        case 2:
          message.address = reader.string();
          break;
        case 3:
          message.pubkey = reader.string();
          break;
        case 4:
          message.status = reader.int32() as any;
          break;
        case 5:
          message.regionId = reader.string();
          break;
        case 6:
          message.kycLevel = reader.int32() as any;
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): DidInfo {
    return {
      did: isSet(object.did) ? String(object.did) : "",
      address: isSet(object.address) ? String(object.address) : "",
      pubkey: isSet(object.pubkey) ? String(object.pubkey) : "",
      status: isSet(object.status) ? didStatusFromJSON(object.status) : 0,
      regionId: isSet(object.regionId) ? String(object.regionId) : "",
      kycLevel: isSet(object.kycLevel) ? kycLevelFromJSON(object.kycLevel) : 0,
    };
  },

  toJSON(message: DidInfo): unknown {
    const obj: any = {};
    message.did !== undefined && (obj.did = message.did);
    message.address !== undefined && (obj.address = message.address);
    message.pubkey !== undefined && (obj.pubkey = message.pubkey);
    message.status !== undefined && (obj.status = didStatusToJSON(message.status));
    message.regionId !== undefined && (obj.regionId = message.regionId);
    message.kycLevel !== undefined && (obj.kycLevel = kycLevelToJSON(message.kycLevel));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<DidInfo>, I>>(object: I): DidInfo {
    const message = createBaseDidInfo();
    message.did = object.did ?? "";
    message.address = object.address ?? "";
    message.pubkey = object.pubkey ?? "";
    message.status = object.status ?? 0;
    message.regionId = object.regionId ?? "";
    message.kycLevel = object.kycLevel ?? 0;
    return message;
  },
};

function createBaseDidDocument(): DidDocument {
  return { info: undefined, vcs: [] };
}

export const DidDocument = {
  encode(message: DidDocument, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.info !== undefined) {
      DidInfo.encode(message.info, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.vcs) {
      Credential.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DidDocument {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDidDocument();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.info = DidInfo.decode(reader, reader.uint32());
          break;
        case 2:
          message.vcs.push(Credential.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): DidDocument {
    return {
      info: isSet(object.info) ? DidInfo.fromJSON(object.info) : undefined,
      vcs: Array.isArray(object?.vcs) ? object.vcs.map((e: any) => Credential.fromJSON(e)) : [],
    };
  },

  toJSON(message: DidDocument): unknown {
    const obj: any = {};
    message.info !== undefined && (obj.info = message.info ? DidInfo.toJSON(message.info) : undefined);
    if (message.vcs) {
      obj.vcs = message.vcs.map((e) => e ? Credential.toJSON(e) : undefined);
    } else {
      obj.vcs = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<DidDocument>, I>>(object: I): DidDocument {
    const message = createBaseDidDocument();
    message.info = (object.info !== undefined && object.info !== null) ? DidInfo.fromPartial(object.info) : undefined;
    message.vcs = object.vcs?.map((e) => Credential.fromPartial(e)) || [];
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
