package debug

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/Fluffy-Bean/lynxie/app"
	"github.com/Fluffy-Bean/lynxie/utils"
	"github.com/bwmarrin/discordgo"
)

func RegisterDebugCommands(a *app.App) {
	a.RegisterCommand("debug", registerDebug(a))
}

func registerDebug(a *app.App) app.Callback {
	return func(h *app.Handler, args []string) app.Error {
		buildTags := "-"
		goVersion := strings.TrimPrefix(runtime.Version(), "go")
		gcCount := runtime.MemStats{}.NumGC
		buildHash, _ := a.Config.CommandExtras["debug_build-hash"]
		buildPipeline, _ := a.Config.CommandExtras["debug_build-pipeline"]
		latency := h.Session.HeartbeatLatency().Milliseconds()

		info, _ := debug.ReadBuildInfo()
		for _, setting := range info.Settings {
			switch setting.Key {
			case "-tags":
				buildTags = strings.ReplaceAll(setting.Value, ",", " ")
			}
		}

		h.Session.ChannelMessageSendComplex(h.Message.ChannelID, &discordgo.MessageSend{
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
				Color: utils.ColorFromRGB(255, 255, 255),
			},
			Reference: h.Reference,
		})

		return app.Error{}
	}
}
