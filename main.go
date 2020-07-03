package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/davidmz/debug-log"

	"github.com/davidmz/mustbe"
	flag "github.com/spf13/pflag"
)

var (
	log    = debug.NewLogger("anon", debug.WithOutput(os.Stderr))
	config *Config
	output *os.File
)

func main() {
	defer mustbe.Catched(func(err error) {
		_, _ = fmt.Fprintln(os.Stderr, "Fatal error:", err)
		os.Exit(1)
	})

	var (
		rulesFile string
		inFile    string
		outFile   string
	)
	flag.StringVarP(&rulesFile, "config", "c", "", "config file name (JSON)")
	flag.StringVarP(&inFile, "input", "i", "", "input file name (STDIN by default)")
	flag.StringVarP(&outFile, "output", "o", "", "output file name (STDOUT by default)")
	flag.Parse()

	if rulesFile == "" {
		_, _ = fmt.Fprintf(os.Stderr, "Flags of %s:\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
		os.Exit(0)
	}

	log.Println("Loading rules from", rulesFile)
	config = loadConfig(rulesFile)

	input := os.Stdin
	if inFile != "" {
		mustbeDone(func() {
			input = mustbe.OKVal(os.Open(inFile)).(*os.File)
		}, "cannot open input file %s: %w", inFile)
		defer input.Close()
	}

	output = os.Stdout
	if outFile != "" {
		mustbeDone(func() {
			output = mustbe.OKVal(os.Create(outFile)).(*os.File)
		}, "cannot create output file %s: %w", outFile)
		defer output.Close()
	}

	log.Println("Start reading input")

	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "COPY ") {
			processTable(line, scanner)
		} else {
			mustbe.OKVal(fmt.Fprintln(output, line))
		}
	}
	mustbe.OK(scanner.Err())
	mustbe.OK(output.Sync())
	log.Println("All done.")
}
