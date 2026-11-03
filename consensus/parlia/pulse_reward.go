package parlia

import (
	"math/big"

	"github.com/ethereum/go-ethereum/params"
	"github.com/holiman/uint256"
)

// PULSE block rewards are released from a genesis-funded reserve. The curve is
// time-based so it remains deterministic even when Parlia block intervals
// change across protocol forks or blocks are delayed.
const (
	pulseRewardYears                 = uint64(128)
	pulseRewardYearSeconds           = uint64(31_557_600) // Julian year: 365.25 days
	pulseRewardSixMonthSeconds       = pulseRewardYearSeconds / 2
	pulseRewardDurationSeconds       = pulseRewardYears * pulseRewardYearSeconds
	pulseRewardPeriods               = int(pulseRewardYears * 2)
	pulseRewardReserveWholePULSE     = int64(110_000_000)
	pulseRewardLinearDropNumerator   = int64(1)
	pulseRewardLinearDropDenominator = int64(20) // 5% over each six-month period
	pulseRewardPrecision             = int64(36)
)

type pulseRewardCurve struct {
	powers       []*big.Int
	cumulative   []*big.Int
	totalWeight  *big.Int
	reserveWei   *big.Int
	periodLength *big.Int
}

var pulseEmissionCurve = newPulseRewardCurve()

func newPulseRewardCurve() *pulseRewardCurve {
	scale := new(big.Int).Exp(big.NewInt(10), big.NewInt(pulseRewardPrecision), nil)
	curve := &pulseRewardCurve{
		powers:       make([]*big.Int, pulseRewardPeriods+1),
		cumulative:   make([]*big.Int, pulseRewardPeriods+1),
		totalWeight:  new(big.Int),
		reserveWei:   new(big.Int).Mul(big.NewInt(pulseRewardReserveWholePULSE), new(big.Int).SetUint64(params.Ether)),
		periodLength: new(big.Int).SetUint64(pulseRewardSixMonthSeconds),
	}
	curve.powers[0] = scale
	curve.cumulative[0] = new(big.Int)
	periodWeightDenominator := big.NewInt(pulseRewardLinearDropDenominator * 2)
	for i := 0; i < pulseRewardPeriods; i++ {
		// The rate falls linearly from 100% to 95% during this period.
		// Its average multiplier is therefore 97.5% of the period's start.
		periodWeight := new(big.Int).Mul(curve.powers[i], curve.periodLength)
		periodWeight.Mul(periodWeight, big.NewInt(39))
		periodWeight.Div(periodWeight, periodWeightDenominator)
		curve.cumulative[i+1] = new(big.Int).Add(curve.cumulative[i], periodWeight)

		// The next period starts at 95% of this period's starting rate.
		curve.powers[i+1] = new(big.Int).Mul(curve.powers[i], big.NewInt(19))
		curve.powers[i+1].Div(curve.powers[i+1], big.NewInt(20))
	}
	curve.totalWeight.Set(curve.cumulative[pulseRewardPeriods])
	return curve
}

// weightAt returns the integrated, fixed-point emission weight from genesis
// through elapsedSeconds. Each six-month segment is linear internally and
// starts 5% below the previous segment's starting rate.
func (curve *pulseRewardCurve) weightAt(elapsedSeconds uint64) *big.Int {
	if elapsedSeconds >= pulseRewardDurationSeconds {
		return new(big.Int).Set(curve.totalWeight)
	}
	period := int(elapsedSeconds / pulseRewardSixMonthSeconds)
	offset := new(big.Int).SetUint64(elapsedSeconds % pulseRewardSixMonthSeconds)

	weight := new(big.Int).Set(curve.cumulative[period])
	if offset.Sign() == 0 {
		return weight
	}

	// Integral of (1 - 0.05*x/D) from 0 to offset:
	// offset - offset^2/(40*D). Keep the division until after multiplying
	// by the period's fixed-point starting-rate multiplier.
	d := curve.periodLength
	numerator := new(big.Int).Mul(new(big.Int).Mul(big.NewInt(40), d), offset)
	numerator.Sub(numerator, new(big.Int).Mul(offset, offset))
	denominator := new(big.Int).Mul(big.NewInt(40), d)
	partial := new(big.Int).Mul(curve.powers[period], numerator)
	partial.Div(partial, denominator)
	return weight.Add(weight, partial)
}

// rewardForInterval returns the reserve-funded reward attributable to the
// timestamp interval (parentTime, blockTime]. Cumulative targets make rewards
// independent of block count and prevent timestamp gaps from losing emissions.
func (curve *pulseRewardCurve) rewardForInterval(genesisTime, parentTime, blockTime uint64, reserveBalance *uint256.Int) *big.Int {
	if blockTime <= genesisTime || blockTime <= parentTime || reserveBalance == nil || reserveBalance.IsZero() {
		return new(big.Int)
	}
	elapsed := func(timestamp uint64) uint64 {
		if timestamp <= genesisTime {
			return 0
		}
		if timestamp-genesisTime >= pulseRewardDurationSeconds {
			return pulseRewardDurationSeconds
		}
		return timestamp - genesisTime
	}
	from := elapsed(parentTime)
	to := elapsed(blockTime)
	if to <= from {
		return new(big.Int)
	}
	// Release all remaining dust at the exact end of the 128-year schedule.
	if to == pulseRewardDurationSeconds && from < pulseRewardDurationSeconds {
		return reserveBalance.ToBig()
	}

	deltaWeight := new(big.Int).Sub(curve.weightAt(to), curve.weightAt(from))
	if deltaWeight.Sign() <= 0 {
		return new(big.Int)
	}
	reward := new(big.Int).Mul(curve.reserveWei, deltaWeight)
	reward.Div(reward, curve.totalWeight)
	if reward.Cmp(reserveBalance.ToBig()) > 0 {
		return reserveBalance.ToBig()
	}
	return reward
}
