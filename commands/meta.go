package commands

import (
	"flag"
	"fmt"
	"runtime"
	"runtime/debug"
	"strings"
	"time"

	"github.com/Fluffy-Bean/lynxie/app"
	"github.com/Fluffy-Bean/lynxie/utils"
	"github.com/bwmarrin/discordgo"
)

func RegisterMetaCommands(a *app.App) {
	a.RegisterCommand("ping", registerPong(a))
	a.RegisterCommand("debug", registerDebug(a))
}

func registerPong(a *app.App) app.Callback {
	return func(h *app.Handler, args []string) app.Error {
		var options struct {
			latency bool
		}

		cmd := flag.NewFlagSet("pong", flag.ContinueOnError)
		cmd.BoolVar(&options.latency, "latency", false, "Display the latency of ping")
		cmd.Parse(args)

		var content string
		if options.latency {
			content = fmt.Sprintf("Pong! %dms", h.Session.HeartbeatLatency().Milliseconds())
		} else {
			content = "Pong!"
		}

		h.Session.ChannelMessageSendComplex(h.Message.ChannelID, &discordgo.MessageSend{
			Embed: &discordgo.MessageEmbed{
				Description: content,
				Color:       utils.ColorFromRGB(255, 255, 255),
			},
			Reference: h.Reference,
		})

		return app.Error{}
	}
}

func registerDebug(a *app.App) app.Callback {
	return func(h *app.Handler, args []string) app.Error {
		modified := false
		revision := "-"
		tags := "-"
		_go := strings.TrimPrefix(runtime.Version(), "go")
		gcCount := runtime.MemStats{}.NumGC
		localTime := time.Now().Local().Format("2006-01-02 15:04:05")

		info, _ := debug.ReadBuildInfo()
		for _, setting := range info.Settings {
			switch setting.Key {
			case "vcs.revision":
				revision = setting.Value
			case "vcs.modified":
				modified = setting.Value == "true"
			case "-tags":
				tags = strings.ReplaceAll(setting.Value, ",", " ")
			}
		}

		if modified {
			revision += " (uncommitted changes)"
		}

		h.Session.ChannelMessageSendComplex(h.Message.ChannelID, &discordgo.MessageSend{
			Embed: &discordgo.MessageEmbed{
				Description: strings.Join(
					[]string{
						"```",
						"Revision:     " + revision,
						"Build Tags:   " + tags,
						"Go version:   " + _go,
						"OS/Arch:      " + runtime.GOOS + "/" + runtime.GOARCH,
						"GC Count:     " + fmt.Sprint(gcCount),
						"Local Time:   " + localTime,
						"```",
					},
					"\n",
				),
				Color: utils.ColorFromRGB(255, 255, 255),
			},
			Reference: h.Reference,
		})

		return app.Error{}
	}
}
