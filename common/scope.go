package common

import (
	"fmt"
)

type DeferredResolver = func(scope *Scope) (interface{}, error)

type SymbolTable map[string]interface{}

func (s SymbolTable) String() string {
	result := "SymbolTable -> {\n"
	for key, val := range s {
		result += fmt.Sprintf("  %q: %v,\n", key, val)
	}
	return result + "}"
}

type Scope struct {
	parent  *Scope
	Imports FileHash
	Symbols SymbolTable
}

func (s *Scope) PredefinedDefaults(defaults SymbolTable) {
	for symbol, value := range defaults {
		if s.DefinedInScope(symbol) == nil {
			s.SetSymbol(symbol, value)
		}
	}
}

func (s *Scope) DefinedInScope(identifier string) *Scope {
	if _, ok := s.Symbols[identifier]; ok {
		return s
	}

	if nil == s.parent {
		return nil
	}

	return s.parent.DefinedInScope(identifier)
}

func (s *Scope) ResolveSymbol(identifier string) interface{} {
	if entry, ok := s.Symbols[identifier]; ok {
		return entry
	}

	if nil != s.parent {
		return s.parent.ResolveSymbol(identifier)
	}

	return nil
}

func (s *Scope) SetSymbol(identifier string, value interface{}) {
	s.Symbols[identifier] = value
}

func (s *Scope) Extend() *Scope {
	return ExtendScope(s)
}

func NewRootScope() *Scope {
	return ExtendScope(nil)
}

func ExtendScope(parentScope *Scope) *Scope {
	return &Scope{parent: parentScope, Imports: make(FileHash), Symbols: make(SymbolTable)}
}

func TransientScope(parentScope *Scope, symbols SymbolTable) *Scope {
	return &Scope{parent: parentScope, Imports: make(FileHash), Symbols: symbols}
}

func (s Scope) String() string {
	if nil == s.parent {
		return fmt.Sprintf(`Scope [ROOT] -> {
	Imports: %v,
	Symbols: %v
}`, s.Imports, s.Symbols)
	} else {
		return fmt.Sprintf(`Scope -> {
	parent: %v,
	Imports: %v,
	Symbols: %v
}`, s.parent, s.Imports, s.Symbols)
	}
}
