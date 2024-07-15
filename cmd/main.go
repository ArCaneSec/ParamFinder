package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ArCaneSec/paramfinder/internal/miner"
	"github.com/ArCaneSec/paramfinder/opts"
	"github.com/projectdiscovery/goflags"
)

func main() {
	opts := &opts.Options{}
	flag := goflags.NewFlagSet()
	flag.SetDescription("Extrcting parameters from different sources with different modes.")

	flag.CreateGroup("inputs", "Inputs",
		flag.StringVarP(&opts.Url, "url", "u", "", "url to extract parameters from"),
		flag.StringVarP(&opts.DirectoryPath, "directory", "d", "", "path to directory containing requests/responses files"),
		flag.StringSliceVarP(&opts.Headers, "headers", "H", nil, "comma separated custom headers for http request, each header must have a value", goflags.CommaSeparatedStringSliceOptions),
	)
	flag.CreateGroup("modes", "Modes",
		flag.BoolVarP(&opts.Silent, "silent", "s", false, "run in silent mode"),
		flag.BoolVarP(&opts.Crawl, "crawl", "c", false, "crawl js files to extract parameters from them as well."),
		flag.BoolVarP(&opts.Headless, "headless", "hs", false, "run with headless browser (usefull when hunting on SPAs)"),
	)
	flag.CreateGroup("output", "Output",
		flag.StringVarP(&opts.OutputPath, "output", "o", "", "path to output file"),
	)

	if err := flag.Parse(); err != nil {
		log.Fatalf("error parsing flags: %v\n", err)
	}
	if err := opts.Validate(); err != nil {
		log.Fatalf("error validating flags: %v\n", err)
	}

	m := &miner.Miner{opts}
	params, err := m.Mine()
	if err != nil {
		m.Flog(err)
	}

	if opts.OutputPath != "" {
		file, err := os.OpenFile(opts.OutputPath, os.O_CREATE|os.O_WRONLY, 0644)

		if err != nil {
			m.Flog(err)
		}

		for _, param := range params {
			file.WriteString(fmt.Sprintf("%s\n", param))
		}
		return
	}
	for _, param := range params {
		fmt.Println(param)
	}
}
