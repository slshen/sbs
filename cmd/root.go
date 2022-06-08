package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/slshen/sb/pkg/boxscore"
	"github.com/slshen/sb/pkg/dataframe"
	"github.com/slshen/sb/pkg/export"
	"github.com/slshen/sb/pkg/game"
	"github.com/slshen/sb/pkg/playbyplay"
	"github.com/slshen/sb/pkg/stats"
	"github.com/slshen/sb/pkg/tournament"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func Root() *cobra.Command {
	root := &cobra.Command{}
	root.SilenceUsage = true
	root.AddCommand(readCommand(), boxCommand(), playByPlayCommand(),
		statsCommand("batting"), statsCommand("pitching"), reCommand(),
		exportCommand(), tournamentCommand(), reAnalysisCommand())
	return root
}

func readCommand() *cobra.Command {
	var home, visitor bool
	c := &cobra.Command{
		Use:   "read",
		Short: "Read and print a score file",
		RunE: func(cmd *cobra.Command, args []string) error {
			games, err := game.ReadGameFiles(args)
			for _, g := range games {
				states := g.GetStates()
				for _, state := range states {
					if home || visitor {
						if visitor && !state.Top() {
							continue
						}
						if home && state.Top() {
							continue
						}
					}
					d, _ := yaml.Marshal(state)
					fmt.Println(string(d))
				}
			}
			return err
		},
	}
	c.Flags().BoolVar(&home, "home", false, "Print only home plays")
	c.Flags().BoolVar(&visitor, "visitor", false, "Print only visitor plays")
	return c
}

func boxCommand() *cobra.Command {
	var (
		yamlFormat   bool
		pdfFormat    bool
		scoringPlays bool
		plays        bool
	)
	c := &cobra.Command{
		Use:   "box",
		Short: "Generate a box score",
		RunE: func(cmd *cobra.Command, args []string) error {
			games, err := game.ReadGameFiles(args)
			if err != nil {
				return err
			}
			var out io.Writer
			if pdfFormat {
				paps := exec.Command("paps", "--format=pdf", "--font=Andale Mono 10",
					"--left-margin=18", "--right-margin=18", "--top-margin=18", "--bottom-margin=18")
				w, err := paps.StdinPipe()
				paps.Stdout = os.Stdout
				paps.Stderr = os.Stderr
				if err != nil {
					return err
				}
				defer w.Close()
				out = w
				if err := paps.Start(); err != nil {
					return err
				}
			} else {
				out = os.Stdout
			}
			for i, g := range games {
				box, err := boxscore.NewBoxScore(g, nil)
				if err != nil {
					return err
				}
				box.IncludeScoringPlays = scoringPlays
				box.IncludePlays = plays
				if yamlFormat {
					dat, err := yaml.Marshal(box)
					if err != nil {
						return err
					}
					if _, err := out.Write(dat); err != nil {
						return err
					}
				} else if err := box.Render(out); err != nil {
					return err
				}
				if i != len(games)-1 {
					if _, err := out.Write([]byte{'\f'}); err != nil {
						return err
					}
				}
			}
			return nil
		},
	}
	c.Flags().BoolVar(&yamlFormat, "yaml", false, "")
	c.Flags().BoolVar(&pdfFormat, "pdf", false, "Run paps to convert output to pdf")
	c.Flags().BoolVar(&scoringPlays, "scoring", false, "Include scoring plays in box")
	c.Flags().BoolVar(&plays, "plays", false, "Include play by play in box")
	return c
}

func playByPlayCommand() *cobra.Command {
	pbp := playbyplay.Generator{}
	c := &cobra.Command{
		Use:   "plays",
		Short: "Generate play by play",
		RunE: func(cmd *cobra.Command, args []string) error {
			games, err := game.ReadGameFiles(args)
			if err != nil {
				return err
			}
			for _, g := range games {
				pbp.Game = g
				if err := pbp.Generate(os.Stdout); err != nil {
					return err
				}
			}
			return nil
		},
	}
	c.Flags().BoolVar(&pbp.ScoringOnly, "scoring", false, "Only show scoring plays")
	return c
}

func statsCommand(statsType string) *cobra.Command {
	var (
		csv bool
		re  reArgs
	)
	mg := stats.NewGameStats(nil)
	c := &cobra.Command{
		Use:     fmt.Sprintf("%s-stats", statsType),
		Aliases: []string{statsType},
		Short:   "Print stats",
		RunE: func(cmd *cobra.Command, args []string) error {
			games, err := game.ReadGameFiles(args)
			if err != nil {
				return err
			}
			mg.RE, err = re.getRunExpectancy()
			if err != nil {
				return err
			}
			for _, g := range games {
				if err := mg.Read(g); err != nil {
					return err
				}
			}
			var data *dataframe.Data
			if statsType == "batting" {
				data = mg.GetBattingData()
			} else {
				data = mg.GetPitchingData()
			}
			if csv {
				return data.RenderCSV(os.Stdout)
			} else {
				fmt.Println(data)
			}
			return nil
		},
	}
	re.registerFlags(c.Flags())
	c.Flags().BoolVar(&csv, "csv", false, "Print in CSV format")
	return c
}

func reCommand() *cobra.Command {
	var (
		csv        bool
		yamlFormat bool
		freq       bool
		pivot      bool
		raw        bool
		bandwidth  float64
	)
	re24 := &stats.ObservedRunExpectancy{}
	c := &cobra.Command{
		Use:   "re",
		Short: "Determine the run expectancy matrix",
		RunE: func(cmd *cobra.Command, args []string) error {
			games, err := game.ReadGameFiles(args)
			if err != nil {
				return err
			}
			for _, g := range games {
				if err := re24.Read(g); err != nil {
					return err
				}
			}
			if yamlFormat {
				return re24.WriteYAML(os.Stdout)
			}
			var data *dataframe.Data
			switch {
			case raw:
				data = re24.GetRunData()
			case pivot || freq:
				rf := re24.GetRunExpectancyFrequency()
				if pivot {
					data = rf.Pivot()
				} else {
					data = &rf.Data
				}
			default:
				data = stats.GetRunExpectancyData(re24)
			}
			if csv {
				return data.RenderCSV(os.Stdout)
			}
			fmt.Println(data)
			return nil
		},
	}
	c.Flags().BoolVar(&csv, "csv", false, "Print in CSV format")
	c.Flags().BoolVar(&yamlFormat, "yaml", false, "Print in YAML format")
	c.Flags().BoolVar(&freq, "freq", false, "Print the frequency of # runs scored per 24-base/out state")
	c.Flags().BoolVar(&pivot, "pivot", false, "Pivot the frequency data by runs")
	c.Flags().BoolVar(&raw, "raw", false, "Get the raw run data")
	c.Flags().Float64Var(&bandwidth, "bandwidth", 0, "KDE bandwidth")
	return c
}

func exportCommand() *cobra.Command {
	var (
		us            string
		league        string
		spreadsheetID string
		dryRun        bool
		re            reArgs
	)
	c := &cobra.Command{
		Use:   "export",
		Short: "Export games and stats to Google sheets",
		RunE: func(cmd *cobra.Command, args []string) error {
			if us == "" {
				return fmt.Errorf("--us is required")
			}
			config, err := export.NewConfig()
			if err != nil {
				return err
			}
			if spreadsheetID != "" {
				config.SpreadsheetID = spreadsheetID
			}
			sheets, err := export.NewSheetExport(config)
			if err != nil {
				return err
			}
			re, err := re.getRunExpectancy()
			if err != nil {
				return err
			}
			export, err := export.NewExport(sheets, re)
			export.DryRun = dryRun
			if err != nil {
				return err
			}
			export.Us = us
			export.League = strings.ToLower(league)
			games, err := game.ReadGameFiles(args)
			if err != nil {
				return err
			}
			return export.Export(games)
		},
	}
	re.registerFlags(c.Flags())
	c.Flags().StringVar(&us, "us", "", "Our `team`")
	c.Flags().StringVar(&league, "league", "", "Include only games in `league`")
	c.Flags().StringVar(&spreadsheetID, "spreadsheet-id", "", "The spreadsheet to use")
	c.Flags().BoolVar(&dryRun, "dry-run", false, "Print the sheets instead of uploading them")
	return c
}

func tournamentCommand() *cobra.Command {
	var (
		us             string
		plays          int
		re             reArgs
		tournamentName string
		playsOnly      bool
	)
	c := &cobra.Command{
		Use:   "tournament",
		Short: "Report on tournament results",
		RunE: func(cmd *cobra.Command, args []string) error {
			if us == "" {
				return fmt.Errorf("--us is required")
			}
			re, err := re.getRunExpectancy()
			if err != nil {
				return err
			}
			games, err := game.ReadGameFiles(args)
			if err != nil {
				return err
			}
			for _, gr := range tournament.GroupByTournament(games) {
				if tournamentName != "" && !strings.Contains(strings.ToLower(gr.Name), tournamentName) {
					continue
				}
				rep := &tournament.Report{
					Us:    us,
					Group: gr,
				}
				if err := rep.Run(re); err != nil {
					return err
				}
				if !playsOnly {
					fmt.Println(rep.GetBattingData())
				}
				topPlays := rep.GetBestAndWorstRE24(plays)
				if playsOnly {
					topPlays = topPlays.Select(
						dataframe.Col("Rnk"),
						dataframe.Col("Game"),
						dataframe.Col("ID"),
						dataframe.Col("O"),
						dataframe.Col("Rnr"),
						dataframe.Col("Play"),
						dataframe.Col("RE24"),
					)
				}
				fmt.Println(topPlays)
			}
			return nil
		},
	}
	re.registerFlags(c.Flags())
	c.Flags().StringVar(&us, "us", "", "Our `team`")
	c.Flags().IntVar(&plays, "plays", 15, "Show the top and bottom `n` plays by RE24")
	c.Flags().StringVar(&tournamentName, "tournament", "", "Show only `tournament`")
	c.Flags().BoolVar(&playsOnly, "plays-only", false, "Only list top plays")
	return c
}

func reAnalysisCommand() *cobra.Command {
	var (
		re reArgs
	)
	c := &cobra.Command{
		Use:   "re-analysis",
		Short: "Analyze run expectancy",
		RunE: func(cmd *cobra.Command, args []string) error {
			re, err := re.getRunExpectancy()
			if err != nil {
				return err
			}
			if re == nil {
				return fmt.Errorf("no RE specified")
			}
			rea := stats.NewREAnalysis(re)
			fmt.Println(rea.Run())
			return nil
		},
	}
	re.registerFlags(c.Flags())
	return c
}
