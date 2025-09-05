package config

import (
	"errors"
	"os"
	"strings"
)

const (
	Claude_Welcome_Reply_Key = "claudeWelcomeReply"
	Claude_Key               = "claudeKey"
	Claude_Url               = "claudeUrl"
	Claude_Model             = "claudeModel"

	// Default values
	DefaultClaudeUrl     = "https://api.anthropic.com"
	DefaultClaudeModel   = "claude-3-5-sonnet-20241022"
	DefaultClaudeWelcome = "我是Claude，开始聊天吧！"
)

// GetClaudeWelcomeReply returns the welcome message for Claude bot
func GetClaudeWelcomeReply() string {
	if reply := os.Getenv(Claude_Welcome_Reply_Key); reply != "" {
		return strings.TrimSpace(reply)
	}
	return DefaultClaudeWelcome
}

// GetClaudeKey returns the Claude API key
func GetClaudeKey() string {
	return strings.TrimSpace(os.Getenv(Claude_Key))
}

// GetClaudeUrl returns the Claude API URL with default fallback
func GetClaudeUrl() string {
	if url := strings.TrimSpace(os.Getenv(Claude_Url)); url != "" {
		return url
	}
	return DefaultClaudeUrl
}

// GetClaudeModel returns the Claude model name with default fallback
func GetClaudeModel() string {
	if model := strings.TrimSpace(os.Getenv(Claude_Model)); model != "" {
		return model
	}
	return DefaultClaudeModel
}

// IsClaudeConfigured checks if Claude is properly configured
func IsClaudeConfigured() bool {
	return GetClaudeKey() != ""
}

// ValidateClaudeConfig validates Claude configuration and returns errors if any
func ValidateClaudeConfig() error {
	if GetClaudeKey() == "" {
		return errors.New("claudeKey is required")
	}
	
	// Validate URL format (basic check)
	url := GetClaudeUrl()
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return errors.New("claudeUrl must be a valid URL starting with http:// or https://")
	}
	
	return nil
}