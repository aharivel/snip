package cli

import (
	"github.com/spf13/cobra"
	"snip/internal/snip"
)

func NewListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List categories",
		RunE: func(cmd *cobra.Command, args []string) error {
			cats, err := snip.ListCategories()
			if err != nil {
				return err
			}
			for _, cat := range cats {
				cmd.Println(cat)
			}
			return nil
		},
	}
	return cmd
}
