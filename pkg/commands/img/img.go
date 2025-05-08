package img

import (
	"bufio"
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"strings"
	"time"

	"git.sr.ht/~sbinet/gg"
	"github.com/Fluffy-Bean/lynxie/_resources"
	"github.com/Fluffy-Bean/lynxie/app"
	"github.com/Fluffy-Bean/lynxie/internal/color"
	"github.com/bwmarrin/discordgo"
)

const maxFileSize = 1024 * 1024 * 10 // 10MB

var client = http.Client{
	Timeout: 10 * time.Second,
}

func RegisterImgCommands(a *app.App) {
	a.RegisterCommand("saveable", registerSaveable(a))
	a.RegisterCommandAlias("gif", "saveable")

	a.RegisterCommand("caption", registerCaption(a))
	a.RegisterCommandAlias("c", "caption")
}

func registerSaveable(a *app.App) app.Callback {
	return func(h *app.Handler, args []string) app.Error {
		fileEndpoint, err := findClosestImage(h)
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

		_, err = h.Session.ChannelMessageSendComplex(h.Message.ChannelID, &discordgo.MessageSend{
			Embed: &discordgo.MessageEmbed{
				Title:       "Saveable",
				Description: "Image converted to GIF :3",
				Image: &discordgo.MessageEmbedImage{
					URL: "attachment://saveable.gif",
				},
				Color: color.RGBToDiscord(1, 1, 1),
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
		if err != nil {
			return app.Error{
				Msg: "failed to send saveable message",
				Err: err,
			}
		}

		return app.Error{}
	}
}

func registerCaption(a *app.App) app.Callback {
	return func(h *app.Handler, args []string) app.Error {
		fileEndpoint, err := findClosestImage(h)
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
				Msg: "failed to read image",
				Err: err,
			}
		}

		img, err := loadImageFromBytes(buff)
		if err != nil {
			return app.Error{
				Msg: "failed to load image",
				Err: errors.New("Failed to load image " + err.Error()),
			}
		}
		imgWidth, imgHeight := img.Bounds().Dx(), img.Bounds().Dy()

		captionSize := float64(imgWidth / 15)
		if captionSize < 16 {
			captionSize = 16
		} else if captionSize > 50 {
			captionSize = 50
		}

		// 8px padding all around
		_, captionHeight := measureText(_resources.FontRoboto, strings.Join(args, " "), captionSize, imgWidth-16)
		captionHeight += 16

		if captionHeight < 128 {
			captionHeight = 128
		}

		canvas := gg.NewContext(imgWidth, imgHeight+captionHeight)
		err = canvas.LoadFontFaceFromBytes(_resources.FontRoboto, captionSize)
		if err != nil {
			return app.Error{
				Msg: "failed to load font",
				Err: err,
			}
		}

		canvas.SetRGBA(1, 1, 1, 1)
		canvas.Clear()

		canvas.SetRGBA(0, 0, 0, 1)
		canvas.DrawStringWrapped(
			strings.Join(args, " "),
			float64(imgWidth/2), float64(captionHeight/2),
			0.5, 0.5, float64(imgWidth),
			1.5,
			gg.AlignCenter,
		)

		canvas.DrawImage(img, 0, captionHeight)

		var export bytes.Buffer
		err = canvas.EncodeJPG(
			bufio.NewWriter(&export),
			&jpeg.Options{Quality: 100},
		)
		if err != nil {
			return app.Error{
				Msg: "failed to encode JPEG",
				Err: err,
			}
		}

		_, err = h.Session.ChannelMessageSendComplex(h.Message.ChannelID, &discordgo.MessageSend{
			Embed: &discordgo.MessageEmbed{
				Title: "Caption",
				Image: &discordgo.MessageEmbedImage{
					URL: "attachment://caption.jpeg",
				},
				Color: color.RGBToDiscord(1, 1, 1),
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
		if err != nil {
			return app.Error{
				Msg: "failed to send caption message",
				Err: err,
			}
		}

		return app.Error{}
	}
}

func loadImageFromBytes(buff []byte) (image.Image, error) {
	var (
		img image.Image
		err error
	)

	contentType := http.DetectContentType(buff)

	switch contentType {
	case "image/png":
		img, err = png.Decode(bytes.NewReader(buff))
		if err != nil {
			return nil, fmt.Errorf("failed to decode png: %s", err)
		}
		break
	case "image/jpeg":
		img, err = jpeg.Decode(bytes.NewReader(buff))
		if err != nil {
			return nil, fmt.Errorf("failed to decode jpeg: %s", err)
		}
		break
	default:
		return nil, fmt.Errorf("unknown or unsupported format: %s", contentType)
	}

	return img, nil
}

func findClosestImage(h *app.Handler) (string, error) {
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

func measureText(font []byte, text string, size float64, width int) (int, int) {
	canvas := gg.NewContext(width, width)
	err := canvas.LoadFontFaceFromBytes(font, size)
	if err != nil {
		return 0, 0
	}

	wrappedText := strings.Join(canvas.WordWrap(text, float64(width)), "\n")

	lineWidth, lineHeight := canvas.MeasureMultilineString(wrappedText, 1.5)

	return int(lineWidth), int(lineHeight)
}
