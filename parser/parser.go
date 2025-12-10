package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

// ParseHNPage parses the HTML of a Hacker News page and extracts stories
func ParseHNPage(htmlContent string, pageNum int) ([]Story, error) {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	var stories []Story

	// Find all story rows (tr with class "athing submission")
	storyRows := findStoryRows(doc)

	for _, row := range storyRows {
		story := parseStoryRow(row, pageNum)
		if story != nil {
			stories = append(stories, *story)
		}
	}

	return stories, nil
}

// findStoryRows finds all <tr class="athing submission"> elements
func findStoryRows(n *html.Node) []*html.Node {
	var rows []*html.Node
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "tr" {
			if hasClass(n, "athing") && hasClass(n, "submission") {
				rows = append(rows, n)
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(n)
	return rows
}

// parseStoryRow parses a story row and its sibling subtext row
func parseStoryRow(row *html.Node, pageNum int) *Story {
	story := &Story{Page: pageNum}

	// Get story ID from the row's id attribute
	story.ID = getAttr(row, "id")

	// Find rank
	rankSpan := findByClass(row, "rank")
	if rankSpan != nil {
		rankText := getTextContent(rankSpan)
		rankText = strings.TrimSuffix(rankText, ".")
		story.Rank, _ = strconv.Atoi(rankText)
	}

	// Find titleline span for headline and URL
	titleline := findByClass(row, "titleline")
	if titleline != nil {
		// First <a> inside titleline has the headline and URL
		link := findFirstElement(titleline, "a")
		if link != nil {
			story.Headline = getTextContent(link)
			story.URL = getAttr(link, "href")
		}
	}

	// Build discussion URL
	if story.ID != "" {
		story.DiscussionURL = fmt.Sprintf("https://news.ycombinator.com/item?id=%s", story.ID)
	}

	// Find the subtext row (next sibling <tr>)
	subtextRow := findNextSibling(row, "tr")
	if subtextRow != nil {
		parseSubtextRow(subtextRow, story)
	}

	// Only return if we got at least the ID and headline
	if story.ID == "" || story.Headline == "" {
		return nil
	}

	return story
}

// parseSubtextRow extracts points, username, age, and comments from the subtext row
func parseSubtextRow(row *html.Node, story *Story) {
	// Find score span
	scoreSpan := findByClass(row, "score")
	if scoreSpan != nil {
		scoreText := getTextContent(scoreSpan)
		// "123 points" -> 123
		parts := strings.Fields(scoreText)
		if len(parts) > 0 {
			story.Points, _ = strconv.Atoi(parts[0])
		}
	}

	// Find username
	userLink := findByClass(row, "hnuser")
	if userLink != nil {
		story.Username = getTextContent(userLink)
	}

	// Find age
	ageSpan := findByClass(row, "age")
	if ageSpan != nil {
		ageLink := findFirstElement(ageSpan, "a")
		if ageLink != nil {
			ageText := getTextContent(ageLink)
			story.AgeValue, story.AgeUnit = parseAge(ageText)
		}
	}

	// Find comments - it's in the last <a> of subline that contains "comment" or "discuss"
	subline := findByClass(row, "subline")
	if subline != nil {
		links := findAllElements(subline, "a")
		for _, link := range links {
			text := getTextContent(link)
			if strings.Contains(text, "comment") {
				// "114 comments" or "1 comment"
				numStr := strings.Fields(text)[0]
				// Handle non-breaking space
				numStr = strings.ReplaceAll(numStr, "\u00a0", "")
				story.Comments, _ = strconv.Atoi(numStr)
				break
			} else if text == "discuss" {
				story.Comments = 0
				break
			}
		}
	}
}

// parseAge parses strings like "4 hours ago" into value and unit
func parseAge(ageStr string) (int, string) {
	// Remove "ago" and extra spaces
	ageStr = strings.TrimSuffix(ageStr, " ago")
	ageStr = strings.TrimSpace(ageStr)

	// Match pattern: number + unit
	re := regexp.MustCompile(`^(\d+)\s+(\w+)$`)
	matches := re.FindStringSubmatch(ageStr)
	if len(matches) == 3 {
		value, _ := strconv.Atoi(matches[1])
		unit := matches[2]
		// Normalize unit (remove trailing 's' for consistency, or keep as-is)
		return value, unit
	}

	return 0, ""
}

// Helper functions for HTML parsing

func hasClass(n *html.Node, class string) bool {
	for _, attr := range n.Attr {
		if attr.Key == "class" {
			classes := strings.Fields(attr.Val)
			for _, c := range classes {
				if c == class {
					return true
				}
			}
		}
	}
	return false
}

func getAttr(n *html.Node, key string) string {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

func findByClass(n *html.Node, class string) *html.Node {
	var result *html.Node
	var f func(*html.Node)
	f = func(n *html.Node) {
		if result != nil {
			return
		}
		if n.Type == html.ElementNode && hasClass(n, class) {
			result = n
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(n)
	return result
}

func findFirstElement(n *html.Node, tag string) *html.Node {
	var result *html.Node
	var f func(*html.Node)
	f = func(n *html.Node) {
		if result != nil {
			return
		}
		if n.Type == html.ElementNode && n.Data == tag {
			result = n
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(n)
	return result
}

func findAllElements(n *html.Node, tag string) []*html.Node {
	var results []*html.Node
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == tag {
			results = append(results, n)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(n)
	return results
}

func findNextSibling(n *html.Node, tag string) *html.Node {
	for s := n.NextSibling; s != nil; s = s.NextSibling {
		if s.Type == html.ElementNode && s.Data == tag {
			return s
		}
	}
	return nil
}

func getTextContent(n *html.Node) string {
	var sb strings.Builder
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.TextNode {
			sb.WriteString(n.Data)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(n)
	return strings.TrimSpace(sb.String())
}
