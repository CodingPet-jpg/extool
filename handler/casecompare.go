package handler

import (
	"extool/base"
	"sync"
)

var parsed chan base.Case
var chain = base.NewCaseChain()

func StartCompare(w2 *sync.WaitGroup) {
	go func() {
		for pcase := range parsed {
			chain.EliAppend(pcase)
		}
		w2.Done()
	}()
}

func FeedCase(fcase base.Case) {
	parsed <- fcase
}
