package db

const (
	other = iota
	wholeTable
	singleColumn
	formula
	alias
)

type selector struct {
	Type    int
	Table   string
	Column  string
	Formula string
	Alias   string
	Value   interface{}
}

func (s *selector) String() string {
	switch s.Type {
	case wholeTable:
		return s.Table + ".*"
	case singleColumn:
		if s.Alias == "" {
			return s.Table + "." + s.Column
		}
		return s.Table + "." + s.Column + " AS " + s.Alias
	case formula:
		return s.Formula + " AS " + s.Alias
	case alias:
		return printString(s.Value) + " AS " + s.Alias
	default:
		return s.Formula
	}
}
