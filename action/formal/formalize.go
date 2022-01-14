package formal

import (
	"github.com/CodingPet-jpg/extool/base"
	"github.com/xuri/excelize/v2"
	"strconv"
)

// 返回一个处理器可以检测连续Step No，并纠正非连续Step No
// 纠正后的数据需要同时反应在row和file里，既可以修改原文件同时还对后续的Action可见

func GetStepNoHandler() func(file *excelize.File, row []string, rowNum *int) {
	var expectStepNo = 1
	const (
		columnStepNoS = "C"
		columnStepNoI = 2
	)
	var R = func(file *excelize.File, row []string, rowNum *int) {
		if len(row) <= columnStepNoI+1 {
			expectStepNo = 1
		} else {
			i, _ := strconv.Atoi(row[columnStepNoI])
			if i == expectStepNo {
				expectStepNo++
			} else {
				axis, _ := excelize.JoinCellName(columnStepNoS, *rowNum)
				row[columnStepNoI] = strconv.Itoa(expectStepNo)
				err := file.SetCellInt(base.Sheet1, axis, expectStepNo)
				if err == nil {
					expectStepNo++
				}
			}
		}
	}
	return R
}
