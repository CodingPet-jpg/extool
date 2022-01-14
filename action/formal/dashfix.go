package formal

import (
	"extool/base"
	"extool/scanner"
	"github.com/xuri/excelize/v2"
	"log"
)

func GetSlashFixHandler() scanner.Action {
	const (
		replaceChar   = "-"
		columnStepNoS = "J"
		columnStepNoI = 9
	)
	return func(ctx *scanner.Context) {
		if len(ctx.Row) > columnStepNoI && ctx.Row[columnStepNoI] == "" {
			ctx.Row[columnStepNoI] = replaceChar
		} else if len(ctx.Row) == columnStepNoI {
			ctx.Row = append(ctx.Row, replaceChar)
		} else {
			return
		}
		if !ctx.IsRead {
			axis, _ := excelize.JoinCellName(columnStepNoS, ctx.RowNum)
			err := ctx.File.SetCellStr(ctx.SheetName, axis, replaceChar)
			base.SetDefaultStyle(ctx.File, axis, ctx.SheetName)
			if err != nil {
				log.Println(err)
			}
		}
	}
}
