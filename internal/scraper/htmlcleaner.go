package scraper

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const maxCleanedHTMLBytes = 120_000

// CleanHTML strips non-content elements and attributes from HTML, producing
// a smaller document suitable for LLM extraction.
func CleanHTML(rawHTML string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(rawHTML))
	if err != nil {
		return "", err
	}

	// Remove elements that add noise without content value.
	doc.Find("script, style, nav, footer, header, svg, noscript, iframe, link, meta").Remove()

	// Strip noisy attributes but keep structural/content ones.
	keepAttrs := map[string]bool{
		"href": true,
		"src":  true,
		"alt":  true,
	}

	doc.Find("*").Each(func(_ int, s *goquery.Selection) {
		for _, attr := range s.Get(0).Attr {
			if !keepAttrs[attr.Key] {
				s.RemoveAttr(attr.Key)
			}
		}
	})

	html, err := doc.Html()
	if err != nil {
		return "", err
	}

	if len(html) > maxCleanedHTMLBytes {
		html = html[:maxCleanedHTMLBytes]
	}

	return html, nil
}
