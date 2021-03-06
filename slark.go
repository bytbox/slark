// slark is a suckless mail archiving program.
package main

import (
	"flag"
	"fmt"
	"os"
)

const VERSION = `0.0.1`

var (
	version = flag.Bool("V", false, "display version information")
	html    = flag.String("html", "html/", "write html output to directory")
	tmpldir = flag.String("templates", "tmpl/", "use template directory")
	statdir = flag.String("static", "static/", "use static directory")
)

func main() {
	flag.Parse()

	if *version {
		fmt.Printf("slark %s\n", VERSION)
		return
	}

	args := flag.Args()
	if len(args) != 1 {
		flag.Usage()
		return
	}

	fname := args[0]
	msgs, err := ReadMboxFile(fname)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading mbox: %s\n", err.Error())
		return
	}

	all, threaded := Thread(msgs)

	writeHtml(*html, *tmpldir, all, threaded)
	copyStatic(*html, *statdir)
}
