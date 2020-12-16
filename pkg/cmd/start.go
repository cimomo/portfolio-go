package cmd

import (
	"github.com/cimomo/portfolio-go/pkg/terminal"
	"github.com/spf13/cobra"
)

var profile string

func newStartCmd() *cobra.Command {
	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Start a terminal window for portfolio",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			startTerminal(profile)
		},
	}

	startCmd.PersistentFlags().StringVar(&profile, "profile", "./examples/profile.yml", "profile for portfolio")

	return startCmd
}

func startTerminal(profile string) error {
	term := terminal.NewTerminal(profile)

	err := term.Start()
	if err != nil {
		return err
	}

	return nil
}
