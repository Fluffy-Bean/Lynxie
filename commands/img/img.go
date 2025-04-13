package img

import (
	"errors"
	"net/http"
	"time"

	"github.com/Fluffy-Bean/lynxie/app"
	"github.com/Fluffy-Bean/lynxie/utils"
	"github.com/bwmarrin/discordgo"
)

const maxFileSize = 1024 * 1024 * 10 // 10MB

var client = http.Client{
	Timeout: 10 * time.Second,
}

func RegisterImgCommands(a *app.App) {
	a.RegisterCommand("saveable", registerSaveable(a))

	a.RegisterCommandAlias("gif", "saveable")
}

func registerSaveable(a *app.App) app.Callback {
	return func(h *app.Handler, args []string) app.Error {
		fileEndpoint, err := getClosestImage(h)
		if err != nil {
			return app.Error{
				Msg: "Could not get image",
				Err: err,
			}
		}

		req, err := http.NewRequest(http.MethodGet, fileEndpoint, nil)
		if err != nil {
			return app.Error{
				Msg: "",
				Err: err,
			}
		}

		if req.ContentLength > maxFileSize {
			return app.Error{
				Msg: "Could not get image",
				Err: errors.New("requested file is too big"),
			}
		}

		res, err := client.Do(req)
		if err != nil {
			return app.Error{
				Msg: "",
				Err: err,
			}
		}
		defer res.Body.Close()

		h.Session.ChannelMessageSendComplex(h.Message.ChannelID, &discordgo.MessageSend{
			Embed: &discordgo.MessageEmbed{
				Title:       "Saveable",
				Description: "Image converted to GIF :3",
				Image: &discordgo.MessageEmbedImage{
					URL: "attachment://image.gif",
				},
				Color: utils.ColorFromRGB(255, 255, 255),
			},
			Files: []*discordgo.File{
				{
					Name:        "image.gif",
					ContentType: "image/gif",
					Reader:      res.Body,
				},
			},
			Reference: h.Reference,
		})

		return app.Error{}
	}
}

func getClosestImage(h *app.Handler) (string, error) {
	// Get message attachments
	if len(h.Message.Attachments) >= 1 {
		if h.Message.Attachments[0].Size > maxFileSize {
			return "", errors.New("file size is too big")
		}

		return h.Message.Attachments[0].ProxyURL, nil
	}

	// If no attachments exist... see if the message is replying to someone
	if h.Message.ReferencedMessage != nil {
		if len(h.Message.ReferencedMessage.Attachments) >= 1 {
			if h.Message.ReferencedMessage.Attachments[0].Size > maxFileSize {
				return "", errors.New("file size is too big")
			}

			return h.Message.ReferencedMessage.Attachments[0].ProxyURL, nil
		}

		// Maybe replying to an embed...?
		if len(h.Message.ReferencedMessage.Embeds) >= 1 {
			//... no file size is provided
			return h.Message.ReferencedMessage.Embeds[0].Image.ProxyURL, nil
		}
	}

	return "", errors.New("no files exists")
}
