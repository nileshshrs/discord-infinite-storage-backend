package bot

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/nileshshrs/infinite-storage/config"
)

func Run(cfg *config.Config) *discordgo.Session {
	dg, err := discordgo.New("Bot " + cfg.DiscordToken)
	if err != nil {
		log.Fatalf("Error creating Discord session: %v", err)
	}

	// Example message handler
	dg.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}
		if m.Content == "!ping" {
			s.ChannelMessageSend(m.ChannelID, "Pong!")
		}
	})

	if err := dg.Open(); err != nil {
		log.Fatalf("Error opening connection: %v", err)
	}

	log.Println("Discord bot is running...")
	return dg
}
