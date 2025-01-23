/* eslint-disable */
export const protobufPackage = "dymensionxyz.dymension.common";

export enum Status {
  PENDING = 0,
  FINALIZED = 1,
  REVERTED = 3,
  UNRECOGNIZED = -1,
}

export function statusFromJSON(object: any): Status {
  switch (object) {
    case 0:
    case "PENDING":
      return Status.PENDING;
    case 1:
    case "FINALIZED":
      return Status.FINALIZED;
    case 3:
    case "REVERTED":
      return Status.REVERTED;
    case -1:
    case "UNRECOGNIZED":
    default:
      return Status.UNRECOGNIZED;
  }
}

export function statusToJSON(object: Status): string {
  switch (object) {
    case Status.PENDING:
      return "PENDING";
    case Status.FINALIZED:
      return "FINALIZED";
    case Status.REVERTED:
      return "REVERTED";
    case Status.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}
