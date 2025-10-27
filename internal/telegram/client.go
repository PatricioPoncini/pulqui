package telegram

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	token      string
	httpClient *http.Client
	baseURL    string
}

func NewClient(token string) *Client {
	return &Client{
		token:      token,
		httpClient: &http.Client{},
		baseURL:    "https://api.telegram.org/bot" + token,
	}
}

func (c *Client) GetUpdates(offset int) ([]Update, error) {
	resp, err := c.httpClient.Get(fmt.Sprintf("%s/getUpdates?offset=%d", c.baseURL, offset))
	if err != nil {
		return nil, fmt.Errorf("failed to get updates: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result struct {
		Ok     bool     `json:"ok"`
		Result []Update `json:"result"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if !result.Ok {
		return nil, fmt.Errorf("telegram API returned ok=false")
	}

	return result.Result, nil
}

func (c *Client) SendMessage(chatID int64, text string, parseMode ...string) error {
	params := url.Values{}
	params.Set("chat_id", fmt.Sprintf("%d", chatID))
	params.Set("text", text)

	if len(parseMode) > 0 {
		params.Set("parse_mode", parseMode[0])
	}

	resp, err := c.httpClient.Post(
		fmt.Sprintf("%s/sendMessage", c.baseURL),
		"application/x-www-form-urlencoded",
		strings.NewReader(params.Encode()),
	)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("telegram API error: %s", string(body))
	}

	return nil
}
