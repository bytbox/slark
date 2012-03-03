package main

import (
	"sort"
	"time"

	. "github.com/bytbox/go-mail"
)

type Threaded struct {
	Message
	Parent   *Threaded
	Children []*Threaded

	parId    []string
}

func (tm *Threaded) addChild(c *Threaded) {
	tm.Children = append(tm.Children, c)
}

func (tm *Threaded) modified() time.Time {
	return tm.Date
}

type Sortable []*Threaded

func (s Sortable) Len() int {
	return len(s)
}

func (s Sortable) Less(i, j int) bool {
	return s[i].Date.Before(s[j].Date)
}

func (s Sortable) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func Thread(msgs []Message) []*Threaded {
	tmap := map[string]*Threaded{}
	for _, msg := range msgs {
		mid := msg.MessageId
		tm := &Threaded{msg, nil, nil, nil}
		// we just use the first parent
		if len(msg.References) > 0 {
			tm.parId = msg.References
		}
		tmap[mid] = tm
	}

	for _, tm := range tmap {
		if tm.parId != nil {
			// The parent should be set using the last id we have
			// registered.
			for _, id := range tm.parId {
				p, ok := tmap[id]
				if ok {
					tm.Parent = p
				}
			}
			if tm.Parent != nil {
				tm.Parent.addChild(tm)
			}
		}
	}

	threaded := []*Threaded{}
	for _, tm := range tmap {
		if tm.Parent == nil {
			threaded = append(threaded, tm)
		}
	}

	sort.Sort(Sortable(threaded))
	return threaded
}
