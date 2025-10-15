/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	pubsub "cloud.google.com/go/pubsub/v2"

	"github.com/spf13/cobra"
)

// publishCmd represents the publish command
var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Publish journeys to Pub/Sub",
	Long:  `Publish all unprocessed journey records from the local database to a Google Cloud Pub/Sub topic. Each journey is sent as a JSON message. After publishing, records are marked as processed to avoid duplicates. Requires GCP_PROJECT_ID and GCP_PUBSUB_TOPIC environment variables.`,
	Run: func(cmd *cobra.Command, args []string) {
		dbConn, err := sql.Open("sqlite3", "oyster.db")
		if err != nil {
			fmt.Printf("error opening database: %v\n", err)
			return
		}
		defer dbConn.Close()

		projectID := os.Getenv("GCP_PROJECT_ID")
		topicName := os.Getenv("GCP_PUBSUB_TOPIC")
		if projectID == "" || topicName == "" {
			fmt.Println("GCP_PROJECT_ID and GCP_PUBSUB_TOPIC environment variables must be set")
			return
		}

		ctx := context.Background()
		client, err := pubsub.NewClient(ctx, projectID)
		if err != nil {
			fmt.Printf("error creating pubsub client: %v\n", err)
			return
		}
		defer client.Close()

		publisher := client.Publisher(topicName)

		rows, err := dbConn.Query("SELECT id, date, start_time, end_time, journey_action, charge, credit, balance, note FROM journeys WHERE processed = 0")
		if err != nil {
			fmt.Printf("error querying journeys: %v\n", err)
			return
		}
		defer rows.Close()

		type Journey struct {
			ID            int    `json:"id"`
			Date          string `json:"date"`
			StartTime     string `json:"start_time"`
			EndTime       string `json:"end_time"`
			JourneyAction string `json:"journey_action"`
			Charge        string `json:"charge"`
			Credit        string `json:"credit"`
			Balance       string `json:"balance"`
			Note          string `json:"note"`
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
		rows.Close()

		published := 0
		for _, j := range journeys {
			data, err := json.Marshal(j)
			if err != nil {
				fmt.Printf("error marshalling journey: %v\n", err)
				continue
			}
			res := publisher.Publish(ctx, &pubsub.Message{Data: data})
			_, err = res.Get(ctx)
			if err != nil {
				fmt.Printf("error publishing journey id %d: %v\n", j.ID, err)
				continue
			}
			_, err = dbConn.Exec("UPDATE journeys SET processed = 1 WHERE id = ?", j.ID)
			if err != nil {
				fmt.Printf("error updating journey id %d: %v\n", j.ID, err)
				continue
			}
			published++
		}
		publisher.Stop()
		fmt.Printf("Published %d journeys to Pub/Sub.\n", published)
	},
}

func init() {
	rootCmd.AddCommand(publishCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// publishCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// publishCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
