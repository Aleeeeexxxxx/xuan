package parser

import (
	"io"
	"xuan/src"

	"github.com/xuri/excelize/v2"
)

type RowBasedSheetParser interface {
	SheetName() string
	ParseRow(index int, row []string) error

	HeaderSize() int

	SetStorage(datastore src.Datastore)
}

type ExcelFileParser struct {
	f         *excelize.File
	datastore src.Datastore

	parsers []RowBasedSheetParser
}

func NewExcelFileParser(file io.Reader, datastore src.Datastore) (*ExcelFileParser, error) {
	f, err := excelize.OpenReader(file)
	if err != nil {
		return nil, err
	}

	return &ExcelFileParser{
		f:         f,
		datastore: datastore,
	}, nil
}

func (parser *ExcelFileParser) AddParser(p RowBasedSheetParser) {
	p.SetStorage(parser.datastore)
	parser.parsers = append(parser.parsers, p)
}

func (parser *ExcelFileParser) Parse() error {
	for _, p := range parser.parsers {
		if err := parser.runParser(p); err != nil {
			return err
		}
	}
	return nil
}

func (parser *ExcelFileParser) runParser(p RowBasedSheetParser) error {
	raw, err := parser.f.GetRows(p.SheetName())
	if err != nil {
		return err
	}

	header := p.HeaderSize()

	for index, row := range raw {
		if index < header {
			continue
		}

		if err := p.ParseRow(index, row); err != nil {
			return err
		}
	}
	return nil
}
