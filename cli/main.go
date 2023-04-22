package main

import (
	"flag"
	"fmt"
	"os"
)

const (
	symbolMessage = "List of comma seperated crypto symbols."
)

var symbols string

var CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

var Usage = func() {
	fmt.Fprintf(CommandLine.Output(), "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
}

func init() {
	fmt.Println("Init ###")
	flag.StringVar(&symbols, "s", "BTC, ETH", symbolMessage)
}

func main() {
	flag.Parse()

	fmt.Println(symbols)
}
