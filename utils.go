package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

func extractJsPath(html string) []string {
	matches := jsFiles.FindAllStringSubmatch(html, -1)

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

func rawRequest(url string) (string, error) {
	client := http.Client{Timeout: 3 * time.Second}
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

func headlessRequest(url string) (string, error) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("ignore-certificate-errors", true),
		chromedp.Flag("disable-web-security", true),
	)
	aCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(
		aCtx,
	)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	defaultHeaders := map[string]interface{}{
		"User-Agent":      "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/114.0",
		"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8",
		"Accept-Language": "en-US,en;q=0.5",
		"Sec-Fetch-Dest":  "document",
		"Sec-Fetch-Mode":  "navigate",
		"Sec-Fetch-Site":  "none",
		"Sec-Fetch-User":  "?1",
		"Referer":         url,
	}

	var htmlContent string
	err := chromedp.Run(ctx,
		network.Enable(),
		network.SetExtraHTTPHeaders(defaultHeaders),
		chromedp.Navigate(url),
		chromedp.WaitReady("body"),
		chromedp.OuterHTML("html", &htmlContent),
	)
	if err != nil {
		return "", fmt.Errorf("[!] Error while running headless browser, check your url and network then try again.\nurl: %s\nerror: %w", url, err)
	}

	return htmlContent, nil
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
