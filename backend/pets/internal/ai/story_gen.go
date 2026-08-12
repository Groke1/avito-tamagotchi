package ai

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/cayman444/avito-gamification-hackathon/backend/pets/internal/domain"

	gigachat "github.com/tigusigalpa/gigachat-go"
)

const systemPrompt = `Ты — виртуальный питомец в приложении-тамагочи. Тебе дают JSON с фактами
о твоём путешествии (локация, события, награда).
Твоя задача — превратить эти факты в короткий, тёплый, немного забавный рассказ
от первого лица, будто питомец сам рассказывает хозяину, что с ним случилось.

Правила:
- Не придумывай события, награды или монеты, которых нет в JSON — используй только то, что дано.
- Пиши по-русски, живо, дружелюбно, 2-4 предложения в поле "story".
- В поле "teaser" намекни на следующее путешествие, не раскрывая деталей.
- Ответ верни СТРОГО в виде одного JSON-объекта, без markdown, без пояснений до или после:
{"title": "...", "story": "...", "teaser": "..."}`

type journeyResultPayload struct {
	Location string   `json:"location"`
	Events   []string `json:"events"`
	Reward   struct {
		Coins int    `json:"coins"`
		Item  string `json:"item,omitempty"`
	} `json:"reward"`
}

type storyResponsePayload struct {
	Title  string `json:"title"`
	Story  string `json:"story"`
	Teaser string `json:"teaser"`
}

type GigaChatStoryGenerator struct {
	client      *gigachat.Client
	model       string
	temperature float64
	maxTokens   int
}

type chatResult struct {
	resp *gigachat.ChatResponse
	err  error
}

func NewGigaChatStoryGenerator(clientID, clientSecret, scope string) (*GigaChatStoryGenerator, error) {
	if clientID == "" || clientSecret == "" {
		return nil, fmt.Errorf("gigachat client id/secret are not configured")
	}

	authKey := base64.StdEncoding.EncodeToString([]byte(clientID + ":" + clientSecret))

	tokenManager := gigachat.NewTokenManager(
		authKey,
		gigachat.WithScope(scope),
	)

	client := gigachat.NewClient(
		tokenManager,
		gigachat.WithDefaultModel(gigachat.GigaChat2Pro),
	)

	return &GigaChatStoryGenerator{
		client:      client,
		model:       gigachat.GigaChat2Pro,
		temperature: 0.8,
		maxTokens:   400,
	}, nil
}

func (g *GigaChatStoryGenerator) Generate(ctx context.Context, journey domain.JourneyResult) (domain.JourneyStory, error) {
	payload := journeyResultPayload{
		Location: journey.Location,
		Events:   journey.Events,
	}
	payload.Reward.Coins = journey.Reward.Coins
	payload.Reward.Item = journey.Reward.Item

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return domain.JourneyStory{}, fmt.Errorf("failed to marshal journey result: %w", err)
	}

	messages := []gigachat.Message{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: string(payloadJSON)},
	}

	resultCh := make(chan chatResult, 1)

	go func() {
		resp, chatErr := g.client.Chat(
			messages,
			gigachat.WithModel(g.model),
			gigachat.WithTemperature(g.temperature),
			gigachat.WithMaxTokens(g.maxTokens),
		)
		resultCh <- chatResult{resp: resp, err: chatErr}
	}()

	select {
	case <-ctx.Done():
		return domain.JourneyStory{}, ctx.Err()

	case res := <-resultCh:
		if res.err != nil {
			return domain.JourneyStory{}, fmt.Errorf("gigachat chat failed: %w", res.err)
		}
		return parseStoryResponse(res.resp)
	}
}

func parseStoryResponse(resp *gigachat.ChatResponse) (domain.JourneyStory, error) {
	if resp == nil || len(resp.Choices) == 0 {
		return domain.JourneyStory{}, errors.New("empty response from gigachat")
	}

	raw := strings.TrimSpace(resp.Choices[0].Message.Content)
	raw = strings.TrimPrefix(raw, "```json")
	raw = strings.TrimPrefix(raw, "```")
	raw = strings.TrimSuffix(raw, "```")
	raw = strings.TrimSpace(raw)

	var parsed storyResponsePayload
	if err := json.Unmarshal([]byte(raw), &parsed); err != nil {
		return domain.JourneyStory{}, fmt.Errorf("failed to parse gigachat story json: %w (raw: %s)", err, raw)
	}

	return domain.JourneyStory{
		Title:  parsed.Title,
		Text:   parsed.Story,
		Teaser: parsed.Teaser,
	}, nil
}
