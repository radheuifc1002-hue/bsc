# PULSE CHAIN system-contract source fork

This directory contains the PULSE-specific StakeHub source fork and the minimal
upstream Solidity files it imports directly. The source baseline is the BSC
genesis-contract Pasteur revision
`041881a02475638b19f3d840871b7621cdebd8f8`.

Upstream source:
https://github.com/bnb-chain/bsc-genesis-contract/tree/041881a02475638b19f3d840871b7621cdebd8f8/contracts

## Local changes

Only StakeHub's validator self-stake parameters and enforcement are changed:

- Initial `minSelfDelegationBNB`: 200,000 PULSE.
- Governance range for that minimum: 200,000–10,000,000 PULSE.
- Maximum operator self-delegation: 10,000,000 PULSE.
- The self-stake cap is checked during validator creation, operator direct
  delegation, and operator self-redelegation into a destination validator.
- Third-party delegation is not subject to the operator self-stake cap.
- `maxElectedValidators` remains 45; other StakeHub defaults and logic are
  inherited from the pinned upstream source.

The cap is checked against the operator's currently pooled self-stake when an
additional operator contribution is made. Rewards accruing to existing shares
are not confiscated or burned if their pooled value later rises above the cap;
new operator contributions are rejected while they would exceed it.

## Build settings

- Solidity: 0.8.17 (explicitly pinned)
- Optimizer: enabled
- Optimizer runs: 200
- Metadata bytecode hash: none
- OpenZeppelin contracts and contracts-upgradeable: 4.9.6, via the upstream
  lockfile at the pinned source revision

Run from this directory:

```sh
npm ci
forge build
forge test -vvv
```

The GitHub Actions workflow runs the build and tests, exports the runtime as
raw hexadecimal, commits it to `core/systemcontracts/pulse/StakeHubContract`,
and runs Go tests for the PULSE system-contract routing.

## Deployment status

The runtime is embedded in the BSC client and selected by the PULSE-specific
Pasteur upgrade map. It is not yet installed into a finalized fresh-network
genesis allocation or validated on a running multi-validator PULSE testnet.
Genesis storage initialization and end-to-end staking tests remain required.

The upstream contract files retain their original SPDX headers. The upstream
license notice is included in this directory.
