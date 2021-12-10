package bump

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockGh struct{}

func (g *mockGh) ViewRepository() (sout, eout bytes.Buffer, err error) {
	const repoDocs = `name:   johnmanjiro13/gh-bump
description:    gh extension for bumping version of a repository`
	sout.WriteString(repoDocs)
	return
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

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			b := &bumper{gh: &mockGh{}}
			assert.NoError(t, b.WithRepository(tt.repository))

			assert.Equal(t, tt.wantRepository, b.repository)
			assert.Equal(t, tt.wantIsCurrent, b.isCurrent)
		})
	}
}

func TestBumper_ResolveRepository(t *testing.T) {
	b := &bumper{gh: &mockGh{}}
	got, err := b.resolveRepository()
	assert.NoError(t, err)
	assert.Equal(t, "johnmanjiro13/gh-bump", got)
}
