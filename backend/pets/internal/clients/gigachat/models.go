package gigachat

type oauthTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresAt   int64  `json:"expires_at"` // unix-время в миллисекундах
}

type chatMessage struct {
	Role    string `json:"role"` // system | user | assistant
	Content string `json:"content"`
}

type chatCompletionRequest struct {
	Model       string        `json:"model"`
	Messages    []chatMessage `json:"messages"`
	Temperature float64       `json:"temperature,omitempty"`
}

type chatCompletionResponse struct {
	Choices []struct {
		Message chatMessage `json:"message"`
	} `json:"choices"`
}

type gigachatAPIError struct {
	StatusCode int    `json:"-"`
	Message    string `json:"message"`
}

func (e *gigachatAPIError) Error() string {
	return e.Message
}
