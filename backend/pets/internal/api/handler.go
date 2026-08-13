package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/cayman444/avito-gamification-hackathon/backend/pets/internal/domain"
	"github.com/cayman444/avito-gamification-hackathon/backend/pets/internal/service"

	"github.com/go-playground/validator/v10"
)

type PetHandler struct {
	petService  *service.PetService
	validator   *validator.Validate
	tripService *service.TripService
}

func NewPetHandler(service *service.PetService, tripService *service.TripService) *PetHandler {
	return &PetHandler{
		petService:  service,
		tripService: tripService,
		validator:   validator.New(),
	}
}

func (ph *PetHandler) GetPet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, err := GetUserIDFromContext(ctx)
	if err != nil {
		writeError(w, ErrUnauthorized)
		return
	}

	pet, err := ph.petService.GetPet(ctx, userID)
	if err != nil {
		writeError(w, ErrPetNotFound)
		return
	}

	petResponse := PetResponse{
		ID:          pet.ID,
		Name:        pet.Name,
		Level:       pet.Level,
		XP:          pet.XP,
		NextLevelXP: pet.NextLevelXP,
		Satiety:     pet.Satiety,
		Happiness:   pet.Happiness,
	}

	writeJSONResponse(w, http.StatusOK, petResponse)
}

func (ph *PetHandler) CreatePet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, err := GetUserIDFromContext(ctx)
	if err != nil {
		writeError(w, ErrUnauthorized)
		return
	}

	var req CreatePetRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, ErrValidationError)
		return
	} else if err = ph.validator.Struct(req); err != nil {
		writeError(w, ErrValidationError)
		return
	}

	pet, err := ph.petService.CreatePet(ctx, userID, req.Name)
	if err != nil {
		writeError(w, ErrPetAlreadyExists)
		return
	}

	petResponse := PetResponse{
		ID:          pet.ID,
		Name:        pet.Name,
		Level:       pet.Level,
		XP:          pet.XP,
		NextLevelXP: pet.NextLevelXP,
		Satiety:     pet.Satiety,
		Happiness:   pet.Happiness,
	}

	writeJSONResponse(w, http.StatusCreated, petResponse)
}

func (ph *PetHandler) FeedPet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, err := GetUserIDFromContext(ctx)
	if err != nil {
		writeError(w, ErrUnauthorized)
		return
	}

	var actionUnavailableError *domain.ActionUnavailableError

	pet, err := ph.petService.FeedPet(ctx, userID)
	if errors.Is(err, domain.ErrPetNotFound) {
		writeError(w, ErrPetNotFound)
		return
	} else if errors.As(err, &actionUnavailableError) {
		writeErrorWithRetryAfter(w, ErrUnavailableAction, int(actionUnavailableError.RetryAfter.Seconds()))
		return
	}

	petResponse := PetResponse{
		ID:          pet.ID,
		Name:        pet.Name,
		Level:       pet.Level,
		XP:          pet.XP,
		NextLevelXP: pet.NextLevelXP,
		Satiety:     pet.Satiety,
		Happiness:   pet.Happiness,
	}

	writeJSONResponse(w, http.StatusOK, petResponse)
}

func (ph *PetHandler) StrokePet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, err := GetUserIDFromContext(ctx)
	if err != nil {
		writeError(w, ErrUnauthorized)
		return
	}

	var actionUnavailableError *domain.ActionUnavailableError

	pet, err := ph.petService.StrokePet(ctx, userID)
	if errors.Is(err, domain.ErrPetNotFound) {
		writeError(w, ErrPetNotFound)
		return
	} else if errors.As(err, &actionUnavailableError) {
		writeErrorWithRetryAfter(w, ErrUnavailableAction, int(actionUnavailableError.RetryAfter.Seconds()))
		return
	}

	petResponse := PetResponse{
		ID:          pet.ID,
		Name:        pet.Name,
		Level:       pet.Level,
		XP:          pet.XP,
		NextLevelXP: pet.NextLevelXP,
		Satiety:     pet.Satiety,
		Happiness:   pet.Happiness,
	}

	writeJSONResponse(w, http.StatusOK, petResponse)
}

func (ph *PetHandler) GetLeaderboard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, err := GetUserIDFromContext(ctx)
	if err != nil {
		writeError(w, ErrUnauthorized)
		return
	}

	limit := 20
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		//nolint:govet // чтобы не жаловался
		l, err := strconv.Atoi(limitStr)
		if err == nil && l >= 1 && l <= 100 {
			limit = l
		}
	}

	leaderboardItems, currentUser, err := ph.petService.GetLeaderboard(ctx, userID, limit)
	if err != nil {
		log.Printf("[GET LEADERBOARD] %v", err)
		writeError(w, ErrInternalError)
		return
	}

	leaderboard := make([]LeaderboardItemResponse, len(leaderboardItems))
	for i, item := range leaderboardItems {
		leaderboard[i] = LeaderboardItemResponse{
			Rank:     item.Rank,
			Level:    item.Level,
			UserName: item.UserName,
			PetName:  item.PetName,
			XP:       item.XP,
		}
	}

	leaderboardResponse := LeaderboardResponse{
		Items: leaderboard,
		CurrentUser: LeaderboardItemResponse{
			Rank:     currentUser.Rank,
			Level:    currentUser.Level,
			UserName: currentUser.UserName,
			PetName:  currentUser.PetName,
			XP:       currentUser.XP,
		},
	}

	writeJSONResponse(w, http.StatusOK, leaderboardResponse)
}

func (ph *PetHandler) DailyBonusForStreak(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req BonusXpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {

		writeError(w, ErrValidationError)
		return
	}

	err := ph.petService.ClaimDailyBonusForStreak(ctx, req.UserID, req.Streak, req.Coins)
	if err != nil {
		log.Printf("[DAILY BONUS] %v", err)
		writeError(w, ErrInternalError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (ph *PetHandler) UpdateXP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req UpdateXPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("[UPDATED XP] %v", err)
		writeError(w, ErrValidationError)
		return
	}

	_, err := ph.petService.GrantXP(ctx, req.UserID, req.XP)
	if err != nil {
		log.Printf("[UPDATED XP] %v", err)
		writeError(w, ErrInternalError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (ph *PetHandler) DailyGainedXP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req struct {
		UserID string `json:"user_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, ErrValidationError)
		return
	}

	gainedXP, err := ph.petService.ClaimDailyGainedXP(ctx, req.UserID)
	if err != nil {
		log.Printf("[DAILY GAINED XP] %v", err)
		writeError(w, ErrInternalError)
		return
	}

	resp := struct {
		DailyGainedXP int `json:"daily_gained_xp"`
	}{
		DailyGainedXP: gainedXP,
	}
	writeJSONResponse(w, http.StatusOK, resp)
}

func (ph *PetHandler) GetNextRewardDescription(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, err := GetUserIDFromContext(ctx)
	if err != nil {
		writeError(w, ErrUnauthorized)
		return
	}

	reward, err := ph.petService.GetNextRewardDescription(ctx, userID)
	if err != nil {
		log.Printf("[NEXT REWARD DESCRIPTION] %v", err)
		writeError(w, ErrInternalError)
		return
	}

	rewardResponse := RewardDescriptionResponse{
		Code:        reward.PromoCode,
		Name:        reward.Name,
		Description: reward.Description,
	}

	writeJSONResponse(w, http.StatusOK, rewardResponse)
}

func (ph *PetHandler) GetWeeklyLeaderboard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, err := GetUserIDFromContext(ctx)
	if err != nil {
		writeError(w, ErrUnauthorized)
		return
	}

	limit := 20
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		//nolint:govet // чтобы не жаловался
		l, err := strconv.Atoi(limitStr)
		if err == nil && l >= 1 && l <= 100 {
			limit = l
		}
	}

	leaderboardItems, currentUser, err := ph.petService.GetWeeklyLeaderboard(ctx, userID, limit)
	if err != nil {
		log.Printf("[GET WEEKLY LEADERBOARD] %v", err)
		writeError(w, ErrInternalError)
		return
	}

	leaderboard := make([]LeaderboardItemResponse, len(leaderboardItems))
	for i, item := range leaderboardItems {
		leaderboard[i] = LeaderboardItemResponse{
			Rank:     item.Rank,
			Level:    item.Level,
			UserName: item.UserName,
			PetName:  item.PetName,
			XP:       item.XP,
		}
	}

	leaderboardResponse := LeaderboardResponse{
		Items: leaderboard,
		CurrentUser: LeaderboardItemResponse{
			Rank:     currentUser.Rank,
			Level:    currentUser.Level,
			UserName: currentUser.UserName,
			PetName:  currentUser.PetName,
			XP:       currentUser.XP,
		},
	}

	writeJSONResponse(w, http.StatusOK, leaderboardResponse)
}

func (ph *PetHandler) MakeTrip(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := CreateTripRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)

	userID, err := GetUserIDFromContext(ctx)
	if err != nil {
		writeError(w, ErrUnauthorized)
		return
	}
	tripDto := service.TripDTO{
		PetID:  req.PetID,
		UserID: userID,
	}
	err = ph.tripService.Generate(ctx, tripDto)
	if err != nil {
		writeError(w, ErrTripGenerationError)
		return
	}

	writeJSONResponse(w, http.StatusAccepted, &TripResponse{
		Status: TripStatusStarted,
	})
}
