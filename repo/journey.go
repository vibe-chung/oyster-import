package repo

import (
	"database/sql"
)

type Journey struct {
	Date          string
	StartTime     string
	EndTime       string
	JourneyAction string
	Charge        float64
	Credit        float64
	Balance       float64
	Note          string
}

func InsertJourney(db *sql.DB, j Journey) error {
	_, err := db.Exec(`INSERT INTO journeys (date, start_time, end_time, journey_action, charge, credit, balance, note) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		j.Date, j.StartTime, j.EndTime, j.JourneyAction, j.Charge, j.Credit, j.Balance, j.Note)
	return err
}
