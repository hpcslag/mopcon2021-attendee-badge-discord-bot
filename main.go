package main

import (
	"flag"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// before using the bot, please use this url paste to browser and accept it to channel or guild
// https://discord.com/api/oauth2/authorize?client_id=[put client id inside]&scope=bot&permissions=8

// Variables used for command line parameters
var (
	BotSecret string
	GuildID   string
)

func init() {

	flag.StringVar(&BotSecret, "t", "", "Bot Token")
	flag.StringVar(&GuildID, "g", "", "855718269092233247")
	flag.Parse()
}

func main() {

	ReadXLSXToMap("test.xlsx", &KeyPairMap)

	RegisterBotFuncAndRun(DiscordAuth{
		BotSecret: BotSecret,
	}, messageCreate)

}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// check message and messages guild role ids
	/*
		fmt.Println("Received msg: ", m.Content)
		fmt.Println("GuildID:", m.GuildID)
		st, _ := s.GuildRoles(m.GuildID)
		for _, v := range st {
			fmt.Println("RoleName and ID: ", v.Name, v.ID)
		}
	*/

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	// https://github-wiki-see.page/m/bwmarrin/discordgo/wiki/FAQ
	isDM, err := ComesFromDM(s, m)
	if err != nil {
		panic(err)
	}
	if isDM {
		//					 global guild id 			role id e.g: "885459800740212817" is for attendee
		s.GuildMemberRoleAdd(GuildID, m.Author.ID, "885459800740212817")
		s.ChannelMessageSend(m.ChannelID, "Done.")
	} else {

		if m.Content == "給我身分" {
			fmt.Println(m.GuildID, m.Author.ID)
			s.GuildMemberRoleAdd(m.GuildID, m.Author.ID, "885459800740212817")

			s.ChannelMessageSend(m.ChannelID, "好")
		}

		if badge, ok := KeyPairMap[m.Content]; ok {
			s.ChannelMessageSend(m.ChannelID, "Badge: "+badge)
		}

		// If the message is "ping" reply with "Pong!"
		if m.Content == "ping" {
			s.ChannelMessageSend(m.ChannelID, "Pong!")
		}

		// If the message is "pong" reply with "Ping!"
		if m.Content == "pong" {
			s.ChannelMessageSend(m.ChannelID, "Ping!")
		}
	}
}

func ComesFromDM(s *discordgo.Session, m *discordgo.MessageCreate) (bool, error) {
	channel, err := s.State.Channel(m.ChannelID)
	if err != nil {
		if channel, err = s.Channel(m.ChannelID); err != nil {
			return false, err
		}
	}

	return channel.Type == discordgo.ChannelTypeDM, nil
}
