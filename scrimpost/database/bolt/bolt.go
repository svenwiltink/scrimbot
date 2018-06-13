package bolt

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/svenwiltink/scrimbot/scrimpost"
	"github.com/svenwiltink/scrimbot/util"
	"log"
)

type Database struct {
	db   *bolt.DB
	path string
}

func (db *Database) GetGuildData(GuildID string) (*scrimpost.GuildData, error) {
	panic("implement me")
}

func (db *Database) GetEventById(EventId string) (*scrimpost.Event, error) {
	panic("implement me")
}

func (db *Database) GetEventByMessage(GuildID string, ChannelID string, MessageID string) (*scrimpost.Event, error) {
	event := &scrimpost.Event{}
	err := db.db.View(func(tx *bolt.Tx) error {

		guildBucket, err := getGuildBucket(tx, GuildID)
		if err != nil {
			return err
		}

		messages := guildBucket.Bucket([]byte("messages"))
		if messages == nil {
			return fmt.Errorf("no messages bucket for guild %s", GuildID)
		}

		message := messages.Get([]byte(ChannelID + "+" + MessageID))
		if message == nil {
			return errors.New("could not get event by message")
		}

		messageStruct := &scrimpost.Message{}
		err = json.Unmarshal(message, messageStruct)

		if err != nil {
			return err
		}

		event, err = getEventById(guildBucket, messageStruct.EventID)

		if err != nil {
			return err
		}

		return nil
	})

	event.GuildID = GuildID

	return event, err
}

func (db *Database) CreateEvent(GuildID string) (*scrimpost.Event, error) {
	event := &scrimpost.Event{}
	err := db.db.Update(func(tx *bolt.Tx) error {
		guilds := tx.Bucket([]byte("guilds"))
		if guilds == nil {
			return errors.New("unable to get guilds bucket")
		}

		guild, err := guilds.CreateBucketIfNotExists([]byte(GuildID))
		if err != nil {
			return err
		}

		events, err := guild.CreateBucketIfNotExists([]byte("events"))
		if err != nil {
			return err
		}

		// get the next eventID and set it in the message
		nextEventId, _ := events.NextSequence()
		event.ID = int(nextEventId)

		saveEvent(event, events)

		return nil
	})

	return event, err
}

func (db *Database) SaveEvent(Event *scrimpost.Event) error {
	//event := &scrimpost.Event{}
	err := db.db.Update(func(tx *bolt.Tx) error {
		guilds := tx.Bucket([]byte("guilds"))
		if guilds == nil {
			return errors.New("unable to get guilds bucket")
		}

		guild, err := guilds.CreateBucketIfNotExists([]byte(Event.GuildID))
		if err != nil {
			return err
		}

		events, err := guild.CreateBucketIfNotExists([]byte("events"))
		if err != nil {
			return err
		}

		saveEvent(Event, events)

		return nil
	})

	return err
}

func saveEvent(event *scrimpost.Event, eventsbucket *bolt.Bucket) error {
	eventBucket, err := eventsbucket.CreateBucketIfNotExists(util.Itob(event.ID))

	eventParticipants := event.Participants
	event.Participants = nil

	eventBytes, err := json.Marshal(event)
	if err != nil {
		return err
	}

	if err = eventBucket.Put([]byte("data"), eventBytes); err != nil {
		return err
	}

	participants, err := eventBucket.CreateBucketIfNotExists([]byte("participants"))
	if err != nil {
		log.Println(err)
		return err
	}

	for _, participant := range eventParticipants {
		participantBytes, err := json.Marshal(participant)
		if err != nil {
			return err
		}

		if err = participants.Put([]byte(participant.ID), participantBytes); err != nil {
			return err
		}

	}

	return nil
}

func (db *Database) SaveParticipation(Event *scrimpost.Event, Participant *scrimpost.Participant) error {
	return db.db.Update(func(tx *bolt.Tx) error {
		guilds := tx.Bucket([]byte("guilds"))
		if guilds == nil {
			return errors.New("unable to get guilds bucket")
		}

		guild, err := guilds.CreateBucketIfNotExists([]byte(Event.GuildID))
		if err != nil {
			return err
		}

		eventsBucket, err := guild.CreateBucketIfNotExists([]byte("events"))
		if err != nil {
			return err
		}

		eventBucket, err := eventsBucket.CreateBucketIfNotExists(util.Itob(Event.ID))
		participants, err := eventBucket.CreateBucketIfNotExists([]byte("participants"))
		if err != nil {
			return err
		}

		participantBytes, err := json.Marshal(Participant)
		if err != nil {
			return err
		}

		if err = participants.Put([]byte(Participant.ID), participantBytes); err != nil {
			return err
		}

		return nil
	})
}

func (db *Database) SaveMessage(GuildID string, Message *scrimpost.Message) error {
	return db.db.Update(func(tx *bolt.Tx) error {
		guilds := tx.Bucket([]byte("guilds"))
		if guilds == nil {
			return errors.New("unable to get guilds bucket")
		}

		guild, err := guilds.CreateBucketIfNotExists([]byte(GuildID))
		if err != nil {
			return err
		}

		messages, err := guild.CreateBucketIfNotExists([]byte("messages"))
		if err != nil {
			return err
		}

		messageBytes, err := json.Marshal(Message)
		if err != nil {
			return err
		}

		if err = messages.Put([]byte(Message.ChannelID+"+"+Message.ID), messageBytes); err != nil {
			return err
		}

		return nil
	})

}

func Load(path string) (*Database, error) {

	database, err := bolt.Open(path, 0600, nil)

	if err != nil {
		return nil, err
	}

	if err = createRootBucket(database); err != nil {
		return nil, err
	}

	return &Database{
		path: path,
		db:   database,
	}, nil
}

func createRootBucket(db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("guilds"))
		if err != nil {
			log.Println(err)
			return err
		}
		return nil
	})
}

func getGuildBucket(tx *bolt.Tx, GuildID string) (*bolt.Bucket, error) {
	guilds := tx.Bucket([]byte("guilds"))
	if guilds == nil {
		return nil, errors.New("unable to get guilds bucket")
	}

	bucket := guilds.Bucket([]byte(GuildID))
	if bucket == nil {
		return nil, errors.New("unable to get bucket for guild")
	}

	return bucket, nil
}

func getEventById(guildBucket *bolt.Bucket, EventID int) (*scrimpost.Event, error) {
	eventsBucket := guildBucket.Bucket([]byte("events"))
	if eventsBucket == nil {
		return nil, errors.New("no events bucket for guild")
	}

	eventBucket := eventsBucket.Bucket(util.Itob(EventID))
	if eventBucket == nil {
		return nil, fmt.Errorf("no bucket found for event %d", EventID)
	}

	data := eventBucket.Get([]byte("data"))

	event := &scrimpost.Event{}
	event.ID = EventID

	err := json.Unmarshal(data, event)
	if err != nil {
		return nil, err
	}

	participantsBucket := eventBucket.Bucket([]byte("participants"))

	participantsBucket.ForEach(func(key, value []byte) error {
		participant := &scrimpost.Participant{}
		json.Unmarshal(value, participant)

		event.Participants = append(event.Participants, participant)
		return nil
	})

	return event, nil
}
