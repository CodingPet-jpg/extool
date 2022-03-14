package main

import (
	"extool/action/formal"
	"extool/config"
	"extool/module"
	"extool/walker"
)

func main() {

	walker.NewWalker(walker.WRITECOPY).
		RegisterAction(module.Sheet1, formal.StepNameFixAction).
		GoWalkDir(config.Cfg.WorkDir)
}
