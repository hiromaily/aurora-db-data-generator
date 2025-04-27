package dbrepo

import "github.com/hiromaily/aurora-db-data-generator/pkg/db"

type SchemaRepository interface {
	Close() error
	GetDBName() string
	GetColumnCount(tableName string) (int, error)
	GetTableInfo(tableName string) ([]db.ColumnInfo, error)
}
