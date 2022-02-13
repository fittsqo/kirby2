package discord

import (
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/ftqo/kirby/database"
	"github.com/ftqo/kirby/files"
)

var s *discordgo.Session
var assets *files.Assets
var adapter *database.Adapter
var cc []*discordgo.ApplicationCommand
var tg string

func Start(token, testGuild string) {
	tg = testGuild
	var err error
	assets = files.GetAssets()
	adapter = database.Open()

	s, err = discordgo.New("Bot " + token)
	if err != nil {
		log.Panicf("failed to initialize discordgo session: %v", err)
	}
	s.AddHandler(ReadyHandler)
	s.AddHandler(GuildCreateEventHandler)
	s.AddHandler(GuildDeleteEventHandler)
	s.AddHandler(GuildMemberAddEventHandler)
	s.AddHandler(ChannelDeleteEventHandler)
	s.AddHandler(InteractionCreateEventHandler)

	s.Identify.Intents = discordgo.IntentsGuildMembers |
		discordgo.IntentsGuilds |
		discordgo.IntentsGuildMessages

	err = s.Open()
	if err != nil {
		log.Fatalf("failed to open the discord session: %v", err)
	}
	cc, err = s.ApplicationCommandBulkOverwrite(s.State.User.ID, tg, commands)
	if err != nil {
		log.Panicf("failed to create command application commands: %v", err)
	}
	log.Print("loaded slash commands !")
}

func Stop() {
	log.Print("removing bot commands !")
	for _, cmd := range cc {
		err := s.ApplicationCommandDelete(s.State.User.ID, tg, cmd.ID)
		if err != nil {
			log.Printf("failed to delete %q command: %v", cmd.Name, err)
		}
	}
	time.Sleep(7 * time.Second) // make sure all temporary messages delete
	log.Print("closing bot connection !")
	s.Close()
	log.Print("closing database connection !")
	adapter.Close()
}
