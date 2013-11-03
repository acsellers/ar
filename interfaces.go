package db

type Queryable interface {
	Identity() Scope

	Cond(column string, condition int, val interface{}) Scope

	EqualTo(column string, val interface{}) Scope
	Between(column string, lower, upper interface{}) Scope
	Where(fragment string, args ...interface{}) Scope

	Limit(limit int) Scope
	Offset(offset int) Scope

	Order(ordering string) Scope
	OrderBy(column, direction string) Scope
	Reorder(ordering string) Scope

	Find(id, val interface{}) error
	Retrieve(val interface{}) error
	RetrieveAll(dest interface{}) error
	Count() (int64, error)

	Delete() error
	UpdateAttribute(column string, val interface{}) error
	UpdateAttributes(values map[string]interface{}) error
	UpdateSql(sql string, vals ...interface{}) error

	Pluck(column string, values interface{}) error
}

type TableInformation interface {
	TableName() string
	PrimaryKeyColumn() string
}

type ScopeInformation interface {
	SelectorSql() string
	ConditionSql() (string, []interface{})
	JoinSql() string
	EndingSql() string
}

type Mapper interface {
	Queryable
	TableInformation
	Initialize(val ...interface{}) error
	SaveAll(val interface{}) error
}

type MapperPlus interface {
	Mapper
	ScopeInformation
	Dupe() MapperPlus
}

type Scope interface {
	Queryable
	TableInformation
	ScopeInformation
}
