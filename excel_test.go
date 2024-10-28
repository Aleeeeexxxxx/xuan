package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExcel(t *testing.T) {
	rq := require.New(t)

	excel := &Excel{
		Path: "test.xlsx",
		Sheets: []Sheet{
			{
				Name: "Sheet1",
				Rows: []ExcelRow{
					&MultiLineRow{
						Cells: []Cell{
							{
								Data: []interface{}{"1"},
							},
							{
								Data: []interface{}{"1", "2"},
							},
							{
								Data: []interface{}{"1", "2", "3"},
							},
						},
					},
				},
			},
		},
	}

	rq.Nil(excel.Write())
}
