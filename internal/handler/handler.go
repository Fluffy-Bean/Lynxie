package handler

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Fluffy-Bean/lynxie/internal/color"
	"github.com/Fluffy-Bean/lynxie/internal/errors"
	"github.com/bwmarrin/discordgo"
)

type Callback func(h *Handler, args []string) errors.Error

type Bot struct {
	Prefix   string
	token    string
	intents  discordgo.Intent
	commands map[string]Callback
	aliases  map[string]string
}

func NewBot(prefix, token string, intents discordgo.Intent) *Bot {
	return &Bot{
		Prefix:   prefix,
		token:    token,
		intents:  intents,
		commands: make(map[string]Callback),
		aliases:  make(map[string]string),
	}
}

func (b *Bot) RegisterCommand(cmd string, f Callback) {
	b.commands[cmd] = f
}

func (b *Bot) RegisterCommandAlias(alias, cmd string) {
	b.aliases[alias] = cmd
}

func (b *Bot) Run() {
	dg, err := discordgo.New("Bot " + b.token)
	if err != nil {
		fmt.Println("Could not create Discord session:", err)

		return
	}

	dg.AddHandler(b.handler)
	dg.Identify.Intents = b.intents

	err = dg.Open()
	if err != nil {
		fmt.Println("Could not connect:", err)

		return
	}

	fmt.Println("Bot is now running. Press CTRL-C to exit.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	dg.Close()
}

type Handler struct {
	Session   *discordgo.Session
	Message   *discordgo.MessageCreate
	Reference *discordgo.MessageReference
}

func (b *Bot) handler(session *discordgo.Session, message *discordgo.MessageCreate) {
	h := &Handler{
		Session: session,
		Message: message,
		Reference: &discordgo.MessageReference{
			ChannelID: message.ChannelID,
			MessageID: message.ID,
		},
	}

	defer func() {
		if r := recover(); r != nil {
			printError(b, h, errors.Error{
				Msg: "But the bot simply refused",
				Err: fmt.Errorf("%v", r),
			})
		}
	}()

	if h.Message.Author.ID == h.Session.State.User.ID {
		return
	}
	if h.Message.Author.Bot {
		return
	}

	var cmd string
	var args string

	cmd = h.Message.Content
	cmd = strings.TrimPrefix(cmd, b.Prefix)
	cmd, args, _ = strings.Cut(cmd, " ")

	alias, ok := b.aliases[cmd]
	if ok {
		cmd = alias
	}

	callback, ok := b.commands[cmd]
	if !ok {
		// Falling back to default help command
		if cmd == "help" {
			printHelp(b, h)
		}

		return
	}

	_ = h.Session.ChannelTyping(h.Message.ChannelID)

	err := callback(h, strings.Split(args, " "))
	if !err.Ok() {
		printError(b, h, err)
	}
}

func printHelp(bot *Bot, h *Handler) {
	var commands []string

	for command := range bot.commands {
		var found []string
		for a, c := range bot.aliases {
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

	_, _ = h.Session.ChannelMessageSendComplex(h.Message.ChannelID, &discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{
			Title:       "Help",
			Description: strings.Join(commands, "\n"),
			Color:       color.RGBToDiscord(255, 255, 255),
			Footer: &discordgo.MessageEmbedFooter{
				Text: "command (aliases...)",
			},
		},
		Reference: h.Reference,
	})
}

func printError(bot *Bot, h *Handler, e errors.Error) {
	fmt.Println(e.Err)

	_, _ = h.Session.ChannelMessageSendComplex(h.Message.ChannelID, &discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{
			Title:       "Error",
			Description: e.Msg,
			Color:       color.RGBToDiscord(255, 0, 0),
		},
		Reference: h.Reference,
	})
}
