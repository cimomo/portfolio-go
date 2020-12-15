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

	return portfolioCmd
}
