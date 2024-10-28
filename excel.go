package main

import (
	"fmt"
	"log"

	"github.com/xuri/excelize/v2"
)

type Excel struct {
	Path   string
	Sheets []Sheet
}

func (excel *Excel) Write() error {
	f := excelize.NewFile()
	var rows []ExcelRow

	for _, sheet := range excel.Sheets {
		if err := sheet.Write(f); err != nil {
			return fmt.Errorf("failed to write sheet(%s): %w", sheet.Name, err)
		}

		rows = append(rows, sheet.DescriptionRows()...)
	}

	if _, err := writeMultiExcelRows(f, "sheet1", 1, rows); err != nil {
		return fmt.Errorf("failed to write sheet1: %w", err)
	}

	if err := f.SaveAs(excel.Path); err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	return nil
}

type NotFoundItem struct {
	Name         string
	RelatedItems []string
}

func (item NotFoundItem) toRow() ExcelRow {
	relatedItems := Cell{}
	for _, i := range item.RelatedItems {
		relatedItems.Data = append(relatedItems.Data, i)
	}

	return &MultiLineRow{
		Cells: []Cell{
			{Data: []interface{}{""}},
			{Data: []interface{}{item.Name}},
			relatedItems,
		},
	}
}

type Description struct {
	Total         int
	Found         int
	NotFoundItems []NotFoundItem
}

func (desc Description) toRows() []ExcelRow {
	var ret []ExcelRow = []ExcelRow{
		&OneLineRow{
			Data: []interface{}{"Total:", desc.Total},
		},
		&OneLineRow{
			Data: []interface{}{"Found:", desc.Found},
		},
		&OneLineRow{
			Data: []interface{}{"Not Found:", len(desc.NotFoundItems)},
		},
	}

	if len(desc.NotFoundItems) > 0 {
		for _, item := range desc.NotFoundItems {
			ret = append(ret, item.toRow())
		}
	}

	return ret
}

type Sheet struct {
	Desc Description
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
	_, err := f.NewSheet(sheet.Name)
	if err != nil {
		return fmt.Errorf("failed to create sheet(%s): %w", sheet.Name, err)
	}

	_, err = writeMultiExcelRows(f, sheet.Name, 1, sheet.Rows)
	return err
}

func (sheet Sheet) DescriptionRows() []ExcelRow {
	ret := []ExcelRow{&OneLineRow{Data: []interface{}{sheet.Name}}}
	ret = append(ret, sheet.Desc.toRows()...)
	ret = append(ret, &OneLineRow{Data: []interface{}{""}})
	return ret
}

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
