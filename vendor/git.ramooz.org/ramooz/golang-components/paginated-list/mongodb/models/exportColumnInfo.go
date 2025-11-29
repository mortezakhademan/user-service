package models

import "github.com/xuri/excelize/v2"

// ExportColumnInfo represents the information of a column to be exported.
type ExportColumnInfo struct {
	// Column is the name of the column in the data source.
	Column string

	// Title is the title of the column to be displayed in the export.
	Title string

	// DataType is the type of data contained in the column.
	DataType string

	// PrepareValueFunc is a function that prepares the value for export.
	// It takes an interface{} as input and returns a string.
	PrepareValueFunc func(interface{}) string

	// SetValueFunc is a function that sets the value in the exported file.
	// It takes an excelize.File, sheet name, row number, column letter, and row data as input.
	// It returns an error if the value could not be set.
	SetValueFunc func(excel *excelize.File, sheetName string, row int, col string, rowData map[string]any) error
}
