package scrimpost

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

const (
	YeaResponse   ScrimResponse = "yea"
	NayResponse   ScrimResponse = "nay"
	MaybeResponse ScrimResponse = "maybe"
)

type ScrimResponse string

func (e ScrimResponse) ToApiName() string {
	switch e {
	case "yea":
		return "yea:454970990963982336"
	case "nay":
		return "nay:454970990833696768"
	case "maybe":
		return "❔"
	default:
		return ""
	}
}

func (e ScrimResponse) String() string {
	return string(e)
}

func ScrimResponseFromEmoji(emoji discordgo.Emoji) ScrimResponse {
	fmt.Println(emoji.Name)
	switch emoji.Name {
	case "yea":
		return "yea"
	case "nay":
		return "nay"
	case "❔":
		return "maybe"
	default:
		return ""
	}
}
