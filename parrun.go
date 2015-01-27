package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
)

func main() {
	verb := flag.Bool("v", true, "Verbose")
	flag.Parse()
	runtime.GOMAXPROCS(8)
	cmd := flag.Arg(0)
	narg := flag.NArg() - 1
	if narg < 1 {
		return
	}
	res := make(chan Out, narg)
	for _, arg := range flag.Args()[1:] {
		arg := arg
		if *verb {
			log.Println("Starting", cmd, arg)
		}
		go func() {
			c := exec.Command(cmd, arg)
			bs, e := c.CombinedOutput()
			res <- Out{arg, bs, e}
		}()
	}
	var haveerr bool
	for i := 0; i < narg; i++ {
		o := <-res
		if o.err != nil {
			fmt.Fprintf(os.Stderr, "%s\n%s%s: %v\n%s\n", blline, erline, o.arg, o.err, blline)
			haveerr = true
		} else {
			fmt.Fprintf(os.Stderr, "%s\n%s%s\n%s\n", blline, okline, o.arg, blline)
		}
		os.Stdout.Write(o.bs)
	}
	fmt.Fprintf(os.Stderr, "\n%s\n", blline)
	if haveerr {
		log.Println("Exiting, there were errors")
		os.Exit(1)
	}
}

const (
	blline = `================================================================================`
	okline = `===================================   OK: `
	erline = `================================   ERROR: `
)

type Out struct {
	arg string
	bs  []byte
	err error
}
