package formal

import (
	"extool/action"
	"github.com/xuri/excelize/v2"
	"log"
)

var StepNameFixAction = func() action.Action {

	textBeforeArray:= [2]string{"//CounterTestTool_Before","//CounterTestTool_After"}
	textAfterArray:= [2]string{"CounterTestTool_Before","CounterTestTool_After"}
	const firstColumn = "A"
	return func(ctx *action.Context) {
		for j,a:=range textBeforeArray{
			if len(ctx.Row)>0 && ctx.Row[0] == a && !ctx.IsReaOnly {
				axis, _ := excelize.JoinCellName(firstColumn, ctx.RowNum)
				err := ctx.File.SetCellStr(ctx.SheetName, axis, textAfterArray[j])
				if err != nil {
						log.Printf("%s : %v\n", ctx.File.Path, err)
						return
					}
				log.Printf("%s : %s\n", ctx.File.Path, "replace success")	
			}
		}
	}
}
