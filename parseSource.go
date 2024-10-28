package main

import (
	"fmt"
	"strings"

	"github.com/xuri/excelize/v2"
)

type Classification struct {
	Criteria string // 产品分类标准
	Level1   string // 一级分类
	Level2   string // 二级分类
	Level3   string // 三级分类
}

type Process struct {
	Name     string // 名称
	Domestic string // 境内/境外
}

type IPCore struct {
	Name     string // 名称
	Type     string // 类型
	Source   string // 来源单位
	Domestic string // 境内/境外
}

type Component struct {
	Important string // 是否核心/重要
	Domestic  string // 境内/境外
	Source    string // 来源单位
}

type Product struct {
	Model string // 型号
	Name  string // 名称

	Classifications []Classification // 产品分类

	ZZKKLevel  string // 自主可控等级
	W          string // 伪国产化
	K          string // 空心国产化
	B          string // 包装国产化
	WKBDetails string // 伪空包说明

	Core         *IPCore    // IP核
	Wafer        *Component // 晶圆
	TubeShell    *Component // 管壳
	Panel        *Component // 盖板
	Frame        *Component // 框架/基板
	BondingWires *Component // 键合丝

	Process *Process // 流片工艺
}

type Source struct {
	FilePath string
	Products map[string]*Product
}

func (source *Source) ParsedProducts() []string {
	var ret []string
	for k := range source.Products {
		ret = append(ret, k)
	}
	return ret
}

type SourceReader struct {
	path string
	f    *excelize.File

	basic    *AllInOneReader
	enhancer *WKBEnhancer
}

func NewSourceReader(path string) (*SourceReader, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	reader := &SourceReader{
		f:     f,
		path:  path,
		basic: &AllInOneReader{f: f},
	}

	return reader, nil
}

func (reader *SourceReader) Parse() (*Source, error) {
	products, err := reader.basic.Parse()
	if err != nil {
		return nil, fmt.Errorf("failed to parse basic data: %w", err)
	}

	reader.enhancer = NewWKBEnhancer(reader.f, products)

	products, err = reader.enhancer.Enhance()
	if err != nil {
		return nil, fmt.Errorf("failed to enhance basic data: %w", err)
	}

	return &Source{
		FilePath: reader.path,
		Products: products,
	}, nil

}

type AllInOneReader struct {
	f *excelize.File
}

func (reader *AllInOneReader) Parse() (map[string]*Product, error) {
	const sheet = "AllInOne数据表"

	raw, err := reader.f.GetRows(sheet)
	if err != nil {
		return nil, fmt.Errorf("failed to get sheet(%s): %w", sheet, err)
	}

	rows, err := validate(raw, 2)
	if err != nil {
		return nil, fmt.Errorf("failed to validate: %w", err)
	}

	products := make(map[string]*Product)

	for _, row := range rows {
		product := reader.parseRow(row)
		products[product.Model] = product
	}

	fmt.Printf("parsed [%s], total products=%d\n", sheet, len(products))

	return products, nil
}

func (reader *AllInOneReader) parseRow(row []string) *Product {
	return &Product{
		Model: row[7],
		Name:  row[8],

		Classifications: []Classification{
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

		Process: &Process{
			Name:     row[33],
			Domestic: "", // TODO
		},
	}
}

type WKBEnhancer struct {
	f            *excelize.File
	base         map[string]*Product
	rowByProduct map[string][][]string
}

func NewWKBEnhancer(f *excelize.File, base map[string]*Product) *WKBEnhancer {
	return &WKBEnhancer{
		f:    f,
		base: base,
	}
}

func (enhancer *WKBEnhancer) Enhance() (map[string]*Product, error) {
	const sheet = "【输出】WKB表2"

	raw, err := enhancer.f.GetRows(sheet)
	if err != nil {
		return nil, fmt.Errorf("failed to get sheet(%s): %w", sheet, err)
	}

	rows, err := validate(raw, 3)
	if err != nil {
		return nil, fmt.Errorf("failed to validate: %w", err)
	}

	enhancer.parse(rows)

	return enhancer.enhance(), nil
}

func (enhancer *WKBEnhancer) parse(rows [][]string) {
	enhancer.rowByProduct = make(map[string][][]string)

	for _, row := range rows {
		product := row[5]
		if product == "#N/A" {
			continue
		}

		if _, ok := enhancer.rowByProduct[product]; !ok {
			enhancer.rowByProduct[product] = make([][]string, 0)
		}
		enhancer.rowByProduct[product] = append(enhancer.rowByProduct[product], row)
	}

	for model, rows := range enhancer.rowByProduct {
		if len(rows) >= 2 {
			fmt.Printf("found multiple rows for product: %s, row=%d", model, len(rows))
		}
	}
}

func (enhancer *WKBEnhancer) enhance() map[string]*Product {
	for model, p := range enhancer.base {
		if rows, ok := enhancer.rowByProduct[model]; ok {
			data := rows[0] // should be only one row

			p.Core = &IPCore{
				Name:     data[6],
				Type:     data[7],
				Source:   data[8],
				Domestic: data[9],
			}

			// 11-13: 晶圆
			p.Wafer = &Component{
				Important: data[11],
				Source:    data[12],
				Domestic:  data[13],
			}
			if data[16] == "/" {
				// 23-25: 框架/基板
				p.Frame = &Component{
					Important: data[23],
					Source:    data[24],
					Domestic:  data[25],
				}
			} else {
				// 15-17: 管壳
				p.TubeShell = &Component{
					Important: data[14],
					Source:    data[15],
					Domestic:  data[16],
				}
				// 19-21: 盖板
				p.Panel = &Component{
					Important: data[19],
					Source:    data[20],
					Domestic:  data[21],
				}

			}
			// 27-29: 键合丝
			p.BondingWires = &Component{
				Important: data[27],
				Source:    data[28],
				Domestic:  data[29],
			}

			// 流片 == 晶圆
			p.Process.Domestic = p.Wafer.Domestic
		} else {
			fmt.Printf("not found product: %s", model)
		}
	}

	return enhancer.base
}

/*
 * Common functions
 */

func validate(raw [][]string, headerLevel int) ([][]string, error) {
	if len(raw) <= headerLevel {
		return nil, fmt.Errorf("invalid data")
	}

	// PrintHeader(raw[headerLevel-1])

	return raw[headerLevel:], nil
}

func PrintHeader(header []string) {
	fmt.Println("[ Header: ===========")
	defer fmt.Println("] Header: ===========")

	for i, h := range header {
		h := strings.ReplaceAll(h, "\n", " ")
		fmt.Printf("%d: %s\n", i, h)
	}
}
