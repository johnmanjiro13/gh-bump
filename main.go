package main

import (
	"fmt"
	"os"
	"os/exec"
	"sort"

	"github.com/Masterminds/semver/v3"
	"github.com/cli/safeexec"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

func newCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bump",
		Short: "bump version of a given repository",
		RunE: func(cmd *cobra.Command, args []string) error {
			r, err := git.PlainOpen(".")
			if err != nil {
				return err
			}
			b := &bumper{
				repo: r,
				cmd:  cmd,
			}
			return b.bump()
		},
	}
	return cmd
}

type bumper struct {
	repo           *git.Repository
	cmd            *cobra.Command
	initialVersion bool
}

func (b *bumper) bump() error {
	current, err := b.currentVersion()
	if err != nil {
		return err
	}
	nextVer, err := b.nextVersion(current)
	ok, err := b.approveBump(nextVer)
	if err != nil {
		return err
	}
	if !ok {
		b.cmd.Println("Bump was canceled.")
		return nil
	}
	if err := gh("release", "create", nextVer.Original()); err != nil {
		return err
	}
	b.cmd.Println("Release created.")
	return nil
}

func main() {
	os.Exit(run())
}

func run() int {
	cmd := newCmd()
	if err := cmd.Execute(); err != nil {
		return 1
	}
	return 0
}

func (b *bumper) currentVersion() (*semver.Version, error) {
	tagrefs, err := b.repo.Tags()
	if err != nil {
		return nil, fmt.Errorf("failed to get tags. err: %w", err)
	}

	var tags []string
	err = tagrefs.ForEach(func(t *plumbing.Reference) error {
		tag := t.Name()
		if tag.IsTag() {
			tags = append(tags, tag.Short())
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	var current *semver.Version
	if len(tags) == 0 {
		current, err = newVersion()
		if err != nil {
			return nil, err
		}
		b.initialVersion = true
		return current, nil
	}

	versions := []*semver.Version{}
	for _, tag := range tags {
		v, err := semver.NewVersion(tag)
		if err != nil {
			continue
		}
		versions = append(versions, v)
	}
	sort.Sort(semver.Collection(versions))

	last := len(tags) - 1
	current = versions[last]
	b.cmd.Println("Tags:")
	for i, v := range versions {
		msg := fmt.Sprintf("- %s", v.Original())
		if i == last {
			msg = fmt.Sprintf("- %s (current)", v.Original())
		}
		b.cmd.Println(msg)
	}

	return current, nil
}

func newVersion() (*semver.Version, error) {
	validate := func(input string) error {
		_, err := semver.NewVersion(input)
		if err != nil {
			return fmt.Errorf("invalid version. err: %w", err)
		}
		return nil
	}

	propt := promptui.Prompt{
		Label:    "New version",
		Validate: validate,
	}
	result, err := propt.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to prompt. err: %w", err)
	}
	return semver.NewVersion(result)
}

func (b *bumper) nextVersion(current *semver.Version) (*semver.Version, error) {
	if b.initialVersion {
		return current, nil
	}
	prompr := promptui.Select{
		Label: "Select next version",
		Items: []string{"Major", "Minor", "Patch"},
	}
	_, bumpType, err := prompr.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to prompt. err: %w", err)
	}

	var next semver.Version
	switch bumpType {
	case "Major":
		next = current.IncMajor()
	case "Minor":
		next = current.IncMinor()
	case "Patch":
		next = current.IncPatch()
	default:
		return nil, fmt.Errorf("invalid type")
	}
	return &next, nil
}

func (b *bumper) approveBump(next *semver.Version) (bool, error) {
	prompt := promptui.Prompt{
		Label:    fmt.Sprintf("Create release %s? [y/n]", next.Original()),
		Validate: func(input string) error { return nil },
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

func gh(args ...string) error {
	ghBin, err := safeexec.LookPath("gh")
	if err != nil {
		return fmt.Errorf("could not find gh. err: %w", err)
	}

	cmd := exec.Command(ghBin, args...)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run gh. err: %w", err)
	}
	return nil
}
