package main

import (
	"net/mail"
	"sort"
	"time"
)

type Threaded struct {
	mail.Message
	Parent   *Threaded
	Children []*Threaded

	parId    []string
}

func (tm *Threaded) addChild(c *Threaded) {
	tm.Children = append(tm.Children, c)
}

func (tm *Threaded) modified() time.Time {
	d := tm.Date
	for _, c := range tm.Children {
		if d.Before(c.Date) {
			d = c.Date
		}
	}
	return d
}

func (tm *Threaded) Root() *Threaded {
	if tm.Parent == nil {
		return tm
	}
	return tm.Parent.Root()
}

type sortable []*Threaded

func (s sortable) Len() int {
	return len(s)
}

func (s sortable) Less(i, j int) bool {
	return s[i].Date.Before(s[j].Date)
}

func (s sortable) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func Thread(msgs []*mail.Message) ([]*Threaded, []*Threaded) {
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

	all := []*Threaded{}

	threaded := []*Threaded{}
	for _, tm := range tmap {
		if tm.Parent == nil {
			threaded = append(threaded, tm)
		}
		all = append(all, tm)
		sort.Sort(sortable(tm.Children))
	}

	sort.Sort(sortable(threaded))
	return all, threaded
}
