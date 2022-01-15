package scanner

import (
	"extool/base"
	"fmt"
	"github.com/xuri/excelize/v2"
	"log"
	"os"
	"path/filepath"
)

// 三种模式，ReadOnly在文件扫描结束后无操作，WRITESOURCE在文件扫描结束后保存原文件，WRITECOPY在文件扫描结束后在当前文件所在目录创建new，并将文件拷贝写道该目录

type OpMode int

const (
	READONLY OpMode = iota
	WRITESOURCE
	WRITECOPY
)

// 对传入的行进行指定行为，包括读取和写入

type Action func(context *Context)

// 封装表格文件和对其采取的行为

type Context struct {
	SheetName string
	File      *excelize.File
	Row       []string
	RowNum    int
	Scase     base.Case
	IsRead    bool
}

type ActionScanner struct {
	mode     OpMode
	abs      string
	workbook *excelize.File
	atos     map[string][]Action
}

func (s *ActionScanner) Scan(resultHandler func(p base.Case)) {
	ctx := Context{File: s.workbook, Scase: base.NewCase(s.abs), IsRead: s.mode == READONLY}
	for sheetName, actions := range s.atos {
		if len(actions) == 0 {
			continue
		}
		ctx.RowNum = 0
		ctx.SheetName = sheetName
		rows, err := s.workbook.Rows(sheetName)
		if err != nil {
			log.Println(err)
		}
		for rows.Next() {
			ctx.RowNum++
			row, err := rows.Columns()
			ctx.Row = row
			if err != nil {
				fmt.Println(err)
			}
			for _, action := range actions {
				action(&ctx)
			}
		}
	}
	if resultHandler != nil {
		resultHandler(ctx.Scase)
	}
	s.finish()
}

func New(filePath string, mode OpMode) *ActionScanner {
	file, err := excelize.OpenFile(filePath)
	if err != nil {
		fmt.Println(err)
	}
	return &ActionScanner{workbook: file, mode: mode, abs: filePath, atos: map[string][]Action{}}
}

func (s *ActionScanner) RegisterAction(sheet string, action Action) {
	s.atos[sheet] = append(s.atos[sheet], action)
}

func (s *ActionScanner) finish() {
	switch s.mode {
	case READONLY:
		return
	case WRITESOURCE:
		err := s.workbook.Save()
		if err != nil {
			log.Println(err)
		}
	case WRITECOPY:
		np := filepath.Join(filepath.Dir(filepath.Dir(s.abs)), "copy")
		errm := os.Mkdir(np, 666)
		if errm != nil {
			log.Println(errm)
		}
		errs := s.workbook.SaveAs(filepath.Join(np, filepath.Base(s.abs)))
		if errs != nil {
			log.Println(errs)
		}
	}
}

func (s *ActionScanner) Valid() bool {
	return s.workbook != nil
}
