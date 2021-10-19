package stacker

import "strings"

func (s *Stacker) ValidBranchName(name string) bool {
	if protectedBranches[name] {
		return false
	}
	if strings.ContainsRune(name, '/') {
		return false
	}
	ok, err := s.git.ValidBranchName(name)
	// TODO: handle err
	return err == nil && ok
}

func (s *Stacker) ValidBranchNames(names ...string) bool {
	for _, name := range names {
		if !s.ValidBranchName(name) {
			return false
		}
	}
	return true
}

var protectedBranches = map[string]bool{
	"main":       true,
	"master":     true,
	"production": true,
	"release":    true,
	"staging":    true,
}
