package main

import (
	"fmt"
	"log"
	"strings"
	"sync"
)

type miner struct {
	url       string
	crawlMode bool
	headless  bool
}

func (m *miner) mine() {
	var (
		contents string
		err error
	)
	if m.crawlMode {
		contents, err = m.runCrawlMode()
	} else {
		contents, err = m.runRawMode()
	}

	if err != nil {
		log.Fatal(err)
	}
	params, _ := findAllParams(contents)

	// loggingDone, err := logError(errors)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	mergedParams := merge(params)

	uniques, _ := uniqueParams(mergedParams)
	// <-loggingDone
	fmt.Println(uniques)
}

func (m *miner) runCrawlMode() (string, error) {
	log.Println("[*] Crawling mode...")

	htmlContent, _ := m.sendRequest()
	jsPaths := extractJsPath(htmlContent)
	urls := make([]string, 0, len(jsPaths))

	for _, path := range jsPaths {
		urls = append(urls, fmt.Sprintf("%s%s", m.url, path))
	}

	responses := make(chan string, len(urls))
	errors := make(chan error, len(urls))
	wg := &sync.WaitGroup{}

	for _, url := range urls {
		wg.Add(1)
		go func() {
			defer wg.Done()

			res, err := rawRequest(url)

			if err != nil {
				errors <- err
				return
			}
			responses <- res
		}()
	}
	go func() {
		wg.Wait()
		close(responses)
		close(errors)
	}()

	var packedBodies []string
	packedBodies = append(packedBodies, htmlContent)

	for res := range responses {
		packedBodies = append(packedBodies, res)
	}

	return strings.Join(packedBodies, ","), nil
}

func (m *miner) runRawMode() (string, error) {

	return "", nil
}

func (m *miner) sendRequest() (string, error) {
	if m.headless {
		return headlessRequest(m.url)
	}
	return rawRequest(m.url)
}
