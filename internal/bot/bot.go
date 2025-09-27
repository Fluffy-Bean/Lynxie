package bot

import (
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	_ "github.com/mattn/go-sqlite3"
)

type CommandContext struct {
	Session *discordgo.Session
	Message *discordgo.MessageCreate
}

type Command func(handler *Handler, context CommandContext) error

type HelpCallback func(handler *Handler, context CommandContext)
type ErrorCallback func(handler *Handler, context CommandContext, err error)

func NewHandler(prefix string, onHelp HelpCallback, onError ErrorCallback) *Handler {
	return &Handler{
		prefix: prefix,

		commands:       make(map[string]Command),
		commandAliases: make(map[string]string),

		onHelp:  onHelp,
		onError: onError,
	}
}

type Handler struct {
	prefix string

	db *sql.DB

	commands       map[string]Command
	commandAliases map[string]string

	onHelp  HelpCallback
	onError ErrorCallback
}

func (h *Handler) GetDB() *sql.DB {
	return h.db
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

func (h *Handler) RegisterCommand(command string, f Command) error {
	if _, ok := h.commands[command]; ok {
		return fmt.Errorf("command already registered")
	}

	h.commands[command] = f

	return nil
}

func (h *Handler) GetCommands() map[string]Command {
	return h.commands
}

func (h *Handler) RegisterCommandAlias(alias, command string) error {
	if _, ok := h.commandAliases[alias]; ok {
		return fmt.Errorf("alias already registered")
	}

	if _, ok := h.commands[command]; !ok {
		return fmt.Errorf("looking for command: %s", command)
	}

	h.commandAliases[alias] = command

	return nil
}

func (h *Handler) GetCommandAliases() map[string]string {
	return h.commandAliases
}

func (h *Handler) FindCommand(name string) (Command, error) {
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

func (h *Handler) ParseArgs(c CommandContext) []string {
	message := c.Message.Content

	message = strings.TrimPrefix(message, h.prefix)
	_, args, _ := strings.Cut(message, " ")

	return strings.Split(args, " ")
}

func (h *Handler) ParseCommand(c CommandContext) string {
	message := c.Message.Content

	message = strings.TrimPrefix(message, h.prefix)
	command, _, _ := strings.Cut(message, " ")

	return command
}

func (h *Handler) handleCommand(session *discordgo.Session, message *discordgo.MessageCreate) {
	context := CommandContext{
		Session: session,
		Message: message,
	}

	defer func() {
		if r := recover(); r != nil {
			h.onError(h, context, fmt.Errorf("panic: %v", r))
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

	commandName := h.ParseCommand(context)
	if commandName == "help" {
		h.onHelp(h, context)

		return
	}

	command, err := h.FindCommand(commandName)
	if err != nil {
		return
	}

	err = command(h, context)
	if err != nil {
		h.onError(h, context, fmt.Errorf("call command: %s", err))
	}
}

func (h *Handler) Run(databasePath, token string, intent discordgo.Intent) error {
	var err error

	h.db, err = sql.Open("sqlite3", databasePath)
	if err != nil {
		return fmt.Errorf("opening database: %v", err)
	}
	defer h.db.Close()

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		return fmt.Errorf("creating discord session: %s", err)
	}

	dg.Identify.Intents = intent
	dg.AddHandler(h.handleCommand)

	err = dg.Open()
	if err != nil {
		return fmt.Errorf("opening discord session: %s", err)
	}
	defer dg.Close()

	fmt.Println("Bot is now running. Press CTRL-C to exit")
	fmt.Println("prefix ............... ", h.prefix)
	fmt.Println("database path ........ ", databasePath)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	<-sc

	return nil
}
