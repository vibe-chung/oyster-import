/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "oyster-import",
	Short: "Oyster card CSV importer CLI",
	Long:  "A command-line tool for importing one or more Oyster card CSV files into a local SQLite database (oyster.db). Supports multi-file import, automatic header skipping, and duplicate journey detection. Each journey and transaction is stored in the 'journeys' table for easy querying and analysis.",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.oyster-import.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// Register subcommands
	rootCmd.AddCommand(importCmd)
	rootCmd.AddCommand(publishCmd)
	rootCmd.AddCommand(exportCmd)
}
