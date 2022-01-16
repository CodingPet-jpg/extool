package approximate

import (
	"extool/action"
	"extool/module"
	"sync"
)

// get case info form context to Parsed channel after one sheet Parsed
// which can be consumed by background comparing schedule

type CaseComparer struct {
	Parsed chan module.Case
	Chain  module.CaseChain
}

// a listener which can collect the case info from context
func (comparer *CaseComparer) feedCase() action.Listener {
	return func(context *action.Context) {
		comparer.Parsed <- context.Scase
	}
}

func NewCaseComparer() *CaseComparer {
	return &CaseComparer{
		Parsed: make(chan module.Case, 8),
		Chain:  module.NewCaseChain(),
	}
}

func (comparer *CaseComparer) GetCaseFeedListener() func() action.Listener {
	return comparer.feedCase
}

func (comparer *CaseComparer) StartCompareBackGround() func() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for pcase := range comparer.Parsed {
			comparer.Chain.EliAppend(pcase)
		}
		wg.Done()
	}()
	return func() {
		wg.Wait()
	}
}

func (comparer *CaseComparer) Close() {
	close(comparer.Parsed)
}
