package user_service

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	rewardCodeDelivery = "DELIVERY_DISCOUNT_10"
	rewardCodeListing  = "FREE_LISTING_PROMOTION"
	rewardCodeAutoteka = "AUTOTEKA_DISCOUNT_20"
)

func TestRewardLifecycle(t *testing.T) {
	cfg := setup(t)
	seedRewardDefinitions(t, cfg)

	suffix := uniqueSuffix(t)
	auth := registerUser(
		t,
		cfg,
		fmt.Sprintf("reward-user-%s", suffix),
		fmt.Sprintf("reward-%s@example.com", suffix),
		testPassword,
	)
	profile := getProfile(t, cfg, auth.AccessToken)

	definitionResp := jsonReq(t, http.MethodGet, cfg.Users.InternalURL+"/rewards/"+rewardCodeDelivery, nil, "")
	require.Equal(t, http.StatusOK, definitionResp.StatusCode)
	definition := decodeBody[rewardDefinitionResponse](t, definitionResp)
	require.Equal(t, rewardCodeDelivery, definition.Code)
	require.NotEmpty(t, definition.Name)
	require.NotEmpty(t, definition.Description)

	rewardsResp := jsonReq(t, http.MethodGet, cfg.Users.APIURL+"/rewards", nil, auth.AccessToken)
	require.Equal(t, http.StatusOK, rewardsResp.StatusCode)
	rewardsBefore := decodeBody[userRewardsResponse](t, rewardsResp)
	require.Empty(t, rewardsBefore.Items)

	activeResp := jsonReq(t, http.MethodGet, cfg.Users.APIURL+"/rewards/active", nil, auth.AccessToken)
	require.Equal(t, http.StatusOK, activeResp.StatusCode)
	activeBefore := decodeBody[userRewardsResponse](t, activeResp)
	require.Empty(t, activeBefore.Items)

	firstReward := grantReward(t, cfg, profile.UserID, rewardCodeDelivery)
	secondReward := grantReward(t, cfg, profile.UserID, rewardCodeListing)
	thirdReward := grantReward(t, cfg, profile.UserID, rewardCodeDelivery)

	require.Equal(t, "active", firstReward.Status)
	require.Equal(t, "active", secondReward.Status)
	require.Equal(t, "active", thirdReward.Status)
	require.NotEqual(t, firstReward.RewardID, thirdReward.RewardID)
	require.NotEqual(t, firstReward.PromoCode, thirdReward.PromoCode)
	require.True(t, strings.HasPrefix(firstReward.PromoCode, rewardCodeDelivery+"-"))
	require.True(t, strings.HasPrefix(secondReward.PromoCode, rewardCodeListing+"-"))
	require.True(t, strings.HasPrefix(thirdReward.PromoCode, rewardCodeDelivery+"-"))

	allRewardsResp := jsonReq(t, http.MethodGet, cfg.Users.APIURL+"/rewards", nil, auth.AccessToken)
	require.Equal(t, http.StatusOK, allRewardsResp.StatusCode)
	allRewards := decodeBody[userRewardsResponse](t, allRewardsResp)
	require.Len(t, allRewards.Items, 3)

	rewardsByID := make(map[string]userRewardResponse, len(allRewards.Items))
	for _, reward := range allRewards.Items {
		rewardsByID[reward.RewardID] = reward
	}

	require.Equal(t, "active", rewardsByID[firstReward.RewardID].Status)
	require.Equal(t, "active", rewardsByID[secondReward.RewardID].Status)
	require.Equal(t, "active", rewardsByID[thirdReward.RewardID].Status)

	activeAfterGrantResp := jsonReq(t, http.MethodGet, cfg.Users.APIURL+"/rewards/active", nil, auth.AccessToken)
	require.Equal(t, http.StatusOK, activeAfterGrantResp.StatusCode)
	activeAfterGrant := decodeBody[userRewardsResponse](t, activeAfterGrantResp)
	require.Len(t, activeAfterGrant.Items, 3)

	singleRewardResp := jsonReq(t, http.MethodGet, cfg.Users.APIURL+"/rewards/"+firstReward.RewardID, nil, auth.AccessToken)
	require.Equal(t, http.StatusOK, singleRewardResp.StatusCode)
	singleReward := decodeBody[userRewardResponse](t, singleRewardResp)
	require.Equal(t, firstReward.RewardID, singleReward.RewardID)
	require.Equal(t, firstReward.PromoCode, singleReward.PromoCode)
	require.Equal(t, definition.Name, singleReward.Name)
	require.Equal(t, definition.Description, singleReward.Description)
	require.Nil(t, singleReward.RedeemedAt)

	redeemResp := jsonReq(t, http.MethodPost, cfg.Users.APIURL+"/rewards/redeem", map[string]any{
		"promo_code": firstReward.PromoCode,
	}, auth.AccessToken)
	require.Equal(t, http.StatusNoContent, redeemResp.StatusCode)
	requireEmptyBody(t, redeemResp)

	redeemedRewardResp := jsonReq(t, http.MethodGet, cfg.Users.APIURL+"/rewards/"+firstReward.RewardID, nil, auth.AccessToken)
	require.Equal(t, http.StatusOK, redeemedRewardResp.StatusCode)
	redeemedReward := decodeBody[userRewardResponse](t, redeemedRewardResp)
	require.Equal(t, "redeemed", redeemedReward.Status)
	require.NotNil(t, redeemedReward.RedeemedAt)

	activeAfterRedeemResp := jsonReq(t, http.MethodGet, cfg.Users.APIURL+"/rewards/active", nil, auth.AccessToken)
	require.Equal(t, http.StatusOK, activeAfterRedeemResp.StatusCode)
	activeAfterRedeem := decodeBody[userRewardsResponse](t, activeAfterRedeemResp)
	require.Len(t, activeAfterRedeem.Items, 2)

	for _, reward := range activeAfterRedeem.Items {
		require.NotEqual(t, firstReward.RewardID, reward.RewardID)
		require.Equal(t, "active", reward.Status)
	}

	allAfterRedeemResp := jsonReq(t, http.MethodGet, cfg.Users.APIURL+"/rewards", nil, auth.AccessToken)
	require.Equal(t, http.StatusOK, allAfterRedeemResp.StatusCode)
	allAfterRedeem := decodeBody[userRewardsResponse](t, allAfterRedeemResp)
	require.Len(t, allAfterRedeem.Items, 3)

	statusByID := make(map[string]string, len(allAfterRedeem.Items))
	for _, reward := range allAfterRedeem.Items {
		statusByID[reward.RewardID] = reward.Status
	}
	require.Equal(t, "redeemed", statusByID[firstReward.RewardID])
	require.Equal(t, "active", statusByID[secondReward.RewardID])
	require.Equal(t, "active", statusByID[thirdReward.RewardID])
}

func TestRewardsAreIsolatedPerUser(t *testing.T) {
	cfg := setup(t)
	seedRewardDefinitions(t, cfg)

	firstSuffix := uniqueSuffix(t)
	firstAuth := registerUser(
		t,
		cfg,
		fmt.Sprintf("reward-owner-%s", firstSuffix),
		fmt.Sprintf("reward-owner-%s@example.com", firstSuffix),
		testPassword,
	)
	firstProfile := getProfile(t, cfg, firstAuth.AccessToken)
	granted := grantReward(t, cfg, firstProfile.UserID, rewardCodeAutoteka)

	secondSuffix := uniqueSuffix(t)
	secondAuth := registerUser(
		t,
		cfg,
		fmt.Sprintf("reward-other-%s", secondSuffix),
		fmt.Sprintf("reward-other-%s@example.com", secondSuffix),
		testPassword,
	)

	resp := jsonReq(t, http.MethodGet, cfg.Users.APIURL+"/rewards/"+granted.RewardID, nil, secondAuth.AccessToken)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)

	apiErr := decodeBody[apiError](t, resp)
	require.Equal(t, "REWARD_NOT_FOUND", apiErr.Code)
}

func TestRewardStatuses(t *testing.T) {
	cfg := setup(t)
	seedRewardDefinitions(t, cfg)

	suffix := uniqueSuffix(t)
	auth := registerUser(
		t,
		cfg,
		fmt.Sprintf("reward-status-%s", suffix),
		fmt.Sprintf("reward-status-%s@example.com", suffix),
		testPassword,
	)
	profile := getProfile(t, cfg, auth.AccessToken)
	granted := grantReward(t, cfg, profile.UserID, rewardCodeDelivery)

	t.Run("unknown definition", func(t *testing.T) {
		resp := jsonReq(t, http.MethodGet, cfg.Users.InternalURL+"/rewards/UNKNOWN_REWARD_CODE", nil, "")
		require.Equal(t, http.StatusNotFound, resp.StatusCode)

		apiErr := decodeBody[apiError](t, resp)
		require.Equal(t, "NOT_FOUND_DEFINITION", apiErr.Code)
	})

	t.Run("redeem unknown promo code", func(t *testing.T) {
		resp := jsonReq(t, http.MethodPost, cfg.Users.APIURL+"/rewards/redeem", map[string]any{
			"promo_code": "UNKNOWN-PROMO-CODE",
		}, auth.AccessToken)
		require.Equal(t, http.StatusNotFound, resp.StatusCode)

		apiErr := decodeBody[apiError](t, resp)
		require.Equal(t, "REWARD_UNAVAILABLE", apiErr.Code)
	})

	t.Run("redeem same reward twice", func(t *testing.T) {
		firstRedeemResp := jsonReq(t, http.MethodPost, cfg.Users.APIURL+"/rewards/redeem", map[string]any{
			"promo_code": granted.PromoCode,
		}, auth.AccessToken)
		require.Equal(t, http.StatusNoContent, firstRedeemResp.StatusCode)
		requireEmptyBody(t, firstRedeemResp)

		secondRedeemResp := jsonReq(t, http.MethodPost, cfg.Users.APIURL+"/rewards/redeem", map[string]any{
			"promo_code": granted.PromoCode,
		}, auth.AccessToken)
		require.Equal(t, http.StatusNotFound, secondRedeemResp.StatusCode)

		apiErr := decodeBody[apiError](t, secondRedeemResp)
		require.Equal(t, "REWARD_UNAVAILABLE", apiErr.Code)
	})
}
