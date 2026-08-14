package domain

type LevelPolicy struct {
	configs map[int]string
}

func NewLevelPolicy() *LevelPolicy {
	return &LevelPolicy{
		configs: map[int]string{
			1: "",
			2: RewardDeliveryDiscount10,
			3: RewardListingDiscount15,
			4: RewardAutotekaDiscount20,
			5: RewardFreeListingHighlight,
			6: RewardFreeListingPromotion,
		},
	}
}

func (lp *LevelPolicy) GetCode(level int) string {
	if level == 1 {
		return ""
	}

	return lp.configs[(level-2)%(len(lp.configs)-1)+2]
}
