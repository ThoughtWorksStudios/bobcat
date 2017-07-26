package interpreter

import (
	fs "path/filepath"
)

type ScopeEntry struct {
	Type  string
	Value interface{}
}

type SymbolTable map[string]*ScopeEntry

type Scope struct {
	parent  *Scope
	imports FileHash
	symbols SymbolTable
}

func (s *Scope) ResolveSymbol(identifier string) *ScopeEntry {
	if entry, ok := s.symbols[identifier]; ok {
		return entry
	}

	if nil != s.parent {
		return s.parent.ResolveSymbol(identifier)
	}

	return nil
}

func (s *Scope) Extend() *Scope {
	return ExtendScope(s)
}

func ExtendScope(parentScope *Scope) *Scope {
	return &Scope{parent: parentScope, imports: make(FileHash), symbols: make(SymbolTable)}
}

// keeps track of canonical file paths, e.g. to determine if we have visited a file before
type FileHash map[string]bool

func (im FileHash) Canonical(path string) (string, error) {
	if p1, e := fs.EvalSymlinks(path); e == nil {
		return fs.Abs(p1)
	} else {
		return "", e
	}
}

func (im FileHash) HaveSeen(path string) (bool, error) {
	if canonical, err := im.Canonical(path); err != nil {
		return false, err
	} else {
		return im[canonical], nil
	}
}

func (im FileHash) MarkSeen(path string) error {
	if canonical, err := im.Canonical(path); err != nil {
		return err
	} else {
		im[canonical] = true
		return nil
	}
}
