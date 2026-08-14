package domain

import "time"

const (
	RewardDeliveryDiscount10   = "DELIVERY_DISCOUNT_10"
	RewardListingDiscount15    = "LISTING_DISCOUNT_15"
	RewardAutotekaDiscount20   = "AUTOTEKA_DISCOUNT_20"
	RewardFreeListingHighlight = "FREE_LISTING_HIGHLIGHT"
	RewardFreeListingPromotion = "FREE_LISTING_PROMOTION"
)

var rewardCodes = []string{
	RewardDeliveryDiscount10,
	RewardListingDiscount15,
	RewardAutotekaDiscount20,
	RewardFreeListingHighlight,
	RewardFreeListingPromotion,
}

type Reward struct {
	ID           string
	PromoCode    string
	Name         string
	Description  string
	Status       string
	ExpiresAt    string
	EarnedReason string
	RedeemedAt   *time.Time
}
