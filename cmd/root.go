package cmd

import (
	"github.com/spf13/cobra"
)

type Bumper interface {
	Bump() error
	WithRepository(repository string) error
	WithTitle(title string)
	WithTarget(target string)
}

func NewCmd(bumper Bumper) *cobra.Command {
	var (
		repository string
		target     string
		title      string
	)
	cmd := &cobra.Command{
		Use:   "bump",
		Short: "bump version of a given repository",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := bumper.WithRepository(repository); err != nil {
				return err
			}
			if target != "" {
				bumper.WithTarget(target)
			}
			if title != "" {
				bumper.WithTitle(title)
			}
			return bumper.Bump()
		},
	}

	cmd.Flags().StringVarP(&repository, "repo", "R", "", "Select another repository using the [HOST/]OWNER/REPO format")
	cmd.Flags().StringVarP(&target, "target", "", "", "Target branch or full commit SHA (default: main branch)")
	cmd.Flags().StringVarP(&title, "title", "t", "", "Release title")
	return cmd
}
