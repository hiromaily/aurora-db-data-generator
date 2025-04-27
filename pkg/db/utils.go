package db

import (
	"fmt"
	"strings"
)

const (
	DBKindMySQL      = "mysql"
	DBKindPostgreSQL = "pgsql"
)

func getDatabaseName(dsn string) (string, error) {
	// Split at '/'
	parts := strings.Split(dsn, "/")
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid DSN: missing '/': %s", dsn)
	}

	// Split the last element with '?'
	dbNameParts := strings.SplitN(parts[len(parts)-1], "?", 2)
	if len(dbNameParts) == 0 {
		return "", fmt.Errorf("invalid DSN: missing database name: %s", dsn)
	}
	return dbNameParts[0], nil
}
