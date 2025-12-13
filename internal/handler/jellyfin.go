package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Fluffy-Bean/lynxie/internal/utils"
)

type JellyfinRegistration struct {
	ID         int
	DiscordID  string
	JellyfinID string
	CreatedAt  time.Time
}

func (h *Handler) FindJellyfinRegistrationByDiscordID(discordID string) (*JellyfinRegistration, error) {
	row := h.db.QueryRow(
		`select id, discord_id, jellyfin_id, created_at from jellyfin_registrations where discord_id = :id`,
		sql.Named("id", discordID),
	)

	var registration JellyfinRegistration
	var createdAt string

	err := row.Scan(
		&registration.ID,
		&registration.DiscordID,
		&registration.JellyfinID,
		&createdAt,
	)
	if err != nil {
		return nil, fmt.Errorf("scan row: %w", err)
	}

	registration.CreatedAt, err = time.Parse(time.RFC3339, createdAt)
	if err != nil {
		return nil, fmt.Errorf("parse createdAt: %w", err)
	}

	return &registration, nil
}

func (h *Handler) FindJellyfinRegistrationByJellyfinID(jellyfinID string) (*JellyfinRegistration, error) {
	row := h.db.QueryRow(
		`select id, discord_id, jellyfin_id, created_at from jellyfin_registrations where jellyfin_id = :id`,
		sql.Named("id", jellyfinID),
	)

	var registration JellyfinRegistration
	var createdAt string

	err := row.Scan(
		&registration.ID,
		&registration.DiscordID,
		&registration.JellyfinID,
		&createdAt,
	)
	if err != nil {
		return nil, fmt.Errorf("scan row: %w", err)
	}

	registration.CreatedAt, err = time.Parse(time.RFC3339, createdAt)
	if err != nil {
		return nil, fmt.Errorf("parse createdAt: %w", err)
	}

	return &registration, nil
}

func (h *Handler) CreateJellyfinRegistration(discordID, discordUsername string) (*JellyfinRegistration, string, string, error) {
	jellyfinUsername := discordUsername
	jellyfinPassword := utils.NewPassword()

	jellyfinUserID, err := h.createJellyfinUser(jellyfinUsername, jellyfinPassword)
	if err != nil {
		return nil, "", "", fmt.Errorf("create jellyfin user: %w", err)
	}

	var registration JellyfinRegistration

	registration.DiscordID = discordID
	registration.JellyfinID = jellyfinUserID
	registration.CreatedAt = time.Now()

	row := h.db.QueryRow(
		`insert into jellyfin_registrations (discord_id, jellyfin_id, created_at) values (:discord_id, :jellyfin_id, :created_at) returning id`,
		sql.Named("discord_id", registration.DiscordID),
		sql.Named("jellyfin_id", registration.JellyfinID),
		sql.Named("created_at", registration.CreatedAt.Format(time.RFC3339)),
	)

	if err := row.Scan(&registration.ID); err != nil {
		return nil, "", "", fmt.Errorf("scan row: %w", err)
	}

	if row.Err() != nil {
		return nil, "", "", fmt.Errorf("insert row: %w", row.Err())
	}

	return &registration, jellyfinUsername, jellyfinPassword, nil
}

type jellyfinNewUserRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type jellyfinNewUserResponse struct {
	ID string `json:"id"`
}

func (h *Handler) createJellyfinUser(username, password string) (string, error) {
	if username == "" || password == "" {
		return "", fmt.Errorf("username or password is empty")
	}

	jellyfinEndpoint := h.GetValue("JellyfinEndpoint")
	jellyfinToken := h.GetValue("JellyfinToken")

	if jellyfinEndpoint == "" || jellyfinToken == "" {
		return "", fmt.Errorf("endpoint or token is empty")
	}

	data, err := json.Marshal(jellyfinNewUserRequest{
		Name:     username,
		Password: password,
	})
	if err != nil {
		return "", fmt.Errorf("marshal data: %w", err)
	}

	url := fmt.Sprintf("%s/Users/New", jellyfinEndpoint)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("MediaBrowser Token=\"%s\"", jellyfinToken))

	res, err := h.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("do request: %w", err)
	}

	switch res.StatusCode {
	case http.StatusUnauthorized:
		return "", fmt.Errorf("unauthorized")

	case http.StatusForbidden:
		return "", fmt.Errorf("forbidden")

	case http.StatusServiceUnavailable:
		return "", fmt.Errorf("service unavailable")
	}

	var newUser jellyfinNewUserResponse
	if err := json.NewDecoder(res.Body).Decode(&newUser); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}

	return newUser.ID, nil
}
