package handler

import (
	"database/sql"
	"fmt"
)

func (h *Handler) doDatabaseMigrations() error {
	migration := getMigrationVersion(h.db)

	var err error
	switch migration {
	case -1:
		fmt.Println("Doing migration 0000")
		_, err = h.db.Exec(`
			create table if not exists meta__migrations (
				version integer not null 
			);
	
			insert into meta__migrations (version) values (0);
		`)

	case 0:
		fmt.Println("Doing migration 0001")
		_, err = h.db.Exec(`
			create table if not exists jellyfin_registrations (
				id          integer not null primary key,
				discord_id  string not null,
				jellyfin_id string not null,
				created_at  string not null
			);
	
			update meta__migrations set version = 1;
		`)

	default:
		return nil
	}
	if err != nil {
		return fmt.Errorf("do migration: %w", err)
	}

	return h.doDatabaseMigrations()
}

func getMigrationVersion(db *sql.DB) int {
	row := db.QueryRow("SELECT version FROM meta__migrations LIMIT 1")

	var version int
	if err := row.Scan(&version); err != nil {
		return -1
	}

	return version
}
