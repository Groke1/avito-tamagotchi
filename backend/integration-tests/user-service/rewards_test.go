package user_service

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

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

func TestCannotRedeemAnotherUsersReward(t *testing.T) {
	cfg := setup(t)
	seedRewardDefinitions(t, cfg)

	ownerSuffix := uniqueSuffix(t)
	ownerAuth := registerUser(
		t,
		cfg,
		fmt.Sprintf("reward-redeem-owner-%s", ownerSuffix),
		fmt.Sprintf("reward-redeem-owner-%s@example.com", ownerSuffix),
		testPassword,
	)
	ownerProfile := getProfile(t, cfg, ownerAuth.AccessToken)
	granted := grantReward(t, cfg, ownerProfile.UserID, rewardCodeDelivery)

	otherSuffix := uniqueSuffix(t)
	otherAuth := registerUser(
		t,
		cfg,
		fmt.Sprintf("reward-redeem-other-%s", otherSuffix),
		fmt.Sprintf("reward-redeem-other-%s@example.com", otherSuffix),
		testPassword,
	)

	resp := jsonReq(t, http.MethodPost, cfg.Users.APIURL+"/rewards/redeem", map[string]any{
		"promo_code": granted.PromoCode,
	}, otherAuth.AccessToken)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)

	apiErr := decodeBody[apiError](t, resp)
	require.Equal(t, "REWARD_UNAVAILABLE", apiErr.Code)

	ownerRewardResp := jsonReq(t, http.MethodGet, cfg.Users.APIURL+"/rewards/"+granted.RewardID, nil, ownerAuth.AccessToken)
	require.Equal(t, http.StatusOK, ownerRewardResp.StatusCode)
	ownerReward := decodeBody[userRewardResponse](t, ownerRewardResp)
	require.Equal(t, "active", ownerReward.Status)
	require.Nil(t, ownerReward.RedeemedAt)
}

func TestExpiredRewardCannotBeRedeemed(t *testing.T) {
	cfg := setup(t)
	seedRewardDefinitions(t, cfg)

	suffix := uniqueSuffix(t)
	auth := registerUser(
		t,
		cfg,
		fmt.Sprintf("reward-expired-%s", suffix),
		fmt.Sprintf("reward-expired-%s@example.com", suffix),
		testPassword,
	)
	profile := getProfile(t, cfg, auth.AccessToken)
	granted := grantReward(t, cfg, profile.UserID, rewardCodeDelivery)
	expireReward(t, cfg, granted.RewardID)

	resp := jsonReq(t, http.MethodPost, cfg.Users.APIURL+"/rewards/redeem", map[string]any{
		"promo_code": granted.PromoCode,
	}, auth.AccessToken)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)

	apiErr := decodeBody[apiError](t, resp)
	require.Equal(t, "REWARD_UNAVAILABLE", apiErr.Code)

	rewardResp := jsonReq(t, http.MethodGet, cfg.Users.APIURL+"/rewards/"+granted.RewardID, nil, auth.AccessToken)
	require.Equal(t, http.StatusOK, rewardResp.StatusCode)
	reward := decodeBody[userRewardResponse](t, rewardResp)
	require.Equal(t, "expired", reward.Status)
	require.NotNil(t, reward.ExpiresAt)
	require.True(t, reward.ExpiresAt.Before(time.Now()))
	require.Nil(t, reward.RedeemedAt)

	activeResp := jsonReq(t, http.MethodGet, cfg.Users.APIURL+"/rewards/active", nil, auth.AccessToken)
	require.Equal(t, http.StatusOK, activeResp.StatusCode)
	active := decodeBody[userRewardsResponse](t, activeResp)
	require.Empty(t, active.Items)

	allResp := jsonReq(t, http.MethodGet, cfg.Users.APIURL+"/rewards", nil, auth.AccessToken)
	require.Equal(t, http.StatusOK, allResp.StatusCode)
	all := decodeBody[userRewardsResponse](t, allResp)
	require.Len(t, all.Items, 1)
	require.Equal(t, "expired", all.Items[0].Status)
}

func TestGetUnknownReward(t *testing.T) {
	cfg := setup(t)

	suffix := uniqueSuffix(t)
	auth := registerUser(
		t,
		cfg,
		fmt.Sprintf("reward-unknown-%s", suffix),
		fmt.Sprintf("reward-unknown-%s@example.com", suffix),
		testPassword,
	)

	resp := jsonReq(t, http.MethodGet, cfg.Users.APIURL+"/rewards/00000000-0000-0000-0000-000000000000", nil, auth.AccessToken)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)

	apiErr := decodeBody[apiError](t, resp)
	require.Equal(t, "REWARD_NOT_FOUND", apiErr.Code)
}

func TestRewardsRequireAccessToken(t *testing.T) {
	cfg := setup(t)

	testCases := []struct {
		name   string
		method string
		path   string
		body   any
	}{
		{name: "all rewards", method: http.MethodGet, path: "/rewards"},
		{name: "active rewards", method: http.MethodGet, path: "/rewards/active"},
		{name: "single reward", method: http.MethodGet, path: "/rewards/00000000-0000-0000-0000-000000000000"},
		{name: "redeem reward", method: http.MethodPost, path: "/rewards/redeem", body: map[string]any{"promo_code": "PROMO"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp := jsonReq(t, tc.method, cfg.Users.APIURL+tc.path, tc.body, "")
			require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

			apiErr := decodeBody[apiError](t, resp)
			require.Equal(t, "UNAUTHORIZED", apiErr.Code)
		})
	}
}

func TestRedeemRewardValidationCases(t *testing.T) {
	cfg := setup(t)

	suffix := uniqueSuffix(t)
	auth := registerUser(
		t,
		cfg,
		fmt.Sprintf("reward-validation-%s", suffix),
		fmt.Sprintf("reward-validation-%s@example.com", suffix),
		testPassword,
	)

	testCases := []struct {
		name string
		body string
	}{
		{name: "empty body", body: `{}`},
		{name: "blank promo code", body: `{"promo_code":"   "}`},
		{name: "unknown field", body: `{"promo_code":"PROMO","extra":1}`},
		{name: "multiple json objects", body: `{"promo_code":"PROMO"}{"promo_code":"OTHER"}`},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp := rawReq(t, http.MethodPost, cfg.Users.APIURL+"/rewards/redeem", tc.body, auth.AccessToken)
			require.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)

			apiErr := decodeBody[apiError](t, resp)
			require.Equal(t, "VALIDATION_ERROR", apiErr.Code)
		})
	}
}
