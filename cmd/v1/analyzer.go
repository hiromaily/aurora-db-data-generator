package main

import (
	"errors"
	"fmt"
	"reflect"
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

// tag info
type TagInfo struct {
	Prefix string
	Fixed  string
}

type TagInfoMap map[string]TagInfo

// Analyze type then create csv ​​in the defined column order
func (a *analyzer) analyze(tableName string, dataStructure any, count int) ([][]string, error) {
	tagInfo := a.toTagInfoMap(tableName, dataStructure)

	// get columun info from target table
	columnInfo, err := a.schemaRepo.GetTableInfo(tableName)
	if err != nil {
		a.logger.Error("failed to call GetTableInfo()", "error", err)
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
			if val, ok := tagInfo[v.ColumnName]; ok {
				switch {
				case val.Prefix != "":
					// a. use prefix with index
					a.logger.Debug("set Prefix", "Prefix", val.Prefix)
					record = append(record, fmt.Sprintf("%s%06d", val.Prefix, idx+1))
				case val.Fixed != "":
					// b. use fixed value
					a.logger.Debug("set Fixed", "Fixed", val.Fixed)
					record = append(record, val.Fixed)
				default:
					// this must be bug for defined tagInfo in struct
					err := errors.New("invalid tagInfo")
					a.logger.Error(
						err.Error(),
						"tagInfo", tagInfo,
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
					// this must be bug、`dataStructure` must have value
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

func (a *analyzer) toTagInfoMap(tableName string, dataStructure any) TagInfoMap {
	// get type info for dataStructure
	t := reflect.TypeOf(dataStructure)
	val := reflect.ValueOf(dataStructure)

	a.logger.Info("dataAnalysis", "tableName", tableName, "type", t)

	// get tag info from struct
	tagInfoMap := make(TagInfoMap)
	for i := range val.NumField() {
		fieldType := t.Field(i)
		tagInfo := TagInfo{
			Prefix: fieldType.Tag.Get("prefix"),
			Fixed:  fieldType.Tag.Get("fixed"),
		}
		tagInfoMap[fieldType.Tag.Get("name")] = tagInfo

		a.logger.Debug(
			"fieldInfo",
			"fieldIndex", i,
			"fieldValue", val.Field(i).Interface(),
			"nameTag", fieldType.Tag.Get("name"),
			"prefixTag", tagInfo.Prefix,
			"fixedTag", tagInfo.Fixed,
		)
	}

	a.logger.Debug("tagInfoMap", "content", tagInfoMap)
	return tagInfoMap
}
