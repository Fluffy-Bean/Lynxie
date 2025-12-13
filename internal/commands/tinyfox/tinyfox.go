package tinyfox

import (
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/Fluffy-Bean/lynxie/internal/handler"
	"github.com/Fluffy-Bean/lynxie/internal/utils"
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
	"leg":          "guara",
	"leggy":        "guara",
	"maned-wolf":   "guara",
	"maney":        "guara",
	"mane":         "guara",
	"lobo":         "guara",
	"painted":      "chi",
	"wad":          "chi",
}

func RegisterCommands(h *handler.Handler) {
	_ = h.RegisterCommand("tinyfox", cmdTinyFox(h))
}

func cmdTinyFox(h *handler.Handler) handler.Callback {
	_ = h.RegisterCommandHelp("tinyfox", "Search for silly animal images")
	_ = h.RegisterCommandAlias("tinyfox", "animal")
	_ = h.RegisterCommandAlias("tinyfox", "a")

	return func(session *discordgo.Session, message *discordgo.MessageCreate) error {
		if len(h.ParseArgs(session, message)) < 1 {
			return fmt.Errorf("animal name not provided")
		}

		animal := h.ParseArgs(session, message)[0]

		alias, ok := animalAliases[animal]
		if ok {
			animal = alias
		}

		if !slices.Contains(animals, animal) {
			return fmt.Errorf("unknown animal %s", animal)
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

		// API responds with JSON rather than image data if there was an incorrect request done
		if res.Header.Get("Content-Type") != "image/jpeg" {
			return fmt.Errorf("content type: possibly unknown animal %s", animal)
		}

		_, err = session.ChannelMessageSendComplex(message.ChannelID, &discordgo.MessageSend{
			Embed: &discordgo.MessageEmbed{
				Title: "Animal",
				Image: &discordgo.MessageEmbedImage{
					URL: "attachment://animal.png",
				},
				Color: utils.RGBToDiscord(255, 255, 255),
			},
			Files: []*discordgo.File{
				{
					Name:        "animal.png",
					ContentType: "",
					Reader:      res.Body,
				},
			},
			Reference: message.Reference(),
		})
		if err != nil {
			return fmt.Errorf("send tinyfox response: %w", err)
		}

		return nil
	}
}
