package fortnite

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/Fluffy-Bean/lynxie/internal/bot"
	"github.com/Fluffy-Bean/lynxie/internal/color"
)

type festivalTrack struct {
	Title        string `json:"title"`
	Artist       string `json:"artist"`
	ReleaseYear  int    `json:"releaseYear"`
	Cover        string `json:"cover"`
	Bpm          int    `json:"bpm"`
	Duration     string `json:"duration"`
	Difficulties struct {
		Vocals        int `json:"vocals"`
		Guitar        int `json:"guitar"`
		Bass          int `json:"bass"`
		Drums         int `json:"drums"`
		PlasticBass   int `json:"plastic-bass"`
		PlasticDrums  int `json:"plastic-drums"`
		PlasticGuitar int `json:"plastic-guitar"`
	} `json:"difficulties"`
	CreatedAt time.Time `json:"createdAt"`
	// LastFeatured interface{} `json:"lastFeatured"`
	// PreviewUrl   interface{} `json:"previewUrl"`
	Featured bool `json:"featured"`
}

var client = http.Client{
	Timeout: 10 * time.Second,
}
var mutex = &sync.RWMutex{}
var festivalTracks = make(map[string]festivalTrack)

func RegisterFortniteCommands(h *bot.Handler) {
	_ = h.RegisterCommand("fn-festival", cmdFestivalTracks(h))
	_ = h.RegisterCommandAlias("fnff", "fn-festival")

	// The actual list only gets updated every 24 hours, but I have not yet figured out a way to do cron jobs on specific hours...
	h.ScheduleTask(updateTrackList, time.Hour*1)
}

func updateTrackList() {
	mutex.Lock()
	defer mutex.Unlock()

	fmt.Println("Updating festival track list...")

	req, err := http.NewRequest(http.MethodGet, "https://raw.githubusercontent.com/FNFestival/fnfestival.github.io/refs/heads/main/data/jam_tracks.json", nil)
	if err != nil {
		fmt.Println("update track list: create request: ", err)

		return
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println("update track list: send request: ", err)

		return
	}
	defer res.Body.Close()

	var updatedTrackList map[string]festivalTrack
	err = json.NewDecoder(res.Body).Decode(&updatedTrackList)
	if err != nil {
		fmt.Println("update track list: decode response: ", err)

		return
	}

	festivalTracks = updatedTrackList
}

func cmdFestivalTracks(h *bot.Handler) bot.Command {
	return func(h *bot.Handler, c bot.CommandContext) error {
		mutex.RLock()
		defer mutex.RUnlock()

		message := &discordgo.MessageSend{
			Embed: &discordgo.MessageEmbed{
				Title:  "Fortnite Festival",
				Fields: []*discordgo.MessageEmbedField{},
				Footer: &discordgo.MessageEmbedFooter{
					Text: "25 randomly picked currently featured tracks",
				},
				Color: color.RGBToDiscord(255, 255, 255),
			},
			Reference: c.Message.Reference(),
		}

		for _, track := range festivalTracks {
			if !track.Featured {
				continue
			}
			if len(message.Embed.Fields) >= 25 {
				break
			}

			message.Embed.Fields = append(message.Embed.Fields, &discordgo.MessageEmbedField{
				Name:   fmt.Sprintf("%s (%s)", track.Title, track.Artist),
				Value:  fmt.Sprintf("%s | %d | %dbpm", track.Duration, track.ReleaseYear, track.Bpm),
				Inline: false,
			})
		}

		_, err := c.Session.ChannelMessageSendComplex(c.Message.ChannelID, message)
		if err != nil {
			return fmt.Errorf("send festival response: %w", err)
		}

		return nil
	}
}
