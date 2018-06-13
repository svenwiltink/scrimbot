package scrimpost

import (
	"github.com/bwmarrin/discordgo"
	"time"
	"log"
)

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
	GuildID      string `json:"-"`
	ID           int
	MessageID    string
	Participants []*Participant `json:",omitempty"`
}

type Participant struct {
	ID        string
	Available ScrimResponse
}

func (event *Event) CreateDiscordEmbed() *discordgo.MessageEmbed {

	yeaList := ""
	nayList := ""
	maybeList := ""

	for _, participant := range event.Participants {
		switch participant.Available {
		case YeaResponse:
			yeaList += participant.ID + "\n"
		case NayResponse:
			nayList += participant.ID + "\n"
		case MaybeResponse:
			maybeList += participant.ID + "\n"
		}
	}

	if yeaList == "" {
		yeaList = "<nobody>"
	}

	if nayList == "" {
		nayList = "<nobody>"
	}

	if maybeList == "" {
		maybeList = "<nobody>"
	}

	return &discordgo.MessageEmbed{
		Author:      &discordgo.MessageEmbedAuthor{},
		Color:       0x00ff00, // Green
		Description: "Lets scrim",
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name:   "yea",
				Value:  yeaList,
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "nay",
				Value:  nayList,
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "maybe",
				Value:  maybeList,
				Inline: true,
			},
		},
		Timestamp: time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
		Title:     "Lets scrim",
	}
}

func (event *Event) HandleReaction(UserID string, response ScrimResponse) (bool, error) {
	for _, participant := range event.Participants {
		if participant.ID == UserID {
			if participant.Available != response {
				participant.Available = response
				err := SaveParticipation(event, participant)
				if err != nil {
					return false, err
				}
				return true, nil
			} else {
				return false, nil
			}
		}
	}

	participant := &Participant{
		ID:        UserID,
		Available: response,
	}

	event.Participants = append(event.Participants, participant)
	err := SaveParticipation(event, participant)
	if err != nil {
		log.Println(err)
		return false, err
	}

	return true, nil
}

func FromMessage(GuildID, ChannelID string, MessageID string) (*Event, error) {
	return database.GetEventByMessage(GuildID, ChannelID, MessageID)
}

func CreateEvent(GuildID string) (*Event, error) {
	return database.CreateEvent(GuildID)
}

func SaveMessage(GuildID string, ChannelId string, MessageId string, event *Event) error {
	return database.SaveMessage(GuildID, &Message{
		ID:        MessageId,
		ChannelID: ChannelId,
		EventID:   event.ID,
	})
}

func SaveEvent(Event *Event) error {
	return database.SaveEvent(Event)
}

func SaveParticipation(Event *Event, Participant *Participant) error {
	return database.SaveParticipation(Event, Participant)
}
