package db

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type MySQLClient struct {
	db     *sqlx.DB
	dbName string
}

func NewMySQLClient(dsn string) (*MySQLClient, error) {
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return nil, err
	}
	// データベース名を取得
	dbName, err := getDatabaseName(dsn)
	if err != nil {
		return nil, err
	}

	return &MySQLClient{
		db:     db,
		dbName: dbName,
	}, nil
}

func (m *MySQLClient) GetDBName() string {
	return m.dbName
}

func (m *MySQLClient) Close() error {
	return m.db.Close()
}

//
// Implementation for SchemaRepository interface
//

func (m *MySQLClient) GetColumnCount(tableName string) (int, error) {
	query := `
		SELECT COUNT(*) 
		FROM information_schema.columns 
		WHERE table_schema = ? 
		AND table_name = ?;	
		`
	var columnCount int
	err := m.db.Get(&columnCount, query, m.dbName, tableName)
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

func (m *MySQLClient) GetTableInfo(tableName string) ([]ColumnInfo, error) {
	query := `
		SELECT COLUMN_NAME, DATA_TYPE, IS_NULLABLE, COLUMN_DEFAULT
		FROM information_schema.columns 
		WHERE table_schema = ? 
		AND table_name = ? 
		ORDER BY ORDINAL_POSITION;
		`
	var columns []ColumnInfo
	err := m.db.Select(&columns, query, m.dbName, tableName)
	if err != nil {
		return nil, err
	}
	return columns, nil
}
