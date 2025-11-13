package cmd

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export journeys as JSON",
	Long:  `Export journey records from the local database as a JSON array to stdout. Supports optional filtering for commute journeys, defined as journeys on Tuesday, Wednesday, or Thursday with a start time between 7:00 and 9:59, and only if an end time is present, using the --commute-only flag. Use the --tail flag to return only the chronologically last record.`,
	Run: func(cmd *cobra.Command, args []string) {
		commuteOnly, _ := cmd.Flags().GetBool("commute-only")
		tail, _ := cmd.Flags().GetBool("tail")

		dbConn, err := sql.Open("sqlite3", "oyster.db")
		if err != nil {
			fmt.Printf("error opening database: %v\n", err)
			return
		}
		defer dbConn.Close()

		query := "SELECT id, date, start_time, end_time, journey_action, charge, credit, balance, note FROM journeys ORDER BY date DESC, start_time DESC"
		if tail {
			query += " LIMIT 1"
		}

		rows, err := dbConn.Query(query)
		if err != nil {
			fmt.Printf("error querying journeys: %v\n", err)
			return
		}
		defer rows.Close()

		type Journey struct {
			ID            int     `json:"id"`
			Date          string  `json:"date"`
			StartTime     string  `json:"start_time"`
			EndTime       string  `json:"end_time"`
			JourneyAction string  `json:"journey_action"`
			Charge        float64 `json:"charge"`
			Credit        float64 `json:"credit"`
			Balance       float64 `json:"balance"`
			Note          string  `json:"note"`
		}

		var journeys []Journey
		for rows.Next() {
			var j Journey
			err := rows.Scan(&j.ID, &j.Date, &j.StartTime, &j.EndTime, &j.JourneyAction, &j.Charge, &j.Credit, &j.Balance, &j.Note)
			if err != nil {
				fmt.Printf("error scanning row: %v\n", err)
				continue
			}
			if commuteOnly {
				if isCommuteJourney(j.Date, j.StartTime, j.EndTime) {
					journeys = append(journeys, j)
				}
			} else {
				journeys = append(journeys, j)
			}
		}

		data, err := json.MarshalIndent(journeys, "", "  ")
		if err != nil {
			fmt.Printf("error marshalling journeys: %v\n", err)
			return
		}
		fmt.Println(string(data))
	},
}

// isCommuteJourney returns true if the journey is on Tue/Wed/Thu between 8am and 10am
func isCommuteJourney(dateStr, timeStr, endTimeStr string) bool {
	if endTimeStr == "" {
		return false
	}

	t, err := parseDateTime(dateStr, timeStr)
	if err != nil {
		return false
	}

	weekday := t.Weekday()
	hour := t.Hour()
	return (weekday == time.Tuesday ||
		weekday == time.Wednesday ||
		weekday == time.Thursday) && (hour >= 7 && hour < 10)
}

// parseDateTime parses date and time strings into a time.Time
func parseDateTime(dateStr, timeStr string) (time.Time, error) {
	layout := "2006-01-02 15:04"
	return time.Parse(layout, dateStr+" "+timeStr)
}

func init() {
	rootCmd.AddCommand(exportCmd)
	exportCmd.Flags().BoolP("commute-only", "c", false, "Filter for commute journeys (Tue/Wed/Thu, 8-10am)")
	exportCmd.Flags().BoolP("tail", "t", false, "Return only the chronologically last record")
}
