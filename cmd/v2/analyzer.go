package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/hiromaily/aurora-db-data-generator/pkg/logger"
	"github.com/hiromaily/aurora-db-data-generator/pkg/repository/dbrepo"
)

type analyzer struct {
	logger      logger.Logger
	schemaRepo  dbrepo.SchemaRepository
	currentDate string
}

func newAnalyzer(logger logger.Logger, schemaRepo dbrepo.SchemaRepository) *analyzer {
	// analyzer
	return &analyzer{
		logger:      logger,
		schemaRepo:  schemaRepo,
		currentDate: time.Now().Format("2006-01-02 15:04:05"),
	}
}

// Analyze type then create csv ​​in the defined column order
func (a *analyzer) analyze(tableName string, dataTagMap DataTagMap, count int) ([][]string, error) {
	// get columun info from target table
	columnInfo, err := a.schemaRepo.GetTableInfo(tableName)
	if err != nil {
		return nil, err
	}

	// create value for csv
	resultCSV := make([][]string, 0, count)
	for idx := range count {
		// data for 1 line
		record := make([]string, 0, len(columnInfo))
		for _, v := range columnInfo {
			if idx == 0 {
				// only first line to logging column name
				a.logger.Info(
					"columnInfo",
					"tableName", tableName,
					"columnName", v.ColumnName,
					"dataType", v.DataType,
					"isNullable", v.IsNullable,
					"columnDefault", v.ColumnDefault,
				)
			}
			// 1. check that column is defined in tagInfoMap or not
			if val, ok := dataTagMap[v.ColumnName]; ok {
				switch {
				case val.Prefix != "":
					// a. use prefix with index
					a.logger.Debug("set Prefix", "Prefix", val.Prefix)
					record = append(record, fmt.Sprintf("%s%06d", val.Prefix, idx+1))
				case val.Fixed != "":
					// b. use fixed value
					a.logger.Debug("set Fixed", "Fixed", val.Fixed)
					record = append(record, val.Fixed)
				case val.Iterator != nil:
					// c. use Iterator
					value, ok := val.Iterator.Next()
					if !ok {
						err := errors.New("generator list doesn't have enough length")
						a.logger.Error(
							err.Error(),
							"dataTagMap", dataTagMap,
							"tableName", tableName,
							"columnName", v.ColumnName,
							"error", err)
						return nil, err
					}
					record = append(record, value)
				default:
					// this must be bug for defined tagInfo in struct
					err := errors.New("invalid dataTagMap: no prefix and fixed value")
					a.logger.Error(
						err.Error(),
						"dataTagMap", dataTagMap,
						"tableName", tableName,
						"columnName", v.ColumnName,
						"error", err)
					return nil, err
				}
			} else {
				// 2. no definition in tagInfoMap, so set value by columnInfo
				switch {
				case v.ColumnName == "created_at" || v.ColumnName == "updated_at":
					// a. set date string for `created_at`, `updated_at`
					a.logger.Debug("set created_at/updated_at", "Fixed", a.currentDate)
					record = append(record, a.currentDate)
				case v.ColumnDefault.Valid:
					// b. use default value if possible
					a.logger.Debug("set ColumnDefault", "ColumnDefault", v.ColumnDefault.String)
					record = append(record, v.ColumnDefault.String)
				case v.IsNullable == "YES":
					// c. use NULL if allowed
					a.logger.Debug("set NULL", "IsNullable", v.IsNullable)
					record = append(record, "\\N")
				default:
					// this must be bug、`dataTagMap` must have value
					err := fmt.Errorf("column: %s is not configurable", v.ColumnName)
					a.logger.Error(
						err.Error(),
						"tableName", tableName,
						"columnName", v.ColumnName,
						"error", err)
					return nil, err
				}
			}
		}
		// validation
		if len(record) != len(columnInfo) {
			err := fmt.Errorf(
				"column length is %d, but stored column length: %d",
				len(columnInfo),
				len(record),
			)
			return nil, err
		}
		// add data for 1 line
		resultCSV = append(resultCSV, record)
	}
	return resultCSV, nil
}
