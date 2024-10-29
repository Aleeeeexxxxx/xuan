package table

import (
	"xuan/src"
	"xuan/src/excel"
)

type TableGeneratorFactory func(args ...interface{}) TableGenerator

type TableGenerator interface {
	SheetName() string
	GenHeader() excel.ExcelRow
	GenBodyForProduct(index int, product *src.Product) []excel.ExcelRow
}
