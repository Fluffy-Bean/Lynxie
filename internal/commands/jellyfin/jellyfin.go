package jellyfin

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"

	"github.com/Fluffy-Bean/lynxie/internal/handler"
	"github.com/Fluffy-Bean/lynxie/internal/utils"
)

func RegisterCommands(h *handler.Handler) {
	_ = h.RegisterCommand("jf-register", cmdRegister(h))
}

func cmdRegister(h *handler.Handler) handler.Callback {
	_ = h.RegisterCommandHelp("jf-register", "Register for the Jellfin instance")

	return func(session *discordgo.Session, message *discordgo.MessageCreate) error {
		if message == nil || message.Author == nil {
			return fmt.Errorf("get message author")
		}

		existingUser, err := h.FindJellyfinRegistrationByDiscordID(message.Author.ID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("find user: %w", err)
		}

		if existingUser != nil {
			_, _ = session.ChannelMessageSendComplex(message.ChannelID, &discordgo.MessageSend{
				Embed: &discordgo.MessageEmbed{
					Title:       "Jellyfin",
					Description: "You already have a registered account! Shew",
					Color:       utils.RGBToDiscord(255, 255, 255),
				},
				Reference: message.Reference(),
			})

			return nil
		}

		_, username, password, err := h.CreateJellyfinRegistration(message.Author.ID, message.Author.Username)
		if err != nil {
			return fmt.Errorf("create jellyfin registration: %w", err)
		}

		_, _ = session.ChannelMessageSendComplex(message.ChannelID, &discordgo.MessageSend{
			Embed: &discordgo.MessageEmbed{
				Title:       "Jellyfin",
				Description: fmt.Sprintf("Registered for an account! Your login is:\nUsername: `%s`\nPassword: `%s`", username, password),
				Color:       utils.RGBToDiscord(255, 255, 255),
			},
			Reference: message.Reference(),
			Flags:     discordgo.MessageFlagsEphemeral,
		})

		return nil
	}
}
