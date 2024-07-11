package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"sync"
)

type miner struct {
	url           string
	crawlMode     bool
	headless      bool
	silent        bool
	directoryPath string
}

func (m *miner) mine() ([]string, error) {
	var (
		contents string
		err      error
	)

	switch {
	case m.directoryPath != "":
		files, err := os.ReadDir(m.directoryPath)
		if err != nil {
			return nil, fmt.Errorf("[!] Error while listing files in provided directory path: %w", err)
		}

		params := make(chan string, 10)
		wg := &sync.WaitGroup{}
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			wg.Add(1)
			go func() {
				defer wg.Done()
				bytes, err := os.ReadFile(fmt.Sprintf("%s/%s", m.directoryPath, file.Name()))
				if err != nil {
					m.fLog(err)
				}

				fileParams := findAllParams(string(bytes))
				for param := range fileParams {
					params <- param
				}
			}()
		}
		go func() {
			wg.Wait()
			close(params)
		}()

		uniques, _ := uniqueParams(params)
		return uniques, nil

	case m.crawlMode:
		contents, err = m.runCrawlMode()
	default:
		contents, err = m.runRawMode()
	}

	if err != nil {
		return nil, err
	}
	params := findAllParams(contents)
	uniques, _ := uniqueParams(params)

	return uniques, nil
}

func (m *miner) runCrawlMode() (string, error) {
	htmlContent, err := m.sendRequest()
	if err != nil {
		return "", err
	}

	jsPaths := extractJsPath(htmlContent)
	urls := make([]string, 0, len(jsPaths))

	urlObj, _ := url.Parse(m.url)
	for _, path := range jsPaths {
		urlObj.Path = path
		urls = append(urls, urlObj.String())
	}

	responses := make(chan string, len(urls))
	wg := &sync.WaitGroup{}

	for _, url := range urls {
		wg.Add(1)
		go func() {
			defer wg.Done()

			res, err := rawRequest(url)

			if err != nil {
				m.log(fmt.Errorf("[!] Error while sending request to %s: %w", url, err))
				return
			}
			responses <- res
		}()
	}
	go func() {
		wg.Wait()
		close(responses)
	}()

	var packedBodies []string
	packedBodies = append(packedBodies, htmlContent)

	for res := range responses {
		packedBodies = append(packedBodies, res)
	}

	return strings.Join(packedBodies, "\n"), nil
}

func (m *miner) runRawMode() (string, error) {
	content, err := m.sendRequest()
	if err != nil {
		return "", err
	}

	return content, nil
}

func (m *miner) sendRequest() (string, error) {
	if m.headless {
		return headlessRequest(m.url)
	}
	return rawRequest(m.url)
}

func findAllParams(html string) <-chan string {
	params := make(chan string, 10)
	var wg sync.WaitGroup

	for _, pattern := range patterns {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for _, param := range pattern.FindAllStringSubmatch(html, -1) {
				params <- param[1]
			}
		}()
	}
	go func() {
		wg.Wait()

		close(params)
	}()

	return params
}

func (m *miner) log(v ...any) {
	if !m.silent {
		log.Println(v...)
	}
}

// fatal log
func (m *miner) fLog(v ...any) {
	if !m.silent {
		log.Fatalln(v...)
	}
	os.Exit(1)
}
