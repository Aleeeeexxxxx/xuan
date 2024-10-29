package sheet

import "xuan/src"

type AllInOneParser struct {
	datastore src.Datastore
}

func (parser *AllInOneParser) SheetName() string {
	return "AllInOne数据表"
}

func (parser *AllInOneParser) HeaderSize() int {
	return 2
}

func (parser *AllInOneParser) SetStorage(datastore src.Datastore) {
	parser.datastore = datastore
}

func (parser *AllInOneParser) ParseRow(index int, row []string) error {
	model := row[7]
	p, _ := GetProductCreateIfNotExist(model, parser.datastore)

	p.Name = row[8]
	p.Classifications = []src.Classification{
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
	}

	p.ZZKKLevel = row[46]
	p.W = row[49]
	p.K = row[50]
	p.B = row[51]
	p.WKBDetails = row[52]

	p.Process = &src.Process{
		Name: row[33],
	}

	return parser.datastore.AddProduct(p)
}
