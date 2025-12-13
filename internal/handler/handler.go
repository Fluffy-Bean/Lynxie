package handler

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	_ "github.com/mattn/go-sqlite3"

	"github.com/Fluffy-Bean/lynxie/internal/utils"
)

type Callback func(session *discordgo.Session, message *discordgo.MessageCreate) error

func NewHandler(prefix string, db *sql.DB) *Handler {
	var handler Handler

	handler.prefix = prefix

	handler.commands = make(map[string]Callback)
	handler.commandHelp = make(map[string]string)
	handler.commandAliases = make(map[string]string)

	handler.values = make(map[string]string)

	handler.httpClient = &http.Client{
		Timeout: 10 * time.Second,
	}

	handler.db = db

	return &handler
}

type Handler struct {
	prefix string

	commands       map[string]Callback
	commandHelp    map[string]string
	commandAliases map[string]string

	values map[string]string

	httpClient *http.Client

	db *sql.DB
	dg *discordgo.Session
}

func (h *Handler) Run(token string) error {
	if err := h.doDatabaseMigrations(); err != nil {
		return fmt.Errorf("do database migrations: %w", err)
	}

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		return fmt.Errorf("create discord session: %s", err)
	}

	dg.Identify.Intents = discordgo.IntentsGuildMessages
	dg.AddHandler(h.handleCommand)

	err = dg.Open()
	if err != nil {
		return fmt.Errorf("opening discord session: %s", err)
	}
	defer dg.Close()

	fmt.Println("Bot is now running. Press CTRL-C to exit.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	<-sc

	return nil
}

func (h *Handler) SetPrefix(prefix string) {
	h.prefix = prefix
}

func (h *Handler) SetValue(key, value string) {
	h.values[key] = value
}

func (h *Handler) GetValue(key string) string {
	v, ok := h.values[key]
	if ok {
		return v
	}

	return ""
}

func (h *Handler) ScheduleTask(do func(), every time.Duration) {
	// Run the command initially
	do()

	// This sucks so bad man
	go func() {
		for range time.Tick(every) {
			func() {
				defer func() {
					if r := recover(); r != nil {
						fmt.Println("panic in scheduled task:", r)
					}
				}()

				do()
			}()
		}
	}()
}

func (h *Handler) RegisterCommand(command string, callback Callback) error {
	if _, ok := h.commands[command]; ok {
		return fmt.Errorf("command already registered")
	}

	h.commands[command] = callback

	return nil
}

func (h *Handler) RegisterCommandHelp(command string, help string) error {
	if _, ok := h.commands[command]; ok {
		return fmt.Errorf("help already registered")
	}

	h.commandHelp[command] = help

	return nil
}

func (h *Handler) RegisterCommandAlias(command, alias string) error {
	if _, ok := h.commandAliases[alias]; ok {
		return fmt.Errorf("alias already registered")
	}

	if _, ok := h.commands[command]; !ok {
		return fmt.Errorf("looking for command: %s", command)
	}

	h.commandAliases[alias] = command

	return nil
}

func (h *Handler) FindCommand(name string) (Callback, error) {
	alias, ok := h.commandAliases[name]
	if ok {
		name = alias
	}

	command, ok := h.commands[name]
	if !ok {
		return nil, fmt.Errorf("looking for command: %s", name)
	}

	return command, nil
}

func (h *Handler) ParseArgs(session *discordgo.Session, message *discordgo.MessageCreate) []string {
	content := strings.TrimPrefix(message.Content, h.prefix)
	_, args, _ := strings.Cut(content, " ")

	return strings.Split(args, " ")
}

func (h *Handler) ParseCommand(session *discordgo.Session, message *discordgo.MessageCreate) string {
	content := strings.TrimPrefix(message.Content, h.prefix)
	command, _, _ := strings.Cut(content, " ")

	return command
}

func (h *Handler) handleCommand(session *discordgo.Session, message *discordgo.MessageCreate) {
	defer func() {
		if r := recover(); r != nil {
			h.handleError(session, message, fmt.Errorf("panic: %v", r))
		}
	}()

	if message.Author.ID == session.State.User.ID {
		return
	}

	if message.Author.Bot {
		return
	}

	if !strings.HasPrefix(message.Content, h.prefix) {
		return
	}

	commandName := h.ParseCommand(session, message)
	switch commandName {
	case "help":
		h.handleHelp(session, message)

	default:
		callback, err := h.FindCommand(commandName)
		if err != nil {
			return
		}

		if err := callback(session, message); err != nil {
			h.handleError(session, message, fmt.Errorf("call callback: %s", err))
		}
	}
}

func (h *Handler) handleHelp(session *discordgo.Session, message *discordgo.MessageCreate) {
	var commands []string

	for command := range h.commands {
		help, ok := h.commandHelp[command]
		if !ok {
			help = "No help message provided"
		}

		commands = append(commands, fmt.Sprintf("`%s`: %s", command, help))
	}

	_, _ = session.ChannelMessageSendComplex(message.ChannelID, &discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{
			Title:       "Help",
			Description: strings.Join(commands, "\n"),
			Color:       utils.RGBToDiscord(255, 255, 255),
			Footer: &discordgo.MessageEmbedFooter{
				Text: "command: help message",
			},
		},
		Reference: message.Reference(),
	})
}

func (h *Handler) handleError(session *discordgo.Session, message *discordgo.MessageCreate, err error) {
	fmt.Println(err)

	_, _ = session.ChannelMessageSendComplex(message.ChannelID, &discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{
			Title:       "Error",
			Description: err.Error(),
			Color:       utils.RGBToDiscord(255, 0, 0),
		},
		Reference: message.Reference(),
	})
}
