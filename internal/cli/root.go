package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"snip/internal/snip"
)

func Execute() error {
	return NewRootCmd().Execute()
}

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "snip",
		Short: "snip is a category-based snippet manager",
		Args:  cobra.ArbitraryArgs,
		RunE:  runCategoryDefault,
	}

	rootCmd.AddCommand(NewListCmd())
	rootCmd.AddCommand(NewCreateCmd())
	rootCmd.AddCommand(NewFindCmd())
	rootCmd.AddCommand(NewAddCmd())
	rootCmd.CompletionOptions.DisableDefaultCmd = false
	rootCmd.ValidArgsFunction = categoryOrActionCompletion

	rootCmd.PersistentFlags().Bool("version", false, "show version")
	rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		cmd.Println("snip - category-based snippets")
		cmd.Println()
		cmd.Println("Usage:")
		cmd.Println("  snip list")
		cmd.Println("  snip create <category>")
		cmd.Println("  snip find <query>")
		cmd.Println("  snip add <category> <headline> [--snippet <text>] [--body <text>] [--lang <lang>]")
		cmd.Println("  snip <category> [list|show|find|clip|edit|delete]")
		cmd.Println()
		cmd.Println("Completion hints:")
		cmd.Println("  categories show as 'category', actions show as 'action', headlines show as 'headline'")
		cmd.Println()
		cmd.Println("Use 'snip completion --help' to generate shell completions.")
	})

	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if cmd.Flags().Changed("version") {
			cmd.Println("snip dev")
			os.Exit(0)
		}
		return nil
	}

	rootCmd.PersistentPreRunE = wrapPersistent(rootCmd.PersistentPreRunE)

	return rootCmd
}

func wrapPersistent(next func(cmd *cobra.Command, args []string) error) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if next != nil {
			if err := next(cmd, args); err != nil {
				return err
			}
		}
		return nil
	}
}

func categoryOrActionCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		return categoryArgCompletion(cmd, args, toComplete)
	}
	if len(args) == 1 {
		return headlineOrActionCompletion(args[0], toComplete)
	}
	if len(args) == 2 {
		if args[1] == "show" || args[1] == "clip" || args[1] == "delete" {
			return headlineMatches(args[0], toComplete)
		}
		if args[1] == "find" {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
	}
	return nil, cobra.ShellCompDirectiveDefault
}

func categoryArgCompletion(_ *cobra.Command, _ []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	cats, err := snip.ListCategories()
	if err != nil {
		return nil, cobra.ShellCompDirectiveDefault
	}
	var filtered []string
	for _, cat := range cats {
		if strings.HasPrefix(cat, toComplete) {
			filtered = append(filtered, cat+"\tcategory")
		}
	}
	return filtered, cobra.ShellCompDirectiveNoFileComp
}

func openEditor(path string) error {
	editor := strings.TrimSpace(os.Getenv("EDITOR"))
	if editor == "" {
		editor = "vim"
	}
	parts := strings.Fields(editor)
	bin := parts[0]
	args := append(parts[1:], path)
	return runCommand(bin, args...)
}

func requireCategory(args []string) (string, error) {
	if len(args) == 0 {
		return "", errors.New("category is required")
	}
	return args[0], nil
}

func categoryFilePath(category string) (string, error) {
	path, err := snip.CategoryFilePath(category)
	if err != nil {
		return "", err
	}
	return filepath.Clean(path), nil
}

func formatEntry(entry snip.Entry) string {
	var builder strings.Builder
	builder.WriteString(snip.HeadingPrefix)
	builder.WriteString(entry.Headline)
	builder.WriteString("\n")
	if entry.Body != "" {
		builder.WriteString(entry.Body)
		builder.WriteString("\n")
	}
	return builder.String()
}

func findMatchesLines(content string, query string) []string {
	return snip.FindMatches(content, query)
}

func errorf(format string, args ...any) error {
	return fmt.Errorf(format, args...)
}

func infof(cmd *cobra.Command, format string, args ...any) {
	cmd.Printf(format+"\n", args...)
}
