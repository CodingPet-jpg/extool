package approximate

import (
	"extool/base"
	"extool/scanner"
)

func GetCaseParseHandler() scanner.Action {
	return func(ctx *scanner.Context) {
		if len(ctx.Row) < 4 {
			return
		}
		var entry = make([]string, 0, 8)
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
