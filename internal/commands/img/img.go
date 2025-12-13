package img

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"strings"
	"time"

	"git.sr.ht/~sbinet/gg"
	"github.com/bwmarrin/discordgo"

	"github.com/Fluffy-Bean/lynxie/_resources"
	"github.com/Fluffy-Bean/lynxie/internal/handler"
	"github.com/Fluffy-Bean/lynxie/internal/utils"
)

const maxFileSize = 1024 * 1024 * 10 // 10MB

var client = http.Client{
	Timeout: 10 * time.Second,
}

func RegisterCommands(h *handler.Handler) {
	_ = h.RegisterCommand("togif", cmdToGif(h))
	_ = h.RegisterCommand("caption", cmdCaption(h))
}

func cmdToGif(h *handler.Handler) handler.Callback {
	_ = h.RegisterCommandHelp("togif", "Convert an image to a GIF to make it save-able within discord")
	_ = h.RegisterCommandAlias("togif", "gif")
	_ = h.RegisterCommandAlias("togif", "saveable")

	return func(session *discordgo.Session, message *discordgo.MessageCreate) error {
		fileEndpoint, err := findClosestImage(session, message)
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

		_, err = session.ChannelMessageSendComplex(message.ChannelID, &discordgo.MessageSend{
			Embed: &discordgo.MessageEmbed{
				Title:       "Saveable",
				Description: "Image converted to GIF :3",
				Image: &discordgo.MessageEmbedImage{
					URL: "attachment://saveable.gif",
				},
				Color: utils.RGBToDiscord(255, 255, 255),
			},
			Files: []*discordgo.File{
				{
					Name:        "saveable.gif",
					ContentType: "image/gif",
					Reader:      res.Body,
				},
			},
			Reference: message.Reference(),
		})
		if err != nil {
			return fmt.Errorf("send togif response: %s", err)
		}

		return nil
	}
}

func cmdCaption(h *handler.Handler) handler.Callback {
	_ = h.RegisterCommandHelp("caption", "Add a caption to a provided image")
	_ = h.RegisterCommandAlias("caption", "c")

	return func(session *discordgo.Session, message *discordgo.MessageCreate) error {
		fileEndpoint, err := findClosestImage(session, message)
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
		captionText := strings.Join(h.ParseArgs(session, message), " ")

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

		_, err = session.ChannelMessageSendComplex(message.ChannelID, &discordgo.MessageSend{
			Embed: &discordgo.MessageEmbed{
				Title: "Caption",
				Image: &discordgo.MessageEmbedImage{
					URL: "attachment://caption.jpeg",
				},
				Color: utils.RGBToDiscord(255, 255, 255),
			},
			Files: []*discordgo.File{
				{
					Name:        "caption.jpeg",
					ContentType: "image/jpeg",
					Reader:      bytes.NewReader(export.Bytes()),
				},
			},
			Reference: message.Reference(),
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

	case "image/jpeg":
		img, err = jpeg.Decode(bytes.NewReader(buff))
		if err != nil {
			return nil, fmt.Errorf("failed to decode jpeg: %s", err)
		}

	default:
		return nil, fmt.Errorf("unknown or unsupported format: %s", contentType)
	}

	return img, nil
}

func findClosestImage(session *discordgo.Session, message *discordgo.MessageCreate) (string, error) {
	if len(message.Attachments) >= 1 {
		if message.Attachments[0].Size > maxFileSize {
			return "", fmt.Errorf("file size is too big")
		}

		return message.Attachments[0].ProxyURL, nil
	}

	if message.ReferencedMessage != nil {
		referenced := message.ReferencedMessage

		if len(referenced.Attachments) >= 1 {
			if referenced.Attachments[0].Size > maxFileSize {
				return "", fmt.Errorf("file size is too big")
			}

			return referenced.Attachments[0].ProxyURL, nil
		}

		if len(referenced.Embeds) >= 1 && referenced.Embeds[0].Image != nil {
			return referenced.Embeds[0].Image.ProxyURL, nil
		}
	}

	history, err := session.ChannelMessages(message.ChannelID, 10, message.ID, "", "")
	if err != nil {
		return "", err
	}
	for _, entry := range history {
		if len(entry.Attachments) >= 1 {
			if entry.Attachments[0].Size > maxFileSize {
				return "", fmt.Errorf("file size is too big")
			}

			return entry.Attachments[0].ProxyURL, nil
		}

		if len(entry.Embeds) >= 1 && entry.Embeds[0].Image != nil {
			return entry.Embeds[0].Image.ProxyURL, nil
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
