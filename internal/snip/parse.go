package snip

import (
	"bufio"
	"strings"
)

const HeadingPrefix = "## "

func ParseCategory(content string) []Entry {
	content = strings.ReplaceAll(content, "\r\n", "\n")
	lines := strings.Split(content, "\n")
	var entries []Entry
	var current *Entry
	var bodyLines []string
	var inFence bool
	var fenceLang string
	var fenceLines []string
	var snippetCaptured bool

	flushCurrent := func() {
		if current == nil {
			return
		}
		current.Body = strings.TrimSpace(strings.Join(bodyLines, "\n"))
		if !snippetCaptured {
			current.Snippet = ""
			current.HasSnippet = false
			current.SnippetLang = ""
		}
		entries = append(entries, *current)
	}

	for _, raw := range lines {
		line := strings.TrimRight(raw, "\r")
		if strings.HasPrefix(line, HeadingPrefix) {
			if inFence && current != nil && !snippetCaptured {
				current.Snippet = strings.Join(fenceLines, "\n")
				current.HasSnippet = true
				current.SnippetLang = fenceLang
			}
			flushCurrent()
			headline := strings.TrimSpace(strings.TrimPrefix(line, HeadingPrefix))
			if headline == "" {
				current = nil
				bodyLines = nil
				inFence = false
				fenceLang = ""
				fenceLines = nil
				snippetCaptured = false
				continue
			}
			current = &Entry{Headline: headline}
			bodyLines = nil
			inFence = false
			fenceLang = ""
			fenceLines = nil
			snippetCaptured = false
			continue
		}

		if current == nil {
			continue
		}

		trimmed := strings.TrimSpace(line)
		if fenceStart := detectFenceStart(trimmed); fenceStart != "" {
			if inFence {
				if !snippetCaptured {
					current.Snippet = strings.Join(fenceLines, "\n")
					current.HasSnippet = true
					current.SnippetLang = fenceLang
					snippetCaptured = true
				}
				inFence = false
				fenceLang = ""
				fenceLines = nil
				bodyLines = append(bodyLines, line)
				continue
			}
			inFence = true
			fenceLang = strings.TrimSpace(strings.TrimPrefix(trimmed, fenceStart))
			fenceLines = nil
			bodyLines = append(bodyLines, line)
			continue
		}

		if inFence {
			fenceLines = append(fenceLines, line)
		}
		bodyLines = append(bodyLines, line)
	}

	if inFence && current != nil && !snippetCaptured {
		current.Snippet = strings.Join(fenceLines, "\n")
		current.HasSnippet = true
		current.SnippetLang = fenceLang
	}
	flushCurrent()

	return entries
}

func detectFenceStart(line string) string {
	if strings.HasPrefix(line, "```") {
		return "```"
	}
	if strings.HasPrefix(line, "~~~") {
		return "~~~"
	}
	return ""
}

func HeadlineList(entries []Entry) []string {
	result := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.Headline != "" {
			result = append(result, entry.Headline)
		}
	}
	return result
}

func FindEntry(entries []Entry, headline string) (Entry, bool) {
	for _, entry := range entries {
		if strings.EqualFold(entry.Headline, headline) {
			return entry, true
		}
	}
	return Entry{}, false
}

func FindMatches(content string, query string) []string {
	content = strings.ReplaceAll(content, "\r\n", "\n")
	scanner := bufio.NewScanner(strings.NewReader(content))
	var matches []string
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(strings.ToLower(line), strings.ToLower(query)) {
			matches = append(matches, line)
		}
	}
	return matches
}
