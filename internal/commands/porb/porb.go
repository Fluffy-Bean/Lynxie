package porb

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/Fluffy-Bean/lynxie/internal/bot"
	"github.com/Fluffy-Bean/lynxie/internal/color"
)

var client = http.Client{
	Timeout: 10 * time.Second,
}

func RegisterPorbCommands(h *bot.Handler) {
	_ = h.RegisterCommand("porb", cmdPorb(h))
	_ = h.RegisterCommandAlias("e621", "porb")

	_ = h.RegisterCommand("yiff", cmdYiff(h))
	_ = h.RegisterCommandAlias("yiff.gay", "yiff")
}

func cmdPorb(h *bot.Handler) bot.Command {
	type post struct {
		Id        int       `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		File      struct {
			Width  int    `json:"width"`
			Height int    `json:"height"`
			Ext    string `json:"ext"`
			Size   int    `json:"size"`
			Md5    string `json:"md5"`
			Url    string `json:"url"`
		} `json:"file"`
		Score struct {
			Up    int `json:"up"`
			Down  int `json:"down"`
			Total int `json:"total"`
		} `json:"score"`
		Tags struct {
			General     []string      `json:"general"`
			Artist      []string      `json:"artist"`
			Contributor []interface{} `json:"contributor"`
			Copyright   []string      `json:"copyright"`
			Character   []interface{} `json:"character"`
			Species     []string      `json:"species"`
			Invalid     []interface{} `json:"invalid"`
			Meta        []string      `json:"meta"`
			Lore        []interface{} `json:"lore"`
		} `json:"tags"`
		Rating       string   `json:"rating"`
		FavCount     int      `json:"fav_count"`
		Sources      []string `json:"sources"`
		Description  string   `json:"description"`
		CommentCount int      `json:"comment_count"`
	}

	return func(h *bot.Handler, c bot.CommandContext) error {
		var options struct {
			Order  string
			Rating string
		}

		cmd := flag.NewFlagSet("", flag.ContinueOnError)

		cmd.StringVar(&options.Order, "order", "random", "Search order")
		cmd.StringVar(&options.Rating, "rating", "e", "Search rating")

		err := cmd.Parse(h.ParseArgs(c))
		if err != nil {
			return fmt.Errorf("parse flags: %w", err)
		}

		url := fmt.Sprintf(
			"https://e621.net/posts.json/?limit=1&tags=order:%s+rating:%s+%s",
			options.Order,
			options.Rating,
			strings.Join(cmd.Args(), "+"),
		)

		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return fmt.Errorf("create request: %w", err)
		}

		req.Header.Add("Accept", "application/json")
		req.Header.Add("Content-Type", "application/json")

		res, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("do request: %w", err)
		}
		defer res.Body.Close()

		var data struct {
			Posts []post `json:"posts"`
		}
		err = json.NewDecoder(res.Body).Decode(&data)
		if err != nil {
			return fmt.Errorf("decode response: %w", err)
		}

		if len(data.Posts) == 0 {
			return fmt.Errorf("no posts found")
		}

		var description string
		if len(data.Posts[0].Description) > 0 {
			description = data.Posts[0].Description
		} else {
			description = "No description provided."
		}

		_, err = c.Session.ChannelMessageSendComplex(c.Message.ChannelID, &discordgo.MessageSend{
			Embed: &discordgo.MessageEmbed{
				Title:       "E621",
				Description: description,
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:  "Score",
						Value: fmt.Sprintf("⬆️ %d | ⬇️ %d", data.Posts[0].Score.Up, data.Posts[0].Score.Down),
					},
					{
						Name:  "Favorites",
						Value: fmt.Sprintf("%d", data.Posts[0].FavCount),
					},
					{
						Name:  "Comments",
						Value: fmt.Sprintf("%d", data.Posts[0].CommentCount),
					},
					{
						Name:   "Source(s)",
						Value:  strings.Join(data.Posts[0].Sources, ", "),
						Inline: false,
					},
				},
				Image: &discordgo.MessageEmbedImage{
					URL: data.Posts[0].File.Url,
				},
				Footer: &discordgo.MessageEmbedFooter{
					Text: fmt.Sprintf("ID: %d | Created: %s", data.Posts[0].Id, data.Posts[0].CreatedAt.Format(time.DateTime)),
				},
				Color: color.RGBToDiscord(255, 255, 255),
			},
			Reference: c.Message.Reference(),
		})
		if err != nil {
			return fmt.Errorf("send e621 response: %w", err)
		}

		return nil
	}
}

func cmdYiff(h *bot.Handler) bot.Command {
	var client = http.Client{
		Timeout: 10 * time.Second,

		// The sites certificates seems to be currently fucked...
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	return func(h *bot.Handler, c bot.CommandContext) error {
		req, err := http.NewRequest(http.MethodGet, "https://yiff.gay", nil)
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
				Title: "yiff.gay",
				Image: &discordgo.MessageEmbedImage{
					URL: "attachment://yiff.png",
				},
				Color: color.RGBToDiscord(255, 255, 255),
			},
			Files: []*discordgo.File{
				{
					Name:        "yiff.png",
					ContentType: "",
					Reader:      res.Body,
				},
			},
			Reference: c.Message.Reference(),
		})
		if err != nil {
			return fmt.Errorf("send yiff response: %w", err)
		}

		return nil
	}
}
