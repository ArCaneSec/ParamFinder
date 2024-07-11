package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	var (
		m          = &miner{}
		outputPath string
	)

	flag.StringVar(&m.url, "u", "", "url to extract parameters from.")
	flag.StringVar(&outputPath, "o", "", "output file.")
	flag.StringVar(&m.directoryPath, "d", "", "path to directory containing requests/responses files.")
	flag.BoolVar(&m.crawlMode, "c", false, "crawl js files to extract parameters from them as well.")
	flag.BoolVar(&m.headless, "h", false, "run with headless browser (usefull when hunting on SPAs).")
	flag.BoolVar(&m.silent, "s", false, "run in silent mode (only output parameters).")

	flag.Parse()
	if m.url == "" && m.directoryPath == "" {
		m.fLog("[!] you must provide a valid url or directory containing requests/respones files")
	}

	params, err := m.mine()
	if err != nil {
		m.fLog(err)
	}

	if outputPath != "" {
		file, err := os.OpenFile(outputPath, os.O_CREATE|os.O_WRONLY, 0644)

		if err != nil {
			m.fLog(err)
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
