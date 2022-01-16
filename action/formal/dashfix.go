package formal

import (
	"extool/action"
	"extool/module"
	"github.com/xuri/excelize/v2"
	"log"
)

var SlashFixAction = func() action.Action {
	const (
		replaceChar   = "-"
		columnStepNoS = "J"
		columnStepNoI = 9
	)
	return func(ctx *action.Context) {
		if len(ctx.Row) > columnStepNoI && ctx.Row[columnStepNoI] == "" {
			ctx.Row[columnStepNoI] = replaceChar
		} else if len(ctx.Row) == columnStepNoI {
			ctx.Row = append(ctx.Row, replaceChar)
		} else {
			return
		}
		if !ctx.IsReaOnly {
			axis, _ := excelize.JoinCellName(columnStepNoS, ctx.RowNum)
			err := ctx.File.SetCellStr(ctx.SheetName, axis, replaceChar)
			module.SetDefaultStyle(ctx.File, axis, ctx.SheetName)
			if err != nil {
				log.Println(err)
			}
		}
	}
}
