package walker

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestRel(t *testing.T) {
	var path = "C:\\Users\\JOKOI\\Desktop\\1"
	var abs = "C:\\Users\\JOKOI"
	s, _ := filepath.Rel(abs, path)
	fmt.Println(s)
}

func TestMakeDir(t *testing.T) {
	err := os.MkdirAll("C:\\Users\\JOKOI\\Desktop\\1\\2\\3\\4\\5\\1.txt", os.ModeDir)
	if err != nil {
		t.Log(err)
	}
}
