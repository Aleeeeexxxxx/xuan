package src

import "errors"

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

type Datastore interface {
	GetProduct(name string) (*Product, error)
	GetProductList() ([]string, error)

	AddProduct(p *Product) error
}

var ErrorProductNotFound = errors.New("product not found")

type InMemoryDatastore struct {
	products map[string]*Product
}

func NewInMemoryDatastore() *InMemoryDatastore {
	return &InMemoryDatastore{
		products: make(map[string]*Product),
	}
}

func (store *InMemoryDatastore) GetProduct(name string) (*Product, error) {
	if p, ok := store.products[name]; ok {
		return p, nil
	}
	return nil, ErrorProductNotFound
}

func (store *InMemoryDatastore) GetProductList() ([]string, error) {
	ret := make([]string, 0, len(store.products))
	for name := range store.products {
		ret = append(ret, name)
	}
	return ret, nil
}

func (store *InMemoryDatastore) AddProduct(p *Product) error {
	store.products[p.Model] = p
	return nil
}
