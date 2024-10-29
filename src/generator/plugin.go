package generator

import "xuan/src/excel"

type Plugin interface {
	Begin(sheetName string)
	End(sheetName string)

	GenOneProduct(name string)
	ProductNotFound(name string)

	Gen(*excel.Excel)
}

type NotFoundItem struct {
	Name    string
	Similar []string
}

func (item NotFoundItem) ToExcelRows() excel.ExcelRow {
	relatedItems := excel.Cell{}
	for _, i := range item.Similar {
		relatedItems.Data = append(relatedItems.Data, i)
	}

	return &excel.MultiLineRow{
		Cells: []excel.Cell{
			{Data: []interface{}{""}},
			{Data: []interface{}{item.Name}},
			relatedItems,
		},
	}
}

type SheetStatistic struct {
	SheetName string
	Generated int
	NotFound  []NotFoundItem
}

func (stat *SheetStatistic) ToExcelRows() []excel.ExcelRow {
	var ret []excel.ExcelRow = []excel.ExcelRow{
		&excel.OneLineRow{Data: []interface{}{stat.SheetName}},

		&excel.OneLineRow{Data: []interface{}{"Total:", stat.Generated + len(stat.NotFound)}},
		&excel.OneLineRow{Data: []interface{}{"Generated:", stat.Generated}},
		&excel.OneLineRow{Data: []interface{}{"Not Found:", len(stat.NotFound)}},
	}

	if len(stat.NotFound) > 0 {
		for _, item := range stat.NotFound {
			ret = append(ret, item.ToExcelRows())
		}
	}
	return ret
}

type Statisticer struct {
	cur   *SheetStatistic
	stats map[string]*SheetStatistic
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
	s.cur.Generated++
}

func (s *Statisticer) ProductNotFound(name string) {
	s.cur.NotFound = append(s.cur.NotFound, NotFoundItem{Name: name})
}

func (s *Statisticer) Gen(ex *excel.Excel) {
	var ret []excel.ExcelRow
	for _, stat := range s.stats {
		ret = append(ret, stat.ToExcelRows()...)
		ret = append(ret, &excel.OneLineRow{Data: []interface{}{""}})
	}

	ex.AddRowsToDefaultSheet(ret)
}
