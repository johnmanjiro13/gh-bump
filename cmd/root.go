package cmd

import (
	"github.com/spf13/cobra"
)

type Bumper interface {
	Bump() error
	WithRepository(repository string) error
	WithDraft()
	WithPrerelease()
	WithNotes(notes string)
	WithNotesFile(filename string)
	WithTitle(title string)
	WithTarget(target string)
}

func New(bumper Bumper) *cobra.Command {
	var (
		repository   string
		isDraft      bool
		isPrerelease bool
		notes        string
		notesFile    string
		target       string
		title        string
	)
	cmd := &cobra.Command{
		Use:   "bump",
		Short: "bump version of a given repository",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := bumper.WithRepository(repository); err != nil {
				return err
			}
			if isDraft {
				bumper.WithDraft()
			}
			if isPrerelease {
				bumper.WithPrerelease()
			}
			if notes != "" {
				bumper.WithNotes(notes)
			}
			if notesFile != "" {
				bumper.WithNotesFile(notesFile)
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
	cmd.Flags().BoolVarP(&isDraft, "draft", "d", false, "Save the release as a draft instead of publishing it")
	cmd.Flags().BoolVarP(&isPrerelease, "prerelease", "p", false, "Mark the release as a prerelease")
	cmd.Flags().StringVarP(&notes, "notes", "n", "", "Release notes")
	cmd.Flags().StringVarP(&notesFile, "notes-file", "F", "", "Read release notes from file")
	cmd.Flags().StringVarP(&target, "target", "", "", "Target branch or full commit SHA (default: main branch)")
	cmd.Flags().StringVarP(&title, "title", "t", "", "Release title")
	return cmd
}
