package chat

type Echo struct{}

func (e *Echo) Chat(userID string, msg string) string {
	return msg
}
