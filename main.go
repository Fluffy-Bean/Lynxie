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

var ConfigBuildHash string
var ConfigBuildPipeline string

func main() {
	a := app.NewApp(app.Config{
		BotPrefix:  ">",
		BotToken:   os.Getenv("TOKEN"),
		BotIntents: discordgo.IntentsGuildMessages,
		CommandExtras: map[string]string{
			"debug_build-hash":     ConfigBuildHash,
			"debug_build-pipeline": ConfigBuildPipeline,
			"e621_username":        os.Getenv("E621_USERNAME"),
			"e621_password":        os.Getenv("E621_PASSWORD"),
		},
	})

	debug.RegisterDebugCommands(a)
	img.RegisterImgCommands(a)
	tinyfox.RegisterTinyfoxCommands(a)
	porb.RegisterPorbCommands(a)

	a.Run()
}
