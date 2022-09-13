package bump

import (
	"errors"
	"fmt"
)

type BumpType string

const (
	Major BumpType = "major"
	Minor BumpType = "minor"
	Patch BumpType = "patch"
	Blank BumpType = ""
)

var ErrInvalidBumpType = errors.New("invalid bump type")

func (b BumpType) String() string {
	return string(b)
}

func (b BumpType) IsBlank() bool {
	return b == Blank
}

func (b BumpType) Valid() error {
	switch b {
	case Major, Minor, Patch:
		return nil
	default:
		return fmt.Errorf("%w: got %s", ErrInvalidBumpType, b)
	}
}

func ParseBumpType(s string) (BumpType, error) {
	b := BumpType(s)
	if err := b.Valid(); err != nil {
		return "", err
	}
	return b, nil
}
