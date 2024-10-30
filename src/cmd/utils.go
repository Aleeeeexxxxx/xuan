package cmd

import (
	"io"

	"xuan/src"
	"xuan/src/generator"
	"xuan/src/generator/table"
	"xuan/src/parser"
	"xuan/src/parser/sheet"
)

func PanicIfNotNil(err error) {
	if err != nil {
		panic(err)
	}
}

func RunGeneratorOrPanic(source io.Reader, output io.Writer, targets []string) {
	datastore := src.NewInMemoryDatastore()

	parser, err := parser.NewExcelFileParser(source, datastore)
	PanicIfNotNil(err)

	parser.AddParser(&sheet.AllInOneParser{})
	parser.AddParser(&sheet.WKB2Parser{})

	PanicIfNotNil(parser.Parse())

	gen := generator.NewExcelGenerator(datastore, targets)
	gen.AddTable(table.NewBasicTableGenerator)
	gen.AddTable(table.NewWKBTableGenerator)

	excel, err := gen.Gen()
	PanicIfNotNil(err)

	PanicIfNotNil(excel.Write(output))
}
