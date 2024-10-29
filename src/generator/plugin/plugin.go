package plugin

import (
	"xuan/src"
	"xuan/src/excel"
	"xuan/src/generator/table"
)

type Plugin interface {
	Begin(sheetName string)
	End(sheetName string)

	GenOneProduct(name string)
	ProductNotFound(name string)

	Gen(*excel.Excel)
}

type PluginMngr struct {
	plugins []Plugin
}

func NewPluginMngr(datastore src.Datastore, genFactory []table.TableGeneratorFactory) *PluginMngr {
	p := &PluginMngr{}

	p.plugins = append(p.plugins, NewStatisticer())
	p.plugins = append(p.plugins, NewPatcher(datastore, genFactory))
	return p
}

func (mngr *PluginMngr) Begin(sheetName string) {
	for _, p := range mngr.plugins {
		p.Begin(sheetName)
	}
}

func (mngr *PluginMngr) End(sheetName string) {
	for _, p := range mngr.plugins {
		p.End(sheetName)
	}
}

func (mngr *PluginMngr) GenOneProduct(name string) {
	for _, p := range mngr.plugins {
		p.GenOneProduct(name)
	}
}

func (mngr *PluginMngr) ProductNotFound(name string) {
	for _, p := range mngr.plugins {
		p.ProductNotFound(name)
	}
}

func (mngr *PluginMngr) Gen(ex *excel.Excel) {
	for _, p := range mngr.plugins {
		p.Gen(ex)
	}
}
