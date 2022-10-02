package cmd

import (
	"github.com/spf13/cobra"
)

func Root() *cobra.Command {
	root := &cobra.Command{}
	root.SilenceUsage = true
	root.AddCommand(readCommand(), boxCommand(), playByPlayCommand(),
		statsCommand("batting"), statsCommand("pitching"), reCommand(),
		tournamentCommand(), reAnalysisCommand(),
		fmtCommand(), altCommand(), webdataCommand(), newGameCommand(),
		battingCountCommand(), battingTimesSeenPitcherCommand(),
		pitchingTimesSeenLineupCommand(),
	)
	return root
}
