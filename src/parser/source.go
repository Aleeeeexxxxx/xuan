package parser

import "xuan/src"

type ExcelFileSource struct {
	filePath string
	version  string

	products map[string]*src.Product
}

func NewExcelFileSource(filePath string) *ExcelFileSource {
	return &ExcelFileSource{
		filePath: filePath, 
		products: make(map[string]*src.Product),
	}
}

func (source *ExcelFileSource) GetProduct(name string) (*src.Product, error) {
	if p, ok := source.products[name]; ok {
		return p, nil
	}
	return nil, src.ErrorProductNotFound
}

func (source *ExcelFileSource) GetProductList() ([]*src.Product, error) {
	ret := make([]*src.Product, 0, len(source.products))
	for _, p := range source.products {
		ret = append(ret, p)
	}
	return ret, nil
}
