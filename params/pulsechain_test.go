package params

import (
	"math/big"
	"testing"
)

func TestPulseChainChainConfig(t *testing.T) {
	if got, want := PulseChainChainConfig.ChainID, big.NewInt(676); got.Cmp(want) != 0 {
		t.Fatalf("Pulse Chain ID = %v, want %v", got, want)
	}
	if got, want := BSCChainConfig.ChainID, big.NewInt(56); got.Cmp(want) != 0 {
		t.Fatalf("BSC Chain ID changed to %v, want %v", got, want)
	}
	if PulseChainChainConfig == BSCChainConfig {
		t.Fatal("Pulse Chain config must not alias the BSC config")
	}
	if PulseChainChainConfig.Parlia == nil {
		t.Fatal("Pulse Chain must retain the BSC Parlia consensus configuration")
	}
	if got := NetworkNames[PulseChainChainConfig.ChainID.String()]; got != "pulse-chain" {
		t.Fatalf("Pulse Chain network name = %q, want %q", got, "pulse-chain")
	}
}
