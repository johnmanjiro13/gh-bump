package cmd

import (
	"bytes"

	"github.com/spf13/cobra"
)

type Bumper interface {
	Bump() error
	WithRepository(repository string) error
}

type Gh interface {
	ViewRepository() (sout, eout bytes.Buffer, err error)
}

func NewCmd(bumper Bumper) *cobra.Command {
	var repository string
	cmd := &cobra.Command{
		Use:   "bump",
		Short: "bump version of a given repository",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := bumper.WithRepository(repository); err != nil {
				return err
			}
			return bumper.Bump()
		},
	}

	cmd.Flags().StringVarP(&repository, "repo", "R", "", "Select another repository using the [HOST/]OWNER/REPO format")
	return cmd
}
