package img

import (
	"bufio"
	"bytes"
	_ "embed"
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
	"github.com/Fluffy-Bean/lynxie/internal/bot"
	"github.com/Fluffy-Bean/lynxie/internal/color"
	"github.com/bwmarrin/discordgo"
)

const maxFileSize = 1024 * 1024 * 10 // 10MB

var client = http.Client{
	Timeout: 10 * time.Second,
}

func RegisterImgCommands(h *bot.Handler) {
	_ = h.RegisterCommand("togif", cmdToGif(h))
	_ = h.RegisterCommandAlias("gif", "togif")
	_ = h.RegisterCommandAlias("saveable", "togif")

	_ = h.RegisterCommand("caption", cmdCaption(h))
	_ = h.RegisterCommandAlias("c", "caption")
}

func cmdToGif(h *bot.Handler) bot.Command {
	return func(h *bot.Handler, c bot.CommandContext) error {
		fileEndpoint, err := findClosestImage(c)
		if err != nil {
			return fmt.Errorf("get image: %s", err)
		}

		req, err := http.NewRequest(http.MethodGet, fileEndpoint, nil)
		if err != nil {
			return fmt.Errorf("create request: %s", err)
		}

		if req.ContentLength > maxFileSize {
			return fmt.Errorf("file size too large: %d", maxFileSize)
		}

		res, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("send request: %s", err)
		}
		defer res.Body.Close()

		_, err = c.Session.ChannelMessageSendComplex(c.Message.ChannelID, &discordgo.MessageSend{
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
			Reference: c.Message.Reference(),
		})
		if err != nil {
			return fmt.Errorf("send togif response: %s", err)
		}

		return nil
	}
}

func cmdCaption(h *bot.Handler) bot.Command {
	return func(h *bot.Handler, c bot.CommandContext) error {
		fileEndpoint, err := findClosestImage(c)
		if err != nil {
			return fmt.Errorf("get image: %s", err)
		}

		req, err := http.NewRequest(http.MethodGet, fileEndpoint, nil)
		if err != nil {
			return fmt.Errorf("create request: %s", err)
		}

		if req.ContentLength > maxFileSize {
			return fmt.Errorf("file size too large: %d", maxFileSize)
		}

		res, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("send request: %s", err)
		}
		defer res.Body.Close()

		buff, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("read response body: %s", err)
		}

		img, err := loadImageFromBytes(buff)
		if err != nil {
			return fmt.Errorf("load image from bytes: %s", err)
		}

		imgWidth, imgHeight := img.Bounds().Dx(), img.Bounds().Dy()

		captionSize := max(min(float64(imgWidth/15), 16), 50)
		captionText := strings.Join(h.ParseArgs(c), " ")

		// 8px padding all around
		_, captionHeight := measureText(_resources.FontRoboto, captionText, captionSize, imgWidth-16)
		captionHeight += 16

		if captionHeight < 128 {
			captionHeight = 128
		}

		canvas := gg.NewContext(imgWidth, imgHeight+captionHeight)
		err = canvas.LoadFontFaceFromBytes(_resources.FontRoboto, captionSize)
		if err != nil {
			return fmt.Errorf("load font: %s", err)
		}

		canvas.SetRGBA(1, 1, 1, 1)
		canvas.Clear()

		canvas.SetRGBA(0, 0, 0, 1)
		canvas.DrawStringWrapped(
			captionText,
			float64(imgWidth/2), float64(captionHeight/2),
			0.5, 0.5, float64(imgWidth),
			1.5,
			gg.AlignCenter,
		)

		canvas.DrawImage(img, 0, captionHeight)

		var export bytes.Buffer
		err = canvas.EncodeJPG(bufio.NewWriter(&export), &jpeg.Options{Quality: 100})
		if err != nil {
			return fmt.Errorf("encode image: %s", err)
		}

		_, err = c.Session.ChannelMessageSendComplex(c.Message.ChannelID, &discordgo.MessageSend{
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
			Reference: c.Message.Reference(),
		})
		if err != nil {
			return fmt.Errorf("send caption response: %s", err)
		}

		return nil
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

func findClosestImage(c bot.CommandContext) (string, error) {
	if len(c.Message.Attachments) >= 1 {
		if c.Message.Attachments[0].Size > maxFileSize {
			return "", fmt.Errorf("file size is too big")
		}

		return c.Message.Attachments[0].ProxyURL, nil
	}

	if c.Message.ReferencedMessage != nil {
		message := c.Message.ReferencedMessage

		if len(message.Attachments) >= 1 {
			if message.Attachments[0].Size > maxFileSize {
				return "", fmt.Errorf("file size is too big")
			}

			return message.Attachments[0].ProxyURL, nil
		}

		if len(message.Embeds) >= 1 && message.Embeds[0].Image != nil {
			return message.Embeds[0].Image.ProxyURL, nil
		}
	}

	history, err := c.Session.ChannelMessages(c.Message.ChannelID, 10, c.Message.ID, "", "")
	if err != nil {
		return "", err
	}
	for _, message := range history {
		if len(message.Attachments) >= 1 {
			if message.Attachments[0].Size > maxFileSize {
				return "", fmt.Errorf("file size is too big")
			}

			return message.Attachments[0].ProxyURL, nil
		}

		if len(message.Embeds) >= 1 && message.Embeds[0].Image != nil {
			return message.Embeds[0].Image.ProxyURL, nil
		}
	}

	return "", fmt.Errorf("no files exists")
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
