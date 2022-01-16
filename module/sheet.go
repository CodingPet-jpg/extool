package module

import (
	"fmt"
	"github.com/xuri/excelize/v2"
)

const (
	Sheet1  = "Sheet1"
	Device  = "device"
	CtrlChe = "ctrls-che"
	CtrlWin = "ctrls-win"
	CtrlWeb = "ctrls-web"
)

func SetDefaultStyle(f *excelize.File, targetCell string, sheetName string) {
	style, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "0000FF", Style: 3},
			{Type: "top", Color: "00FF00", Style: 4},
			{Type: "bottom", Color: "FFFF00", Style: 5},
			{Type: "right", Color: "FF0000", Style: 6},
			{Type: "diagonalDown", Color: "A020F0", Style: 7},
			{Type: "diagonalUp", Color: "A020F0", Style: 8},
		},
	})
	if err != nil {
		fmt.Println(err)
	}
	err = f.SetCellStyle(sheetName, targetCell, targetCell, style)
}
