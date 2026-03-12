package cli

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"snip/internal/snip"
)

func NewAddCmd() *cobra.Command {
	var snippetFlag string
	var bodyFlag string
	var langFlag string

	cmd := &cobra.Command{
		Use:   "add <category> <headline>",
		Short: "Add an entry to a category",
		Long: `Add an entry to a category.

The snippet can be provided via --snippet, or piped from stdin.
If neither is provided, only a headline (and optional --body) is stored.`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			category := args[0]
			headline := args[1]

			snippet := snippetFlag

			// read snippet from stdin if it is a pipe and --snippet was not given
			if snippet == "" {
				fi, err := os.Stdin.Stat()
				if err == nil && (fi.Mode()&os.ModeCharDevice) == 0 {
					data, err := io.ReadAll(os.Stdin)
					if err != nil {
						return err
					}
					snippet = strings.TrimRight(string(data), "\n")
				}
			}

			body := buildEntryBody(bodyFlag, snippet, langFlag)

			entry := snip.Entry{
				Headline:    headline,
				Body:        body,
				Snippet:     snippet,
				HasSnippet:  snippet != "",
				SnippetLang: langFlag,
			}

			if err := snip.AddEntry(category, entry); err != nil {
				if errors.Is(err, snip.ErrEntryExists) {
					return fmt.Errorf("entry already exists: %s", headline)
				}
				if errors.Is(err, os.ErrNotExist) {
					return fmt.Errorf("category not found: %s (use snip create %s first)", category, category)
				}
				return err
			}
			infof(cmd, "added %q to %s", headline, category)
			return nil
		},
	}

	cmd.Flags().StringVar(&snippetFlag, "snippet", "", "snippet text to store")
	cmd.Flags().StringVar(&bodyFlag, "body", "", "description text above the snippet")
	cmd.Flags().StringVar(&langFlag, "lang", "", "language for the snippet fence (e.g. bash, go)")
	return cmd
}

func buildEntryBody(body, snippet, lang string) string {
	if snippet == "" && body == "" {
		return ""
	}
	var b strings.Builder
	if body != "" {
		b.WriteString(body)
		if snippet != "" {
			b.WriteString("\n")
		}
	}
	if snippet != "" {
		fence := "```"
		b.WriteString(fence + lang + "\n")
		b.WriteString(snippet)
		b.WriteString("\n" + fence)
	}
	return b.String()
}
