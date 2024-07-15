package miner

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/ArCaneSec/paramfinder/internal/pattern"
	"github.com/go-rod/rod"
)

func extractJsPath(html string) []string {
	matches := pattern.JsFiles.FindAllStringSubmatch(html, -1)

	allPaths := make([]string, 0, len(matches))

	for _, path := range matches {

		// removing useless js files
		if strings.Contains(path[1], "cdn") || strings.Contains(path[1], "jquery") {
			continue
		}

		allPaths = append(allPaths, path[1])
	}

	return allPaths
}

func rawRequest(url string, headers []string) (string, error) {
	client := http.Client{Timeout: 10 * time.Second}
	redirectCount := 0
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		if len(via) >= 3 {
			return http.ErrUseLastResponse
		}
		redirectCount++
		return nil
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("[!] Failed to construct the request object, url: %s, err: %w", url, err)
	}

	req.Header = http.Header{
		"User-Agent":      {"Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/114.0"},
		"Accept":          {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8"},
		"Accept-Language": {"en-US,en;q=0.5"},
		"Sec-Fetch-Dest":  {"document"},
		"Sec-Fetch-Mode":  {"navigate"},
		"Sec-Fetch-Site":  {"none"},
		"Sec-Fetch-User":  {"?1"},
		"Referer":         {url},
	}
	for i := 0; i < len(headers); i += 2 {
		req.Header.Set(headers[i], headers[i+1])
	}

	res, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("[!] Request to %s url failed with error: %w", url, err)
	}

	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("[!] Failed to read response for %s url, error: %s\n", url, err.Error())
	}

	return string(resBody), nil
}

func headlessRequest(url string, headers []string) (string, error) {
	page := rod.New().MustConnect().MustPage()
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
	page = page.Context(ctx)
	_, err := page.SetExtraHeaders(headers)
	if err != nil {
		return "", fmt.Errorf("headless: error setting custom headers: %w", err)
	}

	err = page.Navigate(url)
	if err != nil {
		return "", fmt.Errorf("headless: error while navigating: %w", err)
	}

	err = page.WaitLoad()
	if err != nil {
		return "", fmt.Errorf("headless: error while waiting for document load: %w", err)
	}
	content, err := page.HTML()
	if err != nil {
		return "", fmt.Errorf("headless: error fetching html content: %w", err)
	}

	return content, nil
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
