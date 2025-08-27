package debug

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/Fluffy-Bean/lynxie/_resources"
	"github.com/Fluffy-Bean/lynxie/internal/bot"
	"github.com/Fluffy-Bean/lynxie/internal/color"
	"github.com/bwmarrin/discordgo"
)

func RegisterDebugCommands(h *bot.Handler) {
	_ = h.RegisterCommand("debug", cmdDebug(h))
	_ = h.RegisterCommand("panic", cmdPanic(h))
}

func cmdDebug(h *bot.Handler) bot.Command {
	return func(h *bot.Handler, c bot.CommandContext) error {
		buildTags := "-"
		goVersion := strings.TrimPrefix(runtime.Version(), "go")
		gcCount := runtime.MemStats{}.NumGC
		buildHash := _resources.BuildHash
		buildPipeline := _resources.BuildPipelineLink
		latency := c.Session.HeartbeatLatency().Milliseconds()

		info, _ := debug.ReadBuildInfo()
		for _, setting := range info.Settings {
			switch setting.Key {
			case "-tags":
				buildTags = strings.ReplaceAll(setting.Value, ",", " ")
			}
		}

		_, err := c.Session.ChannelMessageSendComplex(c.Message.ChannelID, &discordgo.MessageSend{
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
			Reference: c.Message.Reference(),
		})
		if err != nil {
			return fmt.Errorf("send debug response: %s", err)
		}

		return nil
	}
}

func cmdPanic(h *bot.Handler) bot.Command {
	return func(h *bot.Handler, c bot.CommandContext) error {
		panic("we all panic!")
	}
}
