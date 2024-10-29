package sheet

import "xuan/src"

type WKB2Parser struct {
	datastore src.Datastore
}

func (parser *WKB2Parser) HeaderSize() int {
	return 3
}

func (parser *WKB2Parser) SheetName() string {
	return "【输出】WKB表2"
}

func (parser *WKB2Parser) SetStorage(datastore src.Datastore) {
	parser.datastore = datastore
}

func (parser *WKB2Parser) ParseRow(index int, data []string) error {
	model := data[5]
	if model == "#N/A" {
		return nil
	}

	p, _ := GetProductCreateIfNotExist(model, parser.datastore)

	p.Core = &src.IPCore{
		Name:     data[6],
		Type:     data[7],
		Source:   data[8],
		Domestic: data[9],
	}

	// 11-13: 晶圆
	p.Wafer = &src.Component{
		Important: data[11],
		Source:    data[12],
		Domestic:  data[13],
	}
	if data[16] == "/" {
		// 23-25: 框架/基板
		p.Frame = &src.Component{
			Important: data[23],
			Source:    data[24],
			Domestic:  data[25],
		}
	} else {
		// 15-17: 管壳
		p.TubeShell = &src.Component{
			Important: data[14],
			Source:    data[15],
			Domestic:  data[16],
		}
		// 19-21: 盖板
		p.Panel = &src.Component{
			Important: data[19],
			Source:    data[20],
			Domestic:  data[21],
		}

	}
	// 27-29: 键合丝
	p.BondingWires = &src.Component{
		Important: data[27],
		Source:    data[28],
		Domestic:  data[29],
	}

	// 流片 == 晶圆
	if p.Process == nil {
		p.Process = &src.Process{}
	}
	p.Process.Domestic = p.Wafer.Domestic

	return parser.datastore.AddProduct(p)
}

func GetProductCreateIfNotExist(model string, datastore src.Datastore) (*src.Product, error) {
	p, err := datastore.GetProduct(model)
	if err != nil {
		if err == src.ErrorProductNotFound {
			p = &src.Product{
				Model: model,
			}
		} else {
			return nil, err
		}
	}
	return p, nil
}
