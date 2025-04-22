package db

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type DBClient struct {
	db     *sqlx.DB
	dbName string
}

func NewDBClient(dsn string) (*DBClient, error) {
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return nil, err
	}
	// データベース名を取得
	dbName, err := getDatabaseName(dsn)
	if err != nil {
		return nil, err
	}

	return &DBClient{
		db:     db,
		dbName: dbName,
	}, nil
}

func getDatabaseName(dsn string) (string, error) {
	// '/'で分割
	parts := strings.Split(dsn, "/")
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid DSN: missing '/': %s", dsn)
	}

	// 最後の要素を'?'で分割
	dbNameParts := strings.SplitN(parts[len(parts)-1], "?", 2)
	if len(dbNameParts) == 0 {
		return "", fmt.Errorf("invalid DSN: missing database name: %s", dsn)
	}
	return dbNameParts[0], nil
}

func (d *DBClient) GetDBName() string {
	return d.dbName
}

func (d *DBClient) Close() error {
	return d.db.Close()
}

//
// Implementation for SchemaRepository interface
//

func (d *DBClient) GetColumnCount(tableName string) (int, error) {
	query := `
		SELECT COUNT(*) 
		FROM information_schema.columns 
		WHERE table_schema = ? 
		AND table_name = ?;	
		`
	var columnCount int
	err := d.db.Get(&columnCount, query, d.dbName, tableName)
	if err != nil {
		return 0, err
	}
	return columnCount, nil
}

// なぜMappingに`columns.COLUMN_NAME`のようにcolumnsを指定する必要があるのか？
type ColumnInfo struct {
	ColumnName    string         `db:"columns.COLUMN_NAME"`
	DataType      string         `db:"columns.DATA_TYPE"`
	IsNullable    string         `db:"columns.IS_NULLABLE"`
	ColumnDefault sql.NullString `db:"columns.COLUMN_DEFAULT"`
}

func (d *DBClient) GetTableInfo(tableName string) ([]ColumnInfo, error) {
	query := `
		SELECT COLUMN_NAME, DATA_TYPE, IS_NULLABLE, COLUMN_DEFAULT
		FROM information_schema.columns 
		WHERE table_schema = ? 
		AND table_name = ? 
		ORDER BY ORDINAL_POSITION;
		`
	var columns []ColumnInfo
	err := d.db.Select(&columns, query, d.dbName, tableName)
	if err != nil {
		return nil, err
	}
	return columns, nil
}
