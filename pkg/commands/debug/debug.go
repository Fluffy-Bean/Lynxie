package debug

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/Fluffy-Bean/lynxie/_resources"
	"github.com/Fluffy-Bean/lynxie/internal/color"
	"github.com/Fluffy-Bean/lynxie/internal/errors"
	"github.com/Fluffy-Bean/lynxie/internal/handler"
	"github.com/bwmarrin/discordgo"
)

func RegisterDebugCommands(bot *handler.Bot) {
	bot.RegisterCommand("debug", registerDebug(bot))
}

func registerDebug(bot *handler.Bot) handler.Callback {
	return func(h *handler.Handler, args []string) errors.Error {
		buildTags := "-"
		goVersion := strings.TrimPrefix(runtime.Version(), "go")
		gcCount := runtime.MemStats{}.NumGC
		buildHash := _resources.BuildHash
		buildPipeline := _resources.BuildPipelineLink
		latency := h.Session.HeartbeatLatency().Milliseconds()

		info, _ := debug.ReadBuildInfo()
		for _, setting := range info.Settings {
			switch setting.Key {
			case "-tags":
				buildTags = strings.ReplaceAll(setting.Value, ",", " ")
			}
		}

		_, err := h.Session.ChannelMessageSendComplex(h.Message.ChannelID, &discordgo.MessageSend{
			Embed: &discordgo.MessageEmbed{
				Title: "Lynxie",
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:   "Build Tags",
						Value:  buildTags,
						Inline: false,
					},
					{
						Name:   "Go version",
						Value:  goVersion,
						Inline: false,
					},
					{
						Name:   "OS/Arch",
						Value:  runtime.GOOS + "/" + runtime.GOARCH,
						Inline: false,
					},
					{
						Name:   "GC Count",
						Value:  fmt.Sprint(gcCount),
						Inline: false,
					},
					{
						Name:   "Build Hash",
						Value:  fmt.Sprintf("[%s](%s)", buildHash, buildPipeline),
						Inline: false,
					},
					{
						Name:   "Latency",
						Value:  fmt.Sprintf("%dms", latency),
						Inline: false,
					},
				},
				Color: color.RGBToDiscord(255, 255, 255),
			},
			Reference: h.Reference,
		})
		if err != nil {
			return errors.Error{
				Msg: "failed to send debug message",
				Err: err,
			}
		}

		return errors.Error{}
	}
}
