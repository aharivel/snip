package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"snip/internal/snip"
)

func NewListCmd() *cobra.Command {
	var count bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List categories",
		RunE: func(cmd *cobra.Command, args []string) error {
			cats, err := snip.ListCategories()
			if err != nil {
				return err
			}
			for _, cat := range cats {
				if count {
					content, err := snip.ReadCategory(cat)
					if err != nil {
						cmd.Printf("%s (? entries)\n", cat)
						continue
					}
					entries := snip.ParseCategory(content)
					cmd.Println(fmt.Sprintf("%s (%d entries)", cat, len(entries)))
				} else {
					cmd.Println(cat)
				}
			}
			return nil
		},
	}

	cmd.Flags().BoolVarP(&count, "count", "c", false, "show entry count per category")
	return cmd
}
