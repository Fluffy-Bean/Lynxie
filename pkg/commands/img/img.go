package img

import (
	"bufio"
	"bytes"
	_ "embed"
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"strings"
	"time"

	"git.sr.ht/~sbinet/gg"
	"github.com/Fluffy-Bean/lynxie/app"
	"github.com/Fluffy-Bean/lynxie/internal/color"
	"github.com/bwmarrin/discordgo"
)

const maxFileSize = 1024 * 1024 * 10 // 10MB

var client = http.Client{
	Timeout: 10 * time.Second,
}

//go:embed resources/Impact.ttf
var resourceImpactFont []byte

func RegisterImgCommands(a *app.App) {
	a.RegisterCommand("saveable", registerSaveable(a))
	a.RegisterCommandAlias("gif", "saveable")

	a.RegisterCommand("caption", registerCaption(a))
	a.RegisterCommandAlias("c", "caption")
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
					URL: "attachment://saveable.gif",
				},
				Color: color.RGBToDiscord(255, 255, 255),
			},
			Files: []*discordgo.File{
				{
					Name:        "saveable.gif",
					ContentType: "image/gif",
					Reader:      res.Body,
				},
			},
			Reference: h.Reference,
		})

		return app.Error{}
	}
}

func registerCaption(a *app.App) app.Callback {
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

		buff, err := io.ReadAll(res.Body)
		if err != nil {
			return app.Error{
				Msg: "Failed to read image",
				Err: err,
			}
		}

		var img image.Image
		switch http.DetectContentType(buff) {
		case "image/png":
			img, err = png.Decode(bytes.NewReader(buff))
			if err != nil {
				return app.Error{
					Msg: "Failed to decode PNG",
					Err: err,
				}
			}
			break
		case "image/jpeg":
			img, err = jpeg.Decode(bytes.NewReader(buff))
			if err != nil {
				return app.Error{
					Msg: "Failed to decode JPEG",
					Err: err,
				}
			}
			break
		default:
			return app.Error{
				Msg: "Unknown or unsupported image format",
				Err: errors.New("Unknown or unsupported image format " + http.DetectContentType(buff)),
			}
		}

		fontSize := float64(img.Bounds().Dx() / 25)
		if fontSize < 16 {
			fontSize = 16
		} else if fontSize > 50 {
			fontSize = 50
		}

		canvas := gg.NewContext(img.Bounds().Dx(), img.Bounds().Dy()+200)
		err = canvas.LoadFontFaceFromBytes(resourceImpactFont, fontSize)
		if err != nil {
			return app.Error{
				Msg: "Failed to load font",
				Err: err,
			}
		}

		canvas.SetRGBA(255, 255, 255, 255)
		canvas.Clear()

		canvas.SetRGBA(0, 0, 0, 255)
		canvas.DrawStringWrapped(
			strings.Join(args, " "),
			float64(img.Bounds().Dx()/2),
			100,
			0.5,
			0.5,
			float64(img.Bounds().Dx()),
			1.5,
			gg.AlignCenter,
		)

		canvas.DrawImage(img, 0, 200)

		var export bytes.Buffer
		err = canvas.EncodeJPG(bufio.NewWriter(&export), &jpeg.Options{
			Quality: 100,
		})
		if err != nil {
			return app.Error{
				Msg: "Failed to encode JPEG",
				Err: err,
			}
		}

		h.Session.ChannelMessageSendComplex(h.Message.ChannelID, &discordgo.MessageSend{
			Embed: &discordgo.MessageEmbed{
				Title: "Caption",
				Image: &discordgo.MessageEmbedImage{
					URL: "attachment://caption.jpeg",
				},
				Color: color.RGBToDiscord(255, 255, 255),
			},
			Files: []*discordgo.File{
				{
					Name:        "caption.jpeg",
					ContentType: "image/jpeg",
					Reader:      bytes.NewReader(export.Bytes()),
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
