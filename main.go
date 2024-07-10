package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sync"
)

var contentTypeToPattern = make(map[string][]*regexp.Regexp)

func extractParams(pattern *regexp.Regexp, html string, paramCh chan []string) {
	var params []string
	for _, match := range pattern.FindAllStringSubmatch(html, -1) {
		params = append(params, match[1])
	}

	paramCh <- params
}

func findAllParams(html string) (<-chan []string, <-chan error) {
	params := make(chan []string, len(patterns))
	errors := make(chan error, len(patterns))

	var wg sync.WaitGroup

	for _, pattern := range patterns {
		wg.Add(1)
		go func() {
			defer wg.Done()
			extractParams(pattern, html, params)
		}()
	}
	go func() {
		wg.Wait()

		close(params)
		close(errors)
	}()
	return params, errors
}

func merge(params <-chan []string) <-chan string {
	mergedParams := make(chan string)
	var wg sync.WaitGroup

	for paramSlice := range params {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for _, val := range paramSlice {
				mergedParams <- val
			}
		}()
	}
	go func() {
		wg.Wait()
		close(mergedParams)
	}()

	return mergedParams
}

func uniqueParams(params <-chan string) ([]string, map[string]int) {
	seen := make(map[string]int, len(params))
	var uniqueParams []string

	for val := range params {
		if _, ok := seen[val]; !ok {
			seen[val] = 1
			uniqueParams = append(uniqueParams, val)
		} else {
			seen[val] += 1
		}
	}

	return uniqueParams, seen
}

func logError(errors <-chan error) (<-chan struct{}, error) {
	ex, _ := os.Executable()
	dir, _ := filepath.Split(filepath.Dir(ex))
	finalPath := filepath.Join(dir, "logs", "fallparams-logs.txt")

	file, err := os.OpenFile(finalPath, os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("[!] An error occurred while trying to open or create log file, error: %w", err)
	}

	done := make(chan struct{})
	go func() {
		for err := range errors {
			file.WriteString(err.Error() + "\n")
		}
		close(done)
	}()
	return done, nil
}

func main() {
	m := new(miner)

	flag.StringVar(&m.url, "u", "", "url to extract parameters from.")
	flag.BoolVar(&m.crawlMode, "c", false, "crawl js files to extract parameters from them as well.")
	flag.BoolVar(&m.headless, "h", false, "run with headless browser (usefull when hunting on SPAs).")

	flag.Parse()
	if m.url == "" {
		log.Fatalln("[!] you must provide a valid url")
	}

	m.mine()
}
