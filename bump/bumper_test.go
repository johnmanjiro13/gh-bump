package bump_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/johnmanjiro13/gh-bump/bump"
	"github.com/johnmanjiro13/gh-bump/bump/mock_bump"
)

const (
	repoDocs = `name:   johnmanjiro13/gh-bump
description:    gh extension for bumping version of a repository`
	tagList = `v0.2.1  Latest  v0.2.1  2021-12-08T04:19:16Z`
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
	gh.EXPECT().ViewRepository().Return(bytes.NewBufferString(repoDocs), &bytes.Buffer{}, nil)

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			b := bump.New(gh)
			assert.NoError(t, b.WithRepository(tt.repository))

			assert.Equal(t, tt.wantRepository, b.Repository())
			assert.Equal(t, tt.wantIsCurrent, b.IsCurrent())
		})
	}
}

func TestBumper_ResolveRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gh := mock_bump.NewMockGh(ctrl)
	gh.EXPECT().ViewRepository().Return(bytes.NewBufferString(repoDocs), &bytes.Buffer{}, nil)

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
		Return(bytes.NewBufferString(tagList), &bytes.Buffer{}, nil)

	got, err := bump.ListReleases(b)
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("Tags:\n%s", tagList), got)
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
