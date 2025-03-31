package tinyfox

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/Fluffy-Bean/lynxie/app"
	"github.com/Fluffy-Bean/lynxie/utils"
	"github.com/bwmarrin/discordgo"
)

var client = http.Client{
	Timeout: 10 * time.Second,
}

func RegisterTinyfoxCommands(a *app.App) {
	a.RegisterCommand("animal", registerAnimal(a))
}

var animals = []string{
	"fox", "yeen", "dog", "guara", "serval", "ott", "jackal", "bleat", "woof", "chi", "puma", "skunk", "tig", "wah",
	"manul", "snep", "jaguar", "badger", "chee", "racc", "bear", "capy", "bun", "marten", "caracal", "snek",
	"shiba", "dook", "leo", "yote", "poss", "chee", "lynx",
}

func registerAnimal(a *app.App) app.Callback {
	return func(h *app.Handler, args []string) app.Error {
		var options struct {
			Kind string
		}

		cmd := flag.NewFlagSet("", flag.ContinueOnError)

		cmd.StringVar(&options.Kind, "kind", "", "Animal kind to search for")

		cmd.Parse(args)

		if options.Kind == "" {
			return app.Error{
				Msg: "Animal name is required!",
				Err: errors.New("animal name is required"),
			}
		}
		if !slices.Contains(animals, options.Kind) {
			return app.Error{
				Msg: fmt.Sprintf("Animal %s is invalid", options.Kind),
				Err: errors.New("entered invalid animal name"),
			}
		}

		req, err := http.NewRequest(http.MethodGet, "https://api.tinyfox.dev/img?animal="+options.Kind, nil)
		if err != nil {
			return app.Error{
				Msg: "Failed to make request",
				Err: err,
			}
		}

		res, err := client.Do(req)
		if err != nil {
			return app.Error{
				Msg: "Failed to do request",
				Err: err,
			}
		}
		defer res.Body.Close()

		h.Session.ChannelMessageSendComplex(h.Message.ChannelID, &discordgo.MessageSend{
			Embed: &discordgo.MessageEmbed{
				Title: "Animal",
				Image: &discordgo.MessageEmbedImage{
					URL: "attachment://image.png",
				},
				Color: utils.ColorFromRGB(255, 255, 255),
			},
			Files: []*discordgo.File{
				{
					Name:        "image.png",
					ContentType: "",
					Reader:      res.Body,
				},
			},
			Reference: h.Reference,
		})

		return app.Error{}
	}
}
