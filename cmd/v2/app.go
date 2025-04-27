package main

import (
	"errors"
	"fmt"

	"github.com/hiromaily/aurora-db-data-generator/pkg/csv"
	"github.com/hiromaily/aurora-db-data-generator/pkg/db"
	"github.com/hiromaily/aurora-db-data-generator/pkg/logger"
	"github.com/hiromaily/aurora-db-data-generator/pkg/repository/dbrepo"
)

type app struct {
	logger     logger.Logger
	schemaRepo dbrepo.SchemaRepository
	analyzer   *analyzer
}

func newApp(logger logger.Logger, dbKind, dsn string) (*app, error) {
	// DB client
	var schemaRepo dbrepo.SchemaRepository
	var err error

	switch dbKind {
	case db.DBKindPostgreSQL:
		// PostgreSQL
		schemaRepo, err = db.NewPostgreClient(dsn)
	case db.DBKindMySQL:
		// MySQL
		schemaRepo, err = db.NewMySQLClient(dsn)
	default:
		return nil, fmt.Errorf("unsupported DB kind: %s", dbKind)
	}
	if err != nil {
		logger.Error("failed to call NewDBClient()", "error", err)
		return nil, err
	}
	logger.Info("Database connected", "dbName", schemaRepo.GetDBName())

	// Analyzer
	analyzer := newAnalyzer(logger, schemaRepo)

	// App
	return &app{
		logger:     logger,
		schemaRepo: schemaRepo,
		analyzer:   analyzer,
	}, nil
}

func (a *app) close() error {
	return a.schemaRepo.Close()
}

// generate test data for App01
func (a *app) generateApp01(count int) error {
	appName := "app1"
	a.logger.Info("Generate test data", "app_name", appName)

	// user
	userDataTagMap := getUesrDataTagMap()
	iter := newCustomStringIterator("email%06d@example.com")
	userDataTagMap["email"] = DataTag{Iterator: iter}

	// validation
	for key, val := range userDataTagMap {
		if !val.isValid() {
			a.logger.Error(
				"userDataTagMap is invalid",
				"userDataTagMap", userDataTagMap,
				"key", key,
				"val", val,
			)
			return errors.New("smsAuthDataTagMap is invalid")
		}
	}

	tableDataTagMap := []TableDataTagMap{
		{
			tableName:  userTableName,
			dataTagMap: userDataTagMap,
		},
	}
	return a.callGenerateCSVs(appName, tableDataTagMap, count)
}

func (a *app) callGenerateCSVs(appName string, tableDataTagMaps []TableDataTagMap, count int) error {
	// TODO: concurrent execution
	for _, v := range tableDataTagMaps {
		if err := a.generateCSV(appName, v.tableName, v.dataTagMap, count); err != nil {
			return err
		}
	}
	return nil
}

// generate csv for customer table
func (a *app) generateCSV(appName, tableName string, tagInfoMap DataTagMap, count int) error {
	a.logger.Info("Analyze table and given data pattern", "tableName", tableName, "count", count)
	resultCSV, err := a.analyzer.analyze(tableName, tagInfoMap, count)
	if err != nil {
		a.logger.Error("failed to call analyze()", "tableName", tableName, "error", err)
		return err
	}
	if len(resultCSV) == 0 {
		err := errors.New("resultCSV is empty")
		a.logger.Error(err.Error(), "tableName", tableName, "error", err)
		return err
	}

	// generate csv file
	a.logger.Info("Generate csv file",
		"app_name", appName,
		"table_name", tableName,
		"result_count", len(resultCSV),
	)
	csvGen, err := csv.NewCSVGenerator(
		a.logger,
		fmt.Sprintf("testdata/%s/%s_%s.csv", appName, appName, tableName),
	)
	if err != nil {
		return err
	}
	defer csvGen.Close() //nolint: errcheck

	for _, v := range resultCSV {
		a.logger.Debug("result_csv_record", "resultCSV", v)
		if err := csvGen.Generate(v...); err != nil {
			a.logger.Warn("failed to call csvGen.generate()", "error", err)
		}
	}
	return nil
}
