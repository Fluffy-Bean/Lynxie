package commands

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Fluffy-Bean/lynxie/app"
	"github.com/Fluffy-Bean/lynxie/utils"
	"github.com/bwmarrin/discordgo"
)

func RegisterPorbCommands(a *app.App) {
	a.RegisterCommand("e621", registerE621(a))
}

func registerE621(a *app.App) app.Callback {
	username := os.Getenv("E621_USERNAME")
	password := os.Getenv("E621_PASSWORD")

	client := http.Client{
		Timeout: 10 * time.Second,
	}

	return func(h *app.Handler, args []string) app.Error {
		var options struct {
			tags   string
			order  string
			rating string
		}

		cmd := flag.NewFlagSet("", flag.ContinueOnError)
		cmd.StringVar(&options.order, "order", "random", "Search order")
		cmd.StringVar(&options.rating, "rating", "e", "Search rating")
		cmd.StringVar(&options.tags, "tags", "", "Search tags")
		cmd.Parse(args)

		req, err := http.NewRequest(
			http.MethodGet,
			fmt.Sprintf(
				"https://e621.net/posts.json/?limit=1&tags=order:%s+rating:%s+%s",
				options.order,
				options.rating,
				options.tags,
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
		req.Header.Add("User-Agent", fmt.Sprintf("Lynxie/1.0 (by %s on e621)", username))
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
			Posts []struct {
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
				} `json:"preview"`
				Sample struct {
					Has        bool   `json:"has"`
					Height     int    `json:"height"`
					Width      int    `json:"width"`
					Url        string `json:"url"`
					Alternates struct {
					} `json:"alternates"`
				} `json:"sample"`
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
				LockedTags []interface{} `json:"locked_tags"`
				ChangeSeq  int           `json:"change_seq"`
				Flags      struct {
					Pending      bool `json:"pending"`
					Flagged      bool `json:"flagged"`
					NoteLocked   bool `json:"note_locked"`
					StatusLocked bool `json:"status_locked"`
					RatingLocked bool `json:"rating_locked"`
					Deleted      bool `json:"deleted"`
				} `json:"flags"`
				Rating        string   `json:"rating"`
				FavCount      int      `json:"fav_count"`
				Sources       []string `json:"sources"`
				Pools         []int    `json:"pools"`
				Relationships struct {
					ParentId          interface{}   `json:"parent_id"`
					HasChildren       bool          `json:"has_children"`
					HasActiveChildren bool          `json:"has_active_children"`
					Children          []interface{} `json:"children"`
				} `json:"relationships"`
				ApproverId   interface{} `json:"approver_id"`
				UploaderId   int         `json:"uploader_id"`
				Description  string      `json:"description"`
				CommentCount int         `json:"comment_count"`
				IsFavorited  bool        `json:"is_favorited"`
				HasNotes     bool        `json:"has_notes"`
				Duration     interface{} `json:"duration"`
			} `json:"posts"`
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
