package downloader

import (
	"net/url"
	"regexp"
	"strings"
)

func (d *Downloader) extractLinks(content []byte, baseURL string) []string {
	html := string(content)
	var links []string

	// Find all href and src attributes
	patterns := map[string]*regexp.Regexp{
		"href":   regexp.MustCompile(`href\s*=\s*["']([^"']+)["']`),
		"src":    regexp.MustCompile(`src\s*=\s*["']([^"']+)["']`),
		"import": regexp.MustCompile(`@import\s+["']([^"']+)["']`),
	}

	for _, re := range patterns {
		matches := re.FindAllStringSubmatch(html, -1)
		for _, match := range matches {
			if len(match) > 1 {
				link := strings.TrimSpace(match[1])
				if link != "" && d.shouldFollow(link) {
					absolute := d.toAbsolute(link, baseURL)
					if absolute != "" {
						links = append(links, absolute)
					}
				}
			}
		}
	}

	return links
}

func (d *Downloader) shouldFollow(link string) bool {
	if link == "" {
		return false
	}

	skip := []string{"#", "mailto:", "tel:", "javascript:", "data:"}
	for _, prefix := range skip {
		if strings.HasPrefix(link, prefix) {
			return false
		}
	}
	return true
}

func (d *Downloader) toAbsolute(link, baseURL string) string {
	base, err := url.Parse(baseURL)
	if err != nil {
		return ""
	}

	ref, err := url.Parse(link)
	if err != nil {
		return ""
	}

	return base.ResolveReference(ref).String()
}

func (d *Downloader) updateLinks(content []byte, currentURL string) []byte {
	html := string(content)

	// Replace both href and src attributes
	attributes := []string{"href", "src"}
	for _, attr := range attributes {
		pattern := attr + `\s*=\s*["']([^"']+)["']`
		re := regexp.MustCompile(pattern)

		html = re.ReplaceAllStringFunc(html, func(match string) string {
			parts := re.FindStringSubmatch(match)
			if len(parts) < 2 {
				return match
			}

			originalLink := parts[1]
			newLink := d.rewriteLink(originalLink, currentURL)
			return attr + `="` + newLink + `"`
		})
	}

	return []byte(html)
}

func (d *Downloader) rewriteLink(link, currentURL string) string {
	absolute := d.toAbsolute(link, currentURL)

	if d.isSameDomain(absolute) {
		targetPath := d.getLocalPath(absolute)
		currentPath := d.getLocalPath(currentURL)

		relative := d.getRelativePath(currentPath, targetPath)
		if relative != "" {
			return relative
		}
	}

	return link
}