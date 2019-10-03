package goose

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"path/filepath"

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
		if err = dumpSQLFiles(db, dir, topVersion+1, currentVersion); err != nil {
			return errors.Wrapf(err, "dumpSQLFiles: downloading migrations from %d->%d", topVersion+1, currentVersion)
		}

		if err = DownTo(db, dir, topVersion); err != nil {
			return errors.Wrapf(err, "DownTo: down versioning the DB to %d", topVersion)
		}
	} else {
		fmt.Printf("goose: no migrations to run. current version: %d, top most migration: %d\n", currentVersion, topVersion)
	}

	return nil
}

func dumpSQLFiles(db *sql.DB, dir string, fromVersion int64, toVersion int64) error{
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
		fileName := filepath.Join(dir, fmt.Sprintf("%d_down_version.sql", row.VersionID))
		if err = ioutil.WriteFile(fileName, []byte(row.DownData), 0666); err != nil {
			return errors.Wrapf(err, "WriteFile: writing goose down information to %s", fileName)
		}
	}

	return nil
}