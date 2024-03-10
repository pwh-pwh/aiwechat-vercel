package chat

import (
	"context"
	"github.com/google/generative-ai-go/genai"
	"github.com/pwh-pwh/aiwechat-vercel/config"
	"github.com/pwh-pwh/aiwechat-vercel/db"
	"google.golang.org/api/option"
)

const (
	GeminiUser = "user"
	GeminiBot  = "model"
)

type GeminiChat struct {
	BaseChat
	key string
}

func (s *GeminiChat) toDbMsg(msg *genai.Content) db.Msg {
	text := msg.Parts[0].(genai.Text)
	return db.Msg{
		Role: msg.Role,
		Msg:  string(text),
	}
}

func (s *GeminiChat) toChatMsg(msg db.Msg) *genai.Content {
	return &genai.Content{Parts: []genai.Part{genai.Text(msg.Msg)}, Role: msg.Role}
}

func (s *GeminiChat) chat(userId, msg string) string {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(s.key))
	if err != nil {
		return err.Error()
	}
	defer client.Close()
	model := client.GenerativeModel("gemini-pro")
	// Initialize the chat
	cs := model.StartChat()
	var msgs = GetMsgListWithDb(config.Bot_Type_Gemini, userId, &genai.Content{
		Parts: []genai.Part{
			genai.Text(msg),
		},
		Role: GeminiUser,
	}, s.toDbMsg, s.toChatMsg)
	if len(msgs) > 1 {
		cs.History = msgs[:len(msgs)-1]
	}

	resp, err := cs.SendMessage(ctx, genai.Text(msg))
	if err != nil {
		return err.Error()
	}
	text := resp.Candidates[0].Content.Parts[0].(genai.Text)
	msgs = append(msgs, &genai.Content{Parts: []genai.Part{
		text,
	}, Role: GeminiBot})
	SaveMsgListWithDb(config.Bot_Type_Gemini, userId, msgs, s.toDbMsg)
	return string(text)
}

func (g *GeminiChat) Chat(userID string, msg string) string {
	return WithTimeChat(userID, msg, g.chat)

}
