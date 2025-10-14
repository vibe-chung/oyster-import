/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/vibechung/oyster-import/db"
	"github.com/vibechung/oyster-import/repo"

	"github.com/spf13/cobra"
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import Oyster CSV data into a local SQLite database.",
	Long: `Reads an exported Oyster card CSV file and loads all journey and transaction data into a local SQLite database (oyster.db).

Example usage:
  oyster-import import ~/Downloads/565384001.csv

The database will contain a 'journeys' table with columns matching the CSV file header.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		csvPath := args[0]
		dbConn, err := db.InitDB("oyster.db")
		if err != nil {
			fmt.Printf("Error initializing database: %v\n", err)
			return
		}
		defer dbConn.Close()

		file, err := os.Open(csvPath)
		if err != nil {
			fmt.Printf("Error opening file: %v\n", err)
			return
		}
		defer file.Close()

		reader := csv.NewReader(file)
		records, err := reader.ReadAll()
		if err != nil {
			fmt.Printf("Error reading CSV: %v\n", err)
			return
		}

		inserted := 0
		for i, row := range records {
			if i == 0 {
				continue
			}
			if len(row) < 8 {
				fmt.Printf("Skipping incomplete row %d: %v\n", i, row)
				continue
			}
			journey := repo.Journey{
				Date:          row[0],
				StartTime:     row[1],
				EndTime:       row[2],
				JourneyAction: row[3],
				Charge:        parseFloat(row[4]),
				Credit:        parseFloat(row[5]),
				Balance:       parseFloat(row[6]),
				Note:          row[7],
			}
			err := repo.InsertJourney(dbConn, journey)
			if err != nil {
				if err == sql.ErrNoRows {
					// Row already exists, do not count as inserted
					continue
				}
				fmt.Printf("Error inserting row %d: %v\n", i, err)
			} else {
				inserted++
			}
		}
		fmt.Printf("CSV import complete. %d rows inserted.\n", inserted)
	},
}

// parseFloat converts a string to float64, returning 0 if empty or invalid
func parseFloat(s string) float64 {
	if s == "" {
		return 0
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return f
}

func init() {
	rootCmd.AddCommand(importCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// importCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// importCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
