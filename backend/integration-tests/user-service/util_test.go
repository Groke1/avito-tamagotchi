package user_service

import (
	"bytes"
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/caarlos0/env/v10"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

type config struct {
	db *sql.DB

	Users struct {
		BaseURL     string `env:"USERS_BASE_URL" envDefault:"http://users:8080"`
		APIURL      string
		InternalURL string
	}

	Pets struct {
		BaseURL string `env:"PETS_BASE_URL" envDefault:"http://pets:8080"`
		APIURL  string
	}

	PG struct {
		Host     string `env:"POSTGRES_HOST" envDefault:"localhost"`
		Port     string `env:"POSTGRES_PORT" envDefault:"5432"`
		DB       string `env:"POSTGRES_DB" envDefault:"db"`
		User     string `env:"POSTGRES_USER" envDefault:"db_user"`
		Password string `env:"POSTGRES_PASSWORD" envDefault:"12345"`
	}
}

type authResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type profileResponse struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Coins    uint64 `json:"coins"`
}

type updateCoinsResponse struct {
	UserID string `json:"user_id"`
	Coins  uint64 `json:"coins"`
}

type rewardDefinitionResponse struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type userRewardResponse struct {
	RewardID    string     `json:"reward_id"`
	PromoCode   string     `json:"promo_code"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	ExpiresAt   *time.Time `json:"expires_at"`
	RedeemedAt  *time.Time `json:"redeemed_at"`
}

type userRewardsResponse struct {
	Items []userRewardResponse `json:"items"`
}

type streakRecord struct {
	CurrentStreak  int
	LastActiveDate time.Time
}

type apiError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func setup(t *testing.T) *config {
	t.Helper()

	cfg := loadConfig(t)
	waitForUsers(t, cfg.Users.APIURL, 45*time.Second)
	require.NoError(t, cleanDB(t.Context(), cfg.db))
	seedRewardDefinitions(t, cfg)

	t.Cleanup(func() {
		require.NoError(t, cleanDB(context.Background(), cfg.db))
		require.NoError(t, cfg.db.Close())
	})

	return cfg
}

func loadConfig(t *testing.T) *config {
	t.Helper()

	var cfg config
	require.NoError(t, env.Parse(&cfg))

	cfg.Users.BaseURL = normalizeURL(cfg.Users.BaseURL)
	cfg.Users.APIURL = cfg.Users.BaseURL + "/api/v1"
	cfg.Users.InternalURL = cfg.Users.BaseURL + "/internal"
	cfg.Pets.BaseURL = normalizeURL(cfg.Pets.BaseURL)
	cfg.Pets.APIURL = cfg.Pets.BaseURL + "/api/v1"
	cfg.initDB(t)

	return &cfg
}

func (c *config) initDB(t *testing.T) {
	t.Helper()

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		url.QueryEscape(c.PG.User),
		url.QueryEscape(c.PG.Password),
		c.PG.Host,
		c.PG.Port,
		c.PG.DB,
	)

	db, err := sql.Open("postgres", dsn)
	require.NoError(t, err)
	require.NoError(t, db.PingContext(t.Context()))
	c.db = db
}

func normalizeURL(value string) string {
	if strings.HasPrefix(value, "http://") || strings.HasPrefix(value, "https://") {
		return strings.TrimRight(value, "/")
	}

	return "http://" + strings.TrimRight(value, "/")
}

func waitForUsers(t *testing.T, baseURL string, timeout time.Duration) {
	t.Helper()

	ctx, cancel := context.WithTimeout(t.Context(), timeout)
	defer cancel()

	client := &http.Client{Timeout: 2 * time.Second}
	payload := []byte(`{"email":"readiness@example.com","password":"invalid-password"}`)

	for {
		select {
		case <-ctx.Done():
			t.Fatalf("users service %s did not become ready in %v", baseURL, timeout)
		default:
		}

		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodPost,
			baseURL+"/auth/login",
			bytes.NewReader(payload),
		)
		if err != nil {
			time.Sleep(300 * time.Millisecond)
			continue
		}

		req.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			time.Sleep(300 * time.Millisecond)
			continue
		}

		_ = resp.Body.Close()
		if resp.StatusCode/100 == 2 || resp.StatusCode/100 == 4 {
			return
		}

		time.Sleep(300 * time.Millisecond)
	}
}

func jsonReq(
	t *testing.T,
	method string,
	requestURL string,
	body any,
	accessToken string,
) *http.Response {
	t.Helper()

	var payload io.Reader
	if body != nil {
		var buf bytes.Buffer
		require.NoError(t, json.NewEncoder(&buf).Encode(body))
		payload = &buf
	}

	req, err := http.NewRequestWithContext(t.Context(), method, requestURL, payload)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	if accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+accessToken)
	}

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	return resp
}

func rawReq(
	t *testing.T,
	method string,
	requestURL string,
	body string,
	accessToken string,
) *http.Response {
	t.Helper()

	req, err := http.NewRequestWithContext(t.Context(), method, requestURL, strings.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	if accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+accessToken)
	}

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	return resp
}

func decodeBody[T any](t *testing.T, resp *http.Response) T {
	t.Helper()
	defer resp.Body.Close()

	var result T
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
	return result
}

func registerUser(
	t *testing.T,
	cfg *config,
	username string,
	email string,
	password string,
) authResponse {
	t.Helper()

	resp := jsonReq(t, http.MethodPost, cfg.Users.APIURL+"/auth/register", map[string]any{
		"username": username,
		"email":    email,
		"password": password,
	}, "")
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	return decodeBody[authResponse](t, resp)
}

func uniqueSuffix(t *testing.T) string {
	t.Helper()

	buf := make([]byte, 6)
	_, err := rand.Read(buf)
	require.NoError(t, err)

	return hex.EncodeToString(buf)
}

func getProfile(t *testing.T, cfg *config, accessToken string) profileResponse {
	t.Helper()

	resp := jsonReq(t, http.MethodGet, cfg.Users.APIURL+"/profile", nil, accessToken)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	return decodeBody[profileResponse](t, resp)
}

func createPet(t *testing.T, cfg *config, accessToken string, name string) {
	t.Helper()

	resp := jsonReq(t, http.MethodPost, cfg.Pets.APIURL+"/pet/", map[string]any{
		"name": name,
	}, accessToken)
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	requireEmptyBodyOrClose(t, resp)
}

func requireEmptyBody(t *testing.T, resp *http.Response) {
	t.Helper()
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Empty(t, body)
}

func requireEmptyBodyOrClose(t *testing.T, resp *http.Response) {
	t.Helper()
	defer resp.Body.Close()
	_, err := io.Copy(io.Discard, resp.Body)
	require.NoError(t, err)
}

func grantReward(t *testing.T, cfg *config, userID string, code string) userRewardResponse {
	t.Helper()

	resp := jsonReq(t, http.MethodPost, cfg.Users.InternalURL+"/rewards", map[string]any{
		"user_id": userID,
		"code":    code,
	}, "")
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	return decodeBody[userRewardResponse](t, resp)
}

func getStreak(t *testing.T, cfg *config, userID string) streakRecord {
	t.Helper()

	var streak streakRecord
	err := cfg.db.QueryRowContext(
		t.Context(),
		`SELECT current_streak, last_active_date FROM users.user_streaks WHERE user_id = $1`,
		userID,
	).Scan(&streak.CurrentStreak, &streak.LastActiveDate)
	require.NoError(t, err)

	return streak
}

func seedRewardDefinitions(t *testing.T, cfg *config) {
	t.Helper()

	_, err := cfg.db.ExecContext(t.Context(), `
		INSERT INTO users.reward_definitions (code, name, description, earned_description, redeemed_description)
		VALUES
			('DELIVERY_DISCOUNT_10', 'Скидка 10% на доставку', 'Скидка 10% на одну покупку с Авито Доставкой', '', ''),
			('FREE_LISTING_PROMOTION', 'Бесплатное продвижение', 'Одно бесплатное продвижение выбранного объявления', '', ''),
			('AUTOTEKA_DISCOUNT_20', 'Скидка 20% на Автотеку', 'Скидка 20% на один отчёт об истории автомобиля', '', ''),
			('FREE_LISTING_HIGHLIGHT', 'Выделение объявления', 'Бесплатное визуальное выделение одного объявления на ограниченный срок', '', ''),
			('LISTING_DISCOUNT_15', 'Скидка 15% на размещение', 'Скидка 15% на одно платное размещение объявления', '', '')
		ON CONFLICT (code) DO NOTHING
	`)
	require.NoError(t, err)
}

func cleanDB(ctx context.Context, db *sql.DB) error {
	var tables sql.NullString

	err := db.QueryRowContext(ctx, `
		SELECT string_agg(
			quote_ident(schemaname) || '.' || quote_ident(tablename),
			', '
		)
		FROM pg_tables
		WHERE schemaname NOT IN ('pg_catalog', 'information_schema')
		  AND tablename <> 'goose_db_version'
	`).Scan(&tables)
	if err != nil {
		return err
	}

	if !tables.Valid || tables.String == "" {
		return nil
	}

	_, err = db.ExecContext(ctx, "TRUNCATE TABLE "+tables.String+" CASCADE")
	return err
}

func expireReward(t *testing.T, cfg *config, rewardID string) {
	t.Helper()

	result, err := cfg.db.ExecContext(
		t.Context(),
		`UPDATE users.user_rewards SET expires_at = NOW() - INTERVAL '1 minute' WHERE id = $1`,
		rewardID,
	)
	require.NoError(t, err)

	affected, err := result.RowsAffected()
	require.NoError(t, err)
	require.Equal(t, int64(1), affected)
}
