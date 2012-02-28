package main

import (
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/bytbox/go-mail"
)

func mktmpl(tdn, tn string) *template.Template {
	tp := filepath.Join(tdn, tn+".tmpl.html")
	tsrc, err := ioutil.ReadFile(tp)
	if err != nil {
		panic(err)
	}
	tmpl, err := template.New(tn).Parse(string(tsrc))
	if err != nil {
		panic(err)
	}
	return tmpl
}

func writeHtml(odn, tdn string, msgs []mail.RawMessage) {
	os.MkdirAll(odn, os.ModeDir|0755)

	indexTmpl := mktmpl(tdn, "index")
	f, err := os.Create(filepath.Join(odn, "index.html"))
	if err != nil {
		panic(err)
	}
	indexTmpl.Execute(f, msgs)
	f.Close()

/*
	msgTmpl := mktmpl(tdn, "message")

	for _, msg := range msgs {
		msgTmpl.Execute(os.Stdout, msg)
	}
*/
}
