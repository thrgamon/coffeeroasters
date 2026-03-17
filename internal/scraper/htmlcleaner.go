package scraper

import (
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

const maxCleanedHTMLBytes = 120_000

// CleanHTML strips non-content elements and attributes from HTML, producing
// a smaller document suitable for LLM extraction. If contentSelector is
// non-empty, only elements matching that CSS selector are kept.
func CleanHTML(rawHTML string, contentSelector string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(rawHTML))
	if err != nil {
		return "", err
	}

	// If a content selector is configured, extract only that portion.
	if contentSelector != "" {
		selected := doc.Find(contentSelector)
		if selected.Length() > 0 {
			var parts []string
			selected.Each(func(_ int, s *goquery.Selection) {
				h, err := goquery.OuterHtml(s)
				if err == nil {
					parts = append(parts, h)
				}
			})
			doc, err = goquery.NewDocumentFromReader(strings.NewReader(strings.Join(parts, "\n")))
			if err != nil {
				return "", err
			}
		}
	}

	// Remove elements that add noise without content value.
	doc.Find("head, script, style, nav, footer, header, svg, noscript, iframe, link, meta, form, input, button, select, textarea").Remove()

	// Remove HTML comment nodes.
	var removeComments func(*html.Node)
	removeComments = func(n *html.Node) {
		for c := n.FirstChild; c != nil; {
			next := c.NextSibling
			if c.Type == html.CommentNode {
				n.RemoveChild(c)
			} else {
				removeComments(c)
			}
			c = next
		}
	}
	for _, n := range doc.Find("*").Nodes {
		removeComments(n)
	}

	// Remove duplicate/decorative images (keep only the first img per parent).
	doc.Find("div, li, a").Each(func(_ int, s *goquery.Selection) {
		imgs := s.ChildrenFiltered("img, div > img")
		if imgs.Length() > 1 {
			imgs.Slice(1, imgs.Length()).Remove()
		}
	})

	// Strip all attributes except href, src, and alt.
	keepAttrs := map[string]bool{
		"href": true,
		"src":  true,
		"alt":  true,
	}

	doc.Find("*").Each(func(_ int, s *goquery.Selection) {
		node := s.Get(0)
		var kept []html.Attribute
		for _, attr := range node.Attr {
			if keepAttrs[attr.Key] {
				kept = append(kept, attr)
			}
		}
		node.Attr = kept
	})

	// Remove empty elements (no text content, no meaningful children).
	for range 3 {
		doc.Find("div, span, p, li, ul, ol, section, article").Each(func(_ int, s *goquery.Selection) {
			text := strings.TrimSpace(s.Text())
			hasImg := s.Find("img").Length() > 0
			hasLink := s.Find("a").Length() > 0
			if text == "" && !hasImg && !hasLink {
				s.Remove()
			}
		})
	}

	result, err := doc.Html()
	if err != nil {
		return "", err
	}

	// Collapse whitespace.
	result = collapseWhitespace(result)

	if len(result) > maxCleanedHTMLBytes {
		result = result[:maxCleanedHTMLBytes]
	}

	return result, nil
}

var wsRun = regexp.MustCompile(`\s{2,}`)

func collapseWhitespace(s string) string {
	return wsRun.ReplaceAllString(s, " ")
}
