package chat

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/pwh-pwh/aiwechat-vercel/config"
	"github.com/pwh-pwh/aiwechat-vercel/db"
)

const (
	ClaudeUser = "user"
	ClaudeBot  = "assistant"
)

type ClaudeChat struct {
	BaseChat
	key       string
	url       string
	maxTokens int
}

type ClaudeMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ClaudeRequest struct {
	Model         string        `json:"model"`
	Messages      []ClaudeMessage `json:"messages"`
	MaxTokens     int           `json:"max_tokens,omitempty"`
	System        string        `json:"system,omitempty"`
	Temperature   float64       `json:"temperature,omitempty"`
	Stream        bool          `json:"stream"`
}

type ClaudeResponse struct {
	Content []ClaudeContent `json:"content"`
	Model   string          `json:"model"`
}

type ClaudeContent struct {
	Text string `json:"text"`
	Type string `json:"type"`
}

func (s *ClaudeChat) toDbMsg(msg ClaudeMessage) db.Msg {
	return db.Msg{
		Role: msg.Role,
		Msg:  msg.Content,
	}
}

func (s *ClaudeChat) toChatMsg(msg db.Msg) ClaudeMessage {
	return ClaudeMessage{
		Role:    msg.Role,
		Content: msg.Msg,
	}
}

func (s *ClaudeChat) getModel(userID string) string {
	if model, err := db.GetModel(userID, config.Bot_Type_Claude); err == nil && model != "" {
		return model
	} else if model := os.Getenv("claudeModel"); model != "" {
		return model
	}
	return "claude-3-5-sonnet-20241022"
}

func (s *ClaudeChat) chat(userID, msg string) string {
	// Check if user is verified
	if !config.IsUserVerified(userID) {
		return config.GetDevMessage()
	}

	apiUrl := fmt.Sprintf("%s/v1/messages", s.url)
	
	// Get conversation history from database
	var dbMsgs []db.Msg
	if db.ChatDbInstance != nil {
		historyMsgs, err := db.ChatDbInstance.GetMsgList(config.Bot_Type_Claude, userID)
		if err == nil && len(historyMsgs) > 0 {
			dbMsgs = historyMsgs
		}
	}
	
	// Convert database messages to Claude format
	var messages []ClaudeMessage
	for _, dbMsg := range dbMsgs {
		messages = append(messages, s.toChatMsg(dbMsg))
	}
	
	// Add current user message
	messages = append(messages, ClaudeMessage{
		Role:    ClaudeUser,
		Content: msg,
	})
	
	// Create request body
	reqBody := ClaudeRequest{
		Model:     s.getModel(userID),
		Messages:  messages,
		MaxTokens: s.maxTokens,
		System:    config.GetDefaultSystemPrompt(),
		Temperature: 0.7,
		Stream:    false,
	}
	
	// Remove max_tokens if it's 0 or negative
	if s.maxTokens <= 0 {
		reqBody.MaxTokens = 0
	}
	
	// Convert request to JSON
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Sprintf("Error creating request: %v", err)
	}
	
	// Create HTTP request
	req, err := http.NewRequestWithContext(context.Background(), "POST", apiUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Sprintf("Error creating request: %v", err)
	}
	
	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", s.key)
	req.Header.Set("anthropic-version", "2023-06-01")
	
	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Sprintf("Error sending request: %v", err)
	}
	defer resp.Body.Close()
	
	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Sprintf("Error reading response: %v", err)
	}
	
	// Check status code
	if resp.StatusCode != http.StatusOK {
		return fmt.Sprintf("API error (status %d): %s", resp.StatusCode, string(body))
	}
	
	// Parse response
	var claudeResp ClaudeResponse
	if err := json.Unmarshal(body, &claudeResp); err != nil {
		return fmt.Sprintf("Error parsing response: %v", err)
	}
	
	// Extract response text
	if len(claudeResp.Content) == 0 {
		return "Error: Empty response from Claude"
	}
	
	responseText := claudeResp.Content[0].Text
	
	// Save conversation to database
	messages = append(messages, ClaudeMessage{
		Role:    ClaudeBot,
		Content: responseText,
	})
	
	// Convert back to database format for saving
	var saveMsgs []db.Msg
	for _, msg := range messages {
		saveMsgs = append(saveMsgs, s.toDbMsg(msg))
	}
	
	if db.ChatDbInstance != nil {
		db.ChatDbInstance.SetMsgList(config.Bot_Type_Claude, userID, saveMsgs)
	}
	
	return responseText
}

func (c *ClaudeChat) Chat(userID string, msg string) string {
	r, flag := DoAction(userID, msg)
	if flag {
		return r
	}
	return WithTimeChat(userID, msg, c.chat)
}