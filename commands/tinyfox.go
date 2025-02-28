package commands

import (
	"flag"
	"net/http"
	"time"

	"github.com/Fluffy-Bean/lynxie/app"
)

func RegisterTinyfoxCommands(a *app.App) {
	a.RegisterCommand("animal", registerAnimal(a))
}

func registerAnimal(a *app.App) func(h *app.Handler, args []string) {
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	return func(h *app.Handler, args []string) {
		var options struct {
			animal string
		}

		cmd := flag.NewFlagSet("pong", flag.ContinueOnError)
		cmd.StringVar(&options.animal, "animal", "wah", "Get an image of an animal!")
		cmd.Parse(args)

		req, err := http.NewRequest(http.MethodGet, "https://api.tinyfox.dev/img?animal="+options.animal, nil)
		if err != nil {
			return
		}

		res, err := client.Do(req)
		if err != nil {
			return
		}
		defer res.Body.Close()

		h.Session.ChannelFileSend(
			h.Message.ChannelID,
			"animal__"+options.animal+".png",
			res.Body,
		)
	}
}
