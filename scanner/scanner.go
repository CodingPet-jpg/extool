package scanner

import (
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

type Action func(file *excelize.File, rows []string, rowNum *int) int

// 封装表格文件和对其采取的行为

type ActionScanner struct {
	mode      OpMode
	abs       string
	workbook  *excelize.File
	actions   []Action
	sheetName string // current Sheet
}

func (s *ActionScanner) Scan() {
	rows, err := s.workbook.Rows(s.sheetName)
	if err != nil {
		log.Println(err)
	}
	var rowNum = 0
	for rows.Next() {
		rowNum++
		row, err := rows.Columns()
		if err != nil {
			fmt.Println(err)
		}
		for i, action := range s.actions {
			log.Println(i)
			action(s.workbook, row, &rowNum)
		}
	}
	s.finish()
}

func New(sheet string, filePath string, mode OpMode) *ActionScanner {
	file, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil
	}
	return &ActionScanner{workbook: file, sheetName: sheet, mode: mode, abs: filePath}
}

func (s *ActionScanner) RegisterAction(action Action) {
	s.actions = append(s.actions, action)
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
