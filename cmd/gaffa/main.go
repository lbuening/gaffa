package main

import (
	"flag"
	"fmt"
	"github.com/lbuening/gaffa/internal/generate"
	"os"
)

const usage = `
Usage: gaffa [command] [flags]

Commands:
  generate // gaffa code generator
  version  // show gaffa version
`

func main() {
	flag.Usage = func() { _, _ = fmt.Fprint(os.Stderr, usage) }
	flag.Parse()
	switch flag.Arg(0) {
	case "generate":
		flag.Usage = func() { _, _ = fmt.Fprint(os.Stderr, generate.Usage) }
		generateFlags := flag.NewFlagSet("generate", flag.ExitOnError)
		generateFlags.Usage = func() {
			_, _ = fmt.Fprint(os.Stderr, generate.Usage)
		}
		_ = generateFlags.Parse(flag.Args()[1:])
		err := generate.Generate(".", flag.Args()[1:])
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	case "version":
		fmt.Println("gaffa version 0.0.1")
		os.Exit(0)
	default:
		_, _ = fmt.Fprint(os.Stderr, usage)
		os.Exit(1)
	}
}
