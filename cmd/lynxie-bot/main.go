package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/Fluffy-Bean/lynxie/internal/bot"
	"github.com/Fluffy-Bean/lynxie/internal/color"
	"github.com/Fluffy-Bean/lynxie/internal/commands/debug"
	"github.com/Fluffy-Bean/lynxie/internal/commands/fortnite"
	"github.com/Fluffy-Bean/lynxie/internal/commands/img"
	"github.com/Fluffy-Bean/lynxie/internal/commands/porb"
	"github.com/Fluffy-Bean/lynxie/internal/commands/tinyfox"
)

var flags struct {
	Token    *string
	Prefix   *string
	DataPath *string
}

func main() {
	err := parseFlags()
	if err != nil {
		log.Fatal("parse flags:", err)
	}

	h := bot.NewHandler(*flags.Prefix, handleHelp, handleError)

	debug.RegisterDebugCommands(h)
	img.RegisterImgCommands(h)
	tinyfox.RegisterTinyfoxCommands(h)
	porb.RegisterPorbCommands(h)
	fortnite.RegisterFortniteCommands(h)

	databasePath := path.Join(*flags.DataPath, "storage.db")

	log.Fatal(h.Run(databasePath, os.Getenv("TOKEN"), discordgo.IntentsGuildMessages))
}

func parseFlags() error {
	flags.Prefix = flag.String("prefix", ">", "The prefix the bot will search for within messages")
	flags.DataPath = flag.String("datapath", "", "The path to the datapath, defaults to $HOME/.lynxie if empty")
	flags.Token = flag.String("token", "", "The bot token, defaults to checking for TOKEN in env")

	flag.Parse()

	if *flags.DataPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return err
		}

		*flags.DataPath = path.Join(homeDir, ".lynxie", "data")
	}

	if *flags.Token == "" {
		*flags.Token = os.Getenv("TOKEN")
	}

	return nil
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
