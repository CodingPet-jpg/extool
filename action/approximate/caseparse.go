package approximate

import (
	"extool/action"
	"extool/config"
)

// collect case info context

var CaseParseAction = func() action.Action {
	return func(ctx *action.Context) {
		if len(ctx.Row) < 4 {
			return
		}
		var entry = make([]string, 0, 8)
		for i, col := range ctx.Row {
			// append target column each row into string slice
			if config.Cfg.HitIndex(uint64(i)) { //TODO:use config
				entry = append(entry, col)
			}
		}
		if _, ok := ctx.Scase.Contain(entry); !ok {
			ctx.Scase.PushBack(entry)
		}
	}
}
