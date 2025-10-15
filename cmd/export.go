package cmd

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export journeys from the database",
	Long:  `Export journey records from the local database to a CSV file or stdout.`,
	Run: func(cmd *cobra.Command, args []string) {
		dbConn, err := sql.Open("sqlite3", "oyster.db")
		if err != nil {
			fmt.Printf("error opening database: %v\n", err)
			return
		}
		defer dbConn.Close()

		rows, err := dbConn.Query("SELECT id, date, start_time, end_time, journey_action, charge, credit, balance, note FROM journeys ORDER BY date DESC")
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
			journeys = append(journeys, j)
		}

		data, err := json.MarshalIndent(journeys, "", "  ")
		if err != nil {
			fmt.Printf("error marshalling journeys: %v\n", err)
			return
		}
		fmt.Println(string(data))
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)
}
