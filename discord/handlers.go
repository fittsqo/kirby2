package discord

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

func ReadyHandler(s *discordgo.Session, e *discordgo.Ready) {
	usd := discordgo.UpdateStatusData{
		Activities: []*discordgo.Activity{
			{
				Name: "the stars",
				Type: 3,
			},
		},
	}
	err := s.UpdateStatusComplex(usd)
	if err != nil {
		log.Panicln(err)
	}
	log.Println("Bot connected!")
}

func GuildCreateEventHandler(s *discordgo.Session, e *discordgo.GuildCreate) { // bot turns on or joins a guild
	a.InitGuild(e.Guild.ID)
}

func GuildDeleteEventHandler(s *discordgo.Session, e *discordgo.GuildDelete) { // bot leaves a guild
	if !e.Unavailable {
		a.CutGuild(e.Guild.ID)
	}
}

func GuildMemberAddEventHandler(s *discordgo.Session, e *discordgo.GuildMemberAdd) {
	g, err := s.State.Guild(e.GuildID)
	if err != nil {
		log.Printf("Failed to get guild from cache when GuildMemberAdd is fired: %v", err)
		g, err = s.Guild(e.GuildID)
		if err != nil {
			log.Printf("Failed to get guild from direct request: %v", err)
			return
		}
	}
	gw := a.GetGuildWelcome(g.ID)
	if gw.ChannelID != "" {
		wi := welcomeMessageInfo{
			"<@" + e.User.ID + ">",
			e.Nick,
			e.User.Username,
			g.Name,
			e.User.AvatarURL(fmt.Sprint(PfpSize)),
		}
		welcome := generateWelcomeMessage(gw, wi)
		_, err = s.ChannelMessageSendComplex(gw.ChannelID, &welcome)
		if err != nil {
			log.Printf("Failed to send welcome message: %v", err)
		}
	}
}

func ChannelDeleteEventHandler(s *discordgo.Session, e *discordgo.ChannelDelete) {
	gw := a.GetGuildWelcome(e.GuildID)
	if e.Channel.ID == gw.ChannelID {
		a.SetGuildWelcomeChannel(e.GuildID, "")
	}
}
