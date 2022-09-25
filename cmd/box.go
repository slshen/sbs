package cmd

import (
	"io"
	"os"
	"os/exec"

	"github.com/slshen/sb/pkg/boxscore"
	"github.com/slshen/sb/pkg/game"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func boxCommand() *cobra.Command {
	var (
		yamlFormat   bool
		pdfFormat    bool
		scoringPlays bool
		plays        bool
		re           reArgs
	)
	c := &cobra.Command{
		Use:   "box",
		Short: "Generate a box score",
		RunE: func(cmd *cobra.Command, args []string) error {
			re, err := re.getRunExpectancy()
			if err != nil {
				return err
			}
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
				box, err := boxscore.NewBoxScore(g, re)
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
	re.registerFlags(c.Flags())
	return c
}
