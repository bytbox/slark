package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
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

func writeHtml(odn, tdn string, all []*Threaded, msgs []*Threaded) {
	os.MkdirAll(odn, os.ModeDir|0755)

	indexTmpl := mktmpl(tdn, "index")
	f, err := os.Create(filepath.Join(odn, "index.html"))
	if err != nil {
		panic(err)
	}
	err = indexTmpl.Execute(f, msgs)
	f.Close()
	if err != nil {
		panic(err)
	}


	msgTmpl := mktmpl(tdn, "message")

	for _, msg := range all {
		f, err := os.Create(filepath.Join(odn, fmt.Sprintf("%s.html", msg.Id)))
		if err != nil {
			panic(err)
		}
		msgTmpl.Execute(f, msg)
		f.Close()
	}
}

func copyStatic(odn, sdn string) {
	fs, e := ioutil.ReadDir(sdn)
	if e != nil { panic(e) }
	for _, i := range fs {
		if i.IsDir() { continue }
		in := filepath.Join(sdn, i.Name())
		out := filepath.Join(odn, i.Name())
		c, e := ioutil.ReadFile(in)
		if e != nil { continue }
		ioutil.WriteFile(out, c, 0755)
	}
}
