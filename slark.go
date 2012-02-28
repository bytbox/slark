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
)

func main() {
	flag.Parse()

	if *version {
		fmt.Printf("slark %s\n", VERSION)
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
	}

	for _, msg := range msgs {
		fmt.Printf("%#v\n", msg)
	}
}
