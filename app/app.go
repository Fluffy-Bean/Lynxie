package app

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

type Config struct {
	Prefix  string
	Token   string
	Intents discordgo.Intent
}

type App struct {
	Config   Config
	Commands map[string]func(h *Handler, args []string)
}

func NewApp(config Config) *App {
	return &App{
		Config:   config,
		Commands: make(map[string]func(h *Handler, args []string)),
	}
}

func (a *App) RegisterCommand(cmd string, f func(h *Handler, args []string)) {
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
	Session *discordgo.Session
	Message *discordgo.MessageCreate
}

func (a *App) handler(session *discordgo.Session, message *discordgo.MessageCreate) {
	h := &Handler{
		Session: session,
		Message: message,
	}

	if h.Message.Author.ID == h.Session.State.User.ID {
		return
	}

	if h.Message.Author.Bot {
		return
	}

	var command string
	var args string

	command = h.Message.Content
	command = strings.TrimSpace(command)
	command = strings.TrimPrefix(command, a.Config.Prefix)
	command, args, _ = strings.Cut(command, " ")

	callback, ok := a.Commands[command]
	if !ok {
		return
	}

	callback(h, strings.Split(args, " "))
}
