package tinyfox

import (
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/Fluffy-Bean/lynxie/internal/color"
	"github.com/Fluffy-Bean/lynxie/internal/errors"
	"github.com/Fluffy-Bean/lynxie/internal/handler"
	"github.com/bwmarrin/discordgo"
)

var client = http.Client{
	Timeout: 10 * time.Second,
}

var animals = []string{
	"fox",
	"yeen",
	"dog",
	"guara",
	"serval",
	"ott",
	"jackal",
	"bleat",
	"woof",
	"chi",
	"puma",
	"skunk",
	"tig",
	"wah",
	"manul",
	"snep",
	"jaguar",
	"badger",
	"chee",
	"racc",
	"bear",
	"capy",
	"bun",
	"marten",
	"caracal",
	"snek",
	"shiba",
	"dook",
	"leo",
	"yote",
	"poss",
	"lynx",
}

var animalAliases = map[string]string{
	"hyena":        "yeen",
	"serv":         "serval",
	"otter":        "ott",
	"deer":         "bleat",
	"wolf":         "woof",
	"tiger":        "tig",
	"red-panda":    "wah",
	"panda":        "wah",
	"manual":       "manul",
	"palas":        "manul",
	"palas-cat":    "manul",
	"snow-leopard": "snep",
	"jag":          "jaguar",
	"cheetah":      "chee",
	"raccoon":      "racc",
	"rac":          "racc",
	"capybara":     "capy",
	"bunny":        "bun",
	"carac":        "caracal",
	"snake":        "snek",
	"ferret":       "dook",
	"leopard":      "leo",
	"coyote":       "yote",
	"possum":       "poss",
	"opossum":      "poss",
}

func RegisterTinyfoxCommands(bot *handler.Bot) {
	bot.RegisterCommand("animal", registerAnimal(bot))

	bot.RegisterCommandAlias("bot", "animal")
}

func registerAnimal(bot *handler.Bot) handler.Callback {
	return func(h *handler.Handler, args []string) errors.Error {
		if len(args) < 1 {
			return errors.Error{
				Msg: "Animal name is required!",
				Err: fmt.Errorf("animal name is required"),
			}
		}

		animal := args[0]

		if !slices.Contains(animals, animal) {
			alias, ok := animalAliases[animal]
			if !ok {
				return errors.Error{
					Msg: fmt.Sprintf("Animal \"%s\" is invalid. The following animals are supported:\n%s", animal, strings.Join(animals, ", ")),
					Err: fmt.Errorf("entered invalid animal name"),
				}
			}
			animal = alias
		}

		req, err := http.NewRequest(http.MethodGet, "https://api.tinyfox.dev/img?animal="+animal, nil)
		if err != nil {
			return errors.Error{
				Msg: "Failed to make request",
				Err: err,
			}
		}

		res, err := client.Do(req)
		if err != nil {
			return errors.Error{
				Msg: "Failed to do request",
				Err: err,
			}
		}
		defer res.Body.Close()

		_, err = h.Session.ChannelMessageSendComplex(h.Message.ChannelID, &discordgo.MessageSend{
			Embed: &discordgo.MessageEmbed{
				Title: "Animal",
				Image: &discordgo.MessageEmbedImage{
					URL: "attachment://animal.png",
				},
				Color: color.RGBToDiscord(255, 255, 255),
			},
			Files: []*discordgo.File{
				{
					Name:        "animal.png",
					ContentType: "",
					Reader:      res.Body,
				},
			},
			Reference: h.Reference,
		})
		if err != nil {
			return errors.Error{
				Msg: "failed to send tinyfox message",
				Err: err,
			}
		}

		return errors.Error{}
	}
}
