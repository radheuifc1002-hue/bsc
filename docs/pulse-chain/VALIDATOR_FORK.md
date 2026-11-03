# PULSE CHAIN validator-system contract fork notes

Status: PULSE StakeHub Solidity fork implemented; runtime bytecode generated and embedded; Go upgrade routing added. Genesis installation and end-to-end chain validation remain outstanding.

## Upstream StakeHub findings

The upstream `bnb-chain/bsc-genesis-contract` `contracts/StakeHub.sol`
currently defines these initialization values:

| StakeHub parameter | Upstream default |
|---|---:|
| `minSelfDelegationBNB` | 2,000 BNB (upstream default) |
| `minDelegationBNBChange` | 1 BNB |
| `maxElectedValidators` | 45 |
| `unbondPeriod` | 7 days |
| `redelegateFeeRate` | 2 |

Its governance parameter update path constrains `minSelfDelegationBNB`.
The upstream BNB deployment expresses the value in BNB; PULSE CHAIN must use
PULSE-denominated bounds consistent with the PULSE minimum and maximum.

## Exact upstream revision and build settings

The BSC client pins the Fermi StakeHub artifact to source commit
`34618f607f8356cf147dde6a69fae150bd53d5bf` and the Pasteur StakeHub artifact
to source commit `041881a02475638b19f3d840871b7621cdebd8f8`, as recorded in
`core/systemcontracts/upgrade.go`.

The pinned Pasteur source uses Solidity `0.8.17`. PULSE pins Solidity `0.8.17`, optimizer enabled with 200 runs, and `bytecode_hash = "none"`. The upstream package manifest and lockfile are vendored at the pinned revision; OpenZeppelin dependencies are locked.

- Fermi source: https://github.com/bnb-chain/bsc-genesis-contract/commit/34618f607f8356cf147dde6a69fae150bd53d5bf
- Pasteur source: https://github.com/bnb-chain/bsc-genesis-contract/commit/041881a02475638b19f3d840871b7621cdebd8f8
- StakeHub source at Pasteur: https://github.com/bnb-chain/bsc-genesis-contract/blob/041881a02475638b19f3d840871b7621cdebd8f8/contracts/StakeHub.sol

## PULSE CHAIN implementation

The fork is under `pulse-contracts/` with the modified StakeHub, its direct
Solidity dependencies, upstream license notice, Foundry configuration, npm
lockfile, and focused tests.

The generated runtime is embedded at
`core/systemcontracts/pulse/StakeHubContract`. The Go client selects the PULSE
upgrade map by chain ID 676 and substitutes only StakeHub during the Pasteur
transition. Other Pasteur system contracts remain upstream. The embedded
runtime is raw hexadecimal without a `0x` prefix, matching the BSC loader.

Foundry build and three tests passed in CI: initialization values, exact-cap
acceptance, and above-cap rejection. Genesis installation and end-to-end
validator tests remain outstanding.

## PULSE CHAIN requirement

- Initialize the minimum self-delegation/validator eligibility threshold to
  200,000 PULSE (18-decimal base units in the EVM contract).
- Constrain governance updates of the minimum threshold to 200,000–10,000,000
  PULSE.
- Independently enforce a 10,000,000 PULSE cap on each validator's own
  self-delegated stake. This is a per-validator cap, not a cap on total stake
  attracted from third-party delegators.
- Keep the elected-validator capacity at the upstream value of 45.
- Preserve the remaining upstream staking and validator behavior unless a
  separate, explicitly approved change requires otherwise.
- Do not infer the active Parlia block-producer count from
  `maxElectedValidators`; the elected candidate pool and active consensus
  producer set are distinct protocol concepts and must be traced separately.

## Required implementation procedure

1. Pin the exact upstream `bsc-genesis-contract` source revision that matches
   the selected BSC source revision.
2. Vendor or fork the exact Solidity source and its build configuration with
   license notices intact.
3. Change the minimum self-delegation default to 200,000 PULSE; set a separately
   enforced 10,000,000 PULSE maximum on validator self-delegation; and update
   governance validation so the minimum threshold cannot be set below
   200,000 PULSE.
4. Preserve the 45 elected-validator capacity and all unrelated parameters; enforce the self-stake cap on validator creation, direct self-delegation, and self-redelegation without capping third-party delegations.
5. Rebuild the affected system-contract bytecode using the pinned compiler and
   settings; record compiler version, optimizer settings, metadata settings,
   source hashes, and resulting bytecode hashes.
6. Update the matching versioned StakeHub artifact(s) and any genesis storage
   initialization paths. Do not patch hex bytecode by search-and-replace.
7. Verify the deployed runtime bytecode, initialization storage, parameter
   getters, validator election behavior, and governance updates on a local
   testnet before using it in a public testnet genesis.

The BSC fork repository currently carries versioned system-contract bytecode
artifacts under `core/systemcontracts`, but not the Solidity source tree.
The artifact versions and source revision must be matched before modifying any
artifact. No contract artifact has been changed in this phase.
