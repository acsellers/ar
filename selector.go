package db

const (
	WHOLE_TABLE = iota
	SINGLE_COLUMN
	FORMULA
	ALIAS
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
	case WHOLE_TABLE:
		return s.Table + ".*"
	case SINGLE_COLUMN:
		if s.Alias == "" {
			return s.Table + "." + s.Column
		}
		return s.Table + "." + s.Column + " AS " + s.Alias
	case FORMULA:
		return s.Formula + " AS " + s.Alias
	case ALIAS:
		return printString(s.Value) + " AS " + s.Alias
	}

	return ""
}
