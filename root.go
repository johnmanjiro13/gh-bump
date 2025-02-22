package bump

import "github.com/spf13/cobra"

type Bumper interface {
	Bump() error
	WithRepository(repository string) error
	WithDraft()
	WithPrerelease()
	WithDiscussionCategory(category string)
	WithGenerateNotes()
	WithNotes(notes string)
	WithNotesFile(filename string)
	WithTitle(title string)
	WithTarget(target string)
	WithAssetFiles(files []string)
	WithBumpType(bumpType string) error
	WithSuffix(suffix string)
	WithYes()
}

func NewRootCmd(bumper Bumper) *cobra.Command {
	var (
		repository         string
		isDraft            bool
		isPrerelease       bool
		discussionCategory string
		generateNotes      bool
		notes              string
		notesFile          string
		target             string
		title              string
		assetFiles         []string
		bumpType           string
		suffix             string
		yes                bool
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
			if discussionCategory != "" {
				bumper.WithDiscussionCategory(discussionCategory)
			}
			if generateNotes {
				bumper.WithGenerateNotes()
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
			if len(assetFiles) > 0 {
				bumper.WithAssetFiles(assetFiles)
			}
			if bumpType != "" {
				err := bumper.WithBumpType(bumpType)
				if err != nil {
					return err
				}
			}
			if suffix != "" {
				bumper.WithSuffix(suffix)
			}
			if yes {
				bumper.WithYes()
			}
			return bumper.Bump()
		},
	}

	cmd.Flags().StringVarP(&repository, "repo", "R", "", "Select another repository using the [HOST/]OWNER/REPO format")
	cmd.Flags().BoolVarP(&isDraft, "draft", "d", false, "Save the release as a draft instead of publishing it")
	cmd.Flags().BoolVarP(&isPrerelease, "prerelease", "p", false, "Mark the release as a prerelease")
	cmd.Flags().StringVar(&discussionCategory, "discussion-category", "", "Start a discussion of the specified category")
	cmd.Flags().BoolVarP(&generateNotes, "generate-notes", "g", false, "Automatically generate title and notes for the release")
	cmd.Flags().StringVarP(&notes, "notes", "n", "", "Release notes")
	cmd.Flags().StringVarP(&notesFile, "notes-file", "F", "", "Read release notes from file")
	cmd.Flags().StringVar(&target, "target", "", "Target branch or full commit SHA (default: main branch)")
	cmd.Flags().StringVarP(&title, "title", "t", "", "Release title")
	cmd.Flags().StringSliceVar(&assetFiles, "asset-files", []string{}, "Asset files to upload")
	cmd.Flags().StringVar(&bumpType, "bump-type", "", "Bump type (major, minor or patch)")
	cmd.Flags().StringVar(&suffix, "suffix", "", "Suffix for the version")
	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "Answer 'yes' to all questions")
	return cmd
}
