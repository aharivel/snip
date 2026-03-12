package snip

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseCategorySample(t *testing.T) {
	path := filepath.Join("..", "sniptest", "sample.md")
	bytes, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read sample: %v", err)
	}
	entries := ParseCategory(string(bytes))
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	entry, ok := FindEntry(entries, "Login to cluster")
	if !ok {
		t.Fatalf("entry not found")
	}
	if !entry.HasSnippet {
		t.Fatalf("expected snippet")
	}
	if entry.SnippetLang != "bash" {
		t.Fatalf("expected snippet lang bash, got %q", entry.SnippetLang)
	}
	if entry.Snippet == "" {
		t.Fatalf("expected snippet content")
	}
}
