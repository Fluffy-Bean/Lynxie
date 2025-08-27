package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/Fluffy-Bean/lynxie/internal/bot"
	"github.com/Fluffy-Bean/lynxie/internal/color"
	"github.com/Fluffy-Bean/lynxie/internal/commands/debug"
	"github.com/Fluffy-Bean/lynxie/internal/commands/img"
	"github.com/Fluffy-Bean/lynxie/internal/commands/porb"
	"github.com/Fluffy-Bean/lynxie/internal/commands/tinyfox"
)

func main() {
	h := bot.NewHandler(">", handleHelp, handleError)

	debug.RegisterDebugCommands(h)
	img.RegisterImgCommands(h)
	tinyfox.RegisterTinyfoxCommands(h)
	porb.RegisterPorbCommands(h)

	log.Fatal(h.Run(os.Getenv("TOKEN"), discordgo.IntentsGuildMessages))
}

func handleHelp(h *bot.Handler, c bot.CommandContext) {
	var commands []string

	for command := range h.GetCommands() {
		var found []string
		for a, c := range h.GetCommandAliases() {
			if c == command {
				found = append(found, a)
			}
		}

		if len(found) > 0 {
			commands = append(commands, fmt.Sprintf("%s (%s)", command, strings.Join(found, ", ")))
		} else {
			commands = append(commands, command)
		}
	}

	_, _ = c.Session.ChannelMessageSendComplex(c.Message.ChannelID, &discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{
			Title:       "Help",
			Description: strings.Join(commands, "\n"),
			Color:       color.RGBToDiscord(255, 255, 255),
			Footer: &discordgo.MessageEmbedFooter{
				Text: "command (aliases...)",
			},
		},
		Reference: c.Message.Reference(),
	})
}

func handleError(h *bot.Handler, c bot.CommandContext, err error) {
	fmt.Println(err)

	_, _ = c.Session.ChannelMessageSendComplex(c.Message.ChannelID, &discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{
			Title:       "Error",
			Description: err.Error(),
			Color:       color.RGBToDiscord(255, 0, 0),
		},
		Reference: c.Message.Reference(),
	})
}
