/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
	"github.com/vibechung/oyster-import/db"
	"github.com/vibechung/oyster-import/repo"
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import Oyster CSV data into a local SQLite database.",
	Long: `Reads an exported Oyster card CSV file and loads all journey and transaction data into a local SQLite database (oyster.db).

Example usage:
	oyster-import import ~/Downloads/565384001.csv

The database will contain a 'journeys' table with columns matching the CSV file header.`,
	Args: cobra.MinimumNArgs(1),
	Run:  runImportCmd,
}

func runImportCmd(cmd *cobra.Command, args []string) {
	dbConn, err := db.InitDB("oyster.db")
	if err != nil {
		fmt.Printf("error initializing database: %v\n", err)
		return
	}
	defer dbConn.Close()

	totalInserted := 0
	for _, csvPath := range args {
		records, err := readCSV(csvPath)
		if err != nil {
			fmt.Printf("error reading csv '%s': %v\n", csvPath, err)
			continue
		}
		inserted := 0
		for i, row := range records {
			inserted += processJourneyRow(dbConn, row, i, csvPath)
		}
		fmt.Printf("File '%s' import complete. %d rows inserted.\n", csvPath, inserted)
		totalInserted += inserted
	}
	fmt.Printf("All files processed. Total rows inserted: %d\n", totalInserted)
}

// readCSV reads all records from a CSV file
func readCSV(path string) ([][]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading csv: %v", err)
	}
	return records, nil
}

// processJourneyRow parses and inserts a journey row, returning 1 if inserted, 0 otherwise
func processJourneyRow(dbConn *sql.DB, row []string, idx int, csvPath string) int {
	journey, err := parseJourneyRow(row, idx)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	err = repo.InsertJourney(dbConn, journey)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0
		}
		fmt.Printf("error inserting row %d in '%s': %v\n", idx, csvPath, err)
		return 0
	}
	return 1
}

// parseJourneyRow parses a CSV row into a Journey struct, with error handling
func parseJourneyRow(row []string, idx int) (repo.Journey, error) {
	if len(row) < 8 {
		return repo.Journey{}, fmt.Errorf("skipping incomplete row %d: %v", idx, row)
	}
	parsedDate, err := time.Parse("02-Jan-2006", row[0])
	if err != nil {
		return repo.Journey{}, fmt.Errorf("skipping row %d due to invalid date: %v", idx, err)
	}
	formattedDate := parsedDate.Format("2006-01-02")
	return repo.Journey{
		Date:          formattedDate,
		StartTime:     row[1],
		EndTime:       row[2],
		JourneyAction: row[3],
		Charge:        parseFloat(row[4]),
		Credit:        parseFloat(row[5]),
		Balance:       parseFloat(row[6]),
		Note:          row[7],
	}, nil
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
