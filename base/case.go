package base

import (
	"bytes"
	"container/list"
	"encoding/json"
	"fmt"
	"sync"
)

type (
	Case struct {
		Name string
		*list.List
	}

	CaseChain struct {
		*list.List
	}
)

func NewCaseChain() CaseChain {
	return CaseChain{List: list.New()}
}

func NewCase(path string) Case {
	return Case{Name: path, List: list.New()}
}

func (chain CaseChain) EliAppend(c Case) {
	// iterate entry of case waiting append
	for entry := c.Front(); entry != nil; entry = entry.Next() {
		// iterate all case in current chain
		for tcase := chain.Front(); tcase != nil; tcase = tcase.Next() {
			// current case in chain which waiting to be processed
			current := tcase.Value.(Case)
			// if entry can be found in current case,then remove this entry
			if ent, ok := current.Contain(entry.Value.([]string)); ok {
				current.Remove(ent)
				// each time entry being removed,we determine if the case has zero entry
				if current.Len() == 0 {
					chain.Remove(tcase)
				}
				break
			}
		}
	}
	chain.PushBack(c)
}

// Deprecate:low performance
/*func (tcase Case) Contain(sentry []string) (target *list.Element, ok bool) {
	// iterate each entry in current case
	for tentry := tcase.Front(); tentry != nil; tentry = tentry.Next() {
		if strings.Join(sentry, "") == strings.Join(tentry.Value.([]string), "") {
			ok = true
			target = tentry
		}
	}
	return
}*/

// determine if the tentry which identical to sentry can be found in tcase,return tentry and true otherwise nil and false

func (tcase Case) Contain(sentry []string) (target *list.Element, ok bool) {
	// iterate each entry in current case
	for tentry := tcase.Front(); tentry != nil; {
		// compare the target entry and source entry
		// TODO:bugfix out of range
		for i, s := range sentry {
			if s == tentry.Value.([]string)[i] {
				continue
			}
			goto unfind
		}
		ok = true
		target = tentry
		return
	unfind:
		tentry = tentry.Next()
	}
	return
}

type (
	Item             []string
	CaseForSerialize struct {
		Name string
		Flow []Item
	}
)

func (tcase Case) Marshal(twc chan<- *bytes.Buffer, wg *sync.WaitGroup) {
	defer wg.Done()
	buf := new(bytes.Buffer)

	var fcase = CaseForSerialize{Name: tcase.Name, Flow: make([]Item, 0, 4)}
	for entry := tcase.Front(); entry != nil; entry = entry.Next() {
		var item = make([]string, 0, 4)
		item = entry.Value.([]string)
		fcase.Flow = append(fcase.Flow, item)
	}
	by, _ := json.Marshal(fcase)
	_, err := buf.Write(append(by, '\n'))
	if err != nil {
		fmt.Println(err)
	}

	twc <- buf
}

func UnMarshal(b []byte) (ucase Case) {
	var scase CaseForSerialize
	err := json.Unmarshal(b, &scase)
	if err != nil {
		fmt.Println(err)
	}
	ucase = NewCase(scase.Name)
	for _, item := range scase.Flow {
		ucase.PushBack([]string(item))
	}
	return
}
