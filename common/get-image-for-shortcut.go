package common

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"net/url"
)

func GetImageForShortcut(pageURL string) (string, error) {
	parsedURL, err := url.Parse(pageURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %v", err)
	}

	resp, err := http.Get(pageURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch the page: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status: %s", resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to parse the page: %v", err)
	}

	var bestImageURL string

	resolveURL := func(href string) string {
		u, err := url.Parse(href)
		if err != nil || !u.IsAbs() {
			u = parsedURL.ResolveReference(u)
		}
		return u.String()
	}

	doc.Find("link[rel='apple-touch-icon'], link[rel='apple-touch-icon-precomposed']").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			bestImageURL = resolveURL(href)
		}
	})

	if bestImageURL == "" {
		doc.Find("link[rel='icon'], link[rel='shortcut icon']").Each(func(i int, s *goquery.Selection) {
			href, exists := s.Attr("href")
			if exists {
				bestImageURL = resolveURL(href)
			}
		})
	}

	if bestImageURL == "" {
		return "", fmt.Errorf("no suitable image found")
	}

	return bestImageURL, nil
}
