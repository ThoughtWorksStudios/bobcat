package common

import (
	"fmt"
)

// keeps track of canonical file paths, e.g. to determine if we have visited a file before
type FileHash map[string]bool

func (im FileHash) HaveSeen(path string) (bool, error) {
	if canonical, err := canonical(path); err != nil {
		return false, err
	} else {
		seen, inHash := im[canonical]
		return (inHash && seen), nil
	}
}

func (im FileHash) MarkSeen(path string) error {
	if canonical, err := canonical(path); err != nil {
		return err
	} else {
		im[canonical] = true
		return nil
	}
}

func (im FileHash) String() string {
	result := "FileHash -> [\n"
	for f, _ := range im {
		result += fmt.Sprintf("  %q,\n", f)
	}
	return result + "]"
}
