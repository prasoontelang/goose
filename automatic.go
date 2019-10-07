package goose

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

func Automatic(db *sql.DB, dir string) error {
	currentVersion, err := GetDBVersion(db)
	if err != nil {
		return errors.Wrapf(err, "GetDBVersion: getting the applied version in the DB")
	}

	migrations, err := CollectMigrations(dir, minVersion, maxVersion)
	if err != nil {
		return errors.Wrapf(err, "CollectMigrations: retrieving the migration files in the directory")
	}

	migrationsLen := len(migrations)
	topVersion := int64(0)
	if migrationsLen != 0 {
		topVersion = migrations[migrationsLen - 1].Version
	}
	if currentVersion < topVersion {
		if err := Up(db, dir); err != nil {
			return errors.Wrapf(err, "Up: migrating from %d -> %d", currentVersion, topVersion)
		}
	} else if currentVersion > topVersion {
		// dump the SQL files stored in the DB
		if err = downFromDb(db, topVersion+1); err != nil {
			return errors.Wrapf(err, "downFromDb: downloading migrations from %d->%d", topVersion+1, currentVersion)
		}

		if err = DownTo(db, dir, topVersion); err != nil {
			return errors.Wrapf(err, "DownTo: down versioning the DB to %d", topVersion)
		}
	} else {
		fmt.Printf("goose: no migrations to run. current version: %d, top most migration: %d\n", currentVersion, topVersion)
	}

	return nil
}

func downFromDb(db *sql.DB, fromVersion int64) error{
	rows, err := GetDialect().dbVersionQuery(db)
	if err != nil {
		return errors.Wrapf(err, "dbVersionQuery: getting all applied DB versions")
	}
	defer rows.Close()

	for rows.Next() {
		var row MigrationRecord
		if err = rows.Scan(&row.VersionID, &row.IsApplied, &row.DownData); err != nil {
			return errors.Wrapf(err, "Scan: scanning the migration record")
		}
		if row.VersionID < fromVersion || row.VersionID == 0 {
			break
		}

		log.Println("Performing goose down: ", row.DownData)
		r := strings.NewReader(row.DownData)
		statements, useTx, err := parseSQLMigration(r, false)
		if err != nil {
			return errors.Wrapf(err, "ERROR %v: failed to parse SQL migration file from DB")
		}

		if err := runSQLMigration(db, statements, useTx, statements, row.VersionID, false); err != nil {
			return errors.Wrapf(err, "ERROR %v: failed to run down SQL migration from DB")
		}
	}

	return nil
}