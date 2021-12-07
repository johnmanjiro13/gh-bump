package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/cli/safeexec"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

func newCmd() *cobra.Command {
	var repository string
	cmd := &cobra.Command{
		Use:   "bump",
		Short: "bump version of a given repository",
		RunE: func(cmd *cobra.Command, args []string) error {
			b := &bumper{
				repository: repository,
				cmd:        cmd,
			}
			if repository == "" {
				repo, err := resolveRepository()
				if err != nil {
					return err
				}
				b.repository = repo
				b.isCurrent = true
			}
			return b.bump()
		},
	}

	cmd.Flags().StringVarP(&repository, "repo", "R", "", "Select another repository using the [HOST/]OWNER/REPO format")

	return cmd
}

func resolveRepository() (string, error) {
	sout, _, err := gh("repo", "view")
	if err != nil {
		return "", err
	}
	viewOut := strings.Split(sout.String(), "\n")[0]
	repo := strings.TrimSpace(strings.Split(viewOut, ":")[1])
	return repo, nil
}

type bumper struct {
	repository     string
	cmd            *cobra.Command
	initialVersion bool
	isCurrent      bool
}

func (b *bumper) bump() error {
	b.listReleases()
	current, err := b.currentVersion()
	if err != nil {
		return err
	}
	nextVer, err := b.nextVersion(current)
	if err != nil {
		return err
	}
	ok, err := b.approveBump(nextVer)
	if err != nil {
		return err
	}
	if !ok {
		b.cmd.Println("Bump was canceled.")
		return nil
	}

	var sout bytes.Buffer
	if b.isCurrent {
		sout, _, err = gh("release", "create", nextVer.Original())
	} else {
		sout, _, err = gh("release", "create", nextVer.Original(), "-R", b.repository)
	}
	if err != nil {
		return err
	}
	b.cmd.Println("Release was created.")
	b.cmd.Println(sout.String())
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

func (b *bumper) listReleases() error {
	var sout bytes.Buffer
	var err error
	if b.isCurrent {
		sout, _, err = gh("release", "list")
	} else {
		sout, _, err = gh("release", "list", "-R", b.repository)
	}
	if err != nil {
		return err
	}
	b.cmd.Println("Tags:")
	b.cmd.Println(sout.String())
	return nil
}

func (b *bumper) currentVersion() (*semver.Version, error) {
	var sout, eout bytes.Buffer
	var err error
	if b.isCurrent {
		sout, eout, err = gh("release", "view")
	} else {
		sout, eout, err = gh("release", "view", "-R", b.repository)
	}
	if err != nil {
		if strings.Contains(eout.String(), "HTTP 404: Not Found") {
			current, err := newVersion()
			if err != nil {
				return nil, err
			}
			b.initialVersion = true
			return current, nil
		}
		return nil, err
	}
	viewOut := strings.Split(sout.String(), "\n")[1]
	tag := strings.TrimSpace(strings.Split(viewOut, ":")[1])
	current, err := semver.NewVersion(tag)
	if err != nil {
		return nil, fmt.Errorf("invalid version. err: %w", err)
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
		Label: fmt.Sprintf("Select next version. current: %s", current.Original()),
		Items: []string{"patch", "minor", "major"},
	}
	_, bumpType, err := prompr.Run()
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

func (b *bumper) approveBump(next *semver.Version) (bool, error) {
	prompt := promptui.Prompt{
		Label:    fmt.Sprintf("Create release %s ? [y/n]", next.Original()),
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

func gh(args ...string) (sout, eout bytes.Buffer, err error) {
	ghBin, err := safeexec.LookPath("gh")
	if err != nil {
		err = fmt.Errorf("could not find gh. err: %w", err)
		return
	}

	cmd := exec.Command(ghBin, args...)
	cmd.Stdout = &sout
	cmd.Stderr = &eout

	err = cmd.Run()
	if err != nil {
		err = fmt.Errorf("failed to run gh. err: %w, eout: %s", err, eout.String())
		return
	}
	return
}
