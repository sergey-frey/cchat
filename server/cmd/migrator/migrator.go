package migrator

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func NewMigration(storagePath string, migrationPath string) {
	m, err := migrate.New(
		"file://"+migrationPath,
		storagePath)
	if err != nil {
		panic(err)
	}
	// if err := m.Up(); err != nil {
	// 	if errors.Is(err, migrate.ErrNoChange) {
	// 		fmt.Println("no migrations to apply")

	// 		return
	// 	}

	// 	panic(err)

	// }

	if err := m.Migrate(1); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")

			return
		}

		panic(err)

	}
}
