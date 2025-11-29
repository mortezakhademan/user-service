package models

import (
	"go.mongodb.org/mongo-driver/v2/bson"

	"strconv"
	"strings"
	"time"

	"git.ramooz.org/ramooz/golang-components/paginated-list/mongodb/dataType"
)

const (
	COLUMN_INFO_DATA_TYPE_TEXT      = "text"
	COLUMN_INFO_DATA_TYPE_DATE      = "date"
	COLUMN_INFO_DATA_TYPE_INT       = "int"
	COLUMN_INFO_DATA_TYPE_BOOL      = "bool"
	COLUMN_INFO_DATA_TYPE_FLOAT64   = "float64"
	COLUMN_INFO_DATA_TYPE_OBJECT_ID = "objectID"
	COLUMN_INFO_DATA_TYPE_ENUM      = "enum"

	COLUMN_INFO_FILTER_OPERATOR_NONE = "-"
)

type ColumnInfo struct {
	///
	Column           string
	FilterOperator   string
	DataType         string
	PrepareValueFunc func(interface{}) interface{}
}

func NewColumnInfo(column string, filterOperator string, dataType string) *ColumnInfo {
	return newColumnInfoWithExportSettings(column, filterOperator, dataType)
}
func newColumnInfoWithExportSettings(column string, filterOperator string, dataType string) *ColumnInfo {
	return &ColumnInfo{
		Column:         column,
		FilterOperator: filterOperator,
		DataType:       dataType,
	}
}
func NewEnumColumnInfo(column string, mapStringToInt map[string]int32) *ColumnInfo {
	return &ColumnInfo{
		Column:         column,
		FilterOperator: "in",
		DataType:       COLUMN_INFO_DATA_TYPE_ENUM,
		PrepareValueFunc: func(value interface{}) interface{} {
			if value == nil {
				return nil
			}
			return mapStringToInt[value.(string)]
		},
	}
}
func NewTextColumnInfo(column string) *ColumnInfo {
	return NewColumnInfo(column, "like", COLUMN_INFO_DATA_TYPE_TEXT)
}

func NewIntColumnInfo(column string) *ColumnInfo {
	return NewColumnInfo(column, "=", COLUMN_INFO_DATA_TYPE_INT)
}

func NewBetweenIntColumnInfo(column string) *ColumnInfo {
	return NewColumnInfo(column, "between", COLUMN_INFO_DATA_TYPE_INT)
}

func NewBoolColumnInfo(column string) *ColumnInfo {
	return NewColumnInfo(column, "=", COLUMN_INFO_DATA_TYPE_BOOL)
}

func NewFloat64ColumnInfo(column string) *ColumnInfo {
	return NewColumnInfo(column, "=", COLUMN_INFO_DATA_TYPE_FLOAT64)
}

func NewObjectIDColumnInfo(column string) *ColumnInfo {
	return newColumnInfoWithExportSettings(column, "in", COLUMN_INFO_DATA_TYPE_OBJECT_ID)
}

func NewDateColumnInfo(column string) *ColumnInfo {
	return NewColumnInfo(column, "between", COLUMN_INFO_DATA_TYPE_DATE)
}

func (columnInfo *ColumnInfo) SetPrepareValueFunc(prepareValueFunc func(interface{}) interface{}) *ColumnInfo {
	columnInfo.PrepareValueFunc = prepareValueFunc
	return columnInfo
}

func (columnInfo *ColumnInfo) GetFilterVal(filterVal string) bson.E {
	switch strings.ToLower(columnInfo.FilterOperator) {
	case ">":
		return bson.E{Key: columnInfo.Column, Value: bson.M{"$gt": columnInfo.convertValueToDataType(filterVal)}}
	case ">=":
		return bson.E{Key: columnInfo.Column, Value: bson.M{"$gte": columnInfo.convertValueToDataType(filterVal)}}
	case "<":
		return bson.E{Key: columnInfo.Column, Value: bson.M{"$lt": columnInfo.convertValueToDataType(filterVal)}}
	case "<=":
		return bson.E{Key: columnInfo.Column, Value: bson.M{"$lte": columnInfo.convertValueToDataType(filterVal)}}
	case "=":
		value := columnInfo.convertValueToDataType(filterVal)
		if columnInfo.DataType == COLUMN_INFO_DATA_TYPE_INT {
			if v, ok := value.(int); !ok || v == 0 {
				return bson.E{Key: columnInfo.Column, Value: bson.M{"$in": bson.A{nil, 0}}}
			}
		}
		return bson.E{Key: columnInfo.Column, Value: value}
	case "between":
		strValues := strings.Split(strings.TrimSpace(filterVal), ",")
		values := columnInfo.convertValueToArrayDataType(filterVal)
		if len(strValues) == 1 || strValues[1] == "" {
			return bson.E{Key: columnInfo.Column, Value: bson.M{"$gte": values[0]}}
		}

		if strValues[0] == "" {
			return bson.E{Key: columnInfo.Column, Value: bson.M{"$lte": values[1]}}
		} else {
			if values[0] == 0 || values[1] == 0 {
				return bson.E{
					Key: "$or", Value: bson.A{
						bson.M{columnInfo.Column: nil},
						bson.M{columnInfo.Column: bson.M{
							"$gte": values[0],
							"$lte": values[1],
						}},
					},
				}
			}
			return bson.E{
				Key: columnInfo.Column, Value: bson.M{
					"$gte": values[0],
					"$lte": values[1],
				},
			}
		}
	case "in":
		values := columnInfo.convertValueToArrayDataType(filterVal)
		if len(values) == 1 {
			return bson.E{Key: columnInfo.Column, Value: values[0]}
		} else {
			return bson.E{Key: columnInfo.Column, Value: bson.M{"$in": values}}
		}
	case "not in":
		values := columnInfo.convertValueToArrayDataType(filterVal)
		if len(values) == 1 {
			return bson.E{Key: columnInfo.Column, Value: bson.M{"$ne": values[0]}}
		} else {
			return bson.E{Key: columnInfo.Column, Value: bson.M{"$nin": values}}
		}
	case "like":
		return bson.E{Key: columnInfo.Column, Value: bson.Regex{Pattern: filterVal, Options: "i"}}
	case "is":
		if filterVal == "0" {
			return bson.E{Key: columnInfo.Column, Value: nil}
		}
		return bson.E{Key: columnInfo.Column, Value: bson.M{"$ne": nil}}
	case COLUMN_INFO_FILTER_OPERATOR_NONE:
		break
	}
	return bson.E{Key: columnInfo.Column, Value: nil}
}
func (columnInfo *ColumnInfo) getStringValue(value string) string {
	switch columnInfo.DataType {
	case COLUMN_INFO_DATA_TYPE_DATE:
		return dataType.ParseDateTime(value).Format(time.RFC3339)
	case "":
	case COLUMN_INFO_DATA_TYPE_TEXT:
		return value
	}
	return value
}

func (columnInfo *ColumnInfo) convertValueToDataType(value string) interface{} {
	var realValue any
	realValue = value
	switch columnInfo.DataType {
	case COLUMN_INFO_DATA_TYPE_DATE:
		if value == "" {
			realValue = nil
		} else {
			realValue = dataType.ParseDateTime(value)
		}
	case COLUMN_INFO_DATA_TYPE_OBJECT_ID:
		if value == "" {
			realValue = nil
		} else {
			realValue, _ = bson.ObjectIDFromHex(value)
		}
	case COLUMN_INFO_DATA_TYPE_INT:
		if value == "" {
			realValue = nil
		} else {
			realValue, _ = strconv.Atoi(value)
		}
	case COLUMN_INFO_DATA_TYPE_FLOAT64:
		if value == "" {
			return nil
		}
		realValue, _ = strconv.ParseFloat(value, 64)
	case COLUMN_INFO_DATA_TYPE_BOOL:
		if value == "" {
			return nil
		}
		realValue, _ = strconv.ParseBool(value)
	}

	if columnInfo.PrepareValueFunc != nil {
		return columnInfo.PrepareValueFunc(realValue)
	}
	if columnInfo.DataType == COLUMN_INFO_DATA_TYPE_INT {
		i, ok := realValue.(int)
		if !ok {
			return nil
		}
		return i
	}
	if columnInfo.DataType == COLUMN_INFO_DATA_TYPE_BOOL {
		i, ok := realValue.(bool)
		if !ok {
			return nil
		}
		if !i {
			return bson.M{"$in": bson.A{nil, false}}
		}
		return i
	}
	return realValue
}

func (columnInfo *ColumnInfo) convertValueToArrayDataType(value string) []interface{} {

	strValues := strings.Split(strings.TrimSpace(value), ",")
	values := []interface{}{}
	for _, strValue := range strValues {
		values = append(values, columnInfo.convertValueToDataType(strValue))
	}
	return values
}
