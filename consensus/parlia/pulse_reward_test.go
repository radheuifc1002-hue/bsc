package parlia

import (
	"math/big"
	"testing"

	"github.com/holiman/uint256"
)

func TestPulseRewardCurveSixMonthReduction(t *testing.T) {
	curve := newPulseRewardCurve()
	wantFirst := new(big.Int).Mul(curve.powers[0], big.NewInt(19))
	wantFirst.Div(wantFirst, big.NewInt(20))
	if curve.powers[1].Cmp(wantFirst) != 0 {
		t.Fatal("first six-month period does not start at 95%")
	}
	wantSecond := new(big.Int).Mul(wantFirst, big.NewInt(19))
	wantSecond.Div(wantSecond, big.NewInt(20))
	if curve.powers[2].Cmp(wantSecond) != 0 {
		t.Fatal("second six-month period does not compound the 5% reduction")
	}
}

func TestPulseRewardDeclinesLinearlyWithinSixMonths(t *testing.T) {
	curve := newPulseRewardCurve()
	half := pulseRewardSixMonthMilliseconds / 2
	got := curve.weightAt(half)
	d := new(big.Int).SetUint64(pulseRewardSixMonthMilliseconds)
	x := new(big.Int).SetUint64(half)
	numerator := new(big.Int).Mul(new(big.Int).Mul(big.NewInt(40), d), x)
	numerator.Sub(numerator, new(big.Int).Mul(x, x))
	denominator := new(big.Int).Mul(big.NewInt(40), d)
	want := new(big.Int).Mul(curve.powers[0], numerator)
	want.Div(want, denominator)
	if got.Cmp(want) != 0 {
		t.Fatalf("half-period integrated weight mismatch: got %s want %s", got, want)
	}
}

func TestPulseRewardReserveFullyEmitsAt128Years(t *testing.T) {
	curve := newPulseRewardCurve()
	reserve := uint256.MustFromBig(curve.reserveWei)
	start := uint64(1_800_000_000_000)
	got := curve.rewardForInterval(start, start, start+pulseRewardDurationMilliseconds, reserve)
	if got.Cmp(curve.reserveWei) != 0 {
		t.Fatalf("128-year emission did not release full reserve: got %s want %s", got, curve.reserveWei)
	}
}

func TestPulseRewardCumulativeTargetsConserveReserveAcrossGaps(t *testing.T) {
	curve := newPulseRewardCurve()
	reserve := uint256.MustFromBig(curve.reserveWei)
	start := uint64(1_800_000_000_000)
	first := curve.rewardForInterval(start, start, start+pulseRewardSixMonthMilliseconds, reserve)
	remaining := new(big.Int).Sub(curve.reserveWei, first)
	remainingU256 := uint256.MustFromBig(remaining)
	last := curve.rewardForInterval(start, start+pulseRewardSixMonthMilliseconds, start+pulseRewardDurationMilliseconds, remainingU256)
	total := new(big.Int).Add(first, last)
	if total.Cmp(curve.reserveWei) != 0 {
		t.Fatalf("emission conservation failed: got %s want %s", total, curve.reserveWei)
	}
}

func TestPulseRewardDoesNotEmitBeforeGenesisOrAfterReserveEmpty(t *testing.T) {
	curve := newPulseRewardCurve()
	reserve := uint256.MustFromBig(curve.reserveWei)
	if got := curve.rewardForInterval(100, 90, 100, reserve); got.Sign() != 0 {
		t.Fatalf("emitted before genesis: %s", got)
	}
	empty := new(uint256.Int)
	if got := curve.rewardForInterval(100, 100, 101, empty); got.Sign() != 0 {
		t.Fatalf("emitted with empty reserve: %s", got)
	}
}

func TestPulseRewardUsesMillisecondsForBlockIntervals(t *testing.T) {
	curve := newPulseRewardCurve()
	reserve := uint256.MustFromBig(curve.reserveWei)
	start := uint64(1_800_000_000_000)
	got := curve.rewardForInterval(start, start, start+750, reserve)
	want := big.NewInt(268130505097545275) // wei, first 750 ms interval
	if got.Cmp(want) != 0 {
		t.Fatalf("750ms initial reward mismatch: got %s want %s", got, want)
	}
}
