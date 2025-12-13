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

	"github.com/Fluffy-Bean/lynxie/internal/handler"
	"github.com/Fluffy-Bean/lynxie/internal/utils"
)

func RegisterCommands(h *handler.Handler) {
	_ = h.RegisterCommand("porb", cmdPorb(h))
	_ = h.RegisterCommand("yiff", cmdYiff(h))
}

func cmdPorb(h *handler.Handler) handler.Callback {
	_ = h.RegisterCommandHelp("porb", "Search e621")
	_ = h.RegisterCommandAlias("porb", "e621")

	client := http.Client{
		Timeout: 10 * time.Second,
	}

	type apiPost struct {
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
		Preview struct {
			Width  int    `json:"width"`
			Height int    `json:"height"`
			Url    string `json:"url"`
			Alt    string `json:"alt"`
		} `json:"preview"`
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

	return func(session *discordgo.Session, message *discordgo.MessageCreate) error {
		var options struct {
			Order  string
			Rating string
		}

		cmd := flag.NewFlagSet("", flag.ContinueOnError)

		cmd.StringVar(&options.Order, "order", "random", "Search order")
		cmd.StringVar(&options.Rating, "rating", "e", "Search rating")

		err := cmd.Parse(h.ParseArgs(session, message))
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
			Posts []apiPost `json:"posts"`
		}
		err = json.NewDecoder(res.Body).Decode(&data)
		if err != nil {
			return fmt.Errorf("decode response: %w", err)
		}

		if len(data.Posts) == 0 {
			return fmt.Errorf("no posts found")
		}

		post := data.Posts[0]

		msg := discordgo.MessageSend{
			Embed: &discordgo.MessageEmbed{
				Title:       "E621",
				Description: "No description available",
				Fields:      []*discordgo.MessageEmbedField{},
				Color:       utils.RGBToDiscord(255, 255, 255),
			},
			Reference: message.Reference(),
		}

		if post.Description != "" {
			msg.Embed.Description = post.Description
		}

		msg.Embed.Fields = append(msg.Embed.Fields, &discordgo.MessageEmbedField{
			Name:  "Score",
			Value: fmt.Sprintf("⬆️ %d | ⬇️ %d", post.Score.Up, post.Score.Down),
		})

		msg.Embed.Fields = append(msg.Embed.Fields, &discordgo.MessageEmbedField{
			Name:  "Favorites",
			Value: fmt.Sprintf("%d", post.FavCount),
		})

		msg.Embed.Fields = append(msg.Embed.Fields, &discordgo.MessageEmbedField{
			Name:  "Comments",
			Value: fmt.Sprintf("%d", post.CommentCount),
		})

		switch len(post.Sources) {
		case 0:
			// noop

		case 1:
			msg.Embed.Fields = append(msg.Embed.Fields, &discordgo.MessageEmbedField{
				Name:   "Source",
				Value:  post.Sources[0],
				Inline: false,
			})

		default:
			msg.Embed.Fields = append(msg.Embed.Fields, &discordgo.MessageEmbedField{
				Name:   "Sources",
				Value:  strings.Join(post.Sources, ", "),
				Inline: false,
			})
		}

		if post.CreatedAt.IsZero() {
			msg.Embed.Footer = &discordgo.MessageEmbedFooter{
				Text: fmt.Sprintf("ID: %d", post.Id),
			}
		} else {
			msg.Embed.Footer = &discordgo.MessageEmbedFooter{
				Text: fmt.Sprintf("ID: %d | Created: %s", post.Id, post.CreatedAt.Format(time.DateTime)),
			}
		}

		switch post.File.Ext {
		case "jpg", "jpeg", "png", "gif":
			msg.Embed.Image = &discordgo.MessageEmbedImage{
				URL: post.File.Url,
			}

		case "mp4", "webm": // Unfortunately, discord doesn't currently allow bots to embed videos in responses :c
			msg.Embed.Image = &discordgo.MessageEmbedImage{
				URL: post.Preview.Url,
			}

		default:
			return fmt.Errorf("unknown response ext %s", post.File.Ext)
		}

		_, err = session.ChannelMessageSendComplex(message.ChannelID, &msg)
		if err != nil {
			return fmt.Errorf("send e621 response: %w", err)
		}

		return nil
	}
}

func cmdYiff(h *handler.Handler) handler.Callback {
	_ = h.RegisterCommandHelp("yiff", "Get a random image from yiff.gay")
	_ = h.RegisterCommandAlias("yiff", "yiff.gay")

	client := http.Client{
		Timeout: 10 * time.Second,

		// The sites certificates seems to be currently fucked...
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	return func(session *discordgo.Session, message *discordgo.MessageCreate) error {
		req, err := http.NewRequest(http.MethodGet, "https://yiff.gay", nil)
		if err != nil {
			return fmt.Errorf("create request: %w", err)
		}

		res, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("do request: %w", err)
		}
		defer res.Body.Close()

		_, err = session.ChannelMessageSendComplex(message.ChannelID, &discordgo.MessageSend{
			Embed: &discordgo.MessageEmbed{
				Title: "yiff.gay",
				Image: &discordgo.MessageEmbedImage{
					URL: "attachment://yiff.png",
				},
				Color: utils.RGBToDiscord(255, 255, 255),
			},
			Files: []*discordgo.File{
				{
					Name:        "yiff.png",
					ContentType: "",
					Reader:      res.Body,
				},
			},
			Reference: message.Reference(),
		})
		if err != nil {
			return fmt.Errorf("send yiff response: %w", err)
		}

		return nil
	}
}
