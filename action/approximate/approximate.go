package approximate

import (
	"extool/base"
	"extool/scanner"
)

func GetCaseParseHandler() scanner.Action {
	return func(ctx *scanner.Context) {
		if len(ctx.Row) < base.Cfg.Length {
			return
		}
		var entry = make([]string, 0, 4)
		for i, col := range ctx.Row {
			// append target column each row into string slice
			if base.Cfg.HitIndex(uint64(i)) {
				entry = append(entry, col)
			}
		}
		if _, ok := ctx.Scase.Contain(entry); !ok {
			ctx.Scase.PushBack(entry)
		}
	}
}
