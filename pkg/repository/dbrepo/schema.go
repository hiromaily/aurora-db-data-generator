package dbrepo

import (
	"github.com/hiromaily/aurora-db-data-generator/pkg/db"
	"github.com/hiromaily/aurora-db-data-generator/pkg/logger"
)

type schemaRepository struct {
	dbClient *db.MySQLClient
	dbName   string
	logger   logger.Logger
}

func NewSchemaRepository(dbClient *db.MySQLClient, logger logger.Logger) *schemaRepository {
	return &schemaRepository{
		dbClient: dbClient,
		dbName:   dbClient.GetDBName(),
		logger:   logger,
	}
}

func (s *schemaRepository) Close() error {
	return s.dbClient.Close()
}

func (s *schemaRepository) GetColumnCount(tableName string) (int, error) {
	return s.dbClient.GetColumnCount(tableName)
}

func (s *schemaRepository) GetTableInfo(tableName string) ([]db.ColumnInfo, error) {
	return s.dbClient.GetTableInfo(tableName)
}
