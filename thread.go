package main

import (
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
		threaded = append(threaded, tm)
	}
	return threaded
}
