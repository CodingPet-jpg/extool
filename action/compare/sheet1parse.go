package compare

import (
	"extool/action"
	"extool/config"
)

// collect case info context

var Sheet1ParseAction = func() action.Action {
	var header = 0
	return func(ctx *action.Context) {
		// issue:the Sheet1's header be parsed
		if len(ctx.Row) > 8 {
			if header < 2 {
				header++
				return
			}
			var entry = make([]string, 0, 8)
			for i, col := range ctx.Row {
				// append target column each row into string slice
				// issue: the row length < the minimal Bitmap,this situation will produce the entry which length 0
				if config.Cfg.HitIndex(uint64(i)) {
					entry = append(entry, col)
				}
			}
			// issue:prevent the entry that all column is blank
			for _, s := range entry {
				if s != "" {
					break
				}
			}
			// issue:the length of row > limit,but all selected column doesn't hit bitmap
			if len(entry) == 0 {
				return
			}
			if _, ok := ctx.Scase.Contain(entry); !ok {
				ctx.Scase.PushBack(entry)
			}
		}
	}
}
