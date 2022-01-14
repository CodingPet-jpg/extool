package scanner

import (
	"extool/base"
	"fmt"
	"github.com/xuri/excelize/v2"
	"log"
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
	mode      OpMode
	abs       string
	workbook  *excelize.File
	atos      map[string][]Action
	sheetName string // current Sheet
}

func (s *ActionScanner) Scan() base.Case {
	context := Context{File: s.workbook, Scase: base.NewCase(s.abs), IsRead: s.mode == READONLY}
	for sheetName, actions := range s.atos {
		if len(actions) == 0 {
			continue
		}
		context.RowNum = 0
		context.SheetName = sheetName
		rows, err := s.workbook.Rows(sheetName)
		if err != nil {
			log.Println(err)
		}
		for rows.Next() {
			context.RowNum++
			row, err := rows.Columns()
			context.Row = row
			if err != nil {
				fmt.Println(err)
			}
			for _, action := range actions {
				action(&context)
			}
		}
	}
	s.finish()
	return context.Scase
}

func New(filePath string, mode OpMode) *ActionScanner {
	file, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil
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
		err := s.workbook.SaveAs(filepath.Join(filepath.Dir(s.abs), "new", filepath.Base(s.abs)))
		if err != nil {
			log.Println(err)
		}
	}
}
