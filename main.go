package main

import (
	"os"

	"github.com/Fluffy-Bean/lynxie/app"
	"github.com/Fluffy-Bean/lynxie/commands"
)

func main() {
	a := app.NewApp(app.Config{
		Token:  os.Getenv("TOKEN"),
		Prefix: "?",
	})

	commands.RegisterMetaCommands(a)
	commands.RegisterTinyfoxCommands(a)

	a.Run()
}
