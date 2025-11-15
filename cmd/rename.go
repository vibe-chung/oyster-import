/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

// renameCmd represents the rename command
var renameCmd = &cobra.Command{
	Use:   "rename",
	Short: "Rename Oyster CSV files based on their date range",
	Long: `Renames Oyster card CSV files to a standardized format based on the date range of journeys contained within.

The command examines each CSV file and, if valid, renames it to the format:
  YYYY-MM-DD_YYYY-MM-DD.csv (earliest date to latest date)

If a file with the same name already exists, an incremental suffix is added:
  YYYY-MM-DD_YYYY-MM-DD_1.csv
  YYYY-MM-DD_YYYY-MM-DD_2.csv

Example usage:
  oyster-import rename ~/Downloads/565348001.csv ~/Downloads/565348001\ \(1\).csv`,
	Args: cobra.MinimumNArgs(1),
	Run:  runRenameCmd,
}

func runRenameCmd(cmd *cobra.Command, args []string) {
	for _, csvPath := range args {
		if err := renameOysterFile(csvPath); err != nil {
			fmt.Printf("error processing '%s': %v\n", csvPath, err)
			continue
		}
	}
}

// renameOysterFile validates and renames a single Oyster CSV file
func renameOysterFile(csvPath string) error {
	// Read and validate the CSV file
	records, err := readCSV(csvPath)
	if err != nil {
		return fmt.Errorf("failed to read CSV: %w", err)
	}

	// Validate format and extract dates
	earliestDate, latestDate, err := extractDateRange(records)
	if err != nil {
		return fmt.Errorf("not a valid Oyster CSV file: %w", err)
	}

	// Generate new filename
	dir := filepath.Dir(csvPath)
	newName := generateFileName(dir, earliestDate, latestDate)
	newPath := filepath.Join(dir, newName)

	// Rename the file
	if err := os.Rename(csvPath, newPath); err != nil {
		return fmt.Errorf("failed to rename file: %w", err)
	}

	fmt.Printf("Renamed: %s -> %s\n", filepath.Base(csvPath), newName)
	return nil
}

// extractDateRange validates the CSV format and extracts the earliest and latest dates
func extractDateRange(records [][]string) (earliest, latest time.Time, err error) {
	if len(records) < 2 {
		return time.Time{}, time.Time{}, fmt.Errorf("file has insufficient rows")
	}

	// Check header format (basic validation)
	header := records[0]
	if len(header) < 8 {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid header format")
	}

	var dates []time.Time
	for i, row := range records {
		if i == 0 {
			continue // skip header
		}
		if len(row) < 8 {
			continue // skip incomplete rows
		}

		// Parse date in the format "02-Jan-2006" (e.g., "15-Jan-2025")
		parsedDate, err := time.Parse("02-Jan-2006", row[0])
		if err != nil {
			continue // skip rows with invalid dates
		}
		dates = append(dates, parsedDate)
	}

	if len(dates) == 0 {
		return time.Time{}, time.Time{}, fmt.Errorf("no valid journey dates found")
	}

	// Find earliest and latest
	earliest = dates[0]
	latest = dates[0]
	for _, d := range dates {
		if d.Before(earliest) {
			earliest = d
		}
		if d.After(latest) {
			latest = d
		}
	}

	return earliest, latest, nil
}

// generateFileName creates a unique filename in the format YYYY-MM-DD_YYYY-MM-DD.csv
// If a file with that name exists, appends an incremental digit
func generateFileName(dir string, earliest, latest time.Time) string {
	baseName := fmt.Sprintf("%s_%s.csv",
		earliest.Format("2006-01-02"),
		latest.Format("2006-01-02"))

	// Check if file exists
	fullPath := filepath.Join(dir, baseName)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return baseName
	}

	// File exists, try with incremental suffix
	for i := 1; ; i++ {
		nameWithSuffix := fmt.Sprintf("%s_%s_%d.csv",
			earliest.Format("2006-01-02"),
			latest.Format("2006-01-02"),
			i)
		fullPath = filepath.Join(dir, nameWithSuffix)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			return nameWithSuffix
		}
	}
}

func init() {
	// Command-specific flags can be added here if needed
}
