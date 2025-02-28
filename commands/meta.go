package commands

import (
	"flag"
	"fmt"
	"runtime"
	"runtime/debug"
	"strings"
	"time"

	"github.com/Fluffy-Bean/lynxie/app"
)

func RegisterMetaCommands(a *app.App) {
	a.RegisterCommand("ping", registerPong(a))
	a.RegisterCommand("debug", registerDebug(a))
}

func registerPong(a *app.App) func(h *app.Handler, args []string) {
	return func(h *app.Handler, args []string) {
		var options struct {
			latency bool
		}

		cmd := flag.NewFlagSet("pong", flag.ContinueOnError)
		cmd.BoolVar(&options.latency, "latency", false, "Display the latency of ping")
		cmd.Parse(args)

		if options.latency {
			h.Session.ChannelMessageSend(
				h.Message.ChannelID,
				fmt.Sprintf("Pong! %dms", h.Session.HeartbeatLatency().Milliseconds()),
			)
		} else {
			h.Session.ChannelMessageSend(
				h.Message.ChannelID,
				"Pong!",
			)
		}
	}
}

func registerDebug(a *app.App) func(h *app.Handler, args []string) {
	return func(h *app.Handler, args []string) {
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

		h.Session.ChannelMessageSend(
			h.Message.ChannelID,
			fmt.Sprintf(
				"```  Revision :: %s\nBuild Tags :: %s\nGo version :: %s\n   OS/Arch :: %s\n  GC Count :: %d\nLocal Time :: %s```",
				revision,
				tags,
				_go,
				runtime.GOOS+"/"+runtime.GOARCH,
				gcCount,
				localTime,
			),
		)
	}
}
