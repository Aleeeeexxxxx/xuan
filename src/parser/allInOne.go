package parser

import "xuan/src"

type AllInOneParser struct {
	datastore src.Datastore
}

func (parser *AllInOneParser) SheetName() string {
	return "AllInOne数据表"
}

func (parser *AllInOneParser) SetStorage(datastore src.Datastore) {
	parser.datastore = datastore
}

func (parser *AllInOneParser) ParseRow(index int, row []string) error {
	product := &src.Product{
		Model: row[7],
		Name:  row[8],

		Classifications: []src.Classification{
			{
				Criteria: "WKB判定文件3101",
				Level1:   row[1],
				Level2:   row[2],
				Level3:   row[3],
			},
			{
				Criteria: "GJB8118",
				Level1:   row[4],
				Level2:   row[5],
				Level3:   row[6],
			},
		},

		ZZKKLevel:  row[46],
		W:          row[49],
		K:          row[50],
		B:          row[51],
		WKBDetails: row[52],

		Process: &src.Process{
			Name:     row[33],
		},
	}

	return parser.datastore.AddProduct(product)
}
