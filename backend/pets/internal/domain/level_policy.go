package domain

type LevelPolicy struct {
	configs map[int]string
}

func NewLevelPolicy() *LevelPolicy {
	return &LevelPolicy{
		configs: map[int]string{
			1: "",
			2: "DELIVERY_DISCOUNT_10",
			3: "LISTING_DISCOUNT_15",
			4: "AUTOTEKA_DISCOUNT_20",
			5: "FREE_LISTING_HIGHLIGHT",
			6: "FREE_LISTING_PROMOTION",
		},
	}
}

func (lp *LevelPolicy) GetCode(level int) string {
	if level == 1 {
		return ""
	}

	return lp.configs[(level-2)%(len(lp.configs)-1)+2]
}
