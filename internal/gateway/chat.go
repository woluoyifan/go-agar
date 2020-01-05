package gateway

const (
	ChatTypeSystem = "0"
	ChatTypePlayer = "1"
)

type Chat struct {
	Type string
	Data string
}

func NewSystemChat(data string) *Chat {
	return &Chat{
		Type: ChatTypeSystem,
		Data: data,
	}
}

func NewPlayerChat(data string) *Chat {
	return &Chat{
		Type: ChatTypePlayer,
		Data: data,
	}
}
