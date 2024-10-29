package main

import (
	"log"
	"os"

	"xuan/src"
	"xuan/src/cmd"
	"xuan/src/generator"
	"xuan/src/parser"
)

func main() {
	const sourceFilePath = "p1_20240920.xlsx"
	const output = "output.xlsx"

	targets := []string{
		"HWD16T245-B",
		"HWD16T245-C",
		"HWD16T245MCSOP48",
		"HWD16T245ESOP48",
		"HWD16T245-A",
		"HWD16T245FBGA56I",
		"HWD16T245PBGA56M3",
		"HWD16T245PSOP48M3",
		"HWD2210EFBGA256",
		"HWD2210FBGA256I",
		"HWD2210MCFBGA256",
		"HWD2210MCFBGA324",
		"HWD2210MFBGA256",
		"HWD25Q256",
		"HWD25Q256PSOP8I",
		"HWD2V1000-4FG256",
		"HWD2V1000-4FG256N",
		"HWD3490MAA",
		"HWD3490EAA",
		"HWD3490PSOP8M3",
		"HWD4VLX25-10FF668",
		"HWD4VLX25-10FF668N",
		"HWD70345MAG",
		"HWD70345TSSOP24E",
		"HWD7606LQFP64M3",
		"HWD7606CQFP64M1",
		"HWD7606-ALQFP64E",
		"HWD767D301-A",
		"HWD767D301MAG",
		"HWD3490MAA",
		"HWD3490EAA",
		"HWD3490PSOP8M3",
		"HWD3490PSOP8E",
		"HWD32F103ELQFP64",
		"HWD3232-A",
		"HWD3232EESE",
		"HWD3232-B",
		"HWD3232",
		"HWD70351MAG",
		"HWD164245-A",
		"HWD164245EAM",
		"HWD164245MAM",
		"HWD33S2561606B",
		"HWD7490-A",
		"HWD32F407PBGA176I",
		"HWD32F407PBGA176M3",
		"HWD4644PBGA77E",
		"HWD7656CQFP64APM1",
		"HWD164245-A",
		"HWD7656LQFP64M3",
		"HWD708",
		"HWD708R",
		"HWD708S",
		"HWD708T",
		"HWD3485",
		"HWD3485PSOP8E",
		"HWD3485PSOP8M3",
		"HWD18V04MCQFP44",
		"HWD18V04ETQFP44",
		"HWD74401EQFN20",
		"HWD74401MCLCC20",
		"HWD74401MQFN20N",
		"HWD1668",
	}

	datastore := src.NewInMemoryDatastore()

	parser, err := parser.NewExcelFileParser(sourceFilePath, datastore)
	cmd.PanicIfNotNil(err)

	cmd.PanicIfNotNil(parser.Parse())

	gen := generator.NewExcelGenerator(datastore, targets)
	gen.AddTable(generator.NewBasicTableGenerator)
	gen.AddTable(generator.NewWKBTableGenerator)

	excel, err := gen.Gen()
	cmd.PanicIfNotNil(err)

	file, err := os.OpenFile(output, os.O_CREATE|os.O_WRONLY, 0644)
	cmd.PanicIfNotNil(err)
	defer file.Close()

	cmd.PanicIfNotNil(excel.Write(file))

	log.Println("Excel file created successfully!")
}
