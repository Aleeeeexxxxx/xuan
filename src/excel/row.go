package excel

import (
	"log"

	"github.com/xuri/excelize/v2"
)

type ExcelRow interface {
	Write(f *excelize.File, sheetName string, startLine int) (int, error)
}

type OneLineRow struct {
	Data []interface{}
}

func (row *OneLineRow) Write(f *excelize.File, sheetName string, startLine int) (int, error) {
	for i, data := range row.Data {
		cell, err := excelize.CoordinatesToCellName(i+1, startLine)
		if err != nil {
			log.Fatalf("failed to get cell name: %v", err)
		}
		f.SetCellValue(sheetName, cell, data)
	}

	return startLine + 1, nil
}

type Cell struct {
	Data []interface{}
}

type MultiLineRow struct {
	Cells []Cell
}

func (row *MultiLineRow) Write(f *excelize.File, sheetName string, startLine int) (int, error) {
	maxLine := 1

	for i, cell := range row.Cells {
		if len(cell.Data) > maxLine {
			maxLine = len(cell.Data)
		}

		for j, data := range cell.Data {
			excelCell, err := excelize.CoordinatesToCellName(i+1, startLine+j)
			if err != nil {
				log.Fatalf("failed to get cell name: %v", err)
			}
			f.SetCellValue(sheetName, excelCell, data)
		}
	}

	// merge cells
	if maxLine > 1 {
		lastMaxRowNumber := startLine + maxLine - 1

		for i, cell := range row.Cells {
			if len(cell.Data) == maxLine {
				continue
			}

			lastRowNumber := startLine + len(cell.Data) - 1
			lastCell, err := excelize.CoordinatesToCellName(i+1, lastRowNumber)
			if err != nil {
				log.Fatalf("failed to get cell name: %v", err)
			}

			lastExpectedCell, err := excelize.CoordinatesToCellName(i+1, lastMaxRowNumber)
			if err != nil {
				log.Fatalf("failed to get cell name: %v", err)
			}

			if err := f.MergeCell(sheetName, lastCell, lastExpectedCell); err != nil {
				log.Fatalf("failed to merge cells: %v", err)
			}
		}
	}

	return startLine + maxLine, nil
}
