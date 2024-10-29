package excel

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

type Sheet struct {
	Name string
	Rows []ExcelRow
}

func writeMultiExcelRows(f *excelize.File, sheetName string, startLine int, rows []ExcelRow) (int, error) {
	line := startLine

	for _, row := range rows {
		next, err := row.Write(f, sheetName, line)
		if err != nil {
			return -1, fmt.Errorf("failed to write row: %w", err)
		}
		line = next
	}

	return line, nil
}

func (sheet *Sheet) Write(f *excelize.File) error {
	if sheet.Name != defaultSheet {
		_, err := f.NewSheet(sheet.Name)
		if err != nil {
			return fmt.Errorf("failed to create sheet(%s): %w", sheet.Name, err)
		}
	}

	_, err := writeMultiExcelRows(f, sheet.Name, 1, sheet.Rows)
	return err
}
