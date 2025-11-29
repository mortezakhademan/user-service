package componentsList

import (
	"context"
	"fmt"
	"github.com/xuri/excelize/v2"
	ptime "github.com/yaa110/go-persian-calendar"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"reflect"
	"strings"

	"git.ramooz.org/ramooz/golang-components/paginated-list/mongodb/models"
)

func (list *List) prepareExcelResult(cursor *mongo.Cursor, context context.Context, dataSliceType interface{}) error {

	list.ExcelFile = excelize.NewFile()
	sheetName := "Sheet1"
	_, err := list.ExcelFile.NewSheet(sheetName)
	if err != nil {
		fmt.Printf("can't create sheet in export excel : %s", err)
		return err
	}
	if list.exportColumnsInfo == nil {
		list.setDefaultExportColumns(dataSliceType)
	}
	list.setExcelHeaderRow(sheetName)
	rowIndex := 2
	for cursor.Next(context) {
		documentData := map[string]interface{}{}
		cursor.Decode(documentData)
		list.setExcelDataRow(sheetName, rowIndex, documentData)
		rowIndex++
	}

	return nil
}

func (list *List) setDefaultExportColumns(dataSliceType any) {
	tags := getSliceObjectJsonTags(dataSliceType)
	columns := make([]*models.ExportColumnInfo, 0, len(tags))
	for column, columnInfo := range list.fromCollectionColumnsMap {
		if _, ok := tags[column]; ok {
			columns = append(columns, &models.ExportColumnInfo{Column: columnInfo.Column, Title: columnInfo.Column, DataType: columnInfo.DataType})
		}
	}

	for _, dbPipeline := range list.dbPipelines {
		for column, columnInfo := range dbPipeline.ResultColumns {
			if _, ok := tags[column]; ok {
				columns = append(columns, &models.ExportColumnInfo{Column: columnInfo.Column, Title: columnInfo.Column, DataType: columnInfo.DataType})
			}
		}
	}
	list.SetExportColumns(columns)
}

// return map[json tag]bson tag
func getSliceObjectJsonTags(dataSliceType interface{}) map[string]string {
	tagNames := map[string]string{}
	t := reflect.TypeOf(dataSliceType).Elem()
	if t.Kind() == reflect.Slice && t.Elem().Kind() == reflect.Struct {
		// Proceed with reflection if slice is not empty and contains structs
		for i := 0; i < t.Elem().NumField(); i++ {
			field := t.Elem().Field(i)
			jsonTag := field.Tag.Get("json")
			if jsonTag == "" {
				continue
			}
			bsonTag := field.Tag.Get("bson")
			if bsonTag == "" {
				continue
			}
			jsonTag = strings.Split(jsonTag, ",")[0]
			if jsonTag == "" || jsonTag == "-" {
				continue
			}
			bsonTag = strings.Split(bsonTag, ",")[0]
			if bsonTag == "" || bsonTag == "-" {
				continue
			}
			tagNames[jsonTag] = bsonTag
		}
	}
	return tagNames
}

func (list *List) setExcelHeaderRow(sheetName string) {
	columns := list.exportColumnsInfo
	var col rune
	for i, column := range columns {
		title := column.Title
		col = 'A' + int32(i)
		err := list.ExcelFile.SetCellValue(sheetName, fmt.Sprintf("%s%d", string(col), 1), title)
		if err != nil {
			fmt.Printf("can't set cell value in export excel : %s", err)
			return
		}
	}
}

func (list *List) setExcelDataRow(sheetName string, rowIndex int, documentData map[string]interface{}) {
	// Create our map, and retrieve the value for each column from the pointers slice,
	//storing it in the map with the name of the column as the key.
	var col rune
	for i, column := range list.exportColumnsInfo {
		col = 'A' + int32(i)
		if column.SetValueFunc != nil {
			err := column.SetValueFunc(list.ExcelFile, sheetName, rowIndex, string(col), documentData)
			if err != nil {
				fmt.Printf("can't set cell value in export excel : %s", err)
				return
			}
			continue
		}
		var colValStr string
		if column.PrepareValueFunc != nil {
			colValStr = column.PrepareValueFunc(getColumnValue(column.Column, documentData))
		} else {
			colValStr = getColumnStringValue(column, documentData)
		}
		err := list.ExcelFile.SetCellValue(sheetName, fmt.Sprintf("%s%d", string(col), rowIndex), colValStr)
		if err != nil {
			fmt.Printf("can't set cell value in export excel : %s", err)
			return
		}
	}

}

func getColumnStringValue(column *models.ExportColumnInfo, documentData map[string]interface{}) string {
	colVal := getColumnValue(column.Column, documentData)
	if colVal == nil {
		return ""
	}
	switch column.DataType {
	case models.COLUMN_INFO_DATA_TYPE_DATE:
		if data, ok := colVal.(bson.DateTime); ok {
			return ptime.New(data.Time()).Format("yyyy/MM/dd HH:mm")
		}
		return ""
	case models.COLUMN_INFO_DATA_TYPE_OBJECT_ID:
		if data, ok := colVal.(bson.ObjectID); ok {
			return data.Hex()
		}
		return ""
	}
	b, ok := colVal.([]byte)
	if ok {
		return string(b)
	}
	return fmt.Sprint(colVal)
}

func getColumnValue(columnField string, documentData map[string]interface{}) any {
	columnPath := strings.SplitN(columnField, ".", 2)
	if v, ok := documentData[columnPath[0]]; ok {
		if len(columnPath) > 1 {
			if colVal, ok := v.(map[string]any); ok {
				return getColumnValue(columnPath[1], colVal)
			}
			return nil
		}
		return v
	}
	return nil
}
