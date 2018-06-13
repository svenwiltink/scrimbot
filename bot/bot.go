package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/svenwiltink/scrimbot/config"
	"github.com/svenwiltink/scrimbot/scrimpost"
	"log"
	"strings"
)

type ScrimBot struct {
	config         *config.Config
	discordSession *discordgo.Session
	scrimposts     *scrimpost.Database
}

func (bot *ScrimBot) Start() error {

	token := "Bot " + bot.config.DiscordToken
	discord, err := discordgo.New(token)

	if err != nil {
		return err
	}

	bot.discordSession = discord

	bot.discordSession.AddHandler(bot.handleGuildCreate)
	bot.discordSession.AddHandler(bot.handleMessageCreate)
	bot.discordSession.AddHandler(bot.handleReactionAdd)

	// Open a websocket connection to Discord and begin listening.
	err = discord.Open()
	if err != nil {
		return fmt.Errorf("error opening connection: %v", err)
	}

	return nil
}

func (bot *ScrimBot) handleMessageCreate(session *discordgo.Session, message *discordgo.MessageCreate) {

	if message.Author.ID == session.State.User.ID {
		return
	}

	if strings.HasPrefix(message.Content, "!scrim") {
		channel, err := session.Channel(message.ChannelID)
		if err != nil {
			log.Println(err)
			return
		}

		//post := scrimpost.Create(channel.GuildID, channel.ID)
		message, err := session.ChannelMessageSend(channel.ID, "fuck yea bitches")

		if err != nil {
			log.Println(err)
			return
		}

		//post.MessageID = message.ID
		//err = bot.scrimposts.Put(post)
		//if err != nil {
		//	log.Println(err)
		//}

		session.MessageReactionAdd(message.ChannelID, message.ID, scrimpost.YeaResponse.ToApiName())
		session.MessageReactionAdd(message.ChannelID, message.ID, scrimpost.NayResponse.ToApiName())
		session.MessageReactionAdd(message.ChannelID, message.ID, scrimpost.MaybeResponse.ToApiName())
	}
}

func (bot *ScrimBot) handleReactionAdd(session *discordgo.Session, event *discordgo.MessageReactionAdd) {

	if event.UserID == session.State.User.ID {
		return
	}

	message, err := session.ChannelMessage(event.ChannelID, event.MessageID)

	if err != nil {
		log.Printf("Could not get message: %v", err)
		return
	}

	if message.Author.ID != session.State.User.ID {
		return
	}

	channel, err := session.Channel(event.ChannelID)
	if err != nil {
		log.Printf("Could not get channel: %v", err)
		return
	}

	post, err := scrimpost.FromMessage(channel.GuildID, channel.ID, message.ID)
	if err != nil {
		log.Println(err)
	}
	log.Println(post)
	//
	//post, err := scrimpost.FromMessage(channel.GuildID, message)
	//
	//if err != nil {
	//	log.Println(err)
	//	return
	//}
	//
	//shouldUpdate, err := post.HandleReaction(session, event)
	//if err != nil {
	//	log.Println(err)
	//	return
	//}
	//
	//if shouldUpdate {
	//	session.ChannelMessageEditEmbed(message.ChannelID, message.ID, post.CreateEmbed())
	//	bot.scrimposts.Put(post)
	//}
	//
	//session.MessageReactionRemove(message.ChannelID, message.ID, event.Emoji.APIName(), event.UserID)
}

func (bot *ScrimBot) handleGuildCreate(session *discordgo.Session, create *discordgo.GuildCreate) {
	//bot.scrimposts.LoadEntriesForGuild(create.Guild.ID)
}

func (bot *ScrimBot) Close() {
	bot.discordSession.Close()
}

func CreateNewBot(config *config.Config) *ScrimBot {
	return &ScrimBot{
		config: config,
	}
}
