package gigachat

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/cayman444/avito-gamification-hackathon/backend/pets/internal/domain"
)

//nolint:mnd // сколько последних историй реально стоит класть в промпт
const maxRecentStories = 3

type Client struct {
	cfg    *Config
	http   *http.Client
	tokens *tokenManager
}

func NewClient(cfg *Config) *Client {
	httpClient := &http.Client{
		Timeout: cfg.RequestTimeout,
	}

	if cfg.InsecureSkipVerify {
		httpClient.Transport = &http.Transport{
			//nolint:gosec // временный обход недоверенного сертификата НУЦ Минцифры, см. Config.InsecureSkipVerify
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	return &Client{
		cfg:    cfg,
		http:   httpClient,
		tokens: newTokenManager(cfg, httpClient),
	}
}

// Generate реализует service.JourneyStoryGenerator.
func (c *Client) Generate(ctx context.Context, input domain.JourneyGenerationInput) (domain.JourneyStory, error) {
	story, err := c.generate(ctx, input, false)
	if err == nil {
		return story, nil
	}

	if isUnauthorized(err) {
		// токен мог протухнуть между проверкой в кэше и самим запросом — форсим обновление один раз
		return c.generate(ctx, input, true)
	}

	return domain.JourneyStory{}, err
}

func (c *Client) generate(ctx context.Context, input domain.JourneyGenerationInput, forceTokenRefresh bool) (domain.JourneyStory, error) {
	token, err := c.tokens.getToken(ctx, forceTokenRefresh)
	if err != nil {
		return domain.JourneyStory{}, fmt.Errorf("gigachat: failed to get token: %w", err)
	}

	reqBody := chatCompletionRequest{
		Model: c.cfg.Model,
		Messages: []chatMessage{
			{Role: "system", Content: systemPrompt()},
			{Role: "user", Content: userPrompt(input)},
		},
		//nolint:mnd // низкая температура — истории должны точно следовать переданным фактам, а не фантазировать
		Temperature: 0.7,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return domain.JourneyStory{}, fmt.Errorf("gigachat: failed to marshal request: %w", err)
	}

	url := strings.TrimRight(c.cfg.APIBaseURL, "/") + "/v1/chat/completions"

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return domain.JourneyStory{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.http.Do(req)
	if err != nil {
		return domain.JourneyStory{}, fmt.Errorf("gigachat unavailable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return domain.JourneyStory{}, &gigachatAPIError{StatusCode: http.StatusUnauthorized, Message: "unauthorized"}
	}

	if resp.StatusCode != http.StatusOK {
		var apiErr gigachatAPIError
		_ = json.NewDecoder(resp.Body).Decode(&apiErr)
		apiErr.StatusCode = resp.StatusCode
		return domain.JourneyStory{}, fmt.Errorf("gigachat returned status %d: %s", resp.StatusCode, apiErr.Message)
	}

	var completion chatCompletionResponse
	if err := json.NewDecoder(resp.Body).Decode(&completion); err != nil {
		return domain.JourneyStory{}, fmt.Errorf("gigachat: failed to decode response: %w", err)
	}

	if len(completion.Choices) == 0 {
		return domain.JourneyStory{}, fmt.Errorf("gigachat: empty choices in response")
	}

	story, err := parseStory(completion.Choices[0].Message.Content)
	if err != nil {
		return domain.JourneyStory{}, fmt.Errorf("gigachat: failed to parse story: %w", err)
	}

	return story, nil
}

func isUnauthorized(err error) bool {
	var apiErr *gigachatAPIError
	if e, ok := err.(*gigachatAPIError); ok { //nolint:errorlint // локальная проверка без обёрток
		apiErr = e
		return apiErr.StatusCode == http.StatusUnauthorized
	}
	return false
}

func parseStory(content string) (domain.JourneyStory, error) {
	cleaned := strings.TrimSpace(content)
	cleaned = strings.TrimPrefix(cleaned, "```json")
	cleaned = strings.TrimPrefix(cleaned, "```")
	cleaned = strings.TrimSuffix(cleaned, "```")
	cleaned = strings.TrimSpace(cleaned)

	var story domain.JourneyStory
	if err := json.Unmarshal([]byte(cleaned), &story); err != nil {
		return domain.JourneyStory{}, err
	}

	if story.Title == "" || story.Story == "" {
		return domain.JourneyStory{}, fmt.Errorf("gigachat: incomplete story in response")
	}

	return story, nil
}

func systemPrompt() string {
	return "Ты — виртуальный питомец в мобильном приложении. " +
		"Тебе дают JSON с уже посчитанными игровыми фактами о твоем завершённом путешествии: " +
		"локация, список событий, награда, а также память питомца " +
		"(личность, краткое summary истории, персонажи, незавершённые сюжетные линии) и 2-3 последние истории. " +
		"Твоя единственная задача — превратить эти факты в тёплый, живой рассказ от первого лица питомца. " +
		"СТРОГО ЗАПРЕЩЕНО придумывать новые события, награды, предметы или менять переданные факты — " +
		"используй только то, что дано во входном JSON. " +
		"Учитывай память питомца и предыдущие истории, чтобы путешествия ощущались как продолжающийся сериал. " +
		"Ответ верни СТРОГО в виде JSON без markdown-обёртки, без пояснений до или после, " +
		"со следующими и только следующими полями: " +
		`{"title": string, "story": string, "teaser": string}. ` +
		"title — короткий заголовок истории. story — сам рассказ, 2-5 предложений, от первого лица. " +
		"teaser — одна короткая интригующая фраза, зовущая пользователя в следующее путешествие без конкретной локации, без упоминания награды и без спойлеров"
}

func userPrompt(input domain.JourneyGenerationInput) string {
	recent := input.RecentStories
	if len(recent) > maxRecentStories {
		recent = recent[len(recent)-maxRecentStories:]
	}

	journeyJSON, _ := json.Marshal(input.Journey)
	memoryJSON, _ := json.Marshal(input.Memory)
	recentJSON, _ := json.Marshal(recent)

	return fmt.Sprintf(
		"Факты о новом путешествии (JourneyResult):\n%s\n\n"+
			"Память питомца (PetMemory):\n%s\n\n"+
			"Последние истории (для стиля и продолжения сюжета):\n%s\n\n"+
			"Сгенерируй JSON по описанному формату.",
		journeyJSON, memoryJSON, recentJSON,
	)
}
