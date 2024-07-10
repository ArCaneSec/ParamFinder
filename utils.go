package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

func extractJsPath(html string) []string {

	matches := jsFiles.FindAllString(html, -1)

	allPaths := make([]string, 0, len(matches))

	for _, path := range matches {

		// removing useless js files
		if strings.Contains(path, "cdn") || strings.Contains(path, "jquery") {
			continue
		}

		allPaths = append(allPaths, path)
	}

	return allPaths
}

// func sendRawRequest(url string, paths []string) (<-chan string, <-chan error) {
// 	client := http.Client{Timeout: 3 * time.Second}
// 	responses := make(chan string, len(paths))
// 	errors := make(chan error, len(paths))

// 	var wg sync.WaitGroup
// 	for _, url := range paths {
// 		wg.Add(1)
// 		go func() {
// 			defer wg.Done()
// 			log.Printf("[*] Sending request to %s\n", url)
// 			req, err := http.NewRequest("GET", url+"/"+url, nil)
// 			if err != nil {
// 				errors <- fmt.Errorf("[!] Failed to construct the request object, url: %s, err: %w", url, err)
// 				return
// 			}

// 			req.Header = http.Header{
// 				"User-Agent":      {"Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/114.0"},
// 				"Accept":          {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8"},
// 				"Accept-Language": {"en-US,en;q=0.5"},
// 				"Sec-Fetch-Dest":  {"document"},
// 				"Sec-Fetch-Mode":  {"navigate"},
// 				"Sec-Fetch-Site":  {"none"},
// 				"Sec-Fetch-User":  {"?1"},
// 				"Referer":         {url},
// 			}

// 			res, err := client.Do(req)
// 			if err != nil {
// 				errors <- fmt.Errorf("[!] Request to %s url failed with error: %w", url, err)
// 				return
// 			}

// 			defer res.Body.Close()

// 			resBody, err := io.ReadAll(res.Body)
// 			if err != nil {
// 				log.Printf("[!] Failed to read response for %s url, error: %s\n", url, err.Error())
// 			}

//				responses <- string(resBody)
//			}()
//		}
//		go func() {
//			wg.Wait()
//			close(responses)
//			close(errors)
//		}()
//		return responses, errors
//	}
func rawRequest(url string) (string, error) {
	client := http.Client{Timeout: 3 * time.Second}

	log.Printf("[*] Sending request to %s\n", url)
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
		log.Printf("[!] Failed to read response for %s url, error: %s\n", url, err.Error())
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
