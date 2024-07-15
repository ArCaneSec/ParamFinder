package miner

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"sync"
	"github.com/ArCaneSec/paramfinder/opts"
	"github.com/ArCaneSec/paramfinder/internal/pattern"
)

type Miner struct {
	*opts.Options
}

func (m *Miner) Mine() ([]string, error) {
	var (
		contents string
		err      error
	)

	switch {
	case m.DirectoryPath != "":
		files, err := os.ReadDir(m.DirectoryPath)
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
				bytes, err := os.ReadFile(fmt.Sprintf("%s/%s", m.DirectoryPath, file.Name()))
				if err != nil {
					m.Flog(err)
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

	case m.Crawl:
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

func (m *Miner) runCrawlMode() (string, error) {
	htmlContent, err := m.sendRequest()
	if err != nil {
		return "", err
	}

	jsPaths := extractJsPath(htmlContent)
	urls := make([]string, 0, len(jsPaths))

	urlObj, _ := url.Parse(m.Url)
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

			res, err := rawRequest(url, m.Headers)

			if err != nil {
				m.Log(fmt.Errorf("[!] Error while sending request to %s: %w", url, err))
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

func (m *Miner) runRawMode() (string, error) {
	content, err := m.sendRequest()
	if err != nil {
		return "", err
	}

	return content, nil
}

func (m *Miner) sendRequest() (string, error) {
	if m.Headless {
		return headlessRequest(m.Url, m.Headers)
	}
	return rawRequest(m.Url, m.Headers)
}

func findAllParams(html string) <-chan string {
	params := make(chan string, 10)
	var wg sync.WaitGroup

	for _, pattern := range pattern.Patterns {
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

func (m *Miner) Log(v ...any) {
	if !m.Silent {
		log.Println(v...)
	}
}

// fatal Log
func (m *Miner) Flog(v ...any) {
	if !m.Silent {
		log.Fatalln(v...)
	}
	os.Exit(1)
}
