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

	_ "github.com/mattn/go-sqlite3"

	"github.com/spf13/cobra"
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		csvPath := args[0]
		db, err := sql.Open("sqlite3", "oyster.db")
		if err != nil {
			fmt.Printf("Error opening database: %v\n", err)
			return
		}
		defer db.Close()

		createTable := `CREATE TABLE IF NOT EXISTS journeys (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date TEXT,
			start_time TEXT,
			end_time TEXT,
			journey_action TEXT,
			charge REAL,
			credit REAL,
			balance REAL,
			note TEXT
		);`
		_, err = db.Exec(createTable)
		if err != nil {
			fmt.Printf("Error creating table: %v\n", err)
			return
		}

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

		// Skip header
		for i, row := range records {
			if i == 0 {
				continue
			}
			if len(row) < 8 {
				fmt.Printf("Skipping incomplete row %d: %v\n", i, row)
				continue
			}
			_, err := db.Exec(`INSERT INTO journeys (date, start_time, end_time, journey_action, charge, credit, balance, note) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
				row[0], row[1], row[2], row[3], parseFloat(row[4]), parseFloat(row[5]), parseFloat(row[6]), row[7])
			if err != nil {
				fmt.Printf("Error inserting row %d: %v\n", i, err)
			}
		}
		fmt.Println("CSV import complete.")
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
