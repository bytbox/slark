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

func ReadMbox(r io.Reader) (msgs []mail.RawMessage, err error) {
	var mbuf *bytes.Buffer
	var m mail.RawMessage
	br := bufio.NewReaderSize(r, _MAX_LINE_LEN)
	l, _, err := br.ReadLine()
	for err == nil {
		fs := bytes.SplitN(l, []byte{' '}, 3)
		if len(fs) == 3 && string(fs[0]) == "From" {
			// flush the previous message, if necessary
			if mbuf != nil {
				m, err = mail.ParseRaw(mbuf.Bytes())
				if err != nil { return }
				msgs = append(msgs, m)
			} else {
				mbuf = new(bytes.Buffer)
			}
		} else {
			_, err = mbuf.Write(l)
			if err != nil { return }
			_, err = mbuf.Write(crlf)
			if err != nil { return }
		}
		l, _, err = br.ReadLine()
	}
	if err == io.EOF { err = nil }
	return
}

func ReadMboxFile(filename string) ([]mail.RawMessage, error) {
	f, err := os.Open(filename)
	if err != nil { return nil, err }
	msgs, err := ReadMbox(f)
	f.Close()
	return msgs, err
}
