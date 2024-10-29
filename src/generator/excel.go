package generator

import (
	"xuan/src"
	"xuan/src/excel"
)

type TableGeneratorFactory func(args ...interface{}) TableGenerator

type ExcelGenerator struct {
	datastore src.Datastore
	targets   []string

	genFactory []TableGeneratorFactory

	ex     *excel.Excel
	plugin Plugin
}

func NewExcelGenerator(datastore src.Datastore, targets []string) *ExcelGenerator {
	gen := &ExcelGenerator{
		datastore: datastore,
		targets:   targets,
		ex:        excel.NewExcel(),
		plugin:    NewStatisticer(),
	}
	return gen
}

func (gen *ExcelGenerator) AddTable(factory TableGeneratorFactory) {
	gen.genFactory = append(gen.genFactory, factory)
}

func (gen *ExcelGenerator) Gen() (*excel.Excel, error) {
	for _, factory := range gen.genFactory {
		gen.genSheet(factory)
	}

	gen.plugin.Gen(gen.ex)
	return gen.ex, nil
}

func (gen *ExcelGenerator) genSheet(factory TableGeneratorFactory) error {
	tg := factory()
	sheetName := tg.SheetName()

	gen.plugin.Begin(sheetName)
	defer gen.plugin.End(sheetName)

	var records []excel.ExcelRow

	header := tg.GenHeader()
	records = append(records, header)

	for index, target := range gen.targets {
		p, err := gen.datastore.GetProduct(target)
		if err != nil {
			if err == src.ErrorProductNotFound {
				gen.plugin.ProductNotFound(target)
				continue
			} else {
				return err
			}
		}

		records = append(records, tg.GenBodyForProduct(index, p)...)
		gen.plugin.GenOneProduct(target)
	}

	gen.ex.AddSheet(excel.Sheet{
		Name: sheetName,
		Rows: records,
	})
	return nil
}
