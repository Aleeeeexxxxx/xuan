package plugin

import (
	"fmt"

	"xuan/src"
	"xuan/src/excel"
	"xuan/src/generator/table"
)

func FoundRelatedItems(target string, sources []string) []string {
	var result []string

	for _, source := range sources {
		length := CommonPrefixLength(target, source)
		if length > 7 {
			result = append(result, source)
		}
	}

	return result
}

func CommonPrefixLength(s1, s2 string) int {
	minLength := len(s1)
	if len(s2) < minLength {
		minLength = len(s2)
	}

	for i := 0; i < minLength; i++ {
		if s1[i] != s2[i] {
			return i
		}
	}

	return minLength
}

type Patcher struct {
	datastore src.Datastore

	sheets   map[string]*excel.Sheet
	notFound map[string][]*src.Product

	tableGenerators map[string]table.TableGenerator
}

func NewPatcher(datastore src.Datastore, genFactory []table.TableGeneratorFactory) *Patcher {
	s := &Patcher{
		datastore:       datastore,
		sheets:          make(map[string]*excel.Sheet),
		notFound:        make(map[string][]*src.Product),
		tableGenerators: make(map[string]table.TableGenerator),
	}
	s.registerTableGenerator(genFactory)
	return s
}

func (s *Patcher) Begin(sheetName string) {
	s.sheets[sheetName] = nil
}

func (s *Patcher) End(_ string) {}

func (s *Patcher) GenOneProduct(_ string) {}

func (s *Patcher) ProductNotFound(name string) {
	s.notFound[name] = []*src.Product{}
}

func (s *Patcher) Gen(ex *excel.Excel) {
	s.prepareProducts()
	s.prepareSheets(ex)

	for _, sheet := range s.sheets {
		var records []excel.ExcelRow

		records = append(records, &excel.OneLineRow{Data: []interface{}{""}})
		records = append(records, &excel.OneLineRow{Data: []interface{}{""}})
		records = append(records, &excel.OneLineRow{Data: []interface{}{""}})
		records = append(records, &excel.OneLineRow{Data: []interface{}{""}})

		for item, similar := range s.notFound {
			records = append(records,
				&excel.OneLineRow{Data: []interface{}{fmt.Sprintf("Found similar items for [%s]:", item)}})

			for index, product := range similar {
				if gen, ok := s.tableGenerators[sheet.Name]; ok {
					records = append(records, gen.GenBodyForProduct(index, product)...)
				}
			}

			records = append(records, &excel.OneLineRow{Data: []interface{}{""}})
		}

		sheet.Rows = append(sheet.Rows, records...)
	}
}

func (s *Patcher) prepareProducts() {
	products, _ := s.datastore.GetProductList()

	for item, _ := range s.notFound {
		similar := FoundRelatedItems(item, products)

		if len(similar) > 0 {
			s.notFound[item] = make([]*src.Product, 0, len(similar))

			for _, name := range similar {
				product, _ := s.datastore.GetProduct(name) // TODO err handling
				s.notFound[item] = append(s.notFound[item], product)
			}
		}
	}
}

func (s *Patcher) prepareSheets(ex *excel.Excel) {
	for _, sheet := range ex.Sheets {
		if _, ok := s.sheets[sheet.Name]; ok {
			s.sheets[sheet.Name] = sheet
		}
	}
}

func (s *Patcher) registerTableGenerator(genFactory []table.TableGeneratorFactory) {
	for _, factory := range genFactory {
		gen := factory()
		s.tableGenerators[gen.SheetName()] = gen
	}
}
