# Copilot Instructions for oyster-import

## Project Overview

This is a command-line tool written in Go for importing Transport for London (TFL) Oyster card CSV files into a local SQLite database. The application provides import, export, and publish functionality for Oyster journey data.

## Repository Structure

```
oyster-import/
├── main.go              # Entry point - calls cmd.Execute()
├── cmd/                 # Cobra CLI command implementations
│   ├── root.go         # Root command and application setup
│   ├── import.go       # Import CSV files into database
│   ├── export.go       # Export journeys as JSON
│   └── publish.go      # Publish journeys to external services
├── db/                 # Database layer
│   └── database.go     # SQLite initialization and schema
├── repo/               # Repository/data access layer
│   └── journey.go      # Journey struct and insert logic
├── go.mod              # Go module definition
└── oyster.db           # SQLite database (created at runtime, gitignored)
```

## Build and Test Instructions

### Prerequisites
- Go 1.24.2 or later
- CGO enabled (required for go-sqlite3)

### Building
```bash
go build -v
```

This creates an `oyster-import` executable in the project root.

### Testing
```bash
go test ./... -v
```

Note: Currently, there are no test files in the project.

### Running
```bash
# Import one or more CSV files
./oyster-import import ~/Downloads/your-oyster.csv

# Export all journeys as JSON
./oyster-import export

# Export only commute journeys (Tue/Wed/Thu, 7-9:59am)
./oyster-import export --commute-only
```

## Code Style and Conventions

### Go Style
- Follow standard Go conventions and formatting (use `gofmt` or `go fmt`)
- Use meaningful variable and function names
- Error handling: always check and handle errors appropriately
- Comments: Add comments for exported functions and complex logic

### Package Organization
- `cmd/`: CLI commands using Cobra framework - each command in its own file
- `db/`: Database initialization and schema management
- `repo/`: Data models and database operations (repository pattern)

### Naming Conventions
- Exported functions use PascalCase
- Unexported functions use camelCase
- Struct fields are exported (PascalCase) for JSON marshaling
- Database column names use snake_case

## Database Schema

Table: `journeys`

| Column         | Type    | Constraints                     |
|----------------|---------|----------------------------------|
| id             | INTEGER | PRIMARY KEY AUTOINCREMENT       |
| date           | TEXT    | Part of unique constraint       |
| start_time     | TEXT    | Part of unique constraint       |
| end_time       | TEXT    | Part of unique constraint       |
| journey_action | TEXT    |                                 |
| charge         | REAL    |                                 |
| credit         | REAL    |                                 |
| balance        | REAL    |                                 |
| note           | TEXT    |                                 |
| processed      | INTEGER | DEFAULT 0                       |

**Unique Constraint**: (`date`, `start_time`, `end_time`) - prevents duplicate journeys

## Key Dependencies

- **github.com/spf13/cobra**: CLI framework for commands
- **github.com/mattn/go-sqlite3**: SQLite database driver (requires CGO)
- **cloud.google.com/go/pubsub/v2**: Google Cloud Pub/Sub for publishing

## Common Development Tasks

### Adding a New Command
1. Create a new file in `cmd/` (e.g., `cmd/newcommand.go`)
2. Define a new `cobra.Command` variable
3. Register it in `cmd/root.go`'s `init()` function using `rootCmd.AddCommand()`

### Modifying the Database Schema
1. Update the `CREATE TABLE` statement in `db/database.go`
2. Update the `Journey` struct in `repo/journey.go`
3. Update the `INSERT` statement in `repo/InsertJourney()`
4. Consider migration strategy for existing databases

### CSV Import Logic
- Header row is always skipped (index 0)
- Date format in CSV: `02-Jan-2006` (e.g., "15-Jan-2025")
- Date stored in database: `2006-01-02` (ISO format)
- Duplicate detection: uses `ON CONFLICT ... DO NOTHING` SQL clause
- Empty numeric fields default to 0.0

### Export Filtering
Commute journeys are defined as:
- Days: Tuesday, Wednesday, or Thursday
- Time: 7:00 AM to 9:59 AM (start time)
- Must have a non-empty end time

## Important Notes

### File Handling
- CSV files are specified as command-line arguments
- Multiple CSV files can be imported in a single command
- Database file `oyster.db` is created in the current working directory

### Error Handling
- CSV parsing errors are logged but don't stop processing
- Duplicate journeys return `sql.ErrNoRows` (silently skipped, not counted)
- Database connection errors are fatal and stop execution

### Testing
- No existing test infrastructure
- When adding tests, follow Go testing conventions
- Test files should be named `*_test.go`
- Consider using table-driven tests for CSV parsing

## Making Changes

### Minimal Changes Principle
- Make the smallest possible changes to address requirements
- Don't modify working code unless necessary
- Preserve existing functionality
- Maintain compatibility with existing CSV file format

### Before Committing
- Run `go build` to ensure code compiles
- Run `go fmt ./...` to format code
- Test manually with sample CSV files if modifying import/export logic
- Verify database schema changes don't break existing data

### Database Files
- `oyster.db` is gitignored - don't commit it
- `*.csv` files are gitignored - use sample data for testing

## Architecture Patterns

- **Repository Pattern**: Data access logic separated in `repo/` package
- **Cobra CLI Pattern**: Commands are modular and self-contained in `cmd/`
- **Constructor Pattern**: `db.InitDB()` handles database initialization
- **Error Propagation**: Errors bubble up to command handlers for user feedback

## Security Considerations

- SQL injection: Prevented by using parameterized queries (`?` placeholders)
- File path traversal: Command-line arguments treated as-is
- No authentication or authorization (local tool)
- Database file permissions follow OS defaults
