package bump_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/johnmanjiro13/gh-bump/bump"
	"github.com/johnmanjiro13/gh-bump/bump/mock_bump"
)

const (
	repoDocs = `name:   johnmanjiro13/gh-bump
description:    gh extension for bumping version of a repository`
	tagList     = `v0.2.1  Latest  v0.2.1  2021-12-08T04:19:16Z`
	releaseView = `title:  v0.1.0
tag:    v0.1.0`
)

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

	gh := mock_bump.NewMockGh(ctrl)
	gh.EXPECT().ViewRepository().Return(bytes.NewBufferString(repoDocs), nil, nil)

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			b := bump.New(gh)
			assert.NoError(t, b.WithRepository(tt.repository))

			assert.Equal(t, tt.wantRepository, b.Repository())
			assert.Equal(t, tt.wantIsCurrent, b.IsCurrent())
		})
	}
}

func TestBumper_WithDraft(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gh := mock_bump.NewMockGh(ctrl)
	b := bump.New(gh)
	b.WithDraft()

	assert.True(t, b.IsDraft())
}

func TestBumper_WithPrerelease(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gh := mock_bump.NewMockGh(ctrl)
	b := bump.New(gh)
	b.WithPrerelease()

	assert.True(t, b.IsPrerelease())
}

func TestBumper_WithDiscussionCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gh := mock_bump.NewMockGh(ctrl)
	b := bump.New(gh)
	b.WithDiscussionCategory("test")

	assert.Equal(t, "test", b.DiscussionCategory())
}

func TestBumper_WithGenerateNotes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gh := mock_bump.NewMockGh(ctrl)
	b := bump.New(gh)
	b.WithGenerateNotes()

	assert.True(t, b.GenerateNotes())
}

func TestBumper_WithNotes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gh := mock_bump.NewMockGh(ctrl)
	b := bump.New(gh)
	b.WithNotes("note")

	assert.Equal(t, "note", b.Notes())
}

func TestBumper_WithNotesFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gh := mock_bump.NewMockGh(ctrl)
	b := bump.New(gh)
	b.WithNotesFile("filename")

	assert.Equal(t, "filename", b.NotesFilename())
}

func TestBumper_WithTarget(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gh := mock_bump.NewMockGh(ctrl)
	b := bump.New(gh)
	b.WithTarget("target")

	assert.Equal(t, "target", b.Target())
}

func TestBumper_WithTitle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gh := mock_bump.NewMockGh(ctrl)
	b := bump.New(gh)
	b.WithTitle("title")

	assert.Equal(t, "title", b.Title())
}

func TestBumper_ResolveRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gh := mock_bump.NewMockGh(ctrl)
	gh.EXPECT().ViewRepository().Return(bytes.NewBufferString(repoDocs), nil, nil)

	b := bump.New(gh)
	got, err := bump.ResolveRepository(b)
	assert.NoError(t, err)
	assert.Equal(t, "johnmanjiro13/gh-bump", got)
}

func TestBumper_listReleases(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gh := mock_bump.NewMockGh(ctrl)
	b := bump.New(gh)
	gh.EXPECT().ListRelease(b.Repository(), b.IsCurrent()).
		Return(bytes.NewBufferString(tagList), nil, nil)

	got, err := bump.ListReleases(b)
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("Tags:\n%s", tagList), got)
}

func TestBumper_currentVersion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gh := mock_bump.NewMockGh(ctrl)
	b := bump.New(gh)

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
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := bump.IncrementVersion(current, tt.bumpType)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestBumper_createRelease(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gh := mock_bump.NewMockGh(ctrl)
	b := bump.New(gh)

	const version = "v1.0.0"
	gh.EXPECT().CreateRelease(version, b.Repository(), b.IsCurrent(), &bump.ReleaseOption{}).
		Return(bytes.NewBufferString(version), &bytes.Buffer{}, nil)

	got, err := bump.CreateRelease(b, version)
	assert.NoError(t, err)
	assert.Equal(t, "v1.0.0", got)
}
