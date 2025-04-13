package main

import (
	"os"

	"github.com/Fluffy-Bean/lynxie/app"
	"github.com/Fluffy-Bean/lynxie/commands/debug"
	"github.com/Fluffy-Bean/lynxie/commands/img"
	"github.com/Fluffy-Bean/lynxie/commands/porb"
	"github.com/Fluffy-Bean/lynxie/commands/tinyfox"
	"github.com/bwmarrin/discordgo"
)

func main() {
	a := app.NewApp(app.Config{
		Prefix:  ">",
		Token:   os.Getenv("TOKEN"),
		Intents: discordgo.IntentsGuildMessages,
	})

	debug.RegisterDebugCommands(a)
	img.RegisterImgCommands(a)
	tinyfox.RegisterTinyfoxCommands(a)
	porb.RegisterPorbCommands(a)

	a.Run()
}
