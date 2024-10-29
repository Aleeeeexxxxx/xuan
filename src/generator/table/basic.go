package table

import (
	"xuan/src"
	"xuan/src/excel"
)

type BasicTableGenerator struct {
}

func NewBasicTableGenerator(_ ...interface{}) TableGenerator {
	return &BasicTableGenerator{}
}

func (gen BasicTableGenerator) SheetName() string {
	return "电子元器件基本信息表"
}

func (gen BasicTableGenerator) GenHeader() excel.ExcelRow {
	return &excel.OneLineRow{
		Data: []interface{}{
			"序号",
			"元器件名称",
			"型号规格",
			"IP核 名称",
			"IP核 类型",
			"IP核 来源单位",
			"IP核 境内/境外",
			"零部件 名称",
			"零部件 是否“核心”/“关键”",
			"零部件 来源单位",
			"零部件 境内/境外",
			"流片工艺 名称",
			"流片工艺 境内/境外",
			"备注",
		},
	}
}

func (gen BasicTableGenerator) GenBodyForProduct(index int, p *src.Product) []excel.ExcelRow {
	var cells []excel.Cell = []excel.Cell{
		{Data: []interface{}{index + 1}},
		{Data: []interface{}{p.Name}},
		{Data: []interface{}{p.Model}},

		{Data: []interface{}{p.Core.Name}},
		{Data: []interface{}{p.Core.Type}},
		{Data: []interface{}{p.Core.Source}},
		{Data: []interface{}{p.Core.Domestic}},
	}

	var component []excel.Cell
	if p.Frame != nil {
		component = []excel.Cell{
			{Data: []interface{}{"晶圆", "框架/基板", "键合丝"}},
			{Data: []interface{}{p.Wafer.Important, p.Frame.Important, p.BondingWires.Important}},
			{Data: []interface{}{p.Wafer.Source, p.Frame.Source, p.BondingWires.Source}},
			{Data: []interface{}{p.Wafer.Domestic, p.Frame.Domestic, p.BondingWires.Domestic}},
		}
	} else {
		component = []excel.Cell{
			{Data: []interface{}{"晶圆", "管壳/盖板", "键合丝"}},
			{Data: []interface{}{p.Wafer.Important, p.TubeShell.Important, p.BondingWires.Important}},
			{Data: []interface{}{p.Wafer.Source, p.TubeShell.Source, p.BondingWires.Source}},
			{Data: []interface{}{p.Wafer.Domestic, p.TubeShell.Domestic, p.BondingWires.Domestic}},
		}
	}
	cells = append(cells, component...)

	cells = append(cells, []excel.Cell{
		{Data: []interface{}{p.Process.Name}},
		{Data: []interface{}{p.Process.Domestic}},
		{Data: []interface{}{"无"}},
	}...)

	return []excel.ExcelRow{&excel.MultiLineRow{Cells: cells}}
}
