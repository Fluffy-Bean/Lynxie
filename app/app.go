package app

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Fluffy-Bean/lynxie/utils"
	"github.com/bwmarrin/discordgo"
)

type Callback func(h *Handler, args []string) Error

type Config struct {
	Prefix  string
	Token   string
	Intents discordgo.Intent
}

type App struct {
	Config   Config
	Commands map[string]Callback
}

func NewApp(config Config) *App {
	return &App{
		Config:   config,
		Commands: make(map[string]Callback),
	}
}

func (a *App) RegisterCommand(cmd string, f Callback) {
	a.Commands[cmd] = f
}

func (a *App) Run() {
	dg, err := discordgo.New("Bot " + a.Config.Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)

		return
	}

	dg.AddHandler(a.handler)
	dg.Identify.Intents = a.Config.Intents

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)

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

func (a *App) handler(session *discordgo.Session, message *discordgo.MessageCreate) {
	h := &Handler{
		Session: session,
		Message: message,
		Reference: &discordgo.MessageReference{
			ChannelID: message.ChannelID,
			MessageID: message.ID,
		},
	}

	if h.Message.Author.ID == h.Session.State.User.ID {
		return
	}
	if h.Message.Author.Bot {
		return
	}

	var cmd string
	var args string

	cmd = h.Message.Content
	cmd = strings.TrimPrefix(cmd, a.Config.Prefix)
	cmd, args, _ = strings.Cut(cmd, " ")

	callback, ok := a.Commands[cmd]
	if !ok {
		// Falling back to default help command
		if cmd == "help" {
			printHelp(a, h)
		}

		return
	}

	h.Session.ChannelTyping(h.Message.ChannelID)

	err := callback(h, strings.Split(args, " "))
	if !err.Ok() {
		printError(a, h, err)
	}
}

func printHelp(a *App, h *Handler) {
	var commands []string
	for cmd := range a.Commands {
		commands = append(commands, cmd)
	}

	h.Session.ChannelMessageSendComplex(h.Message.ChannelID, &discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{
			Title:       "Help",
			Description: strings.Join(commands, "\n"),
			Color:       utils.ColorFromRGB(255, 255, 255),
		},
		Reference: h.Reference,
	})
}

func printError(a *App, h *Handler, e Error) {
	log.Println(e.Err)

	h.Session.ChannelMessageSendComplex(h.Message.ChannelID, &discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{
			Title:       "Error",
			Description: e.Msg,
			Color:       utils.ColorFromRGB(255, 0, 0),
		},
		Reference: h.Reference,
	})
}
