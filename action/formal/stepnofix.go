package formal

import (
	"extool/action"
	"extool/module"
	"github.com/xuri/excelize/v2"
	"log"
	"strconv"
)

// 返回一个处理器可以检测连续Step No，并纠正非连续Step No
// 纠正后的数据需要同时反应在row和file里，既可以修改原文件同时还对后续的Action可见

var StepNoFixAction = func() action.Action {
	var expectStepNo = 1
	const (
		columnStepNoS = "C"
		columnStepNoI = 2
	)
	return func(ctx *action.Context) {
		if len(ctx.Row) <= columnStepNoI {
			expectStepNo = 1
		} else {
			i, _ := strconv.Atoi(ctx.Row[columnStepNoI])
			if i != expectStepNo {
				ctx.Row[columnStepNoI] = strconv.Itoa(expectStepNo)
				if !ctx.IsReaOnly {
					axis, _ := excelize.JoinCellName(columnStepNoS, ctx.RowNum)
					err := ctx.File.SetCellInt(ctx.SheetName, axis, expectStepNo)
					module.SetDefaultStyle(ctx.File, axis, ctx.SheetName)
					if err != nil {
						log.Println(err)
					}
				}
			}
			expectStepNo++
		}
	}
}
