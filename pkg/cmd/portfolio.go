package cmd

import (
	"github.com/spf13/cobra"
)

// NewPortfolioCmd creates a new root command for Portfolio-Go.
func NewPortfolioCmd() *cobra.Command {
	portfolioCmd := &cobra.Command{
		Use:   "portfolio",
		Short: "Portfolio-Go",
		Long:  ``,
	}

	addCommands(portfolioCmd)

	return portfolioCmd
}

func addCommands(cmd *cobra.Command) {
	cmd.AddCommand(
		newStartCmd(),
	)
}
