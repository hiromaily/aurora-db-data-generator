package db

import (
	"database/sql"
	"log"

	"github.com/jmoiron/sqlx"
)

type PostgreClient struct {
	db     *sqlx.DB
	dbName string
}

func NewPostgreClient(dsn string) (*PostgreClient, error) {
	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		return nil, err
	}
	// データベース名を取得
	dbName, err := getDatabaseName(dsn)
	if err != nil {
		return nil, err
	}

	return &PostgreClient{
		db:     db,
		dbName: dbName,
	}, nil
}

func (m *PostgreClient) GetDBName() string {
	return m.dbName
}

func (m *PostgreClient) Close() error {
	return m.db.Close()
}

//
// Implementation for SchemaRepository interface
//

func (p *PostgreClient) GetColumnCount(tableName string) (int, error) {
	query := `
		SELECT COUNT(*) 
		FROM information_schema.columns
		WHERE table_schema = 'public' 
		AND table_name = $1;
		`
	var columnCount int
	err := p.db.Get(&columnCount, query, tableName)
	if err != nil {
		return 0, err
	}
	return columnCount, nil
}

type PgColumnInfo struct {
	ColumnName    string         `db:"column_name"`
	DataType      string         `db:"data_type"`
	IsNullable    string         `db:"is_nullable"`
	ColumnDefault sql.NullString `db:"column_default"`
}

func (p *PgColumnInfo) Convert() ColumnInfo {
	return ColumnInfo{
		ColumnName:    p.ColumnName,
		DataType:      p.DataType,
		IsNullable:    p.IsNullable,
		ColumnDefault: p.ColumnDefault,
	}
}

func (p *PostgreClient) GetTableInfo(tableName string) ([]ColumnInfo, error) {
	query := `
		SELECT column_name, data_type, is_nullable, column_default
		FROM information_schema.columns
		WHERE table_schema = 'public' 
		AND table_name = $1
		ORDER BY ordinal_position;
		`
	var columns []PgColumnInfo
	err := p.db.Select(&columns, query, tableName)
	if err != nil {
		return nil, err
	}
	log.Println(columns)
	// Convert
	converted := make([]ColumnInfo, 0, len(columns))
	for _, col := range columns {
		converted = append(converted, col.Convert())
	}
	return converted, nil
}
