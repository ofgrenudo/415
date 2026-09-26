package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/ofgrenudo/415/pkg/version"
)

func main() {
	var showVersion bool
	flag.BoolVar(&showVersion, "version", false, "Print version information and exit")
	flag.BoolVar(&showVersion, "v", false, "Shorthand for --version")
	flag.Parse()

	if showVersion {
		fmt.Println(version.String("example"))
		return
	}

	log.Println("Hello World! This is a Test!")
}
