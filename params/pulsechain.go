// Copyright 2026 The PULSE CHAIN Authors
// This file is part of the PULSE CHAIN network configuration.
//
// PULSE CHAIN is derived from BNB Smart Chain. Its initial execution and
// Parlia fork schedule follows the upstream BSC configuration until PULSE
// CHAIN-specific genesis and fork activation parameters are finalized.

package params

import "math/big"

const PulseChainID uint64 = 676

// PulseChainChainConfig is the initial PULSE CHAIN configuration. It copies
// the BSC fork schedule while using PULSE CHAIN's independent replay-protection
// chain ID. The network's genesis allocation and genesis hash must be supplied
// by the PULSE CHAIN genesis specification; this config alone does not create
// or initialize a network.
var PulseChainChainConfig = func() *ChainConfig {
	config := *BSCChainConfig
	config.ChainID = new(big.Int).SetUint64(PulseChainID)
	return &config
}()

func init() {
	NetworkNames[PulseChainChainConfig.ChainID.String()] = "pulse-chain"
}
