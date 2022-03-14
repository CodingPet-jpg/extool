package formal

import (
	"extool/action"
	"github.com/xuri/excelize/v2"
	"log"
	"strings"
)

var StepNameFixAction = func() action.Action {
	const model = "CounterTestTool"
	const firstColumn = "A"
	return func(ctx *action.Context) {
		if len(ctx.Row)>5 && ctx.Row[5] == model {
			axisonebefore, _ := excelize.JoinCellName(firstColumn, ctx.RowNum-1)
			if str,_:=ctx.File.GetCellValue(ctx.SheetName, axisonebefore, excelize.Options{});str!="" && strings.HasPrefix(str, "//"){
				newstr:=strings.TrimPrefix(str, "//")
				err := ctx.File.SetCellStr(ctx.SheetName, axisonebefore, newstr)
				if err != nil {
					log.Printf("%s : %v\n", ctx.File.Path, err)
					return
				}
				log.Printf("%s : %s\n", ctx.File.Path, "replace success")
			}	
		}
	}
}
