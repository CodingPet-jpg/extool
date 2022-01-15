package main

import (
	"extool/action/approximate"
	"extool/action/formal"
	"extool/base"
	"extool/handler"
	"extool/scanner"
	"fmt"
	"time"
)

func main() {
	var start = time.Now()
	var workDir = "C:\\Users\\JOKOI\\Desktop\\1"
	var actionMap = map[string][]handler.Producer{
		base.Sheet1: {
			formal.GetSlashFixHandler,
			formal.GetStepNoHandler,
			approximate.GetCaseParseHandler,
		},
	}

	for i := 0; i < 1; i++ {
		walker := handler.NewWalker(actionMap, scanner.WRITECOPY, false)
		walker.ParWalk(workDir)
		fmt.Println(time.Since(start), "\t", walker.FileCounter)
	}
}
