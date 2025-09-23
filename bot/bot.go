package bot

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/nileshshrs/infinite-storage/config"
)

// SendMessage sends a message to a specific channel
func SendMessage(dg *discordgo.Session, channelID, message string) {
	_, err := dg.ChannelMessageSend(channelID, message)
	if err != nil {
		log.Printf("Failed to send message to channel %s: %v", channelID, err)
	}
}

// Example: Send "Hello world!" after bot starts
func Run(cfg *config.Config) *discordgo.Session {
    dg, err := discordgo.New("Bot " + cfg.DiscordToken)
    if err != nil {
        log.Fatalf("Error creating Discord session: %v", err)
    }

    if err := dg.Open(); err != nil {
        log.Fatalf("Error opening Discord connection: %v", err)
    }

    log.Println("Discord bot is running...")

    // Send message to the fixed channel
    // _, err = dg.ChannelMessageSend(cfg.DiscordChannelID, "Hello, server! Mahiru is online")
    // if err != nil {
    //     log.Printf("Failed to send message: %v", err)
    // }

    return dg
}