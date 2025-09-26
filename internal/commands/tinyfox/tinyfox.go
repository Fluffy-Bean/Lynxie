package tinyfox

import (
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/Fluffy-Bean/lynxie/internal/bot"
	"github.com/Fluffy-Bean/lynxie/internal/color"
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
	"leg":          "gaura",
	"leggy":        "gaura",
	"maned-wolf":   "gaura",
	"maney":        "gaura",
	"mane":         "gaura",
	"lobo":         "gaura",
	"painted":      "chi",
	"wad":          "chi",
}

func RegisterTinyfoxCommands(h *bot.Handler) {
	_ = h.RegisterCommand("tinyfox", cmdTinyFox(h))
	_ = h.RegisterCommandAlias("animal", "tinyfox")
	_ = h.RegisterCommandAlias("a", "tinyfox")
}

func cmdTinyFox(h *bot.Handler) bot.Command {
	return func(h *bot.Handler, c bot.CommandContext) error {
		if len(h.ParseArgs(c)) < 1 {
			return fmt.Errorf("animal name not provided")
		}

		animal := h.ParseArgs(c)[0]

		if !slices.Contains(animals, animal) {
			alias, ok := animalAliases[animal]
			if !ok {
				return fmt.Errorf("unknown animal %s", animal)
			}

			animal = alias
		}

		req, err := http.NewRequest(http.MethodGet, "https://api.tinyfox.dev/img?animal="+animal, nil)
		if err != nil {
			return fmt.Errorf("create request: %w", err)
		}

		res, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("do request: %w", err)
		}
		defer res.Body.Close()

		_, err = c.Session.ChannelMessageSendComplex(c.Message.ChannelID, &discordgo.MessageSend{
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
			Reference: c.Message.Reference(),
		})
		if err != nil {
			return fmt.Errorf("send tinyfox response: %w", err)
		}

		return nil
	}
}
