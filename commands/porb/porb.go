package porb

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Fluffy-Bean/lynxie/app"
	"github.com/Fluffy-Bean/lynxie/utils"
	"github.com/bwmarrin/discordgo"
)

var client = http.Client{
	Timeout: 10 * time.Second,
}
var username = os.Getenv("E621_USERNAME")
var password = os.Getenv("E621_PASSWORD")

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

func RegisterPorbCommands(a *app.App) {
	if username != "" && password != "" {
		a.RegisterCommand("e621", registerE621(a))

		a.RegisterCommandAlias("porb", "e621")
	} else {
		log.Println("Not registering e621 command...")
	}
}

func registerE621(a *app.App) app.Callback {
	return func(h *app.Handler, args []string) app.Error {
		var options struct {
			Tags   string
			Order  string
			Rating string
		}

		cmd := flag.NewFlagSet("", flag.ContinueOnError)

		cmd.StringVar(&options.Order, "order", "random", "Search order")
		cmd.StringVar(&options.Rating, "rating", "e", "Search rating")
		cmd.StringVar(&options.Tags, "tags", "", "Search tags")

		cmd.Parse(args)

		req, err := http.NewRequest(
			http.MethodGet,
			fmt.Sprintf(
				"https://e621.net/posts.json/?limit=1&tags=order:%s+rating:%s+%s",
				options.Order,
				options.Rating,
				options.Tags,
			),
			nil,
		)
		if err != nil {
			return app.Error{
				Msg: "Failed to make request",
				Err: err,
			}
		}

		req.Header.Add("Accept", "application/json")
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("User-Agent", fmt.Sprintf("Lynxie/2.0 (by %s on e621)", username))
		req.SetBasicAuth(username, password)

		res, err := client.Do(req)
		if err != nil {
			return app.Error{
				Msg: "Failed to do request",
				Err: err,
			}
		}
		defer res.Body.Close()

		var data struct {
			Posts []post `json:"posts"`
		}
		json.NewDecoder(res.Body).Decode(&data)

		if len(data.Posts) == 0 {
			return app.Error{
				Msg: "No posts found",
				Err: fmt.Errorf("no posts found"),
			}
		}

		var description string
		if len(data.Posts[0].Description) > 0 {
			description = data.Posts[0].Description
		} else {
			description = "No description provided."
		}

		var generalTags string
		if len(data.Posts[0].Tags.General) > 0 {
			generalTags = strings.Join(data.Posts[0].Tags.General[:20], ", ")
		} else {
			generalTags = "No tags provided."
		}

		h.Session.ChannelMessageSendComplex(h.Message.ChannelID, &discordgo.MessageSend{
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
					{
						Name:   "Tag(s)",
						Value:  generalTags,
						Inline: false,
					},
				},
				Image: &discordgo.MessageEmbedImage{
					URL: data.Posts[0].File.Url,
				},
				Footer: &discordgo.MessageEmbedFooter{
					Text: fmt.Sprintf(
						"ID: %d | Created: %s",
						data.Posts[0].Id,
						data.Posts[0].CreatedAt.Format(time.DateTime),
					),
				},
				Color: utils.ColorFromRGB(255, 255, 255),
			},
			Reference: h.Reference,
		})

		return app.Error{}
	}
}
