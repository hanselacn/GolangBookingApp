package postgres

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/lib/pq"
	"github.com/pressly/goose/v3"
	_ "github.com/ziutek/mymysql/godrv"

	"GolangBookingApp/pkg/logger"
)

var (
	flags   = flag.NewFlagSet("db:migrate", flag.ExitOnError)
	dir     = flags.String("dir", "database/migration", "directory with migration files")
	table   = flags.String("table", "db_migration", "migrations table name")
	verbose = flags.Bool("verbose", false, "enable verbose mode")
	help    = flags.Bool("guide", false, "print help")
	version = flags.Bool("version", false, "print version")
)

func DatabaseMigration(cfg *Config) {

	flags.Usage = usage
	flags.Parse(os.Args[2:])

	if *version {
		fmt.Println(goose.VERSION)
		return
	}
	if *verbose {
		goose.SetVerbose(true)
	}

	goose.SetTableName(*table)

	args := flags.Args()

	if len(args) == 0 || *help {
		flags.Usage()
		return
	}

	switch args[0] {
	case "create":
		if err := goose.Run("create", nil, *dir, args[1:]...); err != nil {
			log.Fatalf("goose run: %v", err)
		}
		return
	case "fix":
		if err := goose.Run("fix", nil, *dir); err != nil {
			log.Fatalf("goose run: %v", err)
		}
		return
	}

	if len(args) < 1 {
		flags.Usage()
		return
	}

	command := args[0]

	dbstring := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)
	db, err := goose.OpenDBWithDriver("postgres", normalizeDBString("postgres", dbstring))

	if err != nil {
		logger.Fatal(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.Fatal(logger.MessageFormat("db migrate: failed to close DB: %v\n", err))
		}
	}()

	arguments := []string{}
	if len(args) > 3 {
		arguments = append(arguments, args[3:]...)
	}

	if err := goose.Run(command, db, *dir, arguments...); err != nil {
		log.Fatalf("db migrate run: %v", err)
	}
}

func usage() {
	fmt.Println(usageCommands)
}

var (
	usageCommands = `
  --dir string     directory with migration files (default "database/migration")
  --guide          print help
  --table string   migrations table name (default "db_migration")
  --verbose        enable verbose mode
  --version        print version

Commands:
    up                   Migrate the DB to the most recent version available
    up-by-one            Migrate the DB up by 1
    up-to VERSION        Migrate the DB to a specific VERSION
    down                 Roll back the version by 1
    down-to VERSION      Roll back to a specific VERSION
    redo                 Re-run the latest migration
    reset                Roll back all migrations
    status               Dump the migration status for the current DB
    version              Print the current version of the database
    create NAME [sql|go] Creates new migration file with the current timestamp
    fix                  Apply sequential ordering to migrations
`
)

func normalizeDBString(driver string, str string) string {
	if driver == "postgres" {
		var err error
		str, err = normalizePostgreSQLURL(str)
		if err != nil {
			log.Fatalf("failed to normalize PostgreSQL connection string: %v", err)
		}
	}
	return str
}

func normalizePostgreSQLURL(url string) (string, error) {
	config, err := pq.ParseURL(url)
	if err != nil {
		return "", err
	}
	return config, nil
}
