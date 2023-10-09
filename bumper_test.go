package bump_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	bump "github.com/johnmanjiro13/gh-bump"
	"github.com/johnmanjiro13/gh-bump/mock"
)

const (
	repoDocs = `name:   johnmanjiro13/gh-bump
description:    gh extension for bumping version of a repository`
	tagList     = `v0.1.0  Latest  v0.1.0  2021-12-08T04:19:16Z`
	releaseView = `title:  v0.1.0
tag:    v0.1.0`
)

func TestBumper_Bump(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gh := mock.NewMockGh(ctrl)
	prompter := mock.NewMockPrompter(ctrl)

	t.Run("bump semver", func(t *testing.T) {
		tests := map[string]struct {
			bumpType string
			next     string
		}{
			"patch": {bumpType: "patch", next: "v0.1.1"},
			"minor": {bumpType: "minor", next: "v0.2.0"},
			"major": {bumpType: "major", next: "v1.0.0"},
		}

		for name, tt := range tests {
			t.Run(name, func(t *testing.T) {
				bumper := bump.NewBumper(gh)
				bumper.SetPrompter(prompter)
				gh.EXPECT().ListRelease(bumper.Repository(), bumper.IsCurrent()).Return(bytes.NewBufferString(tagList), nil, nil)
				gh.EXPECT().ViewRelease(bumper.Repository(), bumper.IsCurrent()).Return(bytes.NewBufferString(releaseView), nil, nil)

				prompter.EXPECT().Select("Select next version. current: v0.1.0", []string{"patch", "minor", "major"}).Return(tt.bumpType, nil)
				prompter.EXPECT().Confirm(fmt.Sprintf("Create release %s ?", tt.next)).Return(true, nil)
				gh.EXPECT().CreateRelease(tt.next, bumper.Repository(), bumper.IsCurrent(), &bump.ReleaseOption{}).Return(nil, nil, nil)
				assert.NoError(t, bumper.Bump())
			})
		}
	})

	t.Run("bump semver with option", func(t *testing.T) {
		tests := map[string]struct {
			bumpType string
			next     string
			hasError bool
		}{
			"patch":   {bumpType: "patch", next: "v0.1.1", hasError: false},
			"minor":   {bumpType: "minor", next: "v0.2.0", hasError: false},
			"major":   {bumpType: "major", next: "v1.0.0", hasError: false},
			"invalid": {bumpType: "invalid", next: "", hasError: true},
		}

		for name, tt := range tests {
			t.Run(name, func(t *testing.T) {
				bumper := bump.NewBumper(gh)
				bumper.SetPrompter(prompter)
				if tt.hasError {
					assert.ErrorIsf(t, bumper.WithBumpType(tt.bumpType), bump.ErrInvalidBumpType, "got %s", tt.bumpType)
					return
				}

				assert.NoError(t, bumper.WithBumpType(tt.bumpType))
				gh.EXPECT().ListRelease(bumper.Repository(), bumper.IsCurrent()).Return(bytes.NewBufferString(tagList), nil, nil)
				gh.EXPECT().ViewRelease(bumper.Repository(), bumper.IsCurrent()).Return(bytes.NewBufferString(releaseView), nil, nil)
				prompter.EXPECT().Confirm(fmt.Sprintf("Create release %s ?", tt.next)).Return(true, nil)
				gh.EXPECT().CreateRelease(tt.next, bumper.Repository(), bumper.IsCurrent(), &bump.ReleaseOption{}).Return(nil, nil, nil)
				assert.NoError(t, bumper.Bump())
			})
		}
	})

	t.Run("cancel bump", func(t *testing.T) {
		bumper := bump.NewBumper(gh)
		bumper.SetPrompter(prompter)
		gh.EXPECT().ListRelease(bumper.Repository(), bumper.IsCurrent()).Return(bytes.NewBufferString(tagList), nil, nil)
		gh.EXPECT().ViewRelease(bumper.Repository(), bumper.IsCurrent()).Return(bytes.NewBufferString(releaseView), nil, nil)

		prompter.EXPECT().Select("Select next version. current: v0.1.0", []string{"patch", "minor", "major"}).Return("patch", nil)
		prompter.EXPECT().Confirm("Create release v0.1.1 ?").Return(false, nil)
		assert.NoError(t, bumper.Bump())
	})

	t.Run("bump another repository", func(t *testing.T) {
		bumper := bump.NewBumper(gh)
		bumper.SetPrompter(prompter)

		const repo = "johnmanjiro13/gh-bump"
		assert.NoError(t, bumper.WithRepository("johnmanjiro13/gh-bump"))

		gh.EXPECT().ListRelease(bumper.Repository(), bumper.IsCurrent()).Return(bytes.NewBufferString(tagList), nil, nil)
		gh.EXPECT().ViewRelease(bumper.Repository(), bumper.IsCurrent()).Return(bytes.NewBufferString(releaseView), nil, nil)

		prompter.EXPECT().Select("Select next version. current: v0.1.0", []string{"patch", "minor", "major"}).Return("patch", nil)
		prompter.EXPECT().Confirm("Create release v0.1.1 ?").Return(true, nil)
		gh.EXPECT().CreateRelease("v0.1.1", repo, false, &bump.ReleaseOption{}).Return(nil, nil, nil)
		assert.NoError(t, bumper.Bump())
	})

	t.Run("bump with -y option", func(t *testing.T) {
		bumper := bump.NewBumper(gh)
		bumper.SetPrompter(prompter)

		bumper.WithYes()
		gh.EXPECT().ListRelease(bumper.Repository(), bumper.IsCurrent()).Return(bytes.NewBufferString(tagList), nil, nil)
		gh.EXPECT().ViewRelease(bumper.Repository(), bumper.IsCurrent()).Return(bytes.NewBufferString(releaseView), nil, nil)

		prompter.EXPECT().Select("Select next version. current: v0.1.0", []string{"patch", "minor", "major"}).Return("patch", nil)
		gh.EXPECT().CreateRelease("v0.1.1", bumper.Repository(), bumper.IsCurrent(), &bump.ReleaseOption{}).Return(nil, nil, nil)
		assert.NoError(t, bumper.Bump())
	})
}

func TestBumper_WithRepository(t *testing.T) {
	tests := map[string]struct {
		repository     string
		wantRepository string
		wantIsCurrent  bool
	}{
		"repository was given": {
			repository:     "johnmanjiro13/gh-bump",
			wantRepository: "johnmanjiro13/gh-bump",
			wantIsCurrent:  false,
		},
		"repository was not given": {
			repository:     "",
			wantRepository: "johnmanjiro13/gh-bump",
			wantIsCurrent:  true,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gh := mock.NewMockGh(ctrl)
	gh.EXPECT().ViewRepository().Return(bytes.NewBufferString(repoDocs), nil, nil)

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			b := bump.NewBumper(gh)
			assert.NoError(t, b.WithRepository(tt.repository))

			assert.Equal(t, tt.wantRepository, b.Repository())
			assert.Equal(t, tt.wantIsCurrent, b.IsCurrent())
		})
	}
}

func TestBumper_WithDraft(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gh := mock.NewMockGh(ctrl)
	b := bump.NewBumper(gh)
	b.WithDraft()

	assert.True(t, b.IsDraft())
}

func TestBumper_WithPrerelease(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gh := mock.NewMockGh(ctrl)
	b := bump.NewBumper(gh)
	b.WithPrerelease()

	assert.True(t, b.IsPrerelease())
}

func TestBumper_WithDiscussionCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gh := mock.NewMockGh(ctrl)
	b := bump.NewBumper(gh)
	b.WithDiscussionCategory("test")

	assert.Equal(t, "test", b.DiscussionCategory())
}

func TestBumper_WithGenerateNotes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gh := mock.NewMockGh(ctrl)
	b := bump.NewBumper(gh)
	b.WithGenerateNotes()

	assert.True(t, b.GenerateNotes())
}

func TestBumper_WithNotes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gh := mock.NewMockGh(ctrl)
	b := bump.NewBumper(gh)
	b.WithNotes("note")

	assert.Equal(t, "note", b.Notes())
}

func TestBumper_WithNotesFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gh := mock.NewMockGh(ctrl)
	b := bump.NewBumper(gh)
	b.WithNotesFile("filename")

	assert.Equal(t, "filename", b.NotesFilename())
}

func TestBumper_WithTarget(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gh := mock.NewMockGh(ctrl)
	b := bump.NewBumper(gh)
	b.WithTarget("target")

	assert.Equal(t, "target", b.Target())
}

func TestBumper_WithTitle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gh := mock.NewMockGh(ctrl)
	b := bump.NewBumper(gh)
	b.WithTitle("title")

	assert.Equal(t, "title", b.Title())
}

func TestBumper_WithAssetFiles(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gh := mock.NewMockGh(ctrl)
	b := bump.NewBumper(gh)
	b.WithAssetFiles([]string{"file1", "file2"})

	assert.Equal(t, []string{"file1", "file2"}, b.AssetFiles())
}

func TestBumper_WithBumpType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gh := mock.NewMockGh(ctrl)
	tests := map[string]struct {
		s       string
		want    bump.BumpType
		wantErr error
	}{
		"major": {
			s:       "major",
			want:    bump.Major,
			wantErr: nil,
		},
		"invalid": {
			s:       "invalid",
			want:    "",
			wantErr: fmt.Errorf("%w: got invalid", bump.ErrInvalidBumpType),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			b := bump.NewBumper(gh)
			err := b.WithBumpType(tt.s)
			assert.Equal(t, tt.want, b.BumpType())
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestBumper_WithYes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gh := mock.NewMockGh(ctrl)
	b := bump.NewBumper(gh)
	b.WithYes()
	assert.True(t, b.Yes())
}

func TestBumper_ResolveRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gh := mock.NewMockGh(ctrl)
	gh.EXPECT().ViewRepository().Return(bytes.NewBufferString(repoDocs), nil, nil)

	b := bump.NewBumper(gh)
	got, err := bump.ResolveRepository(b)
	assert.NoError(t, err)
	assert.Equal(t, "johnmanjiro13/gh-bump", got)
}

func TestBumper_listReleases(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gh := mock.NewMockGh(ctrl)
	b := bump.NewBumper(gh)
	gh.EXPECT().ListRelease(b.Repository(), b.IsCurrent()).
		Return(bytes.NewBufferString(tagList), nil, nil)

	got, err := bump.ListReleases(b)
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("Tags:\n%s", tagList), got)
}

func TestBumper_currentVersion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gh := mock.NewMockGh(ctrl)
	b := bump.NewBumper(gh)

	t.Run("new version", func(t *testing.T) {
		gh.EXPECT().ViewRelease(b.Repository(), b.IsCurrent()).
			Return(bytes.NewBufferString(releaseView), nil, nil)

		got, isInitial, err := bump.CurrentVersion(b)
		assert.NoError(t, err)
		assert.False(t, isInitial)
		want := semver.MustParse("v0.1.0")
		assert.Equal(t, want, got)
	})
}

func TestIncrementVersion(t *testing.T) {
	current := semver.MustParse("v0.1.0")

	tests := map[string]struct {
		bumpType string
		want     *semver.Version
		wantErr  error
	}{
		"major": {
			bumpType: "major",
			want:     semver.MustParse("v1.0.0"),
		},
		"minor": {
			bumpType: "minor",
			want:     semver.MustParse("v0.2.0"),
		},
		"patch": {
			bumpType: "patch",
			want:     semver.MustParse("v0.1.1"),
		},
		"invalid": {
			bumpType: "invalid",
			wantErr:  fmt.Errorf("invalid type"),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := bump.IncrementVersion(current, tt.bumpType)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestBumper_createRelease(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gh := mock.NewMockGh(ctrl)
	b := bump.NewBumper(gh)

	const version = "v1.0.0"
	gh.EXPECT().CreateRelease(version, b.Repository(), b.IsCurrent(), &bump.ReleaseOption{}).
		Return(bytes.NewBufferString(version), nil, nil)

	got, err := bump.CreateRelease(b, version)
	assert.NoError(t, err)
	assert.Equal(t, "v1.0.0", got)
}
