# PULSE CHAIN genesis specification

Status: draft — not a launchable genesis file.

This document records the decisions already made and the inputs that must be
fixed before a deterministic PULSE CHAIN genesis can be generated. Do not use
the BNB Smart Chain mainnet genesis block, allocation, validator extra-data, or
genesis hash as PULSE CHAIN genesis data.

## Network identity

| Parameter | Value | Status |
|---|---:|---|
| Network name | PULSE CHAIN | Confirmed |
| Native asset name | PULSE | Confirmed |
| Native asset symbol | PULSE | Confirmed |
| EVM chain ID | 676 (0x2A4) | Confirmed |
| Client network ID | 676 | Proposed default; must be wired consistently |
| Native asset decimals | 18 | Proposed EVM-compatible denomination |
| Consensus engine | Parlia (BSC-derived PoSA) | Confirmed |

Chain ID is replay-protection data. Network ID is a separate client/network
setting; both should be 676 for the initial public network unless a concrete
networking requirement establishes otherwise.

## Supply and allocation constraints

- Initial supply: 1,100,000,000 PULSE.
- Block-reward reserve: 110,000,000 PULSE (10% of initial supply), held in a
  dedicated, non-spendable-by-ordinary-accounts reserve until released by the
  consensus reward mechanism.
- Remaining supply: 990,000,000 PULSE (90%). Final recipient allocations and
  recipient addresses have not yet been supplied.
- Initial WPULSE liquidity allocation: 50,000,000 PULSE, included within the
  990,000,000 PULSE allocation; it is not additional supply.
- WPULSE is intended to be redeemable 1:1 for locked native PULSE.

The 128-year reward curve, six-month 5% reductions, exact initial per-block
emission, rounding policy, and reserve-release implementation must be specified
and tested before the reserve is placed in a live genesis allocation. No
illustrative reward rate is a genesis parameter.

## Consensus and validator constraints

- Preserve the upstream BSC Parlia consensus and validator-election design.
- Minimum validator self-stake target: 1,000,000 PULSE.
- The minimum self-stake is an eligibility threshold, not a guarantee of
  election.
- Do not alter validator counts, election cadence, delegation, commission,
  slashing, jailing, unbonding, reward, or governance rules as part of this
  threshold change.
- Upstream StakeHub currently sets `maxElectedValidators` to 45; preserve that
  elected-candidate capacity. Do not conflate it with the active Parlia
  block-producer set.
- The exact active producer selection and validator-set behavior must still
  be traced through the versioned validator system contract and Parlia code.
- A production genesis requires initial validator consensus keys, operator
  addresses, self-stake funding, and the initial validator set encoded in
  Parlia genesis extra-data.

The 1,000,000 PULSE threshold must be implemented in the applicable staking
system-contract source and its compiled/deployed bytecode. A Go-only constant
change is insufficient if the contract enforces the threshold.

## EVM fork schedule

A fresh PULSE network must not blindly reuse BSC mainnet's historical block
heights or timestamps. Doing so would cause PULSE to run with a different set
of EVM and BSC protocol rules at genesis than intended.

Before generating genesis, select the exact upstream BSC protocol revision to
inherit and define PULSE activation values for every supported block- and
time-based fork. Fork activation must be consistent across:
- the genesis chain configuration;
- Parlia consensus behavior;
- EVM execution and precompiles;
- versioned system-contract deployments and upgrades;
- RPC and transaction-pool behavior.

The currently added `PulseChainChainConfig` is an identity/configuration
scaffold copied from BSC. Its inherited BSC mainnet fork schedule is
provisional and must not be treated as the final fresh-network schedule.

## Genesis block fields still required

The deterministic genesis artifact still needs:
- launch timestamp;
- final fork schedule and chain configuration;
- genesis gas limit and other header parameters;
- initial Parlia validator keys and validator extra-data;
- final account allocation map summing exactly to 1.1 billion PULSE;
- system-contract code and storage for the selected fork revision;
- reward-reserve account/contract and its initialization state;
- genesis hash, generated from the finalized specification.

No arbitrary BSC validator extra-data, BSC genesis hash, or BSC mainnet
preallocation may be reused.

## Other protocol systems — not yet genesis-ready

The following are agreed design targets but need separate specifications and
audited implementation before being represented as genesis state:

- Base transaction fee target: 0.01 XU per 65,000 gas, converted to PULSE using
  a manipulation-resistant WPULSE/XU TWAP oracle.
- Base-fee distribution: 40% burn, 50% validators, 10% ecosystem.
- Priority fee: user-selected Gwei amount, with 100% burned.
- WPULSE wrapper and Uniswap V2-compatible DEX/pool.
- XU stablecoin, reserve composition, mint/burn permissions, and redemption.
- BNB Chain and Polygon bridges for USDT, USDC, and DAI.

The XU reserve, bridge collateral, treasury, reward reserve, and liquidity
allocation must be separately accounted for. Bridge collateral is not free
treasury balance.

## Required inputs before final genesis generation

1. Initial validator consensus public keys and operator addresses.
2. Allocation recipients and exact amounts for the 990,000,000 PULSE balance.
3. Initial launch timestamp and desired BSC fork baseline.
4. Genesis gas limit and any other network-level block parameters.
5. Final reward contract/reserve mechanism and exact emission arithmetic.
6. System-contract deployment manifest and storage initialization for PULSE.

Until these are finalized, this repository should not ship a default PULSE
genesis function or advertise a genesis hash.
