#!/usr/bin/env bash
# shell script for Mac-only

# set the PG_URL to the postgres instance for testing
PG_URL="postgres://localhost:5432?sslmode=disable"
GEESE="../../bin/goose-darwin64"

display() {
  echo "*** ${1} ***"
}

# ./automatic-up contains the SQL version higher than the DB version
display "Automatic operation chooses Goose Up"
${GEESE} --dir ./automatic-up postgres "${PG_URL}" automatic

# display the contents from the users table
display "Results from the Goose Up - entries in the users table"
psql --command 'SELECT * FROM users;'  "${PG_URL}"

# ./automatic-down contains the SQL version lower than the DB version
display "Automatic operation chooses Goose down"
${GEESE} --dir ./automatic-down postgres "${PG_URL}" automatic

# there are no entries in the users table
display "Results from the Goose Down - empty users table"
psql --command 'SELECT * FROM users;' "${PG_URL}"
