# oyster-import
Import and store TFL Oyster csv files into a local sqlite database

## Features
- Imports Oyster card CSV files into a local SQLite database (`oyster.db`)
- Avoids duplicate journey records using a unique constraint
- Modular code structure (db and repo packages)
- Command-line interface using Cobra

## Usage
1. Export your Oyster card journey history as a CSV file from the TFL website.
2. Run the import command:

	```sh
	go run main.go import ~/Downloads/your-oyster.csv
	```

	Or, if built:
	```sh
	./oyster-import import ~/Downloads/your-oyster.csv
	```

3. The database file `oyster.db` will be created in the project root.
	- Only new rows are inserted; duplicates are ignored.
	- The number of rows inserted is displayed after import.

4. Export journeys as JSON:

	```sh
	./oyster-import export
	```

	To filter for commute journeys only:

	```sh
	./oyster-import export --commute-only
	```
	Commute journeys are defined as journeys on Tuesday, Wednesday, or Thursday with a start time between 7:00 and 9:59, and only if an end time is present.

## Project Structure
- `cmd/import.go` — CLI command logic
- `db/database.go` — Database initialization and schema
- `repo/journey.go` — Journey record insertion logic

## Database Schema
Table: `journeys`
- `id` INTEGER PRIMARY KEY AUTOINCREMENT
- `date` TEXT
- `start_time` TEXT
- `end_time` TEXT
- `journey_action` TEXT
- `charge` REAL
- `credit` REAL
- `balance` REAL
- `note` TEXT
- Unique constraint: (`date`, `start_time`, `end_time`)
