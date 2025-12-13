package debug

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/Fluffy-Bean/lynxie/_resources"
	"github.com/Fluffy-Bean/lynxie/internal/handler"
	"github.com/Fluffy-Bean/lynxie/internal/utils"
)

func RegisterCommands(h *handler.Handler) {
	_ = h.RegisterCommand("debug", cmdDebug(h))
	_ = h.RegisterCommand("panic", cmdPanic(h))
}

func cmdDebug(h *handler.Handler) handler.Callback {
	_ = h.RegisterCommandHelp("debug", "Show debug stats of the bot")

	return func(session *discordgo.Session, message *discordgo.MessageCreate) error {
		buildTags := "-"
		goVersion := strings.TrimPrefix(runtime.Version(), "go")
		gcCount := runtime.MemStats{}.NumGC
		buildHash := _resources.BuildHash
		buildPipeline := _resources.BuildPipelineLink
		latency := session.HeartbeatLatency().Milliseconds()

		info, _ := debug.ReadBuildInfo()
		for _, setting := range info.Settings {
			switch setting.Key {
			case "-tags":
				buildTags = strings.ReplaceAll(setting.Value, ",", " ")
			}
		}

		_, err := session.ChannelMessageSendComplex(message.ChannelID, &discordgo.MessageSend{
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
				Color: utils.RGBToDiscord(255, 255, 255),
			},
			Reference: message.Reference(),
		})
		if err != nil {
			return fmt.Errorf("send debug response: %s", err)
		}

		return nil
	}
}

func cmdPanic(h *handler.Handler) handler.Callback {
	_ = h.RegisterCommandHelp("panic", "Make the bot panic!")

	return func(session *discordgo.Session, message *discordgo.MessageCreate) error {
		panic("we all panic!")
	}
}
