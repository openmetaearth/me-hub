/* eslint-disable */
export const protobufPackage = "dymensionxyz.dymension.sequencer";

/** OperatingStatus defines the operating status of a sequencer */
export enum OperatingStatus {
  /**
   * OPERATING_STATUS_UNBONDED - OPERATING_STATUS_UNBONDED defines a sequencer that is not active and won't
   * be scheduled
   */
  OPERATING_STATUS_UNBONDED = 0,
  /** OPERATING_STATUS_UNBONDING - UNBONDING defines a sequencer that is currently unbonding. */
  OPERATING_STATUS_UNBONDING = 1,
  /**
   * OPERATING_STATUS_BONDED - OPERATING_STATUS_BONDED defines a sequencer that is bonded and can be
   * scheduled
   */
  OPERATING_STATUS_BONDED = 2,
  UNRECOGNIZED = -1,
}

export function operatingStatusFromJSON(object: any): OperatingStatus {
  switch (object) {
    case 0:
    case "OPERATING_STATUS_UNBONDED":
      return OperatingStatus.OPERATING_STATUS_UNBONDED;
    case 1:
    case "OPERATING_STATUS_UNBONDING":
      return OperatingStatus.OPERATING_STATUS_UNBONDING;
    case 2:
    case "OPERATING_STATUS_BONDED":
      return OperatingStatus.OPERATING_STATUS_BONDED;
    case -1:
    case "UNRECOGNIZED":
    default:
      return OperatingStatus.UNRECOGNIZED;
  }
}

export function operatingStatusToJSON(object: OperatingStatus): string {
  switch (object) {
    case OperatingStatus.OPERATING_STATUS_UNBONDED:
      return "OPERATING_STATUS_UNBONDED";
    case OperatingStatus.OPERATING_STATUS_UNBONDING:
      return "OPERATING_STATUS_UNBONDING";
    case OperatingStatus.OPERATING_STATUS_BONDED:
      return "OPERATING_STATUS_BONDED";
    case OperatingStatus.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}
