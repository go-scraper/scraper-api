package services

import (
	"io"
	"net/http"
	"net/url"
	"scraper/logger"
	"scraper/models"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/publicsuffix"
)

func FetchPageInfo(client *http.Client, baseURL string) (*models.PageInfo, error) {
	resp, err := client.Get(baseURL)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	defer resp.Body.Close()

	return ParseHTML(resp.Body, baseURL)
}

func ParseHTML(body io.Reader, baseURL string) (*models.PageInfo, error) {
	pageInfo := &models.PageInfo{HeadingCounts: make(map[string]int)}
	doc, err := html.Parse(body)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	visitNode := func(node *html.Node) {
		switch node.Type {
		case html.ElementNode:
			switch node.Data {
			case "html":
				pageInfo.HTMLVersion = extractHtmlVersion(node)
			case "h1", "h2", "h3", "h4", "h5", "h6":
				pageInfo.HeadingCounts[node.Data]++
			case "a":
				href := extractHref(node)
				if href != "" {
					fullURL := resolveURL(baseURL, href)
					if isInternal(baseURL, fullURL) {
						pageInfo.InternalURLsCount++
					} else {
						pageInfo.ExternalURLsCount++
					}
					pageInfo.URLs = append(pageInfo.URLs, models.URLStatus{URL: fullURL})
				}
			case "form":
				if containsPasswordInput(node) {
					pageInfo.ContainsLoginForm = true
				}
			}
		}
	}

	traverse(doc, visitNode)
	pageInfo.Title = extractTitle(doc)
	return pageInfo, nil
}

func traverse(node *html.Node, visit func(*html.Node)) {
	if node == nil {
		return
	}
	visit(node)
	for child_node := node.FirstChild; child_node != nil; child_node = child_node.NextSibling {
		traverse(child_node, visit)
	}
}

func extractHref(node *html.Node) string {
	for _, attr := range node.Attr {
		if attr.Key == "href" {
			return attr.Val
		}
	}
	return ""
}

func resolveURL(baseURL, href string) string {
	base, _ := url.Parse(baseURL)
	rel, _ := url.Parse(href)
	return base.ResolveReference(rel).String()
}

func isInternal(baseUrl, scrappedUrl string) bool {
	baseUrlParsed, _ := url.Parse(baseUrl)
	scrappedUrlParsed, _ := url.Parse(scrappedUrl)

	baseUrlTld, _ := publicsuffix.EffectiveTLDPlusOne(baseUrlParsed.Host)
	scrappedUrlTld, _ := publicsuffix.EffectiveTLDPlusOne(scrappedUrlParsed.Host)

	return strings.EqualFold(baseUrlTld, scrappedUrlTld)
}

func containsPasswordInput(node *html.Node) bool {
	if node.Type == html.ElementNode && node.Data == "input" {
		for _, attr := range node.Attr {
			if attr.Key == "type" && attr.Val == "password" {
				return true
			}
		}
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if containsPasswordInput(child) {
			return true
		}
	}
	return false
}

func extractTitle(node *html.Node) string {
	if node.Type == html.ElementNode && node.Data == "title" && node.FirstChild != nil {
		return node.FirstChild.Data
	}
	for child_node := node.FirstChild; child_node != nil; child_node = child_node.NextSibling {
		title := extractTitle(child_node)
		if title != "" {
			return title
		}
	}
	return ""
}

func extractHtmlVersion(node *html.Node) string {
	// Check for a "version" attribute
	for _, attr := range node.Attr {
		if attr.Key == "version" {
			return attr.Val
		}
	}

	// Check for DOCTYPE version
	if node.Data == "html" && node.Type == html.ElementNode {
		return "HTML 5" // Default to HTML5
	}

	return "Unknown Version"
}
