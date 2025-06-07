package main

import (
	"os"

	"github.com/Fluffy-Bean/lynxie/internal/handler"
	"github.com/Fluffy-Bean/lynxie/pkg/commands/debug"
	"github.com/Fluffy-Bean/lynxie/pkg/commands/img"
	"github.com/Fluffy-Bean/lynxie/pkg/commands/porb"
	"github.com/Fluffy-Bean/lynxie/pkg/commands/tinyfox"
	"github.com/bwmarrin/discordgo"
)

func main() {
	bot := handler.NewBot("?", os.Getenv("TOKEN"), discordgo.IntentsGuildMessages)

	debug.RegisterDebugCommands(bot)
	img.RegisterImgCommands(bot)
	tinyfox.RegisterTinyfoxCommands(bot)
	porb.RegisterPorbCommands(bot)

	bot.Run()
}
