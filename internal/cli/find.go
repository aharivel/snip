package cli

import (
	"strings"

	"github.com/spf13/cobra"
	"snip/internal/snip"
)

func NewFindCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "find <query>",
		Short: "Search across all categories",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			query := strings.Join(args, " ")
			cats, err := snip.ListCategories()
			if err != nil {
				return err
			}
			for _, cat := range cats {
				content, err := snip.ReadCategory(cat)
				if err != nil {
					continue
				}
				matches := snip.FindMatches(content, query)
				for _, line := range matches {
					cmd.Printf("%s: %s\n", cat, line)
				}
			}
			return nil
		},
	}
	return cmd
}
