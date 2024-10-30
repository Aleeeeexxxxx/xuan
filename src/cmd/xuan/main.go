package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"time"

	"xuan/src/cmd"
)

//go:embed p1_20240920.xlsx
var sourceFile []byte
var targets []string

func input() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("请输入需要生成的产品型号（输入 EOF 结束）：")
	for scanner.Scan() {
		line := scanner.Text()
		if line == "EOF" {
			fmt.Println("")
			break
		}
		targets = append(targets, line)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("读取输入时出错:", err)
		return
	}

	fmt.Println("将生成以下产品报表:")
	for i, line := range targets {
		fmt.Printf("%d: %s\n", i+1, line)
	}
	fmt.Println("")
}

func main() {
	output := fmt.Sprintf("output-%d.xlsx", time.Now().Unix())

	outputFile, err := os.OpenFile(output, os.O_CREATE|os.O_WRONLY, 0644)
	cmd.PanicIfNotNil(err)

	defer outputFile.Close()

	input()

	cmd.RunGeneratorOrPanic(bytes.NewBuffer(sourceFile), outputFile, targets)

	fmt.Println("")
	fmt.Printf("文件保存为:%s,\r\n按回车键结束程序...", output)
	fmt.Scanln()
}
