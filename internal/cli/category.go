package cli

import (
	"errors"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"snip/internal/snip"
)

var categoryActions = map[string]struct{}{
	"list":   {},
	"show":   {},
	"find":   {},
	"clip":   {},
	"edit":   {},
	"delete": {},
}

func runCategoryDefault(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return cmd.Help()
	}
	return handleCategory(cmd, args)
}

func handleCategory(cmd *cobra.Command, args []string) error {
	category, err := requireCategory(args)
	if err != nil {
		return err
	}
	if len(args) == 1 {
		return listHeadlines(cmd, category)
	}

	action := args[1]
	if _, ok := categoryActions[action]; !ok {
		headline := strings.Join(args[1:], " ")
		return showEntry(cmd, category, headline)
	}

	switch action {
	case "list":
		return listHeadlines(cmd, category)
	case "show":
		if len(args) < 3 {
			return errorf("headline is required")
		}
		headline := strings.Join(args[2:], " ")
		return showEntry(cmd, category, headline)
	case "find":
		if len(args) < 3 {
			return errorf("query is required")
		}
		query := strings.Join(args[2:], " ")
		return findInCategory(cmd, category, query)
	case "clip":
		if len(args) < 3 {
			return errorf("headline is required")
		}
		headline := strings.Join(args[2:], " ")
		return clipEntry(cmd, category, headline)
	case "edit":
		return editCategory(cmd, category)
	case "delete":
		if len(args) < 3 {
			return errorf("headline is required")
		}
		headline := strings.Join(args[2:], " ")
		return deleteEntry(cmd, category, headline)
	default:
		return errorf("unknown action: %s", action)
	}
}

func listHeadlines(cmd *cobra.Command, category string) error {
	content, err := snip.ReadCategory(category)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return errorf("category not found: %s", category)
		}
		return err
	}
	entries := snip.ParseCategory(content)
	for _, headline := range snip.HeadlineList(entries) {
		cmd.Println(headline)
	}
	return nil
}

func showEntry(cmd *cobra.Command, category string, headline string) error {
	content, err := snip.ReadCategory(category)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return errorf("category not found: %s", category)
		}
		return err
	}
	entries := snip.ParseCategory(content)
	entry, ok := snip.FindEntry(entries, headline)
	if !ok {
		return errorf("headline not found: %s", headline)
	}
	cmd.Print(formatEntry(entry))
	return nil
}

func findInCategory(cmd *cobra.Command, category string, query string) error {
	content, err := snip.ReadCategory(category)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return errorf("category not found: %s", category)
		}
		return err
	}
	matches := findMatchesLines(content, query)
	for _, line := range matches {
		cmd.Println(line)
	}
	return nil
}

func clipEntry(cmd *cobra.Command, category string, headline string) error {
	content, err := snip.ReadCategory(category)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return errorf("category not found: %s", category)
		}
		return err
	}
	entries := snip.ParseCategory(content)
	entry, ok := snip.FindEntry(entries, headline)
	if !ok {
		return errorf("headline not found: %s", headline)
	}
	if !entry.HasSnippet {
		return errorf("no snippet found for: %s", headline)
	}
	if err := copyToClipboard(entry.Snippet); err != nil {
		if errors.Is(err, ErrNoClipboard) {
			cmd.Print(entry.Snippet)
			return nil
		}
		return err
	}
	infof(cmd, "copied snippet for %s", headline)
	return nil
}

func deleteEntry(cmd *cobra.Command, category string, headline string) error {
	if err := snip.DeleteEntry(category, headline); err != nil {
		if errors.Is(err, snip.ErrEntryNotFound) {
			return errorf("headline not found: %s", headline)
		}
		if errors.Is(err, os.ErrNotExist) {
			return errorf("category not found: %s", category)
		}
		return err
	}
	infof(cmd, "deleted %q from %s", headline, category)
	return nil
}

func editCategory(cmd *cobra.Command, category string) error {
	path, err := categoryFilePath(category)
	if err != nil {
		return err
	}
	if _, err := os.Stat(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return errorf("category not found: %s", category)
		}
		return err
	}
	return openEditor(path)
}

func headlineOrActionCompletion(category string, toComplete string) ([]string, cobra.ShellCompDirective) {
	actions := []string{"list", "show", "find", "clip", "edit", "delete"}
	var choices []string
	for _, action := range actions {
		if strings.HasPrefix(action, toComplete) {
			choices = append(choices, action+"\taction")
		}
	}
	headlines, _ := headlineMatches(category, toComplete)
	choices = append(choices, headlines...)
	return choices, cobra.ShellCompDirectiveNoFileComp
}

func headlineMatches(category string, toComplete string) ([]string, cobra.ShellCompDirective) {
	content, err := snip.ReadCategory(category)
	if err != nil {
		return nil, cobra.ShellCompDirectiveDefault
	}
	entries := snip.ParseCategory(content)
	var filtered []string
	for _, headline := range snip.HeadlineList(entries) {
		if strings.HasPrefix(strings.ToLower(headline), strings.ToLower(toComplete)) {
			filtered = append(filtered, headline+"\theadline")
		}
	}
	return filtered, cobra.ShellCompDirectiveNoFileComp
}
