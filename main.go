package main

import (
	"os"

	"github.com/Fluffy-Bean/lynxie/app"
	"github.com/Fluffy-Bean/lynxie/commands"
	"github.com/bwmarrin/discordgo"
)

func main() {
	a := app.NewApp(app.Config{
		Prefix:  ">",
		Token:   os.Getenv("TOKEN"),
		Intents: discordgo.IntentsGuildMessages,
	})

	commands.RegisterMetaCommands(a)
	commands.RegisterTinyfoxCommands(a)

	a.Run()
}
