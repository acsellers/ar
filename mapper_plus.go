package db

type mapperPlus struct {
	source *source
	query  Scope
}

func (mp *mapperPlus) identity() *mapperPlus {
	if mp.query != nil {
		return mp
	} else {
		return &mapperPlus{
			source: mp.source,
			query:  mp.source.Identity(),
		}
	}
}
func (mp *mapperPlus) Identity() Scope {
	return mp.Dupe()
}
func (mp *mapperPlus) Dupe() MapperPlus {
	return &mapperPlus{source: mp.source, query: mp.query.Identity()}
}

func (mp *mapperPlus) Where(fragment string, args ...interface{}) Scope {
	mp = mp.identity()
	mp.query = mp.query.Where(fragment, args...)
	return mp
}

func (mp *mapperPlus) EqualTo(column string, val interface{}) Scope {
	mp = mp.identity()
	mp.query = mp.query.EqualTo(column, val)
	return mp
}

func (mp *mapperPlus) Between(column string, lower, upper interface{}) Scope {
	mp = mp.identity()
	mp.query = mp.query.Between(column, lower, upper)
	return mp
}

func (mp *mapperPlus) In(column string, items interface{}) Scope {
	mp = mp.identity()
	mp.query = mp.query.In(column, items)
	return mp
}

func (mp *mapperPlus) Cond(column string, condition int, val interface{}) Scope {
	mp = mp.identity()
	mp.query = mp.query.Cond(column, condition, val)
	return mp
}

func (mp *mapperPlus) Limit(limit int) Scope {
	mp = mp.identity()
	mp.query = mp.query.Limit(limit)
	return mp
}

func (mp *mapperPlus) Offset(offset int) Scope {
	mp = mp.identity()
	mp.query = mp.query.Offset(offset)
	return mp
}

func (mp *mapperPlus) OrderBy(column, direction string) Scope {
	mp = mp.identity()
	mp.query = mp.query.OrderBy(column, direction)
	return mp
}

func (mp *mapperPlus) Order(ordering string) Scope {
	mp = mp.identity()
	mp.query = mp.query.Order(ordering)
	return mp
}

func (mp *mapperPlus) Reorder(ordering string) Scope {
	mp = mp.identity()
	mp.query = mp.query.Reorder(ordering)
	return mp
}

func (mp *mapperPlus) Find(id, val interface{}) error {
	return mp.source.Find(id, val)
}

func (mp *mapperPlus) SelectorSql() string {
	return mp.identity().query.SelectorSql()
}
func (mp *mapperPlus) ConditionSql() (string, []interface{}) {
	return mp.identity().query.ConditionSql()
}

func (mp *mapperPlus) JoinsSql() string {
	return mp.identity().query.JoinsSql()
}
func (mp *mapperPlus) EndingSql() string {
	return mp.identity().query.EndingSql()
}

func (mp *mapperPlus) Delete() error {
	return mp.identity().Delete()
}

func (mp *mapperPlus) Retrieve(val interface{}) error {
	return mp.identity().Retrieve(val)
}

func (mp *mapperPlus) RetrieveAll(dest interface{}) error {
	return mp.identity().RetrieveAll(dest)
}

func (mp *mapperPlus) Count() (int64, error) {
	return mp.identity().Count()
}

func (mp *mapperPlus) Pluck(column, vals interface{}) error {
	return mp.identity().Pluck(column, vals)
}
func (mp *mapperPlus) TableName() string {
	return mp.source.TableName()
}

func (mp *mapperPlus) PrimaryKeyColumn() string {
	return mp.source.PrimaryKeyColumn()
}

func (mp *mapperPlus) UpdateAttribute(column string, val interface{}) error {
	return mp.identity().UpdateAttribute(column, val)
}
func (mp *mapperPlus) UpdateAttributes(values Attributes) error {
	return mp.identity().UpdateAttributes(values)
}
func (mp *mapperPlus) UpdateSql(sql string, vals ...interface{}) error {
	return mp.identity().UpdateSql(sql, vals...)
}
func (mp *mapperPlus) Initialize(val ...interface{}) error {
	return mp.source.Initialize(val...)
}
func (mp *mapperPlus) SaveAll(val interface{}) error {
	return mp.source.SaveAll(val)
}
func (mp *mapperPlus) LeftJoin(joins ...interface{}) Scope {
	mp = mp.identity()
	mp.query = mp.query.LeftJoin(joins...)
	return mp
}
func (mp *mapperPlus) InnerJoin(joins ...interface{}) Scope {
	mp = mp.identity()
	mp.query = mp.query.InnerJoin(joins...)
	return mp
}
func (mp *mapperPlus) FullJoin(joins ...interface{}) Scope {
	mp = mp.identity()
	mp.query = mp.query.FullJoin(joins...)
	return mp
}
func (mp *mapperPlus) RightJoin(joins ...interface{}) Scope {
	mp = mp.identity()
	mp.query = mp.query.RightJoin(joins...)
	return mp
}
func (mp *mapperPlus) JoinSql(sql string, args ...interface{}) Scope {
	mp = mp.identity()
	mp.query = mp.query.JoinSql(sql, args...)
	return mp
}

func (mp *mapperPlus) LeftInclude(include ...interface{}) Scope {
	mp = mp.identity()
	mp.query = mp.query.LeftInclude(include...)
	return mp
}
func (mp *mapperPlus) InnerInclude(include ...interface{}) Scope {
	mp = mp.identity()
	mp.query = mp.query.InnerInclude(include...)
	return mp
}
func (mp *mapperPlus) FullInclude(include interface{}, nullRecords interface{}) Scope {
	mp = mp.identity()
	mp.query = mp.query.FullInclude(include, nullRecords)
	return mp
}
func (mp *mapperPlus) RightInclude(include interface{}, nullRecords interface{}) Scope {
	mp = mp.identity()
	mp.query = mp.query.RightInclude(include, nullRecords)
	return mp
}
func (mp *mapperPlus) IncludeSql(il IncludeList, query string, args ...interface{}) Scope {
	mp = mp.identity()
	mp.query = mp.query.IncludeSql(il, query, args...)
	return mp
}
