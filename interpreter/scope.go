package interpreter

import (
	"fmt"
)

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
	imports FileHash
	symbols SymbolTable
}

func (s *Scope) PredefinedDefaults(defaults SymbolTable) {
	for symbol, value := range defaults {
		if s.DefinedInScope(symbol) == nil {
			s.SetSymbol(symbol, value)
		}
	}
}

func (s *Scope) DefinedInScope(identifier string) *Scope {
	if _, ok := s.symbols[identifier]; ok {
		return s
	}

	if nil == s.parent {
		return nil
	}

	return s.parent.DefinedInScope(identifier)
}

func (s *Scope) ResolveSymbol(identifier string) interface{} {
	if entry, ok := s.symbols[identifier]; ok {
		return entry
	}

	if nil != s.parent {
		return s.parent.ResolveSymbol(identifier)
	}

	return nil
}

func (s *Scope) SetSymbol(identifier string, value interface{}) {
	s.symbols[identifier] = value
}

func (s *Scope) Extend() *Scope {
	return ExtendScope(s)
}

func NewRootScope() *Scope {
	return ExtendScope(nil)
}

func ExtendScope(parentScope *Scope) *Scope {
	return &Scope{parent: parentScope, imports: make(FileHash), symbols: make(SymbolTable)}
}

func (s Scope) String() string {
	return fmt.Sprintf(`Scope -> {
	parent: %v,
	imports: %v,
	symbols: %v
}`, s.parent, s.imports, s.symbols)
}
