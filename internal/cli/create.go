package cli

import (
	"errors"

	"github.com/spf13/cobra"
	"snip/internal/snip"
)

func NewCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create <category>",
		Short: "Create a new category file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			category := args[0]
			path, err := snip.CreateCategory(category)
			if err != nil {
				if errors.Is(err, snip.ErrCategoryExists) {
					return errorf("category already exists: %s", category)
				}
				return err
			}
			cmd.Printf("created %s\n", path)
			return nil
		},
	}
	return cmd
}
