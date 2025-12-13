package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/Fluffy-Bean/lynxie/internal/commands/debug"
	"github.com/Fluffy-Bean/lynxie/internal/commands/fortnite"
	"github.com/Fluffy-Bean/lynxie/internal/commands/img"
	"github.com/Fluffy-Bean/lynxie/internal/commands/jellyfin"
	"github.com/Fluffy-Bean/lynxie/internal/commands/porb"
	"github.com/Fluffy-Bean/lynxie/internal/commands/tinyfox"
	"github.com/Fluffy-Bean/lynxie/internal/handler"
)

var flags struct {
	Token        *string
	Prefix       *string
	DataPath     *string
	Experimental *bool
}

func main() {
	err := parseFlags()
	if err != nil {
		log.Fatal("parse flags:", err)
	}

	databasePath := path.Join(*flags.DataPath, "storage.db")

	db, err := sql.Open("sqlite3", databasePath)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	h := handler.NewHandler(*flags.Prefix, db)

	if *flags.Experimental {
		fmt.Println("Running in experimental mode")

		h.SetValue("JellyfinToken", os.Getenv("JELLYFIN_TOKEN"))
		h.SetValue("JellyfinEndpoint", os.Getenv("JELLYFIN_ENDPOINT"))

		jellyfin.RegisterCommands(h)
	}

	debug.RegisterCommands(h)
	img.RegisterCommands(h)
	tinyfox.RegisterCommands(h)
	porb.RegisterCommands(h)
	fortnite.RegisterCommands(h)

	log.Fatal(h.Run(os.Getenv("TOKEN")))
}

func parseFlags() error {
	flags.Prefix = flag.String("prefix", ">", "The prefix the bot will search for within messages")
	flags.DataPath = flag.String("datapath", "", "The path to the datapath, defaults to $HOME/.lynxie if empty")
	flags.Token = flag.String("token", "", "The bot token, defaults to checking for TOKEN in env")
	flags.Experimental = flag.Bool("experimental", false, "Whether or not the bot is in experimental mode")

	flag.Parse()

	if *flags.DataPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return err
		}

		*flags.DataPath = path.Join(homeDir, ".lynxie", "data")
	}

	if *flags.Token == "" {
		*flags.Token = os.Getenv("TOKEN")
	}

	return nil
}
