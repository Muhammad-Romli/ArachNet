# ArachNet

A command-line RSS feed aggregator written in Go, backed by PostgreSQL.

Register users, add RSS feeds, follow and unfollow them, and fetch feed content from the terminal.

> **Status: work in progress.** Fetching works, but `agg` currently only fetches a feed and prints it. Storing posts and a `browse` command are not implemented yet. See [Roadmap](#roadmap).

## Requirements

- [Go](https://go.dev/dl/) 1.26 or newer
- [PostgreSQL](https://www.postgresql.org/download/) (developed against 18.6)

For development only:

- [goose](https://github.com/pressly/goose) for database migrations (developed with v3.28.0)
- [sqlc](https://sqlc.dev/) to generate Go code from SQL queries
- `psql` for inspecting the database

## Installation

Install directly with Go:

```bash
go install github.com/muhammad-romli/ArachNet@latest
```

The `ArachNet` binary is placed in `$GOPATH/bin` (usually `~/go/bin`). Make sure that directory is in your `PATH`.

Or build from source:

```bash
git clone https://github.com/muhammad-romli/ArachNet.git
cd ArachNet
go build -o ArachNet
```

## Setup

### 1. Create the database

```bash
psql -U postgres -c "CREATE DATABASE arachnet;"
```

Your connection string will look like this (adjust user, password, host and port):

```
postgres://USER:PASSWORD@localhost:5432/arachnet?sslmode=disable
```

### 2. Run the migrations

From the project root:

```bash
goose -dir sql/schema postgres "postgres://USER:PASSWORD@localhost:5432/arachnet?sslmode=disable" up
```

This creates the `users`, `feeds` and `feed_follows` tables.

### 3. Create the config file

ArachNet reads a config file that stores two values:

- the PostgreSQL connection string (`DBUrl`)
- the name of the currently logged-in user (set by `register` and `login`)

<!-- TODO: add the config file name, its location, and the exact JSON keys (see internal/config/config.go) -->

TODO: show an example of the file here, for example:

```json
{
  "db_url": "postgres://USER:PASSWORD@localhost:5432/arachnet?sslmode=disable",
  "current_user_name": ""
}
```

The example keys above are a guess, so replace them with what your code actually uses.

## Usage

```bash
ArachNet <command> [arguments]
```

| Command | Arguments | Login required | Description |
|---|---|---|---|
| `register` | `<name>` | No | Create a new user |
| `login` | `<name>` | No | Switch to an existing user |
| `users` | none | No | List all users |
| `reset` | none | No | Reset the database (deletes data) |
| `agg` | `<url>` | No | Fetch an RSS feed and print it |
| `addfeed` | `<name> <url>` | Yes | Add a feed |
| `feeds` | none | No | List all feeds |
| `follow` | `<name or url>` | Yes | Follow a feed |
| `following` | none | Yes | List the feeds the current user follows |
| `unfollow` | `<name or url>` | Yes | Unfollow a feed |

### Example

```bash
ArachNet register alice
ArachNet addfeed "Boot.dev Blog" https://blog.boot.dev/index.xml
ArachNet feeds
ArachNet following
ArachNet agg https://blog.boot.dev/index.xml
```

> Warning: `reset` wipes data. Don't run it against a database you care about.

## Development

Project layout:

```
.
├── main.go              # entry point, command registration
├── registry.go          # command registry and handlers
├── rss_functions.go     # RSS fetching and parsing
├── internal/
│   ├── config/          # config file loading
│   └── database/        # sqlc-generated code (do not edit by hand)
├── sql/
│   ├── schema/          # goose migrations
│   └── queries/         # sqlc query definitions
└── sqlc.yaml
```

After changing anything in `sql/queries/` or `sql/schema/`, regenerate the database code:

```bash
sqlc generate
```

Run the tests:

```bash
go test ./...
```

## Roadmap

- [ ] Store fetched posts in the database
- [ ] `browse` command to read stored posts
- [ ] Run `agg` on an interval

## Built with

- [lib/pq](https://github.com/lib/pq), PostgreSQL driver
- [google/uuid](https://github.com/google/uuid), UUID generation
- [sqlc](https://sqlc.dev/) and [goose](https://github.com/pressly/goose)

## License

This project is licensed under the **MIT License** - see the [LICENSE](LICENSE) file for details.
