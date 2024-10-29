package excel

import (
	"fmt"
	"io"

	"github.com/xuri/excelize/v2"
)

const defaultSheet = "sheet1"

type Excel struct {
	DefaultSheet Sheet
	Sheets       []Sheet
}

func NewExcel() *Excel {
	return &Excel{
		DefaultSheet: Sheet{
			Name: defaultSheet,
		},
	}
}

func (excel *Excel) Write(w io.Writer) error {
	f := excelize.NewFile()

	for _, sheet := range excel.Sheets {
		if err := sheet.Write(f); err != nil {
			return fmt.Errorf("failed to write sheet(%s): %w", sheet.Name, err)
		}
	}

	if err := excel.DefaultSheet.Write(f); err != nil {
		return fmt.Errorf("failed to write to default sheet: %w", err)
	}

	if err := f.Write(w); err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	return nil
}

func (excel *Excel) AddRowsToDefaultSheet(rows []ExcelRow) {
	excel.DefaultSheet.Rows = append(excel.DefaultSheet.Rows, rows...)
}

func (excel *Excel) AddSheet(sheet Sheet) {
	excel.Sheets = append(excel.Sheets, sheet)
}
