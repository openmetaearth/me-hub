# ME-Hub v2.0.18.rc1

Release candidate for upgrading an existing v2.0.18 chain. This is a prerelease, not the latest stable release.

## Changes

- Register the `v2.0.18.rc1` on-chain upgrade handler without reinitializing regiongov or adding its store again.
- Migrate existing regiongov proposal `required_operator` values to their corresponding ordinary account addresses (`me1...`). Preserve address bytes, snapshotted approvers, elector membership, votes, proposal status, and voting queues.
- Store new governance electors and required operators as account addresses; keep staking `Region.OperatorAddress`, validator operator addresses, and replacement `new_operator` values in `mevaloper1...` format.
- Preserve voting authorization across operator replacement and region/equity freezes. Equity withdrawals require the equity threshold plus explicit operator and DAO approvals.
- Use standard protobuf API descriptors for regiongov HTTP gateway enum handling.
- Support uppercase and mixed-case region IDs in region freeze queries consistently with freeze/unfreeze transactions.
- Require an explicit upgrade height when submitting an upgrade proposal with `scripts/upgrade.sh`.

## Upgrade instructions

1. Back up node configuration, keys, and chain data using your normal operational procedure. Do not delete or reset existing node data.
2. Use the exact on-chain plan name `v2.0.18.rc1`. The upgrade height is chosen by the governance proposal, not hard-coded in this release.
3. Stage the new binary, but switch the running node at the scheduled upgrade halt. This SDK rejects starting a binary containing the pending upgrade handler before its trigger height.
4. Verify the downloaded archive against `SHA256SUMS-v2.0.18.rc1.txt`. Check `med version` and `med version --long` before deployment.
5. Restart with the new binary and the existing node home. Confirm the upgrade completes and block production resumes. Do not skip the upgrade height; it performs the required proposal migration.

An old binary reporting `UPGRADE "v2.0.18.rc1" NEEDED` at the scheduled height is the expected upgrade halt. A persistent occurrence after replacement means the running executable must be checked.

If this upgrade name has already been applied on a chain, replacing the executable does not rerun its migration; use a separately planned upgrade for additional migrations.

## Artifacts

Build workflows attach a Linux amd64 static binary archive, SHA256 checksum file, and build version details after verification succeeds. No macOS or Linux arm64 binary is produced by this release workflow.

The separate container workflow publishes:

```sh
docker pull ghcr.io/openmetaearth/me-hub:v2.0.18.rc1
```

This release candidate does not update the Docker `latest` tag. When upgrading containers, ensure the actual container image and any mounted binary are updated; restarting an old container alone does not replace its executable.
