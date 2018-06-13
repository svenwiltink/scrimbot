package scrimpost

type Database interface {
	GetGuildData(GuildID string) (*GuildData, error)
	GetEventById(EventId string) (*Event, error)
	GetEventByMessage(GuildID string, ChannelID string, MessageId string) (*Event, error)

	CreateEvent(GuildID string) (*Event, error)
	SaveParticipation(Event *Event, Participant *Participant) error
	SaveMessage(Message *Message)
}

var database Database

func RegisterDatabase(db Database) {
	database = db
}
