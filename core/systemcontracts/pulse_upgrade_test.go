package systemcontracts

import (
	"encoding/hex"
	"strings"
	"testing"

	pulseContracts "github.com/ethereum/go-ethereum/core/systemcontracts/pulse"
)

func TestPulseStakeHubRuntimeArtifact(t *testing.T) {
	code, err := hex.DecodeString(strings.TrimSpace(pulseContracts.StakeHubContract))
	if err != nil {
		t.Fatalf("PULSE StakeHub runtime is not raw hex: %v", err)
	}
	if len(code) == 0 {
		t.Fatal("PULSE StakeHub runtime is empty")
	}
}

func TestPulsePasteurUpgradeUsesCustomStakeHub(t *testing.T) {
	upgrade := pasteurUpgrade[pulseNet]
	if upgrade == nil {
		t.Fatal("PULSE Pasteur upgrade is not registered")
	}
	if len(upgrade.Configs) != 1 {
		t.Fatalf("PULSE Pasteur upgrade has %d contracts, want only the custom StakeHub", len(upgrade.Configs))
	}
	if got := strings.TrimSpace(upgrade.Configs[0].Code); got != strings.TrimSpace(pulseContracts.StakeHubContract) {
		t.Fatal("PULSE Pasteur upgrade does not use the generated PULSE StakeHub runtime")
	}
}
