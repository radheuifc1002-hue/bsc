# PULSE CHAIN validator-system contract fork notes

Status: source audit in progress — no contract or bytecode changes made.

## Upstream StakeHub findings

The upstream `bnb-chain/bsc-genesis-contract` `contracts/StakeHub.sol`
currently defines these initialization values:

| StakeHub parameter | Upstream default |
|---|---:|
| `minSelfDelegationBNB` | 2,000 BNB |
| `minDelegationBNBChange` | 1 BNB |
| `maxElectedValidators` | 45 |
| `unbondPeriod` | 7 days |
| `redelegateFeeRate` | 2 |

Its governance parameter update path also constrains
`minSelfDelegationBNB` to a range of 1,000–100,000 BNB. Therefore changing
only the initialization default is not sufficient for PULSE CHAIN's requested
2,000,000 PULSE minimum: the governance validation bound must also be reviewed
and updated, otherwise the intended value cannot be set through that path.

Source reference: https://github.com/bnb-chain/bsc-genesis-contract/blob/master/contracts/StakeHub.sol

## PULSE CHAIN requirement

- Set the minimum self-delegation/validator eligibility threshold to
  2,000,000 PULSE (18-decimal base units in the EVM contract).
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
3. Change the minimum self-delegation default to 2,000,000 PULSE and update the
   governance parameter validation range so it accepts the intended value.
4. Preserve the 45 elected-validator capacity and all unrelated parameters.
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
