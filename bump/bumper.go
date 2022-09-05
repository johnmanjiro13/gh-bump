package bump

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/manifoldco/promptui"
)

//go:generate mockgen -source=$GOFILE -package=mock -destination=../mock/mock_${GOPACKAGE}.go
type Gh interface {
	ViewRepository() (sout, eout bytes.Buffer, err error)
	ListRelease(repo string, isCurrent bool) (sout, eout bytes.Buffer, err error)
	ViewRelease(repo string, isCurrent bool) (sout, eout bytes.Buffer, err error)
	CreateRelease(version string, repo string, isCurrent bool, option *ReleaseOption) (sout, eout bytes.Buffer, err error)
}

type ReleaseOption struct {
	IsDraft            bool
	IsPrerelease       bool
	DiscussionCategory string
	GenerateNotes      bool
	Notes              string
	NotesFilename      string
	Target             string
	Title              string
}

type bumper struct {
	gh                 Gh
	repository         string
	isCurrent          bool
	isDraft            bool
	isPrerelease       bool
	discussionCategory string
	generateNotes      bool
	notes              string
	notesFilename      string
	target             string
	title              string
	bumpType           BumpType
	yes                bool
}

func New(gh Gh) *bumper {
	return &bumper{gh: gh}
}

func (b *bumper) WithRepository(repository string) error {
	if repository != "" {
		b.repository = repository
		return nil
	}

	repo, err := b.resolveRepository()
	if err != nil {
		return err
	}
	b.repository = repo
	b.isCurrent = true
	return nil
}

func (b *bumper) WithDraft() {
	b.isDraft = true
}

func (b *bumper) WithPrerelease() {
	b.isPrerelease = true
}

func (b *bumper) WithDiscussionCategory(category string) {
	b.discussionCategory = category
}

func (b *bumper) WithGenerateNotes() {
	b.generateNotes = true
}

func (b *bumper) WithNotes(notes string) {
	b.notes = notes
}

func (b *bumper) WithNotesFile(filename string) {
	b.notesFilename = filename
}

func (b *bumper) WithTarget(target string) {
	b.target = target
}

func (b *bumper) WithTitle(title string) {
	b.title = title
}

func (b *bumper) WithBumpType(s string) error {
	bumpType, err := ParseBumpType(s)
	if err != nil {
		return err
	}
	b.bumpType = bumpType
	return nil
}

func (b *bumper) WithYes() {
	b.yes = true
}

func (b *bumper) Bump() error {
	releases, err := b.listReleases()
	if err != nil {
		return err
	}
	fmt.Println(releases)

	current, isInitial, err := b.currentVersion()
	if err != nil {
		return err
	}
	var nextVer *semver.Version
	if isInitial {
		nextVer = current
	} else if b.bumpType.Valid() == nil && !b.bumpType.IsBlank() {
		nextVer, err = incrementVersion(current, b.bumpType.String())
		if err != nil {
			return err
		}
	} else {
		nextVer, err = nextVersion(current, os.Stdin, os.Stdout)
		if err != nil {
			return err
		}
	}

	// Skip approval if --yes is set
	if !b.yes {
		ok, err := approve(nextVer, os.Stdin, os.Stdout)
		if err != nil {
			return err
		}
		if !ok {
			fmt.Println("Bump was canceled.")
			return nil
		}
	}

	result, err := b.createRelease(nextVer.Original())
	if err != nil {
		return err
	}
	fmt.Println("Release was created.")
	fmt.Println(result)
	return nil
}

func (b *bumper) resolveRepository() (string, error) {
	sout, _, err := b.gh.ViewRepository()
	if err != nil {
		return "", err
	}
	viewOut := strings.Split(sout.String(), "\n")[0]
	repo := strings.TrimSpace(strings.Split(viewOut, ":")[1])
	return repo, nil
}

func (b *bumper) listReleases() (string, error) {
	sout, _, err := b.gh.ListRelease(b.repository, b.isCurrent)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Tags:\n%s", sout.String()), nil
}

func (b *bumper) currentVersion() (current *semver.Version, isInitial bool, err error) {
	sout, eout, err := b.gh.ViewRelease(b.repository, b.isCurrent)
	if err != nil {
		if strings.Contains(eout.String(), "HTTP 404: Not Found") {
			current, err = newVersion(os.Stdin, os.Stdout)
			if err != nil {
				return nil, false, err
			}
			return current, true, nil
		}
		return nil, false, err
	}
	viewOut := strings.Split(sout.String(), "\n")[1]
	tag := strings.TrimSpace(strings.Split(viewOut, ":")[1])
	current, err = semver.NewVersion(tag)
	if err != nil {
		return nil, false, fmt.Errorf("invalid version. err: %w", err)
	}
	return current, false, nil
}

func newVersion(sin io.ReadCloser, sout io.WriteCloser) (*semver.Version, error) {
	validate := func(input string) error {
		_, err := semver.NewVersion(input)
		if err != nil {
			return fmt.Errorf("invalid version. err: %w", err)
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "New version",
		Validate: validate,
		Stdin:    sin,
		Stdout:   sout,
	}
	result, err := prompt.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to prompt. err: %w", err)
	}
	return semver.NewVersion(result)
}

func nextVersion(current *semver.Version, sin io.ReadCloser, sout io.WriteCloser) (*semver.Version, error) {
	prompt := promptui.Select{
		Label:  fmt.Sprintf("Select next version. current: %s", current.Original()),
		Items:  []string{"patch", "minor", "major"},
		Stdin:  sin,
		Stdout: sout,
	}
	_, bumpType, err := prompt.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to prompt. err: %w", err)
	}

	return incrementVersion(current, bumpType)
}

func incrementVersion(current *semver.Version, bumpType string) (*semver.Version, error) {
	var next semver.Version
	switch bumpType {
	case "major":
		next = current.IncMajor()
	case "minor":
		next = current.IncMinor()
	case "patch":
		next = current.IncPatch()
	default:
		return nil, fmt.Errorf("invalid type")
	}
	return &next, nil
}

func approve(next *semver.Version, sin io.ReadCloser, sout io.WriteCloser) (bool, error) {
	validate := func(input string) error {
		if input != "y" && input != "yes" && input != "n" && input != "no" {
			return fmt.Errorf("invalid character. press y/n")
		}
		return nil
	}
	prompt := promptui.Prompt{
		Label:    fmt.Sprintf("Create release %s ? [y/n]", next.Original()),
		Validate: validate,
		Stdin:    sin,
		Stdout:   sout,
	}
	result, err := prompt.Run()
	if err != nil {
		return false, fmt.Errorf("failed to prompt. err: %w", err)
	}
	if result == "y" || result == "yes" {
		return true, nil
	}
	return false, nil
}

func (b *bumper) createRelease(version string) (string, error) {
	option := &ReleaseOption{
		IsDraft:            b.isDraft,
		IsPrerelease:       b.isPrerelease,
		DiscussionCategory: b.discussionCategory,
		GenerateNotes:      b.generateNotes,
		Notes:              b.notes,
		NotesFilename:      b.notesFilename,
		Target:             b.target,
		Title:              b.title,
	}
	sout, _, err := b.gh.CreateRelease(version, b.repository, b.isCurrent, option)
	if err != nil {
		return "", err
	}
	return sout.String(), nil
}
