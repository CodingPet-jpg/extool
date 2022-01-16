package action

import (
	"extool/module"
	"github.com/xuri/excelize/v2"
)

type Context struct {
	SheetName string
	File      *excelize.File
	Row       []string
	RowNum    int
	Scase     module.Case
	IsReaOnly bool
}

type Action func(context *Context)

type Listener func(context *Context)
