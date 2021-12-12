package bump

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/manifoldco/promptui"

	"github.com/johnmanjiro13/gh-bump/cmd"
)

type Gh interface {
	ViewRepository() (sout, eout bytes.Buffer, err error)
	ListRelease(repo string, isCurrent bool) (sout, eout bytes.Buffer, err error)
	ViewRelease(repo string, isCurrent bool) (sout, eout bytes.Buffer, err error)
	CreateRelease(version string, repo string, isCurrent bool, option *ReleaseOption) (sout, eout bytes.Buffer, err error)
}

type ReleaseOption struct {
	IsPrerelease bool
	Target       string
	Title        string
}

type bumper struct {
	gh           Gh
	repository   string
	isCurrent    bool
	isPrerelease bool
	target       string
	title        string
}

func New(gh Gh) cmd.Bumper {
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

func (b *bumper) WithPrerelease() {
	b.isPrerelease = true
}

func (b *bumper) WithTarget(target string) {
	b.target = target
}

func (b *bumper) WithTitle(title string) {
	b.title = title
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
	} else {
		nextVer, err = b.nextVersion(current)
		if err != nil {
			return err
		}
	}

	ok, err := b.approve(nextVer)
	if err != nil {
		return err
	}
	if !ok {
		fmt.Println("Bump was canceled.")
		return nil
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

func (b *bumper) currentVersion() (*semver.Version, bool, error) {
	var isInitial bool
	sout, eout, err := b.gh.ViewRelease(b.repository, b.isCurrent)
	if err != nil {
		if strings.Contains(eout.String(), "HTTP 404: Not Found") {
			current, err := newVersion()
			if err != nil {
				return nil, isInitial, err
			}
			isInitial = true
			return current, isInitial, nil
		}
		return nil, isInitial, err
	}
	viewOut := strings.Split(sout.String(), "\n")[1]
	tag := strings.TrimSpace(strings.Split(viewOut, ":")[1])
	current, err := semver.NewVersion(tag)
	if err != nil {
		return nil, isInitial, fmt.Errorf("invalid version. err: %w", err)
	}
	return current, isInitial, nil
}

func newVersion() (*semver.Version, error) {
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
	}
	result, err := prompt.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to prompt. err: %w", err)
	}
	return semver.NewVersion(result)
}

func (b *bumper) nextVersion(current *semver.Version) (*semver.Version, error) {
	prompt := promptui.Select{
		Label: fmt.Sprintf("Select next version. current: %s", current.Original()),
		Items: []string{"patch", "minor", "major"},
	}
	_, bumpType, err := prompt.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to prompt. err: %w", err)
	}

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

func (b *bumper) approve(next *semver.Version) (bool, error) {
	validate := func(input string) error {
		if input != "y" && input != "yes" && input != "n" && input != "no" {
			return fmt.Errorf("invalid character. press y/n")
		}
		return nil
	}
	prompt := promptui.Prompt{
		Label:    fmt.Sprintf("Create release %s ? [y/n]", next.Original()),
		Validate: validate,
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
		IsPrerelease: b.isPrerelease,
		Target:       b.target,
		Title:        b.title,
	}
	sout, _, err := b.gh.CreateRelease(version, b.repository, b.isCurrent, option)
	if err != nil {
		return "", err
	}
	return sout.String(), nil
}
