package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/ofgrenudo/415/apps/api/internal/api"
	"github.com/ofgrenudo/415/apps/api/internal/db"
	"github.com/ofgrenudo/415/apps/api/internal/store"
	"github.com/ofgrenudo/415/pkg/version"
)

const usageText = `NAME
       415 -- bill tracking API server

SYNOPSIS
       415 serve
       415 -h | --help
       415 -v | --version

DESCRIPTION
       415 is a bill tracking service. It exposes an HTTP CRUD API for
       recording bills (name, description, cost, and the date each
       one was paid) backed by a Postgres database, applying its
       schema migrations automatically on startup.

       Run with no arguments, 415 prints this man page and exits.

COMMANDS
       serve
              Connect to Postgres, apply any pending schema
              migrations, and start the HTTP API server.

OPTIONS
       -h, --help
              Print this help message and exit.

       -v, --version
              Print version information and exit.

API
       The serve command registers the following routes:

              GET    /bills       List all bills.
              POST   /bills       Create a bill.
              GET    /bills/:id   Fetch a single bill.
              PUT    /bills/:id   Update a bill.
              DELETE /bills/:id   Delete a bill.

ENVIRONMENT
       PORT
              Port the HTTP server listens on. Defaults to 8080.

       DB_HOST
              Postgres host. Defaults to localhost.

       DB_PORT
              Postgres port. Defaults to 5332.

       DB_USER
              Postgres user. Defaults to dba.

       DB_PASSWORD
              Postgres password. Defaults to dba.

       DB_NAME
              Postgres database name. Defaults to bank.

       DB_SSLMODE
              Postgres sslmode. Defaults to disable.

       The defaults match infra/db/docker-compose.yaml for local
       development.

EXIT STATUS
       0      Success.

       1      The database connection could not be established,
              migrations could not be applied, or the HTTP server
              failed to start or exited with an error.

EXAMPLES
       Start the API server:

              415 serve

       Start it against a non-default database:

              DB_HOST=db.example.com DB_NAME=bills 415 serve

SEE ALSO
       psql(1)

AUTHOR
       Joshua W.B.
`

func usage() {
	fmt.Fprint(os.Stdout, usageText)
}

func main() {
	var showVersion bool
	var showHelp bool
	flag.Usage = usage
	flag.BoolVar(&showVersion, "version", false, "Print version information and exit")
	flag.BoolVar(&showVersion, "v", false, "Shorthand for --version")
	flag.BoolVar(&showHelp, "help", false, "Print this help message and exit")
	flag.BoolVar(&showHelp, "h", false, "Shorthand for --help")
	flag.Parse()

	if showVersion {
		fmt.Println(version.String("415"))
		return
	}

	if showHelp {
		usage()
		return
	}

	args := flag.Args()
	if len(args) > 0 && args[0] == "serve" {
		if err := serve(); err != nil {
			log.Fatal(err)
		}
		return
	}

	usage()
}

func serve() error {
	conn, err := db.Open(db.ConfigFromEnv())
	if err != nil {
		return fmt.Errorf("connect to db: %w", err)
	}
	defer func() { _ = conn.Close() }()

	if err := db.Migrate(conn); err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}

	billStore := store.NewBillStore(conn)
	billsHandler := api.NewBillsHandler(billStore)

	router := gin.Default()
	billsHandler.Register(router)

	addr := ":" + getEnv("PORT", "8080")
	log.Printf("listening on %s", addr)
	return router.Run(addr)
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}
