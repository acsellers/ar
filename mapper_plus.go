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

func (mp *mapperPlus) JoinSql() string {
	return mp.identity().query.JoinSql()
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

func (mp *mapperPlus) TableName() string {
	return mp.source.TableName()
}

func (mp *mapperPlus) PrimaryKeyColumn() string {
	return mp.source.PrimaryKeyColumn()
}

func (mp *mapperPlus) UpdateAttribute(column string, val interface{}) error {
	return mp.identity().UpdateAttribute(column, val)
}
func (mp *mapperPlus) UpdateAttributes(values map[string]interface{}) error {
	return mp.identity().UpdateAttributes(values)
}
func (mp *mapperPlus) UpdateSql(sql string, vals ...interface{}) error {
	return mp.identity().UpdateSql(sql, vals...)
}

func (mp *mapperPlus) SaveAll(val interface{}) error {
	return mp.source.SaveAll(val)
}
