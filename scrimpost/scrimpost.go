package scrimpost

type GuildData struct {
	ID       string
	Events   []*Event
	Messages []*Message
}

type Message struct {
	ID        string
	ChannelID string
	EventID   int
}

type Event struct {
	ID           int
	MessageID    string
	Participants []*Participant `json:",omitempty"`
}

type Participant struct {
	ID        string
	Available ScrimResponse
}

func FromMessage(GuildID, ChannelID string, MessageID string) (*Event, error) {
	return database.GetEventByMessage(GuildID, ChannelID, MessageID)
}
