package gigachat

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/cayman444/avito-gamification-hackathon/backend/pets/internal/domain"
)

// newTestServer поднимает мок GigaChat: /oauth и /chat/completions на одном
// httptest.Server (у реального GigaChat это разные хосты, но для клиента
// это не важно — оба URL берутся из Config). oauthHandler/completionsHandler
// позволяют настроить поведение под конкретный тест-кейс.
func newTestServer(
	t *testing.T,
	oauthHandler,
	completionsHandler http.HandlerFunc,
) *httptest.Server {
	t.Helper()

	mux := http.NewServeMux()

	mux.HandleFunc("/oauth", func(w http.ResponseWriter, r *http.Request) {
		t.Logf(
			"[MOCK GIGACHAT] %s %s",
			r.Method,
			r.URL.Path,
		)

		t.Logf(
			"[MOCK GIGACHAT] scope=%q",
			r.FormValue("scope"),
		)

		oauthHandler(w, r)
	})

	mux.HandleFunc("/chat/completions", func(w http.ResponseWriter, r *http.Request) {
		t.Logf(
			"[MOCK GIGACHAT] %s %s",
			r.Method,
			r.URL.Path,
		)

		completionsHandler(w, r)
	})

	srv := httptest.NewServer(mux)

	t.Logf(
		"[MOCK GIGACHAT] server started: %s",
		srv.URL,
	)

	t.Cleanup(func() {
		t.Logf(
			"[MOCK GIGACHAT] server stopped: %s",
			srv.URL,
		)

		srv.Close()
	})

	return srv
}

func testConfig(srv *httptest.Server) *Config {
	return &Config{
		AuthKey:        "test-auth-key",
		Scope:          "GIGACHAT_API_PERS",
		Model:          "GigaChat",
		OAuthURL:       srv.URL + "/oauth",
		APIBaseURL:     srv.URL,
		RequestTimeout: 3 * time.Second,
	}
}

func okOAuthHandler(w http.ResponseWriter, _ *http.Request) {
	resp := oauthTokenResponse{
		AccessToken: "test-token",
		ExpiresAt:   time.Now().Add(time.Hour).UnixMilli(),
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

// completionsWithContent возвращает хэндлер, который всегда отвечает 200
// с переданным содержимым message.content — удобно для проверки парсинга.
func completionsWithContent(
	t *testing.T,
	content string,
) http.HandlerFunc {
	t.Helper()

	return func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		t.Logf(
			"[MOCK GIGACHAT] responding with content: %s",
			content,
		)

		resp := chatCompletionResponse{
			Choices: []struct {
				Message chatMessage `json:"message"`
			}{
				{
					Message: chatMessage{
						Role:    "assistant",
						Content: content,
					},
				},
			},
		}

		w.Header().Set(
			"Content-Type",
			"application/json",
		)

		_ = json.NewEncoder(w).Encode(resp)
	}
}

func sampleInput() domain.JourneyGenerationInput {
	return domain.JourneyGenerationInput{
		Journey: domain.JourneyResult{
			Location: "Электронный квартал",
			Events: []string{
				"встретил кота-продавца",
				"заблудился среди коробок",
				"нашёл кассетный плеер",
			},
			Reward: domain.JourneyReward{Coins: 50, Item: "retro_player"},
		},
		Memory: domain.PetMemory{
			Personality: "любопытный и немного трусливый",
			Summary:     "питомец уже несколько раз путешествовал с котом Барсиком",
		},
		RecentStories: []domain.JourneyStory{
			{Title: "День 1", Story: "Я встретил Барсика", Teaser: "..."},
		},
	}
}

const authKey = ""
const clientSecret = ""
const clientID = ""

func TestClient_Generate_Integration(t *testing.T) {
	//if os.Getenv("GIGACHAT_INTEGRATION_TEST") != "1" {
	//	t.Skip("set GIGACHAT_INTEGRATION_TEST=1 to run real GigaChat request")
	//}

	//if authKey == "" {
	//	t.Fatal("GIGACHAT_AUTH_KEY is required")
	//}

	var authKeyReal = base64.StdEncoding.EncodeToString(
		[]byte(clientID + ":" + clientSecret),
	)
	cfg := &Config{
		AuthKey:        authKeyReal,
		Scope:          "GIGACHAT_API_PERS",
		Model:          "GigaChat-2",
		OAuthURL:       "https://ngw.devices.sberbank.ru:9443/api/v2/oauth",
		APIBaseURL:     "https://api.giga.chat/v1",
		RequestTimeout: 30 * time.Second,
	}

	client := NewClient(cfg)

	t.Logf("[REAL GIGACHAT] sending request to %s", cfg.APIBaseURL)

	story, err := client.Generate(
		context.Background(),
		sampleInput(),
	)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	t.Logf(
		"[REAL GIGACHAT RESPONSE]\nTITLE: %s\nSTORY: %s\nTEASER: %s",
		story.Title,
		story.Story,
		story.Teaser,
	)
}

// --- happy path ---

func TestClient_Generate_Success(t *testing.T) {
	wantStory := domain.JourneyStory{
		Title:  "Я вернулся из Электронного квартала!",
		Story:  "Я встретил очень серьёзного кота-продавца и нашёл кассетный плеер.",
		Teaser: "Кажется, рядом осталось ещё одно неизведанное место 👀",
	}
	content, _ := json.Marshal(wantStory)

	srv := newTestServer(t, okOAuthHandler, completionsWithContent(t, string(content)))
	client := NewClient(testConfig(srv))

	got, err := client.Generate(context.Background(), sampleInput())
	if err != nil {
		t.Fatalf("Generate() unexpected error: %v", err)
	}

	if got != wantStory {
		t.Fatalf("Generate() = %+v, want %+v", got, wantStory)
	}
}

// GigaChat иногда заворачивает JSON в ```json ... ``` несмотря на просьбу
// этого не делать — клиент должен всё равно распарсить ответ.
func TestClient_Generate_StripsMarkdownCodeFence(t *testing.T) {
	wantStory := domain.JourneyStory{
		Title:  "Заголовок",
		Story:  "Текст истории",
		Teaser: "Тизер",
	}
	raw, _ := json.Marshal(wantStory)
	wrapped := "```json\n" + string(raw) + "\n```"

	srv := newTestServer(t, okOAuthHandler, completionsWithContent(t, wrapped))
	client := NewClient(testConfig(srv))

	got, err := client.Generate(context.Background(), sampleInput())
	if err != nil {
		t.Fatalf("Generate() unexpected error: %v", err)
	}
	if got != wantStory {
		t.Fatalf("Generate() = %+v, want %+v", got, wantStory)
	}
}

// Явная проверка того, что просили в задаче: в структуре ответа нет mood,
// и лишнее поле "mood" от модели просто игнорируется, а не ломает парсинг.
func TestClient_Generate_IgnoresMoodField(t *testing.T) {
	content := `{
		"title": "Заголовок",
		"story": "Текст истории",
		"teaser": "Тизер",
		"mood": "excited"
	}`

	srv := newTestServer(t, okOAuthHandler, completionsWithContent(t, content))
	client := NewClient(testConfig(srv))

	got, err := client.Generate(context.Background(), sampleInput())
	if err != nil {
		t.Fatalf("Generate() unexpected error: %v", err)
	}

	want := domain.JourneyStory{Title: "Заголовок", Story: "Текст истории", Teaser: "Тизер"}
	if got != want {
		t.Fatalf("Generate() = %+v, want %+v (mood must not leak into JourneyStory)", got, want)
	}
}

// --- ошибки ---

func TestClient_Generate_IncompleteStoryIsError(t *testing.T) {
	// нет обязательного поля "story"
	content := `{"title": "Заголовок", "teaser": "Тизер"}`

	srv := newTestServer(t, okOAuthHandler, completionsWithContent(t, content))
	client := NewClient(testConfig(srv))

	_, err := client.Generate(context.Background(), sampleInput())
	if err == nil {
		t.Fatal("Generate() expected error for incomplete story, got nil")
	}
}

func TestClient_Generate_InvalidJSONIsError(t *testing.T) {
	srv := newTestServer(t, okOAuthHandler, completionsWithContent(t, "это не json, а обычный текст"))
	client := NewClient(testConfig(srv))

	_, err := client.Generate(context.Background(), sampleInput())
	if err == nil {
		t.Fatal("Generate() expected error for non-JSON content, got nil")
	}
}

func TestClient_Generate_APIErrorStatusIsError(t *testing.T) {
	completions := func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(gigachatAPIError{Message: "internal error"})
	}

	srv := newTestServer(t, okOAuthHandler, completions)
	client := NewClient(testConfig(srv))

	_, err := client.Generate(context.Background(), sampleInput())
	if err == nil {
		t.Fatal("Generate() expected error for 500 response, got nil")
	}
}

// --- токен и повторная попытка после 401 ---

// Первый вызов /chat/completions отвечает 401 (протухший/невалидный токен),
// клиент должен один раз форснуть обновление токена и повторить запрос.
func TestClient_Generate_RetriesOnceAfterUnauthorized(t *testing.T) {
	var oauthCalls, completionsCalls int32

	wantStory := domain.JourneyStory{Title: "T", Story: "S", Teaser: "Te"}
	content, _ := json.Marshal(wantStory)

	oauth := func(w http.ResponseWriter, _ *http.Request) {
		atomic.AddInt32(&oauthCalls, 1)
		okOAuthHandler(w, nil)
	}

	completions := func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&completionsCalls, 1)
		if n == 1 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		completionsWithContent(t, string(content))(w, r)
	}

	srv := newTestServer(t, oauth, completions)
	client := NewClient(testConfig(srv))

	got, err := client.Generate(context.Background(), sampleInput())
	if err != nil {
		t.Fatalf("Generate() unexpected error: %v", err)
	}
	if got != wantStory {
		t.Fatalf("Generate() = %+v, want %+v", got, wantStory)
	}

	if calls := atomic.LoadInt32(&completionsCalls); calls != 2 {
		t.Fatalf("expected 2 calls to /chat/completions (initial + retry), got %d", calls)
	}
	if calls := atomic.LoadInt32(&oauthCalls); calls != 2 {
		t.Fatalf("expected 2 calls to /oauth (initial fetch + forced refresh), got %d", calls)
	}
}

// Если токен уже в кэше и не истёк, второй вызов Generate не должен снова
// ходить в /oauth.
func TestClient_Generate_ReusesCachedToken(t *testing.T) {
	var oauthCalls int32

	wantStory := domain.JourneyStory{Title: "T", Story: "S", Teaser: "Te"}
	content, _ := json.Marshal(wantStory)

	oauth := func(w http.ResponseWriter, _ *http.Request) {
		atomic.AddInt32(&oauthCalls, 1)
		okOAuthHandler(w, nil)
	}

	srv := newTestServer(t, oauth, completionsWithContent(t, string(content)))
	client := NewClient(testConfig(srv))

	if _, err := client.Generate(context.Background(), sampleInput()); err != nil {
		t.Fatalf("first Generate() unexpected error: %v", err)
	}
	if _, err := client.Generate(context.Background(), sampleInput()); err != nil {
		t.Fatalf("second Generate() unexpected error: %v", err)
	}

	if calls := atomic.LoadInt32(&oauthCalls); calls != 1 {
		t.Fatalf("expected token to be cached (1 call to /oauth), got %d", calls)
	}
}

// --- проверка того, что реально уходит в модель ---

// Клиент не должен просить модель ничего решать за игровую логику: promt
// должен явно запрещать выдумывать факты и не должен содержать поля mood
// в описании формата ответа.
func TestClient_Generate_RequestDoesNotAskForMood(t *testing.T) {
	var capturedBody chatCompletionRequest

	completions := func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewDecoder(r.Body).Decode(&capturedBody)
		completionsWithContent(t, `{"title":"T","story":"S","teaser":"Te"}`)(w, r)
	}

	srv := newTestServer(t, okOAuthHandler, completions)
	client := NewClient(testConfig(srv))

	if _, err := client.Generate(context.Background(), sampleInput()); err != nil {
		t.Fatalf("Generate() unexpected error: %v", err)
	}

	if len(capturedBody.Messages) < 1 {
		t.Fatal("expected at least a system message to be sent")
	}
	systemMsg := capturedBody.Messages[0].Content

	if strings.Contains(systemMsg, `"mood"`) {
		t.Fatal("system prompt must not describe a mood field in the response format")
	}
	if !strings.Contains(systemMsg, "ЗАПРЕЩЕНО придумывать") {
		t.Fatal("system prompt must forbid the model from inventing new game facts")
	}
}

// --- юнит-тест парсинга без сети ---

func TestParseStory(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    domain.JourneyStory
		wantErr bool
	}{
		{
			name:    "чистый JSON",
			content: `{"title":"A","story":"B","teaser":"C"}`,
			want:    domain.JourneyStory{Title: "A", Story: "B", Teaser: "C"},
		},
		{
			name:    "обёрнут в ```json fence",
			content: "```json\n{\"title\":\"A\",\"story\":\"B\",\"teaser\":\"C\"}\n```",
			want:    domain.JourneyStory{Title: "A", Story: "B", Teaser: "C"},
		},
		{
			name:    "лишнее поле mood игнорируется",
			content: `{"title":"A","story":"B","teaser":"C","mood":"happy"}`,
			want:    domain.JourneyStory{Title: "A", Story: "B", Teaser: "C"},
		},
		{
			name:    "невалидный json",
			content: "просто текст",
			wantErr: true,
		},
		{
			name:    "пустой title",
			content: `{"title":"","story":"B","teaser":"C"}`,
			wantErr: true,
		},
		{
			name:    "пустой story",
			content: `{"title":"A","story":"","teaser":"C"}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseStory(tt.content)
			if tt.wantErr {
				if err == nil {
					t.Fatal("parseStory() expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("parseStory() unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("parseStory() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
