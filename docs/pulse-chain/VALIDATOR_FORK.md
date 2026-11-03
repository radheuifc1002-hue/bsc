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

## Remaining implementation work

1. Add the PULSE StakeHub runtime and storage initialization to the fresh-network
   genesis state. The Pasteur upgrade route is not a substitute for genesis
   predeployment.
2. Confirm the correct first-block/fork activation behavior for a new chain;
   the current PULSE chain config still inherits the provisional BSC schedule.
3. Add integration tests for validator creation, direct self-delegation,
   self-redelegation, third-party delegation, governance parameter bounds,
   jailing below the minimum, and election ranking.
4. Run the Go client tests and a local multi-validator PULSE testnet, including
   state migration and upgrade checks. Solidity unit tests are not
   end-to-end validator tests.
5. Record reproducible source/runtime hashes and compare the runtime embedded
   in Go with the generated Foundry output.

The PULSE-specific runtime is now a committed artifact; it is not yet installed
into a finalized genesis state or validated on a running chain.
