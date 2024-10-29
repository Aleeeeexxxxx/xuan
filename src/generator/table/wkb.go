package table

import (
	"xuan/src"
	"xuan/src/excel"
)

type WKBTableGenerator struct{}

func NewWKBTableGenerator(_ ...interface{}) TableGenerator {
	return &WKBTableGenerator{}
}

func (gen WKBTableGenerator) SheetName() string {
	return "伪、空、包审查信息一览表"
}

func (gen WKBTableGenerator) GenHeader() excel.ExcelRow {
	return &excel.OneLineRow{
		Data: []interface{}{
			"序号",
			"元器件名称",
			"型号规格",
			"生产厂家",
			"自主可控等级",
			"伪国产化",
			"伪国产化备注",
			"空心国产化",
			"空心国产化备注",
			"包装国产化",
			"包装国产化备注",
			"备注",
		},
	}
}

func (gen WKBTableGenerator) GenBodyForProduct(index int, p *src.Product) []excel.ExcelRow {
	return []excel.ExcelRow{&excel.OneLineRow{
		Data: []interface{}{
			index + 1,
			p.Name,
			p.Model,
			"成都华微电子科技股份有限公司",
			p.ZZKKLevel,
			p.W,
			"/",
			p.B,
			"/",
			p.K,
			"/",
			"无",
		},
	}}
}
