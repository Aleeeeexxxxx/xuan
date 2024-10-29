package plugin

import "xuan/src/excel"

type SheetStatistic struct {
	SheetName string

	Generated int
	Total     int
}

func (stat *SheetStatistic) ToExcelRows() []excel.ExcelRow {
	var ret []excel.ExcelRow = []excel.ExcelRow{
		&excel.OneLineRow{Data: []interface{}{stat.SheetName}},

		&excel.OneLineRow{Data: []interface{}{"Total:", stat.Total}},
		&excel.OneLineRow{Data: []interface{}{"Generated:", stat.Generated}},
	}

	return ret
}

type Statisticer struct {
	cur   *SheetStatistic
	stats map[string]*SheetStatistic

	NotFound []string
}

func NewStatisticer() *Statisticer {
	return &Statisticer{
		stats: make(map[string]*SheetStatistic),
	}
}

func (s *Statisticer) Begin(sheetName string) {
	s.cur = &SheetStatistic{SheetName: sheetName}
}

func (s *Statisticer) End(sheetName string) {
	s.stats[sheetName] = s.cur
	s.cur = nil
}

func (s *Statisticer) GenOneProduct(name string) {
	s.cur.Total++
	s.cur.Generated++
}

func (s *Statisticer) ProductNotFound(name string) {
	s.cur.Total++
	s.NotFound = append(s.NotFound, name)
}

func (s *Statisticer) Gen(ex *excel.Excel) {
	var ret []excel.ExcelRow

	ret = append(ret, &excel.OneLineRow{Data: []interface{}{"NotFound:", s.NotFound}})
	if len(s.NotFound) > 0 {
		for _, item := range s.NotFound {
			ret = append(ret, &excel.MultiLineRow{
				Cells: []excel.Cell{
					{Data: []interface{}{""}},
					{Data: []interface{}{item}},
				},
			})
		}
	}

	for _, stat := range s.stats {
		ret = append(ret, &excel.OneLineRow{Data: []interface{}{""}})
		ret = append(ret, stat.ToExcelRows()...)
	}

	ex.AddRowsToDefaultSheet(ret)
}
