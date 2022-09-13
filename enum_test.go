package bump_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	bump "github.com/johnmanjiro13/gh-bump"
)

func TestBumpType_String(t *testing.T) {
	bumpType := bump.BumpType("major")
	assert.Equal(t, "major", bumpType.String())
}

func TestBumpType_IsBlank(t *testing.T) {
	tests := map[string]struct {
		bumpType bump.BumpType
		want     bool
	}{
		"major": {
			bumpType: bump.Major,
			want:     false,
		},
		"blank": {
			bumpType: bump.Blank,
			want:     true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.bumpType.IsBlank())
		})
	}
}

func TestBumpType_Valid(t *testing.T) {
	tests := map[string]struct {
		bumpType string
		want     error
	}{
		"major": {
			bumpType: "major",
			want:     nil,
		},
		"minor": {
			bumpType: "minor",
			want:     nil,
		},
		"patch": {
			bumpType: "patch",
			want:     nil,
		},
		"invalid": {
			bumpType: "invalid",
			want:     fmt.Errorf("%w: got invalid", bump.ErrInvalidBumpType),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			bumpType := bump.BumpType(tt.bumpType)
			assert.Equal(t, tt.want, bumpType.Valid())
		})
	}
}

func TestParseBumpType(t *testing.T) {
	tests := map[string]struct {
		s       string
		want    bump.BumpType
		wantErr error
	}{
		"major": {
			s:       "major",
			want:    bump.BumpType("major"),
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
			got, err := bump.ParseBumpType(tt.s)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
