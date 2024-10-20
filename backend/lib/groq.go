package lib

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Request struct {
	Messages []Message `json:"messages"`
	Model    string    `json:"model"`
}

func MakeRequest(messages []Message, model string) Request {
	return Request{
		Messages: messages,
		Model:    model,
	}
}
