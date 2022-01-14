package main

import (
	"extool/action/approximate"
	"extool/action/formal"
	"extool/base"
	"extool/scanner"
	"fmt"
)

func main() {
	var filePath = "C:\\Users\\JOKOI\\Desktop\\1.xlsx"
	asc := scanner.New(filePath, scanner.WRITESOURCE)
	asc.RegisterAction(base.Sheet1, formal.GetStepNoHandler())
	asc.RegisterAction(base.Sheet1, formal.GetSlashFixHandler())
	asc.RegisterAction(base.Sheet1, approximate.GetCaseParseHandler())
	pcase := asc.Scan()
	fmt.Println(pcase)
}
