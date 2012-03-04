package main

import (
	"bufio"
	"bytes"
	"io"
	"os"

	"github.com/bytbox/go-mail"
)

const _MAX_LINE_LEN = 1024

var crlf = []byte{'\r', '\n'}

func ReadMbox(r io.Reader) (msgs []mail.Message, err error) {
	var mbuf *bytes.Buffer
	var m mail.Message
	lastblank := true
	br := bufio.NewReaderSize(r, _MAX_LINE_LEN)
	l, _, err := br.ReadLine()
	for err == nil {
		fs := bytes.SplitN(l, []byte{' '}, 3)
		if len(fs) == 3 && string(fs[0]) == "From" && lastblank {
			// flush the previous message, if necessary
			if mbuf != nil {
				m, err = mail.Parse(mbuf.Bytes())
				if err != nil {
					println(err.Error())
				} else {
					msgs = append(msgs, m)
				}
			}
			mbuf = new(bytes.Buffer)
		} else {
			_, err = mbuf.Write(l)
			if err != nil {
				return
			}
			_, err = mbuf.Write(crlf)
			if err != nil {
				return
			}
		}
		if len(l) > 0 {
			lastblank = false
		} else {
			lastblank = true
		}
		l, _, err = br.ReadLine()
	}
	if err == io.EOF {
		m, err = mail.Parse(mbuf.Bytes())
		if err != nil {
			return
		}
		msgs = append(msgs, m)
	}
	return
}

func ReadMboxFile(filename string) ([]mail.Message, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	msgs, err := ReadMbox(f)
	f.Close()
	return msgs, err
}
