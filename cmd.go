package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: crudgen [flags] [package directory] structName\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "[package directory] defaults to \".\" if unspecified\n")
		os.Exit(2)
	}

	output := flag.String("o", "", "where to write output; default is to write to stdout")
	flag.Parse()

	path := "."
	var structName string
	switch len(flag.Args()) {
	case 2:
		path = flag.Arg(0)
		structName = flag.Arg(1)
	case 1:
		structName = flag.Arg(0)
	default:
		flag.Usage()
	}

	data, err := Parse(path, structName)
	if err != nil {
		panic(err)
	}

	var outputFile *os.File

	if *output == "" {
		outputFile = os.Stdout
	} else {
		var err error
		outputFile, err = os.Create(*output)
		if err != nil {
			panic(err)
		}
	}

	defer outputFile.Close()
	err = tmpl.Execute(outputFile, data)
	if err != nil {
		panic(err)
	}
}
