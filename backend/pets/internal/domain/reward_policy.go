package domain

import "math/rand/v2"

const (
	negativeXPMin    = 1
	negativeXPMax    = 5
	negativeCoinsMin = 1
	negativeCoinsMax = 3

	positiveXPMin    = 30
	positiveXPMax    = 100
	positiveCoinsMin = 20
	positiveCoinsMax = 50

	jackpotReward = 20
)

type RewardPolicy struct {
	codes []string
}

func NewRewardPolicy() *RewardPolicy {
	return &RewardPolicy{
		codes: rewardCodes,
	}
}

func (rp *RewardPolicy) GetReward() string {
	return rp.codes[rand.IntN(len(rp.codes))]
}

func (rp *RewardPolicy) GenerateReward(isNegative bool) JourneyReward {
	if isNegative {
		return JourneyReward{
			RewardXP:    int32(rand.IntN(negativeXPMax-negativeXPMin+1) + negativeXPMin),
			RewardCoins: int32(rand.IntN(negativeCoinsMax-negativeCoinsMin+1) + negativeCoinsMin),
			RewardCode:  nil,
		}
	}

	reward := JourneyReward{
		RewardXP:    int32(rand.IntN(positiveXPMax-positiveXPMin+1) + positiveXPMin),
		RewardCoins: int32(rand.IntN(positiveCoinsMax-positiveCoinsMin+1) + positiveCoinsMin),
		RewardCode:  nil,
	}

	//nolint:mnd // шанс успеха
	if rand.IntN(100) < jackpotReward {
		code := rp.GetReward()
		reward.RewardCode = &code
	}

	return reward
}
